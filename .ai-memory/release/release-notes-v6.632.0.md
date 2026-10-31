# GitMap v6.632.0

## What's Changed in v6.632.0

- **GitMap Login AI Auto Switch (`gitmap login --use-ai-auto`):** Automatically preserves the original user's Git identity (`user.name`, `user.email`, and active GitMap profile) before configuring global Git identity to `gitmap automation bot - runs autonmously <info@rgitmap.com>` and activating the `autobot` profile.
- **Time Restriction & Dynamic Waiting Defaults (`gitmap login --restriction-time`):** Auto-applies time restriction defaults for anti-credit-waste workflows (20s Pipeline-AI delay, 10s–30s dynamic ETA sleep intervals, 60s runner wait, and rapid polling ban).
- **Revert to Original User (`gitmap login --revert-to-original`):** Safely restores global Git identity and active GitMap profile back to the backed-up original developer credentials.
- **Canonical V7 Workflow Ingestion:** Canonical V7 parent task continuous loop with Tree Swarm Hierarchy, dynamic agent joining, prompt confidentiality, and Git author governance.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.632.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.632.0/install.sh | sh
```
