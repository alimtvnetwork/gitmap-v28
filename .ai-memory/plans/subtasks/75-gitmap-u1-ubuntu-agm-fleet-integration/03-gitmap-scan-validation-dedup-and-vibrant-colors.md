# Subtask 75.3: Scan Validation Warning Deduplication & Vibrant Terminal Palette

- **Parent Plan:** [75-gitmap-u1-ubuntu-agm-fleet-integration.md](../../pending/75-gitmap-u1-ubuntu-agm-fleet-integration.md)
- **Spec Reference:** [02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md](../../../../02-spec/21-app/205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md)
- **Status:** Completed
- **Target Area:** `cli/cmdscanner`, `cli/clicolors`, `cli/cmdvisualizer`, `cli/cmdprojectmanager`

## Objective
Address visual and logging defects identified in user terminal captures: deduplicate repeated scan validation warnings, fix duplicate phrasing in Project Manager error text, and upgrade the dull terminal tree green/yellow palette to vibrant High-Intensity colors.

## Implementation Details
1. In `cli/cmdscanner/`, deduplicate validation error messages so that repositories missing clone URLs (e.g. `.oh-my-zsh`) are only alerted once instead of repeatedly across CSV and JSON export routines.
2. In `cli/cmdprojectmanager/`, correct duplicated phrasing in `project-manager extension storage dir not found near project-manager extension storage dir not found`.
3. In `cli/clicolors/` and tree renderer (`cli/cmdvisualizer/`), upgrade dark green (`\033[32m`) and dull olive (`\033[33m`) to High-Intensity Bright Bold Green (`\033[92;1m` / `#00FF7F`) and Bright Bold Gold (`\033[93;1m` / `#FFD700`).
