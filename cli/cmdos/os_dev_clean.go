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
	paths := DiscoverDevToolCaches(opts)
	if len(paths) == 0 {
		fmt.Println("  No active developer tool caches found.")
		return nil
	}

	return handleDevCleanExecution(paths, opts)
}

func handleDevCleanExecution(paths []DiscoveredCachePath, opts DevCleanOptions) error {
	if opts.IsDryRun {
		renderDevCleanOutput(paths, opts)
		return nil
	}
	if isPromptRequired(opts) && !promptAndConfirmClean(paths, opts) {
		fmt.Println("  Aborted by operator. No files were removed.")
		return nil
	}

	return executeAndRenderClean(paths, opts)
}

func promptAndConfirmClean(paths []DiscoveredCachePath, opts DevCleanOptions) bool {
	dryOpts := opts
	dryOpts.IsDryRun = true
	renderDevCleanOutput(paths, dryOpts)

	return confirmDevClean()
}

func executeAndRenderClean(paths []DiscoveredCachePath, opts DevCleanOptions) error {
	if !opts.IsDryRun {
		purgeDiscoveredPaths(paths)
		_ = osccleanFallbackClean()
		updateSplitDBAfterClean(paths)
	}
	renderDevCleanOutput(paths, opts)
	return nil
}

func updateSplitDBAfterClean(paths []DiscoveredCachePath) {
	sdb, err := store.OpenDevtoolsCacheSplitDB()
	if err != nil || sdb == nil {
		return
	}
	defer sdb.Close()
	var zeroed []DiscoveredCachePath
	for _, p := range paths {
		p.SizeBytes = 0
		p.FilesCount = 0
		p.DirsCount = 0
		zeroed = append(zeroed, p)
	}
	_ = sdb.SaveDiscoveredPaths(ConvertDiscoveredToSplitRecords(zeroed))
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
		return discoverAndPersistDevCaches(sdb, opts)
	}

	cachedPaths := resolveExistingCachedPaths(sdb)
	if len(cachedPaths) == 0 {
		return discoverAndPersistDevCaches(sdb, opts)
	}

	cachedPaths = mergeMissingEcosystemPaths(sdb, cachedPaths, opts)

	return measureDiscoveredPaths(cachedPaths)
}

func mergeMissingEcosystemPaths(sdb *store.DevtoolsCacheSplitDB, cached []DiscoveredCachePath, opts DevCleanOptions) []DiscoveredCachePath {
	rules := getEcosystemDiscoveryRules(opts.OnlyCategories)
	missingRules := findMissingEcosystemRules(cached, rules)
	if len(missingRules) == 0 {
		return cached
	}

	newPaths := filterExistingCachePaths(probeAllEcosystems(missingRules))
	if len(newPaths) == 0 {
		return cached
	}
	_ = sdb.SaveDiscoveredPaths(ConvertDiscoveredToSplitRecords(newPaths))

	return append(cached, newPaths...)
}

func findMissingEcosystemRules(cached []DiscoveredCachePath, rules []EcosystemDiscoveryRule) []EcosystemDiscoveryRule {
	var missing []EcosystemDiscoveryRule
	for _, rule := range rules {
		if !hasEcosystemInPaths(cached, rule.EcosystemID) {
			missing = append(missing, rule)
		}
	}

	return missing
}

func hasEcosystemInPaths(paths []DiscoveredCachePath, ecoID string) bool {
	norm := normalizeEcosystemID(ecoID)
	for _, p := range paths {
		if normalizeEcosystemID(getPathEcosystemKey(p)) == norm {
			return true
		}
	}

	return false
}

func resolveExistingCachedPaths(sdb *store.DevtoolsCacheSplitDB) []DiscoveredCachePath {
	cached, err := sdb.GetCachedPaths("")
	if err != nil || len(cached) == 0 {
		return nil
	}

	valid := filterExistingCachePaths(ConvertSplitRecordsToDiscovered(cached))
	if len(valid) == 0 {
		return nil
	}

	return valid
}

func discoverAndPersistDevCaches(sdb *store.DevtoolsCacheSplitDB, opts DevCleanOptions) []DiscoveredCachePath {
	disc := DiscoverAllDevCaches(DevDiscoveryOptions{TargetEcosystems: opts.OnlyCategories, HasForceRescan: opts.IsForce, IsVerbose: opts.IsVerbose})
	_ = sdb.SaveDiscoveredPaths(ConvertDiscoveredToSplitRecords(disc.Paths))

	return disc.Paths
}

func isPromptRequired(opts DevCleanOptions) bool {
	return !opts.IsDryRun && !opts.HasAutoYes && !opts.IsJSON
}

func confirmDevClean() bool {
	fmt.Print("  Proceed with dev tools cache cleanup? Type 'yes' (or 'y') to continue: ")
	text, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false
	}
	trimmed := strings.ToLower(strings.TrimSpace(text))

	return trimmed == "yes" || trimmed == "y"
}

func renderDevCleanOutput(paths []DiscoveredCachePath, opts DevCleanOptions) {
	switch {
	case opts.IsJSON:
		res := buildDiscoveryResult(paths, 0)
		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(data))
	case opts.IsTree:
		RenderDevTreeOutput(paths, DevTreeRenderOptions{HasColorEnabled: true})
	default:
		RenderDevCleanAnsiCards(paths, opts.IsDryRun)
	}
}
