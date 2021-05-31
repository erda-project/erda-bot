package comment

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/handlers"
)

type issueCommentHandler struct{ handlers.BaseHandler }

func NewIssueCommentHandler(nexts ...handlers.Handler) *issueCommentHandler {
	return &issueCommentHandler{handlers.BaseHandler{Nexts: nexts}}
}

func (h *issueCommentHandler) Execute(ctx context.Context, req *handlers.Request) {
	if req.EventType != events.TypeIssueComment {
		return
	}
	var e events.IssueCommentEvent
	if err := json.Unmarshal(req.EventBytes, &e); err != nil {
		logrus.Warnf("failed to parse event, type: %s, err: %v", events.TypeIssueComment, err)
		return
	}
	req.Event = e
	h.DoNexts(ctx, req)
}
