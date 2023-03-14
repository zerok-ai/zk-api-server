package errors

import "errors"

var (
	ErrAuthenticationFailed = errors.New("rpc error: code = Internal desc = Auth middleware failed: failed to fetch token - unauthenticated")
)
