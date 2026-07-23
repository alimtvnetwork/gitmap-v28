package cmd

// CLI runner for `gitmap startup-add`. Supports Linux/Unix (XDG
// .desktop) and macOS (LaunchAgent .plist); Windows is rejected up
// front with the standard "unsupported OS" message rather than
// letting the startup package's directory resolver do it (gives a
// single consistent error path and matches startup-list /
// startup-remove on the same input).
//
// Output contract (matches startup-list / startup-remove):
//   - Exactly one summary line per outcome.
//   - "Soft" outcomes (already-exists, refused, bad-name) exit 0
//     so users can rerun the same `startup-add` invocation
//     idempotently from a provisioning script.
//   - Real I/O errors and missing required flags exit non-zero.

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/startup"
)

// startupAddFlags bundles the parsed flag values so the orchestrator
// stays small and the resolver helpers can be unit-tested without
// poking at flag.FlagSet plumbing.
type startupAddFlags struct {
	name        string
	exec        string
	displayName string
	comment     string
	workingDir  string
	noDisplay   bool
	force       bool
	backend     string
	output      string
	jsonIndent  int
}

// runStartupAdd is the dispatcher entry point. Windows is now
// supported via the Registry / Startup-folder backends — no early
// short-circuit needed here. The startup package's runtime guards
// surface ErrStartupUnsupportedOS for any remaining unsupported OS.
func runStartupAdd(args []string) {
	checkHelp("startup-add", args)
	cfg := parseStartupAddFlags(args)
	if err := validateStartupOutput(constants.CmdStartupAdd, cfg.output, cfg.jsonIndent); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	exec, ok := resolveStartupAddExec(cfg.exec)
	if !ok {
		fmt.Fprintln(os.Stderr, constants.ErrStartupAddMissingExec)
		os.Exit(2)
	}
	backend, err := startup.ParseBackend(cfg.backend)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	res, err := startup.Add(startup.AddOptions{
		Name: cfg.name, Exec: exec,
		DisplayName: cfg.displayName, Comment: cfg.comment,
		NoDisplay: cfg.noDisplay, Force: cfg.force,
		WorkingDir: cfg.workingDir,
		Backend:    backend,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if cfg.output == constants.OutputJSON {
		_ = emitStartupStatus(cfg.output, cfg.jsonIndent,
			addResultToStatus(cfg.name, cfg.force, res))

		return
	}
	printAddResult(cfg.name, res)
}

// parseStartupAddFlags wires the 6 CLI flags into a startupAddFlags.
// Defaults are deliberately empty strings (not "main", not the
// binary path) so resolveStartupAddExec can distinguish "user
// omitted the flag" from "user passed an empty value".
func parseStartupAddFlags(args []string) startupAddFlags {
	fs := flag.NewFlagSet(constants.CmdStartupAdd, flag.ExitOnError)
	var cfg startupAddFlags
	fs.StringVar(&cfg.name, constants.FlagStartupAddName, "",
		constants.FlagDescStartupAddName)
	fs.StringVar(&cfg.exec, constants.FlagStartupAddExec, "",
		constants.FlagDescStartupAddExec)
	fs.StringVar(&cfg.displayName, constants.FlagStartupAddDisplay, "",
		constants.FlagDescStartupAddDisplay)
	fs.StringVar(&cfg.comment, constants.FlagStartupAddComment, "",
		constants.FlagDescStartupAddComment)
	fs.BoolVar(&cfg.noDisplay, constants.FlagStartupAddNoDisplay, false,
		constants.FlagDescStartupAddNoDisplay)
	fs.BoolVar(&cfg.force, constants.FlagStartupAddForce, false,
		constants.FlagDescStartupAddForce)
	fs.StringVar(&cfg.workingDir, constants.FlagStartupAddWorkingDir, "",
		constants.FlagDescStartupAddWorkingDir)
	fs.StringVar(&cfg.backend, constants.FlagStartupAddBackend, "",
		constants.FlagDescStartupAddBackend)
	fs.StringVar(&cfg.output, constants.FlagStartupOutput, constants.OutputTerminal,
		constants.FlagDescStartupOutput)
	fs.IntVar(&cfg.jsonIndent, constants.FlagStartupJSONIndent,
		constants.StartupListJSONIndentDefault, constants.FlagDescStartupJSONIndent)
	fs.Parse(args)
	if cfg.name == "" {
		fmt.Fprintln(os.Stderr,
			"startup-add: --name is required (e.g. --name myapp)")
		os.Exit(2)
	}

	return cfg
}

// resolveStartupAddExec falls back to os.Executable() when --exec
// was omitted. Returns ok=false only when both inputs fail — the
// caller surfaces that as a clear "must pass --exec" error rather
// than letting startup.Add write a broken Exec= line.
func resolveStartupAddExec(flagValue string) (string, bool) {
	if flagValue != "" {

		return flagValue, true
	}
	bin, err := os.Executable()
	if err != nil || bin == "" {

		return "", false
	}

	return bin, true
}

// printAddResult routes the five AddStatus outcomes to one-line
// messages. Mirrors printRemoveResult in startup.go for symmetry.
func printAddResult(name string, res startup.AddResult) {
	switch res.Status {
	case startup.AddCreated:
		fmt.Printf(constants.MsgStartupAddCreated, res.Path)
	case startup.AddOverwritten:
		fmt.Printf(constants.MsgStartupAddOverwritten, res.Path)
	case startup.AddExists:
		fmt.Printf(constants.MsgStartupAddExists, res.Path)
	case startup.AddRefused:
		fmt.Printf(constants.MsgStartupAddRefused, res.Path)
	case startup.AddBadName:
		fmt.Printf(constants.MsgStartupAddBadName, name)
	}
}
