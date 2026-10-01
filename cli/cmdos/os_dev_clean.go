package cmdos

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/osclean"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunOSDevClean dispatches the developer tools cache cleaner.
func RunOSDevClean(args []string) error {
	if isHelpDevClean(args) {
		printOSDevCleanUsage()
		return nil
	}
	opts := ParseDevCleanOptions(args)
	if isPromptRequired(opts) && !confirmDevClean() {
		fmt.Println("  Aborted by operator.")
		return nil
	}
	paths := DiscoverDevToolCaches(opts)
	return executeAndRenderClean(paths, opts)
}

func executeAndRenderClean(paths []DiscoveredCachePath, opts DevCleanOptions) error {
	if !opts.IsDryRun {
		purgeDiscoveredPaths(paths)
		_ = osccleanFallbackClean()
	}
	renderDevCleanOutput(paths, opts)
	return nil
}

func purgeDiscoveredPaths(paths []DiscoveredCachePath) {
	for _, p := range paths {
		_ = osclean.SweepTarget(p.Path, false, true)
	}
}

func osccleanFallbackClean() error {
	res := osclean.CleanDevCaches(osclean.DevCleanOptions{IsDryRun: false, HasAutoYes: true})
	if res.IsFailure() {
		return res.AppError()
	}
	return nil
}

// DiscoverDevToolCaches discovers active cache paths, leveraging Split-DB persistence.
func DiscoverDevToolCaches(opts DevCleanOptions) []DiscoveredCachePath {
	sdb, err := store.OpenDevtoolsCacheSplitDB()
	if err != nil || sdb == nil {
		disc := DiscoverAllDevCaches(DevDiscoveryOptions{TargetEcosystems: opts.OnlyCategories, HasForceRescan: true, IsVerbose: opts.IsVerbose})
		return filterByOnlyCategories(disc.Paths, opts.OnlyCategories)
	}
	defer sdb.Close()
	return filterByOnlyCategories(resolveCacheFromSplitDB(sdb, opts), opts.OnlyCategories)
}

func filterByOnlyCategories(paths []DiscoveredCachePath, only []string) []DiscoveredCachePath {
	if len(only) == 0 {
		return paths
	}
	var filtered []DiscoveredCachePath
	for _, p := range paths {
		if containsEcosystem(only, getPathEcosystemKey(p)) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func resolveCacheFromSplitDB(sdb *store.DevtoolsCacheSplitDB, opts DevCleanOptions) []DiscoveredCachePath {
	if opts.IsForce {
		_ = sdb.InvalidateCache("")
	} else if cached, err := sdb.GetCachedPaths(""); err == nil && len(cached) > 0 {
		if valid := filterExistingCachePaths(ConvertSplitRecordsToDiscovered(cached)); len(valid) > 0 {
			return measureDiscoveredPaths(valid)
		}
	}
	disc := DiscoverAllDevCaches(DevDiscoveryOptions{TargetEcosystems: opts.OnlyCategories, HasForceRescan: opts.IsForce, IsVerbose: opts.IsVerbose})
	_ = sdb.SaveDiscoveredPaths(ConvertDiscoveredToSplitRecords(disc.Paths))
	return disc.Paths
}

func isPromptRequired(opts DevCleanOptions) bool {
	return !opts.IsDryRun && !opts.HasAutoYes && !opts.IsJSON
}

func confirmDevClean() bool {
	fmt.Print("  Proceed with dev tools cache cleanup? Type 'yes' to continue: ")
	text, err := bufio.NewReader(os.Stdin).ReadString('\n')
	return err == nil && strings.ToLower(strings.TrimSpace(text)) == "yes"
}

func renderDevCleanOutput(paths []DiscoveredCachePath, opts DevCleanOptions) {
	if opts.IsJSON {
		res := buildDiscoveryResult(paths, 0)
		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(data))
	} else if opts.IsTree {
		RenderDevTreeOutput(paths, DevTreeRenderOptions{HasColorEnabled: true})
	} else {
		RenderDevCleanAnsiCards(paths, opts.IsDryRun)
	}
}
