package cmdconflict

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// GetConflictedFiles returns a list of unmerged / conflicted relative file paths in repoPath.
func GetConflictedFiles(repoPath string) ([]string, error) {
	if repoPath == "" {
		repoPath = "."
	}
	cmd := exec.Command("git", "-C", repoPath, "diff", "--name-only", "--diff-filter=U")
	out, err := cmd.Output()
	var files []string
	if err == nil && len(out) > 0 {
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				files = append(files, line)
			}
		}
	}

	statusCmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	statusOut, sErr := statusCmd.Output()
	if sErr == nil && len(statusOut) > 0 {
		for _, line := range strings.Split(string(statusOut), "\n") {
			line = strings.TrimRight(line, "\r")
			if len(line) < 4 {
				continue
			}
			code := line[:2]
			if isUnmergedStatusCode(code) {
				filePath := strings.TrimSpace(line[3:])
				if !containsString(files, filePath) {
					files = append(files, filePath)
				}
			}
		}
	}

	return files, nil
}

func isUnmergedStatusCode(code string) bool {
	switch code {
	case "UU", "AA", "UD", "DU", "DD", "AU", "UA":
		return true
	}
	return false
}

func containsString(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func isInsideGitRepo(path string) bool {
	if path == "" {
		path = "."
	}
	cmd := exec.Command("git", "-C", path, "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(out)) == "true"
}

// ResolveRepoPath resolves a repo query (name, slug, or path) to an existing directory.
func ResolveRepoPath(repoQuery string) string {
	clean := strings.TrimSpace(repoQuery)
	if clean == "" {
		if isInsideGitRepo(".") {
			return "."
		}
		return ""
	}
	if info, err := os.Stat(clean); err == nil && info.IsDir() {
		return clean
	}
	if db, err := store.OpenSpecialReposSplitDB(); err == nil {
		defer db.Close()
		if rec, err := db.GetSpecialRepo(clean); err == nil && rec.LocalPath != "" {
			if info, err := os.Stat(rec.LocalPath); err == nil && info.IsDir() {
				return rec.LocalPath
			}
		}
	}
	if db, err := store.OpenDefault(); err == nil {
		defer db.Close()
		if recs, err := db.FindBySlug(clean); err == nil && len(recs) > 0 {
			for _, r := range recs {
				if r.AbsolutePath != "" {
					if info, err := os.Stat(r.AbsolutePath); err == nil && info.IsDir() {
						return r.AbsolutePath
					}
				}
			}
		}
	}
	return clean
}

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
	if err := LaunchIDEForRepo(repoPath, files); err != nil {
		fmt.Printf("%s Could not launch IDE: %v\n", constants.ColorYellow+"⚠"+constants.ColorReset, err)
	} else {
		bin, name := DetectInstalledIDE()
		_ = bin
		fmt.Printf("%s Launched %s for %s\n", constants.ColorGreen+"✓"+constants.ColorReset, name, repoPath)
	}
	fmt.Print(FormatMergeEditorGuide(repoPath, files))
	return nil
}

func applyAbort(repoPath string) error {
	mergeCmd := exec.Command("git", "-C", repoPath, "merge", "--abort")
	if err := mergeCmd.Run(); err == nil {
		fmt.Printf("%s Merge aborted successfully in %s\n", constants.ColorGreen+"✓"+constants.ColorReset, repoPath)
		return nil
	}
	rebaseCmd := exec.Command("git", "-C", repoPath, "rebase", "--abort")
	if err := rebaseCmd.Run(); err == nil {
		fmt.Printf("%s Rebase aborted successfully in %s\n", constants.ColorGreen+"✓"+constants.ColorReset, repoPath)
		return nil
	}
	cherryCmd := exec.Command("git", "-C", repoPath, "cherry-pick", "--abort")
	if err := cherryCmd.Run(); err == nil {
		fmt.Printf("%s Cherry-pick aborted successfully in %s\n", constants.ColorGreen+"✓"+constants.ColorReset, repoPath)
		return nil
	}
	return fmt.Errorf("no merge, rebase, or cherry-pick in progress to abort in %s", repoPath)
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
		return fmt.Errorf("not inside a git repository and no repository specified.\nUsage: gitmap resolve <repo> [--theirs | --ours | --ide | --abort]")
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

	if _, name := DetectInstalledIDE(); name != "" {
		fmt.Printf("  Auto-launching %s...\n", name)
		_ = LaunchIDEForRepo(repoPath, files)
		fmt.Print(FormatMergeEditorGuide(repoPath, files))
	}

	return nil
}

// RunConflictsCmd implements 'gitmap conflicts [repo]'.
func RunConflictsCmd(args []string) error {
	if hasHelpFlag(args) {
		helpdoc.Print("conflicts")
		return nil
	}

	var target string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			target = a
			break
		}
	}

	if target != "" {
		repoPath := ResolveRepoPath(target)
		if repoPath == "" {
			return fmt.Errorf("repository %q not found or not a valid directory", target)
		}
		files, err := GetConflictedFiles(repoPath)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Printf("%s No merge conflicts detected in %s\n", constants.ColorGreen+"✓"+constants.ColorReset, repoPath)
			return nil
		}
		printConflictTree(target, repoPath, files)
		printQuickActions(target)
		return nil
	}

	if isInsideGitRepo(".") {
		files, err := GetConflictedFiles(".")
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Printf("%s No merge conflicts detected in current repository.\n", constants.ColorGreen+"✓"+constants.ColorReset)
			return nil
		}
		name := filepath.Base(getAbsPath("."))
		printConflictTree(name, ".", files)
		printQuickActions(name)
		return nil
	}

	db, err := store.OpenDefault()
	if err != nil {
		return fmt.Errorf("cannot read repository store: %w", err)
	}
	defer db.Close()

	repos, err := db.ListRepos()
	if err != nil {
		return fmt.Errorf("cannot list repositories: %w", err)
	}

	conflictCount := 0
	for _, r := range repos {
		if r.AbsolutePath == "" {
			continue
		}
		files, _ := GetConflictedFiles(r.AbsolutePath)
		if len(files) > 0 {
			conflictCount++
			printConflictTree(r.RepoName, r.AbsolutePath, files)
			printQuickActions(r.RepoName)
		}
	}

	if conflictCount == 0 {
		fmt.Printf("%s No merge conflicts detected across %d tracked repositories.\n",
			constants.ColorGreen+"✓"+constants.ColorReset, len(repos))
	}

	return nil
}

func printConflictTree(repoName, repoPath string, files []string) {
	fmt.Printf("\n📁 Repository: %s (%s)\n", repoName, repoPath)
	fmt.Printf("   Conflicted Files (%d):\n", len(files))
	for i, f := range files {
		prefix := "├──"
		if i == len(files)-1 {
			prefix = "└──"
		}
		fmt.Printf("   %s %s\n", prefix, f)
	}
}

func printQuickActions(repoName string) {
	fmt.Printf("\n⚡ Quick Actions:\n")
	fmt.Printf("   gitmap resolve %s --ide      Launch VS Code / IDE 3-Way Merge Editor\n", repoName)
	fmt.Printf("   gitmap resolve %s --theirs   Accept incoming remote changes\n", repoName)
	fmt.Printf("   gitmap resolve %s --ours     Accept local current changes\n", repoName)
	fmt.Printf("   gitmap resolve %s --abort    Abort merge or rebase\n", repoName)
	fmt.Printf("   gitmap cpf \"resolve conflicts\" Stage, commit, and push once resolved\n\n")
}

func getAbsPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}
