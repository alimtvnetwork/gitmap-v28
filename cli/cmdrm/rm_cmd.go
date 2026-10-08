package cmdrm

import (
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunSafeRmCLI is the main entry point for task-based safe file removal and rollback.
func RunSafeRmCLI(args []string) error {
	if len(args) == 0 {
		printSafeRmUsage()
		return apperror.NewValidationError("target file pattern required")
	}

	mode, remaining := parseSafeRmSubcommand(args)
	switch mode {
	case "undo":
		return handleSafeRmUndo(remaining)
	case "list":
		return handleSafeRmList()
	case "purge":
		return handleSafeRmPurge(remaining)
	default:
		return handleSafeRmExecute(remaining)
	}
}

func parseSafeRmSubcommand(args []string) (string, []string) {
	first := strings.ToLower(args[0])
	if first == "undo" || first == "--undo" || first == "-u" {
		return "undo", args[1:]
	}
	if first == "list" || first == "--list" {
		return "list", args[1:]
	}
	if first == "purge" || first == "--purge" {
		return "purge", args[1:]
	}

	for _, a := range args {
		if a == "--undo" || a == "-u" {
			return "undo", args
		}
	}
	return "remove", args
}

func handleSafeRmUndo(args []string) error {
	taskID := ""
	for i := 0; i < len(args); i++ {
		a := args[i]
		if (a == "--task" || a == "-t") && i+1 < len(args) {
			taskID = args[i+1]
			i++
			continue
		}
		if !strings.HasPrefix(a, "-") && taskID == "" {
			taskID = a
		}
	}

	opts := RmUndoOptions{TaskId: taskID}
	if err := ExecuteUndo(opts); err != nil {
		return err
	}

	target := taskID
	if target == "" {
		target = "latest session"
	}
	fmt.Printf("%s✔ Restored files successfully for %s%s\n", constants.ColorGreen, target, constants.ColorReset)
	return nil
}

func handleSafeRmList() error {
	manifests, err := ListRmBackups()
	if err != nil {
		return err
	}

	if len(manifests) == 0 {
		fmt.Println("No active recoverable file removal sessions in OS temporary storage.")
		return nil
	}

	fmt.Printf("%-24s  %-20s  %-6s  %-10s  %s\n", "TASK ID", "CREATED AT", "FILES", "SIZE", "REASON")
	fmt.Println(strings.Repeat("-", 78))
	for _, m := range manifests {
		created := m.CreatedAt.Format("2006-01-02 15:04:05")
		sizeStr := formatByteSize(m.TotalSizeBytes)
		fmt.Printf("%-24s  %-20s  %-6d  %-10s  %s\n", m.TaskId, created, m.FileCount, sizeStr, m.Reason)
	}
	return nil
}

func handleSafeRmPurge(args []string) error {
	opts := RmPurgeOptions{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--all" {
			opts.All = true
			continue
		}
		if (a == "--task" || a == "-t") && i+1 < len(args) {
			opts.TaskId = args[i+1]
			i++
			continue
		}
		if (a == "--older-than") && i+1 < len(args) {
			opts.OlderThan = parseOlderThanDuration(args[i+1])
			i++
			continue
		}
	}

	count, err := PurgeRmBackups(opts)
	if err != nil {
		return err
	}

	fmt.Printf("%s✔ Purged %d backup session(s) from temporary storage%s\n", constants.ColorGreen, count, constants.ColorReset)
	return nil
}

func handleSafeRmExecute(args []string) error {
	opts, err := parseExecuteOptions(args)
	if err != nil {
		return err
	}

	manifest, err := StageAndRemoveFiles(opts)
	if err != nil {
		return err
	}

	if opts.DryRun {
		printDryRunManifest(manifest)
		return nil
	}

	printExecutedManifest(manifest)
	return nil
}

func parseExecuteOptions(args []string) (RmOptions, error) {
	opts := RmOptions{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--dry-run" {
			opts.DryRun = true
			continue
		}
		if a == "--force" || a == "-f" {
			opts.Force = true
			continue
		}
		if (a == "--task" || a == "-t") && i+1 < len(args) {
			opts.TaskId = args[i+1]
			i++
			continue
		}
		if (a == "--reason" || a == "-r") && i+1 < len(args) {
			opts.Reason = args[i+1]
			i++
			continue
		}
		if !strings.HasPrefix(a, "-") {
			opts.Patterns = append(opts.Patterns, a)
		}
	}

	if len(opts.Patterns) == 0 {
		return opts, apperror.NewValidationError("target file pattern required")
	}
	return opts, nil
}

func printDryRunManifest(manifest *RmManifest) {
	fmt.Printf("%s[DRY RUN] Would safely remove %d file(s) for task '%s':%s\n", constants.ColorYellow, manifest.FileCount, manifest.TaskId, constants.ColorReset)
	for _, f := range manifest.Files {
		fmt.Printf("  - %s (%s)\n", f.RelativePath, formatByteSize(f.SizeBytes))
	}
}

func printExecutedManifest(manifest *RmManifest) {
	fmt.Printf("%s✔ Safely removed %d file(s) for task '%s'%s\n", constants.ColorGreen, manifest.FileCount, manifest.TaskId, constants.ColorReset)
	for _, f := range manifest.Files {
		fmt.Printf("  - %s (%s)\n", f.RelativePath, formatByteSize(f.SizeBytes))
	}
	fmt.Printf("  Run 'gitmap rm undo %s' to restore.\n", manifest.TaskId)
}

func formatByteSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024.0)
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/(1024.0*1024.0))
}

func printSafeRmUsage() {
	fmt.Println("Usage: gitmap rm <patterns...> [--task <id>] [--reason <text>] [--dry-run] [--force]")
	fmt.Println("       gitmap rm undo [<task_id>]")
	fmt.Println("       gitmap rm list")
	fmt.Println("       gitmap rm purge [--all] [--task <id>] [--older-than <duration>]")
}

func parseOlderThanDuration(val string) time.Duration {
	d, err := time.ParseDuration(val)
	if err != nil {
		return 0
	}
	return d
}
