package cmdsync

import (
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/templates"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// runAddLFSInstall is the entry point dispatched from rootadd.go.
func runAddLFSInstall(args []string) *apperror.AppError {
	checkHelp("add-lfs-install", args)

	flags, err := parseAddLFSInstallFlags(args)
	if err != nil {
		return err
	}

	if !insideGitRepo() {
		fmt.Fprintln(os.Stderr, "  ✗ Not inside a Git repository.")
		fmt.Fprintln(os.Stderr, "    Run this command from the root of a repo (where .git/ lives).")
		cliexit.HandleError(nil, 1)
	}

	if !lfsAvailable() {
		fmt.Fprintln(os.Stderr, "  ✗ Git LFS is not installed or not on PATH.")
		fmt.Fprintln(os.Stderr, "    Install it from https://git-lfs.com and re-run.")
		cliexit.HandleError(nil, 1)
	}

	resolved, err2 := templates.Resolve("lfs", "common")
	if err2 != nil {
		return apperror.WrapSimple(err2, "? Could not resolve lfs/common template")
	}

	target, err3 := gitattributesPath()
	if err3 != nil {
		return apperror.WrapSimple(err3, "? Could not locate repo root")
	}

	printAddLFSInstallBanner(flags.dryRun, resolved.Source, resolved.Path)

	if flags.dryRun {
		printAddLFSInstallDryRun(target, resolved.Content)

		return nil
	}

	if err := runGitLFSInstall(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ git lfs install failed: %v\n", err)
		// Continue anyway — the template merge is independently useful.
	} else {
		fmt.Printf("  %s✓ git lfs install --local%s ran successfully\n",
			constants.ColorGreen, constants.ColorReset)
	}

	res, err4 := templates.Merge(target, addLFSInstallTag, resolved.Content)
	if err4 != nil {
		return apperror.WrapSimple(err4, "? Could not merge template")
	}

	printAddLFSInstallSummary(res)

	return nil
}

// addLFSInstallTag identifies the marker block written into
// .gitattributes. Stable across runs — changing it would orphan blocks
// already on disk.
const addLFSInstallTag = "lfs/common"

// addLFSInstallFlags holds parsed flags for `add lfs-install`.
type addLFSInstallFlags struct {
	dryRun bool
}

// parseAddLFSInstallFlags parses CLI flags. Currently only --dry-run.
func parseAddLFSInstallFlags(args []string) (addLFSInstallFlags, *apperror.AppError) {
	fs := flag.NewFlagSet("add lfs-install", flag.ExitOnError)
	dryRun := fs.Bool("dry-run", false, "preview the merged .gitattributes without writing anything")
	if err := fs.Parse(args); err != nil {
		return addLFSInstallFlags{}, apperror.WrapSimple(err, "✗ Could not parse flags")
	}

	return addLFSInstallFlags{dryRun: *dryRun}, nil
}

// insideGitRepo returns true when the current working directory is part
// of a Git working tree.
func insideGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}

// lfsAvailable returns true when `git lfs` is installed and runnable.
func lfsAvailable() bool {
	cmd := exec.Command("git", "lfs", "version")

	return cmd.Run() == nil
}

// gitTopLevel returns the absolute path of the current repo's top-level dir.
func gitTopLevel() (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitRevParse, "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("empty top-level")
	}

	return filepath.Clean(root), nil
}

// gitattributesPath resolves the absolute path to .gitattributes at the
// current repo's top level. Falls back to CWD when `git rev-parse` fails
// for any reason — but the insideGitRepo check upstream already gates
// that case, so this is defensive.
func gitattributesPath() (string, error) {
	root, err := gitTopLevel()
	if err == nil && len(strings.TrimSpace(root)) > 0 {
		return filepath.Join(root, ".gitattributes"), nil
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return "", cwdErr
	}

	return filepath.Join(cwd, ".gitattributes"), nil
}

// printAddLFSInstallBanner mirrors the `lfs-common` banner so the two
// commands feel like siblings.
func printAddLFSInstallBanner(dryRun bool, source templates.SourceType, path string) {
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

// runGitLFSInstall runs `git lfs install` for the current repo. It is
// safe to call repeatedly — Git LFS treats it as idempotent.
func runGitLFSInstall() error {
	cmd := exec.Command("git", "lfs", "install", "--local")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
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
func outcomeVerb(o templates.MergeOutcomeType) string {
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
