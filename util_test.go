package serrors_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"slices"
	"testing"

	"github.com/Eun/serrors"
)

func CompareErrorStack(t *testing.T, expected, actual []serrors.ErrorStack) {
	encode := func(v any) (string, error) {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "\t")
		if err := enc.Encode(v); err != nil {
			return "", serrors.Wrap(err, "unable to encode").With("obj", fmt.Sprintf("%+v", v))
		}
		return buf.String(), nil
	}

	expectedStack, err := encode(expected)
	Nil(t, err)
	actualStack, err := encode(actual)
	Nil(t, err)

	Equal(t, expectedStack, actualStack)
}

func Equal(t *testing.T, expected, actual any) {
	if expected == nil && actual == nil {
		return
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %+v, but was %+v\n%s", expected, actual, string(debug.Stack()))
	}
}

func NotEqual(t *testing.T, expected, actual any) {
	if expected == nil && actual == nil {
		t.Fatalf("expected not %+v, but was %+v\n%s", expected, actual, string(debug.Stack()))
	}
	if reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected not %+v, but was %+v\n%s", expected, actual, string(debug.Stack()))
	}
}

func isNil(actual any) bool {
	if actual == nil {
		return true
	}

	value := reflect.ValueOf(actual)
	kind := value.Kind()
	isNilableKind := slices.Contains(
		[]reflect.Kind{
			reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice, reflect.UnsafePointer},
		kind)
	if isNilableKind && value.IsNil() {
		return true
	}
	return false
}

func Nil(t *testing.T, actual any) {
	if !isNil(actual) {
		t.Fatalf("expected %+v to be nil\n%s", actual, string(debug.Stack()))
	}
}

func NotNil(t *testing.T, actual any) {
	if isNil(actual) {
		t.Fatalf("expected %+v to be not nil\n%s", actual, string(debug.Stack()))
	}
}
