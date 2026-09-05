# Package `logleveltype`: Standard Logging Severity Level Enum

`coding-guidelines/common/pkg/enum/logleveltype` provides the canonical byte-backed enum for log severity levels and thresholds across all Go applications, conforming to Aukgo (`03-aukgo`) and global enum guidelines.

---

## 1. Overview & Architecture

Following the repository's modular enum design, `logleveltype` provides:
- **Dedicated Package:** `logleveltype` (ends with `type` suffix).
- **Core Type:** `type Variant byte` in `variant.go`.
- **Zero-Value:** `Invalid Variant = iota` (aliased to `Unknown`).
- **Zero Stutter:** Callers write `logleveltype.Variant` and `logleveltype.Debug`, `logleveltype.Info`, etc.

---

## 2. Variant Catalog

| Variant | Value (`byte`) | Name / String | Description |
| :--- | :--- | :--- | :--- |
| `Invalid` / `Unknown` | `0` | `"Unknown"` | Uninitialized or invalid log level (zero value) |
| `Debug` | `1` | `"Debug"` | Verbose diagnostic information |
| `Info` | `2` | `"Info"` | Routine informational operational messages |
| `Warn` | `3` | `"Warn"` | Warning messages for potential issues |
| `Error` | `4` | `"Error"` | Error events requiring attention |
| `Fatal` | `5` | `"Fatal"` | Critical failures leading to immediate termination |

---

## 3. Usage Example

```go
import "coding-guidelines/common/pkg/enum/logleveltype"

func ShouldLog(current, minLevel logleveltype.Variant) bool {
    return current.IsEnabled(minLevel)
}

// Caller usage:
if logleveltype.Info.IsEnabled(logleveltype.Debug) {
    // Info meets Debug threshold
}
```

---

## 4. Helper Methods

- `.Name() string` / `.Label() string` / `.String() string` - PascalCase name.
- `.IsValid() bool` - Checks if variant is within valid defined range (`> Invalid`).
- `.IsInvalid() bool` - Checks if variant is zero-value or undefined.
- `.IsEnabled(threshold Variant) bool` - Checks if current level meets or exceeds threshold (`>= threshold`).
- `.IsDebug() bool`, `.IsInfo() bool`, etc. - Positive boolean predicates.
- `All() []Variant` - Returns slice of all valid variants (`Debug` through `Fatal`).
- `Values() []string` - Returns slice of valid variant names.
- `Parse(s string) result.Wrap[Variant]` - Case-insensitive string parser returning `result.Wrap`.
- `MarshalJSON()` / `UnmarshalJSON()` - PascalCase JSON serialization with string and byte fallback.
