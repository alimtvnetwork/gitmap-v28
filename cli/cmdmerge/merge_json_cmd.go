package cmdmerge

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonenvelope"
)

type mergeJSONOptions struct {
	TargetPattern string
	OutputFile    string
	TargetType    string
	IsYes         bool
	IsDryRun      bool
	IsVerbose     bool
	IsStdout      bool
}

type mergeFileRecord struct {
	FilePath string
	Type     string
	RawData  json.RawMessage
}

type mergeJSONSummary struct {
	TargetType     string
	FilesScanned   int
	TotalItemsIn   int
	DuplicateItems int
	FinalItemCount int
	OutputFile     string
	ImportCommand  string
}

// runMergeJSONCLI dispatches the `gitmap merge-json` and `gitmap json merge` commands.
func RunMergeJSONCLI(args []string) error {
	opts := parseMergeJSONArgs(args)
	files, err := resolveInspectionFiles([]string{opts.TargetPattern}, ".")
	if err != nil || len(files) == 0 {
		printNoMergeFilesFound(opts.TargetPattern)
		return nil
	}

	records := collectMergeableRecords(files, opts.TargetType)
	if len(records) == 0 {
		printNoMatchingRecordsFound(opts.TargetType)
		return nil
	}

	summary, mergedBytes, err := executeRecordsMerge(records, opts)
	if err != nil {
		fmt.Printf("  Merge failed: %v\n", err)
		return nil
	}

	renderMergeSummary(summary, opts, mergedBytes)
	return nil
}

func parseMergeJSONArgs(args []string) mergeJSONOptions {
	opts := mergeJSONOptions{TargetPattern: "*.json"}
	for i := 0; i < len(args); i++ {
		parseSingleMergeFlag(args[i], args, &i, &opts)
	}
	return opts
}

func parseSingleMergeFlag(a string, args []string, i *int, opts *mergeJSONOptions) {
	switch {
	case a == "-y" || a == "--yes":
		opts.IsYes = true
	case a == "-n" || a == "--dry-run":
		opts.IsDryRun = true
	case a == "-v" || a == "--verbose":
		opts.IsVerbose = true
	case a == "--stdout":
		opts.IsStdout = true
	case (a == "-o" || a == "--output") && *i+1 < len(args):
		*i++
		opts.OutputFile = args[*i]
	case (a == "-t" || a == "--type") && *i+1 < len(args):
		*i++
		opts.TargetType = args[*i]
	case !strings.HasPrefix(a, "-") && a != "merge" && a != "json":
		opts.TargetPattern = a
	}
}

func collectMergeableRecords(files []string, targetType string) []mergeFileRecord {
	var records []mergeFileRecord
	for _, f := range files {
		rec, isOk := readAndClassifyFile(f, targetType)
		if isOk {
			records = append(records, rec)
		}
	}
	return records
}

func readAndClassifyFile(filePath, targetType string) (mergeFileRecord, bool) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return mergeFileRecord{}, false
	}
	desc, _, isMatched := jsonenvelope.DetectFormat(content)
	if !isMatched {
		return mergeFileRecord{}, false
	}
	if targetType != "" && !strings.EqualFold(desc.Type, targetType) {
		return mergeFileRecord{}, false
	}
	payload, _, err := jsonenvelope.ExtractPayload(content)
	if err != nil {
		return mergeFileRecord{}, false
	}
	return mergeFileRecord{
		FilePath: filePath,
		Type:     desc.Type,
		RawData:  payload,
	}, true
}

func executeRecordsMerge(records []mergeFileRecord, opts mergeJSONOptions) (mergeJSONSummary, []byte, error) {
	primaryType := records[0].Type
	mergedMap, totalIn := extractAndDeduplicateItems(records, primaryType)
	finalItems := reindexMergedItems(mergedMap)

	outPath := opts.OutputFile
	if outPath == "" {
		outPath = fmt.Sprintf("merged-%s.json", primaryType)
	}

	envelope := buildMergedEnvelope(primaryType, outPath, finalItems)
	data, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return mergeJSONSummary{}, nil, err
	}

	summary := mergeJSONSummary{
		TargetType:     primaryType,
		FilesScanned:   len(records),
		TotalItemsIn:   totalIn,
		DuplicateItems: totalIn - len(finalItems),
		FinalItemCount: len(finalItems),
		OutputFile:     outPath,
		ImportCommand:  envelope.Attributes.ImportCommand,
	}

	if !opts.IsDryRun && !opts.IsStdout {
		_ = os.WriteFile(outPath, data, 0644)
	}
	return summary, data, nil
}

func extractAndDeduplicateItems(records []mergeFileRecord, primaryType string) ([]map[string]any, int) {
	var allItems []map[string]any
	totalIn := 0
	for _, r := range records {
		items := unmarshalItemsSlice(r.RawData)
		totalIn += len(items)
		allItems = append(allItems, items...)
	}

	seen := make(map[string]bool)
	var deduped []map[string]any
	for _, it := range allItems {
		key := extractDedupeKey(it)
		if seen[key] {
			continue
		}
		seen[key] = true
		deduped = append(deduped, it)
	}
	return deduped, totalIn
}

func unmarshalItemsSlice(raw json.RawMessage) []map[string]any {
	var direct []map[string]any
	if err := json.Unmarshal(raw, &direct); err == nil && len(direct) > 0 {
		return direct
	}

	var wrapped struct {
		Nodes     []map[string]any `json:"nodes"`
		Items     []map[string]any `json:"items"`
		Templates []map[string]any `json:"templates"`
		Macros    []map[string]any `json:"macros"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil {
		return pickWrappedSlice(wrapped.Nodes, wrapped.Items, wrapped.Templates, wrapped.Macros)
	}
	return nil
}

func pickWrappedSlice(nodes, items, tmpls, macros []map[string]any) []map[string]any {
	if len(nodes) > 0 {
		return nodes
	}
	if len(items) > 0 {
		return items
	}
	if len(tmpls) > 0 {
		return tmpls
	}
	return macros
}

func extractDedupeKey(it map[string]any) string {
	candidates := []string{"alias", "ipAddress", "ip_address", "slug", "repoId", "name", "id", "workerId"}
	for _, c := range candidates {
		key := resolveCandidateDedupeKey(c, it[c])
		if key != "" {
			return key
		}
	}
	bytes, _ := json.Marshal(it)
	return string(bytes)
}

func resolveCandidateDedupeKey(c string, val any) string {
	if val == nil {
		return ""
	}
	str := fmt.Sprintf("%v", val)
	if str == "" || str == "0" {
		return ""
	}
	return c + ":" + str
}

func reindexMergedItems(items []map[string]any) []map[string]any {
	for idx, it := range items {
		it["id"] = idx + 1
	}
	return items
}

func buildMergedEnvelope(primaryType, outPath string, finalItems []map[string]any) jsonenvelope.Envelope[any] {
	var payload any = finalItems
	if primaryType == jsonenvelope.TypeSSHNodes {
		payload = map[string]any{
			"schemaVersion": "2.0",
			"totalNodes":    len(finalItems),
			"nodes":         finalItems,
		}
	}
	return jsonenvelope.NewEnvelope(primaryType, outPath, "gitmap merge-json", "2.0", payload)
}

func renderMergeSummary(s mergeJSONSummary, opts mergeJSONOptions, data []byte) {
	if opts.IsStdout {
		fmt.Println(string(data))
		return
	}

	fmt.Printf("\n  %s✨ GitMap Multi-JSON Merge Summary%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  • Target Format Type:    %s%s%s\n", constants.ColorGreen, s.TargetType, constants.ColorReset)
	fmt.Printf("  • Source Files Scanned:  %d file(s)\n", s.FilesScanned)
	fmt.Printf("  • Raw Input Records:     %d item(s)\n", s.TotalItemsIn)
	fmt.Printf("  • Duplicates Removed:    %d item(s)\n", s.DuplicateItems)
	fmt.Printf("  • Unique Merged Items:   %s%d item(s)%s (1-based index applied)\n", constants.ColorGreen, s.FinalItemCount, constants.ColorReset)

	if opts.IsDryRun {
		fmt.Printf("  • Output File:           %s (dry-run, not written)\n", s.OutputFile)
	} else {
		fmt.Printf("  • Output File:           %s%s%s (written)\n", constants.ColorYellow, s.OutputFile, constants.ColorReset)
	}

	if s.ImportCommand != "" {
		fmt.Printf("\n  💡 Next Step: Import the merged file:\n")
		fmt.Printf("     %s\n\n", s.ImportCommand)
	}
}

func printNoMergeFilesFound(pattern string) {
	fmt.Printf("\n  %s⚠️  No JSON files found matching pattern: %q%s\n\n", constants.ColorYellow, pattern, constants.ColorReset)
}

func printNoMatchingRecordsFound(targetType string) {
	fmt.Printf("\n  %s⚠️  No JSON files matched the target format: %q%s\n\n", constants.ColorYellow, targetType, constants.ColorReset)
}
