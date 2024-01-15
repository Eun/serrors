package serrors

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Error returns the error string representation including the cause of this error.
func (e *Error) Error() string {
	var parts []string
	if e.message != "" {
		parts = append(parts, e.message)
	}

	if e.cause != nil {
		parts = append(parts, e.cause.Error())
	}

	if len(parts) == 0 {
		return "error"
	}

	return strings.Join(parts, ": ")
}

// Format formats the error according to the format specifier.
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if !s.Flag('+') {
			_, _ = io.WriteString(s, e.Error())
			_, _ = writeFields(s, GetFields(e))
			return
		}
		stack := GetStack(e)
		for i := range stack {
			_, _ = writeError(s, &stack[i])
		}
		return
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

func writeError(w io.Writer, errorStack *ErrorStack) (int, error) {
	ww := &writer{Writer: w}
	n, err := io.WriteString(ww, errorStack.ErrorMessage)
	if err != nil {
		return ww.n, err
	}
	if n > 0 {
		_, err = io.WriteString(ww, "\n")
		if err != nil {
			return ww.n, err
		}
	}
	n, err = writeFields(ww, errorStack.Fields)
	if err != nil {
		return ww.n, err
	}
	if n > 0 {
		_, err = io.WriteString(ww, "\n")
		if err != nil {
			return ww.n, err
		}
	}
	n, err = writeStackFrames(ww, errorStack.StackTrace)
	if err != nil {
		return ww.n, err
	}
	if n > 0 {
		_, err = io.WriteString(ww, "\n")
		if err != nil {
			return ww.n, err
		}
	}
	return ww.n, nil
}

func writeStackFrames(w io.Writer, frames []StackFrame) (int, error) {
	ww := &writer{Writer: w}
	size := len(frames)
	if size == 0 {
		return ww.n, nil
	}
	for i := 0; i < size-1; i++ {
		_, err := writeStackFrame(ww, &frames[i])
		if err != nil {
			return ww.n, err
		}
		_, err = io.WriteString(ww, "\n")
		if err != nil {
			return ww.n, err
		}
	}

	_, err := writeStackFrame(ww, &frames[size-1])
	if err != nil {
		return ww.n, err
	}
	return ww.n, nil
}

func writeStackFrame(w io.Writer, frame *StackFrame) (int, error) {
	return fmt.Fprintf(w, "%s\n\t%s:%d", frame.Func, frame.File, frame.Line)
}

func writeFields(w io.Writer, fields map[string]any) (int, error) {
	sizeOfFields := len(fields)
	if sizeOfFields == 0 {
		return 0, nil
	}

	s := make([]string, 0, sizeOfFields)

	for k := range fields {
		s = append(s, k)
	}
	sort.Strings(s)

	for i, k := range s {
		s[i] = fmt.Sprintf("%s=%v", k, fields[k])
	}
	return fmt.Fprint(w, "[", strings.Join(s, " "), "]")
}

type writer struct {
	io.Writer
	n int
}

func (w *writer) Write(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	if n > 0 {
		w.n += n
	}
	return n, err
}
