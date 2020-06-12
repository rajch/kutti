package clustermanager

import (
	"regexp"
)

// IsValidName checks for the validity of a name.
// Valid names are up to 10 characters long, must start with a lowercase letter, and may contain lowercase letters and digits only.
func IsValidName(name string) bool {
	matched, _ := regexp.MatchString("^[a-z]([a-z0-9]{1,9})$", name)
	return matched
}
