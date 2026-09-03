# Git Accounts & Profiles Management (`gitmap profiles`)

Manage multiple Git accounts, providers (GitHub and GitLab), and organizations with auto-discovery and interactive numeric sequence picking.

---

## Command Overview

```bash
gitmap profiles [subcommand] [flags]
```

Aliases: `gitmap profs`, `gitmap profile git`

### Subcommands Summary

| Subcommand | Alias | Purpose |
|------------|-------|---------|
| `ls` | `list` | List all discovered and registered Git profiles with usage metrics |
| `set-default` | `default` | Set the default profile (supports interactive sequence picking) |
| `switch` | `use` | Switch the active profile (supports sequence picking) |
| `add` | `create` | Register a new user or organization profile |
| `rm` | `remove`, `delete` | Remove a profile by name or sequence number |
| `status` | - | Display current default profile, active profile, and provider |

---

## Interactive Sequence Picker (`1`, `2`, `3`...)

When running in an interactive terminal, omitting arguments will render an interactive selection menu where you can simply type the row number:

```bash
$ gitmap profiles set-default

Select default profile (Enter number 1-2):
  [1] alimtvnetwork (github, user)
  [2] auktvgo (github, organization)
Choice: 2
  ✓ Default Git profile set to: auktvgo (github)
```

In automated or CI/CD environments (`CI=1` or `GITMAP_NON_INTERACTIVE=1`), pass the index or name directly:
```bash
gitmap profiles set-default 1
```

---

## Subcommands & Examples

### 1. `gitmap profiles ls`

List all configured profiles with provider, type, default indicator, and usage frequency:

```bash
gitmap profiles ls
```

Output:
```text
  ● Git Accounts & Profiles (Total: 2)
  --------------------------------------------------------------------------------
  #    Name                 Provider   Type           Default  Usage  Last Used
  --------------------------------------------------------------------------------
  [1]  alimtvnetwork        github     user           * (default) 14     2026-09-03
  [2]  auktvgo              github     organization   -        3      2026-09-02
  --------------------------------------------------------------------------------
  Tip: Switch default with 'gitmap profiles set-default <1-N|name>'
```

Output as JSON:
```bash
gitmap profiles ls --json
```

### 2. `gitmap profiles set-default [1-N|name]`

Switch the default account used for all repository creation and backup operations:

```bash

# Set default by sequence number

gitmap profiles set-default 2

# Set default by name

gitmap profiles set-default alimtvnetwork
```

### 3. `gitmap profiles add <name> [flags]`

Register a new account or organization:

```bash

# Add a GitHub organization

gitmap profiles add riseup-asia --org

# Add a GitLab user account with email

gitmap profiles add dev-lead --gitlab --email dev@example.com
```

### 4. `gitmap profiles rm <1-N|name>`

Remove an existing profile:

```bash

# Remove profile with confirmation

gitmap profiles rm 2

# Remove without prompt in scripts

gitmap profiles rm 2 -y
```

### 5. `gitmap profiles status`

Quick status inspection:

```bash
gitmap profiles status
```

Output:
```text
  ● Active Git Profile:  alimtvnetwork
  ● Default Git Profile: alimtvnetwork
  ● Profiles Configured: 2
```
