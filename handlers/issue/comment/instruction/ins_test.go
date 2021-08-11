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

func TestParseMultiInstructions(t *testing.T) {
	results := parseMultiInstructions(`
/approve
/assign sfwn effet
/cherry-pick release/1.2 release/1.1
`)
	assert.Equal(t, 3, len(results))

	assert.Equal(t, "approve", results[0].Instruction)
	assert.Equal(t, 0, len(results[0].Args))

	assert.Equal(t, "assign", results[1].Instruction)
	assert.Equal(t, "sfwn", results[1].Args[0])
	assert.Equal(t, "effet", results[1].Args[1])

	assert.Equal(t, "cherry-pick", results[2].Instruction)
	assert.Equal(t, "release/1.2", results[2].Args[0])
	assert.Equal(t, "release/1.1", results[2].Args[1])
}
