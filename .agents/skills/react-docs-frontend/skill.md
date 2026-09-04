---
name: react-docs-frontend
description: >-
  Autonomously design, author, refactor, and audit React frontend components in src/ adhering to
  Tailwind CSS tokens, strict TypeScript rules, Result envelopes, and custom hook object returns.
---

# React Docs Frontend Skill

Autonomously implement, refactor, and audit frontend code in `src/` adhering to `spec/07-design-system/`, `spec/08-docs-viewer-ui/`, `spec/24-app-ui-design-system/`, and `spec/17-consolidated-guidelines/31-compiled-simple-coding-guidelines.md`.

## Core Checkpoints & Invariants

1. **Size Limits & Component Modularity:**
   - React component files (`.tsx`): max **100 lines** per file.
   - Non-component files (`.ts`): max **200 lines** per file.
   - Extract child components, dialogs, nav bars, and search modals into dedicated files under `src/components/`.

2. **Strict TypeScript Standards:**
   - Zero `any`, `unknown`, or untyped data structures.
   - String unions (`"light" | "dark"`) are forbidden for persistent states or domain models; use TypeScript Enums.
   - Every Enum name MUST end with the suffix `Type` (e.g., `ThemeType`, `ViewModeType`).
   - Wrap async/fallible operations in `Result<T>` envelopes with explicit `isSuccess` and `isFailure` booleans.

3. **Custom Hook Architecture:**
   - Custom hooks MUST return named object properties, NEVER arrays/tuples (e.g. return `{ theme, isDark, toggleTheme }`, NEVER `[theme, setTheme]`).
   - Hooks must expose positive, affirmative boolean states (`isLoaded`, `isOpen`, `isDark`).

4. **Design System & Theme Tokens:**
   - Adhere to CSS variable-driven tokens: `--primary` is amber gold (`38 92% 50%` light, `41 96% 56%` dark).
   - Use semantic color classes (`bg-primary`, `text-primary`, `--accent-success`); never hardcode raw hex colors or ad-hoc Tailwind colors (e.g., raw `bg-green-500`).
   - Support dark and light theme switching via `useTheme` without flash of unstyled content.

5. **Error & Event Handling:**
   - No silent error swallows in event handlers or `useEffect`.
   - Log errors through centralized notification or error logger with operational context.
   - Trailing newline at EOF; one blank line before `return`.
