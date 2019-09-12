package compile

func makeExitIndex(idx int) CompletionStatus {
	return CompletionStatus((idx << 8) & exitIndexMask)
}

const (
	statusMask    = 15
	exitIndexMask = 0x00000000ffffff00
	unknownIndex  = 0xffffff
)

// CompletionStatus decodes and returns the completion status of the exit.
func (s JITExitSignal) CompletionStatus() CompletionStatus {
	return CompletionStatus(s & statusMask)
}

// Index returns the index to the instruction where the exit happened.
// 0xffffff is returned if the exit was due to normal completion.
func (s JITExitSignal) Index() int {
	return (int(s) & exitIndexMask) >> 8
}
