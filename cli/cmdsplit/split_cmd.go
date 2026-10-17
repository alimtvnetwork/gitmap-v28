// NOTE: Needs refactor — exceeds 300-line target
package cmdsplit


import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

type splitFlags struct {
	dryRun      bool
	autoConfirm bool
	force       bool
	sizeRaw     string
	threshold   int64
}

type largeFileInfo struct {
	RelPath           string
	Size              int64
	Method            uint16
	MethodName        string
	ProjectedChunks   int
	ProjectedPartList []string
}

// RunSplitLargeFilesCmd executes the split-large-files command.
func RunSplitLargeFilesCmd(args []string) error {
	cleanArgs := normalizeSplitArgs(args)
	if hasHelpFlag(cleanArgs) {
		helpdoc.PrintHelpAndExit("split-large-files", cleanArgs)
		return nil
	}

	flags, err := parseSplitFlags(cleanArgs)
	if err != nil {
		return err
	}

	if !isInsideGitRepo() {
		return apperror.NewValidation("split_large_files", "E9001", "not inside a Git repository")
	}

	return processSplitOperation(flags)
}

func normalizeSplitArgs(args []string) []string {
	if len(args) >= 2 && args[0] == "large" && args[1] == "files" {
		return args[2:]
	}

	if len(args) >= 1 && (args[0] == "large-files" || args[0] == "large_files") {
		return args[1:]
	}

	return args
}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "help" {
			return true
		}
	}

	return false
}

func parseSplitFlags(args []string) (splitFlags, error) {
	var sf splitFlags
	fs := flag.NewFlagSet("split-large-files", flag.ContinueOnError)

	fs.BoolVar(&sf.dryRun, "dry-run", false, "Preview without modifying files or git")
	fs.BoolVar(&sf.autoConfirm, "y", false, "Auto-confirm execution prompt")
	fs.BoolVar(&sf.autoConfirm, "yes", false, "Auto-confirm execution prompt")
	fs.BoolVar(&sf.force, "force", false, "Force split execution")
	fs.StringVar(&sf.sizeRaw, "size", "10MB", "Size threshold for splitting")
	fs.StringVar(&sf.sizeRaw, "s", "10MB", "Size threshold for splitting")

	if err := fs.Parse(args); err != nil {
		return sf, apperror.WrapSimple(err, "parse_split_flags")
	}

	threshold, err := parseByteSize(sf.sizeRaw)
	sf.threshold = threshold

	return sf, err
}

func parseByteSize(val string) (int64, error) {
	trimmed := strings.ToUpper(strings.TrimSpace(val))
	if trimmed == "" {
		return 10 * 1024 * 1024, nil
	}

	multiplier, numStr := resolveSizeMultiplier(trimmed)
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, apperror.WrapSimple(err, "invalid_size_format")
	}

	return int64(num * float64(multiplier)), nil
}

func resolveSizeMultiplier(s string) (int64, string) {
	if strings.HasSuffix(s, "GB") {
		return 1024 * 1024 * 1024, strings.TrimSuffix(s, "GB")
	}
	if strings.HasSuffix(s, "G") {
		return 1024 * 1024 * 1024, strings.TrimSuffix(s, "G")
	}
	if strings.HasSuffix(s, "MB") {
		return 1024 * 1024, strings.TrimSuffix(s, "MB")
	}
	if strings.HasSuffix(s, "M") {
		return 1024 * 1024, strings.TrimSuffix(s, "M")
	}
	if strings.HasSuffix(s, "KB") {
		return 1024, strings.TrimSuffix(s, "KB")
	}
	if strings.HasSuffix(s, "K") {
		return 1024, strings.TrimSuffix(s, "K")
	}

	return 1, strings.TrimSuffix(s, "B")
}

func formatByteSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func isInsideGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}

func processSplitOperation(flags splitFlags) error {
	candidates, err := findLargeTrackedFiles(flags.threshold)
	if err != nil {
		return err
	}

	if len(candidates) == 0 {
		fmt.Printf("  %s✓ No files exceeding %s found in repository.%s\n\n",
			constants.ColorGreen, formatByteSize(flags.threshold), constants.ColorReset)
		return nil
	}

	if flags.dryRun {
		return executeSplitDryRun(candidates, flags.threshold)
	}

	return executeSplitActive(candidates, flags)
}

func findLargeTrackedFiles(threshold int64) ([]largeFileInfo, error) {
	cmd := exec.Command("git", "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "ls_files")
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var results []largeFileInfo

	for _, line := range lines {
		relPath := strings.TrimSpace(line)
		if item, ok := inspectLargeFile(relPath, threshold); ok {
			results = append(results, item)
		}
	}

	return results, nil
}

func inspectLargeFile(relPath string, threshold int64) (largeFileInfo, bool) {
	var empty largeFileInfo
	if relPath == "" || isMultipartZipPart(relPath) {
		return empty, false
	}

	fi, err := os.Stat(relPath)
	if err != nil || fi.IsDir() || fi.Size() <= threshold {
		return empty, false
	}

	return buildLargeFileInfo(relPath, fi.Size(), threshold), true
}

func isMultipartZipPart(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if len(ext) >= 3 && ext[1] == 'z' && ext[2] >= '0' && ext[2] <= '9' {
		return true
	}

	return false
}

func buildLargeFileInfo(relPath string, size, threshold int64) largeFileInfo {
	method := ResolveCompressionMethod(relPath)
	chunks := int((size + threshold - 1) / threshold)
	parts := make([]string, chunks)

	for i := 0; i < chunks; i++ {
		parts[i] = filepath.Base(formatPartName(relPath, i))
	}

	return largeFileInfo{
		RelPath:           relPath,
		Size:              size,
		Method:            method,
		MethodName:        CompressionMethodName(method),
		ProjectedChunks:   chunks,
		ProjectedPartList: parts,
	}
}

func executeSplitDryRun(files []largeFileInfo, threshold int64) error {
	fmt.Printf("\n  %s■ gitmap split-large-files — [dry-run] previewing files > %s%s\n\n",
		constants.ColorCyan, formatByteSize(threshold), constants.ColorReset)

	cfg := buildSplitTableConfig(files)
	termout.PrintTable(cfg)

	totalSize := sumLargeFilesSize(files)
	fmt.Printf("\n  %sSummary:%s Found %d file(s) exceeding %s (Total: %s)\n",
		constants.ColorCyan, constants.ColorReset, len(files), formatByteSize(threshold), formatByteSize(totalSize))
	fmt.Printf("  %s[dry-run]%s No files were split and Git index was not modified.\n\n",
		constants.ColorYellow, constants.ColorReset)

	return nil
}

func sumLargeFilesSize(files []largeFileInfo) int64 {
	var total int64
	for _, f := range files {
		total += f.Size
	}

	return total
}

func buildSplitTableConfig(files []largeFileInfo) termout.TableConfig {
	cols := []termout.Column{
		{Title: "File", Align: termout.AlignLeft, MinWidth: 20},
		{Title: "Size", Align: termout.AlignRight, MinWidth: 10},
		{Title: "Method", Align: termout.AlignCenter, MinWidth: 8},
		{Title: "Chunks", Align: termout.AlignRight, MinWidth: 8},
		{Title: "Projected Parts", Align: termout.AlignLeft, MinWidth: 25},
	}

	rows := make([]termout.Row, len(files))
	for i, f := range files {
		partsSummary := strings.Join(f.ProjectedPartList, ", ")
		rows[i] = termout.Row{
			Cells: []string{f.RelPath, formatByteSize(f.Size), f.MethodName, strconv.Itoa(f.ProjectedChunks), partsSummary},
			Color: constants.ColorWhite,
		}
	}

	return termout.TableConfig{Columns: cols, Rows: rows, HeaderColor: constants.ColorYellow}
}

func executeSplitActive(candidates []largeFileInfo, flags splitFlags) error {
	confirmed, err := confirmIrreversibleSplit(flags.autoConfirm || flags.force)
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Printf("  %sOperation cancelled by user.%s\n\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}

	for _, file := range candidates {
		if err := splitAndStageFile(file, flags.threshold); err != nil {
			return err
		}
	}

	printSplitSuccessSummary(len(candidates))

	return nil
}

func confirmIrreversibleSplit(autoConfirm bool) (bool, error) {
	fmt.Println()
	fmt.Printf("  %sWARNING: This operation rewrites Git history and splits files into multipart zip archives. It is irreversible.%s\n",
		constants.ColorYellow, constants.ColorReset)

	if autoConfirm {
		fmt.Println()
		return true, nil
	}

	fmt.Print("  Do you want to proceed? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return false, apperror.WrapSimple(err, "read_confirmation")
	}

	val := strings.TrimSpace(strings.ToLower(text))

	return val == "y" || val == "yes", nil
}

func splitAndStageFile(file largeFileInfo, threshold int64) error {
	res, err := CreateMultipartZip(file.RelPath, threshold)
	if err != nil {
		return err
	}

	if err := os.Remove(file.RelPath); err != nil {
		return apperror.WrapSimple(err, "remove_original_file")
	}

	if err := stageGitSplitChanges(file.RelPath, res.ChunkFiles); err != nil {
		return err
	}

	fmt.Printf("  %s✓ Split%s %s (%s) → %d chunks [%s]\n",
		constants.ColorGreen, constants.ColorReset, file.RelPath, formatByteSize(file.Size), len(res.ChunkFiles), file.MethodName)

	return nil
}

func stageGitSplitChanges(origPath string, chunkFiles []string) error {
	rmCmd := exec.Command("git", "rm", "--cached", "--quiet", "--", origPath)
	_ = rmCmd.Run()

	for _, chunk := range chunkFiles {
		addCmd := exec.Command("git", "add", "--", chunk)
		if err := addCmd.Run(); err != nil {
			return apperror.WrapSimple(err, "git_add_chunk")
		}
	}

	return nil
}

func printSplitSuccessSummary(count int) {
	fmt.Println()
	fmt.Printf("  %s✓ Split complete:%s %d file(s) split into multipart zip archives.\n",
		constants.ColorGreen, constants.ColorReset, count)
	fmt.Println("  Original files removed from Git and multipart archives staged.")
	fmt.Println()
}
