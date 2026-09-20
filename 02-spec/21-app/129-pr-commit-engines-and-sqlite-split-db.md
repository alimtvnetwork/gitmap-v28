# Spec 129: PR Commit Engines, SQLite Split-DB Architecture, and Auto-Merge PR Release Suite

> **Status:** Active  
> **Package:** `cli/committransfer`, `cli/prdb`, `cli/store`, `cli/cmd`, `cli/cmd/commitin`  
> **Related Specs:** [Spec 106](106-commit-left-right-both.md), [Spec 114](114-committransfer-idempotence-and-merge-default.md), [Spec 124](124-polyglot-worker-orchestrator-and-automation-runner.md), [Spec 04](../../02-spec/04-database-conventions/01-naming-conventions.md)

---

## 1. Executive Overview & Objectives

Spec 129 specifies the **PR Commit Engine Family**, the canonical **SQLite Split-DB Directory Standardization**, and the **Deterministic Final Snapshot Synchronization** for GitMap CLI.

The core objectives are:
1. **First-Class `pr` Command Family:** Introduce `gitmap pr` (and `gitmap pull-request`) supporting directional replay (`LEFT RIGHT`), multi-source replay (`in <source> <inputs...>`), branch cleanup (`pr clean` / `pr rm`), and audit listing (`pr list`).
2. **Canonical Split-DB Path Formula:** Enforce `.gitmap/data/<section>/<slug>/sql.db` across all GitMap subsystems (`automation`, `pipeline`, `pr`, `logs`, `installer`, `repo_search`, `schedules`, `startup`, `sites`), deprecating and auto-migrating all legacy fragmented paths.
3. **Automated PR & Feature Branch Engine:** Detect merge commits in the replayed commit stream, automatically fork feature branches (`feature/<slug>` or `pr/<id>`), generate comprehensive Markdown PR descriptions, merge PRs into the mainline target, and record the entire lifecycle into SQLite.
4. **Rich Terminal UI & Visual Execution Graph:** Provide modern, colorized terminal UI with `pterm` headers, progress bars, status badges, and an ASCII execution graph illustrating mainline commits, feature branches, and merge nodes.
5. **Deterministic Final Snapshot Synchronization:** Guarantee that upon completion of `commit-in`, `commit-right`, `commit-left`, or `pr`, the target working tree exactly matches the source repo snapshot byte-for-byte and file-for-file.
6. **PR Branch Pruning (`pr-clean`):** Provide interactive and non-interactive pre-flight inspection, confirmation, and safe deletion of merged and closed PR branches.

---

## 2. Canonical Split-DB Architecture (`.gitmap/data/<section>/<slug>/sql.db`)

### 2.1 The Standard Path Formula

All SQLite databases in GitMap MUST adhere to:
```text
<BaseDataDir>/<section>/<slug>/sql.db
```
Where:
- `<BaseDataDir>` is `.gitmap/data` when inside a repository (or `store.BinaryDataDir()` in global mode).
- `<section>` is the functional category:
  - `automation`: Polyglot runtime probes and search exclusions.
  - `pipeline`: CI/CD pipeline history, logs, and stage traces.
  - `pr`: Pull requests, feature branches, and PR releases.
  - `logs`: System and execution logs.
  - `installer`: Installer scripts and definition state.
  - `repo_search`: Per-repository search index and cache.
  - `schedules`: Scheduled task records and cron history.
  - `startup`: System and user-login autostart execution records.
  - `sites`: Nginx virtual host configurations and domains.
- `<slug>` is the sanitized alphanumeric-and-hyphen repository or item slug (`SanitizeSlug(raw)`).
- `sql.db` is the constant filename for every database.

### 2.2 Transparent Auto-Migration

When `store.ResolveSplitDbPath(section, slug, repoRoot)` is invoked:
1. It computes the target canonical path: `<BaseDataDir>/<section>/<slug>/sql.db`.
2. If the canonical file does not exist, it inspects known legacy locations (e.g. `.gitmap/data/<slug>/<section>/sql.db`, `.gitmap/data/<section>/<slug>/pipeline.db`, `<BinaryDataDir>/<section>.db`).
3. If a legacy database exists, it automatically copies or renames it to the canonical path under an atomic filesystem lock before returning.

---

## 3. PR SQLite Storage Engine Schema (`data/pr/<slug>/sql.db`)

In compliance with Spec 04, the schema utilizes PascalCase table and column names, affirmative booleans, integer primary keys, UTC epoch timestamps, and backward-compatible snake_case views.

```sql
-- PullRequest Table
CREATE TABLE IF NOT EXISTS PullRequest (
    PullRequestId   INTEGER PRIMARY KEY AUTOINCREMENT,
    PrNumber        INTEGER NOT NULL,
    Title           TEXT NOT NULL,
    Description     TEXT NOT NULL DEFAULT '',
    SourceBranch    TEXT NOT NULL,
    TargetBranch    TEXT NOT NULL DEFAULT 'main',
    Status          TEXT NOT NULL DEFAULT 'open',   -- 'open', 'merged', 'closed', 'draft'
    MergeCommitSha  TEXT NULL,
    CreatedAt       INTEGER NOT NULL,               -- UTC epoch seconds
    MergedAt        INTEGER NULL,                   -- UTC epoch seconds
    ClosedAt        INTEGER NULL,                   -- UTC epoch seconds
    Notes           TEXT NULL,
    Comments        TEXT NULL,
    UpdatedAt       INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS IdxPullRequest_PrNumber ON PullRequest (PrNumber);
CREATE INDEX IF NOT EXISTS IdxPullRequest_Status ON PullRequest (Status);
CREATE INDEX IF NOT EXISTS IdxPullRequest_SourceBranch ON PullRequest (SourceBranch);
CREATE INDEX IF NOT EXISTS IdxPullRequest_TargetBranch ON PullRequest (TargetBranch);

-- PrRelease Table
CREATE TABLE IF NOT EXISTS PrRelease (
    PrReleaseId     INTEGER PRIMARY KEY AUTOINCREMENT,
    PullRequestId   INTEGER NOT NULL,
    ReleaseTag      TEXT NOT NULL,
    CommitSha       TEXT NOT NULL,
    Notes           TEXT NULL,
    Comments        TEXT NULL,
    CreatedAt       INTEGER NOT NULL,               -- UTC epoch seconds
    FOREIGN KEY (PullRequestId) REFERENCES PullRequest(PullRequestId) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS IdxPrRelease_PullRequestId ON PrRelease (PullRequestId);
CREATE INDEX IF NOT EXISTS IdxPrRelease_ReleaseTag ON PrRelease (ReleaseTag);

-- PrBranch Table
CREATE TABLE IF NOT EXISTS PrBranch (
    PrBranchId      INTEGER PRIMARY KEY AUTOINCREMENT,
    BranchName      TEXT NOT NULL UNIQUE,
    BranchType      TEXT NOT NULL DEFAULT 'feature', -- 'feature', 'bugfix', 'hotfix', 'release', 'pr'
    IsMerged        INTEGER NOT NULL DEFAULT 0,      -- 1 = merged, 0 = active
    IsDeleted       INTEGER NOT NULL DEFAULT 0,      -- 1 = deleted/pruned, 0 = exists
    CreatedAt       INTEGER NOT NULL,                -- UTC epoch seconds
    MergedAt        INTEGER NULL,                    -- UTC epoch seconds
    DeletedAt       INTEGER NULL,                    -- UTC epoch seconds
    Notes           TEXT NULL,
    Comments        TEXT NULL,
    UpdatedAt       INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS IdxPrBranch_BranchName ON PrBranch (BranchName);
CREATE INDEX IF NOT EXISTS IdxPrBranch_IsMerged ON PrBranch (IsMerged);
CREATE INDEX IF NOT EXISTS IdxPrBranch_IsDeleted ON PrBranch (IsDeleted);

-- Backward-Compatible Views
CREATE VIEW IF NOT EXISTS pull_requests AS
SELECT PullRequestId AS id, PrNumber AS pr_number, Title AS title, Description AS description,
       SourceBranch AS source_branch, TargetBranch AS target_branch, Status AS status,
       MergeCommitSha AS merge_commit_sha, CreatedAt AS created_at, MergedAt AS merged_at,
       ClosedAt AS closed_at, Notes AS notes, Comments AS comments, UpdatedAt AS updated_at
FROM PullRequest;

CREATE VIEW IF NOT EXISTS pr_releases AS
SELECT PrReleaseId AS id, PullRequestId AS pr_id, ReleaseTag AS release_tag,
       CommitSha AS commit_sha, Notes AS notes, Comments AS comments, CreatedAt AS created_at
FROM PrRelease;

CREATE VIEW IF NOT EXISTS pr_branches AS
SELECT PrBranchId AS id, BranchName AS branch_name, BranchType AS branch_type,
       IsMerged AS is_merged, IsDeleted AS is_deleted, CreatedAt AS created_at,
       MergedAt AS merged_at, DeletedAt AS deleted_at, Notes AS notes, Comments AS comments,
       UpdatedAt AS updated_at
FROM PrBranch;
```

---

## 4. PR Commit Replay & Feature Branch Workflow

### 4.1 Merge Detection
During commit walk and plan hydration:
1. For each commit, inspect parent count: `git rev-parse <sha>^@`.
2. If parent count > 1, commit `M` is a merge commit with parents `P1` (target mainline) and `P2` (incoming branch).
3. The branch name is parsed from the commit subject (`Merge pull request #<num> from <branch>`) or generated as `pr/<shortsha>-<slug>`.

### 4.2 Feature Branch Lifecycle
1. **Branch Creation:** Fork a local feature branch at the target's `P1` state: `git checkout -b <branchName> <target-P1>`.
2. **Replay Ingress Commits:** Replay all commits exclusive to `P2` onto the feature branch with conventional message normalization and author preservation.
3. **PR Description Generation:** Compute diffs, component impacts, and commit summaries; generate structured Markdown.
4. **Database Record:** Insert record into `PullRequest` (status: `open`) and `PrBranch` (status: active).
5. **Mainline Merge:** Switch back to mainline (`git checkout <main>`), execute `git merge --no-ff <branchName> -m "<PR Merge Message>"`.
6. **Close PR:** Update `PullRequest` (status: `merged`, `MergedAt = now`) and `PrBranch` (`IsMerged = 1`).

---

## 5. Deterministic Final Snapshot Synchronization

At the end of `commit-in`, `commit-right`, `commit-left`, and `pr`:
1. The target working tree is mirror-pruned against the source repository's latest commit/release tree:
   - Target files not present in the source are removed.
   - All source files are copied into the target.
   - Ignored files (`.git/`, `node_modules/`, `.gitmap/data/`) are preserved.
2. `git add -A` stages any remaining differences.
3. If staged changes exist, a final synchronization commit is created:
   ```text
   chore(sync): synchronize final repository snapshot tree to match source <shortsha>

   gitmap-replay-final-snapshot: <sourceHeadSha>
   gitmap-replay-cmd: <commandName>
   ```
4. This ensures byte-for-byte and file-for-file equality with the latest commit and release.

---

## 6. PR Branch Cleanup Command (`pr-clean` / `pr-rm`)

1. **Pre-flight Check:** Queries `data/pr/<slug>/sql.db` for branches where `IsMerged = 1 AND IsDeleted = 0`.
2. **Pre-flight UI:** Renders a table showing PR number, branch name, merge SHA, and merged timestamp.
3. **Interactive Confirmation:** Prompts `Remove N closed/merged PR branches? [y/N]`. If `-y` / `--yes` is passed, confirmation is skipped.
4. **Pruning Execution:** Deletes local git branches via `git branch -D <branchName>` and updates SQLite `IsDeleted = 1`, `DeletedAt = now()`.

---

## 7. Command Grammar & CLI Options

```text
gitmap pr LEFT RIGHT [flags]
gitmap pr in <source> <inputs...> [flags]
gitmap pr clean [repo] [--yes]
gitmap pr rm [repo] [--yes]
gitmap pr list [repo] [--json]

Flags:
  --pr <mode>            PR mode: merges (default), all, tags, release, off
  --since <ref|date>     Replay commits since ref or date
  --limit <n>            Limit number of commits to replay
  --mirror               Prune files in target that do not exist in source
  --dry-run              Preview actions without writing commits or branches
  -y, --yes              Skip all confirmation prompts
  --json                 Output results as JSON
```
