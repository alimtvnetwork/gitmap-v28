package heatmap

import (
	"bufio"
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// DiscoverHeatmapTasks scans plan repositories and SQLite databases for active and completed tasks.
func DiscoverHeatmapTasks(repoRoot string, limit int) (TaskInventory, error) {
	return DiscoverHeatmapTasksForSlug(repoRoot, limit, "")
}

// DiscoverHeatmapTasksForSlug discovers tasks honoring an explicit active task slug.
func DiscoverHeatmapTasksForSlug(repoRoot string, limit int, activeSlug string) (TaskInventory, error) {
	effectiveLimit := limit
	if effectiveLimit <= 0 {
		effectiveLimit = 10
	}

	compDir := filepath.Join(repoRoot, ".ai-memory", "plans", "completed")
	pendDir := filepath.Join(repoRoot, ".ai-memory", "plans", "pending")
	compTasks, _ := scanPlanDirectory(compDir, true, activeSlug)
	pendTasks, _ := scanPlanDirectory(pendDir, false, activeSlug)
	sqliteTasks := queryTempAgentTasks(repoRoot)

	return assembleTaskInventory(compTasks, pendTasks, sqliteTasks, effectiveLimit), nil
}

func partitionTasksByStatus(tasks []TaskItem) ([]TaskItem, []TaskItem) {
	var comp, pend []TaskItem
	for _, t := range tasks {
		if t.IsCompleted {
			comp = append(comp, t)
			continue
		}

		pend = append(pend, t)
	}

	return comp, pend
}

func assembleTaskInventory(comp, pend, sqliteTasks []TaskItem, limit int) TaskInventory {
	dbComp, dbPend := partitionTasksByStatus(sqliteTasks)
	allComp := mergeTasksBySlug(comp, dbComp)
	allPend := mergeTasksBySlug(pend, dbPend)
	sortedComp := sortCompletedTasks(allComp)
	sortedPend := sortPendingTasks(allPend)

	return TaskInventory{
		Completed: applyTaskLimit(sortedComp, limit),
		Pending:   applyTaskLimit(sortedPend, limit),
	}
}

func scanPlanDirectory(dirPath string, isCompleted bool, activeSlug string) ([]TaskItem, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return []TaskItem{}, nil
	}

	return collectPlanEntries(dirPath, entries, isCompleted, activeSlug), nil
}

func collectPlanEntries(dir string, entries []os.DirEntry, isComp bool, active string) []TaskItem {
	var items []TaskItem
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		if item, hasItem := parsePlanFile(filepath.Join(dir, e.Name()), isComp, active); hasItem {
			items = append(items, item)
		}
	}

	return items
}

func buildParsedTaskItem(slug string, seq int, title string, isComp, isActive bool) TaskItem {
	return TaskItem{
		TaskSequence: seq,
		TaskSlug:     slug,
		Title:        title,
		IsActive:     isActive,
		IsCompleted:  isComp,
	}
}

func parsePlanFile(filePath string, isCompleted bool, activeSlug string) (TaskItem, bool) {
	base := filepath.Base(filePath)
	if !strings.HasSuffix(strings.ToLower(base), ".md") {
		return TaskItem{}, false
	}

	slug := strings.TrimSuffix(base, filepath.Ext(base))
	seq := extractNumericSequence(slug)
	title := extractTitleFromPlan(filePath, slug)
	isActive := !isCompleted && isSlugActive(slug, activeSlug)

	return buildParsedTaskItem(slug, seq, title, isCompleted, isActive), true
}

func isSlugActive(slug, activeSlug string) bool {
	if activeSlug != "" {
		return strings.EqualFold(slug, activeSlug)
	}

	return false
}

func extractTitleFromPlan(filePath string, defaultTitle string) string {
	f, err := os.Open(filePath)
	if err != nil {
		return defaultTitle
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if title, hasTitle := parseTitleLine(line); hasTitle {
			return title
		}
	}

	return defaultTitle
}

func parseTitleLine(line string) (string, bool) {
	if strings.HasPrefix(line, "# ") {
		clean := strings.TrimSpace(strings.TrimPrefix(line, "# "))
		clean = strings.TrimPrefix(clean, "Plan:")
		clean = strings.TrimPrefix(clean, "Master Plan:")

		return strings.TrimSpace(clean), true
	}

	return "", false
}

func extractNumericSequence(name string) int {
	idx := strings.Index(name, "-")
	if idx <= 0 {
		return 0
	}

	seq, err := strconv.Atoi(name[:idx])
	if err != nil {
		return 0
	}

	return seq
}

func sortCompletedTasks(tasks []TaskItem) []TaskItem {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].TaskSequence > tasks[j].TaskSequence
	})

	return tasks
}

func sortPendingTasks(tasks []TaskItem) []TaskItem {
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].IsActive != tasks[j].IsActive {
			return tasks[i].IsActive
		}

		return tasks[i].TaskSequence < tasks[j].TaskSequence
	})

	return tasks
}

func applyTaskLimit(tasks []TaskItem, limit int) []TaskItem {
	if limit <= 0 || len(tasks) <= limit {
		return tasks
	}

	return tasks[:limit]
}

func queryTempAgentTasks(repoRoot string) []TaskItem {
	dbPath := filepath.Join(repoRoot, ".ai-memory", "temp-agents", "ai_agents.db")
	if _, err := os.Stat(dbPath); err != nil {
		return []TaskItem{}
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return []TaskItem{}
	}

	defer conn.Close()

	return scanParentTaskRegistryRows(conn)
}

func scanParentTaskRegistryRows(conn *sql.DB) []TaskItem {
	const q = `SELECT TaskSlug, TaskName, Status, CompletedSteps FROM ParentTaskRegistry ORDER BY UpdatedAt DESC LIMIT 20`
	rows, err := conn.Query(q)
	if err != nil {
		return []TaskItem{}
	}

	defer rows.Close()

	return collectTaskRows(rows)
}

func collectTaskRows(rows *sql.Rows) []TaskItem {
	var items []TaskItem
	for rows.Next() {
		var slug, name, status string
		var steps int
		if err := rows.Scan(&slug, &name, &status, &steps); err == nil {
			items = append(items, buildTaskFromDbRow(slug, name, status, steps))
		}
	}

	return items
}

func buildTaskFromDbRow(slug, name, status string, steps int) TaskItem {
	isComp := strings.EqualFold(status, "COMPLETED") || strings.EqualFold(status, "DONE")
	isActive := strings.EqualFold(status, "ACTIVE") || strings.EqualFold(status, "IN_PROGRESS")

	return TaskItem{
		TaskSequence: extractNumericSequence(slug),
		TaskSlug:     slug,
		Title:        name,
		TotalActions: steps,
		IsActive:     isActive,
		IsCompleted:  isComp,
	}
}

func mergeTasksBySlug(primary, secondary []TaskItem) []TaskItem {
	seen := make(map[string]bool)
	var merged []TaskItem

	for _, item := range primary {
		seen[strings.ToLower(item.TaskSlug)] = true
		merged = append(merged, item)
	}

	for _, item := range secondary {
		key := strings.ToLower(item.TaskSlug)
		if !seen[key] {
			seen[key] = true
			merged = append(merged, item)
		}
	}

	return merged
}

func enrichTasksWithDbMetrics(db *store.HeatmapDB, tasks []TaskItem) []TaskItem {
	if db == nil {
		return tasks
	}

	for i := range tasks {
		if tasks[i].TotalActions == 0 {
			tasks[i].TotalActions = fetchTaskActionCount(db, tasks[i].TaskSlug)
		}
	}

	return tasks
}

func fetchTaskActionCount(db *store.HeatmapDB, taskSlug string) int {
	summary, err := db.GetHeatmapTaskSummary(taskSlug)
	if err != nil || summary == nil {
		return 0
	}

	return summary.TotalActions
}

// CollectConfusionAlerts queries active confusions or unmet expectations for status reporting.
func CollectConfusionAlerts(db *store.HeatmapDB, activeTaskSlug string) []HeatmapAiConfusion {
	if alerts := queryRepoCacheAiConfusion(activeTaskSlug); len(alerts) > 0 {
		return alerts
	}

	if alerts := queryActionLogAlerts(db, activeTaskSlug); len(alerts) > 0 {
		return alerts
	}

	return []HeatmapAiConfusion{}
}

func queryRepoCacheAiConfusion(activeTaskSlug string) []HeatmapAiConfusion {
	dbPath := resolveRepoCacheDbPath()
	if dbPath == "" {
		return []HeatmapAiConfusion{}
	}

	return queryConfusionDbFile(dbPath, activeTaskSlug)
}

func resolveRepoCacheDbPath() string {
	sdb, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return ""
	}

	defer sdb.Close()

	rec, recErr := sdb.GetSpecialRepo("rc")
	if recErr != nil || rec == nil || rec.LocalPath == "" {
		return ""
	}

	return findConfusionDbInPath(rec.LocalPath)
}

func findConfusionDbInPath(basePath string) string {
	p1 := filepath.Join(basePath, ".ai-confusion", "gitmap", "ai-confusion.db")
	if _, err := os.Stat(p1); err == nil {
		return p1
	}

	p2 := filepath.Join(basePath, ".ai-confusion", "ai-confusion.db")
	if _, err := os.Stat(p2); err == nil {
		return p2
	}

	return ""
}

func executeConfusionQuery(conn *sql.DB, activeSlug string) (*sql.Rows, error) {
	const q = `SELECT AiConfusionId, TaskSlug, FilePath, StartLine, EndLine, ReasoningNotes, IsResolved FROM AiConfusion WHERE IsResolved = 0 AND (? = '' OR TaskSlug = ?) ORDER BY AiConfusionId DESC LIMIT 10`

	return conn.Query(q, activeSlug, activeSlug)
}

func queryConfusionDbFile(dbPath, activeTaskSlug string) []HeatmapAiConfusion {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return []HeatmapAiConfusion{}
	}

	defer conn.Close()

	rows, err := executeConfusionQuery(conn, activeTaskSlug)
	if err != nil {
		return []HeatmapAiConfusion{}
	}

	defer rows.Close()

	return scanConfusionRows(rows)
}

func scanSingleConfusion(rows *sql.Rows) (HeatmapAiConfusion, bool) {
	var id, sLine, eLine, isResInt int
	var tSlug, fPath, note string
	if err := rows.Scan(&id, &tSlug, &fPath, &sLine, &eLine, &note, &isResInt); err != nil {
		return HeatmapAiConfusion{}, false
	}

	return HeatmapAiConfusion{
		ConfusionId:   id,
		TaskSlug:      tSlug,
		FilePath:      fPath,
		StartLine:     sLine,
		EndLine:       eLine,
		ReasoningNote: note,
		IsResolved:    isResInt != 0,
	}, true
}

func scanConfusionRows(rows *sql.Rows) []HeatmapAiConfusion {
	var alerts []HeatmapAiConfusion
	for rows.Next() {
		if alert, hasAlert := scanSingleConfusion(rows); hasAlert {
			alerts = append(alerts, alert)
		}
	}

	return alerts
}

func executeActionAlertQuery(conn *sql.DB, activeSlug string) (*sql.Rows, error) {
	const q = `SELECT HeatmapActionLogId, TaskSlug, FilePath, StartLine, EndLine, ReasoningNotes FROM HeatmapActionLog WHERE IsExpectationMet = 0 AND (? = '' OR TaskSlug = ?) ORDER BY HeatmapActionLogId DESC LIMIT 10`

	return conn.Query(q, activeSlug, activeSlug)
}

func queryActionLogAlerts(db *store.HeatmapDB, activeTaskSlug string) []HeatmapAiConfusion {
	if db == nil || db.Conn() == nil {
		return []HeatmapAiConfusion{}
	}

	rows, err := executeActionAlertQuery(db.Conn(), activeTaskSlug)
	if err != nil {
		return []HeatmapAiConfusion{}
	}

	defer rows.Close()

	return scanActionLogAlertRows(rows)
}

func scanSingleActionAlert(rows *sql.Rows) (HeatmapAiConfusion, bool) {
	var id int64
	var sLine, eLine int
	var tSlug, fPath, note string
	if err := rows.Scan(&id, &tSlug, &fPath, &sLine, &eLine, &note); err != nil {
		return HeatmapAiConfusion{}, false
	}

	return HeatmapAiConfusion{
		ConfusionId:   int(id),
		TaskSlug:      tSlug,
		FilePath:      fPath,
		StartLine:     sLine,
		EndLine:       eLine,
		ReasoningNote: note,
		IsResolved:    false,
	}, true
}

func scanActionLogAlertRows(rows *sql.Rows) []HeatmapAiConfusion {
	var alerts []HeatmapAiConfusion
	for rows.Next() {
		if alert, hasAlert := scanSingleActionAlert(rows); hasAlert {
			alerts = append(alerts, alert)
		}
	}

	return alerts
}

func buildSpecificCommands(cleanSlug string) SuggestedCommandsPayload {
	return SuggestedCommandsPayload{
		TreeViewCommand:        "gitmap heatmap tree --task=" + cleanSlug,
		SearchErrorCommand:     "gitmap heatmap search --state=ERROR --task=" + cleanSlug,
		AiConfusionSyncCommand: "gitmap rc ai-confusion view --task=" + cleanSlug,
	}
}

// BuildSuggestedCommands generates actionable CLI commands based on active status.
func BuildSuggestedCommands(activeTaskSlug string) SuggestedCommandsPayload {
	cleanSlug := strings.TrimSpace(activeTaskSlug)
	if cleanSlug == "" {
		return SuggestedCommandsPayload{
			TreeViewCommand:        "gitmap heatmap tree",
			SearchErrorCommand:     "gitmap heatmap search --state=ERROR --limit=10",
			AiConfusionSyncCommand: "gitmap rc ai-confusion sync",
		}
	}

	return buildSpecificCommands(cleanSlug)
}

// ResolveWorkspaceSlug detects the workspace slug name from repo root.
func ResolveWorkspaceSlug(repoRoot string) string {
	if repoRoot == "" || repoRoot == "." {
		cwd, err := os.Getwd()
		if err == nil {
			return filepath.Base(cwd)
		}

		return "gitmap"
	}

	return filepath.Base(repoRoot)
}

// ResolveActiveTaskSlug identifies the primary active task slug.
func ResolveActiveTaskSlug(flagsSlug string, inv TaskInventory) string {
	if strings.TrimSpace(flagsSlug) != "" {
		return strings.TrimSpace(flagsSlug)
	}

	for _, t := range inv.Pending {
		if t.IsActive {
			return t.TaskSlug
		}
	}

	if len(inv.Pending) > 0 {
		return inv.Pending[0].TaskSlug
	}

	return ""
}
