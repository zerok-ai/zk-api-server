package validation

import (
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkObfuscation "github.com/zerok-ai/zk-utils-go/obfuscation/model"
	"regexp"
	"strconv"
)

func ValidateLimit(limit string) (bool, string) {
	if zkCommon.IsEmpty(limit) {
		return false, "Limit cannot be empty"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return false, "Limit must be a valid number"
	}

	if limitInt <= 0 {
		return false, "Limit must be greater than zero"
	}

	return true, ""
}

func ValidateOffset(offset string) (bool, string) {
	if zkCommon.IsEmpty(offset) {
		return false, "Offset cannot be empty"
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return false, "Offset must be a valid number"
	}

	if offsetInt < 0 {
		return false, "Offset cannot be negative"
	}

	return true, ""
}

func ValidateOrgId(orgId string) (bool, string) {
	if zkCommon.IsEmpty(orgId) {
		return false, "OrgId cannot be empty"
	}
	return true, ""
}

func ValidateId(Id string) (bool, string) {
	if zkCommon.IsEmpty(Id) {
		return false, "Obfuscation id cannot be empty"
	}
	return true, ""
}

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
