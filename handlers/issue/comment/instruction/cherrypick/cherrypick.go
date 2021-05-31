package cherrypick

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-bot/conf"
	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/gh"
	"github.com/erda-project/erda-bot/handlers"
	"github.com/erda-project/erda-bot/handlers/issue/comment/instruction"
)

type prCommentInstructionCherryPickHandler struct{ handlers.BaseHandler }

func NewPrCommentInstructionCherryPickHandler(nexts ...handlers.Handler) *prCommentInstructionCherryPickHandler {
	return &prCommentInstructionCherryPickHandler{handlers.BaseHandler{Nexts: nexts}}
}

func (h *prCommentInstructionCherryPickHandler) Execute(ctx context.Context, req *handlers.Request) {
	ins := ctx.Value(instruction.CtxKeyIns).(string)
	if ins != "cherry-pick" {
		return
	}
	args := ctx.Value(instruction.CtxKeyInsArgs).([]string)
	if len(args) == 0 {
		logrus.Warnf("missing cherry-pick target branch, such as release/1.0")
		return
	}
	e := req.Event.(events.IssueCommentEvent)
	pr := ctx.Value(instruction.CtxKeyPR).(events.PR)
	if !pr.Merged {
		logrus.Warnf("pull request not merged, cannot cherry-pick")
		// auto add tip comment
		if err := gh.CreateComment(e.Issue.CommentsURL, "Automated cherry pick can **ONLY** be triggered when this PR is **MERGED**!"); err != nil {
			logrus.Warnf("failed to create tip comment, err: %v", err)
		}
		return
	}
	// auto fork if not forked
	forkedURL, err := gh.EnsureRepoForked(e)
	if err != nil {
		logrus.Warnf("failed to ensure repo forked, err: %v", err)
		return
	}

	// run scripts
	cmd := exec.Command("/scripts/auto_pr.sh")
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)
	cmd.Dir = tmpDir
	const cherryPickFailedDetailFile = "__cherry_pick_failed_detail.txt"
	envs := map[string]string{
		"GITHUB_ACTOR":                   conf.Bot().GitHubActor,
		"GITHUB_EMAIL":                   conf.Bot().GitHubEmail,
		"GITHUB_TOKEN":                   conf.Bot().GitHubToken,
		"FORKED_GITHUB_REPO":             forkedURL,
		"GITHUB_REPO":                    e.Repository.CloneURL,
		"CHERRY_PICK_TARGET_BRANCH":      args[0],
		"GITHUB_PR_NUM":                  fmt.Sprintf("%d", e.Issue.Number),
		"MERGE_COMMIT_SHA":               pr.MergeCommitSha,
		"ORIGIN_ISSUE_BODY":              e.Issue.Body,
		"PR_TITLE":                       e.Issue.Title,
		"CHERRY_PICK_FAILED_DETAIL_FILE": cherryPickFailedDetailFile,
	}
	for k, v := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logrus.Warnf("failed to exec auto_pr.sh, err: %v", err)
		// get cherry-pick failed detail
		cherryPickDetailBytes, err := os.ReadFile(filepath.Join(tmpDir, cherryPickFailedDetailFile))
		if err == nil {
			gh.CreateComment(e.Issue.CommentsURL,
				fmt.Sprintf(""+
					"Automated cherry pick failed, please resolve the confilcts and create PR manually.\n"+
					"Details:\n"+
					"```\n"+
					"%s"+
					"```",
					string(cherryPickDetailBytes)))
		}
		return
	}

	h.DoNexts(ctx, req)
}
