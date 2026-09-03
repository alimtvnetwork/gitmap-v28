# 119 — Git Pull Diagnostics & Display Architecture

## Overview

**Module Number:** 119  
**Version:** 1.0.0  
**Updated:** 2026-09-03  
**Status:** Production-Ready  
**AI Confidence:** Production-Ready  
**Ambiguity Score:** None  

---

## Purpose

Standard `git pull` across multi-repository workspaces frequently fails due to non-unlink working-tree collisions, untracked file conflicts, dirty stashes, or diverged remote tracking branches. This specification formalizes the `gitmap pull` diagnostic engine and the standardized tabular status readout.

---

## User Requirements & Problem Statement

When pulling multiple repositories in batch mode, developers require:
1. Instant identification of failed repositories with categorized root cause diagnosis instead of generic git abort traces.
2. A unified terminal summary table displaying repo identity, current branch, latest remote branch, open PR count, sync status, commit SHA, and operation elapsed time.
3. Explicit actionable guidance on how to resolve untracked collisions before merging.

---

## Terminal Display Contract

### 1. Diagnosis Banner

When a repository fails during pull, the diagnostic classifier intercepts stderr and categorizes the failure:

```text
Please move or remove them before you merge.
Aborting
Diagnosis: non-unlink git pull failure (check auth/merge or run pull manually for full output)
  ── 6 failure(s) total ──
```

### 2. Multi-Repo Status Table Format

```text
  REPO                          BRANCH  LATEST BRANCH    PR/TRACK  STATUS      SHA      TIME
  ──────────────────────────────────────────────────────────────────────────────────────────
  ai-empathy-prompt-tuner-v1    main    main             0 PRs     UP_TO_DATE  f1d94df  1.0s
  prompt-architect-v2           main    release/v1.35.0  0 PRs     UP_TO_DATE  03be798  1.0s
  prompts-connect-v3            main    main             0 PRs     UP_TO_DATE  245dbf9  1.0s
  core-v9                       main    main             1 PRs     BEHIND (2)  a8b712c  1.4s
```

### 3. Status Enums

| Status | Color Code | Condition |
|---|---|---|
| `UP_TO_DATE` | Green | Local HEAD matches remote upstream tracking branch |
| `UPDATED` | Cyan | Successfully fast-forwarded or merged new commits |
| `BEHIND (N)` | Yellow | Remote contains N commits ahead of local HEAD |
| `DIVERGED` | Magenta | Local and remote have distinct commit histories |
| `FAILED` | Red | Non-zero git exit code (conflict, network, auth) |

---

## Acceptance Criteria

### Scenario 1: Clean Fast-Forward Pull
- **Given** a workspace with 10 repositories all cleanly tracking remotes
- **When** `gitmap pull` is executed
- **Then** each repository is pulled concurrently, and the terminal renders the table with green `UP_TO_DATE` or cyan `UPDATED` status with duration timers.

### Scenario 2: Non-Unlink Collision Diagnosis
- **Given** local untracked files that conflict with incoming remote commits
- **When** `gitmap pull` encounters `error: Your local changes to the following files would be overwritten by merge`
- **Then** the diagnostic engine flags a `non-unlink git pull failure`, highlights the colliding files, and outputs the exact cleanup command:
  ```bash
  git stash push -u && gitmap pull && git stash pop
  ```

---

## Cross-References

- Pull Specification: [`101-pull-all.md`](./101-pull-all.md)
- Generic CLI Interface: [`02-cli-interface.md`](./02-cli-interface.md)
- Error Management: [`../03-error-manage/00-overview.md`](../03-error-manage/00-overview.md)
