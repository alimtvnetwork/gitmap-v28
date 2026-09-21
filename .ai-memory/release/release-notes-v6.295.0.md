## Quick Install v6.295.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.295.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.295.0/install.sh | bash
```

## Changelog v6.295.0

- Fixed CI/CD compilation and lint: resolved unused strings imports and eliminated duplicate SJRmCmd declaration in cmdssh
- Fixed Error Management Policy Check: eliminated swallowed SQLite execution error in openSSHHistoryDB with proper apperror.WrapSimple
- Fixed test suites: added resolveRmTarget and validateRmTarget in cli/cmdssh/ssh_rm_target.go for sshjoin_rm_cmd_test.go
- Fixed gitmap pe stack trace leak: removed Go runtime stack trace capture from formatGHFailedError and buildGHFailedErrorMessage
- Fixed golangci-lint strict: removed unused checkHelp in cli/cmdssh/ssh_help_check.go and unused apperror import in cli/cmdpipeline/pipeline_query.go
- Documented 4-part Root Cause Analysis in .ai-memory/cicd-issues/72-ssh-build-failures-and-pe-stacktrace-rca.md
