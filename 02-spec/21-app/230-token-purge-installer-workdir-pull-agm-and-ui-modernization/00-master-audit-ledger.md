# Master Audit Ledger: Spec 230

- **Slug:** `230-token-purge-installer-workdir-pull-agm-and-ui-modernization`
- **Spec Directory:** `02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/`
- **Parent Plan:** `.ai-memory/plans/pending/230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md`
- **Target Version:** `v6.482.0`
- **Created:** 2026-10-06T17:25:00+08:00
- **Status:** In Progress (Phase 1: Spec Authoring & Subtask Generation)

---

## 1. Traceable Task Matrix

| Task ID | Component / Area | Status | Target Files | Evidence / Artifact |
|:---|:---|:---:|:---|:---|
| **Task-01** | Refresh Token Security Purge | **DONE** | `/home/a/.cache/vmware/drag_and_drop/` | Purged `agm_accounts_backup_2026-10-06.json` (exit 0) |
| **Task-02** | Installer `$work` / `$def` & Multi-Workdir | **DONE** | `cli/cmd/cd_workdir_resolver.go`, `cli/cmd/cd_workdir_resolver_test.go` | Keyword resolution for `$work`, `$def`, multi-workdirs; companion tests pass (exit 0) |
| **Task-03** | GitMap Pull Auto-Remediation & Conflict Resolution | **DONE** | `cli/cloner/safe_pull.go`, `cli/cmdpull/pull_remediation.go` | Fast-forward retry, SafeAbortMerge/Rebase, batch remediation prompt (exit 0) |
| **Task-04** | AGM Linux Update Fix & Quiet Execution | **DONE** | `cli/cmdinstall/agm_update.go`, `cli/cmdssh/ssh_update_remote.go` | Silent success, bounded stack trace on failure, batch prompt (exit 0) |
| **Task-05** | Interim Testing from Recent Commits | **DONE** | `cli/cmdpull/`, `cli/cmdpipeline/`, `cli/scanner/` | Targeted unit test execution passes across modified packages (0.010s - 0.107s) |
| **Task-06** | Antigravity Multi-Instance Prompt Query API | **DONE** | `cli/cmdagy/`, `cli/cmdui/ui_server.go` | Multi-instance prompt JSON CLI & REST API endpoints verified via CLI & curl |
| **Task-07** | Settings UI Modernization & Design Specs | **DONE** | `src/pages/Settings.tsx`, `cli/cmdui/ui_assets.go` | 4-Plane depth hierarchy, accessible contrast, curated UI/UX reference catalog |
| **Task-08** | Image Audit, Blur, and Git History Purge Plan | **DONE** | `.ai-memory/plans/`, `assets/screenshots/` | Full 161-image audit, 23 text-only screenshots purged, 4 sanitized assets generated |

---

## 2. Multi-Agent Concurrency & Disjoint Ownership

- **Lead Orchestrator:** Master ledger, parent plan, indexes (`plans/readme.md`, `what-to-read.md`), and final atomic commit.
- **Worker 01 (Core CLI & Pull Remediation):**
  - Specs: `01-architecture-spec.md`
  - Subtasks: `01-refresh-token-purge...`, `02-gitmap-installer...`, `03-pull-auto-remediation...`, `04-agm-linux-update...`
- **Worker 02 (Antigravity API, UI & Image Sanitation):**
  - Specs: `02-component-and-cli-spec.md`
  - Subtasks: `05-antigravity-multi-instance...`, `06-settings-ui-modernization...`, `07-image-audit...`, `08-targeted-interim-tests...`

---

## 3. Precedence & Constraints Ledger

- **R1 Total Ban on Full Test Suites:** No `go test ./...` or `06-cicd-local-runner.py`. File-scoped checks only.
- **R2 Zero Git Trailing Colons:** Use `gitmap cpf "<module> - <summary>"` with hyphens.
- **R3 Strict Relative Git Paths:** All markdown and release links must start with repo relative paths.
- **R4 Zero Solo Execution:** All phases partitioned across `A = 2`, `H = 2` autonomous subagents.
