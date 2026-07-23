package cmd

// Git identity resolution for `gitmap fix-repo`. Mirrors
// scripts/fix-repo/Repo-Identity.ps1: get root, get remote URL,
// parse owner/host/repo, split repo into base + numeric version.

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// fixRepoIdentity is the resolved repo identity used by the command.
type fixRepoIdentity struct {
	root    string
	host    string
	owner   string
	base    string
	current int
}

// resolveFixRepoIdentity runs git, parses the remote URL, and exits
// with the matching FixRepoExit* code on the first failure. Errors
// are written verbatim to os.Stderr per the zero-swallow rule.
func resolveFixRepoIdentity() fixRepoIdentity {
	root := mustGitRoot()
	url := mustGitRemoteURL()
	host, owner, repo := mustParseRemoteURL(url)
	base, current := mustSplitRepoVersion(repo)

	return fixRepoIdentity{root: root, host: host, owner: owner, base: base, current: current}
}

// mustGitRoot returns the repo root or exits E_NOT_A_REPO.
func mustGitRoot() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		fmt.Fprint(os.Stderr, constants.FixRepoErrNotARepo)
		os.Exit(constants.FixRepoExitNotARepo)
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		fmt.Fprint(os.Stderr, constants.FixRepoErrNotARepo)
		os.Exit(constants.FixRepoExitNotARepo)
	}

	return root
}

// mustGitRemoteURL returns origin's URL or exits E_NO_REMOTE.
func mustGitRemoteURL() string {
	out, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		fmt.Fprint(os.Stderr, constants.FixRepoErrNoRemote)
		os.Exit(constants.FixRepoExitNoRemote)
	}
	url := strings.TrimSpace(string(out))
	if url == "" {
		fmt.Fprint(os.Stderr, constants.FixRepoErrNoRemote)
		os.Exit(constants.FixRepoExitNoRemote)
	}

	return url
}

// mustParseRemoteURL extracts host/owner/repo from HTTPS, SSH, and
// `ssh://` forms. Trailing `.git` is stripped first. Exits on
// failure with E_NO_REMOTE since an unparseable URL is operationally
// equivalent to a missing one for this command.
func mustParseRemoteURL(url string) (string, string, string) {
	host, owner, repo, ok := parseRemoteURL(url)
	if !ok {
		fmt.Fprintf(os.Stderr, constants.FixRepoErrParseURLFmt, url)
		os.Exit(constants.FixRepoExitNoRemote)
	}

	return host, owner, repo
}

// parseRemoteURL is the pure-function half of mustParseRemoteURL so
// it can be unit-tested without forking git. Returns (h, o, r, ok).
func parseRemoteURL(url string) (string, string, string, bool) {
	trimmed := strings.TrimSuffix(strings.TrimRight(url, "/"), ".git")
	if h, o, r, ok := matchSSHRemote(trimmed); ok {
		return h, o, r, true
	}
	if h, o, r, ok := matchHTTPSRemote(trimmed); ok {
		return h, o, r, true
	}
	if h, o, r, ok := matchSSHProtoRemote(trimmed); ok {
		return h, o, r, true
	}

	return "", "", "", false
}

var (
	reSSHRemote      = regexp.MustCompile(`^[^@]+@([^:]+):([^/]+)/(.+)$`)
	reHTTPSRemote    = regexp.MustCompile(`^https?://([^/]+)/([^/]+)/(.+)$`)
	reSSHProtoRemote = regexp.MustCompile(`^ssh://[^@]+@([^/:]+)(?::\d+)?/([^/]+)/(.+)$`)
	reRepoVersion    = regexp.MustCompile(`^(.+)-v(\d+)$`)
)

// matchSSHRemote handles `git@host:owner/repo`.
func matchSSHRemote(url string) (string, string, string, bool) {
	m := reSSHRemote.FindStringSubmatch(url)
	if len(m) != 4 {
		return "", "", "", false
	}

	return m[1], m[2], m[3], true
}

// matchHTTPSRemote handles `https?://host/owner/repo`.
func matchHTTPSRemote(url string) (string, string, string, bool) {
	m := reHTTPSRemote.FindStringSubmatch(url)
	if len(m) != 4 {
		return "", "", "", false
	}

	return m[1], m[2], m[3], true
}

// matchSSHProtoRemote handles `ssh://git@host[:port]/owner/repo`.
func matchSSHProtoRemote(url string) (string, string, string, bool) {
	m := reSSHProtoRemote.FindStringSubmatch(url)
	if len(m) != 4 {
		return "", "", "", false
	}

	return m[1], m[2], m[3], true
}

// mustSplitRepoVersion splits `<base>-v<N>` into (base, N) or exits.
// Two distinct exit codes preserve script parity:
// E_NO_VERSION_SUFFIX (4) when no `-vN` is present at all,
// E_BAD_VERSION (5) when N parses but is ≤ 0.
func mustSplitRepoVersion(repo string) (string, int) {
	m := reRepoVersion.FindStringSubmatch(repo)
	if len(m) != 3 {
		fmt.Fprintf(os.Stderr, constants.FixRepoErrNoVerSuffFmt, repo)
		os.Exit(constants.FixRepoExitNoVersionSuffix)
	}
	n, err := strconv.Atoi(m[2])
	if err != nil || n < 1 {
		fmt.Fprint(os.Stderr, constants.FixRepoErrBadVersion)
		os.Exit(constants.FixRepoExitBadVersion)
	}

	return m[1], n
}
