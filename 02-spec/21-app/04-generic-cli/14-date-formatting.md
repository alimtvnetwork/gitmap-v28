# Date Display Format

> **Related specs:**
> - [06-output-formatting.md](06-output-formatting.md) — output formats that display dates
> - [15-constants-reference.md](15-constants-reference.md) — date format layout constants

## Principle

All date/time output passes through a single centralized formatting
function. No command formats dates inline.

## Pipeline

1. **Normalize to UTC** — convert input time to UTC.
2. **Convert to local timezone** — use the machine's local timezone.
3. **Format** — render as `DD-Mon-YYYY hh:mm AM/PM`.

## Layout

```
02-Jan-2006 03:04 PM
```

| Component | Width | Example |
|-----------|-------|---------|
| Day | 2 digits | `06` |
| Month | 3-letter | `Mar` |
| Year | 4 digits | `2026` |
| Time | 12-hour | `03:17 AM` |

## Implementation

```go
const DateDisplayLayout = "02-Jan-2006 03:04 PM"

func FormatDisplayDate(t time.Time) string {
    utc := t.UTC()
    local := utc.Local()

    return local.Format(constants.DateDisplayLayout)
}
```

## Rules

- No `time.Format` calls in command handlers.
- Layout constant lives in `constants`.
- UTC → Local conversion inside the function, not at call site.
- A `FormatDisplayDateUTC` variant adds a `(UTC)` suffix.

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
