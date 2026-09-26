# CI Issue 81: PAE Activity Git Mock, Shared Engine Sync, and Runner Test Isolation RCA

## 1. Reproduction & Symptoms
- **CI Run IDs:** [36235551944](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36235551944), [36236027437](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36236027437)
- **Failing Workflows:**
  - `CI / Go Test Matrix (ubuntu, go 1.24)` & `race-detector`
  - `Lint / Relative Paths & Nested Ifs`
- **Symptoms:**
  1. `pull_efficient_test.go: TestPartitionRecordsByActivity_WithMockData` failed because inactivity evaluation switched to local git log (`git -C repo log -1 --format=%ct`). Dummy test directories without git commits fell back to classifying repos as active.
  2. `race-detector` failed with `cmd/sync_help_test.go: undefined: isSyncHelp` after `isSyncHelp` was renamed to `isCommonHelp` without an alias adapter.
  3. `03-ai-scripts/06-cicd-local-runner.py` and `02-shared-engine.py` were truncated during prompt/skill synchronization from upstream repo, losing `chunk_items` and contract signatures.
  4. `cli/cmd/cmd_invocation_test.go` triggered live AGY pipeline error dispatches during unit test execution when failure logs existed.

## 2. Root Cause Analysis
1. **PAE Inactivity Evaluation Git Requirement**: Commit `645a80d0` shifted PAE inactivity detection from purely database-driven metrics to inspecting `git log -1`. When unit tests mocked directory paths without a genuine `.git` repository and authored commits, `git log` failed, falling back to assuming active state.
2. **Sync Help Naming Drift**: In the transition from `gitmap sync` to `gitmap common`, `isSyncHelp` was renamed to `isCommonHelp` in `cli/cmd/sync.go`, but `cli/cmd/sync_help_test.go` still invoked `isSyncHelp`.
3. **Upstream Script Synchronization Clobber**: Upstream sync scripts overwritten local runners with generic stubs, removing `chunk_items` and breaking Python CI runner tools.
4. **AGY Dispatcher Side-Effects in Unit Tests**: `runPipeline([]string{"error-logs"})` calls `dispatchErrorLogPresentation`, which triggered `cmdpipeline.PipelineAgyFixRunner` when error logs were found.

## 3. Code Fix
1. **Real Git Repos in `pull_efficient_test.go`**: Implemented `setupTestGitRepo` to initialize genuine git repositories with timestamped commits (`GIT_AUTHOR_DATE` set to 5 days ago for inactive, now for active).
2. **Compatibility Adapter in `cli/cmd/sync.go`**: Added `func isSyncHelp(arg string) bool { return isCommonHelp(arg) }`.
3. **Restored Shared Engine & Runner**: Recovered full contract implementations of `03-ai-scripts/02-shared-engine.py` and `06-cicd-local-runner.py`.
4. **Unit Test Isolation**: Wrapped pipeline error test cases in `cli/cmd/cmd_invocation_test.go` with temporary `cmdpipeline.PipelineAgyFixRunner = nil`.
5. **Coding Guidelines Adherence**: Flattened 4 nested-ifs, removed hardcoded absolute path in `cli/helptext/var.md`, and formatted all Go files with `gofmt`.

## 4. Prevention
- Always isolate external CLI/IDE dispatches during unit tests using injectable mock runners.
- When deprecating or renaming CLI utility functions, provide deprecation alias wrappers to prevent breaking test contracts.
- Ensure filesystem mocks for git-inspecting logic use isolated fixture git repositories.
