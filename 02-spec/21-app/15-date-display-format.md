# Date Display Format

## Overview

All CLI date output passes through a centralized formatting function
before display. No command formats dates inline — every date goes
through `gitutil.FormatDisplayDate`.

## Pipeline

1. **Normalize to UTC** — convert the input `time.Time` to UTC
   (`t.UTC()`).
2. **Convert to local timezone** — convert UTC to the machine's
   local timezone (`utc.Local()`).
3. **Format for display** — render using the layout
   `DD-Mon-YYYY hh:mm AM/PM`.

## Display Format

```
02-Jan-2006 03:04 PM
```

| Component | Format   | Example      |
|-----------|----------|--------------|
| Day       | 2 digits | `06`         |
| Month     | 3-letter | `Mar`        |
| Year      | 4 digits | `2026`       |
| Time      | 12-hour  | `03:17 AM`   |

### Examples

| Raw ISO-8601                    | Displayed (UTC+6 machine)    |
|---------------------------------|------------------------------|
| `2026-03-06T03:17:47+00:00`    | `06-Mar-2026 09:17 AM`       |
| `2023-12-14T11:58:26+06:00`    | `14-Dec-2023 11:58 AM`       |
| `2023-12-14T05:48:05+00:00`    | `14-Dec-2023 11:48 AM`       |

## Implementation

| File                          | Purpose                             |
|-------------------------------|-------------------------------------|
| `gitutil/dateformat.go`       | `FormatDisplayDate` function        |
| `constants/constants.go`      | `DateDisplayLayout`, `DateUTCSuffix`|

## Rules

- Commands must **never** call `time.Format` directly for display.
- All display dates flow through `gitutil.FormatDisplayDate`.
- JSON and CSV output also use the formatted date (human-friendly).
- The `FormatDisplayDateUTC` variant appends `(UTC)` when UTC-only
  output is needed.
