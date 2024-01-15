package serrors_test

import (
	"fmt"

	"github.com/Eun/serrors"
)

func ExampleNewBuilder() {
	name := "Joe"
	serr := serrors.NewBuilder().
		With("username", name)

	if name == "alice" {
		panic(serr.New("username cannot be alice"))
	}
	fmt.Println(name)
	// Output:
	// Joe
}

func ExampleErrorBuilder_With() {
	name := "Joe"
	serr := serrors.NewBuilder()
	serr = serr.With("username", name)
	if name == "alice" {
		panic(serr.New("username cannot be alice"))
	}
	fmt.Println(name)
	// Output:
	// Joe
}

func ExampleErrorBuilder_New() {
	name := "Joe"
	serr := serrors.NewBuilder().
		With("username", name)
	if name == "alice" {
		panic(serr.New("username cannot be alice"))
	}
	fmt.Println(name)
	// Output:
	// Joe
}

func ExampleErrorBuilder_Errorf() {
	name := "Joe"
	serr := serrors.NewBuilder().
		With("username", name)
	if name == "alice" {
		panic(serr.Errorf("username cannot be %q", name))
	}
	fmt.Println(name)
	// Output:
	// Joe
}

func ExampleErrorBuilder_Wrap() {
	name := "Joe"
	serr := serrors.NewBuilder().
		With("username", name)
	if err := validateUserName(name); err != nil {
		panic(serr.Wrap(err, "validation of username failed"))
	}
	fmt.Println(name)
	// Output:
	// Joe
}

func ExampleErrorBuilder_Wrapf() {
	name := "Joe"
	serr := serrors.NewBuilder().
		With("username", name)
	if err := validateUserName(name); err != nil {
		panic(serr.Wrapf(err, "validation of username %q failed", name))
	}
	fmt.Println(name)
	// Output:
	// Joe
}
