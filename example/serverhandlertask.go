package example

import (
	"fmt"

	"github.com/kklab-com/gone-core/channel"
	"github.com/kklab-com/gone-http/http"
	websocket "github.com/kklab-com/gone-websocket"
	"github.com/kklab-com/goth-kkutil/value"
)

type ServerHandlerTask struct {
	websocket.DefaultServerHandlerTask
}

func (h *ServerHandlerTask) WSPing(ctx channel.HandlerContext, message *websocket.PingMessage, params map[string]any) {
	println("server WSPing")
	h.DefaultServerHandlerTask.WSPing(ctx, message, params)
	ctx.Channel().Write(h.Builder.Ping(nil, nil)).Sync()
}

func (h *ServerHandlerTask) WSPong(ctx channel.HandlerContext, message *websocket.PongMessage, params map[string]any) {
	println("server WSPong")
}

func (h *ServerHandlerTask) WSText(ctx channel.HandlerContext, message *websocket.DefaultMessage, params map[string]any) {
	println("server WSText")
	println(message.StringMessage())
	pams := map[string]any{}
	for s, a := range params {
		if s == "[gone-http]context_pack" {
			continue
		}

		pams[s] = a
	}

	var obj any = h.Builder.Text(value.JsonMarshal(struct {
		Params  map[string]any `json:"params"`
		Message string         `json:"message"`
	}{
		Params:  pams,
		Message: message.StringMessage(),
	}))

	ctx.Write(obj, nil).Sync()
}

func (h *ServerHandlerTask) WSBinary(ctx channel.HandlerContext, message *websocket.DefaultMessage, params map[string]any) {
	println("server WSBinary")
	println(message.StringMessage())
}

func (h *ServerHandlerTask) WSClose(ctx channel.HandlerContext, message *websocket.CloseMessage, params map[string]any) {
	println(fmt.Sprintf("%s server WSClose %s", ctx.Channel().ID(), message.StringMessage()))
}

func (h *ServerHandlerTask) WSConnected(ch channel.Channel, req *http.Request, resp *http.Response, params map[string]any) {
	println(fmt.Sprintf("%s server WSConnected", ch.ID()))
}

func (h *ServerHandlerTask) WSDisconnected(ch channel.Channel, req *http.Request, resp *http.Response, params map[string]any) {
	println(fmt.Sprintf("%s server WSDisconnected", ch.ID()))
	ch.Parent().Close()
}
