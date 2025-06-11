package web

import "errors"

// shutdown is a type used to help with the graceful termination of
// the service.
type shutdown struct {
	Message string
}

// New returns an errors that causes the framework to signal
// a graceful shutdown.
func New(message string) *shutdown {
	return &shutdown{
		Message: message,
	}
}

// Error implements the error interface.
func (s *shutdown) Error() string {
	return s.Message
}

// IsShutdown is a type assertion for the shutdown type.
// In other words, it checks whether the shutdown type is contained
// in the specific error value.
func IsShutdown(err error) bool {
	var se *shutdown
	return errors.As(err, &se)
}
