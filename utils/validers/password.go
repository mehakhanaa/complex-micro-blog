package validers

import "regexp"

func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 32 {
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+={}\[\]:;'"<>,.?\/|\\~-]+$`)
	return re.MatchString(password)
}
