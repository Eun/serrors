package serrors_test

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/Eun/serrors"
)

func TestError_Error(t *testing.T) {
	testCases := []struct {
		name              string
		error             error
		expectedErrorText string
	}{
		{
			name:              "normal error",
			error:             serrors.New("some error"),
			expectedErrorText: "some error",
		},
		{
			name:              "normal error with fields",
			error:             serrors.New("some error").With("k", "v"),
			expectedErrorText: "some error",
		},
		{
			name:              "wrapped",
			error:             fmt.Errorf("error1: %w", serrors.New("error2").With("k", "v")),
			expectedErrorText: "error1: error2[k=v]",
		},
		{
			name:              "wraps another error",
			error:             serrors.Wrap(errors.New("error2"), "error1"),
			expectedErrorText: "error1: error2",
		},
		{
			name:              "no error text",
			error:             serrors.New(""),
			expectedErrorText: "error",
		},
		{
			name:              "wraps a nil error",
			error:             serrors.Wrap(nil, "some error"),
			expectedErrorText: "some error",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if expect, actual := tc.expectedErrorText, tc.error.Error(); expect != actual {
				t.Fatalf(`expected %q, but was %q`, expect, actual)
			}
		})
	}
}

func TestError_Format(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if expect, actual := true, ok; expect != actual {
		t.Fatalf(`expected %v, but was %v`, expect, actual)
	}

	cause := serrors.New("error 1").With("k1", "v1")       // [TestError_Format01]
	err := serrors.Wrap(cause, "error 2").With("k2", "v2") // [TestError_Format00]

	t.Run("normal string", func(t *testing.T) {
		if expect, actual := "error 2: error 1", fmt.Sprintf("%s", err); expect != actual {
			t.Fatalf(`expected %q, but was %q`, expect, actual)
		}
	})
	t.Run("quoted string", func(t *testing.T) {
		if expect, actual := `"error 2: error 1"`, fmt.Sprintf("%q", err); expect != actual {
			t.Fatalf(`expected %q, but was %q`, expect, actual)
		}
	})
	t.Run("verbose", func(t *testing.T) {
		if expect, actual := "error 2: error 1[k1=v1 k2=v2]", fmt.Sprintf("%v", err); expect != actual {
			t.Fatalf(`expected %q, but was %q`, expect, actual)
		}
	})
	t.Run("extra verbose", func(t *testing.T) {
		expected := fmt.Sprintf("error 2\n[k2=v2]\n%s\nerror 1\n[k1=v1]\n%s\n",
			generateExpectedStack(t, filename, "TestError_Format00"),
			generateExpectedStack(t, filename, "TestError_Format01"),
		)
		if expect, actual := expected, fmt.Sprintf("%+v", err); expect != actual {
			t.Fatalf(`expected %q, but was %q`, expect, actual)
		}
	})
}

func generateExpectedStack(t *testing.T, filename string, markers ...string) string {
	parts := make([]string, len(markers))
	for i, marker := range markers {
		frame := buildStackFrameFromMarker(t, filename, marker)
		parts[i] = fmt.Sprintf("%s\n\t%s:%d", frame.Func, frame.File, frame.Line)
	}
	return strings.Join(parts, "\n")
}
