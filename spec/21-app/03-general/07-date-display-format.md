# Date Display Format Pattern

## Principle

All CLI tools must display dates in a consistent, human-friendly
format. No command should format dates inline — every date passes
through a centralized formatting function.

## Pipeline

1. **Normalize to UTC** — convert input time to UTC.
2. **Convert to local timezone** — use the machine's local timezone.
3. **Format** — render as `DD-Mon-YYYY hh:mm AM/PM`.

## Layout

```
02-Jan-2006 03:04 PM
```

| Component | Width    | Example      |
|-----------|----------|--------------|
| Day       | 2 digits | `06`         |
| Month     | 3-letter | `Mar`        |
| Year      | 4 digits | `2026`       |
| Time      | 12-hour  | `03:17 AM`   |

## Constants

All format strings live in the constants package:

```go
const (
    DateDisplayLayout = "02-Jan-2006 03:04 PM"
    DateUTCSuffix     = " (UTC)"
)
```

## Formatting Function

A single function handles the full pipeline:

```go
func FormatDisplayDate(t time.Time) string {
    utc := t.UTC()
    local := utc.Local()

    return local.Format(constants.DateDisplayLayout)
}
```

## Rules

- No `time.Format` calls in command handlers — always delegate to
  the centralized function.
- The format layout constant lives in `constants`, not in the
  formatting function.
- UTC → Local conversion happens inside the function, not at the
  call site.
- A `FormatDisplayDateUTC` variant exists for UTC-only display with
  a `(UTC)` suffix.
