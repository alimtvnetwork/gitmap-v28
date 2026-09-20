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
	err := json.Unmarshal(line, &step)
	isUserInput := err == nil && step.Type == "USER_INPUT"
	if isUserInput {
		return buildPromptEntry(step, convID, ws)
	}

	return AgyPromptEntry{}, false
}

func buildPromptEntry(step rawTranscriptStep, convID, ws string) (AgyPromptEntry, bool) {
	cleanText := cleanPromptText(step.Content)
	hasText := cleanText != ""
	if hasText {
		return makePromptEntry(step, convID, ws, cleanText), true
	}

	return AgyPromptEntry{}, false
}

func makePromptEntry(step rawTranscriptStep, convID, ws, content string) AgyPromptEntry {
	ts, _ := time.Parse(time.RFC3339, step.CreatedAt)

	return AgyPromptEntry{
		StepIndex: step.StepIndex,
		CreatedAt: ts,
		Content:   content,
		Workspace: ws,
		ConvID:    convID,
	}
}

func cleanPromptText(raw string) string {
	text := stripTagBlock(raw, "<ADDITIONAL_METADATA>", "</ADDITIONAL_METADATA>")
	text = stripTagBlock(text, "<USER_SETTINGS_CHANGE>", "</USER_SETTINGS_CHANGE>")
	text = strings.TrimPrefix(strings.TrimSpace(text), "<USER_REQUEST>")
	text = strings.TrimSuffix(strings.TrimSpace(text), "</USER_REQUEST>")

	return strings.TrimSpace(text)
}

func stripTagBlock(text, openTag, closeTag string) string {
	start := strings.Index(text, openTag)
	hasTag := start >= 0
	if hasTag == false {
		return text
	}

	return removeTagFromStart(text, start, openTag, closeTag)
}

func removeTagFromStart(text string, start int, openTag, closeTag string) string {
	end := strings.Index(text[start:], closeTag)
	hasEnd := end >= 0
	if hasEnd {
		endPos := start + end + len(closeTag)

		return stripTagBlock(text[:start]+text[endPos:], openTag, closeTag)
	}

	return text[:start]
}
