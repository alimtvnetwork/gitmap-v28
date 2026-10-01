# Subtask 01: Multi-Ecosystem Dynamic Discovery Engine

> **Parent Plan:** [62-devtools-cache-discovery-tree-and-split-db](../../pending/62-devtools-cache-discovery-tree-and-split-db.md)  
> **Tracking Spec:** [02-spec/21-app/199-devtools-cache-dynamic-discovery-tree-view-and-split-db.md](../../../../02-spec/21-app/199-devtools-cache-dynamic-discovery-tree-view-and-split-db.md)  
> **Primary File Targets:** `cli/cmdos/os_dev_clean_paths.go`, `cli/cmdos/os_dev_clean_discover.go`  
> **Status:** `PENDING`  
> **Owner:** Worker 01  

---

## 1. Objective & Problem Remediation

### 1.1 Problem Statement
When running `gitmap clear devtools` (or `gitmap clean devtools`, `gitmap devtool clear`), the tool suffered from two fatal architectural defects:
1. **The 0.01 MB Metric Anomaly:** Purge commands (`go clean`, `pnpm store prune`, `npm cache clean`) were triggered *before* directory sweeps counted sizes. By the time directories were traversed, they had already been wiped empty, recording only 0.01 MB freed.
2. **Missing Custom & Dynamic Paths:** Path resolution only checked a small set of static paths (`~/go/pkg/mod`, `%LOCALAPPDATA%\go-build`) and completely missed developer caches located on secondary drives or non-standard paths, such as `C:\dev-tool\go\cache` (752.68 MB) and `C:\dev-tool\go\pkg\mod` (711.71 MB).

### 1.2 Core Mission
Implement a 3-Tier Multi-Ecosystem Dynamic Discovery Engine that probes active CLI environments, environment variables, and filesystem heuristics across all storage drives (`C:\`, `D:\`, `E:\`, `F:\`) before calculating disk metrics or initiating purges.

---

## 2. 3-Tier Dynamic Discovery Architecture

```
┌────────────────────────────────────────────────────────────────────────┐
│             Multi-Ecosystem Dynamic Discovery Engine                   │
├────────────────────────────────────────────────────────────────────────┤
│ Tier 1: Dynamic CLI Queries                                            │
│   • Go:     go env GOCACHE, go env GOMODCACHE, go env GOPATH           │
│   • Node:   pnpm store path, npm config get cache, yarn cache dir      │
│   • Python: pip cache dir                                              │
│   • Rust:   cargo cache --dir                                          │
│   • .NET:   dotnet nuget locals all -l                                 │
├────────────────────────────────────────────────────────────────────────┤
│ Tier 2: Environment Variables                                          │
│   • GOCACHE, GOMODCACHE, GOPATH                                        │
│   • PNPM_HOME, npm_config_cache, YARN_CACHE_FOLDER, BUN_INSTALL        │
│   • PIP_CACHE_DIR, UV_CACHE_DIR, POETRY_CACHE_DIR                      │
│   • CARGO_HOME, RUSTUP_HOME, NUGET_PACKAGES, GRADLE_USER_HOME, M2_HOME │
├────────────────────────────────────────────────────────────────────────┤
│ Tier 3: Drive & Filesystem Heuristics                                  │
│   • Probing <drive>:\dev-tool\* on all available drives (C:, D:, E:)   │
│   • e.g. C:\dev-tool\go\cache, C:\dev-tool\go\pkg\mod, D:\dev-tool\... │
│   • Standard OS caches: %LOCALAPPDATA%\go-build, ~/.cache/go-build     │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Data Structures & Types

All data structures use affirmative booleans (`Is...`, `Has...`, `Can...`) in accordance with repository coding guidelines.

```go
package cmdos

// DiscoveredCachePath represents a single resolved cache directory on the filesystem.
type DiscoveredCachePath struct {
	Path            string `json:"path"`
	Ecosystem       string `json:"ecosystem"`
	SizeBytes       int64  `json:"sizeBytes"`
	FilesCount      int    `json:"filesCount"`
	DirsCount       int    `json:"dirsCount"`
	IsCustom        bool   `json:"isCustom"`
	IsActive        bool   `json:"isActive"`
	DiscoverySource string `json:"discoverySource"` // "cli", "env", "heuristic", "split_db"
}

// DevDiscoveryOptions controls the depth and targets of dynamic path discovery.
type DevDiscoveryOptions struct {
	TargetEcosystems []string `json:"targetEcosystems,omitempty"`
	HasForceRescan   bool     `json:"hasForceRescan"`
	IsVerbose        bool     `json:"isVerbose"`
}

// DevDiscoveryResult aggregates all resolved paths and pre-purge metrics.
type DevDiscoveryResult struct {
	Paths          []DiscoveredCachePath `json:"paths"`
	TotalSizeBytes int64                 `json:"totalSizeBytes"`
	TotalFiles     int                   `json:"totalFiles"`
	TotalDirs      int                   `json:"totalDirs"`
	DurationMs     int64                 `json:"durationMs"`
}

// EcosystemDiscoveryRule encapsulates the 3-tier probe definition for an ecosystem.
type EcosystemDiscoveryRule struct {
	EcosystemID      string
	DisplayName      string
	CliCommands      [][]string
	EnvVariables     []string
	HeuristicPaths   []string
	DevToolSubdirs   []string
}
```

---

## 4. Function Signatures & Modular Design (<= 15 Lines per Function)

### 4.1 Orchestrator & Discovery Flow (`cli/cmdos/os_dev_clean_discover.go`)

```go
// DiscoverAllDevCaches runs the 3-tier discovery across all requested ecosystems.
func DiscoverAllDevCaches(opts DevDiscoveryOptions) DevDiscoveryResult {
	start := time.Now()
	rules := getEcosystemDiscoveryRules(opts.TargetEcosystems)
	rawPaths := probeAllEcosystems(rules)
	validPaths := filterExistingCachePaths(rawPaths)
	measuredPaths := measureDiscoveredPaths(validPaths)

	return buildDiscoveryResult(measuredPaths, time.Since(start).Milliseconds())
}

// probeAllEcosystems iterates through each ecosystem rule and gathers candidate paths.
func probeAllEcosystems(rules []EcosystemDiscoveryRule) []DiscoveredCachePath {
	var results []DiscoveredCachePath
	for _, rule := range rules {
		results = append(results, probeSingleEcosystem(rule)...)
	}
	return results
}

// probeSingleEcosystem evaluates Tier 1 (CLI), Tier 2 (Env), and Tier 3 (Heuristics).
func probeSingleEcosystem(rule EcosystemDiscoveryRule) []DiscoveredCachePath {
	var candidates []DiscoveredCachePath
	candidates = append(candidates, probeTier1CliPaths(rule)...)
	candidates = append(candidates, probeTier2EnvPaths(rule)...)
	candidates = append(candidates, probeTier3HeuristicPaths(rule)...)
	return deduplicateDiscoveredPaths(candidates)
}

// probeTier1CliPaths executes CLI probe commands with safe timeouts.
func probeTier1CliPaths(rule EcosystemDiscoveryRule) []DiscoveredCachePath {
	var items []DiscoveredCachePath
	for _, cmdArgs := range rule.CliCommands {
		if path, ok := executeSafeCliProbe(cmdArgs); ok && len(path) > 0 {
			items = append(items, DiscoveredCachePath{
				Path:            path,
				Ecosystem:       rule.EcosystemID,
				DiscoverySource: "cli",
			})
		}
	}
	return items
}

// probeTier2EnvPaths inspects system environment variables for custom cache roots.
func probeTier2EnvPaths(rule EcosystemDiscoveryRule) []DiscoveredCachePath {
	var items []DiscoveredCachePath
	for _, envKey := range rule.EnvVariables {
		if val := strings.TrimSpace(os.Getenv(envKey)); len(val) > 0 {
			items = append(items, DiscoveredCachePath{
				Path:            val,
				Ecosystem:       rule.EcosystemID,
				DiscoverySource: "env",
			})
		}
	}
	return items
}

// probeTier3HeuristicPaths evaluates drive-wide and OS-standard path locations.
func probeTier3HeuristicPaths(rule EcosystemDiscoveryRule) []DiscoveredCachePath {
	var items []DiscoveredCachePath
	for _, dir := range probeWindowsDevDrives(rule.DevToolSubdirs) {
		items = append(items, DiscoveredCachePath{
			Path:            dir,
			Ecosystem:       rule.EcosystemID,
			IsCustom:        true,
			DiscoverySource: "heuristic",
		})
	}
	for _, raw := range rule.HeuristicPaths {
		items = append(items, DiscoveredCachePath{
			Path:            expandUserPath(raw),
			Ecosystem:       rule.EcosystemID,
			DiscoverySource: "heuristic",
		})
	}
	return items
}
```

### 4.2 Drive Scanning & Metric Measurement (`cli/cmdos/os_dev_clean_paths.go`)

```go
// probeWindowsDevDrives scans all mounted drive letters for dev-tool subpaths.
func probeWindowsDevDrives(subpaths []string) []string {
	var discovered []string
	drives := getAvailableDriveRoots()
	for _, drive := range drives {
		for _, sub := range subpaths {
			candidate := filepath.Join(drive, "dev-tool", sub)
			if isExistingDirectory(candidate) {
				discovered = append(discovered, candidate)
			}
		}
	}
	return discovered
}

// getAvailableDriveRoots returns existing root drives (C:\, D:\, E:\, etc.).
func getAvailableDriveRoots() []string {
	if runtime.GOOS != "windows" {
		return []string{"/"}
	}
	var roots []string
	for _, d := range []string{"C", "D", "E", "F", "G", "X", "Z"} {
		root := d + `:\`
		if isExistingDirectory(root) {
			roots = append(roots, root)
		}
	}
	return roots
}

// measureDiscoveredPaths populates size and file counts before any clean action.
func measureDiscoveredPaths(paths []DiscoveredCachePath) []DiscoveredCachePath {
	var measured []DiscoveredCachePath
	for _, p := range paths {
		bytes, files, dirs := calculateDirectoryMetrics(p.Path)
		p.SizeBytes = bytes
		p.FilesCount = files
		p.DirsCount = dirs
		p.IsActive = true
		measured = append(measured, p)
	}
	return measured
}

// calculateDirectoryMetrics traverses a directory to calculate total byte size and counts.
func calculateDirectoryMetrics(root string) (int64, int, int) {
	var totalBytes int64
	var fileCount, dirCount int
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		if info.IsDir() {
			dirCount++
		} else {
			fileCount++
			totalBytes += info.Size()
		}
		return nil
	})
	return totalBytes, fileCount, dirCount
}

// isExistingDirectory returns true if path exists and is a directory.
func isExistingDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info != nil && info.IsDir()
}
```

---

## 5. Ecosystem Probe Specifications

| Ecosystem | Tier 1 CLI Query | Tier 2 Env Variables | Tier 3 Heuristics (`<drive>:\dev-tool\*`) |
| :--- | :--- | :--- | :--- |
| **Go** | `go env GOCACHE`<br>`go env GOMODCACHE`<br>`go env GOPATH` | `GOCACHE`<br>`GOMODCACHE`<br>`GOPATH` | `go\cache`<br>`go\pkg\mod`<br>`go\testcache`<br>`%LOCALAPPDATA%\go-build` |
| **Node / pnpm** | `pnpm store path` | `PNPM_HOME` | `pnpm\store`<br>`pnpm\cache`<br>`%LOCALAPPDATA%\pnpm\store` |
| **Node / npm** | `npm config get cache` | `npm_config_cache` | `npm\cache`<br>`%LOCALAPPDATA%\npm-cache`<br>`~/.npm` |
| **Node / Yarn** | `yarn cache dir` | `YARN_CACHE_FOLDER` | `yarn\cache`<br>`%LOCALAPPDATA%\Yarn\Cache` |
| **Node / Bun** | `bun pm cache` | `BUN_INSTALL` | `bun\cache`<br>`~/.bun/install/cache` |
| **Python** | `pip cache dir` | `PIP_CACHE_DIR`<br>`UV_CACHE_DIR`<br>`POETRY_CACHE_DIR` | `python\pip\cache`<br>`%LOCALAPPDATA%\pip\cache`<br>`~/.cache/pip` |
| **Rust** | `cargo cache --dir` | `CARGO_HOME`<br>`RUSTUP_HOME` | `cargo\registry\cache`<br>`~/.cargo/registry/cache` |
| **.NET / Nuget** | `dotnet nuget locals all -l` | `NUGET_PACKAGES` | `nuget\cache`<br>`%LOCALAPPDATA%\NuGet\v3-cache` |
| **Java / Gradle** | — | `GRADLE_USER_HOME`<br>`M2_HOME` | `gradle\caches`<br>`~/.gradle/caches`<br>`~/.m2/repository` |

---

## 6. Safety & Protection Gates

1. **Root Directory Protection:** The discovery engine explicitly prevents sweeping root directories (e.g. `C:\`, `D:\dev-tool`), system folders (`C:\Windows`, `C:\Program Files`), or user profiles (`C:\Users\username`).
2. **Affirmative Guard:** `isDevCachePathSafe(path string) bool` checks that the leaf path is a known cache naming token (e.g. `cache`, `testcache`, `mod`, `store`, `repository`, `caches`, `v3-cache`, `npm-cache`, `go-build`).
3. **Pre-Purge Guarantee:** Under no circumstances should `go clean` or `pnpm store prune` run before `measureDiscoveredPaths` returns accurate file counts and byte sizes.

---

## 7. Concrete Verification & Acceptance Criteria

- **AC-1:** On a system with Go installed in `C:\dev-tool\go\`, discovery identifies both `C:\dev-tool\go\cache` and `C:\dev-tool\go\pkg\mod`.
- **AC-2:** Discovery resolves dynamic output from `go env GOCACHE`, `go env GOMODCACHE`, and `pnpm store path` if executables exist.
- **AC-3:** Pre-purge sizing accurately reports hundreds of megabytes (or gigabytes) instead of the previous 0.01 MB metric defect.
- **AC-4:** Multi-drive scanning inspects `C:\`, `D:\`, `E:\`, and `F:\` without throwing errors on missing or unmounted drives.
- **AC-5:** All written functions comply with the repository limit of `<= 15 lines` and use affirmative boolean naming exclusively.
