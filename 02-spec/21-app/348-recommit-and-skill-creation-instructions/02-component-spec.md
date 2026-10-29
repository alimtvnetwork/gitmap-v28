# 348 - Recommit and Skill Creation Instructions: Component Specification

- **Spec ID:** `02-spec/21-app/348-recommit-and-skill-creation-instructions/02-component-spec.md`
- **Subsystem:** `01-prompts/22-letterly`, `01-prompts/14-execute`, `.agents/skills`
- **Status:** APPROVED
- **Target Version:** `v6.629.0`
- **Related Specs:**
  - `02-spec/21-app/348-recommit-and-skill-creation-instructions/readme.md`
  - `02-spec/21-app/348-recommit-and-skill-creation-instructions/01-architecture-spec.md`
  - `02-spec/21-app/348-recommit-and-skill-creation-instructions/03-migration-and-verification-spec.md`
- **Related Plans:**
  - `.ai-memory/plans/pending/348-recommit-and-skill-creation-instructions.md`
  - `.ai-memory/plans/subtasks/348-recommit-and-skill-creation-instructions/02-letterly-and-execution-prompts.md`

---

## 1. Executive Summary & Component Scope

This component specification defines the architectural contracts, structural templates, and programmatic linkages for:
1. **Letterly Recommit & Time-Restricted Prompt Formatter:** `01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md`.
2. **V7 Parent Task N-Step Execution Engine:** `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md`.
3. **On-Call Skill Resolution Contracts:** Explicit relative link mechanics for `.agents/skills/gitmap-time-restriction` and `.agents/skills/gitmap-automation-bot`.
4. **Letterly Index Synchronization:** Registry update for `01-prompts/22-letterly/readme.md`.

### 1.1 Strict Invariants & Non-Negotiable Guardrails
- **V6 Prompt Immutability:** `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` is strictly frozen and immutable. It MUST NEVER be modified, renamed, or overridden.
- **On-Call Basis Only:** The time restriction skill (`gitmap-time-restriction`) and automation bot skill (`gitmap-automation-bot`) are invoked strictly on-call (when the user directive or prompt flags explicitly require them). They are never unconditionally loaded into standard execution loops.
- **Strict Relative Git Paths:** All file links, asset references, and skill triggers MUST use relative paths (`01-prompts/...`, `02-spec/...`, `.ai-memory/...`, `.agents/skills/...`). Absolute paths (`d:\...`, `C:\...`) are strictly prohibited.
- **Positive Booleans Only:** All structural flags, conditional switches, and type contracts MUST begin with `is` or `has` (`isTimeRestricted`, `hasRecommitWindow`, `isAutobotEnabled`, `isValid`). Banned prefixes: `can`, `should`, `was`, `will`, and negative booleans.
- **Leaf Error Sovereignty & 15-Line Limits:** Any auxiliary parsing logic, dispatch routines, or helper functions adhere to leaf error sovereignty (return errors directly, zero exit calls) and a maximum of 15 lines per function (<= 8 preferred).

---

## 2. Subsystem Interaction Architecture

```
+-----------------------------------------------------------------------------+
|                      User Voice Dictation / Raw Transcript                  |
+-----------------------------------------------------------------------------+
                                      |
                                      v
+-----------------------------------------------------------------------------+
|         01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly  |
|   - Lossless verbatim capture                                               |
|   - Fixed Actionable Items 1-8 + Time-Restriction & Recommit Mandates 9-12  |
|   - On-call skill triggers: gitmap-time-restriction, gitmap-automation-bot  |
|   - Execution engine suffix: execute-parent-task-with-n-steps-v7             |
+-----------------------------------------------------------------------------+
                                      |
                                      v
+-----------------------------------------------------------------------------+
|              01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7   |
|   - Multi-agent swarm orchestration (A = 3, H = 2, Capacity = 6)            |
|   - Total self-loop budget: N = 400 (Phase 1: 100, Phase 2: 150, Phase 3: 150) |
|   - Phase 1A: Chat Showcase Gate & Understood Mandate                       |
|   - Phase 1B: On-Call Skill Dispatch Gate                                   |
|   - Phase 2: Parallel execution in disjoint file boxes                      |
|   - Phase 3: Atomic GitMap commit & audit capture                           |
+-----------------------------------------------------------------------------+
                   |                                       |
                   | (on-call: isTimeRestricted)           | (on-call: isAutobotEnabled)
                   v                                       v
+-----------------------------------+   +-----------------------------------+
|  .agents/skills/                   |   |  .agents/skills/                   |
|  gitmap-time-restriction/skill.md  |   |  gitmap-automation-bot/skill.md   |
|  - Time-slot boundaries           |   |  - Profile 'autobot' activation   |
|  - Working window guards          |   |  - AI automation bot identity     |
|  - Recommit timestamp pacing      |   |  - Headless/daemon commit config  |
+-----------------------------------+   +-----------------------------------+
```

---

## 3. Component Contract: Letterly Formatter (`12-execute-recommit-time-restricted-letterly.md`)

- **File Path:** `01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md`
- **Output File Pattern:** Formatted Markdown prompt emitted directly to user or agent context.
- **Successor Skill:** `[execute-parent-task-with-n-steps-v7](file;.agents/skills/execute-parent-task-with-n-steps-v7)`

### 3.1 Structural Contract & Transformation Pipeline
The formatter transforms raw speech transcripts into non-negotiable agent execution prompts following an 8-stage transformation pipeline:

1. **Filler Eradication:** Strips conversational verbal fillers (`um`, `ah`, `uh`, `like`, leading `Okay`, excess `so`) while preserving all technical terminology, flags, and arguments verbatim.
2. **Title & Subject Generation:** Derives a clean Title Case `<task title>` and renders line 1:
   `# <task title>: high priority instruction, non-negotiable task`
3. **Verbatim Input Block:** Embeds `${Input Text Verbatim}` losslessly beneath the title.
4. **Metadata Subheaders:**
   - `## task: <task title>`
   - `## slug: <task-slug>` (lowercase-hyphenated without `XX-` prefix).
5. **Standard Non-Negotiable Items (1–8):**
   - Item 1: Mandatory Git pull pre-flight (`gitmap pull <repo path>` or fallback `git pull`).
   - Item 2: Spec authoring in `02-spec/21-app/XX-<slug>/` and plan enqueueing in `.ai-memory/plans/pending/XX-<slug>.md`.
   - Item 3: Exclusive GitMap search primacy (`gitmap aum search`, `gitmap find`, `gitmap cat`, `gitmap ps`; TOTAL BAN on raw ripgrep/grep/Select-String).
   - Item 4: Strict relative Git paths only (`02-spec/...`, `.ai-memory/...`, `cmd/...`).
   - Item 5: SQL repository pattern across all languages.
   - Item 6: Enum-based action dispatch across all operations.
   - Item 7: Mandatory enum authoring and utility methods (`String`, `Parse`/`TryParse`, `IsValid`, `All`/`Values`).
   - Item 8: Typed error management (`AppError` / `Result[T]`) at repository and service boundaries.
6. **Time-Restriction & Recommit Directives (Items 9–12):**
   - Item 9 (Time Restriction Mandate): When `isTimeRestricted` is active, enforce commit timestamps and execution pacing strictly within permitted working hours using GitMap time restriction facilities (`gitmap time-restrict` / `gitmap commit --time`).
   - Item 10 (Recommit Protocol): When `hasRecommitWindow` is specified, verify prior commit history, re-evaluate working tree state, and execute history adjustments or atomic recommits via `gitmap` semantic commit verbs (`gitmap cpf`, `gitmap cbf`, `gitmap cfr`, `gitmap cpc`, `gitmap cpr`).
   - Item 11 (AI Automation Bot Integration): When `isAutobotEnabled` is specified, configure Git author identity using `gitmap change-user --name "AI automation bot" --email "local@local"` or `gitmap autobot` before modifying code.
   - Item 12: Discrete user tasks extracted from the input dictation.
7. **On-Call Skill References:**
   - Relative links to on-call skills under `## On-Call Skills (Conditional Activation)`:
     - `[gitmap-time-restriction](file;.agents/skills/gitmap-time-restriction)`
     - `[gitmap-automation-bot](file;.agents/skills/gitmap-automation-bot)`
8. **Execution Engine Suffix:**
   - Mandatory trigger: `[execute-parent-task-with-n-steps-v7](file;.agents/skills/execute-parent-task-with-n-steps-v7)`.

### 3.2 Exact Prompt Template Specification

```markdown
# Execute Recommit & Time-Restricted — Letterly Prompt Formatter

Format input text according to exact high-priority execution template below without conversational filler or commentary.

1. Keep the verbatim as is only stripping verbal fillers (`um`, `ah`, `uh`, `like`, leading `Okay`, excess `so`) while strictly preserving everything as said parameters, flags, and file paths.
2. Derive `<task title>` as task subject. Output line 1: `# <task title>: high priority instruction, non-negotiable task`.
3. Output `${Input Text Verbatim}` directly beneath title.
4. Output task name and slug subheaders:
   - `## task: <task title>`
   - `## slug: <task-slug>` (lowercase-hyphenated; when authoring directories and plan files, format as XX-<task-slug> where XX is the zero-padded sequence of that folder).
5. Construct `# Actionable Items Must Follow Non-Negotiable`: include fixed items 1-8 verbatim, time-restriction and recommit items 9-11, followed by discrete technical directives extracted from input text as items 12..N.
6. Append mandatory on-call skill section referencing `gitmap-time-restriction` and `gitmap-automation-bot`.
7. Append mandatory agent invocation suffix pointing to `[execute-parent-task-with-n-steps-v7](file;.agents/skills/execute-parent-task-with-n-steps-v7)`.
8. Output plain markdown only. Never wrap in fenced code blocks.
9. Relative Paths: Only add relative paths (`02-spec/...`, `.ai-memory/...`, `cmd/...`); never add absolute paths.

${Input Text Verbatim} = No change except for removal of filler words.

Output Format:

# <task title>: high priority instruction, non-negotiable task

${Input Text Verbatim}

## task: <task title>
## slug: <task-slug>

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/XX-<slug>/ and enqueue plan task in .ai-memory/plans/pending/XX-<slug>.md (subtasks in .ai-memory/plans/subtasks/XX-<slug>/) first, (xx should be figured out by the folder structure and running AI tools, xx is sequence 01, 02).
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String.
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well.
5. Enforce the SQL repository pattern in every language, no exceptions: all SQL and data access lives behind a repository interface; services and handlers NEVER embed SQL strings directly; repository methods are intention-named.
6. Every action from an enum, no exceptions: every operation and command kind is a member of a named enum; dispatch only by switch/match on the enum; every switch is exhaustive.
7. Enum authoring and util-methods checkup, mandatory: define enums in shared constants/types; provide String(), Parse/TryParse, IsValid, and All/Values.
8. Error management for enums and repositories, mandatory: unknown enum values or failed parses raise typed errors (*apperror.AppError / Result[T]) with the offending value in the message.
9. Time-Restriction Enforcement (On-Call): If the task requires restricted working hours or commit schedules, activate [gitmap-time-restriction](file;.agents/skills/gitmap-time-restriction) to restrict execution and timestamp commits within the authorized window.
10. Recommit & History Adjustment (On-Call): When recommit is requested, inspect git log, verify previous commit hash, and execute atomic recommits via `gitmap` semantic commit verbs (`gitmap cpf`, `gitmap cbf`, `gitmap cfr`, `gitmap cpc`, `gitmap cpr`) without polluting git history.
11. AI Automation Bot Identity (On-Call): If automation agent mode is requested, activate [gitmap-automation-bot](file;.agents/skills/gitmap-automation-bot) and switch Git author credentials via `gitmap change-user --name "AI automation bot" --email "local@local"` before writing commits.
12. ....other steps and more steps sequentially from the input direction.

## On-Call Skills (Conditional Activation)

- [gitmap-time-restriction](file;.agents/skills/gitmap-time-restriction): Active ONLY when time windows, schedule bounds, or restricted commit pacing are specified in the prompt.
- [gitmap-automation-bot](file;.agents/skills/gitmap-automation-bot): Active ONLY when AI bot persona, automation commits, or headless CI workflows are requested.

## Must follow and spawn agent using

[execute-parent-task-with-n-steps-v7](file;.agents/skills/execute-parent-task-with-n-steps-v7)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing work for improved work or less friction, faster execution for AI and coding without errors and following Coding Guidelines.
- [/learn](slashCommand;learn) from [gitmap](file;.agents/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap RC Heatmap (`gitmap rc heatmap status` / `gitmap rc heatmap tree` / `gitmap rc heatmap live` / `gitmap rc heatmap sync`).
- Only add relative paths, never add absolute paths during work; this must be respected on release page and in release notes as well.
```

---

## 4. Component Contract: V7 Execution Prompt (`09-execute-parent-task-with-n-steps-v7.md`)

- **File Path:** `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md`
- **Predecessor (Frozen):** `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`
- **Antigravity Skill Mirror:** `.agents/skills/execute-parent-task-with-n-steps-v7/skill.md`
- **Target Version:** `v7.0.0`

### 4.1 Parameter Budgets & Swarm Hierarchy
V7 maintains identical concurrency capacity while formalizing dynamic join cycles and on-call skill gates:

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A x H = 3 agents x 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks, On-Call Skill Check)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Recommit Gate, Time-Restriction Validation, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

### 4.2 Key Architectural Differentiators (V6 vs V7)

| Feature Dimension | V6 Engine (`02-execute-...-v6.md`) | V7 Engine (`09-execute-...-v7.md`) |
| :--- | :--- | :--- |
| **Immutability Status** | Frozen (No code/prompt modifications allowed) | Active extensible release |
| **Time-Restriction Handling** | Manual user intervention only | First-class On-Call Gate (`gitmap-time-restriction`) |
| **Recommit Verification Gate** | Generic atomic commit | Explicit history audit & recommit protocol |
| **Automation Bot Identity** | Manual `git config` switches | On-Call Autobot Gate (`gitmap-automation-bot`) |
| **Heatmap Synchronization** | Static CLI mentions | Active 3-tier database collision prevention |
| **Dynamic Swarm Resiliency** | Replacements on wave completion | Immediate 0-wait respawn & companion joining |
| **Skill Bootstrap Engine** | Checks workspace skill existence | Supports both folder and API auto-registration |

### 4.3 Phase 1B: On-Call Skill Resolution Gate (New in V7)

Between Phase 1A (Turn-1 Showcase Gate) and Phase 1C (Brain Planning & Roadmap), V7 introduces the **On-Call Skill Resolution Gate**:

```go
// Pseudo-code representation of the V7 On-Call Gate evaluation logic.
// Adheres strictly to positive booleans and <= 15 lines per function.

type OnCallEvaluationInput struct {
	UserPrompt      string
	IsTimeFlagged   bool
	HasRecommitFlag bool
	IsAutobotFlag   bool
}

type OnCallSkillSet struct {
	HasTimeRestriction bool
	HasAutobotProfile  bool
	HasRecommitAction  bool
}

func evaluateOnCallSkills(input OnCallEvaluationInput) OnCallSkillSet {
	hasTime := input.IsTimeFlagged || hasTimeKeywords(input.UserPrompt)
	hasAuto := input.IsAutobotFlag || hasAutobotKeywords(input.UserPrompt)
	hasRecom := input.HasRecommitFlag || hasRecommitKeywords(input.UserPrompt)

	return OnCallSkillSet{
		HasTimeRestriction: hasTime,
		HasAutobotProfile:  hasAuto,
		HasRecommitAction:  hasRecom,
	}
}

func hasTimeKeywords(prompt string) bool {
	return strutil.ContainsAnyFold(prompt, "time restrict", "time window", "restrict hours", "schedule commit")
}

func hasAutobotKeywords(prompt string) bool {
	return strutil.ContainsAnyFold(prompt, "autobot", "automation bot", "ai bot", "headless commit")
}

func hasRecommitKeywords(prompt string) bool {
	return strutil.ContainsAnyFold(prompt, "recommit", "re-commit", "amend commit", "adjust history")
}
```

---

## 5. Component Contract: On-Call Skill Integration

### 5.1 Skill 1: `gitmap-time-restriction`
- **Relative Skill Path:** `.agents/skills/gitmap-time-restriction/skill.md`
- **Invocation Mode:** Conditional ("On-Call")
- **Trigger Conditions:**
  - Presence of flag `isTimeRestricted: true` in task configuration.
  - User directives containing time-slot limits (e.g., "only commit between 09:00 and 18:00").
  - Explicit prompt tag: `[gitmap-time-restriction](file;.agents/skills/gitmap-time-restriction)`.
- **Operational Responsibilities:**
  1. Inspect local and remote timezone configurations.
  2. Validate current timestamp against authorized execution windows before emitting git commit commands.
  3. When outside the window, either pause execution or format commit timestamps using `gitmap commit --time` (when authorized).
  4. Yield structured execution reports indicating time-gate compliance.

### 5.2 Skill 2: `gitmap-automation-bot`
- **Relative Skill Path:** `.agents/skills/gitmap-automation-bot/skill.md`
- **Invocation Mode:** Conditional ("On-Call")
- **Trigger Conditions:**
  - Presence of flag `isAutobotEnabled: true` in task configuration.
  - Headless CI/CD automation or scheduled maintenance runs.
  - User directives specifying bot authoring (e.g., "commit as AI automation bot").
  - Explicit prompt tag: `[gitmap-automation-bot](file;.agents/skills/gitmap-automation-bot)`.
- **Operational Responsibilities:**
  1. Verify Autobot profile availability in `~/.gitmap/data/git_profiles.json` (`ID: 'autobot'`, `Name: 'AI automation bot'`, `Email: 'local@local'`).
  2. Invoke `gitmap change-user autobot` or apply localized Git identity (`user.name "AI automation bot"`, `user.email "local@local"`).
  3. Ensure zero credential leaks and maintain positive boolean state tracking (`isAutobotActive: true`).
  4. Restore previous user identity upon completion if temporary elevation was requested.

---

## 6. Component Contract: Letterly Directory Readme Update

- **Target File:** `01-prompts/22-letterly/readme.md`
- **Modification Type:** Row addition in the Directory Index table.
- **Entry Specification:**

| Column | Value |
| :--- | :--- |
| **#** | `**12**` |
| **Prompt File** | `[`12-execute-recommit-time-restricted-letterly.md`](12-execute-recommit-time-restricted-letterly.md)` |
| **Target Mode** | Recommit & Time-Restricted |
| **Description** | Formats voice dictation with on-call time restrictions, recommit instructions, autobot identity, and V7 execution engine |
| **Canonical Successor Skill** | `execute-parent-task-with-n-steps-v7` |

- **Integrity Rule:** All existing rows (01 through 11) must remain intact without alterations.
- **Cursor Suite Synchronization:** If `01-prompts/23-cursor-prompts/` is synchronized in future tasks, a corresponding cursor mirror (`12-execute-recommit-time-restricted-letterly-cursor.md`) will map to `.cursor/skills/execute-parent-task-with-n-steps-v7`.

---

## 7. Data Models & Type Contracts

All models are authored with positive booleans and PascalCase identifiers:

### 7.1 `TimeRestrictionConfig`
```go
package model

// TimeRestrictionConfig encapsulates scheduling constraints for commit authoring.
type TimeRestrictionConfig struct {
	IsTimeRestricted bool   `json:"isTimeRestricted"`
	HasWindowLimit   bool   `json:"hasWindowLimit"`
	StartHour        int    `json:"startHour"`
	EndHour          int    `json:"endHour"`
	TimeZone         string `json:"timeZone"`
	IsStrictBlock    bool   `json:"isStrictBlock"`
}
```

### 7.2 `RecommitConfig`
```go
package model

// RecommitConfig specifies options for re-evaluating or amending commits.
type RecommitConfig struct {
	HasRecommitAction bool   `json:"hasRecommitAction"`
	IsAmendMode       bool   `json:"isAmendMode"`
	TargetCommitHash  string `json:"targetCommitHash,omitempty"`
	SemanticVerb      string `json:"semanticVerb"`
	HasCleanTree      bool   `json:"hasCleanTree"`
}
```

### 7.3 `AutobotConfig`
```go
package model

// AutobotConfig configures bot identity injection for Git operations.
type AutobotConfig struct {
	IsAutobotEnabled bool   `json:"isAutobotEnabled"`
	IsGlobalScope    bool   `json:"isGlobalScope"`
	ProfileAlias     string `json:"profileAlias"`
	AuthorName       string `json:"authorName"`
	AuthorEmail      string `json:"authorEmail"`
	HasAutoRevert    bool   `json:"hasAutoRevert"`
}
```

---

## 8. Error Management & Leaf Sovereignty

In compliance with repository guidelines:
1. **Direct Error Returns:** Leaf functions that parse prompt options or evaluate on-call triggers return errors directly (`*apperror.AppError` or standard Go `error`).
2. **Zero `os.Exit` Calls:** Leaf modules MUST NEVER invoke `cliexit.HandleError` or `os.Exit`.
3. **Structured Context Wrapping:** Errors bubble to root dispatchers with contextual metadata:

```go
func ValidateTimeWindow(cfg TimeRestrictionConfig, currentHour int) error {
	if !cfg.IsTimeRestricted {
		return nil
	}

	if currentHour < cfg.StartHour || currentHour >= cfg.EndHour {
		return apperror.NewValidationError(
			constants.ErrCodeTimeRestrictionViolation,
			"current time is outside authorized commit window",
		)
	}

	return nil
}
```

---

## 9. Verification & Acceptance Criteria

| ID | Verification Item | Target Standard | Pass Condition |
| :--- | :--- | :--- | :--- |
| **V-01** | V6 Prompt Immutability | `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md` | Unchanged (`git diff` is zero bytes) |
| **V-02** | Letterly Formatter Path | `01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md` | Exists and conforms to template |
| **V-03** | V7 Prompt Path | `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md` | Exists and includes On-Call Gate |
| **V-04** | Letterly Readme Index | `01-prompts/22-letterly/readme.md` | Row 12 registered with V7 successor |
| **V-05** | Relative Path Hygiene | All authored markdown files | Zero drive letters (`C:`, `D:`) or root slashes |
| **V-06** | Positive Boolean Audit | All types, variables, and flags | 100% compliant (`is*`, `has*` only) |
