# Subtask 320-05: cmdcargo Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdcargo` (7 files)
- Leaf Packages to Create:
  - `cli/cmdcargo/enums/`
    - `cargo_subcommand_enum.go`
    - `cargo_command_enum.go`
    - `rust_toolchain_binary_enum.go`
    - `toolchain_env_var_enum.go`
    - `enums_test.go`
  - `cli/cmdcargo/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `CargoSubcommandType`: Gitmap subcommands (`status`, `st`, `info`, `install`, `in`).
2. `CargoCommandType`: Upstream cargo verbs (`build`, `check`, `test`, `clippy`, `run`, `clean`, `fmt`, `doc`, `update`, `init`, `new`, `bench`, `publish`).
3. `RustToolchainBinaryType`: Executables (`cargo`, `rustc`, `rustup`, `clang`).
4. `ToolchainEnvVarType`: Env vars (`LIBCLANG_PATH`, `PATH`, `CARGO_HOME`, `RUSTUP_HOME`, `SYSTEMROOT`, `COMSPEC`, `TMP`, `TEMP`, `USERPROFILE`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (*Type, error)`, `All*() []...Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Flags: `--install`, `-i`, `--help`, `-h`, `--yes`, `-y`, `--version`, `-V`, `--release`, `-r`, `--check`
- Binary names: `BinCargo = "cargo"`, `BinRustc = "rustc"`, `BinRustup = "rustup"`
- Path tokens: `DirCargo = ".cargo"`, `DirBin = "bin"`, `ExtExe = ".exe"`
- UI & suggestions: `StatusHeaderTitle`, `StatusNotInstalled`, `StatusInstalled`, `WingetInstallRust`, `ChocoInstallRust`, `BrewInstallRust`, `AptInstallCargo`.
- Dracula theme tokens: `ColorDraculaPurple`, `ColorDraculaWhite`, etc.
- Error codes: `ErrCodeMissingCargo = "E7100"`.

## Consuming Callers to Refactor
- `cli/cmdcargo/cargo.go`
- `cli/cmdcargo/cargo_exec.go` (also fix leaf error sovereignty)
- `cli/cmdcargo/cargo_flags.go`
- `cli/cmdcargo/cargo_resolve.go`
- `cli/cmdcargo/cargo_status.go`
- `cli/cmdcargo/cargo_suggest.go`

## Verification
- Unit tests in `cli/cmdcargo/enums/enums_test.go` and `cli/cmdcargo/constants/constants_test.go`.
- `go vet ./cli/cmdcargo/...`
- `go test ./cli/cmdcargo/...`
