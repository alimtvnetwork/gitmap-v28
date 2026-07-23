# Resilience Patterns

## Purpose

Standardise fault-tolerance techniques so every network call, external
dependency, and resource-intensive operation degrades gracefully under
failure instead of cascading errors to the user.

## Retry Logic

### Rules

1. Retry only on **transient** failures (network timeout, HTTP 429/5xx).
2. Never retry **permanent** errors (HTTP 400/401/403/404, validation).
3. Use **exponential backoff** with jitter to avoid thundering herd.
4. Cap retries with a **maximum attempt count** (default 3).
5. Log every retry with attempt number, wait duration, and error.

### Go Pattern

```go
func withRetry(ctx context.Context, maxAttempts int, base time.Duration, op func() error) error {
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        err := op()
        if err == nil {
            return nil
        }
        if !isTransient(err) {
            return fmt.Errorf("permanent failure: %w", err)
        }
        if attempt == maxAttempts {
            return fmt.Errorf("failed after %d attempts: %w", maxAttempts, err)
        }
        wait := base * time.Duration(1<<(attempt-1))
        jitter := time.Duration(rand.Int63n(int64(wait / 4)))
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(wait + jitter):
        }
    }
    return nil
}
```

### Constants

| Constant             | Value | Description                   |
|----------------------|-------|-------------------------------|
| `RetryMaxAttempts`   | 3     | Maximum retry attempts        |
| `RetryBaseDelayMs`   | 1000  | Initial backoff in ms         |
| `RetryBackoffFactor` | 2     | Exponential multiplier        |
| `RetryMaxDelayMs`    | 30000 | Backoff cap                   |

## Circuit Breakers

### States

```
CLOSED  ──(failure threshold)──▶  OPEN
  ▲                                  │
  │                            (half-open timer)
  │                                  ▼
  └────(success threshold)───  HALF-OPEN
```

### Rules

1. **Closed** — requests pass through; failures are counted.
2. **Open** — requests fail immediately; no calls to downstream.
3. **Half-Open** — limited probe requests test recovery.
4. Transition thresholds must be constants, not magic numbers.
5. Log every state transition with component name and reason.

### Go Pattern

```go
type CircuitState int

const (
    StateClosed   CircuitState = iota
    StateOpen
    StateHalfOpen
)

type Breaker struct {
    state        CircuitState
    failures     int
    threshold    int
    halfOpenMax  int
    resetTimeout time.Duration
    lastFailure  time.Time
    mu           sync.Mutex
}

func (b *Breaker) Call(op func() error) error {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.state == StateOpen {
        if time.Since(b.lastFailure) < b.resetTimeout {
            return ErrCircuitOpen
        }
        b.state = StateHalfOpen
    }
    b.mu.Unlock()

    err := op()

    b.mu.Lock()
    if err != nil {
        b.failures++
        b.lastFailure = time.Now()
        if b.failures >= b.threshold {
            b.state = StateOpen
        }
        return err
    }
    b.failures = 0
    b.state = StateClosed
    return nil
}
```

### Constants

| Constant                | Value | Description                    |
|-------------------------|-------|--------------------------------|
| `BreakerFailThreshold`  | 5     | Failures before opening        |
| `BreakerResetTimeoutMs` | 30000 | Time before half-open probe    |
| `BreakerHalfOpenMax`    | 2     | Probe requests in half-open    |

## Timeouts

### Rules

1. Every external call must have a **context timeout**.
2. Default timeout: **10 seconds** for HTTP, **5 seconds** for TCP dial.
3. Timeouts must be **constants**, never inline literals.
4. Propagate context cancellation — never ignore `ctx.Done()`.
5. Log timeout events with operation name and configured duration.

### Go Pattern

```go
ctx, cancel := context.WithTimeout(parentCtx, constants.HTTPTimeoutDuration)
defer cancel()

resp, err := client.Do(req.WithContext(ctx))
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        return fmt.Errorf("request timed out after %s: %w", constants.HTTPTimeoutDuration, err)
    }
    return fmt.Errorf("request failed: %w", err)
}
```

### Constants

| Constant              | Value | Description                |
|-----------------------|-------|----------------------------|
| `HTTPTimeoutSec`      | 10    | HTTP request timeout       |
| `TCPDialTimeoutSec`   | 5     | TCP connection timeout     |
| `DBQueryTimeoutSec`   | 15    | Database query timeout     |
| `ShutdownTimeoutSec`  | 30    | Graceful shutdown window   |

## Graceful Degradation

### Rules

1. When a dependency fails, **reduce functionality** — never crash.
2. Return **cached or default data** when the primary source is down.
3. Clearly communicate degraded state to the user.
4. Log degradation events with severity `WARN`.
5. Automatically recover when the dependency comes back.

### Strategies

| Strategy           | When to Use                        | Example                              |
|--------------------|------------------------------------|--------------------------------------|
| Cached fallback    | Read-path failure                  | Serve stale scan results from DB     |
| Default response   | Optional enrichment unavailable    | Skip GitHub avatar, show placeholder |
| Feature toggle     | Entire subsystem down              | Disable release uploads, allow scan  |
| Partial result     | Some repos in batch fail           | Return successes, list failures      |

### Go Pattern

```go
func getRepoData(ctx context.Context, name string) (RepoData, error) {
    data, err := fetchFromAPI(ctx, name)
    if err == nil {
        cacheStore(name, data)
        return data, nil
    }

    cached, cacheErr := cacheLoad(name)
    if cacheErr == nil {
        log.Warn("serving cached data for %s: %v", name, err)
        return cached, nil
    }

    return RepoData{Name: name}, fmt.Errorf("degraded: %w", err)
}
```

## Backpressure Handling

### Rules

1. **Bound all queues and buffers** — unbounded queues cause OOM.
2. When at capacity, **reject or slow** new work — never silently drop.
3. Use buffered channels with explicit capacity constants.
4. Signal callers when backpressure is applied (HTTP 429 / error).
5. Monitor queue depth and alert at 80% capacity.

### Go Pattern

```go
type WorkQueue struct {
    ch chan WorkItem
}

func NewWorkQueue(capacity int) *WorkQueue {
    return &WorkQueue{ch: make(chan WorkItem, capacity)}
}

func (q *WorkQueue) Submit(ctx context.Context, item WorkItem) error {
    select {
    case q.ch <- item:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    default:
        return ErrQueueFull
    }
}
```

### Constants

| Constant             | Value | Description                    |
|----------------------|-------|--------------------------------|
| `QueueCapacity`      | 100   | Maximum buffered work items    |
| `QueueAlertPercent`  | 80    | Alert threshold percentage     |
| `RateLimitPerSec`    | 10    | Max requests per second        |

## Constraints

- Files ≤ 200 lines, functions 8–15 lines.
- No magic numbers — all thresholds in constants.
- Every external call must have a timeout.
- Retry only transient errors — classify before retrying.
- Log all state transitions (circuit breaker, degradation, backpressure).

## References

- [04 Error Handling](../05-coding-guidelines/04-error-handling.md)
- [07 Logging & Observability](../05-coding-guidelines/07-logging-observability.md)
- [09 Performance & Optimization](../05-coding-guidelines/09-performance-optimization.md)
- [15 Monitoring & Alerting](../05-coding-guidelines/15-monitoring-alerting.md)
- [16 Concurrency Patterns](../05-coding-guidelines/16-concurrency-patterns.md)
