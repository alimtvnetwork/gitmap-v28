# Subtask 321-01: cmdconfig Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdconfig` (12 files)
- Leaf Packages to Create:
  - `cli/cmdconfig/enums/`
    - `config_enums.go`
    - `enums_test.go`
  - `cli/cmdconfig/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `ConfigActionType`: Actions (`export`, `import`, `what-configs`, `which-format`, `get`, `set`, `list`, `reset`).
2. `ConfigOutputFormatType`: Output formats (`table`, `json`, `yaml`, `compact`).
3. `ConfigScopeType`: Configuration scope (`local`, `global`, `system`, `project`).
4. `ConfigKeyType`: Known configuration keys (`target_dir`, `default_tool`, `export_dir`, `schema_version`, `encoding`, `auto_import`).
5. `ConfigToolType`: Supported tools (`vscode`, `qtorrent`, `utorrent`, `all`).
6. `ConfigFileFormatType`: Detected file formats (`json`, `yaml`, `db`, `sqlite`, `toml`, `env`, `ini`, `cfg`, `conf`).
7. `ConfigFileEncodingType`: File payload encoding (`utf-8`, `base64`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Command names & aliases: `CmdExportConfig`, `CmdImportConfig`, `CmdWhatConfigs`, `CmdWhichFormat`.
- CLI flags & prefixes: `FlagArgDir`, `FlagArgYes`, `FlagArgHelp`, `FlagArgDirPrefix`.
- OS identifiers & environment tokens: `OSWindows`, `OSDarwin`, `OSLinux`, `EnvAppData`.
- Tool directory subpaths, canonical file names, table columns, starter templates, error codes.
- Zero magic strings.

## Consuming Callers to Refactor
- `cli/cmdconfig/config_model.go`
- `cli/cmdconfig/config_tool_paths.go`
- `cli/cmdconfig/config_export.go`
- `cli/cmdconfig/config_export_writers.go`
- `cli/cmdconfig/config_import.go`
- `cli/cmdconfig/config_import_readers.go`
- `cli/cmdconfig/what_configs_cmd.go`
- `cli/cmdconfig/which_format_cmd.go`
- Test files in `cli/cmdconfig/`

## Verification
- Unit tests in `cli/cmdconfig/enums/enums_test.go` and `cli/cmdconfig/constants/constants_test.go`.
- `go vet ./cli/cmdconfig/...`
- `go test ./cli/cmdconfig/...`
