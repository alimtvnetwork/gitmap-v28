# Release, MINOR bump, MUST enforcement

**Trigger phrases:** `release`, `bump version`, `bump version + add changelog + pin to root readme`, `abump version ...` (typo variants count).

No variables. No prompts to the user. Discover the current version from disk, bump it, do the full ceremony in one turn.

---

## RULE 0, MUST, NON-NEGOTIABLE

1. Read the **canonical version source** for THIS repo (discover it: `version.json`, `package.json` `"version"`, or whatever single file the repo treats as the version of record). Do not guess.

2. Bump MINOR only: `MAJOR.MINOR.PATCH` becomes `MAJOR.(MINOR+1).0`. PATCH MUST reset to `0`.

3. State the previous version and new version explicitly in the reply, before touching any file.

4. Do NOT ask "minor or patch?". Do NOT open plan mode. Do NOT ask for confirmation.

Deviations (only when the trigger explicitly says so):

- **MAJOR** = `(MAJOR+1).0.0` if the user said the change is breaking (storage schema, prompt schema, public SDK, extension contract).

- **PATCH** = `MAJOR.MINOR.(PATCH+1)` only if the user literally said `patch bump` or `patch release`.

When in doubt: MINOR.

## Hard rules (MUST)

- All version pin sites move in lock-step. Partial bumps are rejected.

- The previous version string MUST NOT appear anywhere in the repo after this turn EXCEPT in historic files: `changelog.md`, `release_notes.md`, anything under `.lovable/release/`, and any dated archive folder.

- Changelog entry under the new version heading is MANDATORY. A release without one is INVALID.

- **All markdown filenames MUST be lowercase**: `readme.md`, `changelog.md`, `release_notes.md`, every audit / issue / plan / spec `.md`. Rename any `README.md`, `CHANGELOG.md`, `ReadMe.md`, etc. in the same turn with `mv` (or `git mv` if tracked), and update every reference.

- If ANY step fails or is flagged, log it under `.lovable/release/issues/xx-<new-version>-<slug>.md` AND add an `### Issues` bullet under the new changelog entry linking to that file. Never hide failures.

- Never invent changelog bullets. Only real work since the previous release.

- Never auto-publish unless the user explicitly said `publish`, `deploy`, `ship`, or `go live`.

- No em dashes anywhere.

## Working stance

Past release turns were sloppy: guessed the version, bumped PATCH instead of MINOR, left old versions in `readme.md` install snippets, skipped the changelog, left uppercase markdown filenames, skipped the sync check, buried failures. That is stupid fuck behaviour and it broke installs. Stop it. Read the file, bump the digit, rewrite every pin site, write the changelog, run the sync check, log every failure. Going deep IS the job.

## Pre-flight (before step 1)

- **Idempotency guard:** if the canonical version file already equals the computed new version, STOP. Someone half-ran a release. Detect what is already done, resume from the first incomplete step, do NOT double-bump.

- **Placeholder guard:** if the previous version's changelog entry is empty or a placeholder (`TBD`, `WIP`, no bullets), refuse to release until it is filled or the user overrides.

- **Date source:** the release date is UTC today. Get it from `date -u +%Y-%m-%d`. Do not invent it.

## Mandatory steps (in order, fail-fast)

1. **Read the current version** from the canonical version source. Print previous and new version. Confirm PATCH digit is `0`.

2. **Discover pin sites**, then update every one to the new version in lock-step. Use a single canonical search:

   ```
   rg -n "\b<PREV_MAJOR>\.<PREV_MINOR>\.<PREV_PATCH>\b" -g '!node_modules' -g '!*.lock' -g '!.git'
   rg -n "\b(VERSION|APP_VERSION|EXTENSION_VERSION|SCHEMA_VERSION|CACHE_SCHEMA_VERSION|BUILD_VERSION)\b"
   ```

   Typical pin sites (non-exhaustive):

   - Canonical version file (set `releaseDate` to UTC today if the field exists).

   - Manifests: `manifest.json`, extension / plugin manifests.

   - Source constants named like `VERSION`, `APP_VERSION`, `EXTENSION_VERSION`, `SCHEMA_VERSION`, `CACHE_SCHEMA_VERSION`, `BUILD_VERSION`.

   - `instruction.ts` / `instruction.md` / metadata files with a `Version:` field.

   - Sub-packages under `scripts/`, `standalone-scripts/`, `packages/`, `apps/` carrying their own version constant.

   - Install snippets, badge URLs, zip filenames, release-branch examples, "current version" lines in docs.

   A pin site missed here breaks installs.

3. **Pin the new version in `readme.md`** (lowercase filename, MUST). Rewrite every occurrence of the previous version (`vX.Y.Z` and bare `X.Y.Z`) in badges, install snippets, "current version" lines, release-branch examples, zip filenames, inline references. After this step, `rg "<previous-version>" readme.md` MUST return nothing.

4. **Add a changelog entry** at the top of `changelog.md`, directly under `# Changelog`. Replace `X.Y.Z` with the actual new version and `YYYY-MM-DD` with `date -u +%Y-%m-%d` output:

   ```
   ## [vX.Y.Z] YYYY-MM-DD <short headline>

   ### Added / Changed / Fixed / Removed

   - <one bullet per real change, naming the exact file or behaviour>

   ### Issues (only if any step failed or was flagged)

   - [xx-<new-version>-<slug>](.lovable/release/issues/xx-<new-version>-<slug>.md) short description
   ```

   Use only the subheadings that apply. `### Issues` is REQUIRED whenever any step surfaced a problem, even if worked around.

5. **Rewrite remaining pin sites via the project's stale-version helper if one exists** (discover: `scripts/update-stale-version-refs.*`, `scripts/bump-version.*`, `tools/update-versions.*`). Run it with previous and new version. If no helper exists, use the `rg` output from step 2 and rewrite each match by hand.

6. **Regenerate bundled / aggregated artifacts** (aggregated prompts, generated docs, compiled manifests) if their sources changed this turn. Use whatever generation script the project ships.

7. **Verify version sync**. Run the project's version-sync check script if one exists (discover: `scripts/check-version-sync.*`, `scripts/verify-versions.*`). It MUST exit 0. Non-zero = release is INVALID: log an issue, fix, re-run. If no such script exists, re-run the step 2 `rg` and confirm only historic files (see Hard rules allow-list) still reference the previous version.

8. **Tag and commit** (only if the repo is git-tracked and the user did not disable tagging). Commit message: `release: vX.Y.Z <headline>`. Tag: `git tag vX.Y.Z`. Do NOT push automatically. If the working tree has uncommitted changes UNRELATED to the release, stop and report; do not sweep them into the release commit.

9. **Report** previous version, new version, bump tier, and the exact files changed. No filler.

## Issue logging (MUST, when anything goes wrong)

Path: `.lovable/release/issues/xx-<new-version>-<slug>.md` (lowercase). Body:

- Previous version and new version

- Step that failed (number and name)

- Command run and full error output

- Files involved

- Resolution or workaround, or `unresolved`

Then link it from the `### Issues` bullet under the changelog entry.

## Checklist before you claim done

- [ ] Previous version read from the canonical version source, not memory.

- [ ] New version is a MINOR bump (or explicit MAJOR / PATCH per rules); PATCH digit is `0`.

- [ ] Previous and new version both stated in the reply.

- [ ] Pre-flight passed (idempotency, changelog placeholder, UTC date from `date -u`).

- [ ] Every pin site from step 2 `rg` output matches the new version.

- [ ] Canonical version file's `releaseDate` (if the field exists) is today's UTC date.

- [ ] Changelog entry added at the top of `changelog.md` with real bullets only.

- [ ] All markdown filenames in the repo are lowercase.

- [ ] `rg "<previous-version>"` returns matches ONLY in the historic allow-list (`changelog.md`, `release_notes.md`, `.lovable/release/`, dated archives).

- [ ] `### Issues` block present in the changelog if any step failed or was flagged, with links to `.lovable/release/issues/` files.

- [ ] Stale-version helper (if it exists) ran successfully; otherwise manual rewrite done.

- [ ] Bundled / aggregated artifacts regenerated if their sources changed.

- [ ] Version-sync check (if it exists) exited 0; otherwise manual `rg` confirms allow-list only.

- [ ] Commit + tag created (if git-tracked) with `release: vX.Y.Z <headline>` and `vX.Y.Z`, not pushed.

- [ ] Report includes previous version, new version, tier, and exact file list.

- [ ] No em dashes. No auto-publish.

## Prompt maintenance (meta, run once at end)

Save this prompt's full body into `.lovable/prompts/XX-release.md` (lowercase):

- If any existing file in `.lovable/prompts/` matches `*release*.md` (case-insensitive), OVERWRITE it in place. Do not create a duplicate.

- Otherwise pick `XX` = next 2-digit zero-padded sequence (highest existing `XX` prefix + 1, or `01` if the folder is empty / missing). Create the folder if needed.

- Save the prompt body only, no chat wrapping.

## Must Follow and without negotiation

Listen, past release turns were sloppy. Read the canonical version file, bump MINOR, reset PATCH to zero, pin the new version in `readme.md`, propagate everywhere via the helper (or by hand from the `rg` output), rename every uppercase markdown file to lowercase, write the changelog under the new version, log every failure or flagged issue under `.lovable/release/issues/` with a matching `### Issues` bullet, tag and commit (do not push). Skipping any step = broken installs. Going deep IS the job.

## Ambiguity handling (open questions and answers)

Ambiguity is not a license to guess. It is a file to write.

- Open: `.lovable/ambiguous-questions/01-new-ambiguity/XX-<slug>.md`

- Answered: `.lovable/ambiguous-questions/02-ambiguity-resolved/XX-<slug>.md`

New question file shape:

```
# <one-line question>

Slug: <slug>

Status: open

Raised: <YYYY-MM-DD>

Blocking: release {{version}}

## Question

## Options considered

## Impact if guessed wrong
```

When answered: `mv` from `01-new-ambiguity/` to `02-ambiguity-resolved/`, flip `Status: resolved`, and append a `## Resolution` block (`Answered:`, `Answer:`, `Applied solution:`). Never leave a copy behind. Do NOT confuse ambiguities with release issues: unknown version source, unclear bump policy, or missing changelog target = ambiguity; a failed step during the release run = `.lovable/release/issues/`.
