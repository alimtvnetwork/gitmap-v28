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
	target := repoPath
	if target == "" {
		target = "."
	}
	diffFiles := parseDiffUnmergedFiles(target)
	statusFiles := parseStatusUnmergedFiles(target, diffFiles)
	return statusFiles, nil
}

func parseDiffUnmergedFiles(repoPath string) []string {
	cmd := exec.Command("git", "-C", repoPath, "diff", "--name-only", "--diff-filter=U")
	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		return nil
	}
	var files []string
	for _, line := range strings.Split(string(out), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			files = append(files, trimmed)
		}
	}
	return files
}

func parseStatusUnmergedFiles(repoPath string, existing []string) []string {
	statusCmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	statusOut, sErr := statusCmd.Output()
	if sErr != nil || len(statusOut) == 0 {
		return existing
	}
	files := existing
	for _, line := range strings.Split(string(statusOut), "\n") {
		line = strings.TrimRight(line, "\r")
		if len(line) < 4 {
			continue
		}
		code := line[:2]
		filePath := strings.TrimSpace(line[3:])
		isUnmerged := isUnmergedStatusCode(code) && !containsString(files, filePath)
		if isUnmerged {
			files = append(files, filePath)
		}
	}
	return files
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
	target := path
	if target == "" {
		target = "."
	}
	cmd := exec.Command("git", "-C", target, "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(out)) == "true"
}

// ResolveRepoPath resolves a repo query (name, slug, or path) to an existing directory.
func ResolveRepoPath(repoQuery string) string {
	clean := strings.TrimSpace(repoQuery)
	if clean == "" && isInsideGitRepo(".") {
		return "."
	}
	if clean == "" {
		return ""
	}
	if info, err := os.Stat(clean); err == nil && info.IsDir() {
		return clean
	}
	if found := querySpecialRepoDir(clean); found != "" {
		return found
	}
	if found := queryDefaultRepoDir(clean); found != "" {
		return found
	}
	return clean
}

func querySpecialRepoDir(clean string) string {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return ""
	}
	defer db.Close()
	rec, err := db.GetSpecialRepo(clean)
	if err != nil || rec == nil || rec.LocalPath == "" {
		return ""
	}
	info, statErr := os.Stat(rec.LocalPath)
	if statErr == nil && info.IsDir() {
		return rec.LocalPath
	}
	return ""
}

func queryDefaultRepoDir(clean string) string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer db.Close()
	recs, err := db.FindBySlug(clean)
	if err != nil || len(recs) == 0 {
		return ""
	}
	for _, r := range recs {
		if r.AbsolutePath == "" {
			continue
		}
		info, statErr := os.Stat(r.AbsolutePath)
		if statErr == nil && info.IsDir() {
			return r.AbsolutePath
		}
	}
	return ""
}

// RunConflictsCmd implements 'gitmap conflicts [repo]'.
func RunConflictsCmd(args []string) error {
	if hasHelpFlag(args) {
		helpdoc.Print("conflicts")
		return nil
	}

	target := extractFirstNonFlagArg(args)
	if target != "" {
		return inspectTargetRepoConflicts(target)
	}

	if isInsideGitRepo(".") {
		return inspectCurrentRepoConflicts()
	}

	return scanTrackedReposForConflicts()
}

func extractFirstNonFlagArg(args []string) string {
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			return a
		}
	}
	return ""
}

func inspectTargetRepoConflicts(target string) error {
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

func scanTrackedReposForConflicts() error {
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

func inspectCurrentRepoConflicts() error {
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
