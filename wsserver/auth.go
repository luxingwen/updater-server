package wsserver

import (
	"encoding/json"
	"time"
	"updater-server/model"
	"updater-server/service"
)

type WsAuthController struct {
	ClientService *service.ClientService
}

func (wc *WsAuthController) HandlerRegister(ctx *Context) (err error) {
	var clientInfo ClientInfo
	err = json.Unmarshal(ctx.Message.Data, &clientInfo)
	if err != nil {
		ctx.Logger.Error("ClientHeartBeat: ", err)
		return
	}

	client := &model.Client{
		Uuid:     clientInfo.UUID,
		IP:       clientInfo.HostIP,
		Hostname: clientInfo.HostName,
		OS:       clientInfo.OS,
		Arch:     clientInfo.Arch,
		VMUUID:   clientInfo.Vmuuid,
		SN:       clientInfo.Sn,
		ProxyID:  ctx.Client.UUID,
	}

	err = wc.ClientService.Register(ctx.AppContext(), client)
	if err != nil {
		ctx.Logger.Error("HandlerRegister: ", err)
		msg := ctx.Message
		msg.Data = nil
		msg.To = msg.From
		msg.From = "server"
		msg.Type = "v1/Register"
		ctx.JSONError(CODE_ERROR, msg)
		return
	}

	b, err := json.Marshal(clientInfo)
	if err != nil {
		ctx.Logger.Error("HandlerRegister: ", err)
		msg := ctx.Message
		msg.Data = nil
		msg.To = msg.From
		msg.From = "server"
		msg.Type = "v1/Register"
		ctx.JSONError(CODE_ERROR, msg)
		return
	}

	ctx.Logger.Info("HandlerRegister client: ", string(b))
	msg := ctx.Message
	msg.Type = "v1/Register"
	msg.To = msg.From

	hearBeat := &HeartBeatMsg{
		Time: time.Now().Unix(),
	}
	b1, _ := json.Marshal(hearBeat)

	msg.Data = json.RawMessage(b1)
	ctx.JSONSuccess(msg)
	return
}