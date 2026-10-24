// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdspace/enums"
)

// EmailTaskType is the canonical email trigger task type, defined in enums.
type EmailTaskType = enums.EmailTaskType

const (
	// EmailTaskRelease fires when `gitmap release` completes.
	EmailTaskRelease = enums.EmailTaskRelease
	// EmailTaskPull fires when `gitmap pull` completes.
	EmailTaskPull = enums.EmailTaskPull
	// EmailTaskClone fires when `gitmap clone` completes.
	EmailTaskClone = enums.EmailTaskClone
	// EmailTaskScan fires when `gitmap scan` completes.
	EmailTaskScan = enums.EmailTaskScan
	// EmailTaskPush fires when `gitmap push` completes.
	EmailTaskPush = enums.EmailTaskPush
	// EmailTaskBackup fires when `gitmap backup` completes.
	EmailTaskBackup = enums.EmailTaskBackup
)

// ParseEmailTaskType maps a user-supplied task string to its enum value.
// It returns an enriched app error (see enums.ParseEmailTaskType).
func ParseEmailTaskType(s string) (EmailTaskType, *apperror.AppError) {
	return enums.ParseEmailTaskType(s)
}

// AllEmailTaskTypes returns every supported email trigger task, in the
// canonical order used by help text and validation messages.
func AllEmailTaskTypes() []EmailTaskType {
	return enums.AllEmailTaskTypes()
}
