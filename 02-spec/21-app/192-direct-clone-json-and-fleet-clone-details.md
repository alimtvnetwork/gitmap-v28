# Specification 192: Direct Clone JSON Output & Fleet Clone Node Details Architecture

**Status:** Active
**Author:** Lead Architect
**Category:** Fleet Management & Structured Outputs
**Reference:** [02-spec/21-app/readme.md](readme.md)
**Related Specs:**
- [02-spec/21-app/191-nodes-clone-except-self-windows-runner-and-path.md](191-nodes-clone-except-self-windows-runner-and-path.md)
- [02-spec/21-app/188-json-envelope-v2-terminal-clear-and-deploy-keys.md](188-json-envelope-v2-terminal-clear-and-deploy-keys.md)

---

## 1. Domain Context & Architecture Overview

`gitmap nodes clone`, `gitmap nodes cfr`, and `gitmap nodes cfrp` dispatch parallel repository cloning operations across local and remote nodes over SSH. In prior versions, remote commands emitted standard unstructured terminal text, making parsing of per-node outcomes noisy, fragile, and prone to formatting misalignment.

### 1.1 Problems Addressed
1. **Unstructured Output in Remote Delegations:**
   When remote nodes ran `gitmap clone` / `cfr` / `cfrp`, the parent process received unstructured text logs. Errors, existing repository detections, and success messages had to be parsed heuristically or displayed as truncated raw text.
2. **Missing Granular Node Details Column:**
   In terminal table rendering, remote node rows often displayed generic `ok` or full unformatted stderr messages that degraded column formatting.
3. **Legacy Node Compatibility:**
   Enrolled remote nodes might run older GitMap binaries that do not support the `--json` flag on direct clone verbs. Unconditional `--json` passing on older nodes would fail with `flag provided but not defined: --json`.
4. **Offline Node Misclassification:**
   Host connection timeouts, DNS resolution failures, and connection refusals were sometimes conflated with SSH authentication rejections rather than cleanly classified as `offline`.

---

## 2. Technical Specifications & Architecture Blueprint

### 2.1 Direct Clone Structured JSON Protocol (`DirectCloneJSONResponse`)
All clone subcommands (`RunClone`, `RunCloneFixRepo`, `RunCloneFixRepoPub`) now support `--json` / `-j` flag:
- When `--json` is supplied, standard banner printing is suppressed.
- Outcomes are serialized into a single machine-readable JSON object:
  ```json
  {
    "success": true,
    "repoName": "repo-name",
    "status": "success",
    "message": "already exists on disk",
    "path": "D:\\work\\repo-name"
  }
  ```
- Parser `cmdclone.ParseCloneJSONResponse(out)` safely extracts the trailing JSON envelope even if preceding logs or hooks emitted warnings.

### 2.2 Remote Fleet Execution & Automatic Fallback
- `buildWindowsWorkDirExecString` and `buildUnixWorkDirExecString` automatically append `--json` to the delegated command string.
- If the remote execution returns `flag provided but not defined: -json` or `--json`, `retryWithoutJSON` automatically retries execution with legacy command formatting, preserving backward compatibility across heterogeneous node versions.
- If output contains structured JSON, `RemoteCloneNodeResult.Details` is populated directly with `payload.Message` and `Status` is mapped cleanly.
- For legacy nodes, `extractLegacyCloneDetails` extracts concise summaries (`already exists on disk`, `cloned successfully`) from raw stdout.

### 2.3 Extended Offline Error Heuristics
`isOfflineError(errStr)` classifies node reachability accurately against:
- Network unreachable
- Connection refused / actively refused
- Timeout / i/o timeout
- ConnectEx
- Machine is off / host is down / no route to host
- Name resolution / getaddrinfo failures

### 2.4 Terminal Results Table Formatting
- The table `DETAILS` column surfaces `r.Details` first, falling back to sanitized error strings.
- Local host execution captures stdout/stderr in-memory, parses direct JSON response, and populates `localDetails` cleanly in the `local (current)` row.

---

## 3. Data Contracts & Model Extensions

```go
// DirectCloneJSONResponse captures structured clone outcome for JSON delegates.
type DirectCloneJSONResponse struct {
	Success  bool   `json:"success"`
	RepoName string `json:"repoName"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	Path     string `json:"path"`
}

// RemoteCloneNodeResult adds Details for high-precision table rendering.
type RemoteCloneNodeResult struct {
	Alias      string        `json:"alias"`
	Host       string        `json:"host"`
	Role       string        `json:"role"`
	Status     string        `json:"status"`
	Duration   time.Duration `json:"duration"`
	DurationMs int64         `json:"durationMs"`
	Stdout     string        `json:"stdout,omitempty"`
	Stderr     string        `json:"stderr,omitempty"`
	Error      string        `json:"error,omitempty"`
	Details    string        `json:"details,omitempty"`
}
```

---

## 4. Acceptance Criteria & Quality Gates

1. **AC-DC-001 (Direct Clone JSON Output):** `gitmap clone <url> --json` outputs a valid single-line JSON string conforming to `DirectCloneJSONResponse` with exit code 0.
2. **AC-DC-002 (Fleet Remote JSON Parsing):** `gitmap nodes clone <url>` parses JSON responses from remote nodes and surfaces exact messages in the `DETAILS` column.
3. **AC-DC-003 (Legacy Flag Fallback):** Remote nodes running older versions lacking `--json` automatically fall back via `retryWithoutJSON` and complete without halting the fleet runner.
4. **AC-DC-004 (Offline Error Classification):** Network timeouts and unreachable connections are classified as `offline` with status message `node offline or unreachable`.
5. **AC-DC-005 (Local Execution Details):** The `local (current)` table row displays explicit outcome messages rather than generic fallback strings.
