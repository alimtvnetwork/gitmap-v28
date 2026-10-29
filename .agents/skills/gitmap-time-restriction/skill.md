---
name: gitmap-time-restriction
description: On-call skill for enforcing maintenance time windows, repository date/time standards (DD-MMM-YY, UTC persistence), and historical commit replays via GitMap commitin.
---

# GitMap Time Restriction & Recommit Skill

## 1. Overview & On-Call Invariant

This skill provides operational procedures, automated command sequences, and verification protocols for time-restricted repository workflows and historical commit normalization.

> [!IMPORTANT]
> **On-Call Invariant**: This skill is strictly **on-call**. It remains dormant during normal engineering turns. It activates **only** when explicitly commanded by user prompts, plan directives, or when tasks specify time-window boundaries, commit date overrides, or recommits.
> **Invocation:** `@[gitmap-time-restriction](file;.agents/skills/gitmap-time-restriction)`  

### Trigger Keywords & Directives
- `time restriction`, `time-restriction`, `restrict hours`, `work window`, `cutoff hour`
- `recommit`, `re-commit`, `commitin`, `backdate commit`
- `date format`, `date-only branch`, `DD-MMM-YY`, `--tz`, `--auto-window`
- Prompt reference: `[gitmap-time-restriction](file;.agents/skills/gitmap-time-restriction)`

---

## 2. Universal Date & Time Standard

All repositories, scripts, commands, and prompts adhering to this skill must follow the **Universal Date/Time Standard**:

| Dimension | Standard Format | Example | Invariant / Rule |
| :--- | :--- | :--- | :--- |
| **Storage & Persistence** | UTC Normalized | `2026-10-29T03:24:00Z` | Always call `t.UTC()` before writing to databases, git trees, or JSON logs |
| **Compact Date Token** | `DD-MMM-YY` (3-char upper month) | `29-OCT-26` | Standard date format across commits, tags, and daily manifests |
| **Display / Telemetry** | Local 12-hour AM/PM | `29-Oct-2026 11:24 AM` | Standardized via `FormatDisplayDate`; local workstation time |
| **Branch Naming** | Date-only slug | `feature/29-OCT-26-user-change` | Date portion only; **zero** time-of-day tokens permitted |
| **Time Pollution** | Forbidden | No `11-24-AM` in branches/tags | Prevent timestamp clutter in git references and documentation |

---

## 3. Time Window Enforcement Protocol

When `isTimeRestricted` is active, commits and automated pushes are strictly governed by target timezone perspective and schedule constraints.

### 3.1 Target Timezone Resolution Hierarchy

Commit timestamp evaluation and cutoff checks always resolve timezone in strict order of precedence:
1. **CLI Flag (`--tz`)**: Highest precedence (e.g. `--tz "America/New_York"` or `--tz "UTC"`).
2. **Repository / Task Configuration**: Configured `targetTimeZone` setting in task or repository configuration.
3. **Local Machine Timezone**: Fallback to workstation local timezone if no timezone flag or config is provided.

### 3.2 Operating Window & Schedule Constraints

The maintenance window enforces strict weekday versus weekend scheduling rules calculated in the resolved **target timezone**:

- **Weekday Constraint (Monday – Friday)**:
  - Commits and pushes are strictly permitted **BEFORE 08:00 AM** in the target timezone.
  - Any commit attempt at or after 08:00 AM target local time on weekdays is blocked from immediate execution and must be deferred or scheduled to a valid window.
- **Weekend Rule (Saturday – Sunday)**:
  - Commits and pushes are permitted **ANYTIME** across all 24 hours of Saturday and Sunday without cutoff restrictions.

### 3.3 Operating Window Configuration

```json
{
  "isTimeRestricted": true,
  "hasTargetTimezone": true,
  "targetTimeZone": "America/New_York",
  "weekdayCutoffHour": 8,
  "hasWeekendExemption": true,
  "hasAutoWindow": true,
  "isHeadless": false,
  "isStrictBlock": true
}
```

### 3.4 Interactive Prompt Mode vs. One-Liner Flag Mode

GitMap provides two operational modes for time-restricted workflows:

#### Mode A: Interactive Prompt Mode (Human Operator / Interactive Terminal)
When invoked without flags, GitMap prompts the operator interactively:
1. Prompts for target timezone selection or confirms default configured timezone.
2. Evaluates current target time against the weekday `< 08:00 AM` rule or weekend exemption.
3. If outside the valid window, prompts the operator to choose between scheduling a future timestamp, selecting an earlier pre-08:00 AM window, or aborting.

```bash
# Interactive time-restricted commit flow
gitmap commitin prompt
```

#### Mode B: One-Liner Flag Mode (Headless CI / AI Agent Automation)
For autonomous agents, headless CI/CD pipelines, and scripts, all parameters are passed directly via command flags to eliminate blocking prompts:

- `--tz <timezone>`: Explicitly specifies the target timezone (e.g. `--tz "America/New_York"`).
- `--schedule "<timestamp>"`: Explicitly sets a commit timestamp adhering to the `DD-MMM-YY hh:mm:ss` format.
- `--auto-window`: Automatically resolves and schedules the commit timestamp to the nearest valid window (before 08:00 AM on weekdays or anytime on weekends).
- `--non-interactive`: Disables interactive input prompts, failing fast or executing automatically with provided flags.

```bash
# Headless time-restricted commit with target timezone and auto-window
gitmap cpf "Module - updated feature" --tz "America/New_York" --auto-window --non-interactive

# Explicitly scheduled commit within valid target timezone window
gitmap cpc "Docs - updated readme" --tz "UTC" --schedule "29-OCT-26 07:15:00" --non-interactive
```

### 3.5 Pre-Commit Time Validation Sequence

Before executing any commit verb when `isTimeRestricted` is active:

1. **Query Current Time in Target Timezone**:
   ```bash
   gitmap time-now --tz "America/New_York" --format "DD-MMM-YY hh:mm A"
   ```
2. **Evaluate Schedule Compliance**:
   - If day of week is Saturday or Sunday: commit execution is immediately permitted (`isPermitted: true`).
   - If day of week is Monday through Friday: verify that current target hour is strictly before 08:00 AM.
   - If Monday through Friday and current target hour is 08:00 AM or later:
     - If `--auto-window` is set: automatically shift commit timestamp to valid window before 08:00 AM.
     - If `isStrictBlock` is active without auto-window: halt execution and report schedule block.

---

## 4. Recommit & History Adjustment Workflow

When a task requires recommitting previous changes, modifying commit history, or synchronizing commit trees under time restrictions:

### 4.1 Single Commit Recommit Recipe

When amending or adjusting the immediate HEAD commit without altering file content:

```bash
# 1. Verify working tree is clean
gitmap st

# 2. Check current commit metadata and author timestamp
gitmap log -n 1

# 3. Reset HEAD soft without remote push
gitmap git-reset HEAD~1 --no-push

# 4. Re-execute semantic commit with target timezone and window flags
gitmap cpf "Module - updated feature with standardized timestamps" --tz "America/New_York" --auto-window --non-interactive
```

### 4.2 Batch Historical Normalization (`commitin`)

To replay or normalize commit streams across branches or repositories while respecting time restriction rules:

```bash
# 1. Discover target commits to replay
gitmap commitin inspect --source <source-branch> --target <target-branch>

# 2. Walk and replay commits adhering to target timezone and cutoff rules
gitmap commitin walk --source <source-branch> --target <target-branch> --tz "America/New_York" --auto-window --normalize-dates

# 3. Verify clean state post-walk
gitmap st
```

---

## 5. Coding & Style Guidelines Compliance

When authoring code, scripts, or plans supporting time restriction:

1. **Strict Relative Git Paths**:
   - Always reference `.agents/skills/gitmap-time-restriction/skill.md`, `02-spec/...`, `.ai-memory/...`.
   - Never use absolute paths (`C:\...`, `D:\...`).
2. **Positive Boolean Naming**:
   - Use `isTimeRestricted`, `hasTargetTimezone`, `hasWeekendExemption`, `hasAutoWindow`, `isHeadless`, `isPermitted`, `isWeekday`, `isWeekend`.
   - Ban negative booleans (`isNotPermitted`, `isOutsideWindow`) and non-standard prefixes (`canCommit`, `shouldWait`).
3. **Leaf Error Sovereignty**:
   - Helper functions that validate timestamps return `error` or `*apperror.AppError` directly.
   - Never invoke `os.Exit` or `cliexit.HandleError` inside leaf helpers.
4. **Function Size**:
   - Keep time parsing and formatting functions under 15 lines (<= 8 lines preferred).

---

## 6. Verification & Health Check

Validate the skill installation and compliance:

```bash
# Verify skill file presence
gitmap find --pattern "skill.md" --path ".agents/skills/gitmap-time-restriction"

# Validate YAML frontmatter
gitmap cat .agents/skills/gitmap-time-restriction/skill.md
```
