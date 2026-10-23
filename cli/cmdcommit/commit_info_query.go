package cmdcommit

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/heatmap"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// BuildCommitInfoReport aggregates Git, Split-DB, Heatmap, and Agent metadata into a CommitInfoReport.
func BuildCommitInfoReport(repoPath, fullSha string) (*CommitInfoReport, error) {
	rep := &CommitInfoReport{
		FullSha:  fullSha,
		RepoPath: repoPath,
		RepoSlug: filepath.Base(repoPath),
	}

	if err := queryGitCommitMetadata(repoPath, fullSha, rep); err != nil {
		return nil, err
	}

	populateRemainingCommitData(repoPath, fullSha, rep)

	return rep, nil
}

func populateRemainingCommitData(repoPath, fullSha string, rep *CommitInfoReport) {
	rep.BranchName = queryGitBranchName(repoPath, fullSha)
	queryGitFileDeltas(repoPath, fullSha, rep)
	querySplitDbRecord(repoPath, fullSha, rep)
	queryHeatmapTelemetry(fullSha, rep)
	queryAgentTaskDatabases(repoPath, rep)
	resolveSpecPath(repoPath, rep)
	rep.Provenance = DetectCommitOrigin(rep)
}

func queryGitCommitMetadata(repoPath, fullSha string, rep *CommitInfoReport) error {
	format := "--format=%H%x1f%h%x1f%an%x1f%ae%x1f%ct%x1f%s%x1f%P"
	res := gitutil.RunGitInDir(repoPath, "log", "-1", format, fullSha)
	if res.IsFailure() {
		return apperror.WrapSimple(res.AppError(), "queryGitCommitMetadata failed")
	}

	parts := strings.Split(strings.TrimSpace(string(res.Data)), "\x1f")
	if len(parts) < 6 {
		return apperror.WrapSimple(errors.New("invalid git log format"), "parse git log")
	}

	populateCommitMetadata(rep, parts)

	return nil
}

func populateCommitMetadata(rep *CommitInfoReport, parts []string) {
	rep.FullSha = parts[0]
	rep.ShortSha = parts[1]
	rep.AuthorName = parts[2]
	rep.AuthorEmail = parts[3]
	sec, _ := strconv.ParseInt(parts[4], 10, 64)
	rep.CommittedAt = time.Unix(sec, 0)
	rep.FormattedDate = rep.CommittedAt.Format("2006-01-02 15:04:05 -0700")
	rep.RelativeDate = formatRelativeDate(rep.CommittedAt)
	rep.CommitMessage = parts[5]
	if len(parts) > 6 && len(strings.TrimSpace(parts[6])) > 0 {
		rep.ParentShas = strings.Fields(parts[6])
	}
}

func formatRelativeDate(t time.Time) string {
	diff := time.Since(t)
	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	}
	days := int(diff.Hours() / 24)

	return fmt.Sprintf("%d days ago", days)
}

func queryGitBranchName(repoPath, fullSha string) string {
	res := gitutil.RunGitInDir(repoPath, "branch", "--contains", fullSha)
	if res.IsFailure() {
		return "unknown"
	}

	lines := strings.Split(string(res.Data), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "*"))
		if len(trimmed) > 0 {
			return trimmed
		}
	}

	return "unknown"
}

func queryGitFileDeltas(repoPath, fullSha string, rep *CommitInfoReport) {
	res := gitutil.RunGitInDir(repoPath, "diff-tree", "--no-commit-id", "--numstat", "-r", fullSha)
	if res.IsFailure() {
		return
	}

	lines := strings.Split(strings.TrimSpace(string(res.Data)), "\n")
	for _, line := range lines {
		delta, hasDelta := parseNumstatLine(line)
		if hasDelta {
			rep.FileDeltas = append(rep.FileDeltas, delta)
			rep.TotalInsertions += delta.Insertions
			rep.TotalDeletions += delta.Deletions
		}
	}
	rep.TotalFiles = len(rep.FileDeltas)
}

func parseNumstatLine(line string) (CommitFileDelta, bool) {
	fields := strings.Split(strings.TrimSpace(line), "\t")
	if len(fields) < 3 {
		return CommitFileDelta{}, false
	}

	isBinary := fields[0] == "-" || fields[1] == "-"
	ins, _ := strconv.Atoi(fields[0])
	del, _ := strconv.Atoi(fields[1])

	delta := CommitFileDelta{
		FilePath:   fields[2],
		Insertions: ins,
		Deletions:  del,
		IsBinary:   isBinary,
	}

	return delta, true
}

func querySplitDbRecord(repoPath, fullSha string, rep *CommitInfoReport) {
	dbPath := store.ResolveCommitHistoryDbPath(repoPath)
	if !isFileExisting(dbPath) {
		return
	}

	conn, err := openReadOnlyDb(dbPath)
	if err != nil {
		return
	}

	defer conn.Close()
	scanSplitDbRow(conn, fullSha, dbPath, rep)
}

func scanSplitDbRow(conn *sql.DB, fullSha, dbPath string, rep *CommitInfoReport) {
	q := `SELECT CommitHistoryId, RepoSlug, AiModelName, AiToolName, AgentId, 
		HeatmapSessionId, HeatmapActionId, TaskSlug, SpecPath, PlansCount, 
		SubtasksCount, IsAiCommitted, IsCurrentMachineAi, Notes 
		FROM CommitHistoryRecord WHERE CommitSha = ? LIMIT 1`
	var id, actId int64
	var slug, aiModel, tool, agentId, sessId, task, spec, notes string
	var plans, subs, isAi, isCurrent int
	err := conn.QueryRow(q, fullSha).Scan(
		&id, &slug, &aiModel, &tool, &agentId, &sessId, &actId,
		&task, &spec, &plans, &subs, &isAi, &isCurrent, &notes,
	)
	if err != nil {
		return
	}

	rep.CrossReferences.HasSplitDbRecord = true
	rep.CrossReferences.SplitDbPath = dbPath
	rep.CrossReferences.SplitDbRecordId = id
	rep.Provenance.AiModel = aiModel
	rep.Provenance.ToolName = tool
	rep.Provenance.AgentId = agentId
	rep.Provenance.SessionUuid = sessId
	rep.Provenance.IsAiCommitted = (isAi == 1)
	rep.Provenance.TaskSlug = task
	rep.SpecPath = spec
	rep.CognitiveIntention = notes
	rep.PlanCount = plans
	rep.SubtaskCount = subs
}

func queryHeatmapTelemetry(fullSha string, rep *CommitInfoReport) {
	dbs, err := heatmap.DiscoverHeatmapDatabases()
	if err != nil || len(dbs) == 0 {
		return
	}

	for _, db := range dbs {
		if hasFound := scanHeatmapCommitRecord(db.Path, fullSha, rep); hasFound {
			return
		}
	}
}

func scanHeatmapCommitRecord(dbPath, fullSha string, rep *CommitInfoReport) bool {
	conn, err := openReadOnlyDb(dbPath)
	if err != nil {
		return false
	}

	defer conn.Close()

	var id int64
	var taskSlug string
	q := `SELECT HeatmapCommitLogId, TaskSlug FROM HeatmapCommitLog WHERE CommitSha = ? LIMIT 1`
	if err := conn.QueryRow(q, fullSha).Scan(&id, &taskSlug); err != nil {
		return false
	}

	rep.CrossReferences.HasHeatmapRecord = true
	rep.CrossReferences.HeatmapDbPath = dbPath
	if len(rep.Provenance.TaskSlug) == 0 && len(taskSlug) > 0 {
		rep.Provenance.TaskSlug = taskSlug
	}
	rep.CrossReferences.HeatmapActionCount = countHeatmapActions(conn, taskSlug)

	return true
}

func countHeatmapActions(conn *sql.DB, taskSlug string) int {
	if len(taskSlug) == 0 {
		return 0
	}

	var count int
	q := `SELECT COUNT(*) FROM HeatmapActionLog WHERE TaskSlug = ?`
	if err := conn.QueryRow(q, taskSlug).Scan(&count); err != nil {
		return 0
	}

	return count
}

func queryAgentTaskDatabases(repoPath string, rep *CommitInfoReport) {
	taskSlug := rep.Provenance.TaskSlug
	if len(taskSlug) == 0 {
		return
	}

	taskDbPath := store.ResolveTaskDbPath(repoPath, taskSlug)
	if !isFileExisting(taskDbPath) {
		return
	}

	conn, err := openReadOnlyDb(taskDbPath)
	if err != nil {
		return
	}

	defer conn.Close()
	rep.CrossReferences.HasAgentTask = true
	rep.CrossReferences.AgentTaskDbPath = taskDbPath
	querySubtaskInfo(conn, rep)
}

func querySubtaskInfo(conn *sql.DB, rep *CommitInfoReport) {
	var code, status string
	q := `SELECT SubtaskCode, Status FROM Subtask ORDER BY SubtaskId DESC LIMIT 1`
	if err := conn.QueryRow(q).Scan(&code, &status); err == nil {
		rep.CrossReferences.SubtaskCode = code
		rep.CrossReferences.SubtaskStatus = status
	}
}

func resolveSpecPath(repoPath string, rep *CommitInfoReport) {
	if len(rep.SpecPath) > 0 {
		return
	}

	taskSlug := rep.Provenance.TaskSlug
	if len(taskSlug) == 0 {
		return
	}

	relSpec := filepath.ToSlash(filepath.Join("02-spec", "21-app", taskSlug, "01-architecture-spec.md"))
	fullSpec := filepath.Join(repoPath, relSpec)
	if isFileExisting(fullSpec) {
		rep.SpecPath = relSpec
	}
}
