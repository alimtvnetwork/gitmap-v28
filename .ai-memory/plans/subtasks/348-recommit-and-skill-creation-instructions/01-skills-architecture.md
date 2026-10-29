# Subtask Plan 01: Skills Architecture & Authoring Specification

- **Subtask Slug:** `01-skills-architecture`
- **Parent Plan:** `348-recommit-and-skill-creation-instructions`
- **Specification Reference:** `02-spec/21-app/348-recommit-and-skill-creation-instructions/01-architecture-spec.md`
- **Status:** IN_PROGRESS
- **Target Files:**
  - `.agents/skills/gitmap-time-restriction/skill.md`
  - `.agents/skills/gitmap-automation-bot/skill.md`

---

## 1. Objectives & Scope

This subtask defines the comprehensive operational execution plan for authoring and validating two new on-call Antigravity skills within the repository:
1. **`gitmap-time-restriction`**: Enforces strict maintenance time-window restrictions, historical commit date rewrites via GitMap's `commitin` engine, and the repository-wide date/time standard (`DD-MMM-YY`, UTC persistence, local display, date-only branches, zero time pollution).
2. **`gitmap-automation-bot`**: Provides deterministic identity sandboxing, automated switching to `"AI automation bot <local@local>"` (or configured bot profiles), and headless agent execution guardrails using GitMap's native `autobot` and `change-user` command taxonomy.

Both skills are designed as **on-call** modules: they remain strictly dormant during standard coding turns and activate only when explicitly requested in user instructions, task plans, or designated Letterly directives.

---

## 2. Target Skill Specifications & YAML Frontmatter Contracts

### 2.1 Skill 1: `gitmap-time-restriction`

- **Destination Path:** `.agents/skills/gitmap-time-restriction/skill.md`
- **Exact YAML Frontmatter Contract:**
  ```yaml
  ---
  name: gitmap-time-restriction
  description: On-call skill for enforcing maintenance time windows, repository date/time standards (DD-MMM-YY, UTC persistence), and historical commit replays via GitMap commitin.
  ---
  ```
- **Functional Scope:**
  - Operational time-window validation against UTC constraints.
  - Historical commit timestamp recalculation and normalization.
  - Recommit execution using GitMap's native `commitin` and `recommit` capabilities.
  - Strict adherence to date-only branch naming (`feature/DD-MMM-YY-slug`).
  - Zero time pollution: preventing time-of-day tokens in branches, tags, and commit titles.
- **GitMap Command Taxonomy:**
  - `gitmap recommit --date "<DD-MMM-YY>"`
  - `gitmap commitin walk --source <repo> --target <repo>`
  - `gitmap status --dirty`
  - `gitmap st`
- **Operational Guardrails:**
  - Guardrail TR-1: Commits outside authorized windows are halted until recommit date parameters are supplied.
  - Guardrail TR-2: All persisted dates must be normalized to UTC (`t.UTC()`) before writing to git objects or SQLite databases.
  - Guardrail TR-3: Display dates must use local time formatting via `FormatDisplayDate` (`DD-Mon-YYYY hh:mm AM/PM`).
  - Guardrail TR-4: Branch slugs must never contain hours, minutes, or seconds.

### 2.2 Skill 2: `gitmap-automation-bot`

- **Destination Path:** `.agents/skills/gitmap-automation-bot/skill.md`
- **Exact YAML Frontmatter Contract:**
  ```yaml
  ---
  name: gitmap-automation-bot
  description: On-call skill for switching execution context to AI automation bot identities, headless pipeline orchestration, and credential-isolated git configurations.
  ---
  ```
- **Functional Scope:**
  - Instant transition to automated identity presets (`"AI automation bot <local@local>"`).
  - Explicit scope resolution across `--local`, `--global`, and `--project`.
  - Identity backup and restoration protocols when reverting to human developer context.
  - Headless agent process management and telemetry recording.
- **GitMap Command Taxonomy:**
  - `gitmap autobot` (alias: `gitmap bot`)
  - `gitmap user autobot`
  - `gitmap change-user "<Name> <email>"`
  - `gitmap user-info` (alias: `gitmap whoami`)
  - `gitmap agent <task|subtask|log|diagnose>`
- **Operational Guardrails:**
  - Guardrail AB-1: The bot identity must default to local scope (`--local`) unless `--global` is explicitly commanded, preserving user workstation profiles.
  - Guardrail AB-2: Unauthenticated git pushes are blocked; the agent must verify credentials via `gitmap login --status` prior to push.
  - Guardrail AB-3: When bot mode is active, commit verbs must strictly use GitMap semantic verbs (`gitmap cpc`, `gitmap cbf`, `gitmap cfr`, `gitmap cpf`).
  - Guardrail AB-4: Clean identity restoration upon task completion or explicit `--restore` flag.

---

## 3. Atomic Execution Checklist

### Phase 1: Skill Design & Frontmatter Preparation
- [ ] **1.1**: Verify `.agents/skills/` directory structure and ensure no namespace collisions exist.
- [ ] **1.2**: Confirm YAML frontmatter conforms to Antigravity standard (`name` and single-sentence `description`).
- [ ] **1.3**: Validate relative path references (`.agents/skills/gitmap-time-restriction/skill.md`, `.agents/skills/gitmap-automation-bot/skill.md`).

### Phase 2: Authoring `gitmap-time-restriction`
- [ ] **2.1**: Author `.agents/skills/gitmap-time-restriction/skill.md` with:
  - On-call trigger keywords (`time-restriction`, `recommit`, `commitin`, `backdate`, `work-window`).
  - Detailed operational instructions for maintenance windows.
  - Unified Date/Time Standard section:
    - Format: `DD-MMM-YY` (compact) and `DD-Mon-YYYY hh:mm AM/PM` (display).
    - UTC persistence rule.
    - Local display rule using `FormatDisplayDate`.
    - Date-only branch naming rule (`feature/DD-MMM-YY-slug`).
    - Zero time pollution rule.
  - Step-by-step recommit execution protocol with `gitmap recommit` and `gitmap commitin`.
  - Guardrails checklist and error recovery table.

### Phase 3: Authoring `gitmap-automation-bot`
- [ ] **3.1**: Author `.agents/skills/gitmap-automation-bot/skill.md` with:
  - On-call trigger keywords (`automation-bot`, `autobot`, `bot identity`, `ai bot`, `switch to bot`).
  - Automated identity configuration protocols via `gitmap autobot` and `gitmap change-user`.
  - Scope isolation rules (`--local` default, `--global` workstation protection).
  - Pre-commit and commit execution standards with semantic commit verbs.
  - Reversion protocol to restore human developer identity after turn completion.
  - Guardrails checklist and headless operation constraints.

### Phase 4: Verification, Linting & Hygiene Audit
- [ ] **4.1**: Verify both skill markdown files strictly use relative Git paths (zero absolute paths).
- [ ] **4.2**: Verify all boolean variables and options use positive prefixes (`is*`, `has*`).
- [ ] **4.3**: Verify documentation integrity and preserve all existing skills in `.agents/skills/`.
- [ ] **4.4**: Verify neither skill attempts to execute commands during dormant inspection.

---

## 4. Coding Guidelines & Invariant Rules

| Invariant | Standard | Enforcement Rule |
|:---|:---|:---|
| **Relative Git Paths** | Always relative (e.g. `.agents/skills/...`, `02-spec/...`) | TOTAL BAN on absolute paths (`C:\...`, `/root/...`). Must be strictly relative. |
| **Boolean Naming** | Only `is*` or `has*` | NO negative booleans (`isNotRestricted`). NO banned prefixes (`can*`, `should*`, `will*`). |
| **Function Limits** | <= 15 lines max, <= 8 lines preferred | Code snippets and Go functions in examples must adhere to 15-line limit. |
| **Nested Conditionals** | 0 nested ifs | Early returns and inverted guard clauses only. |
| **Leaf Error Sovereignty** | Leaf functions return `error` directly | BANNED: `cliexit.HandleError` or `os.Exit` inside leaf helpers. Bubble to orchestrator. |
| **GitMap Command Matrix** | GitMap native commands only | TOTAL BAN on raw git commits, ripgrep, or unbuffered shell utilities. |
