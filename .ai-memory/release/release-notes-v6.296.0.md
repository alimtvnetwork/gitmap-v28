## Quick Install v6.296.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.296.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.296.0/install.sh | bash
```

## Changelog v6.296.0

- Fixed SSH join routing: routed 'add' subcommand to RunSSHJoinCLI to support 'gitmap ssh add <user@ip> <alias>'
- Restored unmasked public key display: formatDisplayPublicKey unconditionally outputs full key, and printExistingKeyOnDisk copies to OS clipboard
- Added Memory Rule: enforced TOTAL BAN in .ai-memory/strictly-avoid.md against public key masking and clipboard omission
- Fixed Known Hosts verification: auto-purged stale host keys in known_hosts during join and login to prevent 'REMOTE HOST IDENTIFICATION HAS CHANGED'
- Fixed SSH Password & PAM Auth: supported keyboard-interactive authentication alongside password and enforced auth check before reporting join success
- Preserved Error Stack Traces: populated Stack and Caller in executeClientCmd and enabled stack trace output on E_INTERNAL_ERROR
- Fixed SSH node list and removal: used openSSHDBFunc in fetchSJHosts and auto-confirmed removal in non-interactive terminals
- Documented 4-part Root Cause Analysis in .ai-memory/issues/16-ssh-join-hostkey-and-auth-rca.md
