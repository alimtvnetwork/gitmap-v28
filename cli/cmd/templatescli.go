package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/render"
	"github.com/alimtvnetwork/gitmap-v28/cli/templates"
)

const (
	cmdTemplatesList      = "list"
	cmdTemplatesListAlias = "tl"
	cmdTemplatesShow      = "show"
	cmdTemplatesShowAlias = "ts"
	cmdTemplatesInit      = constants.CmdTemplatesInit
	cmdTemplatesInitAlias = constants.CmdTemplatesInitAlias
	usageTemplatesRoot    = `Usage: gitmap templates <subcommand>

Subcommands:
  list [flags]               List every available template (alias: tl)
  show <kind> <lang>         Print a single template to stdout (alias: ts)
  init <lang> [<lang>...]    Scaffold .gitignore / .gitattributes for languages (alias: ti)
  diff --lang <l> [--kind k] Show what 'add' would change without writing (alias: td)
  ls [--category <c>] [--json] List state templates from gitmap-templates.db
  add [flags]                Add/upsert a state template (--category, --slug, --title, --text)
  edit <id|slug> [flags]     Edit an existing state template
  remove <id|slug>           Remove a state template (aliases: rm, delete)
  import <file.json>         Import state templates & variables with exportId deduplication
  export <dest.json>         Export state templates & referenced variables
  var <set|ls|rm>            Manage template variables ($VAR)
  ui [--port 8787]           Launch the local Templates & Variables Web UI

Kinds:
  ignore | attributes | lfs

Flags (list):
  --kind <ignore|attributes|lfs>   Filter rows to one kind
  --lang <name>                    Filter rows to one language (matches across kinds)

Flags (show):
  --raw                      Disable pretty markdown rendering even on a TTY

Flags (init):
  --lfs                      Also merge lfs/common.gitattributes into .gitattributes
  --dry-run                  Preview every block; do not touch disk
  --force                    Replace existing .gitignore/.gitattributes outright

Flags (diff):
  --lang <name>              Required. Which language to diff.
  --kind <ignore|attributes> Default: both. Restrict to one kind.
  --cwd <path>               Default: current dir. Where to look for the target file.

Examples:
  gitmap templates list
  gitmap templates list --kind ignore
  gitmap templates list --lang go
  gitmap tpl tl --kind attributes
  gitmap templates show ignore go
  gitmap tpl ts attributes node
  gitmap templates show ignore go --raw   # bypass pretty renderer
  gitmap templates init go
  gitmap templates init go node --lfs
  gitmap tpl ti python --dry-run
  gitmap templates diff --lang go
  gitmap tpl td --lang node --kind ignore
`
	headerTemplatesList    = "KIND        LANG            SOURCE  PATH\n"
	fmtTemplatesListRow    = "%-10s  %-14s  %-6s  %s\n"
	labelTemplatesUser     = "user"
	labelTemplatesEmbed    = "embed"
	msgTemplatesEmpty      = "(no templates registered — embedded corpus is empty)\n"
	msgTemplatesFiltered   = "(no templates match the requested filter)\n"
	errTemplatesShowArgs   = "templates show requires <kind> <lang>; e.g. 'templates show ignore go'\n"
	errTemplatesShowFail   = "templates show: %v\n"
	errTemplatesListFail   = "templates list: %v\n"
	errTemplatesListKind   = "templates list: unknown --kind %q (want ignore | attributes | lfs)\n"
	errUnknownTemplatesSub = "unknown 'templates' subcommand: %s\n"
	flagTemplatesShowRaw   = "raw"
	flagDescTemplatesRaw   = "Deprecated alias for --no-pretty (kept for v3.23.x back-compat)"
	flagTemplatesListKind  = "kind"
	flagDescListKind       = "Filter to one kind (ignore | attributes | lfs)"
	flagTemplatesListLang  = "lang"
	flagDescListLang       = "Filter to one language (matches across kinds)"
)

// dispatchTemplates routes `gitmap templates <subcommand>` calls.
func dispatchTemplates(command string) (bool, error) {
	if command != constants.CmdTemplates && command != constants.CmdTemplatesAlias {
		return false, nil
	}

	if len(os.Args) < 3 {
		exitTemplatesRootUsage()

		return true, nil
	}

	routeTemplatesSubcommand(os.Args[2], os.Args[3:])

	return true, nil
}

func exitTemplatesRootUsage() {
	err := apperror.NewWithDetails(
		"cmd.templates.dispatch", "E1103", usageTemplatesRoot,
		"cmd.templates", apperror.ErrorTypeValidation, apperror.SeverityError, nil,
	)
	cliexit.HandleError(err, 1)
}

func routeTemplatesSubcommand(sub string, rest []string) {
	if dispatchStateTemplatesSub(sub, rest) {
		return
	}

	switch sub {
	case cmdTemplatesList, cmdTemplatesListAlias:
		_ = runTemplatesList(rest)
	case cmdTemplatesShow, cmdTemplatesShowAlias:
		_ = runTemplatesShow(rest)
	case cmdTemplatesInit, cmdTemplatesInitAlias:
		runTemplatesInit(rest)
	case cmdTemplatesDiff, cmdTemplatesDiffAlias:
		runTemplatesDiff(rest)
	default:
		exitUnknownTemplatesSub(sub)
	}
}

func exitUnknownTemplatesSub(sub string) {
	err := apperror.NewWithDetails(
		"cmd.templates.dispatch.unknown", "E1104", fmt.Sprintf(errUnknownTemplatesSub, sub),
		"cmd.templates", apperror.ErrorTypeValidation, apperror.SeverityError, map[string]any{"subcommand": sub},
	)
	cliexit.HandleError(err, 1)
}

// runTemplatesList prints every available template grouped by kind.
func runTemplatesList(args []string) error {
	kindFilter, langFilter := parseTemplatesListFlags(args)
	if !isValidKindFilter(kindFilter) {
		exitInvalidKindFilter(kindFilter)

		return nil
	}

	entries, err := templates.List()
	if err != nil {
		exitTemplatesListLoadErr(err)

		return nil
	}

	printFilteredTemplateEntries(entries, kindFilter, langFilter)

	return nil
}

func exitInvalidKindFilter(kindFilter string) {
	err := apperror.NewWithDetails(
		"cmd.templates.list.kindFilter", "E1105", fmt.Sprintf("unknown template kind filter '%s'", kindFilter),
		"cmd.templates", apperror.ErrorTypeValidation, apperror.SeverityError, map[string]any{"kind": kindFilter},
	)
	cliexit.HandleError(err, 1)
}

func exitTemplatesListLoadErr(err error) {
	appErr := apperror.WrapWithDetails(
		err, "cmd.templates.list.load", "E1106", "failed to load templates list",
		"cmd.templates", apperror.ErrorTypeExecution, apperror.SeverityError, nil,
	)
	cliexit.HandleError(appErr, 1)
}

func printFilteredTemplateEntries(entries []templates.Entry, kindFilter, langFilter string) {
	if len(entries) == 0 {
		fmt.Print(msgTemplatesEmpty)

		return
	}

	filtered := filterTemplates(entries, kindFilter, langFilter)
	if len(filtered) == 0 {
		fmt.Print(msgTemplatesFiltered)

		return
	}

	fmt.Print(headerTemplatesList)
	for _, e := range filtered {
		fmt.Printf(fmtTemplatesListRow, e.Kind, e.Lang, sourceLabel(e.Source), e.Path)
	}
}

// parseTemplatesListFlags pulls --kind/--lang out of args. Both are
// optional and case-insensitive on value. Unknown positional args are
// silently ignored — list takes no positional input.
func parseTemplatesListFlags(args []string) (string, string) {
	fs := flag.NewFlagSet(cmdTemplatesList, flag.ExitOnError)
	kind := fs.String(flagTemplatesListKind, "", flagDescListKind)
	lang := fs.String(flagTemplatesListLang, "", flagDescListLang)
	reordered := reorderFlagsBeforeArgs(args)
	_ = fs.Parse(reordered)

	return strings.ToLower(strings.TrimSpace(*kind)),
		strings.ToLower(strings.TrimSpace(*lang))
}

// isValidKindFilter accepts the empty string (no filter) and the three
// canonical kinds. Anything else trips the errTemplatesListKind exit.
func isValidKindFilter(kind string) bool {
	switch kind {
	case "", "ignore", "attributes", "lfs":
		return true
	}

	return false
}

// filterTemplates is a pure helper so tests can pin filter semantics
// without spinning up the full templates.List() embed-walk.
func filterTemplates(in []templates.Entry, kindFilter, langFilter string) []templates.Entry {
	if kindFilter == "" && langFilter == "" {
		return in
	}

	out := make([]templates.Entry, 0, len(in))
	for _, e := range in {
		if isMatchingTemplateEntry(e, kindFilter, langFilter) {
			out = append(out, e)
		}
	}

	return out
}

func isMatchingTemplateEntry(e templates.Entry, kindFilter, langFilter string) bool {
	if kindFilter != "" && e.Kind != kindFilter {
		return false
	}

	if langFilter != "" && !strings.EqualFold(e.Lang, langFilter) {
		return false
	}

	return true
}

// runTemplatesShow prints one template to stdout.
func runTemplatesShow(args []string) error {
	rest, mode := parseTemplatesShowFlags(args)
	if len(rest) < 2 {
		exitTemplatesShowArgsErr()

		return nil
	}

	r, err := templates.Resolve(rest[0], rest[1])
	if err != nil {
		exitTemplatesShowResolveErr(err, rest[0], rest[1])

		return nil
	}

	writeResolvedTemplateStdout(&r, mode)

	return nil
}

func exitTemplatesShowArgsErr() {
	err := apperror.NewWithDetails(
		"cmd.templates.show.args", "E1107", errTemplatesShowArgs,
		"cmd.templates", apperror.ErrorTypeValidation, apperror.SeverityError, nil,
	)
	cliexit.HandleError(err, 1)
}

func exitTemplatesShowResolveErr(err error, kind, lang string) {
	appErr := apperror.WrapWithDetails(
		err, "cmd.templates.show.resolve", "E1108", "failed to resolve template",
		"cmd.templates", apperror.ErrorTypeValidation, apperror.SeverityError, map[string]any{"kind": kind, "lang": lang},
	)
	cliexit.HandleError(appErr, 1)
}

func writeResolvedTemplateStdout(r *templates.Resolved, mode render.PrettyModeType) {
	out := r.Content
	if render.Decide(mode, render.StdoutIsTerminal(), isMarkdownTemplatePath(r.Path)) {
		out = []byte(render.RenderANSI(string(r.Content)))
	}

	if _, err := os.Stdout.Write(out); err != nil {
		appErr := apperror.WrapWithDetails(
			err, "cmd.templates.show.write", "E1109", "failed to write template content to stdout",
			"cmd.templates", apperror.ErrorTypeExecution, apperror.SeverityError, nil,
		)
		cliexit.HandleError(appErr, 1)
	}
}

// parseTemplatesShowFlags extracts --pretty / --no-pretty (preferred) and
// the legacy --raw alias from args, returning the cleaned positional
// slice + the resolved render.PrettyMode. --raw is treated as
// --no-pretty for back-compat with v3.23.x. When both are present,
// --pretty wins (--raw only downgrades when mode is still Auto).
func parseTemplatesShowFlags(args []string) ([]string, render.PrettyModeType) {
	cleaned, mode := ParsePrettyFlag(args)

	fs := flag.NewFlagSet(cmdTemplatesShow, flag.ExitOnError)
	rawFlag := fs.Bool(flagTemplatesShowRaw, false, flagDescTemplatesRaw)
	reordered := reorderFlagsBeforeArgs(cleaned)
	_ = fs.Parse(reordered)

	if *rawFlag && mode == render.PrettyAuto {
		mode = render.PrettyOff
	}

	return fs.Args(), mode
}

// isMarkdownTemplatePath returns true for .md / .markdown files
// (case-insensitive). Templates today are .gitignore / .gitattributes —
// this guard future-proofs the renderer for markdown overlays
// (e.g. ~/.gitmap/templates/notes/*.md) without changing existing UX.
func isMarkdownTemplatePath(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))

	return ext == ".md" || ext == ".markdown"
}

// sourceLabel maps a templates.SourceType to the user-facing column value.
func sourceLabel(s templates.SourceType) string {
	if s == templates.SourceUser {
		return labelTemplatesUser
	}

	return labelTemplatesEmbed
}
