# Update — Overview

**Version:** 2.0.0  
**Updated:** 2026-04-17  
**Status:** Active  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Purpose

Single source of truth for the **CLI update mechanism** used across
all Go CLI tools in the Riseup Asia stack. Covers the full lifecycle:
release production → install → version probing → self-update →
rename-first deploy → handoff → verification → cleanup.

This folder consolidates all update-related specs (formerly split across
`14-generic-update/` and `15-self-update-app-update/`) into a single source
of truth. Do not re-split.

Any CLI that ships an `install` or `update` (or `self-update`)
subcommand MUST follow the contracts in this folder.

---

## Keywords

`self-update` · `cli-update` · `rename-first-deploy` · `handoff` ·
`cleanup` · `deploy-path` · `console-safe` · `winres` ·
`latest-json` · `install-script` · `version-probe` · `code-signing` ·
`updater-binary` · `release-pipeline`

---

## Scoring

| Metric | Value |
|--------|-------|
| AI Confidence | Production-Ready |
| Ambiguity | None |
| Health Score | 100/100 (A+) |

---

## Architecture (One-Liner)

```
release  →  install (probe → download → verify → deploy)
         →  run binary
         →  binary checks for update (latest.json)
         →  download → verify checksum → rename current bin
         →  write new bin to canonical deploy path
         →  exec new bin (handoff) → verify → cleanup → done
```

The defining property: **the running binary is renamed, never
deleted, before the new one is written.** This avoids Windows
file-lock errors and guarantees rollback is always one rename away.

---

## File Inventory

### Core mechanism (01–08)

| # | File | Description | Status |
|---|------|-------------|--------|
| 01 | [01-self-update-overview.md](./01-self-update-overview.md) | Self-update contract: triggers, version comparison, **skip-if-current fast-path**, latest.json schema | ✅ Active |
| 02 | [02-deploy-path-resolution.md](./02-deploy-path-resolution.md) | Canonical deploy path resolution per OS, env-var overrides, PATH registration | ✅ Active |
| 03 | [03-rename-first-deploy.md](./03-rename-first-deploy.md) | Rename-first algorithm, atomicity, Windows file-lock handling | ✅ Active |
| 04 | [04-build-scripts.md](./04-build-scripts.md) | `run.ps1` / `run.sh` build/release scripts (LDFLAGS, asset naming) | ✅ Active |
| 05 | [05-handoff-mechanism.md](./05-handoff-mechanism.md) | Old → new binary handoff: **two-phase summary**, `cmd.Run()` blocking | ✅ Active |
| 06 | [06-cleanup.md](./06-cleanup.md) | Post-handoff cleanup: **mandatory auto-cleanup**, retry on Windows lock | ✅ Active |
| 07 | [07-console-safe-handoff.md](./07-console-safe-handoff.md) | Windows console safety: detached vs attached, stdout/stderr inheritance | ✅ Active |
| 08 | [08-repo-path-sync.md](./08-repo-path-sync.md) | Cross-repo path sync: keeping deploy paths consistent across tools | ✅ Active |

### Verification & metadata (09–12)

| # | File | Description | Status |
|---|------|-------------|--------|
| 09 | [09-version-verification.md](./09-version-verification.md) | Three-branch active-vs-deployed verification, required TRACE/HINT diagnostics | ✅ Active |
| 10 | [10-last-release-detection.md](./10-last-release-detection.md) | Standalone `Get-LastRelease.ps1` / `get-last-release.sh` with tiered fallback | ✅ Active |
| 11 | [11-windows-icon-embedding.md](./11-windows-icon-embedding.md) | `go-winres` + `winres.json` + committed `.syso` strategy | ✅ Active |
| 12 | [12-code-signing.md](./12-code-signing.md) | Authenticode (Win), `codesign` + notarization (macOS), detached GPG signatures | ✅ Active |

### Release & install (13–23)

| # | File | Description | Status |
|---|------|-------------|--------|
| 13 | [13-release-assets.md](./13-release-assets.md) | Asset naming, layout, manifest fields per release | ✅ Active |
| 14 | [14-checksums-verification.md](./14-checksums-verification.md) | SHA-256 generation, `.sha256` files, post-download verification | ✅ Active |
| 15 | [15-release-versioning.md](./15-release-versioning.md) | SemVer rules, tagging, version-bump cadence | ✅ Active |
| 16 | [16-cross-compilation.md](./16-cross-compilation.md) | `GOOS`/`GOARCH` matrix, build host requirements | ✅ Active |
| 17 | [17-release-pipeline.md](./17-release-pipeline.md) | End-to-end release pipeline (build → sign → upload → publish) | ✅ Active |
| 18 | [18-install-scripts.md](./18-install-scripts.md) | `install.ps1` / `install.sh` contract, idempotency, PATH registration | ✅ Active |
| 19 | [19-updater-binary.md](./19-updater-binary.md) | Standalone updater binary architecture (binary-based update path) | ✅ Active |
| 20 | [20-network-requirements.md](./20-network-requirements.md) | HTTP client, retry, proxy support, timeout policy | ✅ Active |
| 21 | [21-config-file.md](./21-config-file.md) | `powershell.json` config schema, defaults, override rules | ✅ Active |
| 22 | [22-update-command-workflow.md](./22-update-command-workflow.md) | End-to-end `<binary> update` workflow (source-based + binary-based) | ✅ Active |
| 23 | [23-install-script-version-probe.md](./23-install-script-version-probe.md) | Latest-version probe in `install.{ps1,sh}` (current+1 .. +20, parallel) | ✅ Active |

### Auxiliary

| File | Description |
|------|-------------|
| [99-consistency-report.md](./99-consistency-report.md) | Per-folder consistency audit |
| [diagrams/](./diagrams/) | Mermaid flow diagrams (handoff, cleanup, deploy) |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| CI/CD self-update mechanism (pipeline side) | `../12-cicd-pipeline-workflows/06-self-update-mechanism.md` |
| Generic CLI blueprint | `../13-generic-cli/00-overview.md` |
| Release pipeline (asset production) | `../16-generic-release/02-release-pipeline.md` |
| Consolidated summary | `../17-consolidated-guidelines/17-self-update-app-update.md` |

---

## Placement Rules

```
AI INSTRUCTION:

1. ALL update-related content (install, update, self-update,
   release, deploy, cleanup, signing) belongs in this folder.
2. App-specific update behavior (UI prompts, app-side gating)
   goes in 21-app/ instead.
3. Pipeline-side concerns (CI workflow YAML, runner setup) go in
   12-cicd-pipeline-workflows/.
4. Each file follows the standard {NN}-{kebab-case-name}.md naming.
5. Add new files to the File Inventory above and update
   99-consistency-report.md.
6. Numbered ranges:
     01–08  Core mechanism
     09–12  Verification & metadata
     13–23  Release & install
     24+    Future additions (do NOT renumber existing files)
```

---

*Overview — v2.0.0 — 2026-04-17 (post-merge)*

---

## Verification

_Auto-generated section — see `spec/14-update/97-acceptance-criteria.md` for the full criteria index._

### AC-UPD-000: Self-update conformance: Overview

**Given** Exercise the rename-first deploy path against a fixture release directory.  
**When** Run the verification command shown below.  
**Then** `latest.json` is written atomically; the old binary is renamed (not deleted) before the new one is moved into place; rollback restores the previous version.

**Verification command:**

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
