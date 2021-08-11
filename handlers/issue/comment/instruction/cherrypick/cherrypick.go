package cherrypick

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v35/github"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-bot/conf"
	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/gh"
	"github.com/erda-project/erda-bot/handlers"
	"github.com/erda-project/erda-bot/handlers/issue/comment"
	"github.com/erda-project/erda-bot/handlers/issue/comment/instruction"
	"github.com/erda-project/erda/pkg/uuid"
)

type prCommentInstructionCherryPickHandler struct{ comment.IssueCommentHandler }

func NewPrCommentInstructionCherryPickHandler(nexts ...handlers.Handler) *prCommentInstructionCherryPickHandler {
	return &prCommentInstructionCherryPickHandler{*comment.NewIssueCommentHandler(nexts...)}
}

func (h *prCommentInstructionCherryPickHandler) Execute(ctx context.Context, req *handlers.Request) {
	multiIns := ctx.Value(instruction.CtxKeyMultiIns).([]events.InstructionWithArgs)
	var filterIns []events.InstructionWithArgs
	for _, ins := range multiIns {
		if ins.Instruction != "cherry-pick" {
			continue
		}
		filterIns = append(filterIns, ins)
	}

	// handle each ins
	for _, insWithArgs := range filterIns {
		if len(insWithArgs.Args) == 0 {
			logrus.Warnf("missing cherry-pick target branch, such as release/1.0")
			return
		}
		e := req.Event.(github.IssueCommentEvent)
		pr, err := gh.GetPR(e.Repo.Owner.GetLogin(), e.Repo.GetName(), e.Issue.GetNumber())
		if err != nil {
			logrus.Warnf("failed to get pr #%d, err: %v", e.Issue.Number, err)
			return
		}
		if !pr.GetMerged() {
			logrus.Warnf("pull request not merged, cannot cherry-pick")
			// auto add tip comment
			if err := gh.CreateComment(e.Issue.GetCommentsURL(), "Automated cherry pick can **ONLY** be triggered when this PR is **MERGED**!"); err != nil {
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

		// do cherry-picks
		for _, arg := range insWithArgs.Args {
			createPR(e, pr, forkedURL, arg)
		}
	}

	h.DoNexts(ctx, req)
}

func createPR(e github.IssueCommentEvent, pr *github.PullRequest, forkedURL string, targetBranch string) {
	// run scripts
	cmd := exec.Command("/scripts/auto_pr.sh")
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)
	cmd.Dir = tmpDir
	const cherryPickFailedDetailFile = "__cherry_pick_failed_detail.txt"
	envs := map[string]string{
		"GITHUB_ACTOR":                       conf.Bot().GitHubActor,
		"GITHUB_EMAIL":                       conf.Bot().GitHubEmail,
		"GITHUB_TOKEN":                       conf.Bot().GitHubToken,
		"FORKED_GITHUB_REPO":                 forkedURL,
		"GITHUB_REPO":                        e.Repo.GetCloneURL(),
		"CHERRY_PICK_TARGET_BRANCH":          targetBranch,
		"POLISHED_CHERRY_PICK_TARGET_BRANCH": strings.ReplaceAll(targetBranch, "/", "-"),
		"GITHUB_PR_NUM":                      fmt.Sprintf("%d", e.Issue.GetNumber()),
		"MERGE_COMMIT_SHA":                   pr.GetMergeCommitSHA(),
		"ORIGIN_ISSUE_BODY":                  e.Issue.GetBody(),
		"PR_TITLE":                           e.Issue.GetTitle(),
		"CHERRY_PICK_FAILED_DETAIL_FILE":     cherryPickFailedDetailFile,
		"UUID":                               uuid.SnowFlakeID(),
	}
	for k, v := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		logrus.Warnf("failed to exec auto_pr.sh, err: %v", err)
		// get cherry-pick failed detail
		cherryPickDetailBytes, err := os.ReadFile(filepath.Join(tmpDir, cherryPickFailedDetailFile))
		if err == nil {
			gh.CreateComment(e.Issue.GetCommentsURL(),
				fmt.Sprintf(""+
					"Automated cherry pick failed, please resolve the confilcts and create PR manually.\n"+
					"Details:\n"+
					"```\n"+
					"%s"+
					"```",
					string(cherryPickDetailBytes)))
		}
	}
}
