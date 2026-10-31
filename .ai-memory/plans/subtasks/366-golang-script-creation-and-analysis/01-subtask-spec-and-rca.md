# Subtask 01: Specification and Root Cause Analysis Documentation

- **Task Code:** `TASK-01`
- **Slug:** `golang-script-creation-and-analysis - Task 01`
- **Parent Plan:** [366-golang-script-creation-and-analysis.md](../../pending/366-golang-script-creation-and-analysis.md)
- **Target Path:** `02-spec/21-app/366-golang-script-creation-and-analysis/`
- **Status:** `pending`
- **Assigned Worker:** Spec Author Subagent

---

## 1. Goal

Author the comprehensive specification suite detailing the 4-part root cause analysis (RCA), failure modes of ad-hoc Python socket scripts, native Go RFC 6455 WebSocket architecture, struct definitions, CLI subcommand contracts, and Tokyo Night Lipgloss visualization formats for Chrome DevTools Protocol (CDP) DOM evaluation and inspection.

---

## 2. In Scope

1. **Root Cause Analysis & Architecture Spec (`01-root-cause-and-architecture.md`):**
   - 4-part RCA detailing why autonomous AI agents wrote ad-hoc 70-line Python scripts using raw TCP sockets (`socket.socket()`).
   - Detailed breakdown of Python approach failure modes: ephemeral port/target ID hardcoding, fragile unbuffered framing failing on payloads >64KB, absence of timeouts and ping/pong handling, and conversation transcript noise.
   - Architectural blueprint for pure Go standard library RFC 6455 WebSocket client in `cli/cmdcdp/` without third-party dependencies.
   - Dynamic port and target discovery mechanics querying `http://127.0.0.1:<port>/json/list`.

2. **CDP Engine & CLI Spec (`02-cdp-engine-and-cli-spec.md`):**
   - Core Go struct contracts: `CDPTarget`, `CDPRequest`, `CDPResponse`, `DOMProbeResult`, `SelectorHit`, and `EvaluatorOptions`.
   - Complete CLI command syntax, arguments, and flag contracts:
     - `gitmap cdp probe [--port <p>] [--target <id>] [--json]`
     - `gitmap cdp eval <expression> [--port <p>] [--target <id>] [--json]`
     - `gitmap cdp targets [--port <p>] [--json]`
   - Default DOM selector ladder specification for targeting active inputs, Monaco editor instances, and contenteditable elements.
   - Tokyo Night Lipgloss design tokens, ANSI escape styling, and structured JSON output schemas.

3. **Master Specification Index (`readme.md`):**
   - Executive index linking specs, deliverables, coding standards, and verification milestones.

---

## 3. Out of Scope

- Writing Go source code in `cli/cmdcdp/` (handled in Tasks 03 and 04).
- Running builds, tests, or version bumping.

---

## 4. Architectural Invariants

- **Relative Paths Only:** All file references must use repo-relative paths (`02-spec/...`, `cli/...`).
- **Positive Booleans Only:** Struct fields and flags must use `is*` or `has*` (`hasTarget`, `isMasked`, `hasExtendedPayload`, `isJSONOutput`).
- **Function Sizing & Zero Nesting:** Design specifications must reflect <= 15 lines max per function and zero nested conditionals.
- **Error Return Sovereignty:** Leaf functions must return `error` or `*apperror.AppError` directly.

---

## 5. Verification Checklist

- [ ] `02-spec/21-app/366-golang-script-creation-and-analysis/01-root-cause-and-architecture.md` exists and covers 4-part RCA, failure modes, and native Go architecture.
- [ ] `02-spec/21-app/366-golang-script-creation-and-analysis/02-cdp-engine-and-cli-spec.md` exists and defines Go types, CLI syntax, selector ladder, and Lipgloss tokens.
- [ ] `02-spec/21-app/366-golang-script-creation-and-analysis/readme.md` exists as master index.
