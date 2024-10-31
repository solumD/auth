package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

var (
	ErrInvalidEmail = fmt.Errorf("invalid email")

	ErrNameContainsSpaces = fmt.Errorf("name can't contain spaces")
	ErrNameTooShort       = fmt.Errorf("minimum name's length is 8 symbols")
	ErrInvalidName        = fmt.Errorf("invalid name")

	ErrPassContainsSpaces = fmt.Errorf("password can't contain spaces")
	ErrPassTooShort       = fmt.Errorf("minimum password's length is 8 symbols")
	ErrInvalidPass        = fmt.Errorf("invalid password")
)

// ValidateName валидирует имя пользователя
func ValidateName(name string) error {
	if len(strings.Fields(name)) > 1 {
		return ErrNameContainsSpaces
	}

	if len(name) < 8 {
		return ErrNameTooShort
	}

	re := regexp.MustCompile("^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$")
	if !re.MatchString(name) {
		return ErrInvalidName
	}

	return nil
}

// ValidateEmail валидирует почту
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}

	return nil
}

// ValidatePassword валидирует пароль
func ValidatePassword(password string) error {
	if len(strings.Fields(password)) > 1 {
		return ErrPassContainsSpaces
	}

	if len(password) < 8 {
		return ErrPassTooShort
	}

	re := regexp.MustCompile("^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$")
	if !re.MatchString(password) {
		return ErrInvalidPass
	}

	return nil
}
