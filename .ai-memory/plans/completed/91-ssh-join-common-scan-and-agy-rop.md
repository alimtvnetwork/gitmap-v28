# Completed Plan 91: SSH Batch Common Join, Subnet Scanner, OS Metadata, AGY Project Re-Read & Optimization, and AUM Benchmarking

Spec Reference: [02-spec/21-app/141-ssh-join-common-scan-and-agy-rop.md](../../02-spec/21-app/141-ssh-join-common-scan-and-agy-rop.md)

## User Request (Verbatim)

```text
gitmap ssh-join-common administrator ip1(aliasing),ip2(aliasing),lastchangedpartoftheIP(aliasing)
gitmap ssh-join-common administrator 192.168.1.3(w1),7(w2),12(w3) --pass <...>

gitmap ssh scan # see how many ips are available in the network

D:\work\03-aukgo\enum\osdetect

gitmap agy reread-optimize-project(rop) N # 5


Okay, fix this third issue for the Git map. Also, at the same time. Oh, sorry. Yeah. Fix this issue for the Git map. Okay, a couple of issues. So what I want using the SSH is that scan command that will try to scan the IPs in the network. For example, you try to connect all these machines, right? Together, connect, join. How can I connect them very easily? What is the process, okay, rather than doing it one at a time? Or probably I can do lots of IPs if they have same username. I could just provide that, and there should be an easier version to connect. And what I do think is that can we get SSH join command SJC, where we could just say administrator, and then provide multiple IPs with shorthand octet notation.
```

## Problem Summary & Architectural Scope
1. **SSH Batch Common Join (`gitmap ssh-join-common` / `sjc`)**:
   - Needed an ergonomic shorthand syntax for enrolling multi-node clusters sharing common credentials and base IP subnets.
   - Suffix/octet expansion: `192.168.1.3(w1),7(w2),12(w3)` dynamically resolves base prefix `192.168.1.` across subsequent octets (`192.168.1.7`, `192.168.1.12`).
   - Automated concurrent joining, credential distribution, host key recording, and SQLite enrollment.
2. **Subnet Discovery Scanner (`gitmap ssh scan`)**:
   - Needed network discovery across the active IPv4 subnet to find hosts listening on SSH port 22.
   - Rerouted default `ssh scan` from stale registered-node pinging to live subnet discovery with automatic `[ENROLLED: <alias>]` vs `[NEW]` identification.
3. **Deep OS Detection & `gitmap os-info`**:
   - During node join, probe remote OS type, kernel/build version, and CPU architecture (`amd64`, `arm64`) and persist into `installation.db` / `gitmap.db`.
   - Standalone CLI command `gitmap os-info [--json]` returning local OS metadata.
4. **Antigravity (AGY) Project Re-Read & Optimize (`gitmap agy rop [N]`)**:
   - Scans last N projects modified in the last 24 hours (default N=5).
   - Backs up conversations into isolated Split-DBs at `data/agy/<repo-slug>/agy.db` before pruning.
   - Trims cache to 10 retained items (`--keep 10`).
   - Hooks into `gitmap storage reset` for complete data lifecycle management.
   - Initiates fresh conversation titled with project name for codebase re-reading.
5. **AUM Benchmarking & E2E Isolation**:
   - Side-by-side performance benchmarks comparing native Go AUM search against legacy Python scripts (Go: 95ms vs Python: 14.35s = 151.1x speedup).
   - Documented in `readme.md`.
   - Local VM tests and heavy benchmarks guarded with `//go:build e2e` to keep CI/CD test runs green and fast.

## Implemented Remediation & Components

### 1. SSH Common Join (`cli/cmdssh/`)
- `sshjoin_common_types.go`: Core data models (`SSHCommonJoinInput`, `ParsedCommonNodeTarget`, `SSHCommonJoinResult`).
- `sshjoin_common_parser.go`: Parses shorthand octet targets (`192.168.1.3(w1),7(w2),12(w3)`) with parenthesized alias extraction and IP prefix state tracking.
- `sshjoin_common.go`: High-speed concurrent worker pool joining nodes via common credentials and enrolling into DB.
- `sshjoin_common_cmd.go`: CLI command wiring for `ssh-join-common` and alias `sjc`.

### 2. Subnet Scanner (`cli/cmdssh/`)
- `ssh_scanner.go`: Implemented `ExecuteSubnetScan` using concurrent goroutines with 1.2s timeout per target across IPv4 `/24` subnets.
- `ssh_scan_cmd.go`: Updated command router to execute `ExecuteSubnetScan` by default with `--cidr`, `--timeout`, and `--ports` flags.

### 3. Remote OS Probing & Local OS Info (`cli/cmdos/` & `cli/cmdssh/`)
- `sshjoin_enroll.go`: Added `probeRemoteOSArch` and `formatOSVersionWithArch` querying `systeminfo` / `uname` over SSH during enrollment.
- `cli/cmdos/os_info_types.go`: `OSInfoRecord` model with `OSType`, `Version`, `Arch`, `Hostname`, `Kernel`.
- `cli/cmdos/os_info.go`, `os_info_windows.go`, `os_info_other.go`: OS detection engine.
- `cli/cmdos/os_info_cmd.go`: CLI command `gitmap os-info [--json]`.

### 4. AGY Project Re-Read & Optimize (`cli/cmdagy/` & `cli/cmd/`)
- `agy_rop_types.go`: Models for `AGYROPConfig`, `AGYROPProjectSummary`, and `AGYROPResult`.
- `agy_rop_backup.go`: Split-DB SQLite backup engine writing to `data/agy/<slug>/agy.db` (`AGYConversationBackup` table) with safe handle closing.
- `agy_rop.go`: Active 24h project scanner, backup coordinator, cache trimmer, and conversation launcher.
- `agy_rop_cmd.go`: CLI command `gitmap agy reread-optimize-project` / `rop [N]`.
- `cli/cmd/storage_reset.go`: Integrated `purgeAGYBackups` into `gitmap storage reset`.

### 5. Documentation & E2E Isolation (`readme.md` & tests)
- `readme.md`: Added benchmark results (Go AUM vs Python), `gitmap sjc`, `gitmap ssh scan`, and `gitmap agy rop` guides.
- `cli/cmdautomation/aum_benchmark_e2e_test.go`: Marked with `//go:build e2e`.
- `cli/cmdssh/ssh_live_vm_e2e_test.go`: Live VM enrollment and execution tests marked with `//go:build e2e`.

## Verification Outcomes
- **Unit Tests**:
  - `cli/cmdssh/ssh_common_join_test.go`: All tests pass (validates octet expansion, alias extraction, error handling).
  - `cli/cmdssh/ssh_scanner_test.go`: Passes.
  - `cli/cmdos/os_info_test.go`: Passes.
  - `cli/cmdagy/agy_rop_test.go`: Passes (validates backup DB lifecycle, file handle release, and reset).
- **Quality Gates**:
  - `python 03-ai-scripts/26-go-code-formatter.py`: 3612 Go files formatted.
  - `python linter-scripts/check-nested-ifs.py`: 0 violations.
  - `python linter-scripts/check-enum-and-boolean.py`: 0 violations across 2764 files.
  - `python linter-scripts/check-error-management.py`: 0 bare panics/exits across 3689 files.
  - `python linter-scripts/check-relative-paths.py`: 0 absolute path / URI violations across 7538 files.
- **Local E2E Tests**:
  - `go test -v -tags=e2e ./cli/cmdautomation`: Passes (Go: 95ms vs Python: 14.35s).
  - `go test -v -tags=e2e ./cli/cmdssh -run TestLiveVMBatchJoinAndExecute`: Passes across live VMs (`w1`, `w2`, `w3`).
