# GitMap v6.515.0

## What's Changed in v6.515.0

- **multi-user project switching, github cli auth status, and git config manager**
  - **Launch Banner Readout**: Bare `gitmap` launch now probes and displays GitHub CLI authorization status (`✔ Authorized as <username>` or `✖ Not logged in`) alongside the active Git author identity without blocking startup.
  - **`gitmap user info`**: Dedicated command (and alias `user-info`) rendering a consolidated Catppuccin Macchiato card with GitHub CLI auth status, local/global Git identities, active GitMap profile, and bound project alias with zero token leakage (`--json` supported).
  - **Multi-User Git Profiles**: Profile registration (`gitmap user add <alias> --name "<name>" --email "<email>"`), listing (`gitmap user list`), and switching (`gitmap user switch <alias>`).
  - **Per-Project Auto-Switching & Bindings**: Bind profiles per-repository (`gitmap user project bind <alias>`), unbind, and auto-sync Git local config (`gitmap user sync`).
  - **Global & Local Git Config Manager**: Inspect and set `user.name` and `user.email` cleanly via `gitmap user config [global|local]`.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.515.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.515.0/install.sh | sh
```
