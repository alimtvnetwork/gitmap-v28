# GitHub Release Standard

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

This document defines the standard format and assembly process for GitHub Releases across all project types. Every release pipeline produces a GitHub Release with a structured body, versioned assets, and install instructions.

---

## Release Body Assembly

The release description is assembled from multiple data sources into a single Markdown file:

```
1. Changelog entry (extracted from changelog.md or generated from commits)
2. Release info table (version, commit SHA, branch, build date, toolchain version)
3. Checksums block (full checksums.txt in a code fence)
4. Install instructions (quick install one-liners for each platform)
5. Asset matrix (table mapping platform/architecture to filenames)
```

### Release Info Table

```markdown
| Field | Value |
|-------|-------|
| Version | `v1.2.3` |
| Commit | `abc1234567` |
| Branch | `release/v1.2.3` |
| Build Date | 2026-04-09 12:00:00 UTC |
```

### Checksums Block

```markdown

### Checksums (SHA256)

\`\`\`
<contents of checksums.txt>
\`\`\`
```

---

## Changelog Extraction

### From changelog.md

Extract the matching version section using `awk`:

```bash
awk -v ver="$VERSION" '
  /^## / {
    if (found) exit
    if (index($0, ver)) found=1
  }
  found { print }
' changelog.md
```

Falls back to a "No changelog entry found" message if the version section does not exist.

### From Commit History (Conventional Commits)

When no `changelog.md` exists, generate release notes from Git history:

```bash
PREV_TAG=$(git tag --sort=-version:refname | grep -E '^v[0-9]' | head -1 || true)
RANGE="${PREV_TAG:+${PREV_TAG}..HEAD}"

# Categorize by conventional commit prefix

FEATS=$(git log $RANGE --pretty=format:"- %s (%h)" --grep="^feat" -i || true)
FIXES=$(git log $RANGE --pretty=format:"- %s (%h)" --grep="^fix" -i || true)
OTHER=$(git log $RANGE --pretty=format:"- %s (%h)" --grep="^refactor\|^chore\|^docs" -i || true)
```

Categories rendered:

- **Features** — `feat:` commits
- **Bug Fixes** — `fix:` commits
- **Maintenance** — `refactor:`, `chore:`, `docs:`, `style:`, `perf:`, `test:`, `ci:`, `build:` commits

---

## Release Type Detection

| Condition | Release Type | `prerelease` | `make_latest` |
|-----------|-------------|--------------|---------------|
| Version contains `-` (e.g., `v1.0.0-beta.1`) | Pre-release | `true` | `false` |
| Version is clean (e.g., `v1.0.0`) | Stable | `false` | `true` |

---

## GitHub Release Action

```yaml
- uses: softprops/action-gh-release@v2
  with:
    tag_name: ${{ steps.version.outputs.version }}
    name: "<Project Name> ${{ steps.version.outputs.version }}"
    body_path: /tmp/release-body.md
    files: dist/*
    draft: false
    prerelease: ${{ contains(steps.version.outputs.version, '-') }}
    make_latest: ${{ !contains(steps.version.outputs.version, '-') }}
```

---

## Install Instructions Template

Every release body includes quick install commands:

### For Binary Projects

**PowerShell (Windows):**
```
irm https://raw.githubusercontent.com/<owner>/<repo>/main/scripts/install.ps1 | iex
```

**Bash (Linux/macOS):**
```
curl -fsSL https://raw.githubusercontent.com/<owner>/<repo>/main/scripts/install.sh | bash
```

**Version-pinned variants** allow downloading a specific version by setting a variable before piping.

### For Browser Extensions

```
1. Download `<extension>-<version>.zip` from this release.
2. Extract it to a local folder.
3. Open `chrome://extensions`.
4. Enable **Developer mode**.
5. Click **Load unpacked** and select the extracted folder.
```

---

## Constraints

- Release body is assembled as a Markdown file, never inline YAML strings
- Changelog extraction must handle missing version sections gracefully
- Pre-release detection is version-string based (contains `-`), not manual
- All release assets pass through checksum generation before upload

---

## Cross-References

- [Shared Conventions](./03-shared-conventions.md) — Version resolution, checksum generation
- [Release Body and Changelog](./14-release-body-and-changelog.md) — Detailed changelog extraction and body template
- [Install Script Generation](./07-install-script-generation.md) — Scripts referenced in install instructions
- [Code Signing](./10-code-signing.md) — Signed binaries in release assets
- [Go Binary Release Pipeline](./05-release-pipeline.md) — Go-specific release implementation

---

*GitHub Release standard — updated: 2026-04-10*
