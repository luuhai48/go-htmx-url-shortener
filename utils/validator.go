package utils

import (
	"errors"

	"github.com/go-passwd/validator"
)

func ValidatePasswordFormat(password string, passwordName ...string) error {
	fieldName := "Password"
	if len(passwordName) > 0 {
		fieldName = passwordName[0]
	}
	validate := validator.New(
		validator.MinLength(8, errors.New(fieldName+" length must be not lower than 8 chars")),
		validator.MaxLength(64, errors.New(fieldName+" length must be not greater than 64 chars")),
		validator.ContainsAtLeast("0123456789", 1, errors.New(fieldName+" must contains number")),
		validator.ContainsAtLeast("abcdefghijklmnopqrstuvwxyz", 1, errors.New(fieldName+" must contains lowercase")),
		validator.ContainsAtLeast("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 1, errors.New(fieldName+" must contains uppercase")),
		validator.ContainsAtLeast("!@#$%^&*.,?", 1, errors.New(fieldName+" must contains special character (!@#$%^&*.,?)")),
	)
	if err := validate.Validate(password); err != nil {
		return err
	}
	return nil
}

func ValidateUsernameFormat(username string) error {
	fieldName := "Username"
	validate := validator.New(
		validator.MinLength(2, errors.New(fieldName+" length must be not lower than 2 chars")),
		validator.MaxLength(64, errors.New(fieldName+" length must be not lower than 64 chars")),
		validator.ContainsOnly("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_", errors.New(fieldName+" must contains only letters, numbers and underscore (_)")),
	)
	if err := validate.Validate(username); err != nil {
		return err
	}
	return nil
}
