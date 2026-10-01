package cmdignore

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// GitMap PAS Formula concurrency constraints:
// Max 2 workers, 2 async operations max per remote node.
const (
	MaxPASWorkers      = 2
	MaxAsyncOpsPerNode = 2
)

func init() {
	RunFixIgnoresAllSSHFn = RunFiasWithPAS
}

// RunFiasWithPAS executes fix-ignore-all across the fleet using the PAS Formula.
func RunFiasWithPAS(args []string) *apperror.AppError {
	remoteCmd := getRemoteCmdString(args)
	err := cmdssh.RunFleetPASCommand("Fix Ignore All", remoteCmd, func() error {
		return executeLocalFix(args)
	})
	hasErr := err != nil
	if hasErr {
		store.LogInternalError("FIX_IGNORE_SSH", "PAS_EXECUTION_FAILED", err.Error(), "", "")
		return apperror.WrapSimple(err, "PAS execution failed")
	}
	return nil
}

func getRemoteCmdString(args []string) string {
	remoteCmd := "gitmap fix-ignore-all"
	isYes := isAutoYesArg(args)
	if isYes {
		remoteCmd += " -y"
	}
	return remoteCmd
}

func executeLocalFix(args []string) error {
	localErr := RunFixIgnoreAll(args)
	hasErr := localErr != nil
	if hasErr {
		return localErr
	}
	return nil
}
