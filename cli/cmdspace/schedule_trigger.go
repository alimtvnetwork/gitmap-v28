// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// settingEmailTaskTriggers is the Setting-table key holding the trigger
// list as a JSON blob. Defined here (not in constants) per subtask 06's
// contract — the sibling worker owns the shared constants files.
const settingEmailTaskTriggers = "email_task_triggers"

const (
	// triggerTimeFormat renders outcome timestamps in email bodies.
	triggerTimeFormat = "2006-01-02 15:04:05"
	// triggerSuccessSummary is the outcome summary for clean completions.
	triggerSuccessSummary = "completed successfully"
	// unknownMachineValue is the machine-info fallback — never an error.
	unknownMachineValue = "unknown"
	// statusOKText and statusFailText label the email status banner.
	statusOKText   = "SUCCESS"
	statusFailText = "FAILED"
)

// EmailTrigger is one registered email notification: when task completes,
// an HTML status email is sent to To.
type EmailTrigger struct {
	Name string        `json:"name"`
	To   string        `json:"to"`
	Task EmailTaskType `json:"task"`
}

// TaskOutcome describes a finished gitmap operation for trigger emails.
type TaskOutcome struct {
	Repo       string `json:"repo"`
	Branch     string `json:"branch"`
	Operation  string `json:"operation"`
	IsSuccess  bool   `json:"is_success"`
	Summary    string `json:"summary"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
}

// openTriggerStore opens the settings database for trigger persistence.
// A variable (not a direct call) so tests can point it at a temp DB.
var openTriggerStore = store.OpenDefault

// sendTriggerEmail sends one rendered trigger email; stubbed in tests.
var sendTriggerEmail = defaultSendTriggerEmail

// LoadEmailTriggers reads all registered triggers (empty when none).
// A corrupt blob surfaces as an error — it is trigger data, not
// unrelated state, so "behave as absent" would hide user config loss.
func LoadEmailTriggers(db *store.DB) ([]EmailTrigger, error) {
	raw := strings.TrimSpace(db.GetSetting(settingEmailTaskTriggers))
	if raw == "" {
		return []EmailTrigger{}, nil
	}

	var triggers []EmailTrigger
	if err := json.Unmarshal([]byte(raw), &triggers); err != nil {
		return nil, apperror.WrapSimple(err, "parse email task triggers:")
	}

	return triggers, nil
}

// SaveEmailTriggers persists the trigger list as a JSON blob.
func SaveEmailTriggers(db *store.DB, triggers []EmailTrigger) error {
	raw, err := json.Marshal(triggers)
	if err != nil {
		return apperror.WrapSimple(err, "encode email task triggers:")
	}

	if err := db.SetSetting(settingEmailTaskTriggers, string(raw)); err != nil {
		return apperror.WrapSimple(err, "save email task triggers:")
	}

	return nil
}

// RegisterEmailTrigger adds a trigger, replacing any existing trigger
// that shares its name.
func RegisterEmailTrigger(db *store.DB, trigger EmailTrigger) error {
	triggers, err := LoadEmailTriggers(db)
	if err != nil {
		return err
	}

	return SaveEmailTriggers(db, replaceOrAppendTrigger(triggers, trigger))
}

// replaceOrAppendTrigger swaps the same-named trigger or appends.
func replaceOrAppendTrigger(triggers []EmailTrigger, trigger EmailTrigger) []EmailTrigger {
	for i, existing := range triggers {
		if existing.Name == trigger.Name {
			triggers[i] = trigger

			return triggers
		}
	}

	return append(triggers, trigger)
}

// RemoveEmailTrigger deletes the trigger with the given name.
func RemoveEmailTrigger(db *store.DB, name string) error {
	triggers, err := LoadEmailTriggers(db)
	if err != nil {
		return err
	}

	kept := filterOutTrigger(triggers, name)
	if len(kept) == len(triggers) {
		return fmt.Errorf("no email trigger named %q", name)
	}

	return SaveEmailTriggers(db, kept)
}

// filterOutTrigger drops the named trigger from the list.
func filterOutTrigger(triggers []EmailTrigger, name string) []EmailTrigger {
	kept := make([]EmailTrigger, 0, len(triggers))
	for _, trigger := range triggers {
		if trigger.Name != name {
			kept = append(kept, trigger)
		}
	}

	return kept
}

// TriggersForTask returns the triggers registered for one task type.
func TriggersForTask(triggers []EmailTrigger, task EmailTaskType) []EmailTrigger {
	var matched []EmailTrigger
	for _, trigger := range triggers {
		if trigger.Task == task {
			matched = append(matched, trigger)
		}
	}

	return matched
}

// InferTaskOutcome builds a TaskOutcome from an operation result,
// detecting repo and branch from the current git work tree. It never
// errors: every detection has a fallback.
func InferTaskOutcome(task EmailTaskType, err error) TaskOutcome {
	summary := triggerSuccessSummary
	isSuccess := true
	if err != nil {
		summary = err.Error()
		isSuccess = false
	}

	now := time.Now().UTC().Format(triggerTimeFormat)

	return TaskOutcome{
		Repo:       detectRepoName(),
		Branch:     detectBranchName(),
		Operation:  string(task),
		IsSuccess:  isSuccess,
		Summary:    summary,
		StartedAt:  now,
		FinishedAt: now,
	}
}

// detectRepoName returns the enclosing git repo's directory name, falling
// back to the cwd basename, then "unknown" — never an error.
func detectRepoName() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err == nil {
		return filepath.Base(strings.TrimSpace(string(out)))
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr == nil {
		return filepath.Base(cwd)
	}

	return unknownMachineValue
}

// detectBranchName returns the current branch, or "" outside a repo.
func detectBranchName() string {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// FireEmailTaskTriggers emails every trigger registered for the task.
// It never panics and never returns an error: a trigger must never break
// the operation that fired it — send failures go to stderr only.
func FireEmailTaskTriggers(task EmailTaskType, outcome TaskOutcome) {
	defer func() {
		_ = recover()
	}()

	fireEmailTaskTriggers(task, outcome)
}

// fireEmailTaskTriggers loads matching triggers and fires each one.
func fireEmailTaskTriggers(task EmailTaskType, outcome TaskOutcome) {
	db, err := openTriggerStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ email triggers skipped: cannot open settings db: %v\n", err)

		return
	}
	defer db.Close()

	triggers, loadErr := LoadEmailTriggers(db)
	if loadErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ email triggers skipped: %v\n", loadErr)

		return
	}

	info := CollectMachineInfo()
	for _, trigger := range TriggersForTask(triggers, task) {
		fireSingleTrigger(trigger, task, outcome, info)
	}
}

// fireSingleTrigger renders and sends one trigger email, reporting
// failures to stderr without disturbing the firing operation.
func fireSingleTrigger(trigger EmailTrigger, task EmailTaskType, outcome TaskOutcome, info MachineInfo) {
	subject := buildTriggerSubject(task, outcome)
	body := BuildTriggerEmailHTML(task, outcome, info)
	if err := sendTriggerEmail(trigger.To, subject, body); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ email trigger %q skipped: %v\n", trigger.Name, err)

		return
	}

	fmt.Printf("  ✉ email trigger %q notified %s\n", trigger.Name, trigger.To)
}

// buildTriggerSubject labels the notification with task + status.
func buildTriggerSubject(task EmailTaskType, outcome TaskOutcome) string {
	status := statusOKText
	if !outcome.IsSuccess {
		status = statusFailText
	}

	return fmt.Sprintf("[gitmap] %s %s: %s", task, status, outcome.Repo)
}
