# Subtask Plan 03: Workgroup CLI Commands, Commit Formatter & Documentation

- **Subtask Slug:** `subtask-03-cli-commands-and-commit-format`
- **Parent Plan:** `317-system-for-folder-specific-work-directories-and-ssh-key-automation`
- **Specification Reference:** `02-spec/21-app/317-system-for-folder-specific-work-directories-and-ssh-key-automation/03-cli-commands-and-commit-format-spec.md`
- **Status:** Completed
- **Target Repository:** `gitmap-v28`

---

## 1. Objectives & Scope

This subtask delivers the user-facing CLI command surface, commit formatting engine, interactive terminal UI menus, markdown help documentation, and AI agent self-training protocols for the folder-specific work directory and SSH automation subsystem:

1. **Workgroup CLI Command Suite**:
   - Implement `gitmap workgroup`, `gitmap wg`, and `gitmap group` commands with subcommands: `add`, `bind`, `unbind`, `ls`, `show`, `test-ssh`, and `rm`.
   - Implement hierarchical ancestor lineage resolution so child repositories automatically inherit parent folder workgroups and bound SSH keys.
2. **Enhanced `cpf` and `cbf` Commit Message Parser & Formatter**:
   - Implement `[module]: <item>` and `type([module]): <item>` parser supporting `feat`, `fix`, `chore`, and `release`.
   - Register `cbf` as a first-class alias for `commit-push-bug` alongside `cpf`.
   - Support duplicate prefix trimming, outer quote sanitization, and backwards compatibility with legacy hyphen-delimited messages.
3. **UI Help Menus & HD Markdown Help Cards**:
   - Render Lipgloss-themed terminal cards in `cli/termout/` for workgroup commands and enhanced commit formatting.
   - Author HD markdown help cards (`workgroup.md`, `commit-format.md`, `cbf.md`) and register them in `cli/helpdoc/catalog.go` and `cli/helpdoc/aliases.go`.
4. **LLM Self-Training Documentation & Skills Sync**:
   - Update the 6-phase curriculum in `cli/cmd/llm/llm_train.go` and `llm_skill.go`.
   - Update `.agents/skills/gitmap/skill.md` with authorized bracketed commit syntax.
   - Synchronize skills across all connected fleet repositories using `03-ai-scripts/38-sync-prompts-skills-scripts.py`.

---

## 2. Target Files & Component Map

| Target File Path | Purpose |
|:---|:---|
| `cli/constants/constants_commands.go` | Declare `CmdWorkgroup`, `CmdWorkgroupAlias`, `CmdGroupAlias`, `CmdCommitPushBugFixAlias` |
| `cli/cmd/roottooling_util.go` | Register command dispatch table entries for `workgroup`, `wg`, `group`, and `cbf` |
| `cli/cmdworkgroup/workgroup_types.go` | Declare data models (`WorkgroupRecord`, `FolderBindingRecord`, `WorkgroupOptions`) |
| `cli/cmdworkgroup/workgroup_flags.go` | Flag parsing for `--ssh-key`, `--desc`, `--user-name`, `--user-email`, `--default` |
| `cli/cmdworkgroup/workgroup_resolver.go`| Hierarchical ancestor directory matching and lineage key lookup |
| `cli/cmdworkgroup/workgroup_test_ssh.go` | Remote SSH handshake execution and host verification |
| `cli/cmdworkgroup/workgroup_ui.go` | Lipgloss tabular view and detail card rendering (`cli/termout`) |
| `cli/cmdworkgroup/workgroup_cmd.go` | Command entrypoint and subcommand dispatcher |
| `cli/cmdworkgroup/workgroup_test.go` | Unit tests for workgroup commands, flag parsing, and lineage matching |
| `cli/cmdcommit/commit_message_parser.go` | Modular parser for `[module]: <item>` and `type([module]): <item>` |
| `cli/cmdcommit/commit_message_parser_test.go` | Comprehensive unit tests for commit message parsing and normalization |
| `cli/cmdcommit/commit_push.go` | Integrate commit message parser and add `RunCommitPushBugFix` / `cbf` handler |
| `cli/cmdcommit/commit_help_menu.go` | Update terminal help card with bracketed module syntax examples |
| `cli/helpdoc/workgroup.md` | HD markdown help card for `workgroup` |
| `cli/helpdoc/commit-format.md` | HD markdown help card for `[module]: <item>` commit formatting |
| `cli/helpdoc/cbf.md` | HD markdown help card for `cbf` command |
| `cli/helpdoc/catalog.go` | Register documentation summaries in static help catalog |
| `cli/helpdoc/aliases.go` | Register documentation aliases (`wg`, `group`, `cbf`) |
| `cli/cmd/llm/llm_train.go` | Update 6-phase curriculum and JSON catalog |
| `cli/cmd/llm/llm_skill.go` | Update embedded Antigravity skill template |
| `.agents/skills/gitmap/skill.md` | Update workspace Antigravity skill with bracketed syntax rules |

---

## 3. Atomic Execution Checklist

### Phase 1: Constants & Command Dispatch Routing
- [x] **1.1**: Register command constants in `cli/constants/constants_commands.go`:
  - `CmdWorkgroup = "workgroup"`
  - `CmdWorkgroupAlias = "wg"`
  - `CmdGroupAlias = "group"`
  - `CmdCommitPushBugFixAlias = "cbf"`
- [x] **1.2**: Register route mappings in `cli/cmd/roottooling_util.go` to route `workgroup`, `wg`, and `group` to `cmdworkgroup.RunWorkgroup(argsTail())`.
- [x] **1.3**: Register route mapping in `cli/cmd/roottooling_util.go` for `cbf` to `cmdcommit.RunCommitPushBug(argsTail())`.
- [x] **1.4**: Add `cbf` to semantic verb lists in `cli/cmd/cmd_sends.go` (`isValidSendsVerb` and `resolveSemanticCommitPrefix`).

### Phase 2: Enhanced Commit Message Parser & Formatter (`cli/cmdcommit`)
- [x] **2.1**: Create `cli/cmdcommit/commit_message_parser.go` implementing:
  - Struct `ParsedCommitMessage` with fields: `RawInput`, `FormattedMessage`, `CommitType`, `Module`, `ItemDescription`, `HasModule`, `HasBracketSyntax`, `IsConventional`.
  - Function `ParseAndFormatCommitMessage(verb string, rawArgs []string) ParsedCommitMessage`.
  - Helper `stripWrappingQuotes(s string) string` removing outer quotes.
  - Helper `tryParseBracketSyntax(verb, s string) (ParsedCommitMessage, bool)` parsing `[module]: <item>`.
  - Helper `tryParseConventionalSyntax(verb, s string) (ParsedCommitMessage, bool)` parsing `type([module]): <item>` and `type(module): <item>`.
  - Helper `tryParseLegacyHyphenSyntax(verb, s string) (ParsedCommitMessage, bool)` normalizing `module - item` to `[module]: item`.
  - Helper `trimDuplicatePrefix(verb, s string) string` stripping redundant leading prefixes.
- [x] **2.2**: Integrate parser into `cli/cmdcommit/commit_push.go`:
  - Update `RunCommitPushFeature` to format via `ParseAndFormatCommitMessage("cpf", cleanArgs)`.
  - Update `RunCommitPushBug` to format via `ParseAndFormatCommitMessage("cbf", cleanArgs)`.
  - Update `RunCommitPushChore` to format via `ParseAndFormatCommitMessage("cpc", cleanArgs)`.
  - Update `RunCommitPushRelease` to format via `ParseAndFormatCommitMessage("cpr", cleanArgs)`.
- [x] **2.3**: Update `cli/cmdcommit/commit_help_menu.go` to present bracketed module syntax examples.
- [x] **2.4**: Create `cli/cmdcommit/commit_message_parser_test.go` verifying:
  - Bracketed modules across all semantic verbs (`cpf`, `cbf`, `cpc`, `cpr`).
  - Conventional types `feat([module]): item` and `fix([module]): item`.
  - Duplicate prefix removal (e.g. `cpf "Feature: [auth]: test"` -> `Feature: [auth]: test`).
  - Legacy hyphen translation (`cpf "auth - test"` -> `Feature: [auth]: test`).
  - Special characters and slashes in module names (e.g. `[cli/cmdworkgroup]`).

### Phase 3: Workgroup Command Suite Implementation (`cli/cmdworkgroup`)
- [x] **3.1**: Create `cli/cmdworkgroup/workgroup_types.go`:
  - Struct `WorkgroupRecord` (`ID`, `Name`, `RootDirectory`, `SSHKeyPath`, `Description`, `GitUserName`, `GitUserEmail`, `IsDefault`, `CreatedAt`, `UpdatedAt`).
  - Struct `FolderBindingRecord` (`ID`, `FolderPath`, `WorkgroupID`, `SSHKeyPath`, `GitUserName`, `GitUserEmail`, `CreatedAt`).
  - Struct `WorkgroupOptions` (`Subcommand`, `Name`, `Path`, `SSHKeyPath`, `Description`, `UserName`, `UserEmail`, `IsDefault`, `IsJSON`).
- [x] **3.2**: Create `cli/cmdworkgroup/workgroup_flags.go` implementing flag parsing with positive booleans (`isDefault`, `isJSON`).
- [x] **3.3**: Create `cli/cmdworkgroup/workgroup_resolver.go` implementing hierarchical ancestor lineage resolution.
- [x] **3.4**: Create `cli/cmdworkgroup/workgroup_test_ssh.go` executing isolated `ssh -T -i <key> git@<host>` verification.
- [x] **3.5**: Create `cli/cmdworkgroup/workgroup_ui.go` rendering Lipgloss tables and formatted detail cards.
- [x] **3.6**: Create `cli/cmdworkgroup/workgroup_cmd.go` implementing `RunWorkgroup(args []string) error` with subcommand dispatcher.
- [x] **3.7**: Create `cli/cmdworkgroup/workgroup_test.go` validating subcommands, resolver, and flag parser.

### Phase 4: UI Help Menus & HD Markdown Help Documentation
- [x] **4.1**: Author `cli/helpdoc/workgroup.md` with full high-definition documentation, synopsis, options, SQLite details, and examples.
- [x] **4.2**: Author `cli/helpdoc/commit-format.md` detailing `[module]: <item>` format, type mappings, and normalization algorithms.
- [x] **4.3**: Author `cli/helpdoc/cbf.md` documenting the `cbf` alias.
- [x] **4.4**: Update `cli/helpdoc/catalog.go` and `cli/helpdoc/aliases.go` with entries for `workgroup`, `wg`, `group`, `commit-format`, and `cbf`.
- [x] **4.5**: Update existing `cli/helpdoc/workdir.md`, `cli/helpdoc/commit-push-feature.md`, and `cli/helpdoc/commit-push-bug.md` cross-referencing the new features.

### Phase 5: LLM Training Integration & Antigravity Skill Sync
- [x] **5.1**: Update `cli/cmd/llm/llm_train.go`:
  - Enhance Phase 1 [Discovery] with `gitmap workgroup ls`.
  - Enhance Phase 4 [Semantic Commit] to mandate `[module]: <item>` bracketed module syntax.
  - Update `outputTrainJSON` catalog phases and `executeSelfLoopSimulation`.
- [x] **5.2**: Update `cli/cmd/llm/llm_skill.go` embedded `SkillTemplate` with workgroup commands and bracketed commit syntax.
- [x] **5.3**: Update `.agents/skills/gitmap/skill.md` with:
  - Replaced commit message rules (authorizing `[module]: <item>` and `type([module]): <item>`).
  - Workgroup and SSH key automation cheat sheet.
- [x] **5.4**: Run `python 03-ai-scripts/38-sync-prompts-skills-scripts.py` to synchronize `.agents/skills/gitmap/skill.md` across all workspace repositories.

### Phase 6: Code Quality, Guideline Verification & Testing
- [x] **6.1**: Run `python go-format-check.py` to guarantee zero coding guideline violations:
  - Positive booleans only (`is`/`has` prefixes).
  - All functions <= 15 lines.
  - Zero nested `if` statements.
  - Mandatory blank lines after `}` and before `return`.
- [x] **6.2**: Run `go test ./cli/cmdcommit/... ./cli/cmdworkgroup/...` to achieve 100% test pass rate.

---

## 4. Acceptance Criteria

- [x] **AC-01**: `gitmap workgroup`, `gitmap wg`, and `gitmap group` dispatch identically without errors.
- [x] **AC-02**: `gitmap wg add <name> <dir> --ssh-key <key>` persists record in SQLite Split-DB without prompting.
- [x] **AC-03**: Ancestor directory matching correctly resolves SSH key for nested subdirectories.
- [x] **AC-04**: `gitmap cpf "[auth]: add jwt"` produces commit message `Feature: [auth]: add jwt`.
- [x] **AC-05**: `gitmap cbf "[db]: fix lock"` produces commit message `Bug: [db]: fix lock`.
- [x] **AC-06**: `gitmap cpf "feat([workgroup]): add ssh binding"` preserves conventional type formatting.
- [x] **AC-07**: `gitmap cpf "Feature: [auth]: test"` strips duplicate `Feature:` cleanly.
- [x] **AC-08**: HD help cards (`workgroup.md`, `commit-format.md`, `cbf.md`) render in terminal with `gitmap hd <topic>`.
- [x] **AC-09**: `gitmap llm train` outputs updated 6-phase curriculum reflecting workgroups and bracketed commits.
- [x] **AC-10**: Antigravity skill `.agents/skills/gitmap/skill.md` is updated and synchronized across all fleet repositories.
