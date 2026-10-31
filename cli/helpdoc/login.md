# gitmap login

Log in and configure authentication for Git providers (GitHub, GitLab) so private-repository clone, pull, and push work seamlessly without credentials prompts.

## Usage

```bash
gitmap login                  # Interactive menu (SSH, PAT, Browser, Status)
gitmap login ssh [alias] [-y] # Generate or reuse SSH key and copy public key to clipboard
gitmap login token [PAT]      # Log in with a Personal Access Token
gitmap login browser-flow     # Log in via browser OAuth / GitHub CLI device flow
gitmap login --status         # Show current login state and active identity
gitmap login --use-ai-auto    # Switch global Git identity to AI bot and backup original user
gitmap login --use-ai-auto --restriction-time # Switch to AI bot and auto-apply wait defaults
gitmap login --revert-to-original # Revert Git identity back to backed-up original user
gitmap logout                 # Remove stored GitHub credentials
```

## Interactive Menu

When invoked with no arguments (`gitmap login`), an interactive menu is displayed:

```
GitMap Authentication:
  [1] SSH Key Setup & Copy (Recommended)
  [2] Personal Access Token (PAT)
  [3] Browser Device Flow
  [4] Check Authentication Status
```

## Subcommands

### 1. SSH Key Authentication (`gitmap login ssh`)

```bash
gitmap login ssh [alias] [-y] [-d <dir>]
```

- **Alias:** Optional name suffix for the SSH key (`id_ed25519_<alias>`). If omitted, defaults to `id_ed25519`.
- **Key Reuse vs Generation:** If the key already exists, GitMap automatically reuses it. Pass `-y`, `--yes`, `-f`, or `--force` to regenerate and overwrite.
- **Clipboard Integration:** Automatically copies the public key (`.pub`) to the OS clipboard using `clip.exe` (Windows), `pbcopy` (macOS), or `wl-copy`/`xclip` (Linux).
- **Settings URLs:** Displays direct browser links to paste your new public key:
  - GitHub: https://github.com/settings/ssh/new
  - GitLab: https://gitlab.com/-/user_settings/ssh_keys
- **Repository Binding:** If executed inside a Git repository, automatically configures repo-local `core.sshCommand`:
  `git config core.sshCommand "ssh -i <SSHKeyPath> -F /dev/null -o IdentitiesOnly=yes"`

### 2. Personal Access Token (`gitmap login token`)

```bash
gitmap login token [PAT] [--no-verify]
gitmap login --token <PAT> [--no-verify]
```

- Accepts a Personal Access Token directly as an argument, or securely prompts for hidden input if omitted.
- Pre-validates the token against `https://api.github.com/user` to verify permissions and extract the account name.
- Stores the token in global git configuration (`github.token`) so CLI and Git operations reuse it.
- Pass `--no-verify` to skip the remote network probe (useful in offline or air-gapped environments).

### 3. Browser Device Flow (`gitmap login browser-flow`)

```bash
gitmap login browser-flow
gitmap login browser
gitmap login --web
```

- Launches GitHub authentication using the GitHub CLI (`gh auth login -w`) when available.
- If `gh` is not installed, opens the token generation page in your browser and prompts for secure paste.

### 4. Authentication Status (`gitmap login --status`)

```bash
gitmap login --status
```

- Inspects system environment, Windows Credential Manager, global git configuration, and GitHub CLI to display current authentication state with masked credentials.

### 5. AI Automation Identity (`gitmap login --use-ai-auto`)

```bash
gitmap login --use-ai-auto [--restriction-time]
```

- Automatically preserves the original user's Git identity (`user.name`, `user.email`, and active GitMap profile) before switching.
- Configures global Git identity to:
  - **Name:** `gitmap automation bot - runs autonmously`
  - **Email:** `info@rgitmap.com`
- Registers and activates the `autobot` profile in GitMap.
- Pass `--restriction-time` to auto-apply anti-credit-waste dynamic waiting defaults (20s Pipeline-AI delay, 10s–30s dynamic ETA sleep intervals, 60s runner wait, and rapid polling ban).

### 6. Revert to Original User (`gitmap login --revert-to-original`)

```bash
gitmap login --revert-to-original
```

- Restores the backed-up `user.name` and `user.email` to global Git configuration.
- Reactivates the original user's GitMap profile alias.

## Flags

| Flag | Description |
|------|-------------|
| `-y`, `--yes`, `-f`, `--force` | Overwrite existing SSH key when running `gitmap login ssh`. |
| `-d`, `--dir <dir>` | Specify custom directory for SSH keys (defaults to `~/.ssh`). |
| `--token <PAT>` | Provide Personal Access Token via flag. |
| `--web`, `--browser` | Trigger browser login flow. |
| `--status` | Show login state and masked token. |
| `--no-verify` | Store token without remote validation. |
| `--use-ai-auto` | Switch global Git identity to AI bot (`gitmap automation bot - runs autonmously <info@rgitmap.com>`) and backup original user. |
| `--restriction-time` | Auto-apply time restriction & dynamic waiting defaults (20s delay, 10s–30s ETA intervals). |
| `--revert-to-original` | Revert Git identity back to original backed-up user info. |
| `-h`, `--help` | Show usage help. |

## Examples

```bash
# Interactive setup
gitmap login

# Set up personal SSH key and copy to clipboard
gitmap login ssh

# Set up work SSH key
gitmap login ssh work

# Force regeneration of work SSH key
gitmap login ssh work -y

# Log in using a Personal Access Token
gitmap login token ghp_xxxxxxxxxxxxxxxxxxxx

# Log in using browser device flow
gitmap login browser-flow

# Verify current authentication state
gitmap login --status

# Switch to AI automation bot and backup original developer identity
gitmap login --use-ai-auto

# Switch to AI automation bot with automatic dynamic waiting defaults applied
gitmap login --use-ai-auto --restriction-time

# Revert back to original developer identity
gitmap login --revert-to-original
```

## See Also

- `gitmap user project bind <alias>` — Bind repository to a Git profile
- `gitmap token list` — Show resolved token and discovery source
- `gitmap logout` — Remove stored GitHub credentials
