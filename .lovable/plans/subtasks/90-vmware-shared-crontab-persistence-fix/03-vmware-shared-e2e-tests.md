# Subtask 03: Author Unit & E2E Tests for Crontab Persistence & Ubuntu

**Parent Plan:** `.lovable/plans/pending/90-vmware-shared-crontab-persistence-fix.md`  
**Target Files:** `gitmap/cmd/vmware_crontab_test.go`, `gitmap/cmd/vmware_shared_e2e_test.go`  

## Objectives
1. Create `gitmap/cmd/vmware_crontab_test.go`:
   - Test `buildUpdatedCrontab` with empty string -> ensures valid crontab syntax starting directly with `@reboot ...`.
   - Test `buildUpdatedCrontab` with existing jobs -> ensures existing jobs are preserved and separated by newline.
   - Test `isCrontabPersisted` with positive and negative inputs.
   - Test `cleanCrontabOutput` when given `"no crontab for a"` -> returns `""`.
2. Create `gitmap/cmd/vmware_shared_e2e_test.go`:
   - End-to-end simulated flow:
     - Step 1: User has no crontab (simulates `crontab -l` exiting with code 1 and `"no crontab for username"`).
     - Step 2: `ensureCrontabPersistence()` successfully writes `@reboot` line to `crontab -` without "bad minute" error.
     - Step 3: Re-running `ensureCrontabPersistence()` is idempotent (returns nil early without writing).
     - Step 4: Existing user crontab entries are preserved.
     - Step 5: `runVmwareSharedStatus` reports `Crontab Persistence: registered=true`.
3. Follow all coding guidelines: functions $\le 15$ lines, blank line before returns, affirmative booleans, zero nested ifs.
