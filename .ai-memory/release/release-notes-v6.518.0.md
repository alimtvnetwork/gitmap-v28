# GitMap v6.518.0

## What's Changed in v6.518.0

- **Dispatch registry** — new `gitmap doctor --check-dispatch`: aggregates all 1000+ command/alias claims across dispatch tables and fails fast on duplicates, so an alias hijack (like the `agm`→gitignore collision) can never ship silently again.
- **Agent collision handling** — `gitmap agent subtask claim-files --task-id <slug> --subtask <id> --files <paths>`: workers declare their file box before modifying (look-ahead); `gitmap agent collisions` reports overlaps across active subtasks (reported, never hard-blocked).
- **Scoped commit staging** — `gitmap cpf/cpb/cpc/cpr --task <slug>` now stages only that task's claimed files instead of `git add -A`; agents never run `git add` directly. This release's own commit was made through the new flow.
- **Agent observability** — new read-only methods: `gitmap agent stats` (recent completions, per-agent counts, collision totals), `gitmap agent heatmap` (file modification ranking), `gitmap agent ps` (running tasks, slugs, agent counts).
- **A08 fix** — swallowed errors in the last-error persistence path (`writeLastErrorFile`, `LogInternalErrorRecord`) now surface as visible stderr warnings instead of silent discards.
- **Version truth unified** — `version.json` carries a single canonical `Version` field; the sync checker now gates the release tag.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.518.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.518.0/install.sh | sh
```
