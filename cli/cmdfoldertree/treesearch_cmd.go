package cmdfoldertree

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// parsedArgs holds flags and positional parameters extracted from CLI arguments.
type parsedArgs struct {
	filePath      string
	dbPath        string
	pattern       string
	isJSON        bool
	limit         int
	dirsOnly      bool
	filesOnly     bool
	caseSensitive bool
	isHelp        bool
}

func parseTreeSearchArgs(args []string) parsedArgs {
	var res parsedArgs
	positionals := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-h" || arg == "--help" || arg == "help" {
			res.isHelp = true
			return res
		}

		if arg == "-f" || arg == "--file" {
			if i+1 < len(args) {
				res.filePath = args[i+1]
				i++
			}
			continue
		}
		if strings.HasPrefix(arg, "-f=") {
			res.filePath = strings.TrimPrefix(arg, "-f=")
			continue
		}
		if strings.HasPrefix(arg, "--file=") {
			res.filePath = strings.TrimPrefix(arg, "--file=")
			continue
		}

		if arg == "--db" {
			if i+1 < len(args) {
				res.dbPath = args[i+1]
				i++
			}
			continue
		}
		if strings.HasPrefix(arg, "--db=") {
			res.dbPath = strings.TrimPrefix(arg, "--db=")
			continue
		}

		if arg == "--json" {
			res.isJSON = true
			continue
		}

		if arg == "-n" || arg == "--limit" {
			if i+1 < len(args) {
				if n, err := strconv.Atoi(args[i+1]); err == nil && n > 0 {
					res.limit = n
				}
				i++
			}
			continue
		}
		if strings.HasPrefix(arg, "--limit=") {
			if n, err := strconv.Atoi(strings.TrimPrefix(arg, "--limit=")); err == nil && n > 0 {
				res.limit = n
			}
			continue
		}

		if arg == "--dirs" || arg == "--dirs-only" {
			res.dirsOnly = true
			continue
		}
		if arg == "--files" || arg == "--files-only" {
			res.filesOnly = true
			continue
		}

		if arg == "--case-sensitive" {
			res.caseSensitive = true
			continue
		}

		// Positional argument
		positionals = append(positionals, arg)
	}

	// If -f was not set with a flag, but we have 2 positional arguments and the first is an existing file
	if res.filePath == "" && len(positionals) >= 2 {
		if fi, err := os.Stat(positionals[0]); err == nil && !fi.IsDir() {
			res.filePath = positionals[0]
			positionals = positionals[1:]
		}
	}

	if len(positionals) > 0 {
		res.pattern = positionals[0]
	}

	return res
}

// RunTreeSearchDispatch routes tree search commands by subcommand name.
func RunTreeSearchDispatch(subcommand string, args []string) error {
	norm := strings.ToLower(strings.TrimSpace(subcommand))
	norm = strings.ReplaceAll(norm, "_", "-")

	switch norm {
	case "tree-search", "treesearch", "ts":
		return RunTreeSearch(args)
	case "tree-search-startswith", "tree-search-starts-with", "tree-search-prefix", "tss":
		return RunTreeSearchStartsWith(args)
	case "tree-search-contains", "tree-search-contain", "tree-search-substr", "tsc":
		return RunTreeSearchContains(args)
	case "tree-search-endswith", "tree-search-ends-with", "tree-search-suffix", "tse":
		return RunTreeSearchEndsWith(args)
	case "tree-search-grep", "tree-search-regex", "tsg":
		return RunTreeSearchGrep(args)
	case "tree-learn", "treelearn", "tl":
		return RunTreeLearn(args)
	default:
		return RunTreeSearch(args)
	}
}

// RunTreeSearch executes wildcard glob matching.
func RunTreeSearch(args []string) error {
	return executeTreeSearch(ModeWildcard, "tree-search", args)
}

// RunTreeSearchStartsWith executes prefix matching with folder path preservation.
func RunTreeSearchStartsWith(args []string) error {
	return executeTreeSearch(ModeStartsWith, "tree-search-startsWith", args)
}

// RunTreeSearchContains executes substring searching.
func RunTreeSearchContains(args []string) error {
	return executeTreeSearch(ModeContains, "tree-search-contains", args)
}

// RunTreeSearchEndsWith executes suffix and extension matching.
func RunTreeSearchEndsWith(args []string) error {
	return executeTreeSearch(ModeEndsWith, "tree-search-endsWith", args)
}

// RunTreeSearchGrep executes regular expression searching.
func RunTreeSearchGrep(args []string) error {
	return executeTreeSearch(ModeGrep, "tree-search-grep", args)
}

func executeTreeSearch(mode SearchMode, cmdName string, args []string) error {
	p := parseTreeSearchArgs(args)
	if p.isHelp {
		ShowTreeSearchHelp(cmdName)
		return nil
	}

	var items []TreeItem
	var err error

	if p.filePath != "" {
		items, err = ParseTreeFile(p.filePath)
		if err != nil {
			return fmt.Errorf("error reading tree file %s: %w", p.filePath, err)
		}
	} else {
		// Fallback to SQLite Split DB
		dbPath := ResolveTreeDBPath(p.dbPath)
		if _, statErr := os.Stat(dbPath); os.IsNotExist(statErr) {
			return fmt.Errorf("no input file specified with -f and SQLite database not found at %s. Run '%s --help'", dbPath, cmdName)
		}
		items, err = QueryTreeDb(dbPath, mode, p.pattern)
		if err != nil {
			return fmt.Errorf("database query error: %w", err)
		}
		// If query returned results and pattern was already applied, we can jump to output
		return renderTreeSearchResults(items, mode, p.pattern, p.isJSON)
	}

	filter := TreeSearchFilter{
		Pattern:       p.pattern,
		SearchMode:    mode,
		DirsOnly:      p.dirsOnly,
		FilesOnly:     p.filesOnly,
		CaseSensitive: p.caseSensitive,
		Limit:         p.limit,
	}

	matched, err := FilterTreeItems(items, filter)
	if err != nil {
		return err
	}

	return renderTreeSearchResults(matched, mode, p.pattern, p.isJSON)
}

func renderTreeSearchResults(items []TreeItem, mode SearchMode, pattern string, isJSON bool) error {
	if isJSON {
		jsonStr, err := FormatTreeItemsJSON(items)
		if err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
		fmt.Println(jsonStr)
		return nil
	}

	queryDesc := fmt.Sprintf("[%s] %q", mode, pattern)
	if strings.TrimSpace(pattern) == "" {
		queryDesc = "(all items)"
	}
	fmt.Print(FormatTreeItemsText(items, queryDesc))
	return nil
}

// RunTreeLearn parses an input tree file and stores entries into the Split SQLite database.
func RunTreeLearn(args []string) error {
	p := parseTreeSearchArgs(args)
	if p.isHelp {
		ShowTreeSearchHelp("tree-learn")
		return nil
	}

	filePath := p.filePath
	if filePath == "" && p.pattern != "" {
		// Check if positional arg is file
		if fi, err := os.Stat(p.pattern); err == nil && !fi.IsDir() {
			filePath = p.pattern
		}
	}

	if filePath == "" {
		return fmt.Errorf("input tree file path is required via -f or --file (see 'gitmap tree-learn --help')")
	}

	start := time.Now()
	items, err := ParseTreeFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse tree file %s: %w", filePath, err)
	}

	dbPath := ResolveTreeDBPath(p.dbPath)
	count, err := IngestTreeItems(dbPath, items)
	if err != nil {
		return fmt.Errorf("failed to ingest tree items into %s: %w", dbPath, err)
	}

	elapsed := time.Since(start)
	fmt.Printf("[OK] Learned %d tree records into SQLite Split DB (%s) in %s\n", count, dbPath, elapsed.Round(time.Millisecond))
	return nil
}
