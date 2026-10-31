# GitMap v6.633.0

## What's Changed in v6.633.0

- **Autonomous AI Git Identity:** Added `gitmap login set ai-auto [agent-name]` and `gitmap login set global ai-auto` (`gitmap login set ai-auto global`) supporting agent names (`muse`, `gemini`, `cursor-model`), setting `user.name` to `gitmap automation bot - (AI Agent Name) - runs autonomously` and `user.email` to `info@rgitmap.com`.
- **Original Identity Safeguard:** Automatic backup of human Git identity into `original_user.json` with multi-tier directory persistence, safe restoration via `gitmap login --revert-to-original`, and protection against bot identities overwriting backup state.
- **Commit AI Agent Provenance:** Added `--ai-agent <name>` flag to semantic commit verbs (`cpf`, `cpb`, `cpc`, `cfr`, `cpr`, `pcp`), auto-formatting commit author identity and recording agent attribution into Split-DB `CommitHistoryRecord` and cognitive heatmap tables (`HeatmapActionLog`, `HeatmapCommitLog`).
- **Heatmap AI Attribution:** Heatmap status, log, tree, and recent commands now accurately track which AI agent made each commit and operation.
- **Canonical Skills and Prompts Synchronization:** Updated `.agents/skills/gitmap/skill.md`, `.agents/skills/gitmap-automation-bot/skill.md`, and V6/V7 execution prompts across both `gitmap` and `coding-guidelines` repositories.
- **V6/V7 Mandatory Completion Checklist & Masked Verification Script:** Added end-of-run table checklists (Tasks, Coding Guidelines, File Ownership, Git Delivery), masked verification script generator, and strict zero-build/zero-test guard unless E2E testing or release is explicitly instructed.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.633.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.633.0/install.sh | sh
```
