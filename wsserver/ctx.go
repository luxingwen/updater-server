package wsserver

import (
	"context"
	"encoding/json"
	"updater-server/pkg/app"
	"updater-server/pkg/config"
	"updater-server/pkg/logger"
	"updater-server/pkg/redisop"
	"updater-server/service"

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

func (ctx *Context) AppContext() *app.Context {
	return &app.Context{
		DB:      ctx.DB,
		Redis:   ctx.Redis,
		Logger:  ctx.Logger,
		Config:  ctx.Config,
		TraceID: ctx.TraceID,
	}
}

func (ctx *Context) JSONSuccess(msg *Message) {
	msg.Method = METHOD_RESPONSE
	msg.Code = CODE_SUCCESS
	err := ctx.Client.SendMessage(msg)
	if err != nil {
		ctx.Logger.Error(err)
	}
}

func (ctx *Context) JSONError(code string, msg *Message) {
	msg.Method = METHOD_RESPONSE
	msg.Code = code
	err := ctx.Client.SendMessage(msg)
	if err != nil {
		ctx.Logger.Error(err)
	}
}

func (ctx *Context) SendRequest(uuid string, typ string, req interface{}) (err error) {
	msg := &Message{
		UUID:   uuid,
		Method: METHOD_REQUEST,
		Type:   typ,
		To:     uuid,
	}

	clientService := &service.ClientService{}

	client, err := clientService.GetClientByUUID(ctx.AppContext(), uuid)
	if err != nil {
		ctx.Logger.Error(err)
		return
	}

	if client == nil {
		ctx.Logger.Error("client not found")
		return
	}

	client, err := ctx.Proxy.GetProxy(client.UUID)

	if err != nil {
		ctx.Logger.Error(err)
		return
	}

	b, _ := json.Marshaler(req)

	msg.Data = json.RawMessage(b)

	client.SendMessage(msg)
	return
}
