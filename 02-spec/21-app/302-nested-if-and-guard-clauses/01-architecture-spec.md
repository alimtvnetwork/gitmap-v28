# Architecture Spec: 302 - Nested If Elimination & Guard Clauses

## 1. Executive Summary & Intent
This specification formalizes the remediation of nested `if` statements across the repository in accordance with `01-prompts/15-cg-execute/03-nested-if-and-guard-clauses.md` and repository guidelines (`02-spec/02-coding-guidelines/`).
Zero tolerance for nested conditionals (`depth > 1`). All branching logic must be flattened using inverted guard clauses, early returns, helper extraction, and single-level control flows.

## 2. Invariants & Rules
1. **Zero Nested Ifs:** No `if` statement may be nested within another `if` block.
2. **Inverted Guard Clauses:** Fail fast and return early at the top of functions.
3. **Function Size Bounds:** Functions <= 8–15 lines.
4. **File Size Bounds:** All files must stay <= 500 lines.
5. **Blank Line Rules:** Blank line after `}` and blank line before `return`.
6. **Zero Builds & Zero Tests:** FastGate AST verification only.

## 3. Disjoint Micro-Batch Allocation
- **Worker 01 (FS & Git Utilities):**
  - Files: `cli/fsutil/recursive_top_level.go`, `cli/fsutil/scan_exclusions.go`, `cli/gitutil/exclude.go`.
  - Action: Flatten nested ifs in `DiscoverTopLevelGitRepos`, `IsValidGitRepo`, and `EnsureGitMapExcluded`.
- **Worker 02 (HelpDoc, Path & Folder CLI):**
  - Files: `cli/helpdoc/examples.go`, `cli/cmd/path_cli.go`, `cli/cmd/folder/folder.go`, `cli/cmd/folder/parse_args.go`.
  - Action: Flatten nested conditionals in `LookupHelpExample`, `checkLocalDirFallback`, `writeOrPrintOutput`, and `extractGlobExtension`.
- **Worker 03 (AutoFix Subsystem):**
  - Files: `cli/cmdautofix/fix.go`, `cli/cmdautofix/scan.go`.
  - Action: Flatten nested ifs in category checks, walk dir filters, and terminal stdin handlers.

## 4. Verification Evidence Contract
- FastGate AST inspection: `python 03-ai-scripts/50-fastgate.py` exit 0.
- Forbidden strings check: `python linter-scripts/check-forbidden-strings.py` exit 0.
