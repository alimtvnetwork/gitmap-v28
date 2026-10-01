package cmdignore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func resolveGroupsConfigPath() string {
	baseDir := store.BinaryDataDir()
	dir := filepath.Join(baseDir, "ignore")
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "groups.json")
}

// LoadGroups reads saved ignore groups from disk.
func LoadGroups() (map[string]IgnoreGroup, error) {
	path := resolveGroupsConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return defaultGroupsMap(), nil
	}
	var groups map[string]IgnoreGroup
	if jsonErr := json.Unmarshal(data, &groups); jsonErr != nil {
		return defaultGroupsMap(), nil
	}
	return groups, nil
}

// SaveGroups writes the ignore groups map to disk.
func SaveGroups(groups map[string]IgnoreGroup) error {
	path := resolveGroupsConfigPath()
	data, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal ignore groups")
	}
	return os.WriteFile(path, data, 0644)
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
func AddGroup(name string, patterns []string) error {
	groups, _ := LoadGroups()
	cleanName := strings.ToLower(strings.TrimSpace(name))
	groups[cleanName] = IgnoreGroup{Name: cleanName, Patterns: patterns}
	return SaveGroups(groups)
}

// RemoveGroup deletes a named ignore group.
func RemoveGroup(name string) error {
	groups, _ := LoadGroups()
	cleanName := strings.ToLower(strings.TrimSpace(name))
	delete(groups, cleanName)
	return SaveGroups(groups)
}

// SetDefaultGroup sets the specified group as default.
func SetDefaultGroup(name string) error {
	groups, _ := LoadGroups()
	cleanName := strings.ToLower(strings.TrimSpace(name))
	for k, g := range groups {
		g.IsDefault = (k == cleanName)
		groups[k] = g
	}
	return SaveGroups(groups)
}

// AddGroupToDefault marks the specified group as also default.
func AddGroupToDefault(name string) error {
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
func ApplyGroupToRepo(groupName, repoDir string) error {
	groups, _ := LoadGroups()
	cleanName := strings.ToLower(strings.TrimSpace(groupName))
	grp, ok := groups[cleanName]
	if !ok {
		return apperror.NewSimple("group not found: "+groupName, "E1002")
	}
	ignorePath := filepath.Join(repoDir, ".gitignore")
	data, _ := os.ReadFile(ignorePath)
	cleaned, isModified := gitignoreagm.DeduplicateAndSanitizeGitignore(string(data), grp.Patterns...)
	if isModified {
		return os.WriteFile(ignorePath, []byte(cleaned), 0644)
	}
	return nil
}
