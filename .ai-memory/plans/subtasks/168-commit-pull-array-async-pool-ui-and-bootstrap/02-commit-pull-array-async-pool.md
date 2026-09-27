# Subtask 02: Implement Array Async Pool in Commit-Pull Staging & Dry-Run
Traceability ID: Task-03
Spec Reference: [02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md](../../../02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md)
Target Files: cli/cmd/commitin/workspace/clone.go, cli/cmd/commitin/orchestrator/pipeline.go, cli/cmd/commitin/orchestrator/dryrun.go
Action:
- Implement `ProbeInputsWithArrayAsyncPool(inputs []ResolvedInput, runDir string, maxWorkers int) []StagedInputResult` in `cli/cmd/commitin/workspace/clone.go`.
- Pre-allocate array of size `len(inputs)`: `slots := make([]StagedInputSlot, len(inputs))`.
- Dispatch bounded concurrent goroutines (semaphore channel `sem := make(chan struct{}, 16)`); each worker `i` evaluates input `inputs[i]`, checks cache or remote, writes directly to `slots[i]`, and sets `slots[i].isReady = true`.
- Run a ticker/cursor loop consuming from index `0` to `len(inputs)-1` in strict sequential order: as `slots[cursor].isReady` becomes true, stream status to terminal and advance `cursor++`.
- In `executePipeline` for dry-run (`ctx.Raw.IsDryRun`): utilize the pre-probed results to skip heavy sequential repo walks when dry-running.
Acceptance Criteria:
- Zero race conditions on slice writing.
- Missing repositories cleanly recorded as `isFound: false, status: "not found"` without out-of-order terminal spew.
- Strict sequential order maintained in terminal display.
Targeted Verification: python 03-ai-scripts/05-guideline-autofixer.py --path cli/cmd/commitin/workspace/clone.go
