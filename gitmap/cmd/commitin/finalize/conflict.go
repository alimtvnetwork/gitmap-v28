package finalize

import (
	"fmt"
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// ConflictDecision is the action a higher-level replay loop takes for
// a single conflicting source commit.
type ConflictDecision int

const (
	ConflictDecisionTakeTheirs ConflictDecision = iota
	ConflictDecisionAbort
)

// Resolve maps the resolved ConflictMode to a ConflictDecision.
// Prompt mode prints the standardized abort banner — callers translate
// ConflictDecisionAbort into exit code CommitInExitConflictAborted.
func Resolve(mode, sourceSha string, out io.Writer) ConflictDecision {
	if mode == constants.CommitInConflictModeForceMerge {
		return ConflictDecisionTakeTheirs
	}
	fmt.Fprintf(out, constants.CommitInErrConflictAborted, sourceSha)
	return ConflictDecisionAbort
}
