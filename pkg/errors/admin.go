package errors

import "errors"

var (
	ErrTooManyAttempts = errors.New("too many attempts")
)
