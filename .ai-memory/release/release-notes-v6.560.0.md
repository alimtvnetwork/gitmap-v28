# GitMap v6.560.0

## What's Changed in v6.560.0

- **Parallel Fix All (`gitmap fix --all`):** Generic typed parallel worker pool dynamically scaling across CPU cores (up to 20 cores).
- **Git Repository Reclone by URL:** `gitmap reclone <url>` directly detects URLs and wipes dirty targets for clean re-cloning.
- **Redesigned Post-Pull Menu:** Option 1 is now "Fix All (Run gitmap fix --all)", followed by local Commit WIP, Stash All, Step-by-step, Fleet, and Skip.
- **Windows GCM Credential Store Fix:** Replaced invalid `wincred` with `wincredman` across credential fixers and git environment.
- **Pull All Concurrency:** Scaled worker pool up to 20 CPU cores while preserving progress bar and ANSI color escape codes.
- **Update Probe Throttle & UI Polish:** Sibling probe throttled to once every 24 hours with local watermark caching and Title Case terminal UI.
- **Repo-Cache Multi-AI Repository:** Created private `repo-cache` repository, migrated Muse 7-day memory digests, and prioritized it for LLM training.
- **Subcommand Hyphenation:** Automatic hyphen-to-space command routing for two-word commands (`train-muse`, `repo-export`, etc.).
- **Heatmap Skill Documentation:** Documented `gitmap heatmap`, `gitmap agent heatmap`, and `full-summary` in GitMap and coding-guidelines skills.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.560.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.560.0/install.sh | sh
```
