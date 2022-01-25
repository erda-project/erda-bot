package label

import (
	`context`
	`strings`

	`github.com/google/go-github/v35/github`
	`github.com/sirupsen/logrus`

	`github.com/erda-project/erda-bot/gh`
	`github.com/erda-project/erda-bot/handlers`
	`github.com/erda-project/erda-bot/handlers/issue/comment`
)

var semanticLabelMap = map[string]string{
	"feat":     "feature",
	"fix":      "bugfix",
	"refactor": "refactor",
	"docs":     "",
	"style":    "",
	"perf":     "",
	"test":     "",
	"build":    "",
	"ci":       "",
	"chore":    "",
	"revert":   "",
}

type prTitleLabelHandler struct{ comment.IssueCommentHandler }

func NewPrTitleLabelHandler(nexts ...handlers.Handler) *prTitleLabelHandler {
	return &prTitleLabelHandler{*comment.NewIssueCommentHandler(nexts...)}
}

func (h *prTitleLabelHandler) Execute(ctx context.Context, req *handlers.Request) {
	e := req.Event.(github.IssueCommentEvent)
	labelsNew, err := findNewLabels(e)
	if err != nil {
		logrus.Errorf("failed to findNewLabels, err: %v", err)
		return
	}
	if len(labelsNew) == 0 {
		h.DoNexts(ctx, req)
	}
	labelsEffect, err := findEffectLabels(e, labelsNew)
	if err != nil {
		logrus.Errorf("failed to findEffectLabels, err: %v", err)
		return
	}
	if len(labelsEffect) != 0 {
		if err := gh.AddLabelsToIssue(e.Repo.Owner.GetLogin(), e.Repo.GetName(), e.Issue.GetNumber(), labelsEffect); err != nil {
			logrus.Errorf("failed to AddLabelsToIssue, err: %v", err)
			return
		}
	}

	h.DoNexts(ctx, req)
}

func findNewLabels(e github.IssueCommentEvent) ([]string, error) {
	labels := getLabelFromTitle(*e.Issue.Title)
	labelsExist, err := gh.ListLabelsByIssue(e.Repo.Owner.GetLogin(), e.Repo.GetName(), e.Issue.GetNumber(), 1, 20)
	if err != nil {
		logrus.Errorf("failed to ListLabelsByIssue, err: %v", err)
		return nil, err
	}
	labelsNew := make([]string, 0)
	labelsExistMap := make(map[string]struct{})
	for _, v := range labelsExist {
		labelsExistMap[v.GetName()] = struct{}{}
	}
	for _, v := range labels {
		if _, ok := labelsExistMap[v]; !ok {
			labelsNew = append(labelsNew, v)
		}
	}
	return labelsNew, nil
}

func findEffectLabels(e github.IssueCommentEvent, labels []string) ([]string, error) {
	labelsAll, err := gh.ListLabel(e.Repo.Owner.GetLogin(), e.Repo.GetName(), 1, 100)
	if err != nil {
		logrus.Errorf("failed to ListLabel, err: %v", err)
		return nil, err
	}
	labelsAllMap := make(map[string]struct{})
	for _, v := range labelsAll {
		labelsAllMap[v.GetName()] = struct{}{}
	}
	labelsNew := make([]string, 0)
	for _, v := range labels {
		if _, ok := labelsAllMap[v]; ok {
			labelsNew = append(labelsNew, v)
		}
	}
	return labelsNew, nil
}

// getLabelFromTitle such as fix(dop,pipeline): fix some bugs ===> [bugfix,dop,pipeline]
func getLabelFromTitle(title string) []string {
	splits := strings.Split(title, ":")
	if len(splits) < 2 {
		return nil
	}
	labels := make([]string, 0)
	prefixTitle := splits[0]
	if !strings.Contains(prefixTitle, "(") || !strings.Contains(prefixTitle, ")") {
		if v, ok := semanticLabelMap[prefixTitle]; ok && v != "" {
			labels = append(labels, v)
		}
		return labels
	}
	if v, ok := semanticLabelMap[strings.TrimSpace(prefixTitle[:strings.Index(prefixTitle, "(")])]; ok && v != "" {
		labels = append(labels, v)
	}
	if !strings.Contains(prefixTitle, ")") {
		return labels
	}
	content := prefixTitle[strings.Index(prefixTitle, "(")+1 : len(prefixTitle)-1]
	if content == "" {
		return labels
	}
	for _, v := range strings.Split(content, ",") {
		if v == "" {
			continue
		}
		labels = append(labels, strings.TrimSpace(v))
	}
	return labels
}
