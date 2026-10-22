package constants

import (
	"strings"
	"testing"
)

func verifyConstant(t *testing.T, name, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s mismatch: got %q, want %q", name, got, want)
	}
}

func TestGitProviderTypeConstants(t *testing.T) {
	verifyConstant(t, "GitProviderGitHub", string(GitProviderGitHub), "github")
	verifyConstant(t, "GitProviderGitLab", string(GitProviderGitLab), "gitlab")
	verifyConstant(t, "GitProviderBitbucket", string(GitProviderBitbucket), "bitbucket")
	verifyConstant(t, "GitProviderAzureDevOps", string(GitProviderAzureDevOps), "azuredevops")
	verifyConstant(t, "GitProviderUnknown", string(GitProviderUnknown), "unknown")
}

func TestGitHostTypeConstants(t *testing.T) {
	verifyConstant(t, "GitHostGitHub", string(GitHostGitHub), "github.com")
	verifyConstant(t, "GitHostGitLab", string(GitHostGitLab), "gitlab.com")
	verifyConstant(t, "GitHostBitbucket", string(GitHostBitbucket), "bitbucket.org")
	verifyConstant(t, "GitHostAzureDevOps", string(GitHostAzureDevOps), "dev.azure.com")
}

func TestGitURLPrefixTypeGitHub(t *testing.T) {
	verifyConstant(t, "GitHubHTTPSPrefix", string(GitHubHTTPSPrefix), "https://github.com/")
	verifyConstant(t, "GitHubHTTPPrefix", string(GitHubHTTPPrefix), "http://github.com/")
	verifyConstant(t, "GitHubSSHPrefix", string(GitHubSSHPrefix), "git@github.com:")
	verifyConstant(t, "GitHubSSHSchemePrefix", string(GitHubSSHSchemePrefix), "ssh://git@github.com/")
}

func TestGitURLPrefixTypeGitLab(t *testing.T) {
	verifyConstant(t, "GitLabHTTPSPrefix", string(GitLabHTTPSPrefix), "https://gitlab.com/")
	verifyConstant(t, "GitLabHTTPPrefix", string(GitLabHTTPPrefix), "http://gitlab.com/")
	verifyConstant(t, "GitLabSSHPrefix", string(GitLabSSHPrefix), "git@gitlab.com:")
	verifyConstant(t, "GitLabSSHSchemePrefix", string(GitLabSSHSchemePrefix), "ssh://git@gitlab.com/")
}

func TestGitURLPrefixTypeBitbucket(t *testing.T) {
	verifyConstant(t, "BitbucketHTTPSPrefix", string(BitbucketHTTPSPrefix), "https://bitbucket.org/")
	verifyConstant(t, "BitbucketHTTPPrefix", string(BitbucketHTTPPrefix), "http://bitbucket.org/")
	verifyConstant(t, "BitbucketSSHPrefix", string(BitbucketSSHPrefix), "git@bitbucket.org:")
	verifyConstant(t, "BitbucketSSHSchemePrefix", string(BitbucketSSHSchemePrefix), "ssh://git@bitbucket.org/")
}

func TestGitPrefixGeneric(t *testing.T) {
	verifyConstant(t, "GitPrefixHTTPS", string(GitPrefixHTTPS), "https://")
	verifyConstant(t, "GitPrefixHTTP", string(GitPrefixHTTP), "http://")
	verifyConstant(t, "GitPrefixSSH", string(GitPrefixSSH), "git@")
	verifyConstant(t, "GitPrefixSSHScheme", string(GitPrefixSSHScheme), "ssh://")
}

func TestGitURLFormatConstants(t *testing.T) {
	verifyConstant(t, "GitHubSSHRepoOwnerFmt", GitHubSSHRepoOwnerFmt, "git@github.com:%s/")
	verifyConstant(t, "GitHubSSHRepoFmt", GitHubSSHRepoFmt, "git@github.com:%s.git")
	verifyConstant(t, "GitHubHTTPSRepoFmt", GitHubHTTPSRepoFmt, "https://github.com/%s")
	verifyConstant(t, "GitHubHTTPSRepoGitFmt", GitHubHTTPSRepoGitFmt, "https://github.com/%s.git")
	verifyConstant(t, "DefaultGitHubOwner", DefaultGitHubOwner, "alimtvnetwork")
}

func TestGitDelimiterAndTokenConstants(t *testing.T) {
	verifyConstant(t, "GitDotGitSuffix", GitDotGitSuffix, ".git")
	verifyConstant(t, "GitPathSeparator", GitPathSeparator, "/")
	verifyConstant(t, "GitColonSeparator", GitColonSeparator, ":")
	verifyConstant(t, "GitHostTokenGitHubSlash", GitHostTokenGitHubSlash, "github.com/")
	verifyConstant(t, "GitHostTokenGitHubColon", GitHostTokenGitHubColon, "github.com:")
}

func TestGitPrefixConsistency(t *testing.T) {
	hasHTTPS := strings.HasPrefix(string(GitHubHTTPSPrefix), string(GitPrefixHTTPS))
	verifyConstant(t, "GitHubHTTPSPrefixHasGitPrefixHTTPS", boolToString(hasHTTPS), "true")

	hasSSH := strings.HasPrefix(string(GitHubSSHPrefix), string(GitPrefixSSH))
	verifyConstant(t, "GitHubSSHPrefixHasGitPrefixSSH", boolToString(hasSSH), "true")
}

func boolToString(isVal bool) string {
	if isVal {
		return "true"
	}

	return "false"
}
