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

func isSameOwnerHTTPSRepo(repoPath string) bool {
	origin := gitutil.GetRemoteOriginURL(repoPath)
	if !strings.HasPrefix(strings.ToLower(origin), string(constants.GitHubHTTPSPrefix)) {
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
	if strings.HasPrefix(lower, string(constants.GitHubHTTPSPrefix)) {
		return extractFirstSlashSegment(trimmed[len(constants.GitHubHTTPSPrefix):])
	}

	if strings.HasPrefix(lower, string(constants.GitHubSSHPrefix)) {
		return extractFirstSlashSegment(trimmed[len(constants.GitHubSSHPrefix):])
	}

	return ""
}

func extractFirstSlashSegment(rest string) string {
	parts := strings.SplitN(rest, constants.GitPathSeparator, 2)
	if len(parts) == 0 {
		return ""
	}

	return strings.TrimSpace(parts[0])
}

func hasSSHRepoForSameOwner(owner string) bool {
	if strings.EqualFold(owner, constants.DefaultGitHubOwner) {
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
	prefix := fmt.Sprintf(constants.GitHubSSHRepoOwnerFmt, strings.ToLower(owner))
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
	_ = exec.Command(constants.GitBin, constants.GitConfigCmd, "--global", "credential.helper", "manager").Run()
	_ = exec.Command(constants.GitBin, constants.GitConfigCmd, "--global", "credential.credentialStore", string(storeType)).Run()
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
	if err != nil {
		return
	}

	defer db.Close()
	_ = db.UpdateRepoTransport(repoPath, sshURL, constants.TransportSSH)
}

func clearResolvedPullError(repoName, repoPath string) {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return
	}

	defer db.Close()
	_ = db.ClearPullErrorsForRepo(repoName, repoPath)
}
