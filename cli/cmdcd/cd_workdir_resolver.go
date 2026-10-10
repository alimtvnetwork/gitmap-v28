// Package cmd — cd_workdir_resolver.go: resolves work directories for cd navigation.
package cmdcd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsetup"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func handleWorkDirOrNotFound(name string, rest []string) error {
	workPath, hasWorkDir := resolveCDWorkDirPath(name)
	if hasWorkDir {
		return dispatchCDWorkPath(workPath, rest)
	}

	suggestions := suggestCDRepos(name)
	msg := formatCDNotFoundMessage(name, suggestions)

	return apperror.NewNotFoundError(msg)
}

func dispatchCDWorkPath(workPath string, rest []string) error {
	if len(rest) > 0 {
		return runCDInner(workPath, rest)
	}

	fmt.Print(workPath)
	WriteShellHandoff(workPath)
	cmdsetup.WarnIfNoWrapper()

	return nil
}

func resolveDefaultWorkDirPath() (string, bool) {
	if path, ok := resolveConfiguredDefaultWorkDir(); ok {
		return path, true
	}

	if path, ok := resolveEnvWorkDir(); ok {
		persistDefaultWorkDirIfPossible(path)
		return path, true
	}

	if path, ok := resolveInferredWorkDirFromRepos(); ok {
		persistDefaultWorkDirIfPossible(path)
		return path, true
	}

	if path, ok := resolveWellKnownWorkDir(); ok {
		persistDefaultWorkDirIfPossible(path)
		return path, true
	}

	if path, ok := resolveCwdWorkDir(); ok {
		persistDefaultWorkDirIfPossible(path)
		return path, true
	}

	return "", false
}

func resolveConfiguredDefaultWorkDir() (string, bool) {
	db, err := store.OpenDefault()
	if err != nil {
		return "", false
	}
	defer db.Close()

	wd, errGet := db.GetDefaultWorkDir()
	if errGet != nil || wd == nil || wd.AbsolutePath == "" {
		return "", false
	}

	info, errStat := os.Stat(wd.AbsolutePath)
	if errStat != nil || !info.IsDir() {
		return "", false
	}

	return wd.AbsolutePath, true
}

func resolveEnvWorkDir() (string, bool) {
	envVars := []string{
		"GITMAP_WORK_DIR",
		"GITMAP_WORK",
		"WORK_DIR",
		"WORK",
	}
	for _, env := range envVars {
		val := strings.TrimSpace(os.Getenv(env))
		if val == "" {
			continue
		}
		if info, err := os.Stat(val); err == nil && info.IsDir() {
			return val, true
		}
	}
	return "", false
}

func resolveInferredWorkDirFromRepos() (string, bool) {
	db, err := store.OpenDefault()
	if err != nil {
		return "", false
	}
	defer db.Close()

	repos, errList := db.ListRepos()
	if errList != nil || len(repos) == 0 {
		return "", false
	}

	counts := make(map[string]int)
	for _, r := range repos {
		if r.AbsolutePath == "" {
			continue
		}
		parent := filepath.Dir(r.AbsolutePath)
		if info, err := os.Stat(parent); err == nil && info.IsDir() {
			counts[parent]++
		}
	}

	topDir := ""
	topCount := 0
	for dir, count := range counts {
		if count > topCount {
			topCount = count
			topDir = dir
		}
	}

	if topDir != "" {
		return topDir, true
	}
	return "", false
}

func resolveWellKnownWorkDir() (string, bool) {
	candidates := []string{
		`D:\work`,
		`C:\work`,
		`D:\git-work`,
		`C:\git-work`,
		`/work`,
		`/git-work`,
	}

	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		candidates = append(candidates,
			filepath.Join(userProfile, "work"),
			filepath.Join(userProfile, "git-work"),
			filepath.Join(userProfile, "projects"),
			filepath.Join(userProfile, "source", "repos"),
		)
	}

	if home := os.Getenv("HOME"); home != "" {
		candidates = append(candidates,
			filepath.Join(home, "work"),
			filepath.Join(home, "git-work"),
			filepath.Join(home, "projects"),
		)
	}

	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c, true
		}
	}

	return "", false
}

func resolveCwdWorkDir() (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}

	base := strings.ToLower(filepath.Base(cwd))
	if base == "work" || base == "git-work" {
		return cwd, true
	}

	parent := filepath.Dir(cwd)
	parentBase := strings.ToLower(filepath.Base(parent))
	if parentBase == "work" || parentBase == "git-work" {
		if info, err := os.Stat(parent); err == nil && info.IsDir() {
			return parent, true
		}
	}

	return "", false
}

func persistDefaultWorkDirIfPossible(path string) {
	db, err := store.OpenDefault()
	if err != nil {
		return
	}
	defer db.Close()

	if errEnsure := db.EnsureWorkDirsTable(); errEnsure == nil {
		label := strings.ToLower(filepath.Base(path))
		if wd, errEnsureWd := db.EnsureWorkDir(path, label, true); errEnsureWd == nil && wd != nil {
			_ = db.SetDefaultWorkDir(wd.AbsolutePath)
		}
	}
}

func resolveCDWorkDirPath(name string) (string, bool) {
	lower := strings.ToLower(name)
	if isWorkDirKeyword(lower) {
		return resolveDefaultWorkDirPath()
	}

	return findWorkDirByNameOrLabel(name)
}

func isWorkDirKeyword(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "work", "$work", `\$work`, "%work%",
		"def", "$def", `\$def`, "%def%",
		"workdir", "$workdir", `\$workdir`, "%workdir%",
		"default", "$default", `\$default`, "%default%",
		"wd", "$wd", `\$wd`, "%wd%":
		return true
	default:
		return false
	}
}

func findWorkDirByNameOrLabel(target string) (string, bool) {
	db, err := store.OpenDefault()
	if err != nil {
		return "", false
	}

	defer db.Close()

	dirs, errList := db.ListWorkDirs()
	if errList != nil || len(dirs) == 0 {
		return "", false
	}

	for _, d := range dirs {
		if matchesWorkDir(d.AbsolutePath, d.Label, target) {
			return d.AbsolutePath, true
		}
	}

	return "", false
}

func matchesWorkDir(absPath, label, target string) bool {
	lowerTarget := strings.ToLower(target)
	if label != "" && strings.EqualFold(label, lowerTarget) {
		return true
	}

	baseName := strings.ToLower(filepath.Base(absPath))
	if baseName == lowerTarget {
		return true
	}

	return strings.Contains(strings.ToLower(absPath), lowerTarget)
}
