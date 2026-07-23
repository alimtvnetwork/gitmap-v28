# gitmap zip-group

Manage named collections of files and folders that are automatically
compressed into ZIP archives during a release.

## Alias

z

## Usage

    gitmap zip-group <subcommand> [arguments]

## Subcommands

| Subcommand | Description |
|------------|-------------|
| create     | Create a named zip group (optionally with paths) |
| add        | Add files or folders to a group |
| remove     | Remove an item from a group |
| list       | List all zip groups |
| show       | Show items in a group (folders expanded dynamically) |
| delete     | Delete a zip group |
| rename     | Set a custom archive name for a group |

## Flags

| Flag | Description |
|------|-------------|
| --archive \<name\> | Custom output filename (used with create/rename) |

## Path Resolution

When adding paths, gitmap resolves them using:
1. **Repo path** — current working directory
2. **Relative path** — the path you provide
3. **Full path** — repo path + relative path combined

If the resolved path is a directory, only the folder reference is stored.
Files within the folder are expanded at runtime during `show` and archive creation.

## Storage

Zip groups are persisted in two locations:
1. **SQLite database** — primary storage with full metadata
2. **.gitmap/zip-groups.json** — JSON mirror for version control

## Prerequisites

- Must be inside a Git repository with release workflow configured (see release.md)

## Examples

### Example 1: Create a group with paths in one step

    gitmap z create "chrome extension" chrome-extension/dist

**Output:**

    ✓ Added chrome-extension/dist to "chrome extension" (folder)
    ✓ Created zip group "chrome extension" with 1 item(s)

### Example 2: Create a group and add items separately

    gitmap z create docs-bundle
    gitmap z add docs-bundle ./README.md ./CHANGELOG.md ./docs/

**Output:**

    ✓ Created zip group "docs-bundle"

    ✓ Added ./README.md to "docs-bundle" (file)
    ✓ Added ./CHANGELOG.md to "docs-bundle" (file)
    ✓ Added ./docs/ to "docs-bundle" (folder)

### Example 3: Create with custom archive name

    gitmap z create extras --archive extra-files.zip
    gitmap z add extras ./config/ ./scripts/deploy.sh

**Output:**

    ✓ Created zip group "extras" (archive: extra-files.zip)

    ✓ Added ./config/ to "extras" (folder)
    ✓ Added ./scripts/deploy.sh to "extras" (file)

### Example 4: List all zip groups

    gitmap z list

**Output:**

    GROUP           ITEMS   ARCHIVE NAME
    docs-bundle     3       docs-bundle.zip
    extras          2       extra-files.zip
    2 zip groups defined

### Example 5: Show items with dynamic folder expansion

    gitmap z show docs-bundle

**Output:**

    Zip group: docs-bundle (3 items):

      📄 ./README.md
        repo:     D:\projects\myapp
        relative: ./README.md
        full:     D:\projects\myapp\README.md
      📁 ./docs/
        repo:     D:\projects\myapp
        relative: ./docs/
        full:     D:\projects\myapp\docs
        Contents (4 files):
          getting-started.md
          api-reference.md
          faq.md
          changelog.md

### Example 6: Use during release

    gitmap release v3.0.0 --zip-group docs-bundle

**Output:**

    Creating tag v3.0.0... done
    ✓ Compressed docs-bundle → docs-bundle_v3.0.0.zip (3 items)
    Uploading to GitHub... done
    ✓ Released v3.0.0

## See Also

- [release](release.md) — Create a release with zip group assets
- [group](group.md) — Manage repository groups

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter zip-group
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
