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
	// approve by bot
	if err := gh.ApprovePR(e.Organization.Login, e.Repository.Name, e.Issue.Number); err != nil {
		logrus.Warnf("failed to approve pr, err: %v", err)
		return
	}

	// async merge until success
	go func() {
		for {
			pr, err := gh.GetPR(e.Organization.Login, e.Repository.Name, e.Issue.Number)
			if err != nil {
				logrus.Warnf("failed to get pr, err: %v, continue", err)
				goto sleep
			}
			if pr.GetMerged() {
				return
			}
			if !pr.GetMergeable() {
				goto sleep
			}
			switch pr.GetMergeableState() {
			case "behind":
				if err := gh.UpdateBranch(e.Organization.Login, e.Repository.Name, e.Issue.Number); err != nil {
					logrus.Warnf("failed to update branch, err: %v, continue", err)
					goto sleep
				}
			case "blocked":
				goto sleep
			default:
				result, err := gh.MergePR(e.Organization.Login, e.Repository.Name, e.Issue.Number)
				if err != nil {
					logrus.Warnf("failed to merge pr, err: %v, continue", err)
					goto sleep
				}
				if result.GetMerged() {
					return
				}
				goto sleep
			}
		sleep:
			time.Sleep(time.Minute * 1)
		}
	}()

	h.DoNexts(ctx, req)
}
