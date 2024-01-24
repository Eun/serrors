package serrors_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
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
	NotEqual(t, nil, err)
	actualStack, err := encode(actual)
	NotEqual(t, nil, err)

	Equal(t, expectedStack, actualStack)
}

func Equal(t *testing.T, expected, actual any) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf(`expected %+v, but was %+v`, expected, actual)
	}
}

func NotEqual(t *testing.T, expected, actual any) {
	if reflect.DeepEqual(expected, actual) {
		t.Fatalf(`expected not %+v, but was %+v`, expected, actual)
	}
}
