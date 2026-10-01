package cmdos

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func getAllDiscoveryRules() []EcosystemDiscoveryRule {
	return []EcosystemDiscoveryRule{
		{EcosystemID: "go", DisplayName: "Go (Golang)", CliCommands: [][]string{{"go", "env", "GOCACHE"}, {"go", "env", "GOMODCACHE"}}, EnvVariables: []string{"GOCACHE", "GOMODCACHE"}, DevToolSubdirs: []string{"go\\cache", "go\\pkg\\mod"}, HeuristicPaths: []string{"%LOCALAPPDATA%\\go-build"}},
		{EcosystemID: "pnpm", DisplayName: "pnpm Store", CliCommands: [][]string{{"pnpm", "store", "path"}}, EnvVariables: []string{"PNPM_HOME"}, DevToolSubdirs: []string{"pnpm\\store", "pnpm\\store\\v10", "pnpm\\cache"}, HeuristicPaths: []string{"%LOCALAPPDATA%\\pnpm\\store"}},
		{EcosystemID: "npm", DisplayName: "npm Cache", CliCommands: [][]string{{"npm", "config", "get", "cache"}}, EnvVariables: []string{"npm_config_cache"}, DevToolSubdirs: []string{"npm\\cache"}, HeuristicPaths: []string{"%LOCALAPPDATA%\\npm-cache", "~/.npm"}},
		{EcosystemID: "yarn", DisplayName: "Yarn Cache", CliCommands: [][]string{{"yarn", "cache", "dir"}}, EnvVariables: []string{"YARN_CACHE_FOLDER"}, DevToolSubdirs: []string{"yarn\\cache"}, HeuristicPaths: []string{"%LOCALAPPDATA%\\Yarn\\Cache"}},
		{EcosystemID: "bun", DisplayName: "Bun Cache", CliCommands: [][]string{{"bun", "pm", "cache"}}, EnvVariables: []string{"BUN_INSTALL"}, DevToolSubdirs: []string{"bun\\cache"}, HeuristicPaths: []string{"~/.bun/install/cache"}},
		{EcosystemID: "python", DisplayName: "Python Pip/UV", CliCommands: [][]string{{"pip", "cache", "dir"}}, EnvVariables: []string{"PIP_CACHE_DIR", "UV_CACHE_DIR"}, DevToolSubdirs: []string{"python\\pip\\cache"}, HeuristicPaths: []string{"%LOCALAPPDATA%\\pip\\cache", "~/.cache/pip"}},
		{EcosystemID: "cargo", DisplayName: "Rust Cargo", CliCommands: [][]string{{"cargo", "cache", "--dir"}}, EnvVariables: []string{"CARGO_HOME"}, DevToolSubdirs: []string{"cargo\\registry\\cache"}, HeuristicPaths: []string{"~/.cargo/registry/cache"}},
		{EcosystemID: "nuget", DisplayName: ".NET NuGet", CliCommands: nil, EnvVariables: []string{"NUGET_PACKAGES"}, DevToolSubdirs: []string{"nuget\\cache"}, HeuristicPaths: []string{"%LOCALAPPDATA%\\NuGet\\v3-cache"}},
		{EcosystemID: "gradle", DisplayName: "Java Gradle", CliCommands: nil, EnvVariables: []string{"GRADLE_USER_HOME"}, DevToolSubdirs: []string{"gradle\\caches"}, HeuristicPaths: []string{"~/.gradle/caches"}},
		{EcosystemID: "maven", DisplayName: "Java Maven", CliCommands: nil, EnvVariables: []string{"M2_HOME"}, DevToolSubdirs: []string{"maven\\repository"}, HeuristicPaths: []string{"~/.m2/repository"}},
	}
}

func getEcosystemDiscoveryRules(targets []string) []EcosystemDiscoveryRule {
	all := getAllDiscoveryRules()
	if len(targets) == 0 {
		return all
	}
	var filtered []EcosystemDiscoveryRule
	for _, rule := range all {
		if containsEcosystem(targets, rule.EcosystemID) {
			filtered = append(filtered, rule)
		}
	}
	return filtered
}

func containsEcosystem(targets []string, id string) bool {
	for _, t := range targets {
		if strings.EqualFold(strings.TrimSpace(t), id) {
			return true
		}
	}
	return false
}

func probeTier1CliPaths(rule EcosystemDiscoveryRule) []DiscoveredCachePath {
	var items []DiscoveredCachePath
	for _, cmdArgs := range rule.CliCommands {
		if path, ok := executeSafeCliProbe(cmdArgs); ok && len(path) > 0 {
			items = append(items, newDiscoveredPath(path, rule.EcosystemID, "cli", false))
		}
	}
	return items
}

func probeTier2EnvPaths(rule EcosystemDiscoveryRule) []DiscoveredCachePath {
	var items []DiscoveredCachePath
	for _, envKey := range rule.EnvVariables {
		if val := strings.TrimSpace(os.Getenv(envKey)); len(val) > 0 {
			items = append(items, newDiscoveredPath(val, rule.EcosystemID, "env", false))
		}
	}
	return items
}

func probeTier3HeuristicPaths(rule EcosystemDiscoveryRule) []DiscoveredCachePath {
	var items []DiscoveredCachePath
	for _, dir := range probeWindowsDevDrives(rule.DevToolSubdirs) {
		items = append(items, newDiscoveredPath(dir, rule.EcosystemID, "heuristic", true))
	}
	for _, raw := range rule.HeuristicPaths {
		items = append(items, newDiscoveredPath(expandUserPath(raw), rule.EcosystemID, "heuristic", false))
	}
	return items
}

func probeWindowsDevDrives(subpaths []string) []string {
	var discovered []string
	for _, drive := range getAvailableDriveRoots() {
		for _, sub := range subpaths {
			target := filepath.Join(drive, "dev-tool", sub)
			if isExistingDirectory(target) {
				discovered = append(discovered, target)
			}
		}
	}
	return discovered
}

func getAvailableDriveRoots() []string {
	if runtime.GOOS != "windows" {
		return []string{"/"}
	}
	var roots []string
	for _, d := range []string{"C", "D", "E", "F", "G"} {
		root := d + `:\`
		if isExistingDirectory(root) {
			roots = append(roots, root)
		}
	}
	return roots
}

func expandUserPath(raw string) string {
	if strings.HasPrefix(raw, "%LOCALAPPDATA%") {
		local := os.Getenv("LOCALAPPDATA")
		return filepath.Join(local, strings.TrimPrefix(raw, "%LOCALAPPDATA%\\"))
	}
	if strings.HasPrefix(raw, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, strings.TrimPrefix(raw, "~/"))
	}
	return filepath.Clean(os.ExpandEnv(raw))
}

func executeSafeCliProbe(cmdArgs []string) (string, bool) {
	if len(cmdArgs) == 0 {
		return "", false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 750*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	res := strings.TrimSpace(string(out))
	return res, len(res) > 0
}
