# Release Pipeline

**Version:** 1.0.0
**Updated:** 2026-04-19

---

## Purpose

Define the GitHub Release artifacts produced by `.github/workflows/release.yml`. The pipeline runs on every `v*` tag push and MUST produce every artifact in the table below — no exceptions.

---

## Trigger

```yaml
on:
  push:
    tags:
      - "v*"

concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false
```

Release builds MUST NEVER be cancelled — every release tag must produce a GitHub Release regardless of subsequent commits.

---

## Required artifacts

Each row below is mandatory. CI MUST `test -f` (or equivalent) before the publish step and fail the build if any artifact is missing.

| # | Artifact | Source | Build step |
|---|----------|--------|-----------|
| 1 | `coding-guidelines-linters-vX.Y.Z.zip` | `linters-cicd/` (excluding `__pycache__/` and `*.pyc`) | `zip -rq` in pipeline |
| 2 | `coding-guidelines-slides-vX.Y.Z.zip` | `slides-app/dist/` after `bun run build && bun run package` | inside `slides-app/` |
| 3 | `linters-install.sh` | `linters-cicd/install.sh` (renamed) | `cp` in pipeline |
| 4 | `install.sh` | repo root | committed |
| 5 | `install.ps1` | repo root | committed |
| 6 | `install-config.json` | repo root | committed |
| 7 | `checksums.txt` | computed | `sha256sum` of every zip |
| 8 | Other `release-artifacts/*.zip` and `*.tar.gz` | `release.sh` outputs | `bash release.sh` |

---

## Job order

```
checkout
  → resolve version (from GITHUB_REF_NAME)
    → bash release.sh                         (builds release-artifacts/)
      → verify release-artifacts/ + checksums.txt
        → setup bun
          → cd slides-app && bun install + bun run build + bun run package
            → verify slides-app/dist.zip
              → rename slides-app zip with version
                → zip linters-cicd/ → coding-guidelines-linters-vX.Y.Z.zip
                  → sha256sum >> checksums.txt
                    → cp install.sh as linters-install.sh
                      → softprops/action-gh-release@v2 publishes everything
```

---

## Release notes

The release body MUST include, at minimum:

1. Install one-liners for Bash and PowerShell.
2. The slide-deck download instructions (filename + double-click `index.html`).
3. The CI/CD linter pack quick-start (composite action one-liner + curl one-liner).
4. A pointer to `checksums.txt` for SHA-256 verification.

The wording in `.github/workflows/release.yml` `body:` is the canonical source.

---

## Build-once rule

Each artifact is compiled **exactly once** per release. Compression, checksums, and publishing operate on already-built artifacts and MUST NOT trigger a rebuild.

---

## Pre-release detection

Tags containing a `-` (e.g. `v3.5.0-beta.1`) MUST be marked `prerelease: true` and MUST NOT be marked `make_latest`.

```yaml
prerelease: ${{ contains(steps.version.outputs.version, '-') }}
make_latest: ${{ !contains(steps.version.outputs.version, '-') }}
```

---

## Failure modes

| Symptom | Cause | Fix |
|---------|-------|-----|
| `release-artifacts/ missing` | `release.sh` failed silently | Read `release.sh` logs; fix and retag |
| `slides-app/dist.zip missing` | Vite build failed | Check `bun run build` output |
| `linters-cicd zip missing` | wrong working directory or `__pycache__` exclude pattern | Verify exclude globs |
| Checksum mismatch downstream | `checksums.txt` not regenerated after re-zip | Always append to `checksums.txt` after zipping |

---

## Anti-requirements

- MUST NOT publish a release that is missing any artifact in the §"Required artifacts" table.
- MUST NOT cancel an in-progress release for a newer tag (concurrency `cancel-in-progress: false`).
- MUST NOT publish source maps inside `coding-guidelines-slides-*.zip` — `slides-app/scripts/package-zip.mjs` must strip them.
- MUST NOT include `node_modules/`, `.git/`, or `__pycache__/` in any release zip.

---

## Cross-references

- [`./00-overview.md`](./00-overview.md) — Distribution overview
- [`./01-install-contract.md`](./01-install-contract.md) — Install contract (consumes these artifacts)
- [`spec/12-cicd-pipeline-workflows/`](../12-cicd-pipeline-workflows/) — CICD conventions
- [`spec/16-generic-release/`](../16-generic-release/) — Generic release standard

---

*Release pipeline — v1.0.0 — 2026-04-19*
