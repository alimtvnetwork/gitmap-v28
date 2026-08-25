// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/glyphs"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/theme"
	"github.com/spf13/cobra"
)

// Run is the main entry point for the CLI.
func Run() {
	initConsole()

	if len(os.Args) < 2 {
		PrintBinaryLocations()
		printUsage()
		os.Exit(1)
	}

	// Strip the global `--theme` palette selector first so it is
	// honored even when no subcommand-specific args are present.
	// theme.Install must run AFTER the env var is set but BEFORE
	// any subcommand writes colored output.
	os.Args = append(os.Args[:1], stripThemeFlag(os.Args[1:])...)
	theme.Install()

	// Strip the global `--glyphs` switch (rich | safe | auto) and
	// install the glyph filter. Runs AFTER theme so the safe-mode
	// ASCII rewrites apply to bytes already past theme's SGR rewrite.
	os.Args = append(os.Args[:1], stripGlyphsFlag(os.Args[1:])...)
	glyphs.Install()

	// Register pipe drainers so cliexit.Fail flushes them before
	// os.Exit. Without this, a failure message written to a
	// theme- or glyphs-wrapped os.Stderr just before os.Exit can be
	// lost on Windows (the forwarder goroutine never gets scheduled
	// to copy bytes from the pipe buffer to the inherited fd).
	cliexit.RegisterFlusher(theme.Drain)
	cliexit.RegisterFlusher(glyphs.Drain)

	// Strip the global `--vscode-sync-disabled` kill switch from argv
	// (and flip the env var) before any subcommand sees its flagset.
	// Done first so even URL-shortcut and alias rewrites operate on
	// already-cleaned args.
	os.Args = append(os.Args[:1], stripVSCodeSyncDisabledFlag(os.Args[1:])...)
	// Strip the global tag-customization flags too, persisting their
	// values into GITMAP_VSCODE_TAG_{ADD,SKIP,MARKER} so every
	// DetectTagsCustom caller — present and future — picks them up.
	os.Args = append(os.Args[:1], stripVSCodeTagFlags(os.Args[1:])...)
	if len(os.Args) < 2 {
		PrintBinaryLocations()
		printUsage()
		os.Exit(1)
	}

	// Skip migration for commands that must produce clean stdout
	cmd := os.Args[1]
	if cmd != constants.CmdVersion && cmd != constants.CmdVersionAlias {
		migrateLegacyDirs()
	}

	// URL shortcut: `gitmap <git-url> [<url2> ...]` (and variants with
	// leading flags like `gitmap --verbose <url>`) is rewritten to
	// `gitmap clone <args...>` so users don't have to remember the
	// subcommand for the most common operation. We trigger when any
	// positional arg looks like an HTTPS / SSH git URL or a comma-list
	// containing one — covering all the forms users actually type:
	//
	//   gitmap https://...                       → gitmap clone https://...
	//   gitmap https://a,https://b,https://c     → gitmap clone https://a,https://b,https://c
	//   gitmap https://a, https://b https://c    → gitmap clone https://a, https://b https://c
	//   gitmap --verbose https://...             → gitmap clone --verbose https://...
	if shouldRewriteToClone(os.Args[1:]) {
		os.Args = append([]string{os.Args[0], constants.CmdClone}, os.Args[1:]...)
	}

	aliasName, cleaned := extractAliasFlag(os.Args[2:])
	if len(aliasName) > 0 {
		resolveAliasContext(aliasName)
		os.Args = append(os.Args[:2], cleaned...)
	}

	runDispatch(os.Args[1])
}

// runDispatch is the single entry point every CLI invocation flows
// through so the theme/glyphs pipe drainers are guaranteed to run
// before the process returns. Centralizing the deferred Drain calls
// here prevents future entry points from re-introducing the Windows
// "last stdout line is lost" bug (the version-mismatch smoke failure
// fixed in v6.74.0). Any new dispatch surface MUST call runDispatch
// rather than dispatch directly.
func runDispatch(command string) {
	// Order matters: glyphs wraps stdout AFTER theme, so drain
	// glyphs first (outermost writer) then theme (inner writer).
	// Defer order runs LIFO, so declaring theme first + glyphs
	// second yields the correct outer→inner drain sequence.
	defer theme.Drain()
	defer glyphs.Drain()
	dispatch(command)
}

// dispatch routes to the correct subcommand handler with audit tracking.
func dispatch(command string) {
	auditID, auditStart, shouldAudit := beginCommandAudit(command, os.Args[2:])

	if dispatchUser(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchCore(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchRelease(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchUtility(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchData(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchTooling(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchProjectRepos(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchDiff(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchMoveMerge(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchAdd(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchSync(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchCommons(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchTemplates(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}
	if dispatchCommitTransfer(command) {
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

		return
	}

	if command == "ip" || command == "ip-change" {
		if err := dispatchIP(context.Background(), os.Args[1:], nil); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)
		return
	}

	if command == "sj" || command == "ssh-join" || command == "ssh-joined" || command == "ssh-joiner" {
		if err := dispatchSJ(context.Background(), os.Args[1:], nil); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)
		return
	}

	if looksLikeURLToken(command) {
		fmt.Fprintf(os.Stderr, constants.ErrUnknownCommandURLHint, command)
	} else {
		fmt.Fprintf(os.Stderr, constants.ErrUnknownCommand, command)
	}
	printUsage()
	os.Exit(1)
}

// shouldRewriteToClone returns true when the args (excluding argv[0])
// describe a bare-URL invocation that should be redirected to `clone`.
// It accepts URLs in any positional slot — not just the first — so
// invocations with leading flags (e.g. `gitmap --verbose <url>`) and
// PowerShell's silent comma-splitting both work.
func shouldRewriteToClone(args []string) bool {
	if len(args) == 0 {
		return false
	}
	// Never rewrite if the first token is already a known subcommand.
	if !looksLikeFlag(args[0]) && !looksLikeURLToken(args[0]) {
		return false
	}
	for _, a := range args {
		if looksLikeURLToken(a) {
			return true
		}
	}
	return false
}

// looksLikeFlag reports whether the token starts with "-" or "--".
func looksLikeFlag(s string) bool {
	return len(s) > 1 && s[0] == '-'
}

// looksLikeURLToken reports whether a token (or any comma-split piece
// of it) is shaped like a git URL. Used by both the shortcut and the
// unknown-command hint so they agree on what counts as a URL.
func looksLikeURLToken(s string) bool {
	for _, part := range splitOnComma(s) {
		if isLikelyURL(part) {
			return true
		}
	}
	return false
}

// splitOnComma is a tiny strings.Split wrapper so root.go doesn't
// need its own strings import for this single use; trims whitespace
// around each piece and drops empties.
func splitOnComma(s string) []string {
	out := make([]string, 0, 4)
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			out = appendTrimmedPiece(out, s, start, i)
			start = i + 1
		}
	}
	return out
}

func appendTrimmedPiece(out []string, s string, start, i int) []string {
	piece := trimSpaces(s[start:i])
	if len(piece) > 0 {
		return append(out, piece)
	}
	return out
}

// trimSpaces removes ASCII whitespace from both ends without pulling
// the strings package into root.go's import surface.
func trimSpaces(s string) string {
	i, j := 0, len(s)
	for i < j && isSpace(s[i]) {
		i++
	}
	for j > i && isSpace(s[j-1]) {
		j--
	}
	return s[i:j]
}

// isSpace reports whether b is ASCII whitespace.
func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func dispatchIP(ctx context.Context, args []string, parent *cobra.Command) error {
	if len(args) == 0 {
		return nil
	}
	switch args[0] {
	case "ip":
		// IPCmd expects to parse args itself, so let Cobra do its thing.
		IPCmd.SetArgs(args[1:])
		return IPCmd.ExecuteContext(ctx)
	case "ip-change":
		IPChangeCmd.SetArgs(args[1:])
		return IPChangeCmd.ExecuteContext(ctx)
	}
	return nil
}

func runSJ(args []string) {
	if err := dispatchSJ(context.Background(), os.Args[1:], nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func dispatchSJ(ctx context.Context, args []string, root *cobra.Command) error {
	if len(args) == 0 {
		return SSHJoinCmd.ExecuteContext(ctx)
	}

	cleanArgs := stripLeadingSJCommand(args)
	if len(cleanArgs) == 0 {
		SSHJoinCmd.SetArgs([]string{})
		return SSHJoinCmd.ExecuteContext(ctx)
	}

	return routeSJSubcommand(ctx, cleanArgs, root)
}

func stripLeadingSJCommand(args []string) []string {
	first := args[0]
	if first == "sj" || first == "ssh-join" || first == "ssh-joined" || first == "ssh-joiner" {
		return args[1:]
	}

	return args
}

func routeSJSubcommand(ctx context.Context, args []string, root *cobra.Command) error {
	subCmd := args[0]
	switch subCmd {
	case "ls", "list":
		return runSJLs(root, args[1:], ctx)
	case "rm", "remove":
		return runSJRm(root, args[1:], ctx)
	case "history", "hist":
		return runSJHistory(root, args[1:], ctx)
	case "add-auth", "auth":
		return runSJAddAuth(root, args[1:], ctx)
	default:
		SSHJoinCmd.SetArgs(args)
		return SSHJoinCmd.ExecuteContext(ctx)
	}
}

