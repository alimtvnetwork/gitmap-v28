package applogger

import (
	"context"
)

// ContextKey defines a custom type for context values to avoid collisions.
type ContextKey string

const (
	// RequestIDKey is the standard context key for request tracking.
	RequestIDKey ContextKey = "request_id"
	// TraceIDKey is the standard context key for distributed trace tracking.
	TraceIDKey ContextKey = "trace_id"
	// UserIDKey is the standard context key for the authenticated user ID.
	UserIDKey ContextKey = "user_id"
)

// ExtractContextFields retrieves standard observability fields from the context.
// It returns a map of available fields, skipping any that are not present.
func ExtractContextFields(ctx context.Context) map[string]any {
	fields := make(map[string]any)

	if reqID, ok := ctx.Value(RequestIDKey).(string); ok && reqID != "" {
		fields[string(RequestIDKey)] = reqID
	}

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		fields[string(TraceIDKey)] = traceID
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		fields[string(UserIDKey)] = userID
	}

	return fields
}

// FromContext creates a child Logger that inherits standard observability fields
// (like request_id, trace_id, and user_id) directly from the provided context.Context.
// If the context is nil, it simply returns the unmodified logger.
func FromContext(ctx context.Context, logger Logger) Logger {
	if ctx == nil || logger == nil {
		return logger
	}

	fields := ExtractContextFields(ctx)
	if len(fields) == 0 {
		return logger
	}

	childLogger := logger.WithFields(fields)
	return childLogger
}
