// Commands `gitmap add ignore` and `gitmap add attributes`.
//
// Both share a near-identical pipeline:
//  1. Validate we are inside a Git working tree.
//  2. Resolve `common.<ext>` plus each language argument from the
//     templates package (overlay > embed).
//  3. Concatenate bodies with a single dedupe pass keyed on the
//     trimmed line — comments included, blank lines preserved.
//  4. Hand the merged body to templates.Merge under a marker block
//     tagged "<kind>/<concatenated-langs>" so re-runs are byte-stable.
//
// The two entry points (`runAddIgnore`, `runAddAttributes`) parameterize
// only the kind / extension / target file / marker tag prefix and reuse
// the shared `addTemplateOp` body. This keeps each entry point under
// the 15-line cap and avoids two diverging copies of the merge dance.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/templates"
)

// addTemplateFlags holds parsed flags shared by ignore + attributes.
type addTemplateFlags struct {
	dryRun bool
}

// addTemplateSpec parameterizes the kind-specific differences for
// `add ignore` vs `add attributes`. Same struct shape as commitTransferSpec
// elsewhere in cmd/ — keeps dispatch readable.
type addTemplateSpec struct {
	kind        string // "ignore" | "attributes"
	subcommand  string // "ignore" | "attributes" (CLI label)
	targetName  string // ".gitignore" | ".gitattributes"
	bannerLabel string // human-friendly banner suffix
}

// runAddIgnore handles `gitmap add ignore [langs...]`. Always merges
// `common` first so OS junk + IDE noise are guaranteed to land regardless
// of which languages the user picked.
func runAddIgnore(args []string) {
	checkHelp("add-ignore", args)
	executeAddTemplate(addTemplateSpec{
		kind:        "ignore",
		subcommand:  "ignore",
		targetName:  ".gitignore",
		bannerLabel: "merge curated .gitignore template block",
	}, args)
}

// runAddAttributes handles `gitmap add attributes [langs...]`. Same
// pipeline as runAddIgnore, just a different file + extension.
func runAddAttributes(args []string) {
	checkHelp("add-attributes", args)
	executeAddTemplate(addTemplateSpec{
		kind:        "attributes",
		subcommand:  "attributes",
		targetName:  ".gitattributes",
		bannerLabel: "merge curated .gitattributes template block",
	}, args)
}

// executeAddTemplate is the shared pipeline. Kept to a single high-level
// flow so the per-step helpers below can stay small and self-documenting.
func executeAddTemplate(spec addTemplateSpec, args []string) {
	flags, langs := parseAddTemplateArgs(spec, args)
	if !insideGitRepo() {
		fmt.Fprintln(os.Stderr, "  ✗ Not inside a Git repository.")
		fmt.Fprintln(os.Stderr, "    Run this from the root of a repo (where .git/ lives).")
		os.Exit(1)
	}

	resolved, err := resolveAddTemplates(spec.kind, langs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	target, err := repoFilePath(spec.targetName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Could not locate repo root: %v\n", err)
		os.Exit(1)
	}

	tag := buildAddTag(spec.kind, langs)
	body := dedupeLines(concatTemplateBodies(resolved))

	printAddTemplateBanner(spec, flags.dryRun, resolved, tag)
	if flags.dryRun {
		printAddTemplateDryRun(target, tag, body)

		return
	}

	res, mergeErr := templates.Merge(target, tag, body)
	if mergeErr != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Could not merge into %s: %v\n", target, mergeErr)
		os.Exit(1)
	}
	printAddTemplateSummary(spec, res)
}

// parseAddTemplateArgs separates `--dry-run` / `-h` from positional langs.
// Empty langs is allowed — the pipeline still merges `common` alone.
func parseAddTemplateArgs(spec addTemplateSpec, args []string) (addTemplateFlags, []string) {
	fs := flag.NewFlagSet("add "+spec.subcommand, flag.ExitOnError)
	dryRun := fs.Bool("dry-run", false, "preview the merged "+spec.targetName+" block without writing anything")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Could not parse flags: %v\n", err)
		os.Exit(1)
	}

	return addTemplateFlags{dryRun: *dryRun}, normalizeLangs(fs.Args())
}

// normalizeLangs lowercases, trims, and de-duplicates the language list
// while preserving the user's first-seen order. `common` is always the
// first language merged so the helper strips any explicit `common`
// passed by the user to avoid double-prepending.
func normalizeLangs(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, raw := range in {
		lang := strings.ToLower(strings.TrimSpace(raw))
		if lang == "" || lang == "common" {
			continue
		}
		if _, dup := seen[lang]; dup {
			continue
		}
		seen[lang] = struct{}{}
		out = append(out, lang)
	}

	return out
}

// resolveAddTemplates loads `common` first, then each requested lang, in
// stable order. A missing template aborts with a clear error so the user
// knows exactly which lang typo caused the failure.
func resolveAddTemplates(kind string, langs []string) ([]templates.Resolved, error) {
	all := append([]string{"common"}, langs...)
	out := make([]templates.Resolved, 0, len(all))
	for _, lang := range all {
		r, err := templates.Resolve(kind, lang)
		if err != nil {
			return nil, fmt.Errorf("template %s/%s: %w", kind, lang, err)
		}
		out = append(out, r)
	}

	return out, nil
}

// concatTemplateBodies stitches every resolved template body together
// with a one-line `# ── <lang> ──` separator so the marker block stays
// human-readable when multiple languages are merged.
func concatTemplateBodies(resolved []templates.Resolved) []byte {
	var b strings.Builder
	for i, r := range resolved {
		if i > 0 {
			b.WriteString("\n")
		}
		fmt.Fprintf(&b, "# ── %s ──\n", r.Lang)
		b.Write(r.Content)
		if len(r.Content) > 0 && r.Content[len(r.Content)-1] != '\n' {
			b.WriteByte('\n')
		}
	}

	return []byte(b.String())
}

// dedupeLines collapses repeated non-blank lines while preserving order
// and keeping every blank line as-is (blank lines are visual separators,
// not duplicates). Comment lines ARE deduped — repeated `# OS junk`
// section headers across templates would otherwise pile up.
func dedupeLines(body []byte) []byte {
	seen := map[string]struct{}{}
	lines := strings.Split(string(body), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			out = append(out, line)

			continue
		}
		if _, dup := seen[trimmed]; dup {
			// Keep position as a blank so visual spacing between
			// retained sections is preserved instead of collapsing.
			out = append(out, "")

			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, line)
	}
	// Collapse runs of 3+ blank lines down to 2 — two blanks is the
	// widest separator we ever want to render after dedupe.
	out = collapseBlankRuns(out, 2)

	return []byte(strings.Join(out, "\n"))
}

// collapseBlankRuns trims any INTERIOR run of blank lines longer than
// maxBlank down to exactly maxBlank consecutive blanks. A run that
// extends to the end of the slice is preserved as-is so trailing
// spacing introduced by dedup'd duplicates survives untouched.
func collapseBlankRuns(lines []string, maxBlank int) []string {
	// Find the index of the last non-blank line. Anything after it is
	// the trailing run and is left alone.
	lastNonBlank := -1
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			lastNonBlank = i
		}
	}
	out := make([]string, 0, len(lines))
	run := 0
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			run++
			if i > lastNonBlank || run <= maxBlank {
				out = append(out, line)
			}

			continue
		}
		run = 0
		out = append(out, line)
	}

	return out
}

// buildAddTag turns a language list into a stable, sorted tag suffix so
// "go,node" and "node,go" share a marker block. `common` is implicit.
func buildAddTag(kind string, langs []string) string {
	if len(langs) == 0 {
		return kind + "/common"
	}
	sorted := make([]string, len(langs))
	copy(sorted, langs)
	sort.Strings(sorted)

	return kind + "/" + strings.Join(sorted, "+")
}

// repoFilePath returns the absolute path of `name` at the repo top
// level. Falls back to CWD when `git rev-parse` fails — the upstream
// insideGitRepo check makes that path defensive only.
func repoFilePath(name string) (string, error) {
	root, err := gitTopLevel()
	if err != nil || strings.TrimSpace(root) == "" {
		cwd, cwdErr := os.Getwd()
		if cwdErr != nil {
			return "", cwdErr
		}

		return filepath.Join(cwd, name), nil
	}

	return filepath.Join(root, name), nil
}

// printAddTemplateBanner mirrors the visual idiom of `add lfs-install`.
func printAddTemplateBanner(spec addTemplateSpec, dryRun bool, resolved []templates.Resolved, tag string) {
	fmt.Println()
	fmt.Printf("  %s■ gitmap add %s —%s %s\n",
		constants.ColorCyan, spec.subcommand, constants.ColorReset, spec.bannerLabel)
	for _, r := range resolved {
		src := "embed"
		if r.Source == templates.SourceUser {
			src = "user"
		}
		fmt.Printf("    %s%s/%s%s  source=%s  (%s)\n",
			constants.ColorDim, r.Kind, r.Lang, constants.ColorReset, src, r.Path)
	}
	fmt.Printf("  block tag: %s%s%s\n", constants.ColorDim, tag, constants.ColorReset)
	if dryRun {
		fmt.Printf("  %s[dry-run]%s no files will be modified\n",
			constants.ColorYellow, constants.ColorReset)
	}
	fmt.Println()
}

// printAddTemplateDryRun shows the marker block exactly as it would
// land on disk. Mirrors printAddLFSInstallDryRun's framing.
func printAddTemplateDryRun(target, tag string, body []byte) {
	fmt.Printf("  %swould write block into:%s %s\n",
		constants.ColorYellow, constants.ColorReset, target)
	fmt.Println()
	fmt.Printf("# >>> gitmap:%s >>>\n", tag)
	os.Stdout.Write(body)
	if len(body) == 0 || body[len(body)-1] != '\n' {
		fmt.Println()
	}
	fmt.Printf("# <<< gitmap:%s <<<\n", tag)
	fmt.Println()
}

// printAddTemplateSummary prints the merge outcome + a `git add` hint.
func printAddTemplateSummary(spec addTemplateSpec, res templates.MergeResult) {
	verb := outcomeVerb(res.Outcome)
	color := constants.ColorGreen
	if !res.Changed {
		verb = "unchanged"
		color = constants.ColorDim
	}
	fmt.Printf("  %s%s%s %s (block: %s)\n",
		color, verb, constants.ColorReset, res.Path, res.BlockTag)
	if res.Changed {
		fmt.Println()
		fmt.Printf("  %sNext step:%s commit the updated %s:\n",
			constants.ColorYellow, constants.ColorReset, spec.targetName)
		fmt.Printf("    git add %s\n", spec.targetName)
		fmt.Printf("    git commit -m \"chore: refresh %s via gitmap template\"\n", spec.targetName)
	}
	fmt.Println()
}
