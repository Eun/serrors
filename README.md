# serrors - Structured Errors
[![Actions Status](https://github.com/Eun/serrors/workflows/push/badge.svg)](https://github.com/Eun/serrors/actions)
[![Coverage Status](https://coveralls.io/repos/github/Eun/serrors/badge.svg?branch=main)](https://coveralls.io/github/Eun/serrors?branch=main)
[![PkgGoDev](https://img.shields.io/badge/pkg.go.dev-reference-blue)](https://pkg.go.dev/github.com/Eun/serrors)
[![go-report](https://goreportcard.com/badge/github.com/Eun/serrors)](https://goreportcard.com/report/github.com/Eun/serrors)
---
*serrors* allows you to add tags/fields/kv-pairs to your errors.

## Usage
```go
package main

import (
	"os"
	"fmt"
	"log/slog"
	
	"github.com/Eun/serrors"
)

func validateUserName(name string) error {
	const maxLength = 10
	if len(name) > maxLength {
		return serrors.New("username is too long").
			With("username", name).
			With("max_length", maxLength)
	}
	return nil
}

func main() {
	user := os.Getenv("USER")
	err := validateUserName(user)
	if err != nil {
		slog.Error("name validation failed",
			"error", err.Error(),
			slog.Group("details", serrors.GetFieldsAsCombinedSlice(err)...),
			"stack", serrors.GetStack(err),
		)
		return
	}
	fmt.Println("Welcome ", user)
}

```


## Problem
We use structured loggers like *slog* to create nice formatted log messages  
and add important context to error messages.  
Take this code as an example:
```go
func validateUserNameLength(name string) error {
	const maxLength = 10
	if len(name) > maxLength {
		slog.Error("username is too long", "username", name, "max_length", maxLength)
		return errors.New("username is too long")
	}
	return nil
}
```
Not only do we return an error, but we also log the error using *slog*.  
Lets look at the calling function:

```go
func addUserToRole(userName, roleName string) error {
	if err := validateUserNameLength(userName); err != nil {
		slog.Error("validation of username failed", "username", name)
		return fmt.Errorf("validation of username failed: %w", err)
	}
	// ...
}
```
Again, we return the error (with the underlying error), and we also
log it - because we need the context in our messages.

In this case we end up with at least two error messages:
1. `slog.Error("username is too long", ...)`
2. `slog.Error("validation of username failed", ...)`
3. when we handle `addUserToRole`: `validation of username failed: username is too long`

The last error that will be logged or printed won't contain any useful information
on why this problem actually occurred.

One possible solution would be to use something like `fmt.Errorf("username is too long [username=%s]", name)`.
However, this could lead to some funny unreadable errors like:
```
validation of username failed [username=MisterDolittle]: username is too long [username=MisterDolittle] [max_length=10]
```

This package attempts to solve this problem by providing methods to add tags/fields/kv-pairs to errors that can later be
retrieved.


## Builder Usage
You could save some code duplication by using the builder functionality:
```go
func validateUserName(name string) error {
	serr := serrors.NewBuilder().
		With("username", name)

	if name == "" {
		return serr.New("username cannot be empty")
	}

	if err := validateUserNameLength(name); err != nil {
		return serr.Wrap(err, "username has invalid length")
	}

	reservedNames := []string{"root", "admin"}
	for _, s := range reservedNames {
		if name == s {
			return serr.Errorf("username cannot be %q, it is reserved", name).
				With("reserved", reservedNames)
		}
	}
	return nil
}
```

## Building without Stack
By default *serrors* collects stack information, this behaviour can be disabled by
setting the build tag `serrors_without_stack`:
```shell
go build -tags serrors_without_stack ...
```
