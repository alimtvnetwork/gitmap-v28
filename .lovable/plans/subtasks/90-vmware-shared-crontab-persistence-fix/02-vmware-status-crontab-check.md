# Subtask 02: Add Crontab Persistence Reporting to VMware Shared Status

**Parent Plan:** `.lovable/plans/pending/90-vmware-shared-crontab-persistence-fix.md`  
**Target Files:** `gitmap/cmd/vmware_status.go`  

## Objectives
1. In `gitmap/cmd/vmware_status.go`:
   - In `runVmwareSharedStatus()`:
     - Print `Crontab Persistence: registered=%t` using `isCrontabPersisted()`.
     - Fulfill the promise in the help text: `shared status Check /mnt/hgfs mount, desktop symlink & crontab persistence`.
2. Follow all coding guidelines: functions $\le 15$ lines, blank line before returns, affirmative booleans, zero nested ifs.
