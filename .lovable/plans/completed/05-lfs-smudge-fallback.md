# LFS Smudge Fallback Implementation Plan

This plan addresses the `smudge filter lfs failed` issue encountered during `git clone` operations when LFS objects are missing on the remote server (404).

## Open Questions
- If there are *multiple* missing LFS files, the git clone stderr might only list the first one that failed. Should we rely on `git lfs ls-files -l` in the fallback to sweep all missing pointers instead of just the one mentioned in the error? (For now, I will implement the regex to parse the file from the stderr as requested in the RCA, but we can enhance it if needed).

## Proposed Changes

### `spec/` 
- **[NEW]** `spec/01-app/88-lfs-smudge-fallback.md`: Documents the LFS Smudge auto-fix architecture. (Created)

### `gitmap/clonefrom`
#### [MODIFY] `execute.go`
- Add `detectLFSSmudgeError(output string) (filename string, isSmudge bool)`
- Modify `runGitClone` to check for this error before returning.
- If detected, trigger the interactive prompt via a callback or local prompt function.
- If confirmed, execute the fallback steps (`restore`, `rm`, `Remove`, `commit`, `push`) inside the cloned destination folder with `GIT_LFS_SKIP_SMUDGE=1`.

#### [NEW] `execute_lfs_fix.go`
- Contains the `executeLFSFix(dest string, file string)` sequence of git commands to repair the broken clone.

#### [NEW] `prompt.go`
- Provide a standardized `confirmYesNo` function (borrowed from `cmd/orphans.go`) but isolated for the `clonefrom` package so it doesn't create circular dependencies.

## Verification Plan

### Automated Tests
- Create a unit test `detectLFSSmudgeError_test.go` to ensure the regex correctly parses the filename from the exact error block provided in the RCA.
- Run `go test ./gitmap/clonefrom/...` to ensure nothing is broken.

### Manual Verification
- Review the logic that constructs the Git commands (`git restore --source=HEAD :/` etc.) to guarantee they match the exact PowerShell snippet provided.
