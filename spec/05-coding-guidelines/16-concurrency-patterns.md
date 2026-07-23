# Concurrency Patterns

## Overview

Standards for goroutine management, mutex usage, channel patterns, and
race condition prevention. Concurrency must be explicit, bounded, and
testable — never fire-and-forget.

---

## 1. Goroutine Management

### Lifecycle Rules

| Rule | Rationale |
|------|-----------|
| Every goroutine has an owner | The caller that spawns it must ensure it exits |
| Use `context.Context` for cancellation | Propagates deadlines and cancellation cleanly |
| Never launch goroutines in `init()` | Unpredictable startup ordering |
| Bounded concurrency with semaphores | Prevent resource exhaustion |

### Spawning Pattern

Always use `errgroup` or a similar structured pattern:

```go
g, ctx := errgroup.WithContext(ctx)

for _, repo := range repos {
    repo := repo // capture loop variable
    g.Go(func() error {
        return processRepo(ctx, repo)
    })
}

if err := g.Wait(); err != nil {
    return fmt.Errorf("processing repos: %w", err)
}
```

### Bounded Concurrency

Limit parallel work with a semaphore channel:

```go
sem := make(chan struct{}, maxWorkers)

for _, item := range items {
    sem <- struct{}{} // acquire
    go func(it Item) {
        defer func() { <-sem }() // release
        process(it)
    }(item)
}
```

### Rules

- Set `maxWorkers` as a constant, never hardcoded inline.
- Default concurrency limit: number of CPUs or a fixed cap (e.g., 10).
- Log goroutine count at debug level for diagnostics.

---

## 2. Mutex Usage

### When to Use Mutexes

| Use Mutex | Use Channel |
|-----------|-------------|
| Protecting shared state (maps, counters) | Coordinating between goroutines |
| Short critical sections (<1µs) | Passing ownership of data |
| Read-heavy workloads (`sync.RWMutex`) | Fan-out/fan-in pipelines |

### Mutex Rules

```go
type RepoCache struct {
    mu    sync.RWMutex
    repos map[string]*Repo
}

func (c *RepoCache) Get(name string) (*Repo, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    r, ok := c.repos[name]
    return r, ok
}

func (c *RepoCache) Set(name string, repo *Repo) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.repos[name] = repo
}
```

| Rule | Rationale |
|------|-----------|
| Always `defer Unlock()` | Prevents deadlocks on early returns or panics |
| Keep critical sections minimal | Hold locks for data access only, not I/O |
| Never hold a mutex across a channel operation | Deadlock risk |
| Use `sync.RWMutex` for read-heavy access | Allows concurrent readers |
| Mutex is a struct field, not a global | Encapsulation and testability |

### Lock Ordering

When multiple mutexes are needed, define and document a strict
acquisition order to prevent deadlocks:

```go
// Lock order: RepoCache.mu → ReleaseCache.mu
// Never acquire ReleaseCache.mu while holding RepoCache.mu
```

---

## 3. Channel Patterns

### Directional Channels

Always declare channel direction in function signatures:

```go
func producer(out chan<- Result) { ... }
func consumer(in <-chan Result)  { ... }
```

### Fan-Out / Fan-In

```go
func fanOut(ctx context.Context, items []Item, workers int) <-chan Result {
    results := make(chan Result, len(items))
    sem := make(chan struct{}, workers)

    var wg sync.WaitGroup
    for _, item := range items {
        wg.Add(1)
        sem <- struct{}{}
        go func(it Item) {
            defer wg.Done()
            defer func() { <-sem }()
            results <- process(ctx, it)
        }(item)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

### Channel Rules

| Rule | Rationale |
|------|-----------|
| Always close channels from the sender side | Prevents send-on-closed panics |
| Use buffered channels for known batch sizes | Prevents goroutine leaks on early exit |
| Select with `ctx.Done()` in every loop | Enables cancellation |
| Never rely on channel buffer size for correctness | Buffers are for performance, not logic |

### Select with Cancellation

```go
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case item, ok := <-in:
        if !ok {
            return nil
        }
        handle(item)
    }
}
```

---

## 4. Race Condition Prevention

### Detection

- Run `go test -race ./...` in CI on every commit.
- The race detector is non-negotiable — no suppressions allowed.
- Test with `-count=5` minimum to surface intermittent races.

### Common Race Patterns

| Pattern | Problem | Fix |
|---------|---------|-----|
| Shared map writes | Concurrent map access panics | `sync.RWMutex` or `sync.Map` |
| Loop variable capture | Goroutine reads stale value | `v := v` or function parameter |
| Unprotected struct fields | Data races on read/write | Mutex or atomic operations |
| Slice append in goroutines | Non-atomic append | Collect via channels |

### Atomic Operations

Use `sync/atomic` for simple counters — avoid mutex overhead:

```go
var processed atomic.Int64

g.Go(func() error {
    // ... work ...
    processed.Add(1)
    return nil
})

verbose.Log("progress", "processed %d items", processed.Load())
```

### sync.Once for Initialization

```go
var (
    dbOnce sync.Once
    dbConn *sql.DB
)

func getDB() *sql.DB {
    dbOnce.Do(func() {
        dbConn = openConnection()
    })
    return dbConn
}
```

---

## 5. Testing Concurrent Code

### Rules

- Use `t.Parallel()` for independent test cases.
- Test with `-race` flag always enabled.
- Use `time.After` or context deadlines — never `time.Sleep`.
- Verify goroutine cleanup with `goleak` in tests.

### Goroutine Leak Detection

```go
func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}
```

### Deterministic Testing

Replace concurrency with sequential execution in unit tests
when testing business logic, not concurrency behavior:

```go
// Test the processing logic, not the concurrency
func TestProcessRepo(t *testing.T) {
    result := processRepo(context.Background(), testRepo)
    assert.NoError(t, result)
}
```

---

## Constraints

- Every goroutine has a clear owner and shutdown path.
- Concurrency limits are constants, not magic numbers.
- `go test -race` runs on every CI build.
- Mutexes use `defer Unlock()` — no exceptions.
- Channels declare direction in function signatures.
- No `time.Sleep` in tests — use synchronization primitives.
- Critical sections contain no I/O or channel operations.

---

## References

- Error Handling: `spec/05-coding-guidelines/04-error-handling.md`
- Testing Patterns: `spec/05-coding-guidelines/06-testing-patterns.md`
- Performance: `spec/05-coding-guidelines/09-performance-optimization.md`
- Logging: `spec/05-coding-guidelines/07-logging-observability.md`
