# SkipMeta Integration Test

## Overview

Verifies that `release-branch` (`rb`) and `release-pending` (`rp`) commands
**never** write `.gitmap/release/vX.Y.Z.json` or update `latest.json`.
These commands process already-existing branches or metadata — only the
primary `release` command writes release metadata files.

The `SkipMeta` flag on `release.Options` controls this behaviour. Both
`ExecuteFromBranch` and `releaseFromMetadata` set `SkipMeta: true`,
which causes `performRelease` to skip the `writeMetadata` call entirely.

---

## What SkipMeta Suppresses

When `opts.SkipMeta` is true, `performRelease` skips:

1. `writeMetadata()` — would create `.gitmap/release/vX.Y.Z.json` and
   call `updateLatestIfStable()` to overwrite `latest.json`.
2. By extension, no `LastMeta` global is set from the original-branch
   metadata write path (finalize still sets it during push).

The auto-commit step (`AutoCommit`) still runs unless `--no-commit` is
also set, but with no new metadata files there is nothing release-specific
to commit.

---

## Test Plan

### Test 1: `ExecuteFromBranch` does not write metadata

**Setup:**
1. Create a temp Git repo with an initial commit (`t.TempDir()`).
2. Override `constants.DefaultReleaseDir` to a temp path.
3. Create branch `release/v9.0.0` from HEAD.
4. Return to `main`.

**Execute:**
```go
release.ExecuteFromBranch("release/v9.0.0", "", "", false, true, true)
```
- `dryRun: true` — avoids real push/tag side effects.
- `noCommit: true` — avoids auto-commit.

**Assert:**
- `release.ReleaseExists(v)` returns `false`.
- No `v9.0.0.json` file exists in `DefaultReleaseDir`.
- `latest.json` either does not exist or is unchanged from before the call.

---

### Test 2: `ExecutePending` does not write metadata (branch-based)

**Setup:**
1. Create a temp Git repo with an initial commit.
2. Override `constants.DefaultReleaseDir` to a temp path.
3. Create branch `release/v8.0.0` from HEAD (no tag).
4. Return to `main`.

**Execute:**
```go
release.ExecutePending("", "", false, true, true)
```
- `dryRun: true`, `noCommit: true`.

**Assert:**
- No `v8.0.0.json` file exists in `DefaultReleaseDir`.
- `latest.json` either does not exist or is unchanged.

---

### Test 3: `ExecutePending` does not write metadata (metadata-based)

**Setup:**
1. Create a temp Git repo with an initial commit.
2. Override `constants.DefaultReleaseDir` to a temp path.
3. Write a seed `.gitmap/release/v7.0.0.json` with the HEAD commit SHA
   (simulating a metadata-only pending release).
4. Ensure no `release/v7.0.0` branch or `v7.0.0` tag exists.

**Execute:**
```go
release.ExecutePending("", "", false, true, true)
```
- `dryRun: true`, `noCommit: true`.

**Assert:**
- The pre-existing `v7.0.0.json` is unchanged (byte-equal).
- No new metadata files were created.
- `latest.json` either does not exist or is unchanged.

---

### Test 4: Primary `release` DOES write metadata (control test)

**Setup:**
1. Create a temp Git repo with an initial commit.
2. Override `constants.DefaultReleaseDir` to a temp path.

**Execute:**
```go
release.Execute(release.Options{
    Version: "v6.0.0",
    DryRun:  false,
    NoCommit: true,
    SkipMeta: false,
})
```
Note: This may need stubbing of push/tag operations or a local-only
test remote to avoid network calls. Alternatively, verify the metadata
write path in isolation by calling `writeMetadata` directly.

**Assert:**
- `v6.0.0.json` exists in `DefaultReleaseDir`.
- `latest.json` references `v6.0.0`.

---

## Test Helpers

### Git repo scaffold

```go
func setupTestRepo(t *testing.T) string {
    t.Helper()
    dir := t.TempDir()
    run(t, dir, "git", "init")
    run(t, dir, "git", "checkout", "-b", "main")
    writeFile(t, dir, "README.md", "test")
    run(t, dir, "git", "add", ".")
    run(t, dir, "git", "commit", "-m", "initial")
    constants.DefaultReleaseDir = filepath.Join(dir, ".gitmap", "release")
    return dir
}
```

### Metadata existence check

```go
func assertNoMetadata(t *testing.T, version string) {
    t.Helper()
    path := filepath.Join(constants.DefaultReleaseDir, version+".json")
    _, err := os.Stat(path)
    if err == nil {
        t.Errorf("metadata file should not exist: %s", path)
    }
}
```

---

## File Layout

| File | Purpose |
|------|---------|
| `tests/release_test/skipmeta_test.go` | All four test cases |

---

## Acceptance Criteria

1. Tests 1–3 confirm zero metadata side effects from `rb` and `rp`.
2. Test 4 confirms the primary `release` command still writes metadata
   (guards against accidental `SkipMeta: true` propagation).
3. All tests use `t.TempDir()` — no shared state, no cleanup needed.
4. Tests pass with `go test ./tests/release_test/ -run TestSkipMeta`.
