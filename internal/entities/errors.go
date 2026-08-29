package entities

import "errors"

var (
	ErrInvalidParams     = errors.New("invalid params")
	ErrRepositoryFailure = errors.New("repository failure")
	ErrMessagingFailure  = errors.New("messaging failure")
	ErrClientFailure     = errors.New("client failure")
	ErrInvalidMessage    = errors.New("invalid message")
)
