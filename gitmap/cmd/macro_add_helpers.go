// Package cmd — macro_add_helpers.go: in-builder helper utilities for interactive macro creation.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/uipref"
)

type interactiveSessionState struct {
	lastInspectedCmd string
	isExecEnabled    bool
}

func newInteractiveState() *interactiveSessionState {
	return &interactiveSessionState{
		isExecEnabled: true,
	}
}

func isInteractiveHelper(line string) bool {
	trimmed := strings.TrimSpace(line)

	return isExactHelper(trimmed) || hasPrefixHelper(trimmed)
}

func isExactHelper(line string) bool {
	switch strings.ToLower(line) {
	case "ls", "dir", ":ls", ":dir", "pwd", ":pwd", "help", ":help", "?", "+add", "exec", ":exec":
		return true
	default:
		return false
	}
}

func hasPrefixHelper(line string) bool {
	prefixes := []string{
		"find ", ":find ", "search ", ":search ", "grep ", ":grep ",
		"replace ", ":replace ", "cd ", ":cd ", "pwd ", ":pwd ", "add ",
		"ls ", "dir ", ":ls ", ":dir ", "exec ", ":exec ",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(strings.ToLower(line), p) {
			return true
		}
	}

	return false
}

func handleInteractiveHelper(line string, state *interactiveSessionState, steps *[]macro.MacroStep, stepNum *int) bool {
	trimmed := strings.TrimSpace(line)
	if handleAddAction(trimmed, state, steps, stepNum) {
		return true
	}
	if handleExecToggle(trimmed, state) {
		return true
	}

	if handleNavigationOrInspection(trimmed, state, steps, stepNum) {
		return true
	}

	return handleSearchOrReplace(trimmed, state)
}

func handleExecToggle(line string, state *interactiveSessionState) bool {
	low := strings.ToLower(line)
	if !strings.HasPrefix(low, "exec") && !strings.HasPrefix(low, ":exec") {
		return false
	}
	arg := strings.ToLower(extractCommandArgument(line))
	if arg == "off" {
		state.isExecEnabled = false
		fmt.Printf("  %s✓ Live command execution disabled.%s\n\n", constants.ColorYellow, constants.ColorReset)

		return true
	}

	return setExecOnOrPrint(arg, state)
}

func setExecOnOrPrint(arg string, state *interactiveSessionState) bool {
	if arg == "on" {
		state.isExecEnabled = true
		fmt.Printf("  %s✓ Live command execution enabled.%s\n\n", constants.ColorGreen, constants.ColorReset)

		return true
	}
	printExecStatus(state.isExecEnabled)

	return true
}

func printExecStatus(isExecEnabled bool) {
	status := "enabled"
	if !isExecEnabled {
		status = "disabled"
	}
	fmt.Printf("  %sLive execution: %s (toggle: 'exec on' / 'exec off').%s\n\n",
		constants.ColorCyan, status, constants.ColorReset)
}

func handleNavigationOrInspection(line string, state *interactiveSessionState, steps *[]macro.MacroStep, stepNum *int) bool {
	if isLsOrDirCmd(line) {
		return handleLsOrDir(line, state)
	}

	if isPwdCmd(line) {
		return handlePwdCmd(line)
	}

	if isCdCmd(line) {
		return handleCdCmd(line, steps, stepNum)
	}

	return checkHelpCmd(line)
}

func checkHelpCmd(line string) bool {
	if isHelpCmd(line) {
		printInteractiveHelp()

		return true
	}

	return false
}

func isLsOrDirCmd(line string) bool {
	low := strings.ToLower(line)

	return low == "ls" || low == "dir" || low == ":ls" || low == ":dir" ||
		strings.HasPrefix(low, "ls ") || strings.HasPrefix(low, "dir ") ||
		strings.HasPrefix(low, ":ls ") || strings.HasPrefix(low, ":dir ")
}

func isPwdCmd(line string) bool {
	low := strings.ToLower(line)

	return low == "pwd" || low == ":pwd" || strings.HasPrefix(low, "pwd ") || strings.HasPrefix(low, ":pwd ")
}

func isCdCmd(line string) bool {
	low := strings.ToLower(line)

	return strings.HasPrefix(low, "cd ") || strings.HasPrefix(low, ":cd ")
}

func isHelpCmd(line string) bool {
	low := strings.ToLower(line)

	return low == "help" || low == ":help" || low == "?"
}

func handleLsOrDir(line string, state *interactiveSessionState) bool {
	state.lastInspectedCmd = line
	parts := strings.Fields(line)
	targetDir := "."
	if len(parts) > 1 && !strings.HasPrefix(parts[1], "-") {
		targetDir = parts[1]
	}

	_ = executeInteractiveLs(targetDir)

	return true
}

func handlePwdCmd(line string) bool {
	low := strings.ToLower(strings.TrimSpace(line))
	if low == "pwd on" || low == ":pwd on" {
		_ = uipref.SetMacroShowPwd(true)
		fmt.Printf("  %s✓ Macro PWD display enabled.%s\n\n", constants.ColorGreen, constants.ColorReset)

		return true
	}

	if low == "pwd off" || low == ":pwd off" {
		_ = uipref.SetMacroShowPwd(false)
		fmt.Printf("  %s✓ Macro PWD display disabled.%s\n\n", constants.ColorYellow, constants.ColorReset)

		return true
	}

	if low == ":pwd" {
		toggleMacroPwd()

		return true
	}

	printCurrentPwdBanner()

	return true
}

func toggleMacroPwd() {
	newState := !uipref.IsMacroPwdVisible()
	_ = uipref.SetMacroShowPwd(newState)
	statusText := "enabled"
	if !newState {
		statusText = "disabled"
	}
	fmt.Printf("  %s✓ Macro PWD display %s.%s\n\n", constants.ColorCyan, statusText, constants.ColorReset)
}

func printCurrentPwdBanner() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("  %s▲ Could not determine PWD: %v%s\n\n", constants.ColorRed, err, constants.ColorReset)

		return
	}

	fmt.Printf("  %s[PWD: %s]%s\n\n", constants.ColorCyan, cwd, constants.ColorReset)
}

func handleCdCmd(line string, steps *[]macro.MacroStep, stepNum *int) bool {
	target := extractCommandArgument(line)
	if target == "" {
		fmt.Printf("  %s▲ Usage: cd <directory>%s\n\n", constants.ColorYellow, constants.ColorReset)

		return true
	}
	if err := executeInteractiveCd(target); err == nil && steps != nil && stepNum != nil {
		recordExplicitCommand(line, steps, stepNum)
	}

	return true
}

func executeInteractiveCd(rawTarget string) error {
	cwd, _ := os.Getwd()
	normTarget := macro.NormalizeTargetPath(rawTarget, cwd)
	if err := os.Chdir(normTarget); err != nil {
		fmt.Printf("  %s▲ cd %s: %v%s\n\n", constants.ColorRed, rawTarget, err, constants.ColorReset)

		return apperror.WrapSimple(err, "change directory")
	}

	newCwd, _ := os.Getwd()
	fmt.Printf("  %s✓ Changed directory to: %s%s\n\n", constants.ColorGreen, newCwd, constants.ColorReset)

	return nil
}

func handleSearchOrReplace(line string, state *interactiveSessionState) bool {
	low := strings.ToLower(line)
	if strings.HasPrefix(low, "find ") || strings.HasPrefix(low, ":find ") {
		return handleFindCmd(line, state)
	}

	if strings.HasPrefix(low, "search ") || strings.HasPrefix(low, ":search ") ||
		strings.HasPrefix(low, "grep ") || strings.HasPrefix(low, ":grep ") {
		return handleSearchCmd(line, state)
	}

	if strings.HasPrefix(low, "replace ") || strings.HasPrefix(low, ":replace ") {
		return handleReplaceCmd(line)
	}

	return false
}

func handleFindCmd(line string, state *interactiveSessionState) bool {
	state.lastInspectedCmd = line
	pattern := extractCommandArgument(line)
	if pattern == "" {
		fmt.Printf("  %s▲ Usage: find <filename-pattern>%s\n\n", constants.ColorYellow, constants.ColorReset)

		return true
	}

	_ = executeInteractiveFind(pattern)

	return true
}

func handleSearchCmd(line string, state *interactiveSessionState) bool {
	state.lastInspectedCmd = line
	query := extractCommandArgument(line)
	if query == "" {
		fmt.Printf("  %s▲ Usage: search <query-text>%s\n\n", constants.ColorYellow, constants.ColorReset)

		return true
	}

	_ = executeInteractiveSearch(query)

	return true
}

func handleReplaceCmd(line string) bool {
	rawArgs := extractCommandArgument(line)
	if rawArgs == "" {
		fmt.Printf("  %s▲ Usage: replace <old> <new> [glob]%s\n\n", constants.ColorYellow, constants.ColorReset)

		return true
	}

	_ = executeInteractiveReplace(rawArgs)

	return true
}

func extractCommandArgument(line string) string {
	idx := strings.Index(line, " ")
	if idx < 0 {
		return ""
	}

	return strings.TrimSpace(line[idx+1:])
}

func handleAddAction(line string, state *interactiveSessionState, steps *[]macro.MacroStep, stepNum *int) bool {
	if strings.EqualFold(line, "+add") {
		return recordLastInspectedCommand(state, steps, stepNum)
	}

	if strings.HasPrefix(strings.ToLower(line), "add ") {
		cmdText := strings.TrimSpace(line[4:])
		recordExplicitCommand(cmdText, steps, stepNum)

		return true
	}

	return false
}

func recordLastInspectedCommand(state *interactiveSessionState, steps *[]macro.MacroStep, stepNum *int) bool {
	if state.lastInspectedCmd == "" {
		fmt.Printf("  %s▲ No inspected command to add.%s\n\n", constants.ColorYellow, constants.ColorReset)

		return true
	}

	cmdToAdd := state.lastInspectedCmd
	state.lastInspectedCmd = ""
	recordExplicitCommand(cmdToAdd, steps, stepNum)

	return true
}

func recordExplicitCommand(cmdText string, steps *[]macro.MacroStep, stepNum *int) {
	*steps = append(*steps, makeMacroStep(*stepNum, cmdText))
	fmt.Printf("  %s✓ Recorded step %d: %s%s\n\n", constants.ColorGreen, *stepNum, cmdText, constants.ColorReset)
	*stepNum++
}

func executeInteractiveLs(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("  %s▲ Cannot read directory %q: %v%s\n\n", constants.ColorRed, dirPath, err, constants.ColorReset)

		return apperror.WrapSimple(err, "read directory for inspection")
	}

	printLsHeader(dirPath)
	dirCount, fileCount, totalSize := printLsEntriesList(dirPath, entries)
	printLsSummary(dirCount, fileCount, totalSize)

	return nil
}

func printLsHeader(dirPath string) {
	absPath, _ := filepath.Abs(dirPath)
	fmt.Println()
	fmt.Printf("  %sDirectory listing: %s%s%s\n", constants.ColorCyan, constants.ColorWhite, absPath, constants.ColorReset)
	fmt.Printf("  %s------------------------------------------------------------------------%s\n", constants.ColorDim, constants.ColorReset)
}

func printLsEntriesList(dirPath string, entries []os.DirEntry) (int, int, int64) {
	dirCount, fileCount := 0, 0
	var totalSize int64

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if entry.IsDir() {
			dirCount++
			printDirectoryEntry(entry.Name(), info.ModTime())
			continue
		}

		fileCount++
		totalSize += info.Size()
		printFileEntry(entry.Name(), info.Size(), info.ModTime())
	}

	return dirCount, fileCount, totalSize
}

func printDirectoryEntry(name string, modTime time.Time) {
	modStr := modTime.Format("2006-01-02 15:04")
	fmt.Printf("  %s[DIR]%s  %-34s  %8s  %s%s%s\n",
		constants.ColorBlue, constants.ColorReset, name, "-",
		constants.ColorDim, modStr, constants.ColorReset)
}

func printFileEntry(name string, size int64, modTime time.Time) {
	modStr := modTime.Format("2006-01-02 15:04")
	sizeStr := formatByteSize(size)
	fmt.Printf("        %-34s  %8s  %s%s%s\n",
		name, sizeStr,
		constants.ColorDim, modStr, constants.ColorReset)
}

func printLsSummary(dirs, files int, totalSize int64) {
	fmt.Printf("  %s------------------------------------------------------------------------%s\n", constants.ColorDim, constants.ColorReset)
	fmt.Printf("  %sTotal:%s %d directories, %d files (%s)\n",
		constants.ColorCyan, constants.ColorReset, dirs, files, formatByteSize(totalSize))
	fmt.Printf("  %s(inspected directory. Enter command for Step, or '+add' to record 'ls')%s\n\n",
		constants.ColorDim, constants.ColorReset)
}

func formatByteSize(size int64) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%d B", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
	default:
		return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
	}
}

func executeInteractiveFind(pattern string) error {
	fmt.Printf("\n  %sSearching for files matching %q...%s\n", constants.ColorCyan, pattern, constants.ColorReset)
	matchCount := 0

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil || matchCount >= 30 {
			return nil
		}
		if d.IsDir() && isIgnoredDir(d.Name()) {
			return filepath.SkipDir
		}
		if isPatternMatched(d.Name(), pattern) {
			matchCount++
			printFindMatch(path, d.IsDir())
		}
		return nil
	})

	printFindSummary(matchCount, err)

	return nil
}

func isIgnoredDir(name string) bool {
	switch strings.ToLower(name) {
	case ".git", "node_modules", "vendor", ".gemini", "dist", "build", "bin":
		return true
	default:
		return false
	}
}

func isPatternMatched(name string, pattern string) bool {
	isMatched, err := filepath.Match(strings.ToLower(pattern), strings.ToLower(name))
	if err == nil && isMatched {
		return true
	}

	return strings.Contains(strings.ToLower(name), strings.ToLower(pattern))
}

func printFindMatch(path string, isDir bool) {
	relPath := filepath.ToSlash(path)
	if isDir {
		fmt.Printf("  %s[DIR]%s %s\n", constants.ColorBlue, constants.ColorReset, relPath)

		return
	}

	fmt.Printf("        %s\n", relPath)
}

func printFindSummary(count int, err error) {
	if err != nil {
		fmt.Printf("  %s▲ Error walking directory: %v%s\n\n", constants.ColorRed, err, constants.ColorReset)

		return
	}

	fmt.Printf("  %sTotal matches: %d%s\n\n", constants.ColorCyan, count, constants.ColorReset)
}

func executeInteractiveSearch(query string) error {
	fmt.Printf("\n  %sSearching for text %q...%s\n", constants.ColorCyan, query, constants.ColorReset)
	matchCount := 0

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil || matchCount >= 25 {
			return nil
		}
		if d.IsDir() && isIgnoredDir(d.Name()) {
			return filepath.SkipDir
		}
		if !d.IsDir() && isTextSearchTarget(d.Name()) {
			searchFileLines(path, query, &matchCount, 25)
		}
		return nil
	})

	printSearchSummary(matchCount, err)

	return nil
}

func isTextSearchTarget(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".go", ".ts", ".js", ".json", ".md", ".yaml", ".yml", ".ps1", ".sh", ".py", ".html", ".css", ".txt":
		return true
	default:
		return false
	}
}

func searchFileLines(path string, query string, count *int, maxCount int) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() && *count < maxCount {
		lineText := scanner.Text()
		if strings.Contains(strings.ToLower(lineText), strings.ToLower(query)) {
			*count++
			printSearchMatch(path, lineNum, strings.TrimSpace(lineText))
		}
		lineNum++
	}
}

func printSearchMatch(path string, lineNum int, snippet string) {
	relPath := filepath.ToSlash(path)
	if len(snippet) > 80 {
		snippet = snippet[:77] + "..."
	}
	fmt.Printf("  %s%s:%d:%s %s\n", constants.ColorCyan, relPath, lineNum, constants.ColorReset, snippet)
}

func printSearchSummary(count int, err error) {
	if err != nil {
		fmt.Printf("  %s▲ Error searching files: %v%s\n\n", constants.ColorRed, err, constants.ColorReset)

		return
	}

	fmt.Printf("  %sTotal matches: %d%s\n\n", constants.ColorCyan, count, constants.ColorReset)
}

func executeInteractiveReplace(rawArgs string) error {
	parts := strings.Fields(rawArgs)
	if len(parts) < 2 {
		fmt.Printf("  %s▲ Usage: replace <old> <new> [glob-pattern]%s\n\n", constants.ColorYellow, constants.ColorReset)

		return apperror.NewValidationError("replace requires <old> and <new> arguments")
	}

	oldStr, newStr := parts[0], parts[1]
	targetGlob := "*"
	if len(parts) > 2 {
		targetGlob = parts[2]
	}

	runFileReplacements(oldStr, newStr, targetGlob)

	return nil
}

func runFileReplacements(oldStr string, newStr string, targetGlob string) {
	replacedCount, filesChanged := 0, 0
	matches, err := filepath.Glob(targetGlob)
	if err != nil {
		fmt.Printf("  %s▲ Invalid glob pattern %q: %v%s\n\n", constants.ColorRed, targetGlob, err, constants.ColorReset)

		return
	}

	for _, file := range matches {
		count, repErr := replaceFileContent(file, oldStr, newStr)
		if repErr == nil && count > 0 {
			replacedCount += count
			filesChanged++
		}
	}

	fmt.Printf("  %s✓ Replaced %d occurrence(s) across %d file(s).%s\n\n",
		constants.ColorGreen, replacedCount, filesChanged, constants.ColorReset)
}

func replaceFileContent(path string, oldStr, newStr string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, apperror.WrapSimple(err, "read file for replacement")
	}

	content := string(data)
	count := strings.Count(content, oldStr)
	if count == 0 {
		return 0, nil
	}

	newContent := strings.ReplaceAll(content, oldStr, newStr)
	if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
		return 0, apperror.WrapSimple(err, "write replaced content")
	}

	return count, nil
}

func printInteractiveHelp() {
	fmt.Println()
	fmt.Printf("  %s● Interactive Macro Builder Helper Commands:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s------------------------------------------------------------------------%s\n", constants.ColorDim, constants.ColorReset)
	fmt.Println("    ls [dir] / dir [dir]       - Inspect directory contents (files & folders)")
	fmt.Println("    +add                       - Record the last inspected command as a step")
	fmt.Println("    add <command>              - Directly record a command without executing")
	fmt.Println("    pwd                        - Show current working directory path")
	fmt.Println("    pwd on / pwd off           - Enable or disable PWD header above step prompt")
	fmt.Println("    cd <directory>             - Change working directory for the session")
	fmt.Println("    find <pattern>             - Search for filenames matching pattern or glob")
	fmt.Println("    search <text> / grep <text>- Search file text contents for matching lines")
	fmt.Println("    replace <old> <new> [glob] - Replace text in matching files")
	fmt.Println("    done / exit / quit         - Save entered steps and create the macro")
	fmt.Println("    cancel / abort             - Abort macro creation without saving")
	fmt.Println("    exec on / exec off         - Toggle live command execution during macro building")
	fmt.Println("    rec / record               - Switch to live terminal recording session")
	fmt.Printf("  %s------------------------------------------------------------------------%s\n\n", constants.ColorDim, constants.ColorReset)
}
