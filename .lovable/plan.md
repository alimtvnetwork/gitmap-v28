## Problem

On Windows, `gitmap fix-repo` fails post-rewrite with:

```
fix-repo: ERROR gofmt failed: fork/exec D:\go-path\bin\gofmt.exe: The filename or extension is too long.
```

Root cause: `gitmap/cmd/fixrepo_gofmt.go` → `invokeGofmt` batches every rewritten `.go` file into a single `gofmt -w <path1> <path2> ...` exec. Windows' `CreateProcess` caps the command line at ~32,767 chars. On a large repo (hundreds of touched Go files with long absolute paths under `D:\wp-work\...\gitmap\...`), the joined args exceed that cap and Windows returns `ERROR_FILENAME_EXCED_RANGE`, which Go surfaces as "The filename or extension is too long."

This is not a gofmt bug and not a PATH issue: it's an argv length overflow. Same binary works on Linux/macOS where `ARG_MAX` is ~2 MB.

## Fix

Chunk the file list in `invokeGofmt` so each `gofmt -w` call stays well under the Windows limit, and prefer repo-relative paths to shrink each arg.

### Changes to `gitmap/cmd/fixrepo_gofmt.go`

1. Convert absolute paths to repo-relative before exec (shorter args, and `gofmt -w` accepts either).
2. Batch files into chunks whose joined length (plus arg overhead) stays under a safe budget (~30,000 chars) — well below Windows' 32,767 cap with headroom for the exe path and `-w`.
3. Loop, calling `exec.Command("gofmt", "-w", chunk...)` per batch. Aggregate output; fail fast on the first non-zero exit and report which batch failed.
4. Keep the single-batch fast path when the joined length fits (no behavior change on small repos / non-Windows).

### New constants (in `gitmap/constants/constants_fixrepo.go`)

- `FixRepoGofmtMaxCmdLen = 30000` — safe budget (documented rationale: Windows CreateProcess cap is 32767).
- Optional: `FixRepoMsgGofmtBatchFmt` — "  gofmt: formatted %d file(s) across %d batch(es)\n" replacing the existing single-count message when batching kicks in.

### Tests

Add `gitmap/cmd/fixrepo_gofmt_test.go` (or extend the existing e2e in `tests/fixrepo_test/gofmt_e2e_test.go`):

- Unit test for a new pure helper `chunkPathsForGofmt(paths []string, budget int) [][]string`:
  - Empty input → empty output.
  - All-fits input → one chunk.
  - Overflow input (synthetic long paths) → multiple chunks, each ≤ budget, no path split across chunks.
  - Single path longer than budget → still emitted in its own chunk (log a warning; gofmt itself will accept it, only Windows CreateProcess would reject — but a single path can't exceed 32k in practice).
- E2E: generate ~500 files with long synthetic names in a tempdir, run `runFixRepoGofmt`, assert exit success and all files formatted.

### Non-goals

- No change to gofmt discovery, dry-run behavior, or error message shape when gofmt is genuinely missing / broken.
- No parallelization across batches (sequential keeps error attribution simple; latency impact is negligible — 2-3 extra fork/exec on huge repos).

## Files touched

```text
gitmap/cmd/fixrepo_gofmt.go              (edit: chunking + relative paths)
gitmap/constants/constants_fixrepo.go    (edit: add budget constant + optional batch msg)
gitmap/cmd/fixrepo_gofmt_test.go         (new: chunker unit tests)
tests/fixrepo_test/gofmt_e2e_test.go     (edit: add long-path overflow case)
```

## Validation

1. `go test ./gitmap/cmd/... ./gitmap/tests/fixrepo_test/...` passes.
2. Manually reproduce on the user's repo path (documented steps): run `gitmap fix-repo --all` on a checkout with >200 touched `.go` files; expect success and a "formatted N file(s) across M batch(es)" line.
3. `gofmt -l gitmap/...` clean after the run.

## Release

Bump patch version (e.g. v6.80.0 → v6.80.1) across `gitmap/constants/constants.go`, `.gitmap/release/latest.json`, `src/constants/index.ts`, and add a changelog entry under `Fixed`: "fix-repo: batch gofmt calls to avoid Windows CreateProcess arg-length overflow."
