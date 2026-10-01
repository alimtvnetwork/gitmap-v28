package cmdos

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DiscoveredCachePath struct {
	Path            string `json:"path"`
	Category        string `json:"category"`
	Ecosystem       string `json:"ecosystem,omitempty"`
	Source          string `json:"source,omitempty"`
	SizeBytes       int64  `json:"sizeBytes"`
	FilesCount      int    `json:"filesCount"`
	DirsCount       int    `json:"dirsCount"`
	IsCustom        bool   `json:"isCustom"`
	IsActive        bool   `json:"isActive"`
	DiscoverySource string `json:"discoverySource"`
}

type DevDiscoveryOptions struct {
	TargetEcosystems []string `json:"targetEcosystems,omitempty"`
	HasForceRescan   bool     `json:"hasForceRescan"`
	IsVerbose        bool     `json:"isVerbose"`
}

type DevDiscoveryResult struct {
	Paths          []DiscoveredCachePath `json:"paths"`
	TotalSizeBytes int64                 `json:"totalSizeBytes"`
	TotalFiles     int                   `json:"totalFiles"`
	TotalDirs      int                   `json:"totalDirs"`
	DurationMs     int64                 `json:"durationMs"`
}

type EcosystemDiscoveryRule struct {
	EcosystemID    string
	DisplayName    string
	CliCommands    [][]string
	EnvVariables   []string
	HeuristicPaths []string
	DevToolSubdirs []string
}

func newDiscoveredPath(path, eco, source string, isCustom bool) DiscoveredCachePath {
	return DiscoveredCachePath{
		Path:            filepath.Clean(path),
		Category:        eco,
		Ecosystem:       eco,
		Source:          source,
		DiscoverySource: source,
		IsCustom:        isCustom,
		IsActive:        true,
	}
}

func DiscoverAllDevCaches(opts DevDiscoveryOptions) DevDiscoveryResult {
	start := time.Now()
	rules := getEcosystemDiscoveryRules(opts.TargetEcosystems)
	rawPaths := probeAllEcosystems(rules)
	uniquePaths := deduplicateDiscoveredPaths(rawPaths)
	validPaths := filterExistingCachePaths(uniquePaths)
	measuredPaths := measureDiscoveredPaths(validPaths)
	return buildDiscoveryResult(measuredPaths, time.Since(start).Milliseconds())
}

func probeAllEcosystems(rules []EcosystemDiscoveryRule) []DiscoveredCachePath {
	var results []DiscoveredCachePath
	for _, rule := range rules {
		results = append(results, probeSingleEcosystem(rule)...)
	}
	return results
}

func probeSingleEcosystem(rule EcosystemDiscoveryRule) []DiscoveredCachePath {
	var candidates []DiscoveredCachePath
	candidates = append(candidates, probeTier1CliPaths(rule)...)
	candidates = append(candidates, probeTier2EnvPaths(rule)...)
	candidates = append(candidates, probeTier3HeuristicPaths(rule)...)
	return candidates
}

func deduplicateDiscoveredPaths(paths []DiscoveredCachePath) []DiscoveredCachePath {
	seen := make(map[string]bool)
	var unique []DiscoveredCachePath
	for _, p := range paths {
		clean := filepath.Clean(strings.ToLower(p.Path))
		if len(clean) > 0 && !seen[clean] {
			seen[clean] = true
			unique = append(unique, p)
		}
	}
	return unique
}

func filterExistingCachePaths(raw []DiscoveredCachePath) []DiscoveredCachePath {
	var valid []DiscoveredCachePath
	for _, p := range raw {
		if isExistingDirectory(p.Path) && isDevCachePathSafe(p.Path) {
			p.IsActive = true
			valid = append(valid, p)
		}
	}
	return valid
}

func measureDiscoveredPaths(paths []DiscoveredCachePath) []DiscoveredCachePath {
	var measured []DiscoveredCachePath
	for _, p := range paths {
		bytes, files, dirs := calculateDirectoryMetrics(p.Path)
		p.SizeBytes = bytes
		p.FilesCount = files
		p.DirsCount = dirs
		measured = append(measured, p)
	}
	return measured
}

func walkMetricsStep(info os.FileInfo, bytes *int64, files, dirs *int) {
	if info.IsDir() {
		*dirs++
		return
	}
	*files++
	*bytes += info.Size()
}

func calculateDirectoryMetrics(root string) (int64, int, int) {
	var totalBytes int64
	var fileCount, dirCount int
	_ = filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
		if err == nil && info != nil {
			walkMetricsStep(info, &totalBytes, &fileCount, &dirCount)
		}
		return nil
	})
	return totalBytes, fileCount, dirCount
}

func isExistingDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info != nil && info.IsDir()
}

func isDevCachePathSafe(path string) bool {
	clean := filepath.Clean(strings.ToLower(path))
	if len(clean) <= 3 || strings.HasPrefix(clean, `c:\windows`) {
		return false
	}
	safeTokens := []string{"cache", "testcache", "mod", "store", "repository", "caches", "v3-cache", "npm-cache", "go-build"}
	for _, tok := range safeTokens {
		if strings.Contains(filepath.Base(clean), tok) || strings.Contains(clean, tok) {
			return true
		}
	}
	return false
}

func sumDiscoveryMetrics(paths []DiscoveredCachePath) (int64, int, int) {
	var bytes int64
	var files, dirs int
	for _, p := range paths {
		bytes += p.SizeBytes
		files += p.FilesCount
		dirs += p.DirsCount
	}
	return bytes, files, dirs
}

func buildDiscoveryResult(paths []DiscoveredCachePath, durationMs int64) DevDiscoveryResult {
	bytes, files, dirs := sumDiscoveryMetrics(paths)
	return DevDiscoveryResult{
		Paths:          paths,
		TotalSizeBytes: bytes,
		TotalFiles:     files,
		TotalDirs:      dirs,
		DurationMs:     durationMs,
	}
}
