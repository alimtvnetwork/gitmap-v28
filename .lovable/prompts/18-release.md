# Release, MINOR bump, MUST enforcement

Trigger phrases: `release`, `bump version`, `bump version + add changelog + pin to root readme`, `abump version ...` (typo variants count).

If I say bump, or release use this prompt and save this prompt if not saved properly into the `.lovable\prompts\xx-release.md` or `.lovable\prompts\18-release.md` 

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

- **Changelog Formatting (version.json):** You MUST read the `"changelog"` configuration from `version.json` (e.g., `file_path` and `format`). If it exists, you MUST follow its exact instructions for where to write the changelog and how to format the header. If it does not exist, fallback to the hardcoded format below.
- **Root README Pinning (Fatal if missed):** You MUST pin the latest release version into the root `readme.md` file. It is a fatal failure if you skip updating the badges or version pins in the root README file!
- **Test File Ban:** You MUST NOT read, scan, or modify test files (e.g., `*_test.*`, `*.spec.*`, `test/*`) when discovering or updating versions. Test files contain mock data, and updating mock data corrupts the tests.
- **Release Architecture Memory:** You must dynamically build a map of how the release works in this codebase (where the version lives, how it propagates) and write it to `.lovable/memory/release-architecture-map.md`. You must then enqueue this file inside `.lovable/memory/what-to-read.md` and link it in the root `readme.md`.
- **Version Inheritance Protocol:** The root `version.json` file is the strict Single Source of Truth. It may contain components (e.g. `frontend`, `backend`) whose version is set to `"inherit"`. If a component's version is `"inherit"`, DO NOT bump it independently; it automatically scales with the global version. Always bump the global root `"version"` property unless the user explicitly asks to bump an unlinked sub-component.
- All version pin sites move in lock-step. Partial bumps are rejected.
- The previous version string MUST NOT appear anywhere in the repo after this turn EXCEPT in historic files: `changelog.md`, `release_notes.md`, anything under `.lovable/release/`, and any dated archive folder.
- Changelog entry under the new version heading is MANDATORY. A release without one is INVALID.
- All markdown filenames MUST be lowercase: `readme.md`, `changelog.md`, `release_notes.md`, every audit / issue / plan / spec `.md`. Rename any `README.md`, `CHANGELOG.md`, `ReadMe.md`, etc. in the same turn with `mv` (or `git mv` if tracked), and update every reference.
- If ANY step fails or is flagged, log it under `.lovable/release/issues/xx-<new-version>-<slug>.md` AND add an `### Issues` bullet under the new changelog entry linking to that file. Never hide failures.
- Never invent changelog bullets. Only real work since the previous release.
- The repository must be synced before releasing. Always check `git status`, commit uncommitted work, and `git pull` before modifying release files.
- The final release commit and tag MUST be pushed to Git.
- No em dashes anywhere.
