package storage

import "errors"

var (
	ErrInvalidPayload = errors.New("invalid payload")
	ErrInvalidSubject = errors.New("invalid subject")
)
