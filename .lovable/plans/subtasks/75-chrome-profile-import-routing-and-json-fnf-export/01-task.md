# Subtask 01: Profile Routing & Constants Update

## Objective
Update `gitmap/constants/constants_profile.go` and `gitmap/cmd/profile.go` to route `gitmap profile` commands for `import`, `import-all`, `export`, `export-all`, `inspect`, `preview`, and `check` directly to Chrome profile handlers.

## Affected Files
- `gitmap/constants/constants_profile.go`
- `gitmap/cmd/profile.go`

## Requirements
1. In `gitmap/constants/constants_profile.go`:
   - Add constants with `// gitmap:cmd skip`:
     - `CmdProfileImport = "import"`
     - `CmdProfileImportAll = "import-all"`
     - `CmdProfileExport = "export"`
     - `CmdProfileExportAll = "export-all"`
     - `CmdProfileInspect = "inspect"`
     - `CmdProfilePreview = "preview"`
     - `CmdProfileCheck = "check"`
     - `CmdProfileImportCheck = "import-check"`
   - Update `HelpProfile` and `ErrProfileUsage`.
2. In `gitmap/cmd/profile.go`:
   - Refactor `routeProfileSub` to delegate to modular sub-routers (`routeGitProfileSub`, `routeDBProfileSub`, `routeChromeProfileSub`).
   - Obey $\le 15$ lines per function limit.
   - Return `*AppError` on unknown subcommands instead of calling `cliexit.HandleError`.
   - Propagate errors without swallowing.
