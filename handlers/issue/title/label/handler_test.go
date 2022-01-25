package label

import (
	`testing`

	`github.com/stretchr/testify/assert`
)

func TestGetLabelFromTitle(t *testing.T) {
	tt := []struct {
		Title  string
		labels []string
	}{
		{Title: "fix(dop,pipeline): fix some bugs", labels: []string{"bugfix", "dop", "pipeline"}},
		{Title: "feat(dop,pipeline): fix some bugs", labels: []string{"feature", "dop", "pipeline"}},
		{Title: "refactor(dop,pipeline): fix some bugs", labels: []string{"refactor", "dop", "pipeline"}},
		{Title: "fix(dop): fix some bugs", labels: []string{"bugfix", "dop"}},
		{Title: "fix: fix some bugs", labels: []string{"bugfix"}},
		{Title: "fix some bugs", labels: nil},
		{Title: "fix(): fix some bugs", labels: []string{"bugfix"}},
		{Title: "fix(: fix some bugs", labels: []string{}},
		{Title: "fix): fix some bugs", labels: []string{}},
		{Title: " fix ( dop,,pipeline ): fix some bugs", labels: []string{"bugfix", "dop", "pipeline"}},
	}
	for _, v := range tt {
		assert.Equal(t, v.labels, getLabelFromTitle(v.Title))
	}
}
