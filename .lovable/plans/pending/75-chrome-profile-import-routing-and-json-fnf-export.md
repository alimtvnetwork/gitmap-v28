# Plan 75: Chrome Profile Import Routing, JSON & FNF Export, and End-to-End Test Suite

## Overview
Fix `gitmap profile import` execution failure (`cmd/profile.go:62 [E9000:EXECUTION]`) by routing profile import, export, and inspection subcommands directly to the Chrome smart profile engine. Furthermore, add structured JSON export and file system writing capabilities (`--json`, `--file <path>`, `--fnf <path>`, `--tempfile <filename>`) to candidate inspection and import previews, and update all corresponding help documentation and end-to-end tests.

---

## 1. Problem Statements & Root Causes

### Problem 1: `gitmap profile import` Fatal Error (`cmd/profile.go:62`)
- **Symptom:**
  ```text
  a@a:~/Desktop/test/chrome-ext$ gitmap profile import
  usage: gitmap profile <create|list|switch|delete|show> [name]
  gitmap: [E9000:EXECUTION] fatal error
    origin: cmd/profile.go:62
  ```
- **Root Cause:**
  `routeProfileSub(sub, args)` in `gitmap/cmd/profile.go` only handled Git/database profile operations (`create`, `list`, `switch`, `delete`, `show`, `git`, `accounts`, `set-default`). When a user typed `gitmap profile import` or `gitmap profile import-all` or `gitmap profile export`, it failed with `E9000` instead of routing to Chrome profile management.
- **Remediation:**
  Decompose `routeProfileSub` into modular sub-routers conforming to the 15-line function limit (`CG-SIZE-002`):
  - `routeGitProfileSub`: `git`, `accounts`, `set-default`
  - `routeDBProfileSub`: `create`, `list`, `switch`, `delete`, `show`
  - `routeChromeProfileSub`:
    - `import`, `import-profile`: `runChromeProfileImport(args)`
    - `import-all`: `runChromeImportAll(args)`
    - `export`, `export-profile`: `runChromeProfileExport(args)`
    - `export-all`: `runChromeExportAll(args)`
    - `inspect`, `preview`, `check`, `import-check`, `ls`: `runChromeProfileImportCheck(args)`
  - Update `constants.ErrProfileUsage` and mark new constants with `// gitmap:cmd skip` for AST parity.

### Problem 2: JSON & FNF Preflight Inspection Export
- **User Request:**
  "I really like the output. That's really nice. Also add options to export into JSON, uh, with FNF and JSON flag, or probably write it to the file system as well. So give these two examples in the help as well."
- **Root Cause:**
  `runChromeProfileImportCheck` only dumped table output or printed JSON to `os.Stdout`. There was no support for saving output to disk via `--file <path>`, `--fnf <path>`, or `--tempfile <filename>`. In addition, `resolveCheckTarget` mistakenly consumed value flags as target directory arguments.
- **Remediation:**
  1. Add flag parsing for `--json`, `--file <path>`, `--fnf <path>`, and `--tempfile <filename>`.
  2. Implement `dispatchPreviewOutput` in `gitmap/cmd/chromeprofile_smart_import_export.go` to cleanly marshal candidates to JSON or formatted text, create parent directories, write to disk, and print clear feedback.
  3. Add usage hints in `renderProfileCandidatesTable` for exporting to JSON and writing with `--file` / `--fnf`.
  4. Update help documentation in `helptext/import-check.md`, `helptext/import-all.md`, and `helptext/profile.md` with explicit examples.

### Problem 3: End-to-End Testing & Verification
- **Requirements:**
  Add E2E tests verifying:
  - `gitmap profile import` routing without crash.
  - `gitmap profile import-all` routing without crash.
  - `gitmap profile check` / `inspect` / `preview` routing.
  - `gitmap chrome import-check --json` producing valid JSON.
  - `gitmap chrome import-check --json --file <path>` writing JSON file to disk.
  - `gitmap chrome import-check --json --fnf <path>` writing JSON file to disk.
  - `gitmap chrome import-check --tempfile <filename>` writing to `.lovable/temp/<filename>`.

---

## 2. Task-Specific Rule Set
1. **Rule 1 (Function Sizing):** All functions strictly $\le 15$ lines. Blank line before every return.
2. **Rule 2 (Zero Swallow):** Return all errors up the call stack; router returns `*AppError` rather than calling `cliexit.HandleError` internally.
3. **Rule 3 (AST Parity):** All new `CmdProfile*` constants must be annotated with `// gitmap:cmd skip`.
4. **Rule 4 (Helptext Syntax):** All code examples in markdown help files must use fenced code blocks (` ``` `), never 4-space indentation.
5. **Rule 5 (Golden Tests):** `go test ./gitmap/helptext/... -run Golden -count=1` must pass cleanly.
