package wsserver

import "encoding/json"

type HeartBeatMsg struct {
	Time int64 `json:"time"`
}

func HandlerProxyHeartBeat(ctx *Context) (err error) {
	var msg HeartBeatMsg
	err = json.Unmarshal(ctx.Message.Data, &msg)
	if err != nil {
		ctx.Logger.Error("ProxyHeartBeat: ", err)
		return
	}
	ctx.Client.Heartbeat = msg.Time
	return
}
