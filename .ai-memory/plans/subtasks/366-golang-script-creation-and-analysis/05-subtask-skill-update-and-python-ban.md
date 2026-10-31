# Subtask 05: GitMap Skill Update and Python CDP Ad-Hoc Script Ban

- **Task Code:** `TASK-05`
- **Slug:** `golang-script-creation-and-analysis - Task 05`
- **Parent Plan:** [366-golang-script-creation-and-analysis.md](../../pending/366-golang-script-creation-and-analysis.md)
- **Target Path:** `.agents/skills/gitmap/skill.md`, `02-spec/02-coding-guidelines/`
- **Status:** `pending`
- **Assigned Worker:** Documentation & Governance Subagent

---

## 1. Goal

Update GitMap agent skill definitions and coding guidelines to document the new `gitmap cdp` CLI capabilities and establish an explicit architectural ban prohibiting autonomous AI agents from writing ad-hoc Python socket scripts for Chrome DevTools Protocol inspection.

---

## 2. In Scope

1. **GitMap Skill Documentation (`.agents/skills/gitmap/skill.md`):**
   - Add dedicated section for Chrome DevTools Protocol (CDP) commands:
     - `gitmap cdp probe [--port <p>] [--target <id>] [--json]`
     - `gitmap cdp eval <expression> [--port <p>] [--target <id>] [--json]`
     - `gitmap cdp targets [--port <p>] [--json]`
   - Include practical usage examples for DOM inspection, active element discovery, and JavaScript expression evaluation in Antigravity IDE and web visualizers.

2. **Architectural Ban on Ad-Hoc Python CDP Scripts:**
   - Enforce an explicit rule in `.agents/skills/gitmap/skill.md`:
     > **STRICT BAN ON AD-HOC PYTHON CDP WEBSOCKET SCRIPTS:**
     > AI agents are strictly prohibited from generating raw Python socket scripts (`socket.socket()`, manual HTTP Upgrade handshakes, `struct.pack`/`unpack` framing) to connect to Chrome DevTools Protocol.
     > All browser and IDE inspection must exclusively execute via `gitmap cdp probe` and `gitmap cdp eval`.
   - Document the technical failure modes of the Python approach (ephemeral ports, payload truncation >64KB, transcript noise).

3. **Coding Guidelines Cross-Referencing (`02-spec/02-coding-guidelines/`):**
   - Update guideline documentation referencing the ban on ad-hoc automation scripts when native GitMap CLI commands exist.

---

## 3. Out of Scope

- Go code implementation (Tasks 03 and 04).
- Building binary or executing release ceremony (Task 06).

---

## 4. Architectural Invariants

- **Relative Paths Only:** Link documents using repository-relative links.
- **Lowercase Filenames:** All files must maintain lowercase casing.
- **Clarity and Precision:** The ban must be unambiguous and provide direct substitute CLI commands.

---

## 5. Verification Checklist

- [ ] `.agents/skills/gitmap/skill.md` contains comprehensive `gitmap cdp` documentation.
- [ ] Explicit ban on ad-hoc Python CDP WebSocket scripts clearly stated in the skill.
- [ ] Alternative commands (`gitmap cdp probe`, `gitmap cdp eval`) prominently featured.
