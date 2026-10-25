# Subtask 329.1 — Version Update Downgrade Guard & Sibling Probe Fix

> **Parent Plan:** `.ai-memory/plans/pending/329-gitmap-context-menu-and-version-fix.md`  
> **Target Package:** `cli/cmdupdate/`  
> **Status:** Completed

---

## 1. Objectives
1. Implement semantic version comparison in `cli/cmdupdate/update_version_target.go` (`CompareSemVer(v1, v2 string) int` and `IsNewerVersion(target, current string) bool`).
2. Update `IsAlreadyUpdated` and `isUpdateSkipped` in `cli/cmdupdate/update_binary_swap.go` to reject version updates where `targetVersion <= currentVersion` unless `IsAllowDowngrade` is set.
3. Add `--allow-downgrade` flag to `UpdateOptions` in `ParseUpdateOptions`.
4. Update `resolveLatestRepoSlug` in `cli/cmdupdate/updateprobe.go` to verify that any probed sibling repository actually has releases newer than the current version before accepting the sibling slug.
5. Add unit tests in `cli/cmdupdate/update_version_target_test.go` covering version comparisons and downgrade guards.
