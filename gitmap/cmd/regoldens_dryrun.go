package cmd

// Dry-run rendering for `gitmap regoldens`. Split out so regoldens.go
// stays under the 200-line file cap and each helper stays under the
// 15-line function cap.

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/goldenguard"
)

// emitRegoldensDryRun prints every invocation that would run, in
// the order it would run, without executing anything. Pre-check
// (when --determinism is set) is printed BEFORE pass 1 because that
// is the actual execution order.
func emitRegoldensDryRun(cfg regoldensFlags) {
	if cfg.determinism {
		emitDryRunPrecheck(cfg)
	}
	emitDryRunMainPasses(cfg)
	emitDryRunDiffNote(cfg)
}

// emitDryRunPrecheck renders the trigger-only command used by the
// determinism pre-check pass: trigger ON, allow-update gate OFF.
func emitDryRunPrecheck(cfg regoldensFlags) {
	cmd := strings.Join(append(
		[]string{goTestUpdateTriggerEnv + "=" + goTestUpdateEnvValue},
		goTestArgv(cfg)...,
	), " ")
	fmt.Fprintf(os.Stdout, "▸ Pre-check (would run first):\n  %s\n", cmd)
}

// emitDryRunMainPasses renders the standard pass-1 (gated) and
// pass-2 (gates stripped) commands using the existing template.
func emitDryRunMainPasses(cfg regoldensFlags) {
	pass1 := strings.Join(append(
		[]string{
			goTestUpdateTriggerEnv + "=" + goTestUpdateEnvValue,
			goldenguard.AllowUpdateEnv + "=" + goTestUpdateEnvValue,
		},
		goTestArgv(cfg)...,
	), " ")
	pass2 := strings.Join(goTestArgv(cfg), " ")
	fmt.Fprintf(os.Stdout, constants.MsgRegoldensDryRun, pass1, pass2)
}

// emitDryRunDiffNote appends a one-liner reminding the user that
// --diff would inject a summary between pass 1 and pass 2.
func emitDryRunDiffNote(cfg regoldensFlags) {
	if !cfg.hasDiff() {
		return
	}
	fmt.Fprintf(os.Stdout,
		"  (--diff=%s: golden diff summary would print between passes)\n",
		cfg.diffMode)
}
