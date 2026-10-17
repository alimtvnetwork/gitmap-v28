package cmdsplit

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

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
		fmt.Printf("  %sOperation canceled by user.%s\n\n", constants.ColorYellow, constants.ColorReset)
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
