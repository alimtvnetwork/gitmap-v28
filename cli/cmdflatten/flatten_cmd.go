package cmdflatten

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
)

// FlattenOptions captures CLI arguments and flags for repo flattening.
type FlattenOptions struct {
	Target   string
	IsDryRun bool
	IsYes    bool
}

// RunFlattenCmd executes `gitmap repo flatten` and `gitmap flatten`.
func RunFlattenCmd(args []string) error {
	if hasHelpFlag(args) {
		helpdoc.Print("flatten")
		return nil
	}

	opts, err := parseFlattenArgs(args)
	if err != nil {
		return err
	}

	absTarget, err := filepath.Abs(opts.Target)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("resolve target path %q", opts.Target))
	}

	targetInfo, err := os.Stat(absTarget)
	if err != nil || !targetInfo.IsDir() {
		return apperror.NewNotFoundError(fmt.Sprintf("target repository %q does not exist or is not a directory", absTarget))
	}

	gitDir := filepath.Join(absTarget, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		return apperror.NewNotFoundError(fmt.Sprintf("target %q is not a git repository (missing .git)", absTarget))
	}

	parentDir := filepath.Dir(absTarget)
	repoName := filepath.Base(absTarget)
	successorPath, successorName := resolveSuccessorPath(parentDir, repoName)
	archivePath := filepath.Join(parentDir, repoName+"-marked-for-deletion")

	commitCount := countCommits(absTarget)
	workingTreeBytes, err := computeWorkingTreeSize(absTarget)
	if err != nil {
		return apperror.WrapSimple(err, "calculate working tree size")
	}

	if opts.IsDryRun {
		printDryRun(absTarget, commitCount, workingTreeBytes, successorName, archivePath)
		return nil
	}

	if confirmErr := ensureConfirmation(opts.IsYes); confirmErr != nil {
		fmt.Println("Flatten operation aborted.")
		return nil
	}

	if err := os.MkdirAll(successorPath, 0755); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("create successor directory %q", successorPath))
	}

	fmt.Printf("Copying working tree files to %s...\n", successorName)
	if err := copyWorkingTree(absTarget, successorPath); err != nil {
		_ = os.RemoveAll(successorPath)
		return apperror.WrapSimple(err, "copy working tree")
	}

	if err := initGitRepo(successorPath); err != nil {
		_ = os.RemoveAll(successorPath)
		return apperror.WrapSimple(err, "initialize git repository")
	}

	if err := commitFlatSnapshot(successorPath); err != nil {
		_ = os.RemoveAll(successorPath)
		return apperror.WrapSimple(err, "commit flat snapshot")
	}

	cwd, _ := os.Getwd()
	if isInsideOrSame(cwd, absTarget) {
		_ = os.Chdir(parentDir)
	}

	if _, err := os.Stat(archivePath); err == nil {
		_ = os.RemoveAll(archivePath)
	}

	fmt.Printf("Archiving bloated repository to %s...\n", filepath.Base(archivePath))
	if err := os.Rename(absTarget, archivePath); err != nil {
		_ = os.RemoveAll(successorPath)
		return apperror.WrapSimple(err, "archive original repository")
	}

	fmt.Printf("Promoting %s to %s...\n", successorName, repoName)
	if err := os.Rename(successorPath, absTarget); err != nil {
		_ = os.Rename(archivePath, absTarget)
		return apperror.WrapSimple(err, "promote successor directory")
	}

	printSuccessSummary(repoName, absTarget, archivePath)
	return nil
}

func ensureConfirmation(isYes bool) error {
	if isYes {
		return nil
	}
	confirmed, err := promptConfirmation()
	if err != nil {
		return err
	}
	if !confirmed {
		return apperror.NewExecutionError("user aborted")
	}
	return nil
}

func parseFlattenArgs(args []string) (FlattenOptions, error) {
	opts := FlattenOptions{
		Target: ".",
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--dry-run", arg == "-dry-run":
			opts.IsDryRun = true
		case arg == "-y", arg == "--yes", arg == "-yes":
			opts.IsYes = true
		case arg == "--target", arg == "-target":
			if i+1 >= len(args) {
				return opts, apperror.NewValidationError("missing value for --target flag")
			}
			opts.Target = args[i+1]
			i++
		case strings.HasPrefix(arg, "--target="):
			opts.Target = strings.TrimPrefix(arg, "--target=")
		case !strings.HasPrefix(arg, "-") && opts.Target == ".":
			opts.Target = arg
		}
	}

	return opts, nil
}

func isInsideOrSame(child, parent string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}

var flattenConfirmTokens = []string{"y", "yes"}

func promptConfirmation() (bool, error) {
	fmt.Println(constants.ColorYellow + "WARNING: This operation squashes repository history into a single flat commit. It is irreversible unless restored via undo." + constants.ColorReset)
	fmt.Print("Proceed with repository flattening? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, apperror.WrapSimple(err, "read confirmation")
	}

	clean := strings.ToLower(strings.TrimSpace(response))
	return slices.Contains(flattenConfirmTokens, clean), nil
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.2f %s", float64(bytes)/float64(div), units[exp])
}

func printDryRun(repoPath string, commits int, sizeBytes int64, successorName, archivePath string) {
	fmt.Printf("\n%s[GitMap Repo Flatten — Dry Run Preview]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Target Repository:      %s\n", repoPath)
	fmt.Printf("  Commit Count:           %d\n", commits)
	fmt.Printf("  Working Tree Size:      %s\n", formatBytes(sizeBytes))
	fmt.Printf("  Successor Staging Dir:  %s\n", successorName)
	fmt.Printf("  Archive Destination:    %s\n", filepath.Base(archivePath))
	fmt.Printf("\n%sDry-run mode active. No files were modified.%s\n\n", constants.ColorYellow, constants.ColorReset)
}

func printSuccessSummary(repoName, absTarget, archivePath string) {
	fmt.Printf("\n%s✓ Repository %s squashed into a single flat commit successfully!%s\n", constants.ColorGreen, repoName, constants.ColorReset)
	fmt.Printf("  Canonical Working Tree: %s\n", absTarget)
	fmt.Printf("  Archived Bloated Repo:  %s\n", archivePath)
	fmt.Printf("\n  %sTo undo: gitmap repo undo-flatten %s%s\n\n", constants.ColorCyan, repoName, constants.ColorReset)
}

var flattenHelpFlagTokens = []string{"-h", "--help", "help"}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if slices.Contains(flattenHelpFlagTokens, arg) {
			return true
		}
	}
	return false
}
