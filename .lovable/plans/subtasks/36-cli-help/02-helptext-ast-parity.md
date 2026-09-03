# Subtask 02 - Helptext AST Parity Verification

## Parent Specification
[36-cli-help-parity-audit.md](.lovable/plans/pending/36-cli-help-parity-audit.md)

## Acceptance Criteria & Requirements
- Execute Go helptext consistency tests: `go test -C gitmap ./helptext/... -count=1`.
- Execute Go AST top-level command tests: `go test -C gitmap ./constants/... -run TestTopLevelCmdRegistryMatchesAST -count=1`.
- Ensure zero drift between Go CLI commands, markdown help texts, and constant definitions.
