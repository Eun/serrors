package serrors_test

import (
	"errors"
	"fmt"
	"net"
	"reflect"
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
	if expect, actual := true, ok; expect != actual {
		t.Fatalf(`expected %v, but was %v`, expect, actual)
	}

	t.Run("Errorf", func(t *testing.T) {
		err := serrors.Errorf("some error"). // [TestErrorErrorf00]
							With("key1", "value1").
							With("key2", "value2")
		if err == nil {
			t.Fatal(`expected not nil`)
		}
		if expect, actual := "some error", err.Error(); expect != actual {
			t.Fatalf(`expected %q, but was %q`, expect, actual)
		}

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
		if expect, actual := expectedFields, serrors.GetFields(err); !reflect.DeepEqual(expect, actual) {
			t.Fatalf(`expected %+v, but was %+v`, expect, actual)
		}
		CompareErrorStack(t, expectedStack, serrors.GetStack(err))
	})

	t.Run("Wrapf", func(t *testing.T) {
		err := testErrorFunc()                  // [TestErrorWrapf11]
		err = serrors.Wrapf(err, "some error"). // [TestErrorWrapf00] [TestErrorWrapf12]
							With("deep.key2", "value2").
							With("key1", "value1").
							With("key2", "value2")
		if err == nil {
			t.Fatal(`expected not nil`)
		}
		if expect, actual := "some error: deep error", err.Error(); expect != actual {
			t.Fatalf(`expected %q, but was %q`, expect, actual)
		}

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
		if expect, actual := expectedFields, serrors.GetFields(err); !reflect.DeepEqual(expect, actual) {
			t.Fatalf(`expected %+v, but was %+v`, expect, actual)
		}
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
			if expect, actual := tc.expectedFields, serrors.GetFields(tc.error); !reflect.DeepEqual(expect, actual) {
				t.Fatalf(`expected %+v, but was %+v`, expect, actual)
			}
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
			if expect, actual := tc.expectedArguments, serrors.GetFieldsAsCombinedSlice(tc.error); !reflect.DeepEqual(expect, actual) {
				t.Fatalf(`expected %+v, but was %+v`, expect, actual)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	t.Run("wrapped error", func(t *testing.T) {
		err1 := errors.New("error1")
		err := serrors.Wrap(err1, "error2")
		if expect, actual := err1, errors.Unwrap(err); !reflect.DeepEqual(expect, actual) {
			t.Fatalf(`expected %+v, but was %+v`, expect, actual)
		}
	})
	t.Run("wrapped no error", func(t *testing.T) {
		err := serrors.Wrap(nil, "error2")
		if errors.Unwrap(err) != nil {
			t.Fatal("expected not nil")
		}
	})
}

func TestIs(t *testing.T) {
	t.Run("wrapped error", func(t *testing.T) {
		err := serrors.Wrap(net.ErrClosed, "error")
		if expect, actual := true, errors.Is(err, net.ErrClosed); expect != actual {
			t.Fatalf(`expected %v, but was %v`, expect, actual)
		}
		if expect, actual := false, errors.Is(err, net.ErrWriteToConnected); expect != actual {
			t.Fatalf(`expected %v, but was %v`, expect, actual)
		}
	})
	t.Run("wrapped no error", func(t *testing.T) {
		err := serrors.Wrap(nil, "error")
		if expect, actual := false, errors.Is(err, net.ErrClosed); expect != actual {
			t.Fatalf(`expected %v, but was %v`, expect, actual)
		}
		if expect, actual := false, errors.Is(err, net.ErrWriteToConnected); expect != actual {
			t.Fatalf(`expected %v, but was %v`, expect, actual)
		}
	})
}

func TestAs(t *testing.T) {
	t.Run("wrapped error", func(t *testing.T) {
		err := serrors.Wrap(&net.AddrError{Addr: "127.0.0.1"}, "error")

		var cause1 *net.AddrError
		if expect, actual := true, errors.As(err, &cause1); expect != actual {
			t.Fatalf(`expected %v, but was %v`, expect, actual)
		}
		if expect, actual := "127.0.0.1", cause1.Addr; expect != actual {
			t.Fatalf(`expected %q, but was %q`, expect, actual)
		}

		var cause2 *net.OpError
		if expect, actual := false, errors.As(err, &cause2); expect != actual {
			t.Fatalf(`expected %v, but was %v`, expect, actual)
		}
		if cause2 != nil {
			t.Fatal("expected not nil")
		}
	})
	t.Run("wrapped no error", func(t *testing.T) {
		err := serrors.Wrap(nil, "error")
		var cause1 *net.AddrError
		if expect, actual := false, errors.As(err, &cause1); expect != actual {
			t.Fatalf(`expected %v, but was %v`, expect, actual)
		}
		if cause1 != nil {
			t.Fatal("expected nil")
		}

		var cause2 *net.OpError

		if expect, actual := false, errors.As(err, &cause2); expect != actual {
			t.Fatalf(`expected %v, but was %v`, expect, actual)
		}
		if cause2 != nil {
			t.Fatal("expected nil")
		}
	})
}
