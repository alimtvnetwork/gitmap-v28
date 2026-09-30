# gitmap os tui

Interactive Bubbletea terminal dashboard for browsing, toggling, and applying OS tweaks, auto-login credentials, display settings, DNS nameservers, and maintenance across Windows, Linux, and macOS.

## Usage

```bash
gitmap os tui [flags]
```

Aliases: `gitmap os menu`, `gitmap os dashboard`

## Flags

| Flag | Description |
|------|-------------|
| `--dry-run`, `-n` | Simulate execution without applying real system modifications |
| `-h`, `--help` | Show command help |

## Tabs & Capabilities

| Tab | Title | Description |
|-----|-------|-------------|
| `[1]` | Tweaks | Telemetry, Activity Feed, Bing Start Search, Classic Context Menu, Ultimate Power Scheme, Dark Mode Theme |
| `[2]` | Auto-Login | Current auto-login status inspection, credential removal, and configuration |
| `[3]` | Display/DM | Display sleep timeouts (never sleep vs 15m default), Wayland session mode toggle, display-manager restart |
| `[4]` | DNS & Net | High-speed secure resolvers (Cloudflare 1.1.1.1, Google 8.8.8.8, Quad9 9.9.9.9, DHCP auto-revert) |
| `[5]` | Clean & Storage | Developer cache purge (Go/npm/pnpm), Antigravity AI logs/scratch cleanup, terminal history cleanup |
| `[6]` | Update | Package repository sync, multi-distro package upgrades (winget/apt/dnf/pacman/brew), regional mirror fixes |

## Keybindings

| Key | Action |
|-----|--------|
| `Tab` / `←` / `→` | Switch tabs sequentially |
| `1` – `6` | Jump directly to tab index 1 through 6 |
| `↑` / `↓` (`k` / `j`) | Move item selection cursor |
| `Space` | Toggle checkbox `[ ]` / `[x]` |
| `a` | Select all items in active tab |
| `n` | Deselect all items in active tab |
| `Enter` | Sequentially execute all selected actions |
| `Esc` / `Backspace` | Return to dashboard after execution |
| `q` / `Ctrl+C` | Quit dashboard |

## Examples

### Launch Interactive OS Dashboard

```bash
# Launch interactive dashboard
gitmap os tui

# Launch in dry-run simulation mode
gitmap os tui --dry-run
```
