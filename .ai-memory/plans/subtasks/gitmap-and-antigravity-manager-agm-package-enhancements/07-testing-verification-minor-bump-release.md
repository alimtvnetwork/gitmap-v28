# Subtask 07: Targeted Verification, Spec/Docs Update & Minor Release Ceremony

**Parent Task:** `gitmap-and-antigravity-manager-agm-package-enhancements`  
**Subtask ID:** `subtask-07`  
**Status:** `IN_PROGRESS`  
**Owned Files:**  
- `version.json`
- `changelog.md`
- `readme.md`
- `cli/version/version.go`
- `02-spec/21-app/readme.md`
- `.ai-memory/plans/readme.md`

## Objective
Execute targeted unit tests on modified packages, update project manifests and documentation indices, perform the minor version bump from `v6.571.4` to `v6.572.0`, and conclude with an atomic GitMap release commit.

## Implementation Details
1. Execute file-scoped unit tests for `cli/cmdpull` and `cli/cmdagy`.
2. Update `version.json`:
   - Bump version to `6.572.0`.
3. Update `cli/version/version.go`:
   - Set version string to `6.572.0`.
4. Update `changelog.md`:
   - Record `## [6.572.0] - 2026-10-18`:
     - Single-repo pull UX simplification (clean progressing message & summary).
     - Fix progress bar rune width calculation & safety right margin (waterfall wrapping fix).
     - Forensic RCA & prompt guardrails against unnecessary AI shell wrapping.
     - AGM package updates review & AGY prompt rerun/review parity.
     - Standalone image URL and content-addressable storage architecture.
     - RC Agent & Heatmap RC 6-stage reconciliation verification.
5. Update `02-spec/21-app/readme.md` and `.ai-memory/plans/readme.md`.
6. Atomic release commit via `gitmap cpf`.
