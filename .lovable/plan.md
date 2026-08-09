## Pipeline Extend V2: AI Release Agent Instructions
All AI agents triggering releases MUST follow `spec/09-pipeline-extend-v2/readme.md` strictly. 
**Rule**: Tag-based releases fail if the source constants aren't updated FIRST.
**Execution**:
1. Run PowerShell sweep to bump versions across all canonical files.
2. Push branch to `main`.
3. Finally push `git tag vX.Y.Z`.

## Problem (verified)

`gitmap clone <url>` (single URL) already keeps the folder name verbatim, including a trailing `-vN` (`gitmap/cmd/clone.go:247-259`).

The multi-URL path does not. `resolveCloneFolder` in `gitmap/cmd/clonemulti.go` calls `clonenext.ParseRepoName` and, when the repo name has a version suffix, returns `parsed.BaseName` — so `codex-june-6-v2` lands in `codex-june-6/`. It is used in two places:

- `gitmap/cmd/clonemulti.go` — `executeDirectCloneOne` (the actual clone destination)
- `gitmap/cmd/clone.go:174` — the VS Code Project Manager pair, mirroring the same resolution

## Fix

1. `gitmap/cmd/clonemulti.go`
   - Rewrite `resolveCloneFolder` to return `folderName` when provided, else `repoName` verbatim. No version parsing, no flattening.
   - Drop the now-unused `clonenext` import.
   - Update the doc comments on `resolveCloneFolder` and `executeDirectCloneOne` (both currently say "versioned URLs flatten via clonenext") to state that base clone never rewrites the folder name; version bumping stays in `clone-next`.
2. `gitmap/cmd/clone.go:169-177` — comment only; behavior follows automatically since it calls the same helper.
3. Tests — add a table test (new `gitmap/cmd/clonemulti_folder_test.go`) covering: plain repo, `-v1`/`-v13` suffix preserved, explicit folder name wins, `.git` and trailing-slash URLs. Also grep and fix any existing test asserting the flattened behavior.

Nothing in `clone-next`, `clone-now`, `reclone`, or `clonefixrepo` is touched — version bumping there stays as is.

## Spec + memory

- `spec/01-app/104-clone-multi.md` — add a "Folder naming" rule: folder = repo name from URL verbatim (`-vN` preserved); explicit second positional arg overrides; no flattening in `clone`.
- `.lovable/memory/` — record the constraint "base `clone` must never strip `-vN` from folder names" so it is not reintroduced.

## Release v6.89.0

- `gitmap/constants/constants.go` → `Version = "6.89.0"`
- `src/constants/index.ts` → `6.89.0`
- `.gitmap/release/v6.89.0.json` (copy the v6.85.0 shape) and update `.gitmap/release/latest.json`
- `changelog.md` → new `v6.89.0` section: maintenance bump
- `readme.md` → repin every `v6.85.0` occurrence to `v6.89.0`

## Verification

- `go build ./... && go vet ./...` and `go test ./gitmap/cmd/... -run Clone -count=1`
- `bunx vitest run src/test/version-sync.test.ts`
- `.github/scripts/check-changelog-version-sync.sh`
