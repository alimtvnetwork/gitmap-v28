# Version Schema (`version.json`)

**Version:** 1.0.0
**Updated:** 2026-04-28
**Status:** Active — **MANDATORY** for every repository in this family
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## §1 — Purpose

Every repository in this family MUST contain a single, authoritative
`version.json` file at the **repository root**. It is the canonical source
of truth for:

- Repository identity (`Title`, `RepoSlug`, `RepoUrl`)
- Current shipped version (`Version`)
- Last commit on the integration branch (`LastCommitSha`)
- Authorship (`Authors[]`)
- Human-readable purpose (`Description`)

All apps, scripts, CI jobs, installers, and contributors MUST read these
values from `version.json`. Hardcoding any of these values inside source
code, docs, or pipelines is a **build-time violation**.

> **Why a single file?** Multiple sources of truth drift. A single file
> with a fixed schema and an automated sync gate cannot.

---

## §2 — File Location

| Item                | Location                                         |
|---------------------|--------------------------------------------------|
| Schema spec (this)  | `spec/01-spec-authoring-guide/17-version-schema.md` |
| Runtime artifact    | `version.json` at the **repository root**        |

The runtime file MUST NOT live inside `spec/`, `src/`, or any subfolder.
Tools locate it by walking up to the nearest directory that contains both
`package.json` (or equivalent) and `version.json`.

---

## §3 — Naming Rules

1. All JSON keys use **PascalCase** (`Version`, `RepoSlug`, `LastCommitSha`).
2. Enum values use **PascalCase** (`PrimaryAuthor`, not `primary_author`).
3. Free-text values (`Title`, `Description`, `Background`) preserve natural
   capitalization and are not coerced.
4. URL strings are preserved verbatim.

These rules align with the project-wide PascalCase mandate
(`spec/02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md`)
and the cross-language Enum standard
(`spec/17-consolidated-guidelines/04-enum-standards.md`).

---

## §4 — Top-Level Schema

| Key             | Type            | Required | Description                                                                 |
|-----------------|-----------------|----------|-----------------------------------------------------------------------------|
| `_purpose`      | String          | optional | Human/AI explanation of version.json as the canonical source of truth.      |
| `_instructions` | Object          | optional | Explicit instructions on how to use/change the version (sync workflow).     |
| `Version`       | String (SemVer) | yes      | Current shipped semantic version (for example `1.4.2`).                     |
| `Title`         | String          | yes      | Human-readable application or project title.                                |
| `RepoSlug`      | String (slug)   | yes      | Repository name in slug form (for example `git-logs-app`).                  |
| `RepoUrl`       | String (URL)    | yes      | Full URL to the repository (HTTPS preferred).                               |
| `LastCommitSha` | String (40-hex) | yes      | SHA of the last commit on the integration branch. Auto-managed (see §7).    |
| `Description`   | String          | yes      | One- or two-sentence summary of the repository's purpose.                   |
| `Authors`       | Array           | yes      | Non-empty array of `Author` objects (see §5). Exactly one `PrimaryAuthor`.  |

> **Forbidden:** Free-text status fields, build timestamps, file counts, or
> any derived statistic. Those belong in separate generated files (such as
> `public/health-score.json`), not in the identity manifest.

---

## §5 — `Author` Object

| Key          | Type            | Required | Description                                                                |
|--------------|-----------------|----------|----------------------------------------------------------------------------|
| `Name`       | String          | yes      | Author's full name.                                                        |
| `Urls`       | Array of String | yes      | Zero or more URLs (website, GitHub, LinkedIn, ORCID, etc.).                |
| `Role`       | Enum (`Role`)   | yes      | One value from the `Role` enum (see §6). Free text is **forbidden**.       |
| `Background` | String          | yes      | Short biography or background note. Single paragraph, plain text.          |

Rules:

1. `Authors` MUST contain **at least one** entry.
2. **Exactly one** entry MUST have `Role = "PrimaryAuthor"`.
3. `Urls` MAY be empty (`[]`) but the key MUST be present.
4. Order is preserved; the first entry is conventionally the `PrimaryAuthor`.

---

## §6 — `Role` Enum

The canonical role set is:

| Value           | Meaning                                                         |
|-----------------|-----------------------------------------------------------------|
| `PrimaryAuthor` | Original creator and primary owner. Exactly one per repo.       |
| `Contributor`   | Made non-trivial contributions (code, spec, docs).              |
| `Maintainer`    | Has merge / release authority on the integration branch.        |
| `Reviewer`      | Has review authority but does not maintain the repo day-to-day. |
| `Sponsor`       | Funded, hosted, or otherwise enabled the work.                  |

- The enum is **closed**: any value outside this list is a schema violation.
- Mirrors the cross-language enum convention in
  `spec/17-consolidated-guidelines/04-enum-standards.md`.
- A single person MAY appear in multiple roles only by appearing as
  multiple `Author` entries.

---

## §6.5 — Sub-Package Tracks & the `"inherit"` Keyword

In polyglot codebases or monorepos containing multiple services or client-server packages (such as `backend`, `frontend`, `cli`, `linters`, `services`, `packages`):

1. **Root Centralization**: All sub-package version tracks MUST be declared in the root `version.json`.
2. **Inheritance (`"inherit"`)**: Setting a sub-package key to `"inherit"` (e.g. `"backend": "inherit"`, `"frontend": "inherit"`) mandates that all build tools, runtimes, and Dockerfiles resolving that sub-package inherit the root global `Version`.
3. **Explicit Tracks**: If a sub-package maintains a separate version cadence (e.g. `"linters": "3.79.0"`), it is declared as an explicit SemVer string while remaining tracked in the single root `version.json`.
4. **Mandatory Codebase Imports**: Every language implementation in the repo (Go, TypeScript, Python, PHP, C#, Rust) MUST import or read the root `version.json` file dynamically rather than hardcoding version strings.

---

## §7 — `LastCommitSha` Automation

`LastCommitSha` MUST be kept in sync with the actual latest commit on the
integration branch. Manual edits are a **fallback only**.

### §7.1 Required automation (this repo)

This repository uses a **husky pre-commit hook**:

1. Before each commit, the hook runs `scripts/sync-version.mjs`, which:
   - Reads `git rev-parse HEAD` for the current branch tip,
   - Writes the resulting SHA into `version.json` as `LastCommitSha`,
   - Re-stages `version.json` so the updated value ships in the commit.
2. The hook is wired in `.husky/pre-commit`.
3. The CI `version-drift` gate (see
   `.lovable/cicd-issues/06-version-drift-after-package-bump.md`) blocks
   any push where `LastCommitSha` does not match the committed `HEAD~1`
   (the parent at the time the hook ran).

### §7.2 Acceptable alternatives

Repositories that cannot use husky MAY substitute:

| Mechanism                            | Notes                                                                                                |
|--------------------------------------|------------------------------------------------------------------------------------------------------|
| Native Git `pre-commit` hook         | Same behavior as §7.1, no husky dependency.                                                          |
| CI workflow on push to `main`        | A bot updates `version.json` post-merge and force-pushes the change back. Higher latency.            |
| Build-time generation                | A release script populates `LastCommitSha` immediately before producing release artifacts.           |

Whichever is chosen, the **acceptance rule** is identical: after any
commit reaches the integration branch, `LastCommitSha` in the published
`version.json` MUST equal the actual SHA of that commit.

---

## §8 — Read Behavior

1. All apps, scripts, and CI jobs MUST read version data from
   `version.json` rather than hardcoding it.
2. If `version.json` is missing or unparseable:
   - Log a clear warning per the error-handling spec
     (`spec/03-error-manage/...`),
   - Fall back to safe defaults (for example `Version = "0.0.0"`,
     `Title = RepoSlug`),
   - Never crash the host application solely because the manifest is
     missing.
3. Readers SHOULD treat the file as **read-only**. Only `sync-version.mjs`
   (or its documented equivalent) writes to it.

---

## §9 — Example `version.json`

```json
{
  "Version": "1.4.2",
  "Title": "Git Logs App",
  "RepoSlug": "git-logs-app",
  "RepoUrl": "https://github.com/org/git-logs-app",
  "LastCommitSha": "a1b2c3d4e5f60718293a4b5c6d7e8f9012345678",
  "Description": "Visualizes Git log activity across repos.",
  "Authors": [
    {
      "Name": "Md. Alim Ul Karim",
      "Urls": [
        "https://riseup.asia",
        "https://github.com/alim-ul-karim"
      ],
      "Role": "PrimaryAuthor",
      "Background": "Founder, Riseup Asia LLC. Focus on developer tooling and spec-first AI workflows."
    },
    {
      "Name": "Jane Doe",
      "Urls": ["https://github.com/janedoe"],
      "Role": "Contributor",
      "Background": "Senior engineer focused on developer tooling."
    }
  ]
}
```

---

## §10 — Coexistence with Legacy Fields (Transitional)

This repository historically shipped a `version.json` with camelCase keys
(`version`, `updated`, `git`, `stats`, `folders`) consumed by the
dashboard, installer, and sync-check drift gate. To migrate without a
breaking change:

1. `sync-version.mjs` writes **both** the new PascalCase keys (§4–§5) and
   the legacy camelCase keys in the same file.
2. New code MUST read PascalCase keys.
3. Legacy camelCase keys are **deprecated** and will be removed once
   every reader has been migrated. Track removal in
   `spec/14-update/` change log.

A repository starting fresh MUST emit only the §4 schema and skip the
legacy keys entirely.

---

## §11 — Acceptance Criteria

| ID         | Statement                                                                                                |
|------------|----------------------------------------------------------------------------------------------------------|
| AC-VS-001  | A `version.json` file exists at the repo root and parses as JSON.                                        |
| AC-VS-002  | All required §4 keys are present and non-null.                                                           |
| AC-VS-003  | All keys use PascalCase exactly as documented.                                                           |
| AC-VS-004  | `Authors` is a non-empty array containing exactly one entry with `Role = "PrimaryAuthor"`.               |
| AC-VS-005  | Every `Author.Role` value is a member of the §6 `Role` enum.                                             |
| AC-VS-006  | `LastCommitSha` matches the actual integration-branch HEAD after any commit lands.                       |
| AC-VS-007  | No source file outside `version.json` and `scripts/sync-version.mjs` hardcodes `Version` or `RepoUrl`.    |
| AC-VS-008  | Readers handle a missing `version.json` with a logged warning and safe fallback (no crash).              |

---

## §12 — Cross-References

- [Key Naming PascalCase](../02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md)
- [Enum Standards](../17-consolidated-guidelines/04-enum-standards.md)
- [Root README Conventions](./11-root-readme-conventions.md)
- [Author Attribution](mem://project/author-attribution)
- [CI/CD Issue 06 — Version Drift](../../.lovable/resolved-issues/06-version-drift-after-package-bump.md)
