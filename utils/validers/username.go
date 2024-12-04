package validers

import "regexp"

func IsValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-z0-9_]+$`)
	return re.MatchString(username)
}
