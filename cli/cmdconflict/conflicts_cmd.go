package cmdconflict

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconflict/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconflict/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
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
	cmd := exec.Command(constants.GitBin, constants.GitFlagDir, repoPath, constants.GitSubCmdDiff, constants.GitArgNameOnly, constants.GitDiffFilterUnmerged)
	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		return nil
	}

	return extractNonEmptyLines(string(out))
}

func extractNonEmptyLines(output string) []string {
	var files []string
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			files = append(files, trimmed)
		}
	}

	return files
}

func parseStatusUnmergedFiles(repoPath string, existing []string) []string {
	statusCmd := exec.Command(constants.GitBin, constants.GitFlagDir, repoPath, constants.GitSubCmdStatus, constants.GitArgPorcelain)
	statusOut, sErr := statusCmd.Output()
	if sErr != nil || len(statusOut) == 0 {
		return existing
	}

	return appendUnmergedStatusFiles(string(statusOut), existing)
}

func appendUnmergedStatusFiles(output string, existing []string) []string {
	files := existing
	for _, line := range strings.Split(output, "\n") {
		files = appendStatusLineIfUnmerged(line, files)
	}

	return files
}

func appendStatusLineIfUnmerged(line string, files []string) []string {
	clean := strings.TrimRight(line, "\r")
	if len(clean) < 4 {
		return files
	}

	code := clean[:2]
	filePath := strings.TrimSpace(clean[3:])
	if isUnmergedStatusCode(code) && !hasString(files, filePath) {
		return append(files, filePath)
	}

	return files
}

func isUnmergedStatusCode(code string) bool {
	statusCode, err := enums.ParseConflictStatusCodeType(code)
	if err != nil {
		return false
	}

	return statusCode.IsUnmerged()
}

func hasString(list []string, item string) bool {
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

	cmd := exec.Command(constants.GitBin, constants.GitFlagDir, target, constants.GitSubCmdRevParse, constants.GitArgIsInsideWorkTree)
	out, err := cmd.Output()

	return err == nil && strings.TrimSpace(string(out)) == constants.GitTrue
}

// ResolveRepoPath resolves a repo query (name, slug, or path) to an existing directory.
func ResolveRepoPath(repoQuery string) string {
	clean := strings.TrimSpace(repoQuery)
	if path, hasResolved := resolveDirectRepoPath(clean); hasResolved {
		return path
	}

	return queryStoredRepoDir(clean)
}

func resolveDirectRepoPath(clean string) (string, bool) {
	if clean == "" && isInsideGitRepo(".") {
		return ".", true
	}

	if clean == "" {
		return "", true
	}

	if info, err := os.Stat(clean); err == nil && info.IsDir() {
		return clean, true
	}

	return "", false
}

func queryStoredRepoDir(clean string) string {
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

	return checkValidDir(rec.LocalPath)
}

func checkValidDir(path string) string {
	info, statErr := os.Stat(path)
	if statErr == nil && info.IsDir() {
		return path
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

	return findFirstValidRepoDir(recs)
}

func findFirstValidRepoDir(recs []model.ScanRecord) string {
	for _, r := range recs {
		if r.AbsolutePath == "" {
			continue
		}

		if valid := checkValidDir(r.AbsolutePath); valid != "" {
			return valid
		}
	}

	return ""
}

// RunConflictsCmd implements 'gitmap conflicts [repo]'.
func RunConflictsCmd(args []string) error {
	if hasHelpFlag(args) {
		helpdoc.Print(constants.CmdConflicts)

		return nil
	}

	target := extractFirstNonFlagArg(args)
	scope := enums.ResolveConflictScope(target, isInsideGitRepo("."))

	return dispatchConflictScope(scope, target)
}

func dispatchConflictScope(scope enums.ConflictScopeType, target string) error {
	switch scope {
	case enums.ConflictScopeTarget:
		return inspectTargetRepoConflicts(target)
	case enums.ConflictScopeCurrent:
		return inspectCurrentRepoConflicts()
	default:
		return scanTrackedReposForConflicts()
	}
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
		return apperror.NewNotFoundError(fmt.Sprintf(constants.ErrRepoNotFoundFmt, target))
	}

	files, err := GetConflictedFiles(repoPath)
	if err != nil {
		return err
	}

	return renderTargetRepoReport(target, repoPath, files)
}

func renderTargetRepoReport(target, repoPath string, files []string) error {
	if len(files) == 0 {
		fmt.Printf(constants.MsgNoConflictsDetectedFmt, constants.ColorGreen+constants.GlyphCheck+constants.ColorReset, repoPath)

		return nil
	}

	printConflictTree(target, repoPath, files)
	printQuickActions(target)

	return nil
}

func scanTrackedReposForConflicts() error {
	repos, err := loadTrackedRepos()
	if err != nil {
		return err
	}

	conflictCount := checkReposConflicts(repos)
	if conflictCount == 0 {
		fmt.Printf(constants.MsgNoConflictsAcrossTrackedFmt,
			constants.ColorGreen+constants.GlyphCheck+constants.ColorReset, len(repos))
	}

	return nil
}

func loadTrackedRepos() ([]model.ScanRecord, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, constants.ErrCannotReadRepoStore)
	}

	defer db.Close()
	repos, err := db.ListRepos()
	if err != nil {
		return nil, apperror.WrapSimple(err, constants.ErrCannotListRepos)
	}

	return repos, nil
}

func checkReposConflicts(repos []model.ScanRecord) int {
	conflictCount := 0
	for _, r := range repos {
		if hasRepoConflict(r) {
			conflictCount++
		}
	}

	return conflictCount
}

func hasRepoConflict(r model.ScanRecord) bool {
	if r.AbsolutePath == "" {
		return false
	}

	files, _ := GetConflictedFiles(r.AbsolutePath)
	if len(files) == 0 {
		return false
	}

	printConflictTree(r.RepoName, r.AbsolutePath, files)
	printQuickActions(r.RepoName)

	return true
}

func inspectCurrentRepoConflicts() error {
	files, err := GetConflictedFiles(".")
	if err != nil {
		return err
	}

	return renderCurrentRepoConflicts(files)
}

func renderCurrentRepoConflicts(files []string) error {
	if len(files) == 0 {
		fmt.Printf(constants.MsgNoConflictsCurrentRepoFmt, constants.ColorGreen+constants.GlyphCheck+constants.ColorReset)

		return nil
	}

	name := filepath.Base(getAbsPath("."))
	printConflictTree(name, ".", files)
	printQuickActions(name)

	return nil
}

func printConflictTree(repoName, repoPath string, files []string) {
	fmt.Printf(constants.HeaderRepoConflictsFmt, repoName, repoPath)
	fmt.Printf(constants.HeaderConflictedFilesCountFmt, len(files))
	for i, f := range files {
		prefix := constants.TreePrefixBranch
		if i == len(files)-1 {
			prefix = constants.TreePrefixLeaf
		}

		fmt.Printf(constants.LineConflictFileFmt, prefix, f)
	}
}

func printQuickActions(repoName string) {
	fmt.Print(constants.HelpSectionQuickActions)
	fmt.Printf(constants.HelpActionIDE, repoName)
	fmt.Printf(constants.HelpActionTheirs, repoName)
	fmt.Printf(constants.HelpActionOurs, repoName)
	fmt.Printf(constants.HelpActionAbort, repoName)
	fmt.Print(constants.HelpActionCommit)
}

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if a == constants.FlagHelpShort || a == constants.FlagHelpLong || a == constants.FlagHelp {
			return true
		}
	}

	return false
}

func getAbsPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	return abs
}
