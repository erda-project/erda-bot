package gh

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHaveWriteAccess(t *testing.T) {
	access, err := HaveWriteAccess("https://api.github.com/repos/erda-project/erda", "sfwn")
	assert.NoError(t, err)
	fmt.Println(access)
}
