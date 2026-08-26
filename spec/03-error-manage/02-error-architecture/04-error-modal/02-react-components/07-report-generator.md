# Error Report Generator

> **Parent:** [React Components Index](./00-overview.md)  
> **Version:** 4.0.0  
> **Updated:** 2026-04-01

---

## Review Compliance

| Rule | Status | Notes |
|------|--------|-------|
| No `Record<string, string[]>` | ✅ Fixed | Uses typed `SuggestedFixMap` interface |
| Function size ≤ 15 lines | ✅ Clean | `getSuggestedFixes` is 3 lines |
| Parameters ≤ 3 | ✅ Clean | Max 2 params |

---

A pure function that converts `CapturedError` to Markdown. No React dependencies.

```typescript
export function generateErrorReport(error: CapturedError, app?: AppInfo): string {
  // Sections (in order):
  // 1. App metadata (name, version, git commit, build time)
  // 2. Error identity (ID, code, level, timestamp)
  // 3. Page context (route + React component name)
  // 4. User interaction (arrow-style summary)
  // 5. Trigger context (component, action, source)
  // 6. Invocation chain (indented tree)
  // 7. User interaction path (numbered steps with routes)
  // 8. Target site URL
  // 9. Session info (ID, type, fetch hint)
  // 10. Error message + details
  // 11. Request info (method, endpoint, status, body)
  // 12. Backend execution logs
  // 13. Backend stack trace (Go)
  // 14. PHP stack trace (table)
  // 15. Frontend execution chain
  // 16. Parsed JS stack frames (table)
  // 17. Error location (file:line)
  // 18. Full context JSON
  // 19. Raw frontend stack trace
}
```

**Sample output:** See [Copy Format Samples](../01-copy-formats/00-overview.md) for complete examples.

---

## Suggested Fixes

```typescript
/** Typed map of error code → fix suggestions */
interface SuggestedFixMap {
  [errorCode: string]: string[];
}

const SUGGESTED_FIXES: SuggestedFixMap = {
  E1001: ["Check backend server is running", "Verify VITE_API_URL", "Check firewall", "Refresh page"],
  E2001: ["Check site credentials", "Verify WordPress accessibility", "Check REST API", "Check plugin"],
  E3001: ["Check plugin files exist", "Verify file permissions", "Check PHP file headers"],
  E4001: ["Check disk space", "Verify PHP upload limits", "Test with smaller plugin"],
  E5001: ["Check plugin for fatal errors", "Verify dependencies", "Check debug.log", "Try manual activation"],
  E9005: ["API returned HTML instead of JSON", "Check backend server", "Verify VITE_API_URL", "Check network tab"],
};

const DEFAULT_FIXES: string[] = ["Check error details", "Review stack trace", "Check backend logs", "Retry"];

export function getSuggestedFixes(code: string): string[] {
  return SUGGESTED_FIXES[code] ?? DEFAULT_FIXES;
}
```

### Violations Fixed (v3.2.0 → v4.0.0)

| Previous | Violation | Fix |
|----------|-----------|-----|
| `Record<string, string[]>` inline | Anonymous type | → `SuggestedFixMap` named interface |
| Inline fallback array | Recreated every call | → `DEFAULT_FIXES` constant |
| `fixes[code] \|\|` | Falsy check instead of nullish | → `??` nullish coalescing |

---

*Report generator — updated: 2026-04-01*
