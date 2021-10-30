package validator

const (
	fullnameMinChars  = 3
	fullnameMaxChars  = 128
	fullnameFieldname = "fullname"
)

func ValidateFullname(fullname string) (string, bool) {
	errMsg, ok := validateStringRequired(fullname, fullnameFieldname)
	if !ok {
		return errMsg, false
	}

	errMsg, ok = validateStringMinChars(fullname, fullnameMinChars, fullnameFieldname)
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
	locationMaxChars  = 128
	locationFieldname = "location"
)

func ValidateLocation(location string) (string, bool) {
	errMsg, ok := validateStringMaxChars(location, locationMaxChars, locationFieldname)
	if !ok {
		return errMsg, false
	}

	return "", true
}

const (
	bioMaxChars  = 255
	bioFieldname = "bio"
)

func ValidateBio(bio string) (string, bool) {
	errMsg, ok := validateStringMaxChars(bio, bioMaxChars, bioFieldname)
	if !ok {
		return errMsg, false
	}

	return "", true
}

const (
	webMaxChars  = 128
	webFieldname = "web"
)

func ValidateWeb(web string) (string, bool) {
	errMsg, ok := validateStringMaxChars(web, webMaxChars, webFieldname)
	if !ok {
		return errMsg, false
	}

	return "", true
}
