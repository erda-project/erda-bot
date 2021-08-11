package instruction

import (
	"context"
	"strings"

	"github.com/google/go-github/v35/github"

	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/handlers"
	"github.com/erda-project/erda-bot/handlers/issue/comment"
	"github.com/erda-project/erda/pkg/strutil"
)

const (
	CtxKeyMultiIns = "multi_ins" // []events.InstructionWithArgs
)

type prCommentInstructionHandler struct{ comment.IssueCommentHandler }

func NewPrCommentInstructionHandler(nexts ...handlers.Handler) *prCommentInstructionHandler {
	return &prCommentInstructionHandler{*comment.NewIssueCommentHandler(nexts...)}
}

func (h *prCommentInstructionHandler) Execute(ctx context.Context, req *handlers.Request) {
	e, ok := req.Event.(github.IssueCommentEvent)
	if !ok {
		return
	}
	// filter pr issue
	if !e.Issue.IsPullRequest() {
		return
	}
	// instructions
	multiIns := parseMultiInstructions(e.GetComment().GetBody())
	if len(multiIns) == 0 {
		return
	}
	ctx = context.WithValue(ctx, CtxKeyMultiIns, multiIns)

	h.DoNexts(ctx, req)
}

// parse instruction from comment
func parseInstructionFromComment(comment string) (string, []string) {
	comment = strings.TrimSpace(comment)
	// comment has leading prefix "/"
	if !strings.HasPrefix(comment, "/") {
		return "", nil
	}
	comment = comment[1:]
	ss := strings.SplitN(comment, " ", -1)
	if len(ss) == 1 {
		return ss[0], nil
	}
	return ss[0], ss[1:]
}

func parseMultiInstructions(text string) []events.InstructionWithArgs {
	var results []events.InstructionWithArgs
	lines := strutil.Split(text, "\n", true)
	for _, line := range lines {
		ins, args := parseInstructionFromComment(line)
		if ins == "" {
			continue
		}
		results = append(results, events.InstructionWithArgs{Instruction: ins, Args: args})
	}
	return results
}
