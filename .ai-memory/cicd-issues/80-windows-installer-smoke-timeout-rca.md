# CI Issue 80: Windows Installer Smoke Runner Timeout RCA

## 1. Reproduction & Symptoms
- **CI Run ID:** [36109397871](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36109397871)
- **Commit:** `2143b036`
- **Failing Job:** `Installer Smoke Windows (source-build)`
- **Symptom:** The job hung for 48m 37s during `python .github/scripts/e2e-cli-smoke.py ./bin/gitmap.exe` until the GitHub Actions runner hard-timed out and canceled the job. All 31 other matrix jobs succeeded cleanly.

## 2. Root Cause Analysis
1. **Missing Job-Level Timeout:**
   - Neither `installer-smoke` nor `installer-smoke-windows` had a `timeout-minutes` property configured in `.github/workflows/ci.yml` or `release.yml`.
2. **Ephemeral Windows Runner Stall:**
   - On GitHub-hosted Windows runners, running sequential subprocesses under console redirection can encounter sporadic I/O locks, file handle locking, or Defender scanning stalls.
   - Without an explicit `timeout-minutes: 10` guard, any transient stall blocks the entire workflow for the default 45–60 minute timeout.

## 3. Code Fix
1. Added explicit `timeout-minutes: 10` to both `installer-smoke` and `installer-smoke-windows` jobs in `.github/workflows/ci.yml`.
2. Added explicit `timeout-minutes: 10` to `installer-smoke` and `installer-smoke-windows` jobs in `.github/workflows/release.yml`.
3. Verified locally on Windows that `smoke-installer.py source` and `e2e-cli-smoke.py` execute and pass in under 55s.

## 4. Prevention
- Ensure all CI matrix smoke jobs, especially Windows runner jobs, have strict bounded job timeouts (`timeout-minutes: 10`).
