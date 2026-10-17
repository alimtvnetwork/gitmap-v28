# Spec 278: Acceptance Criteria & Evidence

## Acceptance Criteria

### AC-01: GCM Credential Store
- `credential.credentialStore` set to `wincredman` or `dpapi` instead of `wincred`.
- `gitmap fix-credential` resolves Windows credential store without fatal errors.

### AC-02: `gitmap reclone <url>`
- Passing a valid URL like `https://github.com/alimtvnetwork/story-telling-v1` to `gitmap reclone` resolves the repo name, cleans any broken/partial directory, and re-clones without "unsupported file extension" error.

### AC-03: Generic Typed Worker Pool for `gitmap fix --all`
- `runFixAll` dispatches remediation steps in parallel up to CPU logical core count (e.g. 20) with minimum 4.
- Real-time progress indicators display active repositories being fixed.

### AC-04: `gitmap pull-all` Concurrency & Menu
- Concurrency scales dynamically up to CPU processor count (up to 20 cores).
- Post-pull menu offers Fix All as option 1, Local Commit as option 2, Stash as option 3, Step-through as option 4, Fleet as option 5.
- Single-repo step-through prints clear tree-view diff and preview steps with color coding.

### AC-05: `gitmap update` Daily Probe Throttle
- Sibling probe skips if already probed within 24 hours.
- Header shows current version in yellow/white with clear transition.

### AC-06: `repo-cache` & Muse Memory Digest
- `repo-cache` repository initialized and synced with `12-muse-ai/`.
- `gitmap llm train muse` loads training digests from `repo-cache` by default.

### AC-07: Subcommand Hyphenated Aliases
- Hyphenated aliases available for all 275 commands and `llm-train-muse`.

### AC-08: Heatmap Skill & Guidelines
- `.agents/skills/gitmap/SKILL.md` updated with heatmap instructions and file contention commands.
