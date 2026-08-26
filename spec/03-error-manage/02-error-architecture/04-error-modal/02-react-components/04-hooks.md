# Hooks

> **Parent:** [React Components Index](./00-overview.md)  
> **Version:** 4.0.0  
> **Updated:** 2026-04-01

---

## Review Compliance

| Rule | Status | Notes |
|------|--------|-------|
| Use React Query | ✅ Fixed | Replaced manual `useState`/`useEffect` with `useQuery` |
| No `unknown` | ✅ Clean | All types explicit |
| Function size ≤ 15 lines | ✅ Clean | Hook body is 12 lines |
| Parameters ≤ 3 | ✅ Clean | 1 param |
| Cognitive complexity ≤ 10 | ✅ Clean | Linear flow, no nesting |

---

## `useSessionDiagnostics`

Fetches session diagnostics and logs in parallel using React Query.

```typescript
import { useQuery } from '@tanstack/react-query';
import { errorApi } from './ErrorApiService';
import type { SessionDiagnostics } from './ErrorModalTypes';

interface SessionDiagnosticsResult {
  diagnostics: SessionDiagnostics | null;
  logs: string | null;
  loading: boolean;
  error: string | null;
  refetch: () => void;
}

async function fetchSessionData(sessionId: string): Promise<{
  diagnostics: SessionDiagnostics | null;
  logs: string | null;
}> {
  const [logsRes, diagRes] = await Promise.all([
    errorApi.getSessionLogs(sessionId),
    errorApi.getSessionDiagnostics(sessionId),
  ]);

  return {
    logs: logsRes.success ? logsRes.data?.logs ?? null : null,
    diagnostics: diagRes.success ? diagRes.data ?? null : null,
  };
}

export function useSessionDiagnostics(sessionId?: string): SessionDiagnosticsResult {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['session-diagnostics', sessionId],
    queryFn: () => fetchSessionData(sessionId!),
    enabled: Boolean(sessionId),
    staleTime: 30_000,
    retry: 1,
  });

  return {
    diagnostics: data?.diagnostics ?? null,
    logs: data?.logs ?? null,
    loading: isLoading,
    error: error instanceof Error ? error.message : null,
    refetch,
  };
}
```

### Violations Fixed (v3.2.0 → v4.0.0)

| Previous | Violation | Fix |
|----------|-----------|-----|
| Manual `useState` + `useEffect` | Project has React Query — redundant state management | → `useQuery` with `queryKey` and `enabled` |
| `fetchData` inside component | Not memoized, recreated each render | → Extracted to standalone `fetchSessionData` function |
| `useEffect` missing `fetchData` in deps | React hooks lint violation | → Eliminated; React Query handles lifecycle |
| `err instanceof Error ? ... : "..."` | Inline ternary with string fallback | → Same pattern but cleaner with single error path |
| No return type annotation | Missing explicit return type | → `SessionDiagnosticsResult` interface |
| No stale time / retry config | Re-fetches unnecessarily | → `staleTime: 30_000`, `retry: 1` |

---

*Hooks — updated: 2026-04-01*
