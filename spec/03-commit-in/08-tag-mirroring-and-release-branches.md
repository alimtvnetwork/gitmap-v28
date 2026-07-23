# 08 — Tag Mirroring & Auto Release Branches

> Status: spec frozen. Implementation forbidden until plan greenlights.
> Builds on §03 (pipeline), §04 (database schema). Read those first.

When `commit-in` walks a source repo, every source commit may have one or
more git tags pointing at it. v1 of `commit-in` ignored tags entirely —
the rewritten history was tagless even when the source carried `v1.2.3`.
This iteration closes that gap and adds an automatic release-branch
side-effect when the mirrored tag looks like a semantic version.

The destination commit SHA is by definition different from the source
SHA (different parent, different tree if files were excluded, different
committer date pin still hashes differently in some edge cases), so we
can never simply re-point the tag — we must **re-create** the tag on the
destination and **persist the (oldSha, newSha, tagName, releaseBranch)
relationship** in SQLite for downstream queries.

---

## 8.1 Vocabulary

| Term                | Meaning                                                                 |
|---------------------|-------------------------------------------------------------------------|
| Source tag          | An **annotated** git tag (`git tag -a`) on the source repo. Lightweight tags are explicitly ignored in v1 — they are local bookmarks, not releases. |
| Mirrored tag        | The destination tag with the same `Name` as the source tag, but pointing at the new (replayed) SHA. Carries the source tag's tagger identity, date, and message verbatim. |
| Version tag         | A mirrored tag whose `Name` matches the canonical semver regex (§8.3). Triggers the auto release-branch side-effect. |
| Auto release branch | A destination branch named `release/<TagName>` (always `release/` prefix, even if the tag itself has no `v`). Created from the same new SHA the version tag points at. Created **on by default**; suppressed by `--no-release-branch`. |

The `release/<…>` naming reuses `constants.ReleaseBranchPrefix` so the
branches are interchangeable with `gitmap-v27 release-branch` and friends.

---

## 8.2 CLI surface (additions to §02)

Three new flags on `commit-in` / `cin`:

| Flag                       | Default       | Description                                                                                                       |
|----------------------------|---------------|-------------------------------------------------------------------------------------------------------------------|
| `--tags <mode>`            | `Annotated`   | `Annotated` (mirror only annotated tags), `All` (mirror lightweight too), `None` (skip tag mirroring entirely).   |
| `--no-release-branch`      | off           | Suppress auto release-branch creation even when a mirrored tag is a version tag.                                  |
| `--release-branch-prefix`  | `release/`    | Override the prefix. Must end with `/`. Empty string is a validation error (would shadow the tag name as branch). |

Validation rules (§02 conformance):

- `--tags=None` + any other tag-mirroring flag set → exit code `2`
  (`BadArgs`) with message `commit-in: --tags=None forbids --no-release-branch / --release-branch-prefix`.
- `--release-branch-prefix` not ending in `/` → exit `2`.

No new exit codes; tag/branch failures during a run are classified as
`PartiallyFailed` (exit `1`) — the replay commit itself is `Created`,
only the side-effect failed.

Auto-completion: `--tags` exposes the three literal values; the other
two flags are plain string/bool. Reuses the v3.0.0 marker-comment
opt-in mechanism — see `mem://features/marker-comments`.

---

## 8.3 Version-tag detection (canonical regex)

The mirrored tag is a "version tag" iff its `Name` matches:

```
^v?\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$
```

This is exactly the regex used by `gitmap-v27 release` for release-branch
discovery (centralized in `constants.VersionTagPattern` — implementation
must reuse that constant; do not introduce a sibling).

Examples:

| Tag        | Version? | Auto release branch (default prefix) |
|------------|----------|--------------------------------------|
| `v1.2.3`   | yes      | `release/v1.2.3`                     |
| `1.2.3`    | yes      | `release/1.2.3`                      |
| `v1.2.3-rc1` | yes    | `release/v1.2.3-rc1`                 |
| `v1.2`     | no       | (mirror tag only, no branch)         |
| `nightly`  | no       | (mirror tag only, no branch)         |
| `v1.2.3+build.42` | yes | `release/v1.2.3+build.42`         |

---

## 8.4 Pipeline placement (extends §03)

The work happens **inline, per-commit**, inside stage 14 (`Commit`),
immediately AFTER `replay.ApplyCommit` returns the new SHA and BEFORE
the next source commit is walked. Rationale: keeps the spec a single
linear stage list, makes partial-failure semantics match real release
workflows (a tag pointing at a valid commit is fine even if a later
commit fails), and lets the SQLite write be a single transaction
together with the `RewrittenCommit` insert.

```
stage 14 — Commit (extended)
  for each plan in plans:
    result   = replay.ApplyCommit(plan, dryRun)            ← unchanged
    rewId    = store.InsertRewrittenCommit(result, plan)   ← unchanged
    tags     = gitutil.AnnotatedTagsAt(srcRepo, plan.SrcSha)
    for each tag in tags filtered by --tags mode:
      gitutil.CreateAnnotatedTag(dstRepo, tag.Name, result.NewSha,
                                 tag.TaggerIdent, tag.TaggerDate,
                                 tag.Message)
      mirroredBranch = ""
      if not --no-release-branch and isVersionTag(tag.Name):
        mirroredBranch = prefix + tag.Name
        gitutil.CreateBranchAt(dstRepo, mirroredBranch, result.NewSha)
      store.UpdateMirroredTagAndBranch(rewId, tag.Name, mirroredBranch)
```

`--dry-run` short-circuits exactly like ApplyCommit: tag detection still
runs (so the run-log can report what *would* have been mirrored), but
neither `git tag` nor `git branch` is invoked, and the SQLite columns
are left NULL (a `--dry-run` row by definition has no NewSha to point a
tag at).

---

## 8.5 Database schema (extends §04 — `RewrittenCommit`)

Two columns added to the existing `RewrittenCommit` table. NO new
table. This matches the "three relationships in one row" picture: a
single row carries `SourceCommitId`-derived old SHA, its own NewSha,
the mirrored tag name, and the auto release branch.

```
ALTER TABLE RewrittenCommit ADD COLUMN MirroredTagName       TEXT NULL;
ALTER TABLE RewrittenCommit ADD COLUMN MirroredReleaseBranch TEXT NULL;
CREATE INDEX IF NOT EXISTS IX_RewrittenCommit_MirroredTagName
    ON RewrittenCommit(MirroredTagName);
CREATE INDEX IF NOT EXISTS IX_RewrittenCommit_MirroredReleaseBranch
    ON RewrittenCommit(MirroredReleaseBranch);
```

| Column                  | Type      | Notes                                                                                              |
|-------------------------|-----------|----------------------------------------------------------------------------------------------------|
| `MirroredTagName`       | TEXT NULL | The tag that was created on the destination. NULL when source commit was tagless or `--tags=None`. |
| `MirroredReleaseBranch` | TEXT NULL | Branch name (e.g. `release/v1.2.3`). NULL when not a version tag, or `--no-release-branch`, or `--dry-run`. |

Migration: idempotent `ADD COLUMN`, file `006_commit_in_tag_mirroring.sql`.
Per §4.5, schema-change discipline requires a NEW migration number — we
do not edit migrations 001–005.

**N-tags-per-commit**: in the rare case a source commit carries multiple
annotated tags, ONE row per `(RewrittenCommitId, MirroredTagName)` pair
is written. The first tag uses the existing RewrittenCommit row; each
additional tag inserts a sibling row with the same `NewSha`,
`SourceCommitId`, `CommitOutcomeId = Skipped`, and a new
`SkipReason.AdditionalTagAlias` (added to the SkipReason enum). Most
repos have ≤1 tag per commit; this is the rarely-hit branch.

---

## 8.6 Profile JSON (extends §05)

Three new optional keys on the profile root, all PascalCase:

```json
{
  "TagsMode": "Annotated",
  "CreateReleaseBranch": true,
  "ReleaseBranchPrefix": "release/"
}
```

Resolution order (matches §05 precedence): CLI flag > profile JSON >
spec default.

`TagsMode` accepts the same `Annotated|All|None` enum as `--tags`. The
enum is mirrored in SQLite as `TagsMode (TagsModeId, Name UNIQUE)` and
seeded in migration 006.

---

## 8.7 Acceptance matrix (extends §07)

| # | Setup                                                | Expected                                                              |
|---|------------------------------------------------------|-----------------------------------------------------------------------|
| T1 | Source commit has annotated tag `v1.2.3`            | Dest gets tag `v1.2.3` at NewSha + branch `release/v1.2.3` at NewSha. RewrittenCommit row has both columns set. |
| T2 | Source commit has lightweight tag `bookmark`         | Default `--tags=Annotated` → ignored. Both columns NULL. With `--tags=All` → tag mirrored, no branch (not semver). |
| T3 | Source commit has annotated tag `nightly`            | Tag mirrored, no branch. `MirroredReleaseBranch` NULL. |
| T4 | `--no-release-branch` + annotated tag `v2.0.0`       | Tag mirrored, NO branch. `MirroredReleaseBranch` NULL. |
| T5 | `--tags=None`                                        | No tag walk, no branch creation, both columns NULL on every row. |
| T6 | `--dry-run` + annotated tag `v1.0.0`                 | No git mutations. Run-log says "would mirror tag v1.0.0, would create release/v1.0.0". RewrittenCommit row written with NewSha NULL (existing rule) and both mirror columns NULL. |
| T7 | Source commit has TWO annotated tags `v1.0.0` and `release-1.0` | Two rows: first uses the canonical RewrittenCommit row (`v1.0.0` + `release/v1.0.0`), second is `Skipped/AdditionalTagAlias` for `release-1.0`. |
| T8 | Branch `release/v1.2.3` already exists on dest       | `git branch release/v1.2.3` fails. Run is `PartiallyFailed` (exit 1). The replay commit row is `Created`, branch column is left NULL, error logged via the standard `os.Stderr` format (zero-swallow rule). |
| T9 | Tag `v1.2.3` already exists on dest with different SHA | `git tag` fails. Same `PartiallyFailed` semantics — never silently overwrite an existing tag. To force, user must delete the dest tag first. |
| T10 | `--release-branch-prefix=releases/`                 | Branch becomes `releases/v1.2.3`. Indexed exactly the same way; existing `release/*` tooling will not pick it up — this is on the user. |

---

## 8.8 Cross-cutting rules

- **Zero-swallow** (Core memory): every `git tag` / `git branch` failure
  is logged to `os.Stderr` in the standardized format with the source
  SHA, dest SHA, tag name, and the underlying git stderr verbatim.
  `errors.Is` is used to recognize `gitutil.ErrTagExists` /
  `gitutil.ErrBranchExists` so the operator gets a one-line "already
  exists, skipping" instead of a stack-trace-shaped error.
- **No magic strings** (Core memory): all flag names, enum values,
  prefix defaults, regex literal — every one of them lives in
  `constants_commit_in_tags.go` (CLI strings in `constants_cli.go`
  per the global rule); the regex is reused from
  `constants.VersionTagPattern`.
- **`<200` lines per file**: the new package
  `gitmap-v27/cmd/commitin/replay/tags.go` holds the orchestration only;
  git plumbing lives in `gitmap-v27/gitutil/tags.go` (annotated-tag
  read/write helpers).
- **Determinism pre-check** (memory `mem://features/determinism-precheck`):
  any new golden fixture for tag-mirroring runs through
  `goldenguard.AssertWriterDeterministic` before the env-var gate.
- **Faithful banner** (memory `mem://features/verify-cmd-faithful-banner`):
  tests asserting mirror-failure paths must call
  `PrintCmdFaithfulReportForTest` so the `[FAIL]`-prefixed banner is
  emitted in the test output.
