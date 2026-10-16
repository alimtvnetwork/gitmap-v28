# Subtask 03 — WS3: commit review → 10 commands + 20 improvements

**Version:** 1.0.0
**Updated:** 2026-10-17
**Owner:** Review worker
**Status:** pending (may run in parallel with WS2; read-only)

---

## Goal

Review the **last 100 commits** and produce **10 new command suggestions**
and **20 codebase improvement proposals** — proposals only, no code.

## Inputs

- `03-commit-review-method.md` — review method (normative).
- `04-improvement-criteria.md` — proposal acceptance + hard-exclusion list (normative).
- `.ai-memory/plans/272-decision-record-improvement-scope.md` — owner-rejected items.

## Steps

1. Walk the last 100 commits newest → oldest; read commit messages, skim
   diffs only where a message is ambiguous.
2. Tag each commit into a theme bucket (`03-commit-review-method.md` §4);
   tally; buckets with ≥3 commits are hotspots.
3. For each hotspot, apply the five questions (§3): repeated manual
   workflow? missing automation? integration gap? recurring bug class?
   documentation that should be behavior?
4. Draft proposals using the template in `04-improvement-criteria.md` §4 —
   each with name, one-line, commit evidence, 3–5 bullet sketch, 2–4
   acceptance criteria, and an explicit exclusion-check line.
5. Self-audit the batch against the hard-exclusion list (§3 of
   `04-improvement-criteria.md`): pipeline DB retention, staging-repo
   philosophy, error-message overhaul, release-tag management changes.
   Any proposal mapping to those four — even renamed — is dropped.
6. Record the final 10 + 20 under `.ai-memory/` (proposal-only; no code, no
   branches).

## Acceptance criteria

- [ ] Exactly 10 command suggestions and exactly 20 improvement proposals.
- [ ] Every proposal cites real commits from the last-100 review as evidence.
- [ ] Every proposal carries the exclusion-check line referencing
      `04-improvement-criteria.md` v1.0.0.
- [ ] Zero proposals touch the hard-exclusion list in any form.
- [ ] No code written, no branches created, no commits made by WS3.

## Out of scope

Implementation; ranking beyond the evidence tally (owner decides priority);
revisiting owner-rejected ideas.
