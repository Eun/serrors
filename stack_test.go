package serrors_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/Eun/serrors"
)

func TestGetStack(t *testing.T) {
	// just generic testing
	// the real test are happening in TestError & TestBuilder.
	t.Run("get stack from nil value", func(t *testing.T) {
		if serrors.GetStack(nil) != nil {
			t.Fatal("expected nil")
		}
	})
	t.Run("get original error from stack", func(t *testing.T) {
		err := serrors.New("some error")

		stack := serrors.GetStack(err)
		if expect, actual := 1, len(stack); expect != actual {
			t.Fatalf(`expected %v, but was %v`, expect, actual)
		}
		if expect, actual := err, stack[0].Error(); !reflect.DeepEqual(expect, actual) {
			t.Fatalf(`expected %+v, but was %+v`, expect, actual)
		}
		if expect, actual := err.Error(), stack[0].ErrorMessage; expect != actual {
			t.Fatalf(`expected %q, but was %q`, expect, actual)
		}
	})
}

func buildStackFrameFromMarker(t *testing.T, fileName, marker string) serrors.StackFrame {
	marker = "[" + marker + "]"
	var result *serrors.StackFrame
	fileSet := token.NewFileSet()

	// Parse the Go file
	file, err := parser.ParseFile(fileSet, fileName, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		t.Fatal(`expected no error`)
	}
	packageName := "github.com/Eun/" + file.Name.Name

	var inspectNode func(n ast.Node) bool
	//nolint:unparam // the ast.Inspect node dictates a bool as a result, even when its true
	inspectNode = func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.File:
			for _, comment := range v.Comments {
				inspectNode(comment)
			}
		case *ast.CommentGroup:
			for _, comment := range v.List {
				inspectNode(comment)
			}
		case *ast.Comment:
			if strings.Contains(v.Text, marker) {
				// Get the function name and line number
				pos := fileSet.Position(v.Slash)
				funcName := findEnclosingFunc(fileSet, file, pos.Offset)
				if result != nil {
					t.Fatalf(`found a marker %q already`, marker)
				}
				result = &serrors.StackFrame{
					File: fileName,
					Func: packageName + "." + funcName,
					Line: pos.Line,
				}
				return true
			}
		}
		return true
	}

	// Traverse the AST
	ast.Inspect(file, inspectNode)
	if result == nil {
		t.Fatalf(`no marker found %q`, marker)
	}
	return *result
}

// helper function to find the enclosing function declaration based on position offset.
func findEnclosingFunc(fileSet *token.FileSet, file *ast.File, offset int) string {
	var result []string
	anonCounter := 1
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			pos := fileSet.Position(node.Pos()).Offset
			end := fileSet.Position(node.End()).Offset

			if pos <= offset && offset <= end {
				result = []string{node.Name.Name}
			}
		case *ast.FuncLit:
			pos := fileSet.Position(node.Pos()).Offset
			end := fileSet.Position(node.End()).Offset

			if pos <= offset && offset <= end {
				result = append(result, "func"+strconv.Itoa(anonCounter))
			}
			anonCounter++
		}

		return true
	})

	return strings.Join(result, ".")
}
