package cmdagy

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// DefaultQueueCheckPrefix is prepended to queued prompts when re-injecting after an active rerun.
const DefaultQueueCheckPrefix = "is it completed properly and released?? Please verify if the following task is already completed; if not, execute and complete it now:"

var imagePathRegex = regexp.MustCompile(`(?i)(?:file:///|[a-zA-Z]:[\\/]|[\w\-\./\\]+\/)[a-zA-Z0-9_\-\./\\]+\.(?:png|jpg|jpeg|gif|webp)`)

// ExtractAllImagePathsFromPrompt combines structured transcript media URIs and inline file paths.
func ExtractAllImagePathsFromPrompt(content string, media []rawTranscriptMedia) []string {
	seen := make(map[string]bool)
	var paths []string

	for _, m := range media {
		clean := normalizeImagePath(m.URI)
		if clean != "" && !seen[clean] {
			seen[clean] = true
			paths = append(paths, clean)
		}
	}

	matches := imagePathRegex.FindAllString(content, -1)
	for _, match := range matches {
		clean := normalizeImagePath(match)
		if clean != "" && !seen[clean] {
			seen[clean] = true
			paths = append(paths, clean)
		}
	}

	return paths
}

func normalizeImagePath(raw string) string {
	clean := strings.TrimSpace(raw)
	clean = strings.TrimPrefix(clean, "file:///")
	clean = strings.TrimPrefix(clean, "file://")
	if clean == "" {
		return ""
	}
	return filepath.FromSlash(clean)
}

// RequeuePendingWithCheckPrefix prepends the completion verification prefix to all queued prompts.
func RequeuePendingWithCheckPrefix(customPrefix string) ([]AgyPromptQueueEntry, error) {
	q, err := LoadPromptQueue()
	if err != nil || len(q.Queued) == 0 {
		return nil, err
	}

	prefix := strings.TrimSpace(customPrefix)
	if prefix == "" {
		prefix = DefaultQueueCheckPrefix
	}

	for i := range q.Queued {
		q.Queued[i].Prompt = ApplyCompletionCheckPrefix(q.Queued[i].Prompt, prefix)
		q.Queued[i].Status = "requeued_with_check_prefix"
	}

	if saveErr := SavePromptQueue(q); saveErr != nil {
		return q.Queued, saveErr
	}
	return q.Queued, nil
}

// ApplyCompletionCheckPrefix prepends the completion check prefix if not already present.
func ApplyCompletionCheckPrefix(promptText, prefix string) string {
	trimmed := strings.TrimSpace(promptText)
	if strings.HasPrefix( strings.ToLower(trimmed), "is it completed") {
		return trimmed
	}
	return fmt.Sprintf("%s\n\n%s", prefix, trimmed)
}

// RestartAndRerunAllActiveProjects reruns the last prompt across up to maxCount active running projects.
func RestartAndRerunAllActiveProjects(maxCount int, isRestart, isDryRun bool, templateID string) error {
	projects, err := loadSortedProjects()
	if err != nil {
		return err
	}

	limit := len(projects)
	if maxCount > 0 && maxCount < limit {
		limit = maxCount
	}

	fmt.Printf("\n%s[AGY RERUN ALL]%s Rerunning active prompts across %d project(s)...\n",
		constants.ColorCyan, constants.ColorReset, limit)

	for i := 0; i < limit; i++ {
		targetSeq := fmt.Sprintf("%d", i+1)
		if runErr := RestartAndRerunProject(targetSeq, isRestart, isDryRun, templateID); runErr != nil {
			return runErr
		}
	}
	return nil
}
