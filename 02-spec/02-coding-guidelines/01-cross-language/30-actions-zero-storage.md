# Rule R18: GitHub Actions Zero-Storage Mandate & Artifact Hygiene (AI Execution Prompt)

> **/goal** Eliminate all persistent GitHub Actions artifact storage consumption in CI workflows, maintaining account-wide usage strictly at 0.0 GB while delegating official releases to GitHub Releases.
> **/learn** Master the ephemeral compilation verification pattern, enforce a total ban on `actions/upload-artifact` in CI/PR pipelines, restrict emergency diagnostic uploads to 1-day retention on failure, and maintain automated nightly purge workflows.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Maintain 0.0 GB persistent artifact storage in GitHub Actions across all routine workflows and pull request pipelines.
- [ ] `/learn` Never use `actions/upload-artifact` for compiled binaries, test logs, coverage files, or inter-job data passing.
- [ ] `/goal` Route official downloadable binary distributions exclusively through GitHub Releases (`release.yml`), which do not count against the 0.5 GB Actions quota.
- [ ] `/learn` Gate any diagnostic artifact upload with `if: failure()` and enforce a strict 1-day retention limit (`retention-days: 1`).

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Status:** Active
**Scope:** Universal (All Repositories, Workflows & Languages)
**Strictness:** Non-Negotiable (Banned Anti-Pattern)

---

## 1. Overview & Core Motivation

GitHub Free accounts provide only **0.5 GB (500 MB)** of total GitHub Actions artifact storage per month across all repositories in the account. When developers upload build binaries (20–30 MB each) across multi-architecture matrix legs (e.g., Linux, macOS, Windows on amd64/arm64) or archive test logs with multi-day retention:
- A single CI run consumes 150–200 MB of storage.
- After just 3–4 runs, the account exceeds 90% (0.45 GB) and triggers account-wide billing warnings or blocks all automated pipeline executions.
- Multi-day retention (`retention-days: 14` or default `90`) keeps obsolete binaries alive for weeks, permanently exhausting the free tier.

**The Golden Rule:** GitHub Actions CI pipelines MUST NOT consume any persistent artifact storage. Actions storage usage must remain strictly at **0.0 GB**.

---

## 2. Strict Prohibitions (What is Strictly Banned)

1. **Total Ban on `actions/upload-artifact` in CI Pipelines:**
   - ❌ **FORBIDDEN:** Using `actions/upload-artifact` in standard CI workflows (`ci.yml`), preflight checks, unit test runs, linter checks, or pull request validations.
   - ❌ **FORBIDDEN:** Uploading compiled application binaries (`*.exe`, Linux binaries, macOS Mach-O files) as workflow artifacts.
   - ❌ **FORBIDDEN:** Uploading test logs, stdout dumps, HTML coverage reports, SARIF files, or benchmark summaries as persistent artifacts.

2. **Total Ban on Long Artifact Retention:**
   - ❌ **FORBIDDEN:** Default 90-day retention or multi-week retention (`retention-days: 14`, `retention-days: 30`).
   - Any temporary emergency diagnostic upload must enforce a strict **1-day floor** (`retention-days: 1`).

3. **Total Ban on Inter-Job Artifact Passing for Matrix Builds:**
   - ❌ **FORBIDDEN:** Compiling a binary in job A, uploading it as an artifact, and downloading it in job B solely to check its file size or print a summary table.

---

## 3. Mandatory Compliant Patterns (How Actions Must Be Written)

### Pattern A: Ephemeral Compilation Verification

CI compilation jobs exist purely to verify that code compiles cleanly with exit status 0 across supported targets (`CGO_ENABLED=0 GOOS=... GOARCH=... go build`).
- The matrix runner builds the binary locally in the ephemeral runner workspace.
- The runner inspects and logs the output size: `echo "Built $OUTPUT ($(du -h "$OUTPUT" | cut -f1))"`.
- The binary is NOT uploaded. When the job finishes, the runner is destroyed and zero bytes of Actions storage are consumed.

### Pattern B: Release Assets Belong Exclusively in GitHub Releases

- Official downloadable binaries and release archives belong exclusively in **GitHub Releases** (via `gh release upload`, `softprops/action-gh-release`, or Goreleaser in `release.yml`).
- Files attached to GitHub Releases are stored in GitHub's release asset storage, which **DOES NOT count** toward the monthly 0.5 GB GitHub Actions artifact quota!

### Pattern C: Emergency Failure Diagnostics (Strict 1-Day Floor)

If an artifact upload is ever strictly needed for debugging flaky CI failures, it MUST:
1. Only run on failure: `if: failure()`
2. Target a single compressed text file (e.g. `failure-log.txt`)
3. Enforce 1-day retention: `retention-days: 1`

### Pattern D: Automated Storage Purge Workflow

Every active repository must maintain a scheduled purge workflow (`.github/workflows/purge-actions-artifacts.yml`) that executes nightly to discover and delete any orphaned artifacts via GitHub API, guaranteeing 0.0 GB storage footprint.

---

## 4. Verification Checklist

- [ ] Zero instances of `actions/upload-artifact` exist in `ci.yml` and routine workflows.
- [ ] CI build matrices verify compilation without uploading binary payloads.
- [ ] Release binaries are uploaded to GitHub Releases, not Actions artifacts.
- [ ] If any diagnostic artifact upload exists, it is gated with `if: failure()` and `retention-days: 1`.
- [ ] Scheduled artifact purge workflow is installed and active.

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-ARCH-030: GitHub Actions Zero-Storage Mandate and Ephemeral Verification

**Given** GitHub Actions workflow YAML files across `.github/workflows/`.
**When** Workflow linters and CI storage audits evaluate action steps, retention configurations, and artifact uploads.
**Then** Standard CI workflows contain zero calls to `actions/upload-artifact`, build matrices execute ephemeral compilation without persistent storage, emergency diagnostic uploads enforce `retention-days: 1`, and account storage remains at 0.0 GB, achieving exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.
