# Release Body and Changelog

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Purpose

Defines the standard patterns for extracting changelog entries and assembling structured GitHub Release descriptions. Every Go binary release MUST produce a consistently formatted release body.

---

## Changelog Extraction

### Convention

Projects maintain a `changelog.md` at the repository root following Keep a Changelog format:

```markdown

## v1.3.0 — Feature Title (2026-04-10)

### Added

- New `map` subcommand for repository visualization

### Fixed

- Crash when no Git remote is configured

## v1.2.0 — Previous Release (2026-04-08)

...
```

### Extraction Pattern

Use `awk` to extract the section matching the current version:

```bash
VERSION="${{ steps.version.outputs.version }}"
ENTRY=$(awk -v ver="$VERSION" '
    /^## / {
        if (found) exit
        if (index($0, ver)) found=1
    }
    found { print }
' changelog.md 2>/dev/null || echo "")

if [ -z "$ENTRY" ]; then
    ENTRY="Release $VERSION"
fi

echo "$ENTRY" > /tmp/changelog-entry.md
```

### How It Works

1. Scan for lines starting with `## `
2. When the version string is found, start capturing
3. When the next `## ` heading is found, stop
4. If no match, fall back to a simple "Release $VERSION" string
5. Write to a temp file for multi-line shell compatibility

### Graceful Fallback

If `changelog.md` does not exist or the version is not found, the release body still generates successfully with a minimal entry. This ensures releases never fail due to missing changelog entries.

---

## Release Body Template

The release body is assembled as a single markdown file written to `/tmp/release-body.md`:

### Structure

```markdown
<changelog entry>

---

## Release Info

| Field | Value |
|-------|-------|
| Version | `<version>` |
| Commit | `<commit-sha-10>` |
| Branch | `<branch>` |
| Build Date | <UTC timestamp> |
| Go Version | <go version> |

## Checksums (SHA256)

```
<contents of checksums.txt>
```

## Install

### Quick install (Windows PowerShell)

```powershell
irm https://github.com/<repo>/releases/download/<version>/install.ps1 | iex
```

### Quick install (Linux / macOS)

```bash
curl -fsSL https://github.com/<repo>/releases/download/<version>/install.sh | bash
```

### Manual download

Download the appropriate archive for your platform from the assets below,
extract, and place the binary in your PATH.

## Assets

| Platform | Architecture | File |
|----------|-------------|------|
| Windows | amd64 | `<binary>-<version>-windows-amd64.zip` |
| Windows | arm64 | `<binary>-<version>-windows-arm64.zip` |
| Linux | amd64 | `<binary>-<version>-linux-amd64.tar.gz` |
| Linux | arm64 | `<binary>-<version>-linux-arm64.tar.gz` |
| macOS | amd64 | `<binary>-<version>-darwin-amd64.tar.gz` |
| macOS | arm64 | `<binary>-<version>-darwin-arm64.tar.gz` |
```

---

## Full Assembly Script

The complete script for assembling the release body in a GitHub Actions step:

```bash
- name: Assemble release body
  run: |
    VERSION="${{ steps.version.outputs.version }}"
    REPO="${{ github.repository }}"
    BINARY_NAME="<binary>"
    COMMIT_SHA="${GITHUB_SHA}"
    BRANCH="${GITHUB_REF_NAME}"
    BUILD_DATE="$(date -u '+%Y-%m-%d %H:%M:%S UTC')"
    GO_VERSION="$(go version | awk '{print $3}')"

    # --- 1. Extract changelog entry ---
    ENTRY=$(awk -v ver="$VERSION" '
        /^## / {
            if (found) exit
            if (index($0, ver)) found=1
        }
        found { print }
    ' changelog.md 2>/dev/null || echo "")

    if [ -z "$ENTRY" ]; then
        ENTRY="Release $VERSION"
    fi

    echo "$ENTRY" > /tmp/changelog-entry.md

    # --- 2. Build the release body ---
    cat > /tmp/release-body.md << RELEASE_EOF
    $(cat /tmp/changelog-entry.md)

    ---

    ## Release Info

    | Field | Value |
    |-------|-------|
    | Version | \`$VERSION\` |
    | Commit | \`${COMMIT_SHA:0:10}\` |
    | Branch | \`$BRANCH\` |
    | Build Date | $BUILD_DATE |
    | Go Version | $GO_VERSION |

    ## Checksums (SHA256)

    \`\`\`
    $(cat dist/checksums.txt)
    \`\`\`

    ## Install

    ### Quick install (Windows PowerShell)

    \`\`\`powershell
    irm https://github.com/$REPO/releases/download/$VERSION/install.ps1 | iex
    \`\`\`

    ### Quick install (Linux / macOS)

    \`\`\`bash
    curl -fsSL https://github.com/$REPO/releases/download/$VERSION/install.sh | bash
    \`\`\`

    ### Install specific version (generic installer)

    \`\`\`powershell
    & { \$Version = "$VERSION"; irm https://raw.githubusercontent.com/$REPO/main/scripts/install.ps1 | iex }
    \`\`\`

    \`\`\`bash
    curl -fsSL https://raw.githubusercontent.com/$REPO/main/scripts/install.sh | bash -s -- --version $VERSION
    \`\`\`

    ### Manual download

    Download the appropriate archive for your platform from the assets below,
    extract, and place the binary in your PATH.

    ## Assets

    | Platform | Architecture | File |
    |----------|-------------|------|
    | Windows | amd64 | \`${BINARY_NAME}-${VERSION}-windows-amd64.zip\` |
    | Windows | arm64 | \`${BINARY_NAME}-${VERSION}-windows-arm64.zip\` |
    | Linux | amd64 | \`${BINARY_NAME}-${VERSION}-linux-amd64.tar.gz\` |
    | Linux | arm64 | \`${BINARY_NAME}-${VERSION}-linux-arm64.tar.gz\` |
    | macOS | amd64 | \`${BINARY_NAME}-${VERSION}-darwin-amd64.tar.gz\` |
    | macOS | arm64 | \`${BINARY_NAME}-${VERSION}-darwin-arm64.tar.gz\` |
    RELEASE_EOF

    # --- 3. Validate release body ---
    test -s /tmp/release-body.md || { echo "::error::Release body is empty"; exit 1; }
    echo "✅ Release body generated ($(wc -l < /tmp/release-body.md) lines)"
```

---

## Pre-Release Body Differences

When the version contains a pre-release suffix (e.g., `v1.3.0-beta.1`), the release body adds a warning banner:

```markdown
> ⚠️ **This is a pre-release.** It may contain breaking changes and is not recommended for production use.
```

Detection logic in the assembly step:

```bash
IS_PRERELEASE="false"
if [[ "$VERSION" == *-* ]]; then
    IS_PRERELEASE="true"
fi

# Prepend warning if pre-release

if [ "$IS_PRERELEASE" = "true" ]; then
    BANNER='> ⚠️ **This is a pre-release.** It may contain breaking changes and is not recommended for production use.\n\n'
    sed -i "1s/^/$BANNER/" /tmp/release-body.md
fi
```

The `gh release create` command must also pass `--prerelease` when a pre-release is detected:

```bash
gh release create "$VERSION" dist/* \
    --repo "$REPO" \
    --title "$VERSION" \
    --notes-file /tmp/release-body.md \
    $([ "$IS_PRERELEASE" = "true" ] && echo "--prerelease")
```

---

## Validation Step

After assembly, validate the release body contains all required sections:

```bash
- name: Validate release body
  run: |
    BODY="/tmp/release-body.md"
    ERRORS=0

    for section in "Release Info" "Checksums" "Install" "Assets"; do
        if ! grep -q "$section" "$BODY"; then
            echo "::error::Missing section: $section"
            ERRORS=$((ERRORS + 1))
        fi
    done

    if ! grep -q "SHA256" "$BODY"; then
        echo "::warning::Checksums section may be empty"
    fi

    LINE_COUNT=$(wc -l < "$BODY")
    if [ "$LINE_COUNT" -lt 20 ]; then
        echo "::warning::Release body seems too short ($LINE_COUNT lines)"
    fi

    if [ "$ERRORS" -gt 0 ]; then
        echo "::error::Release body validation failed with $ERRORS errors"
        exit 1
    fi

    echo "✅ Release body validated successfully ($LINE_COUNT lines)"
```

---

## Generic Installer References

For projects with both version-pinned and generic installers, the release body includes both:

```markdown

### Install specific version (generic installer)

```powershell
irm https://raw.githubusercontent.com/<repo>/main/scripts/install.ps1 | iex

# Or pin this version:

& { $Version = "<version>"; irm https://raw.githubusercontent.com/<repo>/main/scripts/install.ps1 | iex }
```

```bash
curl -fsSL https://raw.githubusercontent.com/<repo>/main/scripts/install.sh | bash

# Or pin this version:

curl -fsSL https://raw.githubusercontent.com/<repo>/main/scripts/install.sh | bash -s -- --version <version>
```
```

---

## Release Info Table Fields

| Field | Source | Example |
|-------|--------|---------|
| Version | `steps.version.outputs.version` | `v1.3.0` |
| Commit | `${GITHUB_SHA:0:10}` (first 10 chars) | `a1b2c3d4e5` |
| Branch | `${GITHUB_REF_NAME}` | `release/v1.3.0` |
| Build Date | `date -u '+%Y-%m-%d %H:%M:%S UTC'` | `2026-04-10 14:30:00 UTC` |
| Go Version | `go version \| awk '{print $3}'` | `go1.25.5` |

---

## Constraints

- **Always include checksums** — the full `checksums.txt` content in a fenced code block
- **Always include an asset matrix** — maps platform/architecture to filename
- **Graceful changelog fallback** — never fail a release because `changelog.md` is missing
- **UTC timestamps** — all dates in the release body use UTC
- **Consistent field order** — Version, Commit, Branch, Build Date, Go Version (always this order)
- **Pre-release banner** — auto-detected from version suffix and prepended to body
- **Validation gate** — release body must pass section validation before `gh release create`

---

## Cross-References

- [GitHub Release Standard](./04-github-release-standard.md) — Pre-release detection, asset rules
- [Shared Conventions](./03-shared-conventions.md) — Version resolution
- [Install Script Generation](./07-install-script-generation.md) — Scripts referenced in the install section
- [Go Binary Release Pipeline](./05-release-pipeline.md) — Full release workflow
- [Release Versioning](../14-update/16-release-versioning.md) — How versions are resolved and tags created

---

*Release body and changelog — v3.2.0 — 2026-04-10*
