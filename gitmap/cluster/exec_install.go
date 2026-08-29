package cluster

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// PackageResult represents the installation result of a single package.
type PackageResult struct {
	PackageName string
	Succeeded   bool
	Stderr      string
}

func detectPackageManager(ctx context.Context) (string, error) {
	isWindows := runtime.GOOS == constants.WindowsOS
	if isWindows {
		return detectWindowsPackageManager(ctx)
	}
	return detectUnixPackageManager(ctx)
}

func detectWindowsPackageManager(ctx context.Context) (string, error) {
	managers := []string{constants.PkgMgrWinget, constants.PkgMgrChocolatey}
	return checkManagers(ctx, managers)
}

func detectUnixPackageManager(ctx context.Context) (string, error) {
	managers := []string{constants.PkgMgrBrew, constants.PkgMgrApt}
	return checkManagers(ctx, managers)
}

func checkManagers(ctx context.Context, managers []string) (string, error) {
	for _, mgr := range managers {
		cmd := exec.CommandContext(ctx, mgr, constants.PkgMgrVersionArg)
		err := runCmdFunc(cmd)
		isFound := err == nil
		if isFound {
			return mgr, nil
		}
	}
	return "", errors.New(constants.ErrNoPackageManager)
}

// ExecInstall installs the specified packages using the detected package manager.
func ExecInstall(ctx context.Context, node ClusterNode, packages []string) ([]PackageResult, error) {
	mgr, err := detectPackageManager(ctx)
	hasError := err != nil
	if hasError {
		return nil, err
	}

	var results []PackageResult

	for _, pkg := range packages {
		var cmdStr string

		isWinget := mgr == constants.PkgMgrWinget
		isChoco := mgr == constants.PkgMgrChocolatey
		isBrew := mgr == constants.PkgMgrBrew
		isApt := mgr == constants.PkgMgrApt

		if isWinget {
			cmdStr = fmt.Sprintf(constants.FormatCmdSpace, mgr, constants.WingetInstallArg, constants.WingetQuietArg, pkg)
		} else if isChoco {
			cmdStr = fmt.Sprintf(constants.FormatCmdSpace, mgr, constants.ChocoInstallArg, constants.ChocoYesArg, pkg)
		} else if isBrew {
			cmdStr = fmt.Sprintf(constants.FormatCmdSpace, mgr, constants.BrewInstallArg, constants.BrewQuietArg, pkg)
		} else if isApt {
			cmdStr = fmt.Sprintf(constants.FormatCmdSpace, mgr, constants.AptInstallArg, constants.AptYesArg, pkg)
		} else {
			continue
		}

		_, stderr, exitCode, cmdErr := ExecCmd(ctx, node, cmdStr)

		isSuccessExit := exitCode == constants.ExitCodeSuccess
		noCmdErr := cmdErr == nil
		succeeded := isSuccessExit && noCmdErr

		res := PackageResult{
			PackageName: pkg,
			Succeeded:   succeeded,
			Stderr:      stderr,
		}
		results = append(results, res)
	}

	return results, nil
}
