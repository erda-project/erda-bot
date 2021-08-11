package pr

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/v35/github"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-bot/events"
	"github.com/erda-project/erda-bot/handlers"
	"github.com/erda-project/erda/pkg/strutil"
)

type prHandler struct{ handlers.BaseHandler }

func NewPRHandler(nexts ...handlers.Handler) *prHandler {
	return &prHandler{handlers.BaseHandler{Nexts: nexts}}
}

func (h *prHandler) Precheck(ctx context.Context, req *handlers.Request) bool {
	if req.EventType != events.TypePR {
		return false
	}
	return true
}

func (h *prHandler) Execute(ctx context.Context, req *handlers.Request) {
	var e github.PullRequestEvent
	if err := json.Unmarshal(req.EventBytes, &e); err != nil {
		logrus.Warnf("failed to parse event, type: %s, err: %v", events.TypePR, err)
		return
	}
	req.Event = e

	// check action
	if !strutil.Exist(events.SupportedPREventActions, e.GetAction()) {
		return
	}

	// set new event for next issue comment related handlers
	req.EventType = events.TypeIssueComment
	req.Event = github.IssueCommentEvent{
		Action: &[]string{"created"}[0],
		Issue: &github.Issue{
			ID:                e.PullRequest.ID,
			Number:            e.PullRequest.Number,
			State:             e.PullRequest.State,
			Locked:            e.PullRequest.Locked,
			Title:             e.PullRequest.Title,
			Body:              e.PullRequest.Body,
			AuthorAssociation: e.PullRequest.AuthorAssociation,
			User:              e.PullRequest.User,
			Labels:            e.PullRequest.Labels,
			Assignee:          e.PullRequest.Assignee,
			Comments:          e.PullRequest.Comments,
			ClosedAt:          e.PullRequest.ClosedAt,
			CreatedAt:         e.PullRequest.CreatedAt,
			UpdatedAt:         e.PullRequest.UpdatedAt,
			ClosedBy:          nil,
			URL:               e.PullRequest.URL,
			HTMLURL:           e.PullRequest.HTMLURL,
			CommentsURL:       e.PullRequest.CommentsURL,
			EventsURL:         e.Repo.EventsURL,
			LabelsURL:         e.Repo.LabelsURL,
			RepositoryURL:     e.Repo.URL,
			Milestone:         e.PullRequest.Milestone,
			PullRequestLinks: &github.PullRequestLinks{
				URL:      e.PullRequest.URL,
				HTMLURL:  e.PullRequest.HTMLURL,
				DiffURL:  e.PullRequest.DiffURL,
				PatchURL: e.PullRequest.PatchURL,
			},
			Repository:       e.Repo,
			Reactions:        nil,
			Assignees:        e.PullRequest.Assignees,
			NodeID:           e.PullRequest.NodeID,
			TextMatches:      nil,
			ActiveLockReason: e.PullRequest.ActiveLockReason,
		},
		Comment:      &github.IssueComment{
			ID:                nil,
			NodeID:            nil,
			Body:              e.PullRequest.Body,
			User:              e.PullRequest.User,
			Reactions:         nil,
			CreatedAt:         e.PullRequest.CreatedAt,
			UpdatedAt:         e.PullRequest.UpdatedAt,
			AuthorAssociation: e.PullRequest.AuthorAssociation,
			URL:               e.PullRequest.URL,
			HTMLURL:           e.PullRequest.HTMLURL,
			IssueURL:          e.PullRequest.IssueURL,
		},
		Changes:      nil,
		Repo:         e.Repo,
		Sender:       e.Sender,
		Installation: e.Installation,
	}
	eventBytes, _ := json.Marshal(req.Event)
	req.EventBytes = eventBytes

	h.DoNexts(ctx, req)
}
