package events

const TypeIssueComment = "issue_comment"

type IssueCommentEvent struct {
	Action string `json:"action,omitempty"` // created, edited, deleted

	Issue struct {
		Number      int    `json:"number,omitempty"`       // 1
		URL         string `json:"url,omitempty"`          // https://api.github.com/repos/erda-project/test-cherry-pick/issues/1, use this url to get pr detail
		HtmlURL     string `json:"html_url,omitempty"`     // https://github.com/erda-project/test-cherry-pick/pull/1
		CommentsURL string `json:"comments_url,omitempty"` // https://github.com/erda-project/test-cherry-pick/pull/1/comments
		Title       string `json:"title,omitempty"`        // Update README.md
		User        struct {
			Login string `json:"login,omitempty"` // actor
		} `json:"user,omitempty"`
		PullRequest *struct {
			URL string `json:"url,omitempty"` // https://api.github.com/repos/erda-project/test-cherry-pick/pulls/1
		} `json:"pull_request,omitempty"`
		Body string `json:"body,omitempty"`
	} `json:"issue,omitempty"`

	Comment struct {
		User struct {
			Login string `json:"login,omitempty"`
		} `json:"user,omitempty"`
		Body string `json:"body,omitempty"` // /cherry-pick release/1.0
	} `json:"comment,omitempty"`

	Repository struct {
		Name          string `json:"name,omitempty"`           // test-cherry-pick
		FullName      string `json:"full_name,omitempty"`      // erda-project/test-cherry-pick
		HtmlURL       string `json:"html_url,omitempty"`       // https://github.com/erda-project/test-cherry-pick
		CloneURL      string `json:"clone_url,omitempty"`      // https://github.com/erda-project/test-cherry-pick.git
		DefaultBranch string `json:"default_branch,omitempty"` // master
		URL           string `json:"url,omitempty"`            // https://api.github.com/repos/erda-project/erda
	} `json:"repository,omitempty"`

	Organization struct {
		Login string `json:"login,omitempty"` // erda-project
	} `json:"organization,omitempty"`
}
