// Package cmdagy — agy_rerun_project_resolve.go resolves active Antigravity projects by recent activity and pinned order.
package cmdagy

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type projectSortEntry struct {
	project     AgyProject
	isCwdMatch  bool
	isRunning   bool
	lastActTime string
	pinnedRank  int
}

func loadActiveSortedProjects() ([]AgyProject, error) {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return nil, apperror.WrapSimple(err, "projects dir")
	}

	allProjects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil || len(allProjects) == 0 {
		return nil, apperror.NewSimple("no Antigravity projects configured in ~/.gemini/config/projects", "E9030")
	}

	validProjects := filterExistingProjects(allProjects)
	if len(validProjects) == 0 {
		return allProjects, nil
	}

	sortProjectsByActivityAndPins(validProjects)
	return validProjects, nil
}

func filterExistingProjects(projects []AgyProject) []AgyProject {
	var valid []AgyProject
	for _, p := range projects {
		ws := p.GetPath()
		if ws == "" {
			continue
		}
		if checkDirExists(ws) {
			valid = append(valid, p)
		}
	}
	return valid
}

func sortProjectsByActivityAndPins(projects []AgyProject) {
	activityMap := fetchProjectLatestActivityMap()
	pinnedMap := fetchPinnedProjectRankMap()
	runningMap := fetchRunningProjectsSet()
	cleanCwd, cleanRealCwd := resolveCurrentCleanCwds()

	entries := make([]projectSortEntry, len(projects))
	for i, p := range projects {
		actTime := resolveProjectActTime(p, activityMap)
		pRank := resolveProjectPinnedRank(p, pinnedMap)
		isCwd := isCwdMatchingProject(p, cleanCwd, cleanRealCwd)
		isRunning := resolveProjectIsRunning(p, runningMap)
		entries[i] = projectSortEntry{
			project:     p,
			isCwdMatch:  isCwd,
			isRunning:   isRunning,
			lastActTime: actTime,
			pinnedRank:  pRank,
		}
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return compareProjectEntries(entries[i], entries[j])
	})

	for i := range entries {
		projects[i] = entries[i].project
	}
}

func resolveCurrentCleanCwds() (string, string) {
	cwd, err := os.Getwd()
	if err != nil || cwd == "" {
		return "", ""
	}

	cleanCwd := cleanProjectWorkspace(cwd)
	cleanRealCwd := resolveCleanRealCwd(cwd)

	return cleanCwd, cleanRealCwd
}

func resolveProjectIsRunning(p AgyProject, runningMap map[string]bool) bool {
	if runningMap[p.ID] {
		return true
	}
	cleanWs := cleanProjectWorkspace(p.GetPath())
	if runningMap[cleanWs] {
		return true
	}
	nameLower := strings.ToLower(p.Name)
	if runningMap[nameLower] {
		return true
	}

	return false
}

func compareProjectEntries(a, b projectSortEntry) bool {
	if a.isCwdMatch != b.isCwdMatch {
		return a.isCwdMatch
	}

	if a.isRunning != b.isRunning {
		return a.isRunning
	}

	if a.lastActTime != b.lastActTime {
		return a.lastActTime > b.lastActTime
	}

	if a.pinnedRank != b.pinnedRank {
		return a.pinnedRank < b.pinnedRank
	}

	if a.project.UpdatedAt != b.project.UpdatedAt {
		return a.project.UpdatedAt > b.project.UpdatedAt
	}

	return strings.ToLower(a.project.Name) < strings.ToLower(b.project.Name)
}

func resolveProjectActTime(p AgyProject, actMap map[string]string) string {
	if actTime, ok := actMap[p.ID]; ok && actTime != "" {
		return actTime
	}

	cleanWs := cleanProjectWorkspace(p.GetPath())
	if actTime, ok := actMap[cleanWs]; ok && actTime != "" {
		return actTime
	}

	return p.UpdatedAt
}

func resolveProjectPinnedRank(p AgyProject, pinMap map[string]int) int {
	if rank, ok := pinMap[p.ID]; ok {
		return rank
	}

	cleanWs := cleanProjectWorkspace(p.GetPath())
	if rank, ok := pinMap[cleanWs]; ok {
		return rank
	}

	nameLower := strings.ToLower(p.Name)
	if rank, ok := pinMap[nameLower]; ok {
		return rank
	}

	return 999999
}

func fetchProjectLatestActivityMap() map[string]string {
	result := make(map[string]string)
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return result
	}

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return result
	}
	defer conn.Close()

	query := `SELECT project_id, workspace_uris, MAX(last_modified_time)
		FROM conversation_summaries
		WHERE (killed IS NULL OR killed = 0)
		GROUP BY COALESCE(NULLIF(project_id, ''), workspace_uris)`

	rows, qErr := conn.Query(query)
	if qErr != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var pid, wsRaw, maxTime sql.NullString
		if scanErr := rows.Scan(&pid, &wsRaw, &maxTime); scanErr != nil {
			continue
		}
		timeVal := maxTime.String
		if pid.Valid && pid.String != "" {
			result[pid.String] = timeVal
		}
		recordWorkspaceActivity(wsRaw, timeVal, result)
	}

	return result
}

func fetchRunningProjectsSet() map[string]bool {
	result := make(map[string]bool)
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return result
	}

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return result
	}
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

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return rankMap
	}
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
	return target == "" || target == "1" || target == "last" || target == "current" || target == "."
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
	if proj, seq, ok := findCwdProject(projects); ok {
		return proj, seq
	}

	return projects[0], 1
}

func isCwdMatchingProject(p AgyProject, cleanCwd, cleanRealCwd string) bool {
	pWs := cleanProjectWorkspace(p.GetPath())
	if isWorkspaceMatchOrSubdir(pWs, cleanCwd, cleanRealCwd) {
		return true
	}
	realWs, err := filepath.EvalSymlinks(p.GetPath())
	if err != nil || realWs == "" {
		return false
	}
	cleanRealWs := cleanProjectWorkspace(realWs)

	return isWorkspaceMatchOrSubdir(cleanRealWs, cleanCwd, cleanRealCwd)
}

func isWorkspaceMatchOrSubdir(ws, cleanCwd, cleanRealCwd string) bool {
	if ws == "" {
		return false
	}
	if ws == cleanCwd || strings.HasPrefix(cleanCwd, ws+"/") {
		return true
	}
	if cleanRealCwd != "" && (ws == cleanRealCwd || strings.HasPrefix(cleanRealCwd, ws+"/")) {
		return true
	}

	return false
}

func resolveCleanRealCwd(cwd string) string {
	realCwd, err := filepath.EvalSymlinks(cwd)
	if err != nil || realCwd == "" {
		return ""
	}
	return cleanProjectWorkspace(realCwd)
}

func findProjectByFlexibleTarget(projects []AgyProject, target string) (AgyProject, error) {
	matched, err := ResolveAgyProjectTargets([]string{target}, projects)
	if err == nil && len(matched) > 0 {
		return matched[0], nil
	}

	cleanTarget := normalizeSlugForComparison(target)
	for _, p := range projects {
		nameNorm := normalizeSlugForComparison(p.Name)
		wsNorm := normalizeSlugForComparison(filepath.Base(p.GetPath()))
		if strings.Contains(nameNorm, cleanTarget) || strings.Contains(wsNorm, cleanTarget) {
			return p, nil
		}
		if cleanTarget != "" && (strings.Contains(cleanTarget, nameNorm) || strings.Contains(cleanTarget, wsNorm)) {
			return p, nil
		}
		if isPhoneticOrAliasMatch(cleanTarget, nameNorm) || isPhoneticOrAliasMatch(cleanTarget, wsNorm) {
			return p, nil
		}
	}

	return AgyProject{}, apperror.NewNotFound("project_target", "E9039", "no project matching target: "+target)
}

func isPhoneticOrAliasMatch(target, candidate string) bool {
	if strings.Contains(target, "xampp") && strings.Contains(candidate, "exam") {
		return true
	}
	if strings.Contains(target, "exam") && strings.Contains(candidate, "xampp") {
		return true
	}
	return false
}

func normalizeSlugForComparison(raw string) string {
	s := strings.ToLower(raw)
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ".", "")
	return s
}
