---
name: gitmap-time-restriction
description: On-call skill for enforcing maintenance time windows, repository date/time standards (DD-MMM-YY, UTC persistence), and historical commit replays via GitMap commitin.
---

# GitMap Time Restriction & Recommit Skill

## 1. Overview & On-Call Invariant

This skill provides operational procedures, automated command sequences, and verification protocols for time-restricted repository workflows and historical commit normalization.

> [!IMPORTANT]
> **On-Call Invariant**: This skill is strictly **on-call**. It remains dormant during normal engineering turns. It activates **only** when explicitly commanded by user prompts, plan directives, or when tasks specify time-window boundaries, commit date overrides, or recommits.

### Trigger Keywords & Directives
- `time restriction`, `time-restriction`, `restrict hours`, `work window`
- `recommit`, `re-commit`, `commitin`, `backdate commit`
- `date format`, `date-only branch`, `DD-MMM-YY`
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

When `isTimeRestricted` is active, commits and automated pushes are governed by configured maintenance windows.

### 3.1 Operating Window Configuration

```json
{
  "isTimeRestricted": true,
  "hasWindowLimit": true,
  "startHourLocal": 9,
  "endHourLocal": 18,
  "timeZone": "Local",
  "isStrictBlock": true
}
```

### 3.2 Pre-Commit Time Validation Sequence

Before executing any commit verb, the agent verifies window compliance:

1. **Query Current Time**:
   ```bash
   gitmap time-now --format "DD-MMM-YY hh:mm A"
   ```
2. **Evaluate Window Condition**:
   - If current hour is between `startHourLocal` and `endHourLocal` (inclusive of start, exclusive of end) on Monday through Friday, proceed normally.
   - If current time is outside the window and `isStrictBlock` is `true`:
     - **Pause / Gate**: Halt commit execution.
     - **Override / Defer**: Apply scheduled commit timestamp using `gitmap commit --time "<DD-MMM-YY hh:mm:ss>"` when explicitly authorized.

---

## 4. Recommit & History Adjustment Workflow

When a task requires recommitting previous changes, modifying commit history, or synchronizing commit trees:

### 4.1 Single Commit Recommit Recipe

When amending or adjusting the immediate HEAD commit without altering file content:

```bash
# 1. Verify working tree is clean
gitmap st

# 2. Check current commit metadata
gitmap log -n 1

# 3. Reset HEAD soft without remote push
gitmap git-reset HEAD~1 --no-push

# 4. Re-execute semantic commit with standardized verb and time
gitmap cpf "Module - updated feature with standardized timestamps"
```

### 4.2 Batch Historical Normalization (`commitin`)

To replay or normalize commit streams across branches or repositories:

```bash
# 1. Discover target commits to replay
gitmap commitin inspect --source <source-branch> --target <target-branch>

# 2. Walk and replay commits adhering to date formatting
gitmap commitin walk --source <source-branch> --target <target-branch> --normalize-dates

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
   - Use `isTimeRestricted`, `hasWindowLimit`, `isWithinWindow`, `isValidTimestamp`.
   - Ban negative booleans (`isNotExpired`, `isOutsideWindow`) and non-standard prefixes (`canCommit`, `shouldWait`).
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
