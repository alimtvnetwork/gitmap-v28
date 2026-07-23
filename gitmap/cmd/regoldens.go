package cmd

// CLI entry point for `gitmap regoldens`. Wraps the two-pass
// golden-fixture regeneration workflow defined in
// spec/05-coding-guidelines/21-golden-fixture-regeneration.md so
// contributors cannot forget the verify pass or leak the gate env
// vars into their shell. The two-key safety gate values come from
// the goldenguard package (single source of truth).

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/goldenguard"
)

// regoldensFlags captures parsed CLI inputs.
type regoldensFlags struct {
	pattern     string
	pkg         string
	skipVerify  bool
	isDryRun    bool
	diffMode    string // "" (off), "short", or "full"
	determinism bool
}

// hasDiff reports whether any diff summary was requested.
func (c regoldensFlags) hasDiff() bool { return c.diffMode != "" }

// goTestUpdateEnvValue mirrors goldenguard.allowUpdateValue (which
// is unexported). Both gate env vars must equal "1" to unlock pass 1.
const goTestUpdateEnvValue = "1"

// goTestUpdateTriggerEnv is the per-test trigger env var checked
// by every golden test in the repo.
const goTestUpdateTriggerEnv = "GITMAP_UPDATE_GOLDEN"

// runRegoldens is the dispatcher entry. checkHelp first so `--help`
// works even if other flags would fail to parse.
func runRegoldens(args []string) {
	checkHelp("regoldens", args)
	cfg := parseRegoldensFlags(args)
	if cfg.pattern == "" {
		fmt.Fprintln(os.Stderr, constants.ErrRegoldensMissingPat)
		os.Exit(2)
	}
	validateDiffMode(cfg.diffMode)
	if cfg.isDryRun {
		emitRegoldensDryRun(cfg)
		return
	}
	executeRegoldens(cfg)
}

// validateDiffMode lives in regoldens_validate.go (file-length cap).

// parseRegoldensFlags wires the flag set. Defaults match the
// constants block so changing a default is a one-line edit there.
func parseRegoldensFlags(args []string) regoldensFlags {
	fs := flag.NewFlagSet("regoldens", flag.ExitOnError)
	cfg := regoldensFlags{pkg: constants.RegoldensDefaultPackageGlob}
	bindRegoldensFlags(fs, &cfg)
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "regoldens: parse flags: %v\n", err)
		os.Exit(2)
	}

	return cfg
}

// bindRegoldensFlags registers every CLI flag against the parser.
// Split out so parseRegoldensFlags stays inside the function-length
// budget — the binding step is mechanical and rarely changes.
func bindRegoldensFlags(fs *flag.FlagSet, cfg *regoldensFlags) {
	fs.StringVar(&cfg.pattern, constants.FlagRegoldensPattern, "",
		constants.FlagDescRegoldensPattern)
	fs.StringVar(&cfg.pkg, constants.FlagRegoldensPackage,
		constants.RegoldensDefaultPackageGlob,
		constants.FlagDescRegoldensPackage)
	fs.BoolVar(&cfg.skipVerify, constants.FlagRegoldensSkipVerify, false,
		constants.FlagDescRegoldensSkipVerify)
	fs.BoolVar(&cfg.isDryRun, constants.FlagRegoldensDryRun, false,
		constants.FlagDescRegoldensDryRun)
	fs.StringVar(&cfg.diffMode, constants.FlagRegoldensDiff, "",
		constants.FlagDescRegoldensDiff)
	fs.BoolVar(&cfg.determinism, constants.FlagRegoldensDeterminism, false,
		constants.FlagDescRegoldensDeterminism)
}

// emitRegoldensDryRun lives in regoldens_dryrun.go (file-length cap).

// goTestArgv returns the `go test ...` argv shared by both passes.
// `-count=1` defeats the test cache so pass 2 actually re-runs.
func goTestArgv(cfg regoldensFlags) []string {
	return []string{"go", "test", cfg.pkg, "-run", cfg.pattern, "-count=1"}
}

// executeRegoldens lives in regoldens_exec.go to keep this file
// under the 200-line cap. The split is purely organizational.

// runRegoldensPass prints the header, runs one pass, and exits 1
// with the supplied error format on failure. withGate toggles the
// gate-vars-injected (pass 1) vs gate-vars-stripped (pass 2) env.
func runRegoldensPass(cfg regoldensFlags, withGate bool, header, errFmt string) {
	fmt.Fprint(os.Stderr, header)
	code := runGoTestPass(cfg, withGate)
	if code == 0 {
		return
	}
	fmt.Fprintf(os.Stderr, errFmt, code)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

// runGoTestPass executes one `go test` invocation and returns its
// exit code. When withGate is false, the two gate env vars are
// explicitly REMOVED from the child env (not just left unset) so a
// developer's leaked shell export cannot silently break pass 2.
func runGoTestPass(cfg regoldensFlags, withGate bool) int {
	argv := goTestArgv(cfg)
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = buildPassEnv(withGate)
	if err := cmd.Run(); err != nil {
		return extractExitCode(err)
	}

	return 0
}

// buildPassEnv returns the child environment for one pass. Gate
// vars are stripped from the parent env first, then re-added only
// when withGate is true.
func buildPassEnv(withGate bool) []string {
	out := stripGoldenGateVars(os.Environ())
	if !withGate {

		return out
	}

	return append(out,
		goTestUpdateTriggerEnv+"="+goTestUpdateEnvValue,
		goldenguard.AllowUpdateEnv+"="+goTestUpdateEnvValue,
	)
}

// stripGoldenGateVars filters parent env entries, dropping any
// gate-related KEY=value pair so a leaked shell export cannot
// influence the child process.
func stripGoldenGateVars(parent []string) []string {
	out := make([]string, 0, len(parent)+2)
	for _, kv := range parent {
		if isGoldenGateVar(kv) {
			continue
		}
		out = append(out, kv)
	}

	return out
}

// isGoldenGateVar reports whether kv is one of the two gate env
// vars in `KEY=value` form. Used to filter the parent environment.
func isGoldenGateVar(kv string) bool {
	return strings.HasPrefix(kv, goTestUpdateTriggerEnv+"=") ||
		strings.HasPrefix(kv, goldenguard.AllowUpdateEnv+"=")
}

// extractExitCode pulls the numeric exit code from an *exec.ExitError.
// Non-ExitError failures (e.g. `go` binary not on PATH) map to 127,
// matching POSIX shell convention for "command not found".
func extractExitCode(err error) int {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	fmt.Fprintf(os.Stderr, "regoldens: failed to invoke `go test`: %v\n", err)

	return 127
}
