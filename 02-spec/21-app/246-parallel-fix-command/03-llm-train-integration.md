# Spec 03 — `gitmap llm train` × Fix Engine Integration (program 246)

## 1. Goal

Give `gitmap llm train` a real, executable 6th phase — **Heal & Fix** — with two
in-process sub-steps: (a) **Heal**: workspace git-state remediation via the existing
`cmdfix` engine; (b) **Fix**: file-content fixing via the NEW `gitmap autofix` engine
(`cli/cmdautofix`, program 246 spec 02 — the parallel port of the 8 AI fixer
scripts: encoding, newlines, naming, paths, gofmt, misspell, markdown, guidelines).
Replace the current print-only 5-phase simulation with a `TrainPhase` interface +
ordered registry, and bring every text surface (JSON catalog, ASCII diagram,
self-loop, chained sequence, `--help`, skill template) to 6-phase consistency.

## 2. Ground truth (verified 2026-10-09)

- `cli/cmd/llm/llm_train.go` — `RunTrain` is sequential print-only. Its only real side
  effect is writing `.agents/skills/gitmap/SKILL.md` via `handleSkillGeneration`.
- The "5 phases" exist ONLY as static text: the `phases` array in `outputTrainJSON()`,
  the ASCII diagram in `llmMarkdownSpec` (`cli/cmd/llm/llm.go`), and prints in
  `executeSelfLoopSimulation()`.
- `cmdfix.RunFix(args []string, aliasOverride string) error` (`cli/cmdfix/fix_cmd.go`)
  is the in-process git-state engine entry — the same engine `gitmap fix` uses
  (`cli/cmd/rootcore.go:79-82` registers `fix`, `stash`, `wip`, `discard` on it).
- `cmdautofix.RunAutofixCmd(args []string) error` (`cli/cmdautofix/fix.go`, NEW in
  program 246 spec 02) is the in-process content-fix engine entry — the same
  engine `gitmap autofix` uses. For the train's report-only needs it is invoked
  as `RunAutofixCmd([]string{<path>})` (dry-run by default; no `--apply`).
  NOTE: `cli/cmd/llm` importing `cli/cmdautofix` must not create an import cycle —
  `cmdautofix` must not import `cmd/llm` (implementer verifies with the build).
- `cmdremediation.LoadRemediationState() []RemediationItem` is public;
  `RemediationItem{RepoPath, RepoName, SummaryReason, Recipes, Files}`
  (`cli/cmdremediation/remediation_box.go:17`).
- Import-cycle check 2026-10-09: nothing in `cmdfix`, `cmdremediation`, `cmdreconcile`,
  `cmdagy`, `cmdignore`, `gitutil`, `constants`, `apperror` imports `cmd/llm` — so
  `cli/cmd/llm` may import `cli/cmdfix` + `cli/cmdremediation`. Implementer re-verifies
  with the central build; if a cycle appears, the heal phase falls back to reporting
  via `cmdremediation` only (no `cmdfix` import).

## 3. `TrainPhase` interface + registry

New file `cli/cmd/llm/llm_phases.go` (≤150 lines):

```go
// TrainPhase is one executable stage of the llm train curriculum.
type TrainPhase interface {
    Name() string
    Run(opts TrainOptions) *apperror.AppError
}

// TrainPhaseResult is the outcome line of one phase for the closing summary.
type TrainPhaseResult struct {
    Name   string
    Status string // "ok" | "skipped" | "error"
    Detail string
}

var trainPhases = []TrainPhase{
    discoveryPhase{},
    refactorPhase{},
    verificationPhase{},
    commitPhase{},
    telemetryPhase{},
    healPhase{},
}

func runTrainPhases(opts TrainOptions) *apperror.AppError
```

- Phases 1–5 are thin adapter structs whose `Run` prints **exactly the same text as
  today** (attribution, stage headers, chained sequence, operational guidelines).
  Behavior-preserving: the curriculum stays print-first.
- `runTrainPhases` iterates the registry in order, prints a `STAGE N — <name>` header
  per phase, collects `TrainPhaseResult`s, and prints a closing summary table.
- `RunTrain` calls `runTrainPhases(opts)` in place of the current inline print
  sequence. `handleSkillGeneration` is unchanged (skill write is phase-agnostic).
- Early-return flags keep current behavior: `--json`, `--url(s)`, `--loop`,
  `--self-loop` return before the registry runs.

## 4. The heal phase

New file `cli/cmd/llm/llm_heal.go` (≤150 lines). `healPhase.Name()` returns
`"6. Heal & Fix"`. `Run` does:

**Sub-step A — Heal (git-state, existing engine):**

1. Print `STAGE 6a: HEAL — workspace git-state remediation (report-only)`.
2. `items := cmdremediation.LoadRemediationState()`. If empty, call
   `cmdfix.RunFix([]string{}, "")` once — its empty-state path performs live discovery
   (`handleLiveDiscoveredIssues`) and saves state — then reload items.
3. Build per-category counts from `item.SummaryReason` via a small keyword bucket map
   in this file (buckets: `dirty-worktree`, `diverged`, `untracked`, `other`).
   Buckets are derived, never invented: bucket from the reason string, fallback `other`.
4. Print a per-category table: `CATEGORY | REPOS | EXAMPLE REPO`.
5. Report the canonical next actions: `gitmap fix ls`, `gitmap fix <repo> <action>`.

**Sub-step B — Fix (content, new autofix engine):**

6. Print `STAGE 6b: FIX — file-content audit via autofix engine (report-only)`.
7. Call `cmdautofix.RunAutofixCmd([]string{"."})` — dry-run by default, so this
   audits without writing. Capture its per-category summary (the engine prints
   `CATEGORY | FILES FLAGGED` lines; the phase re-prints the roll-up).
8. Report the canonical next action: `gitmap autofix --apply` to write the fixes.

Safety rules (non-negotiable):

- **Report-only by default.** Git-state recipes are applied ONLY when
  `--heal-apply <action>` is passed, via `cmdfix.RunFix([]string{"all", action}, "")`,
  where action is one of `stash|wip|discard` (existing `cmdfix` parsing; invalid
  actions error as today). Content fixes are NEVER applied by the train —
  `RunAutofixCmd` is always invoked without `--apply`; the train reports, the
  operator runs `gitmap autofix --apply`.
- `--heal-apply` is **ignored under `--text-only`** (no mutations in text-only mode).
- `RunFix` returns plain `error` → wrap with `apperror.WrapSimple(err, "healPhase.fix")`.
- The phase never deletes files and never prunes aliases (repo hard rules).

`TrainOptions` gains `HealApply string` (`cli/cmd/llm/llm_types.go`).
`parseTrainFlags` gains:

```go
healApply := fs.String("heal-apply", "", "Apply fix recipes non-interactively: stash|wip|discard (default \"\" = report-only)")
```

## 5. Text-level consistency updates (exact locations)

| File | Location | Change |
| :--- | :--- | :--- |
| `llm_train.go` | `outputTrainJSON()` phases array | Append `{"phase": "6. Heal & Fix", "commands": "gitmap fix ls, gitmap autofix", "purpose": "In-process remediation scan + parallel content-fix audit (report-only); optional --heal-apply for git-state recipes"}` |
| `llm_train.go` | `executeSelfLoopSimulation()` | Header `AUTONOMOUS 5-PHASE` → `AUTONOMOUS 6-PHASE`; add line `  Phase 6 [Heal & Fix]:   In-process remediation scan (cmdfix) + autofix content audit (report-only)` |
| `llm.go` | `llmMarkdownSpec` ASCII diagram | Append Phase 6 box after Phase 5 (same box-drawing style): `Phase 6: Heal & Fix (report-only)` / `➔ gitmap fix ls + gitmap autofix — cmdfix + cmdautofix engines in-process` |
| `llm.go` | `llmMarkdownSpec` §1 paragraph | "structured 5-phase lifecycle" → "structured 6-phase lifecycle" |
| `llm_train.go` | `printChainedSequence()` | Append `12. gitmap fix ls  — List repos needing remediation with per-category heal summary` and `13. gitmap autofix — Parallel content-fix audit (encoding, newlines, naming, paths, gofmt, misspell, markdown)` |
| `llm_train.go` | `printTrainUsage()` | Description → 6-phase; add `--heal-apply string` flag line (see §6) |
| `llm_types.go` | `OperationalDirectivesText` | Append `8. Heal Before Commit: Run 'gitmap fix ls' and 'gitmap autofix', clear remediation items before pushing.` |
| `llm_skill.go` | `SkillTemplate` §2 | "Full 4-stage chained curriculum" → "Full 6-phase chained curriculum (Discovery, Refactoring, Verification, Semantic Commit, Telemetry, Heal & Fix)" |
| `llm_skill.go` | `SkillTemplate` | Add §10 "Workspace Heal & Fix" — exact markdown defined in subtask `04-skills-docs.md` |
| `llm_urls.go` | `GetPublicDocLinks()[0].Description` | "Authoritative 5-phase execution lifecycle" → "Authoritative 6-phase execution lifecycle" |

## 6. `gitmap llm train --help` (exact output)

```
Usage: gitmap llm train [flags] (alias: gitmap llm chain, gitmap train, gitmap llm-train)

Executes the autonomous LLM chained onboarding curriculum: a 6-phase self-looping
training run (Discovery, Refactoring, Verification, Semantic Commit, Telemetry,
Heal & Fix) and generates the official Antigravity skill
(.agents/skills/gitmap/SKILL.md). Phase 6 runs the git-state heal scan (cmdfix) and
the parallel content-fix audit (cmdautofix) in-process — report-only; git-state
recipes apply only with --heal-apply.

Flags:
  --heal-apply string  Apply fix recipes non-interactively: stash|wip|discard (default "" = report-only)
  --loop               Execute autonomous 6-phase AI self-looping execution cycle
  --self-loop int      Number of consecutive iterations of the AI self-loop
  --url                Output raw public URL to llm.md instruction specification
  --urls               Output authoritative public documentation links for LLM ingestion
  --json               Output structured machine-readable command specifications
  --skill-path string  Path for generated Antigravity skill (default ".agents/skills/gitmap/SKILL.md")
  --text-only          Print curriculum text without writing skill file
```

## 7. Non-goals

- No behavior change to the `cmdfix` engine itself; `llm` only calls its public entry.
- No new subprocess, no new binary, no new config file.
- Do not touch versioning. Do not prune aliases. Never propose deleting files.

## 8. Acceptance

- `gitmap llm train --text-only` prints 6 stages including the heal report table.
- `gitmap llm train --json` contains 6 phases; phase 6 names the fix commands.
- `gitmap llm train --help` matches §6 verbatim.
- Registry order is Discovery → Refactoring → Verification → Semantic Commit →
  Telemetry → Heal & Fix; each new file ≤300 lines; existing tests in
  `llm_train_test.go` still pass unmodified.
