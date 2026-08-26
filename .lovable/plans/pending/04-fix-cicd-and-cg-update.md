# Plan: Fix CI/CD checkout issue and Enhance `gitmap cg update` UI

## 1. Fix CI/CD Workflow (`.github/workflows/ci.yml`)
- **Root Cause**: The `.github/actions/policy-check` is a local action. Several jobs in `ci.yml` call this action as their first step without checking out the repository first, resulting in `Can't find 'action.yml'`.
- **Solution**: Inject `- uses: actions/checkout@v6` before the `uses: ./.github/actions/policy-check` step in the following jobs:
  - `cmd/ Naming Check`
  - `Legacy Refs Check`
  - `Deploy Layout Check`
  - `constants/ Naming Check`
  - `constants/ Collision Check`
  - `GITMAP_ALLOW_GOLDEN_UPDATE Leak Check`

## 2. Enhance `gitmap cg update` CLI UI
- **Objective**: The user wants a detailed, colorful summary of the coding guidelines update process, showing version transitions (old vs new), files updated, and status for each repository.
- **Solution in `gitmap/cmd/cg_worker.go`**:
  - Update `executeCGWorkers` to collect results instead of just printing "Done".
  - In `runCgWorker`, before running the script, read the current CG version using `ReadCGMetadata(repo)`.
  - Capture standard output and error of the `cmd.Run()`.
  - After running the script, read the CG version again using `ReadCGMetadata(repo)`.
  - Pass the result back to `executeCGWorkers` via a channel.
  - After all workers finish, print a nicely formatted summary block with Lipgloss:
    - Display the target repository.
    - Display `Previous Version -> New Version`.
    - Indicate what files were updated (e.g., `.lovable/coding-guidelines/` and `version.json`).
    - Use bright colors as requested.

## 3. Strict Compliance Checks
- Follow the 15-line function limit for Go.
- Use explicit boolean naming (e.g., `isSuccess`, `hasChanged`).
- Wrap errors using standard error variables or `fmt.Errorf`.
- Output summary list explicitly at the end of the run before bumping the release.
