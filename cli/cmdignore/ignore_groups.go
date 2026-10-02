package cmdignore

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ImmutableDefaultRules defines the mandatory immutable rules for GitMap.
var ImmutableDefaultRules = []string{
	".gitmap/backup/",
}

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
	if err != nil {
		return populateDefaultGroups(db)
	}
	if len(records) == 0 {
		return populateDefaultGroups(db)
	}
	return convertRecordsToGroups(records), nil
}

func convertRecordsToGroups(records map[string]store.IgnoreGroupRecord) map[string]IgnoreGroup {
	groups := make(map[string]IgnoreGroup)
	for k, v := range records {
		patterns := v.Patterns
		if v.IsDefault {
			patterns = EnsureImmutableRules(patterns)
		}
		groups[k] = IgnoreGroup{Name: v.Name, IsDefault: v.IsDefault, Patterns: patterns}
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
		patterns := sanitizeGroupPatterns(g)
		rec := store.IgnoreGroupRecord{Name: g.Name, IsDefault: g.IsDefault, Patterns: patterns}
		if dbErr := db.SaveGroup(rec); dbErr != nil {
			return dbErr
		}
	}
	return nil
}

func sanitizeGroupPatterns(g IgnoreGroup) []string {
	if g.IsDefault || g.Name == "default" {
		return EnsureImmutableRules(g.Patterns)
	}
	return g.Patterns
}

func defaultGroupsMap() map[string]IgnoreGroup {
	defaultPatterns := EnsureImmutableRules(gitignoreagm.DefaultGitmapIgnoreEntries)
	groups := techStackGroups()
	groups["default"] = IgnoreGroup{
		Name:      "default",
		Patterns:  defaultPatterns,
		IsDefault: true,
	}
	return groups
}

func techStackGroups() map[string]IgnoreGroup {
	return map[string]IgnoreGroup{
		"node":   {Name: "node", Patterns: []string{"node_modules/", "npm-debug.log*", "yarn-debug.log*", "yarn-error.log*", ".pnpm-debug.log*"}},
		"go":     {Name: "go", Patterns: []string{"bin/", "*.exe", "*.dll", "*.so", "*.dylib", "*.test", "*.out"}},
		"python": {Name: "python", Patterns: []string{"__pycache__/", "*.py[cod]", "*$py.class", ".venv/", "venv/", "env/"}},
	}
}

// EnsureImmutableRules ensures .gitmap/backup/ is present in patterns.
func EnsureImmutableRules(patterns []string) []string {
	hasBackup := hasExactPattern(patterns, ".gitmap/backup/")
	if hasBackup {
		return patterns
	}
	return append(patterns, ".gitmap/backup/")
}

func hasExactPattern(patterns []string, target string) bool {
	cleanTarget := strings.TrimRight(strings.TrimSpace(target), "/")
	for _, p := range patterns {
		clean := strings.TrimRight(strings.TrimSpace(p), "/")
		if clean == cleanTarget {
			return true
		}
	}
	return false
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
	cleanName := strings.ToLower(strings.TrimSpace(name))
	if cleanName == "default" {
		return apperror.NewSimple("cannot remove immutable default ignore group", "E1014")
	}
	groups, _ := LoadGroups()
	delete(groups, cleanName)
	saveErr := SaveGroups(groups)
	if saveErr != nil {
		return saveErr
	}
	return cleanupGroupBindings(cleanName)
}

func cleanupGroupBindings(groupName string) *apperror.AppError {
	db, err := store.OpenIgnoreSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return db.DeleteBindingsForGroup(groupName)
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
	g, hasGroup := groups[cleanName]
	if !hasGroup {
		return apperror.NewSimple("group not found: "+name, "E1001")
	}
	g.IsDefault = true
	groups[cleanName] = g
	return SaveGroups(groups)
}

// AddPatternsToDefault appends new patterns to the default ignore group.
func AddPatternsToDefault(patterns []string) *apperror.AppError {
	groups, _ := LoadGroups()
	defGrp, hasDef := groups["default"]
	if !hasDef {
		defGrp = IgnoreGroup{Name: "default", IsDefault: true, Patterns: EnsureImmutableRules(nil)}
	}
	defGrp.Patterns = appendUniquePatterns(defGrp.Patterns, patterns)
	defGrp.Patterns = EnsureImmutableRules(defGrp.Patterns)
	groups["default"] = defGrp
	return SaveGroups(groups)
}

func appendUniquePatterns(existing, newPatterns []string) []string {
	seen := buildPatternSet(existing)
	res := existing
	for _, p := range newPatterns {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		seen[trimmed] = true
		res = append(res, trimmed)
	}
	return res
}

func buildPatternSet(patterns []string) map[string]bool {
	seen := make(map[string]bool, len(patterns))
	for _, p := range patterns {
		seen[strings.TrimSpace(p)] = true
	}
	return seen
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
	patterns := resolveGroupPatterns(grp)
	ignorePath := filepath.Join(repoDir, ".gitignore")
	data, _ := os.ReadFile(ignorePath)
	cleaned, isModified := gitignoreagm.DeduplicateAndSanitizeGitignore(string(data), patterns...)
	if !isModified {
		return nil
	}
	if err := os.WriteFile(ignorePath, []byte(cleaned), 0644); err != nil {
		return apperror.WrapSimple(err, "write .gitignore")
	}
	return nil
}

func resolveGroupPatterns(grp IgnoreGroup) []string {
	if grp.IsDefault || grp.Name == "default" {
		return EnsureImmutableRules(grp.Patterns)
	}
	return grp.Patterns
}

// ConnectGroupWithRepo links an ignore group to a repository path.
func ConnectGroupWithRepo(groupName, repoPath string, hasAddWithDefault bool) *apperror.AppError {
	cleanName := strings.ToLower(strings.TrimSpace(groupName))
	if err := validateGroupExists(cleanName); err != nil {
		return err
	}
	if err := saveRepoBindingRecord(repoPath, cleanName); err != nil {
		return err
	}
	_ = ApplyGroupToRepo(cleanName, repoPath)
	applyDefaultBindingIfNeeded(repoPath, hasAddWithDefault)
	return nil
}

func validateGroupExists(groupName string) *apperror.AppError {
	groups, _ := LoadGroups()
	if _, hasGroup := groups[groupName]; !hasGroup {
		return apperror.NewSimple("group not found: "+groupName, "E1002")
	}
	return nil
}

func applyDefaultBindingIfNeeded(repoPath string, hasAddWithDefault bool) {
	if hasAddWithDefault {
		_ = saveRepoBindingRecord(repoPath, "default")
		_ = ApplyGroupToRepo("default", repoPath)
	}
}

func saveRepoBindingRecord(repoPath, groupName string) *apperror.AppError {
	db, err := store.OpenIgnoreSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return db.SaveRepoBinding(repoPath, groupName)
}

// ApplyAllToRepo applies default patterns and bound group patterns to a repo path.
func ApplyAllToRepo(repoPath string) *apperror.AppError {
	_ = ApplyGroupToRepo("default", repoPath)
	bindings, err := LoadBindingsForRepo(repoPath)
	if err != nil {
		return nil
	}
	for _, grpName := range bindings {
		_ = ApplyGroupToRepo(grpName, repoPath)
	}
	return nil
}

// LoadBindingsForRepo returns groups bound to a specific repository.
func LoadBindingsForRepo(repoPath string) ([]string, *apperror.AppError) {
	db, err := store.OpenIgnoreSplitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return db.LoadBindingsForRepo(repoPath)
}

// LoadRepoBindings returns all configured repository to group bindings.
func LoadRepoBindings() ([]store.IgnoreRepoBindingRecord, *apperror.AppError) {
	db, err := store.OpenIgnoreSplitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return db.LoadRepoBindings()
}
