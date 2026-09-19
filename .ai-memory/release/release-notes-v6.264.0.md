## Quick Install v6.264.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.264.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.264.0/install.sh | bash
```

## Changelog v6.264.0

- Restore missing worker pool, chunking, and git tracking functions in 03-ai-scripts/02-shared-engine.py to fix Python linters
- Restore CICD_DIR, normalize_repo_rel, and JobResult 6-parameter constructor in 03-ai-scripts/06-cicd-local-runner.py
- Flatten all depth-2 nested if statements across cli/cmdautomation (changed_files_cmd.go, db_generate.go, db_migrate.go, help_audit.go, run_cmd.go, worker_pool.go, worker_test.go) to comply with maximum depth-1 branching guidelines
- Resolve 52 British English spelling violations across repository documentation and source files via 03-ai-scripts/27-misspell-auditor.py
- Fix result.Result[T] monad method invocations in cli/cmdai/ai_create_test.go (res.IsFailure(), res.AppError(), res.Data)
- Remove hardcoded absolute path to gitmap in 03-ai-scripts/38-sync-prompts-skills-scripts.py to comply with check-relative-paths policy
- Fix redundant newline inside fmt.Println in cli/cmdssh/ssh_auth_key_deploy.go and handle database error in cli/cmdautomation/topology.go
- Replace negative booleans with affirmative isSuccess checks across cmdautomation commands and tests
- Run gofmt across all Go source files to ensure 100% gofmt compliance
