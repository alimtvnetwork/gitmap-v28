package cmdcommit

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// Dry-run preview for the commit suite.
//
// The preview path is deliberately read-only: it runs only git commands
// that cannot mutate the tree, the index, or any remote (status, rev-parse,
// config --get, rev-list --count). It must never call the staging, commit,
// push, pull, or telemetry-recording helpers, so a --dry-run run can never
// sweep files into a commit the way the pre-fix commit-push path did.

// dryRunFile is one line of `git status --porcelain` output.
type dryRunFile struct {
	Path   string
	Status string
}

// dryRunPlan is everything a commit-suite run would do, gathered without
// changing anything.
type dryRunPlan struct {
	Message         string
	Branch          string
	Head            string
	Remote          string
	RemoteURL       string
	TaskSlug        string
	Push            bool
	PullFirst       bool
	Scoped          bool
	Clean           bool
	HasUpstream     bool
	Unpushed        int
	Changed         []dryRunFile
	WouldStage      []dryRunFile
	Claimed         []string
	UnchangedClaims []string
}

// parsePorcelainStatus turns `git status --porcelain` output into file
// entries. Rename entries ("old -> new") report the new path, which is
// what `git add -A` would stage.
func parsePorcelainStatus(out string) []dryRunFile {
	var files []dryRunFile
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 4 {
			continue
		}
		path := strings.TrimSpace(line[3:])
		if idx := strings.Index(path, " -> "); idx >= 0 {
			path = path[idx+4:]
		}
		path = strings.Trim(path, `"`)
		if path == "" {
			continue
		}
		files = append(files, dryRunFile{Path: path, Status: line[:2]})
	}

	return files
}

// collectDryRunPlan gathers the plan with read-only git commands only.
func collectDryRunPlan(message string, meta CommitMetadata, push, pullFirst bool) dryRunPlan {
	plan := dryRunPlan{
		Message:   message,
		Push:      push,
		PullFirst: pullFirst,
		TaskSlug:  meta.TaskSlug,
	}
	branch, err := getCurrentBranchName()
	if err != nil || strings.TrimSpace(branch) == "" {
		plan.Branch = "HEAD"
	} else {
		plan.Branch = strings.TrimSpace(branch)
	}
	plan.Head = resolveShortHeadSha()
	plan.Remote = resolvePushRemote(plan.Branch)
	if plan.Remote != "" {
		if url, err := execGitOutputCP("config", "--get", "remote."+plan.Remote+".url"); err == nil {
			plan.RemoteURL = strings.TrimSpace(url)
		}
	}
	plan.HasUpstream = hasRemoteTracking(plan.Branch)
	if plan.HasUpstream {
		plan.Unpushed = countUnpushedCommitsCP()
	}
	statusOut, err := execGitOutputCP("status", "--porcelain")
	if err == nil {
		plan.Changed = parsePorcelainStatus(statusOut)
	}
	plan.Clean = len(plan.Changed) == 0
	plan.WouldStage = plan.Changed
	if strings.TrimSpace(meta.TaskSlug) != "" {
		if claimed := scopedStageFiles(meta.TaskSlug); len(claimed) > 0 {
			plan.Scoped = true
			plan.Claimed = claimed
			plan.WouldStage, plan.UnchangedClaims = scopeWouldStage(plan.Changed, claimed)
		}
	}

	return plan
}

// scopeWouldStage restricts the would-stage list to task-claimed files
// that actually changed, mirroring stageClaimedFiles; claimed files with
// no changes are reported separately instead of being silently listed.
func scopeWouldStage(changed []dryRunFile, claimed []string) ([]dryRunFile, []string) {
	changedSet := make(map[string]bool, len(changed))
	for _, f := range changed {
		changedSet[f.Path] = true
	}
	var staged []dryRunFile
	for _, f := range changed {
		if containsClaim(claimed, f.Path) {
			staged = append(staged, f)
		}
	}
	var unchanged []string
	for _, c := range claimed {
		if !changedSet[c] {
			unchanged = append(unchanged, c)
		}
	}

	return staged, unchanged
}

func containsClaim(claimed []string, path string) bool {
	for _, c := range claimed {
		if c == path {
			return true
		}
	}

	return false
}

// resolvePushRemote names the remote a push would target: the branch's
// configured remote, else origin when it exists, else empty.
func resolvePushRemote(branch string) string {
	if branch != "" && branch != "HEAD" {
		if remote, err := execGitOutputCP("config", "--get", "branch."+branch+".remote"); err == nil {
			if name := strings.TrimSpace(remote); name != "" {
				return name
			}
		}
	}
	if _, err := execGitOutputCP("config", "--get", "remote.origin.url"); err == nil {
		return "origin"
	}

	return ""
}

// renderDryRunPlan prints the full plan and states the no-mutation
// guarantee. It writes only to stdout; it stages, commits, pushes,
// pulls, and records nothing.
func renderDryRunPlan(plan dryRunPlan) {
	fmt.Println()
	printPaddedInfo("Dry run: no changes will be made - no commit, no push, no staging, no pull.")
	if plan.PullFirst {
		printPaddedInfo("Would pull: latest changes before commit (not performed).")
	}
	printPaddedInfo("Commit message: %s", plan.Message)
	printPaddedInfo("Branch: %s (HEAD %s)", plan.Branch, plan.Head)
	if plan.Push {
		printPaddedInfo("Push target: %s -> %s", plan.pushTarget(), plan.RemoteURL)
	} else {
		printPaddedInfo("Push target: none - local commit only.")
	}
	if plan.Scoped {
		printPaddedInfo("Staging scoped to task %s (%d claimed file(s)).", plan.TaskSlug, len(plan.Claimed))
	}
	if plan.Clean {
		printPaddedInfo("Working tree clean: nothing to stage or commit.")
		if plan.Push && plan.HasUpstream && plan.Unpushed > 0 {
			printPaddedInfo("Would push: %d unpushed commit(s) to %s (not pushed).", plan.Unpushed, plan.pushTarget())
		} else if plan.Push {
			printPaddedInfo("Would push: nothing - no unpushed commits.")
		}
	} else {
		printPaddedInfo("Files that would be staged (%d):", len(plan.WouldStage))
		for _, f := range plan.WouldStage {
			fmt.Printf("    %s %s\n", f.Status, f.Path)
		}
		for _, c := range plan.UnchangedClaims {
			printPaddedInfo("Claimed but unchanged (would not be staged): %s", c)
		}
		if plan.Push && plan.HasUpstream && plan.Unpushed > 0 {
			printPaddedInfo("Plus %d already-unpushed commit(s) would also go to %s.", plan.Unpushed, plan.pushTarget())
		}
	}
	printPaddedInfo("Tag/release actions: none - the commit suite creates no tags or releases.")
	printPaddedSuccess("Dry run complete: tree and index unchanged.")
}

func (plan dryRunPlan) pushTarget() string {
	if plan.Remote == "" {
		return "(no remote configured) " + plan.Branch
	}

	return plan.Remote + "/" + plan.Branch
}

// previewCommitPush is the early-return dry-run path for every
// commit-push variant (cpf, cpb, cpc, cpr, refactor, commit-push, pcp).
// Callers must invoke it before any staging, pulling, committing, or
// pushing, so the dry-run path can never reach the mutation code.
func previewCommitPush(commitMessage string, meta CommitMetadata, pullFirst bool) *apperror.AppError {
	renderDryRunPlan(collectDryRunPlan(commitMessage, meta, true, pullFirst))

	return nil
}

// previewCommit is the dry-run path for the plain commit/cm command.
func previewCommit(cleanArgs []string, meta CommitMetadata) error {
	message := strings.Join(cleanArgs, " ")
	if len(cleanArgs) > 0 && cleanArgs[0] == "-m" {
		message = strings.Join(cleanArgs[1:], " ")
	}
	if strings.TrimSpace(message) == "" {
		message = "(no message - editor would open)"
	}
	renderDryRunPlan(collectDryRunPlan(message, meta, meta.HasPush, false))

	return nil
}
