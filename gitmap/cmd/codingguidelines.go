package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// CodingGuidelinesOpts controls a single Coding Guidelines v24 install run.
type CodingGuidelinesOpts struct {
	WorkingDir string
	Runner     func(name string, args ...string) *exec.Cmd
	LookPath   func(file string) (string, error)
	Stdout     io.Writer
	Stderr     io.Writer
	Stdin      io.Reader
}

// ErrCGShellNotFound is returned when the host lacks the shell required to execute the installer.
var ErrCGShellNotFound = errors.New("coding-guidelines: required shell not found on PATH")

// RunCodingGuidelinesInstall dispatches to the OS-appropriate installer.
func RunCodingGuidelinesInstall(opts CodingGuidelinesOpts) error {
	hasCustomRunner := opts.Runner != nil
	opts = withCGDefaults(opts)
	if runtime.GOOS == "windows" {
		return dispatchCGWindows(opts, hasCustomRunner)
	}

	return dispatchCGUnix(opts, hasCustomRunner)
}

// withCGDefaults fills in real-exec + os stdio for any zero-value fields.
func withCGDefaults(opts CodingGuidelinesOpts) CodingGuidelinesOpts {
	if opts.Runner == nil {
		opts.Runner = exec.Command
	}
	if opts.LookPath == nil {
		opts.LookPath = exec.LookPath
	}
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}
	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}
	if opts.Stdin == nil {
		opts.Stdin = os.Stdin
	}

	return opts
}

// dispatchCGWindows runs the v24 PowerShell installer with syntax compatibility patching.
func dispatchCGWindows(opts CodingGuidelinesOpts, hasCustomRunner bool) error {
	pwsh := resolvePowerShellBinaryWithLookPath(opts.LookPath)
	if pwsh == "" {
		fmt.Fprintf(opts.Stderr, constants.ErrCGShellNotFoundWindows, constants.DefaultCodingGuidelinesURLWindows)

		return ErrCGShellNotFound
	}

	url := constants.DefaultCodingGuidelinesURLWindows
	fmt.Fprintf(opts.Stderr, constants.MsgCGRunningWindows, url)
	cmd, cleanup, err := buildCGWindowsCommand(opts, pwsh, url, hasCustomRunner)
	if err != nil {
		fmt.Fprintf(opts.Stderr, constants.ErrCGCompatPrepareFailed, err)

		return err
	}
	defer cleanup()

	return runCGInstaller(cmd, opts, "windows", url)
}

func buildCGWindowsCommand(
	opts CodingGuidelinesOpts,
	pwsh, url string,
	hasCustomRunner bool,
) (*exec.Cmd, func(), error) {
	if hasCustomRunner {
		script := fmt.Sprintf("irm %s | iex", url)
		cmd := opts.Runner(pwsh, "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)

		return cmd, func() {}, nil
	}

	path, cleanup, err := writeCGCompatScriptWindows(url)
	if err != nil {
		return nil, cleanup, err
	}

	cmd := opts.Runner(pwsh, "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", path)

	return cmd, cleanup, nil
}

// dispatchCGUnix runs the v24 bash installer via curl.
func dispatchCGUnix(opts CodingGuidelinesOpts, hasCustomRunner bool) error {
	if _, err := opts.LookPath("bash"); err != nil {
		fmt.Fprintf(opts.Stderr, constants.ErrCGShellNotFoundUnix, constants.DefaultCodingGuidelinesURLUnix)

		return ErrCGShellNotFound
	}
	if _, err := opts.LookPath("curl"); err != nil {
		fmt.Fprintf(opts.Stderr, constants.ErrCGShellNotFoundUnix, constants.DefaultCodingGuidelinesURLUnix)

		return ErrCGShellNotFound
	}

	url := constants.DefaultCodingGuidelinesURLUnix
	fmt.Fprintf(opts.Stderr, constants.MsgCGRunningUnix, url)
	cmd, cleanup, err := buildCGUnixCommand(opts, url, hasCustomRunner)
	if err != nil {
		fmt.Fprintf(opts.Stderr, constants.ErrCGCompatPrepareFailed, err)

		return err
	}
	defer cleanup()

	return runCGInstaller(cmd, opts, runtime.GOOS, url)
}

func buildCGUnixCommand(
	opts CodingGuidelinesOpts,
	url string,
	hasCustomRunner bool,
) (*exec.Cmd, func(), error) {
	if hasCustomRunner {
		script := fmt.Sprintf("curl -fsSL %s | bash", url)

		return opts.Runner("bash", "-c", script), func() {}, nil
	}

	path, cleanup, err := writeCGCompatScript(url)
	if err != nil {
		return nil, cleanup, err
	}

	return opts.Runner("bash", path), cleanup, nil
}

// runCGInstaller wires stdio + working dir onto the prepared command and executes it.
func runCGInstaller(cmd *exec.Cmd, opts CodingGuidelinesOpts, goos, url string) error {
	cmd.Dir = opts.WorkingDir
	cmd.Stdout = opts.Stdout
	cmd.Stderr = opts.Stderr
	cmd.Stdin = opts.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(opts.Stderr, constants.ErrCGInstallFailed, goos, err)

		return fmt.Errorf("coding-guidelines install (%s, %s): %w", goos, url, err)
	}

	fmt.Fprint(opts.Stderr, constants.MsgCGDone)

	return nil
}
