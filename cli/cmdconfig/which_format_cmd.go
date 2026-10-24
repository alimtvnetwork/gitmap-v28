package cmdconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	cfgconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonx"
)

// FormatInspectionResult records the inspection status of a single JSON file.
type FormatInspectionResult struct {
	FilePath           string
	IsMatched          bool
	Type               string
	TypeName           string
	Source             string
	Version            string
	SuggestedImportCmd string
	SystemImpact       string
	UnmatchedReason    string
}

// inspectJSONFilesOrBanner prints banner if empty or returns inspection results.
func inspectJSONFilesOrBanner(files []string, targetDir string) ([]FormatInspectionResult, bool) {
	if len(files) == 0 {
		printNoJSONFilesBanner(targetDir)

		return nil, false
	}

	return inspectJSONFiles(files), true
}

// resolveAndInspectJSONFiles resolves file targets and runs format inspections.
func resolveAndInspectJSONFiles(opts whichFormatOptions) ([]FormatInspectionResult, bool) {
	files, err := ResolveInspectionFiles(opts.FileArgs, opts.TargetDir)
	if err != nil {
		fmt.Printf("Error resolving JSON files: %v\n", err)

		return nil, false
	}

	return inspectJSONFilesOrBanner(files, opts.TargetDir)
}

// RunWhichFormatCLI inspects JSON files or directory to identify formats and import commands.
func RunWhichFormatCLI(args []string) error {
	opts := ParseWhichFormatArgs(args)
	results, hasResults := resolveAndInspectJSONFiles(opts)
	if hasResults {
		renderFormatInspectionReport(results)
	}

	return nil
}

type whichFormatOptions struct {
	TargetDir string
	FileArgs  []string
	IsYes     bool
}

func isYesFlag(arg string) bool {
	return arg == cfgconst.FlagArgYesShort || arg == cfgconst.FlagArgYes
}

func appendNonFlagArg(arg string, opts *whichFormatOptions) {
	if !strings.HasPrefix(arg, cfgconst.FlagPrefixDash) && !isRoutingKeyword(arg) {
		opts.FileArgs = append(opts.FileArgs, arg)
	}
}

func applyArgToWhichFormatOptions(arg, next string, hasNext bool, opts *whichFormatOptions) int {
	if isYesFlag(arg) {
		opts.IsYes = true

		return 0
	}

	if consumed, skip := tryConsumeDirArg(arg, next, hasNext, &opts.TargetDir); consumed {
		return skip
	}

	appendNonFlagArg(arg, opts)

	return 0
}

// ParseWhichFormatArgs parses CLI arguments into whichFormatOptions.
func ParseWhichFormatArgs(args []string) whichFormatOptions {
	opts := whichFormatOptions{TargetDir: cfgconst.DefaultTargetDir}
	for i := 0; i < len(args); i++ {
		next := ""
		if i+1 < len(args) {
			next = args[i+1]
		}

		i += applyArgToWhichFormatOptions(args[i], next, i+1 < len(args), &opts)
	}

	return opts
}

func isRoutingKeyword(arg string) bool {
	low := strings.ToLower(arg)

	return low == cfgconst.CmdWhich || low == cfgconst.CmdFormat || low == cfgconst.CmdInspect || low == cfgconst.CmdWhichFormat
}

// ResolveInspectionFiles resolves specific file args or discovers directory JSONs.
func ResolveInspectionFiles(fileArgs []string, targetDir string) ([]string, error) {
	if len(fileArgs) > 0 {
		return expandFileGlobs(fileArgs), nil
	}

	return scanDirectoryJSONFiles(targetDir)
}

func expandFileGlobs(patterns []string) []string {
	var collected []string
	for _, p := range patterns {
		collected = appendPatternMatches(collected, p)
	}

	return dedupeAndFilterJSONFiles(collected)
}

func filterJSONDirEntries(dir string, entries []os.DirEntry) []string {
	var matches []string
	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), cfgconst.ExtJSON) {
			matches = append(matches, filepath.Join(dir, e.Name()))
		}
	}

	return matches
}

func scanDirectoryJSONFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	matches := filterJSONDirEntries(dir, entries)
	sort.Strings(matches)

	return matches, nil
}

func dedupeAndFilterJSONFiles(files []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, f := range files {
		clean := filepath.Clean(f)
		if seen[clean] {
			continue
		}

		seen[clean] = true
		out = append(out, clean)
	}

	return out
}

func inspectJSONFiles(files []string) []FormatInspectionResult {
	results := make([]FormatInspectionResult, 0, len(files))
	for _, f := range files {
		results = append(results, inspectSingleFile(f))
	}

	return results
}

func resolveImportCommand(suggested, explicit string) string {
	if explicit != "" {
		return explicit
	}

	return suggested
}

func buildMatchedFormatResult(path string, desc jsonx.TypeDescriptor, attrs jsonx.EnvelopeAttributes) FormatInspectionResult {
	normPath := filepath.ToSlash(path)
	importCmd := resolveImportCommand(fmt.Sprintf(desc.SuggestedImportCmd, jsonx.RelativeBaseName(normPath)), attrs.ImportCommand)

	return FormatInspectionResult{
		FilePath:           normPath,
		IsMatched:          true,
		Type:               desc.Type,
		TypeName:           desc.Name,
		Source:             attrs.Source,
		Version:            attrs.Version,
		SuggestedImportCmd: importCmd,
		SystemImpact:       desc.SystemImpact,
	}
}

func detectAndBuildJSONResult(path string, content []byte) FormatInspectionResult {
	desc, attrs, isMatched := jsonx.DetectFormat(content)
	if !isMatched {
		return FormatInspectionResult{
			FilePath:        path,
			IsMatched:       false,
			UnmatchedReason: "Could not match known GitMap JSON format / unsupported schema (will not import).",
		}
	}

	return buildMatchedFormatResult(path, desc, attrs)
}

func inspectSingleFile(path string) FormatInspectionResult {
	content, err := os.ReadFile(path)
	if err != nil {
		return FormatInspectionResult{
			FilePath:        path,
			IsMatched:       false,
			UnmatchedReason: fmt.Sprintf("Unable to read file: %v", err),
		}
	}

	return detectAndBuildJSONResult(path, content)
}

func printNoJSONFilesBanner(dir string) {
	fmt.Printf("\nNo JSON files found in directory: %s\n", dir)
	fmt.Printf("Usage: gitmap which-format [file1.json file2.json ...]\n")
	fmt.Printf("       gitmap which-format --dir <directory>\n\n")
}

func partitionInspectionResults(results []FormatInspectionResult) ([]FormatInspectionResult, []FormatInspectionResult) {
	var matched, unmatched []FormatInspectionResult
	for _, r := range results {
		if r.IsMatched {
			matched = append(matched, r)
			continue
		}

		unmatched = append(unmatched, r)
	}

	return matched, unmatched
}

func renderFormatInspectionReport(results []FormatInspectionResult) {
	matched, unmatched := partitionInspectionResults(results)
	printInspectionHeader(len(results), len(matched), len(unmatched))
	printMatchedResults(matched)
	printUnmatchedResults(unmatched)
	printBatchCommandSuggestion(matched)
}

func printInspectionHeader(total, matched, unmatched int) {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║ GITMAP JSON FORMAT INSPECTOR (WHICH-FORMAT)                                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════════╝")
	fmt.Printf("  Scanned: %d JSON file(s) | %sMatched: %d%s | %sUnmatched: %d%s\n\n",
		total,
		constants.ColorGreen, matched, constants.ColorReset,
		constants.ColorYellow, unmatched, constants.ColorReset,
	)
}

func printMatchedResults(matched []FormatInspectionResult) {
	if len(matched) == 0 {
		return
	}

	fmt.Printf("%s✔ MATCHED FORMATS (%d):%s\n\n", constants.ColorGreen, len(matched), constants.ColorReset)
	for i, m := range matched {
		printSingleMatchedResult(i+1, m)
	}
}

func printSingleMatchedResult(idx int, m FormatInspectionResult) {
	fmt.Printf("  %d. [%s%s%s]\n", idx, constants.ColorCyan, m.FilePath, constants.ColorReset)
	fmt.Printf("     • Type: %s (%s)\n", m.Type, m.TypeName)
	if m.Source != "" {
		fmt.Printf("     • Source: %s (v%s)\n", m.Source, m.Version)
	}

	fmt.Printf("     • Suggested Import Command:\n")
	fmt.Printf("       %s%s%s\n", constants.ColorYellow, m.SuggestedImportCmd, constants.ColorReset)
	fmt.Printf("     • System Impact:\n")
	fmt.Printf("       %s\n\n", m.SystemImpact)
}

func printUnmatchedResults(unmatched []FormatInspectionResult) {
	if len(unmatched) == 0 {
		return
	}

	fmt.Printf("%s✖ UNMATCHED FORMATS (%d):%s\n\n", constants.ColorRed, len(unmatched), constants.ColorReset)
	for _, u := range unmatched {
		fmt.Printf("  • [%s%s%s]\n", constants.ColorRed, u.FilePath, constants.ColorReset)
		fmt.Printf("    Status: NOT MATCHED\n")
		fmt.Printf("    Notice: %s\n\n", u.UnmatchedReason)
	}
}

func renderBatchSuggestionNotice() {
	fmt.Println("================================================================================")
	fmt.Println("No supported GitMap JSON formats were detected in the inspected files.")
	fmt.Println("================================================================================")
}

func renderBatchSuggestionBanner(batchCmd string) {
	fmt.Println("================================================================================")
	fmt.Println("BATCH IMPORT SUGGESTION (SINGLE-LINE):")
	fmt.Println("================================================================================")
	fmt.Printf("%s%s%s\n\n", constants.ColorGreen, batchCmd, constants.ColorReset)
	fmt.Println("TIP: Append \"-y\" to your import commands to bypass all interactive prompts.")
	fmt.Println("================================================================================")
	fmt.Println()
}

func printBatchCommandSuggestion(matched []FormatInspectionResult) {
	if len(matched) == 0 {
		renderBatchSuggestionNotice()

		return
	}

	var cmds []string
	for _, m := range matched {
		cmds = append(cmds, m.SuggestedImportCmd)
	}

	renderBatchSuggestionBanner(strings.Join(cmds, " && "))
}
