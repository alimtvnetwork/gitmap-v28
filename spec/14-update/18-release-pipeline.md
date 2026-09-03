# Release Pipeline

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Purpose

Define the end-to-end CI/CD workflow that transforms a version tag into a published GitHub Release with downloadable assets, checksums, and install scripts.

---

## Trigger

The release pipeline triggers on:

```yaml
on:
  push:
    tags:
      - 'v*'
```

Optionally also on pushes to `release/**` branches for pre-release builds.

---

## Pipeline Stages

The pipeline executes these stages **in strict order**:

```
1. Checkout → 2. Setup → 3. Resolve Version → 4. Build
  → 5. Compress → 6. Checksum → 7. Generate Scripts
  → 8. Extract Changelog → 9. Publish
```

### Stage Details

| # | Stage | Description | Fails Pipeline |
|---|-------|-------------|----------------|
| 1 | **Checkout** | Clone the repository at the tagged commit | Yes |
| 2 | **Setup** | Install Go toolchain (from `go.mod`) | Yes |
| 3 | **Resolve Version** | Extract version from tag (`refs/tags/v1.2.0` → `1.2.0`) | Yes |
| 4 | **Build** | Cross-compile all 6+ targets into `dist/` | Yes |
| 5 | **Compress** | Archive each binary into `.zip` or `.tar.gz` | Yes |
| 6 | **Checksum** | Generate `checksums.txt` with SHA-256 hashes | Yes |
| 7 | **Generate Scripts** | Create version-pinned `install.ps1` and `install.sh` | Yes |
| 8 | **Extract Changelog** | Pull the relevant `changelog.md` section | No |
| 9 | **Publish** | Create GitHub Release with all assets | Yes |

---

## Version Resolution

Extract the version from the Git ref:

```bash

# From tag: refs/tags/v1.2.0 → 1.2.0

VERSION="${GITHUB_REF#refs/tags/v}"

# From branch: refs/heads/release/1.2.0 → 1.2.0

VERSION="${GITHUB_REF#refs/heads/release/}"
```

The resolved version is used for:

- Build-time `-ldflags` embedding
- Archive file naming
- Install script placeholder substitution
- Release title

---

## Build Stage

See [10-cross-compilation.md](./17-cross-compilation.md) for the full build process. The output is raw binaries in `dist/`.

**Critical rule**: binaries are built **exactly once**. All subsequent stages operate on these artifacts. No stage may trigger a rebuild.

---

## Compress Stage

After building, compress each binary:

```bash
cd dist

for f in <binary>-*; do
  if [[ "$f" == *.exe ]]; then
    # Windows: zip
    zip "${f%.exe}.zip" "$f"
  else
    # Unix: tar.gz (preserves permissions)
    tar czf "${f}.tar.gz" "$f"
  fi
  rm "$f"  # Remove raw binary after archiving
done
```

---

## Checksum Stage

Generate SHA-256 hashes for all archives:

```bash
cd dist
sha256sum *.zip *.tar.gz > checksums.txt
```

Or on macOS:

```bash
shasum -a 256 *.zip *.tar.gz > checksums.txt
```

The `checksums.txt` file is published as a release asset alongside the archives.

---

## Publish Stage

Create a GitHub Release with all assets from `dist/`:

```yaml
- uses: softprops/action-gh-release@v2
  with:
    files: dist/*
    body: ${{ steps.release_body.outputs.body }}
    prerelease: ${{ contains(env.VERSION, '-') }}
    make_latest: ${{ !contains(env.VERSION, '-') }}
```

### Release Body

The release body should include:

1. **Changelog section** — extracted from `changelog.md`
2. **Installation one-liners** — for PowerShell and Bash
3. **Checksum table** — formatted SHA-256 hashes

Example release body structure:

```

## Installation

**Windows (PowerShell)**
irm https://<repo>/install.ps1 | iex

**Linux / macOS**
curl -fsSL https://<repo>/install.sh | bash

## Checksums

| File | SHA-256 |
|------|---------|
| <binary>-linux-amd64.tar.gz | abc123... |
| ... | ... |
```

---

## Prerelease Detection

If the version contains a hyphen (e.g., `v1.2.0-beta.1`), the release is marked as a **prerelease** and is NOT marked as latest:

```yaml
prerelease: ${{ contains(env.VERSION, '-') }}
make_latest: ${{ !contains(env.VERSION, '-') }}
```

---

## Concurrency

Prevent parallel release jobs from conflicting:

```yaml
concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false  # Never cancel release jobs
```

Release branches must NOT use `cancel-in-progress: true` — a cancelled release leaves orphaned tags and partial assets.

---

## Permissions

The workflow needs write access to create releases:

```yaml
permissions:
  contents: write
```

---

## Constraints

- The build stage produces artifacts exactly once — no rebuilds.
- All actions pinned to exact version tags (e.g., `@v6`, not `@main`).
- No interactive prompts in any CI step.
- No notification steps (email, Slack) in the pipeline.
- Release jobs must not be cancellable mid-execution.

---

## Cross-References

- [Cross-Compilation](./17-cross-compilation.md) — Build targets and static linking
- [Release Assets](./14-release-assets.md) — Asset naming and packaging conventions
- [Checksums & Verification](./15-checksums-verification.md) — SHA-256 generation and verification
- [Install Scripts](./19-install-scripts.md) — One-liner installer generation
- [Release Versioning](./16-release-versioning.md) — Version resolution and tagging
- [Updater Binary](./20-updater-binary.md) — Multi-module build for the updater
- [Network Requirements](./21-network-requirements.md) — Download infrastructure for update consumers
- [CI/CD Pipeline Workflows](../12-cicd-pipeline-workflows/01-index.md) — Server-side pipeline patterns
- [Complete Workflow Reference](../12-cicd-pipeline-workflows/02-go-binary-deploy/04-complete-workflow-reference.md) — Full annotated YAML

---

*Release pipeline — v3.2.0 — 2026-04-13*
