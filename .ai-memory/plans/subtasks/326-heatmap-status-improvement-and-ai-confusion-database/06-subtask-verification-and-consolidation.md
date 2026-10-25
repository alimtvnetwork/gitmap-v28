# Subtask Plan 06: Verification, Registry Consolidation and Atomic Completion

- **Subtask Slug:** `06-subtask-verification-and-consolidation`
- **Parent Plan:** `326-heatmap-status-improvement-and-ai-confusion-database`
- **Specification Reference:** `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/03-migration-and-verification-spec.md`
- **Target Repository:** `gitmap-v28`
- **Status:** `PENDING`
- **Created At:** 2026-10-25

---

## 1. Objectives & Scope

This subtask performs end-to-end retrospective verification across all functional deliverables of Plan 326, synchronizes top-level specification and plan registry indexes, reconciles the audit ledger, executes pre-commit security gates, and finalizes the task with an atomic GitMap push:

1. **Retrospective AI Quality Audit:**
   - Audit all 5 core task deliverables against user requirements and acceptance criteria:
     - Heatmap status visual overhaul (padding, Lipgloss cards, colors, completed/pending tasks, tree-view command examples, `--json` flag).
     - Work directory guard & safe stashing (`.git` boundary checks, transparent `[M]`, `[D]`, `[?]` file flags, condensed summary, suggested inspection command).
     - AI Confusion SQLite database in `repo-cache`, CLI commands (`add`, `view`, `complete`, `sync`), JSON file ingestion, safe sync replay engine.
     - Skills and prompt remediation (elimination of malformed commit prefixes, enforcement of 5 verbs and subsystem module taxonomy).
     - Defensive formatting engine updates in `cli/cmdcommit/commit_formatter.go`.
2. **Registry Index Synchronization:**
   - Register Spec 326 in `02-spec/21-app/readme.md` under Section 2 (Active Application Specifications).
   - Register Plan 326 in `.ai-memory/plans/readme.md` under Section 1 (Active Plans).
3. **Master Ledger Reconciliation:**
   - Update `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/00-master-audit-ledger.md` to reflect completed deliverables, verified outcomes, and file hashes.
   - Update `.ai-memory/plans/pending/326-heatmap-status-improvement-and-ai-confusion-database.md` with checked-off subtask items.
4. **Pre-Commit Security Gate:**
   - Execute forbidden string scanning (`check-forbidden-strings.py`).
   - Run AUM pattern search for tokens, private keys, and secrets.
5. **Final Atomic GitMap Commit & Push:**
   - Complete the run using the approved GitMap feature commit:
     ```bash
     gitmap cpf "HEATMAP - status visual overhaul and ai confusion database"
     ```
   - Guarantee the completion invariant: **No Push = Not Done**.

---

## 2. Target Files & Disjoint Ownership

| Target File Path | Purpose | Modification Category |
| :--- | :--- | :--- |
| `02-spec/21-app/readme.md` | Master domain index of application specifications | Spec Registry |
| `.ai-memory/plans/readme.md` | Master execution registry of active/completed plans | Memory Registry |
| `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/00-master-audit-ledger.md` | Master traceability audit ledger | Spec Ledger |
| `.ai-memory/plans/pending/326-heatmap-status-improvement-and-ai-confusion-database.md` | Master plan roadmap | Memory Plan |

---

## 3. Atomic Implementation Checklist

### Phase 1: Retrospective Quality Audit Across All Deliverables

- [ ] **1.1**: Audit Heatmap Status Overhaul:
  - Verify `cli/heatmap/heatmap_cmd_render.go` renders Lipgloss cards with appropriate padding and colors.
  - Verify `cli/heatmap/heatmap_tasks.go` queries and surfaces recent completed and pending tasks.
  - Verify `cli/heatmap/heatmap_cmd.go` supports the `--json` flag and returns structured JSON conforming to `HeatmapStatusPayload`.
  - Verify tree-view command examples are included in terminal output.

- [ ] **1.2**: Audit Work Directory Boundary Guard & Safe Stashing:
  - Verify `cli/cmdfix/fix_nonrepo.go` checks for `.git` presence before staging or committing in workspace folders.
  - Verify `cli/cmdfix/fix_stashing.go` prints transparent file listings with `[M]`, `[D]`, `[?]` flags before stashing.
  - Verify `cli/cmdpull/pull_efficient_render.go` summarizes unchanged repositories as `"Updated all other <N> repositories (no changes)"`.
  - Verify a recommended one-liner command is suggested for dirty or changed repository inspection.

- [ ] **1.3**: Audit AI Confusion SQLite Database & Commands:
  - Verify `cli/store/ai_confusion_store.go` and `cli/store/ai_confusion_schema.go` create normalized `AiConfusion` table in `repo-cache/<repo>/ai_confusion.db`.
  - Verify `gitmap rc ai-confusion` commands (`add`, `view`, `complete`, `sync`) are registered and dispatched in `cli/cmd/rc_ai_confusion.go`.
  - Verify `--file` flag accepts JSON instructions and inserts records into the database.
  - Verify safe pull & sync protocol: exports pending local changes to SQLite, pulls repo-cache with rebase, and replays local records.

- [ ] **1.4**: Audit Commit Format & Skills Remediation:
  - Verify `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` mandates the hyphen format (`MODULE - summary`) and defines `MODULE` as an architectural subsystem.
  - Verify `.agents/skills/gitmap/skill.md` reinforces the 5 semantic verbs and title-case prefix standard.
  - Verify `.agents/skills/muse-master-prompt/skill.md` uses canonical subsystem examples.
  - Verify `cli/cmdcommit/commit_formatter.go` strips redundant prefixes, neutralizes pseudo-module collisions, and maintains $\le 15$ line functions.

---

### Phase 2: Registry Index Updates

- [ ] **2.1**: Update `02-spec/21-app/readme.md`:
  - Locate Section 2 (Active Application Specifications).
  - Add an entry for Spec 326:
    ```markdown
    - [326-heatmap-status-improvement-and-ai-confusion-database](326-heatmap-status-improvement-and-ai-confusion-database/00-master-audit-ledger.md) — Plan 326: Heatmap Status Improvement and AI Confusion Database: Lipgloss status card redesign, completed/pending task display, --json flag, fix-all non-git directory guard, humanized auto-stashing, repo-cache normalized AI confusion SQLite DB, rc CLI commands, JSON instruction ingestion, and commit prefix remediation (Specs: [Ledger](326-heatmap-status-improvement-and-ai-confusion-database/00-master-audit-ledger.md), [Architecture](326-heatmap-status-improvement-and-ai-confusion-database/01-architecture-spec.md), [Component & CLI](326-heatmap-status-improvement-and-ai-confusion-database/02-component-and-cli-spec.md), [Migration & Verification](326-heatmap-status-improvement-and-ai-confusion-database/03-migration-and-verification-spec.md)) (Status: `in_progress`)
    ```

- [ ] **2.2**: Update `.ai-memory/plans/readme.md`:
  - Locate Section 1 (Active Plans).
  - Add an entry for Plan 326:
    ```markdown
    - [326-heatmap-status-improvement-and-ai-confusion-database.md](pending/326-heatmap-status-improvement-and-ai-confusion-database.md) — Heatmap Status Improvement and AI Confusion Database: Visual padding overhaul, task integration, JSON telemetry, safe stashing, repo-cache AI confusion database, and commit prefix remediation (Spec: [326-heatmap-status-improvement-and-ai-confusion-database](../../02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/00-master-audit-ledger.md)) (Status: `active`)
    ```

---

### Phase 3: Audit Ledger Reconciliation

- [ ] **3.1**: Reconcile `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/00-master-audit-ledger.md`:
  - Update traceability matrix status from `Planned` to `Implemented` / `Verified` for all rows.
  - Document verified file paths and key interfaces.

- [ ] **3.2**: Reconcile `.ai-memory/plans/pending/326-heatmap-status-improvement-and-ai-confusion-database.md`:
  - Mark all subtask checklist items as completed (`[x]`).
  - Update status field to `COMPLETED` upon final review.

---

### Phase 4: Pre-Commit Security Gate & Atomic Finish

- [ ] **4.1**: Execute Relative Path Linter:
  ```powershell
  python linter-scripts/check-relative-paths.py
  ```
  - Assert 0 absolute filesystem paths or file URI schemes remain.

- [ ] **4.2**: Execute Secrets Gate:
  ```powershell
  python linter-scripts/check-forbidden-strings.py
  gitmap aum search -r "(BEGIN [A-Z ]*PRIVATE KEY|AKIA[0-9A-Z]{16}|gh[pousr]_[A-Za-z0-9]{36}|sk-[A-Za-z0-9]{20,}|xox[baprs]-[A-Za-z0-9-]{10,})"
  ```
  - Assert zero secrets, keys, or credentials are staged.

- [ ] **4.3**: Execute Final Atomic Commit & Push:
  ```powershell
  gitmap cpf "HEATMAP - status visual overhaul and ai confusion database"
  ```
  - Assert remote push succeeds and clean working tree is maintained.
