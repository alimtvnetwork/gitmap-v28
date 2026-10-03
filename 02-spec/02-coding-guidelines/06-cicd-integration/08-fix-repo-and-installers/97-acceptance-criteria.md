# Acceptance Criteria — Fix-Repo & Installers (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable acceptance criteria registry and verification test matrix for fix-repo, installer, and visibility-change scripts.
> **/learn** Master the AC-CG-INSTALL criteria taxonomy, binary verification proofs, test suite bindings in tests/installer/, and CODE RED limits.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify all script contracts map 1:1 to executable test cases in `tests/installer/` and `linter-scripts/`.
- [ ] `/learn` Maintain zero behavioral divergence between Bash (`.sh`) and PowerShell (`.ps1`) implementations.
- [ ] `/goal` Enforce CODE RED limits (≤ 300 script lines, ≤ 8–15 function lines, zero nested conditionals, positive booleans).
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 1.0.0
**Updated:** 2026-04-28
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None — every criterion is backed by an executable test in `tests/installer/`.

---

## 1. Fix-Repo & Installers Criteria Inventory (`AC-CG-INSTALL-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-INSTALL-001` | Fix-Repo & Installer Scripts Index Conformance | [`readme.md`](readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers --check-only` |
| `AC-CG-INSTALL-002` | Fix-Repo Version Token Rewriter Specification Conformance | [`02-fix-repo-contract.md`](02-fix-repo-contract.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers --check-only` |
| `AC-CG-INSTALL-003` | Installer Scripts Specification Conformance | [`03-installer-contract.md`](03-installer-contract.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers --check-only` |
| `AC-CG-INSTALL-004` | Visibility Change Script Specification Conformance | [`04-visibility-change-contract.md`](04-visibility-change-contract.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers --check-only` |

---

## 2. Binary Test Suite Mappings

### AC-FR — fix-repo

| # | Criterion | Test |
|---|---|---|
| AC-FR-001 | `fix-repo.sh` and `fix-repo.ps1` both exist at repo root and are executable | `tests/installer/check-fix-repo-runner-wiring.sh` |
| AC-FR-002 | Default mode (no flag) replaces last 2 prior versions | `tests/installer/check-fix-repo-arg-forwarding.sh` |
| AC-FR-003 | `--3`, `--5`, `--all` accepted; `--4`, `--6` rejected with exit `6` | `tests/installer/check-fix-repo-arg-forwarding.sh` |
| AC-FR-004 | Numeric-overflow guard: `coding-guidelines-v170` is NOT touched | `tests/installer/check-fix-repo-url-rewrite.sh` |
| AC-FR-005 | URL occurrences ARE rewritten with host preserved | `tests/installer/check-fix-repo-url-rewrite.sh` |
| AC-FR-006 | Idempotent: second run after first changes 0 files | `tests/installer/check-fix-repo-contract-conformance.sh` |
| AC-FR-007 | Exit codes `2/3/4/5/6/7` raised on the documented failures | `tests/installer/check-fix-repo-exit-code-propagation.sh` |
| AC-FR-008 | `--dry-run` writes nothing | `tests/installer/check-fix-repo-debug-preflight.sh` |
| AC-FR-009 | Honors `.gitignore` (uses `git ls-files -z`) | `tests/installer/check-fix-repo-contract-conformance.sh` |

---

### AC-INST — install.sh / install.ps1

| # | Criterion | Test |
|---|---|---|
| AC-INST-001 | Banner printed BEFORE any network call | `tests/installer/check-log-header-env.sh` |
| AC-INST-002 | `-n` / `--no-latest` skips the latest probe | `tests/installer/check-no-latest-api.sh` |
| AC-INST-003 | `--version vX.Y.Z` engages PINNED MODE end-to-end | `tests/installer/check-release-install-acceptance.sh` |
| AC-INST-004 | Unknown flag → exit `1` and prints offending flag to stderr | `tests/installer/check-install-ps1-help.sh` |
| AC-INST-005 | `--help` / `-h` exits `0` after printing usage | `tests/installer/check-install-ps1-help.sh` |
| AC-INST-006 | `--folders` accepts subpaths (e.g. `02-spec/14-update`) | `tests/installer/check-install-folders-config.sh` |
| AC-INST-007 | `--run-fix-repo` invokes `fix-repo.{sh,ps1}` after verify | `tests/installer/check-run-fix-repo-flag.sh` |
| AC-INST-008 | `--log-dir`, `--show-fix-repo-log`, `--max-fix-repo-logs` honored | `tests/installer/check-log-dir-flag.sh`, `check-show-fix-repo-log-flag.sh`, `check-max-fix-repo-logs-flag.sh` |
| AC-INST-009 | Bundle installers use this script unchanged | `tests/installer/check-bundle-installers.sh` |

---

### AC-REL — release-install.sh / .ps1

| # | Criterion | Test |
|---|---|---|
| AC-REL-001 | Always pinned; never queries `/releases/latest` | `tests/installer/check-no-latest-api.sh` |
| AC-REL-002 | Resolution precedence: `--version` > `INSTALLER_VERSION` > baked | `tests/installer/check-release-install-acceptance.sh` |
| AC-REL-003 | Disagreeing sources emit a WARN; higher precedence wins | `tests/installer/check-release-install-acceptance.sh` |
| AC-REL-004 | Invalid semver → exit `2` | `tests/installer/check-release-install-acceptance.sh` |
| AC-REL-005 | Pinned asset 404 at both endpoints → exit `3` | `tests/installer/check-release-install-acceptance.sh` |
| AC-REL-006 | Hand-off to `install.sh` uses `--no-latest --version <tag>` | `tests/installer/check-release-bake.sh` |
| AC-REL-007 | Inner installer rejecting handshake → exit `5` | `tests/installer/check-release-install-acceptance.sh` |

---

### AC-VC — visibility-change

| # | Criterion | Test |
|---|---|---|
| AC-VC-001 | `--visible pub|pri` and toggle (no flag) all parse cleanly | `tests/installer/check-visibility-arg-parsing.sh` |
| AC-VC-002 | Provider detection: GitHub vs GitLab vs unsupported (exit `4`) | `tests/installer/check-visibility-provider-detect.sh` |
| AC-VC-003 | Runner integration `./run.sh visibility` forwards flags 1:1 | `tests/installer/check-visibility-runner-wiring.sh` |
| AC-VC-004 | `private → public` without `--yes` and non-TTY stdin → exit `7` | `tests/installer/check-visibility-arg-parsing.sh` |
| AC-VC-005 | `--dry-run` performs no API calls | `tests/installer/check-visibility-arg-parsing.sh` |

---

### AC-CR — CODE RED compliance (all four scripts)

| # | Criterion | Verification |
|---|---|---|
| AC-CR-001 | Every function ≤ 8–15 effective lines | `python3 linter-scripts/check-function-lengths.py` |
| AC-CR-002 | Every script file ≤ 300 lines | `wc -l fix-repo.sh fix-repo.ps1 install.sh install.ps1 release-install.sh release-install.ps1 visibility-change.sh visibility-change.ps1` |
| AC-CR-003 | Zero nested `if`; guard-and-return | `bash linters-cicd/run-all.sh --path .` (CODE-RED-004) |
| AC-CR-004 | Booleans positively named | `bash linters-cicd/run-all.sh --path .` (CODE-RED-023) |
| AC-CR-005 | Errors never swallowed | code review + linter |

> Note: `install.sh` (702 lines) and `install.ps1` (581 lines) currently
> EXCEED AC-CR-002 and require splitting into `scripts/install/*`
> helpers. Tracked as a follow-up; see `99-troubleshooting.md`.

---

## 3. Detailed Acceptance Criteria Specifications

### AC-CG-INSTALL-001: Fix-Repo and Installer Scripts Index Conformance

- [ ] Root `readme.md` provides complete module navigation, execution prompts, and cross-references.
- [ ] All specification files in `08-fix-repo-and-installers/` adhere to active prompt anatomy, checklist headers, and acceptance criteria blocks.
- [ ] Dual Bash (`.sh`) and PowerShell (`.ps1`) script siblings maintain 100% behavioral parity.

**Given** Installer and auto-fix repository management specifications.
**When** Audited against this installation specification.
**Then** Zero contract or visibility violations are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-INSTALL-002: Fix-Repo Version Token Rewriter Specification Conformance

- [ ] `fix-repo.sh` and `fix-repo.ps1` both exist at repo root, parse origin remote URL, and validate `-v(\d+)$` version suffix.
- [ ] Default mode replaces last 2 prior versions; `--3`, `--5`, `--all` are accepted, and closed set rejects `--4`, `--6` with exit `6`.
- [ ] Numeric-overflow guard ensures non-target versions (`coding-guidelines-v170`) are untouched, while URLs are rewritten with host preserved.
- [ ] Execution is idempotent, honors `.gitignore` via `git ls-files -z`, and supports `--dry-run` without modifying files.
- [ ] Deterministic exit codes `2` (`E_NOT_A_REPO`), `3` (`E_NO_REMOTE`), `4` (`E_NO_VERSION_SUFFIX`), `5` (`E_BAD_VERSION`), `6` (`E_BAD_FLAG`), `7` (`E_WRITE_FAILED`) are enforced.

**Given** Installer scripts and repository fix infrastructure.
**When** Audited against this contract specification.
**Then** Zero contract drift is detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-INSTALL-003: Installer Scripts Specification Conformance

- [ ] Mandatory configuration banner is printed before any network call or remote probe across `install.*` and `release-install.*`.
- [ ] `-n` / `--no-latest` skips latest probe, while `--version <tag>` engages pinned mode end-to-end.
- [ ] Pinned mode never queries `/releases/latest`, never falls back to `main`, and never crosses repository boundaries.
- [ ] Resolution precedence enforces `--version` > `INSTALLER_VERSION` > baked placeholder, emitting `WARN` on disagreement.
- [ ] SHA-256 verification against `checksums.txt` exits `4` on checksum mismatch.
- [ ] Downstream fix-repo execution via `--run-fix-repo` captures logs, rotates history, and propagates exit codes 1:1.
- [ ] Unified exit codes `0` (Success), `1` (Generic failure), `2` (Offline/Invalid version), `3` (Asset 404), `4` (Verify failed), `5` (Handoff rejected) are strictly raised.

**Given** Installer scripts and repository fix infrastructure.
**When** Audited against this contract specification.
**Then** Zero contract drift is detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-INSTALL-004: Visibility Change Script Specification Conformance

- [ ] Provider auto-detection distinguishes GitHub vs GitLab from origin remote URL, exiting `4` on unsupported hosts.
- [ ] `--visible pub|pri` and default toggle parse cleanly; invalid values exit `6`.
- [ ] Interactive confirmation prompt protects `private → public` transitions unless `--yes` / `-y` is provided.
- [ ] Non-interactive stdin without `--yes` aborts with exit `7`.
- [ ] `--dry-run` performs no API calls and prefixes output with `[dry-run]`.
- [ ] Verification re-reads visibility post-apply and raises exit code `8` if unchanged.
- [ ] CODE RED standards enforced across all scripts: functions ≤ 8–15 effective lines, zero nested conditionals, positive booleans.

**Given** Installer scripts and repository fix infrastructure.
**When** Audited against this contract specification.
**Then** Zero contract drift is detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers --check-only
```
**Expected:** exit 0. Zero violations.

---

## 4. How to Run the Verification Matrix

```bash
# All installer + fix-repo + visibility tests
bash tests/installer/run-tests.sh

# CODE RED check on the 8 scripts
bash linters-cicd/run-all.sh --path . --format text
```

A green run of both commands satisfies every AC above except AC-CR-002 (see note).
