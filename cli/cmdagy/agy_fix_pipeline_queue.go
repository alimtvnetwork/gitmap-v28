package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	queuedAgyPromptRelativePath = ".ai-memory/temp/queued-agy-followup-prompt.txt"
	agyPromptQueueRelativePath  = ".ai-memory/temp/agy-prompt-queue.json"
)

const verificationInstructions = `Please perform the following verification steps:
1. Confirm all root causes identified in the 4-part RCA have been properly remediated.
2. Verify that modified files strictly satisfy coding guidelines (<= 8-15 line functions, affirmative booleans, AppError envelopes).
3. Check live pipeline status or local test/linter verification ('gitmap pipeline status').
4. Confirm zero regressions and all quality checks pass before concluding.
`

// BuildVerificationFollowupPrompt creates the secondary verification prompt.
func BuildVerificationFollowupPrompt(repo string, runID uint64, sha string) string {
	var b strings.Builder
	b.WriteString("# Verification Check: Is It Fixed?\n\n")
	appendMetaLine(&b, "Target Repository", repo)
	if runID > 0 {
		appendMetaLine(&b, "Workflow Run ID", fmt.Sprintf("#%d", runID))
	}
	appendMetaLine(&b, "Commit SHA", sha)
	b.WriteString("\n/goal Verify whether the CI/CD pipeline failure and errors have been completely resolved and all quality gates pass.\n\n")
	b.WriteString(verificationInstructions)

	return b.String()
}

// StageVerificationFollowupPrompt stages the follow-up prompt to temp disk and queue file.
func StageVerificationFollowupPrompt(primaryPayload, followupPrompt string) error {
	rootDir := resolveProjectRootDir()
	followupPath := filepath.Join(rootDir, queuedAgyPromptRelativePath)
	_ = os.MkdirAll(filepath.Dir(followupPath), 0755)

	if writeErr := os.WriteFile(followupPath, []byte(followupPrompt), 0644); writeErr != nil {
		return writeErr
	}

	queuePath := filepath.Join(rootDir, agyPromptQueueRelativePath)

	return updatePromptQueueFile(queuePath, primaryPayload, followupPrompt)
}

func makePromptQueueFile(primary, followup, now string) AgyPromptQueueFile {
	return AgyPromptQueueFile{
		Active: &AgyPromptQueueEntry{
			ID: 1, Type: "primary_fix_rca", Title: "Fix CI/CD Pipeline Errors with 4-Part RCA",
			Prompt: primary, Status: "dispatched", CreatedAt: now,
		},
		Queued: []AgyPromptQueueEntry{
			{
				ID: 2, Type: "followup_verification", Title: "Verification Check: Is It Fixed?",
				Prompt: followup, Status: "queued", CreatedAt: now,
			},
		},
		UpdatedAt: now,
	}
}

func updatePromptQueueFile(queuePath, primary, followup string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	queueFile := makePromptQueueFile(primary, followup, now)
	marshaled, err := json.MarshalIndent(queueFile, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(queuePath, marshaled, 0644)
}
