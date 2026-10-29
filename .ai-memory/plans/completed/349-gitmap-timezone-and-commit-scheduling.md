# Plan: 349-gitmap-timezone-and-commit-scheduling

## Status: COMPLETED
## Slug: gitmap-timezone-and-commit-scheduling
## Category: CLI & Timezone Scheduling Engine
## Target Version: v6.630.0

---

## User Request (Verbatim)

# GitMap Timezone and Commit Scheduling: high priority instruction, non-negotiable task

Usually these time restrictions would be timezone specific. So Gitmap should have the ability to put the timezone into perspective as well. And then based on that, it should make the commits in the earlier time. Always before the 8:00 AM time during the weekdays. And during the weekend, it can put any time. You can update the conditions and prompting. That means Gitmap will prompt ask these questions, but also same thing can be provided as one-liner. You need to prepare the help text. You need to have everything into the terminal and also in the UI section as well. And settings UI should also have this feature, so includes that. And also one more thing, in the previous discussion regarding the restriction of the time discussion aspect and how the voice is given wherever you have written, try to make sure that those files are identified and those files are renewed from Git history as well. Do you understand? Can it please follow through?

## task: GitMap Timezone and Commit Scheduling
## slug: gitmap-timezone-and-commit-scheduling

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/XX-gitmap-timezone-and-commit-scheduling/ and enqueue plan task in .ai-memory/plans/pending/XX-gitmap-timezone-and-commit-scheduling.md (subtasks in .ai-memory/plans/subtasks/XX-gitmap-timezone-and-commit-scheduling/) first, (xx should be figured out by the folder structure and running AI tools, xx is sequence 01, 02)
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Implement timezone awareness in GitMap to schedule commits before 8:00 AM on weekdays and any time on weekends.
6. Update GitMap conditions and prompting to include timezone-specific commit scheduling.
7. Prepare help text for terminal and UI sections to reflect new timezone and scheduling features.
8. Ensure settings UI includes timezone and commit scheduling features.
9. Identify and renew files from Git history related to previous time restriction discussions.

## Must follow and spawn agent using

@[\1](file;.agents/skills/\1)

## Additional Instructions

- /plan first before doing work for improved work or less friction, faster execution for AI and coding without errors and following Coding Guidelines.
- /learn from @[\1](file;.agents/skills/\1) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap RC Heatmap (`gitmap rc heatmap status` / `gitmap rc heatmap tree` / `gitmap rc heatmap live` / `gitmap rc heatmap sync`).
- Only add relative paths, never add absolute paths during work; this must be respected on release page and in release notes as well.

---

## 1. High-Level Execution Phases

### Phase 1: Planning, Discovery & Specifications
- Task-01: Deep discovery across timezone utilities, commit command pipelines, UI settings, and Git history inspection. (DONE)
- Task-02: Author comprehensive architecture, component, migration specs, and subtasks:
  - `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/01-architecture-spec.md`
  - `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/02-component-spec.md`
  - `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/03-migration-and-verification-spec.md`
  - `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/readme.md`
  - `.ai-memory/plans/subtasks/349-gitmap-timezone-and-commit-scheduling/01-timezone-scheduling-architecture.md`
  - `.ai-memory/plans/subtasks/349-gitmap-timezone-and-commit-scheduling/02-ui-help-and-settings.md`
  - `.ai-memory/plans/subtasks/349-gitmap-timezone-and-commit-scheduling/03-historical-renewal-and-verification.md`

### Phase 2: Core Timezone Resolution & Schedule Policy Implementation
- Task-03: Implement `cli/timeutil/timezone.go` supporting 3-tier timezone resolution (`--tz` CLI flag > Setting DB `commit_timezone` > system `time.Local`), alias mappings (`est`, `pst`, `utc`, `local`), and UTC normalization (`t.UTC()`).
- Task-04: Implement `cli/timeutil/schedule.go` enforcing the calendar rules:
  - Weekday (Monday–Friday) commits strictly before `08:00 AM` in target timezone.
  - Weekend (Saturday–Sunday) commits permitted anytime across 24 hours.

### Phase 3: Interactive Prompting, Terminal Help Text & Settings UI
- Task-05: Update `cli/cmdcommit` with flags (`--tz`, `--schedule`, `--auto-window`, `--non-interactive`) and interactive prompts when omitted.
- Task-06: Author terminal rich help menu and documentation in `cli/cmdcommit/commit_help_menu.go` and `cli/helpdoc/commit.md`.
- Task-07: Implement Web Settings UI component in `src/components/settings/CommitScheduleCard.tsx` and integrate into `src/pages/Settings.tsx`.

### Phase 4: Historical Discussion File Renewal & Cursor Ecosystem Mirroring
- Task-08: Renew `.agents/skills/gitmap-time-restriction/skill.md` and `.agents/skills/gitmap-automation-bot/skill.md`.
- Task-09: Renew `01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md`, `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md`, and `.agents/skills/execute-parent-task-with-n-steps-v7/skill.md`.
- Task-10: Mirror renewed skills into `.cursor/skills/` and create `01-prompts/23-cursor-prompts/12-execute-recommit-time-restricted-letterly-cursor.md`.

### Phase 5: Verification, Quality Audits & Release Ceremony
- Task-11: Execute mechanical verification matrix (V-01 through V-10).
- Task-12: Staged GitMap semantic commit sequence and version bump release ceremony.

---

## 2. Subtasks Manifest

1. `01-timezone-scheduling-architecture.md`: Target timezone resolution hierarchy, weekday morning cutoff policy, weekend unrestricted execution, and UTC storage invariant.
2. `02-ui-help-and-settings.md`: Component contracts for `cli/timeutil`, CLI flags, interactive terminal prompts, terminal rich help, and React Settings UI.
3. `03-historical-renewal-and-verification.md`: Identification and renewal of historical discussion files, Cursor mirroring, verification matrix (V-01..V-10), and staged GitMap atomic release ceremony.

---

## 3. Verification Acceptance Criteria Matrix (V-01 through V-10)

| ID | Gate / Acceptance Criteria | Target Deliverable | Verification Tool / Command |
| :--- | :--- | :--- | :--- |
| **V-01** | Timezone Resolution Priority: CLI `--tz` > Setting DB > system `time.Local` | `cli/timeutil/timezone.go` | Unit test checking 3-tier precedence |
| **V-02** | Weekday Morning Cutoff: Monday–Friday commits strictly `< 08:00 AM` | `cli/timeutil/schedule.go` | Unit test with weekday boundary timestamps |
| **V-03** | Weekend Unconstrained Rule: Saturday & Sunday allowed anytime | `cli/timeutil/schedule.go` | Unit test with weekend timestamps across 24h |
| **V-04** | Universal UTC Storage Invariant: All Git and DB timestamps saved in UTC | `cli/cmdcommit/`, `cli/timeutil/` | Code audit verifying `t.UTC()` on all writes |
| **V-05** | Dual Interaction Modes: Interactive terminal prompt and CLI one-liner flags | `cli/cmdcommit/commit_prompt.go` | Flag & interactive mock execution test |
| **V-06** | Terminal Help Text: Comprehensive help menus and flag documentation | `cli/cmdcommit/commit_help_menu.go` | Help command output inspection |
| **V-07** | Settings UI Configuration: Interactive timezone and schedule controls | `src/components/settings/CommitScheduleCard.tsx` | React component layout and state inspection |
| **V-08** | Historical File Renewal: 5 historical files updated to erase legacy 9-18 assumptions | 5 identified prompt and skill files | File diff inspection against commit `745bfe7` |
| **V-09** | Cursor Ecosystem Mirroring: Parity across `.cursor/skills/` and cursor prompts | `.cursor/skills/`, `01-prompts/23-cursor-prompts/` | Directory inspection & link validation |
| **V-10** | Strict Coding Guidelines: 100% relative paths, positive booleans, <= 15-line functions | All modified files | AST & pattern audit; zero violations |

---

## 4. Deliverables Inventory (Strict Relative Git Paths)

- Specifications:
  - `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/readme.md`
  - `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/01-architecture-spec.md`
  - `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/02-component-spec.md`
  - `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/03-migration-and-verification-spec.md`
- Subtask Plans:
  - `.ai-memory/plans/subtasks/349-gitmap-timezone-and-commit-scheduling/01-timezone-scheduling-architecture.md`
  - `.ai-memory/plans/subtasks/349-gitmap-timezone-and-commit-scheduling/02-ui-help-and-settings.md`
  - `.ai-memory/plans/subtasks/349-gitmap-timezone-and-commit-scheduling/03-historical-renewal-and-verification.md`
- Master Plan:
  - `.ai-memory/plans/pending/349-gitmap-timezone-and-commit-scheduling.md`
- Historical Files to Renew:
  - `.agents/skills/gitmap-time-restriction/skill.md`
  - `.agents/skills/gitmap-automation-bot/skill.md`
  - `01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md`
  - `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md`
  - `.agents/skills/execute-parent-task-with-n-steps-v7/skill.md`
- Cursor Mirrors to Create:
  - `.cursor/skills/gitmap-time-restriction/skill.md`
  - `.cursor/skills/gitmap-automation-bot/skill.md`
  - `.cursor/skills/execute-parent-task-with-n-steps-v7/skill.md`
  - `01-prompts/23-cursor-prompts/12-execute-recommit-time-restricted-letterly-cursor.md`
