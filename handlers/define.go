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
	Precheck(ctx context.Context, req *Request) bool
	Execute(ctx context.Context, req *Request)
}

type NextHandler interface {
	SetNexts(...Handler)
}

type BaseHandler struct{ Nexts []Handler }

func (b *BaseHandler) Precheck(ctx context.Context, req *Request) bool { return true }
func (b *BaseHandler) Execute(ctx context.Context, req *Request)       { b.DoNexts(ctx, req) }
func (b *BaseHandler) SetNexts(handlers ...Handler)                    { b.Nexts = handlers }
func (b *BaseHandler) DoNexts(ctx context.Context, req *Request) {
	for _, next := range b.Nexts {
		if !next.Precheck(ctx, req) {
			continue
		}
		next.Execute(ctx, req)
	}
}

func NewRootHandler(nexts ...Handler) *BaseHandler { return &BaseHandler{Nexts: nexts} }
