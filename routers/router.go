package routers

import (
	"context"
	"updater-server/controller"
	"updater-server/executor"
	"updater-server/pkg/app"
	"updater-server/service"
	"updater-server/wsserver"
)

func InitRouter(ctx *app.App) {
	InitUserRouter(ctx)
	InitProgramRouter(ctx)
	InitClientRouter(ctx)
	InitTaskRouter(ctx)
	InitPreTaskRouter(ctx)
	InitScriptLibraryRouter(ctx)
	InitFileInfoRouter(ctx)
	InitDangerousCommandRouter(ctx)
	InitWsRouter(ctx)
}

func InitUserRouter(ctx *app.App) {
	v1 := ctx.Group("/api/v1")

	authController := &controller.AuthController{AuthService: &service.AuthService{}}

	userController := &controller.UserController{UserService: &service.UserService{}}
	{
		v1.POST("/user/login", authController.Login)
		v1.POST("/user/info", userController.UserInfo)
	}

}

func InitProgramRouter(ctx *app.App) {
	v1 := ctx.Group("/api/v1")

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
}

func InitClientRouter(ctx *app.App) {
	v1 := ctx.Group("/api/v1")
	clientController := &controller.ClientController{
		Service: &service.ClientService{},
	}
	{
		v1.POST("/client/list", clientController.QueryClients)
	}
}

func InitTaskRouter(ctx *app.App) {
	v1 := ctx.Group("/api/v1")

	taskController := &controller.TaskController{
		Service:                    &service.TaskService{},
		TaskExecutionRecordService: &service.TaskExecutionRecordService{},
	}
	{
		v1.POST("/task/list", taskController.QueryTasks)
		v1.POST("/task/detail", taskController.GetTaskInfo)
		v1.POST("/task/create/single", taskController.CreateSingleTask)
		v1.POST("/task/create/batch", taskController.CreateBatchTask)
		v1.POST("/task/create/multiple", taskController.CreateMultipleTask)
	}

	taskExecRecordDetail := &controller.TaskExecRecordController{
		Service: &service.TaskExecutionRecordService{},
	}

	{
		v1.POST("/task/record/list", taskExecRecordDetail.QueryTaskExecRecords)
		v1.POST("/task/record/detail", taskExecRecordDetail.GetTaskExecRecordInfo)
	}

}

// 预设任务
func InitPreTaskRouter(ctx *app.App) {
	v1 := ctx.Group("/api/v1")
	preTaskController := &controller.PreTaskController{
		PreTaskService: &service.PreTaskService{},
	}
	{
		v1.POST("/pre_task/list", preTaskController.QueryPreTasks)
		v1.POST("/pre_task/create", preTaskController.CreatePreTask)
		v1.POST("/pre_task/update", preTaskController.UpdatePreTask)
		v1.POST("/pre_task/delete", preTaskController.DeletePreTask)
		v1.POST("/pre_task/execute", preTaskController.ExecutePreTask)
		v1.POST("/pre_task/detail", preTaskController.GetPreTaskDetail)
	}
}

// 脚本库
func InitScriptLibraryRouter(ctx *app.App) {
	v1 := ctx.Group("/api/v1")
	scriptLibraryController := &controller.ScriptLibraryController{
		ScriptLibraryService: &service.ScriptLibraryService{},
	}
	{
		v1.POST("/script_library/list", scriptLibraryController.GetScriptLibraryList)
		v1.POST("/script_library/create", scriptLibraryController.CreateScriptLibrary)
		v1.POST("/script_library/update", scriptLibraryController.UpdateScriptLibrary)
		v1.POST("/script_library/delete", scriptLibraryController.DeleteScriptLibrary)
		v1.POST("/script_library/detail", scriptLibraryController.GetScriptLibrary)
	}
}

// 文件库
func InitFileInfoRouter(ctx *app.App) {
	v1 := ctx.Group("/api/v1")
	fileInfoController := &controller.FileInfoController{
		FileInfoService: &service.FileInfoService{},
	}
	{
		v1.POST("/file_info/list", fileInfoController.GetFileList)
		v1.POST("/file_info/upload", fileInfoController.UploadFile)
		v1.POST("/file_info/delete", fileInfoController.DeleteFile)
		v1.POST("/file_info/create_dir", fileInfoController.CreateDir)
	}
}

// 危险指令
func InitDangerousCommandRouter(ctx *app.App) {
	v1 := ctx.Group("/api/v1")
	dangerousCommandController := &controller.DangerousCommandController{
		DangerousCommandService: &service.DangerousCommandService{},
	}
	{
		v1.POST("/dangerous_command/create", dangerousCommandController.CreateDangerousCommand)
		v1.POST("/dangerous_command/update", dangerousCommandController.UpdateDangerousCommand)
		v1.POST("/dangerous_command/delete", dangerousCommandController.DeleteDangerousCommand)
		v1.POST("/dangerous_command/list", dangerousCommandController.GetDangerousCommandList)
		v1.POST("/dangerous_command/detail", dangerousCommandController.GetDangerousCommandInfo)
	}
}

func InitWsRouter(serverApp *app.App) {
	wsContext := &wsserver.Context{
		DB:     serverApp.DB,
		Redis:  serverApp.Redis,
		Logger: serverApp.Logger,
		Config: serverApp.Config,
		Proxy:  wsserver.NewProxyManager(),
	}

	wsController := &controller.WsController{
		MessageHandler: getMessageHandler(wsContext),
	}

	serverApp.GET("/api/v1/ws/:uuid", wsController.Connect)
	serverApp.GET("/api/v1/serverid", wsController.GetServerId)
	serverApp.POST("/api/v1/proxy/info", wsController.GetAllProxy)

	executeServer := executor.ExecutorServer{
		WsContext:                  wsContext,
		TaskService:                &service.TaskService{},
		TaskExecutionRecordService: &service.TaskExecutionRecordService{},
		ClientService:              &service.ClientService{},
	}

	wsController.MessageHandler.RegisterHandler("v1/ExecuteScript", executeServer.HandleResScript)
	wsController.MessageHandler.RegisterHandler("v1/DownloadFile", executeServer.HandleResDownloadFile)

	go executeServer.Worker(context.Background())
}

func getMessageHandler(wsContext *wsserver.Context) *wsserver.MessageHandler {
	msghanlder := wsserver.NewMessageHandler(wsContext, 10)

	wsAuthController := &wsserver.WsAuthController{
		ClientService: &service.ClientService{},
	}

	msghanlder.RegisterHandler("ProxyHeartBeat", wsserver.HandlerProxyHeartBeat)
	msghanlder.RegisterHandler("Heartbeat", wsserver.HandlerClientHeartBeat)
	msghanlder.RegisterHandler("Register", wsAuthController.HandlerRegister)
	msghanlder.RegisterHandler("v1/ClientOffline", wsAuthController.HandleClientOffline)
	msghanlder.PrintRegisteredHandlers()
	return msghanlder
}
