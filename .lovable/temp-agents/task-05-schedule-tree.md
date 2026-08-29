STATUS: DONE

# Task 05: Schedule Tree Output & Help Import/Export Execution Record

## Overview

Implemented `gitmap/cmd/schedule_tree.go` to provide UTF-8 box-drawing tree representation of scheduler tasks during interactive configuration and confirmation, and added the `Import / Export` section to `gitmap help` and compact/filter help outputs.

## Created / Modified Files

- `gitmap/cmd/schedule_tree.go`: Implemented `printScheduleSummaryTree(taskName, interval, shellType string, steps []string)`, `printScheduleTree`, `printScheduleTreeHeader`, `formatScheduleDetails`, `renderScheduleStepList`, and `renderScheduleStepNode`.
- `gitmap/cmd/schedule_cmd.go`: Updated `runScheduleAdd` to render the schedule composition tree using `printScheduleSummaryTree` before database storage; added helper `buildScheduleStepList`.
- `gitmap/cmd/rootusage_groups.go`: Added `printGroupImportExport` rendering the `Import / Export` group.
- `gitmap/cmd/rootusage.go`: Updated `printUsageReleaseAndProjectCategories` to include `printGroupImportExport`.
- `gitmap/cmd/rootusagecompact.go`: Added `HelpGroupImportExport` to `compactGroups()`.
- `gitmap/cmd/rootusagefilter_rows.go`: Added `HelpGroupImportExport` rows to `allHelpRows()`.
- `gitmap/cmd/export.go`: Added `runImportExport` dispatcher.
- `gitmap/cmd/rootdata.go`: Added `import-export` and `ie` routing entries.
- `gitmap/constants/constants_helpgroups.go`: Added `HelpGroupImportExport`, `HelpImportExport`, `HelpExportSummary`, `HelpImportSummary`, and `CompactImportExport` constants, and added `import-export` to `HelpGroupKeys`.

## Validation

- `go build ./...` passed cleanly with 0 errors.
- `go test ./cmd` passed cleanly.
- `go run . help --filter import` verified visible `Import / Export` section.
- `go run . schedule my-task --interval 1h --delay 5m` verified box-drawing schedule tree output.
- All functions adhere strictly to <= 15 lines max.
- Booleans follow positive `is`/`has` naming rules (`isScheduled`, `hasDelay`, `hasInterval`, `hasShell`, `isLastStep`).
- Semantic naming followed without generic terms.
- UTF-8 box-drawing characters (`├──`, `└──`) and ANSI terminal colors properly utilized.
