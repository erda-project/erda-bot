package gh

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/erda-project/erda-bot/conf"
	"github.com/stretchr/testify/assert"
)

func init() {
	conf.Load()
	Init()
}

func TestGetBranchProtection(t *testing.T) {
	protection, err := GetBranchProtection("erda-project", "erda", "master")
	assert.NoError(t, err)
	assert.NotNil(t, protection)
	spew.Dump(protection)
	// assert.True(t, protection.RequiredPullRequestReviews.RequiredApprovingReviewCount > 0)
}
