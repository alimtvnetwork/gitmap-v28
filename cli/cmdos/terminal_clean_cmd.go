package cmdos

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/osclean"
)

// RunTerminalCleanCLI dispatches terminal history cleanup and suggestion reseeding.
func RunTerminalCleanCLI(args []string) error {
	if isHelpTerminalClean(args) {
		printTerminalCleanUsage()
		return nil
	}

	opts := parseTerminalCleanOptions(args)
	if isTerminalCleanPromptRequired(opts) && !confirmTerminalClean() {
		fmt.Println("  Aborted by operator.")
		return nil
	}

	res := osclean.CleanTerminalHistory(opts)
	if res.IsFailure() {
		return res.AppError()
	}

	renderTerminalCleanOutput(res.Value, opts)
	return nil
}

func isHelpTerminalClean(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "-h" || low == "--help" || low == "help" || low == "/?" {
			return true
		}
	}
	return false
}

func parseTerminalCleanOptions(args []string) osclean.TerminalCleanOptions {
	opts := osclean.TerminalCleanOptions{
		ShouldReseed: true,
	}

	for i := 0; i < len(args); i++ {
		a := strings.ToLower(args[i])
		parseTerminalCleanFlag(a, args, &i, &opts)
	}

	return opts
}

func parseTerminalCleanFlag(a string, args []string, i *int, opts *osclean.TerminalCleanOptions) {
	switch {
	case a == "--dry-run" || a == "-n" || a == "-d" || a == "dry-run":
		opts.IsDryRun = true
	case a == "--yes" || a == "-y" || a == "/y" || a == "yes" || a == "--force" || a == "-f":
		opts.HasAutoYes = true
	case a == "--json":
		opts.IsJSON = true
	case a == "--verbose" || a == "-v":
		opts.IsVerbose = true
	case a == "--no-reseed" || a == "--skip-reseed":
		opts.ShouldReseed = false
	case strings.HasPrefix(a, "--only="):
		opts.OnlyShells = strings.Split(strings.TrimPrefix(a, "--only="), ",")
	case a == "--only" && *i+1 < len(args):
		*i++
		opts.OnlyShells = strings.Split(args[*i], ",")
	case a == "terminal" || a == "clear" || a == "clean" || a == "cleanup" || a == "history" ||
		a == "clear-terminal" || a == "clean-terminal" || a == "terminal-clear" || a == "terminal-clean":
		// Ignore action/noun arguments
	}
}

func isTerminalCleanPromptRequired(opts osclean.TerminalCleanOptions) bool {
	return !opts.IsDryRun && !opts.HasAutoYes && !opts.IsJSON
}

func confirmTerminalClean() bool {
	fmt.Print("  Proceed with terminal history cleanup & suggestion reseeding? Type 'yes' to continue: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	return strings.ToLower(strings.TrimSpace(text)) == "yes"
}

func renderTerminalCleanOutput(summary osclean.TerminalCleanSummary, opts osclean.TerminalCleanOptions) {
	if opts.IsJSON {
		data, _ := json.MarshalIndent(summary, "", "  ")
		fmt.Println(string(data))
		return
	}

	osclean.RenderTerminalSummaryTable(summary)
	if opts.IsVerbose {
		osclean.RenderVerboseTerminalNotes(summary.Shells)
	}
}

func printTerminalCleanUsage() {
	fmt.Println("Usage: gitmap clear terminal [flags]")
	fmt.Println("       gitmap clear-terminal [flags]")
	fmt.Println("       gitmap clean terminal [flags]")
	fmt.Println("       gitmap clean-terminal [flags]")
	fmt.Println("       gitmap terminal clear [flags]")
	fmt.Println()
	fmt.Println("Description:")
	fmt.Println("  Clears PowerShell, Bash, Zsh, Sh, and Fish command histories,")
	fmt.Println("  removes old suggestions clutter, and reseeds fresh GitMap suggestions.")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -n, --dry-run      Preview history files and space without deleting")
	fmt.Println("  -y, --yes, -f      Bypass confirmation prompt")
	fmt.Println("      --json         Output structured JSON telemetry")
	fmt.Println("  -v, --verbose      Show individual history file paths")
	fmt.Println("      --no-reseed    Clear history only without re-seeding GitMap suggestions")
	fmt.Println("      --only <list>  Comma-separated shell targets (e.g. powershell,bash,zsh)")
	fmt.Println("  -h, --help         Show this documentation")
}
