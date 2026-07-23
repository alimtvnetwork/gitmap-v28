package cmd

// Config loader + glob matcher for `gitmap fix-repo`. Mirrors
// scripts/fix-repo/Config-Loader.ps1: load JSON, expose ignoreDirs +
// ignorePatterns, and report Test-IsIgnoredPath equivalents.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// fixRepoConfig is the on-disk JSON shape (matches fix-repo.config.json).
type fixRepoConfig struct {
	IgnoreDirs     []string `json:"ignoreDirs"`
	IgnorePatterns []string `json:"ignorePatterns"`
}

// fixRepoIgnore is the resolved in-memory matcher built from the
// JSON config. Held in package state so the rewrite sweep doesn't
// need to thread it through every call.
type fixRepoIgnore struct {
	dirs     []string
	patterns []*regexp.Regexp
}

var fixRepoActiveIgnore fixRepoIgnore

// loadFixRepoConfig resolves the config path (explicit > default >
// missing-is-ok) and populates fixRepoActiveIgnore. An explicit
// --config that doesn't exist or fails to parse exits E_BAD_CONFIG.
func loadFixRepoConfig(explicit, repoRoot string) {
	resolved, err := resolveFixRepoConfigPath(explicit, repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.FixRepoErrBadConfigFmt, err.Error())
		os.Exit(constants.FixRepoExitBadConfig)
	}
	fixRepoActiveIgnore = fixRepoIgnore{}
	if resolved == "" {
		return
	}
	cfg, err := readFixRepoConfig(resolved)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.FixRepoErrBadConfigFmt, err.Error())
		os.Exit(constants.FixRepoExitBadConfig)
	}
	fixRepoActiveIgnore = compileFixRepoIgnore(cfg)
}

// resolveFixRepoConfigPath returns the file path to load, "" when no
// config is present, or an error when an explicit path is missing.
func resolveFixRepoConfigPath(explicit, repoRoot string) (string, error) {
	if explicit != "" {
		info, err := os.Stat(explicit)
		if err != nil || info.IsDir() {
			return "", fmt.Errorf("config file not found: %s", explicit)
		}

		return explicit, nil
	}
	def := filepath.Join(repoRoot, constants.FixRepoConfigFileName)
	info, err := os.Stat(def)
	if err != nil || info.IsDir() {
		return "", nil
	}

	return def, nil
}

// readFixRepoConfig reads + decodes the JSON config file.
func readFixRepoConfig(path string) (fixRepoConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fixRepoConfig{}, err
	}
	var cfg fixRepoConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return fixRepoConfig{}, err
	}

	return cfg, nil
}

// compileFixRepoIgnore turns the JSON config into an executable matcher.
func compileFixRepoIgnore(cfg fixRepoConfig) fixRepoIgnore {
	out := fixRepoIgnore{}
	for _, d := range cfg.IgnoreDirs {
		if d != "" {
			out.dirs = append(out.dirs, strings.Trim(d, "/\\"))
		}
	}
	for _, p := range cfg.IgnorePatterns {
		if p == "" {
			continue
		}
		re, err := regexp.Compile(globToRegex(p))
		if err == nil {
			out.patterns = append(out.patterns, re)
		}
	}

	return out
}

// isFixRepoIgnoredPath returns true when relPath should be skipped
// per the active ignore config. Paths are normalized to forward
// slashes so config patterns are platform-portable.
func isFixRepoIgnoredPath(relPath string) bool {
	norm := strings.ReplaceAll(relPath, "\\", "/")
	for _, d := range fixRepoActiveIgnore.dirs {
		if norm == d || strings.HasPrefix(norm, d+"/") {
			return true
		}
	}
	for _, re := range fixRepoActiveIgnore.patterns {
		if re.MatchString(norm) {
			return true
		}
	}

	return false
}

// globToRegex converts a glob pattern to an anchored regex string.
// `**` matches any depth, `*` matches within one segment, `?`
// matches one non-`/` char. Mirrors ConvertTo-FixRepoRegex.
func globToRegex(pattern string) string {
	var b strings.Builder
	b.WriteByte('^')
	for i := 0; i < len(pattern); i++ {
		i = appendGlobChar(&b, pattern, i)
	}
	b.WriteByte('$')

	return b.String()
}

// appendGlobChar appends one translated token to b and returns the
// new (last-consumed) index. Split from globToRegex to honor the
// 15-line function-length cap.
func appendGlobChar(b *strings.Builder, pattern string, i int) int {
	ch := pattern[i]
	if ch == '*' {
		if i+1 < len(pattern) && pattern[i+1] == '*' {
			b.WriteString(".*")

			return i + 1
		}
		b.WriteString("[^/]*")

		return i
	}
	if ch == '?' {
		b.WriteString("[^/]")

		return i
	}
	if strings.ContainsRune(`.+()[]{}^$|\`, rune(ch)) {
		b.WriteByte('\\')
	}
	b.WriteByte(ch)

	return i
}
