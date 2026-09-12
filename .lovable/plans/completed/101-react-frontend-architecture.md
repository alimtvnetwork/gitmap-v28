# Plan 101: React & Frontend Architecture Audit

## Executive Summary
This master architectural plan establishes repo-wide compliance with authoritative React and frontend architecture guidelines defined in `spec/02-coding-guidelines/02-typescript/`, `spec/07-design-system/`, and `.lovable/coding-guidelines.md`.

## Core Objectives
1. **Named Hook Objects:** Ensure all custom React hooks return named property objects instead of tuple arrays (e.g. `useTheme` returns `{ theme, isDark, source, isSystem, setTheme, toggleTheme }`, `useToast` returns `{ toasts, toast, dismiss }`).
2. **Enum State & Status Values:** Ban raw string unions for state and status in TypeScript; enforce `*Type` suffixed enums (`ThemeType`, `ThemeSourceType`, `TerminalThemeType`).
3. **No Inverted Success Checks:** Ban `!isSuccess` in React and TS components; use affirmative checks or `isFail` / `isFailed`.
4. **Build & Type Safety:** Clean production build via Vite (`npm run build`).

## Verification Results
- `node linter-scripts/check-enum-and-boolean.mjs`: PASS (zero status string unions, zero non-*Type enums, zero inverted success checks)
- `npm run build`: PASS (clean production build in 7.71s, zero errors)
