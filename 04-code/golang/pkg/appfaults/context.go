package appfaults

import (
	"context"

	"coding-guidelines/common/pkg/appfault"
)

type contextKeyType struct{}

var contextKey = contextKeyType{}

// WithFaults returns a derived context carrying a dedicated Collection.
func WithFaults(parent context.Context) (context.Context, *Collection) {
	coll := New()
	ctx := context.WithValue(parent, contextKey, coll)

	return ctx, coll
}

// FromContext retrieves the Collection attached to ctx if present.
func FromContext(ctx context.Context) (*Collection, bool) {
	if ctx == nil {
		return nil, false
	}

	coll, ok := ctx.Value(contextKey).(*Collection)

	return coll, ok
}

// RecordContextError appends err to the context collection if available.
func RecordContextError(ctx context.Context, err *appfault.AppError) bool {
	if coll, ok := FromContext(ctx); ok {
		coll.Add(err)

		return true
	}

	return false
}
