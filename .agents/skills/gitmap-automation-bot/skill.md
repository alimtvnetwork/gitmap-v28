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
> **Invocation:** `@[gitmap-automation-bot](file;.agents/skills/gitmap-automation-bot)`  

### Trigger Keywords & Directives
- `automation bot`, `autobot`, `gitmap automation`, `ai automation bot`
- `change user to bot`, `switch to bot`, `headless commit`, `bot identity`
- Prompt reference: `[gitmap-automation-bot](file;.agents/skills/gitmap-automation-bot)`

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
gitmap whoami
```

Verify that the output displays:
```text
● Active Git: gitmap automation bot - runs autonomously <info@rgitmap.com> (local)
● Global Git: gitmap automation bot - runs autonomously <info@rgitmap.com>
```

### 3.3 Commit Provenance & Subagent Attribution

Attach subagent attribution directly to Git author metadata, Split-DB commit history (`repodb/commit_history.db`), and cognitive heatmaps (`repodb/heatmap.db`):

```bash
# Feature commit with explicit subagent attribution
gitmap cpf "Module - implement autonomous subtask" --ai-agent "worker-01"

# Bug fix commit with explicit subagent attribution
gitmap cbf "Scanner - resolve null pointer probe" --ai-agent "worker-02"

# Chore commit with explicit subagent attribution
gitmap cpc "Docs - align prompt checklists" --ai-agent "worker-03"
```

### 3.4 Identity Reversion & Restoration

When a task requires restoring the original human developer profile after bot execution:

```bash
# Revert to previously saved human profile from original_user.json
gitmap change-user --restore

# Or set back to specific developer profile
gitmap change-user "<Developer Name> <developer@email.com>"
```

### 3.5 Pairing Autobot with Headless Time-Restriction

When autonomous agents operate under time-restricted maintenance windows, pair `gitmap login set global ai-auto` (or `gitmap login set ai-auto [agent-name]`) with headless time-restriction flags:

- `--tz <timezone>`: Target timezone perspective for schedule validation.
- `--auto-window`: Automatically normalizes commit timestamps to the next valid window (strictly before 08:00 AM on weekdays or anytime on weekends).
- `--non-interactive`: Enforces headless execution without interactive terminal prompts.

```bash
# Step 1: Switch workstation git identity to autobot globally
gitmap login set global ai-auto

# Step 2: Execute atomic commit with headless time-restriction and subagent attribution
gitmap cpf "Module - implement autonomous subtask" --ai-agent "worker-01" --tz "America/New_York" --auto-window --non-interactive
```

---

## 4. Headless & Automated Commit Protocols

When operating under bot mode, all commits must adhere to GitMap semantic commit discipline:

1. **Semantic Commit Verbs Only**:
   - `gitmap cpr` - Releases
   - `gitmap cpf` - Features
   - `gitmap cbf` - Bug fixes
   - `gitmap cfr` - Refactorings
   - `gitmap cpc` - Chores
   - **TOTAL BAN** on raw `git commit`, `git add -A`, or arbitrary commit message formats.
2. **Commit Message Format**:
   - Features, fixes, refactors, chores: `"MODULE - summary"` (hyphen-separated; no colons in argument as GitMap provides `Feature: ` or `Bug: `)
   - Releases: `"v<version> <summary>"`
3. **Atomic Scope**:
   - Touch only files assigned to the active subtask.
   - Run pre-flight pull (`gitmap pull`) before writing changes.
4. **Time-Restriction Flag Integration**:
   - In automated bot workflows governed by `[gitmap-time-restriction](file;.agents/skills/gitmap-time-restriction)`, always supply `--tz`, `--auto-window`, and `--non-interactive` to prevent interactive blocks and maintain schedule compliance.

---

## 5. Coding & Style Guidelines Compliance

When authoring code, scripts, or plans supporting the automation bot:

1. **Strict Relative Git Paths**:
   - Always reference `.agents/skills/gitmap-automation-bot/skill.md`, `02-spec/...`, `.ai-memory/...`.
   - Never use absolute paths (`C:\...`, `D:\...`).
2. **Positive Boolean Naming**:
   - Use `isAutobotActive`, `hasLocalScope`, `isRestored`, `isValidProfile`, `hasTimeRestriction`, `hasAutoWindow`, `isHeadless`.
   - Ban negative booleans (`isNotHuman`, `isOutsideWindow`) and non-standard prefixes (`canSwitch`, `shouldRestore`).
3. **Leaf Error Sovereignty**:
   - Identity configuration helpers return `error` or `*apperror.AppError` directly.
   - Never invoke `os.Exit` or `cliexit.HandleError` inside leaf helpers.
4. **Function Size**:
   - Keep profile parsing and command execution helpers under 15 lines (<= 8 lines preferred).

---

## 6. Verification & Health Check

Validate the skill installation and compliance:

```bash
# Verify skill file presence
gitmap find "skill.md" .agents/skills/gitmap-automation-bot

# Validate YAML frontmatter
gitmap cat .agents/skills/gitmap-automation-bot/skill.md
```
