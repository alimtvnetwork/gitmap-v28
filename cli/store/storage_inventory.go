package store

import (
	"os"
	"path/filepath"
	"strings"
)

// CollectAllDatabaseEntries returns an inventory of all SQLite databases managed by Gitmap.
func CollectAllDatabaseEntries() []SplitDatabaseEntry {
	dataDir := BinaryDataDir()
	seen := make(map[string]bool)
	var list []SplitDatabaseEntry

	addEntry := func(e SplitDatabaseEntry) {
		clean := filepath.Clean(e.DatabasePath)
		if !seen[clean] {
			seen[clean] = true
			list = append(list, e)
		}
	}

	collectCoreDatabases(dataDir, addEntry)
	collectSubdirDatabases(dataDir, "schedules", "schedule", addEntry)
	collectSubdirDatabases(dataDir, "pipeline", "pipeline", addEntry)
	collectSubdirDatabases(dataDir, "pipeline_db", "pipeline", addEntry)
	collectSubdirDatabases(dataDir, "repo_search", "repo", addEntry)
	collectSubdirDatabases(dataDir, "repodb", "repodb", addEntry)
	collectParentSubdirDatabases(dataDir, "pipeline", "pipeline", addEntry)
	collectParentSubdirDatabases(dataDir, "repodb", "repodb", addEntry)
	collectUserHomeDatabases(addEntry)
	collectLocalRepoDatabases(addEntry)
	collectLooseDatabases(dataDir, addEntry)

	return list
}

func collectCoreDatabases(dataDir string, add func(SplitDatabaseEntry)) {
	rootPath := DefaultDBPath()
	add(inspectSplitDBFile("master", "root", rootPath, "Primary Gitmap repository tracking database"))

	instPath := filepath.Join(dataDir, "installation.db")
	add(inspectSplitDBFile("split", "installation", instPath, "Tool installer registry database"))

	startupPath := filepath.Join(dataDir, "startup.db")
	add(inspectSplitDBFile("split", "startup", startupPath, "Startup manager and execution log database"))

	sitesPath := filepath.Join(dataDir, "sites.db")
	if _, err := os.Stat(sitesPath); err == nil {
		add(inspectSplitDBFile("split", "sites", sitesPath, "Sites & vhosts configuration database"))
	}
}

func collectSubdirDatabases(baseDir, subName, dbType string, add func(SplitDatabaseEntry)) {
	dir := filepath.Join(baseDir, subName)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".db") {
			continue
		}

		slug := strings.TrimSuffix(e.Name(), ".db")
		path := filepath.Join(dir, e.Name())
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
}
