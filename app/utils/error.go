package utils

import "errors"

var (
	ErrAuthenticationFailed = errors.New("rpc error: code = Internal desc = Auth middleware failed: failed to fetch token - unauthenticated")
	ErrClusterParsingFailed = errors.New("failed to parse cluster info")
	ErrClusterIdEmpty       = errors.New("clusterId cannot be empty")
	ErrPxlStartTimeEmpty    = errors.New("start Time st cannot be empty")
	ErrZkApiKeyEmpty        = errors.New("ZK_API_KEY header cannot be empty")
	ErrNamespaceEmpty       = errors.New("namespace ns cannot be empty")
	ErrServiceNameEmpty     = errors.New("service name cannot be empty")
	ErrPodNameEmpty         = errors.New("pod name cannot be empty")
	ErrInternalServerError  = errors.New("something went wrong, please try again later")
)
