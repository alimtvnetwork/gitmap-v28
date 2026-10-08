package cmd

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/suggestion"
)

func TestResolveCommandSuggestions_Typo(t *testing.T) {
	group := ResolveCommandSuggestions("scann")
	if !group.HasSuggestions() {
		t.Fatalf("expected suggestions for typo 'scann'")
	}
}

func TestBuildUnknownCommandAppError_WithSuggestions(t *testing.T) {
	group := ResolveCommandSuggestions("scann")
	appErr := buildUnknownCommandAppError("scann", group)
	if appErr == nil {
		t.Fatalf("expected non-nil AppError")
	}
	if !appErr.HasSuggestions() {
		t.Errorf("expected suggestions attached to AppError")
	}
}

func TestRenderErrorSuggestions_Nil(t *testing.T) {
	var buf bytes.Buffer
	hasRendered := RenderErrorSuggestions(&buf, nil)
	if hasRendered {
		t.Errorf("expected false for nil AppError")
	}
}

func TestRenderErrorSuggestions_Valid(t *testing.T) {
	var buf bytes.Buffer
	appErr := apperror.NewNotFoundError("unrecognized")
	appErr.WithSuggestions(suggestion.Suggestion{
		Command:     "gitmap scan",
		Description: "Fast scanner",
		Confidence:  0.95,
	})
	hasRendered := RenderErrorSuggestions(&buf, appErr)
	if !hasRendered {
		t.Errorf("expected true for AppError with suggestions")
	}
}
