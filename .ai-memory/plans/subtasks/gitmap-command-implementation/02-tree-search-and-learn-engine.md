# Engineering Subtask Plan: GitMap Tree Search Family & Tree Learn Split-DB Engine

**Subtask ID:** `02-tree-search-and-learn-engine`  
**Subtask Code:** `Task-02-TreeSearch-TreeLearn`  
**Parent Plan:** `gitmap-command-implementation` ([gitmap-command-implementation.md](../../gitmap-command-implementation.md))  
**Run Number:** 83  
**Assigned Agent Role:** Spec Author 02 / Implementation Worker 02  
**Spec Reference:**  
- [02-component-spec.md](../../../../02-spec/21-app/gitmap-command-implementation/02-component-spec.md)  
- [01-architecture-spec.md](../../../../02-spec/21-app/gitmap-command-implementation/01-architecture-spec.md)  
- [Master Plan: GitMap Command Implementation](../../gitmap-command-implementation.md)  
**Status:** `READY FOR IMPLEMENTATION`  

---

## 1. Targeted File Inventory

| Component Area | Primary Implementation Files | Test & Documentation Files |
| :--- | :--- | :--- |
| **Data Models & Universal Parser** | `cli/cmdfoldertree/treesearch_parser.go` | `cli/cmdfoldertree/treesearch_test.go` |
| **Tree Search Match Engine** | `cli/cmdfoldertree/treesearch_match.go` | `cli/cmdfoldertree/treesearch_test.go` |
| **SQLite Split-DB Engine** | `cli/cmdfoldertree/treesearch_db.go` | `cli/cmdfoldertree/treesearch_test.go` |
| **Terminal Help UI Display** | `cli/cmdfoldertree/treesearch_help.go` | `cli/cmdfoldertree/treesearch_test.go` |
| **Command Handlers & Flag Parsing** | `cli/cmdfoldertree/treesearch_cmd.go` | `cli/cmdfoldertree/treesearch_test.go` |
| **CLI Dispatch Integration** | `cli/cmdfoldertree/dispatch_folder_tree.go`, `cli/cmd/roottooling.go` | `cli/cmdfoldertree/treesearch_test.go` |

---

## 2. Disjoint File Ownership & Strict Boundaries

To guarantee multi-agent safety and eliminate concurrency conflicts:
- **Spec Author 02 Owned Files (Authoring Phase):**
  * `02-spec/21-app/gitmap-command-implementation/02-component-spec.md`
  * `.ai-memory/plans/subtasks/gitmap-command-implementation/02-tree-search-and-learn-engine.md`
- **Worker 02 Downstream Owned Implementation Files:**
  * `cli/cmdfoldertree/treesearch_cmd.go`
  * `cli/cmdfoldertree/treesearch_parser.go`
  * `cli/cmdfoldertree/treesearch_match.go`
  * `cli/cmdfoldertree/treesearch_db.go`
  * `cli/cmdfoldertree/treesearch_help.go`
  * `cli/cmdfoldertree/dispatch_folder_tree.go`
  * `cli/cmd/roottooling.go`
  * `cli/cmdfoldertree/treesearch_test.go`
- **Worker 01 Disjoint Owned Files (Hands Off for Worker 02):**
  * `cli/cmd/folder/folder.go`
  * `cli/cmd/folder/parse_args.go`
  * `cli/cmd/folder/filter.go`
  * `cli/cmd/folder/folder_test.go`
- **Strictly Prohibited Actions:**
  * **TOTAL BAN ON GIT COMMANDS:** Subagents MUST NOT execute `git add`, `git commit`, `git status`, `git push`, or `git diff`.
  * **Strictly Relative Paths Only:** All file references across documentation, code comments, and test fixtures must use forward-slash paths relative to repository root (`02-spec/...`, `.ai-memory/...`, `cli/...`). Zero absolute filesystem paths or `file:///` URIs.
  * **Search Exclusively via GitMap:** Use `gitmap aum search`, `gitmap find`, `gitmap lf`, and `gitmap cat`. Strict ban on `grep`, `ripgrep`, `rg`, `Select-String`, or `findstr`.
  * **No Build / No Test in Spec Phase:** Spec authoring requires zero build or test execution.

---

## 3. Step-by-Step Implementation Tasks

### Step 1: Implement Universal Tree Parser in `cli/cmdfoldertree/treesearch_parser.go`
* **Objective:** Parse `.txt` (ASCII tree + flat lists), `.json` (`FolderReport` and `FolderTreeExportDoc`), and `.yaml` formats into normalized `[]TreeFileRecord` structs.
* **Exact Modifications:**
  1. Define data model:
     ```go
     type TreeFileRecord struct {
         ID        int64     `json:"id,omitempty"`
         RelPath   string    `json:"relPath"`
         FileName  string    `json:"fileName"`
         DirPath   string    `json:"dirPath"`
         Extension string    `json:"extension"`
         SizeBytes int64     `json:"sizeBytes"`
         IsDir     bool      `json:"isDir"`
         Depth     int       `json:"depth"`
         UpdatedAt time.Time `json:"updatedAt"`
     }
     ```
  2. Implement `ParseTreeFile(filePath string) ([]TreeFileRecord, error)`:
     - Stat file existence. If not found, return `ErrTreeFileNotFound`.
     - Detect format via `filepath.Ext(filePath)` with content sniffing fallback.
     - Dispatch to `parseTextTreeFile`, `parseJSONTreeFile`, or `parseYAMLTreeFile`.
  3. Implement `parseTextTreeFile(content []byte) ([]TreeFileRecord, error)`:
     - Scan lines. If lines start with visual tree glyphs (`├── `, `└── `, `│   `, indentation):
       - Maintain `dirStack []treeStackEntry` tracking depth.
       - Reconstruct full `RelPath` by combining stack directory path and node name.
     - If lines are flat paths (e.g. `cli/cmd/root.go`):
       - Normalize `\` to `/`.
       - Compute `FileName`, `DirPath`, `Extension`, `Depth`.
  4. Implement `parseJSONTreeFile(content []byte) ([]TreeFileRecord, error)`:
     - Attempt unmarshaling into `FolderReport` (`files: []*FileMeta`). If matched, map records.
     - Attempt unmarshaling into `FolderTreeExportDoc` (`tree: *FolderTreeNode`). If matched, walk tree recursively.
     - Fall back to raw array of path strings `[]string`.
  5. Implement `parseYAMLTreeFile(content []byte) ([]TreeFileRecord, error)`:
     - Parse YAML representation of reports or path documents into `[]TreeFileRecord`.
* **Acceptance Criteria:**
  - Accurately parses ASCII box-drawing trees, preserving deep directory nesting.
  - Correctly extracts files from both `FolderReport` and `FolderTreeExportDoc` JSON payloads.
  - Returns positive boolean `IsDir` correctly on every entry.

---

### Step 2: Implement Tree Search Match Engine in `cli/cmdfoldertree/treesearch_match.go`
* **Objective:** Provide high-speed pattern matching across records for wildcard, prefix (startsWith), substring (contains), suffix (endsWith), and regex (grep) modes.
* **Exact Modifications:**
  1. Define filter configuration:
     ```go
     type TreeSearchFilter struct {
         Pattern       string
         SearchMode    string // "wildcard", "startsWith", "contains", "endsWith", "grep"
         DirsOnly      bool
         FilesOnly     bool
         CaseSensitive bool
         Limit         int
     }
     ```
  2. Implement match evaluation functions:
     - `MatchRecords(records []TreeFileRecord, filter TreeSearchFilter) ([]TreeFileRecord, error)`
     - `matchWildcard(record TreeFileRecord, pattern string, caseSensitive bool) bool`:
       - Handles glob wildcard `*` and `?`.
       - If pattern has `/`, tests `record.RelPath`. If no `/`, tests `record.FileName` or `record.RelPath`.
     - `matchStartsWith(record TreeFileRecord, pattern string, caseSensitive bool) bool`:
       - Prefix match: `strings.HasPrefix(target, pattern)`.
       - Folder path preservation: If prefix matches directory or start of path, matches all descendant records.
     - `matchContains(record TreeFileRecord, pattern string, caseSensitive bool) bool`:
       - Substring match: `strings.Contains(strings.ToLower(record.RelPath), strings.ToLower(pattern))`.
     - `matchEndsWith(record TreeFileRecord, pattern string, caseSensitive bool) bool`:
       - Suffix match: `strings.HasSuffix(target, pattern)` or `strings.EqualFold(record.Extension, pattern)`.
     - `matchGrep(record TreeFileRecord, re *regexp.Regexp) bool`:
       - Regular expression match against `record.RelPath`.
  3. Implement output formatting:
     - `FormatMatchesText(matches []TreeFileRecord, filter TreeSearchFilter) string`:
       - Formats with color glyphs: `[DIR]` (Cyan) and `[FILE]` (Green) with relative paths and byte sizes.
     - `FormatMatchesJSON(matches []TreeFileRecord) (string, error)`:
       - Emits clean JSON array.
* **Acceptance Criteria:**
  - Case-insensitive by default.
  - Full relative path hierarchy preserved on prefix matching.
  - Safe error return on invalid regular expressions without crashing.

---

### Step 3: Implement SQLite Split-DB Storage Engine in `cli/cmdfoldertree/treesearch_db.go`
* **Objective:** Store parsed tree records in `.gitmap/data/treedb/sql.db` (or `.gitmap/treedb/tree_cache.db`), establishing normalized `TreeFile` table, NOCASE indexes, and fast lookup views.
* **Exact Modifications:**
  1. Path resolution `resolveTreeDBPath() string`:
     - Primary: `.gitmap/data/treedb/sql.db`.
     - Fallback: `.gitmap/treedb/tree_cache.db`.
     - User home fallback: `~/.gitmap/data/treedb/sql.db`.
  2. Database opener `OpenTreeDB(dbPath string) (*sql.DB, error)`:
     - Use pure-Go SQLite driver `sql.Open("sqlite", dsn)`.
     - Apply PRAGMAs: WAL mode, busy timeout 5000ms, synchronous NORMAL.
     - Set connection limits: `db.SetMaxOpenConns(1)`.
  3. Schema migration `initTreeDBSchema(db *sql.DB) error`:
     - Create table `TreeFile` with fields `(id, RelPath UNIQUE, FileName, DirPath, Extension, SizeBytes, IsDir, Depth, UpdatedAt)`.
     - Create NOCASE indexes on `RelPath`, `FileName`, `DirPath`, `Extension`, `IsDir`.
     - Create database view `ViewTreeFiles`.
  4. Ingestion routine `IngestTreeRecords(db *sql.DB, records []TreeFileRecord) (*IngestStats, error)`:
     - Wrap in single atomic transaction (`db.Begin()`).
     - Execute batch upserts in chunks of 500:
       `INSERT INTO TreeFile (...) VALUES (...) ON CONFLICT(RelPath) DO UPDATE SET ...`
     - Return `IngestStats`: `TotalRecords`, `DurationMs`, `DBPath`.
  5. Query primitives for DB-backed search:
     - `QueryTreeDB(db *sql.DB, filter TreeSearchFilter) ([]TreeFileRecord, error)`.
* **Acceptance Criteria:**
  - Database schema initializes idempotently.
  - Batch ingestion completes under 80ms for 10,000 records.
  - Queries execute via index scans in $< 10\text{ms}$.

---

### Step 4: Implement Terminal Help UI in `cli/cmdfoldertree/treesearch_help.go`
* **Objective:** Deliver rich ANSI terminal help screens with clear syntax, options, and practical examples for all commands.
* **Exact Modifications:**
  1. Implement `ShowTreeSearchHelp(subcommand string)`:
     - When subcommand is empty or `"all"`, print the master suite overview.
     - When subcommand is `"tree-search"`, print wildcard glob help.
     - When subcommand is `"tree-search-startsWith"`, print prefix matching help.
     - When subcommand is `"tree-search-contains"`, print substring search help.
     - When subcommand is `"tree-search-endsWith"`, print suffix/extension search help.
     - When subcommand is `"tree-search-grep"`, print regex grep search help.
     - When subcommand is `"tree-learn"`, print Split-DB ingestion help.
  2. Adhere to GitMap ANSI styling:
     - Section titles in Cyan/Bold.
     - Command badges in Green.
     - Flag parameters in Yellow.
     - Examples with comments.
* **Acceptance Criteria:**
  - Triggered via `-h`, `--help`, or `help` flag on any command.
  - Formatted cleanly across 80-column and wide terminal windows.

---

### Step 5: Implement CLI Execution Handlers in `cli/cmdfoldertree/treesearch_cmd.go`
* **Objective:** Parse CLI arguments, resolve flags, coordinate parsing/matching/storing, and render output.
* **Exact Modifications:**
  1. Argument and flag parser:
     - Extract `-f` / `--file <path>`.
     - Extract `--db <path>`.
     - Extract `--json`.
     - Extract `-n` / `--limit <int>`.
     - Extract `--dirs` / `--files`.
     - Extract positional search pattern.
  2. Handlers:
     - `RunTreeSearchDispatch(subcommand string, args []string) error`:
       - Check for help flags (`-h`, `--help`). If present, call `ShowTreeSearchHelp(subcommand)`.
       - Map subcommand to search mode (`"wildcard"`, `"startsWith"`, `"contains"`, `"endsWith"`, `"grep"`).
       - If `-f` provided, call `ParseTreeFile`, then `MatchRecords`.
       - If `-f` omitted, check for existing `.gitmap/data/treedb/sql.db` and query DB directly.
       - Output results via text table or JSON.
     - `RunTreeLearn(args []string) error`:
       - Check for help flags.
       - Validate `-f` is provided. If missing, return error explaining `-f` is required.
       - Call `ParseTreeFile`, then `OpenTreeDB`, then `IngestTreeRecords`.
       - Print success badge with record count and elapsed duration.
* **Acceptance Criteria:**
  - Exits with code 0 on successful search or ingestion.
  - Displays user-friendly error messages if file does not exist or pattern is missing.

---

### Step 6: Wire Command Dispatch & Root Tooling
* **Objective:** Integrate all new commands into CLI routing layers.
* **Exact Modifications:**
  1. In `cli/cmdfoldertree/dispatch_folder_tree.go`:
     - Update `DispatchFolderTree(command string) (bool, error)`:
       - Detect `isTreeSearchCommand(command)` $\rightarrow$ call `RunTreeSearchDispatch(command, os.Args[2:])`.
       - Detect `isTreeLearnCommand(command)` $\rightarrow$ call `RunTreeLearn(os.Args[2:])`.
     - Add helper predicates:
       - `isTreeSearchCommand(cmd string) bool` matching `tree-search`, `tree-search-startsWith`, `tree-search-contains`, `tree-search-endsWith`, `tree-search-grep`, and aliases.
       - `isTreeLearnCommand(cmd string) bool` matching `tree-learn`, `treelearn`, `tl`.
  2. In `cli/cmd/roottooling.go`:
     - Add entries to `toolingUtilEntries()`:
       ```go
       {[]string{"tree-search", "treesearch", "ts"}, func() error { return cmdfoldertree.RunTreeSearchDispatch("tree-search", argsTail()) }},
       {[]string{"tree-search-startsWith", "tree-search-startswith", "tree-search-starts-with", "tree-search-prefix"}, func() error { return cmdfoldertree.RunTreeSearchDispatch("tree-search-startsWith", argsTail()) }},
       {[]string{"tree-search-contains", "tree-search-contain", "tree-search-substr"}, func() error { return cmdfoldertree.RunTreeSearchDispatch("tree-search-contains", argsTail()) }},
       {[]string{"tree-search-endsWith", "tree-search-endswith", "tree-search-ends-with", "tree-search-suffix"}, func() error { return cmdfoldertree.RunTreeSearchDispatch("tree-search-endsWith", argsTail()) }},
       {[]string{"tree-search-grep", "tree-search-regex"}, func() error { return cmdfoldertree.RunTreeSearchDispatch("tree-search-grep", argsTail()) }},
       {[]string{"tree-learn", "treelearn", "tl"}, func() error { return cmdfoldertree.RunTreeLearn(argsTail()) }},
       ```
* **Acceptance Criteria:**
  - Both direct dispatch via `root.go` and table lookup via `roottooling.go` resolve correctly.
  - All command variants and casing aliases function identically.

---

### Step 7: Unit Testing & Verification Fixtures in `cli/cmdfoldertree/treesearch_test.go`
* **Objective:** Implement comprehensive unit test coverage without running tests in the spec authoring turn.
* **Exact Modifications:**
  - Test cases for ASCII tree parsing with indentation directory stack.
  - Test cases for flat path list parsing.
  - Test cases for JSON `FolderReport` and `FolderTreeExportDoc` unmarshaling.
  - Test cases for Wildcard, StartsWith, Contains, EndsWith, and Grep matching.
  - Test cases for in-memory SQLite Split-DB creation, batch upsert, and NOCASE queries.
  - Test cases for CLI argument parsing and error propagation.
* **Acceptance Criteria:**
  - 100% logic coverage across parser, matcher, and database modules.

---

## 4. Verification & Scorecard Gates

| Gate ID | Target Command | Acceptance Criteria |
| :--- | :--- | :--- |
| **VG-05** | `gitmap tree-search -f "a.json" "*pattern*"` | Exit code 0, returns files matching wildcard pattern, supports `*` and `?`. |
| **VG-06** | `gitmap tree-search-startsWith -f "a.json" "cli"` | Exit code 0, matches paths starting with `cli`, preserves folder path hierarchy. |
| **VG-07** | `gitmap tree-search-contains -f "a.json" "root"` | Exit code 0, matches paths containing substring `root` case-insensitively. |
| **VG-08** | `gitmap tree-search-endsWith -f "a.json" ".go"` | Exit code 0, matches paths ending in `.go`. |
| **VG-09** | `gitmap tree-search-grep -f "a.json" "cmd.*\.go"` | Exit code 0, executes regular expression match, handles invalid regex with error code. |
| **VG-10** | `gitmap tree-learn -f "a.json"` | Exit code 0, creates `.gitmap/data/treedb/sql.db`, creates `TreeFile` table and indexes, inserts records transactionally, outputs summary badge. |
