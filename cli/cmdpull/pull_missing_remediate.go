package cmdpull

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// LookupRepoRemoteAndID finds the RepoId and remote URL for a given repository path.
func LookupRepoRemoteAndID(repoPath string) (int64, string, error) {
	if repoPath == "" {
		return 0, "", nil
	}

	db, err := store.OpenDefault()
	if err != nil {
		return 0, "", err
	}

	defer db.Close()

	return queryRepoRemoteAndID(db, repoPath)
}

func queryRepoRemoteAndID(db *store.DB, repoPath string) (int64, string, error) {
	return db.GetRepoPrimaryURLByPath(repoPath)
}

// LookupRepoRemoteAndIDBySlug finds the RepoId and remote URL for a given repository slug.
func LookupRepoRemoteAndIDBySlug(slug string) (int64, string, error) {
	if slug == "" {
		return 0, "", nil
	}

	db, err := store.OpenDefault()
	if err != nil {
		return 0, "", err
	}

	defer db.Close()

	return queryRepoRemoteAndIDBySlug(db, slug)
}

func queryRepoRemoteAndIDBySlug(db *store.DB, slug string) (int64, string, error) {
	recs, err := db.FindBySlug(slug)
	if err != nil || len(recs) == 0 {
		return 0, "", err
	}

	rec := recs[0]
	url := rec.SSHUrl
	if url == "" {
		url = rec.HTTPSUrl
	}

	return rec.ID, url, nil
}

// DeleteRepoRecordByPath removes the repository row from gitmap.db.
func DeleteRepoRecordByPath(repoPath string) error {
	if repoPath == "" {
		return nil
	}

	db, err := store.OpenDefault()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.DeleteByPath(repoPath)

	return err
}

func remediateMissingRepoTransport(repoPath, remoteURL string) error {
	db, err := store.OpenDefault()
	if err != nil {
		return err
	}

	defer db.Close()

	transport := constants.TransportSSH
	if strings.HasPrefix(strings.ToLower(remoteURL), string(constants.GitPrefixHTTPS)) {
		transport = constants.TransportHTTPS
	}

	return db.UpdateRepoTransport(repoPath, remoteURL, transport)
}

func printMissingRepoNotice(f PullFailureSummary) {
	relPath := formatRelativeOrCleanPath(f.RepoPath)
	if relPath == "" || relPath == "." {
		relPath = f.RepoName
	}

	fmt.Printf("  • %s%s%s\n", constants.ColorBold, f.RepoName, constants.ColorReset)
	fmt.Printf("    %s Repository directory missing on disk: %s\n", constants.TreeBranch, relPath)
}

func resolveMissingRemoteURL(f PullFailureSummary) string {
	if f.RemoteURL != "" {
		return f.RemoteURL
	}

	_, remoteURL, _ := LookupRepoRemoteAndID(f.RepoPath)
	if remoteURL != "" {
		return remoteURL
	}

	if f.RepoName != "" {
		_, remoteURL, _ = LookupRepoRemoteAndIDBySlug(f.RepoName)
	}

	return remoteURL
}

func printMissingRemoteURLError() {
	fmt.Printf("    %s %s✗%s Cannot clone: no remote URL recorded in database.\n",
		constants.TreeCorner, constants.ColorRed, constants.ColorReset)
}

func printCloneFailure(err error) {
	fmt.Printf("    %s %s✗%s Clone fallback failed: %v\n",
		constants.TreeCorner, constants.ColorRed, constants.ColorReset, err)
}

func reportPostCloneStatus(repoPath, remoteURL string) {
	if errTrans := remediateMissingRepoTransport(repoPath, remoteURL); errTrans != nil {
		fmt.Printf("    %s %s⚠%s Notice: cloned successfully, but failed to record transport in DB: %v\n",
			constants.TreeBranch, constants.ColorYellow, constants.ColorReset, errTrans)
	}

	fmt.Printf("    %s %s✓%s Successfully restored repository via clone fallback.\n",
		constants.TreeCorner, constants.ColorGreen, constants.ColorReset)
}

func purgeInvalidRepoEntry(repoName, repoPath string) error {
	var errReturn error
	if repoPath != "" {
		if errDel := DeleteRepoRecordByPath(repoPath); errDel != nil {
			errReturn = errDel
		}
	}

	if errClear := clearPullSplitDBError(repoName, repoPath); errClear != nil && errReturn == nil {
		errReturn = errClear
	}

	return errReturn
}

func clearPullSplitDBError(repoName, repoPath string) error {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return err
	}

	defer db.Close()

	if appErr := db.ClearPullErrorsForRepo(repoName, repoPath); appErr != nil {
		return appErr
	}

	return nil
}

func execCloneAndRecord(repoName, repoPath, remoteURL string) bool {
	if err := ExecCloneRepo(remoteURL, repoPath); err != nil {
		printCloneFailure(err)
		if strings.Contains(remoteURL, " ") || !strings.Contains(remoteURL, "/") {
			_ = purgeInvalidRepoEntry(repoName, repoPath)
		}

		return false
	}

	reportPostCloneStatus(repoPath, remoteURL)

	return true
}

// RemediateMissingRepo attempts to restore a missing repository using clone fallback.
func RemediateMissingRepo(f PullFailureSummary) bool {
	printMissingRepoNotice(f)

	remoteURL := resolveMissingRemoteURL(f)
	if remoteURL == "" || strings.Contains(remoteURL, " ") {
		printMissingRemoteURLError()
		_ = purgeInvalidRepoEntry(f.RepoName, f.RepoPath)

		return false
	}

	fmt.Printf("    %s Attempting clone fallback from %s...\n", constants.TreeBranch, remoteURL)

	return execCloneAndRecord(f.RepoName, f.RepoPath, remoteURL)
}

func cleanDestinationIfNonRepo(destination string) {
	if destination == "" || gitutil.IsValidGitRepository(destination) {
		return
	}

	info, err := os.Stat(destination)
	if err != nil || !info.IsDir() {
		return
	}

	entries, readErr := os.ReadDir(destination)
	if readErr == nil && len(entries) <= 1 {
		_ = os.RemoveAll(destination)
	}
}

// ExecCloneRepo clones a remote repository into destination path.
func ExecCloneRepo(remoteURL, destination string) error {
	parentDir := filepath.Dir(destination)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return err
	}

	cleanDestinationIfNonRepo(destination)
	cmd := exec.Command(constants.GitBin, constants.GitClone, constants.GitProgressFlag, remoteURL, destination)
	cmd.Env = gitutil.BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)

	return cmd.Run()
}

func isCloneChoice(choice string) bool {
	return choice == "c" || choice == "clone" || choice == "1" || choice == "y" || choice == "yes" || choice == ""
}

func isSkipChoice(choice string) bool {
	return choice == "s" || choice == "skip" || choice == "3"
}

func isQuitChoice(choice string) bool {
	return choice == "q" || choice == "quit"
}

func isDeleteChoice(choice string) bool {
	return choice == "d" || choice == "delete" || choice == "2"
}

func handleCloneOrFallback(choice, remoteURL string, f PullFailureSummary) (bool, bool) {
	if isCloneChoice(choice) {
		f.RemoteURL = remoteURL

		return true, RemediateMissingRepo(f)
	}

	return true, false
}

func dispatchPromptMissingChoice(choice, remoteURL string, f PullFailureSummary) (bool, bool) {
	if isQuitChoice(choice) {
		return false, false
	}

	if isSkipChoice(choice) {
		return true, false
	}

	if isDeleteChoice(choice) {
		return handlePromptDeleteAction(f)
	}

	return handleCloneOrFallback(choice, remoteURL, f)
}

func printPromptRemote(remoteURL string) {
	hasRemote := len(remoteURL) > 0
	if hasRemote {
		fmt.Printf("    Remote: %s\n", remoteURL)

		return
	}

	fmt.Println("    Remote: (none)")
}

func printPromptMissingOptions(remoteURL string) {
	printPromptRemote(remoteURL)

	fmt.Println("    Options:")
	fmt.Println("      [1/c] Clone from remote URL")
	fmt.Println("      [2/d] Delete record from database (clean up stale tracking)")
	fmt.Println("      [s/3] Skip")
	fmt.Println("      [q]   Quit")
	fmt.Printf("  %sChoice [c/d/s/q]:%s ", constants.ColorCyan, constants.ColorReset)
}

// PromptMissingRepoAction provides an interactive prompt for missing repository remediation.
// Returns (shouldContinue, isResolved).
func PromptMissingRepoAction(reader *bufio.Reader, f PullFailureSummary) (bool, bool) {
	remoteURL := resolveMissingRemoteURL(f)
	printPromptMissingOptions(remoteURL)

	line, readErr := reader.ReadString('\n')
	if readErr != nil {
		return false, false
	}

	choice := strings.ToLower(strings.TrimSpace(line))

	return dispatchPromptMissingChoice(choice, remoteURL, f)
}

func handlePromptDeleteAction(f PullFailureSummary) (bool, bool) {
	if err := DeleteRepoRecordByPath(f.RepoPath); err != nil {
		fmt.Printf("  %s✗%s [%s] Failed to remove record from database: %v\n",
			constants.ColorRed, constants.ColorReset, f.RepoName, err)

		return true, false
	}

	fmt.Printf("  %s✓%s [%s] Stale record removed from database.\n",
		constants.ColorGreen, constants.ColorReset, f.RepoName)

	return true, true
}
