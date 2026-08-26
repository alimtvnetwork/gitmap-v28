# Repo Major-Version Migrator (v14 → v15 → v16)

**Version:** 1.0.0
**Status:** Spec — awaiting approval. Implementation BLOCKED until explicit "go".
**Companion:** `25-release-pinned-installer.md` (end-user pinned installer — different tool, different purpose).

---

## Important

- Do not act. Spec only. Implementation begins after explicit approval.
- This script CHANGES repo identity across the codebase. It is destructive by design.
- Never run automatically on `main` without an explicit `--confirm` flag.
- Pinned installers (spec 25) and the generic `install.{ps1,sh}` are out of scope here.

---

## Problem Statement

When a major version is cut, the repo is renamed (e.g. `coding-guidelines-v24` → `coding-guidelines-v24`). Every URL, badge, install one-liner, `package.json` repo field, README link, and CI workflow that hardcodes the old name becomes stale. Manual replacement is error-prone and easy to miss in JSON, YAML, and markdown stamps.

The migrator is a one-shot CI/CD-runnable script that:

1. Discovers every reference to the old repo slug.
2. Rewrites them to the new repo slug.
3. Bumps `package.json` to the new major version.
4. Re-runs sync scripts so derived files (`version.json`, `specTree.json`, `health-score.json`, `readme.md` stamps) are consistent.
5. Produces a single diff/commit that can be reviewed in PR.

---

## Scope

1. New script: `scripts/migrate-repo-major-version.mjs` (Node, matches `sync-*.mjs` style).
2. New CI workflow: `.github/workflows/migrate-repo-version.yml` (manual `workflow_dispatch` only).
3. New spec: this file.
4. Untouched:
   - `release-install.{ps1,sh}` (spec 25 — end-user pinning).
   - Generic `install.{ps1,sh}` (current "always latest" behavior).
   - All bundle installers (`cli-install.*`, `slides-install.*`, etc.) — they are rewritten as data, not edited as logic.

---

## Behavioral Contract

### Inputs

| Flag | Required | Purpose |
|------|----------|---------|
| `--from <slug>` | yes | Old repo slug, e.g. `alimtvnetwork/coding-guidelines-v24`. |
| `--to <slug>` | yes | New repo slug, e.g. `alimtvnetwork/coding-guidelines-v24`. |
| `--new-version <semver>` | yes | New `package.json` version, e.g. `4.0.0`. Must be a major bump. |
| `--dry-run` | no (default ON) | Print planned changes; write nothing. |
| `--confirm` | no | Required to actually write files. Without this, dry-run wins. |
| `--help` | no | Print usage and exit 0. |

### Resolution Order

1. Parse args. Reject if any of `--from`, `--to`, `--new-version` are missing.
2. Validate slugs against `^[A-Za-z0-9-]+/[A-Za-z0-9._-]+$`.
3. Validate `--new-version` against `^\d+\.\d+\.\d+(-[A-Za-z0-9.]+)?$`.
4. Reject if old major == new major (must be a major bump).
5. Reject if `--from` == `--to`.
6. If `--confirm` not present, force dry-run regardless of other flags.

### What Gets Rewritten

The script scans these file globs and replaces exact matches of `--from` with `--to`:

- `readme.md`
- `package.json` (only the `repository.url`, `homepage`, `bugs.url` fields — never `name`)
- `version.json`
- `public/health-score.json`
- `src/data/specTree.json`
- `docs/**/*.md`
- `spec/**/*.md`
- `scripts/**/*.{mjs,js,sh,ps1}`
- `linter-scripts/**/*.{py,sh,ps1}`
- `*.ps1`, `*.sh` at repo root (install scripts)
- `.github/workflows/*.yml`

Excluded:

- `.lovable/`, `node_modules/`, `bun.lock*`, `package-lock.json`, `.release/`, `.git/`.
- Binary files (detected by null-byte scan of first 8KB).

### Pin Enforcement

1. Match must be the EXACT old slug. Substring matches inside other slugs are forbidden (use word-boundary or surrounding-char check).
2. Never resolve `latest`, `main`, `master`, or `HEAD` — operate purely on the `--from` / `--to` strings provided.
3. Never rewrite arbitrary `vNN` mentions in prose. Only rewrite the qualified `owner/repo` slug. Loose version references in docs are NOT touched.

### Post-Rewrite Steps (only with `--confirm`)

1. Update `package.json` `version` field to `--new-version`.
2. Run `node scripts/sync-version.mjs`.
3. Run `node scripts/sync-spec-tree.mjs`.
4. Run `node scripts/sync-readme-stats.mjs`.
5. Run `node scripts/sync-health-score.mjs`.
6. Run `python3 linter-scripts/check-root-readme.py` to verify README stamps still parse.
7. If any sync or lint step exits non-zero, restore from in-memory snapshot and exit 6.

### Output

1. Always print: `Migrating <from> → <to>, version → <new-version>`.
2. Always print a per-file change summary: `<path>: <N> replacement(s)`.
3. In dry-run, end with: `DRY RUN — no files written. Re-run with --confirm to apply.`
4. With `--confirm`, end with: `Migration complete. Review the diff and commit.`

---

## Failure Modes

| Exit | Reason |
|------|--------|
| 1 | Missing required arg (`--from`, `--to`, or `--new-version`). |
| 2 | Invalid slug or version format. |
| 3 | Major version did not increase, or `--from` == `--to`. |
| 4 | Substring collision detected (old slug appears as part of a longer identifier). |
| 5 | Excluded file matched (would have rewritten a forbidden path). |
| 6 | Post-rewrite sync or lint failed; changes rolled back. |
| 7 | Working tree not clean at start (uncommitted changes present). |

---

## Forbidden Behaviors

1. Editing `.lovable/`, `.release/`, `node_modules/`, lockfiles, or `.git/`.
2. Rewriting bare `vNN` tokens outside a qualified slug.
3. Touching the `name` field of `package.json` (Lovable scaffold; not the public repo name).
4. Performing the migration without a clean working tree (must error with exit 7).
5. Skipping the post-rewrite sync chain when `--confirm` is set.
6. Pushing, tagging, or creating a release. The script only edits files locally.

---

## CI/CD Workflow

File: `.github/workflows/migrate-repo-version.yml`

1. Trigger: `workflow_dispatch` ONLY. Never on push, PR, or schedule.
2. Inputs: `from`, `to`, `new_version`, `confirm` (boolean, default false).
3. Steps:
   a. Checkout with `fetch-depth: 0`.
   b. Setup Node, install deps with `npm ci`.
   c. Run `node scripts/migrate-repo-major-version.mjs --from "$FROM" --to "$TO" --new-version "$NEW_VERSION" $( [[ "$CONFIRM" == "true" ]] && echo --confirm )`.
   d. If `confirm == true`, open a PR via `peter-evans/create-pull-request` titled `chore: migrate <from> → <to> (vNEW)`.
   e. Never auto-merge. Human must review and merge.

---

## Validation

1. Slug regex: `^[A-Za-z0-9-]+/[A-Za-z0-9._-]+$`.
2. Version regex: `^\d+\.\d+\.\d+(-[A-Za-z0-9.]+)?$`.
3. Major-bump check: `parseInt(newVersion.split('.')[0]) > parseInt(currentVersion.split('.')[0])`.
4. Word-boundary match: replacement only fires when the old slug is preceded and followed by one of: start-of-string, end-of-string, whitespace, `"`, `'`, `(`, `)`, `<`, `>`, `,`, ``` ` ```, `]`, `[`, newline.

---

## Acceptance Criteria

1. Running with no flags exits 1 with usage.
2. Running with `--from a/b --to a/c --new-version 4.0.0` and no `--confirm` prints the planned change set and exits 0 without modifying any file.
3. Running with `--confirm` updates every in-scope file, bumps `package.json` version, runs all sync scripts, and exits 0.
4. If `package.json` is at v3.x and `--new-version` is v3.62.0, exits 3 (not a major bump).
5. If old slug appears inside a longer identifier (e.g. `coding-guidelines-v24-archive`), the script does NOT rewrite it and reports the collision.
6. If working tree is dirty at start, exits 7 before any read.
7. After a successful `--confirm` run, `npm run lint:readme` and `npm run sync` both pass on the resulting tree.
8. The new spec/14-update/26 file is referenced from `spec/14-update/00-overview.md` and `spec/14-update/readme.md`.

---

## Open Items for User Confirmation

1. **Pre-tag vs post-tag**: should this run BEFORE the new repo is renamed on GitHub, or AFTER? Default assumption: BEFORE — generates the PR that, once merged, makes the repo ready for rename.
2. **Bundle install scripts**: confirm that `*-install.{ps1,sh}` only contain qualified slug references (no bare `vNN`). If they contain bare version tokens, expand scope.
3. **README stamps**: confirm `STAMP:BADGES` and `STAMP:PLATFORM_BADGES` blocks are slug-aware (they should re-emit on next `npm run sync`).
4. **Tag creation**: confirm the migrator stays out of git tag/release creation. Releases remain a separate manual step.

---

## Job Sequence

1. Job 1: Write this spec. Status: Done on approval.
2. Job 2: Implementation — BLOCKED on Open Items confirmation.
   a. Create `scripts/migrate-repo-major-version.mjs`.
   b. Create `.github/workflows/migrate-repo-version.yml`.
   c. Add cross-references in `spec/14-update/00-overview.md` and `spec/14-update/readme.md`.
   d. Add a `npm run migrate:repo` script entry pointing at the new script.
   e. Update root `readme.md` "Power-user flags" section to mention the migrator.

---

## File System References

1. This spec: `spec/14-update/26-repo-major-version-migrator.md`.
2. Companion spec: `spec/14-update/25-release-pinned-installer.md`.
3. Future implementation: `scripts/migrate-repo-major-version.mjs`, `.github/workflows/migrate-repo-version.yml`.
