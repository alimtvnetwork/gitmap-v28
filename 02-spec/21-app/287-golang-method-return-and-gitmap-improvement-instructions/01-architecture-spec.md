# Specification 287: Golang Method Return Types, Zero Swallowed Errors & GitMap Improvement Architecture

## 1. Executive Overview & Problem Statement

### 1.1 The Void Method & Error Swallowing Defect
In production Go codebases, particularly CLI tools and database orchestrators, fallible functions performing mutations, database writes, or external subprocess execution must never be declared with a `void` return signature (`func DoSomething(...)`). Declaring fallible functions without an `error` return forces either:
1. Panic on failure (prohibited).
2. Exit handler invocation inside leaf functions (violating Leaf Error Sovereignty).
3. **Silent error swallowing** (discarding errors with `_ = ...` or empty `if err != nil { return }`).

As identified in `cli/cmdpull/pull_remediation_ssh.go`, `updateRepoTransportInDB` and `remediateCredentialAndConvertToSSH` were declared as void functions, resulting in swallowed database update failures and silent subprocess errors:
```go
// VIOLATION: Void return with swallowed database error
func updateRepoTransportInDB(repoPath, sshURL string) {
    db, err := store.OpenDefault()
    if err != nil {
        return // SWALLOWED ERROR: caller has no awareness of DB failure
    }
    defer db.Close()

    _ = db.UpdateRepoTransport(repoPath, sshURL, constants.TransportSSH) // SWALLOWED ERROR
}
```

### 1.2 The Remediation Architectural Standard
All mutating, I/O, subprocess, and database functions must strictly adhere to the following signature contract:
```go
// REMEDIATED: Explicit error return with leaf error sovereignty
func updateRepoTransportInDB(repoPath, sshURL string) error {
    db, err := store.OpenDefault()
    if err != nil {
        return err
    }
    defer db.Close()

    return db.UpdateRepoTransport(repoPath, sshURL, constants.TransportSSH)
}
```
Callers must capture the returned error and either bubble it up (`return appfault.Wrap(...)`), handle it with graceful user telemetry, or log an explicit warning rather than silently discarding it.

---

## 2. GitMap Systems Architecture & Improvement Horizons

GitMap is an autonomous developer companion and CLI for ultra-fast repository scanning, polyglot automation, cluster/SSH delegation, pipeline self-healing, and coding guideline enforcement. To support future scale and autonomous multi-agent execution, this specification sets forth:
1. **200 Architectural & Functional Improvements** organized across 10 core subsystems.
2. **200 Sequential Implementation Steps** detailing execution actions.
3. **100 Terminal UI & User Experience Enhancements** spanning Lipgloss styling, interactive selectors, tables, error HUDs, and web visualizers.

---

## 3. Subsystem Breakdown for Improvements

| Subsystem | Area | Core Files |
|---|---|---|
| **Subsystem 01** | High-Speed Scanner & Parallel Cloner | `cli/scanpipe/`, `cli/cloner/`, `cli/clonenow/` |
| **Subsystem 02** | Multi-Node SSH Fleet Orchestration | `cli/cmdssh/`, `cli/cmdnodes/`, `cli/cluster/` |
| **Subsystem 03** | Three-Tier Split-DB Architecture | `cli/store/`, `cli/cmddb/`, `cli/dbengine/` |
| **Subsystem 04** | CI/CD Pipeline Telemetry & RCA Engine | `cli/cmdpe/`, `cli/cmdpipeline/`, `cli/cmdpullerror/` |
| **Subsystem 05** | Credential Management & DPAPI Resilience | `cli/cmdfixauth/`, `cli/gitutil/` |
| **Subsystem 06** | Autonomous AI Task Engine & Heatmap | `cli/cmdagent/`, `cli/heatmap/` |
| **Subsystem 07** | Macro Automation & Terminal UX | `cli/macro/`, `cli/cmdmacro/`, `cli/termout/` |
| **Subsystem 08** | Security, Vault & Data Loss Prevention | `cli/secrets/`, `cli/power/` |
| **Subsystem 09** | Polyglot Cache Cleaning & Temp Storage | `cli/osclean/`, `cli/cmdhygiene/` |
| **Subsystem 10** | Release Ceremony & Version Orchestration | `cli/release/`, `cli/cmdversionhistory/` |

---

## 4. Verification & Quality Gates

- **Positive Booleans**: Strict `is*` and `has*` prefixes.
- **Function Limit**: <= 15 lines max per function (<= 8 lines preferred).
- **Zero Nested Ifs**: Flattened conditionals using guard clauses.
- **FastGate AST Linter**: 100% pass across all modified files.
- **Targeted Unit Tests**: Full pass in `cmdpull`, `fsutil`, `heatmap`, `store`.
