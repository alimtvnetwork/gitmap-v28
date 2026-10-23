# Execution Plan: GitMap Mailing and Scheduling Enhancements

**Slug:** `gitmap-mailing-and-scheduling-enhancements`  
**Task Name:** GitMap Mailing and Scheduling Enhancements  
**Status:** COMPLETED  
**Lifecycle Budget:** N = 400 Steps (Phase 1: 1-100, Phase 2: 101-250, Phase 3: 251-400)  
**Concurrency:** A = 3, H = 2, Capacity = 6  

---

## 1. User Request (Verbatim)

```text
# GitMap Mailing and Scheduling Enhancements: high priority instruction, non-negotiable task

Schedule section. `gitmap` schedule section does not have email. We should be able to check other mailboxes, and that should be in the schedule command. We should be able to send schedule email, or we should be able to send email based on some task or activity, let's say creating a file on a specific folder or modifying a file on a specific folder that would trigger and sending email to the notifier emails, whatever that is there. These type of things need to be implemented. And when we do the `gitmap`, it would also share by hiding some of the mail parts. That this is the sender and these are the notifier emails, by hiding some of the parts so that it's not visible all the time. Make sure of that, and try to add more commands related to mailing. There should be a mailing UI where we could send, and there should be a UI for schedule settings as well, where we could schedule or customize the notifier based on folder, file, and all the CLI commands need to be there as well. Can you make a detailed plan of it first? Very detailed plan with all the commands that is possible, all the UI that is possible. First, write the detailed spec.
```

---

## 2. Architecture & Design Blueprint

### Subsystem Topology
```
CLI Dispatch:
  gitmap schedule [mail-check | email | watch-folder | ui]  ->  cmdschedule
  gitmap mail [send | check | test | template | notify | ui] ->  cmdmail

Core Subsystems:
  Multi-Mailbox Checker:
  ├── Inspects multi-mailbox configurations (Split-DB / secrets)
  └── Evaluates unread counts, snippets, and polling intervals

  Folder Watcher & Event Notifier:
  ├── Stat & ModTime periodic polling with IntervalGuard (>= 5s, default 8s)
  ├── Event debouncing window (groups burst file mutations)
  └── Dispatches notification emails to configured recipient lists

  Privacy Masking Engine:
  ├── Standard email masking: user@domain.com -> u***r@d***n.com
  └── Enforced across CLI displays, logs, and public REST/UI envelopes

  Embedded Web Dashboards (Zero-CDN):
  ├── Mailing Composer UI (/mail/ui): Interactive draft, templates, accounts
  └── Schedule Settings UI (/schedule/ui): Watchers, cron jobs, notifier rules
```

---

## 3. Subtask Work Breakdown Matrix

| Subtask ID | Title | Disjoint Files Owned | Assigned Role | Status |
|---|---|---|---|---|
| **Subtask-01** | Git Sync, Preflight & Task Initialization | `ledger.md`, `.ai-memory/plans/` | Lead | DONE |
| **Subtask-02** | Architecture & Specification Authoring | `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/` | Spec Authors (A=3) | DONE |
| **Subtask-03** | Multi-Mailbox Checker, Schedule Email CLI & Privacy Masker | `cli/cmdschedule/schedule_mail_check.go`, `cli/cmdschedule/schedule_email.go`, `cli/cmdschedule/privacy_masker.go`, `cli/cmdschedule/schedule_cmd.go` | Worker 01 | DONE |
| **Subtask-04** | Folder Watcher Event Trigger & Mail CLI Engine | `cli/cmdmail/cmd_mail.go`, `cli/cmdmail/mail_sender.go`, `cli/cmdmail/mail_checker.go`, `cli/cmdmail/folder_watcher.go`, `cli/cmdmail/storage.go` | Worker 02 | DONE |
| **Subtask-05** | Embedded Mailing Composer UI & Schedule Settings UI | `cli/cmdmail/ui_*.go`, `cli/cmdschedule/ui_*.go` | Worker 03 | DONE |
| **Subtask-06** | Targeted Unit Tests & Verification Suite | `cli/cmdmail/*_test.go`, `cli/cmdschedule/*_test.go` | Worker Swarm | DONE |

---

## 4. Verification Evidence & Quality Gates
- **G1 (Multi-Mailbox Polling):** `gitmap schedule mail-check --all` successfully evaluates multiple mailboxes and unread counts.
- **G2 (Scheduled Email Registration):** `gitmap schedule email` stores and triggers automated email delivery tasks.
- **G3 (Folder Event Triggering):** Creating or modifying files in a monitored directory triggers debounced notification emails.
- **G4 (Privacy Masking):** All email addresses are masked by default (`u***r@d***n.com`) across terminal and UI outputs.
- **G5 (Mail CLI Hierarchy):** `gitmap mail send`, `check`, `test`, `template`, `notify` operate properly.
- **G6 (Embedded UIs):** Both `/mail/ui` and `/schedule/ui` serve responsive HTML/JS/CSS without external CDN links.
