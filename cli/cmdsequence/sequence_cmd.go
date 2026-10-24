package cmdsequence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
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

func RunSequence(args []string) error {
	checkHelp("sequence", args)
	if len(args) == 0 {
		_, mode := ParsePrettyFlag(args)
		helpdoc.PrintWithMode("sequence", mode)

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
	re := lazyregex.NumberPrefixRegex.Compiled()
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
