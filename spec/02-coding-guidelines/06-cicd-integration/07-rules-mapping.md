# Rules Mapping — Spec → Check → Severity

> **Version:** 2.0.0
> **Updated:** 2026-04-27

Single source of truth: every CODE RED / STYLE rule, where it is defined
in the spec, which check script enforces it, and what severity each
emits.

---

## CODE RED rules (block merge — SARIF `error`)

| ID | Rule | Spec source | Check script | Phase 1 langs |
|----|------|-------------|--------------|---------------|
| CODE-RED-001 | No nested `if` | `01-cross-language/04-code-style/` | `checks/nested-if/<lang>.py` | go, ts |
| CODE-RED-002 | Boolean naming (Is/Has/Can/Should/Was/Will) | `01-cross-language/02-boolean-principles/` | `checks/boolean-naming/<lang>.py` | go, ts |
| CODE-RED-003 | No magic strings | `01-cross-language/04-code-style/` | `checks/magic-strings/<lang>.py` | go, ts |
| CODE-RED-004 | **Function length hard cap — ≤ 15 lines** (redundant safety net under CODE-RED-005) | `01-cross-language/04-code-style/` | `checks/function-length/<lang>.py` | go, ts, php |
| CODE-RED-005 | **Function length strict — ≤ 8 lines** (binding cap per coding-guidelines.md rule #1) | `01-cross-language/04-code-style/` | `checks/function-length-prefer8/<lang>.py` + `eslint-plugins/coding-guidelines/index.js#preferFunctionLines` | go, ts, php, python, rust, eslint |
| CODE-RED-006 | File length ≤ 300 lines | `01-cross-language/04-code-style/` | `checks/file-length/<lang>.py` | universal |
| CODE-RED-008 | No raw negations in conditions | `01-cross-language/12-no-negatives.md` | `checks/positive-conditions/<lang>.py` | go, ts |

### How CODE-RED-004 and CODE-RED-005 relate (no contradiction)

The two function-length rules form a **coordinated tier** at the same
`error` severity, not competing limits. Effective body line counts map
to outcomes as follows:

| Effective lines | CODE-RED-005 (strict 8) | CODE-RED-004 (hard 15) |
|-----------------|-------------------------|------------------------|
| 0–8             | silent                  | silent                 |
| 9–15            | **error** (build fails) | silent                 |
| 16+             | **error** (build fails) | **error** (redundant)  |

CODE-RED-005 is the **binding** rule — it owns the build-failing
decision for any function over 8 lines. CODE-RED-004 is retained as a
defence-in-depth net so the registry still flags >15-line bodies even
if CODE-RED-005 is ever disabled by a per-file override.

> **Canonical specification of the threshold and counting rules** —
> including what counts as an "effective body line", per-language
> detector scope, deliberate counter divergences, and the verification
> procedure for any future change — lives in
> [`linters-cicd/checks/function-length-prefer8/readme.md`](../../../linters-cicd/checks/function-length-prefer8/readme.md).
> If this table and that README ever disagree, **the README wins** and
> this table is the bug.


---

## Soft preferences (annotate — SARIF `warning`)

*(none currently — CODE-RED-005 was promoted to `error` in v2.0.0 of
this mapping per the strict-8 enforcement decision; this section is
reserved for future soft rules.)*

---

## STYLE rules (annotate — SARIF `warning`)

| ID | Rule | Spec source | Check script | Phase 1 langs |
|----|------|-------------|--------------|---------------|
| STYLE-002 | No `else` after `return`/`throw` | `01-cross-language/04-code-style/` | `checks/no-else-after-return/<lang>.py` | go, ts |

---

## Database rules (block merge — SARIF `error`)

| ID | Rule | Spec source | Check script | Phase 1 langs |
|----|------|-------------|--------------|---------------|
| BOOL-NEG-001 | No Not/No-prefixed boolean columns | `04-database-conventions/02-naming-conventions.md` | `checks/boolean-column-negative/sql.py` | sql |
| DB-FREETEXT-001 | **Presence only** — entity tables need `Description`; transactional tables need `Notes`+`Comments` | `04-database-conventions/03-schema-design.md` §6 | `checks/free-text-columns/sql.py` | sql |
| MISSING-DESC-001 | **Strict superset** — presence (Rules 10/11) + nullability (Rule 12) + waiver mechanism. Recommended for new pipelines. | `04-database-conventions/03-schema-design.md` §6 (v3.4.0) | `checks/missing-desc/sql.py` | sql |

> Both DB-FREETEXT-001 and MISSING-DESC-001 share the same classifier
> via `_lib/free_text_columns.py` so they cannot drift apart. Enable
> **only one** in CI to avoid duplicate findings — MISSING-DESC-001 is
> the recommended choice.

---

## Future rules (Phase 2+)

Added to this table as they ship. Removing a rule requires a major
version bump of the linter pack and a deprecation note in
[`03-language-roadmap.md`](./04-language-roadmap.md).

---

## Changelog

- **v2.0.0** (2026-04-27) — CODE-RED-005 promoted from `warning` to
  `error` per the strict-8 enforcement decision. CODE-RED-004 reframed
  as a redundant >15-line safety net. Added the "How CODE-RED-004 and
  CODE-RED-005 relate" table to make the coordination explicit.
- **v1.0.0** (2026-04-19) — Initial mapping.

---

*Part of [CI/CD Integration](./01-index.md)*
