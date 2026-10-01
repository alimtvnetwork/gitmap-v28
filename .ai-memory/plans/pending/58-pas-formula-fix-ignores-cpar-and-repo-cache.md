# Plan 58: GitMap PAS Formula, Windows Credential Store Fix, Ignore Management Suite, CPAR, and Cache Engine

## Metadata
- **Plan ID**: `58`
- **Spec Reference**: `02-spec/21-app/195-pas-formula-fix-ignores-cpar-and-repo-cache.md`
- **RCA Reference**: `02-spec/22-app-issues/58-windows-git-cache-credential-store-failure-rca.md` / `.ai-memory/issues/24-windows-git-cache-credential-store-failure-rca.md`
- **Status**: In Progress

---

## Subtasks & Orchestration

### Subtask 1: Windows Git Subprocess Credential Store Fix
- Path: `.ai-memory/plans/subtasks/201-pas-formula-ignore-suite-cpar-and-cache/01-windows-git-cache-credential-store-fix.md`
- Strip `GCM_CREDENTIAL_STORE=cache` on Windows across `cli/cloner/cloner.go`, `cli/cloner/safe_pull.go`, `cli/cloner/safe_push.go`, `cli/cmdclone/clonepretty.go`, and `cli/cmdpull/pull.go`.
- Provide centralized `gitutil.BuildSafeGitEnv()` helper function.
- Verify Git for Windows executes without UNIX domain socket error.

### Subtask 2: Pull-First Workflow & GitIgnore Deduplication Sanitizer
- Path: `.ai-memory/plans/subtasks/201-pas-formula-ignore-suite-cpar-and-cache/02-pull-first-and-gitignore-deduplication.md`
- Remove upfront blocking `checkAgmResumeTaskBeforePull` and `checkAgmResumeTaskEfficient` before repo pulls.
- Execute repo pulls immediately; defer ignore checks to background or post-pull summary.
- Implement `.gitignore` deduplication sanitizer that removes duplicate lines while preserving comments and structure.
- Add default ignore rules for `.gitmap/backup/`, `.antigravity_resume_task.json`, and resume task files.

### Subtask 3: Ignore Management Suite & PAS Fleet Extensions
- Path: `.ai-memory/plans/subtasks/201-pas-formula-ignore-suite-cpar-and-cache/03-ignore-management-suite-and-pas-formula.md`
- Implement `gitmap fix ignore all [-y]` (`fia`), `gitmap fix-ignore-all`, `gitmap fix ignores all ssh [-y]`, `gitmap fix-ignores-all-ssh (fias) [-y]` following the GitMap PAS Formula.
- Implement full `gitmap ignore` (`ig`) suite:
  - `add`, `scan`, `scan-ssh` (`ss`), `remove`, `edit`, `action`, `ls`, `help`, `ui`, `app`.
  - Grouping: `add-group`, `remove-group` (`rm-grp`), `set-default-group`, `add-grp-to-default` (`agtd`), `apply`, `connect-group-with-repo` (`cgwp`) `--add-with-default` (`awd`), `export`, `import`.

### Subtask 4: CPAR, Observability Suite (`see`), and Split-DB Repo Cache Engine
- Path: `.ai-memory/plans/subtasks/201-pas-formula-ignore-suite-cpar-and-cache/04-cpar-see-suite-and-split-db-cache.md`
- Implement `gitmap commit-push-all-repos` (`cpar`) with `-y`, `--review` (`-r`), `--review --commit-only` (`-co`).
- Implement `gitmap see` suite: `commit pending`, `git-ignore`/`ig issues`, `errors`, `history`, `errors ssh` (`ses`), `repo-manage ui`.
- Implement `gitmap cache` suite: `create`, `ls`, `add`, `remove` (`rm`), `help`, `search`, `search-multi`, `search-multi-grep`, `recache`/`reconcile`/`sync` using Split-DB SQLite storage.
- Minor release bump and CI/CD verification.
