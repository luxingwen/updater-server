package wsserver

import (
	"context"
	"encoding/json"
	"fmt"
	"updater-server/pkg/app"
	"updater-server/pkg/config"
	"updater-server/pkg/logger"
	"updater-server/pkg/redisop"
	"updater-server/service"

	"github.com/google/uuid"
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

// SendRequest 发送请求
// to: 目标客户端的 UUID
// typ: 请求类型
// traceId: 请求的 traceId
// taskId: 请求的 taskId
// req: 请求的数据
func (ctx *Context) SendRequest(to string, typ string, traceId string, taskId string, req interface{}) (err error) {
	msg := &Message{
		Id:      uuid.New().String(),
		Method:  METHOD_REQUEST,
		Type:    typ,
		To:      to,
		TraceId: traceId,
		TaskId:  taskId,
	}
	ctx.Logger.Info("SendRequest:", to, typ)

	clientService := &service.ClientService{}

	client, err := clientService.GetClientByUUID(ctx.AppContext(), to)
	if err != nil {
		ctx.Logger.Error("get client err:", err)
		return
	}

	if client == nil {
		ctx.Logger.Error("client not found")
		err = fmt.Errorf("client not found:%s", to)
		return
	}

	clientPorxy, err := ctx.Proxy.GetProxy(client.ProxyID)

	if err != nil {
		ctx.Logger.Error("get proxy err:", err)
		return
	}

	b, _ := json.Marshal(req)

	msg.Data = json.RawMessage(b)

	clientPorxy.SendMessage(msg)
	return
}
