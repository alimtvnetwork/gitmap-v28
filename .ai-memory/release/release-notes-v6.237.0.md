## Quick Install v6.237.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.237.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.237.0/install.sh | bash
```

## Changelog v6.237.0

- First-class add-with-pass subcommand (gitmap ssh-join add-with-pass <user>@ip password [alias]) with interactive password prompting and JSON output
- SSH RSA password encryption at rest in SQLite database using RSA-OAEP with SHA-256 via local SSH key (~/.ssh/id_rsa)
- Cross-platform OpenSSH AskPass password auto-supply (SSH_ASKPASS_REQUIRE=force), enabling seamless password logins without terminal prompts
- Comprehensive leaf help, catalog topics, and command markdown documentation parity
