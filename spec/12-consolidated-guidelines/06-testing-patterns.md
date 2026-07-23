# 06 — Testing Patterns

Standards for writing, naming, and organizing tests.

## Test File Placement

Tests live next to the code they verify: `filename_test.go`, `Component.test.tsx`.

## Test Naming

Go: `Test<Function>_<scenario>`. TS: `describe` for unit, `it` for behavior.

## Structure — Arrange, Act, Assert

Every test follows three phases separated by blank lines.

## Table-Driven Tests (Go)

Use for functions with multiple input/output combinations via struct slices and `t.Run`.

## What to Test

| Category | Priority |
|----------|----------|
| Pure functions | High |
| Edge cases | High |
| Error paths | High |
| Constants completeness | Medium |
| Component render | Medium |

Do NOT test: implementation details, third-party internals, trivial getters.

## Test Isolation

Each test is independent. No shared mutable state. Clean up side effects.

## Assertion Style

Go: `require` for hard stops, `assert` for soft checks. TS: specific matchers over `toBe(true)`.

## Coverage Targets

| Scope | Target |
|-------|--------|
| Pure utilities | 90%+ |
| Core business logic | 80%+ |
| CLI commands | Happy + error paths |

## Test Quality

All coding guidelines apply to test code: meaningful names, functions under 25 lines, no magic strings, boolean naming.

---

Source: `spec/05-coding-guidelines/06-testing-patterns.md`
