package cmdconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonenvelope"
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

// runWhichFormatCLI inspects JSON files or directory to identify formats and import commands.
func RunWhichFormatCLI(args []string) error {
	opts := ParseWhichFormatArgs(args)
	files, err := ResolveInspectionFiles(opts.FileArgs, opts.TargetDir)
	if err != nil {
		fmt.Printf("Error resolving JSON files: %v\n", err)
		return nil
	}

	if len(files) == 0 {
		printNoJSONFilesBanner(opts.TargetDir)
		return nil
	}

	results := inspectJSONFiles(files)
	renderFormatInspectionReport(results)
	return nil
}

type whichFormatOptions struct {
	TargetDir string
	FileArgs  []string
	IsYes     bool
}

func ParseWhichFormatArgs(args []string) whichFormatOptions {
	opts := whichFormatOptions{TargetDir: "."}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-y" || arg == "--yes" {
			opts.IsYes = true
			continue
		}
		if (arg == "--dir" || arg == "-d") && i+1 < len(args) {
			opts.TargetDir = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(arg, "--dir=") {
			opts.TargetDir = strings.TrimPrefix(arg, "--dir=")
			continue
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		// Skip leading keyword aliases if invoked as "which format" or "format which"
		if isRoutingKeyword(arg) {
			continue
		}
		opts.FileArgs = append(opts.FileArgs, arg)
	}
	return opts
}

func isRoutingKeyword(arg string) bool {
	low := strings.ToLower(arg)
	return low == "which" || low == "format" || low == "inspect" || low == "which-format"
}

func ResolveInspectionFiles(fileArgs []string, targetDir string) ([]string, error) {
	if len(fileArgs) > 0 {
		return expandFileGlobs(fileArgs), nil
	}
	return scanDirectoryJSONFiles(targetDir)
}

func expandFileGlobs(patterns []string) []string {
	var collected []string
	for _, p := range patterns {
		matches, err := filepath.Glob(p)
		if err == nil && len(matches) > 0 {
			collected = append(collected, matches...)
			continue
		}
		collected = append(collected, p)
	}
	return dedupeAndFilterJSONFiles(collected)
}

func scanDirectoryJSONFiles(dir string) ([]string, error) {
	var matches []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(e.Name()), ".json") {
			matches = append(matches, filepath.Join(dir, e.Name()))
		}
	}
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
		res := inspectSingleFile(f)
		results = append(results, res)
	}
	return results
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

	desc, attrs, isMatched := jsonenvelope.DetectFormat(content)
	if !isMatched {
		return FormatInspectionResult{
			FilePath:        path,
			IsMatched:       false,
			UnmatchedReason: "Could not match known GitMap JSON format / unsupported schema (will not import).",
		}
	}

	normPath := filepath.ToSlash(path)
	baseName := jsonenvelope.RelativeBaseName(normPath)
	importCmd := fmt.Sprintf(desc.SuggestedImportCmd, baseName)
	if attrs.ImportCommand != "" {
		importCmd = attrs.ImportCommand
	}
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

func printNoJSONFilesBanner(dir string) {
	fmt.Printf("\nNo JSON files found in directory: %s\n", dir)
	fmt.Printf("Usage: gitmap which-format [file1.json file2.json ...]\n")
	fmt.Printf("       gitmap which-format --dir <directory>\n\n")
}

func renderFormatInspectionReport(results []FormatInspectionResult) {
	var matched []FormatInspectionResult
	var unmatched []FormatInspectionResult
	for _, r := range results {
		if r.IsMatched {
			matched = append(matched, r)
			continue
		}
		unmatched = append(unmatched, r)
	}

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

func printBatchCommandSuggestion(matched []FormatInspectionResult) {
	if len(matched) == 0 {
		fmt.Println("================================================================================")
		fmt.Println("No supported GitMap JSON formats were detected in the inspected files.")
		fmt.Println("================================================================================")
		return
	}

	var cmds []string
	for _, m := range matched {
		cmds = append(cmds, m.SuggestedImportCmd)
	}
	batchCmd := strings.Join(cmds, " && ")

	fmt.Println("================================================================================")
	fmt.Println("BATCH IMPORT SUGGESTION (SINGLE-LINE):")
	fmt.Println("================================================================================")
	fmt.Printf("%s%s%s\n\n", constants.ColorGreen, batchCmd, constants.ColorReset)
	fmt.Println("TIP: Append \"-y\" to your import commands to bypass all interactive prompts.")
	fmt.Println("================================================================================")
	fmt.Println()
}
