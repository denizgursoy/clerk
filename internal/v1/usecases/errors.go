package usecases

import "errors"

var (
	ErrGroupNotFound  = errors.New("group does not exit")
	ErrMemberNotFound = errors.New("member does not exit")
)
