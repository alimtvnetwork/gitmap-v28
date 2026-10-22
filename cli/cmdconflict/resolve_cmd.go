package cmdconflict

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
)

func parseResolveFlags(args []string) (target string, strategy string) {
	for _, arg := range args {
		low := strings.ToLower(strings.TrimSpace(arg))
		switch low {
		case "--theirs", "-t", "theirs":
			strategy = "theirs"
		case "--ours", "-o", "ours":
			strategy = "ours"
		case "--ide", "-i", "--code", "ide", "code", "resolve", "v":
			strategy = "ide"
		case "--abort", "-a", "abort":
			strategy = "abort"
		case "--rerere", "rerere":
			strategy = "rerere"
		default:
			if !strings.HasPrefix(arg, "-") && target == "" {
				target = arg
			}
		}
	}
	return target, strategy
}

// ResolveConflictAction performs a targeted resolution action on repoPath.
// Returns true if the action was recognized and executed.
func ResolveConflictAction(repoPath string, action string) (bool, error) {
	low := strings.ToLower(strings.TrimSpace(action))
	switch low {
	case "theirs", "t":
		return true, applyTheirs(repoPath)
	case "ours", "o":
		return true, applyOurs(repoPath)
	case "resolve", "ide", "code", "v":
		return true, applyIDE(repoPath)
	case "abort", "a":
		return true, applyAbort(repoPath)
	case "rerere":
		return true, applyRerere(repoPath)
	}
	return false, nil
}

func applyTheirs(repoPath string) error {
	files, err := GetConflictedFiles(repoPath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Printf("%s No conflicted files detected in %s\n", constants.ColorYellow+"⚠"+constants.ColorReset, repoPath)
		return nil
	}
	for _, f := range files {
		_ = exec.Command("git", "-C", repoPath, "checkout", "--theirs", "--", f).Run()
		_ = exec.Command("git", "-C", repoPath, "add", "--", f).Run()
	}
	fmt.Printf("%s Accepted incoming (theirs) changes on %d conflicted file(s) in %s\n",
		constants.ColorGreen+"✓"+constants.ColorReset, len(files), repoPath)
	fmt.Println("  Next: Run 'gitmap cpf \"resolve merge conflicts\"' to commit and push.")
	return nil
}

func applyOurs(repoPath string) error {
	files, err := GetConflictedFiles(repoPath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Printf("%s No conflicted files detected in %s\n", constants.ColorYellow+"⚠"+constants.ColorReset, repoPath)
		return nil
	}
	for _, f := range files {
		_ = exec.Command("git", "-C", repoPath, "checkout", "--ours", "--", f).Run()
		_ = exec.Command("git", "-C", repoPath, "add", "--", f).Run()
	}
	fmt.Printf("%s Accepted current (ours) changes on %d conflicted file(s) in %s\n",
		constants.ColorGreen+"✓"+constants.ColorReset, len(files), repoPath)
	fmt.Println("  Next: Run 'gitmap cpf \"resolve merge conflicts\"' to commit and push.")
	return nil
}

func applyIDE(repoPath string) error {
	files, _ := GetConflictedFiles(repoPath)
	ideErr := LaunchIDEForRepo(repoPath, files)
	if ideErr != nil {
		fmt.Printf("%s Could not launch IDE: %v\n", constants.ColorYellow+"⚠"+constants.ColorReset, ideErr)
	} else {
		_, name := DetectInstalledIDE()
		fmt.Printf("%s Launched %s for %s\n", constants.ColorGreen+"✓"+constants.ColorReset, name, repoPath)
	}
	fmt.Print(FormatMergeEditorGuide(repoPath, files))
	return nil
}

func applyAbort(repoPath string) error {
	mergeCmd := exec.Command("git", "-C", repoPath, "merge", "--abort")
	if mergeCmd.Run() == nil {
		fmt.Printf("%s Merge aborted successfully in %s\n", constants.ColorGreen+"✓"+constants.ColorReset, repoPath)
		return nil
	}
	rebaseCmd := exec.Command("git", "-C", repoPath, "rebase", "--abort")
	if rebaseCmd.Run() == nil {
		fmt.Printf("%s Rebase aborted successfully in %s\n", constants.ColorGreen+"✓"+constants.ColorReset, repoPath)
		return nil
	}
	cherryCmd := exec.Command("git", "-C", repoPath, "cherry-pick", "--abort")
	if cherryCmd.Run() == nil {
		fmt.Printf("%s Cherry-pick aborted successfully in %s\n", constants.ColorGreen+"✓"+constants.ColorReset, repoPath)
		return nil
	}
	return apperror.NewExecutionError(fmt.Sprintf("no merge, rebase, or cherry-pick in progress to abort in %s", repoPath))
}

func applyRerere(repoPath string) error {
	_ = exec.Command("git", "-C", repoPath, "config", "rerere.enabled", "true").Run()
	_ = exec.Command("git", "-C", repoPath, "config", "rerere.autoupdate", "true").Run()
	fmt.Printf("%s Git rerere (Reuse Recorded Resolution) enabled for %s\n",
		constants.ColorGreen+"✓"+constants.ColorReset, repoPath)
	return nil
}

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" || a == "help" {
			return true
		}
	}
	return false
}

// RunResolveCmd implements 'gitmap resolve <repo> [--theirs | --ours | --ide | --abort]'.
func RunResolveCmd(args []string) error {
	if hasHelpFlag(args) {
		helpdoc.Print("resolve")
		return nil
	}

	target, strategy := parseResolveFlags(args)
	repoPath := ResolveRepoPath(target)
	if repoPath == "" {
		return apperror.NewValidationError("not inside a git repository and no repository specified.\nUsage: gitmap resolve <repo> [--theirs | --ours | --ide | --abort]")
	}

	if strategy != "" {
		_, err := ResolveConflictAction(repoPath, strategy)
		return err
	}

	files, err := GetConflictedFiles(repoPath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Printf("%s No merge conflicts detected in %s\n", constants.ColorGreen+"✓"+constants.ColorReset, repoPath)
		return nil
	}

	displayName := target
	if displayName == "" {
		displayName = filepath.Base(getAbsPath(repoPath))
	}
	printConflictTree(displayName, repoPath, files)
	printQuickActions(displayName)

	autoLaunchIDEIfAvailable(repoPath, files)
	return nil
}

func autoLaunchIDEIfAvailable(repoPath string, files []string) {
	_, name := DetectInstalledIDE()
	if name == "" {
		return
	}
	fmt.Printf("  Auto-launching %s...\n", name)
	_ = LaunchIDEForRepo(repoPath, files)
	fmt.Print(FormatMergeEditorGuide(repoPath, files))
}
