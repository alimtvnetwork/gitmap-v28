# Subtask 04: Quality Gates & Verification

## Scope
- Run `go test` across `cmd/`, `store/`, and modified packages.
- Verify `gitmap pipeline errors` and `gitmap pipeline errors --json`.
- Run AST parity check `go test ./constants -run TestCmdConstantsASTParity`.
- Run nested-if linter `python linter-scripts/check-nested-ifs.py`.
- Run boolean & enum linter `python linter-scripts/check-enum-and-boolean.py`.
- Run local CI runner `python 03-ai-scripts/06-cicd-local-runner.py` to ensure all 33 quality gates pass with exit code 0.
