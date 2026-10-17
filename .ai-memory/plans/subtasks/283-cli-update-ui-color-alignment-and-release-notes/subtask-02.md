# Subtask 02: Release Notes Propagation in Update Data Pipeline

## Status: DONE
- **Parent Plan**: `.ai-memory/plans/283-cli-update-ui-color-alignment-and-release-notes.md`
- **Owner**: Worker 02
- **Owned Files**:
  - `cli/cmdupdate/update_fallback_loop.go`
  - `cli/cmdupdate/types.go`
  - `cli/cmdupdate/update_binary_swap.go`

## Objectives
1. In `cli/cmdupdate/update_fallback_loop.go`:
   - Add `Body string `json:"body"`` to `GitHubRelease`.
   - In `fetchExplicitReleaseCandidate` and `probeReleaseCandidate`, populate `cand.ReleaseNotes = ghRel.Body` (or `rel.Body`).
   - Add `ReleaseNotes string` to `ReleaseCandidate` struct.
2. In `cli/cmdupdate/types.go`:
   - Add `ReleaseNotes []string `json:"release_notes,omitempty"`` to `UpdateResult`.
3. In `cli/cmdupdate/update_binary_swap.go`:
   - In `finalizeUpdate`, extract 1–2 bullet lines from `cand.ReleaseNotes` using `ExtractReleaseNotesSummary(cand.ReleaseNotes, 2)` and assign to `res.ReleaseNotes`.
4. Sizing: ensure all modified Go files stay <= 300 lines with nesting depth <= 1.

## Verification
- `cli/cmdupdate/types.go`: 52 lines (<= 300, depth <= 1)
- `cli/cmdupdate/update_binary_swap.go`: 296 lines (<= 300, depth <= 1)
- `cli/cmdupdate/update_fallback_loop.go`: 296 lines (<= 300, depth <= 1)
- Release notes body from GitHubRelease correctly flows into ReleaseCandidate and is summarized into UpdateResult.ReleaseNotes via ExtractReleaseNotesSummary.
