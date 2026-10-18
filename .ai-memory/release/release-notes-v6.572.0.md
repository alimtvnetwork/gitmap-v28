# GitMap v6.572.0

## What's Changed in v6.572.0

### Key Enhancements & Defect Fixes

1. **GitMap Single-Repo Pull UX Simplification (`cli/cmdpull`)**:
   - Replaced verbose 0%-100% counting/compressing lines and 8-column batch table when running in a single repository with a clean, lightweight progressing status (`↓ pulling remote changes...`).
   - Added compact 1-line summary reporting commit delta (`from..to` or `up-to-date`), diffstat (`+X/-Y, N files`), working tree state (`clean` vs `dirty: details`), and auto-merge execution notice.
   - Preserved full batch table display for multi-repo fleet commands (`gitmap pa`, `gitmap pat`).

2. **Progress Bar Column Width & Waterfall Wrapping Fix (`cli/cmdpull`)**:
   - Integrated `github.com/mattn/go-runewidth` to compute accurate terminal column widths for wide Unicode characters and emojis (`⚡`, `🔍`, `🔀`).
   - Enforced a defensive right-margin of `termWidth - 4` preventing terminal emulator auto-wrap cursor escape to row Y+1.
   - Throttled non-TTY progress bar rendering to suppress intermediate milestone spam in background agent tasks and CI logs.

3. **Forensic Root Cause Analysis (RCA) & Guardrails for AI Shell Wrapping**:
   - Conducted deep forensic RCA resolving why AI agents unnecessarily execute shell commands via `gitmap ps "git ..."` instead of native commands.
   - Updated `.ai-memory/strictly-avoid.md`, `.agents/skills/gitmap/SKILL.md`, and `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` with explicit total bans on wrapping tools in `gitmap ps`.
   - Clarified Worker Git Ban: read-only inspections (`git show`, `git log`) are permitted for revision checks, while index-mutating commands remain strictly restricted.
   - Implemented native `gitmap show <rev> [-- <path>]` and `gitmap log [-n <N>] [--stat]` commands in `cli/cmdgit/` to close GitMap's git inspection tooling gap.

4. **Antigravity Manager (AGM) Package Alignment (`cli/cmdagy`)**:
   - Documented recent AGM updates: forensic E9001 elimination (sibling instance account isolation), missing workspace fail-fast validation in `queue.rs`, `>= 30s` background poller sleep floor, and HMAC-SHA256 JWT bridge (`:18188`).
   - Aligned GitMap AGY commands (`gitmap agy rerun`, `backup-running-prompts`, `restore-running-prompts`, `prompt`) with AGM prompt review and rerun options.

5. **Standalone Image URL & Lightweight Prompt Architecture**:
   - Authored specification in `02-spec/21-app/gitmap-and-antigravity-manager-agm-package-enhancements/03-image-url-architecture.md`.
   - Implemented `cli/cmdagy/agy_image_assets.go` extracting inline base64 data URIs into content-addressable SHA-256 disk files (`.ai-memory/assets/img-<sha256>.<ext>`) and rewriting markdown to lightweight standalone URLs (`http://127.0.0.1:18188/assets/img-...`) with dual-reference comments.

6. **RC Agent & Heatmap RC Continuous Reconciliation**:
   - Verified the 6-stage continuous reconciliation protocol (`RunRcHeatmapSyncWithConflictResolution`) across AI namespace directories (`01-antigravity/`, `02-muse/`) ensuring collision-free SQLite sync with `repo-cache`.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.572.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.572.0/install.sh | sh
```
