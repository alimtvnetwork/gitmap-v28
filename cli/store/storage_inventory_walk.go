package store

import (
	"os"
	"path/filepath"
	"strings"
)

// CollectDataTreeDatabases scans a directory tree for SQLite database files.
func CollectDataTreeDatabases(rootDir string, add func(SplitDatabaseEntry)) {
	if !isDirExisting(rootDir) {
		return
	}

	_ = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".db") {
			processDiscoveredDbFile(path, info.Name(), add)
		}
		return nil
	})
}

func processDiscoveredDbFile(path, fileName string, add func(SplitDatabaseEntry)) {
	cleanPath := filepath.Clean(path)
	slashPath := filepath.ToSlash(cleanPath)
	dbType, dbKey, desc := classifyDiscoveredDb(slashPath, fileName)
	add(inspectSplitDBFile(dbType, dbKey, cleanPath, desc))
}

func classifyDiscoveredDb(slashPath, fileName string) (string, string, string) {
	slug := strings.TrimSuffix(fileName, ".db")
	if strings.Contains(slashPath, "/tasks/") || strings.HasSuffix(fileName, "-tasks.db") {
		return resolveTasksDbMetadata(slug)
	}
	if strings.Contains(slashPath, "/pipeline/") {
		return "pipeline", resolveSectionSlug(slashPath, slug), "Pipeline split database"
	}
	if strings.Contains(slashPath, "/repodb/") {
		return "repodb", slug, "Repository split database"
	}
	if strings.Contains(slashPath, "/ssh/") {
		return "ssh", slug, "SSH tasks and history database"
	}
	if strings.Contains(slashPath, "/search/") {
		return "search", "search", "Search index split database"
	}
	if strings.Contains(slashPath, "/ai-instruction/") {
		return "ai-instruction", "ai-instruction", "AI instruction split database"
	}
	if strings.Contains(slashPath, "/installation/") {
		return "split", "installation", "Tool installer registry database"
	}
	if strings.Contains(slashPath, "/startup/") {
		return "split", "startup", "Startup manager database"
	}

	return "split", slug, "Gitmap split database"
}

func resolveSectionSlug(slashPath, slug string) string {
	if slug != "sql" {
		return slug
	}
	parts := strings.Split(slashPath, "/")
	if len(parts) < 2 {
		return slug
	}
	parent := parts[len(parts)-2]
	if parent == "pipeline" || parent == "data" {
		return slug
	}

	return parent
}

func resolveTasksDbMetadata(slug string) (string, string, string) {
	if slug == "sql" {
		return "tasks", "tasks", "Tasks root split database"
	}

	return "tasks", slug, "Section tasks split database"
}
