# Plan 84: Custom Installer Execution Engine & Dual CLI Parity

## 1. Overview & Context

This plan resolves custom installer tool execution and establishes complete dual CLI parity across `gitmap install` and `gitmap installer`:
1. **Custom Installer Execution Engine (`gitmap install <custom-tool>`):**
   - Resolves user-defined custom installer records from SQLite (`installer_scripts`) by slug.
   - Extracts and unmarshals multi-OS JSON instructions (`win`, `unix`, `ubuntu`, or fallback `all`).
   - Dynamically targets the host operating system:
     - **Windows (`windows`):** Dispatches PowerShell script via `powershell -NoProfile -ExecutionPolicy Bypass -Command <script>`.
     - **Linux (`linux`):** Distinguishes Ubuntu/Debian hosts via `/etc/os-release` to prefer Ubuntu-specific commands; otherwise falls back to Unix/All via `bash -c <script>`.
     - **macOS (`darwin`):** Executes Unix or generic script via `bash -c <script>`.
   - Records successful execution into `installation.db` via `splitDB.SaveInstalledTool(slug, version, "custom")`.
   - Refreshes `gitmap install ls` to display installed indicator (`●`) under the `Custom Tools` section.
2. **Dual CLI Command Parity (`gitmap install` & `gitmap installer`):**
   - Resolves argument delegation in `cmd/installer.go` and `cmd/roottooling.go`.
   - Unrecognized commands and tool names passed to `gitmap installer <name>` or shorthand `gitmap in <name>` delegate directly to `runInstall(args)`.
   - Guarantees 100% parity between `gitmap install` and `gitmap installer`.

---

## 2. Task-Specific Rules & Invariants

1. **Strict Relative Paths:** Never use absolute paths or `file:///` URIs in code, markdown, or plans.
2. **Coding Guidelines Adherence:** All Go functions <= 15 lines, blank line before every return, affirmative boolean naming (`is*`, `has*`), zero nested `if` blocks.
3. **Zero Error Swallowing:** All runtime and validation errors wrapped in `*apperror.AppError`.
4. **SQLite Concurrency & Anchoring:** Queries anchored to binary data directory via `store.OpenDefault()`.
5. **No Automatic Releases:** Standard development commit; no version bump or release tagging.

---

## 3. Implementation Details

- `gitmap/cmd/install_custom_exec.go`:
  - `findCustomInstaller(slug)`: Safe lookup in SQLite `installer_scripts`.
  - `hasCustomInstaller(slug)`: Boolean existence validator for tool filters.
  - `executeCustomInstaller(script, opts)`: Unmarshaling, dry-run support, process execution, and installation recording.
  - `resolveOSScriptForPlatform`: Host-aware OS script resolution.
  - `recordCustomInstallSuccess`: Persists tool status in `installation.db`.
- `gitmap/cmd/install.go`:
  - Integrated `findCustomInstaller` in `runInstall` before tool validation.
  - Updated `isKnownInstallTool` to recognize custom tools.
- `gitmap/cmd/installer.go`:
  - Added delegation in `runInstaller` and `RunInstallerCLI` to forward non-Cobra tool names to `runInstall`.
- `gitmap/cmd/roottooling.go`:
  - Cleaned up duplicate alias mapping.
- `gitmap/cmd/install_custom_exec_test.go`:
  - Unit test suite covering instruction parsing, OS resolution, live execution, dry-run, and integration.

---

## 4. Verification

1. Unit tests pass 100%:
   `go test -v ./cmd -run "TestParseInstructionsMap|TestResolveOSScriptForPlatform|TestExecuteCustomInstaller"`
2. End-to-end binary test:
   `gitmap install add alim1` -> `gitmap install alim1` -> prints banner, executes command, records success.
3. Dual parity verification:
   `gitmap installer alim1` executes identically.
4. Listing verification:
   `gitmap install ls` renders `● alim1 1.0.0 dsc` in the `Custom Tools` section.
5. All linters pass with 0 violations.
