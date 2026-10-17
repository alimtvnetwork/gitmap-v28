package store

import (
	"os"
	"path/filepath"
)

// ResolveHeatmapRootDbPath resolves the root SQLite database file path for global heatmap data.
func ResolveHeatmapRootDbPath(repoRoot string) string {
	baseDir := resolveBaseDataDir(repoRoot)
	_ = os.MkdirAll(baseDir, 0755)

	return filepath.ToSlash(filepath.Join(baseDir, "heatmap.db"))
}

// ResolveHeatmapRepoSplitDbPath resolves the repository-specific split SQLite database path for heatmap data.
func ResolveHeatmapRepoSplitDbPath(repoSlug, repoRoot string) string {
	return ResolveSplitDbPath("heatmap", repoSlug, repoRoot)
}

// HeatmapRootDbPath returns the default root heatmap database path using ambient workspace.
func HeatmapRootDbPath() string {
	return ResolveHeatmapRootDbPath("")
}

// HeatmapRepoSplitDbPath returns the default repository split heatmap database path.
func HeatmapRepoSplitDbPath(repoSlug string) string {
	return ResolveHeatmapRepoSplitDbPath(repoSlug, "")
}
