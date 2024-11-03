package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

var (
	errInvalidEmail = fmt.Errorf("invalid email")

	errNameContainsSpaces = fmt.Errorf("name can't contain spaces")
	errNameTooShort       = fmt.Errorf("minimum name's length is 8 symbols")
	errInvalidName        = fmt.Errorf("invalid name")

	errPassContainsSpaces = fmt.Errorf("password can't contain spaces")
	errPassTooShort       = fmt.Errorf("minimum password's length is 8 symbols")
	errInvalidPass        = fmt.Errorf("invalid password")
)

// ValidateName валидирует имя пользователя
func ValidateName(name string) error {
	if len(strings.Fields(name)) > 1 {
		return errNameContainsSpaces
	}

	if len(name) < 8 {
		return errNameTooShort
	}

	re := regexp.MustCompile("^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$")
	if !re.MatchString(name) {
		return errInvalidName
	}

	return nil
}

// ValidateEmail валидирует почту
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errInvalidEmail
	}

	return nil
}

// ValidatePassword валидирует пароль
func ValidatePassword(password string) error {
	if len(strings.Fields(password)) > 1 {
		return errPassContainsSpaces
	}

	if len(password) < 8 {
		return errPassTooShort
	}

	re := regexp.MustCompile("^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$")
	if !re.MatchString(password) {
		return errInvalidPass
	}

	return nil
}
