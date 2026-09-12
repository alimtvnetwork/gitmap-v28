# Rewrite: Interactive Macro Builder PWD Header, In-Builder LS, and Helper Commands

**Title:** Interactive Macro Builder: Dynamic PWD Header, In-Builder File Listing, and Helper Commands (replace, find, search)
**Task ID:** Plan 83
**Target Files:**
- `gitmap/cmd/macro_add_interactive.go`
- `gitmap/cmd/macro_add_helpers.go` [NEW]
- `gitmap/cmd/macro_add_helpers_test.go` [NEW]
- `gitmap/cmd/macro_add_test.go`
- `gitmap/constants/constants.go`
- `gitmap/uipref/uipref.go`

---

## 1. Context & Motivation

When users launch `gitmap macro add <name>` in an interactive terminal session:
1. Currently, the prompt displays `Step N> ` without contextual directory information. Users do not see their current working directory (PWD), making it error-prone to reference relative paths or files.
2. If the user enters `ls`, the interactive builder currently captures `ls` as a macro step rather than allowing the user to inspect directory contents and files interactively during macro composition.
3. Users need built-in in-builder commands to assist macro building:
   - `ls` / `dir`: List files and subdirectories with sizes, file count, and directories.
   - `pwd`: Print full current working directory, or toggle PWD header banner.
   - `pwd on` / `pwd off`: Dynamically enable or disable the PWD header display above the prompt line.
   - `find <pattern>`: Search for file names matching glob or substring.
   - `search <pattern>`: Search file contents for patterns or occurrences.
   - `replace <old> <new> [files...]`: Replace text occurrences across target files.

---

## 2. Requirements & Acceptance Criteria

### A. PWD Display Above Prompt Line (Configurable & Toggleable)
- [ ] Display current working directory above each `Step N> ` prompt:
  ```text
  [PWD: /path/to/current/dir]
  Step 1>
  ```
- [ ] Support toggling PWD visibility on and off:
  - Command `pwd on` or `pwd off` directly inside the interactive prompt toggles display.
  - Persistent preference in `uipref` or config (`macro.show_pwd = true/false`).
  - Flag `--pwd` / `--no-pwd` on `gitmap macro add`.

### B. In-Builder Interactive Commands (Commands that inspect/assist rather than append as macro steps)
- [ ] `ls` / `dir`:
  - When entered in interactive builder as a helper command (or prefixed with `:` or executed directly):
  - Displays formatted directory listing: subdirectories, files, sizes, item counts.
  - Asks user or provides option whether to add `ls` to the macro or simply inspect files.
- [ ] `find <pattern>`:
  - Finds and lists files/directories matching `<pattern>` in current working directory tree.
- [ ] `search <query>`:
  - Searches file contents for `<query>` and displays matching files and lines.
- [ ] `replace <old> <new> [glob]`:
  - Performs safe string replacement in files matching `[glob]` or target file.

### C. Parity, Guidelines & Testing
- [ ] Adhere strictly to Go guidelines: functions $\le 15$ lines, zero nested `if` blocks, affirmative booleans, blank line before every return.
- [ ] Zero error swallowing: wrap errors using `apperror`.
- [ ] Comprehensive unit tests in `gitmap/cmd/macro_add_helpers_test.go` covering PWD toggle, `ls` listing, `find`, `search`, and `replace`.
- [ ] Full CI/CD quality gate pass (`python 03-ai-scripts/06-cicd-local-runner.py`).

---

## 3. Mandatory References

- `.lovable/plan.md`: Failure recovery record and architectural state.
- `mem://01-index.md`: Core memory index and CODE RED guidelines.
- `.lovable/coding-guidelines.md`: Function length, nesting, and boolean naming requirements.
- `gitmap/cmd/macro_add_interactive.go`: Current interactive prompt implementation.
- `gitmap/cmd/macro_types.go`: Macro types and step data structures.

---

## 4. Definition of Done

1. `gitmap macro add <name>` interactively shows `[PWD: ...]` above the step line when enabled.
2. Typing `pwd off` disables the PWD header; typing `pwd on` re-enables it.
3. Typing `ls` previews directory files and contents without aborting the macro session.
4. Typing `find`, `search`, and `replace` executes helper utilities within the session.
5. All unit tests pass 100% green.
6. Linters pass with 0 violations.

---

## 5. Self-Instruction for AI Agent

Before acting, re-read `mem://01-index.md` and `.lovable/coding-guidelines.md`; restate which rules apply. Suggest further improvements to this instruction after execution.
