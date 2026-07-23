package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// CodingGuidelinesOpts controls a single Coding Guidelines v24 install run.
// The Runner factory is injectable so unit tests can substitute a fake
// *exec.Cmd without shelling out to the network. Zero-value opts are valid
// and default to real exec + os stdio.
type CodingGuidelinesOpts struct {
	WorkingDir string
	Runner     func(name string, args ...string) *exec.Cmd
	Stdout     io.Writer
	Stderr     io.Writer
	Stdin      io.Reader
}

// ErrCGShellNotFound is returned when the host lacks the shell required to
// execute the OS-appropriate installer (PowerShell on Windows; bash+curl
// elsewhere). Callers can detect it with errors.Is and surface an
// actionable copy-paste fallback.
var ErrCGShellNotFound = errors.New("coding-guidelines: required shell not found on PATH")

var cgPostfixIncrementPattern = regexp.MustCompile(`\(\(([A-Za-z_][A-Za-z0-9_]*)\+\+\)\)`)

// RunCodingGuidelinesInstall dispatches to the OS-appropriate installer
// (PowerShell on Windows, bash on Unix) and streams stdout/stderr through
// the provided writers. All failures are logged to opts.Stderr per the
// zero-swallow error policy before being returned.
func RunCodingGuidelinesInstall(opts CodingGuidelinesOpts) error {
	hasCustomRunner := opts.Runner != nil
	opts = withCGDefaults(opts)
	if runtime.GOOS == "windows" {
		return dispatchCGWindows(opts)
	}

	return dispatchCGUnix(opts, hasCustomRunner)
}

// withCGDefaults fills in real-exec + os stdio for any zero-value fields.
func withCGDefaults(opts CodingGuidelinesOpts) CodingGuidelinesOpts {
	if opts.Runner == nil {
		opts.Runner = exec.Command
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

// dispatchCGWindows runs the v24 PowerShell installer via `irm | iex`.
func dispatchCGWindows(opts CodingGuidelinesOpts) error {
	pwsh := resolvePowerShellBinary()
	if pwsh == "" {
		fmt.Fprintf(opts.Stderr, constants.ErrCGShellNotFoundWindows, constants.DefaultCodingGuidelinesURLWindows)
		return ErrCGShellNotFound
	}
	url := constants.DefaultCodingGuidelinesURLWindows
	fmt.Fprintf(opts.Stderr, constants.MsgCGRunningWindows, url)
	script := fmt.Sprintf("irm %s | iex", url)
	cmd := opts.Runner(pwsh, "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)

	return runCGInstaller(cmd, opts, "windows", url)
}

// dispatchCGUnix runs the v24 bash installer via `curl -fsSL | bash`.
func dispatchCGUnix(opts CodingGuidelinesOpts, hasCustomRunner bool) error {
	if _, err := exec.LookPath("bash"); err != nil {
		fmt.Fprintf(opts.Stderr, constants.ErrCGShellNotFoundUnix, constants.DefaultCodingGuidelinesURLUnix)
		return ErrCGShellNotFound
	}
	if _, err := exec.LookPath("curl"); err != nil {
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

func buildCGUnixCommand(opts CodingGuidelinesOpts, url string, hasCustomRunner bool) (*exec.Cmd, func(), error) {
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

func writeCGCompatScript(url string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "gitmap-cg-*")
	if err != nil {
		return "", func() {}, err
	}
	path := filepath.Join(dir, "install.sh")
	if err := downloadCGScript(url, path); err != nil {
		_ = os.RemoveAll(dir)
		return "", func() {}, err
	}
	return path, func() { _ = os.RemoveAll(dir) }, patchCGScriptFile(path)
}

func downloadCGScript(url, path string) error {
	cmd := exec.Command("curl", "-fsSL", url, "-o", path)
	return cmd.Run()
}

func patchCGScriptFile(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	patched := patchCGArithmeticIncrements(string(body))
	return os.WriteFile(path, []byte(patched), 0o700)
}

func patchCGArithmeticIncrements(script string) string {
	return cgPostfixIncrementPattern.ReplaceAllString(script, `((${1}+=1))`)
}

// runCGInstaller wires stdio + working dir onto the prepared command and
// executes it. Failures are logged to opts.Stderr in the standardized
// format before being wrapped with %w so callers can errors.Is / unwrap.
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
