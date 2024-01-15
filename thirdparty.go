package serrors

import (
	pkgerrors "github.com/pkg/errors"
)

func buildErrorStackForThirdPartyError(err error) ErrorStack {
	es := ErrorStack{
		error:        err,
		ErrorMessage: err.Error(),
		Fields:       nil,
		StackTrace:   nil,
	}
	// pkgerrors
	{
		type stackTracer interface {
			StackTrace() pkgerrors.StackTrace
		}

		if tracer, ok := err.(stackTracer); ok {
			stackTrace := tracer.StackTrace()
			result := make([]uintptr, len(stackTrace))
			for i, frame := range stackTrace {
				result[i] = uintptr(frame)
			}
			es.StackTrace = resolveStackForStackFrames(result)
		}
	}

	return es
}
