package websocket

import (
	"time"

	"github.com/kklab-com/gone-core/channel"
	"github.com/kklab-com/gone-http/http"
	"github.com/kklab-com/goth-kklogger"
)

type HandlerTask interface {
	WSPing(ctx channel.HandlerContext, message *PingMessage, params map[string]any)
	WSPong(ctx channel.HandlerContext, message *PongMessage, params map[string]any)
	WSClose(ctx channel.HandlerContext, message *CloseMessage, params map[string]any)
	WSBinary(ctx channel.HandlerContext, message *DefaultMessage, params map[string]any)
	WSText(ctx channel.HandlerContext, message *DefaultMessage, params map[string]any)
	WSConnected(ch channel.Channel, req *http.Request, resp *http.Response, params map[string]any)
	WSDisconnected(ch channel.Channel, req *http.Request, resp *http.Response, params map[string]any)
	WSErrorCaught(ctx channel.HandlerContext, req *http.Request, resp *http.Response, msg Message, err error)
}

type ServerHandlerTask interface {
	HandlerTask
	WSUpgrade(req *http.Request, resp *http.Response, params map[string]any) bool
}

type DefaultHandlerTask struct {
	http.DefaultHandlerTask
	Builder DefaultMessageBuilder
}

func (h *DefaultHandlerTask) ErrorCaught(ctx channel.HandlerContext, err error) {
	kklogger.ErrorJ("websocket:DefaultHandlerTask", err.Error())
}

func (h *DefaultHandlerTask) WSPing(ctx channel.HandlerContext, message *PingMessage, params map[string]any) {
	dead := time.Now().Add(time.Minute)
	rtn := &PongMessage{
		DefaultMessage: DefaultMessage{
			MessageType: PongMessageType,
			Message:     message.Message,
			Dead:        &dead,
		},
	}

	ctx.Write(rtn, nil)
}

func (h *DefaultHandlerTask) WSPong(ctx channel.HandlerContext, message *PongMessage, params map[string]any) {
}

func (h *DefaultHandlerTask) WSClose(ctx channel.HandlerContext, message *CloseMessage, params map[string]any) {
}

func (h *DefaultHandlerTask) WSBinary(ctx channel.HandlerContext, message *DefaultMessage, params map[string]any) {
}

func (h *DefaultHandlerTask) WSText(ctx channel.HandlerContext, message *DefaultMessage, params map[string]any) {
}

func (h *DefaultHandlerTask) WSConnected(ch channel.Channel, req *http.Request, resp *http.Response, params map[string]any) {
}

func (h *DefaultHandlerTask) WSDisconnected(ch channel.Channel, req *http.Request, resp *http.Response, params map[string]any) {
}

func (h *DefaultHandlerTask) WSErrorCaught(ctx channel.HandlerContext, req *http.Request, resp *http.Response, msg Message, err error) {
}

type MessageBuilder interface {
	Text(msg string) *DefaultMessage
	Binary(msg []byte) *DefaultMessage
	Close(msg []byte, closeCode CloseCode) *CloseMessage
	Ping(msg []byte, deadline time.Time) *PingMessage
	Pong(msg []byte, deadline time.Time) *PongMessage
}

type DefaultMessageBuilder struct{}

func (b *DefaultMessageBuilder) Text(msg string) *DefaultMessage {
	return &DefaultMessage{
		MessageType: TextMessageType,
		Message:     []byte(msg),
	}
}

func (b *DefaultMessageBuilder) Binary(msg []byte) *DefaultMessage {
	return &DefaultMessage{
		MessageType: BinaryMessageType,
		Message:     msg,
	}
}

func (b *DefaultMessageBuilder) Close(msg []byte, closeCode CloseCode) *CloseMessage {
	return &CloseMessage{
		DefaultMessage: DefaultMessage{
			MessageType: CloseMessageType,
			Message:     msg,
		},
		CloseCode: closeCode,
	}
}

func (b *DefaultMessageBuilder) Ping(msg []byte, deadline *time.Time) *PingMessage {
	return &PingMessage{
		DefaultMessage: DefaultMessage{
			MessageType: PingMessageType,
			Message:     msg,
			Dead:        deadline,
		},
	}
}

func (b *DefaultMessageBuilder) Pong(msg []byte, deadline *time.Time) *PongMessage {
	return &PongMessage{
		DefaultMessage: DefaultMessage{
			MessageType: PongMessageType,
			Message:     msg,
			Dead:        deadline,
		},
	}
}
