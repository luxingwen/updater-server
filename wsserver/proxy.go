package wsserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
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
	MsgIn          chan *Message
}

func NewProxyClient(conn *websocket.Conn, messageHandler *MessageHandler, uid string) *ProxyClient {
	return &ProxyClient{
		UUID:           uid,
		Conn:           conn,
		send:           make(chan []byte),
		MessageHandler: messageHandler,
		MsgIn:          make(chan *Message, 10),
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
			log.Println("read message error:", err, "uuid:", pc.UUID)
			break
		}

		//log.Println("read message:", pc.UUID, string(message))
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			// Handle error
			log.Println("unmarshal message error:", err, "uuid:", pc.UUID)
			continue
		}

		// Process the received message
		pc.SubmitMessage(&msg)
	}
}

func (pc *ProxyClient) SubmitMessage(msg *Message) {
	pc.MsgIn <- msg
}

func (pc *ProxyClient) write() {
	defer func() {
		pc.Conn.Close()
		pc.MessageHandler.Context.Proxy.RemoveClient(pc)
	}()

	for {
		select {
		case message, ok := <-pc.send:
			if !ok {
				// The send channel has been closed
				log.Println("send channel has been closed:", pc.UUID)
				return
			}

			log.Println("write message:", pc.UUID, string(message))
			err := pc.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				// Handle error
				log.Println("write message error:", err, "uuid:", pc.UUID)
				return
			}

		default:
			// No message received
			continue
		}
	}
}

func (pc *ProxyClient) Start() {
	go pc.read()
	go pc.write()

}

func (pc *ProxyClient) Send(msg []byte) {
	//fmt.Println("send message 000:", "uuid:", pc.UUID, "connected", pc.Connected)
	if pc.Connected {
		//log.Println("send message 1111:", string(msg))
		pc.send <- msg
	} else {
		fmt.Println("client not connected")
	}
}

func (pc *ProxyClient) SendMessage(msg *Message) (err error) {

	if msg.To == "" {
		fmt.Println("target client is empty")
		err = errors.New("target client is empty")
		return
	}

	if msg.Id == "" {
		msg.Id = uuid.New().String()
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("SendMessage json marshal error:", err)
		return
	}
	pc.Send(jsonMsg)
	return
}

type ProxyManager struct {
	ProxyClients map[string]*ProxyClient
	mu           sync.Mutex
}

func NewProxyManager() *ProxyManager {
	return &ProxyManager{
		ProxyClients: make(map[string]*ProxyClient),
		mu:           sync.Mutex{},
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
	close(pc.MsgIn)
	pc = nil
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
