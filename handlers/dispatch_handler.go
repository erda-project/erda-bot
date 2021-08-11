package handlers

import (
	"context"
)

////////////////////////

func NewEventTypeParseHandler(nexts ...Handler) *eventTypeParseHandler {
	return &eventTypeParseHandler{BaseHandler{Nexts: nexts}}
}

type eventTypeParseHandler struct{ BaseHandler }

func (h *eventTypeParseHandler) Precheck(ctx context.Context, req *Request) bool {
	if req.HTTPRequest == nil {
		return false
	}
	eventType := req.HTTPRequest.Header.Get("X-GitHub-Event")
	if eventType == "" {
		return false
	}
	req.EventType = eventType
	return true
}

////////////////////////

type eventDispatchHandler struct{ BaseHandler }

func NewEventDispatchHandler(nexts ...Handler) *eventDispatchHandler {
	return &eventDispatchHandler{BaseHandler{Nexts: nexts}}
}
