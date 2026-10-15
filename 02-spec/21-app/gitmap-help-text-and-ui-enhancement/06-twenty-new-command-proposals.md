# 06 — 20 New Command Proposals (template + process)

## 1. Standing rule

WS4 produces **proposals ONLY** — no implementation. Each proposal lives in
`.ai-memory/plans/gitmap-help-text-and-ui-enhancement.md` (or a linked file
under `.ai-memory/plans/`) and follows the template in §2 exactly. The user
picks which proposals become future programs.

## 2. Proposal template (mandatory fields)

```text
### P-NN — `gitmap <verb> [<subverb>]`

- PROBLEM: one paragraph — the job the user cannot do today.
- INVOCATION: one concrete, copy-pasteable command line.
- EXPECTED OUTPUT: representative output, normalized placeholders.
- ALIASES: proposed short aliases (or "none").
- HOST/TARGET: which subsystem (scanner, pull, fleet, agm, db, pipeline...).
- EFFORT: S / M / L (implementation estimate, not a commitment).
- DEPENDS ON: program 273 artifacts it needs (e.g. "help-example registry").
```

## 3. Candidate directions (seed list, not the final 20)

The WS4 worker derives exactly 20 from these directions plus the
missing-command audit gaps (spec 04) — the user-ordered areas are
Lovable/GitHub migration surface, help/UI discoverability, and fleet
ergonomics:

1. `pr` family (`pr list`, `pr view`, `pr merge`) beyond `mr` creation.
2. Migration verbs (`repo export` / `repo import` round-trip with `scan`).
3. Help discoverability (`help search <term>`, `help topic` index).
4. Theme management (`theme list`, `theme preview`).
5. Fleet ergonomics (shortcuts the audit finds missing).
6. Pipeline/CI triage shortcuts on top of `pe`/`pea`.

## 4. Acceptance for the proposal set

- Exactly 20 proposals, each with ALL template fields filled.
- Every proposal includes invocation + expected output (the help-example
  discipline, even though nothing is implemented).
- No proposal duplicates an existing verb/alias.
- Proposals reference the audit gap IDs where applicable.
- File placed under `.ai-memory/plans/`, committed + pushed with the WS4
  batch. No code, no version bump.
