package cmdschedule

import (
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

func matchFlagWithVal(arg string, names ...string) bool {
	for _, n := range names {
		if arg == n || strings.HasPrefix(arg, n+"=") {
			return true
		}
	}

	return false
}

func extractFlagValue(idx *int, args []string) string {
	a := args[*idx]
	if strings.Contains(a, "=") {
		parts := strings.SplitN(a, "=", 2)

		return parts[1]
	}

	if *idx+1 < len(args) {
		*idx++

		return args[*idx]
	}

	return ""
}

func extractMacroNameAndFlags(args []string) (string, []string) {
	return cmdmacro.ExtractMacroNameAndFlags(args)
}

func parseDurationArg(val string, fallback time.Duration) time.Duration {
	return cmdmacro.ParseDurationArg(val, fallback)
}

func parseExecOptions(args []string) macro.ExecOptions {
	return cmdmacro.ParseExecOptions(args)
}

func outputStructuredData(data interface{}, opts macro.ExecOptions) error {
	return cmdmacro.OutputStructuredData(data, opts)
}
