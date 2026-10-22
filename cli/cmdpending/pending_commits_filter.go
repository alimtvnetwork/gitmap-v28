package cmdpending

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func buildPendingCommitsPayload(inspected []RepoPendingCommitRecord, opts PendingCommitsOptions) PendingCommitsPayload {
	dirtyCount, totalUncommitted, totalUnpushed := aggregatePendingCounts(inspected)
	displayed := filterDisplayedPendingCommits(inspected, opts.IsDirtyOnly)
	sortPendingCommits(displayed, opts.SortMode)

	return PendingCommitsPayload{
		Timestamp:             time.Now().UTC(),
		TotalReposScanned:     len(inspected),
		TotalDirtyRepos:       dirtyCount,
		TotalUncommittedFiles: totalUncommitted,
		TotalUnpushedCommits:  totalUnpushed,
		SortMode:              opts.SortMode,
		DetailMode:            opts.DetailMode,
		Repositories:          displayed,
	}
}

func aggregatePendingCounts(inspected []RepoPendingCommitRecord) (int, int, int) {
	dirtyCount, totalUncommitted, totalUnpushed := 0, 0, 0
	for _, item := range inspected {
		if item.IsDirty {
			dirtyCount++
		}
		totalUncommitted += item.TotalUncommitted
		totalUnpushed += item.UnpushedCommitsCount
	}

	return dirtyCount, totalUncommitted, totalUnpushed
}

func filterDisplayedPendingCommits(inspected []RepoPendingCommitRecord, isDirtyOnly bool) []RepoPendingCommitRecord {
	var displayed []RepoPendingCommitRecord
	for _, item := range inspected {
		if isPendingCommitDisplayed(item, isDirtyOnly) {
			displayed = append(displayed, item)
		}
	}

	return displayed
}

func renderLocalPendingCommitsTerminal(payload PendingCommitsPayload, totalScanned int, inspected []RepoPendingCommitRecord, opts PendingCommitsOptions) {
	renderPendingCommitsTerminal(payload, totalScanned, inspected)
	if opts.DetailMode == "all" {
		renderPendingCommitsDetailsAll(payload.Repositories)
	}
	if opts.DetailMode == "1-by-1" {
		renderPendingCommitsInteractive(payload.Repositories)
	}
}

func filterWorkspaceRepositories(records []model.ScanRecord, targetRepo string) ([]model.ScanRecord, error) {
	if targetRepo == "" || strings.EqualFold(targetRepo, "all") {
		return records, nil
	}
	var matched []model.ScanRecord
	for _, r := range records {
		if strings.EqualFold(r.RepoName, targetRepo) ||
			strings.EqualFold(r.Slug, targetRepo) ||
			strings.EqualFold(r.RelativePath, targetRepo) {
			matched = append(matched, r)
		}
	}
	if len(matched) == 0 {
		return nil, apperror.NewSimple(fmt.Sprintf("repository '%s' not found", targetRepo), "E9002")
	}

	return matched, nil
}

func isPendingCommitDisplayed(item RepoPendingCommitRecord, isDirtyOnly bool) bool {
	if !isDirtyOnly {
		return true
	}

	return item.IsDirty || item.HasUnpushed
}

func ResolveWorkspaceRepositories(cwd string) []model.ScanRecord {
	targets := cmdpull.ResolvePullDirectoryTargets(cwd)
	if len(targets) > 0 {
		return targets
	}
	gitPath := filepath.Join(cwd, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		base := filepath.Base(cwd)
		return []model.ScanRecord{
			{
				RepoName:     base,
				Slug:         base,
				AbsolutePath: cwd,
				RelativePath: ".",
			},
		}
	}

	return nil
}

func buildRepoRemediationOptions(rec RepoPendingCommitRecord) []RemediationOption {
	if !rec.IsDirty {
		return []RemediationOption{}
	}
	targetPath := rec.RelativePath
	if targetPath == "" || targetPath == "." || looksLikeHashOrOpaqueID(targetPath) {
		targetPath = rec.RepoName
	}

	return createRemediationOptionsPair(targetPath)
}

func createRemediationOptionsList(targetPath string) []RemediationOption {
	target := quoteRemediationTarget(targetPath)

	return []RemediationOption{
		{OptionNumber: 1, Label: "Commit & Push WIP", Command: fmt.Sprintf(`gitmap sends cp %s "wip: save changes"`, target)},
		{OptionNumber: 2, Label: "Stash WIP", Command: fmt.Sprintf(`gitmap stash %s`, target)},
		{OptionNumber: 3, Label: "All Repos (Commit & Push)", Command: `gitmap sends cp all "wip: save changes"`},
	}
}

func createRemediationOptionsPair(targetPath string) []RemediationOption {
	return createRemediationOptionsList(targetPath)
}

func quoteRemediationTarget(s string) string {
	if strings.Contains(s, " ") && !strings.HasPrefix(s, "\"") {
		return fmt.Sprintf(`"%s"`, s)
	}

	return s
}

func sortPendingCommits(records []RepoPendingCommitRecord, sortMode string) {
	switch sortMode {
	case "count":
		sortByCount(records)
	case "name":
		sortByName(records)
	default:
		sortByPriority(records)
	}
}

func sortByCount(records []RepoPendingCommitRecord) {
	sort.SliceStable(records, func(i, j int) bool {
		countI := resolvePendingSortCount(records[i])
		countJ := resolvePendingSortCount(records[j])
		if countI != countJ {
			return countI > countJ
		}

		return strings.ToLower(records[i].RepoName) < strings.ToLower(records[j].RepoName)
	})
}

func sortByName(records []RepoPendingCommitRecord) {
	sort.SliceStable(records, func(i, j int) bool {
		return strings.ToLower(records[i].RepoName) < strings.ToLower(records[j].RepoName)
	})
}

func sortByPriority(records []RepoPendingCommitRecord) {
	sort.SliceStable(records, func(i, j int) bool {
		pI := calculatePriorityRank(records[i])
		pJ := calculatePriorityRank(records[j])
		if pI != pJ {
			return pI < pJ
		}

		return strings.ToLower(records[i].RepoName) < strings.ToLower(records[j].RepoName)
	})
}

func calculatePriorityRank(r RepoPendingCommitRecord) int {
	if r.IsDirty && r.HasUnpushed {
		return 1
	}
	if r.IsDirty {
		return 2
	}
	if r.HasUnpushed {
		return 3
	}

	return 4
}

func resolvePendingSortCount(r RepoPendingCommitRecord) int {
	total := r.TotalUncommitted
	if total == 0 {
		total = r.UntrackedFilesCount + r.ModifiedFilesCount + r.StagedFilesCount
	}

	return total + r.UnpushedCommitsCount
}

func emitPendingCommitsJSON(payload PendingCommitsPayload) error {
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal pending commits json")
	}
	fmt.Println(string(b))

	return nil
}

func renderPendingCommitsTerminal(payload PendingCommitsPayload, totalScanned int, allInspected []RepoPendingCommitRecord) {
	border := constants.ColorCyan
	reset := constants.ColorReset
	renderTableBanner(border, reset)
	renderTableSummary(payload, border, reset)
	renderTableHeader(border, reset)
	renderTableRows(payload.Repositories, border, reset)
	renderCleanSummaryIfNeeded(payload, totalScanned, allInspected, border, reset)
	renderTableFooter(border, reset)
}

func renderPendingCommitsDetailsAll(records []RepoPendingCommitRecord) {
	fmt.Println()
	fmt.Printf("  %sDetailed Pending Breakdown:%s\n", constants.ColorBold, constants.ColorReset)
	for _, rec := range records {
		if rec.IsClean {
			continue
		}
		renderPendingCommitDetailRecord(rec)
	}
}

func renderPendingCommitDetailRecord(rec RepoPendingCommitRecord) {
	fmt.Printf("\n  %s● %s%s (%s)\n", constants.ColorCyan, rec.RepoName, constants.ColorReset, rec.CurrentBranch)
	printDetailRecordFiles(rec.PendingFiles)
	printDetailRecordSHAs(rec.UnpushedCommitSHAs)
}

func printDetailRecordFiles(files []string) {
	if len(files) == 0 {
		return
	}
	fmt.Printf("    %sModified Files (%d):%s\n", constants.ColorYellow, len(files), constants.ColorReset)
	for _, f := range files {
		fmt.Printf("      %s\n", f)
	}
}

func printDetailRecordSHAs(shas []string) {
	if len(shas) == 0 {
		return
	}
	fmt.Printf("    %sUnpushed Commits (%d):%s\n", constants.ColorYellow, len(shas), constants.ColorReset)
	for _, sha := range shas {
		fmt.Printf("      %s\n", sha)
	}
}

func renderPendingCommitsInteractive(records []RepoPendingCommitRecord) {
	dirtyList := collectDirtyPendingRecords(records)
	if len(dirtyList) == 0 {
		fmt.Fprintln(os.Stderr, "  All repositories are clean. Nothing to step through.")
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	for idx, rec := range dirtyList {
		renderInteractiveRecordStep(rec, idx, len(dirtyList))
		if idx < len(dirtyList)-1 && isDetailPromptTerminated(scanner) {
			break
		}
	}
}

func collectDirtyPendingRecords(records []RepoPendingCommitRecord) []RepoPendingCommitRecord {
	var dirtyList []RepoPendingCommitRecord
	for _, r := range records {
		if !r.IsClean {
			dirtyList = append(dirtyList, r)
		}
	}

	return dirtyList
}

func renderInteractiveRecordStep(rec RepoPendingCommitRecord, idx, total int) {
	fmt.Fprintf(os.Stderr, "\n  [%d/%d] Repository: %s%s%s (branch: %s)\n",
		idx+1, total, constants.ColorCyan, rec.RepoName, constants.ColorReset, rec.CurrentBranch)
	printInteractiveStepFiles(rec.PendingFiles)
	printInteractiveStepSHAs(rec.UnpushedCommitSHAs)
}

func printInteractiveStepFiles(files []string) {
	if len(files) == 0 {
		return
	}
	fmt.Fprintf(os.Stderr, "    Files (%d):\n", len(files))
	for _, f := range files {
		fmt.Fprintf(os.Stderr, "      %s\n", f)
	}
}

func printInteractiveStepSHAs(shas []string) {
	if len(shas) == 0 {
		return
	}
	fmt.Fprintf(os.Stderr, "    Unpushed Commits (%d):\n", len(shas))
	for _, sha := range shas {
		fmt.Fprintf(os.Stderr, "      %s\n", sha)
	}
}

func isDetailPromptTerminated(scanner *bufio.Scanner) bool {
	fmt.Fprintf(os.Stderr, "  Press [Enter] for next, or 'q' to quit: ")
	if !scanner.Scan() {
		return true
	}
	input := strings.TrimSpace(scanner.Text())

	return strings.EqualFold(input, "q") || strings.EqualFold(input, "quit")
}
