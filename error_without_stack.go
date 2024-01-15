//go:build serrors_without_stack
// +build serrors_without_stack

package serrors

func (e *Error) collectStack() []uintptr {
	return nil
}

func resolveStackForStackFrames([]uintptr) []StackFrame {
	return nil
}
