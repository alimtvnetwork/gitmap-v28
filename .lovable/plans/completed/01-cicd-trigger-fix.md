# Execution Plan: CI/CD Trigger & Concurrency Fix

## Task Summary

The GitHub Actions pipelines are failing to trigger on commits, causing concurrency cancellation to fail and badges to permanently display "failing" states. We need to identify the root cause of this lack of triggering and fix it.

## Root Cause Analysis Findings

1. Actionlint identified fatal parsing errors in `.github/workflows/ci.yml`.
   - Missing `run:` or `uses:` directive in `Run installer smoke (source mode)`.
   - Cloudflare email obfuscation artifact: `uses: RubbaBoy/[email protected]` instead of `BYOB@v1.3.0`.
2. Fatal error in `.github/workflows/race-detector.yml`.
   - Missing required input `cache-suffix` for the `setup-go-cached` composite action.
3. Due to these invalid YAML structures, GitHub Actions silently refused to queue the workflows. This resulted in no builds, no concurrency cancellations, and permanently stuck README badges.

## Actionable Items & Execution Steps

1. **[Completed]** Fix `.github/workflows/ci.yml` missing `run` directive.
2. **[Completed]** Fix `.github/workflows/ci.yml` `[email protected]` typo to `BYOB@v1.3.0`.
3. **[Completed]** Fix `.github/workflows/race-detector.yml` missing `cache-suffix: 'race'` parameter.
4. **[Pending]** Write RCA to `.lovable/cicd-issues/16-github-actions-yaml-parse-fatal.md`.
5. **[Pending]** Update `.lovable/cicd-issues/index.md` and `.lovable/strictly-avoid.md`.
6. **[Pending]** Execute local test verification.
7. **[Pending]** Run pre-commit checklist.
8. **[Pending]** Commit with `fix(ci)` convention.
9. **[Pending]** Bump `version.json` MINOR version.
10. **[Pending]** Update `.lovable/memory/release-architecture-map.md`.

## Coding Guidelines Checklist

- [x] Boolean conventions used (is/has prefixes, no negatives).
- [x] No garbage variable names used.
- [x] No magic strings or numbers.
- [x] Error management protocols followed (AppError/AppException).
- [x] Code strictly semantic and formatted.
