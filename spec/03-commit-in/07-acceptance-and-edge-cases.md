# 07 — Acceptance Criteria, Edge Cases, Verbatim Source

## 7.1 Acceptance matrix (numbered, AI-blind testable)

Each AC is a statement an implementer's test suite MUST verify with
exactly one test case. Do not collapse, do not split.

1. `gitmap-v27 commit-in <s> a b c`, `gitmap-v27 commit-in <s> a,b,c`,
   `gitmap-v27 commit-in <s> "a, b, c"`, and
   `gitmap-v27 commit-in <s> "a"  "b"   "c"` produce identical
   `InputRepo.OrderIndex` rows.
2. `<source>` URL is cloned into `CWD/<basename>`; `<source>` empty
   folder is `git init`-ed; `<source>` non-empty non-repo folder is
   `git init`-ed in place; `<source>` missing path is `mkdir -p`'d
   then `git init`-ed. `WasSourceFreshlyInit` is `1` for the last
   three cases.
3. Each input cloned to `<.gitmap>/temp/<runId>/<orderIndex>-<basename>/`
   and recorded in `InputRepo` with the right `InputKind`.
4. `all` from source basename `gitmap-v27` expands to `gitmap-v27, gitmap-v27,
   gitmap-v27, …, gitmap-v27` (ascending; plain `gitmap-v27` first as `v0`).
   `-5` returns the last five sibling versions only.
5. SQLite schema matches §04 byte-for-byte; every PK is
   `INTEGER PRIMARY KEY AUTOINCREMENT` named `<TableName>Id`; every
   classifier is a join table.
6. Per-input commits walked first-parent only, oldest → newest;
   `SourceCommit` and `SourceCommitFile` rows hold zero file content
   and zero file hash columns.
7. Re-encountering a `SourceSha` already in `ShaMap` produces exactly
   one `SkipLog` row with `SkipReason = DuplicateSourceSha` and
   `PreviousRewrittenCommitId` populated.
8. `RewrittenCommit.AppliedAuthorDate == SourceCommit.AuthorDate` AND
   `AppliedCommitterDate == SourceCommit.CommitterDate`, byte-for-byte.
9. Author identity precedence honored per §6.6.
10. `--conflict ForceMerge` keeps the file from the source commit with
    the latest `AuthorDate` when two inputs touch the same path.
    `--conflict Prompt` invokes `code --wait --diff <ours> <theirs>`;
    if `code` is missing, exit `CommitInExitConflictAborted`. (No silent
    fallback to `ForceMerge` — resolves Ambiguity #4.)
11. `Exclusions[Kind=PathFolder]` drops every file whose `RelativePath`
    starts with `Value + "/"`. `Exclusions[Kind=PathFile]` drops only
    the exact match. Already-tracked files in `<source>` from prior
    commits stay untouched (resolves Ambiguity #3).
12. `MessageRules` apply line-level per §6.1 stage 1.
13. `OverrideMessages` random pick is deterministic per §3.4 PRNG seed.
14. Title / body prefix / suffix order: `TitlePrefix + line + TitleSuffix`
    on first line; `MessagePrefix` prepended new line above body;
    `MessageSuffix` appended new line below body.
15. Function-intel block emitted only when `IsEnabled` AND at least one
    detector returned ≥ 1 name.
16. `--save-profile <n>` persists JSON file AND `Profile`+children
    rows in one transaction; `--save-profile-overwrite` required to
    replace; `--set-default` flips `IsDefault` atomically.
17. `ShaMap` row created exactly once per `RewrittenCommit` with
    `Outcome = Created`; never for `Failed` or `Skipped`.
18. `CommitInRun.RunStatus` ends as `Completed` (zero failures),
    `PartiallyFailed` (≥1 commit `Failed`), or `Failed` (global stage
    failure). Summary line on STDOUT matches §2.8.
19. Every enum in §02–§06 lives in its own Go file under
    `gitmap-v27/cmd/commitin/<enumname>.go`, used by both production code
    and tests; mirror table seeded by the same migration that creates
    it (§4.5).
20. Coding guidelines honored: ≤ 8-line functions, ≤ 100-line files,
    no `any` / `interface{}` in new code, no nested `if`, `is/has`
    boolean prefix, every error path logs to `os.Stderr` per
    `mem://tech/code-red-error-management`.

## 7.2 Edge case catalogue

- **Empty input repo** (zero commits) — `WalkCommits` yields nothing;
  `InputRepo` row exists; no error.
- **Detached HEAD in input** — walk the commit graph reachable from
  `HEAD`; same first-parent-only rule.
- **Commit with zero file changes (empty tree diff)** — still walked,
  still produces a `SourceCommit` row; replay uses
  `git commit --allow-empty`.
- **Symlinks in source tree** — replicated as symlinks (no follow).
- **Files outside the working tree** (worktrees, sparse checkout) — out
  of scope; document as "use a fresh full clone".
- **Case-only renames** on case-insensitive FS (macOS default) — drop
  the second file with a WARN; record as `Failed` for that commit.
- **`<source>` is on a different filesystem than `<.gitmap>/temp/`** —
  copy via `io.Copy`; never assume rename works across filesystems.
- **CRLF / LF normalization** — disable git autocrlf for the run
  (`-c core.autocrlf=false`) so source bytes are preserved.
- **`.gitignore` in source repo accidentally hides input file** — bypass
  via `git add -f`. Log a WARN once per commit when this triggers.
- **Concurrent commit-in** — `gitmap.lock` advisory lock; second
  invocation exits `CommitInExitLockBusy` immediately.

## 7.3 Ambiguity resolutions (closes the original prompt's §"Ambiguities Flagged")

| # | Question                                                                 | Resolution                                                               |
|---|--------------------------------------------------------------------------|--------------------------------------------------------------------------|
| 1 | `all` / `-N` discovery scope                                             | Parent directory of `<source>`. (§2.4)                                   |
| 2 | Where does plain `<base>` sort vs `vN`?                                  | Plain `<base>` is treated as `v0` and walked FIRST. (§2.4)               |
| 3 | `--exclude` applied where?                                               | Per-commit BEFORE staging; existing tracked files untouched. (AC #11)    |
| 4 | `Prompt` mode without VS Code on PATH?                                   | Hard fail with `CommitInExitConflictAborted`. No silent downgrade. (AC #10)|
| 5 | Renamed / moved function detection in v1?                                | Out of scope. (§1.5 #1)                                                  |
| 6 | Source repo already has commits sharing a SHA mapping with input?        | Skip via `ShaMap` (idempotency contract §3.5).                           |
| 7 | Profile binding key: path or origin URL?                                 | Absolute symlink-resolved path. (§5.4)                                   |
| 8 | `--save-profile` overwrite default?                                      | Refuse unless `--save-profile-overwrite`. (§5.5)                         |

## 7.4 Conformance test plan (informative)

Implementation MUST ship integration tests under
`gitmap-v27/cmd/commitin/integration_test.go` that build a synthetic input
repo via `git plumbing` (no `git add`, mirror the pattern in
`.github/scripts/smoke-history-pin.sh`), run `commit-in`, and assert
each AC #1–#20 with one `t.Run` per AC.

## 7.99 Verbatim source (frozen — never edit)

The user's original prose is preserved here verbatim for traceability.
If any later iteration disagrees with this section, this section wins.

> See `.lovable/memory/issues/2026-05-06-commit-in-spec.md` for the
> full original prompt (mirrored under the project memory tree to keep
> spec files free of multi-page prose blocks).