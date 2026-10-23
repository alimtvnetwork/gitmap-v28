package heatmap

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// DiscoveredHeatmapDB represents a discovered heatmap database across all 3 tiers.
type DiscoveredHeatmapDB struct {
	Scope        string    `json:"scope"`
	AiName       string    `json:"aiName,omitempty"`
	RepoSlug     string    `json:"repoSlug,omitempty"`
	RepoAlias    string    `json:"repoAlias,omitempty"`
	Path         string    `json:"path"`
	SizeBytes    int64     `json:"sizeBytes"`
	ActionCount  int64     `json:"actionCount"`
	TaskCount    int64     `json:"taskCount"`
	LastModified time.Time `json:"lastModified"`
}

func formatHeatmapFileSize(b int64) string {
	if b < 1024 {
		return fmt.Sprintf("%d B", b)
	}

	if b < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(b)/1024.0)
	}

	return fmt.Sprintf("%.1f MB", float64(b)/(1024.0*1024.0))
}

func queryDbCounts(dbPath string) (int64, int64) {
	dsn := fmt.Sprintf("file:%s?mode=ro&cache=shared", filepath.ToSlash(dbPath))
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return 0, 0
	}

	defer conn.Close()

	var actionCount, taskCount int64
	query := "SELECT COUNT(*), COUNT(DISTINCT TaskSlug) FROM HeatmapActionLog"
	err = conn.QueryRow(query).Scan(&actionCount, &taskCount)
	if err != nil {
		return 0, 0
	}

	return actionCount, taskCount
}

func newDiscoveredDB(scope, ai, repo, path string, info os.FileInfo, act, tsk int64) DiscoveredHeatmapDB {
	return DiscoveredHeatmapDB{
		Scope:        scope,
		AiName:       ai,
		RepoSlug:     repo,
		Path:         path,
		SizeBytes:    info.Size(),
		ActionCount:  act,
		TaskCount:    tsk,
		LastModified: info.ModTime(),
	}
}

func scanTier1LocalRoot() ([]DiscoveredHeatmapDB, error) {
	rootPath := filepath.ToSlash(filepath.Join(store.EffectiveDataDir(), "heatmap.db"))
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, nil
	}

	return []DiscoveredHeatmapDB{buildLocalRootDbEntry(rootPath, info)}, nil
}

func buildLocalRootDbEntry(path string, info os.FileInfo) DiscoveredHeatmapDB {
	actions, tasks := queryDbCounts(path)

	return newDiscoveredDB("local_root", store.DetectCurrentAIName(), "root", path, info, actions, tasks)
}

func scanTier2LocalSplit() ([]DiscoveredHeatmapDB, error) {
	heatmapDir := filepath.Join(store.EffectiveDataDir(), "heatmap")
	entries, err := os.ReadDir(heatmapDir)
	if err != nil {
		return nil, nil
	}

	var dbs []DiscoveredHeatmapDB
	for _, entry := range entries {
		if db, hasDb := probeTier2Entry(heatmapDir, entry); hasDb {
			dbs = append(dbs, db)
		}
	}

	return dbs, nil
}

func probeTier2Entry(baseDir string, entry os.DirEntry) (DiscoveredHeatmapDB, bool) {
	if !entry.IsDir() {
		return DiscoveredHeatmapDB{}, false
	}

	repoSlug := entry.Name()
	dbPath := filepath.ToSlash(filepath.Join(baseDir, repoSlug, "sql.db"))
	info, err := os.Stat(dbPath)
	if err != nil {
		return DiscoveredHeatmapDB{}, false
	}

	actions, tasks := queryDbCounts(dbPath)

	return newDiscoveredDB("local_split_repo", "", repoSlug, dbPath, info, actions, tasks), true
}

func scanTier3RepoCache() ([]DiscoveredHeatmapDB, error) {
	baseDir := store.ResolveRcHeatmapBaseDir()
	aiFolders, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, nil
	}

	var dbs []DiscoveredHeatmapDB
	for _, folder := range aiFolders {
		dbs = append(dbs, scanAiNamespace(baseDir, folder)...)
	}

	return dbs, nil
}

func scanAiNamespace(baseDir string, folder os.DirEntry) []DiscoveredHeatmapDB {
	if !folder.IsDir() {
		return nil
	}

	aiDir := filepath.Join(baseDir, folder.Name())
	aiName := store.NormalizeAIName(folder.Name())
	var result []DiscoveredHeatmapDB
	if rootDb, hasRoot := probeAiRootDb(aiDir, aiName); hasRoot {
		result = append(result, rootDb)
	}

	result = append(result, probeAiRepoDbs(aiDir, aiName)...)

	return result
}

func probeAiRootDb(aiDir, aiName string) (DiscoveredHeatmapDB, bool) {
	rootPath := filepath.ToSlash(filepath.Join(aiDir, "root-heatmap.db"))
	info, err := os.Stat(rootPath)
	if err != nil {
		return DiscoveredHeatmapDB{}, false
	}

	actions, tasks := queryDbCounts(rootPath)

	return newDiscoveredDB("repo_cache_root", aiName, "root", rootPath, info, actions, tasks), true
}

func probeAiRepoDbs(aiDir, aiName string) []DiscoveredHeatmapDB {
	entries, err := os.ReadDir(aiDir)
	if err != nil {
		return nil
	}

	var dbs []DiscoveredHeatmapDB
	for _, entry := range entries {
		if db, hasDb := probeSingleAiRepo(aiDir, aiName, entry); hasDb {
			dbs = append(dbs, db)
		}
	}

	return dbs
}

func probeSingleAiRepo(aiDir, aiName string, entry os.DirEntry) (DiscoveredHeatmapDB, bool) {
	if !entry.IsDir() {
		return DiscoveredHeatmapDB{}, false
	}

	repoSlug := entry.Name()
	dbPath := filepath.ToSlash(filepath.Join(aiDir, repoSlug, "heatmap.db"))
	info, err := os.Stat(dbPath)
	if err != nil {
		return DiscoveredHeatmapDB{}, false
	}

	actions, tasks := queryDbCounts(dbPath)

	return newDiscoveredDB("repo_cache_repo", aiName, repoSlug, dbPath, info, actions, tasks), true
}

// DiscoverHeatmapDatabases surveys all three tiers of heatmap databases on the active host.
func DiscoverHeatmapDatabases() ([]DiscoveredHeatmapDB, error) {
	var allDbs []DiscoveredHeatmapDB
	t1, _ := scanTier1LocalRoot()
	allDbs = append(allDbs, t1...)
	t2, _ := scanTier2LocalSplit()
	allDbs = append(allDbs, t2...)
	t3, _ := scanTier3RepoCache()
	allDbs = append(allDbs, t3...)

	return allDbs, nil
}

// DiscoverAllHeatmapDatabases is an alias for DiscoverHeatmapDatabases.
func DiscoverAllHeatmapDatabases() ([]DiscoveredHeatmapDB, error) {
	return DiscoverHeatmapDatabases()
}

// EnrichHeatmapDatabasesWithAliases matches discovered database slugs against registered aliases.
func EnrichHeatmapDatabasesWithAliases(dbs []DiscoveredHeatmapDB) []DiscoveredHeatmapDB {
	aliasMap := loadAliasMap()
	if len(aliasMap) == 0 {
		return dbs
	}

	for i := range dbs {
		dbs[i] = applyRepoAlias(dbs[i], aliasMap)
	}

	return dbs
}

func loadAliasMap() map[string]string {
	master, err := store.OpenDefault()
	if err != nil {
		return nil
	}

	defer master.Close()

	items, err := master.ListAliasesWithRepo()
	if err != nil {
		return nil
	}

	return buildSlugToAliasIndex(items)
}

func buildSlugToAliasIndex(items []store.AliasWithRepo) map[string]string {
	aliasMap := make(map[string]string, len(items))
	for _, item := range items {
		slug := strings.ToLower(strings.TrimSpace(item.Slug))
		if item.IsPrimary || aliasMap[slug] == "" {
			aliasMap[slug] = item.Alias.Alias
		}
	}

	return aliasMap
}

func applyRepoAlias(db DiscoveredHeatmapDB, aliasMap map[string]string) DiscoveredHeatmapDB {
	slug := strings.ToLower(strings.TrimSpace(db.RepoSlug))
	if alias, hasAlias := aliasMap[slug]; hasAlias {
		db.RepoAlias = alias
	}

	return db
}
