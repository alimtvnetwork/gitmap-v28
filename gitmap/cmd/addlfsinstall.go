// Command `gitmap add lfs-install`.
//
// Two-step operation:
//  1. Run `git lfs install --local` so the current repo has the LFS
//     pre-push / clean / smudge hooks wired up. This is itself
//     idempotent — Git LFS treats repeat invocations as no-ops.
//  2. Resolve the `lfs/common` template (overlay > embed) and merge its
//     body into ./.gitattributes via templates.Merge, which uses a
//     gitmap-managed marker block so the second and later runs are
//     byte-stable no-ops when the template hasn't changed.
//
// Why a separate command from `gitmap lfs-common`? `lfs-common` shells
// out to `git lfs track` per-pattern and writes whatever line format
// `git lfs` decides on. `add lfs-install` is the template-driven path:
// the bytes written to .gitattributes come from the curated, versioned
// `lfs/common.gitattributes` asset (audit-trailed `# source:` header)
// and can be overridden by a user file at ~/.gitmap/templates/lfs/common.gitattributes.
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

// addLFSInstallTag identifies the marker block written into
// .gitattributes. Stable across runs — changing it would orphan blocks
// already on disk.
const addLFSInstallTag = "lfs/common"

// addLFSInstallFlags holds parsed flags for `add lfs-install`.
type addLFSInstallFlags struct {
	dryRun bool
}

// runAddLFSInstall is the entry point dispatched from rootadd.go.
func runAddLFSInstall(args []string) {
	checkHelp("add-lfs-install", args)

	flags := parseAddLFSInstallFlags(args)

	if !insideGitRepo() {
		fmt.Fprintln(os.Stderr, "  ✗ Not inside a Git repository.")
		fmt.Fprintln(os.Stderr, "    Run this command from the root of a repo (where .git/ lives).")
		os.Exit(1)
	}
	if !lfsAvailable() {
		fmt.Fprintln(os.Stderr, "  ✗ Git LFS is not installed or not on PATH.")
		fmt.Fprintln(os.Stderr, "    Install it from https://git-lfs.com and re-run.")
		os.Exit(1)
	}

	resolved, err := templates.Resolve("lfs", "common")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Could not resolve lfs/common template: %v\n", err)
		os.Exit(1)
	}

	target, err := gitattributesPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Could not locate repo root: %v\n", err)
		os.Exit(1)
	}

	printAddLFSInstallBanner(flags.dryRun, resolved.Source, resolved.Path)

	if flags.dryRun {
		printAddLFSInstallDryRun(target, resolved.Content)

		return
	}

	if err := runGitLFSInstall(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ git lfs install failed: %v\n", err)
		// Continue anyway — the template merge is independently useful.
	} else {
		fmt.Printf("  %s✓ git lfs install --local%s ran successfully\n",
			constants.ColorGreen, constants.ColorReset)
	}

	res, err := templates.Merge(target, addLFSInstallTag, resolved.Content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Could not merge template into %s: %v\n", target, err)
		os.Exit(1)
	}

	printAddLFSInstallSummary(res)
}

// parseAddLFSInstallFlags parses CLI flags. Currently only --dry-run.
func parseAddLFSInstallFlags(args []string) addLFSInstallFlags {
	fs := flag.NewFlagSet("add lfs-install", flag.ExitOnError)
	dryRun := fs.Bool("dry-run", false, "preview the merged .gitattributes without writing anything")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Could not parse flags: %v\n", err)
		os.Exit(1)
	}

	return addLFSInstallFlags{dryRun: *dryRun}
}

// gitattributesPath resolves the absolute path to .gitattributes at the
// current repo's top level. Falls back to CWD when `git rev-parse` fails
// for any reason — but the insideGitRepo check upstream already gates
// that case, so this is defensive.
func gitattributesPath() (string, error) {
	root, err := gitTopLevel()
	if err != nil || len(strings.TrimSpace(root)) == 0 {
		cwd, cwdErr := os.Getwd()
		if cwdErr != nil {
			return "", cwdErr
		}

		return filepath.Join(cwd, ".gitattributes"), nil
	}

	return filepath.Join(root, ".gitattributes"), nil
}

// printAddLFSInstallBanner mirrors the `lfs-common` banner so the two
// commands feel like siblings.
func printAddLFSInstallBanner(dryRun bool, source templates.Source, path string) {
	fmt.Println()
	fmt.Printf("  %s■ gitmap add lfs-install —%s LFS hooks + templated .gitattributes block\n",
		constants.ColorCyan, constants.ColorReset)
	src := "embed"
	if source == templates.SourceUser {
		src = "user"
	}
	fmt.Printf("  template source: %s%s%s (%s)\n",
		constants.ColorDim, src, constants.ColorReset, path)
	if dryRun {
		fmt.Printf("  %s[dry-run]%s no files will be modified\n",
			constants.ColorYellow, constants.ColorReset)
	}
	fmt.Println()
}

// printAddLFSInstallDryRun shows the block that would be written, with
// markers, but does not touch disk.
func printAddLFSInstallDryRun(target string, body []byte) {
	fmt.Printf("  %swould write block into:%s %s\n",
		constants.ColorYellow, constants.ColorReset, target)
	fmt.Println()
	fmt.Printf("# >>> gitmap:%s >>>\n", addLFSInstallTag)
	os.Stdout.Write(body)
	if len(body) == 0 || body[len(body)-1] != '\n' {
		fmt.Println()
	}
	fmt.Printf("# <<< gitmap:%s <<<\n", addLFSInstallTag)
	fmt.Println()
}

// printAddLFSInstallSummary renders the merge outcome in the same visual
// idiom as lfs-common's summary block.
func printAddLFSInstallSummary(res templates.MergeResult) {
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
		fmt.Printf("  %sNext step:%s commit the updated .gitattributes:\n",
			constants.ColorYellow, constants.ColorReset)
		fmt.Println("    git add .gitattributes")
		fmt.Println("    git commit -m \"chore: install Git LFS + track common binaries via gitmap template\"")
	}
	fmt.Println()
}

// outcomeVerb maps the structural outcome to a one-word verb for output.
func outcomeVerb(o templates.MergeOutcome) string {
	switch o {
	case templates.MergeCreated:
		return "created"
	case templates.MergeInserted:
		return "inserted block into"
	case templates.MergeUpdated:
		return "updated block in"
	default:
		return "wrote"
	}
}
