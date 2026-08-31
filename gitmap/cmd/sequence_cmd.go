package cmd

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
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/helptext"
)

// SequenceItem represents a single file item with its sequence metadata.
type SequenceItem struct {
	Sequence  int    `json:"sequence"`
	Filename  string `json:"filename"`
	BaseName  string `json:"baseName"`
	Extension string `json:"extension"`
	Path      string `json:"path"`
}

// SequencePayload represents the structured machine output for sequence queries.
type SequencePayload struct {
	Directory      string         `json:"directory"`
	TotalFiles     int            `json:"totalFiles"`
	SequencedFiles int            `json:"sequencedFiles"`
	Files          []SequenceItem `json:"files"`
}

// SequenceRenameOp records an individual file rename operation.
type SequenceRenameOp struct {
	From string `json:"from"`
	To   string `json:"to"`
	Seq  int    `json:"seq"`
}

// SequenceFixReport represents the output report after fixing/re-sequencing.
type SequenceFixReport struct {
	Directory  string             `json:"directory"`
	IsDryRun   bool               `json:"isDryRun"`
	TotalFixed int                `json:"totalFixed"`
	Operations []SequenceRenameOp `json:"operations"`
}

// SequenceFlags holds parsed CLI options for sequence commands.
type SequenceFlags struct {
	IsJson        bool
	IsDryRun      bool
	IsOrderByTime bool
	IsOrderByAZ   bool
	StartNum      int
	ShiftNum      int
	PinMap        map[string]int
}

func runSequence(args []string) error {
	checkHelp("sequence", args)
	if len(args) == 0 {
		_, mode := ParsePrettyFlag(args)
		helptext.PrintWithMode("sequence", mode)
		return nil
	}

	subCmd := args[0]
	subArgs := args[1:]

	switch subCmd {
	case "list", "ls":
		return handleSequenceList(subArgs)
	case "fix", "reorder", "resequence":
		return handleSequenceFix(subArgs)
	case "get":
		return handleSequenceGet(subArgs)
	case "history", "hist":
		return handleSequenceHistory(subArgs)
	default:
		if strings.HasPrefix(subCmd, "-") {
			return handleSequenceList(args)
		}
		return handleSequenceFix(args)
	}
}

func parseSequenceFlags(args []string) (SequenceFlags, []string) {
	flags := SequenceFlags{
		StartNum: 1,
		PinMap:   make(map[string]int),
	}
	var dirs []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--json" || arg == "-json":
			flags.IsJson = true
		case arg == "--dry-run" || arg == "-dry-run":
			flags.IsDryRun = true
		case arg == "--order-by-time" || arg == "-orderbytime":
			flags.IsOrderByTime = true
		case arg == "--order-by-az" || arg == "-orderbyaz":
			flags.IsOrderByAZ = true
		case (arg == "--start" || arg == "-start") && i+1 < len(args):
			if val, err := strconv.Atoi(args[i+1]); err == nil {
				flags.StartNum = val
			}
			i++
		case (arg == "--shift" || arg == "-shift") && i+1 < len(args):
			if val, err := strconv.Atoi(args[i+1]); err == nil {
				flags.ShiftNum = val
			}
			i++
		case (arg == "--pin" || arg == "-pin") && i+1 < len(args):
			parsePinMap(args[i+1], flags.PinMap)
			i++
		case !strings.HasPrefix(arg, "-"):
			dirs = append(dirs, arg)
		}
	}
	return flags, dirs
}

func handleSequenceList(args []string) error {
	flags, dirs := parseSequenceFlags(args)
	targetDir := "."
	if len(dirs) > 0 {
		targetDir = dirs[0]
	}

	payload, err := scanDirectorySequence(targetDir)
	if err != nil {
		return err
	}

	_ = saveSequenceToRepoDB(payload)

	if flags.IsJson {
		data, _ := json.MarshalIndent(payload, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	printSequenceTable(payload)
	return nil
}

func scanDirectorySequence(dir string) (*SequencePayload, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, apperror.Wrap(err, fmt.Sprintf("resolve path %s", dir), nil)
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, apperror.Wrap(err, fmt.Sprintf("read directory %s", dir), nil)
	}

	re := regexp.MustCompile(`^(\d+)[-_](.*)$`)
	var items []SequenceItem
	sequencedCount := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		match := re.FindStringSubmatch(name)
		seq := 0
		base := name

		if len(match) == 3 {
			seq, _ = strconv.Atoi(match[1])
			base = match[2]
			sequencedCount++
		}

		relPath := filepath.Join(dir, name)
		relPath = strings.ReplaceAll(relPath, "\\", "/")

		items = append(items, SequenceItem{
			Sequence:  seq,
			Filename:  name,
			BaseName:  base,
			Extension: filepath.Ext(name),
			Path:      relPath,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Sequence != items[j].Sequence {
			return items[i].Sequence < items[j].Sequence
		}
		return items[i].Filename < items[j].Filename
	})

	normalizedDir := strings.ReplaceAll(dir, "\\", "/")
	return &SequencePayload{
		Directory:      normalizedDir,
		TotalFiles:     len(items),
		SequencedFiles: sequencedCount,
		Files:          items,
	}, nil
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
	targetDir := "."
	if len(dirs) > 0 {
		targetDir = dirs[0]
	}

	absDir, err := filepath.Abs(targetDir)
	if err != nil {
		return apperror.Wrap(err, fmt.Sprintf("invalid directory %s", targetDir), nil)
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		return apperror.Wrap(err, fmt.Sprintf("read dir %s", targetDir), nil)
	}

	parsedFiles := parseSeqFiles(entries, absDir)
	applySequenceOrdering(parsedFiles, flags)

	report := executeSequenceRenames(parsedFiles, absDir, flags.IsDryRun)
	report.Directory = strings.ReplaceAll(targetDir, "\\", "/")

	persistSequenceChanges(report, targetDir, flags.IsDryRun)

	if flags.IsJson {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	printFixReport(report)
	return nil
}

func persistSequenceChanges(report SequenceFixReport, targetDir string, isDryRun bool) {
	if isDryRun {
		return
	}
	_ = recordSequenceHistoryInDB(report)
	updatedPayload, errScan := scanDirectorySequence(targetDir)
	if errScan == nil {
		_ = saveSequenceToRepoDB(updatedPayload)
	}
}

func applySequenceOrdering(files []*seqFile, flags SequenceFlags) {
	if flags.IsOrderByTime {
		sort.Slice(files, func(i, j int) bool { return files[i].Time < files[j].Time })
	} else if flags.IsOrderByAZ {
		sort.Slice(files, func(i, j int) bool {
			return strings.ToLower(files[i].Rest) < strings.ToLower(files[j].Rest)
		})
	}

	pinned, unpinned := partitionFiles(files, flags.PinMap)
	sortKeepOldOrder(unpinned)

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
	maxSeq := findMaxSeq(files)
	digits := 2
	if maxSeq > 99 {
		digits = len(strconv.Itoa(maxSeq))
	}

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
	_, repoDB, err := getRepoDB(ctx)
	if err != nil {
		return err
	}
	defer repoDB.Close()

	targetDir := "."
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		targetDir = args[0]
	}
	targetDir = strings.ReplaceAll(targetDir, "\\", "/")

	rows, err := repoDB.QueryContext(ctx, "SELECT Filename, SequenceNumber, BaseName FROM FileSequence WHERE Directory = ? ORDER BY SequenceNumber ASC", targetDir)
	if err != nil {
		return err
	}
	defer rows.Close()

	var files []SequenceItem
	for rows.Next() {
		var fn, bn string
		var seq int
		if errScan := rows.Scan(&fn, &seq, &bn); errScan == nil {
			files = append(files, SequenceItem{
				Sequence: seq,
				Filename: fn,
				BaseName: bn,
				Path:     filepath.Join(targetDir, fn),
			})
		}
	}

	payload := SequencePayload{
		Directory:      targetDir,
		TotalFiles:     len(files),
		SequencedFiles: len(files),
		Files:          files,
	}

	data, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println(string(data))
	return nil
}

func handleSequenceHistory(args []string) error {
	ctx := context.Background()
	_, repoDB, err := getRepoDB(ctx)
	if err != nil {
		return err
	}
	defer repoDB.Close()

	targetDir := "."
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		targetDir = args[0]
	}
	targetDir = strings.ReplaceAll(targetDir, "\\", "/")

	rows, err := repoDB.QueryContext(ctx, "SELECT Id, Directory, OperationsJson, CreatedAt FROM SequenceHistory WHERE Directory = ? ORDER BY CreatedAt DESC LIMIT 20", targetDir)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Printf("\nSequence History for %s:\n", targetDir)
	for rows.Next() {
		var id int64
		var dir, opsJson string
		var created int64
		if errScan := rows.Scan(&id, &dir, &opsJson, &created); errScan == nil {
			t := time.Unix(created, 0).UTC().Format(time.RFC3339)
			fmt.Printf("  [%s] ID #%d:\n    %s\n", t, id, opsJson)
		}
	}
	fmt.Println()
	return nil
}

func saveSequenceToRepoDB(payload *SequencePayload) error {
	ctx := context.Background()
	_, repoDB, err := getRepoDB(ctx)
	if err != nil {
		return nil
	}
	defer repoDB.Close()

	tx, err := repoDB.BeginTx(ctx, nil)
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	now := time.Now().Unix()
	_, _ = tx.ExecContext(ctx, "DELETE FROM FileSequence WHERE Directory = ?", payload.Directory)
	stmt, err := tx.PrepareContext(ctx, "INSERT OR REPLACE INTO FileSequence (Directory, Filename, SequenceNumber, BaseName, UpdatedAt) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return nil
	}
	defer stmt.Close()

	for _, f := range payload.Files {
		_, _ = stmt.ExecContext(ctx, payload.Directory, f.Filename, f.Sequence, f.BaseName, now)
	}

	return tx.Commit()
}

func recordSequenceHistoryInDB(report SequenceFixReport) error {
	ctx := context.Background()
	_, repoDB, err := getRepoDB(ctx)
	if err != nil {
		return nil
	}
	defer repoDB.Close()

	opsJson, err := json.Marshal(report.Operations)
	if err != nil {
		return nil
	}

	now := time.Now().Unix()
	_, err = repoDB.ExecContext(ctx, "INSERT INTO SequenceHistory (Directory, OperationsJson, CreatedAt) VALUES (?, ?, ?)", report.Directory, string(opsJson), now)
	return err
}

