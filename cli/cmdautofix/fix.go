package cmdautofix

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

// Exit-code contract (spec §2): 0 clean / nothing to do / applied cleanly;
// 1 findings in check-only flow or unfixable findings remain; 2 tool error.
// Signalled through the established cliexit pattern (cliexit.Exit /
// cliexit.Fail / cliexit.HandleUsageError) so the theme/glyphs pipe drainers
// run before process teardown.

// fixSubcommands is the valid subcommand set (spec §1).
var fixSubcommands = []string{
	"encoding", "newlines", "naming", "paths", "gofmt",
	"misspell", "markdown", "guidelines", "all",
}

// RunFixCmd is the `gitmap fix` entry point: parent with per-category
// subcommands. Check-driven flow: Scan → MANDATORY summary → prompt → Apply.
func RunFixCmd(args []string) error {
	sub, rest := splitSubcommand(args)
	if sub == "" {
		RenderFixHelp()
		cliexit.Exit(0)
	}
	if sub == "--help" || sub == "-h" {
		RenderFixHelp()
		cliexit.Exit(0)
	}
	if !isFixSubcommand(sub) {
		if isHelpFlag(sub) {
			renderCategoryOrParentHelp(rest)
			cliexit.Exit(0)
		}
		cliexit.HandleUsageError(fmt.Errorf(
			"unknown subcommand %q — valid subcommands: %s\nNote: git-state remediation moved to 'gitmap stash' / 'gitmap wip' / 'gitmap discard'",
			sub, strings.Join(fixSubcommands, ", ")))
	}
	if hasHelpFlag(rest) {
		RenderFixCategoryHelp(sub)
		cliexit.Exit(0)
	}

	// The flag package stops parsing at the first positional arg, but the
	// spec shape is `fix <category> [path] [flags]` — flags may trail the
	// path. Reorder so flags always precede positionals.
	rest = reorderFixArgs(rest)

	fs := flag.NewFlagSet("fix", flag.ContinueOnError)
	fs.Usage = func() { RenderFixCategoryHelp(sub) }
	yes := fs.Bool("yes", false, "Apply without prompting")
	fs.BoolVar(yes, "y", false, "Apply without prompting (shorthand)")
	workers := fs.Int("workers", runtime.NumCPU(), "Worker threads")
	fs.IntVar(workers, "w", runtime.NumCPU(), "Worker threads (shorthand)")
	exts := fs.String("ext", "", "Comma-separated extension filter (e.g. .go,.md)")
	uriPattern := fs.String("uri-pattern", defaultRepoURI, "REPO_FILE_URI override for the paths category")
	asJSON := fs.Bool("json", false, "Machine-readable report (no prompt unless -y)")
	noCache := fs.Bool("no-cache", false, "Bypass the scan-result cache")

	if err := fs.Parse(rest); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			cliexit.Exit(0)
		}
		cliexit.HandleUsageError(fmt.Errorf("parse fix flags: %w", err))
	}

	if *workers < 1 {
		cliexit.HandleUsageError(fmt.Errorf("--workers must be >= 1 (got %d)", *workers))
	}

	root := "."
	if positional := fs.Args(); len(positional) > 0 {
		root = positional[0]
	}
	info, statErr := os.Stat(root)
	if statErr != nil || !info.IsDir() {
		cliexit.Fail("fix", "scan root", root, fmt.Errorf("unreadable scan path: %w", statErr), 2)
	}

	categories := []string{sub}
	if sub == "all" {
		categories = categoryNames()
	}

	opts := Options{
		Root:       root,
		Workers:    *workers,
		Exts:       normalizeExts(*exts),
		Categories: categories,
		URIPattern: *uriPattern,
		AsJSON:     *asJSON,
		NoCache:    *noCache,
		Yes:        *yes,
	}

	result, err := Scan(opts)
	if err != nil {
		if errors.Is(err, errToolMissing) {
			cliexit.Fail("fix", "tool check", "", err, 2)
		}
		cliexit.HandleUsageError(err)
	}

	fixable := countFixable(result)
	if opts.AsJSON {
		return runJSONFlow(opts, result, fixable)
	}
	return runHumanFlow(opts, result, fixable)
}

// countFixable returns the number of violations in categories that have a Fix.
func countFixable(result *ScanResult) int {
	n := 0
	for _, v := range result.Violations {
		if categoryHasFix(v.Category) {
			n++
		}
	}
	return n
}

func categoryHasFix(name string) bool {
	for _, cat := range categoryRegistry {
		if cat.Name == name {
			return cat.Fix != nil
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Human flow: mandatory summary → prompt → apply
// ---------------------------------------------------------------------------

func runHumanFlow(opts Options, result *ScanResult, fixable int) error {
	printScanSummary(opts, result)

	if len(result.Violations) == 0 {
		fmt.Println("clean: no findings")
		return nil
	}

	if fixable == 0 {
		printViolationList(result)
		fmt.Println("report-only findings — nothing to apply")
		cliexit.Exit(1)
	}

	apply := opts.Yes
	if !apply {
		apply = askApply()
	}
	if !apply {
		fmt.Println("No changes written.")
		return nil
	}

	applied, err := Apply(result, opts)
	if err != nil {
		cliexit.Fail("fix", "apply", "", err, 2)
	}
	fmt.Printf("Modified %d files in %.2f s.\n", applied.FilesModified, applied.Elapsed.Seconds())

	if unfixableRemain(result, applied) {
		printViolationList(result)
		cliexit.Exit(1)
	}
	return nil
}

// unfixableRemain reports whether findings survive the apply pass by design
// (naming is report-only) or by necessity (skipped, never lossy-written).
func unfixableRemain(result *ScanResult, applied *ApplyResult) bool {
	if len(applied.Skipped) > 0 {
		return true
	}
	for _, v := range result.Violations {
		if !categoryHasFix(v.Category) {
			return true
		}
	}
	return false
}

// askApply prompts on stderr and reads one line from the REAL stdin —
// /dev/stdin on unix so the prompt works even when stdout is piped.
// (Windows fallback: os.Stdin, which has no /dev/stdin.)
func askApply() bool {
	fmt.Fprint(os.Stderr, "Apply these fixes? [y/N] ")
	in := os.Stdin
	if runtime.GOOS != "windows" {
		if f, err := os.OpenFile("/dev/stdin", os.O_RDONLY, 0); err == nil {
			defer f.Close()
			in = f
		}
	}
	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil {
		return false
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes"
}

// ---------------------------------------------------------------------------
// JSON flow: machine-readable report; no prompt unless -y
// ---------------------------------------------------------------------------

func runJSONFlow(opts Options, result *ScanResult, fixable int) error {
	var applied *ApplyResult
	if opts.Yes && fixable > 0 {
		var err error
		applied, err = Apply(result, opts)
		if err != nil {
			cliexit.Fail("fix", "apply", "", err, 2)
		}
	}
	printJSONReport(opts, result, applied)

	if len(result.Violations) == 0 {
		return nil
	}
	if applied != nil && !unfixableRemain(result, applied) {
		return nil
	}
	cliexit.Exit(1)
	return nil
}

