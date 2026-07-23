# replace

Repo-wide find/replace across every text file. Two modes: literal text
swap, or version-suffix bump driven by the git remote URL.

**Alias:** `rpl`

Spec: `spec/04-generic-cli/15-replace-command.md`.

## Usage

```
gitmap replace "<old>" "<new>"     # literal text replace
gitmap replace -N                   # bump v(K-N)..v(K-1) → vK
gitmap replace --audit              # report-only scan, no writes
gitmap replace all                  # bump v1..v(K-1) → vK
```

## Flags

| Flag         | Description                                       |
|--------------|---------------------------------------------------|
| `--yes` `-y` | Skip the y/N confirmation prompt                  |
| `--dry-run`  | Print summary only, never write                   |
| `--quiet` `-q` | Suppress per-file diff lines                    |
| `--ext`      | Comma-separated extension allow-list (e.g. `.go,.md`). Leading dot optional. |
| `--ext-case` | `sensitive` or `insensitive` (default `insensitive`). Sensitive preserves the user's casing in `--ext` and matches filenames byte-exact. |

## Examples

### Literal replace with confirmation

```
$ gitmap replace "old-name" "new-name"

  replace: scanning 412 files in /repo
  replace: README.md: 2 matches (old-name -> new-name)
  replace: src/main.go: 1 match (old-name -> new-name)
  replace: 2 files, 3 replacements
  Apply 3 replacements across 2 files? [y/N]: y
  replace: applied 3 replacements across 2 files
```

### Version bump (`-3` on a `gitmap-v28` repo)

```
$ gitmap replace -3

  replace: scanning 412 files in /repo
  replace: go.mod: 1 match (gitmap-v28 -> gitmap-v28)
  replace: docs/upgrade.md: 4 matches (gitmap-v28 -> gitmap-v28)
  replace: 2 files, 5 replacements
  Apply replacements for versions v4..v6 -> v7? [y/N]: y
  replace: applied 5 replacements across 2 files
```

### Audit only

```
$ gitmap replace --audit

  README.md:42: see https://github.com/x/gitmap-v28 for the legacy guide
  go.mod:3: module github.com/x/gitmap-v28
```

## Excluded paths

`.git`, `.gitmap`, `.release`, `node_modules`, `vendor`,
`.gitmap/release`, `.gitmap/release-assets`, and any file whose first
8 KiB contain a null byte (treated as binary).

## See Also

- `release-self` — Bump gitmap's own version
- `clone-next` — Clone the next versioned repo iteration

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter replace
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
