---
name: Version Bump Procedure
description: When user says "bump minor/patch/version" or "release", update ONLY Version const + CHANGELOG. NEVER touch .gitmap/release/.
type: feature
---
# Version Bump Procedure

When the user says "bump the minor", "bump version", "release it", "cut a release", or similar, update ONLY these files:

## Files to update (atomic — batch in parallel)

1. **`gitmap/constants/constants.go`** — `const Version = "X.Y.Z"`. Minor = X.Y+1.0, patch = X.Y.Z+1, major = X+1.0.0.
2. **`src/constants/index.ts`** — `export const VERSION = "vX.Y.Z"` (note the `v` prefix). MUST match the Go literal exactly. Enforced by `src/test/version-sync.test.ts` — CI fails if the two drift. Stale-VERSION drift was the v3.0.0 → v4.22.0 follow-up bug; do not let it recur.
3. **`CHANGELOG.md`** — rename `## Unreleased — ...` to `## vX.Y.Z — (YYYY-MM-DD) — <summary>`.

## ABSOLUTE RULES

- **NEVER create, modify, or delete ANY file under `.gitmap/release/` or `.gitmap/release-assets/`.** These are managed exclusively by the `gitmap` CLI tool itself. The user has corrected this multiple times. The previous "AI must update them" instruction was REVERSED.
- **DO NOT touch ANY file under the `gitmap/` source folder during a bump EXCEPT the single `Version = "X.Y.Z"` line in `gitmap/constants/constants.go`.** No new `.go` files, no edits to other constants, no refactors.
- The user runs `gitmap release` themselves to produce the JSON metadata in `.gitmap/release/`.
- Bump `constants.SchemaVersionCurrent` ONLY when Migrate() gains a new structural step (new CREATE/ALTER/phase/seed/ID rename). Cosmetic changes do NOT require a schema bump.
