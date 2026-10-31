---
name: gitmap-automation-bot
description: Autonomous developer companion skill for system-wide single-user AI bot identity governance, headless pipeline orchestration, and credential-isolated git configurations.
---

# GitMap AI Automation Bot Skill

## 1. Overview & System-Wide Autonomous Bot Invariant

This skill manages automated identity switching, sandboxed bot configurations, and headless commit authoring across local, workspace, and global scopes.

> [!IMPORTANT]
> **System-Wide Autonomous Bot Identity Standard**: Across all autonomous agent runs, AI pairs, and headless multi-agent workflows, GitMap enforces standardized autonomous bot identities:
> - **Global Bot Identity:** `gitmap automation bot - runs autonomously <info@rgitmap.com>`
> - **Agent-Specific Identity:** `gitmap automation bot - <agent-name> - runs autonomously <info@rgitmap.com>`
> 
> Before applying any bot identity, GitMap automatically backs up the original human identity to `original_user.json` to enable non-destructive restoration. AI agents are **STRICTLY FORBIDDEN** from committing under human developer credentials (`MD ALIM UL KARIM`, personal emails, etc.). Commits authored under human accounts in automated sessions cause misplaced attribution, pollute repository history, and violate the autonomous audit invariant.
>
> **Invocation:** `@[gitmap-automation-bot](file;.cursor/skills/gitmap-automation-bot)`  

### Trigger Keywords & Directives
- `automation bot`, `autobot`, `gitmap automation`, `ai automation bot`
- `change user to bot`, `switch to bot`, `headless commit`, `bot identity`
- Prompt reference: `[gitmap-automation-bot](file;.cursor/skills/gitmap-automation-bot)`

---

## 2. Bot Identity Standards

When bot execution is engaged or when any autonomous task executes, GitMap provisions a sandboxed Git identity designed to differentiate autonomous agent contributions from human developer commits.

| Parameter | Standard Global Preset (`ai-auto`) | Agent-Specific Preset (`ai-auto <name>`) [MANDATORY] |
| :--- | :--- | :--- |
| **`user.name`** | `gitmap automation bot - runs autonomously` | `gitmap automation bot - <agent-name> - runs autonomously` |
| **`user.email`** | `info@rgitmap.com` | `info@rgitmap.com` |
| **Mandatory Command** | `gitmap login set global ai-auto` | `gitmap login set ai-auto <agent-name>` (MANDATORY) |
| **Default Scope** | `--global` (machine-global) | `--local` (repository-local) |
| **Human Backup** | Automatically saved to `original_user.json` | Automatically saved to `original_user.json` |
| **Commit Provenance** | `gitmap cpf/cpb/cpc/cfr/cpr --ai-agent <name>` | `gitmap cpf/cpb/cpc/cfr/cpr --ai-agent <name>` |

---

## 3. GitMap Command Taxonomy

GitMap provides dedicated subcommands for rapid identity locking, subagent attribution, and verification:

### 3.1 Instant AI Automation Bot Activation (Mandatory Agent Name & Teaching on the Fly)

To switch the workstation or repository to the standardized bot persona:

```bash
# Lock workstation-global Git identity to default AI automation bot
gitmap login set global ai-auto

# MANDATORY FOR AI AGENTS: Lock Git identity to named AI automation bot
gitmap login set ai-auto muse
gitmap login set ai-auto gemini
gitmap login set ai-auto cursor-model
gitmap login set ai-auto antigravity
gitmap login set ai-auto claude

# On-The-Fly Error Recovery:
# If an agent runs 'gitmap login set ai-auto' without an agent name (and without global),
# GitMap emits an actionable error with on-the-fly teaching examples to prompt immediate self-correction.

# Settings & Settings UI Configuration:
# Default bot name, bot email, and default fallback agent can be configured via CLI or Web Settings UI:
# Mandatory Dynamic AI Identity Rule:
# The AI agent MUST dynamically configure its own agent name (e.g. antigravity, gemini, cursor-model, muse, claude)
# based on the environment it is executing in. NEVER assume or hardcode 'muse' as default.
gitmap settings set ai_auto.default_agent <agent-name>  # e.g. antigravity, gemini, cursor-model, muse
gitmap settings set ai_auto.default_name "gitmap automation bot - runs autonomously"
gitmap settings set ai_auto.default_email "info@rgitmap.com"
# Or open Web Settings UI -> "Speed & Automation" pane -> "ai-auto" settings controls
```

### 3.2 Inspecting Active Identity Context

At Step 0 of every autonomous run, inspect and verify that both `Active Git` and `Global Git` reflect the bot:

```bash
# Query active user and global git identity
gitmap user info

# Fallback identity check
git config --get-all user.name
git config --get-all user.email
```

### 3.3 Reverting to Original Human Credentials

At the conclusion of a session or when restoring interactive development mode:

```bash
# Revert to original human credentials saved in original_user.json
gitmap login --revert-to-original

# Alternative alias
gitmap change-user --restore
```

---

## 4. Operational Invariants for AI Swarms

1. **Mandatory Step 0 Identity Switch:** The very first command of an autonomous task (Phase 1 / Step 0) MUST be `gitmap login set ai-auto <agent-name>` or `gitmap login set global ai-auto`.
2. **Zero Human Commits Invariant:** If `gitmap user info` reveals a personal name or email, DO NOT proceed with commits until bot identity is locked.
3. **Heatmap & Split-DB Provenance:** Always pass `--ai-agent <name>` during commit operations (`gitmap cpf/cpb/cpc/cfr/cpr`) to record the subagent ID in cognitive heatmap traces.
4. **Non-Destructive Restoration:** Never manually delete `original_user.json`; let `gitmap login --revert-to-original` clean up safely upon restoration.
