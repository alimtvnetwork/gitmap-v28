package cmdsync

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/constants"
	globalconstants "github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/templates"
)

// runAddLFSInstall is the entry point dispatched from rootadd.go.
func runAddLFSInstall(args []string) *apperror.AppError {
	checkHelp(constants.CmdAddLFSInstall, args)

	flags, err := parseAddLFSInstallFlags(args)
	if err != nil {
		return err
	}

	if envErr := validateGitLFSEnvironment(); envErr != nil {
		return envErr
	}

	return executeLFSInstall(flags)
}

func validateGitLFSEnvironment() *apperror.AppError {
	if !insideGitRepo() {
		fmt.Fprintln(os.Stderr, "  ✗ Not inside a Git repository.")
		fmt.Fprintln(os.Stderr, "    Run this command from the root of a repo (where .git/ lives).")

		return apperror.NewValidationError("not inside a Git repository")
	}

	if !lfsAvailable() {
		fmt.Fprintln(os.Stderr, "  ✗ Git LFS is not installed or not on PATH.")
		fmt.Fprintln(os.Stderr, "    Install it from https://git-lfs.com and re-run.")

		return apperror.NewExecutionError("Git LFS is not installed or not on PATH")
	}

	return nil
}

func executeLFSInstall(flags addLFSInstallFlags) *apperror.AppError {
	resolved, err := templates.Resolve("lfs", "common")
	if err != nil {
		return apperror.WrapSimple(err, "? Could not resolve lfs/common template")
	}

	target, err2 := gitattributesPath()
	if err2 != nil {
		return apperror.WrapSimple(err2, "? Could not locate repo root")
	}

	printAddLFSInstallBanner(flags.dryRun, resolved.Source, resolved.Path)
	if flags.dryRun {
		printAddLFSInstallDryRun(target, resolved.Content)

		return nil
	}

	return applyLFSInstall(target, resolved.Content)
}

func applyLFSInstall(target string, content []byte) *apperror.AppError {
	if err := runGitLFSInstall(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ git lfs install failed: %v\n", err)
	} else {
		fmt.Printf("  %s✓ git lfs install --local%s ran successfully\n",
			globalconstants.ColorGreen, globalconstants.ColorReset)
	}

	res, mergeErr := templates.Merge(target, constants.AddLFSInstallTag, content)
	if mergeErr != nil {
		return apperror.WrapSimple(mergeErr, "? Could not merge template")
	}

	printAddLFSInstallSummary(res)

	return nil
}

type addLFSInstallFlags struct {
	dryRun bool
}

func parseAddLFSInstallFlags(args []string) (addLFSInstallFlags, *apperror.AppError) {
	fs := flag.NewFlagSet(constants.CmdAddLFSInstall, flag.ExitOnError)
	dryRun := fs.Bool(constants.FlagDryRun, false, "preview the merged .gitattributes without writing anything")
	if err := fs.Parse(args); err != nil {
		return addLFSInstallFlags{}, apperror.WrapSimple(err, "✗ Could not parse flags")
	}

	return addLFSInstallFlags{dryRun: *dryRun}, nil
}

func insideGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}

func lfsAvailable() bool {
	cmd := exec.Command("git", "lfs", "version")

	return cmd.Run() == nil
}

func gitTopLevel() (string, error) {
	cmd := exec.Command(globalconstants.GitBin, globalconstants.GitRevParse, "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", apperror.NewExecutionError("empty top-level")
	}

	return filepath.Clean(root), nil
}

func gitattributesPath() (string, error) {
	root, err := gitTopLevel()
	if err == nil && len(strings.TrimSpace(root)) > 0 {
		return filepath.Join(root, constants.FileGitattributes), nil
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return "", cwdErr
	}

	return filepath.Join(cwd, constants.FileGitattributes), nil
}

func printAddLFSInstallBanner(dryRun bool, source templates.SourceType, path string) {
	fmt.Println()
	fmt.Printf("  %s■ gitmap add lfs-install —%s LFS hooks + templated .gitattributes block\n",
		globalconstants.ColorCyan, globalconstants.ColorReset)
	src := "embed"
	if source == templates.SourceUser {
		src = "user"
	}

	fmt.Printf("  template source: %s%s%s (%s)\n",
		globalconstants.ColorDim, src, globalconstants.ColorReset, path)
	if dryRun {
		fmt.Printf("  %s[dry-run]%s no files will be modified\n",
			globalconstants.ColorYellow, globalconstants.ColorReset)
	}

	fmt.Println()
}

func printAddLFSInstallDryRun(target string, body []byte) {
	fmt.Printf("  %swould write block into:%s %s\n",
		globalconstants.ColorYellow, globalconstants.ColorReset, target)
	fmt.Println()
	fmt.Printf("# >>> gitmap:%s >>>\n", constants.AddLFSInstallTag)
	os.Stdout.Write(body)
	if len(body) == 0 || body[len(body)-1] != '\n' {
		fmt.Println()
	}

	fmt.Printf("# <<< gitmap:%s <<<\n", constants.AddLFSInstallTag)
	fmt.Println()
}

func runGitLFSInstall() error {
	cmd := exec.Command("git", "lfs", "install", constants.FlagLocal)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.WrapSimple(err, "git lfs install").WithContext("output", strings.TrimSpace(string(out)))
	}

	return nil
}

func printAddLFSInstallSummary(res templates.MergeResult) {
	verb := outcomeVerb(res.Outcome)
	color := globalconstants.ColorGreen
	if !res.Changed {
		verb = "unchanged"
		color = globalconstants.ColorDim
	}

	fmt.Printf("  %s%s%s %s (block: %s)\n",
		color, verb, globalconstants.ColorReset, res.Path, res.BlockTag)

	if res.Changed {
		printNextCommitStep()
	}

	fmt.Println()
}

func printNextCommitStep() {
	fmt.Println()
	fmt.Printf("  %sNext step:%s commit the updated .gitattributes:\n",
		globalconstants.ColorYellow, globalconstants.ColorReset)
	fmt.Println("    git add .gitattributes")
	fmt.Println("    git commit -m \"chore: install Git LFS + track common binaries via gitmap template\"")
}

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
