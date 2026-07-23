# Release Pipeline

## Overview

The release pipeline automates binary production, packaging, and GitHub Release creation whenever code is pushed to a `release/**` branch or a `v*` tag. It produces cross-compiled binaries for two tools, platform-specific install scripts, checksums, and a fully formatted release page.

---

## Trigger and Concurrency

### Trigger

```yaml
on:
  push:
    branches:
      - "release/**"
    tags:
      - "v*"
```

### Concurrency

Release pipelines **never cancel** in-progress runs — every release commit must produce complete artifacts:

```yaml
concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false
```

### Permissions

The workflow needs write access to create GitHub Releases:

```yaml
permissions:
  contents: write
```

---

## Version Resolution

The version is extracted from the Git ref:
- **Tags** (`refs/tags/v1.2.3`): use the tag name directly
- **Branches** (`refs/heads/release/v1.2.3`): strip the `release/` prefix

```bash
if [[ "$GITHUB_REF" == refs/tags/* ]]; then
  VERSION="${GITHUB_REF_NAME}"
elif [[ "$GITHUB_REF" == refs/heads/release/* ]]; then
  VERSION="${GITHUB_REF_NAME#release/}"
fi
```

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
VERSION="<resolved-version>"
LDFLAGS="-s -w -X '<module>/constants.Version=$VERSION'"
CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -ldflags "$LDFLAGS" -o "dist/<binary>-${VERSION}-${os}-${arch}${ext}" .
```

### Multiple Binaries

If the project produces multiple binaries (e.g., a main tool and an updater), build each from its own module directory. All outputs go into a single `dist/` folder for unified packaging.

---

## Docs-Site Bundling

If the project includes a documentation site (e.g., React/Vite), the release pipeline builds and bundles it as a release asset so the CLI can serve it locally.

### Build Step

```yaml
- name: Setup Node.js
  uses: actions/setup-node@v4
  with:
    node-version: 22

- name: Build docs-site
  run: |
    npm ci
    npm run build
    cd dist
    zip -r ../dist-output/docs-site.zip .
```

### How It Fits

1. Runs **after** binary compilation but **before** compression/checksums
2. Produces `docs-site.zip` in the `dist/` folder alongside binaries
3. The checksum step covers `docs-site.zip` automatically (it checksums all files in `dist/`)
4. Ships as a standard release asset via `softprops/action-gh-release`

### Install Script Integration

Both install scripts download `docs-site.zip` from the release and extract it next to the binary:

```
<install-dir>/
  <binary>
  docs-site/
    dist/
      index.html
      assets/
        ...
```

The CLI's `help-dashboard` command auto-extracts `docs-site.zip` on first run if the `docs-site/` directory is missing, so the docs work out of the box even if the installer skipped extraction.

### Graceful Fallback

- If the docs-site build fails, the release should still succeed — `docs-site.zip` is optional
- Install scripts treat download failure as non-fatal (print a skip message, continue)
- The CLI falls back to `npm run dev` if no static `dist/` is found

---

## Compression and Checksums

### Compression Rules

| Platform | Format | Command |
|----------|--------|---------|
| Windows | `.zip` | `zip "${name}.zip" "$file"` |
| Linux, macOS | `.tar.gz` | `tar czf "${name}.tar.gz" "$file"` |

After compression, remove the raw binaries — only archives ship as release assets.

### Checksums

Generate SHA-256 checksums for all files in the dist directory:

```bash
sha256sum * > checksums.txt
```

The `checksums.txt` file is included as a release asset for verification.

---

## Install Scripts

The pipeline generates version-pinned install scripts using placeholder substitution. Scripts are written inline in the workflow with `VERSION_PLACEHOLDER` and `REPO_PLACEHOLDER` tokens, then `sed` replaces them:

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

### Common Install Script Patterns

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

Extract the matching version section from `CHANGELOG.md` using `awk`:

```bash
awk -v ver="$VERSION" '
  /^## / {
    if (found) exit
    if (index($0, ver)) found=1
  }
  found { print }
' CHANGELOG.md
```

Falls back to a "No changelog entry found" message if the version section doesn't exist.

---

## Release Body

The GitHub Release description is assembled from multiple sources:

1. **Changelog entry** — extracted from `CHANGELOG.md`
2. **Release info table** — version, commit SHA (truncated), branch, build date, Go version
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
    name: ${{ steps.version.outputs.version }}
    body_path: /tmp/release-body.md
    files: dist/*
    draft: false
    prerelease: ${{ contains(steps.version.outputs.version, '-') }}
    make_latest: ${{ !contains(steps.version.outputs.version, '-') }}
```

---

## Constraints

- **Build once, package once** — binaries are compiled exactly once; compression, checksums, and publishing operate on the already-built artifacts and must never trigger a rebuild
- Every release commit must run to completion — never cancel release pipelines
- Version is resolved from the Git ref, never hardcoded
- All binaries are statically linked (`CGO_ENABLED=0`)
- Checksums are generated for all release assets
- Install scripts use placeholder substitution, not string interpolation
- Install scripts implement rename-first upgrades (never delete-then-write)
- Pre-release vs. stable is determined automatically from the version string
