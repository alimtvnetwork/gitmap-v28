package cmdhooks

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhooks/constants"
)

// commitMsgGuard is the commit-msg hook installed by `hooks install --all`.
// It enforces the house rule: no colons in the commit subject line.
const commitMsgGuard = constants.CommitMsgGuardTemplate

// guardMarker identifies hooks previously installed by this command.
const guardMarker = constants.GuardMarker
