package cmdcommit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/output"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type commitTelemetryParams struct {
	meta          CommitMetadata
	sha           string
	shortSha      string
	repoRoot      string
	repoName      string
	branch        string
	author        string
	email         string
	msg           string
	commitTime    int64
	filesChanged  int
	linesAdded    int
	linesDeleted  int
	plansCount    int
	subtasksCount int
}

// RecordCurrentCommitSafely records the current repository's commit telemetry into Split-DB.
func RecordCurrentCommitSafely(meta CommitMetadata) {
	_ = IngestCommitHistory(".", meta)
}

// RecordChildRepoCommitSafely records a child repository's commit telemetry into Split-DB.
func RecordChildRepoCommitSafely(repoPath, msg string, meta CommitMetadata) {
	childMeta := meta
	if childMeta.TaskSlug == "" {
		childMeta.TaskSlug = extractTaskFromMsg(msg)
	}

	_ = IngestCommitHistory(repoPath, childMeta)
}

func extractTaskFromMsg(msg string) string {
	if strings.Contains(msg, "[") && strings.Contains(msg, "]") {
		start := strings.Index(msg, "[")
		end := strings.Index(msg, "]")
		if end > start+1 {
			return strings.TrimSpace(msg[start+1 : end])
		}
	}

	return ""
}

// IngestCommitHistory extracts git telemetry and records commit history to the external Split-DB.
func IngestCommitHistory(repoPath string, meta CommitMetadata) error {
	rec, touches, err := buildCommitHistoryTelemetry(repoPath, meta)
	if err != nil {
		fmt.Fprintf(output.UIErr(), "Warning: telemetry extraction failed: %v\n", err)

		return nil
	}
	recordToCommitSplitDB(rec, touches)
	updateHeatmapTelemetrySafely(rec)

	return nil
}

func recordToCommitSplitDB(rec store.CommitHistoryRecord, touches []store.CommitFileTouch) {
	db, errDb := store.OpenCommitHistorySplitDB("")
	if errDb != nil {
		fmt.Fprintf(output.UIErr(), "Warning: commit history database open failed: %v\n", errDb)

		return
	}
	defer db.Close()

	if errRec := db.RecordCommit(rec, touches); errRec != nil {
		fmt.Fprintf(output.UIErr(), "Warning: commit history record failed: %v\n", errRec)
	}
}

func buildCommitHistoryTelemetry(repoPath string, meta CommitMetadata) (store.CommitHistoryRecord, []store.CommitFileTouch, error) {
	sha, err := extractGitOutput(repoPath, "rev-parse", "HEAD")
	if err != nil {
		return store.CommitHistoryRecord{}, nil, err
	}
	shortSha := resolveShortSha(sha)
	author, email, commitTime, msg := extractCommitLog(repoPath, sha)
	touches, filesChanged, totalIns, totalDel := extractDiffMetrics(repoPath, sha)
	repoRoot, repoName, branch := resolveRepoContext(repoPath)
	plansCount, subtasksCount := countPlansAndSubtasks(repoRoot)

	p := commitTelemetryParams{
		meta: meta, sha: sha, shortSha: shortSha, repoRoot: repoRoot, repoName: repoName,
		branch: branch, author: author, email: email, msg: msg, commitTime: commitTime,
		filesChanged: filesChanged, linesAdded: totalIns, linesDeleted: totalDel,
		plansCount: plansCount, subtasksCount: subtasksCount,
	}

	return assembleCommitRecord(p), touches, nil
}

func resolveShortSha(sha string) string {
	if len(sha) >= 8 {
		return sha[:8]
	}

	return sha
}

func assembleCommitRecord(p commitTelemetryParams) store.CommitHistoryRecord {
	return store.CommitHistoryRecord{
		CommitSha: p.sha, ShortSha: p.shortSha, RepoName: p.repoName,
		RepoSlug: store.SanitizeSlug(p.repoName), RepoRoot: p.repoRoot, BranchName: p.branch,
		AuthorName: p.author, AuthorEmail: p.email, CommitMessage: p.msg,
		CommittedAtUnix: p.commitTime, AiModelName: p.meta.AIName, AiToolName: "gitmap",
		AgentId: p.meta.AgentID, HeatmapSessionId: p.meta.HeatmapID, TaskSlug: p.meta.TaskSlug,
		SpecPath: p.meta.SpecPath, PlansCount: p.plansCount, SubtasksCount: p.subtasksCount,
		IsAiCommitted: p.meta.IsAIMode, IsCurrentMachineAi: p.meta.IsAIMode,
		FilesChanged: p.filesChanged, LinesAdded: p.linesAdded, LinesDeleted: p.linesDeleted,
		HasPushed: p.meta.HasPush, CreatedAtUnix: time.Now().Unix(),
	}
}

func extractGitOutput(repoPath string, args ...string) (string, error) {
	var cmdArgs []string
	if repoPath != "" && repoPath != "." {
		cmdArgs = append(cmdArgs, "-C", repoPath)
	}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("git", cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		return "", apperror.WrapSimple(err, "cmdcommit.extractGitOutput")
	}

	return strings.TrimSpace(string(out)), nil
}

func extractCommitLog(repoPath, sha string) (string, string, int64, string) {
	out, err := extractGitOutput(repoPath, "log", "-1", "--format=%an%x00%ae%x00%ad%x00%s", "--date=unix", sha)
	if err != nil {
		return "gitmap", "gitmap@local", time.Now().Unix(), "commit"
	}
	parts := strings.Split(out, "\x00")
	if len(parts) < 4 {
		return "gitmap", "gitmap@local", time.Now().Unix(), out
	}
	unixTime, _ := strconv.ParseInt(parts[2], 10, 64)

	return parts[0], parts[1], unixTime, parts[3]
}

func extractDiffMetrics(repoPath, sha string) ([]store.CommitFileTouch, int, int, int) {
	statusOut, _ := extractGitOutput(repoPath, "diff-tree", "--no-commit-id", "--name-status", "-r", sha)
	numstatOut, _ := extractGitOutput(repoPath, "diff-tree", "--no-commit-id", "--numstat", "-r", sha)

	touches, totalIns, totalDel := parseDiffTreeOutputs(statusOut, numstatOut, sha)

	return touches, len(touches), totalIns, totalDel
}

func parseDiffTreeOutputs(statusOut, numstatOut, sha string) ([]store.CommitFileTouch, int, int) {
	statusMap := buildStatusMap(statusOut)
	var touches []store.CommitFileTouch
	totalIns, totalDel := 0, 0

	for _, line := range strings.Split(numstatOut, "\n") {
		touch, ins, del, hasTouch := parseCommitTouchLine(line, sha, statusMap)
		if hasTouch {
			totalIns += ins
			totalDel += del
			touches = append(touches, touch)
		}
	}

	return touches, totalIns, totalDel
}

func buildStatusMap(statusOut string) map[string]string {
	statusMap := make(map[string]string)
	for _, line := range strings.Split(statusOut, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			statusMap[fields[1]] = fields[0]
		}
	}

	return statusMap
}

func parseCommitTouchLine(line, sha string, statusMap map[string]string) (store.CommitFileTouch, int, int, bool) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return store.CommitFileTouch{}, 0, 0, false
	}
	ins, _ := strconv.Atoi(fields[0])
	del, _ := strconv.Atoi(fields[1])
	filePath := fields[2]

	touch := store.CommitFileTouch{
		CommitSha:     sha,
		FilePath:      filePath,
		ChangeType:    statusMap[filePath],
		LinesAdded:    ins,
		LinesDeleted:  del,
		CreatedAtUnix: time.Now().Unix(),
	}

	return touch, ins, del, true
}

func resolveRepoContext(repoPath string) (string, string, string) {
	cleanPath := repoPath
	if cleanPath == "" || cleanPath == "." {
		cleanPath, _ = filepath.Abs(".")
	}
	repoRoot := filepath.ToSlash(cleanPath)
	repoName := filepath.Base(repoRoot)
	branch, _ := extractGitOutput(repoPath, "rev-parse", "--abbrev-ref", "HEAD")

	return repoRoot, repoName, branch
}

func countPlansAndSubtasks(repoRoot string) (int, int) {
	plansDir := filepath.Join(repoRoot, ".ai-memory", "plans")
	plansCount := countFilesWithExt(plansDir, ".md")
	subtasksDir := filepath.Join(plansDir, "subtasks")
	subtasksCount := countSubtaskFiles(subtasksDir)

	return plansCount, subtasksCount
}

func countFilesWithExt(dir, ext string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ext) {
			count++
		}
	}

	return count
}

func countSubtaskFiles(subtasksDir string) int {
	count := 0
	_ = filepath.Walk(subtasksDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			count++
		}
		return nil
	})

	return count
}
