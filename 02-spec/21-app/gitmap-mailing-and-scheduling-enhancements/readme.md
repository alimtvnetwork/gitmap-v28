# GitMap Mailing and Scheduling Enhancements Specification Suite

**Specification Slug:** `gitmap-mailing-and-scheduling-enhancements`  
**Root Path:** `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/`  
**Status:** DRAFT / IN_PROGRESS  
**Associated Plan:** `.ai-memory/plans/gitmap-mailing-and-scheduling-enhancements.md`  

---

## 1. Executive Summary

This specification suite ratifies the architecture, CLI interfaces, event-driven filesystem watcher triggers, privacy-preserving email address masking, dedicated mailing engine, and dual embedded web dashboards for GitMap's mailing and scheduling subsystems.

### Core Capabilities
1. **Multi-Mailbox Checking:** `gitmap schedule mail-check` and `gitmap mail check` supporting multiple accounts, unread detection, and polling intervals.
2. **Scheduled Automated Emails:** `gitmap schedule email` supporting interval timers, cron schedules, file/body templates, and execution telemetry.
3. **Event-Driven Folder/File Watcher Notifier:** `gitmap schedule watch-folder` and `gitmap mail notify watch` monitoring directories for file creation, modification, or deletion, debouncing bursts, and triggering notification emails to configured recipient lists.
4. **Privacy-Preserving Address Masking:** Mandatory address masking across terminal output, logs, and public API envelopes (e.g. `u***r@d***n.com`), protecting sender and recipient PII while preserving verification.
5. **Dedicated Mailing CLI Verb (`gitmap mail`):** `send`, `check`, `test`, `template`, `notify`, `ui`.
6. **Dual Embedded Web Dashboards:**
   - **Mailing Composer UI (`/mail/ui`):** Zero-CDN responsive composer, template library, and delivery monitor.
   - **Schedule & Notifier Settings UI (`/schedule/ui`):** Interactive task manager, folder watcher rule designer, and notifier recipient editor.

---

## 2. Specification Directory Structure

| Specification Document | Focus Area |
|---|---|
| [01-architecture-overview.md](01-architecture-overview.md) | High-level system architecture, scheduler integration, multi-mailbox checker, and event trigger topology. |
| [02-cli-command-contracts.md](02-cli-command-contracts.md) | Full CLI syntax for `gitmap schedule ...` and `gitmap mail ...` verbs, flags, and exit codes. |
| [03-folder-watcher-and-triggers.md](03-folder-watcher-and-triggers.md) | Filesystem polling watcher, `IntervalGuard`, debouncing, and automated notification firing. |
| [04-privacy-and-email-masking.md](04-privacy-and-email-masking.md) | Privacy masking algorithms, redaction standards, and unmasking authorization controls. |
| [05-embedded-web-uis.md](05-embedded-web-uis.md) | Architecture of Mailing Composer UI and Schedule/Notifier Settings UI. |
| [06-storage-and-split-db-schema.md](06-storage-and-split-db-schema.md) | SQLite Split-DB schemas, PascalCase tables, and positive boolean flags. |
| [07-verification-and-acceptance-matrix.md](07-verification-and-acceptance-matrix.md) | Numbered testable acceptance criteria, quality gates, and E2E verification procedures. |
