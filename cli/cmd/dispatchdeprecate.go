package cmd

// Deprecation enforcement for the command dispatch path.
//
// Per docs/deprecation-policy.md a deprecated command keeps working
// (aliases never break), but every invocation prints a one-line stderr
// warning naming the replacement command and the version it will be
// removed in (two minor releases out). The registry is consulted from
// runDispatchTable, so every dispatcher (core/data/tooling/utility/
// release) enforces the policy with no per-site code.

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// deprecationInfo is the registry payload for one deprecated command name.
type deprecationInfo struct {
	replacement string
	removedIn   string
}

// deprecatedCommands maps a deprecated command name to its deprecation
// metadata. Register entries with markDeprecated; warnIfDeprecated
// consults it from runDispatchTable before the handler executes.
var deprecatedCommands = map[string]deprecationInfo{}

// markDeprecated registers oldName as deprecated in favor of
// replacement, removed in removedIn. Use removalTwoReleasesOut for
// the policy default (two minor releases after the current version).
func markDeprecated(oldName, replacement, removedIn string) {
	deprecatedCommands[oldName] = deprecationInfo{replacement: replacement, removedIn: removedIn}
}

// warnIfDeprecated prints the policy-format warning to stderr when
// command is registered as deprecated. The command still executes
// normally afterwards — the warning is advisory, never a gate.
func warnIfDeprecated(command string) {
	info, ok := deprecatedCommands[command]
	if !ok {
		return
	}

	fmt.Fprintf(os.Stderr, "  · %s is deprecated; use %s instead (removal planned in v%s)\n",
		command, info.replacement, info.removedIn)
}

// removalTwoReleasesOut returns the version two minor releases after
// constants.Version — the earliest removal version the deprecation
// policy allows for a newly deprecated command (e.g. 6.550.0 →
// 6.552.0). Falls back to a relative label when the version does not
// parse as major.minor.patch.
func removalTwoReleasesOut() string {
	parts := strings.Split(constants.Version, ".")
	if len(parts) != 3 {
		return constants.Version + " +2 releases"
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return constants.Version + " +2 releases"
	}

	return parts[0] + "." + strconv.Itoa(minor+2) + "." + parts[2]
}
