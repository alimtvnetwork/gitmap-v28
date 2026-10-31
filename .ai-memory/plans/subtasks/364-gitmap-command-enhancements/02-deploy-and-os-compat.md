# Subtask 02: OS-Aware Binary Deploy & Installation Engine

- **Parent Plan:** [363-gitmap-command-enhancements.md](../../pending/363-gitmap-command-enhancements.md)
- **Specification:** [02-deploy-and-os-compat-spec.md](../../../../02-spec/21-app/363-gitmap-command-enhancements/02-deploy-and-os-compat-spec.md)
- **Status:** `pending`
- **Owner:** GitMap Deploy and OS Compatibility Spec Author
- **Date:** 31-Oct-26

---

## 1. Overview & Objective

Upgrade the `gitmap deploy` subsystem from raw file transport into an intelligent, fleet-aware binary distribution and installation engine. The engine must:
1. Detect binary format via magic bytes (`MZ` for Windows PE, `\x7fELF` for Linux ELF, Mach-O for Darwin).
2. Compare binary format against target `node.OS` (and probe fallback) to prevent non-runnable deployments.
3. Automatically skip incompatible nodes (e.g. Windows `.exe` skipped on Ubuntu) while supporting explicit override with `--force` or `--force-copy`.
4. Run fast concurrent TCP liveness probes (`filterOnlineNodesConcurrently`) to gracefully bypass offline nodes without halting fleet execution.
5. Provide actionable post-execution diagnostic dockets and suggestions (cross-compile instructions, ping/wake commands, and status checks).
6. Support automated remote installation hooks with `--install` / `-i`.

---

## 2. Implementation Steps

### Step 1: Magic Byte Detection Engine (`cli/cmdssh/ssh_deploy_magic.go`)
1. Create `DetectBinaryFormat(filePath string) (BinaryFormatType, *apperror.AppError)`.
2. Implement header buffer read (up to 16 bytes).
3. Identify Windows PE by `0x4D, 0x5A` (`"MZ"`).
4. Identify Linux ELF by `0x7F, 0x45, 0x4C, 0x46` (`\x7fELF`).
5. Identify Darwin Mach-O by `0xFE, 0xED, 0xFA, 0xCE`, `0xFE, 0xED, 0xFA, 0xCF`, or `0xCA, 0xFE, 0xBA, 0xBE`.
6. Identify shell scripts by `0x23, 0x21` (`"#!"`).
7. Provide file extension fallback (`.exe`, `.dll`) for universal compatibility.

### Step 2: Node OS Normalization & Compatibility Matrix (`cli/cmdssh/ssh_deploy_compat.go`)
1. Define `OperatingSystemType` enum (`OSTypeWindows`, `OSTypeLinux`, `OSTypeDarwin`, `OSTypeUnknown`).
2. Implement `NormalizeNodeOS(rawOS string) OperatingSystemType`.
3. Implement `CheckBinaryOSCompatibility(binType BinaryFormatType, nodeOS string) (bool, string)`.
4. Ensure positive boolean return value `isCompatible` (never negative).
5. Generate clear human-readable mismatch reasons (e.g. `"Windows PE executable cannot execute natively on Linux/Ubuntu"`).

### Step 3: Fast Concurrent TCP Reachability Probe (`cli/cmdssh/ssh_deploy_liveness.go`)
1. Implement `filterOnlineNodesConcurrently(nodes []db.SSHConnection, timeout time.Duration) ([]db.SSHConnection, []DeploySkippedRecord)`.
2. Run probes in concurrent goroutines with bounded timeouts (default: 1500ms).
3. Safely partition targets into `onlineNodes` and `skippedRecords` (`SkipReasonNodeOffline`).
4. Prevent fleet-wide stalling when one or more cluster nodes are offline or suspended.

### Step 4: Post-Execution Diagnostic Feedback Docket (`cli/cmdssh/ssh_deploy_docket.go`)
1. Define `DeployExecutionDocket`, `DeployNodeSuccessRecord`, and `DeploySkippedRecord`.
2. Implement terminal summary renderer:
   - Results table with Node, IP, OS, Status, and Detail columns.
   - Succeeded, Skipped (OS), Skipped (Offline), and Failed totals.
   - Actionable remediation advice:
     - Cross-compilation commands (e.g. `GOOS=linux GOARCH=amd64 go build ...`).
     - Force transfer override advice (`--force-copy`).
     - Offline node diagnostics (`gitmap ping <node>`, `gitmap ssh-check <node>`).
3. Support `--json` flag to serialize the complete docket structure to standard output.

### Step 5: Command Wiring & Flag Integration (`cli/cmdssh/ssh_deploy_cmd.go` & `ssh_deploy_engine.go`)
1. Add CLI flags:
   - `--force` / `-f`: Override OS compatibility skip.
   - `--force-copy`: Direct alias for `--force`.
   - `--install` / `-i`: Trigger automated post-transfer installation.
   - `--timeout`: Bounded TCP reachability probe timeout.
2. Update `parseDeployFlags` and `DeployOptions` struct with positive booleans:
   - `HasForceOverride bool`
   - `HasInstallHook bool`
3. Update `ExecuteDeploy` to run binary format detection, node liveness filtering, and OS compatibility verification prior to transfer.
4. Execute remote installation script when `HasInstallHook` is `true`.

### Step 6: Unit & Integration Verification (`cli/cmdssh/ssh_deploy_compat_test.go`)
1. Test `DetectBinaryFormat` with mock PE, ELF, Mach-O, and script headers.
2. Test `CheckBinaryOSCompatibility` against complete compatibility matrix.
3. Test `filterOnlineNodesConcurrently` with mock TCP listeners and closed ports.
4. Test post-execution docket formatting and JSON output.

---

## 3. Strict Coding Standards & Invariants

1. **Error Handling:** All leaf functions return `*apperror.AppError` or standard `error`. Never call `cliexit.HandleError` or `os.Exit` inside leaf routines.
2. **Boolean Naming:** Use `is` or `has` prefixes exclusively (`isCompatible`, `isOnline`, `hasForceOverride`, `hasInstallHook`, `isDryRun`). Banned: `can`, `should`, `was`, `will`, and negative booleans (`isNot*`, `isOffline`, `hasNo*`).
3. **Function Sizing:** <= 15 lines max per function, <= 8 lines preferred.
4. **Conditional Flattening:** Zero nested if statements; enforce guard clauses and early returns.
5. **Path Formatting:** 100% relative paths in documentation and code comments.

---

## 4. Acceptance Criteria Checklist

- [ ] `DetectBinaryFormat` successfully differentiates Windows PE (`MZ`) and Linux ELF (`\x7fELF`).
- [ ] Windows `.exe` deployment to Ubuntu node `u1` is automatically skipped by default with `SkipReasonOSIncompatible`.
- [ ] Windows `.exe` deployment to Ubuntu node `u1` succeeds when `--force` or `--force-copy` is passed.
- [ ] Offline nodes are detected within probe timeout and logged as `SkipReasonNodeOffline` without halting the command.
- [ ] Post-execution docket prints formatted summary table, count breakdown, and cross-compilation suggestions.
- [ ] `--install` triggers target OS-specific installation and permission commands.
