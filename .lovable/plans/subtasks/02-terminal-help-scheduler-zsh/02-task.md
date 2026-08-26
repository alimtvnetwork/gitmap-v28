# Task 2: Scheduler CLI Core & SQLite
- Create `gitmap/cmd/schedule_cmd.go` and `gitmap/store/scheduler.go`.
- Implement `gitmap schedule <name> --interval <opt> --delay <opt>`. Support interactive mode configuration.
- Implement `gitmap schedule status` to list running schedulers.
- Add SQLite schema for `scheduler_tasks`.
- Strictly adhere to boolean naming (`isScheduled`, `hasDelay`) and 15-line limit.
