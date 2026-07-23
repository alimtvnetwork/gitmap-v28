# Refactor: cmd/root.go Dispatch Splitting

## Problem

`cmd/root.go` is 388 lines. The `dispatchMisc` function alone is 166 lines
(lines 191–357), routing 28 commands in a single function. The file
contains five dispatch functions that should be split by domain.

---

## Current Structure (388 lines)

| Function | Lines | Commands | Responsibility |
|----------|-------|----------|----------------|
| `Run` | 12–32 | — | Entry point, migration, alias |
| `dispatch` | 35–62 | — | Top-level router with audit |
| `dispatchCore` | 65–99 | 6 | scan, clone, pull, status, exec, has-any-updates |
| `dispatchRelease` | 102–140 | 7 | release, release-branch, release-pending, changelog, clear-release-json, changelog-gen |
| `dispatchUtility` | 143–188 | 7 | update, update-runner, update-cleanup, revert, revert-runner, version, help, docs |
| `dispatchMisc` | 191–357 | 28 | Everything else (data, TUI, profiles, zip, ssh, prune, etc.) |
| `dispatchProjectRepos` | 360–388 | 5 | go-repos, node-repos, react-repos, cpp-repos, csharp-repos |

Existing companion files: `rootusage.go` (171 lines), `rootflags.go` (66 lines).

---

## Target Layout

Split `dispatchMisc` into domain-specific functions in new files. Keep
`root.go` as the slim orchestrator.

### 1. `root.go` — Entry point + top-level router (≈70 lines)

**Retains:**
- `Run()` (lines 12–32)
- `dispatch()` (lines 35–62) — updated to call new dispatch functions

**Updated dispatch chain:**
```go
func dispatch(command string) {
    auditID, auditStart := recordAuditStart(command, os.Args[2:])

    if dispatchCore(command) { ... }
    if dispatchRelease(command) { ... }
    if dispatchUtility(command) { ... }
    if dispatchData(command) { ... }
    if dispatchTooling(command) { ... }
    if dispatchProjectRepos(command) { ... }

    // unknown command
}
```

---

### 2. `rootcore.go` — Core scan/clone/pull commands (≈45 lines)

**Moves from root.go:**
- `dispatchCore()` (lines 65–99)

---

### 3. `rootrelease.go` — Release commands (≈50 lines)

**Moves from root.go:**
- `dispatchRelease()` (lines 102–140)

---

### 4. `rootutility.go` — System utilities (≈55 lines)

**Moves from root.go:**
- `dispatchUtility()` (lines 143–188) — removes the `dispatchMisc` tail call

---

### 5. `rootdata.go` — Data, history, profiles, TUI commands (≈110 lines)

**New function `dispatchData()`**, extracts from `dispatchMisc`:
- `list` / `ls`
- `group` / `g`
- `multi-group` / `mg`
- `history` / `hi`
- `history-reset` / `hr`
- `stats` / `ss`
- `bookmark` / `bk`
- `export` / `ex`
- `import` / `im`
- `profile` / `pf`
- `diff-profiles` / `dp`
- `cd` / `go`
- `watch` / `w`
- `interactive` / `i`
- `db-reset`
- `amend` / `am`
- `amend-list` / `al`

---

### 6. `roottooling.go` — Dev tooling and misc commands (≈80 lines)

**New function `dispatchTooling()`**, extracts from `dispatchMisc`:
- `desktop-sync` / `ds`
- `rescan` / `rs`
- `setup`
- `doctor`
- `latest-branch` / `lb`
- `list-versions` / `lv`
- `list-releases` / `lr`
- `seo-write` / `sw`
- `go-mod` / `gm`
- `completion` / `cmp`
- `zip-group` / `z`
- `alias` / `a`
- `ssh`
- `prune` / `pr`
- `temp-release` / `tr`

---

### 7. `rootprojectrepos.go` — Project type queries (≈35 lines)

**Moves from root.go:**
- `dispatchProjectRepos()` (lines 360–388)

---

## Line Budget

| File | Estimated Lines | Status |
|------|----------------|--------|
| `root.go` | ~70 | ✅ Under 200 |
| `rootcore.go` | ~45 | ✅ Under 200 |
| `rootrelease.go` | ~50 | ✅ Under 200 |
| `rootutility.go` | ~55 | ✅ Under 200 |
| `rootdata.go` | ~110 | ✅ Under 200 |
| `roottooling.go` | ~80 | ✅ Under 200 |
| `rootprojectrepos.go` | ~35 | ✅ Under 200 |
| `rootusage.go` | 171 (unchanged) | ✅ Under 200 |
| `rootflags.go` | 66 (unchanged) | ✅ Under 200 |

**Total:** ~390 lines across 7 files (vs 388 in one file today).

---

## Migration Rules

1. **No behavior changes.** Pure extraction — no renames, no reordering.
2. **Package stays `cmd`.** All files in `gitmap-v27/cmd/`.
3. **Positive-logic `if` chains** preserved (no switch statements).
4. **Blank line before return** rule enforced in all moved functions.
5. **`dispatchMisc` eliminated entirely** — replaced by `dispatchData` + `dispatchTooling`.
6. **Audit tracking pattern unchanged** — each `dispatch*` returns `bool`,
   audit calls stay in the top-level `dispatch()`.

---

## Verification

1. `go build ./...` compiles.
2. `go vet ./cmd/` passes.
3. `go test ./...` passes.
4. No file in `cmd/root*.go` exceeds 200 lines.
5. `grep -c 'func dispatch' cmd/root*.go` shows exactly 7 dispatch functions.
6. `dispatchMisc` no longer exists.

---

## Acceptance Criteria

1. `root.go` ≤ 80 lines — slim entry point only.
2. Six new `root*.go` files, each ≤ 200 lines.
3. `dispatchMisc` removed; replaced by `dispatchData` + `dispatchTooling`.
4. Zero functional changes — identical CLI behavior.
5. All existing tests pass.

---

## See Also

**Same package (`cmd/`) refactors:**

- [62-refactor-seowriteloop.md](62-refactor-seowriteloop.md) — SEO write loop, git ops
- [66-refactor-zipgroupops.md](66-refactor-zipgroupops.md) — zip group CRUD and display
- [68-refactor-aliasops.md](68-refactor-aliasops.md) — alias CRUD and suggest
- [69-refactor-tempreleaseops.md](69-refactor-tempreleaseops.md) — temp release branch ops
- [72-refactor-sshgen.md](72-refactor-sshgen.md) — SSH key generation
- [73-refactor-scanprojects.md](73-refactor-scanprojects.md) — project detection
- [74-refactor-amendexec.md](74-refactor-amendexec.md) — git amend operations
- [75-refactor-status.md](75-refactor-status.md) — status display
- [76-refactor-exec.md](76-refactor-exec.md) — batch execution

**Related `release/` refactors:**
- [58-refactor-workflowfinalize.md](58-refactor-workflowfinalize.md) — pipeline, metadata, zip, GitHub upload
- [64-refactor-workflow.md](64-refactor-workflow.md) — main workflow, validation
