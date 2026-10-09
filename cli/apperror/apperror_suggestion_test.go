package apperror

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/diag"
)

func TestAppError_WithSuggestions(t *testing.T) {
	err := NewNotFoundError("missing")
	if err.HasSuggestions() {
		t.Fatalf("expected HasSuggestions to be false initially")
	}
	err.WithSuggestions(diag.Suggestion{
		Command:    "gitmap scan",
		Confidence: 0.90,
	})
	if !err.HasSuggestions() {
		t.Fatalf("expected HasSuggestions to be true")
	}
	items := err.GetSuggestions()
	if len(items) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(items))
	}
}

func TestAppError_WithSuggestionGroup(t *testing.T) {
	err := NewNotFoundError("missing")
	group := diag.SuggestionGroup{
		Suggestions: []diag.Suggestion{
			{Command: "gitmap status", Confidence: 0.95},
			{Command: "gitmap scan", Confidence: 0.85},
		},
	}
	err.WithSuggestionGroup(group)
	if len(err.GetSuggestions()) != 2 {
		t.Errorf("expected 2 suggestions, got %d", len(err.GetSuggestions()))
	}
}

func TestAppError_NilSafety(t *testing.T) {
	var err *AppError
	if err.HasSuggestions() {
		t.Errorf("expected false for nil AppError")
	}
	if err.GetSuggestions() != nil {
		t.Errorf("expected nil suggestions for nil AppError")
	}
	if err.WithSuggestions() != nil {
		t.Errorf("expected nil for WithSuggestions on nil AppError")
	}
}
