package main

import (
	"updater-server/pkg/app"
	"updater-server/pkg/config"
)

func main() {
	config.InitConfig()
	serverApp := app.NewApp()

	serverApp.Use(app.RequestLogger(), app.ResponseLogger())

	v1 := serverApp.Group("/api/v1")
	v1.POST("/api/v1/client", func(ctx *app.Context) {
		ctx.JSONSuccess("ok")
	})

	serverApp.Router.Run(serverApp.Config.ServerPort)

}
