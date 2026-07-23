# Offline Mode Detection

## Purpose

Detect network unavailability before operations that require internet
access (clone, pull, push, release, update) and fail gracefully with
a clear message instead of cryptic git errors.

## Behavior

1. Before network-dependent operations, call `CheckOnline()`.
2. The check attempts a TCP dial to `github.com:443` with a 5-second timeout.
3. If offline, print a user-friendly message and exit early.
4. Commands with `--force-offline` or local-only operations skip the check.

## Detection Method

```go
conn, err := net.DialTimeout("tcp", "github.com:443", 5*time.Second)
```

- Success → online, proceed normally.
- Failure → offline, print warning and return error.

## Affected Commands

| Command | Behavior When Offline |
|---------|----------------------|
| `clone` | Skip with offline warning |
| `pull` | Skip with offline warning |
| `release` | Abort before push step |
| `update` | Abort with offline warning |
| `push` (internal) | Abort with offline warning |

## Output

```
  ⚠ Network unavailable — cannot reach github.com.
  Offline operations (scan, list, status, group) still work.
```

## Implementation Files

| File | Action | Purpose |
|------|--------|---------|
| `gitutil/network.go` | CREATE | `CheckOnline()` function |
| `constants/constants_messages.go` | MODIFY | Add offline messages |

## Constraints

- Files ≤ 200 lines, functions 8–15 lines.
- No magic strings — all in constants.
- Timeout must be a constant (5 seconds).
