package cmdagy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDecomposePromptContent(t *testing.T) {
	input := `# Section 1: Architecture
Setup the initial configuration.

## Section 2: Implementation
Write the core Go code.

## Section 3: Verification
Run all tests and verify quality.`

	sections := decomposePromptContent(input)
	if len(sections) != 3 {
		t.Fatalf("expected 3 sections, got %d", len(sections))
	}

	if sections[0].Index != 3 || sections[0].Title != "Section 1: Architecture" {
		t.Errorf("unexpected section 0: %+v", sections[0])
	}
	if sections[1].Index != 4 || sections[1].Title != "Section 2: Implementation" {
		t.Errorf("unexpected section 1: %+v", sections[1])
	}
}

func TestRunAgyEnhancePrompts(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-enhance-prompts")
	_ = os.RemoveAll(tmpDir)
	defer os.RemoveAll(tmpDir)

	input := `# Refactor Boolean Flags
Ensure positive booleans across packages.

# Add Unit Tests
Add 100% test coverage.`

	err := runAgyEnhancePrompts(input, tmpDir)
	if err != nil {
		t.Fatalf("unexpected error running enhance-prompts: %v", err)
	}

	indexFile := filepath.Join(tmpDir, "01-index.md")
	if _, statErr := os.Stat(indexFile); statErr != nil {
		t.Errorf("expected 01-index.md to exist")
	}

	verbatimFile := filepath.Join(tmpDir, "02-verbatim-instructions.md")
	if _, statErr := os.Stat(verbatimFile); statErr != nil {
		t.Errorf("expected 02-verbatim-instructions.md to exist")
	}

	sec1File := filepath.Join(tmpDir, "03-refactor-boolean-flags.md")
	if _, statErr := os.Stat(sec1File); statErr != nil {
		t.Errorf("expected 03-refactor-boolean-flags.md to exist")
	}
}

func TestCompactWords(t *testing.T) {
	text := "one two three four five six seven eight nine ten"
	compacted := CompactWords(text, 5)
	if compacted != "one two three four five..." {
		t.Errorf("unexpected compact words: %q", compacted)
	}

	noTrunc := CompactWords(text, 20)
	if noTrunc != text {
		t.Errorf("expected original text, got: %q", noTrunc)
	}
}

func TestDecoratePromptBody(t *testing.T) {
	prefix := "Verify quality first."
	body := "Fix the bug in parser."
	suffix := "Sponsored by Rise Up Asia LLC."

	decorated := decoratePromptBody(body, prefix, suffix)
	expected := "Verify quality first.\n\nFix the bug in parser.\n\nSponsored by Rise Up Asia LLC."
	if decorated != expected {
		t.Errorf("unexpected decorated prompt:\nGot: %q\nWant: %q", decorated, expected)
	}
}
