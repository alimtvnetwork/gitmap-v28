# Issue Domain 10: Macro Automation and UI Hangs

- **Domain:** Interactive Macros, Child Process Spawning, and VMware Mounts
- **Status:** Consolidated Problem & Resolution Matrix

## 1. Non-Idempotent Macro Removal Failures
- **Symptoms:** Macro playback crashed on step 2 when deleting a directory that was already removed.
- **Root Cause:** Macro engine treated non-zero exit codes from `rm` or `Remove-Item` as fatal.
- **Resolution:** Added step preconditions (`absent`, `exists`), safely skipping already-satisfied operations.

## 2. Orphaned Development Server Processes
- **Symptoms:** Subsequent macro runs failed with port 8080 already in use.
- **Root Cause:** Child server processes were detached without process group lifecycle tracking.
- **Resolution:** Implemented process group termination (`killProcessGroup`) on macro cleanup.
