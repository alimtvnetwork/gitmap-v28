---
name: commit-in
description: gitmap commit-in / cin replays commits from N source repos into one destination repo, dedupes by source SHA, replicates author+committer dates, profiles in .gitmap/commit-in/profiles/. Spec at spec/03-commit-in/.
type: feature
---
# commit-in / cin

## What it does
`gitmap commit-in <source> <inputs...>` replays every commit from one
or more input repos (folders or Git URLs) into the source repo, in the
order the inputs were listed and the order their commits originally
happened, while replicating original author + committer dates,
deduping by source SHA, and applying user-defined message / file /
author rules drawn from a saved `Profile`.

## Hard rules (apply to every implementation phase)

- **Source auto-init precedence (no flag, no prompt):**
  URL → `git clone` into `CWD/<basename>`. Existing repo → reuse.
  Existing non-repo folder → `git init` in place. Missing path →
  `mkdir -p && git init`.
- **DB convention:** PascalCase tables/columns/JSON keys/JSON values.
  Every PK is `INTEGER PRIMARY KEY AUTOINCREMENT` named `<TableName>Id`.
  Every classifier (Type/Status/Kind/Mode/Reason/Outcome/Stage/Source)
  is a Go enum AND a SQLite `(Id, Name UNIQUE)` mirror table.
- **Date replication:** Both `AuthorDate` AND `CommitterDate` of the
  rewritten commit equal the source commit's, byte-for-byte. Author
  identity may be overridden; dates may NEVER be.
- **Dedupe via `ShaMap`:** Same `SourceSha` ever seen again → SKIP +
  `SkipLog(DuplicateSourceSha)` with `PreviousRewrittenCommitId`.
- **No file content, no file hash** stored anywhere. Only `RelativePath`
  strings under `SourceCommitFile`.
- **Profile binding key:** Absolute symlink-resolved `SourceRepoPath`,
  never `origin` URL.
- **First-parent only walk:** Oldest → newest per input. Merge commits'
  second-parent history is NOT recursed.
- **Single advisory lock** (`<.gitmap>/gitmap.lock`) per workspace; a
  second concurrent `commit-in` exits `CommitInExitLockBusy`.
- **`--conflict Prompt` without VS Code:** hard-fail with
  `CommitInExitConflictAborted`; never silently downgrade to
  `ForceMerge`.
- **`all` / `-N` discovery scope:** parent directory of `<source>`;
  plain `<base>` is treated as `v0` and walked first.

## File system surface

- `<.gitmap>/db/gitmap.sqlite`         — SQLite DB (shared with rest of gitmap)
- `<.gitmap>/temp/<runId>/`            — per-run input clones
- `<.gitmap>/commit-in/profiles/<n>.json` — strict-decode JSON profile
- `<.gitmap>/logs/commit-in.log`       — run summary log

## Spec & plan pointers

- Spec: `spec/03-commit-in/` (7 iterations: overview, CLI surface,
  pipeline, DB schema + ERD, profiles + JSON, message + function-intel,
  acceptance + edge cases).
- Plan: `.lovable/plan.md` § "commit-in / cin — 2026-05-06" — 7 phased
  implementation steps, gated on the user typing `next`.
- Internal note: `.lovable/memory/issues/2026-05-06-commit-in-spec.md`
  (verbatim user prompt mirrored for traceability per spec §7.99).

## Status

All 7 gated phases complete (2026-05-06, v4.18.0):
1. Constants + typed enums.
2. DB migrations + 8 enum-mirror seeds (SchemaVersion 23→24).
3. Pure CLI parser (`commitin.Parse`).
4. Workspace + source resolution (`commitin/workspace`).
5. Walk + dedupe + replay + runlog (`commitin/{walk,dedupe,replay,runlog}`).
6. Profile JSON load/save + message pipeline + prompt
   (`commitin/{profile,message,prompt}`).
7. Function-intel detectors + finalize + dispatcher entry
   (`commitin/{funcintel,finalize}` + `cmd/commitin.go` + helptext).

Closeout (2026-05-06): full E2E harness + 9 pipeline tests live at
`gitmap/cmd/commitin/e2e/` (happy-path, dedupe, sibling `all`/`-N`,
auto-init, profile precedence, Prompt abort, ForceMerge clobber,
lock-busy). `// gitmap:cmd top-level` marker confirmed already-active
on the parent const block in `constants/constants_cli.go` (line 3 →
line 161); `go generate ./...` populated `commit-in` + `cin` into
`completion/allcommands_generated.go`. CHANGELOG v4.18.0 entry
prepended. Commit-in feature is release-ready.