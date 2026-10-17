package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var supportedAINames = []string{
	"antigravity",
	"muse",
	"cloud-code",
	"codex",
	"grok",
	"cursor",
}

var aiSequencePrefixes = map[string]string{
	"antigravity": "01",
	"muse":        "02",
	"cloud-code":  "03",
	"codex":       "04",
	"grok":        "05",
	"cursor":      "06",
}

// ListSupportedAINames returns the slice of canonical supported AI engine slugs.
func ListSupportedAINames() []string {
	result := make([]string, len(supportedAINames))
	copy(result, supportedAINames)
	return result
}

// NormalizeAIName standardizes raw AI input into a canonical slug.
func NormalizeAIName(raw string) string {
	clean := strings.ToLower(strings.TrimSpace(raw))
	clean = strings.TrimLeft(clean, "0123456789-")
	if clean == "" {
		return "antigravity"
	}
	for _, supported := range supportedAINames {
		if clean == supported || strings.HasSuffix(clean, "-"+supported) {
			return supported
		}
	}
	return clean
}

// AISequencePrefix returns the 2-digit sequence prefix for an AI slug.
func AISequencePrefix(aiName string) string {
	norm := NormalizeAIName(aiName)
	if prefix, exists := aiSequencePrefixes[norm]; exists {
		return prefix
	}
	return "01"
}

// DetectCurrentAIName inspects environment variables to identify the ambient AI.
func DetectCurrentAIName() string {
	if checkEnvAny("ANTIGRAVITY_AGENT_ID", "ANTIGRAVITY_INSTANCE_ID", "GEMINI_CLI_VERSION") {
		return "antigravity"
	}
	if checkEnvAny("MUSE_RUNNER_ID", "MUSE_VERSION", "MUSE_MODE") {
		return "muse"
	}
	if checkEnvAny("CLOUD_CODE_PROJECT", "CLOUD_CODE_VERSION", "GOOGLE_CLOUD_AGENT") {
		return "cloud-code"
	}
	if checkEnvAny("CODEX_ENV", "COPILOT_AGENT", "OPENAI_AGENT_SESSION") {
		return "codex"
	}
	if checkEnvAny("GROK_API_KEY", "GROK_RUNNER_ID") {
		return "grok"
	}
	if checkEnvAny("CURSOR_AGENT", "CURSOR_PROJECT_DIR", "CURSOR_SETTINGS") {
		return "cursor"
	}
	return "antigravity"
}

func checkEnvAny(keys ...string) bool {
	for _, k := range keys {
		if val := strings.TrimSpace(os.Getenv(k)); val != "" {
			return true
		}
	}
	return false
}

// ResolveRepoCacheRoot locates the base directory of the repo-cache repository.
func ResolveRepoCacheRoot() string {
	if info, err := os.Stat(`D:\work\repo-cache`); err == nil && info.IsDir() {
		return `D:\work\repo-cache`
	}
	baseWork := resolveWorkBaseDir()
	candidate := filepath.Join(baseWork, "repo-cache")
	if info, err := os.Stat(candidate); err == nil && info.IsDir() {
		return candidate
	}
	if info, err := os.Stat("repo-cache"); err == nil && info.IsDir() {
		return "repo-cache"
	}
	return candidate
}

func resolveWorkBaseDir() string {
	if info, err := os.Stat(`D:\work`); err == nil && info.IsDir() {
		return `D:\work`
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return filepath.Dir(cwd)
}

// ResolveRcHeatmapBaseDir returns <repo-cache>/.heatmap.
func ResolveRcHeatmapBaseDir() string {
	return filepath.ToSlash(filepath.Join(ResolveRepoCacheRoot(), ".heatmap"))
}

// ResolveRcHeatmapAiDir returns <repo-cache>/.heatmap/01-<ai-name>.
func ResolveRcHeatmapAiDir(aiName string) string {
	norm := NormalizeAIName(aiName)
	folder := fmt.Sprintf("%s-%s", AISequencePrefix(norm), norm)
	return filepath.ToSlash(filepath.Join(ResolveRcHeatmapBaseDir(), folder))
}

// ResolveRcHeatmapRootDbPath resolves the path to root-heatmap.db for the given AI.
func ResolveRcHeatmapRootDbPath(aiName string) string {
	aiDir := ResolveRcHeatmapAiDir(aiName)
	_ = os.MkdirAll(aiDir, 0755)
	return filepath.ToSlash(filepath.Join(aiDir, "root-heatmap.db"))
}

// ResolveRcHeatmapRepoDbPath resolves the path to heatmap.db for a specific repo under the given AI.
func ResolveRcHeatmapRepoDbPath(aiName, repoSlug string) string {
	aiDir := ResolveRcHeatmapAiDir(aiName)
	cleanSlug := SanitizeSlug(repoSlug)
	repoDir := filepath.Join(aiDir, cleanSlug)
	_ = os.MkdirAll(repoDir, 0755)
	return filepath.ToSlash(filepath.Join(repoDir, "heatmap.db"))
}

// ResolveRcHeatmapLocalTrackingDbPath resolves the local tracking database path (<dataDir>/heatmap.db).
func ResolveRcHeatmapLocalTrackingDbPath() string {
	dataDir := EffectiveDataDir()
	_ = os.MkdirAll(dataDir, 0755)
	return filepath.ToSlash(filepath.Join(dataDir, "heatmap.db"))
}

// ResolveAmbientRepoSlug discovers the current repo slug from cwd or fallback.
func ResolveAmbientRepoSlug(cwd string) string {
	target := cwd
	if target == "" {
		if dir, err := os.Getwd(); err == nil {
			target = dir
		}
	}
	base := filepath.Base(target)
	return SanitizeSlug(base)
}
