package cmdautofix

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/output"
)

// ---------------------------------------------------------------------------
// Help — program 243 DRY rule: Go structs are the source of truth.
// helptext/*.md topics are GENERATED from these by the lead, never
// hand-edited here.
// ---------------------------------------------------------------------------

// FixHelpFlag is one flag row in the help screen.
type FixHelpFlag struct {
	Flags string
	Desc  string
}

// FixHelp is the struct source of truth for `gitmap fix --help`.
var FixHelp = struct {
	Usage       string
	Description string
	Subcommands [][2]string
	Flags       []FixHelpFlag
	Flow        string
	ExitCodes   string
	Note        string
}{
	Usage:       "gitmap fix <category> [path] [flags]",
	Description: "Check-driven parallel file-hygiene fixer — native port of the 03-ai-scripts fixers.\nRuns the check, prints a summary, prompts before writing.",
	Subcommands: [][2]string{
		{"encoding", "UTF-8 no-BOM + LF normalization"},
		{"newlines", "CRLF→LF, trim trailing whitespace, one final newline"},
		{"naming", "boolean-comparison style audit (report-only)"},
		{"paths", "Windows absolute path sanitizer (docs)"},
		{"gofmt", "gofmt -w over .go files"},
		{"misspell", "British→American spelling, case-preserving"},
		{"markdown", "collapse 3+ blank lines in .md"},
		{"guidelines", "composite: newlines + naming"},
		{"release tags", "purge broken releases and orphan release tags"},
		{"all", "run every category (one summary, one prompt)"},
	},
	Flags: []FixHelpFlag{
		{"-y, --yes", "Apply without prompting"},
		{"-w, --workers N", "Worker threads (default: CPU count)"},
		{"--ext .go,.md", "Only scan these extensions (default: category scope)"},
		{"--uri-pattern P", "REPO_FILE_URI override for the paths category"},
		{"--json", "Machine-readable report (no prompt unless -y)"},
		{"--no-cache", "Bypass the scan-result cache"},
	},
	Flow:      "scan → summary (files scanned, files modified, per-category fix counts,\ntime taken) → prompt \"Apply these fixes? [y/N]\" → apply.",
	ExitCodes: "0 clean · 1 findings remain · 2 tool error",
	Note:      "git-state remediation (stash/wip/discard) moved to\n'gitmap stash' / 'gitmap wip' / 'gitmap discard'.",
}

// RenderFixHelp prints the parent help screen from the FixHelp struct.
// Exposed so `gitmap help fix` (lead-wired) renders the same content.
// Writes through output.UI() so safe-mode glyph filtering applies.
func RenderFixHelp() {
	w := output.UI()
	fmt.Fprintf(w, "Usage: %s\n\n", FixHelp.Usage)
	fmt.Fprintln(w, FixHelp.Description)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Subcommands:")
	for _, s := range FixHelp.Subcommands {
		fmt.Fprintf(w, "  %-11s %s\n", s[0], s[1])
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	for _, f := range FixHelp.Flags {
		fmt.Fprintf(w, "  %-18s %s\n", f.Flags, f.Desc)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flow: "+FixHelp.Flow)
	fmt.Fprintln(w, "Exit codes: "+FixHelp.ExitCodes)
	fmt.Fprintln(w, "Note: "+FixHelp.Note)
}

// RenderFixCategoryHelp prints help for one subcommand.
func RenderFixCategoryHelp(sub string) {
	w := output.UI()
	desc := ""
	for _, s := range FixHelp.Subcommands {
		if s[0] == sub {
			desc = s[1]
			break
		}
	}
	fmt.Fprintf(w, "Usage: gitmap fix %s [path] [flags]\n\n", sub)
	if desc != "" {
		fmt.Fprintln(w, desc+".")
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w, "Flags:")
	for _, f := range FixHelp.Flags {
		fmt.Fprintf(w, "  %-18s %s\n", f.Flags, f.Desc)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Exit codes: "+FixHelp.ExitCodes)
}

// ---------------------------------------------------------------------------
// Skill surface — spec §7 verbatim. The lead embeds this into
// `gitmap llm train` and .agents/skills/gitmap/SKILL.md.
// ---------------------------------------------------------------------------

// SkillSnippetMD returns the static markdown block teaching AI agents the
// fix workflow. Keep the text byte-identical to spec §7.
func SkillSnippetMD() string {
	return "### gitmap fix — parallel file-hygiene fixer (check-driven)\n" +
		"Check: `gitmap fix <category> [path]` — scans, prints summary, prompts\n" +
		"`Apply these fixes? [y/N]`. Categories: encoding, newlines, naming, paths,\n" +
		"gofmt, misspell, markdown, guidelines, all.\n" +
		"Apply without prompting: `gitmap fix all -y`.\n" +
		"Worker threads: `gitmap fix all -w 8` (default: CPU count).\n" +
		"Summary (always printed): files scanned, files modified, per-category fix\n" +
		"counts, time taken. Exit codes: 0 clean · 1 findings remain · 2 tool error.\n" +
		"Byte-safe: binaries skipped (null-byte probe); files rewritten only when\n" +
		"bytes differ; invalid-UTF-8 files are reported, never lossy-written.\n" +
		"`naming` (boolean-comparison style) is report-only by design.\n" +
		"Note: `gitmap fix` (git-state: stash/wip/discard) moved to\n" +
		"`gitmap stash` / `gitmap wip` / `gitmap discard`.\n" +
		"AI rule: run `gitmap fix all` (check) before committing hygiene-sensitive\n" +
		"work; never hand-roll sed/regex loops for encoding, newlines, or spelling.\n"
}
