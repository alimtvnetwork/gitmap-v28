# `gitmap profiles`

Manage multi-account Git provider profiles (GitHub / GitLab) and organizations with interactive sequence picking.

## Simulation

```
$ gitmap profiles ls
  ● Git Accounts & Profiles (Total: 2)
  --------------------------------------------------------------------------------
  #    Name                 Provider   Type           Default  Usage  Last Used
  --------------------------------------------------------------------------------
  [1]  alimtvnetwork        github     user           * (def)  14     2026-09-03
  [2]  riseup-asia          github     organization   -        8      2026-09-02
```

## Subcommands

```
gitmap profiles ls                             # List all configured accounts
gitmap profiles set-default [1-N|name]         # Set default account / org
gitmap profiles switch [1-N|name]              # Switch active account
gitmap profiles add <name> [flags]             # Register a new profile
gitmap profiles rm <1-N|name>                  # Remove profile
gitmap profiles status                         # Show current active default
```

## Flags

| Flag | Description |
|------|-------------|
| `--provider <github\|gitlab>` | Git provider service (default: `github`) |
| `--org` | Mark account as an organization profile |
| `--email <email>` | Author commit email address |
| `--json` | Output profile catalog in JSON format |
| `-y`, `--yes` | Skip interactive removal confirmation |

## Examples

```bash
# List all configured profiles
gitmap profiles ls

# Set default profile using numeric sequence picker
gitmap profiles set-default 2

# Add a company GitHub organization
gitmap profiles add my-company --org

# Switch active profile by name
gitmap profiles switch alimtvnetwork

# Check current active profile status
gitmap profiles status
```
