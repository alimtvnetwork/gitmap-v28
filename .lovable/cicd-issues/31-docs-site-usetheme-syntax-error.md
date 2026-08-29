# 31-docs-site-usetheme-syntax-error

## Error Summary
In the `Build and Release` workflow (`release.yml`), step `Build docs-site` failed during `npm run build` (`vite build`):
```text
[vite:esbuild] Transform failed with 1 error:
/home/runner/work/gitmap-v28/gitmap-v28/src/hooks/useTheme.ts:72:16: ERROR: Expected ")" but found "res"
file: /home/runner/work/gitmap-v28/gitmap-v28/src/hooks/useTheme.ts:72:16

Expected ")" but found "res"
70 |        // Only follow the OS when the user hasn't explicitly chosen a theme.
71 |        const checkRes = queryWrapperSync(() => localStorage.getItem(THEME_STORAGE_KEY));
72 |        if (check!res.isFail && checkRes.data) return;
   |                  ^
```

---

## 4-Part Root Cause Analysis

### 1. Why it happened
In `src/hooks/useTheme.ts`, line 72 contained a typo `check!res.isFail` where the negation operator `!` was inadvertently placed inside the variable name instead of before it (`!checkRes.isFail`). The TypeScript parser / esbuild encountered invalid syntax `check!res`.

### 2. How it happened
- During theme synchronization logic inside `useTheme.ts`, a check was written to detect if a stored theme preference was present.
- A typo embedded the exclamation point within the identifier (`check!res`), causing Vite's esbuild transform to fail with `Expected ")" but found "res"`.

### 3. Root Cause
Typographical syntax error `check!res` in `src/hooks/useTheme.ts:72`.

### 4. Code Fix
- Rewrote the check in `src/hooks/useTheme.ts` using positive boolean logic complying with project standards:
  ```ts
  const checkRes = queryWrapperSync(() => localStorage.getItem(THEME_STORAGE_KEY));
  const hasUserPreference = checkRes.isSuccess && Boolean(checkRes.data);
  if (hasUserPreference) {
    return;
  }
  ```
- Verified `npm run build` passes cleanly: `built in 6.44s`.
- Added `Docs Site Build` check to `.lovable/ai-fix-scripts/03-cicd-local-runner.py`.
