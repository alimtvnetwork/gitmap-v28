# Master Plan 366: Golang Script Creation and Analysis

- **Sequence ID:** 366
- **Task Slug:** `golang-script-creation-and-analysis`
- **Root Slug:** `golang-script-creation-and-analysis`
- **Specification Reference:** [02-spec/21-app/366-golang-script-creation-and-analysis/readme.md](../../../02-spec/21-app/366-golang-script-creation-and-analysis/readme.md)
- **Status:** `pending`
- **Created At:** 2026-10-31
- **Budget:** 400

---

## 1. Executive Summary & Root Cause Analysis

### 1.1 Problem Statement & Background
During autonomous operations in the Antigravity IDE, Cursor, and web-based visualizers, AI agents frequently require deep runtime introspection of user interface elements, active DOM inputs, contenteditable elements, and JavaScript runtime states. Because the Chrome DevTools Protocol (CDP) provides direct WebSocket endpoints for debugging (`/devtools/page/<target_id>`), autonomous agents have repeatedly authored ad-hoc 70-line Python scripts utilizing raw low-level sockets (`socket.socket()`), manual HTTP/1.1 Upgrade handshakes (`Sec-WebSocket-Key`), custom RFC 6455 frame framing (`send_frame` / `recv_frame` via `struct.pack`/`unpack`), and hardcoded target ports and IDs.

### 1.2 Root Cause Analysis (4-Part RCA)
1. **Tooling Gap:** GitMap lacked a first-class, native CLI command to inspect, probe, and evaluate JavaScript expressions over Chrome DevTools Protocol (CDP) WebSocket sessions.
2. **Ephemeral Hardcoding Failure:** Ad-hoc Python scripts hardcoded ephemeral debug ports (e.g. `49994`, `9222`) and ephemeral target IDs (e.g. `3658C5BD73EEBD1D6A4E52160BBEAFFA`). Any browser restart, tab reload, or IDE crash immediately invalidated these scripts, causing agent failures.
3. **Fragile Framing & Payload Truncation:** The ad-hoc Python framing logic utilized unbuffered `recv(4096)` calls that fail on JSON response payloads larger than 64KB (extended payload lengths `126` and `127`). Furthermore, it lacked connection timeouts, ping/pong heartbeat handling, and TLS/WSS fallback.
4. **Transcript Bloat & Guideline Violations:** Injecting 70+ lines of raw Python socket code repeatedly into conversation turns consumed substantial token budget, generated terminal noise, and violated repository coding standards against unmanaged ad-hoc automation scripts.

### 1.3 Architectural Solution & Native Go Benefits
Plan 366 replaces all ad-hoc Python CDP scripting with a high-performance, native Go implementation packaged directly into the GitMap binary:
- **Pure Go Standard Library WebSocket Client (`cli/cmdcdp/`):** Implements an RFC 6455 compliant WebSocket client using pure Go standard library packages (`net`, `crypto/rand`, `encoding/binary`, `encoding/json`), maintaining zero external dependencies.
- **Automatic Port & Target Discovery:** Automatically probes localhost debug ports (`9222`, `49994`, `9229`) or queries `http://127.0.0.1:<port>/json/list` to dynamically discover active page targets and WebSocket endpoints.
- **First-Class CLI Subcommands:** Exposes `gitmap cdp probe`, `gitmap cdp eval`, and `gitmap cdp targets` with rich Tokyo Night Lipgloss card visualizers for human operators and structured `--json` output for AI agents.
- **Architectural Policy Ban:** Enforces a strict ban in `.agents/skills/gitmap/skill.md` prohibiting AI agents from authoring ad-hoc Python websocket scripts for CDP introspection.

---

## 2. Architectural Invariants & Coding Standards

1. **Leaf Error Sovereignty:** All leaf functions in `cli/cmdcdp/` returning `error` or `*apperror.AppError` must bubble errors directly to outer callers (`return err`). Never invoke exit handlers (`cliexit.HandleError`) or `os.Exit` inside leaf functions.
2. **Positive Boolean Naming:** All booleans must strictly start with `is*` or `has*` (e.g., `hasTarget`, `isMasked`, `hasExtendedPayload`, `isJSONOutput`, `hasFoundInput`). Negative prefixes or banned prefixes (`can*`, `should*`, `was*`, `will*`) are strictly prohibited.
3. **Function Sizing & Zero Nesting:** Functions must remain between 8 and 15 lines max. Guard clauses and early returns must flatten all conditional branches.
4. **Whitespace & Return Formatting:** Exactly one blank line after closing brace `}` and one blank line before `return`.
5. **Relative Git Paths:** Strict adherence to repo-relative paths (`cli/...`, `02-spec/...`, `.ai-memory/...`). Zero hardcoded absolute system paths.
6. **Zero External Dependencies:** The WebSocket CDP engine must leverage standard Go packages (`net`, `crypto/rand`, `encoding/binary`, `encoding/json`, `io`, `time`) without importing heavy external websocket libraries.
7. **Strict Python CDP Script Ban:** AI agents must not generate ad-hoc Python socket scripts for CDP introspection; all CDP interaction must use `gitmap cdp`.

---

## 3. Subtask Breakdown Matrix

| Subtask ID | File | Description | Target Package / Files | Status |
| :--- | :--- | :--- | :--- | :--- |
| **TASK-01** | [01-subtask-spec-and-rca.md](.ai-memory/plans/subtasks/366-golang-script-creation-and-analysis/01-subtask-spec-and-rca.md) | Specification and Root Cause Analysis Documentation | `02-spec/21-app/366-golang-script-creation-and-analysis/` | pending |
| **TASK-02** | [02-subtask-plan-and-subtasks.md](.ai-memory/plans/subtasks/366-golang-script-creation-and-analysis/02-subtask-plan-and-subtasks.md) | Master Execution Plan and Subtask Manifests Creation | `.ai-memory/plans/`, `.ai-memory/plans/subtasks/366-.../` | in-progress |
| **TASK-03** | [03-subtask-native-go-cdp-client.md](.ai-memory/plans/subtasks/366-golang-script-creation-and-analysis/03-subtask-native-go-cdp-client.md) | Native Go WebSocket CDP Client and Discovery Engine | `cli/cmdcdp/` (`types.go`, `ws_client.go`, `discovery.go`, `evaluator.go`) | pending |
| **TASK-04** | [04-subtask-gitmap-cli-dispatch.md](.ai-memory/plans/subtasks/366-golang-script-creation-and-analysis/04-subtask-gitmap-cli-dispatch.md) | GitMap CLI Integration and Subcommand Dispatch | `cli/root.go`, `cli/cmdcdp/cdp_cmd.go` | pending |
| **TASK-05** | [05-subtask-skill-update-and-python-ban.md](.ai-memory/plans/subtasks/366-golang-script-creation-and-analysis/05-subtask-skill-update-and-python-ban.md) | GitMap Skill Update and Python CDP Ad-Hoc Script Ban | `.agents/skills/gitmap/skill.md`, `02-spec/02-coding-guidelines/` | pending |
| **TASK-06** | [06-subtask-build-bump-and-release.md](.ai-memory/plans/subtasks/366-golang-script-creation-and-analysis/06-subtask-build-bump-and-release.md) | Build, Minor Bump v6.640.0, Release and Fleet Verification | `pkg/constants/constants.go`, `changelog.md` | pending |

---

## 4. Subtasks Detail & Deliverables

### TASK-01: Specification and Root Cause Analysis Documentation
- Author comprehensive specification suite in `02-spec/21-app/366-golang-script-creation-and-analysis/`:
  - `01-root-cause-and-architecture.md`: Exhaustive 4-part root cause analysis explaining the emergence of ad-hoc Python scripts, failure modes of raw TCP framing, and architecture of the native Go CDP client.
  - `02-cdp-engine-and-cli-spec.md`: Detailed Go interface contracts, struct definitions (`CDPTarget`, `CDPRequest`, `CDPResponse`, `DOMProbeResult`, `SelectorHit`), command syntax, and Tokyo Night Lipgloss styles.
  - `readme.md`: Master specification index and roadmap.

### TASK-02: Master Execution Plan and Subtask Manifests Creation
- Author master plan `.ai-memory/plans/pending/366-golang-script-creation-and-analysis.md`.
- Create all 6 subtask documents in `.ai-memory/plans/subtasks/366-golang-script-creation-and-analysis/`:
  - `01-subtask-spec-and-rca.md`
  - `02-subtask-plan-and-subtasks.md`
  - `03-subtask-native-go-cdp-client.md`
  - `04-subtask-gitmap-cli-dispatch.md`
  - `05-subtask-skill-update-and-python-ban.md`
  - `06-subtask-build-bump-and-release.md`
- Register the new plan and subtask folder in `.ai-memory/plans/readme.md`.

### TASK-03: Native Go WebSocket CDP Client and Discovery Engine
- Implement the pure Go RFC 6455 WebSocket client in `cli/cmdcdp/`:
  - `types.go`: Core data transfer models (`CDPTarget`, `CDPRequest`, `CDPResponse`, `DOMProbeResult`, `SelectorHit`, `EvaluatorOptions`).
  - `ws_client.go`: Pure Go RFC 6455 client handling `net.DialTimeout`, handshake negotiation with base64 nonces, 4-byte client frame masking, frame fragmentation, and multi-byte extended payload lengths (`126` and `127`).
  - `discovery.go`: Dynamic target resolution querying `http://127.0.0.1:<port>/json/list` across standard ports (`9222`, `49994`, `9229`), auto-selecting the primary inspectable target.
  - `evaluator.go`: Standardized `Runtime.evaluate` protocol runner embedding the canonical DOM selector ladder for input fields, comboboxes, Monaco editors, and contenteditable elements.

### TASK-04: GitMap CLI Integration and Subcommand Dispatch
- Implement `cli/cmdcdp/cdp_cmd.go` registering the `cdp` command tree:
  - `gitmap cdp probe [--port <p>] [--target <id>] [--json]`
  - `gitmap cdp eval <expression> [--port <p>] [--target <id>] [--json]`
  - `gitmap cdp targets [--port <p>] [--json]`
- Wire `cdpCmd` into `cli/root.go`.
- Render results using Tokyo Night Lipgloss styles (`#7aa2f7`, `#bb9af7`, `#7dcfff`, `#414868`), with fallback to structured JSON when `--json` is supplied.

### TASK-05: GitMap Skill Update and Python CDP Ad-Hoc Script Ban
- Update `.agents/skills/gitmap/skill.md` with:
  - Documentation and syntax examples for `gitmap cdp probe`, `gitmap cdp eval`, and `gitmap cdp targets`.
  - An explicit, high-priority rule banning AI agents from authoring ad-hoc Python socket scripts to inspect Chrome DevTools Protocol.
  - Clear troubleshooting instructions for DevTools connectivity and port auto-discovery.

### TASK-06: Build, Minor Bump v6.640.0, Release and Fleet Verification
- Perform compilation verification of `gitmap.exe` ensuring zero compiler errors and strict adherence to coding guidelines.
- Bump minor version to `v6.640.0` in `pkg/constants/constants.go`.
- Update `changelog.md` with release notes detailing native Go CDP client integration and the Python ad-hoc script ban.
- Execute GitMap release commit: `gitmap cpr "v6.640.0 native Go CDP client and DOM probe engine"`.
- Verify local binary deployment and CLI execution via `gitmap cdp --help`.
