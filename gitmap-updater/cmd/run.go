package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// RunUpdate creates a handoff copy and launches the worker to perform the update.
func RunUpdate() {
	fmt.Print(MsgChecking)

	current, err := getInstalledVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, ErrGetVersion, err)
		os.Exit(1)
	}

	fmt.Printf(MsgCurrentVer, current)

	latest, err := fetchLatestTag()
	if err != nil {
		fmt.Fprintf(os.Stderr, ErrFetchRelease, err)
		os.Exit(1)
	}

	fmt.Printf(MsgLatestVer, latest)

	if normalizeVersion(current) == normalizeVersion(latest) {
		fmt.Printf(MsgUpToDate, current)
		return
	}

	fmt.Printf(MsgUpdateAvail, current, latest)
	fmt.Print(MsgHandoff)

	selfPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, ErrCreateCopy, err)
		os.Exit(1)
	}

	copyPath := createHandoffCopy(selfPath)
	launchWorker(copyPath, latest)
}

// createHandoffCopy creates a temporary copy of the updater binary.
func createHandoffCopy(selfPath string) string {
	name := fmt.Sprintf(UpdaterCopy, os.Getpid())
	copyPath := filepath.Join(filepath.Dir(selfPath), name)
	if copyBinary(selfPath, copyPath) == nil {
		return copyPath
	}

	fallbackPath := filepath.Join(os.TempDir(), name)
	if err := copyBinary(selfPath, fallbackPath); err != nil {
		fmt.Fprintf(os.Stderr, ErrCreateCopy, err)
		os.Exit(1)
	}

	return fallbackPath
}

// launchWorker runs the handoff copy with update-worker and the target version.
func launchWorker(copyPath, version string) {
	cmd := exec.Command(copyPath, "update-worker", version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}

		fmt.Fprintf(os.Stderr, ErrLaunchWorker, err)
		os.Exit(1)
	}
}

// copyBinary copies src to dst.
func copyBinary(src, dst string) error {
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
