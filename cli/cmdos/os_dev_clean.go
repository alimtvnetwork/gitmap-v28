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
	_ = DiscoverDevToolCaches(opts)
	return executeAndRenderClean(opts)
}

func executeAndRenderClean(opts DevCleanOptions) error {
	osOpts := osclean.DevCleanOptions{IsDryRun: opts.IsDryRun, HasAutoYes: opts.HasAutoYes, IsJSON: opts.IsJSON, IsVerbose: opts.IsVerbose, OnlyCategories: opts.OnlyCategories}
	res := osclean.CleanDevCaches(osOpts)
	if res.IsFailure() {
		return res.AppError()
	}
	renderDevCleanOutput(res.Value, opts)
	return nil
}

// DiscoverDevToolCaches discovers active cache paths, leveraging Split-DB persistence.
func DiscoverDevToolCaches(opts DevCleanOptions) []store.DevtoolsCacheRecord {
	sdb, err := store.OpenDevtoolsCacheSplitDB()
	if err != nil || sdb == nil {
		return nil
	}
	defer sdb.Close()
	if opts.IsForce {
		_ = sdb.InvalidateCache("")
	} else if cached, err := sdb.GetCachedPaths(""); err == nil && len(cached) > 0 {
		return cached
	}
	records := probeKnownDevCaches(opts.OnlyCategories)
	if len(records) > 0 {
		_ = sdb.SaveDiscoveredPaths(records)
	}
	return records
}

func probeKnownDevCaches(only []string) []store.DevtoolsCacheRecord {
	var records []store.DevtoolsCacheRecord
	for _, p := range []string{`C:\dev-tool\go\cache`, `C:\dev-tool\go\pkg\mod`, `D:\dev-tool\pnpm\store`} {
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			records = append(records, store.DevtoolsCacheRecord{Path: p, Ecosystem: "dev", IsCustom: true, IsActive: true})
		}
	}
	return records
}

func isPromptRequired(opts DevCleanOptions) bool {
	return !opts.IsDryRun && !opts.HasAutoYes && !opts.IsJSON
}

func confirmDevClean() bool {
	fmt.Print("  Proceed with dev tools cache cleanup? Type 'yes' to continue: ")
	text, err := bufio.NewReader(os.Stdin).ReadString('\n')
	return err == nil && strings.ToLower(strings.TrimSpace(text)) == "yes"
}

func renderDevCleanOutput(summary osclean.DevCleanSummary, opts DevCleanOptions) {
	if opts.IsJSON {
		data, _ := json.MarshalIndent(summary, "", "  ")
		fmt.Println(string(data))
	} else if opts.IsTree {
		renderDevCleanTree(summary)
	} else {
		osclean.RenderEnhancedSummaryTable(summary)
	}
}

func renderDevCleanTree(summary osclean.DevCleanSummary) {
	mbTotal := float64(summary.TotalBytesFreed) / (1024 * 1024)
	fmt.Printf("\n  🌳 Developer Tools Cache Tree (Total: %.2f MB, %d files, %d dirs)\n", mbTotal, summary.TotalItemsRemoved, summary.TotalDirsRemoved)
	for i, cat := range summary.Categories {
		branch := "├──"
		if i == len(summary.Categories)-1 {
			branch = "└──"
		}
		fmt.Printf("  %s 📦 %s [%.2f MB - %d files]\n", branch, cat.Label, float64(cat.BytesFreed)/(1024*1024), cat.ItemsRemoved)
	}
}
