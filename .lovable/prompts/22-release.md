# Release, MINOR bump, MUST enforcement

Trigger phrases: `release`, `bump version`, `bump version + add changelog + pin to root readme`, `abump version ...` (typo variants count).

No variables. No prompts to the user. Discover the current version from disk, bump it, do the full ceremony in one turn.

---

## RULE 0, MUST, NON-NEGOTIABLE

1. Read the canonical version source for THIS repo (discover it: `version.json`, `package.json` `"version"`, or whatever single file the repo treats as the version of record). Do not guess.

2. Bump MINOR only: `MAJOR.MINOR.PATCH` becomes `MAJOR.(MINOR+1).0`. PATCH MUST reset to `0`.

3. State the previous version and new version explicitly in the reply, before touching any file.

4. Do NOT ask "minor or patch?". Do NOT open plan mode. Do NOT ask for confirmation.

Deviations (only when the trigger explicitly says so):

- MAJOR = `(MAJOR+1).0.0` if the user said the change is breaking (storage schema, prompt schema, public SDK, extension contract).

- PATCH = `MAJOR.MINOR.(PATCH+1)` only if the user literally said `patch bump` or `patch release`.

When in doubt: MINOR.

## Hard rules (MUST)

- All version pin sites move in lock-step. Partial bumps are rejected.

- The previous version string MUST NOT appear anywhere in the repo after this turn EXCEPT in historic files: `changelog.md`, `release_notes.md`, anything under `.lovable/release/`, and any dated archive folder.

- Changelog entry under the new version heading is MANDATORY. A release without one is INVALID.

- All markdown filenames MUST be lowercase: `readme.md`, `changelog.md`, `release_notes.md`, every audit / issue / plan / spec `.md`. Rename any `README.md`, `CHANGELOG.md`, `ReadMe.md`, etc. in the same turn with `mv` (or `git mv` if tracked), and update every reference.

- If ANY step fails or is flagged, log it under `.lovable/release/issues/xx-<new-version>-<slug>.md` AND add an `### Issues` bullet under the new changelog entry linking to that file. Never hide failures.

- Never invent changelog bullets. Only real work since the previous release.

- The repository must be synced before releasing. Always check `git status`, commit uncommitted work, and `git pull` before modifying release files.

- The final release commit and tag MUST be pushed to Git.

- No em dashes anywhere.
