package gh

import (
	"context"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"

	"github.com/erda-project/erda-bot/conf"
	"github.com/erda-project/erda/pkg/httpclient"
)

func Init() {
	// client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.Bot().GitHubToken})
	tc := oauth2.NewClient(context.Background(), ts)
	client = github.NewClient(tc)

	// httpclient
	hc = httpclient.New(httpclient.WithCompleteRedirect())
}
