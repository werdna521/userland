package validator

import "strings"

const (
	fullnameMaxChars  = 128
	fullnameFieldname = "fullname"
)

func ValidateFullname(fullname string) (string, bool) {
	errMsg, ok := validateStringRequired(fullname, fullnameFieldname)
	if !ok {
		return errMsg, false
	}

	errMsg, ok = validateStringMaxChars(fullname, fullnameMaxChars, fullnameFieldname)
	if !ok {
		return errMsg, false
	}

	return "", true
}

const (
	emailMaxChars  = 128
	emailFieldname = "email"
)

func ValidateEmail(email string) (string, bool) {
	errMsg, ok := validateStringRequired(email, emailFieldname)
	if !ok {
		return errMsg, false
	}

	errMsg, ok = validateStringMaxChars(email, emailMaxChars, emailFieldname)
	if !ok {
		return errMsg, false
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return "invalid email", false
	}

	if strings.LastIndex(email, ".") < strings.Index(email, "@") {
		return "invalid email", false
	}

	return "", true
}

const (
	passwordMinChars  = 8
	passwordMaxChars  = 128
	passwordFieldname = "password"
)

func ValidatePassword(password string) (string, bool) {
	errMsg, ok := validateStringRequired(password, passwordFieldname)
	if !ok {
		return errMsg, false
	}

	errMsg, ok = validateStringMinChars(password, passwordMinChars, passwordFieldname)
	if !ok {
		return errMsg, false
	}

	errMsg, ok = validateStringMaxChars(password, passwordMaxChars, passwordFieldname)
	if !ok {
		return errMsg, false
	}

	// doing this means we'll have a complexity of O(3n). there are other ways to
	// do this that would only cost O(n), but I decided to go with this since it's
	// more readable and easier to follow.
	if !hasLowercase(password) || !hasUppercase(password) || !hasNumber(password) {
		return "password should have at least 1 uppercase character, 1 lowercase character and 1 number", false
	}

	return "", true
}
