package session

import (
	"context"
	"fmt"
)

// Errors
var (
	ErrSessionNotFound = fmt.Errorf("Session not found")
	ErrUnexpected      = fmt.Errorf("Unexpected session error")
)

// Session is collection of behavior of session.
type Session interface {
	Set(ctx context.Context, key string, value []byte) (err error)
	Get(ctx context.Context, key string) (value []byte, err error)
	Update(ctx context.Context, key string, value []byte) (err error)
	Delete(ctx context.Context, key string) (err error)
}
