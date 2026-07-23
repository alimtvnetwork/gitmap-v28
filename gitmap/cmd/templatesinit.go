// Command `gitmap templates init <lang> [<lang>...] [--lfs] [--dry-run] [--force]`.
//
// Scaffolds .gitignore and .gitattributes for one or more languages by
// resolving the corresponding embedded (or user-overlay) templates and
// merging them into the target files via templates.Merge — the same
// idempotent marker-block primitive that powers `add lfs-install`.
//
// Behavior summary:
//   - For each <lang>, ignore/<lang>.gitignore is REQUIRED. Missing →
//     hard error (exit 1) with a one-liner pointing at `templates list`.
//   - attributes/<lang>.gitattributes is OPTIONAL. Missing → soft skip
//     with a dim "no attributes template for <lang>" notice. This matches
//     the embed corpus (every lang has ignore, only some have attributes).
//   - --lfs additionally merges lfs/common.gitattributes into
//     .gitattributes. Reuses templates.Merge with tag "lfs/common" so the
//     block is interchangeable with `gitmap add lfs-install`.
//   - --dry-run prints every block that WOULD be written and exits without
//     touching disk. Outcome verbs reflect what would happen (created /
//     would update / would insert).
//   - --force replaces any pre-existing target file outright with a fresh
//     gitmap-managed block, discarding hand edits OUTSIDE the markers.
//     Without --force, Merge preserves non-marker content and either
//     updates the existing block in place or appends one — see merge.go.
//
// Operates from CWD. Does NOT require being inside a git repo (scaffolding
// before `git init` is a legitimate workflow). The --lfs path also does
// NOT shell out to `git lfs install` — that's `add lfs-install`'s job.
// `templates init --lfs` is purely a template-merge operation.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/templates"
)

const (
	templatesInitLFSTag      = "lfs/common"
	templatesInitIgnoreFile  = ".gitignore"
	templatesInitAttrFile    = ".gitattributes"
	templatesInitKindIgnore  = "ignore"
	templatesInitKindAttribs = "attributes"
	templatesInitKindLFS     = "lfs"
	templatesInitLFSLang     = "common"
)

// templatesInitFlags holds the parsed flag state.
type templatesInitFlags struct {
	lfs    bool
	dryRun bool
	force  bool
	langs  []string
}

// runTemplatesInit is the entry point for `gitmap templates init`.
func runTemplatesInit(args []string) {
	checkHelp("templates-init", args)

	flags, err := parseTemplatesInitFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}
	if len(flags.langs) == 0 {
		fmt.Fprintln(os.Stderr, "  ✗ templates init requires at least one <lang>")
		fmt.Fprintln(os.Stderr, "    Example: gitmap templates init go node --lfs")
		fmt.Fprintln(os.Stderr, "    Run 'gitmap templates list' to see available languages.")
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Could not resolve CWD: %v\n", err)
		os.Exit(1)
	}

	printTemplatesInitBanner(flags, cwd)

	results := executeTemplatesInit(cwd, flags)
	printTemplatesInitSummary(results, flags.dryRun)
}

// parseTemplatesInitFlags extracts --lfs / --dry-run / --force from args
// and returns the remaining positional langs. Uses reorderFlagsBeforeArgs
// (per mem://tech/flag-parsing-logic) so any flag/positional order works:
// `templates init go --lfs` and `templates init --lfs go` are equivalent.
func parseTemplatesInitFlags(args []string) (templatesInitFlags, error) {
	fs := flag.NewFlagSet("templates init", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	lfs := fs.Bool("lfs", false, "also merge lfs/common.gitattributes into .gitattributes")
	dryRun := fs.Bool("dry-run", false, "preview every block that would be written; do not touch disk")
	force := fs.Bool("force", false, "overwrite existing .gitignore/.gitattributes outright (discards hand edits outside the gitmap marker block)")

	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		return templatesInitFlags{}, fmt.Errorf("parse flags: %w", err)
	}

	return templatesInitFlags{
		lfs:    *lfs,
		dryRun: *dryRun,
		force:  *force,
		langs:  fs.Args(),
	}, nil
}

// templatesInitStep is one resolved (kind, lang) → target merge unit.
type templatesInitStep struct {
	kind     string // ignore | attributes | lfs
	lang     string
	tag      string // marker block tag, e.g. "ignore/go"
	target   string // absolute path to .gitignore or .gitattributes
	resolved templates.Resolved
}

// templatesInitResult is the outcome of one step (or a soft skip).
type templatesInitResult struct {
	step       templatesInitStep
	skipped    bool
	skipReason string // dim line printed in summary
	merge      templates.MergeResult
	dryRun     bool
}

// executeTemplatesInit runs every (lang × kind) step plus the optional
// --lfs step. Returns one result per attempted step in deterministic
// order: per lang [ignore, attributes], then the single lfs step.
func executeTemplatesInit(cwd string, flags templatesInitFlags) []templatesInitResult {
	var results []templatesInitResult

	for _, lang := range flags.langs {
		results = append(results, runTemplatesInitStep(templatesInitStep{
			kind:   templatesInitKindIgnore,
			lang:   lang,
			tag:    templatesInitKindIgnore + "/" + lang,
			target: filepath.Join(cwd, templatesInitIgnoreFile),
		}, flags, true))

		results = append(results, runTemplatesInitStep(templatesInitStep{
			kind:   templatesInitKindAttribs,
			lang:   lang,
			tag:    templatesInitKindAttribs + "/" + lang,
			target: filepath.Join(cwd, templatesInitAttrFile),
		}, flags, false))
	}

	if flags.lfs {
		results = append(results, runTemplatesInitStep(templatesInitStep{
			kind:   templatesInitKindLFS,
			lang:   templatesInitLFSLang,
			tag:    templatesInitLFSTag,
			target: filepath.Join(cwd, templatesInitAttrFile),
		}, flags, true))
	}

	return results
}

// runTemplatesInitStep resolves a single (kind, lang) template and either
// merges it into target or returns a soft-skip result. required=true
// means a missing template is fatal; required=false makes it a soft skip.
func runTemplatesInitStep(step templatesInitStep, flags templatesInitFlags, required bool) templatesInitResult {
	res, err := templates.Resolve(step.kind, step.lang)
	if err != nil {
		if !required {
			return templatesInitResult{
				step:       step,
				skipped:    true,
				skipReason: fmt.Sprintf("no %s template for %s", step.kind, step.lang),
				dryRun:     flags.dryRun,
			}
		}
		fmt.Fprintf(os.Stderr, "  ✗ Required template missing: %v\n", err)
		fmt.Fprintln(os.Stderr, "    Run 'gitmap templates list' to see available languages.")
		os.Exit(1)
	}
	step.resolved = res

	if flags.dryRun {
		return templatesInitResult{
			step:   step,
			merge:  simulateTemplatesInitMerge(step, flags),
			dryRun: true,
		}
	}

	if flags.force {
		if err := os.Remove(step.target); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "  ✗ --force could not remove %s: %v\n", step.target, err)
			os.Exit(1)
		}
	}

	merged, err := templates.Merge(step.target, step.tag, res.Content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ merge %s into %s: %v\n", step.tag, step.target, err)
		os.Exit(1)
	}

	return templatesInitResult{step: step, merge: merged}
}

// simulateTemplatesInitMerge produces a MergeResult for --dry-run mode
// without touching disk. Inspects the target file to decide whether the
// outcome would be Created (no file), Updated (block already present),
// or Inserted (file exists but no marker block yet).
func simulateTemplatesInitMerge(step templatesInitStep, flags templatesInitFlags) templates.MergeResult {
	abs, _ := filepath.Abs(step.target)
	res := templates.MergeResult{Path: abs, BlockTag: step.tag, Changed: true}

	prior, err := os.ReadFile(step.target)
	if err != nil || flags.force {
		res.Outcome = templates.MergeCreated

		return res
	}
	marker := []byte("# >>> gitmap:" + step.tag + " >>>")
	if strings.Contains(string(prior), string(marker)) {
		res.Outcome = templates.MergeUpdated

		return res
	}
	res.Outcome = templates.MergeInserted

	return res
}

// printTemplatesInitBanner mirrors the visual idiom used by
// `add lfs-install` so the two scaffolders feel like siblings.
func printTemplatesInitBanner(flags templatesInitFlags, cwd string) {
	fmt.Println()
	fmt.Printf("  %s■ gitmap templates init —%s scaffold .gitignore / .gitattributes\n",
		constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  target dir: %s%s%s\n", constants.ColorDim, cwd, constants.ColorReset)
	fmt.Printf("  langs: %s%s%s",
		constants.ColorDim, strings.Join(flags.langs, ", "), constants.ColorReset)
	if flags.lfs {
		fmt.Printf("   %s+lfs%s", constants.ColorYellow, constants.ColorReset)
	}
	fmt.Println()
	if flags.dryRun {
		fmt.Printf("  %s[dry-run]%s no files will be modified\n",
			constants.ColorYellow, constants.ColorReset)
	}
	if flags.force {
		fmt.Printf("  %s[--force]%s pre-existing .gitignore/.gitattributes will be discarded\n",
			constants.ColorYellow, constants.ColorReset)
	}
	fmt.Println()
}

// printTemplatesInitSummary renders one line per step in the same
// verb / color vocabulary as printAddLFSInstallSummary.
func printTemplatesInitSummary(results []templatesInitResult, dryRun bool) {
	for _, r := range results {
		if r.skipped {
			fmt.Printf("  %s•%s %s\n",
				constants.ColorDim, constants.ColorReset, r.skipReason)

			continue
		}
		printTemplatesInitStepLine(r, dryRun)
	}
	fmt.Println()
	if dryRun {
		fmt.Printf("  %sRe-run without --dry-run to apply.%s\n",
			constants.ColorDim, constants.ColorReset)
	} else {
		fmt.Printf("  %sNext step:%s commit the scaffolded files:\n",
			constants.ColorYellow, constants.ColorReset)
		fmt.Println("    git add .gitignore .gitattributes")
		fmt.Println("    git commit -m \"chore: scaffold ignore/attributes via gitmap templates init\"")
	}
	fmt.Println()
}

// printTemplatesInitStepLine emits one outcome row, choosing verb +
// color from the merge result. Dry-run prefixes get the "would " marker
// so a quick visual scan distinguishes simulation from real output.
func printTemplatesInitStepLine(r templatesInitResult, dryRun bool) {
	verb := outcomeVerb(r.merge.Outcome)
	color := constants.ColorGreen
	if dryRun {
		verb = "would " + verb
		color = constants.ColorYellow
	} else if !r.merge.Changed {
		verb = "unchanged"
		color = constants.ColorDim
	}
	fmt.Printf("  %s%s%s %s (block: %s)\n",
		color, verb, constants.ColorReset, r.merge.Path, r.merge.BlockTag)
}
