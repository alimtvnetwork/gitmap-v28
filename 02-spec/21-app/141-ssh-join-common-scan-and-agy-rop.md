# Specification 141: SSH Batch Common Join, Subnet Scanner, OS Metadata, AGY Project Re-read & Optimization, and AUM Benchmarking

## Status: Draft / Active
## Spec ID: 141-ssh-join-common-scan-and-agy-rop
## Traceability: Task-01, Task-02, Task-03, Task-04, Task-05

---

## 1. User Request (Verbatim)

```text
gitmap ssh-join-common administrator ip1(aliasing),ip2(aliasing),lastchangedpartoftheIP(aliasing)
gitmap ssh-join-common administrator 192.168.1.3(w1),7(w2),12(w3) --pass <...>

gitmap ssh scan # see how many ips are available in the network

D:\work\03-aukgo\enum\osdetect

gitmap agy reread-optimize-project(rop) N # 5


Okay, fix this third issue for the Git map. Also, at the same time. Oh, sorry. Yeah. Fix this issue for the Git map. Okay, a couple of issues. So what I want using the SSH is that scan command that will try to scan the IPs in the network. For example, you try to connect all these machines, right? Together, connect, join. How can I connect them very easily? What is the process, okay, rather than doing it one at a time? Or probably I can do lots of IPs if they have same username. I could just provide that, and there should be an easier version to connect. And what I do think is that can we get SSH join command SJC, where we could just say administrator, and 192.168.1.3(w1),7(w2),12(w3) --pass <...>
```

---

## 2. Overview & Architecture

This specification addresses 5 critical capabilities extending GitMap's SSH Cluster Fleet, Antigravity (AGY) Project Maintenance, and Automation Manager (AUM):

1. **Batch Common SSH Join (`gitmap ssh-join-common` / `sjc`)**:
   - Enables rapid batch onboarding of multiple SSH nodes sharing the same administrative credential.
   - Syntax: `gitmap ssh-join-common <username> <ip-tokens> [--pass <password>] [flags]` (short-form alias: `gitmap sjc`).
   - Supports shorthand octet notation: the first IP provides the base IPv4 subnet prefix (e.g. `192.168.1.3(w1)` sets prefix `192.168.1.`), and subsequent tokens may specify only the modified host octet (e.g. `7(w2)` resolves to `192.168.1.7`, alias `w2`; `12(w3)` resolves to `192.168.1.12`, alias `w3`). If a subsequent token specifies a full IPv4 address, the base prefix is dynamically updated.
   - Captures node alias enclosed in parentheses `(<alias>)`. If omitted, defaults to `node-<IP>`.
   - Uses common username and `--pass` password for all targets, authenticating and persisting hosts in the SQLite database (`installation.db` and `gitmap.db`).

2. **Network Subnet Discovery Scanner (`gitmap ssh scan`)**:
   - Directly scans the local IPv4 subnet (auto-detected `/24` or user-specified CIDR via `gitmap ssh scan [cidr]`) to discover active reachable machines with SSH port 22 open.
   - Concurrently probes targets using a lightweight worker pool with bounded timeouts (e.g. 500-800ms per host).
   - Renders a clean, padded terminal table indicating IP address, open port, online status, probe latency, reverse hostname, enrollment status (`[ENROLLED: <alias>]` vs `[NEW]`), and suggested connection action.

3. **Remote OS Detection & Host Info Engine**:
   - During node onboarding (`ssh-join-common` and standard join), executes automated remote OS detection.
   - Accurately captures OS type (`windows`, `linux`, `macos`, `ubuntu`, `debian`), OS version (e.g. `Windows 11 Pro 10.0.22631`, `Ubuntu 22.04.3 LTS`), and OS architecture (`amd64`, `arm64`).
   - Persists OS metadata into the SQLite database (`SSHConnection` and `ssh_hosts` tables).
   - Provides a standalone CLI command: `gitmap os-info [--json]` for local machine OS profiling and remote node query execution.

4. **Antigravity (AGY) Project Re-Read & Optimization (`gitmap agy reread-optimize-project` / `rop`)**:
   - Syntax: `gitmap agy reread-optimize-project [N]` (short alias: `gitmap agy rop [N]`, default `N = 5`).
   - Filters registered AGY projects to identify the last N projects active/communicated within the last 24 hours. If fewer than N projects were active, operates safely on the available subset.
   - Cleans Antigravity cache using retention policy `--keep 10`.
   - Before purging or replacing conversations, creates an archival backup in a dedicated Split-DB at `data/agy/<repo-slug>/agy.db`.
   - Integrates with `gitmap storage reset` / `gitmap storage reset-errors` to allow clearing and resetting the AGY backup databases.

5. **AUM Go vs Python Search Benchmarking & Documentation**:
   - Side-by-side benchmarking comparing GitMap Go AUM search/find/list against legacy Python search scripts (`03-ai-scripts/12-fast-cached-grep.py`).
   - Documents benchmark results, execution speeds, and memory profiles in `readme.md`.
   - Live VM tests and heavy benchmarking tests are marked with `//go:build e2e` so GitHub Actions CI/CD runs smoothly without missing VM infrastructure.

---

## 3. Data Contracts & Models

### A. SSH Common Join Target (`cli/cmdssh/sshjoin_common_types.go`)

```go
type SSHCommonTarget struct {
    FullIP        string `json:"fullIp"`
    Alias         string `json:"alias"`
    Username      string `json:"username"`
    Port          int    `json:"port"`
    IsSuccess     bool   `json:"isSuccess"`
    DetectedOS    string `json:"detectedOs,omitempty"`
    OSVersion     string `json:"osVersion,omitempty"`
    OSArch        string `json:"osArch,omitempty"`
    ErrorMsg      string `json:"errorMsg,omitempty"`
}

type SSHCommonJoinResult struct {
    Username      string            `json:"username"`
    TotalCount    int               `json:"totalCount"`
    SuccessCount  int               `json:"successCount"`
    FailureCount  int               `json:"failureCount"`
    Targets       []SSHCommonTarget `json:"targets"`
}
```

### B. OS Information Report (`cli/cmdos/os_info_types.go`)

```go
type OSInfoReport struct {
    OSType       string `json:"osType"`
    OSVersion    string `json:"osVersion"`
    Architecture string `json:"architecture"`
    Hostname     string `json:"hostname"`
    NumCPU       int    `json:"numCpu"`
    Kernel       string `json:"kernel,omitempty"`
    Platform     string `json:"platform"`
}
```

### C. AGY ROP Backup Entry (`cli/cmdagy/agy_rop_types.go`)

```go
type AGYConversationBackup struct {
    ConversationID string `json:"conversationId"`
    ProjectSlug    string `json:"projectSlug"`
    ProjectPath    string `json:"projectPath"`
    StepCount      int    `json:"stepCount"`
    TranscriptJSON string `json:"transcriptJson"`
    CreatedAt      string `json:"createdAt"`
}
```

---

## 4. Verification & Quality Gates

1. **Zero Nested Ifs:** All conditionals strictly flattened with early returns and guard clauses.
2. **Positive Booleans:** Boolean names and flags prefixed with `is` or `has` (e.g. `isSuccess`, `hasPassFlag`).
3. **Structured AppError:** All errors wrapped with domain context via `apperror.WrapSimple` or `apperror.New`.
4. **Targeted Linters:** Must pass all 5 repo linters:
   - `python linter-scripts/check-nested-ifs.py`
   - `python linter-scripts/check-enum-and-boolean.py`
   - `python linter-scripts/check-error-management.py`
   - `python linter-scripts/check-relative-paths.py`
   - `python 03-ai-scripts/26-go-code-formatter.py`
5. **E2E Isolation:** Live VM tests and heavy benchmarks guarded with `//go:build e2e` so CI/CD does not fail on remote environments.
