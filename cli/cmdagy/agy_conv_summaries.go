// Package cmdagy provides Antigravity conversation summary discovery and project matching.
package cmdagy

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type summaryRow struct {
	id           string
	title        string
	workspaceURI string
	projectID    string
	depth        int
	stepCount    int
}

type scoredSummary struct {
	info  AgyConvInfo
	score int
}

func getConversationSummariesDBPath() (string, error) {
	home, err := os.UserHomeDir()
	hasErr := err != nil
	if hasErr {
		return "", err
	}
	dbPath := filepath.Join(home, ".gemini", "antigravity", "conversation_summaries.db")
	_, statErr := os.Stat(dbPath)
	hasStatErr := statErr != nil
	if hasStatErr {
		return "", statErr
	}

	return dbPath, nil
}

func scanConversationsFromSummaries(pClean, pName string) ([]AgyConvInfo, error) {
	dbPath, err := getConversationSummariesDBPath()
	hasErr := err != nil
	if hasErr {
		return nil, err
	}
	conn, openErr := store.OpenSQLiteDB(dbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}
	defer conn.Close()

	return queryAndScoreSummaries(conn, pClean, pName)
}

func queryAndScoreSummaries(conn *sql.DB, pClean, pName string) ([]AgyConvInfo, error) {
	rows, err := fetchSummaryRows(conn)
	hasErr := err != nil
	if hasErr {
		return nil, err
	}
	projectIDs := findProjectIDsForRepo(pClean, pName)
	candidates := scoreAllSummaryRows(rows, pClean, pName, projectIDs)
	sortScoredSummaries(candidates)

	return extractConvInfos(candidates), nil
}

func fetchSummaryRows(conn *sql.DB) ([]summaryRow, error) {
	query := "SELECT conversation_id, title, workspace_uris, project_id, nesting_depth, step_count FROM conversation_summaries ORDER BY last_modified_time DESC"
	qRows, err := conn.Query(query)
	hasErr := err != nil
	if hasErr {
		return nil, err
	}
	defer qRows.Close()

	return iterateSummaryRows(qRows)
}

func iterateSummaryRows(rows *sql.Rows) ([]summaryRow, error) {
	var list []summaryRow
	for rows.Next() {
		var r summaryRow
		err := rows.Scan(&r.id, &r.title, &r.workspaceURI, &r.projectID, &r.depth, &r.stepCount)
		hasErr := err != nil
		if hasErr {
			continue
		}
		list = append(list, r)
	}

	return list, nil
}

func scoreAllSummaryRows(rows []summaryRow, pClean, pName string, projectIDs []string) []scoredSummary {
	var scored []scoredSummary
	for _, r := range rows {
		score := computeSummaryScore(r, pClean, pName, projectIDs)
		isCandidate := score > 0
		if isCandidate {
			info := AgyConvInfo{ID: r.id, StepCount: r.stepCount, CleanPath: pClean}
			scored = append(scored, scoredSummary{info: info, score: score})
		}
	}

	return scored
}

func computeSummaryScore(r summaryRow, pClean, pName string, projectIDs []string) int {
	isPathMatch := checkWorkspaceURIMatch(r.workspaceURI, pClean)
	isNameMatch := checkProjectNameMatch(r, pName, projectIDs)
	hasNeither := isPathMatch == false && isNameMatch == false
	if hasNeither {
		return 0
	}
	score := computeBaseMatchScore(isPathMatch, isNameMatch)
	isRoot := r.depth == 0
	if isRoot {
		score += 300
	}

	return score
}

func computeBaseMatchScore(isPathMatch, isNameMatch bool) int {
	isBoth := isPathMatch && isNameMatch
	if isBoth {
		return 1000
	}
	isPathOnly := isPathMatch && isNameMatch == false
	if isPathOnly {
		return 500
	}

	return 200
}

func checkWorkspaceURIMatch(urisJSON, pClean string) bool {
	hasEmpty := urisJSON == "" || pClean == ""
	if hasEmpty {
		return false
	}
	lowerURIs := strings.ToLower(urisJSON)
	matchesSub := fileURIRegex.FindString(lowerURIs)
	cleanPath := cleanURIStringToPath(matchesSub)
	isMatch := cleanPath != "" && isConvPathMatch(pClean, cleanPath)

	return isMatch
}

func checkProjectNameMatch(r summaryRow, pName string, projectIDs []string) bool {
	hasEmptyName := pName == ""
	if hasEmptyName {
		return false
	}
	isIDMatch := isProjectIDInList(r.projectID, projectIDs)
	if isIDMatch {
		return true
	}
	isTitleExact := strings.EqualFold(r.title, pName)
	if isTitleExact {
		return true
	}

	return strings.Contains(strings.ToLower(r.title), strings.ToLower(pName))
}

func isProjectIDInList(id string, list []string) bool {
	hasEmptyID := id == ""
	if hasEmptyID {
		return false
	}
	for _, item := range list {
		isMatch := item == id
		if isMatch {
			return true
		}
	}

	return false
}

func findProjectIDsForRepo(pClean, pName string) []string {
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

	return collectMatchingProjectIDs(projects, pClean, pName)
}

func collectMatchingProjectIDs(projects []AgyProject, pClean, pName string) []string {
	var ids []string
	for _, p := range projects {
		isMatch := isProjectMatch(p, pClean, pName)
		if isMatch {
			ids = append(ids, p.ID)
		}
	}

	return ids
}

func isProjectMatch(p AgyProject, pClean, pName string) bool {
	isNameMatch := strings.EqualFold(p.Name, pName)
	if isNameMatch {
		return true
	}
	pPathClean := cleanProjectWorkspace(p.GetPath())
	isPathMatch := pPathClean != "" && pPathClean == pClean

	return isPathMatch
}

func sortScoredSummaries(items []scoredSummary) {
	sort.Slice(items, func(i, j int) bool {
		hasDiffScore := items[i].score != items[j].score
		if hasDiffScore {
			return items[i].score > items[j].score
		}
		hasDiffSteps := items[i].info.StepCount != items[j].info.StepCount
		if hasDiffSteps {
			return items[i].info.StepCount > items[j].info.StepCount
		}

		return items[i].info.ID < items[j].info.ID
	})
}

func extractConvInfos(scored []scoredSummary) []AgyConvInfo {
	out := make([]AgyConvInfo, 0, len(scored))
	for _, s := range scored {
		out = append(out, s.info)
	}

	return out
}
