package store

import (
	"os"
	"path/filepath"
	"strings"
)

func normalizeStorageDbPath(rawPath string) string {
	clean := filepath.Clean(rawPath)
	if abs, err := filepath.Abs(clean); err == nil {
		clean = abs
	}

	return filepath.ToSlash(strings.ToLower(clean))
}

// CollectAllDatabaseEntries returns an inventory of all SQLite databases managed by Gitmap.
func CollectAllDatabaseEntries() []SplitDatabaseEntry {
	dataDir := BinaryDataDir()
	seen := make(map[string]bool)
	var list []SplitDatabaseEntry

	addEntry := func(e SplitDatabaseEntry) {
		norm := normalizeStorageDbPath(e.DatabasePath)
		if !seen[norm] {
			seen[norm] = true
			list = append(list, e)
		}
	}

	collectAllDataSources(dataDir, addEntry)

	return list
}

func collectAllDataSources(dataDir string, add func(SplitDatabaseEntry)) {
	collectCoreDatabases(dataDir, add)
	collectDataSubdirs(dataDir, add)
	collectParentSubdirDatabases(dataDir, "pipeline", "pipeline", add)
	collectParentSubdirDatabases(dataDir, "repodb", "repodb", add)
	collectUserHomeDatabases(add)
	collectLocalRepoDatabases(add)
	collectLooseDatabases(dataDir, add)
}

func collectDataSubdirs(dataDir string, add func(SplitDatabaseEntry)) {
	collectSubdirDatabases(dataDir, "schedule", "schedule", add)
	collectSubdirDatabases(dataDir, "schedules", "schedule", add)
	collectSubdirDatabases(dataDir, "pipeline", "pipeline", add)
	collectSubdirDatabases(dataDir, "automation", "automation", add)
	collectSubdirDatabases(dataDir, "pipeline_db", "pipeline", add)
	collectSubdirDatabases(dataDir, "repo_search", "repo", add)
	collectSubdirDatabases(dataDir, "repodb", "repodb", add)
}

func collectCoreDatabases(dataDir string, add func(SplitDatabaseEntry)) {
	rootPath := DefaultDBPath()
	add(inspectSplitDBFile("master", "root", rootPath, "Primary Gitmap repository tracking database"))

	instPath := ResolveSplitDbPath(SectionInstallation, "default", "")
	add(inspectSplitDBFile("split", "installation", instPath, "Tool installer registry database"))

	startupPath := ResolveSplitDbPath(SectionStartup, "default", "")
	add(inspectSplitDBFile("split", "startup", startupPath, "Startup manager and execution log database"))

	sitesPath := ResolveSplitDbPath(SectionSites, "default", "")
	if _, err := os.Stat(sitesPath); err == nil {
		add(inspectSplitDBFile("split", "sites", sitesPath, "Sites & vhosts configuration database"))
	}

	tasksPath := ResolveTasksRootDbPath("")
	if _, err := os.Stat(tasksPath); err == nil {
		add(inspectSplitDBFile("tasks", "tasks", tasksPath, "Tasks root split database"))
	}
}

func collectSubdirDatabases(baseDir, subName, dbType string, add func(SplitDatabaseEntry)) {
	dir := filepath.Join(baseDir, subName)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, e := range entries {
		processSubdirEntry(dir, subName, dbType, e, add)
	}
}

func processSubdirEntry(dir, subName, dbType string, e os.DirEntry, add func(SplitDatabaseEntry)) {
	if e.IsDir() {
		processSubdirDirEntry(dir, subName, dbType, e.Name(), add)
		return
	}

	processSubdirFileEntry(dir, subName, dbType, e.Name(), add)
}

func processSubdirDirEntry(dir, subName, dbType, name string, add func(SplitDatabaseEntry)) {
	sqlPath := filepath.Join(dir, name, DbFileName)
	if isFileExisting(sqlPath) {
		desc := "Isolated child split database (" + subName + ")"
		add(inspectSplitDBFile(dbType, name, sqlPath, desc))
	}
}

func processSubdirFileEntry(dir, subName, dbType, name string, add func(SplitDatabaseEntry)) {
	if strings.HasSuffix(name, ".db") {
		slug := strings.TrimSuffix(name, ".db")
		path := filepath.Join(dir, name)
		desc := "Isolated child split database (" + subName + ")"
		add(inspectSplitDBFile(dbType, slug, path, desc))
	}
}

func collectLooseDatabases(dataDir string, add func(SplitDatabaseEntry)) {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".db") {
			continue
		}

		path := filepath.Join(dataDir, e.Name())
		slug := strings.TrimSuffix(e.Name(), ".db")
		add(inspectSplitDBFile("custom", slug, path, "Gitmap custom SQLite database"))
	}
}

func collectLocalRepoDatabases(add func(SplitDatabaseEntry)) {
	localDB := filepath.Join(".gitmap", "gitmap.db")
	if _, err := os.Stat(localDB); err == nil {
		add(inspectSplitDBFile("repo", "local-repo", localDB, "Local repository Gitmap SQLite database"))
	}
	collectSubdirDatabases(".", "repodb", "repodb", add)
	collectSubdirDatabases(".gitmap", "repodb", "repodb", add)
	collectSubdirDatabases(".gitmap", "pipeline", "pipeline", add)
	CollectDataTreeDatabases(".gitmap/data", add)
	CollectDataTreeDatabases("data", add)
}

func collectParentSubdirDatabases(baseDir, subName, dbType string, add func(SplitDatabaseEntry)) {
	parent := filepath.Dir(baseDir)
	if parent == "" || parent == "." || parent == baseDir {
		return
	}

	collectSubdirDatabases(parent, subName, dbType, add)
}

func collectUserHomeDatabases(add func(SplitDatabaseEntry)) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	gitmapDir := filepath.Join(home, ".gitmap")
	collectSubdirDatabases(gitmapDir, "pipeline", "pipeline", add)
	collectSubdirDatabases(gitmapDir, "repodb", "repodb", add)
	collectSubdirDatabases(filepath.Join(gitmapDir, "data"), "pipeline", "pipeline", add)
	collectSubdirDatabases(filepath.Join(gitmapDir, "data"), "automation", "automation", add)
}
