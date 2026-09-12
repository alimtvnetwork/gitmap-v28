package cmdos

import (
	"os"

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
