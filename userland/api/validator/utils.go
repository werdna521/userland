package validator

func hasUppercase(val string) bool {
	for _, v := range val {
		if v >= 'A' && v <= 'Z' {
			return true
		}
	}
	return false
}

func hasLowercase(val string) bool {
	for _, v := range val {
		if v >= 'a' && v <= 'z' {
			return true
		}
	}
	return false
}

func hasNumber(val string) bool {
	for _, v := range val {
		if v >= '0' && v <= '9' {
			return true
		}
	}
	return false
}
