# Subtask 02: Cloner Concurrent & Interactive Worker Parameter Structs

Parent Plan: [143-argument-reduction-and-parameter-structs.md](../../pending/143-argument-reduction-and-parameter-structs.md)

## Goals
1. Refactor `cli/cloner/concurrent.go`:
   - Define `ConcurrentRunParams` value struct for `runConcurrent`.
   - Define `WorkerParams` struct for `startWorkers` and `cloneWorker`.
   - Define `EnqueueJobsParams` struct for `enqueueJobs`.
   - Define `CollectOutcomesParams` struct for `collectOutcomes`.
   - Replace high-arity loose argument signatures (5–6 params) with parameter structs.
2. Refactor `cli/cloner/interactive_clone.go`:
   - Define `InteractiveCloneParams` value struct for `runInteractiveClone` (6 params -> 1 struct).
   - Update call sites in `cloner.go`.

## Acceptance Criteria
- [x] All high-arity signatures (>3 parameters) in `concurrent.go` and `interactive_clone.go` encapsulated into parameter structs.
- [x] All boolean fields prefixed affirmatively (`IsSafePull`, `IsQuiet`).
- [x] `go vet ./...` in `cli/` passes with 0 errors.
