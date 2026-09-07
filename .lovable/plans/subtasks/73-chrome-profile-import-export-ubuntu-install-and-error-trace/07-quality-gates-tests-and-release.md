# Subtask 07: Quality Gates, Linters, CI Verification & Release Orchestration

## Scope
- Run Go unit tests and AST parity verification:
  - `go test ./...`
  - `go test ./constants/... -run TestTopLevelCmdRegistryMatchesAST`
  - `go test ./helptext/... -run TestEveryHelpFileHasExamples`
- Run coding guideline linters:
  - `python 03-ai-scripts/check-nested-ifs.py`
  - `python 03-ai-scripts/check-enum-and-boolean.py`
- Run local CI runner:
  - `python 03-ai-scripts/06-cicd-local-runner.py` exit 0.
- Execute release orchestrator:
  - `python 03-ai-scripts/29-release-orchestrator.py --tier minor`

## Files Touched
- All touched files
- `changelog.md`
- `version.json`
