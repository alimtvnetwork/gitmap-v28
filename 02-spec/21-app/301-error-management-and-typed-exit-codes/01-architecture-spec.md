# Architecture Spec: 301 - Error Management & Strongly-Typed Exit Codes

## 1. Executive Summary & Intent
This specification formalizes the remediation of Error Management & Architecture violations in accordance with `01-prompts/15-cg-execute/02-error-management.md` and `02-spec/03-error-manage/`.
The goal is to eliminate raw magic integer exit codes (such as `1`, `2`) passed into `cliexit.HandleError`, replacing them with strongly-typed constants (`cliexit.ExitCodeGeneralError`, `cliexit.ExitCodeUsageError`, `cliexit.ExitCodeValidationError`) and specialized helpers (`cliexit.HandleGeneralError`, `cliexit.HandleUsageError`, `cliexit.HandleValidationError`), while preserving leaf error return sovereignty and audit log capture.

## 2. Invariants & Rules
1. **No Magic Integer Exit Codes:** Never pass un-typed integers (`1`, `2`) to `cliexit.HandleError`. Always use `int(cliexit.ExitCodeGeneralError)` or specialized helpers.
2. **Specialized Exit Helpers:** Where repetitive errors occur, prefer calling `cliexit.HandleGeneralError(err)`, `cliexit.HandleUsageError(err)`, or `cliexit.HandleValidationError(err)` over verbose parameter passing.
3. **Leaf Error Return Sovereignty:** Leaf functions returning `error` or `Result[T]` must return errors directly upward rather than calling `cliexit.HandleError` and returning `nil`.
4. **File Size Bounds:** All touched files must stay <= 500 lines.
5. **Zero Builds & Zero Tests:** FastGate AST verification only.

## 3. Disjoint Micro-Batch Allocation
- **Micro-Batch 01 (CmdZip Subsystem):**
  - Files: `cli/cmdzip/zip.go`, `cli/cmdzip/unzipcompact.go`, `cli/cmdzip/zipgroup.go`, `cli/cmdzip/zipgroupcreate.go`, `cli/cmdzip/zipgroupshow.go`, `cli/cmdzip/zipgroupops.go`.
  - Action: Replace all 31 occurrences of `cliexit.HandleError(..., 1)` with `cliexit.HandleGeneralError(...)` or `cliexit.HandleError(..., int(cliexit.ExitCodeGeneralError))`.
- **Micro-Batch 02 (CmdVSCode & HelpDoc Subsystems):**
  - Files: `cli/cmdvscode/vscode_cmd.go`, `cli/cmdvscode/vscodeworkspace.go`, `cli/helpdoc/print.go`.
  - Action: Replace magic exit codes with `cliexit.HandleGeneralError(...)` or `cliexit.HandleError(..., int(cliexit.ExitCodeGeneralError))`.
- **Micro-Batch 03 (CmdAuditLegacy Subsystem):**
  - Files: `cli/cmdauditlegacy/auditlegacy.go`.
  - Action: Replace magic exit code 2 with `cliexit.HandleUsageError(walkErr)`.

## 4. Verification Evidence Contract
- FastGate AST inspection: `python 03-ai-scripts/50-fastgate.py` exit 0.
- Forbidden strings check: `python linter-scripts/check-forbidden-strings.py` exit 0.
