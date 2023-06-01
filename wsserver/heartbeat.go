package wsserver

import "encoding/json"

type HeartBeatMsg struct {
	Time int64 `json:"time"`
}

type ClientInfo struct {
	UUID      string `json:"uuid"`
	HostIP    string `json:"hostIp"`
	HostName  string `json:"hostName""`
	Vmuuid    string `json:"vmuuid"`
	Sn        string `json:"sn"`       // 序列号
	OS        string `json:"os"`       //
	Arch      string `json:"arch"`     //
	Heartbeat int64  `json:"hearbeat"` // 心跳时间
	LocalIPs  string `json:"localIps"` // 本地IP地址
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

func HandlerClientHeartBeat(ctx *Context) (err error) {

	var clientInfo ClientInfo
	err = json.Unmarshal(ctx.Message.Data, &clientInfo)
	if err != nil {
		ctx.Logger.Error("ClientHeartBeat: ", err)
		return
	}

	b, _ := json.Marshal(clientInfo)

	ctx.Logger.Info("ClientHeartBeat client: ", string(b), "host ip: ", ctx.Message.ClientIP)
	return
}
