# 09 — Performance & Optimization

Caching, lazy loading, resource management, and efficient execution.

## Avoid N+1 Patterns

Never perform repeated expensive calls inside loops. Pre-compute or batch instead.

## Caching Strategies

In-memory cache for expensive, rarely-changing data. File-based cache for CLI tools with slow computation. Rules: set TTL, invalidate on write, never cache errors, bound size.

## Lazy Loading

Defer expensive initialization until needed. `sync.Once` for Go. `React.lazy` for TS. Don't load what you don't need.

## Resource Management

Close every opened resource (file, DB, HTTP body). Use `defer` immediately after acquisition. Pool expensive resources.

## Batch Operations

Group small operations into batches. Use transactions for multi-statement writes.

## Progress & Timeouts

Long operations must provide `[current/total]` feedback. Every external call must have a `context.WithTimeout`.

## DOM Performance (Frontend)

Minimize re-renders (`React.memo`, `useMemo`). Virtualize long lists (50+ items). Debounce input handlers (200–300ms). Lazy-load images.

## Measurement

Profile before optimizing. Benchmark critical paths. Log slow operations exceeding thresholds.

---

Source: `spec/05-coding-guidelines/09-performance-optimization.md`
