package committransfer

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/committransfer/prdesc"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
	"github.com/alimtvnetwork/gitmap-v28/cli/prdb"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func isPRModeActive(prMode string) bool {
	return prMode != "" && prMode != "off"
}

func isPRRouteEligible(sourceDir, sha, subject, prMode string) bool {
	if !isPRModeActive(prMode) {
		return false
	}
	if isMergeCommit(sourceDir, sha) || prMode == "all" {
		return true
	}
	lower := strings.ToLower(subject)
	isTag := strings.Contains(lower, "tag:") || strings.Contains(lower, "release ") || strings.Contains(lower, "version ")
	isRelease := strings.Contains(lower, "chore(release):") || strings.Contains(lower, "release v")

	return (prMode == "tags" && isTag) || (prMode == "release" && isRelease)
}

func getCommitParents(dir, sha string) []string {
	out, err := gitOut(dir, "rev-parse", sha+"^@")
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}

	return strings.Split(strings.TrimSpace(out), "\n")
}

func isMergeCommit(dir, sha string) bool {
	parents := getCommitParents(dir, sha)

	return len(parents) > 1
}

func runGit(dir string, args ...string) result.Result[string] {
	out, err := gitOut(dir, args...)
	if err != nil {
		return result.Fail[string](apperror.WrapSimple(err, "git "+strings.Join(args, " ")))
	}

	return result.Ok(out)
}

func extractPRNumber(subject string) int {
	re := lazyregex.New(`(?i)#(\d+)`).Compiled()
	if re == nil {
		return 0
	}
	matches := re.FindStringSubmatch(subject)
	if len(matches) < 2 {
		return 0
	}
	num, _ := strconv.Atoi(matches[1])

	return num
}

func parseBranchFromSubject(subject string) string {
	reFrom := lazyregex.New(`(?i)from\s+([^\s]+)`).Compiled()
	if reFrom != nil {
		if m := reFrom.FindStringSubmatch(subject); len(m) > 1 {
			return sanitizeBranchName(m[1])
		}
	}
	reBranch := lazyregex.New(`(?i)merge branch '([^']+)'`).Compiled()
	if reBranch != nil {
		if m := reBranch.FindStringSubmatch(subject); len(m) > 1 {
			return sanitizeBranchName(m[1])
		}
	}

	return ""
}

func sanitizeBranchName(raw string) string {
	parts := strings.Split(raw, "/")
	clean := parts[len(parts)-1]
	if !strings.HasPrefix(raw, "feature/") && !strings.HasPrefix(raw, "pr/") {
		return "feature/" + clean
	}

	return raw
}

func extractBranchName(subject, shortSHA, defaultSlug string) string {
	if branch := parseBranchFromSubject(subject); branch != "" {
		return branch
	}
	if prNum := extractPRNumber(subject); prNum > 0 {
		return fmt.Sprintf("pr/%d", prNum)
	}

	return fmt.Sprintf("feature/%s-%s", shortSHA, defaultSlug)
}

func fetchIngressSHAs(sourceDir, sha string) []string {
	parents := getCommitParents(sourceDir, sha)
	if len(parents) < 2 {
		return []string{sha}
	}
	rangeSpec := fmt.Sprintf("%s..%s", parents[0], parents[1])
	out, err := gitOut(sourceDir, "rev-list", "--reverse", rangeSpec)
	if err != nil || strings.TrimSpace(out) == "" {
		return []string{sha}
	}

	return strings.Split(strings.TrimSpace(out), "\n")
}

func replaySingleIngress(plan ReplayPlan, sha string, opts Options) result.Result[prdesc.PRCommitInfo] {
	sub, body, author, short, when, err := readCommit(plan.SourceDir, sha)
	if err != nil {
		return result.Fail[prdesc.PRCommitInfo](apperror.WrapSimple(err, "read ingress commit"))
	}
	if cErr := checkoutDetached(plan.SourceDir, sha); cErr != nil {
		return result.Fail[prdesc.PRCommitInfo](apperror.WrapSimple(cErr, "checkout detached"))
	}
	if sErr := snapshotCopy(plan.SourceDir, plan.TargetDir, opts); sErr != nil {
		return result.Fail[prdesc.PRCommitInfo](apperror.WrapSimple(sErr, "snapshot copy"))
	}
	_ = addAll(plan.TargetDir)
	commitIngressChanges(plan.TargetDir, sub, body, author, short, when, opts)

	return result.Ok(prdesc.PRCommitInfo{ShortSHA: short, Author: author, Subject: sub})
}

func commitIngressChanges(dir, sub, body, author, short string, when time.Time, opts Options) {
	if !hasStagedChanges(dir) {
		return
	}
	cleaned := CleanMessage(sub, body, opts.Message, short, when).Final
	_, _ = commitWithEnv(dir, cleaned, author, when)
}

func replayIngressCommits(plan ReplayPlan, commit SourceCommit, opts Options) result.Result[[]prdesc.PRCommitInfo] {
	shas := fetchIngressSHAs(plan.SourceDir, commit.SHA)
	var infos []prdesc.PRCommitInfo
	for _, sha := range shas {
		itemRes := replaySingleIngress(plan, sha, opts)
		if itemRes.IsFailure() {
			return result.Fail[[]prdesc.PRCommitInfo](itemRes.Err)
		}
		infos = append(infos, itemRes.Value)
	}

	return result.Ok(infos)
}

func computeDiffSummary(targetDir, baseRef string) prdesc.PRFileDiffSummary {
	out, err := gitOut(targetDir, "diff", "--numstat", baseRef+"...HEAD")
	if err != nil || strings.TrimSpace(out) == "" {
		return prdesc.PRFileDiffSummary{FilesChanged: 1}
	}
	var diff prdesc.PRFileDiffSummary
	lines := strings.Split(strings.TrimSpace(out), "\n")
	diff.FilesChanged = len(lines)
	for _, line := range lines {
		f := strings.Fields(line)
		if len(f) >= 2 {
			add, _ := strconv.Atoi(f[0])
			del, _ := strconv.Atoi(f[1])
			diff.Additions += add
			diff.Deletions += del
		}
	}

	return diff
}

func computeComponentImpacts(targetDir, baseRef string) []prdesc.PRComponentImpact {
	out, err := gitOut(targetDir, "diff", "--name-only", baseRef+"...HEAD")
	if err != nil || strings.TrimSpace(out) == "" {
		return []prdesc.PRComponentImpact{{Component: "Core", Impact: "Replayed", Files: 1}}
	}
	counts := make(map[string]int)
	for _, file := range strings.Split(strings.TrimSpace(out), "\n") {
		parts := strings.Split(filepath.ToSlash(file), "/")
		counts[parts[0]]++
	}
	var impacts []prdesc.PRComponentImpact
	for comp, count := range counts {
		impacts = append(impacts, prdesc.PRComponentImpact{Component: comp, Impact: "Modified", Files: count})
	}

	return impacts
}

func assemblePRMetadata(
	targetDir string,
	commit SourceCommit,
	branchName, mainline string,
	commits []prdesc.PRCommitInfo,
) prdesc.PRMetadata {
	return prdesc.PRMetadata{
		PRNumber:         extractPRNumber(commit.Subject),
		Title:            commit.Subject,
		SourceBranch:     branchName,
		TargetBranch:     mainline,
		Author:           commit.Author,
		Timestamp:        commit.AuthorAt,
		ExecutiveSummary: commit.Body,
		Commits:          commits,
		Impacts:          computeComponentImpacts(targetDir, mainline),
		DiffSummary:      computeDiffSummary(targetDir, mainline),
	}
}

func nextPRNumber(db *prdb.PrSplitDb) int {
	var maxNum int
	row := db.Conn().QueryRow("SELECT COALESCE(MAX(PrNumber), 0) + 1 FROM PullRequest")
	if err := row.Scan(&maxNum); err != nil || maxNum <= 0 {
		return 1
	}

	return maxNum
}

func insertPRRecords(db *prdb.PrSplitDb, branchName, mainline, desc string, prNum int, title string) {
	_ = db.CreatePullRequest(prdb.PullRequestRecord{
		PrNumber: prNum, Title: title, Description: desc,
		SourceBranch: branchName, TargetBranch: mainline, Status: "open",
	})
	_ = db.UpsertPrBranch(prdb.PrBranchRecord{
		BranchName: branchName, BranchType: "feature", IsMerged: false,
	})
}

func recordPRInDatabase(targetDir, branchName, mainline, desc string, meta *prdesc.PRMetadata) result.Result[*prdb.PrSplitDb] {
	slug := prdb.SanitizeRepoSlug(filepath.Base(targetDir))
	dbRes := prdb.OpenPrSplitDb(slug, targetDir)
	if dbRes.IsFailure() {
		return dbRes
	}
	db := dbRes.Value
	if meta.PRNumber <= 0 {
		meta.PRNumber = nextPRNumber(db)
	}
	insertPRRecords(db, branchName, mainline, desc, meta.PRNumber, meta.Title)

	return result.Ok(db)
}

func finalizePRInDatabase(db *prdb.PrSplitDb, prNum int, branchName, mergeSha string) {
	if db == nil {
		return
	}
	_ = db.UpdatePullRequestStatus(prNum, "merged", mergeSha)
	_ = db.UpsertPrBranch(prdb.PrBranchRecord{
		BranchName: branchName, BranchType: "feature", IsMerged: true,
		MergedAt: time.Now().UTC().Unix(),
	})
}

func formatPRMergeMessage(meta prdesc.PRMetadata, branchName string) string {
	if meta.Title != "" && strings.HasPrefix(meta.Title, "Merge ") {
		return meta.Title
	}
	if meta.PRNumber > 0 {
		return fmt.Sprintf("Merge pull request #%d from %s\n\n%s", meta.PRNumber, branchName, meta.Title)
	}

	return fmt.Sprintf("Merge branch '%s'\n\n%s", branchName, meta.Title)
}

func mergeFeatureBranch(targetDir, branchName, msg string) result.Result[string] {
	mergeRes := runGit(targetDir, "merge", "--no-ff", branchName, "-m", msg)
	if mergeRes.IsFailure() {
		return mergeRes
	}

	return runGit(targetDir, "rev-parse", "HEAD")
}

func executePRReplayAndMerge(plan ReplayPlan, commit SourceCommit, opts Options, mainline, branchName string) result.Result[string] {
	ingressRes := replayIngressCommits(plan, commit, opts)
	if ingressRes.IsFailure() {
		_ = runGit(plan.TargetDir, "checkout", mainline)

		return result.Fail[string](ingressRes.Err)
	}
	meta := assemblePRMetadata(plan.TargetDir, commit, branchName, mainline, ingressRes.Value)
	desc := prdesc.GeneratePRDescription(meta)
	dbRes := recordPRInDatabase(plan.TargetDir, branchName, mainline, desc, &meta)
	_ = runGit(plan.TargetDir, "checkout", mainline)
	mergeMsg := formatPRMergeMessage(meta, branchName)
	mergeRes := mergeFeatureBranch(plan.TargetDir, branchName, mergeMsg)
	if dbRes.IsSuccess() && mergeRes.IsSuccess() {
		defer dbRes.Value.Close()
		finalizePRInDatabase(dbRes.Value, meta.PRNumber, branchName, mergeRes.Value)
	}

	return mergeRes
}

// ProcessPR orchestrates local Git PR simulation for a merge or PR-eligible commit.
func ProcessPR(plan ReplayPlan, commit SourceCommit, opts Options) result.Result[string] {
	mainlineRes := runGit(plan.TargetDir, "rev-parse", "--abbrev-ref", "HEAD")
	if mainlineRes.IsFailure() {
		return mainlineRes
	}
	mainline := mainlineRes.Value
	branchName := extractBranchName(commit.Subject, commit.ShortSHA, filepath.Base(plan.TargetDir))
	if res := runGit(plan.TargetDir, "checkout", "-B", branchName); res.IsFailure() {
		return res
	}
	res := executePRReplayAndMerge(plan, commit, opts, mainline, branchName)
	if res.IsSuccess() {
		pterm.Success.Printf("Merged PR (%s) into %s\n", branchName, mainline)
	}

	return res
}
