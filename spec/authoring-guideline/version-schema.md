# VersionSchema (`version.json`)

**Version:** 1.0.0
**Updated:** 2026-04-28
**Status:** Active — **MANDATORY** for every repository in this family
**AI Confidence:** Production-Ready
**Ambiguity:** None

> **Canonical location:** This file is a verbatim mirror published at the
> path requested by the spec author
> (`spec/authoring-guideline/VersionSchema.md`). The numeric, convention-
> aligned canonical copy lives at
> [`spec/01-spec-authoring-guide/17-version-schema.md`](../01-spec-authoring-guide/17-version-schema.md).
> Both files MUST stay byte-identical in their normative sections (§1–§11).
> If they drift, the numeric copy wins.

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

| Item                     | Location                                            |
|--------------------------|-----------------------------------------------------|
| Schema spec (mirror)     | `spec/authoring-guideline/VersionSchema.md`         |
| Schema spec (canonical)  | `spec/01-spec-authoring-guide/17-version-schema.md` |
| Runtime artifact         | `version.json` at the **repository root**           |

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

---

## §4 — Top-Level Schema

| Key             | Type            | Required | Description                                                                 |
|-----------------|-----------------|----------|-----------------------------------------------------------------------------|
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

The canonical role set is **closed**:

| Value           | Meaning                                                         |
|-----------------|-----------------------------------------------------------------|
| `PrimaryAuthor` | Original creator and primary owner. Exactly one per repo.       |
| `Contributor`   | Made non-trivial contributions (code, spec, docs).              |
| `Maintainer`    | Has merge / release authority on the integration branch.        |
| `Reviewer`      | Has review authority but does not maintain the repo day-to-day. |
| `Sponsor`       | Funded, hosted, or otherwise enabled the work.                  |

- Any value outside this list is a schema violation.
- A single person MAY appear in multiple roles only by appearing as
  multiple `Author` entries.

---

## §7 — `LastCommitSha` Automation

`LastCommitSha` MUST be kept in sync with the actual latest commit on the
integration branch. Manual edits are a **fallback only**.

This repository uses a **husky pre-commit hook** (`.husky/pre-commit`)
that invokes `scripts/sync-version.mjs`, which reads `git rev-parse HEAD`,
writes the SHA into `version.json`, and re-stages the file.

Acceptable alternatives for repos without husky: native Git pre-commit
hook, CI workflow on push to `main`, or build-time release script. The
acceptance rule is identical: published `LastCommitSha` MUST equal the
real integration-branch HEAD.

---

## §8 — Read Behavior

1. All apps, scripts, and CI jobs MUST read version data from
   `version.json` rather than hardcoding it.
2. If `version.json` is missing or unparseable, log a clear warning and
   fall back to safe defaults (`Version = "0.0.0"`, `Title = RepoSlug`).
   Never crash solely because the manifest is missing.
3. Readers MUST treat the file as **read-only**. Only `sync-version.mjs`
   writes to it. Multiple writers cause non-deterministic merge corruption — the linter `validate-version-json` MUST fail the build if any non-`sync-version.mjs` source mutates `version.json`.

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

This repository historically shipped `version.json` with camelCase keys
(`version`, `updated`, `git`, `stats`, `folders`). During migration,
`sync-version.mjs` emits **both** PascalCase (§4–§5) and the legacy
camelCase keys. New code MUST read PascalCase. Legacy keys are
**deprecated** and removed once all readers migrate. Fresh repositories
MUST emit only §4.

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
| AC-VS-007  | No source file outside `version.json` and `scripts/sync-version.mjs` hardcodes `Version` or `RepoUrl`.   |
| AC-VS-008  | Readers handle a missing `version.json` with a logged warning and safe fallback (no crash).              |

---

## §12 — Cross-References

- [Canonical schema spec](../01-spec-authoring-guide/17-version-schema.md)
- [Key Naming PascalCase](../02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md)
- [Enum Standards](../17-consolidated-guidelines/04-enum-standards.md)
- [Author Attribution](mem://project/author-attribution)
- [CI/CD Issue 06 — Version Drift](../../.lovable/resolved-issues/06-version-drift-after-package-bump.md)
