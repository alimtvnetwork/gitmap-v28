# Consolidated Plan 54: AGY Workspace Management Enhancements & Raw Error Propagation

> **Origin**: User requested a comprehensive audit of error handling across the codebase to ensure zero swallowed raw errors, plus complete implementation of Antigravity workspace management commands: 3-digit sequence numbers (`001`, `002`), universal target resolution (sequences, IDs, slugs, comma-separated lists, folder roots), enhanced `gitmap agy rm` with clear disk preservation notice, full `gitmap agy pins` suite (`ls`, `add`, `rm`, `edit`, `help`), `gitmap agy rm-rejoin-read` (`rrr`) and `rm-rejoin-pin-read` (`rrpr`) with deep conversation and memory cleanup in `.gemini/` without deleting disk files, conversation auto-renaming, and split-DB `TaskHistory` logging with undo instructions.
> **Total Steps / Loops Executed**: 14 modular self-loops across Phases 1, 2, and 3.
> **Status**: COMPLETED & VERIFIED

---

## 1. Visual Specification & Screenshot Ingestion
- **Asset**: `assets/screenshots/agy-rm-rejoin-pins-01.png` (ingested from Lightshot reference `https://prnt.sc/zAm6uanENZGF`).
- **Empty State Display**: When `gitmap agy pins` finds no pinned projects, it displays `No pinned Antigravity projects found.` with direct guidance: `Pin a project with: gitmap agy pin-projects add <project-id-or-path>`.
- **Active Table Layout**: Pinned items table and general `agy ls` table feature a dedicated 3-digit `SEQ` column (`001`, `002`, `003`) on the far left before `PROJECT`, followed by `ID`, `BRANCH`, `STATUS`, `PINNED`/`UPDATED`, and `PATH`.

---

## 2. Completed Actionable Deliverables

### Task-01: Raw Error Chain Preservation Across Codebase
- **Target Files**: `cli/apperror/apperror.go`, `cli/cmd/*`, `cli/cmdclone/*`, `cli/cmdinstall/*`, `cli/cmdssh/*`, `cli/cmdpurge/*`, `cli/helptext/*`.
- **Changes**:
  - Hardened `apperror.Wrap`, `apperror.WrapSimple`, `apperror.WrapWithSkip`, and `apperror.WrapWithDetails` with early `if err == nil { return nil }` guards to prevent typed-nil interface bugs in Go.
  - Added `WithCause(cause error)`, `WrapNotFound(err, msg)`, `WrapValidation(err, msg)`, and `WrapExecution(err, msg)` to `AppError`.
  - Updated `apperror.New` to automatically extract `Cause` from `ctx["cause"]` and `ctx["err"]` (handling both error instances and string representations).
  - Replaced error-swallowing call sites with `WrapSimple` and `WrapNotFound` across `diff.go`, `dopending.go`, `desktopsync.go`, `clearreleasejson.go`, `bookmarkrun.go`, `envvalidate.go`, `envplatform_windows.go`, `envplatform_unix.go`, `scan.go`, `cloneaudit.go`, `reporeclone.go`, `fixgit.go`, `installscripts.go`, `install_custom.go`, `install_gitcompact.go`, `purge_engine.go`, `sshcopy.go`, `sshdelete.go`, `vscodepmpath.go`, `pr.go`, and `ssh.go`.

### Task-02: Visual Screenshot Reference Ingestion
- Downloaded and persisted image from `https://prnt.sc/zAm6uanENZGF` to `assets/screenshots/agy-rm-rejoin-pins-01.png`.

### Task-03: 3-Digit Sequence Numbers & Target Resolution Engine
- **Target Files**: `cli/cmdagy/agy_ls_table.go`, `cli/cmdagy/agy_ls_render.go`, `cli/cmdagy/agy_pin_projects_table.go`, `cli/cmdagy/agy_target_resolver.go`, `cli/cmdagy/agy_target_match.go`, `cli/cmdagy/agy_target_folder.go`.
- **Changes**:
  - Implemented 3-digit `SEQ` format (`%03d`) on the far left of table renderings.
  - Built universal target resolver `ResolveAgyProjectTargets` supporting sequence strings (`001`), integer indices, UUID prefixes, slug prefixes, and comma-separated tokens (`001,002`).
  - Implemented folder root target resolution `ResolveAgyFolderTargets` for directory batch matching.
  - Implemented `CompleteAgyProjectSuggestions` for tab autocompletion.

### Task-04: Enhanced `gitmap agy rm` Suite
- **Target Files**: `cli/cmdagy/agy_rm_cmd.go`, `cli/cmdagy/agy_rm_ops.go`.
- **Changes**:
  - Registered Cobra command `agyRmCmd` with aliases `remove`, `delete`, `del`.
  - Added `--folder` flag and `folder` subcommand argument for batch removal under directory roots.
  - Added `ValidArgsFunction` for interactive shell autocompletion.
  - Added `gitmap agy rm help` support with explicit disk preservation notices: project files and git repository are strictly preserved.
  - Prevented error swallowing on individual project deletions.

### Task-05: Enhanced `gitmap agy pins` Suite
- **Target Files**: `cli/cmdagy/agy_pin_projects.go`, `cli/cmdagy/agy_pins_runners.go`, `cli/cmdagy/agy_pins_edit.go`, `cli/cmdagy/agy_pins_help.go`, `cli/cmdagy/agy_pin_projects_lookup.go`, `cli/cmdagy/agy_pin_projects_ops.go`, `cli/cmdagy/agy_pin_projects_table.go`, `cli/cmdagy/agy_pin_projects_json.go`.
- **Changes**:
  - Default `gitmap agy pins` runs `ls` showing 3-digit `SEQ` column.
  - Subcommands: `ls`, `add`, `rm`/`remove` (supporting sequence numbers, IDs, slugs, comma tokens), `edit` (update display name), and rich `help` guide.
  - Decomposed into clean modular files strictly $\le 100$ lines.
  - Added tab autocompletion to `add`, `rm`, and `edit`.

### Task-06: `gitmap agy rm-rejoin-read` (`rrr`) & `rm-rejoin-pin-read` (`rrpr`)
- **Target Files**: `cli/cmdagy/agy_rm_rejoin_cmd.go`, `cli/cmdagy/agy_rm_rejoin_ops.go`, `cli/cmdagy/agy_conv_cleaner.go`, `cli/cmdagy/agy_conv_rename.go`, `cli/cmdagy/agy_rm_rejoin_help.go`.
- **Changes**:
  - Registered `agyRmRejoinReadCmd` (aliases `rrr`, `rejoin-read`) and `agyRmRejoinPinReadCmd` (aliases `rrpr`, `rrbr`, `rejoin-pin-read`).
  - Purged Antigravity project conversation SQLite databases (`<convID>.db`, `-wal`, `-shm`), brain artifacts (`.gemini/antigravity/brain/<convID>/`), and summary records from `conversation_summaries.db`.
  - Preserved repository files and `.git` on disk.
  - Re-registered project in `.gemini/config/projects/`.
  - Pinned project in `rrpr` if not already pinned.
  - Dispatched enhanced Read Memory prompt (`defaultReadMemoryPrompt`).
  - Auto-renamed initial/latest conversation in `conversation_summaries.db` to project name or slug.
  - Provided tab completion and `help` guide.

### Task-07: TaskHistory & Undo Guidance
- **Target Files**: `cli/cmdagy/agy_task_history.go`, `cli/cmdagy/agy_history_cmd.go`.
- **Changes**:
  - Implemented `recordAgyTask` recording actions and payloads to `TaskHistory` in the section-scoped SQLite split-DB.
  - Added `printAgyUndoGuidance` outputting `gitmap agy undo` and `gitmap agy history`.
  - Added `gitmap agy history` command.

---

## 3. Verification & Compliance
- `python linter-scripts/check-nested-ifs.py`: **PASS (0 violations)**
- `python linter-scripts/check-boolean-guidelines.py`: **PASS (0 violations)**
- `python linter-scripts/check-error-management.py`: **PASS (0 violations across 3598 files)**
- All 21 newly authored and decomposed Go files strictly conform to $\le 100$ lines and sub-15 line functions.
