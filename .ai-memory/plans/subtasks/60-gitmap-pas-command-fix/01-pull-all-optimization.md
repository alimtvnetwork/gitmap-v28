# Subtask 01: GitMap Pull All Decoupling & PAS Formula Engine

## Objectives
- Ensure `gitmap pull all` (`pa`) prioritizes raw git pull and executes ignore inspections asynchronously (5 repos per worker group).
- Ensure `gitmap pull all ssh` (`pas`) enforces the GitMap PAS Formula: max 2 workers, 2 async operations per node (throttled to 1 worker under high CPU pressure).
- Enqueue all operations into `TaskQueue` table with `queue_id`, `status` ('pending', 'running', 'completed', 'failed'), timestamps.
- Log delegation failures to GitMap Errors DB.

## Target Files
- `cli/cmdpull/pull.go`
- `cli/cmd/rootcore.go`
- `cli/cmd/rootgit.go`
