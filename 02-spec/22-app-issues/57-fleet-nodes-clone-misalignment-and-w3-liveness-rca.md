# 57 — Fleet Nodes Clone Table Misalignment, Stderr Escape Leaks, and Remote Liveness False-Offline RCA

## Status: Resolved
- **Issue ID:** RCA-57
- **Target Subsystems:** `cli/cmdnodes`, `cli/cmdclone`, `cli/cmdssh`
- **Date:** 2026-10-01

---

## 1. Reproduction & Symptoms

### 1.1 Leaked Stderr Text Between Banner & Table
Executing `gitmap nodes clone https://github.com/alimtvnetwork/awansoft-v10` emitted:
```text
  ╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
  ║ GITMAP FLEET NODES CLONE DISPATCH                                                                                ║
  ╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
  ▸ Dispatching 'clone' local host and 5 remote fleet node(s)...

  ↑ cfr: cwd is a git repo — escaping to non-repo ancestor
    from: <workspace>\gitmap
      to: D:\work
  NODE (ALIAS)     HOST                   ROLE       STATUS        DURATION   DETAILS
```
The escape notification `↑ cfr: cwd is a git repo` leaked to `os.Stderr` because `executeLocalClone` only captured `os.Stdout`.

### 1.2 Table Column Misalignment
ANSI color escapes in `statusTag` (`\x1b[32m● success\x1b[0m`) occupy 9 bytes that do not consume visual terminal columns. Go's standard `fmt.Fprintf` calculates width based on byte slice length, causing `%-20s` to add zero padding spaces. Consequently, the `DURATION` column shifted left by 5 to 9 columns, resulting in broken table boundaries.

### 1.3 False Offline & 13.5s Timeout Cascade on Remote Machines (e.g. `w3`)
1. In `cmdssh/ssh_liveness.go`, `CheckConnLiveness` used a 1500ms timeout. Any network jitter or VM latency triggered an offline state, which was cached for 45 seconds.
2. In `cmdssh/ssh_dial_fallback.go`, unreachable nodes sequentially dialed specified keys (4s), default keys (8s), and passwords (4s) before checking liveness, taking 13.5 seconds.
3. In `cmdssh/sshexec.go`, `queryHostPasswordFromDB` queried only `store.OpenDefault()`. When executed with an un-enrolled local database, it never consulted `store.OpenGlobalDefault()`.

---

## 2. Root Cause Analysis (4-Part RCA)

### 2.1 Direct Cause
- `executeLocalClone` captured only stdout, permitting stderr messages to escape.
- `renderRemoteRow` relied on Go's byte-oriented `fmt.Fprintf` formatting rather than visual character width.
- `CheckConnLiveness` used a sub-2s timeout without retry, caching failure for 45 seconds.

### 2.2 Indirect Cause
- Inconsistent database fallback: `fetchAllSSHConnections` checked global fallback, but `queryHostPasswordFromDB` did not.
- Dial fallback lacked a pre-flight connectivity probe, forcing 4 successive SSH timeout iterations.

### 2.3 Environmental Context
- Fleet contains mixed Windows Server 2022 virtual machines (`w1`, `w2`, `w3`, `w4`, `main`) running on VMware Workstation, where occasional SYN latency exceeds 1.5s.

### 2.4 Systemic Vulnerability
- Absence of ANSI-aware terminal string width measurement in table rendering functions across `cmdnodes`.

---

## 3. Corrective & Preventive Actions

1. **Dual Stream Capture**: Redirect both `os.Stdout` and `os.Stderr` in `executeLocalClone`.
2. **Suppress Informational Stderr in Fleet Mode**: Suppress `MsgCFREscapeNested` when `cmdnodes.IsFleetCloneActive()` or JSON mode is true.
3. **ANSI-Aware Visual Padding**: Implement `visibleWidth` and `padRight` to guarantee table alignment irrespective of color escapes.
4. **Resilient Liveness Probing**: Increase timeout to 3000ms, add 1 retry on timeout, and separate cache TTLs (5s for failures, 45s for successes).
5. **Global DB Credential Fallback**: Fall back to `store.OpenGlobalDefault()` in `queryHostPasswordFromDB`.
