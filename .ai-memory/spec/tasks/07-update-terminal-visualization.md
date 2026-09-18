# Update Terminal Visualization Enhancement

Slug: update-terminal-visualization
Status: pending
Created: 2026-08-20

## Intent

Improve the terminal visualization during the `gitmap update` process by clearly displaying the current system version, the target version before the update begins, and a post-update summary detailing the previous version, the new updated version, and the source path or slug it updated from.

## Scope

- `gitmap/cmd/update.go`: Add logic to determine current and target versions before the build phase. Print the pre-update version announcement. Print the post-update summary with versions and source path.
- `gitmap/cmd/updateremote.go` (or similar remote installer paths): Fetch target version before running the script and print the summary.
- `gitmap/constants/constants_messages.go` and `constants_update.go`: Add new formatted message constants for the visualization.

## Acceptance Criteria

1. The command `gitmap update` (or `update --source-rebuild`) prints the current version and target version before the update starts.
2. The summary at the end prints a clear visualization: "Updated from vX.Y.Z to vA.B.C (source: <path or URL>)".
3. The remote installer fallback path also prints the current and target versions and the correct summary.
4. Constants must be added strictly following naming conventions (`Msg` prefix, stored in `constants_update.go` or `constants_messages.go`).
5. All UI output must be properly spaced and styled according to the project's terminal output guidelines.

## Affected Files

- `gitmap/cmd/update.go`
- `gitmap/constants/constants_messages.go`
- `gitmap/constants/constants_update.go`
- `gitmap/cmd/updateremote.go` (if applicable)

## Linked Issues/Commands

- Captured user instruction to improve update terminal visualization.
