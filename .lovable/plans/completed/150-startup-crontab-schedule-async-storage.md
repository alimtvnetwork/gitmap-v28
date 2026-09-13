# Master Spec: 150-startup-crontab-schedule-async-storage [COMPLETED]

## 1. Overview & Goal

Implement the cross-platform unified automation subsystem in GitMap:
1. **Startup Management (`gitmap startup`)**:
   - Top-level command and macro integration (`macro startup add <name>`).
   - Bidirectional integration: startup can run macros, macros can register startup jobs.
   - Run frequencies: `everytime` (every OS login/boot), `once-a-week`.
   - Supports: PowerShell scripts (`.ps1`), Bash scripts (`.sh`), binaries, icon bindings (`--icon`), and macros.
   - Subcommands: `gitmap startup ls`, `gitmap startup add <target>`, `gitmap startup rm <name>`, `gitmap startup run <name>`.
   - Dedicated split DB: `<BinaryDataDir>/startup.db` with PascalCase `StartupItem` and `StartupLog` tables, registered in `SplitDatabaseRegistry`.
   - OS-level native integration: Windows registry / Startup folder shortcuts (`cli/startup/`), Linux XDG desktop entries, macOS LaunchAgents.

2. **Crontab & Schedule Subsystem (`gitmap schedule`, alias `gitmap crontab`)**:
   - Unify `gitmap schedule` and `gitmap crontab` under split DB architecture.
   - Formats: PowerShell (`.ps1`), Bash (`.sh`), shell, JavaScript (`.js`), and macros.
   - Subcommands: `ls`, `add`, `edit`, `rm`, `debug`, `info`, `run`.
   - Per-schedule split database: `<BinaryDataDir>/schedules/<slug>.db` with `ScheduleConfig` and `ScheduleLog`.
   - Debug and info readout showing database file location, execution timestamps, duration, and exit codes.

3. **Interactive Macro Execution (`cli/macro/`, `cli/cmdmacro/`)**:
   - Real-time stdout/stderr streaming during macro step execution.
   - Instant visual terminal feedback when executing commands.

4. **Async Operations & Monitoring (`gitmap async`)**:
   - Detached asynchronous command / monitoring daemon with `-t <seconds>` periodic interval.
   - Monitor tasks, storage, or health.
   - Subcommands: `gitmap async <cmd> [-t <seconds>]`, `gitmap async ls`, `gitmap async stop <id>`.

5. **Storage Calculation & Inspection (`gitmap storage`)**:
   - `gitmap storage`: Inspects current drive/volume (Windows D: / C:, Linux `/`), file system format, total, used, free space, and storage ratio.
   - `gitmap storage ls`: Discovers and lists all Gitmap SQLite databases (root DB, profile DBs, split child DBs) with sizes and record counts.
   - Storage partitioning/swap/LVM: Diagnostic baseline and safe architecture with interactive confirmation safeguards to prevent data loss.

6. **Pipeline Error ETA Polling (`gitmap pipeline errors -t`)**:
   - Reads `.lovable/temp/runner-eta.json` or computes historical duration approximations.
   - Polls dynamically until ETA completes, then emits structured/compact error logs.

7. **Documentation & Help Parity**:
   - Terminal help screen (`gitmap help`, `gitmap help --filter`), LLM documentation (`gitmap llm-docs`), and root `readme.md`.

---

## 2. Completed Subtasks

- `01-split-db-startup-schema.md`: Created `StartupSplitDB` and database models in `cli/store/startup_split_db.go` and schema in `cli/constants/constants_split_db_sql.go`. Updated ERD in `spec/21-app/gitmap-database-erd.mmd`.
- `02-startup-cli-command.md`: Implemented `gitmap startup` command suite (`ls`, `add`, `rm`, `run`) with native OS hooks, weekly deduplication, and execution logging in `cli/cmd/startup_cmd.go` and `startup_run.go`.
- `03-crontab-schedule-unification.md`: Unified `gitmap crontab` alias and per-schedule split DB logging, debugging (`gitmap schedule debug <name>`), and editing in `cli/cmdschedule/schedule_debug.go` and `cli/cmdschedule/schedule_cmd.go`.
- `04-macro-interactive-streaming.md`: Added bidirectional startup-macro registration in `cli/cmdmacro/macro_startup.go` and verified interactive streaming in `cli/macro/execute.go`.
- `05-async-command-runner.md`: Implemented `gitmap async` command with `-t` interval loop, process spawning, and lifecycle tracking in `cli/cmd/async_cmd.go` and `cli/cmd/async_exec.go`.
- `06-storage-inspection.md`: Implemented `gitmap storage` drive metrics calculation (`cli/cmd/storage_drive_windows.go`, `storage_drive_other.go`) and `gitmap storage ls` database inventory (`cli/store/storage_inventory.go`, `cli/cmd/storage_display.go`).
- `07-pipeline-eta-flag.md`: Enhanced `gitmap pipeline errors -t` with runner countdown polling loop in `cli/cmdpipeline/pipeline_runner_eta.go`.
- `08-help-llmdocs-and-readme.md`: Registered commands in CLI constants, help groups, AST parity tests, LLM docs, root dispatchers, and root `readme.md`.
