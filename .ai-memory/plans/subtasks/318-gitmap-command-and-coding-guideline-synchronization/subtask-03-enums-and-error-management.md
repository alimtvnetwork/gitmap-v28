# Subtask 318-03: Zero Magic Strings, Strongly-Typed Enums & V6 Error Management Integration

## Objectives
1. Centralize and replace magic strings across `cmdworkgroup`, `cmdcommit`, `cmdgroup`, and `gitutil` with strongly-typed enums ending in `Type` in `cli/constants/constants_workgroup.go`, `cli/constants/constants_commit.go`, and `cli/constants/constants_error_codes.go`.
2. Fix `cmdcommit.OriginType` values from SCREAMING_SNAKE to PascalCase (`LocalAI`, `ForeignAI`, `HumanOrExternal`).
3. Embed `02-spec/03-error-manage/` principles (Leaf Error Sovereignty, No Magic Exit Codes, Specialized Handlers, Positive Booleans, <= 15-line limit) directly into `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`.
4. Update HD Markdown documentation (`cargo.md`, `rust.md`, `kill.md`, `workgroup.md`, `commit-format.md`), LLM training, and synchronize across all 41 fleet repositories.

## Status: COMPLETED

### Verification Summary
- **Enums Added**:
  - `cli/constants/constants_workgroup.go`: Added `WorkgroupCmdType` and `WorkgroupFlagType`.
  - `cli/constants/constants_commit.go`: Created with `CommitVerbType` and `ConventionalCommitType`.
  - `cli/constants/constants_error_codes.go`: Created with `ErrorCodeType` (`ErrCodeGenericExecution`, `ErrCodeSSHKeyEmpty`, `ErrCodeSSHKeyMissing`, `ErrCodeWorkgroupAddUsage`, `ErrCodeWorkgroupBindUsage`, `ErrCodeWorkgroupNoBoundKey`, `ErrCodeWorkgroupDeleteUsage`).
- **Magic Strings Replaced**:
  - `cli/cmdworkgroup/hook_shim.go`: Replaced raw `"E9001"` and `"E9002"` with `constants.ErrCodeSSHKeyEmpty` and `constants.ErrCodeSSHKeyMissing`.
  - `cli/cmdworkgroup/workgroup.go`: Replaced raw flag strings, subcommands, and `"E9003"`, `"E9004"`, `"E9005"`, `"E9006"` with typed constants; refactored functions to <= 15 lines.
  - `cli/cmdgroup/group.go`: Replaced delegation subcommands and flag strings with `constants.WorkgroupCmd*` and `constants.WorkgroupFlag*`.
  - `cli/cmdcommit/commit_push.go`: Replaced raw verbs and `"E9000"` with typed constants.
  - `cli/cmdcommit/commit_formatter.go`: Replaced raw verbs and conventional types with typed constants.
- **OriginType Converted**:
  - `cli/cmdcommit/commit_info_types.go`: Updated `OriginLocalAI`, `OriginForeignAI`, `OriginHumanOrExternal` to `"LocalAI"`, `"ForeignAI"`, `"HumanOrExternal"`.
- **V6 Skill Updated**:
  - `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` and `.cursor/skills/execute-parent-task-with-n-steps-v6/skill.md`: Embedded 02-spec/03-error-manage/ core principles into Section 1.1.E and Section 13 checklist.
