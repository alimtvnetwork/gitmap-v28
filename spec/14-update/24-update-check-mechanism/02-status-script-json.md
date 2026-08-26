# Status Script JSON Contract

> **Version:** 1.0.0
> **Parent:** [00-overview.md](./00-overview.md)

---

## Purpose

Defines the JSON document that every `Status.ps1` / `Status.sh` script
**MUST** emit on stdout. All keys and string-enum values are PascalCase.

---

## Schema

```json
{
  "Version": "V1.5.0",
  "ReleaseUrl": "https://github.com/{Owner}/{Repo}/releases/tag/v1.5.0",
  "PublishedAt": "2026-04-15T10:00:00Z",
  "Checksum": "Sha256:Abc123...",
  "Install": {
    "Windows": {
      "ScriptUrl": "https://github.com/{Owner}/{Repo}/releases/download/v1.5.0/Install.ps1",
      "Command": "Iwr -UseB <ScriptUrl> | Iex"
    },
    "Unix": {
      "ScriptUrl": "https://github.com/{Owner}/{Repo}/releases/download/v1.5.0/Install.sh",
      "Command": "Curl -FsSL <ScriptUrl> | Bash"
    }
  },
  "Notes": "ReleaseNotesSummary",
  "MinSupportedFrom": "V1.0.0",
  "NewRepoUrl": null
}
```

---

## Field Reference

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `Version` | string | yes | SemVer with leading `V` (e.g., `V1.5.0`) |
| `ReleaseUrl` | string (URL) | yes | GitHub release page |
| `PublishedAt` | string (ISO 8601 UTC) | yes | Release publish time |
| `Checksum` | string | yes | Format `<Algorithm>:<Hex>` — e.g., `Sha256:Abc...` |
| `Install.Windows.ScriptUrl` | string (URL) | yes | Pinned installer per [16-generic-release/08](../../16-generic-release/08-version-pinned-release-installers.md) |
| `Install.Windows.Command` | string | yes | One-line invocation for end users |
| `Install.Unix.ScriptUrl` | string (URL) | yes | Pinned installer |
| `Install.Unix.Command` | string | yes | One-line invocation |
| `Notes` | string | no | Short release-notes summary; `null` allowed |
| `MinSupportedFrom` | string | yes | Lowest version that can directly upgrade to this one |
| `NewRepoUrl` | string (URL) \| null | yes | Set if the project moved to a new repo (migration notice) |

---

## Producer (Release Pipeline)

The script is generated at release time by the pipeline (see
[16-generic-release/08-version-pinned-release-installers.md](../../16-generic-release/08-version-pinned-release-installers.md)).
Templates live next to the pinned installer templates:

```
linter-scripts/installer-templates/
├── Status.ps1.tmpl
└── Status.sh.tmpl
```

Both templates substitute the same placeholders the pinned installers
use (`__EmbeddedVersion__`, `__EmbeddedRepo__`, `__EmbeddedCommit__`,
`__EmbeddedChecksum__`, `__EmbeddedReleaseUrl__`, `__EmbeddedNotes__`,
`__EmbeddedNewRepoUrl__`).

---

## Consumer Behavior

1. The CLI parses the JSON and keeps the entire document as `RawJson`
   in the `UpdateChecker` row.
2. Individual fields are projected into typed columns (see
   [04-database-schema.md](./04-database-schema.md)).
3. If `NewRepoUrl` is non-null, the CLI surfaces a migration notice
   on every subsequent command alongside (or instead of) the standard
   "update available" line.

---

## Validation

A status response is **valid** iff:

1. HTTP status is 200.
2. Body parses as JSON.
3. All `required = yes` fields are present and non-empty.
4. `Version` matches the regex `^V\d+\.\d+\.\d+$`.
5. `PublishedAt` parses as RFC 3339 / ISO 8601.

Any failure → treated as *not-found* in discovery. Reason logged to the
file-system log (`UpdateChecker.log`).

---

*Status Script JSON — v1.0.0 — 2026-04-20*
