# Learned Memory 34: Direct Clone Structured JSON Output, Fleet Remote Fallback, and Offline Reachability Telemetry

Date: 2026-10-01
Related Specs:
- [02-spec/21-app/192-direct-clone-json-and-fleet-clone-details.md](../../../02-spec/21-app/192-direct-clone-json-and-fleet-clone-details.md)
- [02-spec/21-app/191-nodes-clone-except-self-windows-runner-and-path.md](../../../02-spec/21-app/191-nodes-clone-except-self-windows-runner-and-path.md)

## 1. Architectural Lessons & Principles

### 1.1 Structured JSON Output for Subordinate CLI Operations
- Subordinate and delegated CLI commands (e.g. `gitmap clone <url> --json`, `gitmap cfr <url> --json`) must serialize outcomes into a deterministic schema (`DirectCloneJSONResponse`) containing `success`, `repoName`, `status`, `message`, and `path`.
- When parent orchestrators (like `gitmap nodes clone`) invoke remote processes over SSH or in subshells, parsing structured JSON avoids fragile regex scraping of multi-line terminal banners.

### 1.2 Graceful Heterogeneous Version Fallback
- Fleets may contain nodes running different minor versions of GitMap.
- When an orchestrator appends flags (like `--json`), it must handle `flag provided but not defined` errors gracefully by automatically retrying with legacy flagless commands (`retryWithoutJSON`) instead of failing the fleet task.
- Provide legacy output extractors (`extractLegacyCloneDetails`) to harvest high-signal details from unstructured stdout when running against older nodes.

### 1.3 Precise Offline Error Distinction
- Network timeouts, DNS resolution failures, and connection refused errors must not be reported as generic authentication errors.
- Comprehensive pattern matching against OS socket errors (`connectex`, `actively refused`, `no route to host`, `i/o timeout`, `getaddrinfo`) ensures nodes are flagged with `offline` status and user-friendly diagnostics.

### 1.4 In-Memory Local Stream Capture
- When executing local steps in-process during fleet runs, redirecting and capturing stdout/stderr into an in-memory pipe allows direct parsing of the local JSON envelope without file I/O or terminal flickering, surfacing granular messages in the `local (current)` table row.
