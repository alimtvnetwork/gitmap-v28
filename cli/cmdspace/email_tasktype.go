// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// EmailTaskType identifies a gitmap operation that can fire a registered
// email trigger when it completes (subtask 06).
type EmailTaskType string

const (
	// EmailTaskRelease fires when `gitmap release` completes.
	EmailTaskRelease EmailTaskType = "release"
	// EmailTaskPull fires when `gitmap pull` completes.
	EmailTaskPull EmailTaskType = "pull"
	// EmailTaskClone fires when `gitmap clone` completes.
	EmailTaskClone EmailTaskType = "clone"
	// EmailTaskScan fires when `gitmap scan` completes.
	EmailTaskScan EmailTaskType = "scan"
	// EmailTaskPush fires when `gitmap push` completes.
	EmailTaskPush EmailTaskType = "push"
	// EmailTaskBackup fires when `gitmap backup` completes.
	EmailTaskBackup EmailTaskType = "backup"
)

// ParseEmailTaskType maps a user-supplied task string to its enum value.
// Matching is case-insensitive; the error lists every valid value so the
// caller can surface it directly.
func ParseEmailTaskType(s string) (EmailTaskType, error) {
	normalized := strings.ToLower(strings.TrimSpace(s))
	for _, task := range AllEmailTaskTypes() {
		if string(task) == normalized {
			return task, nil
		}
	}

	return "", apperror.NewValidationError(fmt.Sprintf("unknown email task %q: valid values: %s", s, strings.Join(emailTaskTypeNames(), ", ")))
}

// AllEmailTaskTypes returns every supported email trigger task, in the
// canonical order used by help text and validation messages.
func AllEmailTaskTypes() []EmailTaskType {
	return []EmailTaskType{
		EmailTaskRelease,
		EmailTaskPull,
		EmailTaskClone,
		EmailTaskScan,
		EmailTaskPush,
		EmailTaskBackup,
	}
}

// emailTaskTypeNames returns the raw task names for error messages.
func emailTaskTypeNames() []string {
	tasks := AllEmailTaskTypes()
	names := make([]string, len(tasks))
	for i, task := range tasks {
		names[i] = string(task)
	}

	return names
}
