package actions

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/v35/github"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-bot/conf"
	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/gh"
	"github.com/erda-project/erda-bot/handlers"
	"github.com/erda-project/erda-bot/handlers/issue/comment"
	"github.com/erda-project/erda-bot/handlers/issue/comment/instruction"
	"github.com/erda-project/erda/pkg/strutil"
)

// Usage:
//   /erda-actions make-image-and-auto-pr -i=java -i=echo@1.0

const commentKey = "erda-actions"

type prCommentInstructionActionHandler struct{ comment.IssueCommentHandler }

func NewPrCommentInstructionActionHandler(nexts ...handlers.Handler) *prCommentInstructionActionHandler {
	return &prCommentInstructionActionHandler{*comment.NewIssueCommentHandler(nexts...)}
}

func (h *prCommentInstructionActionHandler) Execute(ctx context.Context, req *handlers.Request) {
	multiIns := ctx.Value(instruction.CtxKeyMultiIns).([]events.InstructionWithArgs)
	var filterIns []events.InstructionWithArgs
	for _, ins := range multiIns {
		if ins.Instruction != commentKey {
			continue
		}
		filterIns = append(filterIns, ins)
	}
	e := req.Event.(github.IssueCommentEvent)

	// handle each ins
	for _, insWithArgs := range filterIns {
		cmdType, cmdOpts, err := parseCmdOptsFromArgs(insWithArgs.Args)
		if err != nil {
			logrus.Warnf("failed to parse cmd opts from args, err: %v", err)
			// auto add tip comment
			_ = gh.CreateComment(e.Issue.GetCommentsURL(), fmt.Sprintf("Failed to parse action cmd.\nError:\n%v", err))
			_ = gh.CreateComment(e.Issue.GetCommentsURL(), fmt.Sprintf(`
**Failed to parse action cmd: %s**

**Error:**
%v`, insWithArgs.Instruction, err))
			return
		}
		var cmdErr error
		switch cmdType {
		case CmdTypeOfMakeImageAndAutoPr:
			cmdErr = execCmdTypeOfMakeImageAndAutoPr(e, cmdOpts)
		default:
			// already checked when parseCmdOptsFromArgs
			return
		}
		if cmdErr != nil {
			logrus.Warnf("failed to exec cmd, err: %v", cmdErr)
			// auto add tip comment
			_ = gh.CreateComment(e.Issue.GetCommentsURL(), fmt.Sprintf(`
**Failed to execute cmd: %s**

**Error:**
%v`, cmdType, cmdErr))
			return
		}
	}

	h.DoNexts(ctx, req)
}

type CmdOption struct {
	T OptionType
	V string
}

type OptionType string

var (
	OptionTypeInputAction OptionType = "i"
)

type CmdType string

var (
	CmdTypeOfMakeImageAndAutoPr CmdType = "make-image-and-auto-pr"

	validCmdTypes = []string{string(CmdTypeOfMakeImageAndAutoPr)}
)

func (t CmdType) IsValid() bool {
	switch t {
	case CmdTypeOfMakeImageAndAutoPr:
		return true
	default:
		return false
	}
}

func parseCmdOptsFromArgs(args []string) (CmdType, []CmdOption, error) {
	// check cmd type
	cmdType := CmdType(args[0])
	if !cmdType.IsValid() {
		return "", nil, fmt.Errorf("unsupported command type: %s (Only support: %s)", cmdType, strings.Join(validCmdTypes, " "))
	}

	// check opts
	var cmdOpts []CmdOption

	for _, arg := range args[1:] {
		kv := strings.SplitN(arg, "=", 2)
		if len(kv) < 2 {
			continue
		}
		k, v := kv[0], kv[1]

		// check k
		if !strings.HasPrefix(k, "-") {
			continue
		}
		cmdOpts = append(cmdOpts, CmdOption{
			T: OptionType(k),
			V: v,
		})
	}

	return cmdType, cmdOpts, nil
}

func execCmdTypeOfMakeImageAndAutoPr(e github.IssueCommentEvent, cmdOpts []CmdOption) error {
	// run scripts
	cmd := exec.Command("/scripts/erda-actions/make-image-and-auto-pr.sh")
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)
	cmd.Dir = tmpDir
	envs := map[string]string{
		"DOCKER_REGISTRY_USERNAME": conf.ErdaActionsInfo().DockerRegistryUsername,
		"DOCKER_REGISTRY_PASSWORD": conf.ErdaActionsInfo().DockerRegistryPassword,
		"ISSUE_ID":                 strutil.String(e.Issue.GetID()),
		"TIMESTAMP":                strutil.String(time.Now().Second()),
		"ACTIONS_TO_MAKE": func() string {
			var actionsToMake []string
			for _, opt := range cmdOpts {
				if opt.T != OptionTypeInputAction {
					continue
				}
				actionsToMake = append(actionsToMake, opt.V)
			}
			return strings.Join(actionsToMake, " ")
		}(),
	}
	for k, v := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	var unionBuffer bytes.Buffer
	io.MultiWriter(os.Stdout, &unionBuffer)
	cmd.Stdout = io.MultiWriter(os.Stdout, &unionBuffer)
	cmd.Stderr = io.MultiWriter(os.Stderr, &unionBuffer)
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf(`
failed to exec erda-actions/make-image-and-auto-pr.sh
_err:_ %v
_detail:_ %s`, err, unionBuffer.String())
		logrus.Warn(err)
		return err
	}
	return nil
}
