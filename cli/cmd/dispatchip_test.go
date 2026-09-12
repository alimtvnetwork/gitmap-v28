package cmd

import (
	"context"
	"testing"
)

func TestDispatchIP(t *testing.T) {
	ctx := context.Background()

	// Empty args
	if err := dispatchIP(ctx, []string{}, nil); err != nil {
		t.Errorf("expected no error on empty args, got %v", err)
	}

	// Unrelated args
	if err := dispatchIP(ctx, []string{"other"}, nil); err != nil {
		t.Errorf("expected no error on other args, got %v", err)
	}
}
