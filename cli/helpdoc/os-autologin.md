# gitmap os autologin

Configure native operating system automatic login for Windows and Ubuntu workstations.

## Usage

```bash
gitmap os autologin [subcommand] [flags]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| *(none)* | Launch interactive 3-parameter credential setup (User, Domain, Password) |
| status (st) | Inspect current auto-login status and display manager |
| enable (on) | Enable auto-login with direct command-line arguments |
| disable (off) | Disable auto-login and revert to standard password prompt |
| help | Display help and usage information |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| -u, --user <username> | Current User | Operating system account username |
| -d, --domain <domain> | . | Domain name (defaults to local machine '.') |
| -p, --pass <password> | "" | Account password |

## Examples

### Interactive Setup (Prompts for Username, Domain, and Password)

```bash
gitmap os autologin
```

### Inspect Configuration Status

```bash
gitmap os autologin status
```

### Direct Enablement with Arguments

```bash
gitmap os autologin enable -u developer -d . -p SecretPass123
```

### Disable Auto-Login

```bash
gitmap os autologin disable
```
