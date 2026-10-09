package cmdautofix

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

// Exit-code contract (spec §4): 0 clean · 1 violations found ·
// 2 tool error. Signalled through the established cliexit pattern
// (cliexit.Exit / cliexit.Fail / cliexit.HandleUsageError) so the
// theme/glyphs pipe drainers run before process teardown.

// ---------------------------------------------------------------------------
// Help — program 243 DRY rule: Go structs are the source of truth.
// helptext/*.md topics are GENERATED from these by the lead, never
// hand-edited here.
// ---------------------------------------------------------------------------

// AutofixHelpFlag is one flag row in the help screen.
type AutofixHelpFlag struct {
	Flags string
	Desc  string
}

// AutofixHelp is the struct source of truth for `gitmap autofix --help`.
var AutofixHelp = struct {
	Usage       string
	Alias       string
	Description string
	Flags       []AutofixHelpFlag
	ExitCodes   string
}{
	Usage:       "gitmap autofix [path] [flags]",
	Alias:       "gitmap afx",
	Description: "Fix file hygiene issues in parallel — native port of the 03-ai-scripts autofixers.",
	Flags: []AutofixHelpFlag{
		{"--apply", "Write fixes (default: dry-run audit only)"},
		{"-w, --workers N", "Worker threads (default: CPU count)"},
		{"--ext .go,.md", "Only scan these extensions (default: all text files)"},
		{"--category a,b", "Subset of categories (default: all)"},
		{"--uri-pattern P", "REPO_FILE_URI override for the paths category"},
		{"--json", "Machine-readable report"},
	},
	ExitCodes: "0 clean · 1 violations found · 2 tool error",
}

// RenderAutofixHelp prints the --help screen from the AutofixHelp struct.
// Exposed so `gitmap help autofix` (lead-wired) renders the same content.
func RenderAutofixHelp() {
	fmt.Printf("Usage: %s  (alias: %s)\n\n", AutofixHelp.Usage, AutofixHelp.Alias)
	fmt.Println(AutofixHelp.Description)
	fmt.Println()
	fmt.Println("Flags:")
	for _, f := range AutofixHelp.Flags {
		fmt.Printf("  %-18s %s\n", f.Flags, f.Desc)
	}
	fmt.Println()
	fmt.Println("Categories:")
	for _, cat := range categoryRegistry {
		fmt.Printf("  %-11s %s\n", cat.Name, cat.Desc)
	}
	fmt.Println()
	fmt.Println("Exit codes: " + AutofixHelp.ExitCodes)
}

// ---------------------------------------------------------------------------
// Skill surface — spec §7 verbatim. The lead embeds this into
// `gitmap llm train` and .agents/skills/gitmap/SKILL.md.
// ---------------------------------------------------------------------------

// SkillSnippetMD returns the static markdown block teaching AI agents the
// autofix workflow. Keep the text byte-identical to spec §7.
func SkillSnippetMD() string {
	return "### gitmap autofix — parallel file-hygiene fixer\n" +
		"Dry-run audit (CI gate): `gitmap autofix` — exits 1 on violations.\n" +
		"Apply fixes: `gitmap autofix --apply`.\n" +
		"Worker threads: `gitmap autofix --apply -w 8` (default: CPU count).\n" +
		"Categories (grouping): encoding, newlines, naming, paths, gofmt,\n" +
		"misspell, markdown, guidelines (composite). Select a subset:\n" +
		"`gitmap autofix --category encoding,newlines`.\n" +
		"Scope: `gitmap autofix [path] [--ext .go,.md] [--json]`.\n" +
		"Exit codes: 0 clean · 1 violations found · 2 tool error.\n" +
		"Byte-safe: binaries skipped (null-byte probe); files rewritten only\n" +
		"when bytes differ; invalid-UTF-8 files are reported, never lossy-written.\n" +
		"`naming` (boolean-comparison style) is report-only by design.\n" +
		"AI rule: run `gitmap autofix` (dry-run) before committing hygiene-sensitive\n" +
		"work; never hand-roll sed/regex loops for encoding, newlines, or spelling.\n"
}

// ---------------------------------------------------------------------------
// CLI entry
// ---------------------------------------------------------------------------

// RunAutofixCmd is the `gitmap autofix` / `gitmap afx` entry point.
// Signature is fixed for the LLM train surface; it exits via cliexit with
// the spec §4 exit codes rather than returning them.
func RunAutofixCmd(args []string) error {
	if hasHelpFlag(args) {
		RenderAutofixHelp()
		cliexit.Exit(0)
	}

	fs := flag.NewFlagSet("autofix", flag.ContinueOnError)
	fs.Usage = RenderAutofixHelp
	apply := fs.Bool("apply", false, "Write fixes (default: dry-run audit only)")
	workers := fs.Int("workers", runtime.NumCPU(), "Worker threads")
	fs.IntVar(workers, "w", runtime.NumCPU(), "Worker threads (shorthand)")
	exts := fs.String("ext", "", "Comma-separated extension filter (e.g. .go,.md)")
	categories := fs.String("category", "", "Comma-separated category subset")
	uriPattern := fs.String("uri-pattern", defaultRepoURI, "REPO_FILE_URI override for the paths category")
	asJSON := fs.Bool("json", false, "Machine-readable report")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			cliexit.Exit(0)
		}
		cliexit.HandleUsageError(fmt.Errorf("parse autofix flags: %w", err))
	}

	if *workers < 1 {
		cliexit.HandleUsageError(fmt.Errorf("--workers must be >= 1 (got %d)", *workers))
	}

	root := "."
	if rest := fs.Args(); len(rest) > 0 {
		root = rest[0]
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		cliexit.Fail("autofix", "scan root", root, fmt.Errorf("unreadable scan path: %w", err), 2)
	}

	opts := Options{
		Apply:      *apply,
		Workers:    *workers,
		Exts:       normalizeExts(*exts),
		Categories: splitCSV(*categories),
		URIPattern: *uriPattern,
		AsJSON:     *asJSON,
	}

	result, err := Run(root, opts)
	if err != nil {
		if errors.Is(err, errToolMissing) {
			cliexit.Fail("autofix", "tool check", "", err, 2)
		}
		cliexit.HandleUsageError(err)
	}

	if opts.AsJSON {
		printJSONReport(root, opts, result)
	} else {
		printHumanReport(root, opts, result)
	}

	if len(result.Violations) > 0 {
		cliexit.Exit(1)
	}
	return nil
}

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			return true
		}
	}
	return false
}

// normalizeExts lowercases extensions and ensures a leading dot.
func normalizeExts(raw string) []string {
	parts := splitCSV(raw)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.ToLower(p)
		if !strings.HasPrefix(p, ".") {
			p = "." + p
		}
		out = append(out, p)
	}
	return out
}

func splitCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Reports
// ---------------------------------------------------------------------------

const maxHumanViolations = 50

func fileCount(violations []Violation) int {
	seen := map[string]struct{}{}
	for _, v := range violations {
		seen[v.Path] = struct{}{}
	}
	return len(seen)
}

func printHumanReport(root string, opts Options, result *RunResult) {
	mode := "dry-run audit"
	if opts.Apply {
		mode = "apply mode"
	}
	fmt.Printf("gitmap autofix — %s of '%s' (%d workers, %d files scanned)\n", mode, root, opts.Workers, result.FilesSeen)

	if opts.Apply && len(result.FixedFiles) > 0 {
		fmt.Printf("fixed %d file(s):\n", len(result.FixedFiles))
		for _, f := range result.FixedFiles {
			fmt.Printf("  fixed: %s\n", f)
		}
	}

	if len(result.Violations) == 0 {
		fmt.Println("clean: no violations")
		return
	}

	shown := result.Violations
	truncated := 0
	if len(shown) > maxHumanViolations {
		truncated = len(shown) - maxHumanViolations
		shown = shown[:maxHumanViolations]
	}
	for _, v := range shown {
		loc := v.Path
		if v.Line > 0 {
			loc = fmt.Sprintf("%s:%d", v.Path, v.Line)
		}
		fmt.Printf("  %s [%s] %s\n", loc, v.Category, v.Detail)
	}
	if truncated > 0 {
		fmt.Printf("  ... and %d more (use --json for the full list)\n", truncated)
	}

	if opts.Apply {
		fmt.Printf("%d violation(s) remain in %d file(s)\n", len(result.Violations), fileCount(result.Violations))
	} else {
		fmt.Printf("%d violation(s) in %d file(s) — CI gate FAILED\n", len(result.Violations), fileCount(result.Violations))
	}
}

// jsonReport is the --json machine-readable report shape.
type jsonReport struct {
	Command    string      `json:"command"`
	Path       string      `json:"path"`
	Apply      bool        `json:"apply"`
	Workers    int         `json:"workers"`
	Categories []string    `json:"categories"`
	FilesSeen  int         `json:"files_scanned"`
	Violations []Violation `json:"violations"`
	FixedFiles []string    `json:"fixed_files"`
	ExitCode   int         `json:"exit_code"`
}

func printJSONReport(root string, opts Options, result *RunResult) {
	exitCode := 0
	if len(result.Violations) > 0 {
		exitCode = 1
	}
	cats := opts.Categories
	if len(cats) == 0 {
		cats = categoryNames()
	}
	violations := result.Violations
	if violations == nil {
		violations = []Violation{}
	}
	fixed := result.FixedFiles
	if fixed == nil {
		fixed = []string{}
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		abs = root
	}
	report := jsonReport{
		Command:    "autofix",
		Path:       abs,
		Apply:      opts.Apply,
		Workers:    opts.Workers,
		Categories: cats,
		FilesSeen:  result.FilesSeen,
		Violations: violations,
		FixedFiles: fixed,
		ExitCode:   exitCode,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(report)
}
