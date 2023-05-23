package main

import (
	"updater-server/pkg/app"
	"updater-server/pkg/config"
	"updater-server/controller"
	"updater-server/service"

	"github.com/gin-gonic/gin"

	"net/http"
	"io/ioutil"
)

func main() {
	config.InitConfig()
	serverApp := app.NewApp()

	

	serverApp.Router.GET("/swagger/doc.json", func(c *gin.Context) {
        jsonFile, err := ioutil.ReadFile("./docs/swagger.json") // Replace with your actual json file path
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.Data(http.StatusOK, "application/json", jsonFile)
    })


	serverApp.Router.GET("/swagger/index.html", func(c *gin.Context) {
		b, err := ioutil.ReadFile("swagger.html") // Replace with your actual json file path
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
		c.Data(http.StatusOK, "text/html; charset=utf-8", b)
	})

	serverApp.Use(app.RequestLogger(), app.ResponseLogger())

	programController:=&controller.ProgramController{Service: &service.ProgramService{}}


	v1 := serverApp.Group("/api/v1")
	
	{
		v1.POST("/program/list", programController.GetAllPrograms)
		v1.POST("/program/create", programController.CreateProgram)
		v1.POST("/program/update", programController.UpdateProgram)
		v1.POST("/program/delete", programController.DeleteProgram)
	}
	


	versionController := &controller.VersionController{Service: &service.VersionService{}}
	{
		v1.POST("/program/version/list", versionController.GetAllVersions)
		v1.POST("/program/version/update", versionController.UpdateVersion)
		v1.POST("/program/version/create", versionController.CreateVersion)
		v1.POST("/program/version/delete", versionController.DeleteVersion)
	}

	packageController := &controller.PackageController{Service: &service.PackageService{}}

	{
		v1.POST("/program/package/list", packageController.GetAllPackages)
		v1.POST("/program/package/update", packageController.UpdatePackage)
		v1.POST("/program/package/create", packageController.CreatePackage)
		v1.POST("/program/package/delete", packageController.DeletePackage)
	}


	actionController:=&controller.ProgramActionController{Service:&service.ProgramActionService{}}
	{
		v1.POST("/program/action/list", actionController.GetAllProgramActions)
		v1.POST("/program/action/create", actionController.CreateProgramAction)
		v1.POST("/program/action/update", actionController.UpdateProgramAction)
		v1.POST("/program/action/delete", actionController.DeleteProgramAction)
	}

	serverApp.Router.Run(serverApp.Config.ServerPort)

}
