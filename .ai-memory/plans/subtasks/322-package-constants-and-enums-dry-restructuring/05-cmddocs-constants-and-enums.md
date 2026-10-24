# Subtask 322-05: cmddocs Constants & Enums DRY Restructuring - [COMPLETED]

## Status
- **COMPLETED**: 2026-10-24
- Decoupled `cli/cmddocs` from root `cli/constants`.
- Created leaf packages `cli/cmddocs/enums/` and `cli/cmddocs/constants/`.
- Implemented comprehensive unit tests in `cli/cmddocs/docs_test.go`, `constants_test.go`, and `enums_test.go`.
- Verified 100% green with `go vet` and `go test`.

## Target Scope
- Target Package: `cli/cmddocs` (2 files)
- Leaf Packages to Create:
  - `cli/cmddocs/enums/`
    - `docs_enums.go`
    - `enums_test.go`
  - `cli/cmddocs/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `DocsCategoryType`: Documentation domains (`cli`, `architecture`, `spec`, `prompt`, `guide`, `api`, `all`).
2. `DocsViewerModeType`: Viewer display targets (`browser`, `terminal`, `server`, `web`).
3. `DocsFormatType`: Presentation formats (`markdown`, `html`, `text`, `json`).

Each enum follows:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Command names & aliases: `CmdDocs = "docs"`, `CmdDocsAlias = "d"`.
- URLs & endpoints: `DocsURL = "https://gitmap.dev/docs"`, `DefaultServerHost = "127.0.0.1"`, `DefaultServerPort = 8080`.
- Flags: `--browser`, `--terminal`, `--server`, `--category`, `--format`, `--port`, `--help`.
- OS identifiers & execution tokens: `OSWindows`, `OSDarwin`, `OSLinux`, `CmdWindowsShell`, `CmdArgSlashC`, `CmdArgStart`, `CmdOpen`, `CmdXdgOpen`.
- Output messages & errors: `MsgDocsOpened`, `ErrDocsOpen`, `ErrDocsOpenFailed`.
- Error codes: `ErrCodeDocsLaunchFailed`, `ErrCodeDocsInvalidArgs`, `ErrCodeDocsServerFailed`.

## Consuming Callers Refactored
- `cli/cmddocs/docs.go`: Decoupled completely from root `cli/constants`. Extracted `resolveBrowserCmd(goos, url)` helper and `launchBrowser(url)`. Functions adhere to <= 15 line limit, return styling, and zero nested ifs.
- `cli/cmddocs/hook_shim.go`: Verified compliant with zero root dependencies and <= 15 line functions.
- `cli/cmddocs/docs_test.go`: Added test coverage for `ResolveViewerMode`, OS branches in `resolveBrowserCmd`, `RunDocs` success/failure, and args handling with mock `startCmdFn` and `CheckHelpFn`.

## Verification
- Unit tests: `go test -v ./cmddocs/...` PASS (100% pass across all packages)
- Lint check: `go vet ./cmddocs/...` PASS (0 warnings)
