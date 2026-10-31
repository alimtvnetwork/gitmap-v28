# Subtask 06: Build, Minor Bump v6.640.0, Release and Fleet Verification

- **Task Code:** `TASK-06`
- **Slug:** `golang-script-creation-and-analysis - Task 06`
- **Parent Plan:** [366-golang-script-creation-and-analysis.md](../../pending/366-golang-script-creation-and-analysis.md)
- **Target Path:** `pkg/constants/constants.go`, `changelog.md`
- **Status:** `pending`
- **Assigned Worker:** Release Orchestration Agent

---

## 1. Goal

Conduct static code verification, bump the minor version to `v6.640.0`, update changelog documentation, execute the GitMap semantic release commit, and verify fleet-wide binary operation of the new `gitmap cdp` CLI toolchain.

---

## 2. In Scope

1. **Static Analysis & Compilation Verification:**
   - Verify that all newly created Go files in `cli/cmdcdp/` compile cleanly with zero errors.
   - Verify compliance with coding guidelines (functions <= 15 lines, positive booleans, zero nested ifs).

2. **Version Bump & Changelog Update:**
   - Update `pkg/constants/constants.go` bumping version string to `v6.640.0`.
   - Update `changelog.md` detailing:
     - Pure Go standard library RFC 6455 WebSocket client in `cli/cmdcdp/`.
     - New CLI commands: `gitmap cdp probe`, `gitmap cdp eval`, and `gitmap cdp targets`.
     - Automatic debug port and target discovery over `/json/list`.
     - High-contrast Tokyo Night Lipgloss terminal cards and `--json` support.
     - Strict policy ban on AI agents authoring ad-hoc Python websocket scripts for CDP.

3. **Binary Swap & Local Verification:**
   - Compile `gitmap.exe` and swap binary into `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`.
   - Verify `gitmap cdp --help` outputs registered commands and description.

4. **GitMap Semantic Release Ceremony:**
   - Execute release commit using the mandatory semantic verb:
     ```bash
     gitmap cpr "v6.640.0 native Go CDP client and DOM probe engine"
     ```
   - Push release to remote tracking branch.

---

## 3. Out of Scope

- Core Go implementation of WebSocket or CLI dispatch (handled in Tasks 03 and 04).
- Spec authoring or skill documentation (handled in Tasks 01, 02, and 05).

---

## 4. Architectural Invariants

- **Mandatory Pre-Flight Pull:** Always run `gitmap pull` prior to releasing.
- **Dedicated Semantic Commit Verb:** Must use `gitmap cpr` for releases (`v6.640.0 ...`). Total ban on raw `git commit`, `git add -A`, or brackets.
- **Error Sovereignty:** Clean compilation with zero build errors.

---

## 5. Verification Checklist

- [ ] Version constant in `pkg/constants/constants.go` is `v6.640.0`.
- [ ] `changelog.md` reflects `v6.640.0` features and fixes.
- [ ] `gitmap cdp --help` displays commands properly in local binary.
- [ ] Release commit cleanly executed using `gitmap cpr`.
