package cmd

// CLI runners for `gitmap startup-list` and `gitmap startup-remove`.
// Both commands work on Linux/Unix (XDG `.desktop` files) and macOS
// (LaunchAgent `.plist` files); on Windows they exit with a clear
// "unsupported OS" message rather than silently doing nothing
// (silent no-ops on a different OS would be a UX trap: the user
// would assume their startup entries are managed when they're not).
//
// Output contract (per user requirement):
//   - Clear: every outcome prints exactly one line summary.
//   - Safe no-op: missing files, third-party files, and empty
//     listings all exit 0 with a message — never an error.
//   - Scoped: only entries carrying the gitmap marker are touched.
//     The startup package enforces this; this layer just renders the
//     result. `--dry-run` runs the same classification but skips
//     the actual unlink.

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/startup"
)

// runStartupList enumerates gitmap-managed XDG autostart entries and
// renders them in the chosen format. The default `table` format
// matches the legacy human-readable output; `json` and `csv` exist
// for piping into other tools (jq, spreadsheet imports, etc).
// --backend and --name filter the result set in-memory before the
// renderer runs so format-specific code paths stay untouched.
func runStartupList(args []string) {
	opts, err := parseStartupListFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	entries, err := startup.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	entries = filterStartupList(entries, opts.backend, opts.name)
	dir, _ := startup.AutostartDir()
	if err := renderStartupList(opts.format, opts.jsonIndent, dir, entries); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// startupListOpts bundles the parsed flag values so the dispatcher
// stays a one-liner and adding a future filter (e.g., --since) is
// one struct field + one parser line, never a signature change.
type startupListOpts struct {
	format     string
	jsonIndent int
	backend    string
	name       string
}

// parseStartupListFlags extracts --format, --json-indent, --backend,
// and --name and validates each independently. Unknown values fail
// fast with exit 2 so scripts catch typos immediately rather than
// getting silent fall-through to a default rendering. --json-indent
// is parsed and validated even when the format ignores it, so a
// typo like `--json-indent=99` is caught regardless of the chosen
// format. --backend rejects unknown values with the same shape as
// startup-add so the user sees one consistent message across both
// commands.
func parseStartupListFlags(args []string) (startupListOpts, error) {
	fs := flag.NewFlagSet("startup-list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	format := fs.String(
		constants.FlagStartupListFormat,
		constants.StartupListFormatTable,
		constants.FlagDescStartupListFormat,
	)
	jsonIndent := fs.Int(
		constants.FlagStartupListJSONIndent,
		constants.StartupListJSONIndentDefault,
		constants.FlagDescStartupListJSONIndent,
	)
	backend := fs.String(
		constants.FlagStartupListBackend, "",
		constants.FlagDescStartupListBackend,
	)
	name := fs.String(
		constants.FlagStartupListName, "",
		constants.FlagDescStartupListName,
	)
	if err := fs.Parse(args); err != nil {

		return startupListOpts{}, err
	}
	if *jsonIndent < 0 || *jsonIndent > constants.StartupListJSONIndentMax {

		return startupListOpts{}, fmt.Errorf(constants.ErrStartupListBadJSONIndent, *jsonIndent)
	}
	if err := validateStartupListBackend(*backend); err != nil {

		return startupListOpts{}, err
	}
	switch *format {
	case constants.StartupListFormatTable, constants.OutputTerminal,
		constants.OutputJSON, constants.StartupListFormatJSONL,
		constants.OutputCSV:

		return startupListOpts{
			format: *format, jsonIndent: *jsonIndent,
			backend: *backend, name: *name,
		}, nil
	default:

		return startupListOpts{}, fmt.Errorf(constants.ErrStartupListBadFormat, *format)
	}
}

// validateStartupListBackend rejects unknown --backend values with
// the same error message shape as startup-add. Empty (== no filter)
// is always valid.
func validateStartupListBackend(b string) error {
	switch b {
	case "",
		constants.StartupBackendRegistry,
		constants.StartupBackendRegistryHKLM,
		constants.StartupBackendStartupFolder:

		return nil
	}

	return fmt.Errorf(constants.ErrStartupListBadBackend, b)
}

// runStartupRemove deletes a single managed entry. After --dry-run,
// --backend, --output and --json-indent are parsed off the args,
// exactly one positional name must remain; missing or extra
// positionals trigger the usage error. All four RemoveStatus
// outcomes map to a distinct user-visible message in --output=
// terminal mode (with a `(dry-run)` mirror set when --dry-run is
// active) so the CLI is unambiguous about what happened — or what
// would happen. --output=json emits the shared startupStatus
// object instead, with the same field shape as startup-add for
// uniform downstream parsing.
func runStartupRemove(args []string) {
	cfg, err := parseStartupRemoveFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, constants.ErrStartupRemoveUsage)
		os.Exit(2)
	}
	if err := validateStartupOutput(constants.CmdStartupRemove, cfg.output, cfg.jsonIndent); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	backend, err := startup.ParseBackend(cfg.backend)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	res, err := startup.RemoveWithOptions(cfg.name, startup.RemoveOptions{
		DryRun: cfg.dryRun, Backend: backend,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if cfg.output == constants.OutputJSON {
		_ = emitStartupStatus(cfg.output, cfg.jsonIndent,
			removeResultToStatus(cfg.name, res))

		return
	}
	printRemoveResult(cfg.name, res)
}

// startupRemoveFlags bundles parsed flag values so the dispatcher
// stays small and adding a future flag is one struct field + one
// parser line, never a return-shape change.
type startupRemoveFlags struct {
	name       string
	backend    string
	output     string
	jsonIndent int
	dryRun     bool
}

// parseStartupRemoveFlags pulls --dry-run, --backend, --output, and
// --json-indent off the args and returns the remaining single
// positional name in the struct. Returns an error when the
// positional count is wrong so the caller exits 2 with the usage
// message — matching the pre-flag behavior for malformed
// invocations.
func parseStartupRemoveFlags(args []string) (startupRemoveFlags, error) {
	fs := flag.NewFlagSet(constants.CmdStartupRemove, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var cfg startupRemoveFlags
	fs.BoolVar(&cfg.dryRun, constants.FlagStartupRemoveDryRun, false,
		constants.FlagDescStartupRemoveDryRun)
	fs.StringVar(&cfg.backend, constants.FlagStartupRemoveBackend, "",
		constants.FlagDescStartupRemoveBackend)
	fs.StringVar(&cfg.output, constants.FlagStartupOutput, constants.OutputTerminal,
		constants.FlagDescStartupOutput)
	fs.IntVar(&cfg.jsonIndent, constants.FlagStartupJSONIndent,
		constants.StartupListJSONIndentDefault, constants.FlagDescStartupJSONIndent)
	if err := fs.Parse(args); err != nil {

		return startupRemoveFlags{}, err
	}
	rest := fs.Args()
	if len(rest) != 1 {

		return startupRemoveFlags{}, fmt.Errorf("expected 1 positional name, got %d", len(rest))
	}
	cfg.name = rest[0]

	return cfg, nil
}

// validateStartupOutput rejects unknown --output values and
// out-of-range --json-indent. Lives here (not next to per-command
// flag parsers) so both startup-add and startup-remove get the
// same error wording without duplicating the switch.
func validateStartupOutput(cmd, output string, jsonIndent int) error {
	switch output {
	case constants.OutputTerminal, constants.OutputJSON:
	default:

		return fmt.Errorf(constants.ErrStartupBadOutput, cmd, output)
	}
	if jsonIndent < 0 || jsonIndent > constants.StartupListJSONIndentMax {

		return fmt.Errorf(constants.ErrStartupListBadJSONIndent, jsonIndent)
	}

	return nil
}

// printRemoveResult routes one of four messages depending on the
// status. Dry-run results pick the `(dry-run)` mirror message so a
// preview is visually distinct from a real action — important when
// users pipe both into the same log. Each branch is single-line so
// log-scrapers can classify outcomes without parsing multi-line
// output.
func printRemoveResult(name string, res startup.RemoveResult) {
	if res.DryRun {
		printRemoveResultDryRun(name, res)

		return
	}
	switch res.Status {
	case startup.RemoveDeleted:
		fmt.Printf(constants.MsgStartupRemoveOK, res.Path)
	case startup.RemoveNoOp:
		fmt.Printf(constants.MsgStartupRemoveNoOp, name)
	case startup.RemoveRefused:
		fmt.Printf(constants.MsgStartupRemoveNotOurs, res.Path)
	case startup.RemoveBadName:
		fmt.Printf(constants.MsgStartupRemoveBadName, name)
	}
}

// printRemoveResultDryRun is the `--dry-run` parallel of
// printRemoveResult. Kept as its own function (rather than an
// if/else inside the parent switch) so the live and preview message
// tables are easy to diff line-for-line during code review.
func printRemoveResultDryRun(name string, res startup.RemoveResult) {
	switch res.Status {
	case startup.RemoveDeleted:
		fmt.Printf(constants.MsgStartupRemoveDryOK, res.Path)
	case startup.RemoveNoOp:
		fmt.Printf(constants.MsgStartupRemoveDryNoOp, name)
	case startup.RemoveRefused:
		fmt.Printf(constants.MsgStartupRemoveDryNotOurs, res.Path)
	case startup.RemoveBadName:
		fmt.Printf(constants.MsgStartupRemoveDryBadName, name)
	}
}
