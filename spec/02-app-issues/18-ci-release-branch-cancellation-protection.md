# Post-Mortem #18: CI Release Branch Cancellation Protection

## Problem

The CI pipeline used `cancel-in-progress: true` unconditionally for
all branches, including release branches (`release/**`). When multiple
commits were pushed to a release branch in quick succession, earlier
runs were cancelled — potentially producing incomplete artifacts,
missing binaries, or partial metadata writes.

### Impact

1. **Incomplete release artifacts** — a cancelled release pipeline
   could leave behind partial binary uploads or missing checksums.
2. **Missed metadata writes** — the release metadata JSON and
   `latest.json` update could be skipped if the pipeline was
   cancelled mid-execution.
3. **Silent failures** — cancelled runs appeared as grey in the
   GitHub UI, giving no indication whether the release succeeded
   or was interrupted.
4. **Manual re-runs required** — contributors had to manually
   re-trigger the release pipeline after cancellation.

### Root Cause

```yaml
# BEFORE — unconditional cancellation (broken for releases)
concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: true
```

This applied the same cancellation policy to `main`, feature branches,
and release branches. While cancelling superseded runs on `main` is
desirable (only the latest commit matters), release branches require
every commit to produce complete, validated artifacts.

---

## Solution: Conditional cancel-in-progress

### CI Workflow

Replace the unconditional `true` with a conditional expression that
evaluates to `false` for release branches:

```yaml
# AFTER — conditional cancellation (correct)
concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: ${{ !startsWith(github.ref, 'refs/heads/release/') }}
```

This evaluates to:
- `true` (cancel) for `main`, feature branches, and PR refs.
- `false` (never cancel) for `release/**` branches.

### Release Workflow

Changed to `cancel-in-progress: false` unconditionally, since the
release workflow only triggers on `release/**` branches and `v*` tags:

```yaml
concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false
```

---

## Files Changed

| File | Change |
|------|--------|
| `.github/workflows/ci.yml` | Changed `cancel-in-progress` from `true` to conditional expression excluding `release/**` branches |
| `.github/workflows/release.yml` | Changed `cancel-in-progress` from `true` to `false` |
| `spec/03-general/08-ci-pipeline.md` | Added "Release Branch Protection" section with behavior tables |
| `.lovable/memory/tech/ci-pipeline-architecture.md` | Updated concurrency control section |

---

## Prevention Rules

1. **Never use unconditional `cancel-in-progress: true`** when the
   workflow handles release branches — always use a conditional
   expression to protect them.
2. **Release pipelines must always run to completion** — partial
   builds are unacceptable for release artifacts.
3. **Test cancellation behavior** — after modifying concurrency
   settings, push two rapid commits to a release branch and verify
   both runs complete.

---

## Behavior Matrix

| Scenario | Branch Type | cancel-in-progress | Behavior |
|----------|-------------|-------------------|----------|
| Rapid pushes to `main` | Non-release | `true` | Earlier run cancelled ✅ |
| Rapid pushes to `feature/*` | Non-release | `true` | Earlier run cancelled ✅ |
| Rapid pushes to `release/*` | Release | `false` | Both runs complete ✅ |
| PR update | Non-release | `true` | Earlier run cancelled ✅ |
| `v*` tag push | Release | `false` | Run completes ✅ |

---

## Acceptance Criteria

1. Pushing two commits to `release/**` in quick succession does
   **not** cancel the first run.
2. Pushing two commits to `main` in quick succession **does** cancel
   the first run.
3. The release workflow always runs to completion regardless of
   concurrent pushes.
4. All release artifacts (binaries, checksums, installers) are
   produced for every release commit.

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect.
