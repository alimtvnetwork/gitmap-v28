# Acceptance Criteria — Fix-Repo & Installers

> **Version:** 1.0.0 · **Updated:** 2026-04-28
> Each item is a binary check. The "Test" column points at an executable
> proof in `tests/installer/` (Bash) or `linter-scripts/tests/` (Python).

---

## AC-FR — fix-repo

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

## AC-INST — install.sh / install.ps1

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

## AC-REL — release-install.sh / .ps1

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

## AC-VC — visibility-change

| # | Criterion | Test |
|---|---|---|
| AC-VC-001 | `--visible pub|pri` and toggle (no flag) all parse cleanly | `tests/installer/check-visibility-arg-parsing.sh` |
| AC-VC-002 | Provider detection: GitHub vs GitLab vs unsupported (exit `4`) | `tests/installer/check-visibility-provider-detect.sh` |
| AC-VC-003 | Runner integration `./run.sh visibility` forwards flags 1:1 | `tests/installer/check-visibility-runner-wiring.sh` |
| AC-VC-004 | `private → public` without `--yes` and non-TTY stdin → exit `7` | `tests/installer/check-visibility-arg-parsing.sh` |
| AC-VC-005 | `--dry-run` performs no API calls | `tests/installer/check-visibility-arg-parsing.sh` |

---

## AC-CR — CODE RED compliance (all four scripts)

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

## How to run the matrix

```bash

# All installer + fix-repo + visibility tests

bash tests/installer/run-tests.sh

# CODE RED check on the 8 scripts

bash linters-cicd/run-all.sh --path . --format text
```

A green run of both commands satisfies every AC above except AC-CR-002
(see note).
