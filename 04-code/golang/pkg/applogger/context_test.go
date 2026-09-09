package applogger

import (
	"context"
	"testing"
)

func TestFromContext(t *testing.T) {
	logger := Default()

	ctx := context.WithValue(context.Background(), RequestIDKey, "req-1234")
	ctx = context.WithValue(ctx, UserIDKey, "usr-999")

	childLogger := FromContext(ctx, logger)

	// In testing we can't easily introspect the unexported fields of appLogger,
	// but we can ensure the interface contract is upheld and it doesn't panic.
	if childLogger == nil {
		t.Fatalf("Expected non-nil child logger")
	}

	// Just invoke it to ensure no panics
	childLogger.Info("Context-aware log message")
}
