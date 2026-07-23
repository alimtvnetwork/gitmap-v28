// Package gitutil — latest-branch resolve operations.
package gitutil

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ReadBranchTips reads commit metadata for each remote ref.
func ReadBranchTips(refs []string) ([]RemoteBranchInfo, error) {
	var items []RemoteBranchInfo
	for _, ref := range refs {
		info, ok := readSingleTip(ref)
		if ok {
			items = append(items, info)
		}
	}
	if len(items) == 0 {

		return nil, fmt.Errorf(constants.ErrLatestBranchNoCommits)
	}

	return items, nil
}

// readSingleTip reads commit date, SHA, and subject for one ref.
func readSingleTip(ref string) (RemoteBranchInfo, bool) {
	cmd := exec.Command(constants.GitBin, constants.GitLog, "-1", constants.GitLogTipFormat, ref)
	out, err := cmd.Output()
	if err != nil {

		return RemoteBranchInfo{}, false
	}

	return parseTipLine(strings.TrimSpace(string(out)), ref)
}

// parseTipLine parses a "date|sha|subject" line into RemoteBranchInfo.
func parseTipLine(line, ref string) (RemoteBranchInfo, bool) {
	parts := strings.SplitN(line, constants.GitLogDelimiter, constants.GitLogFieldCount)
	if len(parts) != constants.GitLogFieldCount {

		return RemoteBranchInfo{}, false
	}
	t, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {

		return RemoteBranchInfo{}, false
	}

	return RemoteBranchInfo{
		RemoteRef:  ref,
		CommitDate: t,
		Sha:        parts[1],
		Subject:    parts[2],
	}, true
}

// ResolvePointsAt returns branch names that point exactly at sha.
func ResolvePointsAt(sha, remote string) []string {
	cmd := exec.Command(constants.GitBin, constants.GitForEachRef,
		fmt.Sprintf(constants.GitPointsAtFmt, sha),
		fmt.Sprintf(constants.GitRefsRemotesFmt, remote),
		constants.GitFormatRefnameShort)
	out, err := cmd.Output()
	if err != nil {

		return nil
	}

	return parseUniqueNames(strings.TrimSpace(string(out)), remote)
}

// parseUniqueNames extracts unique branch names from for-each-ref output.
func parseUniqueNames(output, remote string) []string {
	prefix := remote + "/"
	headRef := remote + "/HEAD"
	var names []string
	seen := map[string]bool{}
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 || trimmed == headRef {
			continue
		}
		name := strings.TrimPrefix(trimmed, prefix)
		if seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}

	return names
}

// ResolveContains returns branch names whose history contains sha.
func ResolveContains(sha, remote string) []string {
	cmd := exec.Command(constants.GitBin, constants.GitBranch, constants.GitArgRemote, constants.GitArgContains, sha)
	out, err := cmd.Output()
	if err != nil {

		return nil
	}

	return parseContainsNames(string(out), remote)
}

// parseContainsNames extracts unique branch names from --contains output.
func parseContainsNames(output, remote string) []string {
	prefix := remote + "/"
	var names []string
	seen := map[string]bool{}
	for _, line := range strings.Split(output, "\n") {
		name := extractContainsName(line, prefix)
		if len(name) == 0 || seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}

	return names
}

// extractContainsName extracts a branch name from a single --contains line.
func extractContainsName(line, prefix string) string {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 || strings.Contains(trimmed, constants.HeadPointer) {

		return ""
	}
	if strings.HasPrefix(trimmed, prefix) {
		name := strings.TrimPrefix(trimmed, prefix)
		if name == constants.GitHEAD {

			return ""
		}

		return name
	}

	return ""
}
