package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/erda-project/erda-bot/handlers"
	"github.com/erda-project/erda-bot/handlers/issue/comment"
	"github.com/erda-project/erda-bot/handlers/issue/comment/instruction"
	"github.com/erda-project/erda-bot/handlers/issue/comment/instruction/approve"
	"github.com/erda-project/erda-bot/handlers/issue/comment/instruction/cherrypick"
	"github.com/erda-project/erda/pkg/httpserver"
)

func Webhooks(ctx context.Context, r *http.Request, vars map[string]string) (httpserver.Responser, error) {
	// get event bytes
	eventBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return httpserver.ErrResp(http.StatusBadRequest, "", fmt.Sprintf("%v", err))
	}

	// init request
	req := &handlers.Request{EventBytes: eventBytes, HTTPRequest: r}

	// init handlers chain of responsibility
	h := handlers.NewEventTypeParseHandler(
		handlers.NewEventDispatchHandler(
			comment.NewIssueCommentHandler(
				instruction.NewPrCommentInstructionHandler(
					cherrypick.NewPrCommentInstructionCherryPickHandler(),
					approve.NewPrCommentInstructionApproveHandler(),
				),
			),
			// other event type handler here
		),
	)

	// handle request
	go h.Execute(ctx, req)

	return httpserver.OkResp(nil)
}
