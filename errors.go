package memap

import (
	"errors"
	"fmt"
)

var (
	// ErrClosed is returned when attempting an operation on a closed [Client] connection.
	ErrClosed = errors.New("memap: client connection is closed")

	// ErrNamespaceAlreadyExists is returned when attempting to create a namespace that already exists.
	ErrNamespaceAlreadyExists = errors.New("memap: namespace already exists")

	// ErrNamespaceNotFound is returned when the requested namespace does not exist.
	ErrNamespaceNotFound = errors.New("memap: namespace not found")

	// ErrKeyNotFound is returned when the requested key does not exist or has expired.
	ErrKeyNotFound = errors.New("memap: key not found")

	// ErrKeyAlreadyExists is returned when attempting to initialize a data structure with a key that already exists.
	ErrKeyAlreadyExists = errors.New("memap: key already exists")

	// ErrBufferEmpty is returned when attempting to pop or peek from an empty ring buffer.
	ErrBufferEmpty = errors.New("memap: buffer is empty")

	// ErrIndexOutOfBounds is returned when attempting to access a ring buffer element at an invalid index.
	ErrIndexOutOfBounds = errors.New("memap: index out of bounds")

	// ErrLimitExceeded is returned when a counter increment/decrement exceeds its limit or drops below zero.
	ErrLimitExceeded = errors.New("memap: limit exceeded")

	// ErrFieldNotFound is returned when the requested field does not exist in the hash map.
	ErrFieldNotFound = errors.New("memap: field not found")

	// errNilRequest is returned internally when a nil request is passed to do.
	errNilRequest = errors.New("memap: request cannot be nil")
)

// serverErrorMap maps raw error messages from the Memap server to client sentinel errors.
var serverErrorMap = map[string]error{
	"namespace already exists": ErrNamespaceAlreadyExists,
	"namespace not found":      ErrNamespaceNotFound,
	"key not found":            ErrKeyNotFound,
	"key already exists":       ErrKeyAlreadyExists,
	"buffer is empty":          ErrBufferEmpty,
	"index out of bounds":      ErrIndexOutOfBounds,
	"limit exceeded":           ErrLimitExceeded,
	"field not found":          ErrFieldNotFound,
}

// ServerError represents a command-level failure returned by the Memap server.
// If the server message matches a known sentinel error (such as [ErrKeyNotFound]),
// [ServerError.Unwrap] returns that sentinel error, enabling checks with [errors.Is].
type ServerError struct {
	Message string
	Err     error
}

// Error formats the server failure message as a human-readable string.
func (e *ServerError) Error() string {
	return fmt.Sprintf("memap server error: %s", e.Message)
}

// Unwrap returns the underlying client sentinel error, if matched.
func (e *ServerError) Unwrap() error {
	return e.Err
}

func newServerError(msg string) error {
	if msg == "" {
		return errors.New("memap: unknown server error")
	}
	sentinel := serverErrorMap[msg]
	return &ServerError{
		Message: msg,
		Err:     sentinel,
	}
}
