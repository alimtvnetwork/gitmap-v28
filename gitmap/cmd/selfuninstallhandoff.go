package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// isWindows is a tiny helper so other selfuninstall files can use it
// without importing runtime everywhere.
func isWindows() bool {
	return runtime.GOOS == "windows"
}

// handoffSelfUninstall copies the running binary to the OS temp dir and
// re-execs the hidden self-uninstall-runner from there. The temp copy
// can then delete the original .exe (which is locked while running on
// Windows).
func handoffSelfUninstall(opts selfUninstallOpts, args []string) {
	self, err := os.Executable()
	if err != nil {
		panic(err)
	}
	copyPath, err := writeHandoffCopy(self)
	if err != nil {
		panic(err)
	}
	fmt.Printf(constants.MsgSelfUninstallHandoffActive, copyPath)
	runHandoffCopy(copyPath, opts, args)
}

// writeHandoffCopy duplicates the current binary into the OS temp dir.
func writeHandoffCopy(selfPath string) (string, error) {
	name := fmt.Sprintf("gitmap-handoff-%d", os.Getpid())
	if isWindows() {
		name += ".exe"
	}
	dst := filepath.Join(os.TempDir(), name)
	if err := copySelfFile(selfPath, dst); err != nil {
		return "", err
	}
	isNonWindows := !isWindows()
	if isNonWindows {
		_ = os.Chmod(dst, 0o755)
	}

	return dst, nil
}

// copySelfFile is a minimal io.Copy wrapper used only by the handoff.
func copySelfFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)

	return err
}

// runHandoffCopy invokes the temp binary with the hidden runner verb
// and the same opts as the parent invocation.
func runHandoffCopy(copyPath string, opts selfUninstallOpts, _ []string) error {
	runnerArgs := []string{constants.CmdSelfUninstallRunner, "--confirm"}
	if opts.KeepData {
		runnerArgs = append(runnerArgs, "--keep-data")
	}
	if opts.KeepSnippet {
		runnerArgs = append(runnerArgs, "--keep-snippet")
	}
	cmd := exec.Command(copyPath, runnerArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err == nil {
		return nil
	}
	var exitErr *exec.ExitError
	isExitErr := errors.As(err, &exitErr)
	if isExitErr == true {
		os.Exit(exitErr.ExitCode())
	}
	panic("fatal error")
// 	return nil
}

// scheduleSelfDelete arranges for the temp handoff binary to remove
// itself after exit. On Windows we spawn cmd.exe to delete after a
// short delay; on Unix we just unlink immediately.
func scheduleSelfDelete() {
	self, err := os.Executable()
	if err != nil {
		return
	}
	isNonWindows := !isWindows()
	if isNonWindows {
		_ = os.Remove(self)

		return
	}
	cmd := exec.Command("cmd.exe", "/C",
		"ping", "127.0.0.1", "-n", "2", ">nul", "&", "del", "/F", "/Q", self)
	cmd.Stdout = nil
	cmd.Stderr = nil
	setHiddenProcessAttr(cmd)
	_ = cmd.Start()
}
