# Canonical Size Tier (Single Source of Truth)

> **This file is the ONLY authoritative source for function-length, file-length, and component-size limits.**
> All other locations (`.cursorrules`, `eslint.config.js`, `linters-cicd/`, `spec/13-generic-cli/08-code-style.md`, `spec/17-consolidated-guidelines/31-compiled-simple-coding-guidelines.md`) mirror this table and MUST reference it. If any of those drift, this file wins and the others get patched.

## Tier

| Metric                    | Limit         | Enforcement           | Rule ID       |
|---------------------------|---------------|-----------------------|---------------|
| Function body (preferred) | ≤ 8 lines     | warn                  | CODE-RED-005  |
| Function body (hard cap)  | ≤ 15 lines    | error (build-fails)   | CODE-RED-004  |
| File length               | ≤ 300 lines   | error                 | CODE-RED-006  |
| React component file      | ≤ 100 lines   | error (`.tsx` only)   | CODE-RED-006R |
| Struct / class            | ≤ 120 lines   | error                 | CODE-RED-017  |
| Parameters per function   | ≤ 3           | error                 | CODE-RED-008  |
| Cognitive complexity      | ≤ 10          | error                 | CODE-RED-CC10 |

## Counting rules

- Line counts **skip** blank lines and pure-comment lines.
- Function signature line is **not** counted; body lines are.
- Error-handling scaffold (`if err != nil { return apperror.Wrap(err) }` in Go, `catch (e) { throw apperror.wrap(e) }` in TS) is **not** counted.
- There is no soft "hard-max 400" file limit. 300 is the single file cap.

## Waivers

Use only when a genuine domain reason exists (e.g. exhaustive switch on a codegen enum).

```ts
// lint-allow: function-length reason="exhaustive switch over generated enum" max=42
function mapKind(k: Kind): Label { ... }
```

- `reason="..."` is **required** and must be human-readable.
- `max=N` bounds the waiver; the function still fails if it grows past `N`.
- Waivers apply to both `max-function-lines` (15 cap) and `prefer-function-lines` (8 warn).

## Cross-references (must match this tier)

| File | Role |
|------|------|
| `.cursorrules` (Quick Rule 4) | Editor / AI reminder |
| `eslint.config.js` | JS/TS enforcement (`max-lines`, `max-function-lines`, `prefer-function-lines`, `.tsx` override) |
| `linters-cicd/checks/file-length/` | Language-agnostic CI check |
| `linters-cicd/checks/function-length-prefer8/` | Language-agnostic prefer-8 check |
| `spec/13-generic-cli/08-code-style.md` | CLI code-style ref |
| `spec/17-consolidated-guidelines/31-compiled-simple-coding-guidelines.md` §"File Size Limits" | Consolidated guidelines mirror |

## Change protocol

1. Edit this file first.
2. Update the mirrors above in the **same commit**.
3. Bump patch version and add an entry to `changelog.md` under "Canonical tier change".
