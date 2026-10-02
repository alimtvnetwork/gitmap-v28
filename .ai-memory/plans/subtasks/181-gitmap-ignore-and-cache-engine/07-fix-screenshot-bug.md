---
status: PENDING
---

# Subtask: Fix Screenshot Bug (Ignore Latency)

## Objective
Resolve the critical performance bug where `gitmap scan <group>` takes ~140 seconds due to heavy regex evaluation on every file.

## RCA (Root Cause Analysis)
- **Symptom**: `gitmap scan <group>` takes around 140 seconds to complete.
- **Root Cause**: The WalkDir function evaluates ignore patterns (heavy regexes) on EVERY file it encounters.
- **Compounding Factor**: The cache table currently filters out ignored files on write, meaning if a file is ignored, it isn't cached, so its ignore status has to be re-evaluated on subsequent scans.
- **Solution**: The cache table should store both ignored and non-ignored paths and filter them on **read** rather than write. The ignore boolean result should be cached in Split-DB (`repodb.db`). The evaluation should NEVER re-run unless the `.gitmapignore` file's hash changes.

## Requirements
- Store true/false for the ignore bit in the cache engine.
- Update `WalkDir` to evaluate ignores only when not cached or when `.gitmapignore` hash is updated.
- Update cache reading logic to filter out files with `ignored = true`.
- Confirm 140-second latency drops significantly on cached scans.

## Next Steps
1. Implement the RCA fix in the core scanning and ignore evaluation logic.
2. Update the cache table writes to include ignored paths with `ignored = true`.
3. Add a check to match `.gitmapignore` hash before doing any heavy regex evaluations.
