# File & Folder Naming — TypeScript / JavaScript (AI Execution Prompt)

> **/goal** Enforce uniform TypeScript and JavaScript file and folder naming (`kebab-case.ts`, `kebab-case/` folders, test suffixes).
> **/learn** Master module path ergonomics, barrel file conventions (`index.ts`), test suffixes (`.test.ts`, `.spec.ts`), and lowercase directory hygiene.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce `kebab-case.ts` for all utilities, hooks (`use-*.ts`), types, and test files (`*.test.ts`).
- [ ] `/learn` Never use PascalCase or camelCase directory names; all folders MUST be lowercase kebab-case.
- [ ] `/goal` Maintain consistent barrel exports via `index.ts` within bounded modular directories.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Overview

TypeScript/JavaScript projects follow framework conventions (React, Node.js, etc.) with consistent kebab-case for most files and PascalCase for React components.

---

## File Naming Rules

### 1. General Files — `kebab-case.ts`

```
✅ api-client.ts
✅ use-auth.ts
✅ date-utils.ts
❌ apiClient.ts
❌ api_client.ts
❌ ApiClient.ts
```

### 2. React Components — `PascalCase.tsx`

```
✅ UserCard.tsx
✅ AdminSettings.tsx
✅ NavigationMenu.tsx
❌ user-card.tsx
❌ userCard.tsx
```

### 3. Hooks — `use-{name}.ts`

```
✅ use-auth.ts
✅ use-mobile.ts
✅ use-toast.ts
❌ useAuth.ts
❌ UseAuth.ts
```

### 4. Test Files — `*.test.ts` or `*.spec.ts`

```
✅ api-client.test.ts
✅ UserCard.test.tsx
✅ api-client.spec.ts
```

### 5. Type/Interface Files — `kebab-case.types.ts`

```
✅ api.types.ts
✅ user.types.ts
✅ error-codes.types.ts
```

### 6. Constants/Config — `kebab-case.ts`

```
✅ app-config.ts
✅ route-constants.ts
✅ error-messages.ts
```

---

## Folder Naming Rules

### 1. All Folders — `kebab-case`

```
✅ src/components/
✅ src/hooks/
✅ src/lib/
✅ src/pages/
❌ src/Components/
❌ src/myHooks/
```

### 2. Standard React/Vite Layout

```
src/
├── components/              ← kebab-case folders
│   ├── ui/                  ← shadcn components
│   │   ├── button.tsx
│   │   └── dialog.tsx
│   ├── NavLink.tsx          ← PascalCase component files
│   └── UserCard.tsx
├── hooks/                   ← custom hooks
│   ├── use-auth.ts
│   └── use-mobile.ts
├── lib/                     ← utilities
│   └── utils.ts
├── pages/                   ← route pages
│   ├── Index.tsx
│   └── NotFound.tsx
├── types/                   ← shared types
│   └── api.types.ts
└── App.tsx
```

### 3. Index Files

Use `index.ts` for barrel exports:

```
✅ components/ui/index.ts
✅ hooks/index.ts
```

---

## Forbidden Patterns

| Pattern | Why |
|---------|-----|
| `snake_case.ts` | Not JS/TS convention |
| `SCREAMING_CASE.ts` | Reserved for env files only |
| PascalCase folders | `Components/`, `Hooks/` — always lowercase kebab-case |
| Spaces in filenames | Breaks imports |
| `.jsx` for TypeScript | Use `.tsx` when using TypeScript |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| TypeScript Standards | [../02-typescript/readme.md](../02-typescript/readme.md) |
| Cross-Language Rules | [./02-cross-language.md](./02-cross-language.md) |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/08-file-folder-naming/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-FILE-005: TypeScript and JavaScript File Naming Standards

**Given** Repository file and directory structures across polyglot stacks.
**When** Audited against this file and folder naming specification.
**Then** Zero uppercase or invalid naming patterns are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only
```
**Expected:** exit 0. Zero violations.
