// Package cmdagy — agy_rerun_project_resolve.go resolves active Antigravity projects by recent activity and pinned order.
package cmdagy

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type projectSortEntry struct {
	project         AgyProject
	isCwdMatch      bool
	isRunning       bool
	hasActivePrompt bool
	lastActTime     string
	pinnedRank      int
}

func loadActiveSortedProjects() ([]AgyProject, error) {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return nil, apperror.WrapSimple(err, "projects dir")
	}

	allProjects, loadErr := DiscoverUnifiedAgyProjects(dirPath)
	if loadErr != nil || len(allProjects) == 0 {
		return nil, apperror.NewSimple("no Antigravity projects configured", "E9030")
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
		_, _, hasActivePrompt := checkActivePromptFile(p.GetPath())
		if hasActivePrompt {
			isRunning = true
			actTime = resolvePromptActiveTime(p.GetPath(), actTime)
		}
		entries[i] = projectSortEntry{
			project:         p,
			isCwdMatch:      isCwd,
			isRunning:       isRunning,
			hasActivePrompt: hasActivePrompt,
			lastActTime:     actTime,
			pinnedRank:      pRank,
		}
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return compareProjectEntries(entries[i], entries[j])
	})

	for i := range entries {
		projects[i] = entries[i].project
	}
}

func resolvePromptActiveTime(ws, currentAct string) string {
	promptModTime := resolvePromptFileModTime(ws)
	if promptModTime > currentAct {
		return promptModTime
	}
	return currentAct
}

func resolvePromptFileModTime(ws string) string {
	paths := []string{
		filepath.Join(ws, ".ai-memory", "temp", "active-agy-pipeline-fix-prompt.txt"),
		filepath.Join(ws, "active-agy-pipeline-fix-prompt.txt"),
		filepath.Join(ws, "agy-prompt-queue.json"),
	}
	var latest time.Time
	for _, p := range paths {
		fi, err := os.Stat(p)
		if err != nil {
			continue
		}
		if fi.ModTime().After(latest) {
			latest = fi.ModTime()
		}
	}
	if !latest.IsZero() {
		return latest.UTC().Format(time.RFC3339)
	}
	return ""
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
	if a.hasActivePrompt != b.hasActivePrompt {
		return a.hasActivePrompt
	}

	if a.isRunning != b.isRunning {
		return a.isRunning
	}

	if a.lastActTime != b.lastActTime {
		return a.lastActTime > b.lastActTime
	}

	if a.isCwdMatch != b.isCwdMatch {
		return a.isCwdMatch
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
	var maxTime string
	if actTime, ok := actMap[p.ID]; ok && actTime != "" {
		maxTime = actTime
	}
	cleanWs := cleanProjectWorkspace(p.GetPath())
	if actTime, ok := actMap[cleanWs]; ok && actTime > maxTime {
		maxTime = actTime
	}
	nameLower := strings.ToLower(p.Name)
	if actTime, ok := actMap[nameLower]; ok && actTime > maxTime {
		maxTime = actTime
	}
	if p.UpdatedAt > maxTime {
		maxTime = p.UpdatedAt
	}
	return maxTime
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

	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return result
	}
	conn := openRes.Value
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
