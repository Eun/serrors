package serrors_test

import (
	"log/slog"
	"slices"

	"github.com/Eun/serrors"
)

func validateUserNameLength(name string) error {
	const maxLength = 10
	if len(name) > maxLength {
		return serrors.New("username is too long").
			With("username", name).
			With("max_length", maxLength)
	}
	return nil
}

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

func addUserToRole(userName, roleName string) error {
	if err := validateUserName(userName); err != nil {
		return serrors.Wrap(err, "validation of username failed").With("username", userName)
	}

	if roleName == "" {
		return serrors.New("rolename cannot be empty")
	}
	availableRoles := []string{"admin", "user"}
	if !slices.Contains(availableRoles, roleName) {
		return serrors.Errorf("unknown role %q", roleName).
			With("username", userName).
			With("available_roles", availableRoles)
	}

	// todo: add user to role
	// ...

	return nil
}

func ExampleError() {
	if err := addUserToRole("joe", "guest"); err != nil {
		slog.Error("name validation failed",
			"error", err.Error(),
			slog.Group("details", serrors.GetFieldsAsCombinedSlice(err)...),
			"stack", serrors.GetStack(err),
		)
		return
	}
}
