package approve

import (
	"context"
	"time"

	"github.com/google/go-github/v35/github"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/gh"
	"github.com/erda-project/erda-bot/handlers"
	"github.com/erda-project/erda-bot/handlers/issue/comment"
	"github.com/erda-project/erda-bot/handlers/issue/comment/instruction"
)

type prCommentInstructionApproveHandler struct{ comment.IssueCommentHandler }

func NewPrCommentInstructionApproveHandler(nexts ...handlers.Handler) *prCommentInstructionApproveHandler {
	return &prCommentInstructionApproveHandler{*comment.NewIssueCommentHandler(nexts...)}
}

func (h *prCommentInstructionApproveHandler) Execute(ctx context.Context, req *handlers.Request) {
	multiIns := ctx.Value(instruction.CtxKeyMultiIns).([]events.InstructionWithArgs)
	find := false
	for _, ins := range multiIns {
		if ins.Instruction == "approve" {
			find = true
			break
		}
	}
	if !find {
		return
	}
	e := req.Event.(github.IssueCommentEvent)
	// check pr author
	if e.Issue.User.GetLogin() == e.Comment.User.GetLogin() {
		gh.CreateComment(e.Issue.GetCommentsURL(), "Pull request authors can't approve their own pull request.")
		return
	}
	// check write access
	haveWriteAccess, err := gh.HaveWriteAccess(e.Repo.GetURL(), e.Comment.User.GetLogin())
	if err != nil {
		logrus.Warnf("failed to check write access, err: %v", err)
		return
	}
	if !haveWriteAccess {
		// send no permission comment
		gh.CreateComment(e.Issue.GetCommentsURL(), "You have no write access to use /approve instruction.")
		return
	}
	// merge
	// auto add approved label
	if err := gh.AddApprovedLabel(e.Issue.GetURL()); err != nil {
		logrus.Warnf("failed to add approved label, err: %v", err)
		return
	}
	pr, err := gh.GetPR(e.Repo.Owner.GetLogin(), e.Repo.GetName(), e.Issue.GetNumber())
	if err != nil {
		logrus.Warnf("failed to get pr, err: %v, continue", err)
		return
	}
	// master branch need at least one approval
	if *pr.Base.Ref == "master" {
		// approve by bot
		if err := gh.ApprovePR(e.Repo.Owner.GetLogin(), e.Repo.GetName(), e.Issue.GetNumber()); err != nil {
			logrus.Warnf("failed to approve pr, err: %v", err)
			return
		}
	}

	// async merge until success
	go func() {
		for {
			pr, err := gh.GetPR(e.Repo.Owner.GetLogin(), e.Repo.GetName(), e.Issue.GetNumber())
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
				if err := gh.UpdateBranch(e.Repo.Owner.GetLogin(), e.Repo.GetName(), e.Issue.GetNumber()); err != nil {
					logrus.Warnf("failed to update branch, err: %v, continue", err)
					goto sleep
				}
			case "blocked":
				goto sleep
			default:
				result, err := gh.MergePR(e.Repo.Owner.GetLogin(), e.Repo.GetName(), e.Issue.GetNumber())
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
