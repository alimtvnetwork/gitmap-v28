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
func LookupRepoRemoteAndID(repoPath string) (int64, string) {
	if repoPath == "" {
		return 0, ""
	}

	db, err := store.OpenDefault()
	if err != nil {
		return 0, ""
	}

	defer db.Close()

	if db.Conn() == nil {
		return 0, ""
	}

	var repoID int64
	var httpsUrl, sshUrl, discUrl string
	row := db.Conn().QueryRow(
		"SELECT RepoId, COALESCE(HttpsUrl, ''), COALESCE(SshUrl, ''), COALESCE(DiscoveredUrl, '') FROM Repo WHERE AbsolutePath = ?",
		repoPath,
	)
	if err := row.Scan(&repoID, &httpsUrl, &sshUrl, &discUrl); err != nil {
		return 0, ""
	}

	if httpsUrl != "" {
		return repoID, httpsUrl
	}

	if sshUrl != "" {
		return repoID, sshUrl
	}

	return repoID, discUrl
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

// RemediateMissingRepo attempts to restore a missing repository using clone fallback.
func RemediateMissingRepo(f PullFailureSummary) bool {
	fmt.Printf("  %s⚠%s [%s] Directory missing on disk: %s\n",
		constants.ColorYellow, constants.ColorReset, f.RepoName, f.RepoPath)

	remoteURL := f.RemoteURL
	if remoteURL == "" {
		_, remoteURL = LookupRepoRemoteAndID(f.RepoPath)
	}

	if remoteURL == "" {
		fmt.Printf("  %s✗%s [%s] Cannot clone: no remote URL recorded in database.\n",
			constants.ColorRed, constants.ColorReset, f.RepoName)

		return false
	}

	fmt.Printf("  %s→%s Attempting clone fallback from %s...\n",
		constants.ColorCyan, constants.ColorReset, remoteURL)

	if err := ExecCloneRepo(remoteURL, f.RepoPath); err != nil {
		fmt.Printf("  %s✗%s [%s] Clone fallback failed: %v\n",
			constants.ColorRed, constants.ColorReset, f.RepoName, err)

		return false
	}

	fmt.Printf("  %s✓%s [%s] Successfully restored repository via clone fallback.\n",
		constants.ColorGreen, constants.ColorReset, f.RepoName)

	return true
}

// ExecCloneRepo clones a remote repository into destination path.
func ExecCloneRepo(remoteURL, destination string) error {
	parentDir := filepath.Dir(destination)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", "--progress", remoteURL, destination)
	cmd.Env = gitutil.BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)

	return cmd.Run()
}

// PromptMissingRepoAction provides an interactive prompt for missing repository remediation.
// Returns (shouldContinue, isResolved).
func PromptMissingRepoAction(reader *bufio.Reader, f PullFailureSummary) (bool, bool) {
	remoteURL := f.RemoteURL
	if remoteURL == "" {
		_, remoteURL = LookupRepoRemoteAndID(f.RepoPath)
	}

	if remoteURL != "" {
		fmt.Printf("    Remote: %s\n", remoteURL)
	} else {
		fmt.Println("    Remote: (none)")
	}

	fmt.Println("    Options:")
	fmt.Println("      [1/c] Clone from remote URL")
	fmt.Println("      [2/d] Delete record from database (clean up stale tracking)")
	fmt.Println("      [s/3] Skip")
	fmt.Println("      [q]   Quit")
	fmt.Printf("  %sChoice [c/d/s/q]:%s ", constants.ColorCyan, constants.ColorReset)

	line, _ := reader.ReadString('\n')
	choice := strings.ToLower(strings.TrimSpace(line))

	if choice == "q" || choice == "quit" {
		return false, false
	}

	if choice == "s" || choice == "skip" || choice == "3" {
		return true, false
	}

	if choice == "d" || choice == "delete" || choice == "2" {
		return handlePromptDeleteAction(f)
	}

	// Default to clone if choice is c / clone / 1 / y / enter
	if choice == "c" || choice == "clone" || choice == "1" || choice == "y" || choice == "yes" || choice == "" {
		f.RemoteURL = remoteURL

		return true, RemediateMissingRepo(f)
	}

	return true, false
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
