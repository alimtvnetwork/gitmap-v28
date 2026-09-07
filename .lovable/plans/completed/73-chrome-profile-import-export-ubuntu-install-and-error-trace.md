# 73-chrome-profile-import-export-ubuntu-install-and-error-trace.md: Chrome Profile Multi-Directory/ZIP Discovery, Ubuntu Chrome Installer & Rich Error Diagnostics

## 1. Executive Summary

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

## 2. Task-Specific Rules & Constraints

1. **Rule 1 (Strict Relative Paths):** Zero `file:///` URIs or machine-specific absolute paths in markdown files, artifacts, or code. All paths must be relative to repository root (`/`).
2. **Rule 2 (Coding Guidelines Compliance):** Functions <= 15 lines (prefer <= 8), blank line before every return statement, affirmative booleans (`is_*`, `has_*`), no magic numbers, and zero swallowed errors (`err != nil` must be handled or wrapped via `apperror`).
3. **Rule 3 (Strict 200-Line File Limit):** All new and modified Go source files must remain <= 200 lines. Decompose logic into focused single-responsibility files.
4. **Rule 4 (AST Parity & Constants):** All CLI verbs and help entries must be synchronized with `gitmap/constants/constants_cli.go` and verified with `TestTopLevelCmdRegistryMatchesAST` and `TestEveryHelpFileHasExamples`.
5. **Rule 5 (CI/CD Local Runner Validation):** Must run `python 03-ai-scripts/06-cicd-local-runner.py` exit 0 before release orchestration.

## 3. Subtasks Breakdown

- [01-chrome-profile-polymorphic-discovery-and-preflight-ls.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/01-chrome-profile-polymorphic-discovery-and-preflight-ls.md)
- [02-chrome-profile-zip-multi-profile-extraction.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/02-chrome-profile-zip-multi-profile-extraction.md)
- [03-chrome-helptext-and-alias-routing.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/03-chrome-helptext-and-alias-routing.md)
- [04-ubuntu-chrome-installer-and-special-handler.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/04-ubuntu-chrome-installer-and-special-handler.md)
- [05-enhanced-error-reporting-and-stack-traces.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/05-enhanced-error-reporting-and-stack-traces.md)
- [06-vmware-help-and-dryrun-integration.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/06-vmware-help-and-dryrun-integration.md)
- [07-quality-gates-tests-and-release.md](../subtasks/73-chrome-profile-import-export-ubuntu-install-and-error-trace/07-quality-gates-tests-and-release.md)
