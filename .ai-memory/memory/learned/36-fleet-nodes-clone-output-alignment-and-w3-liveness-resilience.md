# Learned Memory 36: Fleet Nodes Clone Output Alignment, Stderr Escape Suppression, and Resilient Liveness Probing

Date: 2026-10-01
Related Specs:
- [02-spec/21-app/194-fleet-nodes-clone-output-alignment-and-w3-liveness-resilience.md](../../../02-spec/21-app/194-fleet-nodes-clone-output-alignment-and-w3-liveness-resilience.md)
- [02-spec/22-app-issues/57-fleet-nodes-clone-misalignment-and-w3-liveness-rca.md](../../../02-spec/22-app-issues/57-fleet-nodes-clone-misalignment-and-w3-liveness-rca.md)
- [.ai-memory/plans/completed/57-fleet-nodes-clone-output-alignment-and-w3-liveness-resilience.md](../../plans/completed/57-fleet-nodes-clone-output-alignment-and-w3-liveness-resilience.md)

## 1. Architectural Lessons & Principles

### 1.1 ANSI-Aware Terminal Visual Width Alignment
Standard Go string formatters (`fmt.Sprintf("%-20s", val)`) measure raw byte count. ANSI escape sequences (`\x1b[32m...\x1b[0m`) add non-printing bytes (e.g. 9 bytes), leading to under-padding when calculating column widths.
- Implement an ANSI-stripping visible width calculator (`visibleWidth(s)`) that counts only visual runes.
- Use explicit visual right-padding (`padRight(s, width)`) across all table rendering in `cmdnodes` to maintain strict column alignment regardless of color escapes.

### 1.2 Dual Stream Capture & Information Leak Suppression
When executing nested operations (e.g. `executeLocalClone` inside `gitmap nodes clone`):
- Both `os.Stdout` and `os.Stderr` must be intercepted into dedicated buffers to prevent informative messages (such as `↑ cfr: cwd is a git repo`) from leaking into stdout above structured tables.
- Suppress nested repo escape notifications when fleet mode (`IsFleetCloneActive()`) or structured JSON output is active.

### 1.3 Asymmetric Liveness Cache TTLs & Dynamic Retry
- Fast-fail liveness probes on busy nodes (e.g. Windows Server VMs on VMware) can fail transiently under load.
- Implement a 3000ms probe with 1 automatic retry on timeout.
- Decouple cache TTLs: 45 seconds for successful probes, but only 5 seconds for failed probes so recovered nodes are not locked out for extended periods.

### 1.4 Hierarchical SQLite Credential Fallback
- Local databases (`store.OpenDefault()`) may be un-enrolled when invoked from repository workspaces.
- Always fall back to the global installation vault (`store.OpenGlobalDefault()`) when resolving stored SSH node credentials.
