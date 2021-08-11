package events

// InstructionWithArgs
// -> cherry-pick release/1.2 release/1.1
// -> approve
type InstructionWithArgs struct {
	Instruction string
	Args        []string
}
