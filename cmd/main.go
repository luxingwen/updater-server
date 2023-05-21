package main

import "updater-server/pkg/app"

func main() {
	serverApp := app.NewApp()
	serverApp.Use(app.RequestLogger(), app.ResponseLogger())

	serverApp.Router.Run(app.Config.ServerPort)

}
