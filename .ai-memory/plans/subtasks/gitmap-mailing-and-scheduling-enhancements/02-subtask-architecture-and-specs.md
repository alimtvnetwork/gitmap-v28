# Subtask 02: Architecture & Specification Authoring

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/gitmap-mailing-and-scheduling-enhancements.md`  
> **Status:** DONE  
> **Target Subsystems:** `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/`  
> **Owned Files:**  
> - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/readme.md`  
> - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/01-architecture-overview.md`  
> - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/02-cli-command-contracts.md`  
> - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/03-folder-watcher-and-triggers.md`  
> - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/04-privacy-and-email-masking.md`  
> - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/05-embedded-web-uis.md`  
> - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/06-storage-and-split-db-schema.md`  
> - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/07-verification-and-acceptance-matrix.md`  

---

## 1. Concrete Objectives

1. **Specification Suite Index & Overview (`readme.md`):**
   - Provide an executive overview of the mailing and scheduling enhancements.
   - Index all 7 domain specifications with clear scope boundaries and deliverables.

2. **System Architecture Overview (`01-architecture-overview.md` - Spec Author 01):**
   - Author system topology with Mermaid diagram integrating `cmdschedule` and `cmdmail`.
   - Specify multi-mailbox checking architecture, account registries, and polling intervals.
   - Formulate scheduled email lifecycle from registration, cron evaluation, template compilation, to dispatch and logging.
   - Detail filesystem folder watcher architecture with `IntervalGuard` (>= 5s, default 8s) and event debouncing window.
   - Outline pervasive privacy masking architecture (`u***r@d***n.com`).
   - Ratify Split-DB persistence rules: PascalCase schemas, positive boolean hygiene (`IsActive`, `HasCompleted`, `IsMasked`), and single-writer safety (`MaxOpenConns(1)`).

3. **CLI Command Contracts Specification (`02-cli-command-contracts.md` - Spec Author 01):**
   - Formulate exact syntax, flags, default values, and exit codes for schedule commands:
     - `gitmap schedule mail-check` / `gitmap schedule mail`
     - `gitmap schedule email`
     - `gitmap schedule watch-folder`
     - `gitmap schedule ui`
   - Formulate exact syntax, flags, sub-actions, and exit codes for mail commands:
     - `gitmap mail send`
     - `gitmap mail check`
     - `gitmap mail test`
     - `gitmap mail template` (`list`, `get`, `save`, `del`)
     - `gitmap mail notify` (`list`, `trigger`, `history`)
     - `gitmap mail ui`
   - Define Universal JSON Response Envelope schema with PascalCase attributes.

4. **Folder Watcher & Event Triggers Specification (`03-folder-watcher-and-triggers.md` - Spec Author 02):**
   - Author polling-based filesystem watcher engine using file metadata snapshot caching.
   - Document `IntervalGuard` enforcement (minimum 5s, default 8s) for OS safety.
   - Formulate burst debouncing window (`1s`–`30s`) and change digest construction.
   - Map event triggers to notifier email recipient dispatches.

5. **Privacy & Email Masking Specification (`04-privacy-and-email-masking.md` - Spec Author 02):**
   - Detail privacy masking transformation algorithms: `user@domain.com` -> `u***r@d***n.com`.
   - Define boundary rules for single/double letter usernames and multi-subdomain formats.
   - Document redaction rules across terminal outputs, console log files, and public REST responses.
   - Specify administrative unmasking controls (`--reveal` / `--no-mask`).

6. **Embedded Web Dashboards Specification (`05-embedded-web-uis.md` - Spec Author 03):**
   - Specify zero-CDN responsive Mailing Composer UI (`/mail/ui`).
   - Specify zero-CDN responsive Schedule & Notifier Settings UI (`/schedule/ui`).
   - Define embedded asset bundling, REST endpoints, and security boundaries.

7. **Storage & Split-DB Schema Specification (`06-storage-and-split-db-schema.md` - Spec Author 03):**
   - Formulate DDL for tables: `ScheduledEmail`, `MailboxConfig`, `MailboxStateRecord`, `FolderWatchRule`, `MailNotificationLog`.
   - Enforce pure-Go CGO-free driver, PascalCase columns, and single-writer concurrency.

8. **Verification & Acceptance Matrix (`07-verification-and-acceptance-matrix.md` - Spec Author 03):**
   - Ratify testable criteria for Quality Gates G1 through G6.
   - Outline unit, integration, and E2E verification test cases.

---

## 2. Work Breakdown Matrix (Spec Authors A=3)

| Role | Assigned Document(s) | Subsystem Focus | Status |
|---|---|---|---|
| **Spec Author 01** | `01-architecture-overview.md` | Architecture topology, scheduler integration, lifecycle, safety | DONE |
| **Spec Author 01** | `02-cli-command-contracts.md` | CLI syntax, flags, exit codes, Universal JSON Envelope | DONE |
| **Spec Author 01** | `01-subtask-git-sync-and-preflight.md` | Task preflight, ledger setup, budget allocation | DONE |
| **Spec Author 01** | `02-subtask-architecture-and-specs.md` | Specification authoring subtask decomposition | DONE |
| **Spec Author 02** | `03-folder-watcher-and-triggers.md` | Filesystem poller, IntervalGuard, debouncing, notification triggers | IN_PROGRESS |
| **Spec Author 02** | `04-privacy-and-email-masking.md` | Masking algorithms, redaction standards, reveal authorization | IN_PROGRESS |
| **Spec Author 03** | `05-embedded-web-uis.md` | Zero-CDN Mailing Composer UI and Schedule Settings UI | IN_PROGRESS |
| **Spec Author 03** | `06-storage-and-split-db-schema.md` | SQLite Split-DB schemas, PascalCase tables, positive booleans | IN_PROGRESS |
| **Spec Author 03** | `07-verification-and-acceptance-matrix.md` | Acceptance gates G1-G6, test plans, verification matrix | IN_PROGRESS |

---

## 3. Action Items Checklist

- [x] 1. Author specification suite index in `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/readme.md`.
- [x] 2. Author architecture overview specification in `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/01-architecture-overview.md` (Spec Author 01).
- [x] 3. Author CLI command contracts specification in `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/02-cli-command-contracts.md` (Spec Author 01).
- [x] 4. Author Subtask 01 plan in `.ai-memory/plans/subtasks/gitmap-mailing-and-scheduling-enhancements/01-subtask-git-sync-and-preflight.md` (Spec Author 01).
- [x] 5. Author Subtask 02 plan in `.ai-memory/plans/subtasks/gitmap-mailing-and-scheduling-enhancements/02-subtask-architecture-and-specs.md` (Spec Author 01).
- [ ] 6. Author folder watcher and triggers specification in `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/03-folder-watcher-and-triggers.md` (Spec Author 02).
- [ ] 7. Author privacy and email masking specification in `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/04-privacy-and-email-masking.md` (Spec Author 02).
- [ ] 8. Author embedded web UIs specification in `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/05-embedded-web-uis.md` (Spec Author 03).
- [ ] 9. Author storage and Split-DB schema specification in `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/06-storage-and-split-db-schema.md` (Spec Author 03).
- [ ] 10. Author verification and acceptance matrix in `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/07-verification-and-acceptance-matrix.md` (Spec Author 03).

---

## 4. Verification Evidence

- `01-architecture-overview.md` authored (145 lines, <= 300 lines limit verified).
- `02-cli-command-contracts.md` authored (256 lines, <= 300 lines limit verified).
- `01-subtask-git-sync-and-preflight.md` authored (77 lines, <= 300 lines limit verified).
- `02-subtask-architecture-and-specs.md` authored (verified <= 300 lines).
- All files enforce strict relative Git paths, lowercase filenames, PascalCase tables/columns, positive booleans, single-writer SQLite connection safety, and zero git commands.
