# Subtask 04: GitMap CLI Integration and Subcommand Dispatch

- **Task Code:** `TASK-04`
- **Slug:** `golang-script-creation-and-analysis - Task 04`
- **Parent Plan:** [366-golang-script-creation-and-analysis.md](../../pending/366-golang-script-creation-and-analysis.md)
- **Target Path:** `cli/cmdcdp/cdp_cmd.go`, `cli/root.go`
- **Status:** `pending`
- **Assigned Worker:** GitMap CLI Engineer

---

## 1. Goal

Integrate the CDP engine into the GitMap CLI command tree by wiring `gitmap cdp [probe|eval|targets]`, adding argument and flag parsing, and rendering formatted Tokyo Night Lipgloss cards for human developers and structured JSON for AI agents.

---

## 2. In Scope

1. **Cobra Command Tree Implementation (`cli/cmdcdp/cdp_cmd.go`):**
   - Implement `NewCdpCmd()` returning the parent `cdp` command.
   - Implement subcommands:
     - `gitmap cdp probe`: Runs discovery, connects to active target WebSocket, runs the default DOM input selector ladder, and reports matches.
     - `gitmap cdp eval <expression>`: Evaluates an arbitrary JavaScript expression against the active page target and outputs the result.
     - `gitmap cdp targets`: Lists all inspectable browser/IDE targets discovered via `/json/list` with their title, URL, type, and WebSocket URL.
   - Support flags:
     - `--port <int>` (explicit port override, e.g. 9222, 49994).
     - `--target <string>` (target ID substring or index).
     - `--json` (emit raw JSON envelope for AI consumption).

2. **Root Command Registration (`cli/root.go`):**
   - Register `cmdcdp.NewCdpCmd()` into `rootCmd.AddCommand()`.
   - Ensure clean routing and consistent CLI help annotations.

3. **Tokyo Night Lipgloss Terminal Output:**
   - Design high-contrast visual summary cards:
     - Header badge with target title and port.
     - Selector ladder table showing match status (`HIT` / `MISS`), element tag, aria-label, placeholder, and preview text.
     - Active element indicator and count of contenteditable elements.
   - Guard rendering so `--json` outputs clean JSON directly to stdout without ANSI escape sequences.

---

## 3. Out of Scope

- Core WebSocket framing logic (handled in Task 03).
- Updating documentation skills (handled in Task 05).
- Release ceremony and binary swap (handled in Task 06).

---

## 4. Architectural Invariants

- **Leaf Error Sovereignty:** Commands return `error` up to Cobra execution handlers.
- **Positive Booleans Only:** Flags named with positive semantics (`isJSONOutput`, `hasTargetOverride`).
- **Function Limit:** Helper rendering functions <= 15 lines max with guard clauses.
- **Terminal UI Standards:** Conform to GitMap Tokyo Night Lipgloss palette guidelines.

---

## 5. Verification Checklist

- [ ] `gitmap cdp` command tree compiles and is wired into `cli/root.go`.
- [ ] Subcommands `probe`, `eval`, and `targets` handle flags and arguments correctly.
- [ ] `--json` produces valid, parseable JSON on stdout.
- [ ] Terminal Lipgloss UI outputs legible cards without formatting glitches.
- [ ] Zero unhandled errors or raw `os.Exit` calls in leaf functions.
