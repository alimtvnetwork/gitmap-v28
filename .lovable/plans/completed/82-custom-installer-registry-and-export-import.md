# Plan 82: Custom Installer Interactive Registry, Dual CLI Parity, JSON/ZIP Export & Import, and Dynamic Install LS

## 1. Overview & Context

This plan delivers the requested custom installer management system for Gitmap:
1. **Interactive Installer Creation (`gitmap install add <name> [version]` & `gitmap installer add <name> [version]`):**
   - Macro-style interactive prompt flow:
     - Installer description.
     - Windows commands/script (PowerShell).
     - Unix generic commands/script (bash/sh).
     - Ubuntu-specific commands/script (apt/bash).
     - Confirmation prompt: "Do you want to add or edit Unix or Ubuntu instructions later? (y/n)".
   - Saves record to SQLite (`installer_scripts`) with multi-OS payload in `model.InstallerScript.Scripts`.
   - Also supports non-interactive CLI flags (`--desc`, `--win`, `--unix`, `--ubuntu`, `--yes`).
2. **Dual CLI Command Parity (`install` and `installer`):**
   - `gitmap install add` <-> `gitmap installer add` (and alias `create`).
   - `gitmap install export` <-> `gitmap installer export` (and `export-all`).
   - `gitmap install import` <-> `gitmap installer import`.
3. **Flexible Export & Import Formats (Single JSON & ZIP):**
   - Export:
     - JSON single file format: `gitmap install export <slug> -o installer.json` or `--format json`.
     - Export all to JSON: `gitmap install export --all -o all-installers.json`.
     - ZIP bundle format: `gitmap install export <slug> -o bundle.zip` or `--format zip`.
   - Import:
     - Auto-detects `.json` vs `.zip`:
       - Single installer JSON or JSON array.
       - ZIP archive containing installer JSON files.
4. **Dynamic `gitmap install ls` Custom Section:**
   - Automatically queries registered custom installers from SQLite.
   - Dynamically renders a `Custom Tools` / `Custom Installers` section in `gitmap install ls`.
   - Shows status dot, slug/name, version, and description.

---

## 2. Task-Specific Rules & Invariants

1. **Strict Relative Paths:** Never use absolute paths or `file:///` URIs in code, markdown, or plans.
2. **Coding Guidelines Adherence:** All Go functions must be <= 15 lines with a mandatory blank line before every return statement, affirmative booleans (`is*`, `has*`), and zero nested `if` blocks.
3. **Error Management:** All errors must use `*apperror.AppError` with distinct domain error codes (`E_INSTALLER_*`). Zero bare panics, `os.Exit`, or swallowed errors.
4. **SQLite Concurrency & Anchoring:** Use `store.OpenDefault()` with `EvalSymlinks` path anchoring and `SetMaxOpenConns(1)` for clean thread-safe DB transactions.
5. **Interactive & Scripted Fallbacks:** Interactive prompts must check `isTerminalInput()` and allow piped/flag non-interactive execution without hanging CI runners or scripts.

---

## 3. Subtask Decomposition

- [01-task-interactive-installer-add.md](../subtasks/82-custom-installer-registry-and-export-import/01-task-interactive-installer-add.md): Implement interactive prompt recorder for `install add` / `installer add` capturing description, win, unix, ubuntu instructions, and confirmation.
- [02-task-export-import-json-zip.md](../subtasks/82-custom-installer-registry-and-export-import/02-task-export-import-json-zip.md): Implement single JSON and ZIP export/import in `install` and `installer` commands with format auto-detection.
- [03-task-custom-install-ls-integration.md](../subtasks/82-custom-installer-registry-and-export-import/03-task-custom-install-ls-integration.md): Integrate custom installers from SQLite into `gitmap install ls` under a dedicated `Custom Tools` category block.
- [04-task-verification-and-ci-gates.md](../subtasks/82-custom-installer-registry-and-export-import/04-task-verification-and-ci-gates.md): Unit tests for add, export, import, and install ls; compile and verify local CI runner exits 0.

---

## 4. Verification Plan

1. Test CLI interactive and non-interactive `gitmap install add "my-tool" v1.2 --desc "Test description" --win "winget install my-tool" --ubuntu "sudo apt install my-tool"`.
2. Test `gitmap install export my-tool -o my-tool.json` and verify JSON payload.
3. Test `gitmap install export --all -o all-tools.zip` and verify zip entries.
4. Test `gitmap install import my-tool.json`.
5. Test `gitmap install ls` to verify `Custom Tools` category block shows `my-tool` with version and description.
6. Verify unit tests pass: `go test -v ./cmd -run "TestInstallAdd|TestInstallExport|TestInstallImport|TestCustomInstallList"`.
7. Verify linters pass (`check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-relative-paths.py`, `check-error-management.py`).
