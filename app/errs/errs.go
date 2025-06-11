package errs

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func New(code ErrCode, err error) Error {
	return Error{
		Code:    code,
		Message: err.Error(),
	}
}

func Newf(code ErrCode, format string, a ...any) Error {
	return Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}

// Error implements the error interface.
func (e Error) Error() string {
	return e.Message
}

// IsError checks if the error is an Error.
func IsError(err error) bool {
	var e Error
	return errors.As(err, &e)
}

// GetError returns a copy of the Error.
func GetError(err error) Error {
	var e Error
	if !errors.As(err, &e) {
		return Error{}
	}
	return e
}
