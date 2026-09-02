package cmd

import (
	"database/sql"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type AgyConvInfo struct {
	ID        string
	StepCount int
	UserSteps int
	CleanPath string
}

type AgyProjectConvs struct {
	Project   AgyProject
	Convs     []AgyConvInfo
	HasActive bool
}

var fileURIRegex = regexp.MustCompile(`file:///[^\x00-\x1f\x7f-\xff"'\s]+`)

func getConversationsDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gemini", "antigravity", "conversations"), nil
}

func scanAllConversations() ([]AgyConvInfo, error) {
	dir, err := getConversationsDirPath()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []AgyConvInfo
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".db") {
			full := filepath.Join(dir, e.Name())
			if info, ok := readSingleConvDB(full, e.Name()); ok {
				out = append(out, info)
			}
		}
	}
	return out, nil
}

func readSingleConvDB(dbPath, fileName string) (AgyConvInfo, bool) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return AgyConvInfo{}, false
	}
	defer conn.Close()

	steps := querySingleCount(conn, "SELECT COUNT(*) FROM steps")
	userSteps := querySingleCount(conn, "SELECT COUNT(*) FROM steps WHERE step_type = 1")
	cleanPath := extractWorkspaceFromConv(conn)
	cid := strings.TrimSuffix(fileName, ".db")

	return AgyConvInfo{
		ID:        cid,
		StepCount: steps,
		UserSteps: userSteps,
		CleanPath: cleanPath,
	}, true
}

func extractWorkspaceFromConv(conn *sql.DB) string {
	var blob []byte
	row := conn.QueryRow("SELECT data FROM trajectory_metadata_blob WHERE id='main'")
	if err := row.Scan(&blob); err != nil || len(blob) == 0 {
		return ""
	}
	match := fileURIRegex.Find(blob)
	if len(match) == 0 {
		return ""
	}
	return cleanURIStringToPath(string(match))
}

func cleanURIStringToPath(rawURI string) string {
	trimmed := strings.TrimPrefix(rawURI, "file:///")
	decoded, err := url.PathUnescape(trimmed)
	if err != nil {
		decoded = trimmed
	}
	clean := filepath.Clean(filepath.FromSlash(decoded))
	return strings.ToLower(clean)
}

func mapProjectsToConversations(projects []AgyProject, convs []AgyConvInfo) []AgyProjectConvs {
	var results []AgyProjectConvs
	for _, p := range projects {
		if p.ID == "outside-of-project" {
			continue
		}
		pClean := cleanProjectWorkspace(p.GetPath())
		matching, hasActive := findMatchingConvs(pClean, convs)
		results = append(results, AgyProjectConvs{
			Project:   p,
			Convs:     matching,
			HasActive: hasActive,
		})
	}
	return results
}

func cleanProjectWorkspace(rawPath string) string {
	if rawPath == "" {
		return ""
	}
	return strings.ToLower(filepath.Clean(rawPath))
}

func findMatchingConvs(pClean string, convs []AgyConvInfo) ([]AgyConvInfo, bool) {
	var matched []AgyConvInfo
	hasActive := false
	for _, c := range convs {
		if isConvPathMatch(pClean, c.CleanPath) {
			matched = append(matched, c)
			if isConvActive(c) {
				hasActive = true
			}
		}
	}
	return matched, hasActive
}

func isConvPathMatch(pClean, cClean string) bool {
	if pClean == "" || cClean == "" {
		return false
	}
	return pClean == cClean || strings.HasPrefix(pClean, cClean) || strings.HasPrefix(cClean, pClean)
}

func isConvActive(c AgyConvInfo) bool {
	return c.StepCount > 2 || c.UserSteps > 0
}
