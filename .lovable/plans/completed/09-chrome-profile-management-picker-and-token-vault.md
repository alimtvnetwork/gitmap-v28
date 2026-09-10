# Milestone Summary: Chrome Profile Management, Picker & Token Vault

## 1. Executive Overview & Consolidated Tasks

- **Milestone Domain:** Chromium Local State Schema, Graphical Profile Picker Recognition, JSON/FNF Preflight & Vault
- **Total Original Plans Merged:** 3 plans
  - `63-chrome-profile-picker-visibility.md`
  - `73-chrome-profile-import-export-ubuntu-install-and-error-trace.md`
  - `75-chrome-profile-import-routing-and-json-fnf-export.md`
- **Associated Subtask Folders Folded:** 2 folders
  - `73-chrome-profile-import-export-ubuntu-install-and-error-trace`
  - `75-chrome-profile-import-routing-and-json-fnf-export`
- **Status:** `COMPLETED`
- **Core Architecture & Invariants:** Populate all 13 Chromium UI attributes in Local State. Sanitize Preferences on import. Preflight inspection supports --json, --file, --fnf, and --tempfile. Reversible 2-pass Base64 + Caesar cipher vault.

## 2. Key Architectural Decisions & Spec Implementations

- **Authoritative Specifications Implemented:**
  - spec/25-chrome-profile-management/01-profile-registration-and-picker.md — 13-attribute UI schema and ordering.
  - spec/25-chrome-profile-management/02-import-export-and-vault.md — Multi-profile discovery, ZIP extraction, and token cipher.
- **Core Architecture Contracts:**
  - Populate all 13 Chromium UI attributes in Local State. Sanitize Preferences on import. Preflight inspection supports --json, --file, --fnf, and --tempfile. Reversible 2-pass Base64 + Caesar cipher vault.

## 3. Deep Dive into Consolidated Plans & Subtask Chronicles

Every individual plan and subtask merged into this milestone is preserved below in full technical detail, ensuring 100% fidelity, zero truncation, and complete traceability.

### Merged Plan: `63-chrome-profile-picker-visibility.md`

#### 63 — Chrome Profile Picker Recognition & Local State Registration

**Status:** COMPLETED
**Priority:** High (CODE RED)
**Spec Reference:** [`spec/25-chrome-profile-management/01-profile-registration-and-picker.md`](../../../spec/25-chrome-profile-management/01-profile-registration-and-picker.md)
**Related Issue:** [`.lovable/memory/issues/2026-09-05-chrome-profile-picker-visibility-desync.md`](../../memory/issues/2026-09-05-chrome-profile-picker-visibility-desync.md)

---

##### 1. Visual Problem Evidence

###### Figure 1: Import CLI Execution (`Profile 5` Allocated)
![Import Terminal Output](../../assets/screenshots/19-chrome-profile-import-terminal.png)
*CLI logs confirming that `gitmap chrome profile import "erfan.office.n@gmail.com"` allocated `Profile 5` on disk.*

###### Figure 2: Chrome Profile Picker UI (`Profile 5` Missing)
![Chrome Profile Picker UI](../../assets/screenshots/20-chrome-profile-picker-missing-profile5.png)
*Chrome Profile Picker showing existing profiles: "Main Account", "Alim Mum / Md Karim", "Riseup Veo 1", "Alim Personal Veo 1 / Business P v1", "Riseup Team v4", and "Lov". `Profile 5` is missing!*

---

##### 2. Root Cause Analysis Summary

1. **Active Chrome Process In-Memory Overwrite:**
   - Chrome was actively running during import (`chrome.exe` processes active).
   - Chrome maintains `Local State` in-memory and flushes periodically or on exit, completely wiping external disk modifications.
2. **Incomplete Chromium Profile Picker Schema:**
   - `registerImportedProfileToLocalState` only saved basic fields (`name`, `user_name`), omitting avatar styling, color codes, `active_time`, `shortcut_name`, and default avatar flags required by Chromium's UI parser.
3. **Preferences & Local State Desynchronization:**
   - Profile's internal `Preferences` contained stale metadata (`Person 1`, old custodian records) rather than the allocated profile name and sanitized identity.

---

##### 3. Implementation Plan & Technical Invariants

###### Step 1: Upgrade Chrome Local State Registration Engine
- In `gitmap/cmd/chromeprofile_register.go` and `gitmap/cmd/chromeprofile_smart_import.go`:
  - Consolidate profile registration into a unified, hardened function `registerChromeProfileWithFullSchema(dstDir, displayName, email string) error`.
  - Populate all 13 Chromium UI attributes:
    - `name`: Display name
    - `shortcut_name`: Display name
    - `user_name`: Email address
    - `active_time`: Float64 timestamp
    - `avatar_icon`: `"chrome://theme/IDR_PROFILE_AVATAR_26"`
    - `default_avatar_fill_color`: `-13625057`
    - `default_avatar_stroke_color`: `-1786428`
    - `profile_highlight_color`: `-13625057`
    - `is_using_default_avatar`: `true`
    - `is_using_default_name`: `false`
    - `is_ephemeral`: `false`
    - `is_consented_primary_account`: `false`
    - `signin.with_credential_provider`: `false`
  - Ensure `profiles_order` contains the directory without duplicate entries.

###### Step 2: Preferences Sanitization on Import
- Before registration, sanitize `<ProfileDir>/Preferences`:
  - Set `profile.name` = `displayName`.
  - Set `profile.using_default_name` = `false`.
  - Remove stale `account_info`, `signin`, `google`, `gaia_cookie`, `custodian_*`.

###### Step 3: Process Concurrency Guard & User Advisory
- Check `isChromeRunning(runtime.GOOS)`.
- If running:
  - Output high-visibility terminal alert informing user that Chrome must be closed or restarted for profile picker recognition.
- If closed:
  - Confirm instant profile picker readiness.

###### Step 4: Implement Orphan Reconciliation Command (`gitmap chrome profile reconcile`)
- Add `gitmap chrome profile reconcile` (alias `gitmap chrome profile repair`).
- Scans `%LOCALAPPDATA%\Google\Chrome\User Data` for any profile folder (like `Profile 5`) not listed in `Local State`.
- Automatically synthesizes missing metadata and registers the orphaned profiles.

###### Step 5: Immediate Remediation of Existing `Profile 5`
- Register `Profile 5` (`erfan.office.n@gmail.com`) immediately so it appears in the user's Chrome Profile Picker.

---

##### 4. Acceptance Criteria & Validation

1. `go test -v ./cmd/... -run TestChromeProfile` passes 100%.
2. Zero nested if statements repository-wide (`python linter-scripts/check-nested-ifs.py`).
3. Zero error management violations (`python linter-scripts/check-error-management.py`).
4. Full CI/CD local runner suite green (`python 03-ai-scripts/06-cicd-local-runner.py`).

### Merged Plan: `73-chrome-profile-import-export-ubuntu-install-and-error-trace.md`

#### 73-chrome-profile-import-export-ubuntu-install-and-error-trace.md: Chrome Profile Multi-Directory/ZIP Discovery, Ubuntu Chrome Installer & Rich Error Diagnostics

##### 1. Executive Summary

This plan addresses critical runtime issues discovered during Chrome profile imports, tool installations on Ubuntu Linux, and diagnostic error trace visibility:
1. **Chrome Profile Import & Multi-Profile Discovery (`gitmap chrome import-all`)**:
   - Resolve single-dummy profile import bug where directory containing multiple profiles (e.g. `chrome-ext/` with `manifest.json` + `Profile 1/`, `Profile 2/`, `Default/`) only imports `manifest.json` as a profile named `"manifest"`.
   - Prevent `manifest.json` (metadata index) from being parsed or treated as a profile payload.
   - Support deep recursive subdirectory profile discovery (`Default/`, `Profile */`, `Guest Profile/`, custom directories) containing `<profile>.json`, `Preferences`, `Bookmarks`, or `Extensions/`.
   - Support multi-profile ZIP archive extraction & discovery: when importing a `.zip` (e.g. `chrome-ext.zip`), unpack into a safe sandbox, inspect `manifest.json` or profile folders, discover all profiles, and batch import them into Chrome User Data.
   - Preflight inspection commands: `gitmap chrome import ls [path]`, `gitmap chrome import-all ls [path]`, `gitmap chrome import-check [path]` displaying tabular profile previews (name, display name, email, extensions count, bookmarks count, target directory) with `--json` and `--verbose` support.
2. **Chrome CLI Helptext & Alias Parity**:
   - Provide dedicated markdown help files for `import-all`, `export-all`, `copy-all`, `import-check`.
   - Add alias routing in `helptext/print.go` to prevent `E1071: No help available for 'import-all'` errors.
   - Ensure all markdown help files contain `## Examples` with fenced code blocks satisfying AST tests.
3. **Ubuntu Google Chrome Installation (`gitmap chrome install` / `gitmap install chrome`)**:
   - Fix APT failure `Unable to locate package google-chrome-stable` (exit status 100) on Ubuntu/Debian.
   - Download official `.deb` (`https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb`) to `/tmp/gitmap-stage-*`, install via `sudo apt install -y <deb>`, and fall back to Google keyring and APT repository configuration.
   - Register `ToolChrome` and `ToolGoogleChrome` in `constants/constants_install.go`.
4. **Enhanced Error Diagnostics & Operational Stack Traces**:
   - Fix bug in `cliexit/handle.go`: line 73 `if e.Cause != nil && e.Message != ""` silences `e.Cause` when `e.Message == ""` .
   - Enhance error output formatting with origin caller (`e.Caller`), command context, exit code, log path, and stderr snippet.
   - Print diagnostic operational stack traces for fatal execution errors (`E9000:EXECUTION`) or when in debug/verbose mode.
   - Fix `cmd/installtools.go` (`handleInstallError`) to wrap rich `AppError` instead of empty `fatal error`.
5. **VMware Mount & Commands Verification**:
   - Wire `checkHelp("vmware", args)` into `cmd/vmware.go` and `cmd/vmware_shared.go`.
   - Register `HelpVmware` in `constants/constants_cli.go`, `cmd/rootusage_groups.go`, and `helptext/catalog.go`.
   - Support `--dry-run` and auto-probe for `open-vm-tools` / `vmhgfs-fuse` in `runVmwareSharedEnable`.
   - Display verified commands (`gitmap vmware shared enable`, `status`, `gitmap os fix-link`) and render full help text.

##### 2. Task-Specific Rules & Constraints

1. **Rule 1 (Strict Relative Paths):** Zero `file:///` URIs or machine-specific absolute paths in markdown files, artifacts, or code. All paths must be relative to repository root (`/`).
2. **Rule 2 (Coding Guidelines Compliance):** Functions <= 15 lines (prefer <= 8), blank line before every return statement, affirmative booleans (`is_*`, `has_*`), no magic numbers, and zero swallowed errors (`err != nil` must be handled or wrapped via `apperror`).
3. **Rule 3 (Strict 200-Line File Limit):** All new and modified Go source files must remain <= 200 lines. Decompose logic into focused single-responsibility files.
4. **Rule 4 (AST Parity & Constants):** All CLI verbs and help entries must be synchronized with `gitmap/constants/constants_cli.go` and verified with `TestTopLevelCmdRegistryMatchesAST` and `TestEveryHelpFileHasExamples`.
5. **Rule 5 (CI/CD Local Runner Validation):** Must run `python 03-ai-scripts/06-cicd-local-runner.py` exit 0 before release orchestration.

##### 3. Subtasks Breakdown

- [01-chrome-profile-polymorphic-discovery-and-preflight-ls.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/01-chrome-profile-polymorphic-discovery-and-preflight-ls.md)
- [02-chrome-profile-zip-multi-profile-extraction.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/02-chrome-profile-zip-multi-profile-extraction.md)
- [03-chrome-helptext-and-alias-routing.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/03-chrome-helptext-and-alias-routing.md)
- [04-ubuntu-chrome-installer-and-special-handler.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/04-ubuntu-chrome-installer-and-special-handler.md)
- [05-enhanced-error-reporting-and-stack-traces.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/05-enhanced-error-reporting-and-stack-traces.md)
- [06-vmware-help-and-dryrun-integration.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/06-vmware-help-and-dryrun-integration.md)
- [07-quality-gates-tests-and-release.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/07-quality-gates-tests-and-release.md)

#### Granular Subtask Execution Details for `73-chrome-profile-import-export-ubuntu-install-and-error-trace`

##### Subtasks Folder: `73-chrome-profile-import-export-ubuntu-install-and-error-trace` (7 subtask files incorporated)
###### Subtask File: `01-chrome-profile-polymorphic-discovery-and-preflight-ls.md`

#### Subtask 01: Chrome Profile Polymorphic Discovery & Preflight `ls`

##### Scope
- Update `cmd/chromeprofile_smart_import.go` and add `cmd/chromeprofile_smart_import_discovery.go`:
  - Discard/skip `manifest.json` from being treated as a single profile.
  - Implement recursive profile discovery across subdirectories (`Default/`, `Profile 1/`, `Profile 2/`, etc.).
  - Recognize profile folders containing `<folder>.json`, `Preferences`, `Bookmarks`, or `Extensions/`.
  - Add preflight inspection handlers for:
    - `gitmap chrome import ls [path]`
    - `gitmap chrome import-all ls [path]`
    - `gitmap chrome import-check [path]`
  - Render clear tabular inspection view (Profile Directory Name, Display Name, Email, Bookmarks Count, Extensions Count, Source Path) with `--json` support.

##### Files Touched
- `gitmap/cmd/chromeprofile_smart_import.go`
- `gitmap/cmd/chromeprofile_smart_import_discovery.go`
- `gitmap/cmd/chromeprofile_import_handlers.go`
- `gitmap/cmd/chrome.go`

###### Subtask File: `02-chrome-profile-zip-multi-profile-extraction.md`

#### Subtask 02: Chrome Profile ZIP Multi-Profile Extraction & Batch Import

##### Scope
- Update `cmd/chromeprofile_zip_import.go`:
  - Prevent ZIP archives from being treated as a single generic profile named after the archive file.
  - Safely extract `.zip` archive into temporary staging sandbox (`/tmp/gitmap-stage-chrome-*` or `os.TempDir()`).
  - Read `manifest.json` if present in root or probe top-level directories for profiles.
  - For each discovered profile in the archive:
    - Extract preferences, bookmarks, extensions.
    - Dispatch to `importChromeJSON` / `applyChromeExport`.
  - Clean up sandbox upon completion.

##### Files Touched
- `gitmap/cmd/chromeprofile_zip_import.go`
- `gitmap/cmd/chromeprofile_smart_import.go`

###### Subtask File: `03-chrome-helptext-and-alias-routing.md`

#### Subtask 03: Chrome CLI Helptext & Alias Routing Parity

##### Scope
- Author missing help markdown files in `gitmap/helptext/`:
  - `import-all.md`: Detailed help and examples for `gitmap chrome import-all`.
  - `export-all.md`: Detailed help and examples for `gitmap chrome export-all`.
  - `copy-all.md`: Detailed help and examples for `gitmap chrome copy-all`.
  - `import-check.md`: Detailed help and examples for `gitmap chrome import-check`.
- Update `gitmap/helptext/print.go`:
  - Add alias resolution mapping for `import-all`, `export-all`, `copy-all`, `import-check`, and `import-ls`.
  - Ensure all markdown files contain `## Examples` with fenced code blocks satisfying AST tests.

##### Files Touched
- `gitmap/helptext/import-all.md`
- `gitmap/helptext/export-all.md`
- `gitmap/helptext/copy-all.md`
- `gitmap/helptext/import-check.md`
- `gitmap/helptext/print.go`

###### Subtask File: `04-ubuntu-chrome-installer-and-special-handler.md`

#### Subtask 04: Ubuntu Google Chrome Installer & Special Handler

##### Scope
- Create `gitmap/cmd/install_chrome_linux.go`:
  - Download official `.deb` (`https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb`) into staging `/tmp`.
  - Execute `sudo apt install -y /tmp/google-chrome-stable_current_amd64.deb` to handle dependencies and configure repository automatically.
  - Fallback: configure Google signing key (`/etc/apt/keyrings/google-chrome.gpg`) and sources list (`/etc/apt/sources.list.d/google-chrome.list`), then `sudo apt update && sudo apt install -y google-chrome-stable`.
- Update `gitmap/cmd/install.go`:
  - Wire `chrome` and `google-chrome` to `specialInstallHandler` invoking `runInstallChromeLinux` on Linux.
- Update `gitmap/cmd/chrome.go`:
  - In `runChromeInstall`, call `runInstallChromeLinux` when on Linux instead of generic package manager.
- Update `gitmap/constants/constants_install.go`:
  - Register `ToolChrome` and `ToolGoogleChrome` in `InstallToolDescriptions` and `InstallToolCategories`.

##### Files Touched
- `gitmap/cmd/install_chrome_linux.go`
- `gitmap/cmd/install.go`
- `gitmap/cmd/chrome.go`
- `gitmap/constants/constants_install.go`

###### Subtask File: `05-enhanced-error-reporting-and-stack-traces.md`

#### Subtask 05: Enhanced Error Diagnostics & Operational Stack Traces

##### Scope
- Update `gitmap/cliexit/handle.go`:
  - Fix bug on line 73: change `if e.Cause != nil && e.Message != ""` to `if e.Cause != nil`.
  - Render origin caller (`e.Caller`), command context, exit code, log path, and stderr snippet.
  - Print diagnostic operational stack traces for fatal execution errors (`E9000:EXECUTION`) or when running with `--debug` / `GITMAP_DEBUG=1`.
- Update `gitmap/cmd/installtools.go`:
  - Refactor `handleInstallError` to wrap a rich `AppError` with exit code, command, and stderr rather than returning bare `apperror.NewSimple("fatal error", "E9000")`.

##### Files Touched
- `gitmap/cliexit/handle.go`
- `gitmap/cmd/installtools.go`

###### Subtask File: `06-vmware-help-and-dryrun-integration.md`

#### Subtask 06: VMware Help & Dry-Run Integration

##### Scope
- Update `gitmap/cmd/vmware.go` and `gitmap/cmd/vmware_shared.go`:
  - Wire `checkHelp("vmware", args)` so `gitmap vmware help` and `gitmap vmware shared help` show structured help.
  - Support `--dry-run` flag in `runVmwareSharedEnable`.
  - Auto-probe `open-vm-tools` and `vmhgfs-fuse` with prompt/auto-install if missing.
- Update `gitmap/constants/constants_cli.go`, `gitmap/cmd/rootusage_groups.go`, and `gitmap/helptext/catalog.go`:
  - Register `HelpVmware` in constants and catalog.
- Verify commands (`gitmap vmware shared enable`, `status`, `gitmap os fix-link`) and render full help text.

##### Files Touched
- `gitmap/cmd/vmware.go`
- `gitmap/cmd/vmware_shared.go`
- `gitmap/constants/constants_cli.go`
- `gitmap/cmd/rootusage_groups.go`
- `gitmap/helptext/catalog.go`

###### Subtask File: `07-quality-gates-tests-and-release.md`

#### Subtask 07: Quality Gates, Linters, CI Verification & Release Orchestration

##### Scope
- Run Go unit tests and AST parity verification:
  - `go test ./...`
  - `go test ./constants/... -run TestTopLevelCmdRegistryMatchesAST`
  - `go test ./helptext/... -run TestEveryHelpFileHasExamples`
- Run coding guideline linters:
  - `python 03-ai-scripts/check-nested-ifs.py`
  - `python 03-ai-scripts/check-enum-and-boolean.py`
- Run local CI runner:
  - `python 03-ai-scripts/06-cicd-local-runner.py` exit 0.
- Execute release orchestrator:
  - `python 03-ai-scripts/29-release-orchestrator.py --tier minor`

##### Files Touched
- All touched files
- `changelog.md`
- `version.json`


### Merged Plan: `75-chrome-profile-import-routing-and-json-fnf-export.md`

#### Plan 75: Chrome Profile Import Routing, JSON & FNF Export, and End-to-End Test Suite

##### Overview
Fix `gitmap profile import` execution failure (`cmd/profile.go:62 [E9000:EXECUTION]`) by routing profile import, export, and inspection subcommands directly to the Chrome smart profile engine. Furthermore, add structured JSON export and file system writing capabilities (`--json`, `--file <path>`, `--fnf <path>`, `--tempfile <filename>`) to candidate inspection and import previews, and update all corresponding help documentation and end-to-end tests.

---

##### 1. Problem Statements & Root Causes

###### Problem 1: `gitmap profile import` Fatal Error (`cmd/profile.go:62`)
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

###### Problem 2: JSON & FNF Preflight Inspection Export
- **User Request:**
  "I really like the output. That's really nice. Also add options to export into JSON, uh, with FNF and JSON flag, or probably write it to the file system as well. So give these two examples in the help as well."
- **Root Cause:**
  `runChromeProfileImportCheck` only dumped table output or printed JSON to `os.Stdout`. There was no support for saving output to disk via `--file <path>`, `--fnf <path>`, or `--tempfile <filename>`. In addition, `resolveCheckTarget` mistakenly consumed value flags as target directory arguments.
- **Remediation:**
  1. Add flag parsing for `--json`, `--file <path>`, `--fnf <path>`, and `--tempfile <filename>`.
  2. Implement `dispatchPreviewOutput` in `gitmap/cmd/chromeprofile_smart_import_export.go` to cleanly marshal candidates to JSON or formatted text, create parent directories, write to disk, and print clear feedback.
  3. Add usage hints in `renderProfileCandidatesTable` for exporting to JSON and writing with `--file` / `--fnf`.
  4. Update help documentation in `helptext/import-check.md`, `helptext/import-all.md`, and `helptext/profile.md` with explicit examples.

###### Problem 3: End-to-End Testing & Verification
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

##### 2. Task-Specific Rule Set
1. **Rule 1 (Function Sizing):** All functions strictly $\le 15$ lines. Blank line before every return.
2. **Rule 2 (Zero Swallow):** Return all errors up the call stack; router returns `*AppError` rather than calling `cliexit.HandleError` internally.
3. **Rule 3 (AST Parity):** All new `CmdProfile*` constants must be annotated with `// gitmap:cmd skip`.
4. **Rule 4 (Helptext Syntax):** All code examples in markdown help files must use fenced code blocks (` ``` `), never 4-space indentation.
5. **Rule 5 (Golden Tests):** `go test ./gitmap/helptext/... -run Golden -count=1` must pass cleanly.

#### Granular Subtask Execution Details for `75-chrome-profile-import-routing-and-json-fnf-export`

##### Subtasks Folder: `75-chrome-profile-import-routing-and-json-fnf-export` (4 subtask files incorporated)
###### Subtask File: `01-task.md`

#### Subtask 01: Profile Routing & Constants Update

##### Objective
Update `gitmap/constants/constants_profile.go` and `gitmap/cmd/profile.go` to route `gitmap profile` commands for `import`, `import-all`, `export`, `export-all`, `inspect`, `preview`, and `check` directly to Chrome profile handlers.

##### Affected Files
- `gitmap/constants/constants_profile.go`
- `gitmap/cmd/profile.go`

##### Requirements
1. In `gitmap/constants/constants_profile.go`:
   - Add constants with `// gitmap:cmd skip`:
     - `CmdProfileImport = "import"`
     - `CmdProfileImportAll = "import-all"`
     - `CmdProfileExport = "export"`
     - `CmdProfileExportAll = "export-all"`
     - `CmdProfileInspect = "inspect"`
     - `CmdProfilePreview = "preview"`
     - `CmdProfileCheck = "check"`
     - `CmdProfileImportCheck = "import-check"`
   - Update `HelpProfile` and `ErrProfileUsage`.
2. In `gitmap/cmd/profile.go`:
   - Refactor `routeProfileSub` to delegate to modular sub-routers (`routeGitProfileSub`, `routeDBProfileSub`, `routeChromeProfileSub`).
   - Obey $\le 15$ lines per function limit.
   - Return `*AppError` on unknown subcommands instead of calling `cliexit.HandleError`.
   - Propagate errors without swallowing.

###### Subtask File: `02-task.md`

#### Subtask 02: JSON, File, and FNF Export Options for Inspection Preview

##### Objective
Implement `--json`, `--file <path>`, `--fnf <path>`, and `--tempfile <filename>` export options for Chrome profile inspection preflight checks.

##### Affected Files
- `gitmap/cmd/chromeprofile_transfer_opts.go`
- `gitmap/cmd/chromeprofile_smart_import_export.go` (new)
- `gitmap/cmd/chromeprofile_smart_import.go`
- `gitmap/cmd/chromeprofile_smart_import_preview.go`

##### Requirements
1. In `gitmap/cmd/chromeprofile_transfer_opts.go`:
   - Extend `chromeTransferOptions` with `IsJSON`, `FilePath`, `TempFile`, `Fnf`.
   - Update `parseChromeTransferOptions` to safely parse `--file`, `-f`, `-o`, `--fnf`, `--tempfile`, `--json`.
2. In `gitmap/cmd/chromeprofile_smart_import_export.go`:
   - Create `dispatchPreviewOutput` to handle console table, raw JSON, or writing to disk (`--file`, `--fnf`, `--tempfile`).
   - Create parent directories automatically and print confirmation on write.
3. In `gitmap/cmd/chromeprofile_smart_import.go`:
   - Update `resolveCheckTarget` to skip valued flags so `--file <path>` does not pollute the target path.
   - Wire `runChromeProfileImportCheck` to use `dispatchPreviewOutput`.
4. In `gitmap/cmd/chromeprofile_smart_import_preview.go`:
   - Update `renderProfileCandidatesTable` usage hints to include `--json`, `--file`, and `--fnf` examples.

###### Subtask File: `03-task.md`

#### Subtask 03: Helptext Documentation & Golden Tests

##### Objective
Update helptext markdown documentation across relevant profile and import commands, adding examples for JSON export, `--file`, and `--fnf`, and ensuring all helptext golden tests pass.

##### Affected Files
- `gitmap/helptext/import-check.md`
- `gitmap/helptext/import-all.md`
- `gitmap/helptext/profile.md`
- `gitmap/helptext/chrome-profile-import.md`

##### Requirements
1. Document `--json`, `--file <path>`, `--fnf <path>`, and `--tempfile <filename>`.
2. Add explicit examples showing both `--json` and `--json --file <path>`.
3. Use fenced code blocks (` ``` `), strictly no 4-space indented blocks.
4. Pass `go test ./gitmap/helptext/... -run Golden -count=1`.

###### Subtask File: `04-task.md`

#### Subtask 04: End-to-End Tests & CI Verification

##### Objective
Author and execute comprehensive unit and end-to-end tests validating profile subcommand routing, JSON export, file writing, and preflight inspect behavior.

##### Affected Files
- `gitmap/cmd/profile_route_test.go` (new)
- `gitmap/cmd/chromeprofile_export_flags_test.go` (new)

##### Requirements
1. Test `gitmap profile import`, `import-all`, `export`, `export-all`, `inspect`, `preview`, `check` routes correctly.
2. Test `import-check` with `--json`, `--file <path>`, `--fnf <path>`, and `--tempfile <filename>`.
3. Verify files are created with valid JSON on disk.
4. Run all local tests and CI quality gates.


## 4. Unified Quality Gates & Verification Checklist

- [x] **Zero Concept Loss:** All source plans, code modifications, and execution steps preserved in full.
- [x] **Subtasks Inlined:** All associated subtasks folded directly into this document.
- [x] **Strict Relative Paths:** All citations use repository-relative paths without drive letters or file:/// URIs.
- [x] **Function Sizing:** All referenced codebase functions conform to <= 15 lines body cap.
- [x] **Coding Guidelines:** Affirmative booleans, zero nested ifs, and universal AppError wrapping verified.
- [x] **CI/CD Quality Gates:** All component tests pass legitimately under the local CI/CD runner.

## 5. Root Cause Analyses & Bug Fixes Referenced

- [`.lovable/memory/issues/2026-09-05-chrome-profile-picker-visibility-desync.md`](.lovable/memory/issues/2026-09-05-chrome-profile-picker-visibility-desync.md)
- [`.lovable/memory/learned/11-chrome-profile-import-routing-and-fnf-export.md`](.lovable/memory/learned/11-chrome-profile-import-routing-and-fnf-export.md)
