## Quick Install v6.274.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.274.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.274.0/install.sh | bash
```

## Changelog v6.274.0

- Fix Go vet errors in `cli/cmdagy/agy_ping.go` (import `cli/result`, fix `res.Err` field access)
- Add missing SSH repository forwarding functions `GetSSHHostByAlias`, `GetSSHHostByIP`, and `ListSSHHosts` in `cli/store/ssh_repo.go`
- Unpack `(*SSHConnection, error)` in `cli/cmdssh/ssh_pass_cmd.go` and flatten conditional nesting
- Enforce affirmative booleans across `cli/cmdssh/ssh_pass_cmd.go`, `cli/cmdssh/sshjoin_enroll.go`, `cli/cmdagy/agy_fix_pipeline_inject.go`, and `cli/cmdagy/agy_transcript_parse.go`
- Resolve SQLite datetime parsing errors in `cli/db/sshconnection.go` by removing `COALESCE` on timestamps and scanning `NullTime` directly
- Flatten nested `if` statements across `cli/cmdagy/agy_fix_pipeline_inject.go` and `cli/cmdagy/agy_conv_scanner.go`
- Replace hardcoded Windows absolute paths in `cli/helptext/agy.md` with portable relative placeholders
- Fix Go typed nil interface evaluation in `cli/cmdmacro/macro_add_interactive.go`
- Remove unused table layout calculation helpers in `cli/cmdpull/pull_table_layout.go`
- Update unit tests and test DB reconnection mocks in `cli/cmdpull/pull_table_test.go`, `cli/cmdssh/ssh_exec_install_test.go`, and `cli/cmdssh/ssh_login_cmd_test.go`
- Optimize repository name width with lead truncation and remove commit range and changes columns in `gitmap pull`
