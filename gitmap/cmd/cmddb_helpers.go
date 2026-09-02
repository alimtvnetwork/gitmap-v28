package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// DBFileInfo holds descriptive metadata for any discovered SQLite database.
type DBFileInfo struct {
	Name     string
	Path     string
	Size     int64
	Category string
	Purpose  string
	RepoID   int64
	RepoSlug string
}

func findSplitDBDirs() []string {
	binDataDir := store.BinaryDataDir()
	binRoot := filepath.Dir(binDataDir)
	raw := []string{
		filepath.Join(binDataDir, "repo_search"),
		filepath.Join(binRoot, constants.DefaultOutputDir, "repo_search"),
		filepath.Join(".", "data", "repo_search"),
		filepath.Join(".", constants.DefaultOutputDir, "repo_search"),
	}
	return dedupeDirs(raw)
}

func dedupeDirs(raw []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, d := range raw {
		clean := filepath.Clean(d)
		if !seen[clean] {
			seen[clean] = true
			if info, err := os.Stat(clean); err == nil && info.IsDir() {
				out = append(out, clean)
			}
		}
	}
	return out
}

func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024.0)
	}
	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes)/(1024.0*1024.0))
	}
	return fmt.Sprintf("%.2f GB", float64(bytes)/(1024.0*1024.0*1024.0))
}

func promptConfirm(msg string) (bool, error) {
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	ans := strings.ToLower(strings.TrimSpace(line))
	hasConfirmed := ans == "y" || ans == "yes"
	return hasConfirmed, nil
}

func hasConfirmFlag(args []string) bool {
	for _, a := range args {
		trimmed := strings.TrimSpace(a)
		if trimmed == "-y" || trimmed == "--yes" || trimmed == "-f" || trimmed == "--force" || trimmed == "--confirm" {
			return true
		}
	}
	return false
}

func collectMainDBInfo() (DBFileInfo, bool) {
	mainPath := store.DefaultDBPath()
	info, err := os.Stat(mainPath)
	if err != nil {
		return DBFileInfo{
			Name:     filepath.Base(mainPath),
			Path:     mainPath,
			Category: "Primary Master DB",
			Purpose:  "Central SQLite database storing global tracked repositories, scan history, configurations, and profiles.",
		}, false
	}
	return DBFileInfo{
		Name:     filepath.Base(mainPath),
		Path:     mainPath,
		Size:     info.Size(),
		Category: "Primary Master DB",
		Purpose:  "Central SQLite database storing global tracked repositories, scan history, configurations, and profiles.",
	}, true
}
