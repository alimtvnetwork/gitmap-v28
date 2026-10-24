# Subtask 05: Rollout Package 5 (`cli/cmdcommit`) & Cross-Package Verification

> **Parent Plan:** [.ai-memory/plans/pending/319-package-constants-and-enums-dry-restructuring.md](../../pending/319-package-constants-and-enums-dry-restructuring.md)  
> **Target Package:** `cli/cmdcommit`  
> **Status:** Ready for Execution (Next Turn)  

---

## 1. Scope & Technical Deliverables
1. **Enum: `CommitVerbType`** (`cli/cmdcommit/enums/commit_verb_enum.go`):
   - Values: `CommitVerbInvalid` (0), `CommitVerbFeat`, `CommitVerbFix`, `CommitVerbChore`, `CommitVerbRefactor`, `CommitVerbDocs`, `CommitVerbTest`, `CommitVerbRelease`, `CommitVerbBug`.
   - Methods: `IsValid() bool`, `String() string`, `IsFeature() bool`, `IsFix() bool`, `IsRelease() bool`.
2. **Enum: `CommitFormatType`** (`cli/cmdcommit/enums/commit_format_enum.go`):
   - Values: `CommitFormatInvalid` (0), `CommitFormatStandardBracket`, `CommitFormatHyphenSeparated`, `CommitFormatConventional`.
   - Methods: `IsValid() bool`, `String() string`, `IsHyphenSeparated() bool`.
3. **Constants: `cli/cmdcommit/constants/constants.go`**:
   - Flags: `FlagModule`, `FlagType`, `FlagSignoff`, `FlagBranch`.
   - Format strings: `FmtModulePrefix`, `FmtHyphenSummary`.
   - Validation errors: `ErrMsgRequired`, `ErrInvalidVerb`.
4. **Call-Site Refactoring:**
   - Update `cli/cmdcommit/commit.go`, `cli/cmdcommit/commit_formatter.go`, `cli/cmdcommit/types.go`.
5. **Cross-Package Verification:**
   - Execute FastGate AST validation across all 5 packages.
   - Run enum guidelines linter and forbidden string linter.
   - Run `gofmt -w` on all modified directories.

---

## 2. Coding Guidelines Invariants
- Files $\le 300$ lines.
- Functions $\le 15$ lines.
- Positive booleans only (`is*`, `has*`).
- Zero nested `if` statements.

---

## 3. Verification Evidence Contract
- Run `gofmt -w cli/cmdcommit/`.
- AST pass via FastGate.
- Automated commit via `gitmap cpf "constants - centralize enums and constants across 5 packages"`.
