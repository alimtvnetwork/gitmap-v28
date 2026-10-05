# RCA-100: Linter Policy Violations, Error Management, and Unused Symbols Remediation

**Date:** 2026-10-05
**Status:** ✅ Resolved
**Severity:** High (CI/CD Pipeline Failure across Policy Linters and Full Suite Guard)
**Affected Workflows:** `CI (#37228760212)`
**Run URLs:**
- `https://github.com/alimtvnetwork/gitmap-v28/actions/runs/37228760212`

---

## 1. Symptom

On commit `676258d7`, GitHub Actions CI run `#37228760212` failed on multiple jobs:

1. **Boolean Guidelines Linter:**
   - 5 inverted success checks (`!isSuccess`) across `cli/cmdagy/agy_deploy_cmd.go`, `cli/cmdagy/agy_prompt_dispatch.go`, `cli/cmdinstall/install_tools_flag.go`, `cli/cmdnodes/nodes_send_projects.go`, and `cli/cmdnodes/nodes_sync_settings.go`.
2. **Error Management Linter:**
   - Swallowed database query errors detected in `cli/cmdagent/agent_diagnose.go:295` and `cli/cmdagent/agent_diagnose.go:296` without proper annotation or handling.
3. **Full Suite Guard (golangci-lint & staticcheck):**
   - Unused package-level variables in `cli/firewall/windows.go` (`winPortRegex`, `winRuleRegex`).
   - Unused dead functions in `cli/cmddb/cmddb_reset.go` (`performDbReset`, `clearSplitDbFiles`) and `cli/cmd/reset.go` (`parseResetFlags`, `executeReset`, `removeActiveDbFile`, `activeDbPath`, `reseedFromJSON`).
   - Deprecated `strings.Title` call (SA1019) in `cli/cmdapps/apps.go:358`.
4. **Nested If Linter:**
   - 90 nested-if / anti-compression violations (depth > 1) across 27 source files.

---

## 2. Root Cause

1. **Inverted Boolean Success Checks:**
   - New fleet and prompt dispatch logic checked failure states using negative boolean conditions (`if !res.IsSuccess`), violating repository guidelines requiring positive boolean flags.
2. **Unannotated Best-Effort Database Queries:**
   - Telemetry count queries in `agent_diagnose.go` discarded database scan errors via blank identifiers (`_ = db.QueryRow(...)`) without the required `// lint-allow: ignore-db-error` annotation.
3. **Dead Code and Deprecated Standard Library API:**
   - After migrating database reset operations to unified handlers, obsolete helper functions in `cli/cmd/reset.go` and `cli/cmddb/cmddb_reset.go` remained uncalled.
   - Deprecated `strings.Title` was flagged by staticcheck in Go 1.24/1.25.
4. **Deeply Nested Conditionals and Single-Line Ifs:**
   - Command dispatchers and JSON error formatters nested `if isJson` and error checking inside parent condition blocks instead of utilizing early returns and guard clauses.
   - Embedded HTML template scripts in `cli/cmdagy/agy_ui_html.go` contained single-line `if` statements.

---

## 3. Resolution

1. **Enforced Positive Boolean Logic:**
   - Refactored `if !res.IsSuccess` to positive branching (`if res.IsSuccess { ... } else { ... }` or explicit failure booleans) across all 5 affected files.
2. **Annotated Best-Effort Database Queries:**
   - Added `// lint-allow: ignore-db-error` comments above non-critical count queries in `cli/cmdagent/agent_diagnose.go`, `agent_cleanup.go`, `agent_log.go`, and `cli/store/agent_store.go`.
3. **Purged Dead Symbols and Replaced Deprecated Calls:**
   - Removed unused regex variables and dead reset helpers.
   - Implemented rune-aware `toTitleCase()` in `cli/cmdapps/apps.go` to eliminate SA1019.
4. **Flattened Nested Conditionals:**
   - Refactored nested if statements across 27 files into top-level guard clauses and isolated helper functions (`printAppsListJSON`, `printAppsUninstallJSON`, `handleCDScanFallback`, etc.).
   - Expanded single-line `if` statements in `cli/cmdagy/agy_ui_html.go` into multiline blocks.
5. **Symbol and Type Parity:**
   - Corrected parameter types for `handleInstallToolFailure` in `cli/cmdinstall/install_tools_flag.go`.
   - Renamed conflicting package-level helper in `cli/cmdssh/ssh_client.go` to `resolveTargetHostAlias`.

---

## 4. Verification

- `python linter-scripts/check-nested-ifs.py --all`: ✅ 0 violations across 4,181 files.
- `python linter-scripts/check-boolean-guidelines.py --all`: ✅ 0 violations across 4,181 files.
- `python linter-scripts/check-enum-and-boolean.py`: ✅ 0 violations across 3,191 files.
- `python linter-scripts/check-error-management.py --all`: ✅ 0 violations across 4,241 files.
- `python linter-scripts/check-relative-paths.py`: ✅ 0 violations across 8,658 files.
- `go vet ./...` in `cli/`: ✅ Clean (exit code 0).
- `golangci-lint run --issues-exit-code=1 ./...` in `cli/`: ✅ Clean (exit code 0).
