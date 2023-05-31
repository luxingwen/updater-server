package wsserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"
)

const (
	METHOD_REQUEST  = "request"
	METHOD_RESPONSE = "response"

	CODE_SUCCESS = "success"
	CODE_ERROR   = "error"
	CODE_TIMEOUT = "timeout"
)

type Message struct {
	From    string          `json:"from"`
	To      string          `json:"to"`
	Id      string          `json:"id"`
	Type    string          `json:"type"`
	Method  string          `json:"method"`
	Data    json.RawMessage `json:"data"`
	Code    string          `json:"code"`
	Msg     string          `json:"msg"` // 新增 Msg 字段
	TraceId string          `json:"traceId"`
	Timeout time.Duration   // 添加 Timeout 字段
}

type HandlerFunc func(ctx *Context) error

type MessageHandler struct {
	Context  *Context
	handlers map[string]HandlerFunc
	in       chan *Message
}

func NewMessageHandler(ctx *Context, bufferSize int) *MessageHandler {
	return &MessageHandler{
		Context:  ctx,
		handlers: make(map[string]HandlerFunc),
		in:       make(chan *Message, bufferSize),
	}
}

func (h *MessageHandler) RegisterHandler(messageType string, handler HandlerFunc) {
	if _, exists := h.handlers[messageType]; exists {
		log.Fatalf("Handler already registered for message type: %s", messageType)
	}

	h.handlers[messageType] = handler
}

func (h *MessageHandler) PrintRegisteredHandlers() {
	fmt.Println("Registered Handlers:")
	for messageType, handler := range h.handlers {
		handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		fmt.Printf("Message Type: %s\nHandler: %s\n", messageType, handlerName)
	}
	fmt.Println("------------------------")
}

func (h *MessageHandler) HandleMessages(client *ProxyClient, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in HandleMessages: %v", r)
				}
			}()

			for msg := range h.in {
				ctx := context.Background()
				if msg.Timeout > 0 {
					ctx, _ = context.WithTimeout(ctx, msg.Timeout)
				}

				ctxWithCancel, cancel := context.WithCancel(ctx)

				context := &Context{
					Client:  client,
					Message: msg,
					Ctx:     ctxWithCancel,
					Cancel:  cancel,
					Extra:   make(map[string]interface{}),
					Redis:   h.Context.Redis,
					DB:      h.Context.DB,
					Logger:  h.Context.Logger,
					Config:  h.Context.Config,
				}

				if handler, ok := h.handlers[msg.Type]; ok {
					err := handler(context)
					if err != nil {
						log.Printf("Error handling message: %s", err)
					}
				} else {
					log.Printf("No handler registered for message type: %s", msg.Type)
				}
			}
		}()
	}
}

func (h *MessageHandler) SubmitMessage(msg *Message) {
	h.in <- msg
}
