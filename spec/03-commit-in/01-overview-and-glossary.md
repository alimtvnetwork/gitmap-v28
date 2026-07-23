# 01 — Overview and Glossary

## 1.1 Purpose (one sentence)

`gitmap-v27 commit-in` (`cin`) replays every commit from one or more **input
repos** (local folders or Git URLs) into a single **source repo**, in the
order the inputs were listed and the order their commits originally
happened, while replicating original author + committer dates, deduping
by source SHA, and applying user-defined message / file / author rules
drawn from a saved `Profile`.

## 1.2 Verbatim user intent (frozen)

The user's original prose is preserved unmodified in
`07-acceptance-and-edge-cases.md` §7.99 ("Verbatim Source"). Every rule
in this folder MUST be traceable to a sentence in that section. If a
later iteration disagrees with §7.99, §7.99 wins and the iteration is
wrong.

## 1.3 Glossary (canonical terms — use these, only these)

| Term                  | Meaning                                                                                              |
|-----------------------|------------------------------------------------------------------------------------------------------|
| **Source repo**       | First positional arg `<source>`. Destination of every rewritten commit.                              |
| **Input repo**        | Any positional after `<source>`. Provider of commits. Folder OR Git URL.                             |
| **Run**               | One invocation of `gitmap-v27 commit-in`. One row in `CommitInRun`.                                      |
| **Source SHA**        | Original commit SHA inside an input repo. Stable across runs.                                        |
| **New SHA**           | SHA assigned by `git commit` after replay into the source repo. Never equal to source SHA.           |
| **Profile**           | Saved bundle of answers (conflict mode, exclusions, message rules, …) bound to a `<source>` path.    |
| **Sandbox**           | `<.gitmap>/temp/<runId>/<orderIndex>-<basename>/` — where each input repo is cloned for the run.     |
| **Workspace root**    | Closest ancestor of CWD containing `.gitmap/` (or CWD if none). Anchors `temp/`, `db/`, `commit-in/`.|
| **Walk**              | Per-input traversal of commits, oldest → newest, first-parent only.                                  |
| **Replay**            | Materialize files at a source SHA into the source working tree, then `git commit` with original dates.|

## 1.4 Hard invariants (referenced from every other iteration)

- INV-01 — Commits walked **oldest → newest** per input, in the exact
  positional order the user supplied inputs.
- INV-02 — `AuthorDate` AND `CommitterDate` of the rewritten commit
  EQUAL the source commit's, byte-for-byte (RFC-3339 with timezone).
- INV-03 — File payload of the rewritten commit equals the source
  tree at that SHA, MINUS files matched by any `ProfileExclusion` of
  kind `PathFolder` or `PathFile`.
- INV-04 — Same `SourceSha` already present in `ShaMap` → SKIP; insert
  `SkipLog` row with `SkipReason = DuplicateSourceSha` and
  `PreviousRewrittenCommitId` filled.
- INV-05 — Source repo precedence (no flag, no prompt):
  `URL → clone into CWD/<basename>` ▸ `existing repo → reuse` ▸
  `existing non-repo folder → git init in place` ▸
  `missing path → mkdir -p && git init`.
- INV-06 — No file content, no blob hash, no diff payload is ever
  written to SQLite. Only `RelativePath` strings.
- INV-07 — Every classifier (`Type`, `Status`, `Kind`, `Category`,
  `Mode`, `Reason`, `Outcome`, `Stage`, `Source`) is an enum in code AND
  a mirror table in SQLite with `(Id INTEGER PK AI, Name TEXT UNIQUE)`.
- INV-08 — All JSON keys AND JSON enum values use PascalCase.
- INV-09 — `commit-in` never rewrites history of the source repo. It
  only appends new commits to the current branch tip.
- INV-10 — A single `gitmap.lock` (advisory file lock under `.gitmap/`)
  is held for the duration of the run. Concurrent `commit-in` invocations
  in the same workspace fail fast with exit `CommitInExitLockBusy`.

## 1.5 Out of scope for v1 (explicit non-goals)

1. Renamed / moved function detection in `--function-intel`.
2. Cross-language function detection in a single file (one language per file).
3. Pushing the source repo to any remote. `commit-in` only writes locally.
4. Rewriting any pre-existing commit in the source repo.
5. Submodules — input repos with submodules are walked at the superproject
   level only; submodule commits are not recursed into.
6. LFS pointer expansion — pointers are copied verbatim as files.

## 1.6 Success criterion (single sentence the agent can grep for)

After a successful run, every non-skipped source commit across all input
repos has exactly one `RewrittenCommit` row with a populated `NewSha`,
the source repo's git log shows those commits in walk order with
identical `AuthorDate` and `CommitterDate`, and `ShaMap` contains
`(SourceSha → RewrittenCommitId)` for each.