# GitMap v6.519.0

## What's Changed in v6.519.0

- **Benchmark-first-line** — `agent collisions`, `stats`, `heatmap`, `ps`, `ls`, and `show` now print `[Nms]` as the first line; `--json` carries `elapsed_ms` first. The cached SQLite reads measure 1–7ms.
- **Collision cache confirmed** — SQLite Tier-1 central DB (`.ai-memory/temp-agents/ai_agents.db`): `FileClaim` rows written at claim time, released on subtask complete/fail. `agent collisions` is one DB open + one query (no Tier-2 scans).
- **New `gitmap agent ls`** — last 12 tasks (configurable) with IDs, statuses, agent counts, subtask tree view; broken DBs show `BROKEN`, never crash.
- **New `gitmap agent show <id>`** — full task detail with subtask subtree, file-box counts, evidence excerpts.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.519.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.519.0/install.sh | sh
```
