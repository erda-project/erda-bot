package instruction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInstructionFromComment(t *testing.T) {
	ins, args := parseInstructionFromComment("/assign sfwn @sfwn")
	assert.Equal(t, "assign", ins)
	assert.Equal(t, 2, len(args))

	ins, args = parseInstructionFromComment("/assign sfwn")
	assert.Equal(t, "assign", ins)
	assert.Equal(t, 1, len(args))
}
