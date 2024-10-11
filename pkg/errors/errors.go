package errors

import "errors"

var (
	ErrGeneric         error = errors.New("something went wrong")
	ErrNotFound        error = errors.New("not found")
	ErrInvalidPassword error = errors.New("password is invalid")
)
