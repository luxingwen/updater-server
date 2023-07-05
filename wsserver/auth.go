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
		msg.To = clientInfo.UUID
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
		msg.To = clientInfo.UUID
		msg.From = "server"
		msg.Type = "v1/Register"
		ctx.JSONError(CODE_ERROR, msg)
		return
	}

	ctx.Logger.Info("HandlerRegister client: ", string(b))
	msg := ctx.Message
	msg.Type = "v1/Register"
	msg.To = clientInfo.UUID

	hearBeat := &HeartBeatMsg{
		Time: time.Now().Unix(),
	}
	b1, _ := json.Marshal(hearBeat)

	ctx.Logger.Info("HandlerRegiste msg: ", string(b1))
	msg.Data = json.RawMessage(b1)
	ctx.JSONSuccess(msg)

	// err = ctx.SendRequest(client.Uuid, "v1/Register", hearBeat)
	// if err != nil {
	// 	ctx.Logger.Error("HandlerRegister: ", err)
	// }
	return
}
