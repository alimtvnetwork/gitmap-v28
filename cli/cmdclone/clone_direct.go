// Package cmdclone — clone_direct.go handles direct URL cloning.
package cmdclone

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/usercontext"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

// isDirectURL returns true when source is a git URL (not a file path).
func isDirectURL(source string) bool {
	lower := strings.ToLower(source)
	if strings.HasPrefix(lower, constants.PrefixHTTPS) ||
		strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, constants.PrefixSSH) {
		return true
	}

	if strings.HasPrefix(lower, "git@") {
		at := strings.Index(lower, "@")
		colon := strings.Index(lower[at:], ":")

		return colon > 0
	}

	return false
}

// repoNameFromURL derives the repository name from a clone URL.
func repoNameFromURL(url string) string {
	name := strings.TrimRight(url, "/\\")
	name = strings.TrimSuffix(name, ".git")
	name = strings.TrimRight(name, "/\\")
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}

	if idx := strings.LastIndex(name, ":"); idx >= 0 {
		name = name[idx+1:]
	}

	return name
}

func cleanDirectCloneTarget(absPath string, isClean bool) {
	if !isClean || absPath == "" {
		return
	}

	if _, err := os.Stat(absPath); err == nil {
		fmt.Printf("  • Force/clean removing existing directory: %s\n", absPath)
		_ = os.RemoveAll(absPath)
	}
}

func maybeCleanMultiFolder(url string, isClean bool) {
	if !isClean {
		return
	}

	multiFolder := resolveCloneFolder(repoNameFromURL(url), "")
	absFolder, errFolder := filepath.Abs(multiFolder)
	if errFolder == nil {
		cleanDirectCloneTarget(absFolder, true)
	}
}

func handleExistingGitRepo(params DirectCloneParams, repoName, absPath string) {
	if params.IsJSON {
		emitDirectCloneJSON(DirectCloneJSONResponse{
			Success:  true,
			RepoName: repoName,
			Status:   "already_exists",
			Message:  fmt.Sprintf("already exists on disk (%s)", absPath),
			Path:     absPath,
		})

		return
	}

	fmt.Printf("~ %s already exists on disk (%s), skipping clone and workspace re-registration.\n", repoName, absPath)
}

func printCloneStartup(params DirectCloneParams, url, absPath, repoName, folderName string) {
	if params.IsJSON {
		return
	}

	printCloneTermBlockForURL(params.Output, 1, url, absPath)
	fmt.Printf(constants.MsgCloneURLCloning, repoName, folderName)
}

func emitDirectCloneFailureJSON(repoName, absPath string, cloneErr error) {
	emitDirectCloneJSON(DirectCloneJSONResponse{
		Success:  false,
		RepoName: repoName,
		Status:   "failed",
		Message:  cloneErr.Error(),
		Path:     absPath,
	})
}

func handleDirectCloneFailure(params DirectCloneParams, taskDB *store.DB, taskID int64, url, repoName, absPath string, cloneErr error) error {
	failPendingTask(taskDB, taskID, fmt.Sprintf(constants.ErrCloneURLFailed, url, cloneErr))
	closeTaskDB(taskDB)
	if params.IsJSON {
		emitDirectCloneFailureJSON(repoName, absPath, cloneErr)

		return cloneErr
	}

	fmt.Fprintf(os.Stderr, constants.ErrCloneURLFailed, url, cloneErr)
	printRepoSlugSuggestions(url)

	return cloneErr
}

func handleDirectCloneSuccess(params DirectCloneParams, repoName, absPath string) {
	if params.IsJSON {
		emitDirectCloneJSON(DirectCloneJSONResponse{
			Success:  true,
			RepoName: repoName,
			Status:   "cloned",
			Message:  fmt.Sprintf("cloned %s successfully", repoName),
			Path:     absPath,
		})

		return
	}

	fmt.Printf(constants.MsgCloneURLDone, repoName)
}

func resolveDirectFolderName(folderName, repoName string) string {
	if len(folderName) == 0 {
		return repoName
	}

	return folderName
}

func resolveDirectAbsPath(folderName string) string {
	absPath, err := filepath.Abs(folderName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not resolve absolute path for %s: %v\n", folderName, err)

		return folderName
	}

	return absPath
}

func checkDirectTargetCollision(absPath string, noReplace bool) error {
	_, statErr := os.Stat(absPath)
	if noReplace && statErr == nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneURLExists, absPath)

		return apperror.NewExecutionError(fmt.Sprintf(constants.ErrCloneURLExists, absPath))
	}

	return nil
}

func validateDirectCloneFolder(folderName string, noReplace bool) (string, error) {
	if isDirectURL(folderName) {
		fmt.Fprintf(os.Stderr, constants.ErrCloneStaleBinaryFolderURL, folderName, constants.Version)

		return "", apperror.NewExecutionError(fmt.Sprintf(constants.ErrCloneStaleBinaryFolderURL, folderName, constants.Version))
	}

	absPath := resolveDirectAbsPath(folderName)
	if err := checkDirectTargetCollision(absPath, noReplace); err != nil {
		return "", err
	}

	return absPath, nil
}

func resolveDirectCloneURL(rawURL string, useSSH bool, repoName string) string {
	url := coerceURLToStoredTransport(rawURL)
	url = resolveDirectSSHURL(url, useSSH)
	if isSSHCloneURL(url) {
		return url
	}

	return resolveDirectHTTPSURL(url, repoName)
}

func initDirectCloneTask(absPath string) (int64, *store.DB) {
	workDir, _ := os.Getwd()
	cmdArgs := buildCommandArgs(append([]string{"clone"}, os.Args[2:]...))

	return createPendingTask(constants.TaskTypeClone, absPath, workDir, "clone", cmdArgs)
}

func applyWinningProfile(absPath string, winningProfile *model.GitProfile) {
	if winningProfile == nil {
		return
	}

	_ = usercontext.BindProject(absPath, winningProfile.Name)
	_ = usercontext.SyncProjectUser(absPath, *winningProfile)
}

func postCloneRegistration(params DirectCloneParams, url, repoName, folderName, absPath string) {
	persistRecloneTransport(url)
	handleDirectCloneSuccess(params, repoName, absPath)
	upsertDirectClone(url, repoName, folderName, absPath)
	registerSingleDesktop(repoName, absPath)
	_ = workspacesync.SyncAntigravity(absPath, repoName)
	WriteShellHandoff(absPath)
}

func postCloneEditors(params DirectCloneParams, repoName, absPath string) {
	if !params.IsJSON {
		openInVSCode(absPath)
	}

	syncSingleClonedRepoToVSCodePM(absPath, repoName, params.NoVSCodeSync)
}

func performDirectCloneRun(params DirectCloneParams, repoName, folderName, absPath string) error {
	taskID, taskDB := initDirectCloneTask(absPath)
	url := resolveDirectCloneURL(params.URL, params.UseSSH, repoName)
	printCloneStartup(params, url, absPath, repoName, folderName)

	winningProfile, cloneErr := runDirectCloneWithProfileFallback(params, url, absPath)
	if cloneErr != nil {
		return handleDirectCloneFailure(params, taskDB, taskID, url, repoName, absPath, cloneErr)
	}

	applyWinningProfile(absPath, winningProfile)
	postCloneRegistration(params, url, repoName, folderName, absPath)
	postCloneEditors(params, repoName, absPath)
	completePendingTask(taskDB, taskID)
	closeTaskDB(taskDB)

	return nil
}

// executeDirectClone clones a single repo from a direct URL.
func executeDirectClone(params DirectCloneParams) error {
	escapeNestedGitRepo()
	params.URL = ResolveRepoSlug(params.URL)
	repoName := repoNameFromURL(params.URL)
	folderName := resolveDirectFolderName(params.FolderName, repoName)

	absPath, err := validateDirectCloneFolder(folderName, params.NoReplace)
	if err != nil {
		return err
	}

	cleanDirectCloneTarget(absPath, params.IsClean)
	if isGitRepo(absPath) {
		handleExistingGitRepo(params, repoName, absPath)

		return nil
	}

	return performDirectCloneRun(params, repoName, folderName, absPath)
}

func upsertDirectClone(url, repoName, folderName, absPath string) {
	rec := model.ScanRecord{
		Slug:         strings.ToLower(repoName),
		RepoName:     repoName,
		RelativePath: folderName,
		AbsolutePath: absPath,
	}

	populateDirectCloneURLs(&rec, url)

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not open database: %v\n", err)

		return
	}

	defer db.Close()

	if upsertErr := db.UpsertRepos([]model.ScanRecord{rec}); upsertErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not save repo to database: %v\n", upsertErr)

		return
	}

	if markErr := db.MarkCloned(absPath); markErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not stamp clone time: %v\n", markErr)
	}
}

func registerSingleDesktop(name, absPath string) {
	records := []model.ScanRecord{{
		RepoName:     name,
		AbsolutePath: absPath,
	}}
	result := desktop.AddRepos(records)
	if result.Added > 0 {
		fmt.Printf(constants.MsgDesktopSummary, result.Added, result.Failed)
	}
}
