# CI Release Automation

The CI release pipeline (`.github/workflows/release.yml`) automates the production of 6 cross-compiled binary targets (Windows, Linux, Darwin for amd64/arm64) for both `gitmap` and `gitmap-updater` when pushing to `release/**` branches or `v*` tags. Artifacts use a versioned naming convention (e.g., `gitmap-v28.56.0-windows-amd64.zip`).

## Pipeline Steps

1. **Resolve version** — from tag name or branch name (`release/v2.56.0` → `v2.56.0`).
2. **Build binaries** — cross-compiles 12 binaries (6 gitmap + 6 gitmap-updater) with `-ldflags` version embedding.
3. **Compress** — Windows binaries → `.zip`; Linux/macOS → `.tar.gz`.
4. **Generate checksums** — SHA256 for all dist files → `checksums.txt`.
5. **Generate install scripts** — version-pinned `install.ps1` (Windows) and `install.sh` (Linux/macOS) are created with placeholder substitution and attached as release assets.
6. **Extract changelog** — matching section from `CHANGELOG.md`.
7. **Build release body** — changelog, metadata table, checksums, install instructions (PowerShell + Bash one-liners), and asset matrix.
8. **Create GitHub Release** — `softprops/action-gh-release@v2`; pre-release when version contains `-`; `make_latest` for stable releases.

## Install Scripts

Both install scripts are generated inline in the workflow with `VERSION_PLACEHOLDER` and `REPO_PLACEHOLDER` tokens replaced via `sed`:

| Script | Platform | Features |
|--------|----------|----------|
| `install.ps1` | Windows | Checksum verification, versioned binary detection, rename-first upgrade, registry PATH update |
| `install.sh` | Linux/macOS | Checksum verification (`sha256sum`/`shasum`), `.tar.gz` extraction, rename-first upgrade, shell-aware PATH append (bash/zsh/fish) |

## Concurrency

`concurrency: release-${{ github.ref }}, cancel-in-progress: true` — cancels superseded runs on the same ref.
