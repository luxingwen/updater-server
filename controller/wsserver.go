package controller

import (
	"log"
	"updater-server/pkg/app"
	"updater-server/wsserver"

	"github.com/gorilla/websocket"
)

type WsController struct {
	MessageHandler *wsserver.MessageHandler
}

func (ws *WsController) Connect(c *app.Context) {

	uid := c.Param("uuid")
	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if err != nil {
		// Handle the error
		log.Println("upgrade error:", err)
		return
	}

	client := wsserver.NewProxyClient(conn, ws.MessageHandler, uid)

	//log.Println("new client:", client.UUID)
	client.MessageHandler.Context.Proxy.AddClient(client)
	client.Connected = true
	go client.Start()
	go ws.MessageHandler.HandleMessages(client, 100)
}

func (ws *WsController) GetServerId(c *app.Context) {
	c.JSONSuccess(c.AppId)
}

func (ws *WsController) GetAllProxy(c *app.Context) {
	r := ws.MessageHandler.Context.Proxy.GetAllProxy()
	type porxyInfo struct {
		UUID      string `json:"uuid"`
		Connected bool   `json:"connected"`
		Heartbeat int64  `json:"heartbeat"`
	}

	var proxyInfos []porxyInfo
	for _, v := range r {
		proxyInfos = append(proxyInfos, porxyInfo{
			UUID:      v.UUID,
			Connected: v.Connected,
			Heartbeat: v.Heartbeat,
		})
	}
	c.JSONSuccess(proxyInfos)
}
