package serrors_test

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

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
	require.True(t, ok)

	t.Run("Errorf", func(t *testing.T) {
		err := serrors.Errorf("some error"). // [TestErrorErrorf00]
							With("key1", "value1").
							With("key2", "value2")

		require.NotNil(t, err)
		require.Equal(t, "some error", err.Error())

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
		require.Equal(t, expectedFields, serrors.GetFields(err))
		CompareErrorStack(t, expectedStack, serrors.GetStack(err))
	})

	t.Run("Wrapf", func(t *testing.T) {
		err := testErrorFunc()                  // [TestErrorWrapf11]
		err = serrors.Wrapf(err, "some error"). // [TestErrorWrapf00] [TestErrorWrapf12]
							With("deep.key2", "value2").
							With("key1", "value1").
							With("key2", "value2")
		require.NotNil(t, err)
		require.Equal(t, "some error: deep error", err.Error())

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
		require.Equal(t, expectedFields, serrors.GetFields(err))
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
			require.Equal(t, tc.expectedFields, serrors.GetFields(tc.error))
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
			require.Equal(t, tc.expectedArguments, serrors.GetFieldsAsCombinedSlice(tc.error))
		})
	}
}

func TestUnwrap(t *testing.T) {
	t.Run("wrapped error", func(t *testing.T) {
		err1 := errors.New("error1")
		err := serrors.Wrap(err1, "error2")
		require.Equal(t, err1, errors.Unwrap(err))
	})
	t.Run("wrapped no error", func(t *testing.T) {
		err := serrors.Wrap(nil, "error2")
		require.Nil(t, errors.Unwrap(err))
	})
}

func TestIs(t *testing.T) {
	t.Run("wrapped error", func(t *testing.T) {
		err := serrors.Wrap(net.ErrClosed, "error")
		require.True(t, errors.Is(err, net.ErrClosed))
		require.False(t, errors.Is(err, net.ErrWriteToConnected))
	})
	t.Run("wrapped no error", func(t *testing.T) {
		err := serrors.Wrap(nil, "error")
		require.False(t, errors.Is(err, net.ErrClosed))
		require.False(t, errors.Is(err, net.ErrWriteToConnected))
	})
}

func TestAs(t *testing.T) {
	t.Run("wrapped error", func(t *testing.T) {
		err := serrors.Wrap(&net.AddrError{Addr: "127.0.0.1"}, "error")

		var cause1 *net.AddrError
		require.True(t, errors.As(err, &cause1))
		require.Equal(t, cause1.Addr, "127.0.0.1")

		var cause2 *net.OpError
		require.False(t, errors.As(err, &cause2))
		require.Nil(t, cause2)
	})
	t.Run("wrapped no error", func(t *testing.T) {
		err := serrors.Wrap(nil, "error")
		var cause1 *net.AddrError
		require.False(t, errors.As(err, &cause1))
		require.Nil(t, cause1)

		var cause2 *net.OpError
		require.False(t, errors.As(err, &cause2))
		require.Nil(t, cause2)
	})
}
