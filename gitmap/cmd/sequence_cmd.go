package cmd

import (
	"context"
	"database/sql"
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
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
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

	return dispatchSequenceCommand(args)
}

func dispatchSequenceCommand(args []string) error {
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
	}

	return dispatchSequenceFallback(subCmd, args)
}

func dispatchSequenceFallback(subCmd string, args []string) error {
	if strings.HasPrefix(subCmd, "-") {
		return handleSequenceList(args)
	}

	return handleSequenceFix(args)
}

func parseSequenceFlags(args []string) (SequenceFlags, []string) {
	flags := SequenceFlags{
		StartNum: 1,
		PinMap:   make(map[string]int),
	}

	var dirs []string

	for i := 0; i < len(args); i++ {
		advanced, dir := parseOneSeqArg(args, i, &flags)
		i += advanced
		if dir != "" {
			dirs = append(dirs, dir)
		}
	}

	return flags, dirs
}

func parseOneSeqArg(args []string, i int, flags *SequenceFlags) (int, string) {
	arg := args[i]
	if parseSeqBoolFlags(arg, flags) {
		return 0, ""
	}

	adv := parseSeqValFlags(args, i, flags)
	if adv > 0 {
		return adv, ""
	}

	if !strings.HasPrefix(arg, "-") {
		return 0, arg
	}

	return 0, ""
}

func parseSeqBoolFlags(arg string, flags *SequenceFlags) bool {
	switch arg {
	case "--json", "-json":
		flags.IsJson = true

		return true
	case "--dry-run", "-dry-run":
		flags.IsDryRun = true

		return true
	case "--order-by-time", "-orderbytime":
		flags.IsOrderByTime = true

		return true
	case "--order-by-az", "-orderbyaz":
		flags.IsOrderByAZ = true

		return true
	}

	return false
}

func parseSeqValFlags(args []string, i int, flags *SequenceFlags) int {
	arg := args[i]
	if i+1 >= len(args) {
		return 0
	}

	nextVal := args[i+1]
	switch arg {
	case "--start", "-start":
		if val, err := strconv.Atoi(nextVal); err == nil {
			flags.StartNum = val
		}

		return 1
	case "--shift", "-shift":
		if val, err := strconv.Atoi(nextVal); err == nil {
			flags.ShiftNum = val
		}

		return 1
	case "--pin", "-pin":
		parsePinMap(nextVal, flags.PinMap)

		return 1
	}

	return 0
}

func handleSequenceList(args []string) error {
	flags, dirs := parseSequenceFlags(args)
	targetDir := resolveSequenceDir(dirs)

	payload, err := scanDirectorySequence(targetDir)
	if err != nil {
		return err
	}

	if err := saveSequenceToRepoDB(payload); err != nil {
		return err
	}

	return outputSequenceList(payload, flags)
}

func resolveSequenceDir(dirs []string) string {
	if len(dirs) > 0 {
		return dirs[0]
	}

	return "."
}

func outputSequenceList(payload *SequencePayload, flags SequenceFlags) error {
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

	items, seqCount := extractSequenceItems(dir, entries)
	sortSequenceItems(items)

	return buildSequencePayload(dir, items, seqCount), nil
}

func extractSequenceItems(dir string, entries []os.DirEntry) ([]SequenceItem, int) {
	re := regexp.MustCompile(`^(\d+)[-_](.*)$`)
	var items []SequenceItem
	sequencedCount := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		item, isSeq := buildSequenceItem(dir, entry.Name(), re)
		if isSeq {
			sequencedCount++
		}

		items = append(items, item)
	}

	return items, sequencedCount
}

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

func queryAndPrintSequenceGet(ctx context.Context, repoDB *sql.DB, targetDir string) error {
	query := "SELECT Filename, SequenceNumber, BaseName FROM FileSequence WHERE Directory = ? ORDER BY SequenceNumber ASC"
	rows, err := repoDB.QueryContext(ctx, query, targetDir)
	if err != nil {
		return err
	}

	defer rows.Close()

	files, err := scanSequenceRows(rows, targetDir)
	if err != nil {
		return err
	}

	payload := buildSequencePayload(targetDir, files, len(files))
	data, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println(string(data))

	return nil
}

func scanSequenceRows(rows *sql.Rows, targetDir string) ([]SequenceItem, error) {
	var files []SequenceItem
	for rows.Next() {
		var fn, bn string
		var seq int
		if err := rows.Scan(&fn, &seq, &bn); err != nil {
			return nil, apperror.WrapSimple(err, "scan sequence row")
		}

		files = append(files, SequenceItem{
			Sequence: seq,
			Filename: fn,
			BaseName: bn,
			Path:     filepath.Join(targetDir, fn),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "iterate sequence rows")
	}

	return files, nil
}

func handleSequenceHistory(args []string) error {
	ctx := context.Background()
	mainDB, repoDB, err := getRepoDB(ctx)
	if err != nil {
		return err
	}

	defer mainDB.Close()
	defer repoDB.Close()

	targetDir := resolveSequenceDir(args)
	targetDir = strings.ReplaceAll(targetDir, "\\", "/")

	return queryAndPrintSequenceHistory(ctx, repoDB, targetDir)
}

func queryAndPrintSequenceHistory(ctx context.Context, repoDB *sql.DB, targetDir string) error {
	query := "SELECT SequenceHistoryId, Directory, OperationsJson, CreatedAt FROM SequenceHistory WHERE Directory = ? ORDER BY CreatedAt DESC LIMIT 20"
	rows, err := repoDB.QueryContext(ctx, query, targetDir)
	if err != nil {
		return err
	}

	defer rows.Close()

	return printSequenceHistoryRows(rows, targetDir)
}

func printSequenceHistoryRows(rows *sql.Rows, targetDir string) error {
	fmt.Printf("\nSequence History for %s:\n", targetDir)
	for rows.Next() {
		var id int64
		var dir, opsJson string
		var created int64
		if err := rows.Scan(&id, &dir, &opsJson, &created); err != nil {
			return apperror.WrapSimple(err, "scan sequence history")
		}

		t := time.Unix(created, 0).UTC().Format(time.RFC3339)
		fmt.Printf("  [%s] ID #%d:\n    %s\n", t, id, opsJson)
	}

	if err := rows.Err(); err != nil {
		return apperror.WrapSimple(err, "iterate sequence history")
	}

	fmt.Println()

	return nil
}

func saveSequenceToRepoDB(payload *SequencePayload) error {
	ctx := context.Background()
	mainDB, repoDB, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, "get repo db for sequence")
	}

	defer mainDB.Close()
	defer repoDB.Close()

	wrapper, appErr := dbengine.WrapDb(repoDB, dbengine.DbSQLite)
	if appErr != nil {
		return appErr
	}

	if txErr := wrapper.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return executeSaveSequenceTx(ctx, tx, payload)
	}); txErr != nil {
		return txErr
	}

	return nil
}

func executeSaveSequenceTx(ctx context.Context, tx *dbengine.TxWrapper, payload *SequencePayload) *apperror.AppError {
	now := time.Now().Unix()
	_, err := tx.Exec(ctx, "DELETE FROM FileSequence WHERE Directory = ?", payload.Directory)
	if err != nil {
		return apperror.WrapSimple(err, "delete existing file sequence")
	}

	return insertSequenceFilesTx(ctx, tx, payload, now)
}

func insertSequenceFilesTx(ctx context.Context, tx *dbengine.TxWrapper, payload *SequencePayload, now int64) *apperror.AppError {
	stmt, err := tx.Prepare(ctx, "INSERT OR REPLACE INTO FileSequence (Directory, Filename, SequenceNumber, BaseName, UpdatedAt) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return apperror.WrapSimple(err, "prepare insert file sequence")
	}

	defer stmt.Close()

	for _, f := range payload.Files {
		if _, execErr := stmt.ExecContext(ctx, payload.Directory, f.Filename, f.Sequence, f.BaseName, now); execErr != nil {
			return apperror.WrapSimple(execErr, "exec insert file sequence item")
		}
	}

	return nil
}

func recordSequenceHistoryInDB(report SequenceFixReport) error {
	ctx := context.Background()
	mainDB, repoDB, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, "get repo db for history")
	}

	defer mainDB.Close()
	defer repoDB.Close()

	return insertSequenceHistoryTx(ctx, repoDB, report)
}

func insertSequenceHistoryTx(ctx context.Context, repoDB *sql.DB, report SequenceFixReport) error {
	opsJson, err := json.Marshal(report.Operations)
	if err != nil {
		return apperror.WrapSimple(err, "marshal sequence operations")
	}

	now := time.Now().Unix()
	query := "INSERT INTO SequenceHistory (Directory, OperationsJson, CreatedAt) VALUES (?, ?, ?)"
	_, err = repoDB.ExecContext(ctx, query, report.Directory, string(opsJson), now)
	if err != nil {
		return apperror.WrapSimple(err, "insert sequence history")
	}

	return nil
}
