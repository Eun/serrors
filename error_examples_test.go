package serrors_test

import (
	"fmt"

	"github.com/Eun/serrors"
)

func ExampleNew() {
	name := "Joe"
	if name == "alice" {
		panic(serrors.New("invalid name").With("name", name))
	}
	fmt.Println(name)
	// Output:
	// Joe
}

func ExampleErrorf() {
	name := "Joe"
	if name == "alice" {
		panic(serrors.Errorf("invalid name %q", name).With("name", name))
	}
	fmt.Println(name)
	// Output:
	// Joe
}

func ExampleWrap() {
	name := "Joe"
	if err := validateUserName(name); err != nil {
		panic(serrors.Wrap(err, "invalid name").With("name", name))
	}
	fmt.Println(name)
	// Output:
	// Joe
}

func ExampleWrapf() {
	name := "Joe"
	if err := validateUserName(name); err != nil {
		panic(serrors.Wrapf(err, "invalid name %q", name).With("name", name))
	}
	fmt.Println(name)
	// Output:
	// Joe
}

func ExampleWith() {
	name := "Joe"
	if name == "alice" {
		panic(serrors.New("invalid name").With("name", name))
	}
	fmt.Println(name)
	// Output:
	// Joe
}
