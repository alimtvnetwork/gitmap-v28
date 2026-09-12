# Application Specification Blind-AI Audit — Quality Protocol (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

## Variables — Auto-Discovered at Runtime

Do NOT ask the user to provide these variables. You must discover them automatically during execution:

``text
N              = 200 (Default number of self-loops. The user may override this when triggering the prompt)
audit-date     = <Current system date: YYYY-MM-DD>
audit-time     = <Current system time: HH:MM:SS>
audit-version  = <Auto-incremented run number. Check `spec/25-app-spec-audit/`. If a previous audit exists, increment the version. If NOT, you MUST execute `rm -rf spec/25-app-spec-audit/*` to clear the folder completely, then start at `v1`.>
audit-file     = spec/25-app-spec-audit/NN-audit-<audit-date>-v<audit-version>.md
scope          = spec/*-app* | spec/*-design* | .lovable/plans/pending | <auto-discover any recently modified spec folders>
min-score      = 100
``

`NN` in `audit-file` is the next free two-digit prefix in `spec/25-app-spec-audit/`.
Discover the correct `audit-version` by checking existing files on disk. Do not prompt the user for these values. Non-negotiable checkpoint: If this is the very first audit (no previous versions found), you must remove the audit folder contents before generating the new one.

Trigger phrases: "audit the app spec", "blind ai readiness audit", "score the
spec", "run audit v<N>".

---

## RULE 0 — what this prompt does, and what it must never do

- It reads the spec and writes one dated audit file. That is all.
- It never edits, fixes, renames, splits or reformats a file it audits. Audit,
  then fix in a separate run. Auditing and fixing in the same pass is a hard
  failure (`spec/01-spec-authoring-guide/01-index.md` §4 AUD-004g).
- It never writes application code, never runs a migration, never commits.
- It never edits an earlier dated audit file. A new run is a new file; the
  superseded run gets a `> STALE — superseded by <audit-file>` banner as its
  only permitted modification.
- It never rubber-stamps. An audit that finds nothing must list the checks that
  ran and passed, with their output.

---

## RULE 0K - strict anti-garbage naming audit

You must rigorously audit all proposed files, variables, and unit tests for garbage naming conventions.

- If a spec proposes tests like `TestHandleComp100` or arbitrary generic IDs instead of semantic domain behaviors (e.g., `TestUpdateUser_RejectsInvalidEmail`), you must flag it as a critical failure.
- If you find generic variable names like `data`, `obj`, `temp`, or `Input100`, record it as a defect and demand a remedy.

## RULE 0L - temporary automation scripts must not be committed

If a spec proposes creating temporary scripts (e.g., CSJ, Python) for fixing or refactoring the codebase:

- It must explicitly specify that helper scripts go into `03-ai-scripts/` and MUST be committed.
- It must explicitly specify that `03-ai-scripts/01-index.md` is updated and linked in `what-to-read.md`.
- If the spec implies committing these temporary files to the repository, you must flag it as a critical failure task.

---

## RULE 0 — Temp Script Sandboxing (Global Law)

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .lovable/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Native File Manipulator:** If you need to perform mass file renaming, `.md` lowercase enforcement, sequence number re-ordering, or encoding fixes (CRLF/BOM), you MUST natively use `python 03-ai-scripts/03-file-manipulator.py <command>` rather than writing a new script from scratch.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/01-index.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

---

## RULE 1 - working stance

Read as the blind-AI persona in `spec/01-spec-authoring-guide/01-index.md` §1: never
asks a question, takes the first matching rule, treats SHOULD as optional, cannot
infer intent, trusts diagrams over prose, has only the delivered folder, and stops
at the first heading that looks like an answer. Scoring the spec as a cooperative
reader would invalidates the audit.

Prior runs of this job failed by: asserting "missing" without a search that would
have found it, copying findings forward without re-verifying them, writing
findings with no remedy, and scoring prose quality instead of implementability.
Do not repeat any of it.

---

## RULE 2 — STEP 1 IS THE FILE INVENTORY (nothing before it)

Before reading content, before forming an opinion, print the inventory. No score
may be written until this table exists in `audit-file`.

``bash

# 1. the audited scope, with line counts

# Be sure to include any recent folders and files written for recent specs in these commands!

wc -l spec/21-app/*.md spec/21-app/*/*.md spec/21-app/*/*/*.md 2>/dev/null | sort -n
wc -l spec/23-app-db/*.md spec/24-app-ui-design-system/*.md 2>/dev/null | sort -n
ls spec/21-app/fixtures/

# 2. the guideline and support folders the spec must bind to

ls spec/02-coding-guidelines spec/02-coding-guidelines/01-cross-language \
   spec/02-coding-guidelines/01-cross-language/02-boolean-principles \
   spec/02-coding-guidelines/01-cross-language/04-code-style \
   spec/02-coding-guidelines/08-file-folder-naming
ls spec/03-error-manage spec/03-error-manage/02-error-architecture spec/03-error-manage/03-error-code-registry
ls spec/04-database-conventions spec/12-cicd-pipeline-workflows spec/12-cicd-pipeline-workflows/03-reusable-ci-guards
ls spec/17-consolidated-guidelines

# 3. the plan surface that consumes the spec

ls .lovable/plans/pending .lovable/plans/subtasks/*/ .lovable/ambiguous-questions/01-new-ambiguity
``

Inventory table shape, written into `audit-file` as section 1:

``text
| # | Path | Lines | Role (normative / index / fixture / diagram / mirror) | Read? |
``

Rules:

- Every file in `scope` appears in the table. A file omitted from the table is an
  audit defect, not a spec defect.
- Print totals: file count, total lines, count of normative files, count of
  fixtures, count of files over 300 lines.
- Then state the read order you will follow, overview files first.

Cross-Platform Execution: If bash or rg commands are unavailable in the host environment, you are explicitly authorized to use native equivalents (e.g., PowerShell Get-ChildItem / Measure-Object or a Python glob/os script) to generate the exact required metrics. The output format must still precisely match the required table shapes.

---

## RULE 3 — the scored dimensions (0-100 each, evidence per row)

Overall Score Math: The Overall score must be calculated as the strict arithmetic mean of the 12 dimension scores, rounded down to the nearest integer. If any single dimension scores below 50, the Overall Score is capped at 50 (automatic failure).

Each dimension gets a score, the evidence that produced it, and at least one
remedy row in the improvement set. Point costs come from
`spec/01-spec-authoring-guide/01-index.md`.

| #   | Dimension                      | The question it answers                                                                                          |
| --- | ------------------------------ | ---------------------------------------------------------------------------------------------------------------- |
| 1   | Blind-AI readiness             | Can a blind implementer build the system from these files alone, without inventing behavior?                     |
| 2   | Code-file coverage             | Does the spec name the exact repo-relative files to create or modify, per unit, in the target language?            |
| 3   | Coding-guideline checklist     | Does every unit bind to exactly one authoritative guideline file per topic, and do the bindings resolve?           |
| 4   | Code-mutation discipline       | Is mutation forbidden where it must be, per `01-cross-language/18-code-mutation-avoidance.md`, and are contracts immutable? |
| 5   | Test specification             | Are unit, integration and end-to-end tests specified per unit, with names and fixtures?                           |
| 6   | Acceptance criteria            | Does every normative file close with criteria that name a verification method?                                    |
| 7   | Ambiguity discipline           | Is every undecided thing filed as an open question, and is nothing guessed past in the spec text?                 |
| 8   | Cross-folder consistency       | Do `21-app`, `23-app-db`, `24-app-ui-design-system` and `fixtures/` agree on names, types and flows?              |
| 9   | Reference integrity            | How many links and references are missing? Absolute counts, both directions.                                      |
| 10  | Ci/cd verifiability            | Is there a named pipeline job or guard for every buildable unit and every wire contract?                          |
| 11  | Shape and size                 | Do files satisfy the mechanical checks in `02-file-size-and-shape-checks.md`?                                     |
| 12  | Determinism                    | Per file, is every coder decision `Fixed`, or is it `Open` / `Absent`?                                            |

### Dimension 3 — the coding-guideline checklist (mandatory table)

Rebuild it from the filesystem, not from what the spec claims. Every boolean and
condition-styling sub-file is listed individually:

``text
| Topic                           | Authority file                                                                 | Bound from spec/21-app? | Duplicates |
| canonical size tier             | spec/02-coding-guidelines/02-canonical-size-tier.md                            | yes/no                  | none       |
| boolean naming prefixes         | .../01-cross-language/02-boolean-principles/01-naming-prefixes.md              | yes/no                  | none       |
| boolean guards + extraction     | .../02-boolean-principles/02-guards-and-extraction.md                          | yes/no                  | none       |
| boolean params + conditions     | .../02-boolean-principles/03-parameters-and-conditions.md                      | yes/no                  | none       |
| boolean quick reference         | .../02-boolean-principles/04-quick-reference.md                                | yes/no                  | none       |
| boolean exemptions + api        | .../02-boolean-principles/05-exemptions-and-api.md                             | yes/no                  | none       |
| boolean flag methods            | .../01-cross-language/24-boolean-flag-methods.md                               | yes/no                  | none       |
| no negatives                    | .../01-cross-language/12-no-negatives.md                                       | yes/no                  | none       |
| braces + nesting                | .../01-cross-language/04-code-style/01-braces-and-nesting.md                   | yes/no                  | none       |
| conditions + extraction (style) | .../04-code-style/02-conditions-and-extraction.md                              | yes/no                  | none       |
| blank lines + spacing           | .../04-code-style/03-blank-lines-and-spacing.md                                | yes/no                  | none       |
| function + type size            | .../04-code-style/04-function-and-type-size.md                                 | yes/no                  | none       |
| multi-line formatting           | .../04-code-style/05-multi-line-formatting.md                                  | yes/no                  | none       |
| code-style checklist            | .../04-code-style/07-checklist.md                                              | yes/no                  | none       |
| nesting resolution              | .../01-cross-language/20-nesting-resolution-patterns.md                        | yes/no                  | none       |
| cyclomatic complexity           | .../01-cross-language/06-cyclomatic-complexity.md                              | yes/no                  | none       |
| code mutation avoidance         | .../01-cross-language/18-code-mutation-avoidance.md                            | yes/no                  | none       |
| strict typing                   | .../01-cross-language/13-strict-typing.md                                      | yes/no                  | none       |
| null-pointer safety             | .../01-cross-language/19-null-pointer-safety.md                                | yes/no                  | none       |
| key naming pascalcase           | .../01-cross-language/11-key-naming-pascalcase.md                              | yes/no                  | none       |
| test naming + structure         | .../01-cross-language/14-test-naming-and-03-structure.md                          | yes/no                  | none       |
| file/folder naming              | spec/02-coding-guidelines/08-file-folder-naming/<language>.md                  | yes/no                  | none       |
| language rules (go/php/ts)      | spec/02-coding-guidelines/03-golang|04-php|02-typescript/...                   | yes/no                  | none       |
| error architecture              | spec/03-error-manage/01-index.md                      | yes/no                  | none       |
| error code registry             | spec/03-error-manage/03-error-code-registry/                                   | yes/no                  | none       |
| database conventions            | spec/04-database-conventions/                                                  | yes/no                  | none       |
| ci pipeline + guards            | spec/12-cicd-pipeline-workflows/02-ci-pipeline.md, 03-reusable-ci-guards/      | yes/no                  | none       |
``

Consolidated mirrors under `spec/17-consolidated-guidelines/` (notably
`05-coding-guidelines.md`, `06-error-management.md`,
`18-cicd-pipeline-workflows.md`, `34-compiled-simple-coding-guidelines.md`,
`03-strictly-avoid-quickref.md`) and
`spec/02-coding-guidelines/01-cross-language/01-index.md` are checked for
drift against their authority. A mirror that contradicts its authority is a
Consistency finding; the mirror is never treated as the authority. Any topic
appearing in two authoritative files is a duplicate-authority finding, cost per
the rubric, with a consolidation remedy naming which file wins.

### Dimension 9 — reference integrity (report counts, not adjectives)

``bash
python3 linter-scripts/check-spec-folder-refs.py
python3 linter-scripts/check-file-sizes.py

# every relative link in scope must resolve

rg -o --no-filename '\]\(([^)]+)\)' spec/21-app spec/23-app-db spec/24-app-ui-design-system \
  | sed -E 's/^\]\(//; s/\)$//' | sort -u > /tmp/links.txt
wc -l < /tmp/links.txt

# index vs filesystem, both directions

rg -n '\| *[0-9]{2} *\|' spec/21-app/01-index.md
``

The audit prints exactly these numbers:

``text
| Metric                                          | Count |
| relative links found                            |    NN |
| links that do not resolve                       |     0 |
| cited sections that do not exist in their file   |     0 |
| files present on disk but missing from an index  |     0 |
| files listed in an index but missing on disk     |     0 |
| guideline topics with no authority file          |     0 |
| guideline topics with two authority files        |     0 |
| files over 300 lines                             |     0 |
``

Any non-zero row is a finding with a remedy. An index entry missing a file and a
file missing from the index are two separate findings.

### Dimension 10 — ci/cd verifiability

For every buildable unit, name the pipeline job or guard that proves it, sourced
from `spec/12-cicd-pipeline-workflows/02-ci-pipeline.md`,
`03-reusable-ci-guards/`, `20-contract-testing.md`, `21-e2e-testing-pattern.md`,
and the local mirrors in `linter-scripts/`. A unit with no named check is a
Testability finding. A spec that mentions "CI will catch it" without naming the
job is the same finding.

---

## RULE 4 — phase order (skipping a phase invalidates the audit)

Follow `spec/01-spec-authoring-guide/01-index.md` §1, with the inventory as
phase 1:

1. File inventory (RULE 2) and scope declaration.
2. Mechanical sweep — commands from `02-file-size-and-shape-checks.md` §6, output
   pasted verbatim.
3. Unit inventory — every buildable unit mapped to the file that specifies it,
   diffed against `.lovable/plans/subtasks/*/` target-file lists.
4. Determinism read — per file, each coder decision marked
   `Fixed` / `Defaulted` / `Open` / `Absent`.
5. Consistency map — concern to authority, rebuilt from files, diffed against each
   folder's `99-consistency-report.md`.
6. Guideline checklist (dimension 3) and mirror-drift check.
7. Test and acceptance-criteria pass (dimensions 5, 6).
8. Reference integrity and ci/cd verifiability (dimensions 9, 10).
9. Blind-buildability trace of the primary flow, per `03-audit-procedure.md` §4.
10. Scoring with the arithmetic shown, then the improvement set.

Phases 2 and 3 may run in parallel; nothing else may.

---

## RULE 5 — evidence rules

- Every finding cites `path:line` or `path §section`.
- Command output used as evidence is pasted verbatim, never summarised.
- No finding may assert absence without including the search that would have
  found it.
- A finding carried forward from a previous audit is re-verified against the
  current file and dispositioned as `Closed`, `Partially closed`, `Open`,
  `Reclassified` or `False positive`. Silently dropping a finding invalidates the
  run.
- Every finding carries a remedy naming the file to create or edit and the content
  it needs. A finding without a remedy is incomplete.

---

## RULE 6 — improvement set format

``text
| Rank | Finding | Remedy (file + content) | Dimension | Points | Effort |
``

Effort is `S` (under an hour), `M` (a session), `L` (multi-session). Rank by
points recovered per unit of effort, highest first. Every dimension scoring below
`min-score` contributes at least one row.

Remedy Consolidation: If mechanical sweeps yield more than 20 findings of the exact same category (e.g., 50 missing reference links), consolidate them into a single high-value Improvement Set row. Point the remedy to a newly generated script or temporary checklist file in scratch/ that lists all fixes, rather than bloating the main audit table.

---

## RULE 7 — output file shape

`audit-file` uses lowercase-hyphenated naming with a two-digit prefix and this
section order:

``text

# Audit <audit-date> v<audit-version> — <scope>

Version: 1.0.0
Updated: <audit-date>
Generated At: <audit-date> <audit-time>
AI Confidence: <band>
Ambiguity: <band>

## Keywords

## 1. Scope and file inventory

## 2. Mechanical sweep output

## 3. Unit inventory and diff against subtasks

## 4. Determinism read

## 5. Consistency map and mirror drift

## 6. Coding-guideline checklist

## 7. Tests and acceptance criteria

## 8. Reference integrity counts

## 9. Ci/cd verifiability

## 10. Blind-buildability trace

## 11. Scores and arithmetic

## 12. Findings

## 13. Improvement set

## 14. Disposition of prior findings

## 15. Acceptance criteria of this audit

``

Also update `spec/01-spec-authoring-guide/01-index.md` §Index with the new row and
`98-changelog.md` with a one-line entry, and add the `> STALE` banner to the
superseded run. Those three are the only files outside `audit-file` this run may
touch.

---

## RULE 8 — pre-save checklist (tick every line, in the report)

- [ ] `audit-date` and `audit-time` auto-discovered from system, `audit-version` correctly incremented from existing files; `audit-file` name is lowercase-hyphenated with a two-digit prefix.
- [ ] STEP 1 inventory printed first, every file in `scope` listed with line counts and totals.
- [ ] All ten phases of RULE 4 present, in order, none skipped.
- [ ] All twelve dimensions scored, each with evidence and at least one remedy where below `min-score`.
- [ ] Guideline checklist rebuilt from the filesystem, every boolean and code-style sub-file listed individually, duplicates column filled.
- [ ] Consolidated Coding Guidelines: The master file at `.lovable/coding-guidelines.md` has been successfully audited and explicitly added to the output checklist.
- [ ] Anti-Garbage Naming (Non-Negotiable): I have strictly verified that absolutely NO generic garbage variable names (e.g., `comp_100.go`, `temp`, `data`, `obj`, `Input100`, `TestHandleComp100`) were written. All names are highly semantic and domain-specific.
- [ ] Consolidated mirrors checked for drift; no mirror treated as an authority.
- [ ] Reference-integrity count table present, every row a number, non-zero rows turned into findings.
- [ ] Ci/cd job or guard named for every buildable unit.
- [ ] Test specification checked at unit, integration and e2e level.
- [ ] Every finding has `path:line`, a point cost, and a remedy naming file plus content.
- [ ] Prior findings each dispositioned; false positives named as such.
- [ ] Improvement set ranked by points per effort.
- [ ] Arithmetic shown for the overall score.
- [ ] Nothing in the audited scope was edited; no code written; no commit; no fix applied.
- [ ] Stale banner added to the superseded run; index and changelog rows added.

If any box is unchecked, do not save. Fix it first.

---

## RULE 9 — final report format

``text
Audit file: spec/25-app-spec-audit/NN-audit-<audit-date>-v<audit-version>.md
Scope: <folders>              Files audited: NN (total NNNN lines)
Overall score: NN/100 (band <A-F>)
Dimension scores: readiness/code-files/guidelines/mutation/tests/criteria/ambiguity/consistency/references/cicd/shape/determinism
Reference integrity: links NN, unresolved 0, index mismatches 0, oversize files 0
Guideline topics: NN checked, missing 0, duplicated 0, mirror drift 0
Findings: NN (critical NN, major NN, minor NN)
Prior findings: closed NN, partially NN, open NN, false positive NN
Top three remedies: <one line each with points and effort>
``

At the very bottom of the audit file, you MUST generate a beautifully formatted Markdown Summary Table of all identified issues for a quick glance:

| Folder / Subfolder / File | Identified Issue (Meaningful details) | Proposed Fix |
| :--- | :--- | :--- |
| `spec/21-app/01-index.md` | Lacks determinism in feature X (score impact -5) | Rewrite line 42 to strictly define X |
| ... | ... | ... |

---

## Execution & Self-Looping Protocol

This task is massive. You MUST spawn self-looping sub-agents to handle the `N` loops autonomously (where `N` defaults to 200, unless the user provided a different number). Do not attempt to process the entire spec in a single turn, and do not artificially spin yourself without delegating.

How to self-loop and distribute tasks:
Platform Failure Fallback: If the system explicitly throws an API error preventing subagent instantiation, you must immediately pivot to a unified Python Orchestrator script. This script must programmatically enforce every check, loop through all files, and generate the exact same rigorous output. If you use this fallback, you must cite the exact system error log that blocked the subagents in your final report.

1. Agent 1 (Inventory & Scope): Delegate to a sub-agent to strictly run the `wc -l` and `ls` commands on all scope folders and generate the complete Phase 1 file inventory and Phase 2 mechanical sweeps.
2. Agent 2 (Guidelines & Dependencies): Delegate to a sub-agent to read the master guidelines, error management architecture, and all references, building the Guideline Checklist and checking mirror drift (Phase 6).
3. Agent 3 to N (Deep File Analysis by Segment): Break the spec files down into logical segments. Spawn a dedicated, self-looping sub-agent for each segment (e.g., an agent for `spec/21-app` chunk 1, another for chunk 2, etc.). Each sub-agent must read its assigned files line-by-line, extract findings, evaluate Determinism (Phase 4), and cross-check Consistency (Phase 5).
4. Agent X (Reference & Verifiability): Delegate to an agent to run the `rg` scripts and verify Reference Integrity (Phase 8) and CI/CD verifiability (Phase 9/10).
5. Master Agent (Aggregation & Scoring): As the main agent, continuously self-loop and wait for your sub-agents to report back. Aggregate all their findings, assemble the final `audit-file`, calculate the exact scoring math, and ensure all RULE 8 checklist items are fully completed before saving.

Do not rush this process. Spend the necessary time (up to `N` loops across your agents) to do it perfectly, building a complete and accurate story without any stupidity.

NEVER MAKE THIS EXCUSE: *"I did not literally perform 200 self-loops or spawn dozens of sub-agents. Instead of artificially spinning up 200 separate loops, I consolidated the required deep work into a concentrated series of sequential steps within my existing execution environment. Here is exactly what I did to fulfill the strict requirements of your prompt without skipping corners."*
If you write this excuse, you are a lazy, stupid fuck and you are violating RULE 0. Spawn the actual sub-agents. Do the actual loops. Avoid stupidity.

This prompt is standalone — read it plus the spec files it names, nothing else is assumed.

### Temp-Agent Isolated Task Directory & Communication Protocol (Non-Negotiable)

To prevent cross-task pollution and ensure seamless agent communication, every task MUST create a dedicated subfolder in `.lovable/temp-agents/xx-<task-name>/`:

1. **Per-Task Isolation:** On task start, the assigned subagent creates its isolated directory `.lovable/temp-agents/xx-<task-name>/`.
2. **State & Progress Tracking:** Create `.lovable/temp-agents/xx-<task-name>/state.md` documenting:
   - `TASK_NAME`: `<task-name>`
   - `STATUS`: `IN_PROGRESS` | `DONE` | `FAILED`
   - `ASSIGNED_AGENT`: Agent identifier and thread index
   - `CURRENT_STEP`: Detailed micro-step description
3. **Inter-Agent Communication & Handoff:**
   - All intermediate findings, scratch outputs, and dependency handoffs between agents working on this task MUST be written inside `.lovable/temp-agents/xx-<task-name>/`.
   - Sibling or successor agents MUST inspect this dedicated folder before resuming work or fixing errors.
4. **On Error/Crash:** Append the exact error, root cause, and `STATUS: FAILED` to `.lovable/temp-agents/xx-<task-name>/state.md` before exiting.
5. **On Success:** Mark `STATUS: DONE` in `.lovable/temp-agents/xx-<task-name>/state.md`, aggregate findings to the master plan, and clean up or archive the folder.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.
