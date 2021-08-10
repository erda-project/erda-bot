package star

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CatchZeng/dingtalk"
	"github.com/google/go-github/v35/github"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-bot/conf"
	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/handlers"
)

type event struct {
	github.StarEvent
	Repository github.Repository `json:"repository"`
	Sender     github.User       `json:"sender"`
}

type starHandler struct{ handlers.BaseHandler }

func NewStarHandler(nexts ...handlers.Handler) *starHandler {
	return &starHandler{handlers.BaseHandler{Nexts: nexts}}
}

func (h *starHandler) Execute(ctx context.Context, req *handlers.Request) {
	if req.EventType != events.TypeStar {
		return
	}
	var e event
	if err := json.Unmarshal(req.EventBytes, &e); err != nil {
		logrus.Warnf("failed to parse event, type: %s, err: %v", events.TypeStar, err)
		return
	}
	req.Event = e

	// print star info
	if e.GetAction() != "created" {
		return
	}
	msg := fmt.Sprintf("%s was stared by %s (star: %d)", e.Repository.GetFullName(), e.Sender.GetLogin(), e.Repository.GetStargazersCount())

	if conf.DingTalk().AccessToken == "" {
		return
	}

	client := dingtalk.NewClient(conf.DingTalk().AccessToken, conf.DingTalk().Secret)

	if _, err := client.Send(dingtalk.NewTextMessage().SetContent(msg)); err != nil {
		logrus.Errorf("failed to send dingtalk message, msg: %s, err: %v", msg, err)
	}

	h.DoNexts(ctx, req)
}
