# Cross-Spec Contradiction Check Process

> **Version:** 1.0.0
> **Updated:** 2026-03-09
> **Purpose:** Prevent conflicting rules between spec files

---

## 1. Problem Statement

Specifications are authored across dozens of files by multiple contributors (human and AI). Rules in one file can silently contradict rules in another file, creating confusion and inconsistent code. Examples of past contradictions:

| Issue | File A says... | File B says... | Tracking |
|-------|---------------|---------------|----------|
| P7 inline assignment | "All inline `if` assignments banned" | "`if err := fn(); err != nil` is correct Go" | Issue #08 <!-- external: spec/23-how-app-issues-track/08-p7-inline-assignment-contradiction.md --> |
| Log key casing | "camelCase log keys" (v9.x) | "PascalCase for ALL keys" (v10.x) | Fixed in key-naming-pascalcase.md |

---

## 2. Mandatory Contradiction Check Categories

Every consistency report scan **MUST** include these contradiction checks:

### 2.1 — Naming Convention Conflicts

| Check | What to Compare |
|-------|----------------|
| Key casing | All files defining JSON/log/config key format must agree on PascalCase |
| Abbreviation casing | All files using `Id`/`ID`, `Url`/`URL`, `Db`/`DB` must use first-letter-only caps |
| Variable naming | Boolean prefix rules (`is`/`has`) must be consistent across Go, PHP, TS specs |
| File naming | PascalCase source files must be enforced in all language specs |

### 2.2 — Control Flow Conflicts

| Check | What to Compare |
|-------|----------------|
| Inline `if` statements | P7 rule (banned) vs code examples using `if err := fn(); err != nil` (exempt) |
| Type assertions | §7.2 (banned in business logic) vs code examples using `.(Type)` |
| Error handling | `apperror.Result[T]` mandate vs code examples returning `(T, error)` |
| Negation rules | P1 (positive naming) — now enforced: `isDefined()`, `hasContent()` replace old `isNotNil()`, `isNotEmpty()` |

### 2.3 — Type Safety Conflicts

| Check | What to Compare |
|-------|----------------|
| `interface{}` / `any` | Strict typing spec (banned) vs code examples using `map[string]any` |
| Multi-return | Single return rule vs code examples returning `(T, error)` |
| Raw `error` | I-2 invariant vs struct fields using raw `error` type |

### 2.4 — Architecture Conflicts

| Check | What to Compare |
|-------|----------------|
| Function size | 15-line limit vs code examples exceeding it |
| Parameter count | 3-param max vs code examples with more |
| Guard patterns | `IsDefined()` / `IsDefinedAndValid()` mandate vs raw nil checks |

---

## 3. How to Run a Contradiction Scan

### Step 1: Identify the Rule Source Files

These are the **canonical rule sources** (if a conflict is found, these win):

| Rule Domain | Canonical Source |
|-------------|-----------------|
| Naming (all languages) | `spec/02-coding-guidelines/01-cross-language/15-master-coding-guidelines/00-overview.md` |
| Key naming (PascalCase) | `spec/02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md` |
| Boolean logic | `spec/02-coding-guidelines/03-golang/02-boolean-standards.md` + `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/00-overview.md` |
| Strict typing | `spec/02-coding-guidelines/01-cross-language/13-strict-typing.md` |
| Error handling | `spec/03-error-manage/01-error-resolution/10-apperror-package/01-apperror-reference.md` |
| Code style | `spec/02-coding-guidelines/01-cross-language/04-code-style/` |
| No negatives | `spec/02-coding-guidelines/01-cross-language/12-no-negatives.md` |
| Function naming | `spec/02-coding-guidelines/01-cross-language/10-function-naming.md` |
| Database naming | `spec/02-coding-guidelines/01-cross-language/07-database-naming.md` |

### Step 2: Search for Contradictions

For each rule domain, search ALL spec files for code examples or prose that contradict the canonical source:

```bash

# Example: find camelCase log keys that should be PascalCase

grep -rn '"[a-z][a-zA-Z]*"' spec/ --include="*.md" | grep -i "logkey\|log.*key\|context.*key"

# Example: find type assertions in business logic examples

grep -rn '\.\(\*\?[A-Z]' spec/ --include="*.md" | grep -v "EXEMPTED\|Exempt\|stdlib\|test"

# Example: find raw nil checks that should use IsDefined()

grep -rn 'if.*!= nil' spec/ --include="*.md" | grep -v "err\|error\|recover\|ok"

# Example: find negative boolean helpers

grep -rn 'isNot[A-Z]\|hasNo[A-Z]' spec/ --include="*.md" | grep -v "FORBIDDEN\|WRONG\|❌"
```

### Step 3: Classify Each Finding

| Classification | Action |
|---------------|--------|
| **True contradiction** — rule A says X, rule B says Y | Fix the non-canonical file. Create issue in `spec/23-how-app-issues-track/` |
| **Stale example** — code example predates a new rule | Update the example to match current rules |
| **Missing exemption** — valid pattern not covered by rule | Add exemption to the canonical source with `// EXEMPTED:` annotation |
| **False positive** — appears in "❌ WRONG" / "FORBIDDEN" example | No action needed |

### Step 4: Document & Track

1. Create issue write-up at `spec/23-how-app-issues-track/{NN}-{slug}.md`
2. Update the canonical source if an exemption is needed
3. Fix all non-canonical files
4. Update the global `spec/99-consistency-report.md`
5. Update `.lovable/memories/workflow/03-mistake-remediation-protocol.md`

---

## 4. When to Run

| Trigger | Required? |
|---------|-----------|
| New rule added to any canonical source | ✅ Mandatory |
| Consistency report refresh (any `99-consistency-report.md`) | ✅ Mandatory |
| Before marking a remediation issue as "Done" | ✅ Mandatory |
| After large-scale code example updates (e.g., abbreviation sweep) | ✅ Mandatory |
| Ad-hoc review by AI or human | Recommended |

---

## 5. Contradiction Prevention (New Rule Checklist)

Before adding any new rule to a canonical source:

- [ ] **Search existing code examples** — `grep -rn` for patterns the new rule would ban
- [ ] **Count impact** — How many files contain the pattern?
- [ ] **Check exemptions** — Does the pattern have legitimate uses? (e.g., `if err := fn(); err != nil` is idiomatic Go)
- [ ] **Cross-validate** — Read ALL canonical sources for conflicting statements
- [ ] **Update all prompts** — `spec/02-spec-management-software/12-prompts/` must reflect the new rule
- [ ] **Version bump** — Canonical source and global consistency report

---

## 6. Known Contradiction-Prone Areas

These areas have historically produced contradictions and need extra scrutiny:

1. **Go inline `if` patterns** — P7 bans inline statements, but Go has 4 exempt idioms
2. **Boolean helper naming** — P1 mandates positive naming; resolved: `isDefined()` / `hasContent()` replace all negative helpers
3. **Log key casing** — Recently changed from camelCase to PascalCase; old examples persist
4. **Type assertions** — §7.2 bans them, but many spec code examples still use `.(Type)`
5. **Error return patterns** — Single-return mandate vs `(T, error)` examples
6. **`ctx.Value()` casting** — Raw `ctx.Value("key").(string)` vs typed accessor requirement

---

## Cross-References

- Global Consistency Report <!-- external: spec/99-consistency-report.md -->
- Mistake Remediation Protocol <!-- external: .lovable/memories/workflow/03-mistake-remediation-protocol.md -->
- Issue Template <!-- external: spec/23-how-app-issues-track/01-issue-template.md -->
- [Master Coding Guidelines](./15-master-coding-guidelines/00-overview.md)
- [PascalCase Key Naming](./11-key-naming-pascalcase.md)
- [Boolean Standards](../03-golang/02-boolean-standards.md)
