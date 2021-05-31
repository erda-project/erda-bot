package handlers

import (
	"context"
	"net/http"
)

type Request struct {
	EventBytes  []byte        `json:"eventBytes,omitempty"`
	EventType   string        `json:"eventType,omitempty"`
	Event       interface{}   `json:"event,omitempty"`
	HTTPRequest *http.Request `json:"httpRequest,omitempty"`
}

type Handler interface {
	LogicHandler
	NextHandler
}

type LogicHandler interface {
	Execute(ctx context.Context, req *Request)
}

type NextHandler interface {
	SetNexts(...Handler)
}

type BaseHandler struct{ Nexts []Handler }

func (b *BaseHandler) SetNexts(handlers ...Handler) { b.Nexts = handlers }
func (b *BaseHandler) DoNexts(ctx context.Context, req *Request) {
	for _, next := range b.Nexts {
		next.Execute(ctx, req)
	}
}

