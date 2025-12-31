package web

import "errors"

// Shutdown is a type used to help with the graceful termination of
// the service.
type Shutdown struct {
	Message string
}

// New returns an errors that causes the framework to signal
// a graceful Shutdown.
func New(message string) *Shutdown {
	return &Shutdown{
		Message: message,
	}
}

// Error implements the error interface.
func (s *Shutdown) Error() string {
	return s.Message
}

// IsShutdown is a type assertion for the Shutdown type.
// In other words, it checks whether the Shutdown type is contained
// in the specific error value.
func IsShutdown(err error) bool {
	var se *Shutdown
	return errors.As(err, &se)
}
