# 11 Endpoint Resolver DB Upgrade

## Parent Task Goal

Enhance the global endpoint resolution logic in Gitmap so that commands (especially `commit-right`, `commit-left`, `commit-both`, `open`, `inject`, etc.) can seamlessly resolve repositories using five modalities:
1. ID of the repo (from the SQLite DB)
2. Alias (from the SQLite DB)
3. Relative path
4. Absolute path
5. URL (HTTPS, SSH, git@)

## Architectural Plan

1. **Core Resolution Logic:**
   - Create `gitmap/cmd/resolver.go`.
   - Implement `func ResolveEndpointString(raw string) string` which orchestrates the resolution.
   - It will check if `raw` is a URL. If so, return unchanged.
   - It will check if `raw` is an existing path (absolute or relative resolving to absolute). If so, return it.
   - It will query the SQLite database (`store.DB`) to resolve by ID, Alias, or Slug. If found, return the resolved `AbsolutePath`.
   - If nothing is found, return `raw` and let downstream functions (like `movemerge`) throw errors.

2. **Integration:**
   - Modify `gitmap/cmd/committransfer.go` `resolveCommitEndpoints` to intercept `leftRaw` and `rightRaw` using `ResolveEndpointString` before passing to `movemerge.ResolveEndpoint`.
   - Ensure the new logic relies on `openDB()` to fetch the SQLite database.

3. **Subtasks Execution:**
   - **Subtask 1: DB & SQL Constants Update**
     - Add `SQLSelectRepoByID` to `gitmap/constants/constants_store.go`.
     - Implement `FindByID` in `gitmap/store/repo.go`.
   - **Subtask 2: Implement Core Resolver**
     - Write `gitmap/cmd/resolver.go`.
     - Include logic to parse numeric IDs, query aliases, and fallback to slugs.
   - **Subtask 3: Integration & Testing**
     - Update `gitmap/cmd/committransfer.go`.
     - Run `go test ./...` and ensure functionality works.
     - Group commits and push.
   - **Subtask 4: Documentation**
     - Update memory and specifications.

## Code Review Guide

- Do not swallow errors; use `apperror`.
- Ensure no cyclic dependencies between `cmd`, `store`, and `movemerge`.
- Boolean variables must be prefixed with `is`, `has`, `should`, or `can`.
- No generic variable names (`temp`, `data`, `res`).
