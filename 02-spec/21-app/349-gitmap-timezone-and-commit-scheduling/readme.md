# 349 - GitMap Timezone and Commit Scheduling Specification

This specification suite governs the implementation of timezone-aware commit scheduling in GitMap (strictly before 08:00 AM on weekdays, anytime on weekends), dual interactive prompt and CLI one-liner flag execution, Settings UI integration, terminal help text, renewal of historical time restriction discussion files, and Cursor ecosystem mirroring.

## Documents
- `01-architecture-spec.md`: Architectural blueprint for target timezone resolution (CLI flag > config Setting Split-DB > system local), calendar-aware commit window policy, universal UTC storage invariant, and system state/sequence workflows.
- `02-component-spec.md`: Component contracts for `cli/timeutil/timezone.go`, `cli/timeutil/schedule.go`, CLI commit flags (`--tz`, `--schedule`, `--auto-window`), interactive terminal prompts, terminal help documentation, and Web Settings UI (`CommitScheduleCard.tsx`).
- `03-migration-and-verification-spec.md`: Migration and renewal of historical discussion files, Cursor ecosystem mirroring, mechanical verification matrix (V-01 through V-10), and staged GitMap atomic release ceremony.
