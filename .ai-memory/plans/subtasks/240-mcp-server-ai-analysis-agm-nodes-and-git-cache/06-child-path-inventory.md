# Subtask 06: Child-Path Fast Inventory & Search Subsystem

- **Parent Task:** `240-mcp-server-ai-analysis-agm-nodes-and-git-cache`
- **Subtask ID:** `06-child-path-inventory`
- **Spec Reference:** `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/02-component-and-cli-spec.md` (Section 2)
- **Status:** `PENDING`
- **Assigned Subagent:** Implementation Subagent (Phase 2)

---

## 1. Objective & Scope

Implement the **Child-Path Fast Inventory and AST-Aware Search CLI Command** (`gitmap child-path`, with shorthands `gitmap cp` and `gitmap childpath`).

### The Problem
During autonomous code refactoring and speculative audits, AI coding agents and CLI scripts must frequently inventory directories, discover matching source files, and filter by extension or depth. Historically, workflows relied on:
1. **PowerShell `Get-ChildItem`:** Slow, heavy .NET object allocations (multi-second delays in deep trees), cross-platform discrepancies between Windows PowerShell 5.1 and pwsh Core, localized text formatting, and encoding breaks (documented in `assets/screenshots/240-get-childitem-filter.png`).
2. **`git grep`:** Incapable of scanning untracked files, local scratch directories, scaffolding templates, or `.ai-memory` paths excluded by `.gitignore` (documented in `assets/screenshots/240-git-grep-func-runfix.png`).

### The Solution
Implement `gitmap child-path [dir] [glob]` in pure Go utilizing low-level kernel directory walking (`os.ReadDir`), glob matching, and unified JSON/terminal output:
- Pure Go recursion with bounded depth (`--depth <N>`).
- Filtering by extension (`--ext .go,.ts,.json`) and entry type (`--type f|d|all`).
- Bounded emission limit (`--limit <N>`).
- Built-in omission of heavy non-source directories (`.git`, `node_modules`, `.cache`) unless explicitly requested.
- Forward-slash normalized relative paths across all operating systems.
- Structured JSON telemetry envelope (`--json`) for direct AI agent ingestion.

---

## 2. Concrete Files to Create / Modify

| File | Nature | Purpose |
| :--- | :--- | :--- |
| `cli/cmdchildpath/childpath_types.go` | New | Data models: `EntryTypeFilter`, `ChildPathOptions`, `ChildPathItem`, and `ChildPathResponse`. |
| `cli/cmdchildpath/childpath_walker.go` | New | High-speed directory traversal engine utilizing `os.ReadDir`, depth limits, ignore lists, and glob filters. |
| `cli/cmdchildpath/childpath_cmd.go` | New | CLI argument parsing, flags extraction (`--depth`, `--ext`, `--type`, `--limit`, `--json`, `--size`), table formatting via `termtable`, and entry point `RunChildPathCommand()`. |
| `cli/cmdchildpath/childpath_test.go` | New | Unit and benchmark tests validating filtering, depth truncation, JSON output, and performance (<15ms for 10k items). |
| `cli/cmd/rootdispatch.go` | Modify | Register `child-path`, `childpath`, and `cp` aliases in root command dispatcher. |
| `cli/constants/constants_cli.go` | Modify | Add command constants: `CmdChildPath = "child-path"` and aliases `cp`, `childpath`. |

---

## 3. Detailed Implementation Requirements

### 3.1 CLI Argument & Flag Contract
```bash
gitmap child-path [dir] [glob] [flags]
gitmap cp [dir] [glob] [flags]
```
- Flags:
  - `--depth <N>` (`-d`): Maximum recursion depth (0 = unlimited).
  - `--ext <list>` (`-e`): Comma-delimited extensions (e.g. `.go,.ts`).
  - `--type <f|d|all>` (`-t`): Filter type (`f` files, `d` dirs, `all`).
  - `--limit <N>` (`-n`, `-l`): Maximum records to output.
  - `--json` (`-j`): Emit structured JSON telemetry.
  - `--include-hidden`: Include dotfiles and hidden directories.
  - `--size` (`-s`): Display file sizes and timestamps.

### 3.2 Walker Engine (`cli/cmdchildpath/childpath_walker.go`)
```go
package cmdchildpath

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func WalkChildPath(opts ChildPathOptions) (*ChildPathResponse, error) {
	start := time.Now()
	var items []ChildPathItem
	baseDir := filepath.Clean(opts.TargetDir)

	var walk func(currentDir string, currentDepth int) error
	walk = func(currentDir string, currentDepth int) error {
		if opts.MaxDepth > 0 && currentDepth > opts.MaxDepth {
			return nil
		}

		entries, err := os.ReadDir(currentDir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			name := entry.Name()
			if !opts.IsIncludeHidden && isDefaultExcluded(name) {
				continue
			}

			fullPath := filepath.Join(currentDir, name)
			relPath, _ := filepath.Rel(baseDir, fullPath)
			relPath = filepath.ToSlash(relPath)
			isDir := entry.IsDir()

			if matchesFilter(entry, relPath, opts) {
				info, _ := entry.Info()
				var size int64
				var modTime time.Time
				if info != nil {
					size = info.Size()
					modTime = info.ModTime()
				}

				items = append(items, ChildPathItem{
					RelPath:   relPath,
					Name:      name,
					IsDir:     isDir,
					SizeBytes: size,
					ModTime:   modTime,
					Extension: strings.ToLower(filepath.Ext(name)),
					Depth:     currentDepth,
				})

				if opts.Limit > 0 && len(items) >= opts.Limit {
					return nil
				}
			}

			if isDir {
				if err := walk(fullPath, currentDepth+1); err != nil {
					return err
				}
				if opts.Limit > 0 && len(items) >= opts.Limit {
					return nil
				}
			}
		}
		return nil
	}

	err := walk(baseDir, 1)
	return &ChildPathResponse{
		TargetDir:    opts.TargetDir,
		TotalMatched: len(items),
		DurationMs:   time.Since(start).Milliseconds(),
		Items:        items,
	}, err
}
```

---

## 4. Acceptance Criteria

- [ ] `gitmap child-path` correctly inventories filesystem trees with normalized forward slashes.
- [ ] `--depth`, `--ext`, `--type`, `--limit`, and `--json` flags function as specified.
- [ ] Performance benchmark passes (<15ms scan time for typical repository structures).
- [ ] Excludes `.git`, `node_modules`, and `.cache` directories by default.
- [ ] Unit tests in `childpath_test.go` achieve >90% coverage with zero failures.
- [ ] Strict relative path hygiene enforced (zero absolute path leakage in output).

---

## 5. Verification Commands

```powershell
# 1. Run unit tests
go test -v ./cli/cmdchildpath/...

# 2. Test execution in table mode
gitmap child-path . "*.go" --depth 2 --type f

# 3. Test execution in JSON mode
gitmap child-path cli "*.go" --depth 1 --json --limit 5

# 4. Verify coding guidelines
python 03-ai-scripts/05-guideline-autofixer.py cli/cmdchildpath --check-only
```
