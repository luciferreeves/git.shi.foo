package validate

import "regexp"

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func Email(email string) bool {
	return emailPattern.MatchString(email)
}
