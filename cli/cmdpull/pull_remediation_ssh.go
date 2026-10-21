package cmdpull

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func readRemoteOriginURL(repoPath string) string {
	originOut, err := exec.Command("git", "-C", repoPath, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(originOut))
}

func isSameOwnerHTTPSRepo(repoPath string) bool {
	origin := readRemoteOriginURL(repoPath)
	if !strings.HasPrefix(strings.ToLower(origin), "https://github.com/") {
		return false
	}

	owner := extractGitHubOwnerFromURL(origin)
	if owner == "" {
		return false
	}

	return hasSSHRepoForSameOwner(owner)
}

func extractGitHubOwnerFromURL(url string) string {
	trimmed := strings.TrimSpace(url)
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "https://github.com/") {
		return extractFirstSlashSegment(trimmed[len("https://github.com/"):])
	}

	if strings.HasPrefix(lower, "git@github.com:") {
		return extractFirstSlashSegment(trimmed[len("git@github.com:"):])
	}

	return ""
}

func extractFirstSlashSegment(rest string) string {
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) == 0 {
		return ""
	}

	return strings.TrimSpace(parts[0])
}

func hasSSHRepoForSameOwner(owner string) bool {
	if strings.EqualFold(owner, "alimtvnetwork") {
		return true
	}

	repos := loadStoreRepos()
	if len(repos) == 0 {
		return false
	}

	return hasSSHRepoOwner(repos, owner)
}

func loadStoreRepos() []model.ScanRecord {
	db, err := store.OpenDefault()
	if err != nil {
		return nil
	}

	defer db.Close()
	repos, _ := db.ListRepos()

	return repos
}

func hasSSHRepoOwner(repos []model.ScanRecord, owner string) bool {
	prefix := "git@github.com:" + strings.ToLower(owner) + "/"
	for _, r := range repos {
		if strings.HasPrefix(strings.ToLower(r.SSHUrl), prefix) ||
			strings.HasPrefix(strings.ToLower(r.DiscoveredURL), prefix) {
			return true
		}
	}

	return false
}

func remediateCredentialAndConvertToSSH(repoPath, repoName string) {
	storeType := gitutil.ResolveWindowsCredentialStore()
	_ = exec.Command("git", "config", "--global", "credential.helper", "manager").Run()
	_ = exec.Command("git", "config", "--global", "credential.credentialStore", storeType).Run()
	applySSHTransportRemediation(repoPath)
}

func applySSHTransportRemediation(repoPath string) {
	hasChanged, oldURL, newURL, err := ApplyTransportFlag(repoPath, true, false)
	if err != nil || !hasChanged {
		return
	}

	fmt.Printf("    %s Converted remote to SSH: %s -> %s\n", constants.TreeBranch, oldURL, newURL)
	updateRepoTransportInDB(repoPath, newURL)
}

func updateRepoTransportInDB(repoPath, sshURL string) {
	db, err := store.OpenDefault()
	if err != nil || db.Conn() == nil {
		return
	}

	defer db.Close()
	_, _ = db.Conn().Exec(
		`UPDATE Repo SET SshUrl = ?, DiscoveredUrl = ?, IdentifiedTransport = 'ssh' WHERE AbsolutePath = ?`,
		sshURL, sshURL, repoPath,
	)
}

func clearResolvedPullError(repoName, repoPath string) {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return
	}

	defer db.Close()
	_ = db.ClearPullErrorsForRepo(repoName, repoPath)
}
