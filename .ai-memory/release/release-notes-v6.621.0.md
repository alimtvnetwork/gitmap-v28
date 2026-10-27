## Quick Install v6.621.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.621.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.621.0/install.sh | bash
```

## Changelog v6.621.0

- Show-Log Non-Git Guard: Running `gitmap show-log` outside a git repository emits a clean guidance error without dumping any Go runtime stack traces.
- Process Killer (`proc kill`): Native process management command `gitmap proc kill <pid> [--force]` with self/caller/system PID protection and two-stage graceful termination.
- Heatmap Query Shortcuts: Instant telemetry views (`gitmap heatmap recent`, `gitmap heatmap failures`, `gitmap heatmap hotfiles`, `gitmap heatmap tasks`) executing sub-millisecond SQLite queries.
- Show-Log UI Modernization: Added embedded SVG favicon, compact icon buttons with hover tooltips, external diff tool configuration modal (Beyond Compare, VS Code), and searchable branch selector.
- Clone Bare Repo Resolution: `gitmap clone <name>` automatically resolves bare repository names to the highest semantic version release and converts to SSH when requested.
