// Package cmdagy — agy_running_prompts_render.go formats terminal output and help text for running prompts.
package cmdagy

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// TruncateWords trims text to maxWords and returns the truncated string and total word count.
func TruncateWords(text string, maxWords int) (string, int) {
	words := strings.Fields(text)
	total := len(words)
	if maxWords <= 0 || total <= maxWords {
		return text, total
	}
	return strings.Join(words[:maxWords], " ") + "...", total
}

func formatBytes(sizeBytes int64) string {
	if sizeBytes < 1024 {
		return fmt.Sprintf("%d B", sizeBytes)
	}
	kb := float64(sizeBytes) / 1024.0
	if kb < 1024.0 {
		return fmt.Sprintf("%.1f KB", kb)
	}
	return fmt.Sprintf("%.1f MB", kb/1024.0)
}

func countTotalBackedPrompts(batches []store.PromptBackupBatchRecord) int {
	total := 0
	for _, b := range batches {
		total += b.TotalPrompts
	}
	return total
}

// RenderBackupBatchesTable displays the backup registry in styled box and columns.
func RenderBackupBatchesTable(dbPath string, dbSize int64, batches []store.PromptBackupBatchRecord) {
	totalPrompts := countTotalBackedPrompts(batches)
	printBackupBoxHeader(dbPath, dbSize, len(batches), totalPrompts)
	printBackupTableHeaders()
	if len(batches) == 0 {
		fmt.Printf("  %sNo backup batches recorded.%s\n\n", constants.ColorDim, constants.ColorReset)
		return
	}
	for _, b := range batches {
		printBackupBatchRow(b)
	}
	fmt.Println()
}

func printBackupBoxHeader(dbPath string, dbSize int64, batchCount, promptCount int) {
	fmt.Println()
	fmt.Printf("  %s┌── Antigravity Running Prompts Backup Registry ────────────────────────┐%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s│%s Database: %-58s %s│%s\n", constants.ColorCyan, constants.ColorReset, fmt.Sprintf("%s (%s)", dbPath, formatBytes(dbSize)), constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s│%s Total Batches: %d · Total Backed Prompts: %-32d %s│%s\n", constants.ColorCyan, constants.ColorReset, batchCount, promptCount, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s└───%s────────────────────────────────────────────────────────────────────%s┘%s\n", constants.ColorCyan, constants.ColorCyan, constants.ColorCyan, constants.ColorReset)
}

func printBackupTableHeaders() {
	fmt.Printf("  %s%-14s %-21s %-6s %-8s %-9s %s%s\n", constants.ColorWhite, "Batch ID", "Created At", "Total", "Running", "Enqueued", "Status", constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("─", 72), constants.ColorReset)
}

func printBackupBatchRow(b store.PromptBackupBatchRecord) {
	statusColor := constants.ColorGreen
	if b.Status == "restored" {
		statusColor = constants.ColorYellow
	}
	created := b.CreatedAt
	if len(created) > 19 {
		created = created[:19]
	}
	fmt.Printf("  %-14s %-21s %-6d %-8d %-9d %s%s%s\n", b.BatchID, created, b.TotalPrompts, b.RunningCount, b.EnqueuedCount, statusColor, b.Status, constants.ColorReset)
}

// RenderRunningPromptsTable prints active and enqueued prompt items to the console.
func RenderRunningPromptsTable(items []store.RunningPromptRecord, isFull bool) {
	fmt.Println()
	fmt.Printf("  %s● Antigravity Running & Queued Prompts (%d items)%s\n", constants.ColorCyan, len(items), constants.ColorReset)
	fmt.Printf("  %s%-20s %-10s %-8s %s%s\n", constants.ColorWhite, "PROJECT", "STATUS", "WORDS", "PROMPT", constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("─", 80), constants.ColorReset)
	if len(items) == 0 {
		fmt.Printf("  %sNo active or queued prompts found.%s\n\n", constants.ColorDim, constants.ColorReset)
		return
	}
	for _, it := range items {
		printRunningPromptRow(it, isFull)
	}
	fmt.Println()
}

func printRunningPromptRow(it store.RunningPromptRecord, isFull bool) {
	statusColor := constants.ColorGreen
	if it.Status == store.PromptStatusEnqueued {
		statusColor = constants.ColorYellow
	}
	proj := it.ProjectName
	if len(proj) > 20 {
		proj = proj[:17] + "..."
	}
	text := it.Snippet
	if isFull {
		text = it.Prompt
	}
	fmt.Printf("  %-20s %s%-10s%s %-8d %s\n", proj, statusColor, strings.ToUpper(string(it.Status)), constants.ColorReset, it.WordCount, text)
}

// RenderRunningPromptsHelp prints synopsis and usage guidelines for running-prompts commands.
func RenderRunningPromptsHelp() {
	printRunningPromptsHelpUsage()
	printRunningPromptsHelpCommands()
	printRunningPromptsHelpFlags()
}

func printRunningPromptsHelpUsage() {
	fmt.Printf(`%sAntigravity Running Prompts Management%s

Usage:
  gitmap agy running-prompts <command> [flags]

`, constants.ColorCyan, constants.ColorReset)
}

func printRunningPromptsHelpCommands() {
	fmt.Printf(`Commands:
  backup   Snapshot running and queued prompts to Split-DB
  restore  Restore backed-up prompts into project queue
  clean    Prune old restored backup batches
  ls       Inspect and list running/queued prompts
  export   Export prompts to SQLite .db or .json
  import   Import prompts from SQLite .db or .json

`)
}

func printRunningPromptsHelpFlags() {
	fmt.Printf(`Flags:
  -l, --limit int       Limit output entries (default 8)
      --wc int          Word count truncation limit (default 100)
      --full            Show full prompt without word truncation
      --json            Output in JSON format
  -f, --file string     Custom backup file path
  -k, --keep            Keep restored batch from retention pruning
      --force           Force delete all backup batches during clean
`)
}

// RenderBackupHelp prints usage guidelines for backup-running-prompts.
func RenderBackupHelp() {
	fmt.Printf(`%sAntigravity Running Prompts Backup%s

Usage:
  gitmap agy backup-running-prompts [flags]
  gitmap agy backup-running-prompts ls [flags]

Flags:
  -f, --file string   Target backup database path (default: data/backup-prompts/sql.db)
      --json          Output backup summary or batches in JSON format
`, constants.ColorCyan, constants.ColorReset)
}

// RenderRestoreHelp prints usage guidelines for restore-running-prompts.
func RenderRestoreHelp() {
	fmt.Printf(`%sAntigravity Running Prompts Restore%s

Usage:
  gitmap agy restore-running-prompts [flags]
  gitmap agy running-prompts restore [flags]

Flags:
  -k, --keep          Exclude restored batch from 1-day auto-pruning
  -f, --file string   Source database file path
      --json          Output restoration result in JSON format
`, constants.ColorCyan, constants.ColorReset)
}
