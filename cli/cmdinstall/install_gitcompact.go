package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const (
	gitCompactWinScript  = "irm https://raw.githubusercontent.com/alimtvnetwork/git-compact/main/install.ps1 | iex"
	gitCompactUnixScript = "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/git-compact/main/install.sh | sh -s -- --dir ~/.local/bin"
)

func handleGitCompactInstall(opts installOptions) {
	if err := runInstallGitCompact(opts); err != nil {
		cliexit.HandleError(apperror.WrapSimple(err, "cmdinstall.installGitCompact"), 1)
	}
}

func runInstallGitCompact(opts installOptions) error {
	if opts.DryRun {
		printGitCompactDryRun()

		return nil
	}

	if opts.Check {
		return checkGitCompactStatus()
	}

	return performGitCompactInstall(opts)
}

func printGitCompactDryRun() {
	cmdStr := gitCompactUnixScript
	if runtime.GOOS == "windows" {
		cmdStr = gitCompactWinScript
	}

	fmt.Printf(constants.MsgInstallDryCmd, cmdStr)
}

func checkGitCompactStatus() error {
	bin, ver := resolveToolProbeCommand(constants.ToolGitCompact)
	if bin != "" {
		fmt.Printf(constants.MsgInstallFound, constants.ToolGitCompact, ver)

		return nil
	}

	fmt.Printf(constants.MsgInstallNotFound, constants.ToolGitCompact)

	return nil
}

func performGitCompactInstall(opts installOptions) error {
	fmt.Printf(constants.MsgInstallInstalling, constants.ToolGitCompact)
	start := time.Now()
	err := executeGitCompactInstall(opts)
	durMs := time.Since(start).Milliseconds()
	isSuccess := err == nil
	recordGitCompactTelemetry(isSuccess, durMs, err)
	if err != nil {
		return err
	}

	printGitCompactSuccess()

	return nil
}

func printGitCompactSuccess() {
	fmt.Printf(constants.MsgInstallSuccess, constants.ToolGitCompact)
}

func executeGitCompactInstall(opts installOptions) error {
	cmd := buildGitCompactCommand()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return apperror.NewWithDetails("cmdinstall.gitcompact", "E4001",
			fmt.Sprintf("git-compact install failed: %v", err),
			"cmdinstall", apperror.ErrorTypeExecution, apperror.SeverityError, nil).WithCause(err)
	}

	return nil
}

func buildGitCompactCommand() *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("powershell", "-NoProfile", "-Command", gitCompactWinScript)
	}

	return exec.Command("sh", "-c", gitCompactUnixScript)
}

func recordGitCompactTelemetry(isSuccess bool, durMs int64, cmdErr error) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return
	}

	defer splitDB.Close()
	saveGitCompactInstalledRecord(splitDB, isSuccess)
	saveGitCompactAuditLog(splitDB, isSuccess, durMs, cmdErr)
}

func saveGitCompactInstalledRecord(splitDB *store.InstallationSplitDB, isSuccess bool) {
	if isSuccess {
		recordGitCompactToolVersion(splitDB)
	}
}

func recordGitCompactToolVersion(splitDB *store.InstallationSplitDB) {
	_, ver := resolveToolProbeCommand(constants.ToolGitCompact)
	if ver == "" {
		ver = "installed"
	}

	_ = splitDB.SaveInstalledTool(constants.ToolGitCompact, ver, "script")
}

func saveGitCompactAuditLog(splitDB *store.InstallationSplitDB, isSuccess bool, durMs int64, cmdErr error) {
	exitCode, errMsg := resolveGitCompactErrorDetails(cmdErr)
	_ = splitDB.RecordLog(store.InstallationLogRecord{
		Tool:           constants.ToolGitCompact,
		Action:         "install",
		PackageManager: "script",
		DurationMs:     durMs,
		IsSuccess:      isSuccess,
		ExitCode:       exitCode,
		Stderr:         errMsg,
	})
}

func resolveGitCompactErrorDetails(cmdErr error) (int, string) {
	if cmdErr != nil {
		return 1, cmdErr.Error()
	}

	return 0, ""
}
