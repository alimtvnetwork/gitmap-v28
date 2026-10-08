# Spec: 244-e2e-verify-program-243-changes — Test Scope and Plan

Slug: `244-e2e-verify-program-243-changes`
Owner: Spec Agent 1 (mode: self, read-write)
Binary under test: built from pristine HEAD `f127f7a` (committed program-243 state), expected at `~/workspace/gitmap-e2e-244/gitmap`

## User Request (Verbatim)

```
# High Priority Instruction

Can you please run `gitmap` space PA first, and then you start with pulling the code. And also can you please run and test the recent changes in terms of writing some end-to-end scripts and run it here and see if you can find a way to verify that, does it work or not? Can you please do that for me?

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Run `gitmap` space PA first
5. Pull the code
6. Run and test the recent changes by writing some end-to-end scripts
7. Verify if the scripts work or not

## Must follow and spawn agent using

[execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing the work to reduce the credits.
- [/learn](slashCommand;learn) from [gitmap](file;.agents/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
```

## Objective

Write and run end-to-end scripts that verify the **committed** program-243 changes
work as specified, using a scratch-built binary, without touching the real repo or
pushing anything to origin.

## Scope Boundary (Non-Negotiable)

- IN SCOPE: the committed program-243 state — HEAD `f127f7a` on main: `space
  backup-branch` (`cli/cmdspace/`), help-displayer core (`cli/helpdisplay/`),
  DRY migration for scan/pull-all/space, wave-D splits, `cli/enums/`, pre-commit
  file-size gate (`03-ai-scripts/50-fastgate.py` calling `check_go_file_sizes`
  from `03-ai-scripts/52-file-size-check.py`, effective `DEFAULT_LINE_LIMIT=500` at committed HEAD `f127f7a`).
- OUT OF SCOPE: the in-flight 243b-consolidation uncommitted work (broken import
  cycles, ~629 uncommitted changes, actively worked by another coordinator).
  Nothing under test is built from the dirty working tree.
- No real pushes to origin. No commits of test artifacts into the repo.
- Known discrepancy (not a 243 failure): the 300-line cap exists only in docs
  (mindset file, spec 243) at committed HEAD `f127f7a`; the committed code still
  enforces 500 (`52-file-size-check.py` default and `50-fastgate.py`
  docstring/message agree). The 243b working tree already flips the default to
  300 — that alignment is 243b's scope.

## Test Groups and Assertions

### A. backup-branch contract (10 assertions)

A1. `gitmap space backup-branch --help` — exit 0; stdout contains a plain-text
    usage const including the literal `--no-push`.
A2. `gitmap space --help backup-branch` — exit 0; stdout contains the Displayer
    group headers `About`, `Behavior`, `Flags`, `Examples`.
A3. `gitmap space backup-branch` (no args) — exit 1; stderr mentions "missing
    task string".
A4. `gitmap space backup-branch "!!!"` — exit 1; stderr mentions "empty backup
    slug".
A5. Dirty scratch repo + `--no-push "X"` — exit 1; stderr mentions "uncommitted
    changes". (Scratch repo: `git init` + one commit, then an uncommitted edit.)
A6. Clean scratch repo without origin:
    `gitmap space backup-branch --no-push "E2E Smoke Test_42"` — exit 0; stdout
    matches `^backup/e2e-smoke-test-42 @ [0-9a-f]+$`; `git rev-parse` confirms
    the `backup/e2e-smoke-test-42` branch exists in the scratch repo.
A7. Same clean scratch repo without origin, WITHOUT `--no-push` — exit 1;
    stderr mentions push failure (no origin configured).
A8. Pre-create branch `backup/dup-case` in scratch repo; run
    `gitmap space backup-branch --no-push "dup case"` — exit 1; stderr mentions
    "already exists". Re-run with `--force` — exit 0 and the branch points at
    current HEAD.
A9. Task string of 80 chars (e.g. 80 × `a`) — exit 0 with `--no-push`; resulting
    branch name slug portion is capped at 60 chars.
A10. Informational: in a scratch clone, `git branch -r` lists
    `origin/backup/2026-10-08-pre-help-displayer` (the real backup branch created
    by program 243).

### B. help displayer rendering (7 invocations)

Each: exit 0, a header line present, group headers present.

| #  | Command                              | Suggestion `* ` lines expected |
|----|--------------------------------------|---------------------------------|
| B1 | `gitmap scan --help`                 | no                              |
| B2 | `gitmap s --help`                    | no                              |
| B3 | `gitmap space --help`                | yes                             |
| B4 | `gitmap space --help common`         | yes                             |
| B5 | `gitmap space --help backup-branch`  | n/a (see A2)                    |
| B6 | `gitmap pull-all --help`             | yes                             |
| B7 | `gitmap pa --help`                   | yes                             |

### C. unknown-command suggestions (4 inputs)

C1. `gitmap scna` — exit 1; stderr contains "Did you mean?" and "gitmap scan".
C2. `gitmap docker` — exit 1; stderr contains "gitmap cluster".
C3. `gitmap releas` — exit 1; stderr contains "gitmap release".
C4. `gitmap somethingrandom` — exit 1; stderr contains "unrecognized command"
    (case-insensitive, E1001); no suggestion box rendered.

### D. enforcement + regression

D1. `gitmap --version` — exit 0; stdout prints exactly `gitmap v6.515.0`.
D2. `03-ai-scripts/52-file-size-check.py` defines `DEFAULT_LINE_LIMIT = 500` (committed reality at `f127f7a`; see discrepancy note).
D3. `cli/enums/` directory exists.
D4. Read-only `gitmap aum search <term>` and `gitmap lf` — exit 0.

### E. build verification

E1. `go build ./...` run from `cli/` — exit 0.

## Harness Design

- Bash scripts live UNDER `~/workspace/gitmap-e2e-244/` (OUTSIDE the repo, never
  committed): `e2e-backup-branch.sh`, `e2e-help-displayer.sh`,
  `e2e-suggestions.sh`, `e2e-enforcement-regression.sh`, and `run-all.sh`.
- `run-all.sh` orchestrates the four group scripts, printing per-assertion
  `PASS`/`FAIL` lines, summary counts, and exits non-zero on any failure.
- All destructive cases (branch creation, dirty-tree simulation) run ONLY in
  disposable scratch git repos under `~/workspace/gitmap-e2e-244/work/`
  (`git init` + one commit, no origin). Test branches are cleaned with
  `git branch -D` inside the scratch repos only.
- NEVER run destructive cases against the real repo
  (`~/workspace/repos/gitmap-v28`).
- Scripts use `set -u -o pipefail`, capture exit codes explicitly, and use the
  scratch binary path via env var `GITMAP_BIN` (default
  `~/workspace/gitmap-e2e-244/gitmap`).

## Non-Goals

- No testing of the in-flight 243b uncommitted work.
- No real pushes to origin.
- No commits of test artifacts into the repo.
- No performance or concurrency testing; single sequential runs only.

## Related Paths (relative)

- Spec: `02-spec/21-app/244-e2e-verify-program-243-changes/01-test-scope-and-plan.md`
- Plan: `.ai-memory/plans/244-e2e-verify-program-243-changes.md`
- Subtask: `.ai-memory/plans/subtasks/244-e2e-verify-program-243-changes/01-author-e2e-scripts.md`
