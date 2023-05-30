package main

import (
	"updater-server/controller"
	"updater-server/pkg/app"
	"updater-server/pkg/config"
	"updater-server/service"

	"github.com/gin-gonic/gin"

	"io/ioutil"
	"net/http"
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

	v1 := serverApp.Group("/api/v1")

	authController := &controller.AuthController{AuthService: &service.AuthService{}}

	userController := &controller.UserController{UserService: &service.UserService{}}
	{
		v1.POST("/user/login", authController.Login)
		v1.POST("/user/info", userController.UserInfo)
	}

	programController := &controller.ProgramController{Service: &service.ProgramService{}}

	{
		v1.POST("/program/list", programController.GetAllPrograms)
		v1.POST("/program/create", programController.CreateProgram)
		v1.POST("/program/update", programController.UpdateProgram)
		v1.POST("/program/delete", programController.DeleteProgram)
		v1.POST("/program/detail", programController.GetProgramDetail)
	}

	versionController := &controller.VersionController{Service: &service.VersionService{}}
	{
		v1.POST("/program/version/list", versionController.GetAllVersions)
		v1.POST("/program/version/update", versionController.UpdateVersion)
		v1.POST("/program/version/create", versionController.CreateVersion)
		v1.POST("/program/version/delete", versionController.DeleteVersion)
		v1.POST("/program/version/detail", versionController.GetVersionByUUID)
	}

	packageController := &controller.PackageController{Service: &service.PackageService{}}

	{
		v1.POST("/program/package/list", packageController.GetAllPackages)
		v1.POST("/program/package/update", packageController.UpdatePackage)
		v1.POST("/program/package/create", packageController.CreatePackage)
		v1.POST("/program/package/delete", packageController.DeletePackage)
		v1.POST("/program/package/file/upload/:programUuid", packageController.UploadFile)
		v1.POST("/program/package/file/delete/:programUuid", packageController.DeleteFile)
	}

	actionController := &controller.ProgramActionController{
		Service:                    &service.ProgramActionService{},
		ClientService:              &service.ClientService{},
		TaskService:                &service.TaskService{},
		TaskExecutionRecordService: &service.TaskExecutionRecordService{},
	}
	{
		v1.POST("/program/action/list", actionController.GetAllProgramActions)
		v1.POST("/program/action/create", actionController.CreateProgramAction)
		v1.POST("/program/action/update", actionController.UpdateProgramAction)
		v1.POST("/program/action/delete", actionController.DeleteProgramAction)
		v1.POST("/program/action/execute", actionController.CreateActionTask)
		v1.POST("/program/action/detail", actionController.GetProgramActionByUUID)
	}

	clientController := &controller.ClientController{
		Service: &service.ClientService{},
	}
	{
		v1.POST("/client/list", clientController.QueryClients)
	}

	taskController := &controller.TaskController{
		Service: &service.TaskService{},
	}
	{
		v1.POST("/task/list", taskController.QueryTasks)
	}

	serverApp.Router.Run(serverApp.Config.ServerPort)
}
