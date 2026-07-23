# 02 — Release Pipeline

## Purpose

Define the end-to-end CI/CD workflow that transforms a version tag into
a published GitHub Release with downloadable assets, checksums, and
install scripts.

---

## Trigger

The release pipeline triggers on:

```yaml
on:
  push:
    tags:
      - 'v*'
```

Optionally also on pushes to `release/**` branches for pre-release
builds.

---

## Pipeline Stages

The pipeline executes these stages **in strict order**:

```
1. Checkout  →  2. Setup  →  3. Resolve Version  →  4. Build
     →  5. Compress  →  6. Checksum  →  7. Generate Scripts
     →  8. Extract Changelog  →  9. Publish
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
| 8 | **Extract Changelog** | Pull the relevant `CHANGELOG.md` section | No |
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

See [01-cross-compilation.md](01-cross-compilation.md) for the full
build process. The output is raw binaries in `dist/`.

**Critical rule**: binaries are built **exactly once**. All subsequent
stages operate on these artifacts. No stage may trigger a rebuild.

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
    rm "$f"   # Remove raw binary after archiving
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

The `checksums.txt` file is published as a release asset alongside the
archives.

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

1. **Changelog section** — extracted from `CHANGELOG.md`
2. **Installation one-liners** — for PowerShell and Bash
3. **Checksum table** — formatted SHA-256 hashes

```markdown
## Installation

**Windows (PowerShell)**
```powershell
irm https://<repo-raw-url>/install.ps1 | iex
```

**Linux / macOS**
```bash
curl -fsSL https://<repo-raw-url>/install.sh | bash
```

## Checksums

| File | SHA-256 |
|------|---------|
| <binary>-linux-amd64.tar.gz | abc123... |
| ... | ... |
```

---

## Prerelease Detection

If the version contains a hyphen (e.g., `v1.2.0-beta.1`), the release
is marked as a **prerelease** and is NOT marked as latest:

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
  cancel-in-progress: false   # Never cancel release jobs
```

Release branches must NOT use `cancel-in-progress: true` — a cancelled
release leaves orphaned tags and partial assets.

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

## Application-Specific References

| App Spec | Covers |
|----------|--------|
| [02-powershell-build-deploy.md](../03-general/02-powershell-build-deploy.md) | Cross-compilation and checksums in build pipeline |
| [05-code-signing.md](../03-general/05-code-signing.md) | SignPath.io integration, signing step placement, SmartScreen reputation |
| [12-release-command.md](../01-app/12-release-command.md) | gitmap-v28 release CLI, version resolution, CI workflow |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
