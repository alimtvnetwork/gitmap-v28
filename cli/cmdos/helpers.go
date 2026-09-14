package cmdos

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

func expandHome(p string) string {
	return cmdscan.ExpandHome(p)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func openDB() (*store.DB, error) {
	return store.OpenDefault()
}

func runPowerNeverSleep() error {
	if RunPowerNeverSleepFn != nil {
		return RunPowerNeverSleepFn()
	}

	return nil
}

func runPowerSet(args []string) error {
	if RunPowerSetFn != nil {
		return RunPowerSetFn(args)
	}

	return nil
}

func runPowerReset() error {
	if RunPowerResetFn != nil {
		return RunPowerResetFn()
	}

	return nil
}

func hasFlag(args []string, flagName string) bool {
	for _, arg := range args {
		isExact := arg == flagName
		isPrefix := strings.HasPrefix(arg, flagName+"=")
		if isExact || isPrefix {
			return true
		}
	}
	return false
}

func extractArgFlag(args []string, flagName string) string {
	for i, arg := range args {
		val := resolveArgFlagValue(args, i, arg, flagName)
		if len(val) > 0 {
			return val
		}
	}
	return ""
}

func resolveArgFlagValue(args []string, i int, arg, flagName string) string {
	if arg == flagName && i+1 < len(args) {
		return args[i+1]
	}
	prefix := flagName + "="
	if strings.HasPrefix(arg, prefix) {
		return strings.TrimPrefix(arg, prefix)
	}
	return ""
}

func extractArgFlagWithDefault(args []string, flagName, defaultVal string) string {
	val := extractArgFlag(args, flagName)
	hasVal := len(val) > 0
	if hasVal {
		return val
	}
	return defaultVal
}
