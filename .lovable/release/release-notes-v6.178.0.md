## Quick Install v6.178.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.178.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.178.0/install.sh | bash
```

## Changelog v6.178.0

- Displayed dirty and modified files list during interactive reconciliation (`gitmap reconcile` and interactive prompt in `gitmap pull`)
- Added comprehensive file status categories (`modified:`, `untracked:`, `deleted:`, `staged:`) with 10-file display cap and overflow summary
- Refactored porcelain dirty state inspection in `gitutil/dirty_inspect.go` to strictly adhere to coding guidelines (<=15 lines per function)
- Flattened nested conditionals in `cmd/reconcile_prompt.go` using early return guard clauses to achieve 0 nested if violations
- Fixed Windows `cmd.exe` escape character issue in `.lovable/ai-fix-scripts/06-cicd-local-runner.py` by quoting `-run="^$"` regex in Compile Gate
- Added unit test suites in `gitutil/dirty_inspect_test.go` and `cmd/reconcile_cmd_test.go`
- Validated all 25 CI/CD quality gates with 100% green pass (exit code 0)
