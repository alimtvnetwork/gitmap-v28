# Subtask 02: Master Execution Plan and Subtask Manifests Creation

- **Task Code:** `TASK-02`
- **Slug:** `golang-script-creation-and-analysis - Task 02`
- **Parent Plan:** [366-golang-script-creation-and-analysis.md](../../pending/366-golang-script-creation-and-analysis.md)
- **Target Path:** `.ai-memory/plans/pending/`, `.ai-memory/plans/subtasks/366-golang-script-creation-and-analysis/`
- **Status:** `in-progress`
- **Assigned Worker:** Plan Author Subagent

---

## 1. Goal

Author and publish the master execution plan and individual subtask manifests for Task 366 (Golang Script Creation and Analysis), and register them into `.ai-memory/plans/readme.md` following repository conventions.

---

## 2. In Scope

1. **Master Plan Authoring (`366-golang-script-creation-and-analysis.md`):**
   - Author `.ai-memory/plans/pending/366-golang-script-creation-and-analysis.md` containing executive summary, 4-part root cause analysis, benefits of native Go implementation, architectural invariants, and subtask breakdown.

2. **Subtask Manifests Creation:**
   - Author individual markdown files under `.ai-memory/plans/subtasks/366-golang-script-creation-and-analysis/`:
     - `01-subtask-spec-and-rca.md`
     - `02-subtask-plan-and-subtasks.md`
     - `03-subtask-native-go-cdp-client.md`
     - `04-subtask-gitmap-cli-dispatch.md`
     - `05-subtask-skill-update-and-python-ban.md`
     - `06-subtask-build-bump-and-release.md`
   - Ensure every subtask file includes a `Slug:` matching `golang-script-creation-and-analysis - Task NN`.

3. **Master Registry Updates (`.ai-memory/plans/readme.md`):**
   - Register `366-golang-script-creation-and-analysis.md` under `## 3. Pending Plans`.
   - Register `[366-golang-script-creation-and-analysis]` under `## 4. Active Subtasks Directories`.

---

## 3. Out of Scope

- Go code implementation of the CDP client (Task 03) or CLI dispatchers (Task 04).
- Git commits, builds, or test suite execution.

---

## 4. Architectural Invariants

- **Relative Paths Only:** All file links and paths must be relative to the Git repository root.
- **Lowercase Filenames:** Strictly lowercase filenames across all created plans and subtasks.
- **Zero Build / Zero Test Rule:** No compilation or test runs during plan authoring.

---

## 5. Verification Checklist

- [x] `.ai-memory/plans/pending/366-golang-script-creation-and-analysis.md` created with complete 6-subtask matrix.
- [x] Subtasks `01` through `06` created in `.ai-memory/plans/subtasks/366-golang-script-creation-and-analysis/`.
- [x] Every subtask document contains `Slug: golang-script-creation-and-analysis - Task NN`.
- [ ] Registered in `.ai-memory/plans/readme.md` under sections 3 and 4.
