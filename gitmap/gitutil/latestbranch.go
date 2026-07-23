// Package gitutil — latest-branch core operations.
package gitutil

import (
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// RemoteBranchInfo holds commit metadata for a remote-tracking branch.
type RemoteBranchInfo struct {
	RemoteRef  string
	CommitDate time.Time
	Sha        string
	Subject    string
}

// IsInsideWorkTree checks if the current directory is inside a git repo.
func IsInsideWorkTree() bool {
	cmd := exec.Command(constants.GitBin, constants.GitRevParse, constants.GitArgInsideWorkTree)
	err := cmd.Run()

	return err == nil
}

// FetchAllPrune runs git fetch --all --prune.
func FetchAllPrune() error {
	cmd := exec.Command(constants.GitBin, constants.GitFetch, constants.GitArgAll, constants.GitArgPrune)

	return cmd.Run()
}

// ListRemoteBranches returns trimmed remote-tracking branch names,
// excluding HEAD pointer lines.
func ListRemoteBranches() ([]string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitBranch, constants.GitArgRemote)
	out, err := cmd.Output()
	if err != nil {

		return nil, err
	}

	return parseRemoteBranchLines(string(out)), nil
}

// parseRemoteBranchLines extracts branch refs from git branch -r output.
func parseRemoteBranchLines(output string) []string {
	var refs []string
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		if strings.Contains(trimmed, constants.HeadPointer) {
			continue
		}
		refs = append(refs, trimmed)
	}

	return refs
}

// FilterByRemote keeps only refs starting with "<remote>/".
func FilterByRemote(refs []string, remote string) []string {
	prefix := remote + "/"
	var filtered []string
	for _, r := range refs {
		if strings.HasPrefix(r, prefix) {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

// FilterByPattern keeps only refs whose branch name matches
// the given glob or substring pattern.
func FilterByPattern(refs []string, pattern string) []string {
	var filtered []string
	for _, r := range refs {
		name := StripRemotePrefix(r)
		if matchesPattern(name, pattern) {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

// matchesPattern checks glob first, then substring fallback.
func matchesPattern(name, pattern string) bool {
	matched, err := filepath.Match(pattern, name)
	if err == nil && matched {

		return true
	}

	return strings.Contains(name, pattern)
}

// SortByDateDesc sorts items by CommitDate descending.
func SortByDateDesc(items []RemoteBranchInfo) {
	sort.Slice(items, func(i, j int) bool {

		return items[i].CommitDate.After(items[j].CommitDate)
	})
}

// SortByNameAsc sorts items by branch name ascending.
func SortByNameAsc(items []RemoteBranchInfo) {
	sort.Slice(items, func(i, j int) bool {

		return items[i].RemoteRef < items[j].RemoteRef
	})
}

// StripRemotePrefix removes the "<remote>/" prefix from a ref.
func StripRemotePrefix(ref string) string {
	idx := strings.Index(ref, "/")
	if idx >= 0 {

		return ref[idx+1:]
	}

	return ref
}

// TruncSha returns the first N characters of a SHA for display.
func TruncSha(sha string) string {
	if len(sha) > constants.ShaDisplayLength {

		return sha[:constants.ShaDisplayLength]
	}

	return sha
}
