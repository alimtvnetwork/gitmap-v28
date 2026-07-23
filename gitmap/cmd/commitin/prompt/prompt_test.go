package prompt

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestAskStringNoPromptReturnsError(t *testing.T) {
	var out bytes.Buffer
	a := &Asker{NoPrompt: true, In: strings.NewReader(""), Out: &out}
	_, err := a.AskString("ConflictMode", "Conflict mode?", "")
	if !errors.Is(err, ErrNoPrompt) {
		t.Fatalf("want ErrNoPrompt, got %v", err)
	}
	if !strings.Contains(out.String(), "ConflictMode") {
		t.Fatalf("error message lacks field name: %q", out.String())
	}
}

func TestAskStringFallbackOnEmpty(t *testing.T) {
	a := &Asker{In: strings.NewReader("\n"), Out: &bytes.Buffer{}}
	got, err := a.AskString("X", "X?", "default-val")
	if err != nil || got != "default-val" {
		t.Fatalf("want fallback, got %q err=%v", got, err)
	}
}

func TestAskStringReturnsTypedAnswer(t *testing.T) {
	a := &Asker{In: strings.NewReader("hello\n"), Out: &bytes.Buffer{}}
	got, _ := a.AskString("X", "X?", "")
	if got != "hello" {
		t.Fatalf("want hello, got %q", got)
	}
}

func TestAskEnumRetriesUntilValid(t *testing.T) {
	a := &Asker{In: strings.NewReader("Bogus\nForceMerge\n"), Out: &bytes.Buffer{}}
	got, err := a.AskEnum("ConflictMode", "Mode", []string{"ForceMerge", "Prompt"}, "")
	if err != nil || got != "ForceMerge" {
		t.Fatalf("want ForceMerge, got %q err=%v", got, err)
	}
}
