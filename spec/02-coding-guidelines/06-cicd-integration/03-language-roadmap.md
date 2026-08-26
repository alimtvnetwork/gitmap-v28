# Language Roadmap

> **Version:** 1.0.0
> **Updated:** 2026-04-19

Phased rollout per user direction. Each phase is independently shippable
and adds **only** new plugins under `linters-cicd/checks/<rule>/<lang>.py`
plus registry entries.

---

## Phase 1 — Go + TypeScript + PHP (✅ shipping in v3.9.0)

All three currently-specced languages ship together so consumers with
mixed Go/TS/PHP repos (the common Riseup Asia stack) get full coverage
on day one — no waiting for a follow-up release.

| Check | Go | TypeScript | PHP |
|-------|----|------------|-----|
| nested-if (CODE-RED-001) | ✅ regex+AST hybrid | ✅ regex+AST hybrid | ✅ regex+`phply` hybrid |
| function-length (CODE-RED-004) | ✅ | ✅ | ✅ |
| file-length (CODE-RED-006) | ✅ | ✅ | ✅ |
| magic-strings (CODE-RED-003) | ✅ | ✅ | ✅ (WordPress-aware allowlist) |
| boolean-naming (CODE-RED-002) | ✅ | ✅ | ✅ |
| positive-conditions (CODE-RED-008) | ✅ | ✅ | ✅ |
| no-else-after-return (STYLE-002) | ✅ | ✅ | ✅ |

**Why these three first:**
- **Go + TS** are the languages used in this repo, so the checks self-test
  against the spec's own corpus.
- **PHP** rules already exist in `spec/02-coding-guidelines/04-php/` and
  the WordPress/Laravel projects under Riseup Asia LLC need parity with
  the Go/TS enforcement on day one. PHP plugins use `phply` for AST and
  regex fallbacks for legacy WordPress code that fails to parse.

---

## Phase 2 — Python + Rust (planned)

Python uses the standard library `ast` module — trivial. Rust uses
`tree-sitter-rust` Python bindings.

---

## Phase 4+ — On request

Any additional language (Java, Kotlin, Swift, C#, …) is added on user
request following [`02-plugin-model.md`](./02-plugin-model.md). The
orchestrator and SARIF contract remain unchanged.

---

## Promotion criteria (todo → shipping)

A language graduates from "planned" to "shipping" when:

1. All 7 Phase 1 checks have a working plugin with fixtures.
2. `linters-cicd/checks/<rule>/fixtures/<lang>/` has ≥ 1 bad and ≥ 1 good
   fixture per check.
3. CI runs `validate-sarif.py` on every emission.
4. The rule appears in `06-rules-mapping.md` with status `shipping`.

---

*Part of [CI/CD Integration](./00-overview.md)*
