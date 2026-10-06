# 07 — Antigravity Integration and Multi-IDE Synchronization

- **Subsystem:** AI Agent Tooling & IDE Synchronization
- **Status:** Authoritative Reference

## 1. Google Antigravity (AGY) Integration
- Interacts with AGY agent daemon, injecting prompt templates (`is-done`) and error diagnostics.
- Discovers running IDE processes across Windows and Linux (`Antigravity.exe` / `Antigravity`).

## 2. Session State Preservation
- Backs up agent transcripts (`transcript.jsonl`) and decision memory to `~/.gitmap/backup/agy/`.
- Re-seeds restored workspaces upon IDE restart.

## 3. Multi-IDE Profile Synchronization
- Synchronizes workspace configuration, extensions, and theme presets across Cursor, VS Code, and Antigravity IDE.
