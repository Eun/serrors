package serrors

import (
	"errors"
	"fmt"
)

// Error is the error that will be returned by ErrorBuilder and the functions
// New, Errorf, Wrap and Wrapf.
// It implements the stdlib error interface.
type Error struct {
	message string
	cause   error
	fields  map[string]any
	stack   []uintptr
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *Error) Unwrap() error { return e.cause }

// Cause returns the cause of this error.
func (e *Error) Cause() error { return e.cause }

// With adds the field key with the value to the error fields.
func (e *Error) With(key string, value any) *Error {
	if e.fields == nil {
		e.fields = make(map[string]any)
	}
	e.fields[key] = value
	return e
}

// New creates a new Error with the supplied message.
func New(message string) *Error {
	return &Error{
		message: message,
		cause:   nil,
		fields:  nil,
		stack:   collectStack(),
	}
}

// Errorf creates a new Error with the supplied message formatted according to a format specifier.
func Errorf(format string, a ...any) *Error {
	return New(fmt.Sprintf(format, a...))
}

// Wrap creates a new Error with the supplied message.
// The passed in error will be added as a cause for this error.
func Wrap(err error, message string) *Error {
	return &Error{
		message: message,
		cause:   err,
		fields:  nil,
		stack:   collectStack(),
	}
}

// Wrapf creates a new Error with the supplied message formatted according to a format specifier.
// The passed in error will be added as a cause for this error.
func Wrapf(err error, format string, a ...any) *Error {
	return Wrap(err, fmt.Sprintf(format, a...))
}

// GetFields will return all fields that are added to the specified error.
func GetFields(err error) map[string]any {
	if err == nil {
		return nil
	}

	// errors are nested, and we don't want nested errors to
	// overwrite fields of the errors above
	// We need to collect the fields first and then
	// reverse iterate them and add the key-values to the final map
	var collectedFields []map[string]any
	for err != nil {
		if e, ok := err.(*Error); ok {
			if len(e.fields) > 0 {
				collectedFields = append(collectedFields, e.fields)
			}
		}
		err = errors.Unwrap(err)
	}

	if len(collectedFields) == 0 {
		return nil
	}

	fields := make(map[string]any)
	for i := len(collectedFields) - 1; i >= 0; i-- {
		for k, v := range collectedFields[i] {
			fields[k] = v
		}
	}
	return fields
}

// GetFieldsAsCombinedSlice will return all fields as a slice that are added to the specified error.
// The format will be [key1, value1, key2, value2, ..., keyN, valueN].
func GetFieldsAsCombinedSlice(err error) []any {
	fields := GetFields(err)
	if fields == nil {
		return nil
	}
	args := make([]any, 0, len(fields))
	for s, a := range fields {
		args = append(args, s, a)
	}
	return args
}
