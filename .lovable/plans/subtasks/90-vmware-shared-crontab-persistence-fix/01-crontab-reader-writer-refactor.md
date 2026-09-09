# Subtask 01: Extract and Implement Safe Crontab Reader & Writer (`vmware_crontab.go`)

**Parent Plan:** `.lovable/plans/pending/90-vmware-shared-crontab-persistence-fix.md`  
**Target Files:** `gitmap/cmd/vmware_crontab.go`, `gitmap/cmd/vmware_shared.go`  

## Objectives
1. Create `gitmap/cmd/vmware_crontab.go`:
   - `readCurrentCrontab() string`: Runs `crontab -l` using `Output()` (stdout only).
   - If execution fails or output contains `"no crontab for"`, returns clean `""`.
   - `isCrontabPersisted() bool`: Checks if `vmhgfs-fuse` and `/mnt/hgfs` are in current crontab.
   - `buildUpdatedCrontab(current, rebootLine string) string`: Appends `rebootLine` with valid newlines without prepending error text.
   - `writeCrontab(content string) error`: Pipes clean crontab into `crontab -`.
   - `ensureCrontabPersistence() error`: Coordinates the flow and returns typed `AppError` on failure.
   - Command seams `crontabCommandFunc = exec.Command` for testability.
2. Remove `ensureCrontabPersistence` from `gitmap/cmd/vmware_shared.go`, reducing its file length well under 250 lines.
3. Follow all coding guidelines: $\le 15$ line functions, blank line before returns, affirmative booleans, zero nested ifs.
