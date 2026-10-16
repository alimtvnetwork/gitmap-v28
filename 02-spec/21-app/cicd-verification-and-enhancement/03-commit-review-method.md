# 03 — Commit Review Method

**Version:** 1.0.0
**Updated:** 2026-10-17

---

## 1. Purpose

How WS3 reviews the **last 100 commits** to surface feature ideas — the raw
material for the 10 command suggestions and 20 improvement proposals. This is
a **read-only** review: no code changes, no refactors, proposals only.

## 2. Scope and method

1. Walk the last 100 commits newest → oldest, reading commit messages (and
   diffs only where the message is ambiguous — skim, don't deep-read).
2. Tag each commit into a theme bucket (see §4). A commit gets one primary
   bucket.
3. Tally the buckets. Buckets with ≥3 commits are **hotspots** — they are
   where the owner actually spends effort, and where automation pays off.
4. For each hotspot, ask the five questions in §3 and draft proposals that
   pass `04-improvement-criteria.md`.
5. Never propose anything on the hard-exclusion list
   (`.ai-memory/plans/272-decision-record-improvement-scope.md`).

## 3. The five questions (per hotspot)

1. **Repeated manual workflow?** — Did the owner (or agents) do the same
   multi-step thing by hand more than once? → command candidate that folds
   the steps into one verb.
2. **Missing automation?** — Is there a check, guard, or generation step that
   humans keep remembering to do? → enforcement or codegen candidate
   (mindset: enforcement beats documentation).
3. **Integration gap?** — Do two features exist but never talk to each other
   (e.g. a report a dashboard ignores, a flag a sibling command lacks)? →
   wiring candidate.
4. **Recurring bug class?** — Same kind of fix appearing repeatedly (path
   handling, ANSI leakage, off-by-one in tables)? → guard/linter candidate.
5. **Documentation that should be behavior?** — A rule written in a doc that
   keeps getting violated? → candidate to enforce it in code or CI.

## 4. Theme buckets (starting set — extend as the review demands)

| Bucket | Signals |
|--------|---------|
| `ci-fix` | workflow edits, gate tweaks, flaky-test workarounds |
| `help-text` | helpdoc regeneration, example additions, wording fixes |
| `theme-ui` | theme/color/table rendering changes |
| `release` | bump scripts, changelog, asset, installer edits |
| `new-command` | fresh verbs, flags, aliases |
| `refactor` | file splits, renames, dedup |
| `bugfix` | fixes with RCA references |
| `spec-plan` | spec, plan, subtask, memory-only commits |
| `infra` | runner scripts, DB migrations, telemetry |
| `debt` | TODO removals, dead-code cleanup |

## 5. Output format (per proposal)

Each of the 10 command suggestions and 20 improvement proposals is recorded
in `.ai-memory/` with:

- **name** — proposed verb/area (kebab-case, lowercase).
- **one-line** — what it does, in the owner's vocabulary.
- **evidence** — which commits (shas or message quotes) prove the need; a
  proposal without commit evidence is dropped.
- **sketch** — 3–5 bullet behavior sketch (flags, output, safety). No code.
- **acceptance criteria** — 2–4 testable lines for the future implementer.
- **exclusion check** — explicit line: "not on the hard-exclusion list"
  with the list version/date referenced.

## 6. Non-goals

- No ranking beyond the evidence tally — the owner decides priority.
- No implementation, no POC code, no branch creation.
- No revisiting owner-rejected ideas (see `04-improvement-criteria.md` §3);
  "never approved" is a hard stop, not a backlog.
- No inventing commits: every evidence line must be a real commit from the
  last 100.
