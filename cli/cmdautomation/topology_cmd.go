package cmdautomation

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	topologyCmd = &cobra.Command{
		Use:     "topology [dir]",
		Aliases: []string{"topo", "codebase-topology"},
		Short:   "Discover and cache codebase topology (subsystems, schemas, workflows, languages)",
		RunE:    runTopologyCmd,
	}

	topologyOpts TopologyOptions
)

func runTopologyCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		topologyOpts.Dir = args[0]
	}
	monad := RunTopology(topologyOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	res := monad.Value
	if topologyOpts.Query != "" {
		renderTopologyQuery(res, topologyOpts.Query)
		return nil
	}
	renderTopologyResult(res, topologyOpts.IsJson)
	return nil
}

func renderTopologyResult(res TopologyResult, isJson bool) {
	if isJson {
		printTopologyJson(res)
		return
	}
	printTopologyTerminal(res)
}

func printTopologyJson(res TopologyResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printTopologyTerminal(res TopologyResult) {
	fmt.Printf("\n%s[Polyglot Codebase Topology Map]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Root:       %s\n", res.RootPath)
	fmt.Printf("  Files:      %d\n", res.TotalFiles)
	fmt.Printf("  Scan Time:  %.2fms\n", res.DurationMs)
	fmt.Printf("  Expires:    %s\n\n", res.ExpiresAt)

	printLanguageBreakdown(res.Languages, res.Manifests)
	printSubsystemBreakdown(res.Subsystems)
}

func printLanguageBreakdown(langs map[string]int, manifests map[string][]string) {
	fmt.Printf("%s📊 Language Breakdown:%s\n", constants.ColorCyan, constants.ColorReset)
	keys := make([]string, 0, len(langs))
	for k := range langs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		count := langs[k]
		manList := manifests[k]
		manStr := ""
		if len(manList) > 0 {
			manStr = fmt.Sprintf(" (manifests: %s)", strings.Join(manList, ", "))
		}
		fmt.Printf("  • %-12s : %5d files%s\n", k, count, manStr)
	}
	fmt.Println()
}

func printSubsystemBreakdown(subsystems map[string]SubsystemData) {
	fmt.Printf("%s🏛️ Subsystem Map:%s\n", constants.ColorCyan, constants.ColorReset)
	keys := make([]string, 0, len(subsystems))
	for k := range subsystems {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		data := subsystems[k]
		display := formatSubsystemRoots(data.Roots)
		fmt.Printf("  • %-10s : %s\n", k, display)
	}
	fmt.Println()
}

func formatSubsystemRoots(roots []string) string {
	if len(roots) == 0 {
		return "(none detected)"
	}
	limit := 3
	if len(roots) <= limit {
		return strings.Join(roots, ", ")
	}
	head := strings.Join(roots[:limit], ", ")
	return fmt.Sprintf("%s (+%d more)", head, len(roots)-limit)
}

func renderTopologyQuery(res TopologyResult, query string) {
	q := strings.ToLower(strings.TrimSpace(query))
	fmt.Printf("\n%s⚡ Topology Routing Query for: `%s`%s\n\n", constants.ColorBold, query, constants.ColorReset)
	hasSubsys := printMatchedSubsystem(res.Subsystems, q)
	hasLang := printMatchedLanguage(res.Languages, res.Manifests, res.LangRoots, q)
	if !hasSubsys && !hasLang {
		fmt.Printf("%s⚠️ No direct matches found for '%s'.%s\n\n", constants.ColorYellow, query, constants.ColorReset)
	}
}

func printMatchedSubsystem(subsystems map[string]SubsystemData, q string) bool {
	data, exists := subsystems[q]
	if !exists {
		return false
	}
	fmt.Printf("%s📦 Subsystem: [%s]%s\n", constants.ColorGreen, q, constants.ColorReset)
	for _, r := range data.Roots {
		fmt.Printf("   📁 Root: %s\n", r)
	}
	for _, s := range data.SchemaFiles {
		fmt.Printf("   📄 Schema: %s\n", s)
	}
	for _, e := range data.Entrypoints {
		fmt.Printf("   🚀 Entry: %s\n", e)
	}
	fmt.Println()
	return true
}

func printMatchedLanguage(langs map[string]int, manifests map[string][]string, roots map[string][]string, q string) bool {
	count, exists := langs[q]
	if !exists {
		return false
	}
	fmt.Printf("%s🔤 Language: [%s] — %d files%s\n", constants.ColorGreen, q, count, constants.ColorReset)
	for _, m := range manifests[q] {
		fmt.Printf("   📋 Manifest: %s\n", m)
	}
	for _, r := range roots[q] {
		fmt.Printf("   📁 Directory: %s\n", r)
	}
	fmt.Println()
	return true
}

func initTopologyFlags() {
	topologyCmd.Flags().BoolVarP(&topologyOpts.IsRefresh, "refresh", "r", false, "Force refresh topology discovery cache")
	topologyCmd.Flags().BoolVarP(&topologyOpts.IsJson, "json", "j", false, "Output results as raw JSON")
	topologyCmd.Flags().StringVarP(&topologyOpts.Query, "query", "q", "", "Search specific subsystem or language routing")
	topologyCmd.Flags().IntVar(&topologyOpts.TtlSec, "ttl", 1800, "Cache TTL in seconds")
}

func init() {
	AutomationCmd.AddCommand(topologyCmd)
	initTopologyFlags()
}
