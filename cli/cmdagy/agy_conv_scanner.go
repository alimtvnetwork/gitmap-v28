package cmdagy

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	convScanMu     sync.RWMutex
	lastConvScanAt time.Time
	cachedConvList []AgyConvInfo
)

const convCacheTTL = 10 * time.Second

// InvalidateConvCache clears the in-memory conversation and transcript TTL caches.
func InvalidateConvCache() {
	convScanMu.Lock()
	lastConvScanAt = time.Time{}
	cachedConvList = nil
	convScanMu.Unlock()

	invalidateTranscriptCache()
}

func copyConvList(src []AgyConvInfo) []AgyConvInfo {
	if src == nil {
		return nil
	}
	dst := make([]AgyConvInfo, len(src))
	copy(dst, src)
	return dst
}

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

var fileURIRegex = lazyregex.FileUriRegex

func getConversationsDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".gemini", "antigravity", "conversations"), nil
}

func scanAllConversations() ([]AgyConvInfo, error) {
	convScanMu.RLock()
	if !lastConvScanAt.IsZero() && time.Since(lastConvScanAt) < convCacheTTL && cachedConvList != nil {
		cached := copyConvList(cachedConvList)
		convScanMu.RUnlock()
		return cached, nil
	}
	convScanMu.RUnlock()

	dir, err := getConversationsDirPath()
	hasDirErr := err != nil
	if hasDirErr {
		return nil, err
	}

	convs, err := readConversationsFromDir(dir)
	if err != nil {
		return nil, err
	}

	convScanMu.Lock()
	lastConvScanAt = time.Now()
	cachedConvList = copyConvList(convs)
	convScanMu.Unlock()

	return convs, nil
}

func readConversationsFromDir(dir string) ([]AgyConvInfo, error) {
	entries, err := os.ReadDir(dir)
	hasReadErr := err != nil
	if hasReadErr {
		return nil, err
	}

	return collectConvEntries(dir, entries), nil
}

func collectConvEntries(dir string, entries []os.DirEntry) []AgyConvInfo {
	var out []AgyConvInfo
	for _, e := range entries {
		info, isReadSuccess := tryReadConvEntry(dir, e)
		if isReadSuccess {
			out = append(out, info)
		}
	}

	return out
}

func tryReadConvEntry(dir string, e os.DirEntry) (AgyConvInfo, bool) {
	if e.IsDir() || !strings.HasSuffix(e.Name(), ".db") {
		return AgyConvInfo{}, false
	}

	full := filepath.Join(dir, e.Name())

	return readSingleConvDB(full, e.Name())
}

func readSingleConvDB(dbPath, fileName string) (AgyConvInfo, bool) {
	convID := strings.TrimSuffix(fileName, ".db")
	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return readConvFallback(convID), true
	}
	conn := openRes.Value
	defer conn.Close()

	return buildConvInfo(conn, convID), true
}

func readConvFallback(convID string) AgyConvInfo {
	cleanPath := extractWorkspaceFromTranscript(convID)

	return AgyConvInfo{
		ID:        convID,
		StepCount: 0,
		UserSteps: 0,
		CleanPath: cleanPath,
	}
}

func buildConvInfo(conn *sql.DB, convID string) AgyConvInfo {
	steps := querySingleCount(conn, "SELECT COUNT(*) FROM steps")
	userSteps := querySingleCount(conn, "SELECT COUNT(*) FROM steps WHERE step_type = 1")
	cleanPath := resolveConvPath(conn, convID)

	return AgyConvInfo{
		ID:        convID,
		StepCount: steps,
		UserSteps: userSteps,
		CleanPath: cleanPath,
	}
}

func resolveConvPath(conn *sql.DB, convID string) string {
	cleanPath := extractWorkspaceFromConv(conn)
	hasEmptyPath := cleanPath == ""
	if hasEmptyPath {
		return extractWorkspaceFromTranscript(convID)
	}

	return cleanPath
}

func extractWorkspaceFromConv(conn *sql.DB) string {
	blob := queryTrajectoryBlob(conn)
	hasEmptyBlob := len(blob) == 0
	if hasEmptyBlob {
		return ""
	}
	match := fileURIRegex.Find(blob)
	hasMatch := len(match) > 0
	if hasMatch {
		return cleanURIStringToPath(string(match))
	}

	return ""
}

func queryTrajectoryBlob(conn *sql.DB) []byte {
	var blob []byte
	row := conn.QueryRow("SELECT data FROM trajectory_metadata_blob WHERE id='main'")
	err := row.Scan(&blob)
	if err != nil {
		return nil
	}

	return blob
}

func mapProjectsToConversations(projects []AgyProject, convs []AgyConvInfo) []AgyProjectConvs {
	var results []AgyProjectConvs
	for _, p := range projects {
		isOutside := p.ID == "outside-of-project"
		if isOutside {
			continue
		}
		results = append(results, buildProjectConvs(p, convs))
	}

	return results
}

func buildProjectConvs(p AgyProject, convs []AgyConvInfo) AgyProjectConvs {
	pClean := cleanProjectWorkspace(p.GetPath())
	matching, hasActive := findMatchingConvs(pClean, convs)

	return AgyProjectConvs{
		Project:   p,
		Convs:     matching,
		HasActive: hasActive,
	}
}

func cleanProjectWorkspace(rawPath string) string {
	if rawPath == "" {
		return ""
	}

	return strings.ToLower(filepath.Clean(rawPath))
}

func appendMatchingConv(matched []AgyConvInfo, c AgyConvInfo, hasActive bool) ([]AgyConvInfo, bool) {
	matched = append(matched, c)
	if isConvActive(c) {
		return matched, true
	}

	return matched, hasActive
}

func findMatchingConvs(pClean string, convs []AgyConvInfo) ([]AgyConvInfo, bool) {
	var matched []AgyConvInfo
	hasActive := false
	for _, c := range convs {
		isMatch := isConvPathMatch(pClean, c.CleanPath)
		if isMatch {
			matched, hasActive = appendMatchingConv(matched, c, hasActive)
		}
	}

	return matched, hasActive
}

func isConvPathMatch(pClean, cClean string) bool {
	hasEmptyPath := pClean == "" || cClean == ""
	if hasEmptyPath {
		return false
	}
	pNorm := strings.ReplaceAll(filepath.ToSlash(strings.ToLower(pClean)), "\\", "/")
	cNorm := strings.ReplaceAll(filepath.ToSlash(strings.ToLower(cClean)), "\\", "/")
	isExact := pNorm == cNorm
	isSubdir := strings.HasPrefix(pNorm, cNorm+"/")
	isParent := strings.HasPrefix(cNorm, pNorm+"/")

	return isExact || isSubdir || isParent
}

func isConvActive(c AgyConvInfo) bool {
	return c.StepCount > 2 || c.UserSteps > 0
}

func querySingleCount(conn *sql.DB, query string) int {
	var count int
	row := conn.QueryRow(query)
	if err := row.Scan(&count); err != nil {
		return 0
	}

	return count
}
