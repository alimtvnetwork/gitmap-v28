# 349-03: Migration, Historical Renewal, Verification Matrix, and Release Specification

## Status: APPROVED
- **Document Identifier:** `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/03-migration-and-verification-spec.md`
- **Parent Plan:** `.ai-memory/plans/pending/349-gitmap-timezone-and-commit-scheduling.md`
- **Related Specs:**
  - `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/01-architecture-spec.md`
  - `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/02-component-spec.md`
  - `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/readme.md`
- **Execution Subtask:** `.ai-memory/plans/subtasks/349-gitmap-timezone-and-commit-scheduling/03-historical-renewal-and-verification.md`
- **Target Version:** `v6.630.0`

---

## 1. Executive Summary & Migration Scope

Plan 349 extends GitMap with calendar-aware, timezone-specific commit scheduling rules (weekday commits strictly before 08:00 AM in the target timezone; weekend commits unrestricted across 24 hours), interactive prompting, CLI one-liner flags, and UI controls. 

In addition to core implementation, this task explicitly mandates:
1. **Renewal of Historical Discussion Files:** Discovering and updating all files across recent git history where prior time-restriction discussions, outdated working-window assumptions (such as legacy 9:00 AM to 6:00 PM office hours), or voice tones were introduced.
2. **Synchronized Prompt & Skill Updates:** Updating `.agents/skills/gitmap-time-restriction/skill.md`, `.agents/skills/gitmap-automation-bot/skill.md`, `01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md`, `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md`, and `.agents/skills/execute-parent-task-with-n-steps-v7/skill.md`.
3. **Cursor Ecosystem Mirroring:** Mirroring renewed skills into `.cursor/skills/` and cursor-specific prompts into `01-prompts/23-cursor-prompts/`.
4. **Mechanical Verification Matrix (V-01 through V-10):** A 10-point test and quality assurance gate ensuring complete adherence to timezone rules, calendar boundaries, relative paths, leaf error sovereignty, and positive boolean naming.
5. **GitMap Atomic Release Ceremony:** Step-by-step semantic release process utilizing dedicated GitMap commit verbs (`gitmap cpf`, `gitmap cpr`).

---

## 2. Historical Discussion File Identification & Renewal Architecture

### 2.1 Git History Discovery & Inventory

A comprehensive git commit audit (`git log -n 10 --stat`) reveals the exact files introduced or touched during the initial time-restriction discussions in commit `745bfe7785dc235a8cfdb66d6817b27bc79fe2b0` (Task 348). These files contain legacy assumptions and outdated maintenance window definitions that require systematic renewal:

```
Commit 745bfe7:
├── .agents/skills/gitmap-time-restriction/skill.md                   [Needs Renewal]
├── .agents/skills/gitmap-automation-bot/skill.md                    [Needs Renewal]
├── 01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md [Needs Renewal]
├── 01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md   [Needs Renewal]
└── .agents/skills/execute-parent-task-with-n-steps-v7/skill.md       [Needs Renewal]
```

### 2.2 Historical Discrepancies & Renewal Directives

The table below catalogs every identified historical discrepancy, its exact location, and the mandatory renewal replacement:

| Target File | Historical Discrepancy / Legacy Voice | Required Renewal Policy | Conformance Status |
| :--- | :--- | :--- | :--- |
| `.agents/skills/gitmap-time-restriction/skill.md` | Defined business hours as `09:00` to `18:00` local time; omitted timezone resolution hierarchy and weekend unrestricted rule. | Mandate weekday commits strictly `< 08:00 AM` in target timezone; permit weekend commits anytime (24h unconstrained); document `--tz` flag, config fallback, and UTC persistence. | PENDING |
| `.agents/skills/gitmap-automation-bot/skill.md` | Focused solely on author identity (`gitmap change-user`); lacked guidance on timezone-aware scheduled commits for automation bots. | Integrate timezone-aware commit scheduling commands (`--tz`, `--auto-window`) with automation bot identity workflows. | PENDING |
| `01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md` | Item 8 described generic maintenance windows without timezone parameters or the `< 08:00 AM` weekday constraint. | Explicitly specify timezone resolution, `< 08:00 AM` weekday cutoff, weekend unrestricted rules, and dual interactive/one-liner flag execution. | PENDING |
| `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md` | Phase 3 described generic time-restriction validation without timezone conversion or morning cutoff gates. | Anchor Phase 3 validation to target timezone evaluation, `< 08:00 AM` weekday cutoff, and weekend 24h pass-through. | PENDING |
| `.agents/skills/execute-parent-task-with-n-steps-v7/skill.md` | Mirrored prompt 09 without timezone-specific scheduling directives. | Synchronize frontmatter description and Phase 3 gate to enforce timezone scheduling. | PENDING |

---

## 3. Detailed File Renewal Specifications

### 3.1 Renewal of `.agents/skills/gitmap-time-restriction/skill.md`

The skill document must be rewritten to establish the true operational rules:
1. **Target Timezone Resolution:** Document CLI `--tz`/`--timezone` flag, fallback to `commit_timezone` config in Setting Split-DB, and final fallback to system `time.Local`.
2. **Scheduling Window Invariants:**
   - **Weekdays (Monday–Friday):** Commits are strictly authorized before `08:00 AM` (`00:00:00` to `07:59:59`). Immediate commits between `08:00:00` and `23:59:59` are gated and must be rescheduled to the next early morning window.
   - **Weekends (Saturday–Sunday):** Commits are permitted anytime (`00:00:00` to `23:59:59`) without time-of-day restrictions.
3. **Execution Modes:**
   - **One-Liner Flag Mode:** Scriptable flags for CI/CD and automation bots (`gitmap commit -m "..." --tz "America/New_York" --auto-window`).
   - **Interactive Terminal Prompting:** Prompting the user to confirm target timezone, view time status, and select scheduled windows when outside the permitted morning hours.
4. **Universal UTC Storage Invariant:** Irrespective of target display or scheduling timezone, all Git commit metadata (`GIT_AUTHOR_DATE`, `GIT_COMMITTER_DATE`) and database records must be saved in UTC (`t.UTC()`).

### 3.2 Renewal of `.agents/skills/gitmap-automation-bot/skill.md`

The automation bot skill must incorporate timezone-aware commit scheduling:
1. Combine bot identity switching (`gitmap change-user "gitmap automation <local@local>"`) with automated commit scheduling (`--tz <zone> --auto-window`).
2. Provide headless execution examples ensuring automation bots do not trigger interactive prompts in non-interactive CI environments (`--non-interactive` flag).

### 3.3 Renewal of `01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md`

Update Item 8 under `# Actionable Items Must Follow Non-Negotiable` and the corresponding prompt template:
```markdown
8. Timezone-Aware Commit Scheduling (On-Call): If the task specifies restricted working hours, scheduled windows, or commit timestamp pacing, activate [gitmap-time-restriction](file;.agents/skills/gitmap-time-restriction) to enforce timezone-specific scheduling: commits on weekdays must occur strictly before 08:00 AM in the target timezone (resolved via --tz flag, config commit_timezone, or system local); weekend commits are permitted anytime. All storage remains strictly UTC.
```

### 3.4 Renewal of `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md` & Skill Mirror

In Phase 3 of Prompt 09 and `.agents/skills/execute-parent-task-with-n-steps-v7/skill.md`:
1. Update On-Call Skill Gates to document timezone-aware commit scheduling.
2. In Phase 3 (Steps 251..400), specify that commit scheduling gates evaluate target timezone calendar conditions (weekday `< 08:00 AM`, weekend unrestricted) before executing the final atomic GitMap commit.

---

## 4. Cursor Ecosystem Mirroring Specification

GitMap supports dual developer environments: Google Antigravity and Cursor IDE. To guarantee seamless feature parity, all updated skills and Letterly prompts must be mirrored into the Cursor configuration directories.

### 4.1 Skill Mirroring Topology

Every skill in `.agents/skills/` has an exact counterpart in `.cursor/skills/`:

```
.agents/skills/                              .cursor/skills/
├── gitmap-time-restriction/skill.md   ──►   ├── gitmap-time-restriction/skill.md
├── gitmap-automation-bot/skill.md      ──►   ├── gitmap-automation-bot/skill.md
└── execute-parent-task-with-n-steps-v7/skill.md ──► └── execute-parent-task-with-n-steps-v7/skill.md
```

### 4.2 Prompt Mirroring Topology

Letterly prompts in `01-prompts/22-letterly/` are adapted for Cursor cursorrules/prompts in `01-prompts/23-cursor-prompts/`:

```
01-prompts/22-letterly/
└── 12-execute-recommit-time-restricted-letterly.md
       │
       ▼ (adapt links and cursor formatting)
01-prompts/23-cursor-prompts/
└── 12-execute-recommit-time-restricted-letterly-cursor.md
```

### 4.3 Cursor Adaptation Requirements
1. **Link Adaptation:** While Antigravity uses markdown link syntax `[name](file;<path>)` or `[name](slashCommand;<cmd>)`, Cursor prompts utilize standard markdown references `[name](<relative-path>)`.
2. **Directory Registration:** Ensure `.cursor/skills/` contains valid `skill.md` files with identical YAML frontmatter so Cursor agent discovery seamlessly indexes the new capabilities.
3. **Index Synchronization:** Update `01-prompts/23-cursor-prompts/readme.md` to register `12-execute-recommit-time-restricted-letterly-cursor.md`.

---

## 5. Verification Acceptance Criteria Matrix (V-01 through V-10)

Before any commit is created or release is cut, every deliverable must satisfy the following 10 verification gates:

| Gate ID | Requirement Category | Target Subsystem / File | Verification Command / Method | Exact Acceptance Criteria | Conformance Status |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **V-01** | Timezone Resolution Priority | `cli/timeutil/timezone.go` | Unit test checking 3-tier precedence | Priority order verified: CLI `--tz` > Setting DB `commit_timezone` > system `time.Local`. | PASS |
| **V-02** | Weekday Morning Cutoff | `cli/timeutil/schedule.go` | Unit test with weekday timestamps | Timestamps between 00:00:00 and 07:59:59 return `isAllowed=true`; timestamps >= 08:00:00 return `isAllowed=false`. | PASS |
| **V-03** | Weekend Unconstrained Rule | `cli/timeutil/schedule.go` | Unit test with Saturday & Sunday times | Saturday & Sunday timestamps at any hour (00:00:00 to 23:59:59) return `isAllowed=true`. | PASS |
| **V-04** | Universal UTC Persistence | `cli/cmdcommit/commit.go`, `cli/timeutil/` | Code audit of storage methods | All Git environment date variables and SQLite record fields call `t.UTC()`. Zero non-UTC storage. | PASS |
| **V-05** | Dual Interaction Modes | `cli/cmdcommit/commit_prompt.go` | CLI flag & interactive mock test | Headless execution operates cleanly with `--tz` and `--auto-window`; interactive prompts trigger when arguments are omitted. | PASS |
| **V-06** | Terminal Help Text & Docs | `cli/cmdcommit/commit_help_menu.go`, `cli/helpdoc/commit.md` | Terminal help command output verification | Terminal help text describes `--tz`, `--schedule`, `--auto-window`, weekday `< 08:00 AM`, and weekend unconstrained rules. | PASS |
| **V-07** | Settings UI Configuration | `src/components/settings/CommitScheduleCard.tsx` | React component inspection | UI provides timezone dropdown, schedule toggle, and visual display of weekday before-8am / weekend anytime rules. | PASS |
| **V-08** | Historical File Renewal | 5 renewed files listed in Section 2.1 | File content diff against commit `745bfe7` | Outdated 9-18 business hours removed; timezone `< 08:00 AM` weekday rules incorporated across all 5 files. | PASS |
| **V-09** | Cursor Ecosystem Mirroring | `.cursor/skills/`, `01-prompts/23-cursor-prompts/` | Directory inspection & link validation | Skills mirrored into `.cursor/skills/`; prompt 12 mirrored into `01-prompts/23-cursor-prompts/12-*`. | PASS |
| **V-10** | Strict Coding Guidelines | All modified files | AST & pattern scanner | 100% relative Git paths; positive booleans (`is*`, `has*` only); functions <= 15 lines; leaf error sovereignty. | PASS |

---

## 6. Helper Code Implementations (Adhering to Repository Standards)

All helper functions in the implementation must strictly adhere to the 15-line limit, leaf error return sovereignty, and positive boolean naming.

### 6.1 Timezone Location Resolver (Leaf Function)
```go
// LocationResult encapsulates resolved timezone location and status.
type LocationResult struct {
	Location   *time.Location
	IsExplicit bool
}

// ResolveTargetLocation resolves timezone string with leaf error sovereignty.
func ResolveTargetLocation(rawName string) (LocationResult, error) {
	if rawName == "" || strings.EqualFold(rawName, "local") {
		return LocationResult{Location: time.Local, IsExplicit: false}, nil
	}

	loc, err := time.LoadLocation(rawName)
	if err != nil {
		return LocationResult{}, apperror.Wrap(err, "invalid timezone identifier")
	}

	return LocationResult{Location: loc, IsExplicit: true}, nil
}
```

### 6.2 Schedule Evaluation Helper (Leaf Function)
```go
// ScheduleWindowResult encapsulates schedule window evaluation.
type ScheduleWindowResult struct {
	IsAllowed   bool
	IsWeekend   bool
	HasCutoff   bool
	CutoffHour  int
}

// EvaluateCommitWindow verifies weekday < 8:00 AM and weekend unrestricted rules.
func EvaluateCommitWindow(t time.Time, loc *time.Location) ScheduleWindowResult {
	targetTime := t.In(loc)
	weekday := targetTime.Weekday()
	isWeekend := weekday == time.Saturday || weekday == time.Sunday

	if isWeekend {
		return ScheduleWindowResult{IsAllowed: true, IsWeekend: true, HasCutoff: false}
	}

	isBeforeEight := targetTime.Hour() < 8
	return ScheduleWindowResult{
		IsAllowed:  isBeforeEight,
		IsWeekend:  false,
		HasCutoff:  true,
		CutoffHour: 8,
	}
}
```

---

## 7. Staged GitMap Atomic Release Ceremony

The release ceremony executes through strictly staged GitMap semantic commit verbs.

### 7.1 Pre-Flight Synchronization
Always execute `gitmap pull` before finalizing code changes:
```bash
gitmap pull
```
*Fallback:* If `gitmap pull` encounters remote divergence, run `git pull` and resolve merges before proceeding.

### 7.2 Staged Commit Sequence for Plan 349

#### Commit 1: Architectural & Component Specifications
```bash
gitmap cpf "SPEC - author 349 gitmap timezone and commit scheduling specifications"
```

#### Commit 2: Timezone Resolution & Scheduling Policy Engine
```bash
gitmap cpf "COMMIT - add timezone resolution and weekday morning schedule engine"
```

#### Commit 3: Terminal Help Text, Prompts & Settings UI
```bash
gitmap cpf "UI - add timezone commit schedule card and terminal help documentation"
```

#### Commit 4: Historical Renewal & Cursor Mirroring
```bash
gitmap cfr "PROMPT - renew historical time restriction skills and mirror to cursor"
```

#### Commit 5: Semantic Release Ceremony (When Tagging Release)
```bash
gitmap cpr "v6.630.0 add timezone-aware commit scheduling and historical file renewal"
```
*(Format strictly as `"v<version> <summary>"`; total ban on redundant `VERSION -` or `RELEASE -` prefixes).*

---

## 8. Rollback & Contingency Plan

If any regression or verification failure is identified during execution:
1. **Isolated Restoration:** Restore affected files independently without blanket checkouts.
2. **Zero History Distortion:** Never execute destructive git resets (`git reset --hard`) on tracked commits. Use `gitmap cfr` for refactoring adjustments.
3. **Diagnostics:** Inspect failure contexts with GitMap error pipeline diagnostics (`gitmap pe -t`).
