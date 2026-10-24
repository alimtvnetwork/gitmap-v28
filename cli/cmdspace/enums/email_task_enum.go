package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// EmailTaskType identifies a gitmap operation that can fire a registered
// email trigger when it completes.
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

// IsValid reports whether the email task matches a known task.
func (t EmailTaskType) IsValid() bool {
	switch t {
	case EmailTaskRelease, EmailTaskPull, EmailTaskClone, EmailTaskScan, EmailTaskPush, EmailTaskBackup:
		return true
	default:
		return false
	}
}

// String returns the string representation of the email task.
func (t EmailTaskType) String() string {
	return string(t)
}

// ParseEmailTaskType maps a user-supplied task string to its enum value.
// Matching is case-insensitive; the error lists every valid value so the
// caller can surface it directly.
func ParseEmailTaskType(s string) (EmailTaskType, *apperror.AppError) {
	normalized := strings.ToLower(strings.TrimSpace(s))
	for _, task := range AllEmailTaskTypes() {
		if task.String() == normalized {
			return task, nil
		}
	}

	return "", apperror.NewValidation("ParseEmailTaskType", "E1000", fmt.Sprintf("unknown email task %q: valid values: %s", s, strings.Join(EmailTaskTypeNames(), ", "))).WithContext("value", s)
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

// EmailTaskTypeNames returns the raw task names for error messages.
func EmailTaskTypeNames() []string {
	tasks := AllEmailTaskTypes()
	names := make([]string, len(tasks))
	for i, task := range tasks {
		names[i] = task.String()
	}

	return names
}

// IsRelease reports whether the task is release.
func (t EmailTaskType) IsRelease() bool { return t == EmailTaskRelease }

// IsPull reports whether the task is pull.
func (t EmailTaskType) IsPull() bool { return t == EmailTaskPull }

// IsClone reports whether the task is clone.
func (t EmailTaskType) IsClone() bool { return t == EmailTaskClone }

// IsScan reports whether the task is scan.
func (t EmailTaskType) IsScan() bool { return t == EmailTaskScan }

// IsPush reports whether the task is push.
func (t EmailTaskType) IsPush() bool { return t == EmailTaskPush }

// IsBackup reports whether the task is backup.
func (t EmailTaskType) IsBackup() bool { return t == EmailTaskBackup }
