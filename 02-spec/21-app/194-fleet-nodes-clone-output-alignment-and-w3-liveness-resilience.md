# 194 — Fleet Nodes Clone Output Alignment & W3 Liveness Resilience

## Status: Active
- **Spec ID:** 194
- **Scope:** Fleet Orchestration, Table Formatting, SSH Liveness Caching, Credential Fallback
- **Created At:** 2026-10-01

---

## 1. User Request (Verbatim)

```text
Output is still wrong for

gitmap nodes clone https://github.com/alimtvnetwork/awansoft-v10


I can still acess to w3 machine fix it and release minor please and
```

---

## 2. Root Cause & Architectural Gaps

1. **Terminal Stderr Leakage**: In `executeLocalClone`, only `os.Stdout` was redirected to an internal buffer. Any stderr output (such as `↑ cfr: cwd is a git repo — escaping to non-repo ancestor` from `escapeNestedGitRepo`) leaked directly to standard error between the dispatch banner and table.
2. **Table Alignment Failure via ANSI Byte Measurement**: Go's `fmt.Fprintf` format verbs measure byte slice length rather than terminal visual column width. ANSI color sequences (`\x1b[32m...\x1b[0m`) consume 9–10 bytes, causing the `%-20s` column formatting to under-pad colored statuses (`● success`, `○ offline`), shifting the `DURATION` column 5–9 characters leftward and distorting table columns.
3. **Hyper-Aggressive 1500ms Probe & 45s Dead Caching**: `CheckConnLiveness` used a rigid 1.5s timeout. High-latency VM networks or temporarily busy Windows nodes (such as `w3`) that drop a single SYN packet are instantly marked offline, and that failure was cached for 45 seconds across all commands.
4. **Local Database Isolation**: `queryHostPasswordFromDB` queried only `store.OpenDefault()`. When executed from the repo root containing an un-enrolled local `data/gitmap.db`, it failed to inspect `store.OpenGlobalDefault()`, causing password fallback failure.
5. **13.5s Timeout Cascade**: Remote node worker routines sequentially attempted key dials, multiple default keys, and password dials without an initial fast-fail probe, delaying unreachable node reporting by 13.5 seconds.

---

## 3. Data Contracts & Visual Specification

### 3.1 Fleet Nodes Clone Clean Table Layout
```text
  ╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
  ║ GITMAP FLEET NODES CLONE DISPATCH                                                                                ║
  ╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
  ▸ Dispatching 'clone' local host and 5 remote fleet node(s)...

  NODE (ALIAS)     HOST                   ROLE       STATUS        DURATION   DETAILS
  --------------------------------------------------------------------------------------------------------------
  local (current)  127.0.0.1              master     ● success     in-process already exists on disk (./awansoft-v10)
  w1               node-w1            worker     ● success     4120ms     already exists on disk
  w4               node-w4           worker     ● success     4280ms     already exists on disk
  w3               node-w3           worker     ○ offline     3010ms     node offline or unreachable
  main             node-main           worker     ○ offline     3015ms     node offline or unreachable
  u1               node-u1           worker     ○ offline     3020ms     node offline or unreachable
  --------------------------------------------------------------------------------------------------------------

  ✔ Fleet Clone Summary: 3/6 node(s) completed successfully (3 failed)
```

### 3.2 Liveness Probe Parameters
- Default Probe Timeout: 3000ms.
- Retry: 1 retry before marking offline.
- Failure Cache TTL: 5s (prevents stale locking).
- Success Cache TTL: 45s.

---

## 4. Acceptance Criteria

- **AC-01:** `gitmap nodes clone <repo>` produces zero leaked text above the results table.
- **AC-02:** Table columns (`NODE (ALIAS)`, `HOST`, `ROLE`, `STATUS`, `DURATION`, `DETAILS`) align strictly to column boundaries regardless of ANSI escapes.
- **AC-03:** `queryHostPasswordFromDB` queries `OpenGlobalDefault()` if `OpenDefault()` contains no match.
- **AC-04:** Liveness check uses a 3s timeout with 1 retry and 5s failure TTL.
- **AC-05:** Release minor version bump is verified without build or test execution in adherence to R1.
