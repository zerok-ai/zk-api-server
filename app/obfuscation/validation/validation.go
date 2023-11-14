package validation

import (
	zkObfuscation "github.com/zerok-ai/zk-utils-go/obfuscation/model"
	"regexp"
)

func isValidRegex(expr string) bool {
	_, err := regexp.Compile(expr)
	return err == nil
}

func ValidateObfuscationRule(rule zkObfuscation.Rule) (bool, string) {
	if rule.Analyzer.Type != "regex" {
		return false, "Unsupported analyzer type. Only regex is supported."
	}
	validateRegex := isValidRegex(rule.Analyzer.Pattern)
	if !validateRegex {
		return false, "Invalid regex pattern."
	}
	if rule.Anonymizer.Operator != "replace" {
		return false, "Invalid anonymizer operator. Only replace is supported."
	}
	return true, ""
}
