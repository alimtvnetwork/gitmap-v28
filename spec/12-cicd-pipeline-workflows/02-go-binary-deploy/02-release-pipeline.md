# Go Binary — Release Pipeline

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

The release pipeline automates binary production, packaging, and GitHub Release creation whenever code is pushed to a `release/**` branch or a `v*` tag. It produces cross-compiled binaries, platform-specific install scripts, checksums, and a fully formatted release page.

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

Release pipelines **never cancel** in-progress runs — every release commit must produce complete artifacts.

---

## Version Resolution

See [Shared Conventions — Version Resolution](../01-shared-conventions.md#version-resolution).

---

## Binary Building

### Targets

Build for 6 platform/architecture combinations:

| Platform | Architecture | Extension |
|----------|-------------|-----------|
| Windows | amd64 | `.exe` |
| Windows | arm64 | `.exe` |
| Linux | amd64 | (none) |
| Linux | arm64 | (none) |
| macOS | amd64 | (none) |
| macOS | arm64 | (none) |

### Build Command

```bash
VERSION="<resolved version>"
COMMIT="${GITHUB_SHA:0:10}"
BUILD_DATE="$(date -u '+%Y-%m-%d')"
LDFLAGS="-s -w \
    -X '<module>/version.Version=$VERSION' \
    -X '<module>/version.Commit=$COMMIT' \
    -X '<module>/version.BuildDate=$BUILD_DATE'"
CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -ldflags "$LDFLAGS" -o "dist/<binary>-${VERSION}-${os}-${arch}${ext}" .
```

### LDFLAGS Variables

| Variable | Source | Purpose |
|----------|--------|---------|
| `Version` | Git ref resolution | Semantic version string |
| `Commit` | `${GITHUB_SHA:0:10}` | Short commit hash for traceability |
| `BuildDate` | `date -u '+%Y-%m-%d'` | UTC build date |

### Windows Icon & Version Embedding

Before compilation, generate Windows resource files with `go-winres make` (not `init`). The `make` command injects the release version into the binary's Properties dialog:

```yaml
- name: Generate Windows resources
  working-directory: <module-dir>
  run: |
    go install github.com/tc-hib/go-winres@v0.3.3
    CLEAN_VERSION="${VERSION#v}"
    go-winres make --product-version "$CLEAN_VERSION" --file-version "$CLEAN_VERSION"
    ls -la rsrc_windows_*.syso
```

This generates `.syso` resource files that `go build` links automatically when `GOOS=windows`. The `winres.json` in the module root defines the icon and manifest. See [Binary Icon Branding](../09-binary-icon-branding.md) for full `winres.json` schema.

### Multiple Binaries (Multi-Module Build)

When the project produces multiple binaries (e.g., main tool + updater), each is a **separate Go module** with its own `go.mod`. The CI pipeline builds each from its module directory. All outputs converge into the **main module's `dist/` folder**.

#### Directory Layout

```
<repo-root>/
├── <binary>/                  # Main binary Go module
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── winres.json
│   └── dist/                  # ALL build outputs go here
├── <binary>-updater/          # Updater Go module (separate go.mod)
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── winres.json
└── docs-site/                 # Optional documentation site
    └── package.json
```

#### Build Commands

```bash

# Main binary — builds to its own dist/

cd <binary>
CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
  go build -ldflags "$LDFLAGS_MAIN" -o "dist/<binary>-${os}-${arch}${ext}" .

# Updater — builds to the MAIN module's dist/ (note: ../<binary>/dist/)

cd <binary>-updater
CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
  go build -ldflags "$LDFLAGS_UPDATER" -o "../<binary>/dist/<binary>-updater-${os}-${arch}${ext}" .
```

#### LDFLAGS Per Binary

| Binary | Embedded Constants |
|--------|-------------------|
| Main | `Version`, `Commit`, `BuildDate`, `RepoPath` |
| Updater | `Version`, `Commit`, `RepoPath`, `BinaryName` |

#### Go Toolchain Resolution

When using `actions/setup-go`, each module resolves its Go version from its own `go.mod`:

```yaml
- uses: actions/setup-go@v5
  with:
    go-version-file: <binary>/go.mod
    cache-dependency-path: <binary>/go.sum
```

For the updater build step, the same Go version is reused (already installed).

See [Updater Binary](../../14-update/19-updater-binary.md) for the full updater architecture.

---

## Docs-Site Bundling

If the project includes a documentation site (e.g., a Node.js/Vite app), it is built and packaged as a release asset alongside the binaries.

### Build Steps

```yaml
- name: Setup Node.js
  uses: actions/setup-node@v4
  with:
    node-version: "20"

- name: Build docs-site
  working-directory: docs-site
  run: |
    npm ci
    npm run build

- name: Package docs-site
  run: |
    cd docs-site/dist
    zip -r "../../<binary>/dist/docs-site.zip" .
```

### Install Script Integration

Both `install.ps1` and `install.sh` download and extract `docs-site.zip` alongside the binary:

```bash

# In install.sh

DOCS_URL="https://github.com/$REPO/releases/download/$VERSION/docs-site.zip"
curl -fsSL "$DOCS_URL" -o "$TMP_DIR/docs-site.zip" 2>/dev/null || true
if [ -f "$TMP_DIR/docs-site.zip" ]; then
  unzip -qo "$TMP_DIR/docs-site.zip" -d "$INSTALL_DIR/docs-site"
fi
```

The docs-site download is **optional** — if it fails, the install continues without it.

### Checksums

`docs-site.zip` is included in `checksums.txt` alongside the binary archives

---

## Compression and Checksums

### Compression Rules

| Platform | Format | Command |
|----------|--------|---------|
| Windows | `.zip` | `zip "${name}.zip" "$file"` |
| Linux, macOS | `.tar.gz` | `tar czf "${name}.tar.gz" "$file"` |

After compression, remove the raw binaries — only archives ship as release assets.

### Checksums

```bash
cd dist/
sha256sum * > checksums.txt
```

The `checksums.txt` file is included as a release asset for verification.

---

## Code Signing

See [Code Signing](../05-code-signing.md) for the full SignPath integration pattern. Code signing occurs **after** binary building and **before** compression.

---

## Install Scripts

See [Install Script Generation](../04-install-script-generation.md) for the full pattern. The pipeline generates version-pinned install scripts using placeholder substitution:

```bash
sed -i "s|VERSION_PLACEHOLDER|$VERSION|g" dist/install.ps1
sed -i "s|REPO_PLACEHOLDER|$REPO|g" dist/install.ps1
chmod +x dist/install.sh
```

### PowerShell Installer (`install.ps1`)

Features:
- Auto-detect CPU architecture (amd64/arm64)
- Download versioned `.zip` from GitHub Releases
- Verify SHA-256 checksum against `checksums.txt`
- Rename-first upgrade strategy (move existing binary to `.old` before replacing)
- Extract binary from zip archive
- Add install directory to user `PATH` via Windows registry

### Bash Installer (`install.sh`)

Features:
- Auto-detect OS (`linux`/`darwin`) and architecture (`amd64`/`arm64`)
- Download versioned `.tar.gz` from GitHub Releases
- Verify SHA-256 checksum (`sha256sum` or `shasum -a 256` fallback)
- Rename-first upgrade strategy
- Shell-aware PATH configuration (bash → `.bashrc`, zsh → `.zshrc`, fish → `config.fish`)
- Support `curl` or `wget` for downloads
- CLI flags: `--version`, `--dir`, `--arch`, `--no-path`

### Common Install Script Pattern

Both scripts follow the same flow:

```
1. Detect platform and architecture
2. Download binary archive + checksums.txt
3. Verify checksum (fail on mismatch)
4. Rename existing binary to .old (if upgrading)
5. Extract and install new binary
6. Clean up .old file
7. Add to PATH (if not already present)
8. Print installed version
```

---

## Changelog Extraction

See [Release Body and Changelog](../07-release-body-and-changelog.md) for the full extraction pattern and release body template.

---

## Release Body

The GitHub Release description is assembled from multiple sources:

1. **Changelog entry** — extracted from `changelog.md`
2. **Release info table** — version, commit SHA, branch, build date, Go version
3. **Checksums block** — full `checksums.txt` in a code fence
4. **Install instructions** — quick install one-liners for PowerShell and Bash, plus version-pinned alternatives
5. **Asset matrix** — table mapping platform/architecture to filenames

### Release Type Detection

| Condition | Release Type | `prerelease` | `make_latest` |
|-----------|-------------|--------------|---------------|
| Version contains `-` (e.g., `v1.0.0-beta.1`) | Pre-release | `true` | `false` |
| Version is clean (e.g., `v1.0.0`) | Stable | `false` | `true` |

---

## GitHub Release Creation

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

## Build-Once Rule

Binaries are compiled **exactly once**. Compression, checksums, and publishing operate on the already-built artifacts and must **never** trigger a rebuild.

---

## Error Recovery

| Failure Point | Behavior | Fallback |
|---------------|----------|----------|
| Version resolution fails | Pipeline exits — invalid Git ref | Fix ref and re-push |
| One cross-compile target fails | Pipeline fails — **no partial releases** | Fix build error, re-tag |
| Code signing times out (600s) | Signing step fails; pipeline stops | Set `SIGNPATH_SIGNING_ENABLED=false`, release unsigned |
| Changelog extraction finds nothing | Uses fallback: `"No changelog entry found for $VERSION."` | Add changelog entry |
| Docs-site build fails | Pipeline fails | Fix docs or remove docs-site step |
| GitHub Release creation fails | Pipeline fails | Check `contents: write` permission |

**Policy**: Release builds are **all-or-nothing**. A partial release (e.g., 5 of 6 targets) is never published. If any stage fails, the entire pipeline fails and must be re-run after fixing the root cause.

### Changelog Fallback

```bash
CHANGELOG_ENTRY=$(awk -v ver="$VERSION" '
  /^## / { if (found) exit; if (index($0, ver)) found=1 }
  found { print }
' changelog.md 2>/dev/null || echo "")

if [ -z "$CHANGELOG_ENTRY" ]; then
  CHANGELOG_ENTRY="No changelog entry found for $VERSION."
fi
```

---

## Constraints

- **Build once, package once** — no rebuild during packaging
- Every release commit must run to completion — never cancel
- Version is resolved from the Git ref, never hardcoded
- All binaries are statically linked (`CGO_ENABLED=0`)
- Checksums are generated for all release assets (including `docs-site.zip`)
- Install scripts use placeholder substitution, not string interpolation
- Raw binaries are removed after compression — only archives ship
- Multi-module builds output to a **single `dist/` directory**
- Docs-site is built and packaged as `docs-site.zip` if present

---

## Cross-References

- [Shared Conventions](../01-shared-conventions.md) — Version resolution, checksums, permissions
- [GitHub Release Standard](../02-github-release-standard.md) — Release body assembly, pre-release detection
- [Complete Workflow Reference](./03-complete-workflow-reference.md) — Full annotated YAML assembling all patterns
- [Install Script Generation](../04-install-script-generation.md) — Placeholder strategy, PS1/Bash installers
- [Code Signing](../05-code-signing.md) — SignPath integration, pipeline placement
- [Release Body and Changelog](../07-release-body-and-changelog.md) — Changelog extraction, body template
- [Binary Icon Branding](../09-binary-icon-branding.md) — `go-winres make` and version injection
- [Self-Update Mechanism](../06-self-update-mechanism.md) — How CLI tools consume release assets
- [Self-Update Full Specs](../../14-update/00-overview.md) — Client-side update implementation
- [Updater Binary](../../14-update/19-updater-binary.md) — Standalone updater architecture

---

*Go binary release pipeline — updated: 2026-04-13*
