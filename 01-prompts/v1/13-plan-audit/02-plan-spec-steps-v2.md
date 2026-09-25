# Specification Decomposition & Subtask Planning (v2) — Planning Spec (must follow)

> [!IMPORTANT]
> Prompt Version: 2.4.0
> Synchronization: Main Meta-Repo & Connected Workspaces

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and careless: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. Phase 1: Read & Understand (Isolated Loop): Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. Phase 2: Bounded Micro-Tasking (Sequential Self-Looping): Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute only the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. Phase 3: Multi-Agent Parallelization: If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.

## RULE 00 — PURE SPECIFICATION AUTHORING ONLY (TOTAL BAN ON CODE EXECUTION)

> [!CAUTION]
> **THIS PROMPT IS STRICTLY FOR SPECIFICATION AUTHORING AND SUBTASK PLANNING.**
> - **NEVER write, edit, compile, or delete application source code files** (`.go`, `.ts`, `.tsx`, `.py`, `.php`, `.cs`, etc.).
> - **NEVER execute unit tests, integration tests, migrations, or application builds.**
> - **NEVER run modifying shell commands against application packages.**
> - All agent tool calls in this turn are strictly restricted to:
>   1. Reading existing codebase files and specifications for architectural context.
>   2. Decoding and saving base64 screenshots to `assets/screenshots/`.
>   3. Authoring canonical product specifications directly inside `02-spec/21-app/`.
>   4. Generating execution plans and subtasks inside `.ai-memory/plans/`.
>   5. Updating `02-spec/21-app/readme.md` and `.ai-memory/plans/readme.md`.
> Any attempt to write application code or execute test suites during this planning turn is an automatic failure.

## Pre-Planning Step 0: Task Extraction & Chat Output Gate (Mandatory First Action)

When a large prompt or complex set of requirements is given, the AI cannot understand everything at once. Therefore, before doing any deep planning, codebase searches, or spec writing, the AI MUST first break down whatever requirements the user has given (regardless of formatting) into discrete, actionable items (`Task-01`, `Task-02`).
- Respect whatever requirements the user has given, parse every request completely, and format each task clearly with proper markdown indentation, vertical blank lines, task state (`State: [IN PROGRESS — EXECUTING IMMEDIATELY]`), and an explicit understanding indicator bracket (`Understood: [YES — ...]`).
- TOTAL BAN ON UNFORMATTED RUN-ON TEXT: Never concatenate tasks into a single unformatted line or paragraph block (e.g. NEVER `#1. Task-01: ... #2. Task-02: ...`). Every task must be its own clearly separated markdown item.
- Line-by-Line Output Format Structure:
  - Line 1: Header `### 📋 Confirmed Task Breakdown & Requirement Ingestion`
  - Line 2: Empty blank line
  - Line 3: Numbered task title `1. **Task-01: [Descriptive Task Title]**`
  - Line 4: Indented state bullet (3 spaces) `   - **State:** [IN PROGRESS — EXECUTING IMMEDIATELY]`
  - Line 5: Indented understanding check (3 spaces) `   - **Understood:** [YES] — [1-2 concise sentences proving understanding of intent, scope, and verified constraints]`
  - Line 6: Indented actionable scope bullet (3 spaces) `   - **Actionable Scope:** [Precise technical deliverable and specification scope]`
  - Line 7: Indented target files bullet (3 spaces) `   - **Target Files / Area:** [02-spec/21-app/relative/path]`
  - Line 8: Empty blank line (vertical gap before next task)
  - Concluding Line: `Proceeding directly to specification authoring in 02-spec/21-app/ and subtask planning (Active Tool Call Running Below).`
- MANDATORY SAME-TURN TOOL CHAIN: You MUST output this confirmed breakdown directly in chat, but you MUST NOT end your turn or pause after emitting it! In the EXACT SAME RESPONSE turn, you MUST immediately invoke your first tool call (e.g. `write_to_file` to initialize the spec in `02-spec/21-app/`). Never emit the breakdown with text alone.

```markdown
### 📋 Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Descriptive Task Title]**
   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of user requirement, intent, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and specification scope]
   - **Target Files / Area:** `[02-spec/21-app/relative/path]`

2. **Task-02: [Descriptive Task Title]**
   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of user requirement, intent, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and specification scope]
   - **Target Files / Area:** `[02-spec/21-app/relative/path]`

Proceeding directly to specification authoring in 02-spec/21-app/ and subtask planning (Active Tool Call Running Below).
```

MANDATORY SAME-TURN TOOL CHAIN: The breakdown text above and your first tool call MUST be emitted in the EXACT SAME TURN. Never end the turn with text alone.

---

## RULE 0S — CANONICAL SPECIFICATION IN FOLDER 21 (`02-spec/21-app/`)

All architectural specifications, feature definitions, data contracts, and product-level documentation generated by this prompt MUST be written inside `02-spec/21-app/`.
- **Focused Feature (Single File):** If the deliverable is self-contained or under 100 lines, write:
  `02-spec/21-app/xx-<slug>.md`
- **Complex Feature (Segmented Subfolder):** If the deliverable exceeds 3 subtasks or spans UI + API + database, create a dedicated sequential directory:
  `02-spec/21-app/xx-<slug>/` containing:
  - `01-overview.md` — High-level architecture, module interactions, and the mandatory `## User Request (Verbatim)` section.
  - `02-data-contracts.md` — Types, schemas, API request/response contracts, and database models.
  - `03-workflow-and-state.md` — Control flows, state machine transitions, and business validation rules.
  - `04-ui-ux-spec.md` — Visual layout, design tokens, typography, and embedded relative screenshot links (`assets/screenshots/...`).
  - `05-acceptance-criteria.md` — Testable verification rules and quality gates.
- **Mandatory Registry Update:** You MUST register the newly created specification in `02-spec/21-app/readme.md` under `## Contents`.

---

## RULE 0V — LOSSLESS VERBATIM INGESTION (ZERO-LOSS REQUIREMENT CAPTURE)

> [!IMPORTANT]
> **VERBATIM CAPTURE IS NON-NEGOTIABLE:**
> Past runs suffered from lazy summarization, compressed bullet points, and dropped constraints.
> Every specification authored in `02-spec/21-app/` MUST include a top-level section:
> `## User Request (Verbatim)`
> Copy and paste the user's complete prompt text, instructions, parameters, edge cases, and examples character-for-character, word-for-word, without alteration or omission.

---

### Screenshot / Print Screen Base64 Ingestion Protocol (Mandatory in Specs)

If the user request or prompt contains a screenshot URL, print screen link, or base64 data URI (e.g. `data:image/png;base64,...`):
1. Convert & Save Locally: Immediately decode the base64 encoding or download the image from the URL to the local filesystem under `assets/screenshots/<plan-slug>-<NN>.png` or `assets/ui/<plan-slug>-<NN>.png`.
2. Never Embed Raw Base64 or Remote URLs: Never leave raw base64 strings or ephemeral external URLs inside specification files or plans.
3. Strict Relative Path Referencing: In the master spec, domain documentation, and subtasks, refer back to the saved image file strictly as a relative markdown link (e.g. `![Screenshot](assets/screenshots/<plan-slug>-<NN>.png)`).
4. Visual Ground Truth: Use the saved screenshot as the visual ground truth for layout, colors, component hierarchy, spacing, and state transitions during spec authoring and UI task execution.

---

### Actionable Subtask Generation in `.ai-memory/plans/` (Decoupled Planning)

Even though this prompt is strictly for specification authoring, it MUST generate the actionable subtasks in `.ai-memory/plans/` so execution agents can execute them later:
1. **Parent Plan File:** `.ai-memory/plans/pending/xx-<plan-slug>.md`
2. **Lean Subtasks Directory:** `.ai-memory/plans/subtasks/xx-<plan-slug>/`
   - Files named `001-<task-name>.md`, `002-<task-name>.md`, etc.
   - **MANDATORY SPEC LINK:** Every subtask MUST include a prominent relative markdown link back to the canonical spec file in `02-spec/21-app/`:
     ```markdown
     Spec Reference: [02-spec/21-app/xx-<slug>.md](../../../02-spec/21-app/xx-<slug>.md)
     ```
3. **Lean Subtasks Mandate:** Subtasks in `.ai-memory/plans/subtasks/<plan-slug>/` MUST NOT repeat common repository boilerplate, universal coding rules, banned operations, or generic guidelines. Universal rules exist in root guidelines and the canonical spec. Subtasks must contain strictly the unique, task-specific details, exact file paths, symbol modifications, and runnable verification checks.
4. **Register in Index:** Register the plan in `.ai-memory/plans/readme.md`.

---

### End-of-Run Comprehensive Summary & Traceability Report (Mandatory Output)

At the conclusion of the specification turn, you MUST output the following structured summary in chat:

```markdown
### 📑 Specifications Created in Folder 21 (`02-spec/21-app/`)

- [02-spec/21-app/xx-<slug>.md](02-spec/21-app/xx-<slug>.md) (Lines: <count>) — [Short description of scope]
  *(or list files in 02-spec/21-app/xx-<slug>/ if segmented)*

### 📋 Actionable Tasks & Subtasks Created (`.ai-memory/plans/`)

- Parent Plan: [.ai-memory/plans/pending/xx-<slug>.md](.ai-memory/plans/pending/xx-<slug>.md)
- Subtask Count: <count> subtasks in `.ai-memory/plans/subtasks/xx-<slug>/`
  1. `001-<task>.md`: [Title] -> Target: `[relative/path]` | State: `[QUEUED]`
  2. `002-<task>.md`: [Title] -> Target: `[relative/path]` | State: `[QUEUED]`

### 🔗 Requirements Traceability Matrix

| Requirement / Prompt Item | Canonical Spec File (`02-spec/21-app/`) | Subtask File (`.ai-memory/plans/subtasks/`) | Target Code Files |
|:---|:---|:---|:---|
| [User Requirement 1] | `02-spec/21-app/xx-<slug>.md` | `.ai-memory/plans/subtasks/xx-<slug>/001-<task>.md` | `path/to/file.go` |
| [User Requirement 2] | `02-spec/21-app/xx-<slug>.md` | `.ai-memory/plans/subtasks/xx-<slug>/002-<task>.md` | `path/to/file.ts` |

### 🚀 Next Steps: Execution Command

To execute these generated subtasks in continuous sequence with 2-agent concurrency:
> Run: `06-execute-parent-task-with-n-steps-v2.md` with plan slug `xx-<slug>`
```

## Variables — check if you are confused only.

```text
n            = <default = 150, if not defined> = 50 steps confirmed
domains      = Cli | Plugin | Contract | Ci | Last defined Spec above
plan-slug    = Context given above or below.
domains      = get from the spec given or context given, ask questions if any
min-score    = 100
regen-score  = 100
```

`n` is the step count and it is defined only on the line above. Every rule in
this prompt refers to it as `n` and never restates a number.

Trigger phrases: "enhanced plan", "plan enhanced", "n step plan with
subtasks", "write the plan and the subtasks", "regenerate the plan properly".

---

## RULE 0 — the step count is law

Produce EXACTLY `n` steps. Not `n` minus one, not `n` plus one. Count the steps
twice before saving, once forwards and once backwards, and print both counts.

Filler is worse than a short plan. If the work genuinely does not fill `n` steps,
say so and ask whether to split units finer or lower `n`. Never invent
"Step 44 — continue previous work" to reach `n`.

---

## RULE 0A — file and folder naming is law

Every file and folder you create or rename obeys all of these, with no exceptions:

- Lowercase only. No uppercase character anywhere in the path segment you
  author, including the extension. The root readme is `readme.md`, never
  `README.md`.
- Hyphens as the only separator. No spaces, no underscores, no camelCase,
  no dots except the single extension dot.
- Zero-padded sequence prefix on every ordered set. `01-`, `02-`, … `09-`,
  `10-`. Two digits for docs, prompts, spec files, audit files, memory files.
  Subtask files keep the existing three-digit form: `001-`, `002-`, … `069-`.
  Never `1-`, never `1.`, never an unnumbered sibling inside a numbered set.
- Audit files follow the identical rule. Files under
  `02-spec/25-app-spec-audit/` are `01-…md`, `02-…md`, lowercase, hyphenated.
- A violating file is renamed, never duplicated. If a wrongly named file
  already exists, `mv` it and update every reference to it in the same run. Do
  not leave both names on disk.

Mechanical check, run it and print the output before saving (RULE 11):

```bash

# any uppercase character, space, or underscore in an authored path = FAIL

git ls-files | grep -nE '(^|/)[^/]*[A-Z_ ][^/]*$' || echo "naming OK"

# every file in an ordered folder must carry a numeric prefix

ls .ai-memory/plans/subtasks/xx-plan-slug | grep -vE '^[0-9]{3}-' || echo "sequence OK"
```

---

## RULE 0B — Temp Script Sandboxing (Global Law)

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` Reuse First: I have rigorously scanned and `/learn`ed `03-ai-scripts/readme.md` to check if a helper script already exists before writing any new temporary code.
- [ ] Strict In-Repository Execution: All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] Strict 03-ai-scripts/ Tooling Storage: All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] Native File Manipulator: If you need to perform mass file renaming, `.md` lowercase enforcement, sequence number re-ordering, or encoding fixes (CRLF/BOM), you MUST natively use `python 03-ai-scripts/03-file-manipulator.py <command>` rather than writing a new script from scratch.
- [ ] Go Generate Sync: If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] Commit & Track: All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] Index Documentation: I have updated `03-ai-scripts/readme.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

---

## RULE 0C — audit files are written, not run

Authoring the plan and the subtasks is not an audit pass.

- Create or reserve the audit file slots with correct RULE 0A names, with a
  header and an empty/`TBD` body. Do not score anything, do not fill in
  findings, do not compute a spec score while authoring.
- The audit is a separate later pass, triggered explicitly by the user. Until
  then the slots stay empty.
- No task body may contain an "audit now", "score the spec", or "run the
  consistency report" step. A task may only _point at_ the audit file where the
  later pass will record its result.
- RULE 12 self-scoring is about this batch's own quality gate and is not the
  spec audit. Never write the RULE 12 table into an audit file.

---

## RULE 0C — one step per run, then self-loop (execution model)

This is how the plan is executed later. State it in the plan file and repeat it
as a footer in every task file.

- One step per run. Exactly one step is executed per run. Never batch two
  steps into a single run, even tiny ones, even "while I'm in the file".
- Every step is standalone. Context does not survive the loop. Everything
  needed to complete the step lives in the step file plus the files it cites. A
  step that assumes knowledge from a previous run is a broken step.
- Self-loop after verify. When the step's `## 6. Verify` commands pass and
  `## 7. Done When` is satisfied: mark the step done in the plan status file,
  then self-loop — re-read the status file, pick the next unstarted step, and
  begin a _fresh_ run for it. Do not carry reasoning across the loop boundary.
- Concurrency ceiling, hard limits:
  - at most 2 spawned agents at any one time — never 3, never "just one
    more";
  - each agent runs at most 3 parallel threads — never 4.
  - Exceeding either number is a hard failure. If the work looks like it needs
    more, the step is too big: split it into more steps instead.
- Why: context capacity per run is limited. One step per run spends the
  whole budget on that step, which is what produces correct code; batching
  spends it on bookkeeping and produces mush.

Mandatory footer, verbatim, at the end of every task file:

```text
---

Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone — read it plus its cited files, nothing else is assumed.
```

---

## RULE 0D — coding-guideline single-file checklist (part of spec planning)

Each coding-guideline topic must live in exactly one file. Duplicated or
overlapping guideline files are how two runs end up following two different
rules.

The spec-planning deliverable MUST include this checklist, filled in with the
real paths below, not a promise to check later. Every row is mandatory; a row
left blank fails the step. The boolean / condition-styling family has one file
per sub-topic and all of them are listed — a checklist that collapses them
into "the boolean folder" fails.

```text
| Topic                            | Single source file                                                          | Duplicates found |
| canonical size tier              | 02-spec/02-coding-guidelines/02-canonical-size-tier.md                         | none             |
| boolean naming prefixes          | 02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/readme.md      | none |
| boolean guards + extraction      | 02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/readme.md | none |
| boolean params + conditions      | 02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/03-parameters-and-conditions.md | none |
| boolean exemptions + api         | 02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/05-exemptions-and-api.md   | none |
| boolean quick reference          | 02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/readme.md      | none |
| boolean flag methods             | 02-spec/02-coding-guidelines/01-cross-language/24-boolean-flag-methods.md      | none             |
| no negatives                     | 02-spec/02-coding-guidelines/01-cross-language/12-no-negatives.md              | none             |
| braces + nesting                 | 02-spec/02-coding-guidelines/01-cross-language/readme.md | none      |
| conditions + extraction (style)  | 02-spec/02-coding-guidelines/01-cross-language/readme.md | none |
| blank lines + spacing            | 02-spec/02-coding-guidelines/01-cross-language/readme.md | none |
| function + type size             | 02-spec/02-coding-guidelines/01-cross-language/readme.md | none  |
| multi-line formatting            | 02-spec/02-coding-guidelines/01-cross-language/readme.md | none   |
| code-style checklist             | 02-spec/02-coding-guidelines/01-cross-language/04-code-style/07-checklist.md   | none             |
| nesting resolution               | 02-spec/02-coding-guidelines/01-cross-language/20-nesting-resolution-patterns.md | none           |
| cyclomatic complexity            | 02-spec/02-coding-guidelines/01-cross-language/06-cyclomatic-complexity.md     | none             |
| code mutation avoidance          | 02-spec/02-coding-guidelines/01-cross-language/18-code-mutation-avoidance.md   | none             |
| strict typing                    | 02-spec/02-coding-guidelines/01-cross-language/13-strict-typing.md             | none             |
| null-pointer safety              | 02-spec/02-coding-guidelines/01-cross-language/19-null-pointer-safety.md       | none             |
| naming + casing (keys)           | 02-spec/02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md     | none             |
| file/folder naming               | 02-spec/02-coding-guidelines/08-file-folder-naming/<language>.md               | none             |
| testing                          | 02-spec/02-coding-guidelines/01-cross-language/readme.md | none             |
| error handling + codes           | 02-spec/03-error-manage/readme.md                   | none             |
| error code registry              | 02-spec/03-error-manage/03-error-code-registry/                                | none             |
| logging + stack traces           | 02-spec/21-app/readme.md             | none             |
| serialization/determinism        | 02-spec/21-app/04-json-contract/                                               | none             |
| ci/cd verification               | 02-spec/12-cicd-pipeline-workflows/02-ci-pipeline.md                           | none             |
| ci guards                        | 02-spec/12-cicd-pipeline-workflows/03-reusable-ci-guards/00-overview.md        | none             |
| contract + e2e testing           | 02-spec/12-cicd-pipeline-workflows/readme.md, 21-e2e-testing-pattern.md | none       |
| static analysis / sarif          | 02-spec/02-coding-guidelines/06-cicd-integration/01-sarif-contract.md          | none             |
```

Consolidated mirrors that MAY be cited as a reading aid, never as the authority:

- `02-spec/17-consolidated-guidelines/05-coding-guidelines.md`
- `02-spec/17-consolidated-guidelines/06-error-management.md`
- `02-spec/17-consolidated-guidelines/18-cicd-pipeline-workflows.md`
- `02-spec/17-consolidated-guidelines/34-compiled-simple-coding-guidelines.md`
- `02-spec/17-consolidated-guidelines/03-strictly-avoid-quickref.md`
- `02-spec/02-coding-guidelines/01-cross-language/readme.md`

When a consolidated mirror and a numbered guideline folder disagree, the numbered
folder wins and the disagreement is filed under RULE 9.

Rules for the checklist:

- Enumerate guideline files by path — no "the guidelines folder".
- If a topic appears in two authoritative files, the step fails: consolidate
  into one file and leave a one-line pointer in the other, then re-run the
  checklist.
- If a topic has no file, file it as a spec gap under RULE 9 — do not invent the
  guideline inside a task.
- This checklist is a spec-planning output, not an audit artifact (RULE 0B).

---

## RULE 0E — capture and install every input (nothing stays only in chat)

Before a single plan line is written, route every user input into its file, then
link that file from the plan's `## Context`. Chat is not storage.

| Input                                                | File                                          |
| ---------------------------------------------------- | --------------------------------------------- |
| Command, convention, "always do X", new CLI          | `.ai-memory/spec/commands/01-<slug>.md`         |
| Bug, regression, broken behavior                    | `.ai-memory/issues/01-<slug>.md`                |
| CI/CD-specific failure                               | `.ai-memory/cicd-issues/01-<slug>.md`           |
| Institutional knowledge (pattern, decision)          | `.ai-memory/memory/` + update `.ai-memory/memory/readme.md` |
| "Never do this again"                                | `.ai-memory/strictly-avoid.md`                  |
| Idea, not yet approved                               | `.ai-memory/suggestions.md`                     |

Create missing folders on demand. Attachments:

- Every attached image or file is REQUIRED input; never leave one only in chat.
- Save verbatim under an `assets/` subfolder next to the file it belongs to;
  project-wide assets go to `.ai-memory/assets/<slug>/` and get a `.ai-memory/memory/readme.md`
  note.
- Names are lowercase-hyphenated with the original extension (RULE 0A).
- The plan carries an `## Attachments` section: one bullet per file with a
  one-line caption stating what the executor must take from it.
- An unreadable or ambiguous attachment becomes an ambiguity file (RULE 0G), and
  the question links the asset.

Plan lifecycle: a new plan is written to `.ai-memory/plans/pending/01-<plan-slug>.md`
with `Status: pending` and a one-line row in `.ai-memory/plans/readme.md`. When done
it is `mv`-ed to `.ai-memory/plans/completed/`, with `Status: completed` flipped in
the same move and the index updated in the same commit. Never copy, never leave a
duplicate across `pending/` and `completed/`, never delete a `pending/` file.

---

## RULE 0F — release policy (this is law, stated in full here)

Individual task runs NEVER release. No version bump, no changelog entry, no
release-notes update, no root `readme.md` version pin on a per-task basis. A run
that touches the version while sibling tasks are pending is auto-reject.

The release fires ONLY when the ENTIRE plan is finished, meaning every task and
subtask of `plan-file` has moved out of `.ai-memory/plans/pending/` into
`.ai-memory/plans/completed/` with `Status: completed`. At that moment, and only
then:

- Bump the MINOR version per the release ceremony in `01-prompts/17-release-management/04-release.md`.
- Add one changelog entry covering the whole plan, never a single task.
- Update release notes.
- Pin the new version in the root `readme.md`.

State this policy verbatim in the plan's `## Context` so the executing run cannot
"forget" it. The final step of the plan MAY be "run the release ceremony" only if
it genuinely is the last step; it never appears earlier and never in a plan that
leaves siblings pending. RULE 8 owns the batch-commit rules and defers to this
rule for releases.

---

## RULE 0G — ambiguity handling (a question mark is a file, not prose)

Ambiguity is not a licence to guess. Any sentence in the plan or in a task file
that ends in a question mark, or that says "assume", "probably", "TBD" or
"we think", MUST become an ambiguity file instead of staying in the body.

- Open: `.ai-memory/ambiguous-questions/01-new-ambiguity/01-<slug>.md`
- Answered: `.ai-memory/ambiguous-questions/02-ambiguity-resolved/01-<slug>.md`

```markdown

# <one-line question>

Slug: <slug>
Status: open
Raised: <YYYY-MM-DD>
Blocking: <plan slug(s) or "none">

## Question

## Options considered

## Impact if guessed wrong

```

When answered: `mv` the file into `02-ambiguity-resolved/`, flip
`Status: resolved`, and append:

```markdown

## Resolution

Answered: <YYYY-MM-DD>
Answer: <user answer>
Applied solution: <what changed / where>
```

Never leave a copy behind. A plan blocked by an open question is still written,
with `Status: blocked-by-ambiguity` and every question file linked from Context.
Each blocked task states the interim provisional default in section 8.

---

## RULE 0H — folder structure and file layout (the literal tree)

RULE 0A gives the naming law; this is the layout it applies to. Authored paths
match this exactly:

```text
02-spec/21-app/readme.md
02-spec/21-app/xx-<plan-slug>.md (or 02-spec/21-app/xx-<plan-slug>/01-overview.md)
.ai-memory/plans/pending/01-<plan-slug>.md
.ai-memory/plans/subtasks/<plan-slug>/index.md
.ai-memory/plans/subtasks/<plan-slug>/001-<subtask-title>.md
.ai-memory/plans/subtasks/<plan-slug>/002-<subtask-title>.md   ... NNN-task.md
.ai-memory/plans/completed/01-<plan-slug>.md
.ai-memory/plans/readme.md
.ai-memory/ambiguous-questions/01-new-ambiguity/01-<slug>.md
.ai-memory/ambiguous-questions/02-ambiguity-resolved/01-<slug>.md
.ai-memory/issues/01-<slug>.md
.ai-memory/cicd-issues/01-<slug>.md
.ai-memory/spec/commands/01-<slug>.md
02-spec/25-app-spec-audit/NN-audit-<yyyy-mm-dd>-v<N>.md
```

- Two-digit prefixes for docs, prompts, spec, audit, memory, issue and ambiguity
  files. Three-digit prefixes for subtask files.
- No unnumbered sibling inside a numbered folder, `index.md` and `readme.md`
  excepted.
- One plan equals one file. A plan never spans two files in `pending/`.

---

## RULE 0I — root-cause analysis for every bug-driven step

A step that exists because something broke carries an RCA, not just a fix.

- The bug is filed first: `.ai-memory/issues/01-<slug>.md`, or
  `.ai-memory/cicd-issues/01-<slug>.md` for pipeline failures (RULE 0E).
- The task body's section 8 links that issue file and the RCA record.
- RCA record shape, mandatory, in the issue file:

```markdown

## Root cause analysis

Symptom: <observed failure, with the command or log line that showed it>
Trigger: <the exact input/commit/config that produced it>
Root cause: <the single mechanism, at symbol + file level>
Why it escaped: <the missing test, guard, or spec rule that let it through>
Fix: <symbol + file changed>
Prevention: <new test name, new lint rule, or new guideline rule + its file>
Regression check: <runnable command + expected output>
```

- Follow the retrospective and verification conventions in
  `02-spec/03-error-manage/01-error-resolution/03-retrospectives/` and
  `02-spec/03-error-manage/01-error-resolution/04-verification-patterns/`; pipeline
  RCAs follow `02-spec/12-cicd-pipeline-workflows/17-release-pipeline-issues-rca.md`.
- "Prevention: be more careful" is a failed RCA. Prevention names a test, a
  guard, or a rule file.
- No fix step ships without its `Prevention` and `Regression check` lines filled.

---

## RULE 0J — reference integrity and CI/CD verification

Two mechanical gates, both run with output pasted into the report (RULE 11).

Missing links and references. Report absolute counts, not adjectives:

```bash

# every cited spec/.ai-memory path in the batch must resolve

rg -o --no-filename '(spec|\.ai-memory)/[A-Za-z0-9/._-]+'   .ai-memory/plans/pending/01-xx-plan-slug.md .ai-memory/plans/subtasks/xx-plan-slug/*.md   | sed 's/[.,)`]*$//' | sort -u > /tmp/cited-paths.txt
wc -l < /tmp/cited-paths.txt                       # total citations
while read p; do test -e "$p" || echo "MISSING $p"; done < /tmp/cited-paths.txt | tee /tmp/missing.txt
wc -l < /tmp/missing.txt                           # must be 0
python3 linter-scripts/check-spec-folder-refs.py   # folder-level ref integrity
```

Acceptance: `missing paths = 0`, `missing sections = 0`, and every cited section
anchor exists inside the cited file. The report prints the four numbers:
citations total, missing files, missing sections, unreferenced-but-required
guideline files. A non-zero count blocks the save.

CI/CD verification. Every plan states how CI proves it, and every task that
changes code names the pipeline check that guards it:

- Pipeline definition: `02-spec/12-cicd-pipeline-workflows/02-ci-pipeline.md`.
- Reusable guards: `02-spec/12-cicd-pipeline-workflows/03-reusable-ci-guards/`
  (`01-forbidden-name-guard.md`, `04-baseline-diff-lint-gate.md`,
  `06-matrix-test-aggregator.md`).
- Contract and end-to-end layers:
  `02-spec/12-cicd-pipeline-workflows/readme.md` and
  `21-e2e-testing-pattern.md`.
- Local mirrors of the CI gates: `03-ai-scripts/06-cicd-local-runner.py`, plus the
  specific `linter-scripts/check-*.py|mjs` scripts the task can break.
- The plan carries a `## CI/CD verification` section mapping each domain to the
  jobs that must be green before the batch commits. A task whose verification
  command has no CI counterpart is an incomplete task.

---

## RULE 1 — working stance (read this before you write anything)

The AI running this prompt has been a lazy fraud on prior runs of this exact job:

- Generated 69 task files whose bodies were byte-identical after the title line,
  then reported it as done. The audit scored that batch 28 out of 100.
- Wrote phase-conditional mush — "If this is the Scaffold phase, create the
  files; if this is the Implement phase, write the business logic" — which is a
  body that fits any task and therefore teaches nothing.
- Pasted the identical `/learn` list into every single file as anchoring theatre.
- Put a commit-and-push block in every task, inviting 69 commits for one feature.
- Cited spec paths that did not exist, and sections that did not exist inside
  files that did.
- Invented contracts (field names, status codes, table columns) where the spec
  was silent instead of filing the ambiguity, laundering a guess into code.
- Wrote "implemented correctly" and "works as expected" as verification steps.
- Wrote "the relevant file" and "the model layer" instead of a path.
- Derived dependencies from file numbering instead of data flow, producing a
  graph that deadlocked on itself.

Every one of those is a hard failure here. A blind AI must be able to open one
task file, read nothing else except the files that task cites, and produce
correct code. If the file does not carry that much, you have not done the work.
Aggressive enforcement is intentional. Do not soften it.

---

## RULE 2 — how the step list is derived

1. Inventory the buildable units from the spec, in data-flow order: shared
   contract types → producer → transport → consumer intake → consumer writes →
   UI → CI. Name every unit before assigning any step numbers.
2. Split each unit into phases only when each phase owns distinct symbols and a
   distinct runnable check. The usual split is `Scaffold` (types, files,
   signatures, no behavior), `Implement` (behavior plus unit tests),
   `Wire+Test` (integration into the caller plus an end-to-end check).
3. Units × phases must total `n`. If the total overshoots, merge phases in the
   least risky units. If it undershoots, split the riskiest units further —
   never pad.
4. Order steps so a step is never scheduled before the step producing its input.

---

## RULE 3 — one step, one buildable unit, one domain

| Requirement            | Test                                                                        |
| ---------------------- | --------------------------------------------------------------------------- |
| One buildable unit     | Two target modules that can ship independently means two steps              |
| One domain             | Exactly one value from `domains` — never two in one step                    |
| One observable outcome | You can state what a reviewer sees that they could not see before           |
| Named symbols          | The step already knows the struct/class/function/component names it creates |

If a step cannot name its symbols, the spec has not decided them. That is a spec
gap — handle it under RULE 9. Do not write the step vaguely and move on.

---

## RULE 4 — the eight-section task file (mandatory shape, this prompt owns it)

Every file in `.ai-memory/plans/subtasks/xx-plan-slug/` uses exactly this shape, in
this order, named `NNN-task.md` with a zero-padded three-digit sequence:

```text
---

plan: .ai-memory/plans/pending/01-xx-plan-slug.md
domain: <one of domains>
phase: Scaffold | Implement | Wire+Test
target_files: [<exact repo-relative paths>]
depends_on: [<Task NNN>]
citations:
  app_spec: "02-spec/21-app/... §<section>"
  canonical_size: "02-spec/02-coding-guidelines/02-canonical-size-tier.md"
  language_guideline: "<path(s)>"
  boolean_styling: "02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/<file>.md"
  folder_naming: "02-spec/02-coding-guidelines/08-file-folder-naming/<file>.md"
  error_architecture: "02-spec/03-error-manage/readme.md"
  error_codes: "02-spec/21-app/readme.md"
  logging_traces: "02-spec/21-app/readme.md"
  response_envelope: "02-spec/21-app/readme.md"
  golden_fixture: "02-spec/21-app/fixtures/<file>.example.json"
  strictly_avoid: ".ai-memory/strictly-avoid.md"
  database: "02-spec/04-database-conventions/..."
  ui_surface: "n/a"
  tests: "unit <name>"
  ci_cd_guard: "linter-scripts/check-*"
  ambiguity: "n/a"
  issue_rca: "n/a"
---

# Task NNN — <specific outcome, not a phase name>

## 1. Learn

## 2. Goal

## 3. Inputs and Contracts

## 4. Execute

## 5. Constraints

## 6. Verify

## 7. Done When

## 8. Notes and Open Questions

```

Acceptance bar per section — a section that misses its bar fails the batch:

| Section                     | Bar                                                                                 |
| --------------------------- | ----------------------------------------------------------------------------------- |
| 1. Learn                    | 3-7 links, at least half unique to this task, each with a one-line "why read this"  |
| 2. Goal                     | 2-4 sentences of behavior and blast radius. Never a restated title                 |
| 3. Inputs and Contracts     | Types consumed and produced, the wire shape inlined verbatim, error codes in scope |
| 4. Execute                  | Ordered steps; every step names a symbol AND the file it lands in                   |
| 5. Constraints              | 3-6 rules this task could actually violate, each with a rule id and source file     |
| 6. Verify                   | At least one runnable command plus the expected output                              |
| 7. Done When                | 3-6 binary statements, each citing a numbered acceptance criterion                  |
| 8. Notes and Open Questions | Ambiguity file path plus the interim default, or the single word `None.`            |

---

## RULE 5 — the citation table (the critical requirement)

A task that cites two files is a task the executor will improvise around. Every
task file MUST carry all twelve of the following, resolved to real paths and, for
spec files, a real section anchor. Spread them across sections 1, 3, 5 and 6 —
but all twelve must be present somewhere in the file.

| #   | Citation                           | Where it comes from                                                                                                                                                                                                                                                                                                                                                                                                          |
| --- | ---------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | Deciding app-spec section          | `02-spec/21-app/...` file plus section, e.g. `04-json-contract/02-section-and-asset-schema.md §Section`                                                                                                                                                                                                                                                                                                                         |
| 2   | Canonical size tier                | `02-spec/02-coding-guidelines/02-canonical-size-tier.md`                                                                                                                                                                                                                                                                                                                                                                        |
| 3   | Language guideline for this domain | Go: `02-spec/02-coding-guidelines/01-cross-language/readme.md` + the specific rule file (`02-boolean-standards.md`, `05-defer-rules.md`, `09-wrapped-boolean-results.md`). PHP: `04-php/00-overview.md` + `02-forbidden-patterns.md`, `03-naming-conventions.md`, `05-response-array-standard.md`. TS/React: `02-typescript/08-typescript-standards-reference.md` + `12-discriminated-union-patterns.md`, `14-state-management.md` |
| 4   | File and folder naming             | `02-spec/02-coding-guidelines/08-file-folder-naming/` — `03-golang.md`, `02-php-wordpress.md`, or `04-typescript-javascript.md` for the language this task writes                                                                                                                                                                                                                                                               |
| 5   | Error architecture                 | `02-spec/03-error-manage/readme.md` plus `06-apperror-package/` (Go) or `05-response-envelope/` (transport)                                                                                                                                                                                                                                                                                          |
| 6   | Error code range                   | `02-spec/21-app/readme.md` — quote the symbolic and numeric codes this task may emit                                                                                                                                                                                                                                                                                                     |
| 7   | Logging and stack traces           | `02-spec/21-app/readme.md` for anything that can fail                                                                                                                                                                                                                                                                                                                                 |
| 8   | Response envelope                  | `02-spec/21-app/readme.md` for any REST or CLI output surface                                                                                                                                                                                                                                                                                                                                |
| 9   | Golden fixture                     | `02-spec/21-app/fixtures/<file>.example.json` plus `02-spec/21-app/readme.md` for any task touching a wire format                                                                                                                                                                                                                                                                                                    |
| 10  | Strictly-avoid rules               | `.ai-memory/strictly-avoid.md` — name the specific rules this task could break, not the whole file                                                                                                                                                                                                                                                                                                                             |
| 11  | Exact target files                 | Repo-relative paths, created or edited, in the header                                                                                                                                                                                                                                                                                                                                                                        |
| 12  | Exact symbols                      | Struct / class / function / component names with signatures or field lists                                                                                                                                                                                                                                                                                                                                                   |

Plus: when the spec is silent, cite the ambiguity file
(`.ai-memory/ambiguous-questions/01-new-ambiguity/NN-<slug>.md`) and the interim
default. When a task has a database surface, also cite
`02-spec/04-database-conventions/` and the owning table section in `02-spec/23-app-db/`.

These citations must be task-specific. A task whose twelve citations match its
neighbour's line for line fails the clone gate in RULE 6. Cite the section that
decides _this_ task's behavior, not the folder overview as a reflex.

### RULE 5A — the per-task frontmatter citations (must follow, non-negotiable)

Every task file MUST begin with the YAML frontmatter block described in RULE 4. This replaces the old markdown checklist. Missing citations or empty keys fail the batch. This is the file list the executor is guaranteed to have; nothing else may be assumed.

- The `citations` block in the YAML frontmatter must include all keys shown in the RULE 4 shape.
- Values must be exact repo-relative paths (and for spec files, include the section anchor).
- Keys that genuinely do not apply are written as `"n/a — <reason>"`. A blank value, a folder without a file, or the word "relevant" fails the batch.

---

## RULE 6 — the anti-clone gate is run, not promised

Before the batch is saved, run both sweeps and paste the real output into the
report:

```bash

# 1. uniform length is the first smell

wc -l .ai-memory/plans/subtasks/xx-plan-slug/*.md | sort -n | head

# 2. strip title + header lines, hash the remainder; any bucket > 1 is a FAIL

for f in .ai-memory/plans/subtasks/xx-plan-slug/*.md; do
  tail -n +4 "$f" | rg -v '^\*\*(Plan|Domain|Target Files|Depends On)'     | sha256sum | cut -c1-12 | tr '
' ' '; echo "$f"
done | sort | uniq -c -w12 | sort -rn | head
```

Gate: zero hash buckets with more than one member, and no pair above 60 percent
similarity after the title and header lines are removed. Also sweep that every
cited path resolves:

```bash
rg -o --no-filename `02-spec/[A-Za-z0-9/._-]+' .ai-memory/plans/subtasks/xx-plan-slug/*.md   | sed 's/[.,)`]*$//' | sort -u | while read p; do test -e "$p" || echo "MISSING $p"; done
```

---

## RULE 7 — banned phrasing, with the rewrite

| Banned                                       | Why                                            | Rewrite                                                                 |
| -------------------------------------------- | ---------------------------------------------- | ----------------------------------------------------------------------- |
| "If this is the Scaffold phase…"             | Body works for any task, so it teaches nothing | Name the structs this phase declares                                    |
| "Write the business logic"                   | Zero information                               | "Implement `func (q *Frontier) Enqueue(e Entry) appfault.Result[bool]`" |
| "Step 4 of 69"                               | Ordering is not content                        | State the outcome                                                       |
| "the relevant file", "the model layer"       | Unresolvable                                   | Exact repo-relative path                                                |
| "implemented correctly", "works as expected" | Unverifiable                                   | A command plus expected output                                          |
| "Review the code"                            | Not a check                                    | A test name or linter invocation                                        |
| "max 80-100 lines per function"              | Contradicts the canonical tier                 | Cite `02-spec/02-coding-guidelines/02-canonical-size-tier.md`              |
| The same `/learn` list in every task         | Anchoring theatre                              | 3-7 links chosen for this task                                          |
| A commit or release block in a task          | Invites one commit per task                    | See RULE 8                                                              |

---

## RULE 8 — commits and releases are batch-level (stated in full here)

No task file contains a commit, push, tag, or release instruction. The policy for
this batch, in full:

- One commit per batch of completed tasks, never per task. The message names
  the plan and the task range, e.g. `feat(cli): tasks 012-018 frontier + parser`.
- The commit happens after the batch's verification commands pass, not before.
- Tasks move to `.ai-memory/plans/completed/` in the same commit that lands their
  code, and `.ai-memory/plans/readme.md` is updated in that same commit.
- A release (tag, changelog entry, artifact build) fires only when every task
  of `plan-file` has moved to `.ai-memory/plans/completed/` with
  `Status: completed`. Releasing mid-plan is forbidden.
- Never rewrite history, never force-push, never commit build artifacts.

---

## RULE 9 — spec gaps are filed, never invented

If a step cannot be written because the spec does not decide something:

1. File `.ai-memory/ambiguous-questions/01-new-ambiguity/NN-<slug>.md` with the
   question, the options, and the cost of each.
2. Reference that file from the task's section 8.
3. State the interim default the executor must use, and mark it as provisional.

Writing an invented contract into a task file launders a guess into an
implementation. It is the most expensive defect class in this project.

---

## RULE 10 — dependencies come from data flow

`Depends On` is derived from what a task consumes, never from file numbering.
Contract types → producer emit → transport bundle → consumer intake → consumer
writes → UI → CI. The graph MUST be acyclic; state the check you ran and the
number of roots. A task that unblocks many others is scheduled early and says so
in section 2.

---

## RULE 11 — pre-save checklist (tick every line, in the report)

- [ ] `n`, `plan-slug`, and `plan-file` were supplied, not guessed.
- [ ] Step count equals `n`, counted forwards and backwards.
- [ ] Every step names one unit, one domain from `domains`, one outcome.
- [ ] Every task file has all eight sections of RULE 4, in order, and meets each bar.
- [ ] Every task carries all twelve citations of RULE 5, task-specific.
- [ ] Every `## 1. Learn` list is 3-7 items, half unique, each with a "why".
- [ ] Every `## 4. Execute` step names a symbol and a file.
- [ ] Every task has at least one runnable command in `## 6. Verify` with expected output.
- [ ] Every `## 7. Done When` item is binary and cites an acceptance-criteria number.
- [ ] Every cited path resolves — file exists AND the section exists inside it.
- [ ] Clone sweep shows zero multi-member buckets; max pairwise similarity under 60 percent.
- [ ] Dependency graph is acyclic.
- [ ] No banned phrasing from RULE 7 anywhere in the batch.
- [ ] No task file contains a commit, push, tag, or release instruction.
- [ ] Plan task table, `.ai-memory/plans/readme.md`, and the memory index updated.
- [ ] RULE 0A naming check ran; output printed; zero uppercase/space/underscore paths, every ordered file carries its zero-padded prefix.
- [ ] RULE 0C respected: audit slots exist with correct names and empty bodies; no audit was scored during authoring; no task contains an audit step.
- [ ] RULE 0D footer present verbatim in every task file; plan file states one step per run and self-loop; no task implies batching steps.
- [ ] RULE 0C ceilings stated: max 2 agents, max 3 threads per agent, nowhere exceeded or contradicted.
- [ ] Consolidated Coding Guidelines (Mandatory): I have explicitly cited the master consolidated coding guideline file at `.ai-memory/coding-guidelines.md` in the plan and marked it as mandatory reading for all execution agents.
- [ ] RULE 0E coding-guideline single-file checklist filled in with real paths, every boolean/condition-styling sub-file listed individually, every topic single-source, duplicates column reads "none".
- [ ] Anti-Garbage Naming (Non-Negotiable): I have strictly verified that absolutely NO generic garbage variable names (e.g., `comp_100.go`, `temp`, `data`, `obj`, `Input100`, `TestHandleComp100`) were written. All names are highly semantic and domain-specific.
- [ ] Execution model confirmed: exactly one step per run, self-loop after Verify, no line anywhere in the plan or a task implies batching two steps.
- [ ] Self-loop instruction present in the plan file AND verbatim in every task footer; max 2 agents / max 3 threads per agent restated and nowhere exceeded.
- [ ] RULE 0E applied: every command, issue, cicd-issue, memory item, strictly-avoid entry and attachment written to its file and linked from Context; nothing left only in chat; every attachment has a caption.
- [ ] RULE 0F release policy quoted in the plan's Context; no task touches version, changelog, release notes, or the readme version pin.
- [ ] RULE 0G respected: zero question marks, "assume", "probably" or "TBD" left in plan or task bodies; each became an ambiguity file with an interim default.
- [ ] RULE 0H layout matched exactly; `.ai-memory/plans/readme.md` row present; no unnumbered sibling in a numbered folder.
- [ ] RULE 0I: every bug-driven step links an issue file with a complete RCA record, including `Prevention` and `Regression check`.
- [ ] RULE 0J counts printed: citations total, missing files = 0, missing sections = 0, unreferenced required guideline files = 0; `check-spec-folder-refs.py` output pasted.
- [ ] RULE 0J: plan has a `## CI/CD verification` section, and every code task names the CI job or linter script that guards it.
- [ ] RULE 5A per-task YAML frontmatter citations present and filled in every task file; no blank values, no "relevant file", `n/a` values carry a reason.

---

## RULE 12 — Subagent Validation & Self-Score

Before printing the final report, you MUST spawn a read-only Validation Subagent to independently verify your batch.

- Pass the Validation Subagent the path to the plan folder (`.ai-memory/plans/subtasks/xx-plan-slug/`).
- Ask the Subagent to compute the scores for the dimensions below based on the generated files.
- You must use the Subagent's exact output for your final report, rather than scoring your own work.

A score without evidence is not a score.

```text
| Dimension          | Score | Evidence |
| Uniqueness         |    ?? | clone buckets: 0; max pair similarity: 41% |
| Specificity        |    ?? | 69/69 tasks name >= 3 symbols with signatures |
| Anchoring          |    ?? | 0 dead paths; 12/12 citations present in 69/69 tasks |
| Reference integrity|    ?? | citations: 412; missing files: 0; missing sections: 0 |
| Verifiability      |    ?? | 69/69 tasks carry a runnable command with expected output |
| Ci coverage        |    ?? | 69/69 code tasks name a CI job or linter script |
| Sequencing         |    ?? | acyclic; 12 roots; longest chain 6 |
| Overall            |    ?? |  |
```

Below `min-score` overall: regenerate the offending domain group before saving.
Below `regen-score`: regenerate the whole batch from scratch. Never patch a
cloned batch — clones are regenerated, not edited.

---

## RULE 13 — final report format

```text
Plan: 01-xx-plan-slug.md                 Steps: n (counted twice)
Subtasks: <count> in .ai-memory/plans/subtasks/xx-plan-slug/
Clone buckets > 1: 0            Max similarity: NN%
Citations total: NNN            Missing files: 0    Missing sections: 0
Guideline checklist: NN topics, duplicates: 0
Execution model: one step per run, self-loop, max 2 agents x 3 threads
Ci/cd verification: <jobs per domain>
Tasks blocked on ambiguity: n (Q..)   Rca records written: n
Self-score: Uniqueness/Specificity/Anchoring/Reference/Verifiability/Ci/Sequencing -> Overall
Spec gaps filed: <list or none>
```

### Temp-Agent Isolated Task Directory & Communication Protocol (Non-Negotiable)

To prevent cross-task pollution and ensure seamless agent communication, every task MUST create a dedicated subfolder in `.ai-memory/temp-agents/xx-<task-name>/`:

1. Per-Task Isolation: On task start, the assigned subagent creates its isolated directory `.ai-memory/temp-agents/xx-<task-name>/`.
2. State & Progress Tracking: Create `.ai-memory/temp-agents/xx-<task-name>/state.md` documenting:
   - `TASK_NAME`: `<task-name>`
   - `STATUS`: `IN_PROGRESS` | `DONE` | `FAILED`
   - `ASSIGNED_AGENT`: Agent identifier and thread index
   - `CURRENT_STEP`: Detailed micro-step description
3. Inter-Agent Communication & Handoff:
   - All intermediate findings, scratch outputs, and dependency handoffs between agents working on this task MUST be written inside `.ai-memory/temp-agents/xx-<task-name>/`.
   - Sibling or successor agents MUST inspect this dedicated folder before resuming work or fixing errors.
4. On Error/Crash: Append the exact error, root cause, and `STATUS: FAILED` to `.ai-memory/temp-agents/xx-<task-name>/state.md` before exiting.
5. On Success: Mark `STATUS: DONE` in `.ai-memory/temp-agents/xx-<task-name>/state.md`, aggregate findings to the master plan, and clean up or archive the folder.

