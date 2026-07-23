---
name: marker-comments
description: Decentralized opt-in for completion generator using // gitmap:cmd top-level on const blocks and // gitmap:cmd skip on individual specs
type: feature
---

# Marker-Comment Opt-In for Completion Generator (v3.0.0)

## Why
Replaces the centralized `sourceFiles` allowlist + `skipNames` map in `gitmap/completion/internal/gencommands/main.go`. Domain owners now control inclusion **locally**, without ever editing the generator.

## How

1. In any `gitmap/constants/*.go`, add `// gitmap:cmd top-level` to the **doc comment** immediately above a `const (...)` block.
2. To exclude a single spec inside an opted-in block, append `// gitmap:cmd skip` as a trailing line comment.

```go
// gitmap:cmd top-level
// Bookmark commands.
const (
    CmdBookmarkAdd    = "add"    // gitmap:cmd skip
    CmdBookmarkList   = "list"
    CmdBookmarkRemove = "remove"
)
```

3. Run `go generate ./...` in `gitmap/` to regenerate `allcommands_generated.go`.

## CI enforcement

A `generate-check` job in `.github/workflows/ci.yml` runs `go generate ./...` and fails via `git diff --exit-code` if anything drifts. Wired into `test-summary`'s `needs` so the SHA-passthrough cache cannot mark a green run unless the drift check passes.

## Migration done in v3.0.0
- 40 const blocks across 34 constants files annotated `// gitmap:cmd top-level`.
- 52 `// gitmap:cmd skip` annotations mirror the previous policy exactly.
- `allcommands_generated.go` regenerates byte-for-byte identically (143 entries).
- `gitmap/completion/completion.go::manualExtras` is now empty with a doc comment pointing future contributors at the marker convention.

## Migration guide
Full guide for external contributors lives at the top of `CHANGELOG.md` under **Migration guide — v2.x → v3.0.0**.
