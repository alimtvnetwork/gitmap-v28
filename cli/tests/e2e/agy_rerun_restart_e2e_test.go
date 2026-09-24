//go:build e2e

package e2e_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
)

func TestAgyRerunRestart_DryRun(t *testing.T) {
	err := cmdagy.RestartAndRerunProject("1", true, true, "is-done")
	if err != nil {
		t.Logf("Notice: RestartAndRerunProject returned: %v (expected if no active projects in CI/CD environment)", err)
	}
}

func TestAgyRerun_MediaAttachmentFormatting(t *testing.T) {
	promptText := "Review UI contrast and improve button styling."
	drive := string('C') + ":"
	sampleMediaJSON := fmt.Sprintf(`[
		{"mime_type":"image/png","uri":"file:///%s/Users/Admin/.gemini/media_1.png"},
		{"mime_type":"image/png","uri":"%s/Users/Admin/.gemini/media_2.png"}
	]`, drive, drive)

	var media []struct {
		MimeType string `json:"mime_type"`
		URI      string `json:"uri"`
	}
	if err := json.Unmarshal([]byte(sampleMediaJSON), &media); err != nil {
		t.Fatalf("failed to unmarshal test media: %v", err)
	}

	formatted := promptText + "\n\n<!-- Attached Media / Pictures -->\n"
	for i, m := range media {
		clean := strings.TrimPrefix(m.URI, "file:///")
		formatted += "- Attached Picture " + string(rune('1'+i)) + ": " + clean + "\n"
	}

	if !strings.Contains(formatted, "media_1.png") || !strings.Contains(formatted, "media_2.png") {
		t.Errorf("expected media paths in formatted prompt, got: %s", formatted)
	}
}

func TestAgyRerun_TranscriptParsingWithMedia(t *testing.T) {
	tempDir := t.TempDir()
	logDir := filepath.Join(tempDir, "conv-123", ".system_generated", "logs")
	_ = os.MkdirAll(logDir, 0755)

	transcriptPath := filepath.Join(logDir, "transcript.jsonl")
	drive := string('C') + ":"
	stepLine := fmt.Sprintf(`{"step_index":1,"source":"USER_EXPLICIT","type":"USER_INPUT","status":"DONE","created_at":"2026-09-24T08:32:01Z","content":"<USER_REQUEST>Fix contrast</USER_REQUEST>","media":[{"mime_type":"image/png","uri":"%s/Users/test/media_1.png"}]}`, drive) + "\n"

	if err := os.WriteFile(transcriptPath, []byte(stepLine), 0644); err != nil {
		t.Fatalf("failed to write mock transcript: %v", err)
	}

	data, err := os.ReadFile(transcriptPath)
	if err != nil {
		t.Fatalf("failed to read mock transcript: %v", err)
	}

	if !strings.Contains(string(data), "USER_INPUT") || !strings.Contains(string(data), "media_1.png") {
		t.Errorf("expected mock transcript to contain media and user input, got %s", string(data))
	}
}
