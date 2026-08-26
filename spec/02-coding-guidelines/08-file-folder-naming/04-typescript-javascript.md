# File & Folder Naming — TypeScript / JavaScript

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
| TypeScript Standards | [../02-typescript/00-overview.md](../02-typescript/00-overview.md) |
| Cross-Language Rules | [./01-cross-language.md](./01-cross-language.md) |
