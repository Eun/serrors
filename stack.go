package serrors

import (
	"errors"
	"strings"
)

// StackFrame represents a single stack frame of an error.
type StackFrame struct {
	File string `json:"file" yaml:"file"`
	Func string `json:"func" yaml:"func"`
	Line int    `json:"line" yaml:"line"`
}

// ErrorStack holds an error and its relevant information.
type ErrorStack struct {
	error        error
	ErrorMessage string         `json:"error_message" yaml:"error_message"`
	Fields       map[string]any `json:"fields" yaml:"fields"`
	StackTrace   []StackFrame   `json:"stack_trace" yaml:"stack_trace"`
}

// Error returns the main error.
func (es *ErrorStack) Error() error {
	return es.error
}

// GetStack returns the errors that are present in the provided error.
func GetStack(err error) []ErrorStack {
	if err == nil {
		return nil
	}
	var collectedErrors []ErrorStack
	for err != nil {
		errorsToAdd := buildErrorStack(err)
		collectedErrors = append(collectedErrors, errorsToAdd)
		err = errors.Unwrap(err)
	}

	return cleanStack(collectedErrors)
}

func buildErrorStack(err error) ErrorStack {
	if serr, ok := err.(*Error); ok {
		return ErrorStack{
			error:        err,
			ErrorMessage: serr.message,
			Fields:       serr.fields,
			StackTrace:   resolveStackForStackFrames(serr.stack),
		}
	}
	return buildErrorStackForThirdPartyError(err)
}

func cleanStack(stackFrames []ErrorStack) []ErrorStack {
	// the following code only exists for cleaning up error formats.
	// e.g. when using pkg/errors.Wrap function two errors are added to the stack.
	// See github.com/pkg/errors@v0.9.1/errors.go:184
	// Before cleanup the ErrorStack would look like this:
	// [
	// 	{
	// 		"error_message": "serrors",
	// 		"fields": null,
	// 		"stack_trace": [{...}]
	// 	},
	// 	{
	// 		"error_message": "pkgerrors: errors",
	// 		"fields": null,
	// 		"stack_trace": [{...}]
	// 		]
	// 	},
	// 	{
	// 		"error_message": "pkgerrors: errors",
	// 		"fields": null,
	// 		"stack_trace": null
	// 	},
	// 	{
	// 		"error_message": "errors",
	// 		"fields": null,
	// 		"stack_trace": null
	// 	}
	// ]

	// remove errors from the list when they have no stack and have the same message as the next one
	for i := len(stackFrames) - 1; i > 0; i-- {
		if len(stackFrames[i].StackTrace) == 0 && len(stackFrames[i-1].StackTrace) > 0 &&
			stackFrames[i].ErrorMessage == stackFrames[i-1].ErrorMessage {
			appendToFields(&stackFrames[i-1].Fields, stackFrames[i].Fields)
			stackFrames = append(stackFrames[:i], stackFrames[i+1:]...)
		}
	}

	// clean up duplicate messages
	//nolint:gomnd // we only need to clear up when we have more than 1 stack frame
	if size := len(stackFrames); size > 1 {
		fullErrorText := ": " + stackFrames[size-1].ErrorMessage
		//nolint:gomnd //iterate backwards and start at the second last stack frame
		for i := size - 2; i >= 0; i-- {
			s, found := strings.CutSuffix(stackFrames[i].ErrorMessage, fullErrorText)
			if found {
				stackFrames[i].ErrorMessage = s
			}
			fullErrorText = ": " + stackFrames[i].ErrorMessage
		}
	}
	return stackFrames
}

//nolint:gocritic // allow passing a pointer to a map, we do this so we don't create a new map on every run.
func appendToFields(dst *map[string]any, src map[string]any) {
	if len(src) == 0 {
		return
	}
	if *dst == nil {
		*dst = make(map[string]any)
	}
	for s, a := range src {
		(*dst)[s] = a
	}
}
