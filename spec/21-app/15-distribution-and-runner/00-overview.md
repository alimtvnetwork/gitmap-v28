# Distribution and Runner

**Version:** 1.0.0
**Updated:** 2026-04-19
**Owner:** Md. Alim Ul Karim
**Status:** Authoritative — implementation must match this spec.

---

## Purpose

This folder defines the **end-user distribution surface** of `coding-guidelines-v24`:

1. The **install scripts** (`install.sh`, `install.ps1`) that pull spec + linters from a GitHub release into a user's repo.
2. The **runner scripts** (`run.sh`, `run.ps1`) at the repo root that update the local clone, build artifacts, and dispatch sub-commands (e.g. `slides`).
3. The **CI/CD release pipeline** (`.github/workflows/release.yml`) that packages every distributable artifact: `linters-cicd/`, `slides-app/dist/`, install scripts, and SHA-256 checksums.
4. The **`install-config.json`** contract that drives the default folder list.

Everything in this folder is **end-user-facing**. If a non-developer can't follow the README and have a working install in 60 seconds, this spec has failed.

---

## Reading order

| # | File | Purpose |
|---|------|---------|
| 00 | [00-overview.md](./00-overview.md) | This document |
| 01 | [01-install-contract.md](./01-install-contract.md) | What `install.sh` / `install.ps1` install, in what layout, and from where |
| 02 | [02-runner-contract.md](./02-runner-contract.md) | Root `run.sh` / `run.ps1` sub-command surface (`<no args>`, `slides`, `lint`) |
| 03 | [03-release-pipeline.md](./03-release-pipeline.md) | Which artifacts the GitHub Release publishes and their naming |
| 04 | [04-install-config.md](./04-install-config.md) | Schema and defaults for `install-config.json` |

---

## Distributable artifacts (canonical list)

Every GitHub Release MUST publish all of the following. Missing any one is a release blocker.

| Artifact | Source | Filename pattern | Purpose |
|----------|--------|------------------|---------|
| Spec + linters tree (download-on-demand) | `spec/`, `linters/`, `linter-scripts/`, `linters-cicd/` in main branch | (sourced via `codeload.github.com` archive — not a release asset) | Powers `install.sh`/`install.ps1` |
| Linters CI/CD pack | `linters-cicd/` | `coding-guidelines-linters-vX.Y.Z.zip` | Drop-in CI artifact; consumed by `linters-install.sh` |
| Slides deck | `slides-app/dist/` | `coding-guidelines-slides-vX.Y.Z.zip` | Offline trainer deck (double-click `index.html`) |
| Bash installer | `install.sh` | `install.sh` | Linux/macOS one-liner |
| PowerShell installer | `install.ps1` | `install.ps1` | Windows one-liner |
| Linters quick-installer | `linters-cicd/install.sh` (renamed) | `linters-install.sh` | CI one-liner that installs only `linters-cicd/` |
| Default install config | `install-config.json` | `install-config.json` | Authoritative folder list shipped with installers |
| Checksums | computed in CI | `checksums.txt` | SHA-256 of every zip |

---

## Sub-command surface (root runner)

The repo-root `run.sh` / `run.ps1` MUST implement this contract. Sub-commands are positional; flags are forwarded to the inner script.

| Invocation | Effect |
|------------|--------|
| `./run.ps1` (no args) | `git pull` → run the Go coding-guidelines validator on `src/` (legacy default; preserved for back-compat) |
| `./run.ps1 lint [path]` | Same as the no-args form, but explicit. Forwards `--Path` etc. to `linter-scripts/run.ps1`. |
| `./run.ps1 slides` | `git pull` → `cd slides-app && bun install && bun run build && bun run preview` → open the preview URL in the default browser |
| `./run.ps1 help` | Print the sub-command table |

> **Default behavior is preserved.** Existing users who type `./run.ps1` with no args MUST get the same Go validator behavior they had before this spec.

---

## Default install layout

After running `install.sh` / `install.ps1` with no flags, the user's repo MUST contain:

```
<dest>/
├── spec/                  ← full coding-guidelines spec tree
├── linters/               ← per-language lint plugins (eslint configs, etc.)
├── linter-scripts/        ← orchestrator scripts (legacy validator)
└── linters-cicd/          ← Python check suite, run-all.sh, registry, baseline
```

The default folder list lives in `install-config.json` and MUST equal:

```json
["spec", "linters", "linter-scripts", "linters-cicd"]
```

This list is the contract. It MUST stay in sync with §"Distributable artifacts" above.

---

## Cross-references

- Slides app spec: [`spec-slides/00-overview.md`](../../spec-slides/00-overview.md)
- CICD pipeline conventions: [`spec/12-cicd-pipeline-workflows/`](../12-cicd-pipeline-workflows/)
- Generic CLI conventions: [`spec/13-generic-cli/`](../13-generic-cli/)
- Generic release standard: [`spec/16-generic-release/`](../16-generic-release/)

---

*Distribution-and-runner overview — v1.0.0 — 2026-04-19*

---

## Verification

_Auto-generated section — see `spec/15-distribution-and-runner/97-acceptance-criteria.md` for the full criteria index._

### AC-DIST-000: Distribution & runner conformance: Overview

**Given** Validate the install contract and runner contract against a clean machine fixture.  
**When** Run the verification command shown below.  
**Then** Install script is idempotent; runner detects missing deps and exits with a stable error code; PATH entries are deduped.

**Verification command:**

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
