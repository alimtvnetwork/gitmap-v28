package cmdfoldertree

import (
	"fmt"
	"strings"
)

// ShowTreeSearchHelp renders terminal help text for the tree-search suite or specific subcommand.
func ShowTreeSearchHelp(subcommand string) {
	norm := strings.ToLower(strings.TrimSpace(subcommand))
	norm = strings.ReplaceAll(norm, "_", "-")

	switch norm {
	case "tree-search-startswith", "tree-search-starts-with", "startswith", "prefix":
		printHelpStartsWith()
	case "tree-search-contains", "tree-search-contain", "contains", "substr":
		printHelpContains()
	case "tree-search-endswith", "tree-search-ends-with", "endswith", "suffix":
		printHelpEndsWith()
	case "tree-search-grep", "tree-search-regex", "grep", "regex":
		printHelpGrep()
	case "tree-learn", "treelearn", "learn":
		printHelpTreeLearn()
	case "tree-search", "treesearch", "wildcard", "all", "":
		printHelpMaster()
	default:
		printHelpMaster()
	}
}

func printHelpMaster() {
	fmt.Print(`
  ╔═════════════════════════════════════════════════════════════════════════════╗
  ║              GITMAP TREE SEARCH & LEARN COMMAND SUITE                       ║
  ╚═════════════════════════════════════════════════════════════════════════════╝

  High-speed pattern search, regular expression grep, and SQLite Split-DB
  indexing across exported directory trees, hierarchies, and path reports.

  USAGE:
    gitmap tree-search -f <file> <pattern*ends*>       Wildcard pattern search (* and ?)
    gitmap tree-search-startsWith -f <file> <prefix>   Prefix search (preserves folder path)
    gitmap tree-search-contains -f <file> <substring>  Case-insensitive substring search
    gitmap tree-search-endsWith -f <file> <suffix>     Suffix or file extension search
    gitmap tree-search-grep -f <file> <regex>          Regular expression grep
    gitmap tree-learn -f <file>                        Ingest tree into SQLite Split DB

  OPTIONS:
    -f, --file <path>        Path to input tree file (.json, .yaml, .txt)
    --db <path>              Override SQLite Split-DB path (default: .gitmap/data/treedb/sql.db)
    --json                   Format output as clean JSON array
    -n, --limit <int>        Limit maximum number of results (default: 1000)
    --dirs                   Filter to directory entries only
    --files                  Filter to file entries only
    -h, --help               Display command help information

  SUPPORTED TREE FORMATS:
    .json                    FolderReport, FolderTreeExportDoc, []TreeItem, or path array
    .yaml, .yml              YAML document, report, or hierarchy tree
    .txt, .log               Indented ASCII tree (box glyphs) or flat relative path list

  EXAMPLES:
    gitmap tree-search -f "tree.json" "*root*.go"
    gitmap tree-search-startsWith -f "tree.txt" "cli/cmd"
    gitmap tree-search-contains -f "export.yaml" "pipeline"
    gitmap tree-search-endsWith -f "report.json" ".go"
    gitmap tree-search-grep -f "tree.json" "^cli/.*_test\.go$"
    gitmap tree-learn -f "report.json"
    gitmap tree-search "*cmd*"                         (queries Split DB if -f omitted)

`)
}

func printHelpStartsWith() {
	fmt.Print(`
  COMMAND: gitmap tree-search-startsWith

  Evaluates prefix matching against directory hierarchies. If a folder path is
  specified (e.g. "cli/cmd"), all descendant files and directories are matched
  and their full relative path hierarchy is preserved.

  USAGE:
    gitmap tree-search-startsWith -f <file> <prefix> [flags]

  ALIASES:
    tree-search-startswith, tree-search-prefix, tss

  FLAGS:
    -f, --file <path>        Input tree file (.json, .yaml, .txt)
    --json                   Emit JSON array output
    -n, --limit <int>        Limit number of results
    --dirs                   Match directories only
    --files                  Match files only
    -h, --help               Show this help screen

  EXAMPLES:
    gitmap tree-search-startsWith -f "tree.json" "cli/cmd"
    gitmap tree-search-startsWith -f "export.txt" "docs"
    gitmap tree-search-startsWith -f "tree.yaml" "src/components" --dirs

`)
}

func printHelpContains() {
	fmt.Print(`
  COMMAND: gitmap tree-search-contains

  Performs case-insensitive substring searching across relative paths and file names.

  USAGE:
    gitmap tree-search-contains -f <file> <substring> [flags]

  ALIASES:
    tree-search-contain, tree-search-substr, tsc

  FLAGS:
    -f, --file <path>        Input tree file (.json, .yaml, .txt)
    --json                   Emit JSON array output
    -n, --limit <int>        Limit number of results
    --dirs                   Match directories only
    --files                  Match files only
    -h, --help               Show this help screen

  EXAMPLES:
    gitmap tree-search-contains -f "tree.json" "root"
    gitmap tree-search-contains -f "tree.txt" "sqlite"
    gitmap tree-search-contains -f "export.yaml" "test" --files

`)
}

func printHelpEndsWith() {
	fmt.Print(`
  COMMAND: gitmap tree-search-endsWith

  Matches files and directories ending with the specified suffix or file extension.

  USAGE:
    gitmap tree-search-endsWith -f <file> <suffix> [flags]

  ALIASES:
    tree-search-endswith, tree-search-suffix, tse

  FLAGS:
    -f, --file <path>        Input tree file (.json, .yaml, .txt)
    --json                   Emit JSON array output
    -n, --limit <int>        Limit number of results
    --dirs                   Match directories only
    --files                  Match files only
    -h, --help               Show this help screen

  EXAMPLES:
    gitmap tree-search-endsWith -f "tree.json" ".go"
    gitmap tree-search-endsWith -f "tree.txt" "_test.go"
    gitmap tree-search-endsWith -f "export.yaml" ".json"

`)
}

func printHelpGrep() {
	fmt.Print(`
  COMMAND: gitmap tree-search-grep

  Evaluates regular expression patterns against relative paths and file names.

  USAGE:
    gitmap tree-search-grep -f <file> <regex> [flags]

  ALIASES:
    tree-search-regex, tsg

  FLAGS:
    -f, --file <path>        Input tree file (.json, .yaml, .txt)
    --json                   Emit JSON array output
    -n, --limit <int>        Limit number of results
    --dirs                   Match directories only
    --files                  Match files only
    -h, --help               Show this help screen

  EXAMPLES:
    gitmap tree-search-grep -f "tree.json" "cmd.*\.go$"
    gitmap tree-search-grep -f "tree.txt" "^cli/(cmd|cmdfoldertree)/"
    gitmap tree-search-grep -f "export.yaml" "v[0-9]+\.json$"

`)
}

func printHelpTreeLearn() {
	fmt.Print(`
  COMMAND: gitmap tree-learn

  Parses any supported directory tree export file and persists all entries into
  a normalized Split SQLite database (.gitmap/data/treedb/sql.db). Enables
  instant sub-millisecond querying without re-reading source files.

  USAGE:
    gitmap tree-learn -f <file> [flags]

  ALIASES:
    treelearn, tl

  FLAGS:
    -f, --file <path>        Required input tree file (.json, .yaml, .txt)
    --db <path>              Override SQLite database path
    -h, --help               Show this help screen

  DATABASE SCHEMA:
    Table: TreeFile (TreeFileId, RootPath, RelPath, FileName, DirPath, Extension, SizeBytes, IsDir, Depth, UpdatedAt)
    Indexes: RelPath (NOCASE), FileName (NOCASE), DirPath (NOCASE), Extension (NOCASE), IsDir

  EXAMPLES:
    gitmap tree-learn -f "tree.json"
    gitmap tree-learn -f "tree.yaml" --db ".gitmap/data/treedb/custom.db"
    gitmap tree-learn -f "tree.txt"

`)
}
