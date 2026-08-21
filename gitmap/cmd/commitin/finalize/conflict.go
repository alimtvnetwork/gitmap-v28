package finalize

import (
	"fmt"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ConflictDecisionType is the action a higher-level replay loop takes for
// a single conflicting source commit.
type ConflictDecisionType int

const (
	ConflictDecisionTakeTheirs ConflictDecisionType = iota
	ConflictDecisionAbort
)

// Resolve maps the resolved ConflictMode to a ConflictDecisionType.
// Prompt mode prints the standardized abort banner — callers translate
// ConflictDecisionAbort into exit code CommitInExitConflictAborted.
func Resolve(mode, sourceSha string, out io.Writer) ConflictDecisionType {
	if mode == constants.CommitInConflictModeForceMerge {
		return ConflictDecisionTakeTheirs
	}
	fmt.Fprintf(out, constants.CommitInErrConflictAborted, sourceSha)

	return ConflictDecisionAbort
}
