package pion_clients

import "errors"

var (
	ErrAccessKeyNotFound = errors.New("accessKey not found")
	ErrInternalError     = errors.New("STS internal error")
)
