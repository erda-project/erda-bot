package approve

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/gh"
	"github.com/erda-project/erda-bot/handlers"
	"github.com/erda-project/erda-bot/handlers/issue/comment/instruction"
)

type prCommentInstructionApproveHandler struct{ handlers.BaseHandler }

func NewPrCommentInstructionApproveHandler(nexts ...handlers.Handler) *prCommentInstructionApproveHandler {
	return &prCommentInstructionApproveHandler{handlers.BaseHandler{Nexts: nexts}}
}

func (h *prCommentInstructionApproveHandler) Execute(ctx context.Context, req *handlers.Request) {
	ins := ctx.Value(instruction.CtxKeyIns).(string)
	if ins != "approve" {
		return
	}
	e := req.Event.(events.IssueCommentEvent)
	pr := ctx.Value(instruction.CtxKeyPR).(events.PR)
	// check write access
	haveWriteAccess, err := gh.HaveWriteAccess(e.Repository.URL, e.Comment.User.Login)
	if err != nil {
		logrus.Warnf("failed to check write access, err: %v", err)
		return
	}
	if !haveWriteAccess {
		// send no permission comment
		gh.CreateComment(e.Issue.CommentsURL, "You have no write access to use /merge instruction.")
		return
	}
	// merge
	// auto add lgtm label
	if err := gh.AddLGTMLabel(e.Issue.URL); err != nil {
		logrus.Warnf("failed to add lgtm label, err: %v", err)
		return
	}
	// async merge until success
	go func() {
		for {
			// TODO when to rebase branch
			//if pr.Rebaseable {
			//	if err := gh.UpdateBranch(e.Issue.PullRequest.URL); err != nil {
			//		logrus.Warnf("failed to update branch, err: %v, continue", err)
			//		goto sleep
			//	}
			//}
			pr, err = gh.GetPullRequest(e.Issue.PullRequest.URL)
			if err != nil {
				logrus.Warnf("failed to get issue, err: %v, continue", err)
				goto sleep
			}
			if pr.Mergeable {
				if err := gh.MergePR(e.Issue.PullRequest.URL); err != nil {
					logrus.Warnf("failed to merge pr, err: %v, continue", err)
					goto sleep
				}
			}
		sleep:
			time.Sleep(time.Minute * 1)
		}
	}()

	h.DoNexts(ctx, req)
}
