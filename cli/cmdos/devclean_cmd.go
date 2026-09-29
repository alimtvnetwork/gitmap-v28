package cmdos

import (
	"fmt"
	"strings"
)

// RunDevToolCLI routes devtool and dt commands.
func RunDevToolCLI(args []string) error {
	if len(args) == 0 || isHelpDevTool(args[0]) {
		printDevToolUsage()
		return nil
	}

	subCmd := strings.ToLower(args[0])
	if subCmd == "clear" || subCmd == "clean" || subCmd == "cleanup" || subCmd == "cache" {
		return RunOSDevClean(args[1:])
	}

	return RunOSDevClean(args)
}

// RunDevToolClear runs the dev tool cache clearer directly.
func RunDevToolClear(args []string) error {
	return RunOSDevClean(args)
}

func isHelpDevTool(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

func printDevToolUsage() {
	fmt.Println("Usage: gitmap devtool clear [flags]")
	fmt.Println("       gitmap devtools clear [flags]")
	fmt.Println("       gitmap clear devtools [flags]")
	fmt.Println("       gitmap clear dev-tools [flags]")
	fmt.Println("       gitmap clear dev-tools-cache [flags]")
	fmt.Println("       gitmap devtools-cache clear [flags]")
	fmt.Println("       gitmap clean-dev [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  clear, clean       Scan and sweep 10 categories of developer caches")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -n, --dry-run      Preview space reclaimed without deleting")
	fmt.Println("  -y, --yes          Bypass confirmation prompt")
	fmt.Println("      --json         Output structured JSON summary")
	fmt.Println("  -v, --verbose      Display individual notes and warnings")
	fmt.Println("      --only <cats>  Comma-separated categories to clean (e.g. go,node,python)")
}
