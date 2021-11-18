package validator

import (
	"fmt"
)

func validateStringRequired(val string, fieldname string) (string, bool) {
	if val == "" {
		return fmt.Sprintf("%s is required", fieldname), false
	}

	return "", true
}

func validateStringMinChars(val string, minChars int, fieldname string) (string, bool) {
	if len(val) < minChars {
		return fmt.Sprintf("%s should be at least %d characters", fieldname, minChars), false
	}

	return "", true
}

func validateStringMaxChars(val string, maxChars int, fieldname string) (string, bool) {
	if len(val) > maxChars {
		return fmt.Sprintf("%s should be at most %d characters", fieldname, maxChars), false
	}

	return "", true
}
