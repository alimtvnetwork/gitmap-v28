package cmdagy

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"time"
)

func scanTranscriptLines(f *os.File, convID, ws string) []AgyPromptEntry {
	var list []AgyPromptEntry
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		entry, hasPrompt := parseTranscriptLine(scanner.Bytes(), convID, ws)
		if hasPrompt {
			list = append(list, entry)
		}
	}

	return list
}

func parseTranscriptLine(line []byte, convID, ws string) (AgyPromptEntry, bool) {
	var step rawTranscriptStep
	if err := json.Unmarshal(line, &step); err != nil || step.Type != "USER_INPUT" {
		return AgyPromptEntry{}, false
	}

	return buildPromptEntry(step, convID, ws)
}

func buildPromptEntry(step rawTranscriptStep, convID, ws string) (AgyPromptEntry, bool) {
	cleanText := cleanPromptText(step.Content)
	if cleanText == "" {
		return AgyPromptEntry{}, false
	}
	ts, _ := time.Parse(time.RFC3339, step.CreatedAt)

	return AgyPromptEntry{
		StepIndex: step.StepIndex,
		CreatedAt: ts,
		Content:   cleanText,
		Workspace: ws,
		ConvID:    convID,
	}, true
}

func cleanPromptText(raw string) string {
	text := strings.TrimPrefix(raw, "<USER_REQUEST>")
	text = strings.TrimSuffix(text, "</USER_REQUEST>")

	return strings.TrimSpace(text)
}
