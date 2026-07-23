# 04 — Error Handling

Consistent error handling for predictable behavior and clear diagnostics.

## Fail Fast

Check preconditions at the top. Return or throw early rather than nesting the happy path.

## Go — Check Errors Immediately

Handle the error on the very next line after a call. Never use the result before checking the error.

## Go — Wrap Errors with Context

Use `fmt.Errorf` with `%w` to add context while preserving the original error for `errors.Is` / `errors.As`.

```go
result, err := db.Query(query)
if err != nil {
    return fmt.Errorf("querying repos: %w", err)
}
```

## Never Swallow Errors (Code Red)

Every error must be explicitly logged to `os.Stderr`. Empty catch blocks are prohibited. Generic "file not found" messages without paths are prohibited.

## Mandatory Path Context in File Errors

Every error involving a file/directory must include:
1. The exact resolved path
2. The operation attempted
3. The specific failure reason

```
Error: [message] at [path]: [error] (operation: [op], reason: [reason])
```

## Format Verb Compliance

All `fmt.Fprintf`, `fmt.Printf`, and `fmt.Errorf` calls must match format verb count to argument count. Enforce via `go vet` in CI.

## Go — No `panic` for Expected Conditions

Reserve `panic` for truly unrecoverable programmer errors. All expected failures use `error` returns.

---

Source: `spec/05-coding-guidelines/04-error-handling.md`
