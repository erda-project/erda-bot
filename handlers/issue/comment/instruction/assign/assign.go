package assign

import (
	"context"
	"fmt"

	"github.com/google/go-github/v35/github"

	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/gh"
	"github.com/erda-project/erda-bot/handlers"
	"github.com/erda-project/erda-bot/handlers/issue/comment"
	"github.com/erda-project/erda-bot/handlers/issue/comment/instruction"
	"github.com/erda-project/erda/pkg/strutil"
)

type prCommentInstructionAssignHandler struct{ comment.IssueCommentHandler}

func NewPrCommentInstructionAssignHandler(nexts ...handlers.Handler) *prCommentInstructionAssignHandler {
	return &prCommentInstructionAssignHandler{*comment.NewIssueCommentHandler(nexts...)}
}

func (h *prCommentInstructionAssignHandler) Execute(ctx context.Context, req *handlers.Request) {
	multiIns := ctx.Value(instruction.CtxKeyMultiIns).([]events.InstructionWithArgs)
	var filterIns []events.InstructionWithArgs
	for _, ins := range multiIns {
		if ins.Instruction != "assign" {
			continue
		}
		filterIns = append(filterIns, ins)
	}
	e := req.Event.(github.IssueCommentEvent)

	// handle each ins
	for _, insWithArgs := range filterIns {
		// get reviewers
		if len(insWithArgs.Args) == 0 {
			gh.CreateComment(e.Issue.GetCommentsURL(), "No assignee specified!")
			return
		}
		// add reviewers
		assignees := strutil.TrimSlicePrefixes(insWithArgs.Args, "@") // @sfwn -> sfwn
		if err := gh.AddPRReviewers(e.Repo.Owner.GetLogin(), e.Repo.GetName(), e.Issue.GetNumber(), assignees); err != nil {
			gh.CreateComment(e.Issue.GetCommentsURL(), fmt.Sprintf("Add assignees failed, err: %v", err))
			return
		}
	}

	h.DoNexts(ctx, req)
}
