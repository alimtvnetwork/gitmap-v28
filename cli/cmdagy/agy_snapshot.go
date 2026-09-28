package cmdagy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ExportAGYRestoreSnapshot scans Antigravity workspaces and conversations and writes a snapshot JSON.
func ExportAGYRestoreSnapshot(outPath string) (AGYRestoreSnapshot, error) {
	targetPath, pathErr := resolveSnapshotOutputPath(outPath)
	hasPathErr := pathErr != nil
	if hasPathErr {
		return AGYRestoreSnapshot{}, pathErr
	}
	snap := buildAGYRestoreSnapshot(collectSnapshotProjects(), collectSnapshotConversations())
	writeErr := writeSnapshotFile(targetPath, snap)
	hasWriteErr := writeErr != nil
	if hasWriteErr {
		return snap, writeErr
	}
	return snap, nil
}

func buildAGYRestoreSnapshot(projects []AGYSnapshotProject, convs []AGYSnapshotConv) AGYRestoreSnapshot {
	return AGYRestoreSnapshot{
		CreatedAt:     time.Now().UTC(),
		Timestamp:     time.Now().Unix(),
		TotalProjects: len(projects),
		TotalConvs:    len(convs),
		Projects:      projects,
		Conversations: convs,
	}
}

func resolveSnapshotOutputPath(outPath string) (string, error) {
	hasCustom := outPath != ""
	if hasCustom {
		return filepath.Abs(outPath)
	}
	home, err := os.UserHomeDir()
	hasErr := err != nil
	if hasErr {
		return "", err
	}
	gitmapDir := filepath.Join(home, ".gitmap")
	_ = os.MkdirAll(gitmapDir, 0755)
	fileName := fmt.Sprintf("agy-snapshot-%s.json", time.Now().Format("20060102-150405"))
	return filepath.Join(gitmapDir, fileName), nil
}

func collectSnapshotProjects() []AGYSnapshotProject {
	dirPath, err := getProjectsDirPath()
	hasErr := err != nil
	if hasErr {
		return nil
	}
	projects, loadErr := loadAllAgyProjects(dirPath)
	hasLoadErr := loadErr != nil
	if hasLoadErr {
		return nil
	}
	var out []AGYSnapshotProject
	for _, p := range projects {
		out = append(out, convertProjectToSnapshot(p))
	}
	return out
}

func convertProjectToSnapshot(p AgyProject) AGYSnapshotProject {
	return AGYSnapshotProject{
		ID:        p.ID,
		Name:      p.Name,
		Workspace: p.GetPath(),
		Branch:    p.GetBranch(),
	}
}

func collectSnapshotConversations() []AGYSnapshotConv {
	seen := make(map[string]bool)
	var all []AGYSnapshotConv
	fromDB := collectConvsFromDB(seen)
	all = append(all, fromDB...)
	fromBrain := collectConvsFromBrain(seen)
	all = append(all, fromBrain...)
	return all
}

func collectConvsFromDB(seen map[string]bool) []AGYSnapshotConv {
	dbPath, err := getConversationSummariesDBPath()
	hasErr := err != nil
	if hasErr {
		return nil
	}
	conn, openErr := store.OpenSQLiteDB(dbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil
	}
	defer conn.Close()

	return readSummaryRowsForSnapshot(conn, seen)
}

func readSummaryRowsForSnapshot(conn *sql.DB, seen map[string]bool) []AGYSnapshotConv {
	query := "SELECT conversation_id, title, workspace_uris, project_id, step_count FROM conversation_summaries ORDER BY last_modified_time DESC"
	rows, err := conn.Query(query)
	hasErr := err != nil
	if hasErr {
		return nil
	}
	defer rows.Close()

	return scanSnapshotRows(rows, seen)
}

func scanSnapshotRows(rows *sql.Rows, seen map[string]bool) []AGYSnapshotConv {
	var list []AGYSnapshotConv
	for rows.Next() {
		item, isOk := scanSingleSnapshotRow(rows, seen)
		if isOk {
			list = append(list, item)
		}
	}
	return list
}

func scanSingleSnapshotRow(rows *sql.Rows, seen map[string]bool) (AGYSnapshotConv, bool) {
	var id, title, uri, projID string
	var steps int
	err := rows.Scan(&id, &title, &uri, &projID, &steps)
	hasErr := err != nil
	if hasErr {
		return AGYSnapshotConv{}, false
	}
	seen[id] = true
	return AGYSnapshotConv{
		ID:            id,
		Title:         title,
		ProjectID:     projID,
		WorkspacePath: cleanURIStringToPath(uri),
		StepCount:     steps,
	}, true
}

func collectConvsFromBrain(seen map[string]bool) []AGYSnapshotConv {
	brainDir, err := GetBrainLogsDirPath()
	hasErr := err != nil
	if hasErr {
		return nil
	}
	entries, readErr := os.ReadDir(brainDir)
	hasReadErr := readErr != nil
	if hasReadErr {
		return nil
	}
	return scanBrainEntriesForSnapshot(brainDir, entries, seen)
}

func scanBrainEntriesForSnapshot(brainDir string, entries []os.DirEntry, seen map[string]bool) []AGYSnapshotConv {
	var list []AGYSnapshotConv
	for _, e := range entries {
		hasSeen := seen[e.Name()]
		if hasSeen {
			continue
		}
		hasDir := e.IsDir()
		if !hasDir {
			continue
		}
		item, isOk := buildBrainConvSnapshot(brainDir, e.Name())
		if isOk {
			seen[e.Name()] = true
			list = append(list, item)
		}
	}
	return list
}

func buildBrainConvSnapshot(brainDir, convID string) (AGYSnapshotConv, bool) {
	ws := resolveConvWorkspace(convID)
	prompts := readConvPrompts(brainDir, convID)
	title := convID
	hasPrompts := len(prompts) > 0
	if hasPrompts {
		title = truncateSnapshotTitle(prompts[0].Content)
	}
	return AGYSnapshotConv{
		ID:            convID,
		Title:         title,
		WorkspacePath: ws,
		StepCount:     len(prompts),
	}, true
}

func truncateSnapshotTitle(content string) string {
	clean := strings.ReplaceAll(content, "\n", " ")
	clean = strings.TrimSpace(clean)
	hasShort := len(clean) <= 60
	if hasShort {
		return clean
	}
	return clean[:60] + "..."
}

func writeSnapshotFile(outPath string, snap AGYRestoreSnapshot) error {
	dir := filepath.Dir(outPath)
	_ = os.MkdirAll(dir, 0755)
	data, err := json.MarshalIndent(snap, "", "  ")
	hasErr := err != nil
	if hasErr {
		return err
	}
	return os.WriteFile(outPath, data, 0644)
}
