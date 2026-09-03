package memap

import (
	"errors"
	"fmt"
)

var (
	// ErrClosed is returned when attempting an operation on a closed Client connection.
	ErrClosed = errors.New("memap: client connection is closed")

	// errNilRequest is returned internally when a nil request is passed to do.
	errNilRequest = errors.New("memap: request cannot be nil")
)

// ServerError represents a command-level failure returned by the Memap server.
type ServerError struct {
	Message string
}

// Error formats the server failure message as a human-readable string.
func (e *ServerError) Error() string {
	return fmt.Sprintf("memap server error: %s", e.Message)
}

func newServerError(msg string) error {
	if msg == "" {
		return errors.New("memap: unknown server error")
	}
	return &ServerError{Message: msg}
}
