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

func resolveQueuedAgyPromptPath() string {
	return filepath.Join(resolveProjectRootDir(), queuedAgyPromptRelativePath)
}

func resolveAgyPromptQueuePath() string {
	return filepath.Join(resolveProjectRootDir(), agyPromptQueueRelativePath)
}

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
	followupPath := resolveQueuedAgyPromptPath()
	_ = os.MkdirAll(filepath.Dir(followupPath), 0755)

	if writeErr := os.WriteFile(followupPath, []byte(followupPrompt), 0644); writeErr != nil {
		return writeErr
	}

	queuePath := resolveAgyPromptQueuePath()

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

// LoadPromptQueue loads the prompt queue file from disk.
func LoadPromptQueue() (AgyPromptQueueFile, error) {
	queuePath := resolveAgyPromptQueuePath()
	data, err := os.ReadFile(queuePath)
	if err != nil {
		return AgyPromptQueueFile{Queued: make([]AgyPromptQueueEntry, 0)}, nil
	}

	var q AgyPromptQueueFile
	if unmarshalErr := json.Unmarshal(data, &q); unmarshalErr != nil {
		return AgyPromptQueueFile{Queued: make([]AgyPromptQueueEntry, 0)}, unmarshalErr
	}

	return q, nil
}

// SavePromptQueue persists the prompt queue file to disk.
func SavePromptQueue(q AgyPromptQueueFile) error {
	queuePath := resolveAgyPromptQueuePath()
	_ = os.MkdirAll(filepath.Dir(queuePath), 0755)
	q.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	marshaled, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(queuePath, marshaled, 0644)
}

func computeNextQueueID(q AgyPromptQueueFile) int {
	nextID := 1
	if q.Active != nil && q.Active.ID >= nextID {
		nextID = q.Active.ID + 1
	}
	for _, item := range q.Queued {
		if item.ID >= nextID {
			nextID = item.ID + 1
		}
	}

	return nextID
}

// EnqueuePrompt appends a new prompt entry to the queue.
func EnqueuePrompt(promptType, title, promptText string) (*AgyPromptQueueEntry, error) {
	q, err := LoadPromptQueue()
	if err != nil {
		return nil, err
	}
	entry := AgyPromptQueueEntry{
		ID: computeNextQueueID(q), Type: promptType, Title: title,
		Prompt: promptText, Status: "queued",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	q.Queued = append(q.Queued, entry)

	return &entry, SavePromptQueue(q)
}

// GetQueueStatus returns the current prompt queue file.
func GetQueueStatus() (AgyPromptQueueFile, error) {
	return LoadPromptQueue()
}

// PopNextQueuedPrompt pops the next queued prompt and marks it active.
func PopNextQueuedPrompt() (*AgyPromptQueueEntry, error) {
	q, err := LoadPromptQueue()
	if err != nil {
		return nil, err
	}
	if len(q.Queued) == 0 {
		return nil, nil
	}

	popped := q.Queued[0]
	q.Queued = q.Queued[1:]
	popped.Status = "dispatched"
	q.Active = &popped

	if saveErr := SavePromptQueue(q); saveErr != nil {
		return nil, saveErr
	}
	_ = stagePoppedPromptToActiveFile(popped.Prompt)

	return &popped, nil
}

func stagePoppedPromptToActiveFile(promptText string) error {
	activePath := resolveActiveAgyPromptPath()
	_ = os.MkdirAll(filepath.Dir(activePath), 0755)
	copyClipboardIfNotSkipped(promptText, false)

	return os.WriteFile(activePath, []byte(promptText), 0644)
}

// ClearPromptQueue empties all queued and active items from the prompt queue.
func ClearPromptQueue() error {
	q, err := LoadPromptQueue()
	if err != nil {
		return err
	}
	q.Queued = make([]AgyPromptQueueEntry, 0)
	q.Active = nil

	return SavePromptQueue(q)
}
