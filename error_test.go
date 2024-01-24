package serrors_test

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"testing"

	"github.com/Eun/serrors"
)

var _ error = &serrors.Error{} // make sure we implement the error interface

func testErrorFunc() error {
	return serrors.New("deep error"). // [TestErrorWrapf10]
						With("deep.key1", "value1").
						With("deep.key2", "should be overwritten")
}

func TestError(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	Equal(t, true, ok)

	t.Run("Errorf", func(t *testing.T) {
		err := serrors.Errorf("some error"). // [TestErrorErrorf00]
							With("key1", "value1").
							With("key2", "value2")
		NotEqual(t, nil, err)
		Equal(t, "some error", err.Error())

		expectedFields := map[string]any{
			"key1": "value1",
			"key2": "value2",
		}
		expectedStack := []serrors.ErrorStack{
			{
				ErrorMessage: "some error",
				Fields:       expectedFields,
				StackTrace: []serrors.StackFrame{
					buildStackFrameFromMarker(t, filename, "TestErrorErrorf00"),
				},
			},
		}
		Equal(t, expectedFields, serrors.GetFields(err))
		CompareErrorStack(t, expectedStack, serrors.GetStack(err))
	})

	t.Run("Wrapf", func(t *testing.T) {
		err := testErrorFunc()                  // [TestErrorWrapf11]
		err = serrors.Wrapf(err, "some error"). // [TestErrorWrapf00] [TestErrorWrapf12]
							With("deep.key2", "value2").
							With("key1", "value1").
							With("key2", "value2")
		NotEqual(t, nil, err)
		Equal(t, "some error: deep error", err.Error())

		expectedFields := map[string]any{
			"deep.key1": "value1",
			"deep.key2": "value2",
			"key1":      "value1",
			"key2":      "value2",
		}
		expectedStack := []serrors.ErrorStack{
			{
				ErrorMessage: "some error",
				Fields: map[string]any{
					"deep.key2": "value2",
					"key1":      "value1",
					"key2":      "value2",
				},
				StackTrace: []serrors.StackFrame{
					buildStackFrameFromMarker(t, filename, "TestErrorWrapf00"),
				},
			},
			{
				ErrorMessage: "deep error",
				Fields: map[string]any{
					"deep.key1": "value1",
					"deep.key2": "should be overwritten",
				},
				StackTrace: []serrors.StackFrame{
					buildStackFrameFromMarker(t, filename, "TestErrorWrapf10"),
					buildStackFrameFromMarker(t, filename, "TestErrorWrapf11"),
				},
			},
		}
		Equal(t, expectedFields, serrors.GetFields(err))
		CompareErrorStack(t, expectedStack, serrors.GetStack(err))
	})
}

func TestGetFields(t *testing.T) {
	testCases := []struct {
		name           string
		error          error
		expectedFields map[string]any
	}{
		{
			name:           "error is nil",
			error:          nil,
			expectedFields: nil,
		},
		{
			name:           "not an error",
			error:          errors.New("some error"),
			expectedFields: nil,
		},
		{
			name:           "not containing fields",
			error:          serrors.New("some error"),
			expectedFields: nil,
		},
		{
			name:           "regular",
			error:          serrors.New("").With("k", "v"),
			expectedFields: map[string]any{"k": "v"},
		},
		{
			name:           "wrapped",
			error:          fmt.Errorf("some error: %w", serrors.New("some error").With("k", "v")),
			expectedFields: map[string]any{"k": "v"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			Equal(t, tc.expectedFields, serrors.GetFields(tc.error))
		})
	}
}

func TestGetFieldsAsArguments(t *testing.T) {
	testCases := []struct {
		name              string
		error             error
		expectedArguments []any
	}{
		{
			name:              "error is nil",
			error:             nil,
			expectedArguments: nil,
		},
		{
			name:              "not an error",
			error:             errors.New("some error"),
			expectedArguments: nil,
		},
		{
			name:              "not containing fields",
			error:             serrors.New("some error"),
			expectedArguments: nil,
		},
		{
			name:              "regular",
			error:             serrors.New("").With("k", "v"),
			expectedArguments: []any{"k", "v"},
		},
		{
			name:              "wrapped",
			error:             fmt.Errorf("some error: %w", serrors.New("some error").With("k", "v")),
			expectedArguments: []any{"k", "v"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			Equal(t, tc.expectedArguments, serrors.GetFieldsAsCombinedSlice(tc.error))
		})
	}
}

func TestUnwrap(t *testing.T) {
	t.Run("wrapped error", func(t *testing.T) {
		err1 := errors.New("error1")
		err := serrors.Wrap(err1, "error2")
		Equal(t, err1, errors.Unwrap(err))
	})
	t.Run("wrapped no error", func(t *testing.T) {
		err := serrors.Wrap(nil, "error2")
		NotEqual(t, nil, errors.Unwrap(err))
	})
}

func TestIs(t *testing.T) {
	t.Run("wrapped error", func(t *testing.T) {
		err := serrors.Wrap(net.ErrClosed, "error")
		Equal(t, true, errors.Is(err, net.ErrClosed))
		Equal(t, false, errors.Is(err, net.ErrWriteToConnected))
	})
	t.Run("wrapped no error", func(t *testing.T) {
		err := serrors.Wrap(nil, "error")
		Equal(t, false, errors.Is(err, net.ErrClosed))
		Equal(t, false, errors.Is(err, net.ErrWriteToConnected))
	})
}

func TestAs(t *testing.T) {
	t.Run("wrapped error", func(t *testing.T) {
		err := serrors.Wrap(&net.AddrError{Addr: "127.0.0.1"}, "error")

		var cause1 *net.AddrError
		Equal(t, true, errors.As(err, &cause1))
		Equal(t, "127.0.0.1", cause1.Addr)

		var cause2 *net.OpError
		Equal(t, false, errors.As(err, &cause2))
		NotEqual(t, nil, cause2)
	})
	t.Run("wrapped no error", func(t *testing.T) {
		err := serrors.Wrap(nil, "error")
		var cause1 *net.AddrError
		Equal(t, false, errors.As(err, &cause1))
		NotEqual(t, nil, cause1)

		var cause2 *net.OpError
		Equal(t, false, errors.As(err, &cause2))
		NotEqual(t, nil, cause2)
	})
}
