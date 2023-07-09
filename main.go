package main

import (
	"updater-server/pkg/app"
	"updater-server/pkg/config"
	"updater-server/routers"

	"github.com/gin-gonic/gin"

	"io/ioutil"
	"net/http"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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
		b, err := ioutil.ReadFile("./swagger/swagger.html") // Replace with your actual json file path
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", b)
	})

	serverApp.Router.Static("/api/v1/pkg", "./public")
	serverApp.Use(app.Cors())
	serverApp.Use(app.RequestLogger(), app.ResponseLogger())
	routers.InitRouter(serverApp)
	serverApp.Router.Run(serverApp.Config.ServerPort)
}
