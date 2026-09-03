# Consolidated: CI/CD Pipeline Workflows — Complete Reference

**Version:** 3.3.0
**Updated:** 2026-04-22
**Source Module:** [`02-spec/12-cicd-pipeline-workflows/`](../12-cicd-pipeline-workflows/01-index.md)

---

## Purpose

This is the **standalone consolidated reference** for all CI/CD pipeline specifications. An AI reading only this file must be able to implement any pipeline described here without consulting source specs.

### Source-Folder Coverage Map

The source module `02-spec/12-cicd-pipeline-workflows/` contains 13 root specs + 3 subfolders + diagrams. All are summarized below.

| Source File / Folder | Section | Status |
|----------------------|---------|--------|
| `03-shared-conventions.md`, `02-ci-pipeline.md` | §§Platform, Triggers, CI Pipeline | ✅ Full |
| `04-github-release-standard.md`, `05-release-pipeline.md` | §Release Pipeline | ✅ Full |
| `06-vulnerability-scanning.md` | §Vulnerability Scanning | ✅ Full |
| `03-reusable-ci-guards/` | §Reusable CI Guards | ✅ Full (added) |
| `07-install-script-generation.md`, `08-installation-flow.md` | §Install Scripts & Flow | ✅ Full |
| `09-changelog-integration.md`, `10-code-signing.md` | §Changelog & Code Signing | ✅ Full |
| `11-self-update-mechanism.md`, `12-version-and-help.md` | §Self-Update & Version | ✅ Reference (full in `20-self-update-app-update.md`) |
| `13-environment-variable-setup.md`, `14-release-body-and-changelog.md` | §Env Setup & Release Body | ✅ Full |
| `15-terminal-output-standards.md` | §Terminal Output | ✅ Full |
| `16-binary-icon-branding.md` | §Binary Icon Branding | ✅ Full |
| `17-release-pipeline-issues-rca.md` | §Release Issues RCA | ✅ Reference |
| `01-browser-extension-deploy/` | §Browser Extension Deploy | ✅ Reference |
| `02-go-binary-deploy/` | §Go Binary Deploy | ✅ Reference |
| `images/` (mermaid diagrams) | Diagram links throughout | ✅ |

---

## Platform & Shared Conventions

| Convention | Rule |
|-----------|------|
| Platform | GitHub Actions, `ubuntu-latest` |
| Action versions | Pinned to exact tags (`@v6`), **never** `@latest` or `@main` |
| Tool versions | Pinned to exact versions, **never** `@latest` |
| CI concurrency | Cancel superseded runs, except release branches |
| Release concurrency | **Never cancel** — every commit runs to completion |
| Version resolution | Derived from Git ref (tag or branch), never hardcoded |
| Checksums | SHA-256 for all release assets |
| Permissions | Minimum required (`read` for CI, `write` for releases) |
| Working directory | Use `working-directory:` key, never `cd` |
| Build-once rule | Compile once, then compress/checksum/publish — no rebuilds |

### Trigger Patterns

```yaml

# CI Pipeline

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

# Release Pipeline

on:
  push:
    branches: ["release/**"]
    tags: ["v*"]

# Scheduled Scans

on:
  schedule:
    - cron: "0 9 * * 1"  # Weekly Monday 9:00 UTC
  workflow_dispatch:       # Manual trigger
```

### Concurrency Control

```yaml

# CI — cancel old runs on same branch (except release branches)

concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: ${{ !startsWith(github.ref, 'refs/heads/release/') }}

# Release — NEVER cancel

concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false
```

### Permissions

| Pipeline | Permission | Reason |
|----------|-----------|--------|
| CI | `contents: read` | Only reads source code |
| Release | `contents: write` | Creates GitHub Releases, uploads assets |

---

## Pipeline Archetype 1: Browser Extension Deploy

For Node.js/pnpm projects that produce Chrome/Chromium extensions.

### CI Pipeline

```
Lint → Test → SDK Build → Parallel Module Builds → Extension Assembly
```

- Diamond dependency graph: SDK builds first, modules build in parallel, extension assembles last
- Each module is a separate build target with its own `pnpm build`
- Source maps included in CI builds

### Release Pipeline

```
Same build graph → Source Map Removal → ZIP Packaging → GitHub Release
```

- Source maps stripped from production builds
- `.zip` archives of extension and standalone components
- Asset naming: `{component}-v{version}.zip`

---

## Pipeline Archetype 2: Go Binary Deploy

For Go projects that produce cross-compiled CLI binaries.

### CI Pipeline — 4-Phase Structure

```
Phase 1: SHA Dedup Gate → Phase 2: Lint + Vulncheck → Phase 3: Test Matrix → Phase 4: Cross-Compile + Summary
```

#### Phase 1: SHA-Based Dedup Gate

```yaml
- name: Check for changes
  id: dedup
  run: |
    CURRENT_SHA=$(git rev-parse HEAD)
    LAST_SHA=$(cat .last-ci-sha 2>/dev/null || echo "")
    if [ "$CURRENT_SHA" = "$LAST_SHA" ]; then
      echo "skip=true" >> $GITHUB_OUTPUT
    else
      echo "skip=false" >> $GITHUB_OUTPUT
    fi
```

Skips redundant validation via **step-level conditionals** (not job-level `if`) to ensure the job always succeeds for branch protection rules.

#### Phase 2: Lint + Vulnerability Check

- `golangci-lint` with pinned version
- `govulncheck` for Go vulnerability scanning
- Stdlib vulns → warn only; third-party vulns → fail

#### Phase 3: Test Matrix

- Unit tests across multiple Go versions
- Race detector enabled: `go test -race ./...`
- Coverage reporting

#### Phase 4: Cross-Compile

- 6 targets: `windows/amd64`, `windows/arm64`, `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`
- `CGO_ENABLED=0` for fully static binaries
- `-ldflags -X` embeds `Version`, `RepoPath`, `CommitSHA`, `BuildDate`
- Output to `dist/` directory

### Release Pipeline — 9-Stage Strict Order

```
1. Checkout → 2. Setup Go → 3. Resolve Version → 4. Build All Targets →
5. Compress → 6. Checksum → 7. Generate Install Scripts → 8. Extract Changelog → 9. Publish
```

#### Version Resolution

```bash

# From tag: refs/tags/v1.2.0 → 1.2.0

VERSION="${GITHUB_REF#refs/tags/v}"
```

#### Compression

| OS | Format | Tool |
|----|--------|------|
| Windows | `.zip` | `zip` |
| Linux/macOS | `.tar.gz` | `tar czf` (preserves permissions) |

#### Checksum Generation

```bash
cd dist && sha256sum *.zip *.tar.gz > checksums.txt
```

Generated from **compressed archives**, not raw binaries (TOCTOU prevention).

#### Prerelease Detection

Versions containing `-` (e.g., `1.2.0-beta.1`) are marked prerelease and **NOT** set as latest.

#### Release Asset Set (Standard: 9 files)

| # | File | Description |
|---|------|-------------|
| 1-6 | `<binary>-<os>-<arch>.<ext>` | 6 compressed archives |
| 7 | `checksums.txt` | SHA-256 checksums |
| 8 | `install.ps1` | PowerShell installer |
| 9 | `install.sh` | Bash installer |

---

## Install Script Generation

### PowerShell (`install.ps1`) Flow

```
1. Print banner → 2. Detect arch → 3. Download .zip + checksums.txt →
4. Verify SHA-256 → 5. Rename existing (rename-first) →
6. Extract to install dir → 7. Clean .old file →
8. Register PATH (Registry + Profile + Git Bash) → 9. Print summary
```

- Install location: `$env:LOCALAPPDATA\<binary>\`
- PATH registration: Windows Registry (User), PowerShell `$PROFILE`, Git Bash profiles

### Bash (`install.sh`) Flow

```
1. Print banner → 2. Detect OS/arch → 3. Select download tool (curl/wget) →
4. Download .tar.gz + checksums.txt → 5. Verify SHA-256 →
6. Rename existing → 7. Extract → 8. Clean .old →
9. Register PATH in shell profiles → 10. Print summary
```

- Install location: `$HOME/.local/bin` (user) or `/usr/local/bin` (root)
- Shell-aware PATH: bash (`~/.bashrc`), zsh (`~/.zshrc`), fish (`~/.config/fish/config.fish`)
- SHA-256 fallback: `sha256sum` → `shasum -a 256` (macOS)

### Placeholder Strategy

Both scripts use `VERSION_PLACEHOLDER` and `REPO_PLACEHOLDER` tokens replaced via `sed` after heredoc write during release pipeline.

---

## Code Signing

Windows binary signing via SignPath:

| Aspect | Detail |
|--------|--------|
| Feature flag | `vars.SIGNPATH_SIGNING_ENABLED` repository variable |
| Pipeline placement | After build, before compression/checksums |
| In-place replacement | Signed binaries overwrite unsigned in `dist/` |
| Scope | Only `.exe` files; Linux/macOS not signed |
| Verification | Confirms signing completed without error |

---

## Vulnerability Scanning

| Mode | Trigger | Purpose |
|------|---------|---------|
| In-CI | Every push/PR | Gate code changes |
| Standalone | Weekly schedule + manual | Catch newly disclosed CVEs |

**Classification:**

- Third-party vulnerabilities → **fail** the pipeline
- Stdlib vulnerabilities → **warn** only (logged, not blocking)

---

## Release Body & Changelog

### Changelog Extraction

```bash
awk '/^## \['"${VERSION}"'\]/{flag=1;next}/^## \[/{flag=0}flag' changelog.md
```

Graceful fallback if version section not found.

### Release Body Template

```markdown

## What's Changed

{changelog content}

## Build Info

| Key | Value |
|-----|-------|
| Version | {version} |
| Commit | {sha} |
| Branch | {branch} |
| Build Date | {date} |
| Go Version | {go_version} |

## Checksums

\`\`\`
{checksums.txt content}
\`\`\`

## Install

**PowerShell:** `irm .../install.ps1 | iex`
**Bash:** `curl -fsSL .../install.sh | bash`

## Assets

| File | OS | Arch |
|------|-----|------|
| binary-windows-amd64.zip | Windows | amd64 |
| ... | ... | ... |
```

---

## Terminal Output Standards

### Formatting Rules

| Rule | Detail |
|------|--------|
| Indentation | 2-space from left margin, never tabs |
| Progress | `[1/4] Step description` |
| Success | `[OK] Check passed` |
| Failure | `[FAIL] Check failed` |
| Warning | `[WARN] Non-fatal issue` |
| Errors | `Error: <message> (<operation>)` to stderr |
| Tables | Fixed-width columns with header separators |

### Status Icons

| Icon | Meaning | Used In |
|------|---------|---------|
| `[OK]` | Success | Doctor checks, verification |
| `[FAIL]` | Failure | Doctor checks |
| `[WARN]` | Warning | Deprecation, stdlib vulns |
| `[+]` | Added | PATH registration |
| `[-]` | Removed | PATH removal |
| `[=]` | No change | Already exists |
| `[SKIP]` | Skipped | Already installed |

**ASCII-only** for scripts and CI output. Unicode (`✓`, `✗`) only in Go-compiled binaries.

### Error Output Format

```
  Error: Failed to open database at /path/to/db.sqlite
    Reason:     file is locked by another process
    Operation:  store.Open
    Suggestion: Close other instances and retry
```

---

## Binary Icon Branding

Windows icon embedding via `go-winres`:

- `winres.json` in project root
- Icon files in `assets/` folder
- Embedded during build via `go generate`
- Distinct icon for updater binary vs main binary

---

## Version & Help System

### Version Display

```
<binary> version
  Version:    1.2.0
  Commit:     abc1234
  Build Date: 2026-04-16
  Go Version: go1.22.0
```

### Help System

- Root `--help` shows command list
- Subcommand `--help` shows usage, flags, examples
- CI verification: `<binary> version` and `<binary> --help` must exit 0

---

## Environment Variable Setup

The `env` command manages persistent PATH and environment:

- `<binary> env` — auto-detect and register
- Windows: User PATH via Registry + PowerShell profile + Git Bash profiles
- Unix: Shell-aware PATH injection (bash/zsh/fish)
- Auto-home directory creation

---

## Key Rules

1. ALL CI/CD content belongs in `02-spec/12-cicd-pipeline-workflows/`
2. App-specific deployment notes go in `21-app/` instead
3. Shared patterns (version resolution, checksums, releases) are in root-level files
4. Archetype-specific patterns are in their respective subfolders
5. New pipeline types get their own subfolder with `01-index.md`
6. Client-side update logic lives in `02-spec/14-update/`
7. Never use `@latest` — pin all versions
8. Never cancel release pipeline jobs
9. Build once — never rebuild after the build stage

---

## File Inventory

| # | File | Description |
|---|------|-------------|
| 01 | `03-shared-conventions.md` | Platform, triggers, concurrency, version resolution |
| 02 | `04-github-release-standard.md` | Release body, pre-release detection, asset matrix |
| 03 | `06-vulnerability-scanning.md` | Standalone and in-CI vuln scanning |
| 04 | `07-install-script-generation.md` | PS1+Bash installer pattern, placeholder strategy |
| 05 | `10-code-signing.md` | SignPath integration, feature-flag gating |
| 06 | `11-self-update-mechanism.md` | CLI self-update: deploy path, rename-first, handoff |
| 07 | `14-release-body-and-changelog.md` | Changelog extraction, release body template |
| 08 | `08-installation-flow.md` | End-to-end install: one-liners, upgrade, uninstall |
| 09 | `09-changelog-integration.md` | Changelog format, CI extraction, release body assembly |
| 10 | `10-version-and-help.md` | Version display, help system, CI verification |
| 11 | `11-environment-variable-setup.md` | `env` command: persistent PATH, auto-home, registry |
| 12 | `12-terminal-output-standards.md` | Output formatting: icons, tables, progress, errors |
| 13 | `13-binary-icon-branding.md` | Windows icon embedding via go-winres |
| — | `01-browser-extension-deploy/` | Browser extension CI + release (3 files) |
| — | `02-go-binary-deploy/` | Go binary CI + release (3 files) |

---

## Reusable CI Guards

**Source:** `12-cicd-pipeline-workflows/03-reusable-ci-guards/`

The reusable-guards subfolder defines six pluggable workflow modules every pipeline composes from:

| Guard | Purpose | Failure Action |
|-------|---------|----------------|
| Forbidden Names | Block files/identifiers matching project blocklist | Fail build with file paths |
| Grandfather Baseline | Allow legacy violations from a baseline file but reject new ones | Fail with diff vs baseline |
| Collision Audit | Detect duplicate IDs across enums, error codes, route slugs | Fail with collision pairs |
| Lint Diff Gate | Lint only changed files in PRs (full lint on main) | Fail if delta introduces violations |
| Lint Suggestions | Auto-comment fixable suggestions on PRs | No fail; advisory |
| Test Aggregator | Merge multi-language test reports into one JUnit XML | Fail if any underlying suite failed |

Each guard ships as a reusable workflow callable via `uses: ./.github/workflows/guard-<name>.yml`. Inputs are documented per-guard in the source folder.

---

## Browser Extension Deploy (Subfolder)

**Source:** `12-cicd-pipeline-workflows/01-browser-extension-deploy/` (3 files)

- CI workflow lints, builds, and packages a `.zip` for Chrome / Edge / Firefox stores.
- Release workflow uploads the zip via the Chrome Web Store API + Edge Add-ons + AMO using API tokens stored as GitHub secrets (`CHROME_REFRESH_TOKEN`, `EDGE_API_KEY`, `AMO_JWT_ISSUER`).
- Version is derived from `manifest.json`'s `version` field; CI fails if it doesn't match the git tag.

---

## Go Binary Deploy (Subfolder)

**Source:** `12-cicd-pipeline-workflows/02-go-binary-deploy/` (3 files)

- 6-target cross-compilation matrix: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64`, `windows/arm64`.
- Each binary is built with `-trimpath -ldflags "-s -w -X main.Version=$VERSION -X main.Commit=$COMMIT -X main.BuildTime=$BUILD_TIME"`.
- Windows binaries get the icon embedded via `go-winres` before final compile (see Binary Icon Branding section).
- Outputs are SHA-256 checksummed and uploaded as release assets with consistent naming: `<binary>-<version>-<os>-<arch>[.exe]`.

---

## Release Pipeline Issues RCA

**Source:** `12-cicd-pipeline-workflows/17-release-pipeline-issues-rca.md`

Living document tracking every release-pipeline incident with: symptom, root cause, fix, regression test added, spec update. New incidents append a dated entry. The first entry chronologically lists the **build-once rule** as the lesson learned from a 2026-Q1 incident where assets diverged because of inadvertent rebuilds between checksum and upload.

---

## Vulnerability Scanning (Detailed)

**Source:** `12-cicd-pipeline-workflows/06-vulnerability-scanning.md`

| Mode | When | Tool | Action on Finding |
|------|------|------|-------------------|
| In-CI scan | Every push/PR | Language-specific (`govulncheck`, `npm audit`, `cargo audit`, `composer audit`) | Fail PR on high/critical |
| Scheduled scan | Weekly Monday 09:00 UTC | All scanners + SBOM diff | Open issue tagged `vulnerability` |
| Manual scan | `workflow_dispatch` | All scanners | Comment results to triggering user |

`govulncheck` version is **pinned to `v1.1.4`** (per the CI pipeline diagram). Upgrading requires a one-line PR + retrospective if the new version raises new findings.

---

## Install Script Generation (Detailed)

**Source:** `12-cicd-pipeline-workflows/07-install-script-generation.md`

- Two scripts generated per release: `install.ps1` (Windows) and `install.sh` (Linux/macOS).
- Both use a **placeholder strategy**: a template file with `{{VERSION}}`, `{{CHECKSUM_<TARGET>}}`, `{{ASSET_<TARGET>}}` tokens replaced at release time.
- Scripts are themselves checksummed; the README quick-install one-liner pins the script SHA-256 so it can never be silently swapped.
- Generated install scripts are uploaded as additional release assets (`install.ps1`, `install.sh`, `install.ps1.sha256`, `install.sh.sha256`).

---

## Code Signing (Detailed)

**Source:** `12-cicd-pipeline-workflows/10-code-signing.md`

- Windows binaries: SignPath integration via OIDC, no long-lived secrets in GitHub.
- macOS binaries: notarization via `notarytool` using App Store Connect API key (stored in GitHub secrets).
- Linux binaries: detached GPG signature alongside the binary; public key published to release notes.
- All signing is **feature-flag gated** by `ENABLE_CODE_SIGNING` repo variable so forks/dev builds skip it cleanly.

---

## Terminal Output Standards

**Source:** `12-cicd-pipeline-workflows/15-terminal-output-standards.md`

| Element | Rule |
|---------|------|
| Status icons | `✓` success, `✗` failure, `⚠` warning, `ℹ` info, `→` action |
| Tables | UTF-8 box drawing, fixed-width columns, never raw piped output |
| Progress | Spinner for indeterminate, `[NN/TOTAL]` prefix for determinate |
| Errors | Red bold prefix `Error:` followed by code + path on next line |
| Colors | Honor `NO_COLOR` env var; fall back to plain text in non-TTY |
| Width | Wrap at terminal width; never assume 80 cols |

---

## Binary Icon Branding (Detailed)

**Source:** `12-cicd-pipeline-workflows/16-binary-icon-branding.md`

- Windows binaries get a `.ico` embedded via `go-winres make --in winres.json` run **before** final `go build`.
- `winres.json` defines: icon path (multi-resolution `.ico`), version info block (matching the build version), file description, copyright, original filename.
- macOS binaries: `.icns` bundled into `.app` only when distributing as an app bundle (CLI tools skip this).
- Linux: no embedding; icon ships as a separate `.desktop` file in install scripts when applicable.

---

## Environment Variable Setup (`env` Command)

**Source:** `12-cicd-pipeline-workflows/13-environment-variable-setup.md`

The CLI exposes an `env` subcommand that ensures the binary directory is on `PATH` persistently:

| OS | Mechanism | Persistence |
|----|-----------|-------------|
| Windows | Edits `HKCU\Environment\Path` via Win32 registry, broadcasts `WM_SETTINGCHANGE` | Per-user, survives reboot |
| macOS | Appends `export PATH=...` to `~/.zshrc` (or `~/.bash_profile` if Bash) | Per-user shell profile |
| Linux | Appends to `~/.bashrc` and `~/.profile`; detects `fish` and updates `~/.config/fish/config.fish` | Per-user |

The `env` command is **idempotent** — running it twice never duplicates the entry. Detection is by exact-path match of the binary's parent directory.

---

## Release Body Format

**Source:** `12-cicd-pipeline-workflows/14-release-body-and-changelog.md`

Standard release body sections (assembled by CI from changelog.md):

1. **Summary** — one-paragraph what + why.
2. **Highlights** — 3–5 bullet headlines.
3. **Changes** — full changelog entries since previous tag.
4. **Install** — one-liner per OS with pinned script SHA.
5. **Verify** — checksum verification command per OS.
6. **Assets table** — every asset with size, SHA-256, signing status.
7. **Upgrade Notes** — only present if breaking changes.

Pre-releases (tags matching `v*-rc.*`, `v*-beta.*`) are flagged via the GitHub Releases API `prerelease: true` field automatically.

---

## Diagrams

**Source:** `12-cicd-pipeline-workflows/images/`

- `ci-pipeline-flow.mmd` — CI/release/vulnerability scan lifecycle.
- `unified-architecture.mmd` — Hierarchy of all CI/CD specs across Foundation / Distribution / Validation / Delivery layers.

Both are Mermaid sources rendered inline in the docs viewer. Do not copy as PNGs — render from source so updates stay in sync.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Full Source | `../12-cicd-pipeline-workflows/01-index.md` |
| Self-Update & App Update | `../14-update/01-index.md` |
| Folder Structure Rules | `../01-spec-authoring-guide/02-folder-structure.md` |

---

*Consolidated CI/CD pipeline workflows — v3.2.0 — 2026-04-16*
