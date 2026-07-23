# Naming Conventions — Cross-Language Reference

## Overview

Consistent naming across TypeScript and Go for variables, functions,
constants, types, and booleans.

---

## 1. Variables — Meaningful Names Only

Never use single-character or cryptic names.

| ❌ Wrong | ✅ Correct |
|----------|-----------|
| `s` | `source`, `section` |
| `x` | `index`, `xCoordinate` |
| `d` | `directory`, `duration` |
| `t` | `target`, `timestamp` |
| `cb` | `onComplete`, `handleClick` |
| `a`, `r` | `accumulator`, `repo` |

Exception: `i` in a simple `for` loop is acceptable.

---

## 2. Booleans — `is` or `has` Prefix

Every boolean variable, constant, parameter, and function must start
with `is` or `has`.

### TypeScript

```ts
const isActive = user.status === UserStatus.Active;
function isEligible(user: User): boolean { ... }
const hasPermission = checkHasPermission(user, resource);
```

### Go

```go
isValid := validate(input)
func hasRole(userID string, role Role) bool { ... }
```

---

## 3. Exported Constants — PascalCase

### TypeScript

```ts
// ✅ Correct
export const ProjectTypes: Record<ProjectType, ProjectTypeConfig> = { ... };

// ❌ Wrong
export const PROJECT_TYPES = { ... };
export const project_types = { ... };
```

### Go

```go
// Go constants are always PascalCase when exported
const DefaultBranch = "main"
const ModeHTTPS = "https"
```

---

## 4. Functions — Verb-Led

| Language | Exported | Unexported |
|----------|----------|------------|
| TypeScript | `camelCase` | `camelCase` |
| Go | `PascalCase` | `camelCase` |

All function names start with a verb: `get`, `set`, `build`, `format`,
`check`, `is`, `has`, `parse`, `resolve`, `write`, `read`.

---

## 5. Types and Interfaces

### TypeScript

- Interfaces: PascalCase, noun-based (`ProjectTypeConfig`, `TierStyle`).
- Enums: PascalCase with PascalCase members.
- Type aliases: PascalCase (`RepoStatus`, `CommandCategory`).

### Go

- Structs: PascalCase (`ScanRecord`, `CloneResult`).
- Interfaces: PascalCase, `-er` suffix when applicable (`Formatter`).

---

## 6. Files

| Language | Convention | Example |
|----------|-----------|---------|
| TypeScript | camelCase or kebab-case | `specData.ts`, `type-badge.tsx` |
| Go | lowercase, single word | `terminal.go`, `csv.go` |

---

## References

- Code Quality Improvement: `spec/05-coding-guidelines/01-code-quality-improvement.md`
- Go Code Style: `spec/05-coding-guidelines/02-go-code-style.md`
