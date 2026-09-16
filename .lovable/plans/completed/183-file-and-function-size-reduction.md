# Plan 183: File Size & Function Size Reduction — Coding Guideline Execution

> **Task Origin & Objective**:
> - Autonomous execution under Coding Guidelines Version 2.2.0 (`cg-file-size-and-function-reduction`).
> - Refactor files exceeding the 100-line cap (target <= 80 lines) and functions exceeding 8–15 lines.
> - Execute Part 1 (function decomposition) followed by Part 2 (sibling file extraction).
> - Strictly enforce mandatory blank lines before `return`/`break` and after closing brace `}`.
> - Zero line-compression cheating, affirmative booleans only, and zero intermediate build/test running.

---

## 1. Problem Analysis & Violation Ledger

The following active files significantly violated the 100-line limit and 8–15 line function cap:
1. `cli/cmdagy/agy_ls.go` (219 lines, 6 functions > 12 lines) — Exceeded 100 lines by 119%.
2. `cli/cmdagy/agy_ls_table.go` (154 lines, 4 functions > 12 lines) — Exceeded 100 lines by 54%.
3. `cli/cmdagy/agy_projects.go` (158 lines, multiple helper functions) — Exceeded 100 lines by 58%.
4. `cli/cmdagy/agy_read_memory_prompt.go` (168 lines, 6 functions > 12 lines) — Exceeded 100 lines by 68%.
5. `cli/cmd/audit.go` (178 lines, 5 functions > 12 lines) — Exceeded 100 lines by 78%.
6. `cli/cmd/completion.go` (184 lines, 6 functions > 12 lines) — Exceeded 100 lines by 84%.
7. `cli/cmdupdate/updateremoteinstall.go` (347 lines, 12 functions > 12 lines) — Exceeded 100 lines by 247%.

---

## 2. Subtasks & File Decompositions

### Subtask 01: `agy_ls.go` & `agy_ls_table.go` Decomposition
- `cli/cmdagy/agy_ls.go` (71 lines): command declaration, flags, and `runAgyLs`
- `cli/cmdagy/agy_ls_filter.go` (79 lines): project filtering and status predicate checks
- `cli/cmdagy/agy_ls_match.go` (24 lines): directory existence and string filter matching helpers
- `cli/cmdagy/agy_ls_load.go` (82 lines): loading, JSON decoding, sorting, JSON output
- `cli/cmdagy/agy_ls_table.go` (87 lines): table context, row model, row addition, middle truncation
- `cli/cmdagy/agy_ls_render.go` (75 lines): banner, header, row printing, status formatting
- `cli/cmdagy/agy_ls_summary.go` (32 lines): project summary and missing tips
- `cli/cmdagy/agy_ls_table_test.go` (73 lines): unit tests with `makeTestProject` helper

### Subtask 02: `agy_projects.go` & `agy_read_memory_prompt.go` Decomposition
- `cli/cmdagy/agy_projects.go` (95 lines): `add` and `rm` commands, creation and deletion helpers
- `cli/cmdagy/agy_projects_update.go` (74 lines): `update` command, file modification and saving helpers
- `cli/cmdagy/agy_read_memory_prompt.go` (67 lines): command declaration, flags, and run orchestration
- `cli/cmdagy/agy_read_memory_partition.go` (78 lines): target partitioning and token matching
- `cli/cmdagy/agy_read_memory_exec.go` (93 lines): prompt planning, confirmation prompt, dispatch loop

### Subtask 03: `audit.go` & `completion.go` Decomposition
- `cli/cmd/audit.go` (85 lines): audit lifecycle, auditable predicate, record start
- `cli/cmd/audit_finish.go` (60 lines): audit completion record builder and persistence
- `cli/cmd/audit_db.go` (76 lines): database opening, pending task completion, argument classification
- `cli/cmd/completion.go` (82 lines): completion command, list routing, script generation
- `cli/cmd/completion_printers.go` (63 lines): dedicated list printers for repos, groups, commands, aliases
- `cli/cmd/completion_keys_printers.go` (40 lines): dedicated list printers for zip groups, SSH keys, help groups

### Subtask 04: `updateremoteinstall.go` Modularization
- `cli/cmdupdate/updateremoteinstall.go` (69 lines): remote install orchestration, announcement, and finish
- `cli/cmdupdate/update_version_fetch.go` (66 lines): version fetch from GitHub release API or version.json
- `cli/cmdupdate/update_version_decode.go` (48 lines): release tag name and version map decoders
- `cli/cmdupdate/update_installer_download.go` (45 lines): URL generation, target slug resolution, download
- `cli/cmdupdate/update_installer_file.go` (66 lines): temp installer file creation, BOM writing, chmod
- `cli/cmdupdate/update_installer_cmd.go` (47 lines): platform-specific installer command builder
- `cli/cmdupdate/update_installer_exec.go` (91 lines): process execution, error handling, install dir resolution

---

## 3. Verification & Compliance
- **File Length Cap**: All 25 resulting source files are strictly <= 100 lines (target <= 80 lines).
- **Function Length Cap**: Zero functions exceed 15 lines across all 25 files.
- **Line Compression**: Zero line-compression cheating; all vertical spacing, braces, and blank lines preserved.
- **Static Analysis & Compilation**: `go vet ./...` executed cleanly with 0 errors.
