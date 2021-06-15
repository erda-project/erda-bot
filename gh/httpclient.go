package gh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/erda-project/erda-bot/conf"
	"github.com/erda-project/erda-bot/events"

	"github.com/erda-project/erda/pkg/httpclient"
)

const (
	authorizationHeader = "Authorization"
	bearer              = "Bearer "
)

var hc *httpclient.HTTPClient

func EnsureRepoForked(e events.IssueCommentEvent) (string, error) {
	forkedURL := fmt.Sprintf("https://github.com/%s/%s", conf.Bot().GitHubActor, e.Repository.Name)
	exist := GetRepo(conf.Bot().GitHubActor, e.Repository.Name)
	if exist {
		return forkedURL, nil
	}
	// not exist, do fork
	err := CreateRepoFork(e.Organization.Login, e.Repository.Name)
	if err != nil {
		return forkedURL, err
	}
	return forkedURL, nil
}

func GetRepo(owner, repo string) bool {
	hc := httpclient.New(httpclient.WithCompleteRedirect())
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	resp, err := hc.Get(url).Header(authorizationHeader, bearer+conf.Bot().GitHubToken).Do().DiscardBody()
	if err != nil {
		return false
	}
	if !resp.IsOK() {
		return false
	}
	return true
}

func CreateRepoFork(fromOrganization, repo string) error {
	hc := httpclient.New(httpclient.WithCompleteRedirect())
	targetURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/forks", fromOrganization, repo)
	var respBody bytes.Buffer
	resp, err := hc.Post(targetURL).Header(authorizationHeader, bearer+conf.Bot().GitHubToken).Do().Body(&respBody)
	if err != nil {
		return err
	}
	if !resp.IsOK() {
		return fmt.Errorf("%s", respBody.String())
	}
	// invoke success, wait success
	stopWait := false
	doneCh := make(chan bool)
	go func() {
		getted := GetRepo(conf.Bot().GitHubActor, repo)
		for !getted && !stopWait {
			time.Sleep(time.Second * 10)
			getted = GetRepo(conf.Bot().GitHubActor, repo)
		}
		if getted {
			doneCh <- true
		}
	}()
	timer := time.NewTimer(time.Minute * 5)
	select {
	case <-timer.C:
		stopWait = true
		return fmt.Errorf("wait forked repo failed")
	case <-doneCh:
		stopWait = true
		return nil
	}
}

func CreateComment(commentsURL, comment string) error {
	hc := httpclient.New(httpclient.WithCompleteRedirect())
	var respBody bytes.Buffer
	resp, err := hc.Post(commentsURL).Header(authorizationHeader, bearer+conf.Bot().GitHubToken).
		JSONBody(map[string]string{"body": comment}).Do().Body(&respBody)
	if err != nil {
		return err
	}
	if !resp.IsOK() {
		return fmt.Errorf("%s", respBody.String())
	}
	return nil
}

func AddApprovedLabel(issueURL string) error {
	hc := httpclient.New(httpclient.WithCompleteRedirect())
	var respBody bytes.Buffer
	resp, err := hc.Post(issueURL+"/labels").Header(authorizationHeader, bearer+conf.Bot().GitHubToken).
		JSONBody(map[string][]string{"labels": {"approved"}}).Do().Body(&respBody)
	if err != nil {
		return err
	}
	if !resp.IsOK() {
		return fmt.Errorf("%s", respBody.String())
	}
	return nil
}

func HaveWriteAccess(repoURL string, login string) (bool, error) {
	hc := httpclient.New(httpclient.WithCompleteRedirect())
	var respBody bytes.Buffer
	resp, err := hc.Get(repoURL+"/collaborators/"+login+"/permission").
		Header(authorizationHeader, bearer+conf.Bot().GitHubToken).
		Do().Body(&respBody)
	if err != nil {
		return false, err
	}
	if !resp.IsOK() {
		return false, fmt.Errorf("%s", respBody.String())
	}
	var permission struct {
		Permission string `json:"permission,omitempty"`
	}
	if err := json.NewDecoder(&respBody).Decode(&permission); err != nil {
		return false, err
	}
	switch permission.Permission {
	case "write", "admin":
		return true, nil
	default:
		return false, nil
	}
}
