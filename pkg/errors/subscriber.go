package errors

import "errors"

var (
	ErrAlreadySubscribed = errors.New("already subscribed")
	ErrSubscribing       = errors.New("error subscribing")
)
