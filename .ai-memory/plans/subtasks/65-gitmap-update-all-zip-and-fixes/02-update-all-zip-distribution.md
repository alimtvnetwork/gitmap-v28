# Subtask 65.2: Multi-Node Update Engine & Zip SCP Distribution

> **Parent Plan:** [65-gitmap-update-all-zip-and-fixes.md](../../pending/65-gitmap-update-all-zip-and-fixes.md)  
> **Tracking Spec:** [02-component-spec.md](../../../../02-spec/21-app/65-gitmap-update-all-zip-and-fixes/02-component-spec.md)  
> **Status:** Pending  
> **Primary File Targets:** `cli/cmd/rootutility.go`, `cli/cmdupdate/update_fleet.go`, `cli/cmdupdate/update_fleet_test.go`  

---

## 1. Objective

Implement multi-node zip packaging and streaming SCP distribution for fleet updates across cluster nodes:
1. Support commands:
   - `gitmap update all --include-others`
   - `gitmap update all zip --include-others`
   - `gitmap update-all-zip --include-others`
   - `gitmap uaz --include-others` (shorthand alias)
2. Package updated binary into an in-memory zip archive locally (for `gitmap` or `agm`).
3. Stream the archive across SSH sessions to target remote temporary directory via `cmdssh.StreamFileToRemote`.
4. Trigger remote installation scripts (PowerShell for Windows, Shell/Bash for POSIX) to extract and install the update.
5. Ingest JSON telemetry from remote nodes and render in the unified summary table.

---

## 2. Implementation Scope

### 2.1 CLI Routing (`cli/cmd/rootutility.go`)
- In `utilityCoreEntries()`, register dispatch entries for:
  - `"uaz"`
  - `"update-all-zip"`
  - `"updateallzip"`
  routing to `runUpdateHelp`.
- In `runUpdateHelp()`, ensure `cmdupdate.IsFleetUpdateCommand(cmdName, args)` handles `uaz` and `update-all-zip` tokens.

### 2.2 Fleet Update Engine (`cli/cmdupdate/update_fleet.go`)
- **Extend `FleetUpdateOptions`:**
  - Add `IsZip bool`
  - Add `IncludeOthers bool`
- **Flag Parsing (`processUpdateFlag`):**
  - Check for `--include-others` / `--includeothers` -> set `opts.IncludeOthers = true`.
  - Check for `zip` positional argument or `--zip` flag -> set `opts.IsZip = true`.
- **Target Resolution:**
  - When `opts.IncludeOthers` is enabled, merge targets from `cmdssh.FetchAllSSHConnections()` with discovered cluster nodes from DB, deduplicating by IP address.
- **Packaging Pipeline (`createUpdatePackageZip`):**
  - Read local binary (`gitmap` or `agm`).
  - Create zip archive in memory with executable and install launcher script.
- **Remote Execution (`executeSSHFleetZipUpdate`):**
  - Stream zip to `C:\Windows\Temp\` (Windows) or `/tmp/` (POSIX) using `cmdssh.StreamFileToRemote`.
  - Execute remote PowerShell or Shell command to expand archive, overwrite binary, and output JSON telemetry.
  - Return JSON telemetry string.

### 2.3 Unit Testing (`cli/cmdupdate/update_fleet_test.go`)
- Add `TestExecuteFleetUpdate_ZipDistribution`:
  - Verify that passing `zip` / `uaz` / `update-all-zip` triggers zip packaging and remote stream execution.
  - Verify mock targets receive the streaming payload and return success telemetry.
- Add `TestExecuteFleetUpdate_IncludeOthers`:
  - Verify that `--include-others` expands target set to include additional cluster nodes.

---

## 3. Coding Guidelines & Constraints
- All boolean fields and flags must be positive (`IsZip`, `IncludeOthers`, `IsSuccess`, `IsOnline`).
- Functions must remain small (<= 15 lines per function) and focused.
- Wrap all errors using `apperror.WrapSimple` or `apperror.NewSimple`.
- Zero raw unindexed commands; adhere strictly to repository standards.

---

## 4. Verification Steps
- Run `go test ./cli/cmdupdate/...` to verify existing and new tests pass.
- Verify CLI command routing in dry-run mode: `gitmap uaz --dry-run`.
