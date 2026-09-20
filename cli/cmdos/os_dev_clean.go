package cmdos

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/osclean"
)

// RunOSDevClean dispatches the developer tools cache cleaner.
func RunOSDevClean(args []string) error {
	if isHelpDevClean(args) {
		printOSDevCleanUsage()
		return nil
	}

	opts := parseDevCleanOptions(args)
	if isPromptRequired(opts) && !confirmDevClean() {
		fmt.Println("  Aborted by operator.")
		return nil
	}

	res := osclean.CleanDevCaches(opts)
	if res.IsFailure() {
		return res.AppError()
	}

	renderDevCleanOutput(res.Value, opts)
	return nil
}

func isHelpDevClean(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "-h" || low == "--help" || low == "help" || low == "/?" {
			return true
		}
	}

	return false
}

func parseDevCleanOptions(args []string) osclean.DevCleanOptions {
	opts := osclean.DevCleanOptions{}
	for i := 0; i < len(args); i++ {
		a := strings.ToLower(args[i])
		parseDevCleanFlag(a, args, &i, &opts)
	}

	return opts
}

func parseDevCleanFlag(a string, args []string, i *int, opts *osclean.DevCleanOptions) {
	switch {
	case a == "--dry-run" || a == "-n" || a == "-d" || a == "dry-run":
		opts.IsDryRun = true
	case a == "--yes" || a == "-y" || a == "/y" || a == "yes":
		opts.HasAutoYes = true
	case a == "--json":
		opts.IsJSON = true
	case a == "--verbose" || a == "-v":
		opts.IsVerbose = true
	case strings.HasPrefix(a, "--only="):
		opts.OnlyCategories = strings.Split(strings.TrimPrefix(a, "--only="), ",")
	case a == "--only" && *i+1 < len(args):
		*i++
		opts.OnlyCategories = strings.Split(args[*i], ",")
	}
}

func isPromptRequired(opts osclean.DevCleanOptions) bool {
	return !opts.IsDryRun && !opts.HasAutoYes && !opts.IsJSON
}

func confirmDevClean() bool {
	fmt.Print("  Proceed with dev tools cache cleanup? Type 'yes' to continue: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	return strings.ToLower(strings.TrimSpace(text)) == "yes"
}

func renderDevCleanOutput(summary osclean.DevCleanSummary, opts osclean.DevCleanOptions) {
	if opts.IsJSON {
		data, _ := json.MarshalIndent(summary, "", "  ")
		fmt.Println(string(data))
		return
	}

	printDevCleanBanner(opts.IsDryRun)
	for _, cat := range summary.Categories {
		renderCategoryStats(cat, opts.IsDryRun, opts.IsVerbose)
	}
	renderDevCleanSummaryTotals(summary)
}

func printDevCleanBanner(isDryRun bool) {
	fmt.Println()
	fmt.Println("  OS Dev-Cleanup (Developer Tools Cache Remover)")
	fmt.Println("  ==============================================")
	if isDryRun {
		fmt.Println("  [DRY-RUN] Preview only. No cache files will be deleted.")
	}
	fmt.Println()
}

func renderCategoryStats(cat osclean.CategoryCleanStats, isDryRun, isVerbose bool) {
	mb := float64(cat.BytesFreed) / (1024 * 1024)
	action := "Cleaned"
	if isDryRun {
		action = "Would clean"
	}
	fmt.Printf("  • %-18s %s %d files, %d dirs (%.2f MB freed)\n",
		cat.Category+":", action, cat.ItemsRemoved, cat.DirsRemoved, mb)

	if isVerbose {
		printVerboseNotes(cat.Notes, cat.Errors)
	}
}

func printVerboseNotes(notes, errors []string) {
	for _, n := range notes {
		fmt.Printf("      note: %s\n", n)
	}
	for _, e := range errors {
		fmt.Printf("      warning: %s\n", e)
	}
}

func renderDevCleanSummaryTotals(summary osclean.DevCleanSummary) {
	mb := float64(summary.TotalBytesFreed) / (1024 * 1024)
	mode := "Cleaned"
	if summary.IsDryRun {
		mode = "[Dry-Run] Reclaimable"
	}
	fmt.Println()
	fmt.Printf("  ✔ %s Total: %.2f MB across %d files and %d dirs (%dms)\n\n",
		mode, mb, summary.TotalItemsRemoved, summary.TotalDirsRemoved, summary.DurationMs)
}

func printOSDevCleanUsage() {
	fmt.Println("Usage: gitmap clean-dev [flags]")
	fmt.Println("       gitmap os dev-clean [flags]")
	fmt.Println("       gitmap os cleanup dev [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -n, --dry-run      Preview space reclaimed without deleting")
	fmt.Println("  -y, --yes          Bypass confirmation prompt")
	fmt.Println("      --json         Output structured JSON summary")
	fmt.Println("  -v, --verbose      Display individual subpaths and command notes")
	fmt.Println("      --only <cats>  Comma-separated categories to clean (e.g. go,npm,pnpm)")
	fmt.Println("  -h, --help         Show this documentation")
}
