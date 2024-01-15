//go:build !serrors_without_stack
// +build !serrors_without_stack

package serrors

import (
	"runtime"
	"strings"
)

func collectStack() []uintptr {
	const depth = 64
	var pcs [depth]uintptr
	n := runtime.Callers(0, pcs[:])
	return pcs[0:n]
}

func resolveStackForStackFrames(stackFrames []uintptr) []StackFrame {
	var result []StackFrame
	frames := runtime.CallersFrames(stackFrames)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		if strings.HasPrefix(frame.Function, "runtime.") ||
			strings.HasPrefix(frame.Function, "testing.") ||
			strings.HasPrefix(frame.Function, "github.com/Eun/serrors.") {
			continue
		}

		result = append(result, StackFrame{
			File: frame.File,
			Func: frame.Function,
			Line: frame.Line,
		})
	}
	return result
}
