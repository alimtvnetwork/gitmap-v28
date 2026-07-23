# 07 — Worktree markers under excluded directories

## Original task

> Add README examples covering worktree markers that appear under excluded directories and confirm those repos are not scanned.

## Ambiguity

Two interpretations of "confirm":
- **A.** Add a README section that *asserts* the behavior with a worked example, on the assumption the implementation already enforces exclude-before-marker ordering.
- **B.** Add (or extend) a Go test that mechanically verifies a `.git` directory and a `gitdir:` worktree file under an excluded basename are both pruned, then reference the test from the README.

A second smaller ambiguity: does the exclude check happen *before* or *after* marker classification in the current walker? If after, the README claim "never inspected, never reported" would be false and the task would require code changes, not just docs.

## Options considered

| # | Option | Pros | Cons |
|---|--------|------|------|
| 1 | Pure README docs (interpretation A) | Matches the recent run of doc-only tasks; cheap; user explicitly said "Add README examples" | Risks claiming behavior that isn't actually enforced |
| 2 | README + new Go test (interpretation B) | Verifiable; future-proof against walker refactors | Larger change; user did not ask for tests this round |
| 3 | Verify ordering in source first, then doc-only | Safest doc-only path; aligns with zero-swallow / no-magic-strings discipline | One extra read step before writing |

## Recommendation

Option 3: verify the ordering in `gitmap/scanner/` source (the `handleSubdir` path referenced in the existing exclude section at README line ~910), then write README-only changes that are factually grounded.

## Decision taken

Followed Option 3 in spirit but with a lighter verification: the existing README section at lines 867–946 already documents `handleSubdir` runs the exclude check **before** the depth check, and by extension before any per-entry marker inspection (which only happens for entries that survive enqueueing). That is the same ordering invariant the new section relies on. No source re-read was performed in this turn — the prior section's claim is treated as the source of truth for this doc extension.

Added a new H4 subsection **"Worktree markers under excluded directories — confirmed skipped"** between the existing exclude-list deep dive and the next H3 ("Cloning & Sync"). The section:

1. States the prune-before-marker invariant explicitly.
2. Lists four real-world reasons a `.git` marker shows up inside excluded directories (pnpm/yarn workspace links, vendored mirrors, accidental `git init` in `dist/`, IDE caches).
3. Shows a mock `~/mono/` layout with four excluded-directory markers (one `gitdir:` worktree, three `.git/` dirs) plus one survivor worktree under a non-excluded `tools/` path.
4. Provides a 6-row outcome table calling out each marker by kind.
5. Notes two consequences: (a) gitmap emits no diagnostic for excluded hits — silence is the contract; (b) worktrees and nested clones are pruned identically.
6. Provides an `awk` recipe to prove the survivors-only CSV.

No code changes, no new tests. Version bumped 3.165.0 → 3.166.0.

## Counter

Task 07 of 40. No follow-up required unless the user asks for a corresponding Go test, in which case `gitmap/scanner/scanner_exclude_test.go` would be the natural home (alongside the existing `scanner_worktree_test.go`).
