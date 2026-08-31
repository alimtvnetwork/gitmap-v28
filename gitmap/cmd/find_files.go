package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type MatchKind int

const (
	MatchExact MatchKind = iota
	MatchContains
	MatchStartsWith
	MatchEndsWith
	MatchWildcard
)

type FindFilesOptions struct {
	Query     string
	Kind      MatchKind
	Exts      []string
	Limit     int
	IsJson    bool
	TargetDir string
}

func runFindFiles(args []string) error {
	checkHelp(constants.CmdFindFiles, args)
	return executeFindFiles(args, MatchExact)
}

func runFindFilesAny(args []string) error {
	checkHelp(constants.CmdFindFilesAny, args)
	return executeFindFiles(args, MatchContains)
}

func runFindFilesStartsWith(args []string) error {
	checkHelp(constants.CmdFindFilesStartsWith, args)
	return executeFindFiles(args, MatchStartsWith)
}

func runFindFilesEndsWith(args []string) error {
	checkHelp(constants.CmdFindFilesEndsWith, args)
	return executeFindFiles(args, MatchEndsWith)
}

func runListFiles(args []string) error {
	checkHelp(constants.CmdListFiles, args)
	return executeFindFiles(args, MatchWildcard)
}

func executeFindFiles(args []string, defaultKind MatchKind) error {
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

	return outputFindResults(matches, opts.IsJson)
}

func parseFindFilesOptions(args []string, defaultKind MatchKind) FindFilesOptions {
	opts := FindFilesOptions{
		Kind:      defaultKind,
		TargetDir: ".",
	}
	var nonFlags []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case (arg == "-ext" || arg == "--ext") && i+1 < len(args):
			opts.Exts = parseExtensionList(args[i+1])
			i++
		case strings.HasPrefix(arg, "-ext="):
			opts.Exts = parseExtensionList(strings.TrimPrefix(arg, "-ext="))
		case strings.HasPrefix(arg, "--ext="):
			opts.Exts = parseExtensionList(strings.TrimPrefix(arg, "--ext="))
		case (arg == "--limit" || arg == "-l") && i+1 < len(args):
			opts.Limit, _ = strconv.Atoi(args[i+1])
			i++
		case arg == "--json" || arg == "-json":
			opts.IsJson = true
		case (arg == "--dir" || arg == "-d") && i+1 < len(args):
			opts.TargetDir = args[i+1]
			i++
		case !strings.HasPrefix(arg, "-"):
			nonFlags = append(nonFlags, arg)
		}
	}

	if len(nonFlags) > 0 {
		opts.Query = nonFlags[0]
	}
	return opts
}

func scanAndMatchFiles(opts FindFilesOptions) ([]string, error) {
	absRoot, err := filepath.Abs(opts.TargetDir)
	if err != nil {
		return nil, err
	}

	var results []string
	errWalk := filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, errIn error) error {
		if errIn != nil {
			return errIn
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
	if name == ".git" {
		return fs.SkipDir
	}
	return nil
}

func outputFindResults(matches []string, isJson bool) error {
	if !isJson {
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
