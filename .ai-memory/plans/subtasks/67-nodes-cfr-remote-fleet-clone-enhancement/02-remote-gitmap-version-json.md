# Subtask 02: Remote GitMap Version Probing via JSON Protocol

> **Parent Plan:** `67-nodes-cfr-remote-fleet-clone-enhancement.md`  
> **Status:** `PENDING`  
> **Target Files:**  
> - `cli/cmd/rootutility.go`  
> - `cli/cmd/rootutility_test.go`  
> - `cli/cmdnodes/nodes_clone_remote.go`  
> - `cli/cmdnodes/nodes_clone_types.go`  
> - `cli/cmdnodes/nodes_clone_test.go`  

---

## 1. Technical Context & Scope

When coordinating fleet-wide operations, operators have no prior indication of what versions of `gitmap` are installed across remote SSH machines. Version mismatches can trigger silent incompatibilities, unrecognized command flags, or parsing failures. 

Currently:
1. `gitmap version` in `cli/cmd/rootutility.go` only emits human-readable text (`gitmap v6.453.0\n`) and does not support structured JSON flags.
2. The remote fleet runner in `cli/cmdnodes/nodes_clone_remote.go` connects without pre-flight version validation.
3. Unreachable or slow-to-connect nodes cause execution delays rather than failing fast in a dedicated pre-flight phase.

This subtask implements structured `--json` emission for `gitmap version`, along with a concurrent, 3-second bounded SSH probing pipeline in `cmdnodes` to inspect and report remote GitMap versions before any clone task is dispatched.

---

## 2. Technical Specification & Implementation Details

### 2.1 Structured Local Version Emission (`cli/cmd/rootutility.go`)

- **Command Syntax:** `gitmap version --json`, `gitmap version -j`, `gitmap -v --json`.
- **Payload Data Structure:**
  ```go
  type GitMapVersionPayload struct {
      Version   string `json:"version"`
      Commit    string `json:"commit,omitempty"`
      OS        string `json:"os"`
      Arch      string `json:"arch"`
      BuildDate string `json:"buildDate,omitempty"`
  }
  ```
- **Execution Logic:**
  - In `printVersionBlock()`, check whether arguments contain `--json` or `-j` via `isJSONVersionRequest(args)`.
  - If requested, assemble `GitMapVersionPayload` using:
    - Version: `constants.Version`
    - Commit: `resolveGitmapCommitSHA()`
    - OS: `runtime.GOOS`
    - Arch: `runtime.GOARCH`
    - BuildDate: resolved build timestamp
  - Serialize to standard output with 2-space indentation and trailing newline.

### 2.2 Remote SSH Probing Pipeline (`cli/cmdnodes/nodes_clone_remote.go`)

- **Function:** `probeRemoteNodeVersion(conn db.SSHConnection, timeout time.Duration) NodePreFlightInfo`
- **Context & Timeout:**
  - Create bounded context: `ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)`.
  - Connect client via SSH using `cmdssh.ConnectSSHClientWithErr(conn)`.
- **Remote Execution:**
  - Execute remote command `gitmap version --json`.
  - On Windows targets: `gitmap version --json` via PowerShell session.
  - On Unix/Linux targets: `gitmap version --json` via Bash/sh session.
- **Output Parsing & Fallback:**
  - Attempt to parse JSON into `GitMapVersionPayload`.
  - If remote returns non-JSON (older GitMap binary without `--json`), fallback to scraping the version string using SemVer regular expression `v[0-9]+\.[0-9]+\.[0-9]+`.
  - If remote binary is missing, report `"gitmap not found in PATH"`.
- **Connection Error Classification:**
  - Offline/Unreachable: `isOfflineError(err.Error())` -> status `"offline"`.
  - Auth failure: `isAuthError(err.Error())` -> status `"auth_failed"`.

### 2.3 Fleet-Wide Probing Aggregator

- **Function:** `probeFleetNodesVersion(conns []db.SSHConnection, timeout time.Duration, targetDir string) FleetPreFlightReport`
- **Concurrency:**
  - Probe all candidate connections concurrently via goroutines and `sync.WaitGroup`.
  - Synchronize updates to the results slice using `sync.Mutex`.
- **Output Aggregation:**
  - Compute `TotalCount`, `OnlineCount`, and `OfflineCount`.
  - Associate resolved target destination directory for each node based on OS group.

---

## 3. Function Decomposition Plan (<= 15 Lines Per Function)

1. `isJSONVersionRequest(args []string) bool` (<= 8 lines)
2. `buildVersionPayload() GitMapVersionPayload` (<= 10 lines)
3. `emitVersionJSON(out io.Writer, payload GitMapVersionPayload) error` (<= 8 lines)
4. `printVersionBlock() error` (<= 12 lines)
5. `probeRemoteNodeVersion(conn db.SSHConnection, timeout time.Duration) NodePreFlightInfo` (<= 14 lines)
6. `executeRemoteVersionCommand(client *ssh.Client, isWin bool) (string, error)` (<= 10 lines)
7. `parseRemoteVersionPayload(raw string) (GitMapVersionPayload, bool)` (<= 12 lines)
8. `extractFallbackVersionString(raw string) string` (<= 10 lines)
9. `classifyProbeError(err error, dur time.Duration) NodePreFlightInfo` (<= 14 lines)
10. `probeFleetNodesVersion(conns []db.SSHConnection, timeout time.Duration, dir string) FleetPreFlightReport` (<= 14 lines)
11. `dispatchProbeWorker(wg *sync.WaitGroup, mu *sync.Mutex, conn db.SSHConnection, timeout time.Duration, report *FleetPreFlightReport)` (<= 12 lines)
12. `calculatePreFlightCounts(report *FleetPreFlightReport)` (<= 10 lines)

---

## 4. Coding Guidelines & Invariant Rules

- **Positive Boolean Fields:** Use `isOnline`, `hasValidJSON`, `hasJSONFlag`, `isWin`.
- **Error Wrapping:** Wrap system errors with `apperror.WrapSimple(err, "probe remote node version")`.
- **Context Discipline:** Every network and SSH operation must respect context cancellation and the 3-second ceiling.
- **Zero-Git Command Rule:** Do not invoke git commands during coding or testing.

---

## 5. Verification & Testing Protocol

1. **Unit Tests for Local JSON Version (`cli/cmd/rootutility_test.go`):**
   - Test `gitmap version --json` emits valid JSON containing `version`, `os`, and `arch`.
   - Test `gitmap version` without flags continues to emit standard human-readable text.
2. **Unit Tests for Probing & Fallback (`cli/cmdnodes/nodes_clone_test.go`):**
   - Test `parseRemoteVersionPayload` with valid JSON string.
   - Test `extractFallbackVersionString` with legacy `gitmap v6.450.0` output.
   - Test `classifyProbeError` with mock timeout, connection refused, and auth rejection errors.
3. **Quality Gate:**
   - Confirm all functions are 15 lines or fewer.
   - Ensure zero compiler warnings and zero linter regressions.
