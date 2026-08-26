# Browser Extension — Release Pipeline

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

The release pipeline automates extension packaging and GitHub Release creation whenever code is pushed to a `release/**` branch or a `v*` tag. It produces versioned zip archives for the extension and all standalone components.

---

## Trigger and Concurrency

```yaml
on:
  push:
    branches:
      - "release/**"
    tags:
      - "v*"

concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false

permissions:
  contents: write
```

Release builds must **NEVER** be cancelled — every release commit must produce a GitHub Release regardless of newer commits.

---

## Job Graph

```
setup (lint + test + version resolution)
  └── build-sdk
        ├── build-module-a (parallel)
        ├── build-module-b (parallel)
        └── build-module-c (parallel)
              └── release (build extension + package + GitHub Release)
```

The build graph is identical to CI, except the final job also packages assets and creates the release.

---

## Version Resolution

Version is resolved in the setup job and passed downstream via `outputs`:

```yaml
setup:
  outputs:
    version: ${{ steps.version.outputs.version }}
```

See [Shared Conventions — Version Resolution](../01-shared-conventions.md#version-resolution) for the resolution logic.

---

## Packaging

### Source Map Removal

Source maps MUST be removed before packaging — they must never ship in releases:

```bash
MAP_COUNT=$(find dist/ -name '*.map' | wc -l | tr -d ' ')
find dist/ -name '*.map' -delete
echo "Removed ${MAP_COUNT} source map files"
```

### Extension Archive

The extension `dist/` contents are zipped **without** the `dist/` prefix so the archive extracts cleanly:

```bash
cd extension/dist
zip -r "../../release-assets/<extension>-${VERSION}.zip" .
cd ../..
```

### Standalone Component Archives

Each standalone module is zipped from its own directory:

```bash
cd standalone-scripts/<module-name>
zip -r "../../release-assets/<module-name>-${VERSION}.zip" .
cd ../..
```

### Additional Release Assets

| Asset | Source |
|-------|--------|
| `install.sh` | `scripts/install.sh` — Bash installer |
| `install.ps1` | `scripts/install.ps1` — PowerShell installer |
| `VERSION.txt` | Contains the resolved version string |
| `changelog.md` | Full project changelog (if exists) |
| `checksums.txt` | SHA-256 checksums of all assets |

---

## Release Notes Generation

Release notes are assembled from conventional commit history. See [GitHub Release Standard](../02-github-release-standard.md) for the full format.

The release body includes:
1. Version header
2. Release info table (version, commit SHA, branch, build date)
3. Categorized commit log (Features, Bug Fixes, Maintenance)
4. SHA-256 checksums block
5. Asset description table
6. Quick install instructions (PowerShell and Bash one-liners)
7. Manual install steps for Chromium browsers

---

## GitHub Release Creation

```yaml
- uses: softprops/action-gh-release@v2
  with:
    tag_name: ${{ needs.setup.outputs.version }}
    name: "<Extension Name> ${{ needs.setup.outputs.version }}"
    body_path: release-assets/RELEASE_NOTES.md
    files: release-assets/*
    draft: false
    prerelease: ${{ contains(needs.setup.outputs.version, '-') }}
    make_latest: ${{ !contains(needs.setup.outputs.version, '-') }}
```

---

## Build-Once Rule

Binaries and extension bundles are compiled **exactly once**. Compression, checksums, and publishing operate on the already-built artifacts and must **never** trigger a rebuild.

---

## Constraints

- Source maps are removed before packaging — never ship `.map` files
- Extension zip must not include a `dist/` prefix directory
- All assets pass through checksum generation
- Release notes are written to a file, not inline YAML
- The same build graph as CI is used — no shortcuts in release builds
- Every release commit runs to completion — never cancelled

---

*Browser extension release pipeline — updated: 2026-04-09*
