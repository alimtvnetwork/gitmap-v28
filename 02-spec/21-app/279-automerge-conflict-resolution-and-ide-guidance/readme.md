# Spec 279: GitMap Automerge Optimization, IDE Merge Guidance & Conflict CLI

## Overview
Comprehensive enhancement of GitMap's merge conflict handling, automerge resilience, and developer experience across:
1. **Automerge Optimization & Prevention:** Pre-flight branch divergence detection, histogram diff algorithm (`-X diff-algorithm=histogram`), rename detection, and automated Git `rerere` (Reuse Recorded Resolution) activation.
2. **IDE Integration & Manual VSCode Guidance:** Multi-IDE detection (`code`, `cursor`, `windsurf`, `codium`), one-click IDE launch commands, and structured step-by-step walkthroughs for VS Code's 3-way Merge Editor.
3. **GitMap Conflict Resolution CLI:** First-class `gitmap resolve <repo> [--theirs | --ours | --ide | --abort]`, `gitmap conflicts [repo]`, and expanded `gitmap fix <repo>` sub-actions (`resolve`, `theirs`, `ours`, `abort`, `rerere`).
4. **Post-Pull Interactive Conflict Menu:** Dedicated Lipgloss visual conflict card during `gitmap pull` / `gitmap pa`, listing conflicted files with one-key resolution actions.
5. **Helptext & UI Parity:** CLI markdown help documents (`cli/helpdoc/resolve.md`, `cli/helpdoc/conflicts.md`), terminal `howto` entries, and web dashboard troubleshooting registry.

## Specifications
- [01-scope-and-plan.md](01-scope-and-plan.md)
- [02-acceptance-criteria-and-evidence.md](02-acceptance-criteria-and-evidence.md)
