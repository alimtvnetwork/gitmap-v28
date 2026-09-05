# Package `openfiletype`: Standard File Open Mode Enum

`coding-guidelines/common/pkg/enum/openfiletype` provides the canonical byte-backed enum for POSIX file opening modes and flags across all Go applications, conforming to Aukgo (`03-aukgo`) and global enum guidelines.

---

## 1. Overview & Architecture

Following the repository's modular enum design, `openfiletype` provides:
- **Dedicated Package:** `openfiletype` (ends with `type` suffix).
- **Core Type:** `type Variant byte` in `variant.go`.
- **Zero-Value:** `Invalid Variant = iota` (always first).
- **Zero Stutter:** Callers write `openfiletype.Variant` and `openfiletype.CreateAppend`.

---

## 2. Variant Catalog

| Variant | Value (`byte`) | Open Flags (`os.O_*`) | Description |
| :--- | :--- | :--- | :--- |
| `Invalid` | `0` | `os.O_RDONLY` | Uninitialized / invalid state (zero value) |
| `ReadOnly` | `1` | `os.O_RDONLY` | Open existing file for read-only access |
| `WriteOnly` | `2` | `os.O_WRONLY` | Open file for write-only access |
| `ReadWrite` | `3` | `os.O_RDWR` | Open file for both reading and writing |
| `Append` | `4` | `os.O_WRONLY \| os.O_APPEND` | Open file for writing at the end |
| `CreateAppend` | `5` | `os.O_CREATE \| os.O_WRONLY \| os.O_APPEND` | Create file if missing and append data |
| `CreateTruncate` | `6` | `os.O_CREATE \| os.O_WRONLY \| os.O_TRUNC` | Create or truncate file for overwriting |
| `CreateNew` | `7` | `os.O_CREATE \| os.O_EXCL \| os.O_WRONLY` | Create new file atomically; fail if exists |
| `ReadOrCreateOnly` | `8` | `os.O_RDONLY \| os.O_CREATE` | Open existing for read or create if missing |
| `WriteOrCreateOnly` | `9` | `os.O_WRONLY \| os.O_CREATE` | Open existing for write or create if missing |
| `ReadWriteOrCreateOnly` | `10` | `os.O_RDWR \| os.O_CREATE` | Open existing for read/write or create if missing |

---

## 3. Usage Example

```go
import "coding-guidelines/common/pkg/enum/openfiletype"

func OpenTarget(path string, mode openfiletype.Variant) (*os.File, error) {
    if mode.IsInvalid() {
        return nil, fmt.Errorf("invalid mode")
    }

    return os.OpenFile(path, mode.Flags(), 0644)
}

// Caller usage:
file, err := OpenTarget("data.log", openfiletype.CreateAppend)
```

---

## 4. Helper Methods

- `.Flags() int` - Maps variant to `os.OpenFile` flags.
- `.IsValid() bool` - Checks if variant is within valid defined range (`> Invalid`).
- `.IsInvalid() bool` - Checks if variant is zero-value or undefined.
- `.IsReadOnly() bool`, `.IsCreateAppend() bool`, etc. - Positive boolean predicates.
- `All() []Variant` - Returns slice of all valid variants.
- `Parse(s string) result.Wrap[Variant]` - Case-insensitive string parser returning `result.Wrap`.
- `MarshalJSON()` / `UnmarshalJSON()` - PascalCase JSON serialization.
