# `gitmap-v27 scan` output schema

Stable, automation-grade reference for the artifacts produced by:

```sh
gitmap-v27 scan [DIR] [--output json|csv|both|terminal] [--out-file FILE]
```

Default artifact paths (relative to the scan root, controlled by `outputDir`
in `config.json`, default `.gitmap/output/`):

| Artifact | Default path | Format | Schema |
|---|---|---|---|
| JSON | `.gitmap/output/gitmap.json` | UTF-8, pretty-printed JSON array | [`scan-output.schema.json`](./scan-output.schema.json) |
| CSV  | `.gitmap/output/gitmap.csv`  | RFC 4180, UTF-8, **CRLF** line endings | column table below |

> **Stability contract.** Field names and JSON keys listed here are
> permanent. New fields are additive only and are **appended** at the
> end of the CSV column list — never inserted in the middle, never
> renamed, never repurposed. The CSV reader tolerates legacy
> 8/9/10/12/13-column layouts so older exports still round-trip.

---

## JSON shape

Top level is a JSON array. Each element is a `ScanRecord`:

```json
[
  {
    "id": 42,
    "slug": "acme-api",
    "repoId": "123456789",
    "repoName": "acme-api",
    "httpsUrl": "https://github.com/acme/acme-api.git",
    "sshUrl": "git@github.com:acme/acme-api.git",
    "discoveredUrl": "git@github.com:acme/acme-api.git",
    "branch": "main",
    "branchSource": "head",
    "relativePath": "services/acme-api",
    "absolutePath": "/home/me/code/services/acme-api",
    "cloneInstruction": "git clone -b main https://github.com/acme/acme-api.git services/acme-api",
    "notes": "",
    "depth": 2,
    "transport": "ssh"
  }
]
```

See [`scan-output.schema.json`](./scan-output.schema.json) for the full
JSON Schema (Draft 2020-12) — usable directly with `ajv`, `jsonschema`,
`check-jsonschema`, etc.

---

## CSV column order (current, 13 columns)

The CSV header row is written verbatim from
`gitmap-v27/constants/constants_terminal.go::ScanCSVHeaders`:

| # | Column            | JSON key           | Type    | Notes |
|---|-------------------|--------------------|---------|-------|
| 1 | `repoName`        | `repoName`         | string  | Working-tree basename. |
| 2 | `httpsUrl`        | `httpsUrl`         | string  | Canonical HTTPS clone URL; empty if none. |
| 3 | `sshUrl`          | `sshUrl`           | string  | Canonical SSH clone URL; empty if none. |
| 4 | `branch`          | `branch`           | string  | Branch to clone/operate on. |
| 5 | `branchSource`    | `branchSource`     | enum    | `head` \| `config-default` \| `flag-default` \| `fallback` \| `""` |
| 6 | `relativePath`    | `relativePath`     | string  | Forward-slash path from scan root. |
| 7 | `absolutePath`    | `absolutePath`     | string  | Native absolute path at scan time. |
| 8 | `cloneInstruction`| `cloneInstruction` | string  | Ready-to-run `git clone …` command. |
| 9 | `notes`           | `notes`            | string  | Free-form annotation. |
|10 | `depth`           | `depth`            | integer | Depth from scan root; base-10. |
|11 | `repoId`          | `repoId`           | string  | Provider-side repo id (e.g. GitHub numeric id); may be empty. |
|12 | `discoveredUrl`   | `discoveredUrl`    | string  | Raw `origin` URL before HTTPS/SSH normalization. |
|13 | `transport`       | `transport`        | enum    | `ssh` \| `https` \| `other` \| `""` |

> JSON also exposes two stable fields not present in the CSV:
> `id` (local SQLite row id) and `slug` (filesystem-safe id).
> These are intentionally omitted from CSV so positional column
> indices stay backwards-compatible with all prior layouts.

### Legacy CSV layouts (still parseable)

`formatter.ParseCSV` auto-detects these older shapes by column count:

| Cols | Layout                             |
|------|------------------------------------|
| 8    | pre-`branchSource`                 |
| 9    | pre-`depth` (branchSource present) |
| 10   | pre-`repoId` (depth present)       |
| 12   | pre-`transport`                    |
| 13   | **current**                        |

Missing fields parse as zero values. New columns are always appended
at the right.

---

## Recipes

```sh
# All SSH-cloned repos as a flat list
jq -r '.[] | select(.transport=="ssh") | .sshUrl' .gitmap/output/gitmap.json

# Same query against the CSV (column 13 = transport)
awk -F, 'NR>1 && $13=="ssh"' .gitmap/output/gitmap.csv

# Repos at the max-depth boundary — candidates for a deeper rescan
jq '.[] | select(.depth==4)' .gitmap/output/gitmap.json

# Validate a JSON file against the schema
check-jsonschema --schemafile spec/08-json-schemas/scan-output.schema.json \
  .gitmap/output/gitmap.json
```

---

## Source of truth

- Struct + tags: [`gitmap-v27/model/record.go`](../../gitmap-v27/model/record.go) (`ScanRecord`)
- CSV header constant: [`gitmap-v27/constants/constants_terminal.go`](../../gitmap-v27/constants/constants_terminal.go) (`ScanCSVHeaders`)
- CSV writer / parser: [`gitmap-v27/formatter/csv.go`](../../gitmap-v27/formatter/csv.go)
- JSON writer: `gitmap-v27/formatter/json.go` (uses standard `encoding/json` with the struct tags above)

Any drift between this document and those files is a bug — please open
an issue or a PR that updates both sides in lockstep.
