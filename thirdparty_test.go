package serrors_test

import (
	"errors"
	"fmt"
	"runtime"
	"testing"

	pkgerrors "github.com/pkg/errors"

	"github.com/Eun/serrors"
)

func TestGetStack_WithPkgErrors(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	Equal(t, true, ok)

	err := serrors.Wrap( //  [TestGetStack_WithPkgErrors00]
		pkgerrors.Wrap( //  [TestGetStack_WithPkgErrors10]
			errors.New("errors"),
			"pkgerrors"),
		"serrors",
	)

	expectedStack := []serrors.ErrorStack{
		{
			ErrorMessage: "serrors",
			StackTrace: []serrors.StackFrame{
				buildStackFrameFromMarker(t, filename, "TestGetStack_WithPkgErrors00"),
			},
		},
		{
			ErrorMessage: "pkgerrors",
			StackTrace: []serrors.StackFrame{
				buildStackFrameFromMarker(t, filename, "TestGetStack_WithPkgErrors10"),
			},
		},
		{
			ErrorMessage: "errors",
			StackTrace:   nil,
		},
	}

	CompareErrorStack(t, expectedStack, serrors.GetStack(err))
}

func TestError_Format_WithPkgErrors(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	Equal(t, true, ok)

	err := serrors.Wrap( //  [TestError_Format_WithPkgErrors01]
		pkgerrors.Wrap( //  [TestError_Format_WithPkgErrors02]
			errors.New("errors"),
			"pkgerrors"),
		"serrors",
	)

	expected := fmt.Sprintf("serrors\n%s\npkgerrors\n%s\nerrors\n",
		generateExpectedStack(t, filename, "TestError_Format_WithPkgErrors01"),
		generateExpectedStack(t, filename, "TestError_Format_WithPkgErrors02"),
	)
	Equal(t, expected, fmt.Sprintf("%+v", err))
}

func TestError_Error_WithPkgErrors(t *testing.T) {
	err := serrors.Wrap(
		pkgerrors.Wrap(
			errors.New("errors"),
			"pkgerrors"),
		"serrors",
	)

	Equal(t, "serrors: pkgerrors: errors", err.Error())
}
