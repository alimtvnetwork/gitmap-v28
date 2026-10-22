# Plan 299: Result Wrapper Store SSH, Database Queries, Templates & Visibility Monadic Results

**Plan ID:** `299-result-wrapper-store-ssh-and-templates`
**Status:** `completed`
**Budget:** 400 steps
**Created:** 2026-10-23
**Completed:** 2026-10-23
**Spec Reference:** [02-spec/21-app/299-result-wrapper-store-ssh-and-templates/01-architecture-spec.md](../../02-spec/21-app/299-result-wrapper-store-ssh-and-templates/01-architecture-spec.md)

---

## 1. Context & Motivation
Following the completion of Plans 296, 297, and 298, Plan 299 continues the migration of Go functions returning multi-value error tuples (`([]T, error)`, `(T, error)`) to canonical single-return monadic containers (`result.WrapSlice[T]`, `result.Wrap[T]`).

This wave focused on database query repositories in `cli/store/` (SSH hosts, SSH keys, known hosts, command history, tasks, aliases, bookmarks, releases, workdirs), template discovery in `cli/templates/`, and pattern parsing in `cli/visibility/`.

---

## 2. Discrete Subtasks Breakdown & Outcomes
- **Subtask 01 (Store SSH & Task History):** `[DONE]`
  - Added monadic `*Wrap` methods for `ListHosts`, `ListHostsByRole`, `ListHostsByTarget`, `ListSSHHosts` in `cli/store/ssh_repo.go` and `cli/store/ssh_target_resolver.go`.
  - Added monadic `*Wrap` methods for `ListSSHKeys`, `SSHKeyNames` in `cli/store/sshkey.go`.
  - Added monadic `*Wrap` methods for `ListSSHKnownHosts` in `cli/store/ssh_known_hosts_repo.go`.
  - Added monadic `*Wrap` methods for `ListSSHHistory` in `cli/store/ssh_hist_repo.go`.
  - Added monadic `*Wrap` methods for `ListTaskHistory` in `cli/store/tasks_split_db.go`.
  - Decomposed `cli/store/ssh_repo.go` into `ssh_target_resolver.go` and `ssh_enroll.go`, keeping all files under the 500-line cap.
- **Subtask 02 (Store Aliases, Bookmarks, Releases, WorkDirs & Types):** `[DONE]`
  - Centralized types in `cli/store/types.go` (`SSHHostSliceWrap`, `SSHKeySliceWrap`, `SSHKnownHostSliceWrap`, `SSHHistorySliceWrap`, `TaskHistoryRecordSliceWrap`, `AliasSliceWrap`, `BookmarkRecordSliceWrap`, `ReleaseRecordSliceWrap`, `WorkDirSliceWrap`).
  - Added monadic `*Wrap` methods in `cli/store/alias.go`, `bookmark.go`, `release.go`, `workdir_list.go`.
- **Subtask 03 (Templates & Visibility):** `[DONE]`
  - Created `cli/templates/types.go` (`EntrySliceWrap`, `EntryWrap`).
  - Added `ListWrap()` in `cli/templates/list.go`.
  - Created `cli/visibility/types.go` (`PatternSliceWrap`, `PatternWrap`).
  - Added `ParsePatternListWrap()` in `cli/visibility/pattern.go`.
- **Subtask 04 (FastGate Validation & Atomic GitMap Commit):** `[DONE]`
  - Executed `python 03-ai-scripts/50-fastgate.py`: 100% PASS in 0.034s across all candidate files.
  - Finalized with atomic GitMap commit.

---

## 3. Disjoint File Ownership
- **Worker 01:** `cli/store/ssh_repo.go`, `cli/store/sshkey.go`, `cli/store/ssh_known_hosts_repo.go`, `cli/store/ssh_hist_repo.go`, `cli/store/tasks_split_db.go`
- **Worker 02:** `cli/store/types.go`, `cli/store/alias.go`, `cli/store/bookmark.go`, `cli/store/release.go`, `cli/store/workdir_list.go`
- **Worker 03:** `cli/templates/*`, `cli/visibility/*`
- **Lead Agent:** `.ai-memory/plans/*`, `02-spec/*`, `.ai-memory/temp-agents/*`

> **Closed 2026-10-23 (task 190, subtask 06 drain):** plan complete per its own status record and the spec 02 audit (substantively complete with matching commits); moved from pending/ to completed/ so the pending folder reflects only genuinely open work.
