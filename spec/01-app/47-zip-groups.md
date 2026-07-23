# Zip Groups — Structured Release Asset Bundling

## Overview

Zip groups allow users to define named collections of files and folders
that are automatically compressed into a single archive during a release.
Groups can be persisted in the database for reuse across releases, or
defined ad-hoc via flags on `gitmap-v27 release`.

---

## Goals

1. **Persistent groups** — save named file/folder collections in SQLite
   for repeated use across releases.
2. **Ad-hoc groups** — define one-off bundles inline during a release.
3. **Best compression** — always use the highest compression level
   available (deflate level 9 for zip).
4. **Custom naming** — optionally name the output archive; defaults to
   `<group-name>_<version>.zip`.
5. **Release integration** — zip groups are resolved and compressed
   during the release workflow, and the resulting archives are attached
   as release assets.
6. **Metadata persistence** — zip group definitions are stored in the
   `.gitmap/release/vX.Y.Z.json` metadata under a `zipGroups` key.

---

## Command: `gitmap-v27 zip-group`

Alias: `z`

Manages persistent zip group definitions stored in SQLite.

### Subcommands

| Subcommand | Alias | Description |
|------------|-------|-------------|
| `create`   | —     | Create a named zip group |
| `add`      | —     | Add files or folders to a group |
| `remove`   | —     | Remove an item from a group |
| `list`     | —     | List all zip groups |
| `show`     | —     | Show items in a group |
| `delete`   | —     | Delete a zip group |
| `rename`   | —     | Rename the output archive |

### Usage

```
gitmap-v27 zip-group create <name> [--archive <filename.zip>]
gitmap-v27 zip-group add <name> <path> [<path>...]
gitmap-v27 zip-group remove <name> <path>
gitmap-v27 zip-group list
gitmap-v27 zip-group show <name>
gitmap-v27 zip-group delete <name>
gitmap-v27 zip-group rename <name> --archive <filename.zip>
```

### Examples

```bash
# Create a group and add items
gitmap-v27 z create docs-bundle
gitmap-v27 z add docs-bundle ./README.md ./CHANGELOG.md ./docs/

# Create with custom archive name
gitmap-v27 z create extras --archive extra-files.zip
gitmap-v27 z add extras ./config/ ./scripts/deploy.sh

# List groups
gitmap-v27 z list

# Show group contents
gitmap-v27 z show docs-bundle

# Delete a group
gitmap-v27 z delete extras
```

---

## Release Integration

### Flag: `--zip-group`

Activates one or more persistent zip groups during a release. Each group
produces a single `.zip` archive attached as a release asset.

```bash
gitmap-v27 release v3.0.0 --zip-group docs-bundle
gitmap-v27 release v3.0.0 --zip-group docs-bundle --zip-group extras
```

### Flag: `-Z` / `--zip`

Defines ad-hoc items to zip during release. Behavior is controlled by
the `--bundle` flag.

```bash
# Each item becomes its own archive
gitmap-v27 release v3.0.0 -Z ./dist/report.pdf -Z ./dist/manual.pdf

# Bundle all ad-hoc items into one archive
gitmap-v27 release v3.0.0 -Z ./dist/report.pdf -Z ./dist/manual.pdf --bundle docs.zip
```

### Flag: `--bundle`

Controls how ad-hoc `-Z` items are packaged.

| Value | Behavior |
|-------|----------|
| (absent) | Each `-Z` item is zipped individually |
| `<name>.zip` | All `-Z` items bundled into a single named archive |

### Combined Usage

Persistent groups and ad-hoc items can be used together:

```bash
gitmap-v27 release v3.0.0 --zip-group docs-bundle -Z ./extras/notes.txt
```

---

## Compression Strategy

All archives use **ZIP format with Deflate level 9** (maximum compression).
This applies to both persistent groups and ad-hoc items.

When a `-Z` target is a folder, the entire directory tree is added
recursively to the archive, preserving relative paths.

When a `-Z` target is a single file, it produces `<filename>.zip`
unless `--bundle` is used.

---

## Data Model

### SQLite Tables

#### `ZipGroups`

| Column      | Type    | Constraints                    |
|-------------|---------|--------------------------------|
| Id          | TEXT    | PRIMARY KEY                    |
| Name        | TEXT    | NOT NULL UNIQUE                |
| ArchiveName | TEXT    | DEFAULT '' (custom output name)|
| CreatedAt   | TEXT    | DEFAULT CURRENT_TIMESTAMP      |

#### `ZipGroupItems`

| Column      | Type    | Constraints                        |
|-------------|---------|-------------------------------------|
| GroupId     | TEXT    | NOT NULL REFERENCES ZipGroups(Id) ON DELETE CASCADE |
| Path        | TEXT    | NOT NULL                            |
| IsFolder    | INTEGER | DEFAULT 0                           |
| PRIMARY KEY | —       | (GroupId, Path)                      |

### Release Metadata (`.gitmap/release/vX.Y.Z.json`)

The `zipGroups` key records which groups were used and the resulting
archive filenames:

```json
{
  "version": "3.0.0",
  "tag": "v3.0.0",
  "zipGroups": [
    {
      "name": "docs-bundle",
      "archive": "docs-bundle_v3.0.0.zip",
      "items": ["README.md", "CHANGELOG.md", "docs/"]
    }
  ],
  "assets": [
    "gitmap_v3.0.0_windows_amd64.exe.zip",
    "docs-bundle_v3.0.0.zip"
  ]
}
```

### Model Structs

```go
// ZipGroup represents a named collection of files/folders for archiving.
type ZipGroup struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    ArchiveName string `json:"archiveName"`
    CreatedAt   string `json:"createdAt"`
}

// ZipGroupItem links a file or folder path to a zip group.
type ZipGroupItem struct {
    GroupID  string `json:"groupId"`
    Path     string `json:"path"`
    IsFolder bool   `json:"isFolder"`
}
```

---

## Archive Naming

| Scenario | Output filename |
|----------|-----------------|
| Persistent group, no custom name | `<group-name>_<version>.zip` |
| Persistent group, custom name | `<archive-name>` (as-is, version not appended) |
| Ad-hoc single file, no `--bundle` | `<filename>.zip` |
| Ad-hoc folder, no `--bundle` | `<foldername>.zip` |
| Ad-hoc with `--bundle name.zip` | `name.zip` |

---

## Workflow Integration

The zip group step runs **after** Go cross-compilation (step 7a) and
**before** asset upload (step 7b) in the release lifecycle:

### Step 7a — Cross-Compile (existing)
### Step 7a′ — Zip Groups
1. Resolve all `--zip-group` references from the database.
2. Resolve all `-Z` ad-hoc items.
3. Validate all paths exist (warn and skip missing items).
4. For each group/bundle, create a ZIP archive at deflate level 9.
5. Place archives in the `release-assets/` staging directory.

### Step 7b — Upload Assets (existing, unchanged)

---

## Dry-Run Support

```
[dry-run] Would create 2 zip archive(s):
  → docs-bundle_v3.0.0.zip (3 items: README.md, CHANGELOG.md, docs/)
  → report.pdf.zip (1 item: dist/report.pdf)
[dry-run] Would upload 8 assets + checksums.txt
```

---

## Error Handling

- Missing file/folder in a group → warn, skip item, continue.
- Empty group (all items missing) → warn, skip group entirely.
- Group name not found in DB → error, abort release.
- Zip creation failure → warn, continue with remaining groups.
- Never abort the entire release because of a zip group failure.

---

## Package Structure

| File | Responsibility |
|------|----------------|
| `cmd/zipgroup.go` | Subcommand dispatch (create/add/remove/list/show/delete) |
| `cmd/zipgroupops.go` | Subcommand implementation |
| `release/ziparchive.go` | ZIP creation with max compression |
| `store/zipgroup.go` | Database CRUD for ZipGroups/ZipGroupItems |
| `model/zipgroup.go` | Data structs |
| `constants/constants_zipgroup.go` | Messages, SQL, flag descriptions |
| `helptext/zip-group.md` | Command help |

---

## Acceptance Criteria

1. `gitmap-v27 z create docs` creates a persistent zip group in the database.
2. `gitmap-v27 z add docs ./README.md ./docs/` adds items to the group.
3. `gitmap-v27 z show docs` displays all items with file/folder indicators.
4. `gitmap-v27 z list` lists all zip groups with item counts.
5. `gitmap-v27 z delete docs` removes the group and all its items.
6. `gitmap-v27 release v1.0.0 --zip-group docs` creates `docs_v1.0.0.zip` and attaches it.
7. `gitmap-v27 release v1.0.0 -Z ./file.txt` creates `file.txt.zip` and attaches it.
8. `gitmap-v27 release v1.0.0 -Z ./a.txt -Z ./b.txt --bundle extras.zip` bundles both into one archive.
9. `--dry-run` lists planned archives without creating them.
10. Missing items produce warnings but do not abort the release.
11. Zip archives use maximum compression (deflate level 9).
12. Release metadata JSON includes `zipGroups` array.
13. Persistent groups and ad-hoc items can be combined in a single release.
