package wsserver

import (
	"context"
	"updater-server/pkg/config"
	"updater-server/pkg/logger"
	"updater-server/pkg/redisop"

	"gorm.io/gorm"
)

type Context struct {
	DB      *gorm.DB
	Redis   *redisop.RedisClient
	Logger  *logger.Logger
	Config  *config.Config
	TraceID string
	Proxy   *ProxyManager // 新增 Proxy 字段
	Client  *ProxyClient

	Message *Message
	Extra   map[string]interface{}
	Ctx     context.Context
	Cancel  context.CancelFunc
}
