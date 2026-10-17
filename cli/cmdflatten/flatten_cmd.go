package cmdflatten

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// FlattenOptions captures CLI arguments and flags for repo flattening.
type FlattenOptions struct {
	Target   string
	IsDryRun bool
	IsYes    bool
}

// RunFlattenCmd executes `gitmap repo flatten` and `gitmap flatten`.
func RunFlattenCmd(args []string) error {
	opts, err := parseFlattenArgs(args)
	if err != nil {
		return err
	}

	absTarget, err := filepath.Abs(opts.Target)
	if err != nil {
		return fmt.Errorf("resolve target path %q: %w", opts.Target, err)
	}

	targetInfo, err := os.Stat(absTarget)
	if err != nil || !targetInfo.IsDir() {
		return fmt.Errorf("target repository %q does not exist or is not a directory", absTarget)
	}

	gitDir := filepath.Join(absTarget, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		return fmt.Errorf("target %q is not a git repository (missing .git)", absTarget)
	}

	parentDir := filepath.Dir(absTarget)
	repoName := filepath.Base(absTarget)
	successorPath, successorName := resolveSuccessorPath(parentDir, repoName)
	archivePath := filepath.Join(parentDir, repoName+"-marked-for-deletion")

	commitCount := countCommits(absTarget)
	workingTreeBytes, err := computeWorkingTreeSize(absTarget)
	if err != nil {
		return fmt.Errorf("calculate working tree size: %w", err)
	}

	if opts.IsDryRun {
		printDryRun(absTarget, commitCount, workingTreeBytes, successorName, archivePath)

		return nil
	}

	if !opts.IsYes {
		confirmed, err := promptConfirmation()
		if err != nil {
			return err
		}

		if !confirmed {
			fmt.Println("Flatten operation aborted.")

			return nil
		}
	}

	if err := os.MkdirAll(successorPath, 0755); err != nil {
		return fmt.Errorf("create successor directory %q: %w", successorPath, err)
	}

	fmt.Printf("Copying working tree files to %s...\n", successorName)
	if err := copyWorkingTree(absTarget, successorPath); err != nil {
		_ = os.RemoveAll(successorPath)

		return fmt.Errorf("copy working tree: %w", err)
	}

	if err := initGitRepo(successorPath); err != nil {
		_ = os.RemoveAll(successorPath)

		return fmt.Errorf("initialize git repository: %w", err)
	}

	if err := commitFlatSnapshot(successorPath); err != nil {
		_ = os.RemoveAll(successorPath)

		return fmt.Errorf("commit flat snapshot: %w", err)
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

		return fmt.Errorf("archive original repository: %w", err)
	}

	fmt.Printf("Promoting %s to %s...\n", successorName, repoName)
	if err := os.Rename(successorPath, absTarget); err != nil {
		_ = os.Rename(archivePath, absTarget)

		return fmt.Errorf("promote successor directory: %w", err)
	}

	printSuccessSummary(repoName, absTarget, archivePath)

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
				return opts, fmt.Errorf("missing value for --target flag")
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

func resolveSuccessorPath(parentDir, repoName string) (string, string) {
	suffix := 1
	for {
		name := fmt.Sprintf("%s-flat-v%d", repoName, suffix)
		path := filepath.Join(parentDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return path, name
		}

		suffix++
	}
}

func countCommits(repoPath string) int {
	cmd := exec.Command("git", "-C", repoPath, "rev-list", "--count", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}

	count, parseErr := strconv.Atoi(strings.TrimSpace(string(out)))
	if parseErr != nil {
		return 0
	}

	return count
}

func computeWorkingTreeSize(root string) (int64, error) {
	var totalBytes int64

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}

		if info.IsDir() {
			if isExcludedDir(info.Name()) && path != root {
				return filepath.SkipDir
			}

			return nil
		}

		totalBytes += info.Size()

		return nil
	})

	return totalBytes, err
}

func isExcludedDir(name string) bool {
	switch strings.ToLower(name) {
	case ".git", "node_modules", "dist", "temp", ".pnpm-store":
		return true

	default:
		return false
	}
}

func copyWorkingTree(srcRoot, dstRoot string) error {
	return filepath.Walk(srcRoot, func(srcPath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if info.IsDir() {
			if isExcludedDir(info.Name()) && srcPath != srcRoot {
				return filepath.SkipDir
			}

			return nil
		}

		relPath, err := filepath.Rel(srcRoot, srcPath)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstRoot, relPath)
		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		return copyFile(srcPath, dstPath, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return nil
}

func initGitRepo(repoPath string) error {
	cmdInit := exec.Command("git", "init", "-b", "main")
	cmdInit.Dir = repoPath
	if err := cmdInit.Run(); err != nil {
		fallbackCmd := exec.Command("git", "init")
		fallbackCmd.Dir = repoPath
		if fallbackErr := fallbackCmd.Run(); fallbackErr != nil {
			return fallbackErr
		}
	}

	return nil
}

func commitFlatSnapshot(repoPath string) error {
	cmdAdd := exec.Command("git", "add", "-A")
	cmdAdd.Dir = repoPath
	if err := cmdAdd.Run(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	cmdCommit := exec.Command("git", "-c", "user.name=GitMap", "-c", "user.email=gitmap@local",
		"commit", "-m", "flat: initial squashed repository snapshot")
	cmdCommit.Dir = repoPath
	if err := cmdCommit.Run(); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	return nil
}

func isInsideOrSame(child, parent string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	return !strings.HasPrefix(rel, "..")
}

func promptConfirmation() (bool, error) {
	fmt.Println(constants.ColorYellow + "WARNING: This operation squashes repository history into a single flat commit. It is irreversible unless restored via undo." + constants.ColorReset)
	fmt.Print("Proceed with repository flattening? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("read confirmation: %w", err)
	}

	clean := strings.ToLower(strings.TrimSpace(response))

	return clean == "y" || clean == "yes", nil
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
