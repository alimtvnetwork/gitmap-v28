# 03 — Naming Conventions

Cross-language naming standards for variables, functions, constants, types, and files.

## Variables — Meaningful Names Only

Never use single-character or cryptic names.

| Wrong | Correct |
|-------|---------|
| `s` | `source`, `section` |
| `x` | `index`, `xCoordinate` |
| `d` | `directory`, `duration` |
| `t` | `target`, `timestamp` |
| `a`, `r` | `accumulator`, `repo` |

Exception: `i` in a simple `for` loop is acceptable.

## Booleans — `is` or `has` Prefix

Every boolean variable, constant, parameter, and function must start with `is` or `has`.

## Exported Constants — PascalCase

TS: `export const ProjectTypes = { ... }` — never `PROJECT_TYPES` or `project_types`.
Go: `const DefaultBranch = "main"`.

## Functions — Verb-Led

| Language | Exported | Unexported |
|----------|----------|------------|
| TypeScript | camelCase | camelCase |
| Go | PascalCase | camelCase |

Start with: `get`, `set`, `build`, `format`, `check`, `is`, `has`, `parse`, `resolve`, `write`, `read`.

## Types and Interfaces

TS: PascalCase nouns (`ProjectTypeConfig`, `TierStyle`). Go: PascalCase structs, `-er` suffix interfaces.

## Files

| Language | Convention | Example |
|----------|-----------|---------|
| TypeScript components | PascalCase | `TypeBadge.tsx` |
| Go source files | lowercase, single word | `terminal.go` |
| Spec documents | kebab-case with numeric prefix | `01-overview.md` |

## Database Naming

All table names and column names use **PascalCase** — never snake_case. Example: `DetectedProjects`, `ProjectTypeId`, `CreatedAt`.

---

Source: `spec/05-coding-guidelines/03-naming-conventions.md`
