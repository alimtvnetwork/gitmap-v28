package cmdignore

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// LoadGroups reads saved ignore groups from the split DB.
func LoadGroups() (map[string]IgnoreGroup, *apperror.AppError) {
	db, err := store.OpenIgnoreSplitDB()
	if err != nil {
		return defaultGroupsMap(), nil
	}
	defer db.Close()
	return loadGroupsFromDB(db)
}

func loadGroupsFromDB(db *store.IgnoreSplitDB) (map[string]IgnoreGroup, *apperror.AppError) {
	records, err := db.LoadGroups()
	if err != nil || len(records) == 0 {
		return populateDefaultGroups(db)
	}
	return convertRecordsToGroups(records), nil
}

func convertRecordsToGroups(records map[string]store.IgnoreGroupRecord) map[string]IgnoreGroup {
	groups := make(map[string]IgnoreGroup)
	for k, v := range records {
		groups[k] = IgnoreGroup{Name: v.Name, IsDefault: v.IsDefault, Patterns: v.Patterns}
	}
	return groups
}

func populateDefaultGroups(db *store.IgnoreSplitDB) (map[string]IgnoreGroup, *apperror.AppError) {
	defaults := defaultGroupsMap()
	for _, g := range defaults {
		rec := store.IgnoreGroupRecord{Name: g.Name, IsDefault: g.IsDefault, Patterns: g.Patterns}
		_ = db.SaveGroup(rec)
	}
	return defaults, nil
}

// SaveGroups writes the ignore groups map to the split DB.
func SaveGroups(groups map[string]IgnoreGroup) *apperror.AppError {
	db, err := store.OpenIgnoreSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return saveGroupsToDB(db, groups)
}

func saveGroupsToDB(db *store.IgnoreSplitDB, groups map[string]IgnoreGroup) *apperror.AppError {
	_ = db.ClearGroups()
	for _, g := range groups {
		rec := store.IgnoreGroupRecord{Name: g.Name, IsDefault: g.IsDefault, Patterns: g.Patterns}
		if dbErr := db.SaveGroup(rec); dbErr != nil {
			return dbErr
		}
	}
	return nil
}

func defaultGroupsMap() map[string]IgnoreGroup {
	return map[string]IgnoreGroup{
		"default": {
			Name:      "default",
			Patterns:  gitignoreagm.DefaultGitmapIgnoreEntries,
			IsDefault: true,
		},
		"node": {
			Name:     "node",
			Patterns: []string{"node_modules/", "npm-debug.log*", "yarn-debug.log*", "yarn-error.log*", ".pnpm-debug.log*"},
		},
		"go": {
			Name:     "go",
			Patterns: []string{"bin/", "*.exe", "*.dll", "*.so", "*.dylib", "*.test", "*.out"},
		},
		"python": {
			Name:     "python",
			Patterns: []string{"__pycache__/", "*.py[cod]", "*$py.class", ".venv/", "venv/", "env/"},
		},
	}
}

// AddGroup creates or updates a named ignore group.
func AddGroup(name string, patterns []string) *apperror.AppError {
	groups, _ := LoadGroups()
	cleanName := strings.ToLower(strings.TrimSpace(name))
	groups[cleanName] = IgnoreGroup{Name: cleanName, Patterns: patterns}
	return SaveGroups(groups)
}

// RemoveGroup deletes a named ignore group.
func RemoveGroup(name string) *apperror.AppError {
	groups, _ := LoadGroups()
	cleanName := strings.ToLower(strings.TrimSpace(name))
	delete(groups, cleanName)
	return SaveGroups(groups)
}

// SetDefaultGroup sets the specified group as default.
func SetDefaultGroup(name string) *apperror.AppError {
	groups, _ := LoadGroups()
	cleanName := strings.ToLower(strings.TrimSpace(name))
	for k, g := range groups {
		g.IsDefault = (k == cleanName)
		groups[k] = g
	}
	return SaveGroups(groups)
}

// AddGroupToDefault marks the specified group as also default.
func AddGroupToDefault(name string) *apperror.AppError {
	groups, _ := LoadGroups()
	cleanName := strings.ToLower(strings.TrimSpace(name))
	if g, ok := groups[cleanName]; ok {
		g.IsDefault = true
		groups[cleanName] = g
		return SaveGroups(groups)
	}
	return apperror.NewSimple("group not found: "+name, "E1001")
}

// ApplyGroupToRepo appends the group patterns to a repository's .gitignore.
func ApplyGroupToRepo(groupName, repoDir string) *apperror.AppError {
	groups, _ := LoadGroups()
	cleanName := strings.ToLower(strings.TrimSpace(groupName))
	grp, ok := groups[cleanName]
	if !ok {
		return apperror.NewSimple("group not found: "+groupName, "E1002")
	}
	return applyGroupPatterns(grp, repoDir)
}

func applyGroupPatterns(grp IgnoreGroup, repoDir string) *apperror.AppError {
	ignorePath := filepath.Join(repoDir, ".gitignore")
	data, _ := os.ReadFile(ignorePath)
	cleaned, isModified := gitignoreagm.DeduplicateAndSanitizeGitignore(string(data), grp.Patterns...)
	if isModified {
		if err := os.WriteFile(ignorePath, []byte(cleaned), 0644); err != nil {
			return apperror.WrapSimple(err, "write .gitignore")
		}
	}
	return nil
}
