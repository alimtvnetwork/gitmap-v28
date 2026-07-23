# 05 — Profiles and JSON Shape

## 5.1 File system layout

```
<workspaceRoot>/.gitmap/
├── db/
│   └── gitmap.sqlite
├── temp/
│   └── <commitInRunId>/
│       └── <orderIndex>-<basename>/   # one cloned input
├── commit-in/
│   └── profiles/
│       └── <profileName>.json
└── logs/
    └── commit-in.log
```

- `<workspaceRoot>` = nearest ancestor of CWD containing `.gitmap/`,
  else CWD.
- Profile filename MUST equal `Profile.Name` exactly (case-sensitive).
- `temp/<commitInRunId>/` is deleted on `Finalize` unless `--keep-temp`
  was passed.

## 5.2 Canonical JSON shape (`<profileName>.json`)

All keys AND all enum values are PascalCase. Unknown keys are an error
(strict decode). Order of keys is fixed (matches the order below) so
files diff cleanly.

```json
{
  "Name": "Default",
  "SchemaVersion": 1,
  "SourceRepoPath": "/abs/path/to/source",
  "IsDefault": true,
  "ConflictMode": "ForceMerge",
  "Author": {
    "Name": "Jane Doe",
    "Email": "jane@example.com"
  },
  "Exclusions": [
    { "Kind": "PathFolder", "Value": "node_modules" },
    { "Kind": "PathFile",   "Value": "secrets.env" }
  ],
  "MessageRules": [
    { "Kind": "StartsWith", "Value": "Signed-off-by:" },
    { "Kind": "Contains",   "Value": "[skip ci]" }
  ],
  "MessagePrefix":  ["chore:", "feat:"],
  "MessageSuffix":  [],
  "TitlePrefix":    "",
  "TitleSuffix":    " — via gitmap-v27",
  "OverrideMessages": ["Improve module", "Refine implementation"],
  "OverrideOnlyWeak": true,
  "WeakWords": ["change", "update", "updates"],
  "FunctionIntel": {
    "IsEnabled": true,
    "Languages": ["Go", "TypeScript", "Python"]
  }
}
```

## 5.3 JSON Schema (informative — store in `spec/08-json-schemas/commit-in-profile.schema.json` when implementation starts)

- `SchemaVersion: 1` is REQUIRED. A loader that sees an unknown
  `SchemaVersion` exits `CommitInExitProfileMissing`.
- `Author` is optional. Absent → use OS git default → fall back to
  source commit's author.
- `FunctionIntel.Languages` MUST be a subset of the
  `FunctionIntelLanguage` enum (`Go`, `JavaScript`, `TypeScript`,
  `Rust`, `Python`, `Php`, `Java`, `CSharp`).

## 5.4 Binding rule (resolves Ambiguity #7)

Profiles are bound to `<source>` by **absolute symlink-resolved
filesystem path**, not by `origin` URL. Rationale:

- A user can have multiple clones of the same upstream URL with
  different commit-in policies.
- The path is known at every stage; the `origin` URL is only known
  after `EnsureSource` runs and may not exist (freshly init-ed repos).

`--default` semantics: load the single `Profile` row where
`SourceRepoPath = <source>` AND `IsDefault = 1`. Zero matches →
`CommitInExitProfileMissing`. Two matches → impossible (unique partial
index, see §4.3).

## 5.5 Save semantics (resolves Ambiguity #8)

- `--save-profile <name>` writes the JSON file AND inserts/updates
  `Profile` + child rows in one transaction.
- If a file or row with `Name` already exists: refuse with exit
  `CommitInExitBadArgs` UNLESS `--save-profile-overwrite` is also
  passed.
- `--set-default` flips `IsDefault = 1` for this profile and `0` for
  every other profile bound to the same `SourceRepoPath` (atomic in
  the same transaction).

## 5.6 Profile load order (highest precedence wins)

1. Explicit flags on the current invocation (`--conflict`, `--exclude`, …).
2. `--profile <name>` JSON file (and mirrored DB rows).
3. `--default` profile (`Profile.IsDefault = 1`).
4. Built-in defaults (see §02 flag-default column).

An unset value at every layer triggers an interactive prompt unless
`--no-prompt` was passed (then exit `CommitInExitMissingAnswer`).