# Specification: Zero Magic Strings, Strongly-Typed Enums & V6 Error Management Integration

## 1. Zero Magic Strings Policy
Per `02-spec/17-consolidated-guidelines/07-enum-standards.md`, all string literals representing known domain sets (subcommands, flags, types, action codes) must be defined as strongly-typed PascalCase enums ending in `Type` in dedicated `constants/` files.

### Standardized Enum Families
1. **Workgroup Commands (`WorkgroupCmdType`):**
   `WorkgroupCmdList`, `WorkgroupCmdAdd`, `WorkgroupCmdBind`, `WorkgroupCmdUnbind`, `WorkgroupCmdShow`, `WorkgroupCmdTestSSH`, `WorkgroupCmdRemove`.
2. **Workgroup Flags (`WorkgroupFlagType`):**
   `WorkgroupFlagSSHKey`, `WorkgroupFlagDesc`, `WorkgroupFlagHost`, `WorkgroupFlagUserName`, `WorkgroupFlagUserEmail`, `WorkgroupFlagDefault`, `WorkgroupFlagJSON`.
3. **Commit Verbs (`CommitVerbType`):**
   `CommitVerbFeature` (`"cpf"`), `CommitVerbBug` (`"cbf"`), `CommitVerbRefactor` (`"cfr"`), `CommitVerbChore` (`"cpc"`), `CommitVerbRelease` (`"cpr"`).
4. **Conventional Commit Types (`ConventionalCommitType`):**
   `ConventionalFeat` (`"FEAT"`), `ConventionalFix` (`"FIX"`), `ConventionalRefactor` (`"REFACTOR"`), `ConventionalChore` (`"CHORE"`), `ConventionalRelease` (`"RELEASE"`).
5. **Error Codes (`ErrorCodeType`):**
   `ErrCodeGenericExecution` (`"E9000"`), `ErrCodeSSHKeyEmpty` (`"E9001"`), `ErrCodeSSHKeyMissing` (`"E9002"`), `ErrCodeWorkgroupAddUsage` (`"E9003"`), `ErrCodeWorkgroupBindUsage` (`"E9004"`), `ErrCodeWorkgroupNoBoundKey` (`"E9005"`), `ErrCodeWorkgroupDeleteUsage` (`"E9006"`).

---

## 2. Embedding `02-spec/03-error-manage/` into V6 Execute Prompt

The canonical autonomous execution prompt (`.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`) must be updated with the following core architectural invariants:

1. **Leaf Error Sovereignty:**
   Leaf functions returning an `error` or `*appfault.AppError` must return the error directly (`return err` or `return appfault.Wrap(...)`). Leaf functions must NEVER invoke exit handlers (`cliexit.HandleError`) or `os.Exit`.
2. **Outer Caller Handles Exits:**
   Only root command dispatchers receive errors and decide how to terminate, ensuring every error is captured by database audit logging (`finishCommandAudit`).
3. **No Magic Integer Exit Codes:**
   Raw integers (`1`, `2`) are strictly banned in exit handlers. Strongly-typed enums ending in `Type` (`constants.ExitGeneralError`, `constants.ExitUsageError`) must be passed.
4. **Parameter Reduction via Specialized Helpers:**
   Use `cliexit.HandleValidationError`, `cliexit.HandleUsageError`.
5. **Positive Booleans Only:**
   All boolean variables and fields must begin with `is` or `has` exclusively. Negative booleans (`notValid`, `disabled`) and banned prefixes (`can`, `should`, `was`, `will`) are strictly forbidden.
6. **15-Line Function Limit & Zero Nested Ifs:**
   Functions <= 15 lines max (<= 8 lines preferred). Flatten all conditionals with early return guard clauses.
