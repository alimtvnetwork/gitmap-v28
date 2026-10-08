# Architecture Specification: Forward-Slash Path Hygiene, Pre-Pull Remediation & Interactive Fix Prompt

- **Feature Slug:** `243-pull-forward-slash-paths-pre-pull-and-interactive-remediation-prompt`
- **Module:** `cli/cmdpull` & `cli/cmd`
- **Specification Version:** 1.0.0
- **Status:** APPROVED

---

## 1. Problem Statement & System Context

During large batch repository pull operations (e.g. `gitmap pa` / `gitmap pull-all`), users encountered two critical usability issues:

1. **Windows Escaped Backslashes in Output:**
   - Command recommendations and remediation options printed escaped Windows backslashes:
     ```text
     Option 2 (Stash Changes): git -C "D:\\work\\gitmap" stash
     Option 1 (Track / Stage): git -C "D:\\work\\letsmarknow" add .
     ```
   - This causes ugly formatting in terminal UIs and complicates cross-platform copy-pasting into bash, zsh, and PowerShell environments. Paths must universally use clean forward slashes (`/`) without double-backslash escapes (`D:/work/gitmap`).

2. **Missing Remediation Prompt in `gitmap pa`:**
   - When `gitmap pa` completes with dirty (e.g. 4 repos) or failed (e.g. 2 repos) repositories, `handlePullRemediationForRecords` explicitly checked:
     ```go
     if len(remItems) == 0 || isConcisePullOutput(opts.all, opts.showStatus) {
         return
     }
     ```
   - Because `opts.all` was true, the function exited immediately without prompting the user to remediate the dirty or failed repositories.
   - The user expects an interactive prompt after batch pulling:
     ```text
     Remediate dirty / failed repository(ies) now?
       [a/1] Fix all
       [s/2] Fix single / select repository
       [k/q] Skip / Exit
     Choice [a/s/k]:
     ```

3. **Pull Before Changes Invariant:**
   - When remediation actions (such as committing WIP, stashing, or resolving branch states) are triggered, GitMap must always ensure that remote updates are pulled before local changes are finalized, avoiding stale state or upstream desynchronization.

---

## 2. Architecture & Design Principles

```
┌────────────────────────────────────────────────────────────────────────┐
│                        gitmap pa (pull-all)                            │
└──────────────────────────────────┬─────────────────────────────────────┘
                                   │
                                   ▼
              ┌────────────────────────────────────────┐
              │      Concurrent Batch Pull Work        │
              └────────────────────┬───────────────────┘
                                   │
                                   ▼
              ┌────────────────────────────────────────┐
              │     Concise Active Results Render      │
              │  • Updated (green)                     │
              │  • Dirty (yellow, forward slashes)     │
              │  • Failed (red, forward slashes)       │
              └────────────────────┬───────────────────┘
                                   │
                                   ▼
              ┌────────────────────────────────────────┐
              │   Interactive Remediation Gate         │
              │   Has Dirty or Failed Repositories?    │
              └────────┬──────────────────────┬────────┘
                       │ Yes (TTY)            │ No / Non-TTY
                       ▼                      ▼
        ┌───────────────────────────┐   ┌───────────────────────────┐
        │  Prompt User:             │   │ Complete / Finalize Task  │
        │  [a] Fix all              │   │ Return cleanly            │
        │  [s] Fix single / select  │   └───────────────────────────┘
        │  [k] Skip                 │
        └──────────────┬────────────┘
                       │
        ┌──────────────┼──────────────┐
        ▼              ▼              ▼
   [a] Fix All    [s] Single     [k] Skip
        │              │              │
        ▼              ▼              ▼
 Pull Before   Pull Before       Exit Prompt
 Changes       Changes
```

### Invariants:
1. **Zero Double-Backslashes:** All repository directories, file paths, and command options emitted in terminal output MUST pass through `filepath.ToSlash(path)`.
2. **Interactive TTY Safeguard:** In non-interactive environments (CI/CD, piped stdout, `--json`, or explicit `--no-fix`), bypass interactive prompts automatically without blocking.
3. **Pre-Pull Safety:** Any fix flow must pull remote changes prior to applying stashes, commits, or branch rebases.
4. **Clean Exit on Skip:** Choosing `k` / `q` / `skip` must exit the remediation flow with code 0.
