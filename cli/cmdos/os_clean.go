package cmdos

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/osclean"
)

func runOSClean(args []string) error {
	if len(args) > 0 && isDevTarget(args[0]) {
		return RunOSDevClean(args[1:])
	}

	if isHelpClean(args) {
		printOSCleanUsage()
		return nil
	}

	opts := parseCleanOptions(args)
	res := osclean.CleanTempDirectories(opts)
	if res.IsFailure() {
		return res.AppError()
	}

	printCleanSummary(res.Value, opts.IsDryRun)
	return nil
}

func isHelpClean(args []string) bool {
	if len(args) == 0 {
		return false
	}
	sub := strings.ToLower(args[0])
	return sub == "help" || isOSHelpArg(sub)
}

func parseCleanOptions(args []string) osclean.CleanOptions {
	return osclean.CleanOptions{
		IsDryRun:   hasFlag(args, "--dry-run") || hasFlag(args, "-n"),
		IsVerbose:  hasFlag(args, "--verbose") || hasFlag(args, "-v"),
		IsTempOnly: true,
	}
}

func printCleanSummary(stats osclean.CleanStats, isDryRun bool) {
	mbFreed := float64(stats.FreedBytes) / (1024 * 1024)
	if isDryRun {
		fmt.Printf("✔ [Dry-run] Would clean %d files, %d dirs (%.2f MB freed)\n",
			stats.RemovedFilesCount, stats.RemovedDirsCount, mbFreed)
		return
	}

	fmt.Printf("✔ Cleaned %d files, %d dirs (%.2f MB freed)\n",
		stats.RemovedFilesCount, stats.RemovedDirsCount, mbFreed)
}

func printOSCleanUsage() {
	fmt.Println("Usage: gitmap os clean [temp|dev] [flags]")
	fmt.Println("       gitmap os clear temp [flags]")
	fmt.Println("       gitmap os dev-clean [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -n, --dry-run    Simulate cleanup without deleting files")
	fmt.Println("  -v, --verbose    Show detailed cleanup output")
}
