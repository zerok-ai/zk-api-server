package validation

import (
	zkObfuscation "github.com/zerok-ai/zk-utils-go/obfuscation/model"
	"regexp"
)

// isValidRegex validates the regex pattern
// Returns true if the regex pattern is valid, else returns false
func isValidRegex(expr string) bool {
	_, err := regexp.Compile(expr)
	return err == nil
}

// ValidateObfuscationRule validates the obfuscation rule
// Returns true if the rule is valid, else returns false with the error message
func ValidateObfuscationRule(rule zkObfuscation.Rule) (bool, string) {

	// Only regex analyzer is supported
	if rule.Analyzer.Type != "regex" {
		return false, "Unsupported analyzer type. Only regex is supported."
	}

	//Checking if the regex pattern is valid
	validateRegex := isValidRegex(rule.Analyzer.Pattern)
	if !validateRegex {
		return false, "Invalid regex pattern."
	}

	//Only replace operator is supported for anonymizer.
	if rule.Anonymizer.Operator != "replace" {
		return false, "Invalid anonymizer operator. Only replace is supported."
	}

	return true, ""
}
