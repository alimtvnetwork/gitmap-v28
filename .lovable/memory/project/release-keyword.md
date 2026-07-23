---
name: Release Keyword Shortcut
description: When user says "release" (alone or as a verb), it means the full bump procedure — bump minor, add CHANGELOG entry, pin version in root README. Same as version-bump-procedure.md.
type: preference
---

# "release" = full bump procedure

When the user says **"release"**, **"mark this as release"**, or similar shorthand, treat it as a request to:

1. Bump MINOR in `gitmap/constants/constants.go` (`Version`)
2. Mirror to `src/constants/index.ts` (`VERSION`, with `v` prefix)
3. Add a new `## vX.Y.Z — (YYYY-MM-DD) — <summary>` entry at the top of `CHANGELOG.md`
4. Pin the new version everywhere it appears in root `README.md` (sed replace of the previous pinned version)

Do NOT touch `.gitmap/release/` — that is managed by the `gitmap` CLI itself (see `version-bump-procedure.md` and the Core constraint).

**How to apply:** Always batch the 4 edits in parallel. Use `sed -i 's/vOLD/vNEW/g' README.md` for the README pin refresh.
