package committransfer

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/committransfer/graph"
)

// RunRight is the public entry point for `commit-right`.
func RunRight(sourceDir, targetDir string, opts Options) error {
	err := runOneDirection(sourceDir, targetDir, opts)
	if err == nil {
		RenderRunGraph(targetDir, opts)
	}

	return err
}

// RunPR is the public entry point for `pr`.
func RunPR(sourceDir, targetDir string, opts Options) error {
	opts.PRMode = "merges"
	err := runOneDirection(sourceDir, targetDir, opts)
	if err == nil {
		RenderRunGraph(targetDir, opts)
	}

	return err
}

// maybePush runs `git push` unless --no-push is set, the target is not
// a git repo, or there are no new commits. Returns true on success.
func maybePush(targetDir string, opts Options, newCount int) bool {
	defer RenderRunGraph(targetDir, opts)
	if opts.NoPush || opts.NoCommit || newCount == 0 || opts.DryRun {
		return false
	}
	if _, err := pushHEAD(targetDir); err != nil {
		fmt.Fprintf(os.Stderr, "%s push failed: %v\n", opts.LogPrefix, err)

		return false
	}

	return true
}

// RenderRunGraph extracts recent graph events from targetDir and renders the visual execution graph.
func RenderRunGraph(targetDir string, opts Options) {
	if opts.DryRun {
		return
	}
	events := CollectGraphEvents(targetDir, 10)
	if len(events) == 0 {
		return
	}
	out := graph.RenderExecutionGraph(events)
	if out != "" {
		fmt.Fprintf(os.Stdout, "\n%s Visual Execution Graph:\n%s\n", opts.LogPrefix, out)
	}
}

// CollectGraphEvents reads recent commit history from targetDir and maps them to graph events.
func CollectGraphEvents(targetDir string, limit int) []graph.GraphEvent {
	if limit <= 0 {
		limit = 10
	}
	out, err := gitOut(targetDir, "log", fmt.Sprintf("-n%d", limit), "--reverse", "--format=%H%x1f%h%x1f%s%x1f%P")
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	events := parseLogLines(targetDir, out)

	return enrichBranchEvents(targetDir, events)
}

func parseLogLines(targetDir, out string) []graph.GraphEvent {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	var events []graph.GraphEvent
	for _, line := range lines {
		if ev := parseCommitLineToEvent(targetDir, line); ev != nil {
			events = append(events, *ev)
		}
	}

	return events
}

func parseCommitLineToEvent(targetDir, line string) *graph.GraphEvent {
	parts := strings.Split(line, "\x1f")
	if len(parts) < 3 {
		return nil
	}
	parents := ""
	if len(parts) >= 4 {
		parents = parts[3]
	}
	isMerge := len(strings.Fields(parents)) > 1

	return buildParsedEvent(targetDir, parts[0], parts[1], parts[2], isMerge)
}

func buildParsedEvent(targetDir, sha, shortSha, subject string, isMerge bool) *graph.GraphEvent {
	tag := resolveCommitTag(targetDir, sha)
	branch := resolveCommitBranch(subject, shortSha, isMerge)

	return &graph.GraphEvent{
		CommitSha:  sha,
		BranchName: branch,
		IsMerge:    isMerge,
		PRNumber:   extractPRNumber(subject),
		ReleaseTag: tag,
		Message:    subject,
	}
}

func resolveCommitTag(targetDir, sha string) string {
	tag, err := gitOut(targetDir, "tag", "--points-at", sha)
	if err == nil && strings.TrimSpace(tag) != "" {
		return strings.Split(strings.TrimSpace(tag), "\n")[0]
	}

	return ""
}

func resolveCommitBranch(subject, shortSha string, isMerge bool) string {
	if !isMerge {
		return "main"
	}
	if branch := parseBranchFromSubject(subject); branch != "" {
		return branch
	}
	if prNum := extractPRNumber(subject); prNum > 0 {
		return fmt.Sprintf("pr/%d", prNum)
	}

	return "main"
}

func enrichBranchEvents(targetDir string, events []graph.GraphEvent) []graph.GraphEvent {
	for i := range events {
		enrichSingleMergeEvent(targetDir, &events[i], events)
	}

	return events
}

func enrichSingleMergeEvent(targetDir string, ev *graph.GraphEvent, events []graph.GraphEvent) {
	if !ev.IsMerge {
		return
	}
	parents := getCommitParents(targetDir, ev.CommitSha)
	if len(parents) < 2 {
		return
	}
	branch := extractBranchName(ev.Message, ev.CommitSha, "")
	markFeatureCommits(targetDir, parents[0], parents[1], branch, events)
}

func markFeatureCommits(targetDir, p1, p2, branchName string, events []graph.GraphEvent) {
	shas, err := gitOut(targetDir, "rev-list", p1+".."+p2)
	if err != nil || strings.TrimSpace(shas) == "" {
		return
	}
	shaSet := make(map[string]bool)
	for _, s := range strings.Split(strings.TrimSpace(shas), "\n") {
		shaSet[s] = true
	}
	for j := range events {
		if shaSet[events[j].CommitSha] {
			events[j].BranchName = branchName
		}
	}
}
