package assign

import (
	"context"
	"fmt"

	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/gh"
	"github.com/erda-project/erda-bot/handlers"
	"github.com/erda-project/erda-bot/handlers/issue/comment/instruction"
	"github.com/erda-project/erda/pkg/strutil"
)

type prCommentInstructionAssignHandler struct{ handlers.BaseHandler }

func NewPrCommentInstructionAssignHandler(nexts ...handlers.Handler) *prCommentInstructionAssignHandler {
	return &prCommentInstructionAssignHandler{handlers.BaseHandler{Nexts: nexts}}
}

func (h *prCommentInstructionAssignHandler) Execute(ctx context.Context, req *handlers.Request) {
	ins := ctx.Value(instruction.CtxKeyIns).(string)
	if ins != "assign" {
		return
	}
	e := req.Event.(events.IssueCommentEvent)
	// check pr author
	if e.Issue.User.Login == e.Comment.User.Login {
		gh.CreateComment(e.Issue.CommentsURL, "Pull request authors can't review their own pull request.")
		return
	}
	// get reviewers
	assignees := ctx.Value(instruction.CtxKeyInsArgs).([]string)
	if len(assignees) == 0 {
		gh.CreateComment(e.Issue.CommentsURL, "No assignee specified!")
		return
	}
	// add reviewers
	assignees = strutil.TrimSlicePrefixes(assignees, "@") // @sfwn -> sfwn
	if err := gh.AddPRReviewers(e.Organization.Login, e.Repository.Name, e.Issue.Number, assignees); err != nil {
		gh.CreateComment(e.Issue.CommentsURL, fmt.Sprintf("Add assignees failed, err: %v", err))
		return
	}

	h.DoNexts(ctx, req)
}
