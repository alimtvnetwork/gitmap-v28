# Upload Retry Logic

## Purpose

Add exponential backoff retry logic for GitHub release asset uploads
to handle transient network failures and API rate limits gracefully.

## Behavior

1. Each asset upload is attempted up to **3 times** (configurable).
2. On failure, wait with **exponential backoff**: 1s, 2s, 4s.
3. HTTP 429 (rate limit) and 5xx errors trigger retries.
4. HTTP 4xx errors (except 429) fail immediately (no retry).
5. Network errors (timeout, DNS) always trigger retries.
6. Each retry logs a warning with attempt number and wait duration.

## Constants

| Constant              | Value | Description                    |
|-----------------------|-------|--------------------------------|
| `RetryMaxAttempts`    | 3     | Maximum upload attempts        |
| `RetryBaseDelayMs`    | 1000  | Initial delay in milliseconds  |
| `RetryBackoffFactor`  | 2     | Multiplier per retry           |

## Output Messages

```
  ⟳ Retry 1/3 for asset.zip (waiting 1s)...
  ⟳ Retry 2/3 for asset.zip (waiting 2s)...
  ✗ Upload failed for asset.zip after 3 attempts: <error>
  ✓ Uploaded asset.zip (attempt 2)
```

## Implementation Files

| File                             | Action | Purpose                        |
|----------------------------------|--------|--------------------------------|
| `release/retry.go`              | CREATE | Generic retry wrapper function |
| `release/assetsupload.go`       | MODIFY | Use retry wrapper              |
| `constants/constants_assets.go` | MODIFY | Add retry constants/messages   |

## Constraints

- Files ≤ 200 lines, functions 8–15 lines.
- No magic strings — all in constants.
- Positive conditionals only.
