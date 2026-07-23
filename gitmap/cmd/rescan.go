package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// runRescan handles the "rescan" subcommand.
func runRescan() {
	cache, err := loadScanCache()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRescanNoCache, err)
		os.Exit(1)
	}
	fmt.Printf(constants.MsgRescanReplay, cache.Dir)
	runScanFromCache(cache)
}

// loadScanCache reads the last-scan.json from the output folder.
func loadScanCache() (model.ScanCache, error) {
	path := filepath.Join(constants.DefaultOutputFolder, constants.DefaultScanCacheFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return model.ScanCache{}, err
	}
	var cache model.ScanCache
	err = json.Unmarshal(data, &cache)

	return cache, err
}

// saveScanCache writes the current scan flags to last-scan.json.
func saveScanCache(outputDir string, cache model.ScanCache) {
	path := filepath.Join(outputDir, constants.DefaultScanCacheFile)
	data, err := json.MarshalIndent(cache, "", constants.JSONIndent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to marshal scan cache: %v\n", err)

		return
	}
	if err := os.MkdirAll(filepath.Dir(path), constants.DirPermission); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not create cache directory: %v\n", err)

		return
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not write scan cache to %s: %v\n", path, err)

		return
	}
	fmt.Printf(constants.MsgScanCacheSaved, filepath.Base(path))
}

// runScanFromCache replays a scan using cached flags.
func runScanFromCache(c model.ScanCache) {
	args := buildScanArgs(c)
	runScan(args)
}

// buildScanArgs reconstructs CLI args from a ScanCache.
func buildScanArgs(c model.ScanCache) []string {
	var args []string
	args = appendStringFlag(args, "--config", c.ConfigPath)
	args = appendStringFlag(args, "--mode", c.Mode)
	args = appendStringFlag(args, "--output", c.Output)
	args = appendStringFlag(args, "--out-file", c.OutFile)
	args = appendStringFlag(args, "--output-path", c.OutputPath)
	args = appendBoolFlag(args, "--github-desktop", c.GithubDesktop)
	args = appendBoolFlag(args, "--open", c.OpenFolder)
	args = appendBoolFlag(args, "--quiet", c.Quiet)
	if len(c.Dir) > 0 && c.Dir != constants.DefaultDir {
		args = append(args, c.Dir)
	}

	return args
}

// appendStringFlag appends a flag pair if the value is non-empty.
func appendStringFlag(args []string, flag, value string) []string {
	if len(value) > 0 {
		return append(args, flag, value)
	}

	return args
}

// appendBoolFlag appends a flag if the condition is true.
func appendBoolFlag(args []string, flag string, enabled bool) []string {
	if enabled {
		return append(args, flag)
	}

	return args
}
