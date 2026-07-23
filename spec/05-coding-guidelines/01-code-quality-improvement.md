# Code Quality Improvement — Universal Coding Guidelines

## Process Name

**Code Quality Improvement** — a systematic approach to enforcing
readability, maintainability, and consistency across all languages
(TypeScript, Go, and any future additions).

---

## 1. No Magic Strings

Every string literal used for comparison, defaults, labels, or keys
must live in a dedicated constants file.

- **Comparison groups** (e.g., tiers, statuses, roles) → use an `enum`
  (TypeScript) or `const` group (Go).
- **Standalone values** (defaults, labels, formats) → named constants.

### Before

```ts
if (tier === "free") {
  return "basic";
}
```

### After

```ts
// constants/tiers.ts
export enum Tier {
  Free = "free",
  Pro = "pro",
}

// usage
if (tier === Tier.Free) {
  return PlanLabel.Basic;
}
```

---

## 2. Exported Object Constants — PascalCase

All exported constant objects use PascalCase names.

### Before

```ts
export const ws_tier_labels = { free: { label: "Free" } };
```

### After

```ts
export const WsTierLabels: Record<Tier, TierStyle> = {
  [Tier.Free]: { label: "Free", bg: "bg-muted", fg: "text-muted-foreground" },
};
```

---

## 3. No Inline Type Definitions — Extract Named Types

Never define types inline. Create a separate, reusable named type.

### Before

```ts
export const TierLabels: Record<string, { label: string; bg: string; fg: string }> = { ... };
```

### After

```ts
// types/tier.ts
export interface TierStyle {
  label: string;
  bg: string;
  fg: string;
}

// constants/tiers.ts
export const TierLabels: Record<Tier, TierStyle> = { ... };
```

---

## 4. Function Length — 8 to 25 Lines

- Target: **8–25 lines** of code (excluding blanks and comments).
- If a function exceeds 25 lines, extract a helper.
- Do **not** cram multiple statements onto one line to bypass the limit.
- Each function does one thing.

### Before

```ts
function processData(items: Item[]) {
  // 40+ lines of mixed filtering, mapping, formatting, and rendering
}
```

### After

```ts
function processData(items: Item[]) {
  const filtered = filterActiveItems(items);
  const mapped = mapToViewModel(filtered);

  return formatOutput(mapped);
}

function filterActiveItems(items: Item[]): Item[] {
  return items.filter((item) => item.isActive);
}
```

---

## 5. Simple Conditionals — No Negation, No Complexity

- **No negation** in `if` conditions: no `!`, no `!=`, no negative
  function names like `isNotValid`.
- **No complex compound conditions** inline. Extract them into a
  well-named boolean function or variable.

### Before

```ts
if (!user.isActive && !(role === "admin" || role === "moderator")) {
  deny();
}
```

### After

```ts
const isRestricted = checkIsRestricted(user, role);
if (isRestricted) {
  deny();
}

function checkIsRestricted(user: User, role: Role): boolean {
  const isInactive = user.isActive === false;
  const isPrivileged = role === Role.Admin || role === Role.Moderator;

  return isInactive && isPrivileged === false;
}
```

> **Note:** Even `=== false` is preferred over `!` for positive-logic
> readability.

---

## 6. Boolean Naming — `is` or `has` Prefix

Every boolean variable, constant, parameter, and function must start
with `is` or `has`.

| Element | Convention | Example |
|---------|-----------|---------|
| Variable | `is` / `has` prefix | `isActive`, `hasPermission` |
| Function | `is` / `has` prefix | `isEligible()`, `hasRole()` |
| Constant | `Is` / `Has` prefix | `IsDebugMode`, `HasFeatureFlag` |

### Before

```ts
const active = user.status === "active";
function eligible(u: User): boolean { ... }
```

### After

```ts
const isActive = user.status === UserStatus.Active;
function isEligible(user: User): boolean { ... }
```

---

## 7. Meaningful Variable Names

Never use single-character or cryptic variable names like `s`, `x`,
`d`, `t`. Every name must convey intent.

| ❌ Wrong | ✅ Correct |
|----------|-----------|
| `s` | `source`, `section` |
| `x` | `index`, `xCoordinate` |
| `d` | `directory`, `duration` |
| `t` | `target`, `timestamp` |
| `cb` | `onComplete`, `handleClick` |

Exception: `i` in a simple `for` loop index is acceptable.

---

## 8. Blank Line Before `return`

Always add a blank line before `return`, unless the `return` is the
only line inside an `if` block.

### Before

```ts
function getLabel(tier: Tier): string {
  const style = TierLabels[tier];
  return style.label;
}
```

### After

```ts
function getLabel(tier: Tier): string {
  const style = TierLabels[tier];

  return style.label;
}
```

---

## 9. Self-Documenting Code

- Code should explain itself through naming and structure.
- If you need a comment to explain a section, that section should be
  its own function.
- Avoid inline comments that restate what the code does.

---

## 10. File Length — Max 200 Lines

No source file should exceed 200 lines. When it does, split by
responsibility into focused files.

---

## Summary Table

| # | Rule | Scope |
|---|------|-------|
| 1 | No magic strings — enums or constants | All languages |
| 2 | Exported object constants — PascalCase | TypeScript |
| 3 | No inline type definitions — extract named types | TypeScript |
| 4 | Function length — 8 to 25 lines | All languages |
| 5 | Simple conditionals — no negation, no complex logic | All languages |
| 6 | Boolean naming — `is` / `has` prefix | All languages |
| 7 | Meaningful variable names | All languages |
| 8 | Blank line before `return` | All languages |
| 9 | Self-documenting code | All languages |
| 10 | File length — max 200 lines | All languages |

---

## Real Refactoring Examples

### A. Inline type → Named type (TypeBadge.tsx)

**Before:**
```ts
export const PROJECT_TYPES: Record<ProjectType, { label: string; color: string; icon: typeof Code2 }> = { ... };
```

**After:**
```ts
// types.ts
export interface ProjectTypeConfig {
  label: string;
  color: string;
  icon: typeof Code2;
}

// TypeBadge.tsx
export const ProjectTypes: Record<ProjectType, ProjectTypeConfig> = { ... };
```

### B. Magic strings → Constants (Watch.tsx)

**Before:**
```ts
const MOCK_REPOS = [
  { name: "api-gateway", status: "clean", ... },
  { name: "frontend-app", status: "dirty", ... },
];
const statusColor = (s: string) => s === "dirty" ? "text-yellow-400" : "text-primary";
```

**After:**
```ts
// constants/index.ts
export const RepoStatus = { Clean: "clean", Dirty: "dirty" } as const;
export type RepoStatus = (typeof RepoStatus)[keyof typeof RepoStatus];

// Watch.tsx
const isDirty = (status: RepoStatus): boolean => status === RepoStatus.Dirty;
const statusColor = (status: RepoStatus) => isDirty(status) ? "text-yellow-400" : "text-primary";
```

### C. Duplicate constant → Single import (ProjectDetailDialog.tsx)

**Before:** `PROJECT_TYPES` was defined identically in both `TypeBadge.tsx` and `ProjectDetailDialog.tsx`.

**After:** Defined once in `TypeBadge.tsx`, imported in `ProjectDetailDialog.tsx`:
```ts
import { ProjectTypes } from "@/components/projects/TypeBadge";
```

### D. Single-char variables → Meaningful names

**Before:**
```ts
const filtered = SAMPLE_PROJECTS.filter((p) => { const q = search.toLowerCase(); ... });
const stash = MOCK_REPOS.reduce((a, r) => a + r.stash, 0);
```

**After:**
```ts
const filtered = SAMPLE_PROJECTS.filter((project) => { const searchLower = search.toLowerCase(); ... });
const stash = MOCK_REPOS.reduce((acc, repo) => acc + repo.stash, 0);
```

### E. Switch → Positive if/else (TerminalDemo.tsx)

**Before:**
```ts
const colorFor = (type?: string) => {
  switch (type) {
    case "input": return "text-[hsl(var(--terminal-foreground))]";
    case "header": return "text-primary font-bold";
    default: return "text-[hsl(var(--foreground))]/70";
  }
};
```

**After:**
```ts
const colorFor = (type?: TerminalLineType) => {
  if (type === TerminalLineType.Input) return "text-[hsl(var(--terminal-foreground))]";
  if (type === TerminalLineType.Header) return "text-primary font-bold";
  if (type === TerminalLineType.Accent) return "text-primary";

  return "text-[hsl(var(--foreground))]/70";
};
```

---

## References

- Go-specific rules: `spec/03-general/06-code-style-rules.md`
- Go CLI rules: `spec/04-generic-cli/08-code-style.md`
- Compliance audit: `spec/01-app/18-compliance-audit.md`
- TypeScript constants: `src/constants/index.ts`
- Project types: `src/components/projects/types.ts`
- Command types: `src/data/commands.ts` (CommandFlag, CommandExample, CommandCategory)
