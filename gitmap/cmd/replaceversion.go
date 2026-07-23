package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// remoteSlugRe matches the trailing `<base>-vN` segment of a git remote
// URL. The base may contain dashes; only the last `-vN` counts.
var remoteSlugRe = regexp.MustCompile(`^(?P<base>.+)-v(?P<num>\d+)$`)

// detectVersion returns (base, K) parsed from `git remote get-url
// origin`. Exits 1 with a descriptive error on any failure mode.
func detectVersion() (string, int) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrReplaceNoRemote, err)
		os.Exit(1)
	}
	slug := slugFromRemote(strings.TrimSpace(string(out)))
	m := remoteSlugRe.FindStringSubmatch(slug)
	if m == nil {
		fmt.Fprintf(os.Stderr, constants.ErrReplaceVersionParse, slug)
		os.Exit(1)
	}
	num, _ := strconv.Atoi(m[2])
	return m[1], num
}

// slugFromRemote extracts the last path segment from a remote URL,
// trimming a trailing `.git`.
func slugFromRemote(url string) string {
	url = strings.TrimSuffix(url, ".git")
	if i := strings.LastIndex(url, "/"); i >= 0 {
		return url[i+1:]
	}
	if i := strings.LastIndex(url, ":"); i >= 0 {
		return url[i+1:]
	}
	return url
}

// versionTargets returns the ascending list of target ints to replace.
// When n==0 the caller wants `all` mode (1..K-1).
func versionTargets(k, n int) []int {
	if n == 0 || n >= k {
		n = k - 1
	}
	if n < 1 {
		return nil
	}
	start := k - n
	if start < 1 {
		start = 1
	}
	out := make([]int, 0, k-start)
	for i := start; i < k; i++ {
		out = append(out, i)
	}
	return out
}

// pairsForTarget returns the two replace pairs for one target version.
func pairsForTarget(base string, target, k int) []replacePair {
	return []replacePair{
		{old: fmt.Sprintf("%s-v%d", base, target), new: fmt.Sprintf("%s-v%d", base, k)},
		{old: fmt.Sprintf("%s/v%d", base, target), new: fmt.Sprintf("%s/v%d", base, k)},
	}
}
