package cmdfind

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/cmdfind/enums"

// MatchKindType aliases the strongly-typed enum from the enums package.
type MatchKindType = enums.MatchKindType

const (
	MatchExact      = enums.MatchExact
	MatchContains   = enums.MatchContains
	MatchStartsWith = enums.MatchStartsWith
	MatchEndsWith   = enums.MatchEndsWith
	MatchWildcard   = enums.MatchWildcard
)

type MatchKind = enums.MatchKindType

type FindFilesOptions struct {
	Query     string
	Kind      MatchKind
	Exts      []string
	Limit     int
	IsJson    bool
	TargetDir string
}

func RunFindFiles(args []string) error {
	checkHelp(constants.CmdFindFiles, args)

	return executeFindFiles(args, MatchExact)
}

func RunFindFilesAny(args []string) error {
	checkHelp(constants.CmdFindFilesAny, args)

	return executeFindFiles(args, MatchContains)
}

func RunFindFilesStartsWith(args []string) error {
	checkHelp(constants.CmdFindFilesStartsWith, args)

	return executeFindFiles(args, MatchStartsWith)
}

func RunFindFilesEndsWith(args []string) error {
	checkHelp(constants.CmdFindFilesEndsWith, args)

	return executeFindFiles(args, MatchEndsWith)
}

func RunListFiles(args []string) error {
	checkHelp(constants.CmdListFiles, args)

	return executeFindFiles(args, MatchWildcard)
}

func executeFindFiles(args []string, defaultKind MatchKind) error {
	start := time.Now()
	opts := parseFindFilesOptions(args, defaultKind)
	if opts.Query == "" && defaultKind != MatchWildcard {
		showFindFilesHelp(defaultKind)

		return nil
	}

	if opts.Query == "" && defaultKind == MatchWildcard {
		showFindFilesHelp(defaultKind)

		return nil
	}

	matches, err := scanAndMatchFiles(opts)
	if err != nil {
		return apperror.WrapSimple(err, "find-files failed")
	}

	elapsed := time.Since(start)

	return outputFindResults(matches, opts.IsJson, elapsed)
}

func parseFindFilesOptions(args []string, defaultKind MatchKind) FindFilesOptions {
	opts := FindFilesOptions{
		Kind:      defaultKind,
		TargetDir: ".",
	}
	var nonFlags []string
	opts, nonFlags = processFindArgs(args, opts)
	if opts.Query == "" && len(nonFlags) > 0 {
		opts.Query = nonFlags[0]
	}

	return opts
}

func processFindArgs(args []string, opts FindFilesOptions) (FindFilesOptions, []string) {
	var nonFlags []string
	for i := 0; i < len(args); i++ {
		opts, nonFlags, i = evaluateFindArg(args, i, opts, nonFlags)
	}

	return opts, nonFlags
}

func evaluateFindArg(args []string, i int, opts FindFilesOptions, nonFlags []string) (FindFilesOptions, []string, int) {
	arg := args[i]
	if isDirFlag(arg) && i+1 < len(args) {
		opts.TargetDir = args[i+1]
		return opts, nonFlags, i + 1
	}
	if isQueryFlag(arg) && i+1 < len(args) {
		opts.Query = args[i+1]
		return opts, nonFlags, i + 1
	}
	if isExtFlag(arg) && i+1 < len(args) {
		opts.Exts = parseExtensionList(args[i+1])
		return opts, nonFlags, i + 1
	}
	if slices.Contains([]string{"--limit", "-l"}, arg) && i+1 < len(args) {
		opts.Limit, _ = strconv.Atoi(args[i+1])
		return opts, nonFlags, i + 1
	}

	return evaluateFindInlineArg(args, i, opts, nonFlags)
}

func evaluateFindInlineArg(args []string, i int, opts FindFilesOptions, nonFlags []string) (FindFilesOptions, []string, int) {
	arg := args[i]
	if slices.Contains([]string{"--json", "-json"}, arg) {
		opts.IsJson = true
		return opts, nonFlags, i
	}
	if val, hasVal := extractPrefixedFlag(arg, "--dir=", "--root=", "--path="); hasVal {
		opts.TargetDir = val
		return opts, nonFlags, i
	}
	if val, hasVal := extractPrefixedFlag(arg, "--query=", "--pattern="); hasVal {
		opts.Query = val
		return opts, nonFlags, i
	}
	if val, hasVal := extractPrefixedFlag(arg, "-ext=", "--ext="); hasVal {
		opts.Exts = parseExtensionList(val)
		return opts, nonFlags, i
	}
	if !strings.HasPrefix(arg, "-") {
		nonFlags = append(nonFlags, arg)
	}

	return opts, nonFlags, i
}

func isDirFlag(arg string) bool {
	return slices.Contains([]string{"--dir", "-d", "--root", "--path", "-p"}, arg)
}

func isQueryFlag(arg string) bool {
	return slices.Contains([]string{"--query", "-q", "--pattern"}, arg)
}

func isExtFlag(arg string) bool {
	return slices.Contains([]string{"-ext", "--ext"}, arg)
}

func extractPrefixedFlag(arg string, prefixes ...string) (string, bool) {
	for _, p := range prefixes {
		if strings.HasPrefix(arg, p) {
			return strings.TrimPrefix(arg, p), true
		}
	}

	return "", false
}

func scanAndMatchFiles(opts FindFilesOptions) ([]string, error) {
	absRoot, err := filepath.Abs(opts.TargetDir)
	if err != nil {
		return nil, err
	}

	var results []string
	errWalk := filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, errIn error) error {
		if err := handleFindWalkErr(errIn); err != nil {
			return err
		}

		if d.IsDir() {
			return handleFindDirSkip(d.Name())
		}

		rel, _ := filepath.Rel(absRoot, path)
		relNorm := filepath.ToSlash(rel)
		baseName := d.Name()

		if isFileMatching(opts, baseName, relNorm) {
			results = append(results, relNorm)
		}

		if opts.Limit > 0 && len(results) >= opts.Limit {
			return fs.SkipAll
		}

		return nil
	})

	return results, errWalk
}

func handleFindDirSkip(name string) error {
	if name == ".git" || name == ".tmp" || name == "node_modules" || name == ".gitmap" {
		return fs.SkipDir
	}

	return nil
}

func handleFindWalkErr(err error) error {
	if err == nil || os.IsNotExist(err) {
		return nil
	}

	return err
}

func outputFindResults(matches []string, isJson bool, elapsed time.Duration) error {
	if !isJson {
		fmt.Printf("Search completed: %d hit(s) in %v\n", len(matches), elapsed)
		for _, m := range matches {
			fmt.Println(m)
		}

		return nil
	}

	b, err := json.MarshalIndent(matches, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func showFindFilesHelp(kind MatchKind) {
	switch kind {
	case MatchExact:
		printFindFilesExactHelp()
	case MatchContains:
		printFindFilesAnyHelp()
	case MatchStartsWith:
		printFindFilesStartHelp()
	case MatchEndsWith:
		printFindFilesEndHelp()
	case MatchWildcard:
		printFindWildcardHelp()
	}
}

func printFindFilesExactHelp() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap find-files <filename> [-ext <exts>] [--limit <n>] [--json] [--dir <path>]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Description:" + constants.ColorReset)
	fmt.Println("  Find files matching exact filename with optional extension filtering.")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Aliases:" + constants.ColorReset)
	fmt.Println("  gitmap ff <filename>")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println("  gitmap find-files \"constants.go\"")
	fmt.Println("  gitmap find-files \"package.json\" --json")
	fmt.Println("  gitmap find-files \"index\" -ext \"ts,tsx\" --limit 10")
	fmt.Println("  gitmap ff \"README.md\"")
}

func printFindFilesAnyHelp() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap find-files-any <substring> [-ext <exts>] [--limit <n>] [--json] [--dir <path>]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Description:" + constants.ColorReset)
	fmt.Println("  Find files containing substring in their name with optional extension filtering.")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Aliases:" + constants.ColorReset)
	fmt.Println("  gitmap ffa <substring>")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println("  gitmap find-files-any \"runner\" -ext \"py,go\"")
	fmt.Println("  gitmap find-files-any \"config\" --json")
	fmt.Println("  gitmap ffa \"pipeline\" -ext \"go\"")
}

func printFindFilesStartHelp() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap find-files-startswith <prefix> [-ext <exts>] [--limit <n>] [--json] [--dir <path>]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Description:" + constants.ColorReset)
	fmt.Println("  Find files whose name starts with the given prefix.")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Aliases:" + constants.ColorReset)
	fmt.Println("  gitmap ffs <prefix>")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println("  gitmap find-files-startswith \"01-\" -ext \"md\"")
	fmt.Println("  gitmap find-files-startswith \"test_\" -ext \"py\"")
	fmt.Println("  gitmap ffs \"build\" --json")
}

func printFindFilesEndHelp() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap find-files-endswith <suffix> [-ext <exts>] [--limit <n>] [--json] [--dir <path>]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Description:" + constants.ColorReset)
	fmt.Println("  Find files whose name ends with the given suffix.")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Aliases:" + constants.ColorReset)
	fmt.Println("  gitmap ffe <suffix>")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println("  gitmap find-files-endswith \"_test.go\" -ext \"go\"")
	fmt.Println("  gitmap find-files-endswith \".spec.ts\"")
	fmt.Println("  gitmap ffe \"_test\" -ext \"py,go\"")
}

func printFindWildcardHelp() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap find <wildcard-pattern> [-ext <exts>] [--limit <n>] [--json] [--dir <path>]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Description:" + constants.ColorReset)
	fmt.Println("  Universal wildcard and glob file search (*pattern, pattern*, *pattern*).")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Aliases:" + constants.ColorReset)
	fmt.Println("  gitmap f <wildcard-pattern>")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println("  gitmap find \"*config*.json\"")
	fmt.Println("  gitmap find \"*_test.go\" -ext \"go\"")
	fmt.Println("  gitmap f \"*docker*\" --limit 5")
	fmt.Println("  gitmap f \"*.yaml\" --json")
}
