# 04 — Improvement Criteria

**Version:** 1.0.0
**Updated:** 2026-10-17

---

## 1. Purpose

What makes a command suggestion or improvement proposal acceptable for program
274 — and what is permanently off the table. Every WS3 proposal must pass §2
and must not touch §3.

## 2. Acceptance criteria for a proposal

A proposal is ACCEPTED only if all of these hold:

1. **Owner-shaped.** It serves the staging-repo mindset: velocity,
   fix-forward, small additive changes. It does not slow the loop for
   process or add ceremony.
2. **Evidence-grounded.** It cites real commits from the last-100 review
   (shas or message quotes). No evidence → rejected.
3. **Small.** Implementable within the repo's file-size discipline (~300
   lines max per file); a proposal that needs a big-bang rewrite is rejected.
4. **Additive.** It adds a verb, flag, guard, or wiring — it does not remove,
   rename, or prune existing commands. Aliases are a feature, never sprawl
   to be cleaned.
5. **Single purpose.** One verb does one thing; no kitchen-sink commands.
6. **Enforcement-flavored where possible.** If the idea is "people should
   remember X", prefer a check/gate/generator over a doc paragraph
   (enforcement beats documentation; DRY single source of truth).
7. **Testable future.** Carries 2–4 concrete acceptance criteria a future
   implementer can verify.
8. **Safe by default.** Dry-run or read-only where feasible; destructive or
   external actions need explicit flags and confirmation.

## 3. HARD EXCLUSION list — owner-rejected, never propose

Copied from `.ai-memory/plans/272-decision-record-improvement-scope.md`
(owner decision 2026-10-15; brief adds release-tag management as #4).
An explicit "never approved" is a hard stop, not a backlog. These stay parked
until the owner reopens them — do not propose them, do not work around them,
do not re-litigate them.

| # | Rejected idea | Owner decision |
|---|---------------|----------------|
| 1 | **Pipeline DB retention policy** — retention/expiry rules for pipeline databases | NOT approved (program 272 item 4) |
| 2 | **Staging-repo philosophy decision** — revisiting or formalizing the staging-repo model | NOT approved (program 272 item 5) |
| 3 | **Actionable error-message overhaul** — repo-wide rewrite of error messages | NOT approved (program 272 item 6) |
| 4 | **Release-tag management changes** — any change to how release tags are created, named, or managed | NOT approved (program brief) |

Proposing any of the four — or a renamed/restructured variant of any of the
four — is an automatic rejection of the whole WS3 batch. When in doubt,
leave it out.

## 4. Proposal template (WS3 output)

```markdown
## <kebab-case-name>

- one-line: <what it does, owner's vocabulary>
- evidence: <commit shas / message quotes from the last-100 review>
- sketch:
  - <behavior bullet 1>
  - <behavior bullet 2>
  - <behavior bullet 3>
- acceptance:
  - <testable criterion 1>
  - <testable criterion 2>
- exclusion-check: not on the hard-exclusion list (04-improvement-criteria.md v1.0.0, 2026-10-17)
```

## 5. Batch acceptance (lead verification)

- Exactly 10 command suggestions + exactly 20 improvement proposals.
- Every proposal carries evidence lines and an exclusion-check line.
- Zero proposals map to §3 in any form.
- Recorded under `.ai-memory/` (proposal-only; no code, no branches).
