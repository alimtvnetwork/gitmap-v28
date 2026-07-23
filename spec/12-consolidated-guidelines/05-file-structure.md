# 05 — File & Project Structure

Directory layout, module boundaries, file sizing, and import ordering.

## Group by Responsibility

Organize files by what they do, not what they are.

### Go CLI Layout

```
cmd/        — CLI routing and flag parsing
config/     — Config loading and merging
constants/  — All shared string literals
model/      — Shared data structures
store/      — Database access
```

## One Responsibility Per File

Each file owns a single concern. Split signals:
- File exceeds 200 lines
- 2+ unrelated function groups
- Mixed types and logic
- Mixed constants and components

## Import Ordering

Group imports separated by blank lines:
1. Standard library
2. External packages
3. Internal packages

## Leaf Packages

`model`, `constants`, and `types` are leaf packages. They are imported by many but never import project-specific packages. Breaking this creates circular dependencies.

## File Naming

| Language | Convention | Example |
|----------|-----------|---------|
| TS components | PascalCase | `TypeBadge.tsx` |
| Go files | lowercase, single word | `terminal.go` |
| Specs | kebab-case with prefix | `01-overview.md` |

---

Source: `spec/05-coding-guidelines/05-file-project-structure.md`
