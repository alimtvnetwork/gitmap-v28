package osclean

import (
	"os"
	"path/filepath"
	"strings"
)

func getEnhancedCategoryDefs() []EnhancedCategoryDef {
	return []EnhancedCategoryDef{
		{"go-cache", "Go build & module caches", []string{"go", "gobuild", "gocache"}, cleanGoEnhanced},
		{"node-cache", "Node / npm / pnpm / yarn / bun", []string{"node", "npm", "pnpm", "yarn", "bun"}, cleanNodeEnhanced},
		{"python-cache", "Python / pip / uv / poetry", []string{"python", "pip", "uv", "poetry"}, cleanPythonEnhanced},
		{"bundler-cache", "Vite / webpack / turbopack", []string{"vite", "webpack", "turbopack", "parcel"}, cleanWebBundlerEnhanced},
		{"antigravity-logs", "Antigravity logs & transcripts", []string{"agy", "antigravity", "gemini"}, cleanAntigravityEnhanced},
		{"vscode-cache", "VS Code server & workspace caches", []string{"vscode", "code"}, cleanVSCodeEnhanced},
		{"chrome-cache", "Chrome / Chromium dev caches", []string{"chrome", "chromium"}, cleanChromeDevEnhanced},
		{"git-cache", "Git dangling objects & caches", []string{"git", "gitcache"}, cleanGitEnhanced},
		{"temp-artifacts", "System temp dev artifacts", []string{"temp", "tmp"}, cleanSystemTempEnhanced},
		{"test-binaries", "Stale test binaries (*.test)", []string{"test", "tests", "testbin"}, cleanStaleTestBinaries},
	}
}

func cleanGoEnhanced(isDryRun bool) CategoryCleanStats {
	return cleanGoCache(isDryRun)
}

func cleanNodeEnhanced(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{Category: "node-cache", Label: "Node / npm / pnpm / yarn / bun"}
	mergeCleanStats(&res, cleanPnpmCache(isDryRun))
	mergeCleanStats(&res, cleanNpmCache(isDryRun))
	mergeCleanStats(&res, cleanYarnCache(isDryRun))
	mergeCleanStats(&res, cleanBunCache(isDryRun))
	return res
}

func cleanPythonEnhanced(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{Category: "python-cache", Label: "Python / pip / uv / poetry"}
	mergeCleanStats(&res, cleanPipCache(isDryRun))
	home := resolveHomeDir()
	local := resolveLocalAppData()
	paths := []string{
		filepath.Join(local, "uv", "cache"),
		filepath.Join(home, ".cache", "uv"),
		filepath.Join(local, "pypoetry", "Cache"),
		filepath.Join(home, ".cache", "pypoetry"),
	}
	for _, p := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(p, isDryRun, false))
	}
	return res
}

func cleanWebBundlerEnhanced(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{Category: "bundler-cache", Label: "Vite / webpack / turbopack"}
	home := resolveHomeDir()
	local := resolveLocalAppData()
	paths := []string{
		filepath.Join(local, "vite"), filepath.Join(home, ".cache", "vite"),
		filepath.Join(local, "webpack"), filepath.Join(home, ".cache", "webpack"),
		filepath.Join(local, "turbopack"), filepath.Join(home, ".cache", "turbopack"),
		filepath.Join(local, "parcel"), filepath.Join(home, ".cache", "parcel"),
	}
	for _, p := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(p, isDryRun, false))
	}
	return res
}

func cleanAntigravityEnhanced(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{Category: "antigravity-logs", Label: "Antigravity logs & transcripts"}
	home := resolveHomeDir()
	local := resolveLocalAppData()
	paths := []string{
		filepath.Join(home, ".gemini", "antigravity", "logs"),
		filepath.Join(home, ".gemini", "antigravity", "crashes"),
		filepath.Join(home, ".gemini", "antigravity", "scratch"),
		filepath.Join(local, "antigravity", "Cache"),
		filepath.Join(local, "antigravity", "Code Cache"),
	}
	for _, p := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(p, isDryRun, false))
	}
	return res
}

func cleanVSCodeEnhanced(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{Category: "vscode-cache", Label: "VS Code server & workspace caches"}
	home := resolveHomeDir()
	appData := os.Getenv("APPDATA")
	paths := []string{
		filepath.Join(appData, "Code", "CachedData"),
		filepath.Join(appData, "Code", "CachedExtensionVSIXs"),
		filepath.Join(appData, "Code", "logs"),
		filepath.Join(home, ".cache", "vscode"),
	}
	for _, p := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(p, isDryRun, false))
	}
	return res
}

func cleanChromeDevEnhanced(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{Category: "chrome-cache", Label: "Chrome / Chromium dev caches"}
	local := resolveLocalAppData()
	home := resolveHomeDir()
	paths := []string{
		filepath.Join(local, "Google", "Chrome", "User Data", "Default", "Cache"),
		filepath.Join(local, "Google", "Chrome", "User Data", "Default", "Code Cache"),
		filepath.Join(local, "Google", "Chrome", "User Data", "Default", "GPUCache"),
		filepath.Join(home, ".cache", "google-chrome"),
	}
	for _, p := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(p, isDryRun, false))
	}
	return res
}

func cleanGitEnhanced(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{Category: "git-cache", Label: "Git dangling objects & caches"}
	local := resolveLocalAppData()
	paths := []string{
		filepath.Join(local, "Git", "cache"),
	}
	for _, p := range filterUniquePaths(paths) {
		mergeCleanStats(&res, SweepTarget(p, isDryRun, false))
	}
	return res
}

func cleanSystemTempEnhanced(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{Category: "temp-artifacts", Label: "System temp dev artifacts"}
	tempDir := os.TempDir()
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return res
	}
	for _, e := range entries {
		if isDevTempArtifact(e.Name()) {
			full := filepath.Join(tempDir, e.Name())
			mergeCleanStats(&res, SweepTarget(full, isDryRun, false))
		}
	}
	return res
}

func isDevTempArtifact(name string) bool {
	low := strings.ToLower(name)
	prefixes := []string{"chocolatey", "npm-", "yarn-", "pip-", "go-build", "cargo-", "tmp-"}
	for _, pre := range prefixes {
		if strings.HasPrefix(low, pre) {
			return true
		}
	}
	return false
}

func cleanStaleTestBinaries(isDryRun bool) CategoryCleanStats {
	res := CategoryCleanStats{Category: "test-binaries", Label: "Stale test binaries (*.test)"}
	tempDir := os.TempDir()
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return res
	}
	for _, e := range entries {
		if isTestBinary(e.Name()) {
			full := filepath.Join(tempDir, e.Name())
			mergeCleanStats(&res, SweepTarget(full, isDryRun, false))
		}
	}
	return res
}

func isTestBinary(name string) bool {
	low := strings.ToLower(name)
	return strings.HasSuffix(low, ".test") ||
		strings.HasSuffix(low, ".test.exe") ||
		strings.HasPrefix(low, "__debug_bin")
}
