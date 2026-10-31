---
name: gitmap-automation-bot
description: Autonomous developer companion skill for system-wide single-user AI bot identity governance, headless pipeline orchestration, and credential-isolated git configurations.
---

# GitMap AI Automation Bot Skill

## 1. Overview & System-Wide Single-User Invariant

This skill manages automated identity switching, sandboxed bot configurations, and headless commit authoring across local, workspace, and global scopes.

> [!IMPORTANT]
> **System-Wide Single Autonomous Bot Invariant**: Across all autonomous agent runs, AI pairs, and headless multi-agent workflows, GitMap enforces **ONE UNIFIED AUTONOMOUS BOT USER** across the entire workstation and all connected repositories:
> - **Git Author Name:** `gitmap automation bot - runs autonmously`
> - **Git Author Email:** `info@rgitmap.com`
> 
> AI agents are **STRICTLY FORBIDDEN** from committing under human developer credentials (`MD ALIM UL KARIM`, personal emails, etc.). Commits authored under human accounts in automated sessions cause misplaced attribution, pollute repository history, and violate the autonomous audit invariant.
>
> **Invocation:** `@[gitmap-automation-bot](file;.agents/skills/gitmap-automation-bot)`  

### Trigger Keywords & Directives
- `automation bot`, `autobot`, `gitmap automation`, `ai automation bot`
- `change user to bot`, `switch to bot`, `headless commit`, `bot identity`
- Prompt reference: `[gitmap-automation-bot](file;.agents/skills/gitmap-automation-bot)`

---

## 2. Bot Identity Standards (Universal Single Identity)

When bot execution is engaged or when any autonomous task executes, GitMap locks the Git identity machine-wide to guarantee that all autonomous agent contributions are cleanly attributed to the autonomous bot engine:

| Parameter | Standard Value | Scope & Enforcement |
| :--- | :--- | :--- |
| **`user.name`** | `gitmap automation bot - runs autonmously` | Workstation `--global` and Repository `--project` |
| **`user.email`** | `info@rgitmap.com` | Workstation `--global` and Repository `--project` |
| **Default Scope** | `--global` + `--project` | Machine-wide invariant enforced at Step 0 preflight |
| **Audit Tag** | `[bot]` | Autonomous execution provenance |

---

## 3. GitMap Command Taxonomy

GitMap provides dedicated subcommands for rapid identity locking and verification:

### 3.1 System-Wide Identity Enforcement

To lock the machine-wide global and project identity to the single autonomous bot:

```bash
# Enforce globally across workstation
gitmap user change-user "gitmap automation bot - runs autonmously" "info@rgitmap.com" --global

# Enforce in current project repository
gitmap user change-user "gitmap automation bot - runs autonmously" "info@rgitmap.com" --project
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
● Active Git: gitmap automation bot - runs autonmously <info@rgitmap.com> (local)
● Global Git: gitmap automation bot - runs autonmously <info@rgitmap.com>
```

### 3.3 Scope-Specific Identity Binding

To configure or re-bind bot identity across scopes:

```bash
# Repository-local only
gitmap change-user --local "gitmap automation bot - runs autonmously <info@rgitmap.com>"

# Workstation-global
gitmap change-user --global "gitmap automation bot - runs autonmously <info@rgitmap.com>"

# Project-level binding
gitmap user project bind autobot
```

### 3.4 Pairing Autobot with Headless Time-Restriction

When autonomous agents operate under time-restricted maintenance windows, pair the autonomous bot with headless time-restriction flags:

- `--tz <timezone>`: Target timezone perspective for schedule validation.
- `--auto-window`: Automatically normalizes commit timestamps to the next valid window (strictly before 08:00 AM on weekdays or anytime on weekends).
- `--non-interactive`: Enforces headless execution without interactive terminal prompts.

```bash
# Step 1: Switch repository git identity to autobot globally
gitmap user change-user "gitmap automation bot - runs autonmously" "info@rgitmap.com" --global

# Step 2: Execute atomic commit with headless time-restriction flags
gitmap cpf "Module - implement autonomous subtask" --tz "America/New_York" --auto-window --non-interactive
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
