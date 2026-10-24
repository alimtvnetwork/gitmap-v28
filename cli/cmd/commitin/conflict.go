package commitin

import (
	"fmt"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/commitin/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ConflictDecisionType is the canonical replay-loop decision type, defined in enums.
type ConflictDecisionType = enums.ConflictDecisionType

const (
	ConflictDecisionTakeTheirs = enums.ConflictDecisionTakeTheirs
	ConflictDecisionAbort      = enums.ConflictDecisionAbort
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
