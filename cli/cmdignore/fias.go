package cmdignore

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
)

func init() {
	RunFixIgnoresAllSSHFn = runFiasWithPAS
}

func runFiasWithPAS(args []string) *apperror.AppError {
	remoteCmd := getRemoteCmdString(args)
	err := cmdssh.RunFleetPASCommand("Fix Ignore All", remoteCmd, func() error {
		return executeLocalFix(args)
	})
	hasErr := err != nil
	if hasErr {
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
