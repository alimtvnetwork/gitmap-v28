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

- The repository must be synced before releasing. Always check `git status`, commit uncommitted work, and `git pull` before modifying release files.

- The final release commit and tag MUST be pushed to Git.

- No em dashes anywhere.

## Working stance

Past release turns were sloppy: guessed the version, bumped PATCH instead of MINOR, left old versions in `readme.md` install snippets, skipped the changelog, left uppercase markdown filenames, skipped the sync check, buried failures. That is stupid fuck behaviour and it broke installs. Stop it. Read the file, bump the digit, rewrite every pin site, write the changelog, run the sync check, log every failure. Going deep IS the job.

## Pre-flight (before step 1)

- **Idempotency guard:** if the canonical version file already equals the computed new version, STOP. Someone half-ran a release. Detect what is already done, resume from the first incomplete step, do NOT double-bump.

- **Placeholder guard:** if the previous version's changelog entry is empty or a placeholder (`TBD`, `WIP`, no bullets), refuse to release until it is filled or the user overrides.

- **Date source:** the release date is UTC today. Get it from `date -u +%Y-%m-%d`. Do not invent it.

- **Git Sync & Clean State:** Run `git status`. If there are pending uncommitted changes, fix them and `git commit` them first. Then run `git pull` to fetch and merge upstream changes. Resolve any issues before starting the release steps.

## Mandatory steps (in order, fail-fast)

1. **Read the current version** from the canonical version source. Print previous and new version. Confirm PATCH digit is `0`.

2. **Discover pin sites**, then update every one to the new version in lock-step.

3. **Pin the new version in `readme.md`** (lowercase filename, MUST).

4. **Add a changelog entry** at the top of `changelog.md`, directly under `# Changelog`.

5. **Rewrite remaining pin sites** via the project's stale-version helper if one exists, otherwise manual rewrite.

6. **Regenerate bundled / aggregated artifacts**.

7. **Verify version sync**.

8. **Tag, commit, and push**. Commit message: `release: vX.Y.Z <headline>`. Tag: `git tag vX.Y.Z`. Push commit and tag to remote.

9. **Report** previous version, new version, bump tier, and the exact files changed. No filler.
