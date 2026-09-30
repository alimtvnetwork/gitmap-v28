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
	if isDevCleanIgnoredToken(a) {
		return
	}
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

func isDevCleanIgnoredToken(a string) bool {
	return a == "clear" || a == "clean" || a == "cleanup" || a == "cache" ||
		a == "caches" || a == "dev" || a == "devs" || a == "developer" ||
		a == "devtool" || a == "devtools" || a == "dev-tool" || a == "dev-tools" ||
		a == "tool" || a == "tools" || a == "devtool-cache" || a == "devtools-cache" ||
		a == "dev-tool-cache" || a == "dev-tools-cache"
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

	osclean.RenderEnhancedSummaryTable(summary)
	if opts.IsVerbose {
		renderVerboseDevNotes(summary.Categories)
	}
	printDevCleanOptimizationSuggestions()
}

func printDevCleanOptimizationSuggestions() {
	fmt.Println("  💡 Clean & Cache Optimization Suggestions:")
	fmt.Println("    • Preview space without deleting:    gitmap clear devtools --dry-run")
	fmt.Println("    • Clean specific ecosystems only:    gitmap clear devtools --only go,npm,pnpm -y")
	fmt.Println("    • Clear shell history & suggestions: gitmap clear terminal -y")
	fmt.Println("    • Clean Antigravity IDE caches:      gitmap agy clean-cache")
	fmt.Println()
}

func renderVerboseDevNotes(categories []osclean.CategoryCleanStats) {
	for _, cat := range categories {
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

func printOSDevCleanUsage() {
	fmt.Println("Usage: gitmap clear devtools [flags]")
	fmt.Println("       gitmap clear dev-tools [flags]")
	fmt.Println("       gitmap clear dev-tools-cache [flags]")
	fmt.Println("       gitmap devtools-cache clear [flags]")
	fmt.Println("       gitmap clean-dev [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -n, --dry-run      Preview space reclaimed without deleting")
	fmt.Println("  -y, --yes          Bypass confirmation prompt")
	fmt.Println("      --json         Output structured JSON summary")
	fmt.Println("  -v, --verbose      Display individual subpaths and command notes")
	fmt.Println("      --only <cats>  Comma-separated categories to clean (e.g. go,npm,pnpm)")
	fmt.Println("  -h, --help         Show this documentation")
	fmt.Println()
	printDevCleanOptimizationSuggestions()
}
