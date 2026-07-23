# 15 — Concurrency Patterns

Goroutine management, mutex usage, channel patterns, and race prevention.

## Goroutine Rules

Every goroutine has an owner that ensures it exits. Use `context.Context` for cancellation. Never launch in `init()`. Bound concurrency with semaphores.

Use `errgroup` for structured patterns.

## Bounded Concurrency

```go
sem := make(chan struct{}, maxWorkers)
```

`maxWorkers` is a constant, never hardcoded inline. Default: CPU count or fixed cap (10).

## Mutex Rules

Use `sync.RWMutex` for read-heavy access. Always `defer Unlock()`. Keep critical sections minimal — no I/O. Never hold a mutex across channel operations. Mutex is a struct field, not global.

## Channel Rules

Always close from sender side. Use buffered for known batch sizes. Select with `ctx.Done()` in every loop. Declare direction in signatures.

## Race Prevention

Run `go test -race ./...` in CI — non-negotiable, no suppressions. Test with `-count=5`.

Common fixes: `v := v` for loop capture, `sync.RWMutex` for shared maps, channels for slice collection.

## Testing Concurrent Code

Use `t.Parallel()`. Never `time.Sleep` — use synchronization. Detect leaks with `goleak`.

---

Source: `spec/05-coding-guidelines/16-concurrency-patterns.md`
