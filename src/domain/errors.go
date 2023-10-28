package domain

import "errors"

var (
	ErrInternalServer  = errors.New("internal server error")
	ErrDatabase        = errors.New("database error")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
)
