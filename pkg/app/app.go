package app

import (
	"updater-server/pkg/config"
	"updater-server/pkg/db"
	"updater-server/pkg/logger"
	"updater-server/pkg/redisop"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Redis  *redisop.RedisClient
	Logger *logger.Logger
	Config *config.Config
	Router *gin.Engine
}

type AppRouterGroup struct {
	*gin.RouterGroup
	App *App
}

func NewApp() *App {
	app := &App{}
	app.Config = config.GetConfig()
	app.DB = db.GetDB(app.Config.MySQL)
	app.Logger = logger.NewLogger(app.Config.LogConfig)
	app.Router = gin.Default()

	return app
}

func (app *App) Group(relativePath string, handlers ...gin.HandlerFunc) *AppRouterGroup {
	return &AppRouterGroup{
		RouterGroup: app.Router.Group(relativePath, handlers...),
		App:         app,
	}
}

func (app *App) Use(handlers ...HandlerFunc) {
	for _, hf := range handlers {
		app.Router.Use(app.Wrap(hf))
	}
}

func (rg *AppRouterGroup) GET(relativePath string, hf HandlerFunc) {
	rg.RouterGroup.GET(relativePath, rg.App.Wrap(hf))
}

func (rg *AppRouterGroup) POST(relativePath string, hf HandlerFunc) {
	rg.RouterGroup.POST(relativePath, rg.App.Wrap(hf))
}

func (rg *AppRouterGroup) PUT(relativePath string, hf HandlerFunc) {
	rg.RouterGroup.PUT(relativePath, rg.App.Wrap(hf))
}

func (rg *AppRouterGroup) DELETE(relativePath string, hf HandlerFunc) {
	rg.RouterGroup.DELETE(relativePath, rg.App.Wrap(hf))
}

func (rg *AppRouterGroup) PATCH(relativePath string, hf HandlerFunc) {
	rg.RouterGroup.PATCH(relativePath, rg.App.Wrap(hf))
}

// Similarly, define other HTTP method handlers like POST, PUT, DELETE...
