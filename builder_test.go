package serrors_test

import (
	"runtime"
	"testing"

	"github.com/Eun/serrors"
)

func testBuilderErrorFunc() error {
	return serrors.New("deep error"). // [TestBuilderWrapf10]
						With("deep.key1", "value1").
						With("deep.key2", "should be overwritten")
}

func TestBuilder(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	Equal(t, true, ok)

	t.Run("Errorf", func(t *testing.T) {
		errorBuilder := serrors.NewBuilder().
			With("key1", "value1").
			With("key2", "should be overwritten")
		err := errorBuilder.Errorf("some error"). // [TestBuilderErrorf00]
								With("key2", "value2").
								With("key3", "value3")
		NotEqual(t, nil, err)
		Equal(t, "some error", err.Error())

		expectedFields := map[string]any{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		expectedStack := []serrors.ErrorStack{
			{
				ErrorMessage: "some error",
				Fields:       expectedFields,
				StackTrace: []serrors.StackFrame{
					buildStackFrameFromMarker(t, filename, "TestBuilderErrorf00"),
				},
			},
		}
		Equal(t, expectedFields, serrors.GetFields(err))
		CompareErrorStack(t, expectedStack, serrors.GetStack(err))
	})

	t.Run("Wrapf", func(t *testing.T) {
		errorBuilder := serrors.NewBuilder().
			With("key1", "value1").
			With("key2", "should be overwritten")

		err := testBuilderErrorFunc()                // [TestBuilderWrapf11]
		err = errorBuilder.Wrapf(err, "some error"). // [TestBuilderWrapf00] [TestBuilderWrapf12]
								With("deep.key2", "value2").
								With("key2", "value2").
								With("key3", "value3")
		NotEqual(t, nil, err)
		Equal(t, "some error: deep error", err.Error())

		expectedFields := map[string]any{
			"deep.key1": "value1",
			"deep.key2": "value2",
			"key1":      "value1",
			"key2":      "value2",
			"key3":      "value3",
		}
		expectedStack := []serrors.ErrorStack{
			{
				ErrorMessage: "some error",
				Fields: map[string]any{
					"deep.key2": "value2",
					"key1":      "value1",
					"key2":      "value2",
					"key3":      "value3",
				},
				StackTrace: []serrors.StackFrame{
					buildStackFrameFromMarker(t, filename, "TestBuilderWrapf00"),
				},
			},
			{
				ErrorMessage: "deep error",
				Fields: map[string]any{
					"deep.key1": "value1",
					"deep.key2": "should be overwritten",
				},
				StackTrace: []serrors.StackFrame{
					buildStackFrameFromMarker(t, filename, "TestBuilderWrapf10"),
					buildStackFrameFromMarker(t, filename, "TestBuilderWrapf11"),
				},
			},
		}
		Equal(t, expectedFields, serrors.GetFields(err))
		CompareErrorStack(t, expectedStack, serrors.GetStack(err))
	})
}
