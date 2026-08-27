# Subtask 1: CLI Registry

1. In `gitmap/constants/constants_cli.go`:
   - Add `CmdFind = "find"`, `CmdFindRegex = "find-regex"`, etc.
   - Add `FlagLimit = "limit"`, `FlagLimitShort = "l"`.
2. In `gitmap/constants/cmd_constants_test.go`:
   - Register them in `topLevelCmds()`.
3. In `gitmap/cmd/roottooling.go`:
   - Map them to `runFind`, `runFindRegex`, `runFindRead`, `runFindReadJson`, `runFindRegexRead`, `runFindRegexReadJson`, `runFindHelp`, `runSearchHelp`, `runRegexHelp`.
4. In `gitmap/helptext/`:
   - Create dummy help markdown files for all new commands to pass `TestEveryCmdIDHasHelpFile`.
- [x] Task incomplete
