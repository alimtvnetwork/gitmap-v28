# Internationalization

## Purpose

Standardise how applications detect locale, externalise user-facing
strings, format dates and numbers, and support right-to-left (RTL)
scripts so that adding a new language never requires code changes
beyond a translation file.

## Locale Detection

### Rules

1. Detect locale from environment (`LANG`, `LC_ALL`) or browser
   (`navigator.language`).
2. Fall back to `en-US` when detection fails.
3. Store the resolved locale in a single, importable constant — never
   re-detect per call site.
4. Allow explicit override via `--locale` flag (CLI) or config key.

### Go Pattern

```go
func DetectLocale() string {
    for _, key := range []string{"LC_ALL", "LANG"} {
        if val := os.Getenv(key); val != "" {
            return normalise(val)
        }
    }
    return constants.DefaultLocale // "en-US"
}
```

### TypeScript Pattern

```ts
export function detectLocale(): string {
  return navigator.language ?? "en-US";
}
```

## String Externalisation

### Rules

1. **No user-visible literal strings in logic files** — all text lives
   in a locale resource (JSON, constants file, or message catalogue).
2. Keys use dot-separated namespaces: `section.action.detail`
   (e.g. `clone.error.offline`).
3. Interpolation uses named placeholders, never positional:
   `"Cloned {repoName} in {duration}"`.
4. Pluralisation uses ICU-style rules or explicit plural keys
   (`items.one`, `items.other`).
5. Constants files group messages by feature module.

### Go Pattern

```go
// constants/messages_clone.go
const (
    MsgCloneStart   = "Cloning %s into %s"
    MsgCloneSuccess = "Cloned %s (%s)"
    MsgCloneOffline = "Cannot clone — network unavailable"
)
```

### TypeScript Pattern

```ts
// locales/en-US.json
{
  "clone": {
    "start": "Cloning {repoName} into {targetDir}",
    "success": "Cloned {repoName} ({duration})",
    "offline": "Cannot clone — network unavailable"
  }
}
```

## Date & Number Formatting

### Rules

1. Store all dates in **UTC** internally.
2. Display dates in the user's **local timezone** using the project's
   standard format: `DD-Mon-YYYY hh:mm AM/PM`
   (Go layout: `02-Jan-2006 03:04 PM`).
3. Use `Intl.DateTimeFormat` (TS) or `time.LoadLocation` (Go) —
   never manual offset arithmetic.
4. Numbers use locale-aware grouping (`Intl.NumberFormat` or
   `message.NewPrinter`).
5. Formatting functions live in a single utility module
   (`gitutil/dateformat.go`, `lib/format.ts`).

### Go Pattern

```go
func FormatLocalDate(t time.Time) string {
    return t.Local().Format(constants.DateDisplayFormat)
}
```

### TypeScript Pattern

```ts
export function formatDate(date: Date, locale: string): string {
  return new Intl.DateTimeFormat(locale, {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  }).format(date);
}
```

## RTL Support

### Rules

1. Set `dir="rtl"` on the root element when locale is RTL.
2. Use **logical CSS properties** (`margin-inline-start` instead of
   `margin-left`; `padding-inline-end` instead of `padding-right`).
3. Icons that imply direction (arrows, chevrons) must flip in RTL.
4. Never use `text-align: left` — use `text-align: start`.
5. Test every layout in at least one RTL locale before release.

### RTL Locale Detection

```ts
const rtlLocales = new Set(["ar", "he", "fa", "ur"]);

export function isRTL(locale: string): boolean {
  return rtlLocales.has(locale.split("-")[0]);
}
```

### Tailwind RTL

```html
<!-- Use Tailwind logical utilities -->
<div class="ms-4 me-2 text-start">
  <!-- ms = margin-inline-start, me = margin-inline-end -->
</div>
```

## Translation Workflow

| Step | Action                                    |
|------|-------------------------------------------|
| 1    | Developer adds key to default locale file |
| 2    | CI extracts new/changed keys              |
| 3    | Translator fills target locale files      |
| 4    | CI validates all keys present in all locales |
| 5    | Missing keys fall back to `en-US`         |

## Constants

| Constant           | Value     | Description                  |
|--------------------|-----------|------------------------------|
| `DefaultLocale`    | `"en-US"` | Fallback locale              |
| `DateDisplayFormat`| (see above) | Standard date layout       |
| `MaxKeyLength`     | 80        | Message key character limit  |

## Constraints

- Files ≤ 200 lines, functions 8–15 lines.
- No user-visible strings outside locale resources / constants.
- All dates stored as UTC, displayed as local.
- Logical CSS properties only — no physical `left`/`right` for layout.
- Every new locale must pass CI key-completeness check.

## References

- [03 Naming Conventions](../05-coding-guidelines/03-naming-conventions.md)
- [07 Logging & Observability](../05-coding-guidelines/07-logging-observability.md)
- [12 Documentation Standards](../05-coding-guidelines/12-documentation-standards.md)
