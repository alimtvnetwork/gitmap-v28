# Performance & Optimization

Universal guidelines for caching strategies, lazy loading, resource management, and efficient execution patterns across all languages.

---

## 1. Avoid N+1 Patterns

Never perform repeated expensive calls (DB queries, subprocess spawns, network requests) inside loops. Pre-compute or batch instead.

### Bad

```go
for _, repo := range repos {
    author, _ := git.GetAuthor(repo.Path) // subprocess per repo
    fmt.Println(repo.Name, author)
}
```

### Good

```go
authors := git.GetAuthorsForPaths(repoPaths) // single batch call
for _, repo := range repos {
    fmt.Println(repo.Name, authors[repo.Path])
}
```

```typescript
// Bad: N queries
const users = await Promise.all(ids.map(id => db.getUser(id)));

// Good: single query
const users = await db.getUsersByIds(ids);
```

---

## 2. Caching Strategies

### In-Memory Cache

Use for data that is expensive to compute but rarely changes within a session:

```go
var cache = map[string]string{}

func getExpensive(key string) string {
    if val, ok := cache[key]; ok {
        return val
    }
    val := computeExpensive(key)
    cache[key] = val
    return val
}
```

### Cache Invalidation Rules

| Rule | Guidance |
|---|---|
| Time-based expiry | Set a TTL; re-fetch after expiry |
| Event-based invalidation | Clear on write/mutation |
| Never cache errors | Failed lookups must not be stored |
| Bound cache size | Use LRU or fixed-size maps to prevent memory leaks |

### File-Based Cache

For CLI tools, cache results to disk when computation is slow and data is stable:

```go
const CacheTTL = 24 * time.Hour

func loadCachedOrCompute(path string, compute func() []byte) []byte {
    info, err := os.Stat(path)
    if err == nil && time.Since(info.ModTime()) < CacheTTL {
        data, _ := os.ReadFile(path)
        return data
    }
    data := compute()
    os.WriteFile(path, data, 0o644)
    return data
}
```

---

## 3. Lazy Loading

Defer expensive initialization until the resource is actually needed.

### Go — Lazy Initialization

```go
var dbOnce sync.Once
var dbConn *sql.DB

func getDB() *sql.DB {
    dbOnce.Do(func() {
        dbConn, _ = sql.Open("sqlite3", dbPath)
    })
    return dbConn
}
```

### TypeScript — Code Splitting

```typescript
// Lazy-load heavy components
const HeavyChart = lazy(() => import("./components/HeavyChart"));

// Lazy-load modules
async function runExport() {
    const { generateCSV } = await import("./utils/export");
    return generateCSV(data);
}
```

### Principles

- **Don't load what you don't need**: Parse config files only when accessed
- **Defer filesystem scans**: Walk directories only when results are requested
- **Lazy-open databases**: Connect on first query, not at startup

---

## 4. Resource Management

### Close What You Open

Every opened resource (file, DB connection, HTTP body) must be closed:

```go
f, err := os.Open(path)
if err != nil {
    return err
}
defer f.Close()
```

```typescript
const controller = new AbortController();
try {
    const res = await fetch(url, { signal: controller.signal });
    return await res.json();
} finally {
    controller.abort();
}
```

### Pooling

Reuse expensive resources instead of creating them per-call:

- Database connection pools (bounded size)
- HTTP clients (reuse `http.Client` in Go; share `fetch` sessions)
- Buffer pools for serialization (`sync.Pool` in Go)

---

## 5. Batch Operations

Group small operations into batches to reduce overhead:

```go
// Bad: one insert per record
for _, rec := range records {
    db.Exec(sqlInsert, rec.Name, rec.Value)
}

// Good: single transaction
tx, _ := db.Begin()
for _, rec := range records {
    tx.Exec(sqlInsert, rec.Name, rec.Value)
}
tx.Commit()
```

```typescript
// Bad: sequential awaits
for (const item of items) {
    await saveItem(item);
}

// Good: parallel batch
await Promise.all(items.map(item => saveItem(item)));
```

---

## 6. Progress & Timeout Patterns

Long-running operations must provide feedback and respect timeouts:

```go
// Progress counter for batch jobs
for i, item := range items {
    fmt.Fprintf(os.Stderr, "\r[%d/%d] Processing %s", i+1, len(items), item.Name)
    process(item)
}
```

```go
// Context timeout for network calls
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()
resp, err := client.Do(req.WithContext(ctx))
```

---

## 7. DOM & Rendering Performance (Frontend)

- **Minimize re-renders**: Use `React.memo`, `useMemo`, and `useCallback` for expensive computations
- **Virtualize long lists**: Use windowing libraries for lists exceeding 50 items
- **Debounce input handlers**: Delay search/filter callbacks by 200–300ms
- **Lazy-load images**: Use `loading="lazy"` on `<img>` tags below the fold
- **Reduce DOM depth**: Keep component nesting shallow; avoid wrapper-heavy layouts

---

## 8. Measurement Rules

- **Profile before optimizing**: Never guess at bottlenecks — measure first
- **Benchmark critical paths**: Track timing for operations users wait on
- **Set performance budgets**: Define maximum acceptable latency for key operations
- **Log slow operations**: Emit warnings when operations exceed thresholds

```go
start := time.Now()
result := expensiveOperation()
elapsed := time.Since(start)
if elapsed > 500*time.Millisecond {
    log.Printf("[perf] expensiveOperation took %s", elapsed)
}
```

---

## References

- [Code Quality Improvement](./01-code-quality-improvement.md)
- [Logging & Observability](./07-logging-observability.md)
- [Security & Secrets](./08-security-secrets.md)
- [Subprocess Optimization](../01-app/03-scanner.md)

---

**Contributors**: Alim Ul Karim · [Riseup Labs](https://riseuplabs.com)
