# RCA: GitHub Actions Silently Failing to Trigger Due to YAML Fatal Errors

## Error Summary
The GitHub Actions pipelines (`ci.yml`, `race-detector.yml`, etc.) stopped triggering entirely on each commit. Furthermore, concurrency groups failed to cancel older jobs, and the README badges remained stuck at "failing". There were no logs in GitHub Actions indicating failure because the files themselves could not be parsed.

## Root Cause Analysis
1. A previous update to `.github/workflows/ci.yml` introduced a Cloudflare email obfuscation artifact: `uses: RubbaBoy/[email protected]` instead of `uses: RubbaBoy/BYOB@v1.3.0`.
2. A separate step `Run installer smoke (source mode)` on Windows lacked a `run:` or `uses:` directive.
3. The `.github/workflows/race-detector.yml` was missing a required input (`cache-suffix`) for the local composite action `setup-go-cached`.
4. Because the workflow files contained invalid YAML and missing required schema fields, GitHub Actions immediately aborted workflow processing on push. This bypassed the `concurrency` cancellation logic, causing the "commit builds are not happening" behavior described by the user.

## Solution Applied
1. Fixed `.github/workflows/ci.yml`: Addressed the missing `run: ./.github/scripts/smoke-installer.ps1 source` directive.
2. Fixed `.github/workflows/ci.yml`: Restored `RubbaBoy/BYOB@v1.3.0` action reference.
3. Fixed `.github/workflows/race-detector.yml`: Appended `with: cache-suffix: 'race'`.
4. Verified syntax correctness locally using `actionlint`.

## What NOT to Repeat
- Never copy-paste GitHub Action references without verifying the tag. Cloudflare email obfuscation can mutate `@v1.2.3` into `[email protected]`.
- Never commit a workflow file without verifying that every `name:` step has a corresponding `run:` or `uses:`.
- Always run `actionlint` locally before committing `.github/workflows/` modifications.
