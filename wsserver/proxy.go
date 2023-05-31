package wsserver

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type ProxyClient struct {
	Conn           *websocket.Conn
	UUID           string
	Connected      bool
	HostIP         string
	HostName       string
	send           chan []byte
	Heartbeat      int64
	MessageHandler *MessageHandler
}

func NewProxyClient(conn *websocket.Conn, messageHandler *MessageHandler, uid string) *ProxyClient {
	return &ProxyClient{
		UUID:           uid,
		Conn:           conn,
		send:           make(chan []byte),
		MessageHandler: messageHandler,
	}
}

func (pc *ProxyClient) read() {
	defer func() {
		pc.Conn.Close()
		pc.MessageHandler.Context.Proxy.RemoveClient(pc)
	}()

	for {
		_, message, err := pc.Conn.ReadMessage()
		if err != nil {
			// Handle error
			break
		}

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			// Handle error
			continue
		}

		// Process the received message
		pc.MessageHandler.SubmitMessage(&msg)
	}
}

func (pc *ProxyClient) write() {
	defer pc.Conn.Close()

	for {
		select {
		case message, ok := <-pc.send:
			if !ok {
				// The send channel has been closed
				return
			}

			err := pc.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				// Handle error
				return
			}
		}
	}
}

func (pc *ProxyClient) Start() {
	go pc.read()
	go pc.write()

}

// 向服务器发送消息
func (pc *ProxyClient) SendMessage(msg []byte) {
	if pc.Connected {
		pc.send <- msg
	}
}

type ProxyManager struct {
	ProxyClients map[string]*ProxyClient
	mu           sync.Mutex
}

func NewProxyManager() *ProxyManager {
	return &ProxyManager{
		ProxyClients: make(map[string]*ProxyClient),
	}
}

func (pm *ProxyManager) AddClient(pc *ProxyClient) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.ProxyClients[pc.UUID] = pc
}

func (pm *ProxyManager) RemoveClient(pc *ProxyClient) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.ProxyClients, pc.UUID)
}

func (pm *ProxyManager) GetProxy(uuid string) (r *ProxyClient, err error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	r, ok := pm.ProxyClients[uuid]
	if !ok {
		return nil, errors.New("proxy not found:" + uuid)
	}
	return r, nil
}

func (pm *ProxyManager) GetAllProxy() (r []*ProxyClient) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for _, v := range pm.ProxyClients {
		r = append(r, v)
	}
	return r
}
