# Subtask 01: Terminal UI Overhaul & Pre-Flight Status Dashboard

> **Parent Plan:** `67-nodes-cfr-remote-fleet-clone-enhancement.md`  
> **Status:** `PENDING`  
> **Target Files:**  
> - `cli/cmdnodes/nodes_clone_table.go`  
> - `cli/cmdnodes/nodes_clone_types.go`  
> - `cli/cmdnodes/nodes_clone_help.go`  
> - `cli/cmdnodes/nodes_clone_table_test.go`  

---

## 1. Technical Context & Scope

The current terminal output in `cli/cmdnodes/nodes_clone_table.go` uses an obsolete ASCII border banner that fails to provide observability during fleet operations. Operators are unable to see:
1. Total number of registered machines and their live reachability status prior to dispatch.
2. The specific remote destination directories where requests are being routed (`alias -> ip:dest`).
3. The remote GitMap version running on each target node.
4. Clean, colorized status badges and actionable command suggestions in the footer.

This subtask overhauls the terminal rendering engine to implement a high-clarity, modern dashboard with pre-flight probing summaries, dispatch route telemetry, aligned column tables, and actionable next-step suggestions.

---

## 2. Technical Specification & Component Contracts

### 2.1 Pre-Flight Readiness Dashboard
- **Function:** `renderFleetPreFlightTable(out io.Writer, report FleetPreFlightReport)`
- **Counters Line:**
  `  ▸ Fleet Readiness: %d registered node(s) | %d online | %d offline`
- **Columns (Widths):**
  - `NODE (ALIAS)`: 16 chars (Bold White)
  - `HOST`: 18 chars (Plain White)
  - `OS`: 10 chars (Cyan)
  - `VERSION`: 14 chars (White)
  - `DESTINATION`: 32 chars (Dim White / Path)
  - `STATUS`: 14 chars (Status badge)
- **Status Badges:**
  - Online: `constants.ColorGreen + "● online" + constants.ColorReset`
  - Offline: `constants.ColorYellow + "○ offline" + constants.ColorReset`
  - Auth Failed: `constants.ColorMagenta + "▲ auth_failed" + constants.ColorReset`
  - Unreachable: `constants.ColorRed + "✗ unreachable" + constants.ColorReset`

### 2.2 Modernized Fleet Start Banner
- **Function:** `renderFleetStartBanner(out io.Writer, opts NodesCloneOptions, report FleetPreFlightReport)`
- **Header Box:**
  Clean single-line border box using Unicode box characters (`┌───`, `│`, `└───`).
- **Metadata Lines:**
  - `• Mode:` Operation kind (`cfr`, `clone`, `cfrp`).
  - `• Target:` Target repository URL or staged manifest filename.
  - `• Workdir:` Active base directory with preserved relative path.
  - `• Scope:` Local master plus count of active online worker nodes.
  - `• Dispatch:` Destination routing mapping for all online targets:
    `  w1 -> 192.168.1.10:D:\work`
    `  w2 -> 192.168.1.11:~/work`

### 2.3 Overhauled Results Telemetry Table
- **Function:** `renderFleetResultsTable(out io.Writer, results []RemoteCloneNodeResult, isLocalSuccess bool, localDetails string, localDuration time.Duration, opts NodesCloneOptions)`
- **Columns:**
  - `NODE (ALIAS)` (16)
  - `HOST` (18)
  - `ROLE` (10)
  - `STATUS` (14)
  - `VERSION` (12)
  - `DURATION` (12)
  - `DETAILS` (Rest)
- **Row Formats:**
  - Local host row marked with `local (current)` and role `master`.
  - Remote rows dynamically displaying remote GitMap version and millisecond durations.
  - Aligned using `padVisual()` with ANSI escape code stripping (`stripANSI`).

### 2.4 Actionable Next-Step Footer Suggestions
- **Function:** `renderFleetFooterSuggestions(out io.Writer, opts NodesCloneOptions)`
- Renders standard GitMap `[tip]` footer block:
  ```text
  [tip] Fleet Operations & Suggested Commands:
    • Ping Fleet Nodes:          gitmap nodes ping
    • Inspect Node Connection:   gitmap ssh test <alias>
    • Query Machine Telemetry:   gitmap machine --ssh
    • Rerun Except Local Host:   gitmap nodes cfr <repo> --except-self
  ```

---

## 3. Function Decomposition Plan (<= 15 Lines Per Function)

To strictly adhere to repository coding standards, all rendering routines are decomposed into modular functions:

1. `renderFleetPreFlightTable(out io.Writer, report FleetPreFlightReport)` (<= 10 lines)
2. `renderPreFlightCounters(out io.Writer, total, online, offline int)` (<= 8 lines)
3. `renderPreFlightHeader(out io.Writer)` (<= 6 lines)
4. `renderPreFlightRow(out io.Writer, node NodePreFlightInfo)` (<= 12 lines)
5. `resolvePreFlightBadge(isOnline bool, errStr string)` (<= 12 lines)
6. `renderFleetStartBanner(out io.Writer, opts NodesCloneOptions, report FleetPreFlightReport)` (<= 12 lines)
7. `renderBannerHeaderBox(out io.Writer, kind NodesCloneKind)` (<= 8 lines)
8. `renderBannerMetadata(out io.Writer, opts NodesCloneOptions, onlineCount int)` (<= 12 lines)
9. `renderBannerDispatchRoutes(out io.Writer, nodes []NodePreFlightInfo)` (<= 14 lines)
10. `renderFleetResultsTable(...)` (<= 12 lines)
11. `renderResultsTableHeader(out io.Writer)` (<= 6 lines)
12. `renderLocalResultRow(out io.Writer, isLocalSuccess, isSkipLocal bool, details string, dur time.Duration)` (<= 14 lines)
13. `renderRemoteResultRow(out io.Writer, r RemoteCloneNodeResult)` (<= 14 lines)
14. `renderFleetSummaryFooter(out io.Writer, results []RemoteCloneNodeResult, isLocalSuccess, isSkipLocal bool)` (<= 12 lines)
15. `renderFleetFooterSuggestions(out io.Writer, opts NodesCloneOptions)` (<= 12 lines)

---

## 4. Coding Guidelines & Invariant Rules

- **Positive Booleans Only:** Use `isOnline`, `isLocalSuccess`, `isSkipLocal`, `hasFile`. Avoid `isOffline` or `notSuccess`.
- **Zero Raw Magic Strings:** Reference color and formatting constants from `cli/constants`.
- **ANSI Width Safety:** Never compute raw byte or rune lengths on colorized strings; always use `visualLen()` / `padVisual()`.
- **No Git Commands:** Implementation and verification must strictly avoid running any git commands.

---

## 5. Verification & Testing Protocol

1. **Unit Test Coverage (`cli/cmdnodes/nodes_clone_table_test.go`):**
   - Test `renderFleetPreFlightTable` with empty, single-node, and multi-node reports.
   - Test visual alignment with varying alias and path lengths.
   - Verify `resolvePreFlightBadge` returns expected ANSI colors for `isOnline = true` vs `isOnline = false`.
   - Verify `renderFleetFooterSuggestions` outputs proper command tips.
2. **Quality Checks:**
   - Verify function line count does not exceed 15 lines.
   - Verify absence of compiler warnings or lint errors.
