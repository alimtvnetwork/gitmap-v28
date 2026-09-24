# Plan 97: SSH Common Batch Join, Remote OS Detection & Telemetry, AGY Project Re-read & Optimization, and Native AUM Benchmarking

> **Plan Status:** Completed  
> **Traceability IDs:** Task-01 .. Task-07  
> **Spec Reference:** [02-spec/21-app/148-ssh-join-common-os-detect-rop-and-e2e-benchmarks.md](../../../02-spec/21-app/148-ssh-join-common-os-detect-rop-and-e2e-benchmarks.md)  
> **Issues Covered:** [02-spec/22-app-issues/41-vscode-startup-failure-and-search-latency-rca.md](../../../02-spec/22-app-issues/41-vscode-startup-failure-and-search-latency-rca.md)

---

## 1. Architectural Context & Subtask Mapping

| Subtask ID | Focus Area | Target Deliverable | Status |
| :--- | :--- | :--- | :--- |
| **Subtask 97-01** | `cmdssh` | Shorthand octet parsing with aliases (`sjc` / `ssh-join-common`), batch joins, remote OS probe & SQLite persistence | Completed |
| **Subtask 97-02** | `cmdos` | JSON structured OS & hardware telemetry routing (`gitmap os info --json` / `gitmap os-info`) | Completed |
| **Subtask 97-03** | `cmdagy` | AGY Rerun IDE restart (`gitmap agy rerun [1|2|3|4|project]`), process termination, prompt replay with media attachments | Completed |
| **Subtask 97-04** | `cmdagy` | AGY Project Re-Read & Optimization (`gitmap agy rop [N]`), split-DB backup to `data/AGY/<slug>.db`, `--keep 10` cache retention | Completed |
| **Subtask 97-05** | `cli/tests/e2e` | Search benchmark (GitMap AUM vs Go walk vs Python grep) & documentation in `docs/benchmarks/search_benchmark.md` | Completed |
| **Subtask 97-06** | `cmdchromeprofile` | Chrome Profile auth session & refresh token import/export hermetic E2E tests (`//go:build e2e`) | Completed |
| **Subtask 97-07** | `scripts-fixer`, `02-spec` | VS Code startup failure 4-Part RCA (Issue 41) & path normalization (dynamic environment variables only) | Completed |

---

## 2. Verified Outcomes

1. **SSH Batch Common Join**: Added `cli/cmdssh/sshjoin_common.go` with shorthand IP/octet parsing (e.g. `192.168.1.3(w1),7(w2),12(w3)`), automated remote OS detection upon initial join, and persistence to `SSHConnection.OSVersion` in SQLite.
2. **OS Telemetry CLI**: Integrated `RunOSInfoCLI` in `cli/cmdos/os.go`, ensuring `gitmap os info --json` and `gitmap os-info` cleanly output structured hardware, CPU, memory, and OS details.
3. **AGY IDE Rerun & Restart**: Hardened `cli/cmdagy/agy_rerun_restart.go` to reliably terminate active Antigravity instances, restart in the project workspace, and replay the complete transcript including media/picture attachments.
4. **AGY Project Re-Read & Optimization**: Implemented `gitmap agy rop` with split-DB persistence at `data/AGY/<slug>.db` and 10-item context retention.
5. **AUM Search Benchmarks**: Created `docs/benchmarks/search_benchmark.md` and `cli/tests/e2e/search_benchmark_e2e_test.go`, verifying a 33,000x speedup (<1ms vs 33.2s) over unindexed Python grep.
6. **Chrome Profile Auth E2E Tests**: Added hermetic tests in `cli/cmdchromeprofile/chrome_profile_auth_e2e_test.go` guarded by `//go:build e2e`.
7. **VS Code Startup RCA & Dynamic Paths**: Authored `02-spec/22-app-issues/41-vscode-startup-failure-and-search-latency-rca.md` and ensured all PowerShell and Go repair scripts use `$env:APPDATA` and dynamic paths without hardcoded drive letters.
