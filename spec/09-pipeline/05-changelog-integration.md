# Changelog Integration

## Overview

The changelog is a critical piece of the release pipeline. It must be synchronized across three surfaces: the `CHANGELOG.md` file, the GitHub Release description, and (optionally) a documentation site. This document describes the format, automation, and integration points.

---

## CHANGELOG.md Format

Use [Keep a Changelog](https://keepachangelog.com/) conventions with SemVer headers:

```markdown
## v1.3.0 — 2026-04-09

### Added
- New `env` command for cross-platform environment variable management.
- Shell completion for Fish.

### Fixed
- Fixed PATH registration on Git Bash when `.bash_profile` does not exist.

### Changed
- Installer now prints version at both start and end of installation.

---

## v1.2.3 — 2026-04-01

### Fixed
- Fixed checksum verification failure on ARM64 macOS.
```

### Rules

- Version headers use `## v<version> — <YYYY-MM-DD>` format
- Subsections: `### Added`, `### Fixed`, `### Changed`, `### Removed`
- Each bullet is a single, descriptive line (no multi-line bullets)
- No duplicate content between versions — each entry reflects only that release's changes
- The topmost entry is always the latest version
- Date is the release date, not the development date

---

## Adding a Changelog Entry

### When to Add

Add a changelog entry **before** creating a release branch or tag. The entry must already exist in `CHANGELOG.md` on the commit that gets tagged.

### How to Add

1. Open `CHANGELOG.md`
2. Insert a new `## v<version>` section at the top (below any header)
3. Add bullets under the appropriate subsection (`Added`, `Fixed`, `Changed`)
4. Commit the changelog update as part of the release preparation

### Example Workflow

```bash
# 1. Bump version constant in code
# 2. Update CHANGELOG.md with new entry
# 3. Commit both changes
git add constants/constants.go CHANGELOG.md
git commit -m "Prepare v1.3.0 release"

# 4. Create release branch
git checkout -b release/v1.3.0
git push origin release/v1.3.0
```

---

## Changelog Extraction in CI

The release pipeline extracts the changelog entry for the current version and includes it in the GitHub Release body.

### Extraction Script

```bash
VERSION="${1:?Usage: extract-changelog.sh <version>}"

# Strip 'v' prefix for matching
CLEAN_VERSION="${VERSION#v}"

CHANGELOG=$(awk -v ver="$CLEAN_VERSION" '
  /^## / {
    if (found) exit
    if (index($0, ver)) found=1
  }
  found { print }
' CHANGELOG.md)

if [ -z "$CHANGELOG" ]; then
  CHANGELOG="No changelog entry found for $VERSION."
fi

echo "$CHANGELOG"
```

### Integration Point

In the release workflow, call this after version resolution:

```yaml
- name: Extract changelog
  id: changelog
  run: |
    CHANGELOG=$(bash .github/scripts/extract-changelog.sh "${{ steps.version.outputs.version }}")
    # Write to file for multi-line output
    echo "$CHANGELOG" > /tmp/changelog-entry.md
```

---

## GitHub Release Body Structure

The release description assembles multiple sections:

```markdown
## What's New

<extracted changelog entry>

---

## Release Info

| Field | Value |
|-------|-------|
| Version | v1.3.0 |
| Commit | abc1234 |
| Branch | release/v1.3.0 |
| Built | 2026-04-09 14:30 UTC |
| Go | 1.24.2 |

---

## Install

**PowerShell (Windows)**
\`\`\`powershell
irm https://github.com/<owner>/<repo>/releases/download/v1.3.0/install.ps1 | iex
\`\`\`

**Bash (Linux / macOS)**
\`\`\`bash
curl -fsSL https://github.com/<owner>/<repo>/releases/download/v1.3.0/install.sh | bash
\`\`\`

---

## Checksums

\`\`\`
<sha256 checksums>
\`\`\`

---

## Assets

| Platform | Architecture | File |
|----------|-------------|------|
| Windows | amd64 | <tool>-v1.3.0-windows-amd64.zip |
| Windows | arm64 | <tool>-v1.3.0-windows-arm64.zip |
| Linux | amd64 | <tool>-v1.3.0-linux-amd64.tar.gz |
| Linux | arm64 | <tool>-v1.3.0-linux-arm64.tar.gz |
| macOS | amd64 | <tool>-v1.3.0-darwin-amd64.tar.gz |
| macOS | arm64 | <tool>-v1.3.0-darwin-arm64.tar.gz |
```

### Assembly Script

```bash
VERSION="$1"
REPO="$2"

cat > /tmp/release-body.md << EOF
## What's New

$(cat /tmp/changelog-entry.md)

---

## Release Info

| Field | Value |
|-------|-------|
| Version | $VERSION |
| Commit | ${GITHUB_SHA::10} |
| Branch | ${GITHUB_REF_NAME} |
| Built | $(date -u +"%Y-%m-%d %H:%M UTC") |
| Go | $(go version | awk '{print $3}' | sed 's/go//') |

---

## Install

**PowerShell (Windows)**
\`\`\`powershell
irm https://github.com/$REPO/releases/download/$VERSION/install.ps1 | iex
\`\`\`

**Bash (Linux / macOS)**
\`\`\`bash
curl -fsSL https://github.com/$REPO/releases/download/$VERSION/install.sh | bash
\`\`\`

---

## Checksums

\`\`\`
$(cat dist/checksums.txt)
\`\`\`

---

## Assets

| Platform | Architecture | File |
|----------|-------------|------|
$(for f in dist/*.{zip,tar.gz}; do
  name=$(basename "$f")
  # Parse platform and arch from filename
  echo "| ... | ... | $name |"
done)
EOF
```

---

## Documentation Site Changelog

If the project has a documentation site, maintain a parallel changelog data file:

```typescript
// src/data/changelog.ts
export const changelog = [
  {
    version: "1.3.0",
    date: "2026-04-09",
    changes: [
      { type: "added", description: "New `env` command for environment variable management." },
      { type: "fixed", description: "Fixed PATH registration on Git Bash." },
    ],
  },
  // ... older versions
];
```

Both `CHANGELOG.md` and the data file must be updated together during release preparation.

---

## Terminal Output — Release Printing

When the CLI has a `list-releases` or `release info` command, format the output consistently:

```
$ <tool> list-releases

  VERSION    DATE          BRANCH              TAG
  v1.3.0     2026-04-09    release/v1.3.0      v1.3.0
  v1.2.3     2026-04-01    release/v1.2.3      v1.2.3
  v1.2.2     2026-03-28    release/v1.2.2      v1.2.2

$ <tool> release info v1.3.0

  Version:    v1.3.0
  Branch:     release/v1.3.0
  Tag:        v1.3.0
  Commit:     abc1234def5
  Date:       2026-04-09

  Changelog:
    - New `env` command for cross-platform environment variable management.
    - Shell completion for Fish.
    - Fixed PATH registration on Git Bash.
```

---

## Constraints

- Changelog entry MUST exist before the release branch/tag is created
- The CI pipeline MUST extract and print the changelog in the release body
- No hardcoded version strings in the changelog extraction script
- Changelog format must be parseable by `awk` (line-based, `## v` headers)
- Date format is `YYYY-MM-DD` (ISO 8601)
- Each version entry is self-contained — no references to "same as above"
