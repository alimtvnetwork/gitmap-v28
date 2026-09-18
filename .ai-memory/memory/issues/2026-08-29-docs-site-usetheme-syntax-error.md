# Memory Issue: 2026-08-29 Docs Site useTheme Syntax Error

## 1. Why it happened

`src/hooks/useTheme.ts:72` had invalid syntax `check!res.isFail` (exclamation point inside variable name), which crashed esbuild during `vite build` in `Build docs-site`.

---

## 2. How it happened

`release.yml` invoked `npm run build` to build the docs-site bundle, which halted when Vite's esbuild transform encountered `check!res`.

---

## 3. Root Cause

Typographical syntax error `check!res` in `useTheme.ts`.

---

## 4. Code Fix

Extracted to positive boolean `const hasUserPreference = checkRes.isSuccess && Boolean(checkRes.data); if (hasUserPreference) return;`. Added `Docs Site Build` to local runner.
