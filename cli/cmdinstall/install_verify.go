// Package cmdinstall — install_verify.go: `gitmap install verify`.
//
// Post-install self-check: binary version, PATH registration, symlink
// integrity, deploy contamination, and install source. Read-only.
//
// Wiring (lead): inside runInstall's subcommand chain, before flag parsing:
//  1. match with isInstallVerifyCommand(args)
//  2. call RunInstallVerify(args[1:]) with an explicit nil check —
//     *appfault.AppError is a typed nil when there is no error, so a bare
//     `return RunInstallVerify(...)` would convert it to a non-nil error.
package cmdinstall

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)

// isInstallVerifyCommand reports whether args invoke `install verify`.
func isInstallVerifyCommand(args []string) bool {
	return len(args) > 0 && strutil.EqualFoldAnyTrim(args[0], "verify")
}

// verifyCheck is one checklist row.
type verifyCheck struct {
	name     string
	isPassed bool
	detail   string
}

// RunInstallVerify executes `gitmap install verify` — read-only checklist.
// args is the tail after the `verify` token.
func RunInstallVerify(args []string) *appfault.AppError {
	if isInstallHelpRequested(args) {
		fmt.Print(installVerifyHelpText)
		cliexit.Exit(0)
	}

	checks := runInstallVerifyChecks()
	failed := 0
	for _, check := range checks {
		mark := "✓"
		if !check.isPassed {
			mark = "✗"
			failed++
		}
		if check.detail != "" {
			fmt.Printf("%s %s (%s)\n", mark, check.name, check.detail)
		} else {
			fmt.Printf("%s %s\n", mark, check.name)
		}
	}

	if failed > 0 {
		return appfault.NewValidationError(
			fmt.Sprintf("install verify: %d of %d checks failed", failed, len(checks)))
	}

	return nil
}

// runInstallVerifyChecks builds the checklist rows.
func runInstallVerifyChecks() []verifyCheck {
	exe, exeErr := os.Executable()

	checks := []verifyCheck{
		checkBinaryVersion(exe, exeErr),
		checkPathRegistered(exe, exeErr),
		checkSymlinksIntact(exe, exeErr),
		checkNoDeployContamination(exe, exeErr),
		checkInstallSource(exe, exeErr),
	}

	return checks
}

// checkBinaryVersion verifies the running binary exists and reports its version.
func checkBinaryVersion(exe string, exeErr error) verifyCheck {
	if exeErr != nil || exe == "" {
		return verifyCheck{name: "binary version", isPassed: false, detail: "os.Executable failed"}
	}
	if _, statErr := os.Stat(exe); statErr != nil {
		return verifyCheck{name: "binary version", isPassed: false, detail: exe}
	}

	return verifyCheck{name: "binary version", isPassed: true, detail: "v" + constants.Version}
}

// checkPathRegistered verifies the binary's directory is on PATH.
func checkPathRegistered(exe string, exeErr error) verifyCheck {
	if exeErr != nil || exe == "" {
		return verifyCheck{name: "PATH registered", isPassed: false, detail: "unknown binary path"}
	}

	dir := filepath.Clean(filepath.Dir(exe))
	for _, entry := range filepath.SplitList(os.Getenv("PATH")) {
		cleaned := filepath.Clean(entry)
		if runtime.GOOS == "windows" {
			if strings.EqualFold(cleaned, dir) {
				return verifyCheck{name: "PATH registered", isPassed: true, detail: dir}
			}
		} else if cleaned == dir {
			return verifyCheck{name: "PATH registered", isPassed: true, detail: dir}
		}
	}

	return verifyCheck{name: "PATH registered", isPassed: false, detail: dir + " not on PATH"}
}

// checkSymlinksIntact verifies the executable path resolves cleanly and any
// symlink in the chain points at an existing target.
func checkSymlinksIntact(exe string, exeErr error) verifyCheck {
	if exeErr != nil || exe == "" {
		return verifyCheck{name: "symlinks intact", isPassed: false, detail: "unknown binary path"}
	}

	resolved, evalErr := filepath.EvalSymlinks(exe)
	if evalErr != nil {
		return verifyCheck{name: "symlinks intact", isPassed: false, detail: evalErr.Error()}
	}
	if _, statErr := os.Stat(resolved); statErr != nil {
		return verifyCheck{name: "symlinks intact", isPassed: false, detail: "target missing: " + resolved}
	}
	if resolved != filepath.Clean(exe) {
		return verifyCheck{name: "symlinks intact", isPassed: true, detail: "symlink → " + resolved}
	}

	return verifyCheck{name: "symlinks intact", isPassed: true, detail: "no symlinks in chain"}
}

// checkNoDeployContamination flags deploy staging artifacts (run.sh) left
// next to the installed binary — the tar-extract-dir contamination class.
func checkNoDeployContamination(exe string, exeErr error) verifyCheck {
	if exeErr != nil || exe == "" {
		return verifyCheck{name: "no deploy contamination", isPassed: false, detail: "unknown binary path"}
	}

	dir := filepath.Dir(exe)
	for _, marker := range []string{"run.sh", "deploy", "staging"} {
		if _, statErr := os.Stat(filepath.Join(dir, marker)); statErr == nil {
			return verifyCheck{name: "no deploy contamination", isPassed: false, detail: marker + " next to binary"}
		}
	}

	return verifyCheck{name: "no deploy contamination", isPassed: true, detail: ""}
}

// checkInstallSource classifies the install origin by binary location.
// Informational: always passes, the detail carries the classification.
func checkInstallSource(exe string, exeErr error) verifyCheck {
	if exeErr != nil || exe == "" {
		return verifyCheck{name: "install source", isPassed: true, detail: "unknown"}
	}

	lower := strings.ToLower(exe)
	home, _ := os.UserHomeDir()
	switch {
	case strings.Contains(lower, "gitmap-v28") || hasNearbyGitDir(exe):
		return verifyCheck{name: "install source", isPassed: true, detail: "local build"}
	case home != "" && strings.HasPrefix(lower, strings.ToLower(filepath.Join(home, ".local", "bin"))):
		return verifyCheck{name: "install source", isPassed: true, detail: "github release"}
	case strings.HasPrefix(lower, "/usr/local/bin") || strings.HasPrefix(lower, "/opt/"):
		return verifyCheck{name: "install source", isPassed: true, detail: "github release"}
	default:
		return verifyCheck{name: "install source", isPassed: true, detail: "other"}
	}
}

// hasNearbyGitDir reports whether a .git directory sits near the binary
// (source-checkout heuristic for the install-source classification).
func hasNearbyGitDir(exe string) bool {
	dir := filepath.Dir(exe)
	for i := 0; i < 4; i++ {
		if _, statErr := os.Stat(filepath.Join(dir, ".git")); statErr == nil {
			return true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return false
		}
		dir = parent
	}

	return false
}
