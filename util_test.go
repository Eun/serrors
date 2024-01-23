package serrors_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/Eun/serrors"
)

func NewSLogLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			return a
		},
	}))
}

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
	if err != nil {
		t.Fatal(`expected no error`)
	}
	actualStack, err := encode(actual)
	if err != nil {
		t.Fatal(`expected no error`)
	}

	if expect, actual := expectedStack, actualStack; !reflect.DeepEqual(expect, actual) {
		t.Fatalf(`expected %+v, but was %+v`, expect, actual)
	}
}
