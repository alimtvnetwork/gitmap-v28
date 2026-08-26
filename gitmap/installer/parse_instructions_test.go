// Package installer — parse_instructions_test.go tests instruction parsing.
package installer

import (
	"testing"
)

func TestParseInstructions(t *testing.T) {
	jsonBlob := `[{"command": "echo hello", "description": "greeting"}]`
	res := ParseInstructions(jsonBlob)
	if len(res) != 1 || res[0].Command != "echo hello" {
		t.Fatalf("unexpected parsed json: %+v", res)
	}

	linesBlob := "echo 1\necho 2"
	resLines := ParseInstructions(linesBlob)
	if len(resLines) != 2 {
		t.Fatalf("unexpected parsed lines: %+v", resLines)
	}

	if empty := ParseInstructions(""); empty != nil {
		t.Fatalf("expected nil on empty string")
	}
}
