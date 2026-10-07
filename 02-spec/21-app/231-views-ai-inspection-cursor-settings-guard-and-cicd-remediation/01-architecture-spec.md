# Architecture Specification: Views AI Inspection, Cursor Settings Guard & CI/CD Remediation

> **Document Version:** 1.0.0  
> **Status:** Active  
> **Scope:** GitMap Core CLI & Web UI (`cli/cmdcursor`, `cli/cmdagy`, `src/pages/Settings.tsx`, `.github/workflows`)  
> **Traceability:** Task-231  

---

## 1. Executive Summary & Problem Statement

Following recent rapid iterations across multi-instance Antigravity discovery, UI modernization, and GitHub login/logout authentication, another automated agent implemented:
1. **Views AI & Settings UI Modernization**:
   - `src/pages/Settings.tsx`: Introduced multi-instance Antigravity discovery panels, active instance pills, and CSS3 depth catalogs.
   - `cli/cmdagy/agy_instance_discovery.go`: Multi-instance process scanning and prompt query engines.
   - `cli/cmdui/ui_server.go`: REST endpoints for `/api/instances` and `/api/prompts/instances`.
2. **Cursor Settings Guard**:
   - Uncommitted modifications in `cli/cmdcursor/cursor_settings.go` and `cli/cmdcursor/cursor_settings_test.go` introducing `requireCursorSettingsFile` to safeguard `gitmap cursor settings sync` against missing settings files with actionable remediation guidance.

During codebase audit and linter inspection, several issues were identified:
- In `src/pages/Settings.tsx`, lines 181-187 contained a nested `if` statement violating repository coding standards:
  ```ts
  if (response.ok) {
    const json = await response.json();
    if (json.isSuccess && Array.isArray(json.instances)) {
      setInstances(json.instances);
      return;
    }
  }
  ```
- In `cli/cmdcursor/cursor_settings.go` and `cursor_settings_test.go`, vertical newline conventions (blank line after `}` before statements and before `return`) were missing.
- In `cli/cmdagy/agy_instance_discovery.go`, vertical newline spacing after closing braces in `unmarshalWinProcItems` and `truncateSliceByLimit` needed enforcement.

This task resolves these coding guideline violations, commits the Cursor settings guard, validates all unit tests and linters locally, performs a minor version bump (v6.505.0 -> v6.506.0), executes the release ceremony, and confirms 100% green CI/CD pipeline health.

---

## 2. Component Boundaries & Scope

### 2.1 Cursor Settings Guard (`cli/cmdcursor/cursor_settings.go`)
- Adds `requireCursorSettingsFile(path string) error`:
  - Validates that the Cursor `settings.json` file exists on the local filesystem before initiating remote synchronization to target cluster nodes (`u1`, `w1`, etc.).
  - Returns structured `apperror.NewWithDetails` ("cmd.cursor.settings.sync", "E1037") suggesting the user run `gitmap cursor settings apply` if missing.
  - Accompanied by unit tests in `cli/cmdcursor/cursor_settings_test.go`:
    - `TestRequireCursorSettingsFileMissing`
    - `TestRequireCursorSettingsFilePresent`

### 2.2 Web UI Guard Inversion (`src/pages/Settings.tsx`)
- Inverts `response.ok` check into an early guard clause:
  ```ts
  if (!response.ok) {
    return;
  }
  const json = await response.json();
  if (json.isSuccess && Array.isArray(json.instances)) {
    setInstances(json.instances);
    return;
  }
  ```

### 2.3 Style & Vertical Spacing Remediation
- Ensures blank lines after `}` blocks across `cli/cmdcursor/` and `cli/cmdagy/agy_instance_discovery.go`.
- Validates via `linter-scripts/check-nested-ifs.py` and `linter-scripts/check-enum-and-boolean.py`.

---

## 3. Verification & Quality Gates

1. **Unit Test Verification**:
   - `go test -v ./cli/cmdcursor/...` (100% pass)
   - `go test -v ./cli/cmdagy/...` (100% pass)
   - `go test -v ./cli/cmd...` (100% pass)
2. **Linter Verification**:
   - `python linter-scripts/check-nested-ifs.py` (0 violations)
   - `python linter-scripts/check-enum-and-boolean.py` (0 violations)
   - `python linter-scripts/check-relative-paths.py` (0 violations)
3. **Release & CI/CD Telemetry**:
   - Minor version bump: `python 03-ai-scripts/37-bump-version.py --tier minor`
   - Atomic commit & push: `gitmap cpf`
   - Automated 5-step release branching & tag creation
   - Post-release pipeline monitoring via `gitmap pe -t`
