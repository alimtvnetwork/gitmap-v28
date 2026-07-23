# 17 — Resilience Patterns

Retry logic, circuit breakers, timeouts, graceful degradation, and backpressure.

## Retry Logic

Retry only transient failures (network, HTTP 429/5xx). Never retry permanent errors (400/401/403/404). Exponential backoff with jitter. Max 3 attempts. Log every retry.

## Circuit Breakers

States: Closed → Open → Half-Open. Transition thresholds must be constants. Log every state transition. Defaults: 5 failures to open, 30s reset timeout.

## Timeouts

Every external call must have a context timeout. Defaults: HTTP 10s, TCP dial 5s, DB query 15s, shutdown 30s. Timeouts must be constants. Propagate `ctx.Done()`.

## Graceful Degradation

When a dependency fails, reduce functionality — never crash. Serve cached or default data. Communicate degraded state. Auto-recover when dependency returns.

| Strategy | When |
|----------|------|
| Cached fallback | Read-path failure |
| Default response | Optional enrichment unavailable |
| Feature toggle | Entire subsystem down |
| Partial result | Some items in batch fail |

## Backpressure

Bound all queues and buffers. Reject or slow new work at capacity — never silently drop. Signal callers (HTTP 429 / error).

## Constraints

No magic numbers — all thresholds in constants. Every external call has a timeout. Log all state transitions.

---

Source: `spec/05-coding-guidelines/18-resilience-patterns.md`
