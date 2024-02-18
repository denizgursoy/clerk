package client

import "errors"

var (
	ErrEmptyFunction = errors.New("notify function cannot be empty")
	ErrEmptyGroup    = errors.New("group cannot be empty")
)
