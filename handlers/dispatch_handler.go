package handlers

import (
	"context"

	"github.com/erda-project/erda-bot/events"
)

////////////////////////

func NewEventTypeParseHandler(nexts ...Handler) *eventTypeParseHandler {
	return &eventTypeParseHandler{BaseHandler{Nexts: nexts}}
}

type eventTypeParseHandler struct{ BaseHandler }

func (h *eventTypeParseHandler) Execute(ctx context.Context, req *Request) {
	if req.HTTPRequest == nil {
		return
	}
	eventType := req.HTTPRequest.Header.Get("X-GitHub-Event")
	if eventType == "" {
		return
	}
	req.EventType = eventType
	h.DoNexts(ctx, req)
}

////////////////////////

type eventDispatchHandler struct{ BaseHandler }

func NewEventDispatchHandler(nexts ...Handler) *eventDispatchHandler {
	return &eventDispatchHandler{BaseHandler{Nexts: nexts}}
}

func (h *eventDispatchHandler) Execute(ctx context.Context, req *Request) {
	switch req.EventType {
	case events.TypeIssueComment:
		h.DoNexts(ctx, req)
		return
	default:
		return
	}
}

