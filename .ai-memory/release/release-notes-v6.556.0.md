# GitMap v6.556.0

## What's Changed in v6.556.0

- **`gitmap os --json` / `gitmap os json`**: Fixed subcommand dispatch to output structured `MachineIdentity` JSON without throwing `E_INVALID_OS_SUBCMD`.
- **`gitmap os`**: Enhanced 3-part layout displaying Modern OS menu, Active Host Machine Identity, and sample real-world commands at the bottom.
- **`gitmap cpu`, `processor`, `ram`, `disk`**: Instant hardware metrics across Windows (PowerShell CIM average query) and Linux (`/proc/stat`, `/proc/meminfo`, `syscall.Statfs`) with 5-second cross-process caching (`gitmap_hw_metrics_cache.json`) to eliminate CPU spikes.
- **`gitmap os network` / `gitmap network`**: Network adapter discovery, MAC addresses, IPv4/IPv6 addresses, link status (`isUp`), link speed, and byte/packet counters with `--json` support.
- **`gitmap nodes tasks`**: Distributed task queue interrogation across cluster nodes, reporting CPU, RAM, active/pending task counts, heavy-duty task classification (`isHeavyDuty`), and role-based scheduling placeholder (`?`).
- **Unified Node Discovery**: Fixed node enumeration in `nodes_resources.go` by integrating unified node discovery (`collectUnifiedNodes`).
- **Terminal Help Documentation**: Added documentation in `cli/helpdoc/` (`cpu.md`, `ram.md`, `disk.md`, `network.md`, `nodes-tasks.md`) and command palette definitions in `src/data/commands.ts`.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.556.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.556.0/install.sh | sh
```
