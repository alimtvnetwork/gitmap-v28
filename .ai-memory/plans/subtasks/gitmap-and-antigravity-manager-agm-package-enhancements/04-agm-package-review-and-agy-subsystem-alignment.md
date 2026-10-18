# Subtask 04: AGM Package Review & AGY Subsystem Alignment

**Parent Task:** `gitmap-and-antigravity-manager-agm-package-enhancements`  
**Subtask ID:** `subtask-04`  
**Status:** `IN_PROGRESS`  
**Owned Files:**  
- `cli/cmdagy/agy_rerun.go`
- `cli/cmdagy/agy_prompt_subcmds.go`
- `cli/cmdagy/agy_running_prompts_backup.go`

## Objective
Review recent changes in Antigravity Manager (AGM v4.223.0 - v4.227.3) including E9001 elimination, 30s broadcast floor, fail-fast missing workspaces, and prompt rerun/review. Align GitMap's AGY subsystem commands (`gitmap agy rerun`, `backup-running-prompts`, `restore-running-prompts`, `prompt`) to ensure full synchronization with AGM.

## Implementation Details
1. Document AGM improvements:
   - Sibling account isolation prevents concurrent running instances from stealing accounts.
   - Missing workspace directory detection in `queue.rs` transitions status to `failed` immediately.
   - Background pollers enforced `>= 30s` sleep floor.
   - HMAC-SHA256 JWT bridge over `:18188`.
2. Inspect and verify GitMap AGY commands:
   - `gitmap agy running-prompts ls`
   - `gitmap agy rerun` (`rr`, `rra`, `rwi`, `rwc`)
   - `gitmap agy backup-running-prompts` (`brp`) and `restore-running-prompts` (`rrp`)
3. Ensure image path extraction (`ExtractAllImagePathsFromPrompt`) and verification suffixes match AGM's prompt rerun options.
