# Issue Domain 07: Antigravity and IDE Integration Glitches

- **Domain:** Antigravity Agent Daemon, Prompt Injection, and Terminal Freezes
- **Status:** Consolidated Problem & Resolution Matrix

## 1. Win32 Terminal ConHost Freeze
- **Symptoms:** Interactive terminal prompts hung indefinitely when executed inside agent sub-shells.
- **Root Cause:** Blocking standard input read calls on Windows ConHost handles without timeout.
- **Resolution:** Re-architected input reader using non-blocking channels and a 30-second bounded timeout context.

## 2. Agent Brain State Loss across Updates
- **Symptoms:** Updating Antigravity binary wiped active conversation sessions and scratchpad memory.
- **Root Cause:** Installer script purged the parent config directory during version overwrites.
- **Resolution:** Added automated backup hook copying `transcript.jsonl` to `~/.gitmap/backup/agy/` before upgrades.
