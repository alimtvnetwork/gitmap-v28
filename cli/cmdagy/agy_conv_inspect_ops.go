package cmdagy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AgyConvInspectRow represents a conversation row for inspection tables.
type AgyConvInspectRow struct {
	Seq       int    `json:"seq"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Messages  int    `json:"messages"`
	Status    string `json:"status"`
	Queued    int    `json:"queued"`
	Preview   string `json:"preview,omitempty"`
	IsRunning bool   `json:"isRunning"`
	LastMod   string `json:"lastModified"`
}

var (
	xmlTagRegex  = regexp.MustCompile(`(?s)<[^>]+>.*?</[^>]+>`)
	tagLineRegex = regexp.MustCompile(`(?m)^\s*<[^>]+>\s*$`)
)

// TrimBoilerplatePrompt cleans system tags and extracts up to 200 chars of user prompt.
func TrimBoilerplatePrompt(raw string) string {
	cleaned := strings.TrimSpace(raw)
	if cleaned == "" {
		return ""
	}

	cleaned = xmlTagRegex.ReplaceAllString(cleaned, " ")
	cleaned = tagLineRegex.ReplaceAllString(cleaned, " ")

	lines := strings.Split(cleaned, "\n")
	var useful []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "#") {
			continue
		}
		useful = append(useful, trimmed)
	}

	joined := strings.Join(useful, " ")
	runes := []rune(joined)
	if len(runes) > 200 {
		return string(runes[:200]) + "..."
	}

	return string(runes)
}

func detectCurrentProjectContext() (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}

	gitDir := filepath.Join(cwd, ".git")
	if fi, statErr := os.Stat(gitDir); statErr == nil && fi.IsDir() {
		return cwd, true
	}

	return cwd, false
}

// FetchConversationInspectRows loads and formats conversation inspection rows.
func FetchConversationInspectRows(n int, projectFilter string) ([]AgyConvInspectRow, bool, error) {
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return nil, false, err
	}

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return nil, false, openErr
	}
	defer conn.Close()

	isDepthMode := false
	var rows *sql.Rows
	var qErr error

	if projectFilter != "" {
		pattern := "%" + filepath.Base(projectFilter) + "%"
		query := `SELECT conversation_id, title, workspace_uris, step_count, not_fully_idle, preview, last_modified_time 
			FROM conversation_summaries 
			WHERE workspace_uris LIKE ? OR project_id LIKE ? 
			ORDER BY last_modified_time DESC LIMIT ?`
		rows, qErr = conn.Query(query, pattern, pattern, n)
		isDepthMode = true
	} else {
		query := `SELECT conversation_id, title, workspace_uris, step_count, not_fully_idle, preview, last_modified_time 
			FROM conversation_summaries 
			ORDER BY last_modified_time DESC LIMIT ?`
		rows, qErr = conn.Query(query, n)
	}

	if qErr != nil {
		return nil, false, apperror.WrapSimple(qErr, "query conversations")
	}
	defer rows.Close()

	var result []AgyConvInspectRow
	for rows.Next() {
		var cid, title, wsRaw, preview, lastMod sql.NullString
		var stepCount, notIdle sql.NullInt64

		if scanErr := rows.Scan(&cid, &title, &wsRaw, &stepCount, &notIdle, &preview, &lastMod); scanErr != nil {
			continue
		}

		wsPath := extractCleanWorkspaceFromURIs(wsRaw.String)
		isRunning := notIdle.Int64 != 0
		status := "IDLE"
		if isRunning {
			status = "RUNNING"
		}

		queuedCount := countWorkspaceQueued(wsPath)

		displayName := strings.TrimSpace(title.String)
		if displayName == "" {
			displayName = cid.String
		}

		trimmedPreview := TrimBoilerplatePrompt(preview.String)

		result = append(result, AgyConvInspectRow{
			ID:        cid.String,
			Name:      displayName,
			Path:      wsPath,
			Messages:  int(stepCount.Int64),
			Status:    status,
			Queued:    queuedCount,
			Preview:   trimmedPreview,
			IsRunning: isRunning,
			LastMod:   lastMod.String,
		})
	}

	sortInspectRows(result)
	assignInspectRowSequences(result)

	return result, isDepthMode, nil
}

func sortInspectRows(rows []AgyConvInspectRow) {
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].IsRunning != rows[j].IsRunning {
			return rows[i].IsRunning
		}
		return rows[i].LastMod > rows[j].LastMod
	})
}

func assignInspectRowSequences(rows []AgyConvInspectRow) {
	for i := range rows {
		rows[i].Seq = i + 1
	}
}

func countWorkspaceQueued(wsPath string) int {
	if wsPath == "" {
		return 0
	}
	qPath := filepath.Join(wsPath, ".ai-memory", "temp", "agy-prompt-queue.json")
	data, err := os.ReadFile(qPath)
	if err != nil {
		return 0
	}
	var q AgyPromptQueueFile
	if jsonErr := json.Unmarshal(data, &q); jsonErr != nil {
		return 0
	}

	items := filterPendingQueuedItems(q)
	return len(items)
}

// RenderInspectRowsTable renders formatted table output for conversation inspections.
func RenderInspectRowsTable(rows []AgyConvInspectRow, isDepthMode bool) {
	fmt.Println()
	fmt.Printf("  %s%-5s  %-28s  %-35s  %-8s  %-10s  %-6s%s\n",
		constants.ColorWhite,
		"SEQ", "NAME/ID", "PROJECT PATH", "MESSAGES", "STATUS", "QUEUED",
		constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("─", 100), constants.ColorReset)

	for _, r := range rows {
		statusColor := constants.ColorDim
		if r.IsRunning {
			statusColor = constants.ColorGreen + "\033[1m"
		}

		nameTrunc := r.Name
		if len(nameTrunc) > 28 {
			nameTrunc = nameTrunc[:25] + "..."
		}

		pathTrunc := r.Path
		if len(pathTrunc) > 35 {
			pathTrunc = "..." + pathTrunc[len(pathTrunc)-32:]
		}

		fmt.Printf("  %s%03d%s    %-28s  %-35s  %-8d  %s%-10s%s  %-6d\n",
			constants.ColorYellow, r.Seq, constants.ColorReset,
			nameTrunc,
			pathTrunc,
			r.Messages,
			statusColor, r.Status, constants.ColorReset,
			r.Queued,
		)

		if isDepthMode && r.Preview != "" {
			fmt.Printf("       %s↳ Prompt: %s%s\n", constants.ColorDim, r.Preview, constants.ColorReset)
		}
	}
	fmt.Println()
}

// CacheInspectRowSequences saves current inspected rows to sequence cache.
func CacheInspectRowSequences(rows []AgyConvInspectRow, category string) {
	entries := make([]CachedSequenceEntry, 0, len(rows))
	for _, r := range rows {
		entries = append(entries, CachedSequenceEntry{
			Seq:      r.Seq,
			ID:       r.ID,
			Name:     r.Name,
			Path:     r.Path,
			Messages: r.Messages,
			Status:   r.Status,
			Queued:   r.Queued,
			Preview:  r.Preview,
			Category: category,
		})
	}
	_ = SaveSequenceCache(entries)
}

// ExportInspectRowsToJSONFile writes inspection rows to a formatted JSON file.
func ExportInspectRowsToJSONFile(rows []AgyConvInspectRow, targetFile string) error {
	data, mErr := json.MarshalIndent(rows, "", "  ")
	if mErr != nil {
		return apperror.WrapSimple(mErr, "marshal inspect rows JSON")
	}
	if writeErr := os.WriteFile(targetFile, data, 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write inspect rows file")
	}
	fmt.Printf("Exported %d conversation records to %s\n", len(rows), targetFile)
	return nil
}
