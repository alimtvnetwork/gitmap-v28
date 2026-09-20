package cmdagy

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// AgyPromptEntry stores a captured user prompt from an Antigravity transcript.
type AgyPromptEntry struct {
	StepIndex int       `json:"stepIndex"`
	CreatedAt time.Time `json:"createdAt"`
	Content   string    `json:"content"`
	Workspace string    `json:"workspace"`
	ConvID    string    `json:"convId"`
}

type rawTranscriptStep struct {
	StepIndex int    `json:"step_index"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	Content   string `json:"content"`
}

// GetBrainLogsDirPath returns the root path where Antigravity brain transcripts are stored.
func GetBrainLogsDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".gemini", "antigravity", "brain"), nil
}

// CollectPromptsForWorkspace reads prompts matching a workspace directory.
func CollectPromptsForWorkspace(workspace string) []AgyPromptEntry {
	all := CollectAllPrompts()
	hasEmptyWs := workspace == ""
	if hasEmptyWs {
		return all
	}

	return filterPromptsByWorkspace(all, workspace)
}

func filterPromptsByWorkspace(all []AgyPromptEntry, ws string) []AgyPromptEntry {
	normTarget := strings.ToLower(filepath.Clean(ws))
	var matched []AgyPromptEntry
	for _, p := range all {
		isMatch := strings.Contains(strings.ToLower(p.Workspace), normTarget)
		if isMatch {
			matched = append(matched, p)
		}
	}

	return matched
}

// CollectAllPrompts scans all conversation brain transcripts for user prompts.
func CollectAllPrompts() []AgyPromptEntry {
	brainDir, err := GetBrainLogsDirPath()
	hasErr := err != nil
	if hasErr {
		return nil
	}
	all := scanBrainDirPrompts(brainDir)
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	return all
}

func scanBrainDirPrompts(brainDir string) []AgyPromptEntry {
	entries, readErr := os.ReadDir(brainDir)
	hasErr := readErr != nil
	if hasErr {
		return nil
	}
	var all []AgyPromptEntry
	for _, e := range entries {
		if e.IsDir() {
			all = append(all, readConvPrompts(brainDir, e.Name())...)
		}
	}

	return all
}

func readConvPrompts(brainDir, convID string) []AgyPromptEntry {
	f, err := openConvTranscript(brainDir, convID)
	hasErr := err != nil
	if hasErr {
		return nil
	}
	defer f.Close()

	ws := resolveConvWorkspace(convID)

	return scanTranscriptLines(f, convID, ws)
}

func openConvTranscript(brainDir, convID string) (*os.File, error) {
	logDir := filepath.Join(brainDir, convID, ".system_generated", "logs")
	fullPath := filepath.Join(logDir, "transcript_full.jsonl")
	f, err := os.Open(fullPath)
	hasFull := err == nil
	if hasFull {
		return f, nil
	}
	shortPath := filepath.Join(logDir, "transcript.jsonl")

	return os.Open(shortPath)
}

func openConvTranscriptFile(convID string) (*os.File, error) {
	brainDir, err := GetBrainLogsDirPath()
	hasDirErr := err != nil
	if hasDirErr {
		return nil, err
	}

	return openConvTranscript(brainDir, convID)
}
