package cmdrelease

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/release"
)

func filterTagsBetweenBounds(
	allTags []string,
	lo, hi release.Version,
	pinned []string,
) ([]string, []string) {
	var toDelete []string
	var kept []string

	for _, tag := range allTags {
		v, err := release.Parse(tag)
		if err != nil {
			kept = append(kept, tag)
			continue
		}

		if isVersionEqual(v, lo) || isVersionEqual(v, hi) {
			kept = append(kept, tag)
			continue
		}

		if IsTagPinned(tag, pinned) {
			kept = append(kept, tag)
			continue
		}

		if isTagStrictlyBetween(v, lo, hi) {
			toDelete = append(toDelete, tag)
		} else {
			kept = append(kept, tag)
		}
	}

	return toDelete, kept
}

func isVersionEqual(a, b release.Version) bool {
	return a.Major == b.Major && a.Minor == b.Minor && a.Patch == b.Patch && a.PreRelease == b.PreRelease
}

func isTagStrictlyBetween(v, lo, hi release.Version) bool {
	if v.GreaterThan(lo) && hi.GreaterThan(v) {
		return true
	}

	isSameMinor := v.Major == lo.Major && v.Minor == lo.Minor
	if isSameMinor && v.Patch == lo.Patch && v.IsPreRelease() {
		return true
	}
	if isSameMinor && v.Patch > lo.Patch && (v.Major < hi.Major || v.Minor < hi.Minor || v.Patch < hi.Patch) {
		return true
	}

	return false
}

func printAutoDryRunEvaluation(lo, hi release.Version, toDelete, toKeep, pinned []string) {
	fmt.Printf("\x1b[1;36m▶ release remove-between\x1b[0m: [%s .. %s]\n", lo.String(), hi.String())
	fmt.Printf("  • Milestone anchors preserved : \x1b[1;32m%s\x1b[0m and \x1b[1;32m%s\x1b[0m\n", lo.String(), hi.String())
	if len(pinned) > 0 {
		fmt.Printf("  • Pinned tags preserved       : \x1b[1;33m%s\x1b[0m\n", strings.Join(pinned, ", "))
	}
	fmt.Printf("  • Candidate tags to remove    : \x1b[1;31m%d\x1b[0m\n", len(toDelete))
	fmt.Printf("  • Total tags retained         : \x1b[1;32m%d\x1b[0m\n", len(toKeep))

	if len(toDelete) > 0 {
		fmt.Println("\n\x1b[1;33m📋 Candidate intermediate tag(s) queued for removal:\x1b[0m")
		for _, t := range toDelete {
			fmt.Printf("    \x1b[31m✖ %s\x1b[0m\n", t)
		}
	}
}

func confirmBatchDeletion(count int, lo, hi string) bool {
	fmt.Printf("\n\x1b[1;31m⚠ Permanently delete %d intermediate tag(s) between %s and %s locally and on origin? [y/N]: \x1b[0m", count, lo, hi)
	reader := bufio.NewReader(os.Stdin)
	reply, err := reader.ReadString('\n')
	if err != nil && len(reply) == 0 {
		return false
	}
	t := strings.TrimSpace(reply)
	return strings.EqualFold(t, "y") || strings.EqualFold(t, "yes")
}

func applyBatchTagRemoval(tags []string, keepRemote bool) []string {
	var done []string

	for _, t := range tags {
		_ = exec.Command("gh", "release", "delete", t, "-y", "--cleanup-tag=false").Run()
	}
	done = append(done, "gh releases")

	if !keepRemote {
		done = append(done, deleteRemoteOriginTags(tags)...)
	}

	args := append([]string{"tag", "-d"}, tags...)
	_ = exec.Command("git", args...).Run()
	done = append(done, "local git tags")

	for _, t := range tags {
		v, err := release.Parse(t)
		if err == nil {
			jsonPath := filepath.Join(constants.DefaultReleaseDir, v.ReleaseFileName())
			_ = os.Remove(jsonPath)
		}
	}
	done = append(done, "metadata sidecars")

	return done
}

func deleteRemoteOriginTags(tags []string) []string {
	batchSize := 40
	for i := 0; i < len(tags); i += batchSize {
		end := min(i+batchSize, len(tags))
		chunk := tags[i:end]
		args := append([]string{"push", "origin", "--delete"}, chunk...)
		_ = exec.Command("git", args...).Run()
	}

	return []string{"remote origin tags"}
}

func optimizeRepositoryPostPrune(keepRemote bool) {
	fmt.Println("  • Pruning references & optimizing repository pack files...")
	if !keepRemote {
		_ = exec.Command("git", "fetch", "--prune", "--prune-tags", "origin").Run()
	}
	_ = exec.Command("git", "gc", "--prune=now").Run()
}
