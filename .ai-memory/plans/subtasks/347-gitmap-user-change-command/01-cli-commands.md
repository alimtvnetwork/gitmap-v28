# Subtask Plan 01: CLI Commands, Argument Parsers & Error Sovereignty

- **Subtask Slug:** `01-cli-commands`
- **Parent Plan:** `347-gitmap-user-change-command`
- **Specification Reference:** `02-spec/21-app/347-gitmap-user-change-command/01-architecture-spec.md`
- **Status:** IN_PROGRESS
- **Target Subsystem:** `cli/cmduser`, `cli/cmd`, `cli/constants`

---

## 1. Objectives & Scope

This subtask implements the user-facing CLI command surface, input parsing pipeline, error return sovereignty architecture, and unit test suite for the `change-user`, `change-login`, and `autobot` commands in GitMap:

1. **CLI Command Routing**:
   - Register top-level verbs `change-user`, `change-login`, `chuser`, `chlogin`, `autobot`, and `bot` in `cli/cmd/root.go` and `cli/cmduser/user_cmd.go`.
   - Register subcommands `gitmap user change`, `gitmap user change-login`, and `gitmap user autobot`.
2. **Zero-Regex Input Parsing Engine**:
   - Implement RFC 5322 angle bracket notation parser: `"Jane Doe <jane@example.com>"`.
   - Implement dual positional argument parser: `"Jane Doe" jane@example.com` and inverted tolerance `jane@example.com "Jane Doe"`.
   - Implement single positional argument parser: `jane@example.com` (email) or `"Jane Doe"` (name).
   - Implement explicit flag support: `--name` / `-n`, `--email` / `-e`, and `--empty-email`.
   - Implement scope resolution: `--global` / `-g`, `--local` / `-l`, `--project` / `-p`, with automatic fallback to local (inside repo) or global (outside repo).
3. **Autobot Preset Engine**:
   - Implement automated identity provisioning for `"AI automation bot"`.
   - Handle optional email address, default fallback (`local@local`), and empty/unset configuration.
4. **Error Return Sovereignty**:
   - Ensure all parsing and execution leaf functions return `error` or typed structs/results directly.
   - Enforce total ban on `cliexit.HandleError` and `os.Exit` inside leaf functions so errors bubble cleanly up to `root.go` for database audit capture.
5. **Unit Test Suite**:
   - Author comprehensive unit tests in `cli/cmduser/user_change_test.go` verifying all parsing variations, edge cases, scope permutations, and error boundaries.

---

## 2. Target Files & Component Map

| Target File Path | Purpose |
|:---|:---|
| `cli/constants/constants_commands_extra.go` | Declare command string constants (`CmdChangeUser`, `CmdChangeLogin`, `CmdAutobot`, aliases) |
| `cli/cmduser/user_change_types.go` | Declare domain models (`UserChangeOptions`, `ParsedIdentity`, `UserChangeScopeType`) |
| `cli/cmduser/user_change_parser.go` | Implement zero-regex <= 15 line parsing helpers and scope resolvers |
| `cli/cmduser/user_change_cmd.go` | Implement command runners `RunChangeUser` and `RunAutobot` |
| `cli/cmduser/user_cmd.go` | Extend `DispatchUser` and `dispatchUserSubcommand` route tables |
| `cli/cmduser/user_change_ui.go` | ANSI status card and tabular change confirmation rendering |
| `cli/cmduser/user_change_test.go` | Comprehensive unit tests for argument parsing and identity application |

---

## 3. Atomic Execution Checklist

### Phase 1: Constants & Domain Type Modeling
- [ ] **1.1**: Register command constants in `cli/constants/constants_commands_extra.go`:
  - `CmdChangeUser = "change-user"`
  - `CmdChangeUserAlias = "chuser"`
  - `CmdChangeLogin = "change-login"`
  - `CmdChangeLoginAlias = "chlogin"`
  - `CmdAutobot = "autobot"`
  - `CmdAutobotAlias = "bot"`
- [ ] **1.2**: Create `cli/cmduser/user_change_types.go` declaring:
  - Enum `UserChangeScopeType` with values `ScopeLocal`, `ScopeGlobal`, `ScopeProject`.
  - Struct `ParsedIdentity` with fields: `Name`, `Email`, `HasName`, `HasEmail`, `IsAngleBracket`.
  - Struct `UserChangeOptions` with fields: `Name`, `Email`, `Scope`, `IsGlobal`, `IsLocal`, `IsProject`, `IsAutobot`, `IsEmptyEmail`, `HasExplicitEmail`, `HasExplicitName`.

### Phase 2: Zero-Regex Input Parsing Engine (`cli/cmduser/user_change_parser.go`)
- [ ] **2.1**: Implement `stripWrappingQuotes(s string) string` removing outer single or double quotes without regex.
- [ ] **2.2**: Implement `tryParseAngleBracket(raw string) (ParsedIdentity, bool)` using `strings.IndexByte(raw, '<')` and `strings.LastIndexByte(raw, '>')`.
- [ ] **2.3**: Implement `parsePositionalPair(first, second string) ParsedIdentity` with intelligent `@` symbol detection.
- [ ] **2.4**: Implement `parseSinglePositional(arg string) ParsedIdentity` mapping email if `@` is present, name otherwise.
- [ ] **2.5**: Implement `extractScopeFlags(args []string) (ScopeFlags, []string)` parsing `--global`, `-g`, `--local`, `-l`, `--project`, `-p`.
- [ ] **2.6**: Implement `resolveTargetScope(flags ScopeFlags, isInsideRepo bool) UserChangeScopeType`.
- [ ] **2.7**: Implement master parsing entrypoint `ParseUserChangeArgs(args []string, isInsideRepo bool) (UserChangeOptions, error)`:
  - <= 15 lines max.
  - Returns strongly-typed error if identity cannot be determined.

### Phase 3: Autobot Feature Engine
- [ ] **3.1**: Implement `ParseAutobotArgs(args []string, isInsideRepo bool) (UserChangeOptions, error)`:
  - Default Name: `"AI automation bot"`.
  - Optional email parsing from positional arg or `--email` flag.
  - Default email fallback: `"local@local"` when unprovided.
  - Support `--empty-email` flag to leave `user.email` unset.
- [ ] **3.2**: Implement `buildAutobotOptions(email string, scope UserChangeScopeType, isEmpty bool) UserChangeOptions`.

### Phase 4: Command Handlers & CLI Routing
- [ ] **4.1**: Create `cli/cmduser/user_change_cmd.go`:
  - Implement `RunChangeUser(args []string) error`.
  - Implement `RunAutobot(args []string) error`.
  - Implement `applyGitIdentity(opts UserChangeOptions) error`.
- [ ] **4.2**: Update `cli/cmduser/user_cmd.go`:
  - Extend `DispatchUser(command string) (bool, error)` to route `change-user`, `change-login`, `chuser`, `chlogin`, `autobot`, `bot`.
  - Extend `dispatchUserSubcommand(sub string, args []string) error` to route `change`, `change-login`, and `autobot`.
- [ ] **4.3**: Create `cli/cmduser/user_change_ui.go`:
  - Render ANSI confirmation card detailing updated Name, Email, and Scope.

### Phase 5: Error Return Sovereignty & DB Audit Verification
- [ ] **5.1**: Audit all newly added functions in `cli/cmduser/` to ensure zero calls to `cliexit.HandleError` or `os.Exit`.
- [ ] **5.2**: Verify errors return directly as `error` or `*apperror.AppError` so `root.go` captures them in `finishCommandAudit`.
- [ ] **5.3**: Verify non-zero exit codes bubble cleanly to top-level command router without intermediate silent failure.

### Phase 6: Unit Test Suite & Code Hygiene
- [ ] **6.1**: Create `cli/cmduser/user_change_test.go` covering:
  - Angle bracket parsing: `"Jane Doe <jane@example.com>"`, `Jane <jane@dev>`, unclosed bracket syntax.
  - Dual positional: `"Jane Doe" jane@example.com`, `jane@example.com "Jane Doe"`.
  - Single positional email and name variations.
  - Explicit flags: `--name`, `--email`, `-n`, `-e`, `--empty-email`.
  - Scope resolutions: `--global`, `--local`, `--project`, default repository inference.
  - Autobot command with no args, custom email, `--empty-email`, and `--global`.
  - Error boundaries on invalid/missing arguments.
- [ ] **6.2**: Code hygiene validation:
  - Verify every function <= 15 lines.
  - Verify zero nested `if` statements.
  - Verify all booleans prefixed with `is*` or `has*`.
  - Verify zero regex package imports.

---

## 4. Coding Guidelines & Invariant Rules

| Invariant | Standard | Enforcement Rule |
|:---|:---|:---|
| **Relative Git Paths** | Always relative (e.g. `cli/cmduser/user_change_cmd.go`) | TOTAL BAN on absolute filesystem paths (`/abs/path`) |
| **Boolean Naming** | Only `is*` or `has*` | NO `can*`, `should*`, `was*`, `will*`, NO negative booleans |
| **Function Length** | <= 15 lines max, <= 8 lines preferred | Guard clauses, decompose into helper functions |
| **Nested Conditionals** | 0 nested ifs | Early returns and inverted condition checks |
| **Error Sovereignty** | Leaf functions return `error` | BANNED: calling `cliexit.HandleError` in leaf functions |
| **Exit Codes** | Strongly-typed enums only | BANNED: magic integer exit codes (`os.Exit(1)`) |
| **Regex Policy** | Zero regex | BANNED: `regexp.MustCompile`, `lazyregex` in user change parser |
| **Spacing Style** | Vertical whitespace discipline | Blank line after `}`, blank line before `return` |
