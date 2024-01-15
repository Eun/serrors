package serrors

import "fmt"

// ErrorBuilder is a type that provides a way to build errors.
type ErrorBuilder struct {
	msg    string
	fields map[string]any
}

// NewBuilder creates a new ErrorBuilder.
// To build the error use either ErrorBuilder.New, ErrorBuilder.Errorf, ErrorBuilder.Wrap
// or ErrorBuilder.Wrapf.
func NewBuilder() *ErrorBuilder {
	return &ErrorBuilder{
		msg:    "",
		fields: nil,
	}
}

// With adds the field key with the value to the error fields.
func (eb *ErrorBuilder) With(key string, value any) *ErrorBuilder {
	if eb.fields == nil {
		eb.fields = make(map[string]any)
	}
	eb.fields[key] = value
	return eb
}

// New creates a new Error with the supplied message. The error
// will contain all fields that were previously passed to ErrorBuilder.
func (eb *ErrorBuilder) New(message string) *Error {
	eb.msg = message
	err := New(message)
	err.fields = eb.fields
	return err
}

// Errorf creates a new Error with the supplied message formatted according to a format specifier.
// The error will contain all fields that were previously passed to ErrorBuilder.
func (eb *ErrorBuilder) Errorf(format string, a ...any) *Error {
	return eb.New(fmt.Sprintf(format, a...))
}

// Wrap creates a new Error with the supplied message.
// The passed in error will be added as a cause for this error.
// The error will contain all fields that were previously passed to ErrorBuilder.
func (eb *ErrorBuilder) Wrap(err error, message string) *Error {
	serr := Wrap(err, message)
	serr.fields = eb.fields
	return serr
}

// Wrapf creates a new Error with the supplied message formatted according to a format specifier.
// The passed in error will be added as a cause for this error.
// The error will contain all fields that were previously passed to ErrorBuilder.
func (eb *ErrorBuilder) Wrapf(err error, format string, a ...any) *Error {
	return eb.Wrap(err, fmt.Sprintf(format, a...))
}
