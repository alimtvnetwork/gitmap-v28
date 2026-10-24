package cmdsequence

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func buildSequenceItem(dir, name string, re *regexp.Regexp) (SequenceItem, bool) {
	match := re.FindStringSubmatch(name)
	seq := 0
	base := name
	isSeq := false
	if len(match) == 3 {
		seq, _ = strconv.Atoi(match[1])
		base = match[2]
		isSeq = true
	}

	relPath := strings.ReplaceAll(filepath.Join(dir, name), "\\", "/")

	return SequenceItem{
		Sequence:  seq,
		Filename:  name,
		BaseName:  base,
		Extension: filepath.Ext(name),
		Path:      relPath,
	}, isSeq
}

func sortSequenceItems(items []SequenceItem) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Sequence != items[j].Sequence {
			return items[i].Sequence < items[j].Sequence
		}

		return items[i].Filename < items[j].Filename
	})
}

func buildSequencePayload(dir string, items []SequenceItem, seqCount int) *SequencePayload {
	return &SequencePayload{
		Directory:      strings.ReplaceAll(dir, "\\", "/"),
		TotalFiles:     len(items),
		SequencedFiles: seqCount,
		Files:          items,
	}
}

func printSequenceTable(payload *SequencePayload) {
	fmt.Printf("\nDirectory Sequence: %s (%d files, %d sequenced)\n", payload.Directory, payload.TotalFiles, payload.SequencedFiles)
	fmt.Printf("%-6s  %-35s  %-30s\n", "SEQ", "FILENAME", "BASE NAME")
	fmt.Println(strings.Repeat("-", 75))

	for _, f := range payload.Files {
		seqStr := "--"
		if f.Sequence > 0 {
			seqStr = fmt.Sprintf("%02d", f.Sequence)
		}

		fmt.Printf("%-6s  %-35s  %-30s\n", seqStr, f.Filename, f.BaseName)
	}

	fmt.Println()
}

func handleSequenceFix(args []string) error {
	flags, dirs := parseSequenceFlags(args)
	targetDir := resolveSequenceDir(dirs)

	absDir, err := filepath.Abs(targetDir)
	if err != nil {
		return apperror.Wrap(err, fmt.Sprintf("invalid directory %s", targetDir), nil)
	}

	report, err := processSequenceFix(absDir, targetDir, flags)
	if err != nil {
		return err
	}

	return outputSequenceFixReport(report, flags)
}

func processSequenceFix(absDir, targetDir string, flags SequenceFlags) (SequenceFixReport, error) {
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return SequenceFixReport{}, apperror.Wrap(err, fmt.Sprintf("read dir %s", targetDir), nil)
	}

	parsedFiles := parseSeqFiles(entries, absDir)
	applySequenceOrdering(parsedFiles, flags)

	report := executeSequenceRenames(parsedFiles, absDir, flags.IsDryRun)
	report.Directory = strings.ReplaceAll(targetDir, "\\", "/")
	if err := persistSequenceChanges(report, targetDir, flags.IsDryRun); err != nil {
		return SequenceFixReport{}, err
	}

	return report, nil
}

func outputSequenceFixReport(report SequenceFixReport, flags SequenceFlags) error {
	if flags.IsJson {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))

		return nil
	}

	printFixReport(report)

	return nil
}

func persistSequenceChanges(report SequenceFixReport, targetDir string, isDryRun bool) error {
	if isDryRun {
		return nil
	}

	if err := recordSequenceHistoryInDB(report); err != nil {
		return err
	}

	updatedPayload, errScan := scanDirectorySequence(targetDir)
	if errScan != nil {
		return errScan
	}

	return saveSequenceToRepoDB(updatedPayload)
}

func applySequenceOrdering(files []*seqFile, flags SequenceFlags) {
	sortSeqFiles(files, flags)
	pinned, unpinned := partitionFiles(files, flags.PinMap)
	sortKeepOldOrder(unpinned)
	assignNewSequences(unpinned, pinned, flags)
}

func sortSeqFiles(files []*seqFile, flags SequenceFlags) {
	if flags.IsOrderByTime {
		sort.Slice(files, func(i, j int) bool { return files[i].Time < files[j].Time })
	} else if flags.IsOrderByAZ {
		sort.Slice(files, func(i, j int) bool {
			return strings.ToLower(files[i].Rest) < strings.ToLower(files[j].Rest)
		})
	}
}

func assignNewSequences(unpinned, pinned []*seqFile, flags SequenceFlags) {
	usedSeqs := buildUsedSeqs(pinned)
	currentSeq := flags.StartNum

	for _, pf := range unpinned {
		for usedSeqs[currentSeq] {
			currentSeq++
		}

		newSeq := currentSeq + flags.ShiftNum
		pf.NewSeq = newSeq
		usedSeqs[newSeq] = true
		currentSeq++
	}
}

func executeSequenceRenames(files []*seqFile, absDir string, isDryRun bool) SequenceFixReport {
	digits := calcSeqDigits(files)
	var ops []SequenceRenameOp
	for _, pf := range files {
		format := fmt.Sprintf("%%0%dd-%%s", digits)
		if op := processSeqRename(pf, absDir, format, isDryRun); op != nil {
			ops = append(ops, *op)
		}
	}

	return SequenceFixReport{
		IsDryRun:   isDryRun,
		TotalFixed: len(ops),
		Operations: ops,
	}
}

func calcSeqDigits(files []*seqFile) int {
	maxSeq := findMaxSeq(files)
	if maxSeq > 99 {
		return len(strconv.Itoa(maxSeq))
	}

	return 2
}

func processSeqRename(pf *seqFile, absDir, format string, isDryRun bool) *SequenceRenameOp {
	newName := fmt.Sprintf(format, pf.NewSeq, pf.Rest)
	if newName == pf.BaseName {
		return nil
	}

	oldPath := pf.OriginalPath
	newPath := filepath.Join(absDir, newName)
	if !isDryRun {
		executeRename(oldPath, newPath, pf.BaseName, newName)
	}

	return &SequenceRenameOp{
		From: pf.BaseName,
		To:   newName,
		Seq:  pf.NewSeq,
	}
}

func printFixReport(report SequenceFixReport) {
	statusPrefix := "✓ Renamed"
	if report.IsDryRun {
		statusPrefix = "[DRY RUN] Would rename"
	}

	fmt.Printf("\nSequence Fix Report: %s (%d files)\n", report.Directory, report.TotalFixed)
	fmt.Println(strings.Repeat("-", 60))
	for _, op := range report.Operations {
		fmt.Printf("  %s %s -> %s (seq %02d)\n", statusPrefix, op.From, op.To, op.Seq)
	}

	if report.TotalFixed == 0 {
		fmt.Println("  No files needed re-sequencing (already cleanly ordered).")
	}

	fmt.Println()
}

func handleSequenceGet(args []string) error {
	ctx := context.Background()
	mainDB, repoDB, err := getRepoDB(ctx)
	if err != nil {
		return err
	}

	defer mainDB.Close()
	defer repoDB.Close()

	targetDir := resolveSequenceDir(args)
	targetDir = strings.ReplaceAll(targetDir, "\\", "/")

	return queryAndPrintSequenceGet(ctx, repoDB, targetDir)
}
