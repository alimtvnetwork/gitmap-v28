package cmdos

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
)

type aiCleanOptions struct {
	isDryRun bool
	isYes    bool
	isJSON   bool
	isHelp   bool
}

type aiCleanJSONPayload struct {
	Status     string            `json:"status"`
	IsDryRun   bool              `json:"is_dry_run"`
	TotalFiles int               `json:"total_files"`
	TotalBytes int64             `json:"total_bytes"`
	HumanBytes string            `json:"human_bytes"`
	FreedFiles int               `json:"freed_files"`
	FreedBytes int64             `json:"freed_bytes"`
	Categories []AICleanCategory `json:"categories"`
}

var promptAICleanConfirmFn = defaultPromptAICleanConfirm

// RunOSAICleanCLI handles scanning and purging AI cache directories.
func RunOSAICleanCLI(args []string) error {
	opts := parseAICleanOptions(args)
	if opts.isHelp {
		printAICleanUsage()

		return nil
	}

	return dispatchAICleanMode(opts)
}

func dispatchAICleanMode(opts aiCleanOptions) error {
	cats := DiscoverAllAICleanTargets()
	files, bytes := calculateAICleanTotals(cats)
	if opts.isJSON {
		return handleAICleanJSON(cats, files, bytes, opts)
	}

	return executeAICleanTerminal(cats, files, bytes, opts)
}

func parseAICleanOptions(args []string) aiCleanOptions {
	isHelp := hasFlag(args, "--help") || hasFlag(args, "-h") || hasSubcommandArg(args, "help")
	isDryRun := hasFlag(args, "--dry-run") || hasFlag(args, "-n")
	isYes := hasFlag(args, "--yes") || hasFlag(args, "-y")
	isJSON := hasFlag(args, "--json")

	return aiCleanOptions{
		isDryRun: isDryRun,
		isYes:    isYes,
		isJSON:   isJSON,
		isHelp:   isHelp,
	}
}

func hasSubcommandArg(args []string, target string) bool {
	if len(args) == 0 {
		return false
	}

	return strings.ToLower(args[0]) == target
}

func calculateAICleanTotals(categories []AICleanCategory) (int, int64) {
	totalFiles := 0
	var totalBytes int64

	for _, cat := range categories {
		totalFiles += cat.FileCount
		totalBytes += cat.TotalBytes
	}

	return totalFiles, totalBytes
}

func handleAICleanJSON(cats []AICleanCategory, files int, bytes int64, opts aiCleanOptions) error {
	payload := buildAICleanJSONPayload(cats, files, bytes, opts)
	applyJSONPurge(&payload, cats, files, opts)

	return printJSONPayload(payload)
}

func applyJSONPurge(p *aiCleanJSONPayload, cats []AICleanCategory, files int, opts aiCleanOptions) {
	if !opts.isDryRun && files > 0 && opts.isYes {
		freedFiles, freedBytes := purgeAICleanFiles(cats)
		p.FreedFiles = freedFiles
		p.FreedBytes = freedBytes
		p.Status = "cleaned"
	}
}

func printJSONPayload(payload aiCleanJSONPayload) error {
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "json-marshal")
	}

	fmt.Println(string(data))

	return nil
}

func buildAICleanJSONPayload(cats []AICleanCategory, files int, bytes int64, opts aiCleanOptions) aiCleanJSONPayload {
	status := "clean"
	if files > 0 {
		status = "ready"
	}

	return aiCleanJSONPayload{
		Status:     status,
		IsDryRun:   opts.isDryRun,
		TotalFiles: files,
		TotalBytes: bytes,
		HumanBytes: cmddb.FormatBytes(bytes),
		Categories: cats,
	}
}

func executeAICleanTerminal(cats []AICleanCategory, files int, bytes int64, opts aiCleanOptions) error {
	fmt.Print(renderAICleanPreflightTable(cats, files, bytes))
	if files == 0 {
		fmt.Println("✔ No AI cache files found. System is clean.")

		return nil
	}
	if opts.isDryRun {
		fmt.Println("ℹ [dry-run] Preview mode only; no files removed.")

		return nil
	}

	return confirmAndPurgeAIClean(cats, opts)
}

func confirmAndPurgeAIClean(cats []AICleanCategory, opts aiCleanOptions) error {
	isAllowed, err := verifyUserConfirmation(opts.isYes)
	if err != nil {
		return apperror.WrapSimple(err, "confirm-prompt")
	}
	if !isAllowed {
		fmt.Println("Cleanup aborted.")

		return nil
	}

	return executePurgeAndReport(cats)
}

func executePurgeAndReport(cats []AICleanCategory) error {
	freedFiles, freedBytes := purgeAICleanFiles(cats)
	printAICleanSuccess(freedFiles, freedBytes)

	return nil
}

func verifyUserConfirmation(isYes bool) (bool, error) {
	if isYes {
		return true, nil
	}

	return promptAICleanConfirm("Proceed with AI cache cleanup? [y/N]: ")
}

func promptAICleanConfirm(msg string) (bool, error) {
	if promptAICleanConfirmFn != nil {
		return promptAICleanConfirmFn(msg)
	}

	return false, nil
}

func defaultPromptAICleanConfirm(msg string) (bool, error) {
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	trimmed := strings.ToLower(strings.TrimSpace(line))
	hasConfirmed := trimmed == "y" || trimmed == "yes"

	return hasConfirmed, nil
}

func purgeAICleanFiles(categories []AICleanCategory) (int, int64) {
	var totalFreedBytes int64
	freedCount := 0

	for _, cat := range categories {
		count, bytes := purgeCategoryFiles(cat.Paths)
		freedCount += count
		totalFreedBytes += bytes
	}

	return freedCount, totalFreedBytes
}

func purgeCategoryFiles(paths []string) (int, int64) {
	freedCount := 0
	var freedBytes int64

	for _, p := range paths {
		count, bytes := removeSingleFileSafely(p)
		freedCount += count
		freedBytes += bytes
	}

	return freedCount, freedBytes
}

func removeSingleFileSafely(filePath string) (int, int64) {
	if isProtectedFile(filePath) {
		return 0, 0
	}
	fi, err := os.Stat(filePath)
	if err != nil {
		return 0, 0
	}
	if removeErr := os.Remove(filePath); removeErr != nil {
		return 0, 0
	}

	return 1, fi.Size()
}

func printAICleanSuccess(freedCount int, freedBytes int64) {
	humanSize := cmddb.FormatBytes(freedBytes)
	fmt.Printf("✔ Successfully cleaned %d AI cache files (%s freed).\n", freedCount, humanSize)
}

func renderAICleanPreflightTable(cats []AICleanCategory, totalFiles int, totalBytes int64) string {
	var sb strings.Builder
	appendTableBorder(&sb)
	appendTableHeader(&sb)
	appendTableBorder(&sb)
	appendCategoryRows(&sb, cats)
	appendTableBorder(&sb)
	appendTableRow(&sb, "Total", totalFiles, cmddb.FormatBytes(totalBytes))
	appendTableBorder(&sb)

	return sb.String()
}

func appendCategoryRows(sb *strings.Builder, cats []AICleanCategory) {
	for _, cat := range cats {
		appendTableRow(sb, cat.Name, cat.FileCount, cmddb.FormatBytes(cat.TotalBytes))
	}
}

func appendTableBorder(sb *strings.Builder) {
	sb.WriteString("+------------------------------+-------------+-------------+\n")
}

func appendTableHeader(sb *strings.Builder) {
	sb.WriteString("| Category                     | Files       | Size        |\n")
}

func appendTableRow(sb *strings.Builder, name string, count int, humanSize string) {
	sb.WriteString(fmt.Sprintf("| %-28s | %11d | %11s |\n", name, count, humanSize))
}

const aiCleanUsageText = `Usage: gitmap os ai-clean [flags]
       gitmap ai-clean [flags]

Scan and purge Antigravity brain caches, system generated tasks,
OS temp AI dumps, and GitMap installer temp caches.

Flags:
  -n, --dry-run    Preview cache files to be deleted without deleting
  -y, --yes        Bypass interactive confirmation prompt
      --json       Output summary and discovered paths in JSON format
  -h, --help       Show this help message
`

func printAICleanUsage() {
	fmt.Print(aiCleanUsageText)
}
