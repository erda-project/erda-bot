package gh

import (
	"github.com/google/go-github/v35/github"

	"github.com/erda-project/erda/pkg/httpclient"
)

func init() {
	client = github.NewClient(nil)
	hc = httpclient.New(httpclient.WithCompleteRedirect())
}
