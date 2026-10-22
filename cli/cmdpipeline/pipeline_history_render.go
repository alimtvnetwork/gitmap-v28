package cmdpipeline

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func extractReleaseFromTitle(title string) string {
	if len(title) == 0 {
		return ""
	}

	lower := strings.ToLower(title)
	hasReleaseKeyword := strings.Contains(lower, "release")
	hasVTag := strings.Contains(lower, "v")
	if hasReleaseKeyword == false && hasVTag == false {
		return ""
	}

	clean := sanitizeTitleSeparators(title)
	tokens := strings.Fields(clean)
	for _, tok := range tokens {
		if isSemverToken(tok) {
			return ensureVPrefix(tok)
		}
	}

	return ""
}

func sanitizeTitleSeparators(title string) string {
	clean := strings.ReplaceAll(title, ":", " ")
	clean = strings.ReplaceAll(clean, "(", " ")
	clean = strings.ReplaceAll(clean, ")", " ")
	clean = strings.ReplaceAll(clean, "/", " ")
	clean = strings.ReplaceAll(clean, ",", " ")

	return clean
}

func isSemverToken(tok string) bool {
	raw := strings.TrimPrefix(strings.ToLower(tok), "v")
	parts := strings.Split(raw, ".")
	if len(parts) < 2 {
		return false
	}
	_, err0 := strconv.Atoi(parts[0])
	_, err1 := strconv.Atoi(parts[1])

	return err0 == nil && err1 == nil
}

func ensureVPrefix(s string) string {
	if strings.HasPrefix(strings.ToLower(s), "v") {
		return s
	}

	return "v" + s
}

func extractReleaseFromSha(sha string) string {
	isShort := len(sha) < 4
	if isShort {
		return ""
	}

	tagOut, err := exec.Command("git", "tag", "--points-at", sha).Output()
	tag := parseFirstVTagFromOutput(tagOut, err)
	if len(tag) > 0 {
		return tag
	}

	return extractReleaseFromRepoPath(sha)
}

func parseFirstVTagFromOutput(output []byte, err error) string {
	if err != nil {
		return ""
	}

	return parseFirstVTag(string(output))
}

func extractReleaseFromRepoPath(sha string) string {
	repoDir := findLocalRepoDirForSlug(resolveCurrentRepoSlug())
	if len(repoDir) == 0 {
		return ""
	}

	tagOutRepo, errRepo := exec.Command("git", "-C", repoDir, "tag", "--points-at", sha).Output()

	return parseFirstVTagFromOutput(tagOutRepo, errRepo)
}

func findLocalRepoDirForSlug(slug string) string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer db.Close()

	repo := store.NewRepoRepository(db.Conn())
	pathRes := repo.FindPathBySlug(slug)
	if pathRes.HasError() {
		logRepoLookupFailure(slug, pathRes.Err)

		return ""
	}

	return pathRes.Value
}

func logRepoLookupFailure(slug string, appErr *apperror.AppError) {
	rec := store.InternalErrorRecord{
		ErrorCode:  "E_PIPELINE_REPO_LOOKUP",
		ErrorType:  "DB_QUERY_ERROR",
		Command:    "pipeline",
		Message:    fmt.Sprintf("failed to find repo path for slug %s: %v", slug, appErr),
		SourceFile: "cli/cmdpipeline/pipeline_history_render.go",
	}
	store.LogInternalErrorRecord(rec)
}

func parseFirstVTag(output string) string {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		low := strutil.NormalizeLowerTrim(line)
		hasVPrefix := strings.HasPrefix(low, "v")
		if hasVPrefix {
			return line
		}
	}

	return ""
}

func formatGroupWorkflowsSummary(workflows []CommitWorkflowItem, maxWidth int) string {
	if len(workflows) == 0 {
		return "-"
	}
	ordered := prioritizeGroupWorkflows(workflows)
	full := summarizeGroupWorkflows(ordered)
	if visibleLen(full) <= maxWidth {
		return full
	}

	return formatWorkflowsWithRemaining(ordered, maxWidth)
}

func prioritizeGroupWorkflows(workflows []CommitWorkflowItem) []CommitWorkflowItem {
	if len(workflows) <= 1 {
		return workflows
	}
	var failing []CommitWorkflowItem
	var active []CommitWorkflowItem
	var passing []CommitWorkflowItem

	for _, wf := range workflows {
		if isWorkflowFailure(wf) {
			failing = append(failing, wf)
			continue
		}
		if isWorkflowActive(wf) {
			active = append(active, wf)
			continue
		}
		passing = append(passing, wf)
	}

	result := make([]CommitWorkflowItem, 0, len(workflows))
	result = append(result, failing...)
	result = append(result, active...)
	result = append(result, passing...)

	return result
}

func formatWorkflowsWithRemaining(workflows []CommitWorkflowItem, maxWidth int) string {
	var parts []string
	for i, wf := range workflows {
		item := fmt.Sprintf("%s [%s]", wf.Name, formatWorkflowShortStatus(wf))
		remaining := len(workflows) - (i + 1)
		candidate := buildWorkflowCandidate(parts, item, remaining)
		if visibleLen(candidate) > maxWidth {
			return finalizeTruncatedSummary(workflows, parts, remaining, maxWidth)
		}
		parts = append(parts, item)
	}

	return strings.Join(parts, ", ")
}

func buildWorkflowCandidate(parts []string, nextItem string, remaining int) string {
	suffix := ""
	if remaining > 0 {
		suffix = fmt.Sprintf(" (+%d)", remaining)
	}
	if len(parts) == 0 {
		return nextItem + suffix
	}

	return strings.Join(parts, ", ") + ", " + nextItem + suffix
}

func finalizeTruncatedSummary(workflows []CommitWorkflowItem, parts []string, remaining int, maxWidth int) string {
	remainingUnprinted := len(workflows) - len(parts)
	if len(parts) > 0 && remainingUnprinted > 0 {
		return strings.Join(parts, ", ") + fmt.Sprintf(" (+%d)", remainingUnprinted)
	}
	if len(parts) > 0 {
		return strings.Join(parts, ", ")
	}
	if len(workflows) == 0 {
		return "-"
	}
	sfx := ""
	if len(workflows) > 1 {
		sfx = fmt.Sprintf(" (+%d)", len(workflows)-1)
	}
	first := fmt.Sprintf("%s [%s]", workflows[0].Name, formatWorkflowShortStatus(workflows[0]))
	if visibleLen(first+sfx) > maxWidth && maxWidth > visibleLen(sfx)+3 {
		avail := maxWidth - visibleLen(sfx) - 3
		return first[:avail] + "..." + sfx
	}

	return first + sfx
}

func formatCommitOffsetLabel(index int) string {
	if index == 0 {
		return "latest"
	}

	return fmt.Sprintf("-%d", index)
}

func summarizeGroupWorkflows(workflows []CommitWorkflowItem) string {
	var parts []string
	for _, wf := range workflows {
		statusShort := formatWorkflowShortStatus(wf)
		parts = append(parts, fmt.Sprintf("%s [%s]", wf.Name, statusShort))
	}

	return strings.Join(parts, ", ")
}

func formatWorkflowShortStatus(wf CommitWorkflowItem) string {
	lower := strings.ToLower(strings.TrimSpace(wf.Conclusion))
	if strings.HasPrefix(lower, "cancel") {
		return "CANCEL"
	}
	if isFailingConclusion(wf.Conclusion) {
		return "FAIL"
	}
	if wf.Conclusion == "success" {
		return "PASS"
	}
	if isWorkflowStatusActive(wf.Status) {
		return "RUN"
	}
	if len(wf.Conclusion) > 0 {
		return strings.ToUpper(wf.Conclusion)
	}
	if len(wf.Status) > 0 {
		return strings.ToUpper(wf.Status)
	}

	return "UNKNOWN"
}

func renderCachedErrorsList(errors []pipelinedb.PipelineErrorRecord) {
	if len(errors) == 0 {
		fmt.Println("    (No error diagnostic logs cached for this run)")

		return
	}

	for _, e := range errors {
		fmt.Printf("    %s %s\n", termout.TokyoErrorPill("STEP"), e.StepName)
		fmt.Printf("      %s%s%s\n", constants.ColorRed, e.ErrorText, constants.ColorReset)
	}
}

func buildFailureCardHeader(r pipelinedb.PipelineRunRecord, idx, total int) string {
	tag := termout.TokyoTag("FAIL", fmt.Sprintf("[%d/%d] #%d", idx, total, r.RunId))
	title := termout.SafeTokyoStyle(termout.TokyoPrimary).Bold(true).Render(r.WorkflowName)

	return fmt.Sprintf("%s %s", tag, title)
}

func buildFailureCardMeta(r pipelinedb.PipelineRunRecord) string {
	bTag := termout.TokyoTag("Branch", r.Branch)
	sTag := termout.TokyoTag("SHA", r.Sha)
	dTag := termout.TokyoTag("Duration", formatDurationSeconds(r.DurationSeconds))
	tTag := termout.TokyoTag("Date", formatRunTimestamp(r.CreatedAt))

	return fmt.Sprintf("%s  %s\n%s  %s", bTag, sTag, dTag, tTag)
}

func renderSingleCachedFailureCard(db *pipelinedb.PipelineSplitDb, r pipelinedb.PipelineRunRecord, idx, total int, isDetailed bool) {
	header := buildFailureCardHeader(r, idx, total)
	meta := buildFailureCardMeta(r)
	boxContent := fmt.Sprintf("%s\n\n%s", header, meta)
	fmt.Printf("\n%s\n", termout.TokyoBoxStyle(boxContent))
	renderRunErrorsByDetailMode(db, r.RunId, isDetailed)
}

func renderRunErrorsByDetailMode(db *pipelinedb.PipelineSplitDb, runId uint64, isDetailed bool) {
	if isDetailed {
		renderDetailedOrLegacyErrors(db, runId)

		return
	}

	renderCompactOrLegacyErrors(db, runId)
}

func renderDetailedOrLegacyErrors(db *pipelinedb.PipelineSplitDb, runId uint64) {
	detailRes := db.QueryDetailedErrorLogsByRunId(runId)
	if detailRes.HasRecord() {
		renderCachedErrorsList(detailRes.Data)

		return
	}

	legacyRes := db.QueryErrorLogsByRunId(runId)
	renderCachedErrorsList(legacyRes.Data)
}

func renderCompactOrLegacyErrors(db *pipelinedb.PipelineSplitDb, runId uint64) {
	compactRes := db.QueryCompactErrorLogsByRunId(runId)
	if compactRes.HasRecord() {
		renderCachedCompactErrorsList(compactRes.Data)

		return
	}

	legacyRes := db.QueryErrorLogsByRunId(runId)
	renderCachedErrorsList(legacyRes.Data)
}

func renderCachedCompactErrorsList(errors []pipelinedb.PipelineCompactErrorRecord) {
	if len(errors) == 0 {
		fmt.Println("    (No error diagnostic logs cached for this run)")

		return
	}

	for _, e := range errors {
		fmt.Printf("    %s %s\n", termout.TokyoErrorPill("STEP"), e.StepName)
		fmt.Printf("      %s%s%s\n", constants.ColorRed, e.ErrorText, constants.ColorReset)
	}
}

func renderCachedFailuresTerminal(db *pipelinedb.PipelineSplitDb, repo string, runs []pipelinedb.PipelineRunRecord, isDetailed bool) {
	relDb := filepath.ToSlash(FormatRelativeDbPath(db.Path))
	if len(runs) == 0 {
		printEmptyCachedFailures(repo, db.Path)

		return
	}

	printCachedFailuresList(db, relDb, runs, isDetailed)
}

func printEmptyCachedFailures(repo, dbPath string) {
	fmt.Printf("\n  No cached pipeline failures found in SQLite for %s.\n", repo)
	fmt.Printf("  Pipeline DB: %s\n", filepath.ToSlash(FormatRelativeDbPath(dbPath)))
	fmt.Printf("  DB Size:     %s\n", ResolveDbFileSize(dbPath))
	fmt.Printf("  Cleanup:     gitmap pipeline clear -y\n\n")
}

func printCachedFailuresList(db *pipelinedb.PipelineSplitDb, relDb string, runs []pipelinedb.PipelineRunRecord, isDetailed bool) {
	fmt.Printf("\n  %s● Last %d Cached Pipeline Failure(s) from SQLite (%s):%s\n",
		constants.ColorRed, len(runs), relDb, constants.ColorReset)
	for i, r := range runs {
		renderSingleCachedFailureCard(db, r, i+1, len(runs), isDetailed)
	}
}

func fetchLastCachedFailures(repo string, count int) (*pipelinedb.PipelineSplitDb, []pipelinedb.PipelineRunRecord, error) {
	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return nil, nil, apperror.WrapSimple(err, "fetch_last_cached_failures")
	}

	runs, qErr := queryCachedFailedRuns(db, count)
	if qErr != nil {
		_ = db.Close()

		return nil, nil, qErr
	}

	return db, runs, nil
}

func queryCachedFailedRuns(db *pipelinedb.PipelineSplitDb, count int) ([]pipelinedb.PipelineRunRecord, error) {
	limit := resolveMaxSyncLimit(count)
	runRes := db.QueryLastFailedRuns(limit)
	if runRes.IsFailure() {
		return nil, runRes.AppError()
	}

	return runRes.Data, nil
}

// RenderLastCachedFailures displays cached failure logs from SQLite up to count.
func RenderLastCachedFailures(repo string, count int, isJSON bool, isDetailed ...bool) error {
	detailed := isDetailedRequested(isDetailed)
	db, runs, err := fetchLastCachedFailures(repo, count)
	if err != nil {
		return err
	}
	defer db.Close()

	if isJSON {
		return printCachedRunsJSON(runs)
	}

	renderCachedFailuresTerminal(db, repo, runs, detailed)

	return nil
}

func isDetailedRequested(isDetailed []bool) bool {
	return len(isDetailed) > 0 && isDetailed[0]
}

func printCachedRunsJSON(runs []pipelinedb.PipelineRunRecord) error {
	err := printJSON(runs)
	if err != nil {
		return apperror.WrapSimple(err, "render_last_cached_failures")
	}

	return nil
}

// HandlePipelineHistoryErrors inspects args for negative offset or --last-failures flag and executes them.
func HandlePipelineHistoryErrors(args []string) (bool, error) {
	flags := ParsePipelineErrorFlags(args)
	repo := flags.RepoTarget
	if len(repo) == 0 {
		repo = resolveCurrentRepoSlug()
	}
	isJSON := flags.IsJSON

	if isHandled, err := tryHandleOffsetInspection(repo, args, isJSON); isHandled {
		return true, err
	}

	return tryHandleLastFailuresInspection(repo, args, isJSON)
}

func tryHandleOffsetInspection(repo string, args []string, isJSON bool) (bool, error) {
	offset, hasOffset := extractNegativeOffset(args)
	if hasOffset {
		runs := resolveCachedRunsOrFetch(repo)

		return true, InspectPositionalRun(repo, runs, offset, isJSON)
	}

	return false, nil
}

func tryHandleLastFailuresInspection(repo string, args []string, isJSON bool) (bool, error) {
	count, hasLastFailures := extractLastFailuresFlag(args)
	if hasLastFailures {
		isDetailed := hasDetailedArg(args)

		return true, RenderLastCachedFailures(repo, count, isJSON, isDetailed)
	}

	return false, nil
}
