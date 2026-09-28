# Plan 184: Antigravity Manager (`gitmap agm update`) Fix, Parallel 1-to-1 Running-Prompts SQLite Backup with Media, and E2E Account Switch (`98%` Test -> `15%` Production Reset)

**Status:** Completed (`2026-09-28`)
**Canonical Spec:** [02-spec/21-app/89-agm-update-fix-and-account-switch-e2e.md](file:///d:/work/gitmap/02-spec/21-app/89-agm-update-fix-and-account-switch-e2e.md)

---

## 1. Consolidated Subtasks & Verified Outcomes

- [x] **Subtask 01 (`01-agm-update-version-pin-and-verify-fix.md`)**:
  - Removed `ConsoleHost_history.txt` and `PSConsoleReadLine` history scraping from `d:\work\Antigravity-Manager\install.ps1` (`cc2171c6`) and scoped `Resolve-PinnedVersion` strictly to `Antigravity-Manager` URLs.
  - Implemented `resolveLatestAgManagerReleaseVersion()` in `cli/cmdinstall/installagmanager.go` and passed `-Version '<latest>' -Update -NoLaunch` plus `AGM_VERSION=<latest>` on a cache-busted `install.ps1` URL.
  - Updated `findAgManagerWindowsPath()` (`cli/cmdinstall/installagmanager_windows.go`) and `isAgManagerInstalled()` (`cli/cmdinstall/installagmanager_exec.go`) to prioritize `C:\Users\Administrator\AppData\Local\Programs\agm-alim\agm-alim.exe` (`agm-alim`).
  - Live verified `gitmap agm update` upgrading `v4.85.0 -> v4.89.0` with `✔ Antigravity Manager (4.89.0) completed successfully.`

- [x] **Subtask 02 (`02-parallel-running-prompts-backup-with-media.md`)**:
  - Extended `RunningPromptRecord` with `MediaPaths []string` and `PromptBackupSummary` with `ProjectNames []string` in `cli/store/backup_prompts_split_types.go`, `cli/store/backup_prompts_split_db.go`, and `cli/store/backup_prompts_split_ops.go`.
  - Refactored `CollectActiveAndQueuedPrompts` in `cli/cmdagy/agy_running_prompts_backup.go` to remove `cwd` injection, enforce `!IsRestrictedSystemOrHomeDir`, run parallel 1-to-1 goroutine backup per active project (`sync.WaitGroup` + `sync.Mutex`), and extract attached image/screenshot paths (`extractPromptMediaPaths`).
  - Ensured all CLI and notification summaries list project counts and human-readable project names (`scripts-fixer`, `Antigravity-Manager`, `gitmap`) with zero raw project UUIDs.

- [x] **Subtask 03 (`03-account-switch-threshold-refresh-email-supabase-lock.md`)**:
  - Implemented `gitmap agy account-switch` (and top-level `gitmap account-switch` / `gitmap asw`) across `cli/cmdagy/agy_account_switch_cmd.go`, `cli/cmdagy/agy_account_switch_engine.go`, `cli/cmdagy/agy_account_switch_notify.go`, `cli/cmdagy/agy_cmd.go`, and `cli/cmd/root.go`.
  - Implemented candidate ranking by highest credit, candidate refresh verification (`refreshedCredit == initialCredit`), Email and Supabase distributed VM lock checks (`isAccountLockedByOtherVM`), parallel SQLite running-prompts backup, fast-forward delegation, post-switch restore & liveness verification, and Telegram + Email notifications.

- [x] **Subtask 04 (`04-rich-cli-help-instance-e2e-and-minor-release.md`)**:
  - Added rich boxed CLI help menus with examples to `gitmap agm help`, `gitmap agy running-prompts help`, and `gitmap agy account-switch help`.
  - Executed live E2E test in temporary instance (`gitmap agy account-switch test --threshold 98 --instance test-e2e-inst`), confirmed full 98% switch + backup/restore + notification + instance cleanup, and verified default threshold reset to `15%` (`DefaultAccountSwitchThreshold = 15`).
