package cmdagy

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func fetchRunningProjectsSet() map[string]bool {
	result := make(map[string]bool)
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return result
	}

	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return result
	}
	conn := openRes.Value
	defer conn.Close()

	query := `SELECT project_id, workspace_uris, title
		FROM conversation_summaries
		WHERE (killed IS NULL OR killed = 0) AND not_fully_idle != 0`

	rows, qErr := conn.Query(query)
	if qErr != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var pid, wsRaw, title sql.NullString
		if scanErr := rows.Scan(&pid, &wsRaw, &title); scanErr != nil {
			continue
		}
		if pid.Valid && pid.String != "" {
			result[pid.String] = true
		}
		recordRunningWorkspaceUri(wsRaw, result)
		if title.Valid && title.String != "" {
			result[strings.ToLower(title.String)] = true
		}
	}

	return result
}

func recordRunningWorkspaceUri(wsRaw sql.NullString, result map[string]bool) {
	if !wsRaw.Valid || wsRaw.String == "" {
		return
	}
	wsClean := extractCleanWorkspaceFromURIs(wsRaw.String)
	if wsClean != "" {
		result[wsClean] = true
	}
}

func recordWorkspaceActivity(wsRaw sql.NullString, timeVal string, result map[string]string) {
	if !wsRaw.Valid || wsRaw.String == "" {
		return
	}
	wsClean := extractCleanWorkspaceFromURIs(wsRaw.String)
	if wsClean != "" {
		result[wsClean] = timeVal
	}
}

func extractCleanWorkspaceFromURIs(wsRaw string) string {
	lowerURIs := strings.ToLower(wsRaw)
	matchesSub := fileURIRegex.FindString(lowerURIs)
	if matchesSub == "" {
		return ""
	}
	return cleanURIStringToPath(matchesSub)
}

func fetchPinnedProjectRankMap() map[string]int {
	rankMap := make(map[string]int)
	convOrder := fetchPinnedConversationsOrder()
	if len(convOrder) == 0 {
		return rankMap
	}

	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return rankMap
	}

	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return rankMap
	}
	conn := openRes.Value
	defer conn.Close()

	for rank, convID := range convOrder {
		populateRankFromConv(conn, convID, rank, rankMap)
	}

	return rankMap
}

func populateRankFromConv(conn *sql.DB, convID string, rank int, rankMap map[string]int) {
	query := "SELECT project_id, workspace_uris, title FROM conversation_summaries WHERE conversation_id = ? LIMIT 1"
	var pid, wsRaw, title sql.NullString
	if err := conn.QueryRow(query, convID).Scan(&pid, &wsRaw, &title); err != nil {
		return
	}

	assignProjectRank(pid, rank, rankMap)
	assignWorkspaceRank(wsRaw, rank, rankMap)
	assignTitleRank(title, rank, rankMap)
}

func assignProjectRank(pid sql.NullString, rank int, rankMap map[string]int) {
	if !pid.Valid || pid.String == "" {
		return
	}
	if _, exists := rankMap[pid.String]; !exists {
		rankMap[pid.String] = rank
	}
}

func assignWorkspaceRank(wsRaw sql.NullString, rank int, rankMap map[string]int) {
	if !wsRaw.Valid || wsRaw.String == "" {
		return
	}
	wsClean := extractCleanWorkspaceFromURIs(wsRaw.String)
	if wsClean == "" {
		return
	}
	if _, exists := rankMap[wsClean]; !exists {
		rankMap[wsClean] = rank
	}
}

func assignTitleRank(title sql.NullString, rank int, rankMap map[string]int) {
	if !title.Valid || title.String == "" {
		return
	}
	tLower := strings.ToLower(title.String)
	if _, exists := rankMap[tLower]; !exists {
		rankMap[tLower] = rank
	}
}

func fetchPinnedConversationsOrder() []string {
	appStoragePath := getAntigravityAppStoragePath()
	if appStoragePath == "" {
		return nil
	}

	data, err := os.ReadFile(appStoragePath)
	if err != nil {
		return nil
	}

	var rawMap map[string]interface{}
	if unmarshalErr := json.Unmarshal(data, &rawMap); unmarshalErr != nil {
		return nil
	}

	rawOrder, hasOrder := rawMap["pinned_conversations_order"]
	if !hasOrder {
		return nil
	}

	orderStr, isStr := rawOrder.(string)
	if !isStr || orderStr == "" {
		return nil
	}

	var convIDs []string
	if err := json.Unmarshal([]byte(orderStr), &convIDs); err != nil {
		return nil
	}

	return convIDs
}

func getAntigravityAppStoragePath() string {
	appData := os.Getenv("APPDATA")
	winPath := filepath.Join(appData, "Antigravity", "app_storage.json")
	if appData != "" && checkFileExists(winPath) {
		return winPath
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	candidates := []string{
		filepath.Join(home, ".config", "Antigravity", "app_storage.json"),
		filepath.Join(home, "Library", "Application Support", "Antigravity", "app_storage.json"),
	}

	for _, c := range candidates {
		if checkFileExists(c) {
			return c
		}
	}

	return ""
}

func resolveClosestActiveProject(projects []AgyProject, target string) (AgyProject, int) {
	if len(projects) == 0 {
		return AgyProject{}, 1
	}

	trimmed := strings.TrimSpace(target)
	if isHelpKeyword(trimmed) {
		return AgyProject{}, 0
	}

	if isCwdTarget(trimmed) {
		return findCwdProjectOrDefault(projects)
	}

	seq, err := strconv.Atoi(trimmed)
	if err == nil && seq >= 1 && seq <= len(projects) {
		return projects[seq-1], seq
	}

	matched, matchErr := findProjectByFlexibleTarget(projects, trimmed)
	if matchErr == nil {
		return matched, findProjectSequence(projects, matched.ID)
	}

	return projects[0], 1
}

func isCwdTarget(target string) bool {
	return target == "" || target == "last" || target == "current" || target == "."
}

func findCwdProject(projects []AgyProject) (AgyProject, int, bool) {
	cleanCwd, cleanRealCwd := resolveCurrentCleanCwds()
	if cleanCwd == "" {
		return AgyProject{}, 0, false
	}

	for i, p := range projects {
		if isCwdMatchingProject(p, cleanCwd, cleanRealCwd) {
			return p, i + 1, true
		}
	}

	return AgyProject{}, 0, false
}

func findCwdProjectOrDefault(projects []AgyProject) (AgyProject, int) {
	cwdProj, cwdSeq, ok := findCwdProject(projects)
	if ok && isProjectActiveOrRunning(cwdProj) {
		return cwdProj, cwdSeq
	}
	if len(projects) > 0 && isProjectActiveOrRunning(projects[0]) {
		return projects[0], 1
	}
	if ok {
		return cwdProj, cwdSeq
	}

	return projects[0], 1
}

func isProjectActiveOrRunning(p AgyProject) bool {
	_, _, hasActivePrompt := checkActivePromptFile(p.GetPath())
	if hasActivePrompt {
		return true
	}
	runningMap := fetchRunningProjectsSet()

	return resolveProjectIsRunning(p, runningMap)
}
