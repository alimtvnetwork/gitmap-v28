# Distribution (AI Execution Prompt)

> **/goal** Package, version, verify, and distribute the portable linter pack (`linters-cicd/`) via release archives and GitHub composite Actions.
> **/learn** Master versioned release asset generation, checksum verification (`checksums.txt`), GitHub composite Action wiring, and zero-storage release compliance.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Package the standalone `linters-cicd/` directory into versioned release archives (`coding-guidelines-linters-vX.Y.Z.zip`).
- [ ] `/learn` Ensure `install.sh` verifies SHA-256 checksums against `checksums.txt` before extraction.
- [ ] `/goal` Support single-line workflow integration via GitHub composite Action `linters-cicd/action.yml`.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version:** 1.0.0
> **Updated:** 2026-04-19

The linter pack ships **two ways**, both produced by the same release job
in `.github/workflows/release.yml`.

---

## 1. Versioned ZIP (universal)

Built into every GitHub Release as
`coding-guidelines-linters-vX.Y.Z.zip`. Contains:

```
linters-cicd/
├── checks/
├── ci/
├── configs/
├── action.yml
├── run-all.sh
├── install.sh
├── readme.md
└── VERSION
```

Install one-liner (Linux / macOS):

```bash
curl -fsSL https://github.com/alimtvnetwork/coding-guidelines-v24/releases/latest/download/install.sh | bash
```

The installer:

1. Downloads the matching `coding-guidelines-linters-<latest>.zip`
2. Extracts to `./linters-cicd/`
3. Verifies SHA-256 against `checksums.txt`
4. Prints next-step commands

Flags:

- `-d <dir>` install destination (default: `./linters-cicd`)
- `-v <version>` install a specific version (default: latest)
- `-n` skip checksum verification (not recommended)

---

## 2. GitHub composite Action

`linters-cicd/action.yml` lets GitHub users skip the install entirely:

```yaml
- uses: alimtvnetwork/coding-guidelines-v24/linters-cicd@v3.9.0
```

GitHub clones the repo at the specified ref and runs `action.yml`. Zero
maintenance for consumers — just bump the version pin.

---

## Release pipeline integration

`.github/workflows/release.yml` runs on every `v*` tag and:

1. Zips `linters-cicd/` → `coding-guidelines-linters-vX.Y.Z.zip`
2. Computes SHA-256, appends to `checksums.txt`
3. Uploads as a Release asset alongside slides-app and core artifacts

---

*Part of [CI/CD Integration](./readme.md)*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/06-cicd-integration/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-CI-006: Distribution Packaging and Release Asset Governance

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.
