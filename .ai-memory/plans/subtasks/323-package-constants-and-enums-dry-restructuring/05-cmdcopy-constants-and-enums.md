# Subtask 323-05: cmdcopy Constants & Enums DRY Restructuring

**Status:** COMPLETED
**Completed At:** 2026-10-24T21:58:00+08:00

## Target Scope
- Target Package: `cli/cmdcopy` (3 files: `copy_paste.go`, `helpers_shim.go`, `copy_paste_test.go`)
- Leaf Packages Created:
  - `cli/cmdcopy/enums/`
    - `copy_enums.go`
    - `enums_test.go`
  - `cli/cmdcopy/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Implemented
1. `CopyActionType`: Copy action (`copy`, `paste`).
2. `CopyTargetType`: Destination target (`clipboard`, `memory`, `stdout`, `file`).
3. `CopyModeType`: Execution mode (`text`, `file`, `pipe`).
4. `CopySourceType`: Source location (`args`, `file`, `stdin`, `clipboard`, `memory`).

Each enum follows:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.
- Leaf package zero internal dependencies (standard library only).

## Specific Constants Implemented
- Command names & aliases: `CmdCopy = "copy"`, `CmdCopyAliasMem = "cp-mem"`, `CmdCopyAliasToMem = "copy-to-memory"`, `CmdPaste = "paste"`, `CmdPasteAliasMem = "paste-mem"`.
- Flags: `FlagFile = "--file"`, `FlagFileShort = "-f"`, `FlagFilepath = "--filepath"`, `FlagOut = "--out"`, `FlagOutShort = "-o"`, `FlagOutput = "--output"`, `FlagPrefixFile = "--file="`, `FlagPrefixOut = "--out="`.
- Directories, files & permissions: `DirGitmap = ".gitmap"`, `DirMemory = "memory"`, `FileClipboard = "clipboard.txt"`, `DefaultFilePerm = 0644`, `DefaultDirPerm = 0755`, `DotDirSentinel = "."`.
- Git CLI: `GitBin = "git"`, `GitRevParse = "rev-parse"`, `GitArgShowTopLevel = "--show-toplevel"`.
- ANSI colors & glyphs: `ColorGreen`, `ColorYellow`, `ColorReset`, `GlyphClipboard = "📋"`, `GlyphSuccess = "✅"`, `GlyphWarning = "⚠️ "`.
- Output messages: `MsgCopiedSuccessFmt`, `MsgPastedSuccessFmt`, `MsgEmptyPasteNoticeFmt`.
- Help & usage: `HelpCopyUsage`, `HelpCopyDesc`, `HelpPasteUsage`, `HelpPasteDesc`.
- Error codes & messages: `ErrCodeCopyEmpty`, `ErrCodeReadFileFailed`, `ErrCodeReadStdinFailed`, `ErrCodeWriteMemoryFailed`, `ErrMsgNothingToCopy`, `ErrMsgEmptyTopLevel`, `ErrFmtReadFile`, `ErrFmtReadStdin`, `ErrFmtWriteMemoryClipboard`, `ErrFmtCreateMemoryDir`, `ErrFmtWritePastedContent`.
- Exit codes: `ExitSuccess = 0`, `ExitGeneralError = 1`, `ExitUsageError = 2`.
- Helper function: `AllCopyFileFlags() []string`.

## Consuming Callers Refactored
- `cli/cmdcopy/copy_paste.go`: Decoupled completely from root `cli/constants`. Converted magic strings to package-local `constants.*` and `enums.*`. Preserved error sovereignty: returned `error` or `*apperror.AppError` directly from leaf functions. Decomposed `fetchClipboardOrMemoryContent` into `fetchClipboardOrMemoryContent` and `readMemoryFallbackContent` to enforce the <= 15 lines rule. Added `determineCopySource` and `determinePasteTarget` using enums.
- `cli/cmdcopy/helpers_shim.go`: Replaced magic strings with `constants.GitBin`, `constants.GitRevParse`, `constants.GitArgShowTopLevel`, `constants.ErrMsgEmptyTopLevel`, and `constants.AllCopyFileFlags()`.
- `cli/cmdcopy/copy_paste_test.go`: Refactored to test with `constants.*` and `enums.*`. Decomposed tests to ensure every test function is <= 15 lines.
- Caller integration: `cli/cmd/rootutility.go`: Imported `cmdcopyconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdcopy/constants"` and used `cmdcopyconstants.CmdCopy`, `cmdcopyconstants.CmdCopyAliasMem`, `cmdcopyconstants.CmdCopyAliasToMem`, `cmdcopyconstants.CmdPaste`, `cmdcopyconstants.CmdPasteAliasMem`.

## Verification
- Unit tests in `cli/cmdcopy/enums/enums_test.go` and `cli/cmdcopy/constants/constants_test.go` (100% PASS).
- Unit tests in `cli/cmdcopy/copy_paste_test.go` (100% PASS).
- `go -C cli vet ./cmdcopy/...` (PASS, 0 issues).
- `go -C cli test -count=1 -v ./cmdcopy/...` (PASS, 100% pass across all packages).
- `go -C cli vet ./cmd` (PASS, 0 issues).
