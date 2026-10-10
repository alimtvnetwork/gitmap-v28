# Plan: GitMap Tree Commands & Tree Search/Learn Engine Implementation

## User Request (Verbatim)
```text
# GitMap Command Implementation: high priority instruction, non-negotiable task

make sure these commands are implemented and also

gitmap tree --file/f "a.txt"
gitmap tree --file/f "a.json"
gitmap tree --file/f "a.yaml"

gitmap tree-search -f "a.txt/a.json" "search pattern*ends*" # * meaning wildcard
gitmap tree-search-startsWith -f "a.txt/a.json" "search pattern" # meaning applying text and ends with * that means don't care, also respect fodler path if given
gitmap tree-search-contains -f "a.txt/a.json" "search pattern"
gitmap tree-search-endsWith -f "a.txt/a.json" "search pattern"
gitmap tree-search-grep -f "a.txt/a.json" "search pattern"
gitmap tree-learn -f "a.txt/a.json"  # save as split db with proper info to find files fast, clear???

Please also include the help text

## slug: gitmap-command-implementation
```

---

## 1. Executive Summary & Goals
This project implements the requested GitMap tree export, tree search, and tree learn command suites, and resolves the crash when running `gitmap tree "*.go"`:

1. **`gitmap tree` Enhancements:**
   - Add `--file` / `-f` flags supporting export to `.txt`, `.json`, and `.yaml` with automatic format deduction.
   - Fix Windows path syntax crash by intercepting wildcard/glob positional arguments (e.g. `gitmap tree "*.go"`) and routing them to `FilterConfig` while defaulting `opts.TargetDir` to `.`.
2. **`gitmap tree-search` Family:**
   - `tree-search`: Wildcard pattern matching (`*`).
   - `tree-search-startsWith`: Prefix matching with folder path preservation.
   - `tree-search-contains`: Substring matching.
   - `tree-search-endsWith`: Suffix matching.
   - `tree-search-grep`: Regular expression matching.
   - Universal parser for `-f` input supporting plain text tree/flat list, JSON reports, and YAML.
3. **`gitmap tree-learn` Split SQLite DB Engine:**
   - Ingests parsed tree files into `.gitmap/data/treedb/sql.db` (or `.gitmap/treedb/tree_cache.db`).
   - Creates `TreeFile` table with NOCASE indexes on `RelPath`, `FileName`, `DirPath`, and `Extension`.
   - Fast retrieval queries and index telemetry.
4. **CLI Help & Dispatch Wiring:**
   - Wire all commands into `cli/cmd/roottooling.go` and `cli/cmdfoldertree/dispatch_folder_tree.go`.
   - Comprehensive help text menus and examples for all new commands.

---

## 2. Architecture & File Ownership Matrix

| Phase / Worker | Strictly Owned Files | Responsibilities |
|:---|:---|:---|
| **Spec Author 01** | `02-spec/21-app/gitmap-command-implementation/01-architecture-spec.md`<br>`.ai-memory/plans/subtasks/gitmap-command-implementation/01-tree-flags-and-glob-fix.md` | Architecture specification & subtask plan for tree export flags & glob fix |
| **Spec Author 02** | `02-spec/21-app/gitmap-command-implementation/02-component-spec.md`<br>`.ai-memory/plans/subtasks/gitmap-command-implementation/02-tree-search-and-learn-engine.md` | Component specification & subtask plan for tree search, tree learn, and DB |
| **Worker 01** | `cli/cmd/folder/folder.go`<br>`cli/cmd/folder/parse_args.go`<br>`cli/cmd/folder/filter.go`<br>`cli/cmd/folder/folder_test.go` | Implement `--file`/`-f`, fix glob crash, support IncludeGlobs, unit tests |
| **Worker 02** | `cli/cmdfoldertree/treesearch_cmd.go`<br>`cli/cmdfoldertree/treesearch_parser.go`<br>`cli/cmdfoldertree/treesearch_match.go`<br>`cli/cmdfoldertree/treesearch_db.go`<br>`cli/cmdfoldertree/treesearch_help.go`<br>`cli/cmdfoldertree/dispatch_folder_tree.go`<br>`cli/cmd/roottooling.go` | Implement tree-search family, tree-learn Split DB engine, dispatch wiring, help text |
| **Lead Orchestrator** | `.ai-memory/plans/gitmap-command-implementation.md`<br>`.ai-memory/plans/readme.md`<br>`02-spec/21-app/readme.md`<br>`.ai-memory/temp-agents/83-gitmap-command-implementation/` | Orchestration, targeted checks, git commits (`gitmap cpf`) |

---

## 3. Subtask Breakdown

- **Subtask 1 (`Task-01`):** `gitmap tree` `--file/-f` export and glob/pattern crash fix in `cli/cmd/folder/`.
- **Subtask 2 (`Task-02`):** `gitmap tree-search` suite (`tree-search`, `startsWith`, `contains`, `endsWith`, `grep`) and multi-format input parser in `cli/cmdfoldertree/`.
- **Subtask 3 (`Task-03`):** `gitmap tree-learn` Split SQLite DB storage and high-speed query indexing.
- **Subtask 4 (`Task-04`):** CLI dispatch integration, comprehensive help text, end-to-end verification, and atomic GitMap commit.

---

## 4. Verification & Scorecard Gates
- **VG-01:** `gitmap tree "*.go"` executes without crash (exit 0) and filters `.go` files.
- **VG-02:** `gitmap tree --file "a.txt"` creates plain text tree output.
- **VG-03:** `gitmap tree -f "a.json"` creates valid JSON tree report.
- **VG-04:** `gitmap tree -f "a.yaml"` creates valid YAML tree report.
- **VG-05:** `gitmap tree-search -f "a.json" "*pattern*"` executes wildcard match.
- **VG-06:** `gitmap tree-search-startsWith -f "a.json" "cli"` matches paths starting with `cli/`.
- **VG-07:** `gitmap tree-search-contains -f "a.json" "root"` matches paths containing `root`.
- **VG-08:** `gitmap tree-search-endsWith -f "a.json" ".go"` matches paths ending in `.go`.
- **VG-09:** `gitmap tree-search-grep -f "a.json" "cmd.*\.go"` executes regex search.
- **VG-10:** `gitmap tree-learn -f "a.json"` writes records to SQLite Split DB with exit code 0.
