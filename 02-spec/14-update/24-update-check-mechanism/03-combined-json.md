# Combined Discovery JSON

> **Version:** 1.0.0
> **Parent:** [00-overview.md](./00-overview.md)

---

## Purpose

After the six parallel probes complete, the CLI assembles a **single
combined JSON document** that captures the entire discovery result.
This is the document persisted as `RawJson` in the `UpdateChecker`
table and printed by `update-check` (sync mode).

---

## Schema

```json
{
  "CurrentVersion": "V1.0.0",
  "LatestVersion": "V1.5.0",
  "HasUpdate": true,
  "CheckedAt": "2026-04-20T12:00:00Z",
  "OwnerKind": "User",
  "Owner": "AlifKarim",
  "CurrentRepo": "Movie-Cli-V15",
  "Candidates": [
    { "Version": "V1.1.0", "Found": true,  "Url": "https://raw.githubusercontent.com/AlifKarim/Movie-Cli-V16/main/Status.ps1" },
    { "Version": "V1.2.0", "Found": true,  "Url": "https://raw.githubusercontent.com/AlifKarim/Movie-Cli-V17/main/Status.ps1" },
    { "Version": "V1.3.0", "Found": false, "Url": null },
    { "Version": "V1.4.0", "Found": true,  "Url": "https://raw.githubusercontent.com/AlifKarim/Movie-Cli-V19/main/Status.ps1" },
    { "Version": "V1.5.0", "Found": true,  "Url": "https://raw.githubusercontent.com/AlifKarim/Movie-Cli-V20/main/Status.ps1" }
  ],
  "Selected": {
    "Version": "V1.5.0",
    "ReleaseUrl": "...",
    "PublishedAt": "...",
    "Checksum": "...",
    "Install": { "Windows": { "...": "..." }, "Unix": { "...": "..." } },
    "Notes": "...",
    "MinSupportedFrom": "V1.0.0",
    "NewRepoUrl": null
  }
}
```

---

## Field Reference

| Field | Type | Notes |
|-------|------|-------|
| `CurrentVersion` | string | Installed CLI version |
| `LatestVersion` | string | Highest version returned by a successful probe |
| `HasUpdate` | boolean | `LatestVersion > CurrentVersion` |
| `CheckedAt` | string (ISO 8601 UTC) | When discovery completed |
| `OwnerKind` | enum string | `User` \| `Organization` |
| `Owner` | string | GitHub login |
| `CurrentRepo` | string | The repo the installed CLI was built from |
| `Candidates` | array (length 5) | One entry per **lookahead** probe (V+1..V+5) |
| `Candidates[].Version` | string | Probed version label |
| `Candidates[].Found` | boolean | True if HTTP 200 + valid JSON |
| `Candidates[].Url` | string \| null | The status URL probed; null when not found |
| `Selected` | object | Full status JSON from the winning probe — schema defined in [02-status-script-json.md](./02-status-script-json.md) |

> The current-version (V) probe is **not** listed in `Candidates` — its
> result is captured in `CurrentVersion` and is the comparison baseline.

---

## Examples

### No update available

```json
{
  "CurrentVersion": "V1.5.0",
  "LatestVersion": "V1.5.0",
  "HasUpdate": false,
  "CheckedAt": "2026-04-20T12:00:00Z",
  "Candidates": [
    { "Version": "V1.6.0", "Found": false, "Url": null },
    { "Version": "V1.7.0", "Found": false, "Url": null },
    { "Version": "V1.8.0", "Found": false, "Url": null },
    { "Version": "V1.9.0", "Found": false, "Url": null },
    { "Version": "V2.0.0", "Found": false, "Url": null }
  ],
  "Selected": null
}
```

### Migration to a new repo

`Selected.NewRepoUrl` is non-null. The CLI surfaces a migration banner
in addition to the standard update warning.

---

*Combined Discovery JSON — v1.0.0 — 2026-04-20*
