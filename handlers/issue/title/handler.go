package title

import (
	`context`
	`encoding/json`

	`github.com/google/go-github/v35/github`
	`github.com/sirupsen/logrus`

	`github.com/erda-project/erda-bot/events`
	`github.com/erda-project/erda-bot/handlers`
	`github.com/erda-project/erda-bot/handlers/issue/comment`
)

type prTitleHandler struct{ comment.IssueCommentHandler }

func NewPrTitleHandler(nexts ...handlers.Handler) *prTitleHandler {
	return &prTitleHandler{*comment.NewIssueCommentHandler(nexts...)}
}

func (h *prTitleHandler) Precheck(ctx context.Context, req *handlers.Request) bool {
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
