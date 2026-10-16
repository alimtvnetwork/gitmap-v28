# Plan — cicd-verification-and-enhancement (program 274)

**Version:** 1.0.0
**Updated:** 2026-10-17
**Status:** Active (Wave 1 — WS1 in progress)

---

## Program charter

> "CI/CD verification + enhancement program: confirm CI green, fix failures,
> minor bump + release until green; 10 new command suggestions; review last
> 100 commits; 20 codebase improvements excluding rejected ideas; E2E testing
> of recent features; critical feedback."

Definition of done: release published AND CI green on the release commit
(see `02-spec/21-app/cicd-verification-and-enhancement/01-overview.md` §5).

## Workstream checklist

| WS | Workstream | Owner | Spec | Subtask file | Status |
|----|-----------|-------|------|--------------|--------|
| WS1 | Spec + plan + subtasks | Worker 1 (this wave) | `02-spec/21-app/cicd-verification-and-enhancement/` | `.ai-memory/plans/subtasks/cicd-verification-and-enhancement/01-ws1-spec-plan.md` | in progress |
| WS2 | CI green loop | CI worker | `02-ci-green-loop-procedure.md` | `02-ws2-ci-green-loop.md` | pending |
| WS3 | Commit review: 10 commands + 20 improvements | Review worker | `03-commit-review-method.md`, `04-improvement-criteria.md` | `03-ws3-commit-review.md` | pending |
| WS4 | E2E of recent features | E2E worker | `05-e2e-test-plan.md` | `04-ws4-e2e-feedback.md` | pending |
| WS5 | Minor bump + release until green | Lead | `01-overview.md` §5 (DoD) | `05-ws5-release.md` | BLOCKED until WS2 green |

## File boxes

### WS1 owns (write these, nothing else)

- [ ] `02-spec/21-app/cicd-verification-and-enhancement/01-overview.md`
- [ ] `02-spec/21-app/cicd-verification-and-enhancement/02-ci-green-loop-procedure.md`
- [ ] `02-spec/21-app/cicd-verification-and-enhancement/03-commit-review-method.md`
- [ ] `02-spec/21-app/cicd-verification-and-enhancement/04-improvement-criteria.md`
- [ ] `02-spec/21-app/cicd-verification-and-enhancement/05-e2e-test-plan.md`
- [ ] `.ai-memory/plans/cicd-verification-and-enhancement.md` (this file)
- [ ] `.ai-memory/plans/subtasks/cicd-verification-and-enhancement/01-ws1-spec-plan.md`
- [ ] `.ai-memory/plans/subtasks/cicd-verification-and-enhancement/02-ws2-ci-green-loop.md`
- [ ] `.ai-memory/plans/subtasks/cicd-verification-and-enhancement/03-ws3-commit-review.md`
- [ ] `.ai-memory/plans/subtasks/cicd-verification-and-enhancement/04-ws4-e2e-feedback.md`
- [ ] `.ai-memory/plans/subtasks/cicd-verification-and-enhancement/05-ws5-release.md`

### Reserved for a later wave (do NOT create now)

- [ ] `02-spec/21-app/cicd-verification-and-enhancement/06-e2e-critical-feedback.md`

### WS3 outputs (proposals only, later wave)

- [ ] 10 new command suggestions → `.ai-memory/`
- [ ] 20 improvement proposals → `.ai-memory/`

## Sequencing and gates

1. WS1 completes → lead reviews spec + plan.
2. WS2 runs the CI green loop until 100% green; WS3 and WS4 may run in
   parallel with WS2 (read-only work), but WS2's fix commits take sequencing
   priority — lead coordinates.
3. WS2 reports green on an explicit HEAD sha → WS5 unblocked.
4. WS5: minor bump (from 6.557.0) + release; if the release commit is red,
   repeat bump+release until green. NEVER cut a new bump/release while a
   Release run is in flight.
5. Lead verifies everything personally before declaring DONE (standing rule:
   never trust self-reported "green").

## Risks / open items

- Standing open items are OUT of scope unless the owner approves: Send/Queue
  buttons on owner's machine; PR #6 'Sync v5' do-not-merge; release-pipeline
  durable fix proposed 2026-10-10/11 (unanswered — do not implement).
- `06-e2e-critical-feedback.md` must not exist until the later wave writes it.
