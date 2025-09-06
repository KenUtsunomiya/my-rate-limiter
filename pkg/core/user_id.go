package core

import "regexp"

type UserID string

var appIDPattern = regexp.MustCompile(`^[0-9]{1,19}$`)

// String ...
func (id UserID) String() string {
	return string(id)
}

// IsValid ...
func (id UserID) IsValid() bool {
	return appIDPattern.MatchString(string(id))
}
