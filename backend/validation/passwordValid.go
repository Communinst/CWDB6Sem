package validation

import (
	"errors"
	"regexp"
)

func IsPasswordValid(password string) error {
	length := len(password)
	if length < 6 || length > 31 {
		return errors.New("Invalid password length: at least 6, at peak 31.")
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString
	hasLower := regexp.MustCompile(`[a-z]`).MatchString
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString
	hasSpecial := regexp.MustCompile(`[#!@$%^&*-]`).MatchString

	if hasUpper(password) && hasLower(password) && hasDigit(password) && hasSpecial(password) {
		return nil
	}
	return errors.New("Invalid password value.")

}
