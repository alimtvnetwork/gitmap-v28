package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/render"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/templates"
)

const (
	cmdTemplatesDiff      = constants.CmdTemplatesDiff
	cmdTemplatesDiffAlias = constants.CmdTemplatesDiffAlias

	flagDiffLang     = "lang"
	flagDescDiffLang = "Language to diff (e.g. go, node). Required."
	flagDiffKind     = "kind"
	flagDescDiffKind = "Kind to diff (ignore | attributes). Default: both."
	flagDiffCwd      = "cwd"
	flagDescDiffCwd  = "Working directory to diff against (default: current dir)."

	errDiffLangRequired = "templates diff: --lang <name> is required\n"
	errDiffBadKind      = "templates diff: unknown --kind %q (want ignore | attributes)\n"
	errDiffResolve      = "templates diff: resolve %s/%s: %v\n"
	errDiffRun          = "templates diff: %v\n"

	msgDiffNoChange = "no changes for %s/%s in %s\n"

	// Exit codes mirror standard diff(1):
	//   0 = no changes
	//   1 = differences found
	//   2 = error
	exitDiffNoChange = 0
	exitDiffChanged  = 1
	exitDiffError    = 2
)

// runTemplatesDiff implements `gitmap templates diff --lang <l> [--kind <k>]`.
// Compares the on-disk gitmap-managed block (if any) against what
// `add <kind> <lang>` would write, prints unified-style hunks, and exits
// 0 / 1 / 2 per standard diff(1) conventions for script friendliness.
func runTemplatesDiff(args []string) {
	lang, kind, cwd := parseTemplatesDiffFlags(args)
	if lang == "" {
		fmt.Fprint(os.Stderr, errDiffLangRequired)
		os.Exit(exitDiffError)
	}
	kinds, ok := resolveDiffKinds(kind)
	if !ok {
		fmt.Fprintf(os.Stderr, errDiffBadKind, kind)
		os.Exit(exitDiffError)
	}

	anyChanged := runDiffForKinds(kinds, lang, cwd)
	if anyChanged {
		os.Exit(exitDiffChanged)
	}
	os.Exit(exitDiffNoChange)
}

// parseTemplatesDiffFlags pulls --lang/--kind/--cwd out of args. All
// values are lowered + trimmed; --cwd defaults to "." (resolved against
// the process working directory by templates.Diff -> filepath.Abs).
func parseTemplatesDiffFlags(args []string) (lang, kind, cwd string) {
	fs := flag.NewFlagSet(cmdTemplatesDiff, flag.ExitOnError)
	langPtr := fs.String(flagDiffLang, "", flagDescDiffLang)
	kindPtr := fs.String(flagDiffKind, "", flagDescDiffKind)
	cwdPtr := fs.String(flagDiffCwd, ".", flagDescDiffCwd)
	reordered := reorderFlagsBeforeArgs(args)
	_ = fs.Parse(reordered)

	return strings.ToLower(strings.TrimSpace(*langPtr)),
		strings.ToLower(strings.TrimSpace(*kindPtr)),
		strings.TrimSpace(*cwdPtr)
}

// resolveDiffKinds expands "" -> {ignore, attributes}. A specific kind
// passes through after validation; anything else fails.
func resolveDiffKinds(kind string) ([]string, bool) {
	switch kind {
	case "":
		return []string{"ignore", "attributes"}, true
	case "ignore", "attributes":
		return []string{kind}, true
	}

	return nil, false
}

// runDiffForKinds runs diffOneKind over each kind, returning true if
// any kind reported a change. Errors during resolution are surfaced
// with exit 2; missing-template-for-lang is treated as a hard error
// because the user explicitly asked for it.
func runDiffForKinds(kinds []string, lang, cwd string) bool {
	any := false
	sort.Strings(kinds)
	for _, k := range kinds {
		if diffOneKind(k, lang, cwd) {
			any = true
		}
	}

	return any
}

// diffOneKind resolves the template, computes the diff against the
// matching target file, and prints hunks (or a no-change line). Returns
// true when the on-disk content differs from the template.
func diffOneKind(kind, lang, cwd string) bool {
	r, err := templates.Resolve(kind, lang)
	if err != nil {
		fmt.Fprintf(os.Stderr, errDiffResolve, kind, lang, err)
		os.Exit(exitDiffError)
	}
	target := filepath.Join(cwd, targetFileFor(kind))
	tag := kind + "/" + lang
	res, err := templates.Diff(target, tag, r.Content)
	if err != nil {
		fmt.Fprintf(os.Stderr, errDiffRun, err)
		os.Exit(exitDiffError)
	}
	if res.Status == templates.DiffNoChange {
		fmt.Printf(msgDiffNoChange, kind, lang, res.Path)

		return false
	}
	printDiffHunks(res.Hunks)

	return true
}

// targetFileFor maps a template kind to its on-disk file name. Kept
// in sync with the addTemplateSpec table in addignoreattrs.go.
func targetFileFor(kind string) string {
	if kind == "attributes" {
		return ".gitattributes"
	}

	return ".gitignore"
}

// printDiffHunks writes hunks to stdout, colorized via the existing
// render package when stdout is a TTY. `+` lines are green-ish (we
// reuse cyan tokens since the pretty renderer hasn't allocated green
// yet), `-` lines are yellow; `@@` banners are muted.
func printDiffHunks(hunks []string) {
	useColor := render.StdoutIsTerminal()
	for _, h := range hunks {
		fmt.Println(decorateDiffLine(h, useColor))
	}
}

// decorateDiffLine wraps one hunk line in pretty-render tokens (then
// substitutes ANSI codes) when color is requested. No-op otherwise so
// piped output stays a clean unified-style diff that other tools can
// re-parse.
func decorateDiffLine(line string, useColor bool) string {
	if !useColor || line == "" {
		return line
	}
	switch line[0] {
	case '+':
		return render.HighlightQuotesANSI(constants.ColorCyan + line + constants.ColorReset)
	case '-':
		return render.HighlightQuotesANSI(constants.ColorYellow + line + constants.ColorReset)
	case '@':
		return render.HighlightQuotesANSI(constants.ColorDim + line + constants.ColorReset)
	}

	return line
}
