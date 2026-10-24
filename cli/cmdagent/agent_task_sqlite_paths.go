// Package cmdagent — agent_task_sqlite_paths.go handles status queries, path resolution, and secrets URLs.
package cmdagent

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var runPrefixRe = lazyregex.New(`^(\d+)-`)

func buildTaskInitOutput(repoRoot, dbPath, relRunDir, slug string, parentId int64, budget int) map[string]any {
	absDbPath, _ := filepath.Abs(dbPath)
	relDbPath := toRelativePath(repoRoot, dbPath)
	return map[string]any{
		"action":           "INITIALIZED",
		"parentTaskId":     parentId,
		"slug":             slug,
		"runDirectory":     relRunDir,
		"database":         relDbPath,
		"databasePath":     relDbPath,
		"relativeDbPath":   relDbPath,
		"absoluteDbPath":   filepath.ToSlash(absDbPath),
		"repoSecretsUrl":   resolveAgentRepoSecretsURL(),
		"totalStepsBudget": budget,
		"status":           "ACTIVE",
	}
}

func resolveNextRunPrefix(tempDir string) int {
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return 1
	}

	maxPrefix := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		matches := runPrefixRe.FindStringSubmatch(entry.Name())
		if len(matches) <= 1 {
			continue
		}
		val, convErr := strconv.Atoi(matches[1])
		if convErr == nil && val > maxPrefix {
			maxPrefix = val
		}
	}

	return maxPrefix + 1
}

func runTaskSQLiteStatus(args []string) error {
	fs := flag.NewFlagSet("task status", flag.ContinueOnError)
	dbFlag := fs.String("db", "", "Path to SQLite database")
	_ = fs.Bool("json", true, "Output results in JSON format")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	dbPath := resolveDbPath(*dbFlag)
	db, openErr := OpenSqliteDB(dbPath)
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	var parentSlug, parentName, parentStatus string
	row := db.QueryRow("SELECT TaskSlug, TaskName, Status FROM ParentTask ORDER BY ParentTaskId DESC LIMIT 1;")
	scanErr := row.Scan(&parentSlug, &parentName, &parentStatus)
	if scanErr != nil && scanErr != sql.ErrNoRows {
		return appfault.WrapExecution(scanErr, "failed to scan parent task record")
	}

	rows, queryErr := db.Query(`
		SELECT SubtaskId, TaskCode, Title, COALESCE(AssignedAgentRole, ''), Status, COALESCE(Evidence, '')
		FROM Subtask
		ORDER BY SubtaskId ASC;
	`)
	if queryErr != nil {
		return appfault.WrapExecution(queryErr, "failed to query subtasks")
	}
	defer rows.Close()

	type subtaskOut struct {
		SubtaskId         int64  `json:"SubtaskId"`
		TaskCode          string `json:"TaskCode"`
		Title             string `json:"Title"`
		AssignedAgentRole string `json:"AssignedAgentRole"`
		Status            string `json:"Status"`
		Evidence          string `json:"Evidence"`
	}

	var subtasks []subtaskOut
	counts := make(map[string]int)

	for rows.Next() {
		var s subtaskOut
		if err := rows.Scan(&s.SubtaskId, &s.TaskCode, &s.Title, &s.AssignedAgentRole, &s.Status, &s.Evidence); err == nil {
			subtasks = append(subtasks, s)
			counts[s.Status]++
		}
	}

	total := len(subtasks)
	done := counts["DONE"]
	var percent float64
	if total > 0 {
		percent = math.Round((float64(done)/float64(total)*100)*10) / 10
	}

	out := buildTaskStatusOutput(parentSlug, parentName, parentStatus, dbPath, total, done, counts, percent, subtasks)
	return printJson(out)
}

func resolveTaskDbPath(repoRoot, raw string) string {
	cleaned := filepath.Clean(strings.TrimSpace(raw))
	if filepath.IsAbs(cleaned) {
		return cleaned
	}
	return filepath.Join(repoRoot, cleaned)
}

func resolveDbPath(flagVal string) string {
	clean := strings.TrimSpace(flagVal)
	repoRoot := FindRepoRoot()

	if clean != "" {
		return resolveTaskDbPath(repoRoot, clean)
	}

	tempDir := ResolveAgentTempDir("")
	latestRun, err := findLatestRunDir(tempDir)
	if err == nil && latestRun != "" {
		return filepath.Join(latestRun, Tier2TaskDBName)
	}

	return filepath.Join(tempDir, "01-agent", Tier2TaskDBName)
}

func buildTaskStatusOutput(parentSlug, parentName, parentStatus, dbPath string, total, done int, counts map[string]int, percent float64, subtasks any) map[string]any {
	repoRoot := FindRepoRoot()
	absDbPath, _ := filepath.Abs(dbPath)
	relDbPath := toRelativePath(repoRoot, dbPath)
	return map[string]any{
		"taskSlug":        parentSlug,
		"taskName":        parentName,
		"status":          parentStatus,
		"database":        relDbPath,
		"relativeDbPath":  relDbPath,
		"absoluteDbPath":  filepath.ToSlash(absDbPath),
		"repoSecretsUrl":  resolveAgentRepoSecretsURL(),
		"totalSubtasks":   total,
		"completed":       done,
		"pending":         counts["PENDING"],
		"inProgress":      counts["IN_PROGRESS"],
		"failed":          counts["FAILED"],
		"percentComplete": percent,
		"subtasks":        subtasks,
	}
}

func toRelativePath(root, target string) string {
	absTarget := target
	if !filepath.IsAbs(target) {
		absTarget = filepath.Join(root, target)
	}
	absRoot, _ := filepath.Abs(root)
	absTarget, _ = filepath.Abs(absTarget)
	rel, err := filepath.Rel(absRoot, absTarget)
	if err != nil {
		return filepath.ToSlash(target)
	}
	return filepath.ToSlash(rel)
}

func resolveAgentRepoSecretsURL() string {
	if url := queryAgentSpecialRepoSecretsURL(); len(url) > 0 {
		return url
	}
	return queryAgentGitSecretsRemoteURL()
}

func queryAgentSpecialRepoSecretsURL() string {
	sdb, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return ""
	}
	defer sdb.Close()

	rec, recErr := sdb.GetSpecialRepo("rs")
	if recErr == nil && rec != nil && len(strings.TrimSpace(rec.RemoteURL)) > 0 {
		return strings.TrimSpace(rec.RemoteURL)
	}
	return ""
}

func queryAgentGitSecretsRemoteURL() string {
	secretsRoot := secrets.ResolveRepoSecretsRoot()
	if len(secretsRoot) == 0 {
		return ""
	}
	cmd := exec.Command("git", "-C", secretsRoot, "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		return strings.TrimSpace(string(out))
	}
	return ""
}

func printJson(data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}

func readTasksRawJson(jsonStr, jsonFile string) (string, error) {
	raw := strings.TrimSpace(jsonStr)
	if raw != "" {
		return raw, nil
	}
	trimmedFile := strings.TrimSpace(jsonFile)
	if trimmedFile == "" {
		return "", nil
	}
	fileBytes, err := os.ReadFile(trimmedFile)
	if err != nil {
		return "", appfault.WrapExecution(err, "failed to read tasks-file")
	}
	return string(fileBytes), nil
}

func resolveParentTaskId(db *sql.DB, explicitId int64) (int64, error) {
	if explicitId > 0 {
		return explicitId, nil
	}
	return queryLatestParentTaskId(db)
}

func queryLatestParentTaskId(db *sql.DB) (int64, error) {
	var parentId int64
	row := db.QueryRow("SELECT ParentTaskId FROM ParentTask ORDER BY ParentTaskId DESC LIMIT 1;")
	scanErr := row.Scan(&parentId)
	if scanErr != nil {
		return 0, appfault.WrapExecution(scanErr, "ParentTask record missing in database")
	}
	return parentId, nil
}
