# Plan 83: Interactive Macro Builder PWD Header, In-Builder LS Listing & Helper Commands

**Status:** complete  
**Created:** 2026-09-09  
**Specification:** `spec/21-macro-builder-enhancements/`  
**Target:** `gitmap/cmd/macro_add_interactive.go`, `gitmap/cmd/macro_add_helpers.go`, `gitmap/cmd/macro_add_helpers_test.go`, `gitmap/uipref/`

---

## 1. Problem Statement & User Intent

During interactive macro creation (`gitmap macro add <name>`), users enter steps one by one into prompt `Step N> `.
1. **Missing PWD Context:** Users cannot see their current working directory (`PWD`), making it difficult to write commands operating on relative paths without switching terminals or aborting.
2. **In-Builder `ls` Inspection:** When a user types `ls` into the prompt, the builder simply registers `ls` as a step rather than listing directory files and folders to assist composition.
3. **Helper Commands (`find`, `search`, `replace`):** Users frequently need to find filenames, search file text, or perform replacements while composing commands for automated macros.

---

## 2. Architecture & Design

### A. Dynamic PWD Display Header
- Above each step prompt (`Step %d> `), render the PWD banner:
  ```text
    [PWD: /path/to/current/workdir]
    Step 1> 
  ```
- Toggleable via:
  - In-prompt commands: `pwd on` and `pwd off` (or `:pwd`).
  - Persistent preference stored in `gitmap/uipref/` (`IsMacroPwdVisible`).
  - Flag `--pwd` / `--no-pwd` on `gitmap macro add`.

### B. In-Builder File Listing (`ls` / `dir`)
- Recognize `ls`, `dir`, `:ls`, `:dir` in the interactive builder loop.
- Formats and displays directory files, subdirectories, file sizes, and item counts cleanly.
- After printing, prompts the user:
  `[ls executed for inspection. Press Enter to continue macro, or type '+add' to record 'ls' as step]`

### C. Helper Utilities (`find`, `search`, `replace`)
- Implement `gitmap/cmd/macro_add_helpers.go`:
  - `executeInteractiveFind(pattern string)`: scans current directory for matching files.
  - `executeInteractiveSearch(query string)`: greps matching lines across files.
  - `executeInteractiveReplace(oldStr, newStr, targetGlob string)`: safe token replacement in files.
- Integrated into the interactive loop with zero impact on standard macro step addition.

---

## 3. Subtasks Breakdown

1. [01-task-pwd-header-and-toggle.md](../subtasks/83-interactive-macro-builder-pwd-ls-search/01-task-pwd-header-and-toggle.md) — Implement PWD header rendering above `Step N> ` and toggle commands (`pwd on`/`pwd off`).
2. [02-task-interactive-ls-command.md](../subtasks/83-interactive-macro-builder-pwd-ls-search/02-task-interactive-ls-command.md) — Implement in-builder `ls`/`dir` file inspection without exiting builder session.
3. [03-task-search-find-replace-commands.md](../subtasks/83-interactive-macro-builder-pwd-ls-search/03-task-search-find-replace-commands.md) — Implement `find`, `search`, and `replace` in-builder helpers.
4. [04-task-unit-and-e2e-tests.md](../subtasks/83-interactive-macro-builder-pwd-ls-search/04-task-unit-and-e2e-tests.md) — Unit tests, coding guidelines validation, and CI/CD runner verification.

---

## 4. Verification & Testing

- `go test -v ./cmd/... -run TestInteractiveMacro`
- `python linter-scripts/check-nested-ifs.py`
- `python linter-scripts/check-boolean-guidelines.py`
- `python linter-scripts/check-relative-paths.py`
- `python linter-scripts/check-error-management.py`
- `python 03-ai-scripts/06-cicd-local-runner.py`
