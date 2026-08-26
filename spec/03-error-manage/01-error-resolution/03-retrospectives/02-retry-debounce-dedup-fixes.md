# Retry, Debounce & Deduplication Fixes — Full Retrospective

> **Created:** 2026-02-12  
**Version:** 3.2.0  
> **Purpose:** Document every retry/debounce/dedup issue that was found and fixed, with root causes, symptoms, and exact code solutions — for AI training.

---

## Table of Contents

1. [QueryClient Automatic Retry Loop](#1-queryclient-automatic-retry-loop)
2. [Window Focus Refetch Storm](#2-window-focus-refetch-storm)
3. [Publish Double-Invocation (API-Level Dedup Lock)](#3-publish-double-invocation-api-level-dedup-lock)
4. [Post-Publish Cooldown Guard](#4-post-publish-cooldown-guard)
5. [WebSocket Event Listener Duplication (PublishProgressDialog)](#5-websocket-event-listener-duplication-publishprogressdialog)
6. [Toast Notification Deduplication](#6-toast-notification-deduplication)
7. [Circuit Breaker for Failing Endpoints](#7-circuit-breaker-for-failing-endpoints)
8. [Snapshot Query Retry Suppression](#8-snapshot-query-retry-suppression)
9. [Anti-Patterns Summary (Never Do This)](#9-anti-patterns-summary-never-do-this)

---

## 1. QueryClient Automatic Retry Loop

### Symptom

Failed API requests were automatically retried 3 times (React Query default), causing:
- Triple error toasts for a single failure
- Unnecessary load on backend during outages
- Confusing UX where errors appeared multiple times

### Root Cause

React Query's `QueryClient` defaults to `retry: 3` and `refetchOnWindowFocus: true`. When the backend was down or returned errors, every query would silently retry 3 times, and every time the user tabbed back to the app, all stale queries would refetch.

### Fix — `src/App.tsx`

```typescript
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      retry: false,              // ← CRITICAL: No automatic retries
      refetchOnWindowFocus: false, // ← CRITICAL: No background refetches
    },
  },
});
```

### Key Principle

**Data should only refresh through explicit user actions** (button clicks, manual refresh, navigation). Never rely on React Query's automatic retry/refetch — it creates invisible network storms and duplicate errors.

---

## 2. Window Focus Refetch Storm

### Symptom

Every time the user switched tabs and returned, ALL queries would refetch simultaneously, causing:
- Sudden burst of 10-20+ API requests
- Stale error modals reappearing
- Backend rate limiting

### Root Cause

`refetchOnWindowFocus: true` (React Query default) triggers a refetch for every mounted query when `document.visibilitychange` fires.

### Fix

Set `refetchOnWindowFocus: false` globally (see Section 1) and also per-query for sensitive hooks:

```typescript
// src/hooks/useRemoteSnapshots.ts
{
  retry: false,
  refetchOnWindowFocus: false,
  meta: { suppressGlobalError: true },
}
```

### Key Principle

Pair `refetchOnWindowFocus: false` with `suppressGlobalError: true` in query meta for background/polling queries that should never trigger the global error modal.

---

## 3. Publish Double-Invocation (API-Level Dedup Lock)

### Symptom

Users could trigger the publish function multiple times by:
- Double-clicking the publish button
- Auto-publish triggering while a manual publish was in-flight
- WebSocket reconnection re-triggering the publish

This caused duplicate ZIP uploads, duplicate activity logs, and race conditions on the remote WordPress site.

### Root Cause

No guard existed at the API method level. The publish function was a simple `request()` call that could be invoked any number of times concurrently.

### Fix — `src/lib/api/methods.ts`

```typescript
publishPlugin: (() => {
  const inFlight = new Set<string>();
  const cooldowns = new Map<string, number>();
  const COOLDOWN_MS = 30_000; // 30s cooldown after success

  return (pluginId, siteId, options) => {
    const key = `${pluginId}:${siteId}`;

    // Guard 1: Block if already in-flight
    if (inFlight.has(key)) {
      console.warn(`[api.publishPlugin] BLOCKED duplicate in-flight request`);

      return Promise.resolve({
        success: false,
        error: { code: "E_DEDUP", message: "A publish is already in progress" },
      });
    }

    // Guard 2: Block if within cooldown period
    const lastSuccess = cooldowns.get(key);

    if (lastSuccess && Date.now() - lastSuccess < COOLDOWN_MS) {
      const secsLeft = Math.ceil((COOLDOWN_MS - (Date.now() - lastSuccess)) / 1000);

      return Promise.resolve({
        success: false,
        error: { code: "E_COOLDOWN", message: `Cooldown active (${secsLeft}s remaining)` },
      });
    }

    inFlight.add(key);

    return request(endpoint, { method: HttpMethod.Post, body }).then(response => {
      if (response.success) cooldowns.set(key, Date.now());

      return response;
    }).finally(() => {
      inFlight.delete(key);
    });
  };
})(),
```

### Architecture Details

- **IIFE closure** wraps the function so `inFlight` and `cooldowns` are module-scoped singletons
- **`E_DEDUP` error code** — returned locally, never reaches the network
- **`E_COOLDOWN` error code** — prevents rapid re-publish after success
- **30-second cooldown** — prevents the auto-publish re-triggering immediately after a manual publish completes
- **`.finally()`** — ensures `inFlight` is always cleaned up, even on network errors

### Key Principle

**Critical mutating operations must have an API-layer dedup lock**, not just UI-level button disabling. UI guards can be bypassed by programmatic callers (auto-publish, WebSocket event handlers).

---

## 4. Post-Publish Cooldown Guard

### Symptom

After a successful publish, the auto-publish watcher would detect the new files (uploaded by the publish itself) and immediately trigger another publish cycle, creating an infinite re-publish loop.

### Root Cause

The file watcher detected the ZIP files and version.json changes created by the publish as "new changes" and queued another auto-publish.

### Fix

The 30-second cooldown in the `publishPlugin` IIFE (Section 3) blocks any re-invocation within the cooldown window. Combined with the `E_COOLDOWN` error code, the auto-publish handler can silently discard the blocked attempt.

---

## 5. WebSocket Event Listener Duplication (PublishProgressDialog)

### Symptom

The PublishProgressDialog showed:
- Duplicate or triple log entries for the same publish step
- Progress bar jumping erratically (e.g., 30% → 60% → 30% → 90%)
- Late-arriving "complete" events from previous publishes resetting the UI

### Root Cause

The `useEffect` that subscribed to WebSocket events had **unstable dependencies** (callbacks, computed labels) that changed on every render, causing:
1. Effect re-runs → re-subscriptions → multiple active listeners for the same event
2. No deduplication of log entries
3. No completion lock — late events from finished publishes were still processed

### Fix — `src/components/plugins/PublishProgressDialog.tsx`

**Pattern 1: Stabilize dependencies with `useRef`**
```typescript
// Move unstable deps into refs — they update every render but don't trigger effect re-runs
const onCompleteRef = useRef(onComplete);
onCompleteRef.current = onComplete;
const pluginNameRef = useRef(pluginName);
pluginNameRef.current = pluginName;
const siteNameRef = useRef(siteName);
siteNameRef.current = siteName;
```

**Pattern 2: Completion lock**
```typescript
const publishCompletedRef = useRef(false);

// In PUBLISH_COMPLETE handler:
if (!publishCompletedRef.current) {
  publishCompletedRef.current = true; // PERMANENTLY lock
  // ... process completion
}

// In all other event handlers, gate on the lock:
if (payload.pluginId === pluginId && !publishCompletedRef.current) { ... }
```

**Pattern 3: Log deduplication**
```typescript
const seenLogKeysRef = useRef(new Set<string>());

const addLog = useCallback((log: PublishLogEntry) => {
  const key = `${log.timestamp}|${log.step}|${log.message}`;

  if (seenLogKeysRef.current.has(key)) return; // Skip duplicates
  seenLogKeysRef.current.add(key);
  setLogs(prev => [...prev, { ...log }]);
}, []);
```

**Pattern 4: Force-unsubscribe on close + cleanup**
```typescript
const unsubsRef = useRef<Array<() => void>>([]);

const forceUnsubAll = useCallback(() => {
  unsubsRef.current.forEach(fn => fn());
  unsubsRef.current = [];
}, []);

// In effect: store unsub handles
const unsub1 = wsClient.on(WS_EVENTS.PUBLISH_STARTED, handler);
unsubsRef.current.push(unsub1);

// On dialog close:
forceUnsubAll();
```

**Pattern 5: Reset all state on dialog open**
```typescript
useEffect(() => {
  if (open) {
    publishCompletedRef.current = false;
    seenLogKeysRef.current.clear();
    setLogs([]);
    // ... reset all other state
  } else {
    forceUnsubAll(); // Clean up when dialog closes
  }

}, [open, pluginId, siteId]);
```

### Key Principles

1. **Never put callbacks or computed strings in `useEffect` dependency arrays** — use `useRef` to read their latest value inside the effect
2. **WebSocket listeners must have a completion lock** — once a lifecycle event (publish, restore, backup) is "done," ignore all subsequent events for that session
3. **Deduplicate event-driven state updates** — WebSocket events can arrive out-of-order or multiple times; always deduplicate by a composite key
4. **Force-unsubscribe on unmount AND on lifecycle completion** — don't rely on cleanup alone

---

## 6. Toast Notification Deduplication

### Symptom

Multiple identical toast notifications appeared simultaneously, e.g.:
- "Publish complete" × 3 (from WebSocket event + local state change + query refetch)
- "Connection test passed" × 2 (from WebSocket + API response)

### Root Cause

Both the WebSocket event handler (`useWsToastNotifications`) AND the local UI handler would fire toasts for the same event. No dedup existed.

### Fix — `src/lib/dedupToast.ts`

```typescript
const DEDUP_WINDOW_MS = 3000; // 3 seconds
const recentToasts = new Map<string, number>();

function isDuplicate(key: string): boolean {
  const lastShown = recentToasts.get(key);

  if (lastShown && Date.now() - lastShown < DEDUP_WINDOW_MS) return true;
  recentToasts.set(key, Date.now());
  // Periodic cleanup
  if (recentToasts.size > 50) {
    const now = Date.now();

    for (const [k, t] of recentToasts) {
      if (now - t > DEDUP_WINDOW_MS) recentToasts.delete(k);
    }
  }

  return false;
}

export const dedupToast = {
  success: createDedupMethod("success"),
  error: createDedupMethod("error"),
  warning: createDedupMethod("warning"),
  info: createDedupMethod("info"),
  message: createDedupMethod("message"),
};
```

### Usage

```typescript
// Replace: import { toast } from "sonner";
// With:    import { dedupToast as toast } from "@/lib/dedupToast";
```

### Key Principle

**Always use `dedupToast` instead of raw `sonner.toast`** in any component that might fire toasts from multiple sources (WebSocket + API response + state change).

---

## 7. Circuit Breaker for Failing Endpoints

### Symptom

When a backend endpoint was consistently failing (e.g., site offline), every poll/check would generate error toasts and error log entries indefinitely.

### Root Cause

No mechanism existed to stop calling a persistently failing endpoint. Polling intervals would keep firing, generating noise.

### Fix — `src/lib/circuitBreaker.ts`

```typescript
const DEFAULT_CONFIG = {
  failureThreshold: 5,    // Open circuit after 5 failures
  cooldownMs: 60000,      // Wait 1 minute before trying again
  failureWindowMs: 60000, // Count failures within 1-minute window
};

// Usage:
const result = await withCircuitBreaker('api.getSites', () => api.getSites());
```

States: `closed` (normal) → `open` (blocked) → `half-open` (test one request after cooldown) → `closed` (if test succeeds).

### Key Principle

**Wrap polling and health-check calls in a circuit breaker** to prevent error storms against failing endpoints.

---

## 8. Snapshot Query Retry Suppression

### Symptom

Snapshot list queries for disconnected sites (siteId=0) would trigger the global error modal on every poll cycle.

### Root Cause

Queries with `siteId: 0` are "background" queries that run speculatively. Their failures should be silent.

### Fix

```typescript
// src/hooks/useRemoteSnapshots.ts
{
  retry: false,
  refetchOnWindowFocus: false,
  meta: { suppressGlobalError: true }, // ← Prevents error modal
}

// src/components/settings/SnapshotSettingsTab.tsx
{
  retry: false,
  meta: { suppressGlobalError: true },
}
```

The `suppressGlobalError` flag is checked in the global `QueryCache.onError` handler and skips the error modal.

---

## 9. Anti-Patterns Summary (Never Do This)

### ❌ NEVER: Use React Query defaults for retry/refetch

```typescript
// BAD — causes retry storms and background refetches
const queryClient = new QueryClient(); // uses retry:3, refetchOnWindowFocus:true
```
```typescript
// GOOD
const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false, refetchOnWindowFocus: false } },
});
```

### ❌ NEVER: Put callbacks/labels in useEffect deps for WebSocket listeners

```typescript
// BAD — re-subscribes on every render
useEffect(() => {
  const unsub = wsClient.on("event", () => onComplete(pluginName));
  return () => unsub();
}, [onComplete, pluginName]); // ← these change every render!
```
```typescript
// GOOD — stable deps, read latest from ref
const onCompleteRef = useRef(onComplete);
onCompleteRef.current = onComplete;
useEffect(() => {
  const unsub = wsClient.on("event", () => onCompleteRef.current("..."));

  return () => unsub();
}, []); // ← stable, only subscribes once
```

### ❌ NEVER: Rely only on UI-level button disabling for mutation dedup

```typescript
// BAD — can be bypassed by programmatic callers
<Button disabled={isLoading} onClick={publish}>Publish</Button>
```
```typescript
// GOOD — dedup at the API method level
publishPlugin: (() => {
  const inFlight = new Set<string>();

  return (pluginId, siteId, opts) => {
    const key = `${pluginId}:${siteId}`;

    if (inFlight.has(key)) return Promise.resolve({ success: false, error: { code: "E_DEDUP" } });
    inFlight.add(key);

    return request(...).finally(() => inFlight.delete(key));
  };
})(),
```

### ❌ NEVER: Use raw `toast()` from sonner in WebSocket-connected components

```typescript
// BAD — duplicate toasts from WS + local state
import { toast } from "sonner";
toast.success("Done");
```
```typescript
// GOOD — automatic 3-second dedup window
import { dedupToast as toast } from "@/lib/dedupToast";
toast.success("Done");
```

### ❌ NEVER: Process WebSocket events without a completion lock

```typescript
// BAD — late events from previous sessions leak through
wsClient.on("complete", (data) => { setIsComplete(true); });
```
```typescript
// GOOD — ignore events after lifecycle is done
const completedRef = useRef(false);
wsClient.on("complete", (data) => {
  if (completedRef.current) return;
  completedRef.current = true;
  setIsComplete(true);
});
```

### ❌ NEVER: Fire polling queries for background/speculative targets without suppressGlobalError

```typescript
// BAD — error modal pops up for disconnected sites
useQuery({ queryKey: ["snapshots", 0], queryFn: () => api.getSnapshots(0) });
```
```typescript
// GOOD — silent failure for background queries
useQuery({
  queryKey: ["snapshots", 0],
  queryFn: () => api.getSnapshots(0),
  retry: false,
  meta: { suppressGlobalError: true },
});
```

---

## Files Involved

| File | What It Does |
|------|-------------|
| `src/App.tsx` | Global QueryClient config (`retry: false`, `refetchOnWindowFocus: false`) |
| `src/lib/api/methods.ts` | Publish dedup lock (IIFE with `inFlight` Set + cooldown Map) |
| `src/lib/dedupToast.ts` | Toast deduplication wrapper (3s window) |
| `src/lib/circuitBreaker.ts` | Circuit breaker pattern for failing endpoints |
| `src/components/plugins/PublishProgressDialog.tsx` | WS listener stabilization (useRef, completion lock, log dedup) |
| `src/hooks/useWsToastNotifications.ts` | Uses `dedupToast` for WS-driven notifications |
| `src/hooks/useRemoteSnapshots.ts` | Per-query retry suppression + `suppressGlobalError` |
| `src/components/settings/SnapshotSettingsTab.tsx` | Per-query retry suppression + `suppressGlobalError` |
