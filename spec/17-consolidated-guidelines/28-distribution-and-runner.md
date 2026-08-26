# 28 — Distribution and Runner (Consolidated)

> **Version:** 1.0.0
> **Updated:** 2026-04-22
> **Type:** Consolidated Guideline — Phase 6B promotion
> **Source module:** [`spec/15-distribution-and-runner/`](../15-distribution-and-runner/)
> **Status:** Authoritative digest — a blind AI can implement the install + runner surface from this file alone.

This file consolidates everything a distributing/consuming agent needs in
order to (a) ship a GitHub Release of `coding-guidelines-v24` and (b)
install + run it inside a target repository — without reading the source
folder.

It complements:

- `15-cicd-pipeline-workflows.md` — generic CI/CD conventions
- `17-self-update-app-update.md` — self-update / rename-first deploy mechanics
- `23-generic-cli.md` — sub-command dispatcher pattern reused by `run.sh` / `run.ps1`

---

## §1 — Purpose & Scope

The distribution surface has **four** moving parts. All four are
release-blockers if missing or out of sync:

1. **Install scripts** — `install.sh` (Bash) and `install.ps1`
   (PowerShell) pull spec + linters from a GitHub release into a user's
   repository.
2. **Runner scripts** — repo-root `run.sh` / `run.ps1` update the local
   clone, build artifacts, and dispatch sub-commands (`lint`, `slides`,
   `help`).
3. **CI/CD release pipeline** — `.github/workflows/release.yml` packages
   every distributable artifact (linters, slides deck, install scripts,
   SHA-256 checksums).
4. **`install-config.json`** — the authoritative folder list shipped
   with installers.

Acceptance test for the surface: a non-developer following the README
MUST have a working install in **≤ 60 seconds**. If not, this spec has
failed.

---

## §2 — Distributable Artifacts (Canonical List)

Every GitHub Release MUST publish all of the following. Missing any one
is a release blocker.

| Artifact | Source | Filename pattern | Purpose |
|----------|--------|------------------|---------|
| Spec + linters tree | `spec/`, `linters/`, `linter-scripts/`, `linters-cicd/` (main) | sourced via `codeload.github.com` archive (not a release asset) | Powers `install.sh` / `install.ps1` |
| Linters CI/CD pack | `linters-cicd/` | `coding-guidelines-linters-vX.Y.Z.zip` | Drop-in CI artifact; consumed by `linters-install.sh` |
| Slides deck | `slides-app/dist/` | `coding-guidelines-slides-vX.Y.Z.zip` | Offline trainer deck (double-click `index.html`) |
| Bash installer | `install.sh` | `install.sh` | Linux/macOS one-liner |
| PowerShell installer | `install.ps1` | `install.ps1` | Windows one-liner |
| Linters quick-installer | `linters-cicd/install.sh` (renamed) | `linters-install.sh` | CI one-liner, installs only `linters-cicd/` |
| Default install config | `install-config.json` | `install-config.json` | Authoritative folder list shipped with installers |
| Checksums | computed in CI | `checksums.txt` | SHA-256 of every zip — verified by installers |

**Drift contract:** the artifact list above MUST equal the `assets:`
list in `.github/workflows/release.yml`. CI fails the release if the
two diverge.

---

## §3 — Install Contract (`install.sh` / `install.ps1`)

### §3.1 One-liners (must work verbatim)

```bash

# Linux / macOS

curl -fsSL https://github.com/<org>/coding-guidelines-v24/releases/latest/download/install.sh | bash
```

```powershell

# Windows

irm https://github.com/<org>/coding-guidelines-v24/releases/latest/download/install.ps1 | iex
```

### §3.2 Behavior — both installers MUST

1. **Detect** an existing install (presence of any folder in
   `install-config.json`'s `folders` list) and prompt before overwriting.
2. **Download** the spec + linters tree via
   `codeload.github.com/<org>/coding-guidelines-v24/zip/refs/heads/main`.
3. **Verify** every downloaded zip against `checksums.txt` (SHA-256).
   Mismatch → exit code `2`, no files written.
4. **Extract** only the folders listed in `install-config.json` (default:
   `spec`, `linters`, `linter-scripts`, `linters-cicd`).
5. **Be idempotent.** Re-running with no flags MUST converge to the
   same on-disk state.
6. **Print a summary**: folders installed, version pinned, next-step
   commands.

### §3.3 Exit codes (stable contract)

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | User abort / declined overwrite |
| `2` | Checksum or download failure (no files written) |
| `3` | Missing prerequisite (`curl`/`unzip`/`tar`) |
| `4` | Invalid `install-config.json` |

---

## §4 — Runner Contract (`run.sh` / `run.ps1`)

The repo-root runner MUST implement this sub-command surface.
Sub-commands are **positional**; remaining flags are forwarded to the
inner script.

| Invocation | Effect |
|------------|--------|
| `./run.ps1` (no args) | `git pull` → run the Go validator on `src/` (legacy default; preserved for back-compat) |
| `./run.ps1 lint [path]` | Same as the no-args form, but explicit. Forwards `--Path` etc. to `linter-scripts/run.ps1` |
| `./run.ps1 slides` | `git pull` → `cd slides-app && bun install && bun run build && bun run preview` → open the preview URL in the default browser |
| `./run.ps1 help` | Print the sub-command table |

> **Back-compat is mandatory.** Existing users typing `./run.ps1` with
> no args MUST get the same Go validator behavior they had before this
> spec. Breaking this is a release blocker.

---

## §5 — Default Install Layout

After `install.sh` / `install.ps1` runs with no flags, the user's repo
MUST contain:

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

This list is **the contract**. It MUST stay in sync with §2's
"Distributable artifacts" table.

---

## §6 — `install-config.json` Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "version": "1.0.0",
  "folders": ["spec", "linters", "linter-scripts", "linters-cicd"],
  "source": {
    "owner": "<org>",
    "repo": "coding-guidelines-v24",
    "ref": "main"
  },
  "verify": {
    "checksums": "checksums.txt",
    "algorithm": "sha256"
  }
}
```

Validation rules:

- `version` MUST be SemVer.
- `folders` MUST be a non-empty array of repo-relative directory names.
- Every entry in `folders` MUST exist in the source repo's `main` ref
  at release time. The release pipeline fails otherwise.
- `source.ref` MAY be a branch, tag, or commit SHA.

---

## §7 — Release Pipeline (`.github/workflows/release.yml`)

Triggered on `push` of a tag matching `v[0-9]+.[0-9]+.[0-9]+`. Steps:

1. Checkout at the tag.
2. Build `slides-app/dist/` (`cd slides-app && bun install && bun run build`).
3. Zip `linters-cicd/` → `coding-guidelines-linters-vX.Y.Z.zip`.
4. Zip `slides-app/dist/` → `coding-guidelines-slides-vX.Y.Z.zip`.
5. Compute SHA-256 of every zip → `checksums.txt`.
6. Upload **all** artifacts from §2 to the GitHub Release.
7. Verify `assets:` list matches §2 (drift check). Fail otherwise.

---

## §8 — Verify & Uninstall (Consumer-Side)

Verify a downloaded installer before piping to a shell:

```bash
curl -fsSLO https://github.com/<org>/coding-guidelines-v24/releases/latest/download/install.sh
curl -fsSLO https://github.com/<org>/coding-guidelines-v24/releases/latest/download/checksums.txt
sha256sum -c checksums.txt --ignore-missing
bash install.sh
```

Uninstall:

```bash

# from the repo root that was installed-into

rm -rf spec linters linter-scripts linters-cicd
```

---

## §9 — Checklist for a Blind AI

Before claiming the distribution surface is correct, verify **all** of
the following:

- [ ] Both one-liners in §3.1 work against `releases/latest`.
- [ ] Installer exit codes match §3.3.
- [ ] Re-running the installer is idempotent (§3.2 #5).
- [ ] `install-config.json` validates against §6.
- [ ] `./run.sh` and `./run.ps1` implement every row of §4.
- [ ] No-arg `./run.ps1` still triggers the legacy Go validator.
- [ ] Release workflow uploads every artifact in §2.
- [ ] `checksums.txt` covers every zip uploaded.
- [ ] Slides zip extracts to a folder whose `index.html` opens offline.

---

## §10 — Cross-References

- Source folder: [`spec/15-distribution-and-runner/`](../15-distribution-and-runner/)
- Slides app spec: [`spec-slides/00-overview.md`](../../spec-slides/00-overview.md)
- CI/CD pipeline conventions: [`spec/12-cicd-pipeline-workflows/`](../12-cicd-pipeline-workflows/)
- Generic CLI conventions: [`spec/13-generic-cli/`](../13-generic-cli/)
- Generic release standard: [`spec/16-generic-release/`](../16-generic-release/)
- Self-update mechanics: [`17-self-update-app-update.md`](./17-self-update-app-update.md)

---

*Consolidated digest of `15-distribution-and-runner/` — v1.0.0 — 2026-04-22*