package memap

import (
	"errors"
	"testing"
)

func TestErrorMapping(t *testing.T) {
	tests := []struct {
		serverMsg string
		targetErr error
	}{
		{"namespace already exists", ErrNamespaceAlreadyExists},
		{"namespace not found", ErrNamespaceNotFound},
		{"key not found", ErrKeyNotFound},
		{"key already exists", ErrKeyAlreadyExists},
		{"buffer is empty", ErrBufferEmpty},
		{"index out of bounds", ErrIndexOutOfBounds},
		{"limit exceeded", ErrLimitExceeded},
		{"field not found", ErrFieldNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.serverMsg, func(t *testing.T) {
			err := newServerError(tt.serverMsg)
			if !errors.Is(err, tt.targetErr) {
				t.Fatalf("expected errors.Is(err, %v) to be true, got false", tt.targetErr)
			}

			var sErr *ServerError
			if !errors.As(err, &sErr) {
				t.Fatalf("expected errors.As to succeed for *ServerError")
			}
			if sErr.Message != tt.serverMsg {
				t.Fatalf("expected Message == %q, got %q", tt.serverMsg, sErr.Message)
			}
		})
	}
}

func TestUnknownServerError(t *testing.T) {
	err := newServerError("something unexpected happened")
	if errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("expected errors.Is with ErrKeyNotFound to be false")
	}

	var sErr *ServerError
	if !errors.As(err, &sErr) {
		t.Fatalf("expected errors.As to succeed")
	}
	if sErr.Message != "something unexpected happened" {
		t.Fatalf("unexpected message: %s", sErr.Message)
	}
}
