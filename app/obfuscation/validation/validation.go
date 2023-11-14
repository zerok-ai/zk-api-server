package validation

import "regexp"

func IsValidRegex(expr string) bool {
	_, err := regexp.Compile(expr)
	return err == nil
}
