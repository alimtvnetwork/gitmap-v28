# Subtask 321-03: cmdsearch Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdsearch` (5 files)
- Leaf Packages to Create:
  - `cli/cmdsearch/enums/`
    - `search_enums.go`
    - `enums_test.go`
  - `cli/cmdsearch/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `SearchScopeType`: Search domain (`content`, `symbol`, `path`, `commit`, `repo`).
2. `SearchMatchModeType`: Match engines (`keyword`, `regex`, `exact`, `fuzzy`).
3. `SearchOutputFormatType`: Presentation format (`text`, `json`, `table`, `compact`).
4. `SearchCategoryType`: Query origins (`user`, `ai`).
5. `SearchActionType`: Operational actions (`query`, `history`, `top`, `stats`, `dh2d`, `clean`, `clear`, `purge`, `reset`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Command names & subcommands: `CmdSearch = "search"`, `CmdRepoSearch`, `CmdRepoRegex`, `CmdRepoSearchJSON`, `CmdReplaceRegex`, `CmdSearchReplaceAll`.
- CLI flags & prefixes: `FlagArgQuery`, `FlagArgPattern`, `FlagArgDir`, `FlagArgLimit`, `FlagArgAI`, `FlagArgJSON`, etc.
- Default limit `DefaultLocalSearchLimit = 20`, ignored directories (`node_modules`, `vendor`, `bin`), output paths (`.gitmap/output/repo_search`), empty JSON array `"[]"`, format strings, error messages.

## Consuming Callers to Refactor
- `cli/cmdsearch/search.go`
- `cli/cmdsearch/search_entry.go`
- `cli/cmdsearch/search_help_menu.go`
- `cli/cmdsearch/helpers_shim.go`

## Verification
- Unit tests in `cli/cmdsearch/enums/enums_test.go` and `cli/cmdsearch/constants/constants_test.go`.
- `go vet ./cli/cmdsearch/...`
- `go test ./cli/cmdsearch/...`
