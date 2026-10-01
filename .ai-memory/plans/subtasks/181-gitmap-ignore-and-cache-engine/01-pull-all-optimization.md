# Subtask 01: Pull All Optimization

## Description
Optimize the `gitmap pull all` command and its SSH fleet counterpart (`gitmap pull all ssh` / PAS formula) to adhere to concurrency and tracking limits.

## Objectives
- Integrate `pull all` and `pull all ssh` with the `TaskQueue` structure.
- Ensure the PAS Formula concurrency limits (2 workers, 2 async ops max per SSH node) are respected.
- Record task statuses (pending, running, completed, failed) in `TaskQueue`.
- Log all SSH delegation errors to the local GitMap errors DB.

## Acceptance Criteria
- `gitmap pull all ssh` correctly queues tasks without exceeding the 2-worker limit.
- Task completion marks `Status` as 'completed' in `TaskQueue`.
- Errors are robustly aggregated and inserted into the errors DB.
