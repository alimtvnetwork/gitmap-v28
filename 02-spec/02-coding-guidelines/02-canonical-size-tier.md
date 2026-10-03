# Canonical Size Tier (Single Source of Truth) (AI Execution Prompt)

> **/goal** Enforce canonical sizing limits across all codebases (≤8 lines preferred / ≤15 lines hard cap per function, ≤300 lines per file, ≤100 lines per React component, ≤120 lines per struct/class, ≤3 parameters, ≤10 cognitive complexity).
> **/learn** Master single-responsibility decomposition, early return guard clauses, options parameter structs, and strict avoidance of monolithic files or unconstrained function bloat.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Keep all function bodies within ≤8 lines (preferred) and never exceed ≤15 lines (hard cap build error).
- [ ] `/learn` Never allow file lengths to exceed ≤300 lines (or ≤100 lines for React TSX components); immediately decompose oversized units.
- [ ] `/goal` Limit function parameters to ≤3; refactor 4+ arguments into dedicated parameter options structs or objects.
- [ ] `/learn` Verify zero sizing and complexity drift across linters and CI/CD quality gates with explicit waivers where domain-justified.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **This file is the ONLY authoritative source for function-length, file-length, and component-size limits.**
> All other locations (`.cursorrules`, `eslint.config.js`, `linters-cicd/`, `02-spec/13-generic-cli/08-code-style.md`, `02-spec/17-consolidated-guidelines/34-compiled-simple-coding-guidelines.md`) mirror this table and MUST reference it. If any of those drift, this file wins and the others get patched.

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
| `02-spec/13-generic-cli/08-code-style.md` | CLI code-style ref |
| `02-spec/17-consolidated-guidelines/34-compiled-simple-coding-guidelines.md` §"File Size Limits" | Consolidated guidelines mirror |

## Change protocol

1. Edit this file first.
2. Update the mirrors above in the **same commit**.
3. Bump patch version and add an entry to `changelog.md` under "Canonical tier change".

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-ROOT-002: Canonical Size Tier Enforcement

**Given** Source files in the repository.
**When** Codebases are audited by coding guideline scanners and lint rules.
**Then** Function lengths (≤8 preferred, ≤15 hard cap), file lengths (≤300 lines, TSX ≤100 lines), struct lengths (≤120 lines), and parameter counts (≤3) are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only
```
**Expected:** exit 0. Zero violations.
