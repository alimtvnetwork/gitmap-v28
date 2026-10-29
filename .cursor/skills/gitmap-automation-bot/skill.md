---
name: gitmap-automation-bot
description: On-call skill for switching execution context to AI automation bot identities, headless pipeline orchestration, and credential-isolated git configurations.
---

# GitMap AI Automation Bot Skill

## 1. Overview & On-Call Invariant

This skill manages automated identity switching, sandboxed bot configurations, and headless commit authoring across local, workspace, and global scopes.

> [!IMPORTANT]
> **On-Call Invariant**: This skill is strictly **on-call**. It remains dormant during standard engineering workflows. It activates **only** when explicitly requested by user instructions, task plans, CI/CD orchestration scripts, or designated Letterly directives.

### Trigger Keywords & Directives
- `automation bot`, `autobot`, `gitmap automation`, `ai automation bot`
- `change user to bot`, `switch to bot`, `headless commit`
- Prompt reference: `[gitmap-automation-bot](file;.cursor/skills/gitmap-automation-bot)`

---

## 2. Bot Identity Standards

When bot execution is engaged, GitMap provisions a sandboxed Git identity designed to differentiate autonomous agent contributions from human developer commits.

| Parameter | Standard Value | Alternate Preset |
| :--- | :--- | :--- |
| **`user.name`** | `gitmap automation` | `AI automation bot` |
| **`user.email`** | `local@local` (or empty) | `bot@local` |
| **Default Scope** | `--local` (repository level) | `--global` (workstation level when requested) |
| **Audit Tag** | `[bot]` | `[agent]` |

---

## 3. GitMap Command Taxonomy

GitMap provides dedicated subcommands for rapid identity transitions:

### 3.1 Instant Autobot Activation

To switch the current repository to the default bot persona:

```bash
# Set repository local git identity to autobot
gitmap autobot

# Or explicitly via change-user
gitmap change-user "gitmap automation <local@local>"
```

### 3.2 Inspecting Current Identity

Verify that the active author and committer match expected bot credentials:

```bash
# Query active user identity
gitmap user-info

# Fallback identity check
gitmap whoami
```

### 3.3 Scope-Specific Identity Binding

To configure bot identity without affecting other repositories or global profiles:

```bash
# Repository-local only (recommended default)
gitmap change-user --local "gitmap automation <local@local>"

# Workstation-global (use only when instructed)
gitmap change-user --global "gitmap automation <local@local>"

# Project-level binding
gitmap user project bind autobot
```

### 3.4 Identity Reversion & Restoration

When a task requires restoring the human developer profile after bot execution:

```bash
# Revert to previously saved user profile
gitmap change-user --restore

# Or set back to specific developer profile
gitmap change-user "<Developer Name> <developer@email.com>"
```

### 3.5 Pairing Autobot with Headless Time-Restriction

When autonomous agents operate under time-restricted maintenance windows, pair `gitmap autobot` with headless time-restriction flags:

- `--tz <timezone>`: Target timezone perspective for schedule validation.
- `--auto-window`: Automatically normalizes commit timestamps to the next valid window (strictly before 08:00 AM on weekdays or anytime on weekends).
- `--non-interactive`: Enforces headless execution without interactive terminal prompts.

```bash
# Step 1: Switch repository git identity to autobot
gitmap autobot

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
   - Features, fixes, refactors, chores: `"MODULE - summary"`
   - Releases: `"v<version> <summary>"`
3. **Atomic Scope**:
   - Touch only files assigned to the active subtask.
   - Run pre-flight pull (`gitmap pull`) before writing changes.
4. **Time-Restriction Flag Integration**:
   - In automated bot workflows governed by `[gitmap-time-restriction](file;.cursor/skills/gitmap-time-restriction)`, always supply `--tz`, `--auto-window`, and `--non-interactive` to prevent interactive blocks and maintain schedule compliance.

---

## 5. Coding & Style Guidelines Compliance

When authoring code, scripts, or plans supporting the automation bot:

1. **Strict Relative Git Paths**:
   - Always reference `.cursor/skills/gitmap-automation-bot/skill.md`, `02-spec/...`, `.ai-memory/...`.
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
gitmap find --pattern "skill.md" --path ".cursor/skills/gitmap-automation-bot"

# Validate YAML frontmatter
gitmap cat .cursor/skills/gitmap-automation-bot/skill.md
```
