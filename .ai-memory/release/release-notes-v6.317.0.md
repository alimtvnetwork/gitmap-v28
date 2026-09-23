## Quick Install v6.317.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.317.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.317.0/install.sh | bash
```

## Changelog v6.317.0

- Chrome Profile Auth Session & Cookies Preservation: Preserved active OAuth refresh tokens (`token_service`), binary session cookies (`Network/Cookies`), preferences credentials (`signin.allowed: true`), and Local State GAIA mappings across export/import cycles (Plan 92 / Spec 142).
- CI/CD Misspell Standardization: Standardized on American English `canceled` in `pipeline_commit_groups_test.go` and `pipeline_history_test.go` to satisfy `golangci-lint` `misspell`.
- Boolean Guidelines Compliance: Converted inverted success checks in `cli/cmdssh/sshjoin_common.go` to affirmative positive conditions and decomposed table rendering to satisfy function length limits (RCA 77).
