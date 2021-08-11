package comment

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/v35/github"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/handlers"
)

type IssueCommentHandler struct{ handlers.BaseHandler }

func NewIssueCommentHandler(nexts ...handlers.Handler) *IssueCommentHandler {
	return &IssueCommentHandler{handlers.BaseHandler{Nexts: nexts}}
}

func (h *IssueCommentHandler) Precheck(ctx context.Context, req *handlers.Request) bool {
	if req.EventType != events.TypeIssueComment {
		return false
	}
	var e github.IssueCommentEvent
	if err := json.Unmarshal(req.EventBytes, &e); err != nil {
		logrus.Errorf("failed to parse event, type: %s, err: %v", events.TypeIssueComment, err)
		return false
	}
	req.Event = e
	return true
}
