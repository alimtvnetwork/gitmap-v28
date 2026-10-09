# Subtask 03 — TrainPhase Interface + Registry + Heal Phase (program 246)

## Objective

Implement spec `02-spec/21-app/246-parallel-fix-command/03-llm-train-integration.md`:
a `TrainPhase` interface with an ordered registry in `cli/cmd/llm`, plus a real
executable 6th phase ("Heal & Fix") with two in-process sub-steps — (a) Heal:
git-state remediation via the existing `cmdfix` engine; (b) Fix: content audit via
the NEW `cmdautofix` engine (program 246 spec 02) — and 6-phase consistency across
all text surfaces.

## Read first

- `02-spec/21-app/246-parallel-fix-command/03-llm-train-integration.md` (full spec)
- `.ai-memory/what-to-read.md` — "The Mindset" (§5: ~300 lines max per file)

## Owned files (ONLY these — do not touch any other file)

New:

- `cli/cmd/llm/llm_phases.go` (≤150 lines)
- `cli/cmd/llm/llm_heal.go` (≤150 lines)

Edit:

- `cli/cmd/llm/llm_train.go` — `RunTrain` iterates the registry; `parseTrainFlags`
  gains `--heal-apply`; `outputTrainJSON`/`executeSelfLoopSimulation`/
  `printChainedSequence`/`printTrainUsage` text updates per spec §5 table
- `cli/cmd/llm/llm.go` — ASCII diagram Phase 6 box + "6-phase lifecycle" wording
- `cli/cmd/llm/llm_types.go` — `TrainOptions.HealApply string`; `OperationalDirectivesText` item 8
- `cli/cmd/llm/llm_skill.go` — §2 "6-phase" line only (the §10 fix section is subtask 04)
- `cli/cmd/llm/llm_urls.go` — one word: "5-phase" → "6-phase" in the first link description
- `cli/cmd/llm/llm_train_test.go` — extend only (see checklist)

## Implementation checklist

- [ ] `llm_phases.go`: `TrainPhase` interface (`Name() string`,
      `Run(TrainOptions) *apperror.AppError`), `TrainPhaseResult`, ordered
      `trainPhases` registry (Discovery → Refactoring → Verification →
      Semantic Commit → Telemetry → Heal & Fix), `runTrainPhases` iterator with
      per-stage headers + closing summary table.
- [ ] Phases 1–5 are adapter structs printing byte-identical text to today's
      `RunTrain` sequence (no behavior change to the curriculum).
- [ ] `llm_heal.go`: `healPhase` sub-step A (Heal) loads
      `cmdremediation.LoadRemediationState()`; on empty state calls
      `cmdfix.RunFix([]string{}, "")` once for live discovery, then reloads;
      prints the per-category table from `SummaryReason` buckets
      (`dirty-worktree`, `diverged`, `untracked`, `other`).
- [ ] `llm_heal.go`: `healPhase` sub-step B (Fix) calls
      `cmdautofix.Scan(cmdautofix.Options{Workers: runtime.NumCPU()})`
      directly (check-only — NEVER `RunFixCmd`, which would prompt on stdin;
      never apply from the train); prints the per-category roll-up from the
      `ScanResult` and the canonical next action `gitmap fix all -y`.
      Verify no import cycle: `cli/cmdautofix` must not import `cli/cmd/llm`.
- [ ] Report-only by default. `--heal-apply stash|wip|discard` → 
...[truncated 1703 chars]