package instruction

import (
	"context"
	"strings"

	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/handlers"
)

const (
	CtxKeyIns     = "ins" // instruction
	CtxKeyInsArgs = "ins_args"
)

type prCommentInstructionHandler struct{ handlers.BaseHandler }

func NewPrCommentInstructionHandler(nexts ...handlers.Handler) *prCommentInstructionHandler {
	return &prCommentInstructionHandler{handlers.BaseHandler{Nexts: nexts}}
}

func (h *prCommentInstructionHandler) Execute(ctx context.Context, req *handlers.Request) {
	e, ok := req.Event.(events.IssueCommentEvent)
	if !ok {
		return
	}
	// filter pr issue
	if e.Issue.PullRequest == nil {
		return
	}
	// instruction
	ins, args := parseInstructionFromComment(e.Comment.Body)
	if ins == "" {
		return
	}
	ctx = context.WithValue(ctx, CtxKeyIns, ins)
	ctx = context.WithValue(ctx, CtxKeyInsArgs, args)

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
