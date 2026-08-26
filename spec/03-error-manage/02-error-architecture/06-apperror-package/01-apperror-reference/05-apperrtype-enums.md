# AppError Package Reference — Domain Error Type Enums (`apperrtype`)

> **Parent:** [AppError Package Reference](./00-overview.md)  
> **Version:** 2.0.0  
> **Updated:** 2026-04-02  
> **Breaking:** Migrated from per-domain `byte` enums to single `uint16` Variation enum (inspired by [evatix-go/errorwrapper/errtype](https://gitlab.com/auk-go/errorwrapper/-/tree/develop/errtype))

---

## Purpose

Canonical reference for the `apperrtype` package. All error types are defined as variants of a single `Variation` enum (`uint16`), with a global registry map providing `Name`, `Code`, and `Message` for each variant. This eliminates raw string codes and centralizes error metadata.

---

## Package Structure

```
types/apperrtype/
├── variation.go              ← Variation uint16 type + all variant constants
├── variant_structure.go      ← VariantStructure struct + display/error methods
├── error_type_interface.go   ← ErrorType interface
├── variant_registry.go       ← Global variantRegistry map[Variation]VariantStructure
├── string_to_variant_map.go  ← Reverse-lookup map: string name → Variation
├── code_to_variant_map.go    ← Reverse-lookup map: string code → Variation
├── consts.go                 ← minValue / maxValue bounds
└── format_consts.go          ← Format strings for display methods
```

---

## Core Types

### Variation (uint16 enum)

Single enum type for all error variants. Uses `uint16` to support 400+ variants across all domains.

```go
// apperrtype/variation.go
package apperrtype

type Variation uint16

const (
    NoError Variation = iota  // 0 — zero value, "no error"

    // ── E1xxx — Configuration ──────────────────────────
    ConfigFileMissing          // E1001
    ConfigParseFailure         // E1002
    ConfigKeyMissing           // E1003
    EnvVarMissing              // E1004

    // ── E2xxx — Database / Site / Plugin ───────────────
    DBConnectionFailed         // E2001
    DBQueryFailed              // E2002
    DBRecordNotFound           // E2003
    DBDuplicateKey             // E2004
    DBMigrationFailed          // E2005
    SiteNotFound               // E2010
    SiteBlocked                // E2011
    PluginSlugMissing          // E2012
    PluginNotFound             // E2013
    PluginAlreadyActive        // E2014

    // ── E3xxx — WordPress API ──────────────────────────
    WPConnectionFailed         // E3001
    WPAuthFailed               // E3002
    WPEndpointNotFound         // E3003
    WPRateLimited              // E3004
    WPResponseInvalid          // E3005

    // ── E4xxx — File System ────────────────────────────
    FileNotFound               // E4001
    FileReadFailed             // E4002
    FileWriteFailed            // E4003
    DirCreateFailed            // E4004
    PermissionDenied           // E4005
    FileNotExist               // E4006
    FileAccessFailed           // E4007
    FileAppendFailed           // E4008
    DirNotExist                // E4009
    DirAccessFailed            // E4010
    EmptyFilePath              // E4011
    PathNotFound               // E4012
    PathInvalid                // E4013
    PathStatFailed             // E4014
    SymlinkFailed              // E4015
    PathMissing                // E4016
    PathFailedToCreate         // E4017
    PathFailedToRead           // E4018
    PathFailedToWrite          // E4019
    PathFailedToDelete         // E4020

    // ── E5xxx — Sync ──────────────────────────────────
    SyncConflict               // E5001
    SyncTimeout                // E5002
    SyncChecksumFail           // E5003
    SyncLockAcquire            // E5004
    SyncStateMismatch          // E5005
    SyncOutOfSync              // E5006

    // ── E6xxx — Backup ────────────────────────────────
    BackupCreateFailed         // E6001
    BackupRestoreFailed        // E6002
    BackupNotFound             // E6003
    BackupCorrupted            // E6004
    BackupQuotaExceeded        // E6005

    // ── E7xxx — Git ───────────────────────────────────
    GitCloneFailed             // E7001
    GitPushFailed              // E7002
    GitPullFailed              // E7003
    GitMergeConflict           // E7004
    GitRepoNotFound            // E7005
    GitCommitFailed            // E7006

    // ── E8xxx — Build ─────────────────────────────────
    BuildCompileFailed         // E8001
    BuildDependencyMissing     // E8002
    BuildTimeout               // E8003
    BuildArtifactFailed        // E8004
    BuildConfigInvalid         // E8005
    BuildTranspileFailed       // E8006

    // ── E9xxx — General ───────────────────────────────
    InternalError              // E9001
    ValidationFailed           // E9002
    NotImplemented             // E9003
    Unauthorized               // E9004
    RateLimited                // E9005
    InvalidInput               // E9006
    InvalidOutput              // E9007
    InvalidCondition           // E9008
    UnexpectedValue            // E9009
    UnexpectedType             // E9010
    OutOfRangeValue            // E9011
    CastingFailed              // E9012
    NullOrEmpty                // E9013
    MismatchExpectation        // E9014
    NotFound                   // E9015
    CrudOperationFailed        // E9016
    MappingFailed              // E9017
    ParsingFailed              // E9018
    SerializationFailed        // E9019

    // ── E10xxx — E2E Test ─────────────────────────────
    E2ESetupFailed             // E10001
    E2EAssertFailed            // E10002
    E2ETimeoutFailed           // E10003
    E2EFixtureFailed           // E10004

    // ── E11xxx — Publish ──────────────────────────────
    PublishFailed              // E11001
    PublishConflict            // E11002
    PublishRollback            // E11003
    PublishTargetDown          // E11004

    // ── E12xxx — Version ──────────────────────────────
    VersionNotFound            // E12001
    VersionConflict            // E12002
    VersionParseFail           // E12003
    VersionLocked              // E12004

    // ── E13xxx — Session ──────────────────────────────
    SessionExpired             // E13001
    SessionNotFound            // E13002
    SessionInvalid             // E13003
    SessionLimitHit            // E13004

    // ── E14xxx — Crypto ───────────────────────────────
    CryptoEncryptFailed        // E14001
    CryptoDecryptFailed        // E14002
    CryptoKeyInvalid           // E14003
    CryptoHashMismatch         // E14004
    CryptoChecksumFailed       // E14005

    // ── E15xxx — Network / Connection ─────────────────
    NetworkOffline             // E15001
    ConnectionFailed           // E15002
    ConnectionTimeout          // E15003
    Disconnected               // E15004
    RequestFailed              // E15005

    // ── E16xxx — Process / Execution ──────────────────
    ProcessFailed              // E16001
    CommandExecutionFailed     // E16002
    ScriptFailed               // E16003
    LockFailed                 // E16004
    StepFailed                 // E16005
    CompletionFailed           // E16006

    // ── E17xxx — Encoding / Conversion ────────────────
    EncodingFailed             // E17001
    DecodingFailed             // E17002
    MarshalFailed              // E17003
    UnmarshalFailed            // E17004
    ConversionFailed           // E17005

    // ── E18xxx — Permission / Authorization ───────────
    PermissionFailed           // E18001
    AuthorizationFailed        // E18002
    AccessDenied               // E18003
    ResourceFrozen             // E18004

    // ── Sentinel ──────────────────────────────────────
    MaxError                   // must remain last
)
```

### VariantStructure (display + error creation)

Rich metadata struct for each variant, matching the `evatix-go/errorwrapper` pattern:

```go
// apperrtype/variant_structure.go
package apperrtype

import (
    "errors"
    "fmt"
    "strings"
)

// VariantStructure holds the name, code, message, and variant for each error type.
type VariantStructure struct {
    Name    string     // PascalCase variant name (e.g. "SiteNotFound")
    Code    string     // String code (e.g. "E2010")
    Message string     // Human-readable default message
    Variant Variation  // The enum value itself
}

// String returns "Name (Code - N) : Message"
func (it VariantStructure) String() string {
    return it.TypeNameCodeMessage()
}

// TypeNameCodeMessage returns formatted: "Name (Code - N) : Message"
func (it VariantStructure) TypeNameCodeMessage() string {
    return fmt.Sprintf("%s (Code - %d) : %s", it.Name, it.Variant, it.Message)
}

// CodeTypeName returns "(#N - Name)"
func (it VariantStructure) CodeTypeName() string {
    return fmt.Sprintf("(#%d - %s)", it.Variant, it.Name)
}

// CodeTypeNameWithMessage returns "(#N - Name) message"
func (it VariantStructure) CodeTypeNameWithMessage(msg string) string {
    return fmt.Sprintf("(#%d - %s) %s", it.Variant, it.Name, msg)
}

// Error creates a stdlib error with additional context message.
func (it VariantStructure) Error(additionalMessage string) error {
    msg := fmt.Sprintf("[%s] %s: %s", it.Code, it.Message, additionalMessage)
    return errors.New(strings.ToLower(msg))
}

// ErrorNoRefs creates a stdlib error with no additional context.
func (it VariantStructure) ErrorNoRefs() error {
    msg := fmt.Sprintf("[%s] %s", it.Code, it.Message)
    return errors.New(strings.ToLower(msg))
}

// Panic panics with formatted message + additional context.
func (it VariantStructure) Panic(additionalMessage string) {
    panic(it.Error(additionalMessage))
}
```

### ErrorType Interface

```go
// apperrtype/error_type_interface.go
package apperrtype

// ErrorType is the interface all error type enums must implement.
type ErrorType interface {
    Code() string
    Message() string
    Name() string
}
```

### Variation Methods (implements ErrorType)

```go
// Methods on Variation — delegates to the global registry

func (v Variation) Code() string {
    if vs, ok := variantRegistry[v]; ok {
        return vs.Code
    }
    return "E0000"
}

func (v Variation) Message() string {
    if vs, ok := variantRegistry[v]; ok {
        return vs.Message
    }
    return "unknown error"
}

func (v Variation) Name() string {
    if vs, ok := variantRegistry[v]; ok {
        return vs.Name
    }
    return "Unknown"
}

func (v Variation) String() string {
    if vs, ok := variantRegistry[v]; ok {
        return vs.TypeNameCodeMessage()
    }
    return fmt.Sprintf("Unknown (Code - %d)", v)
}

func (v Variation) CodeTypeName() string {
    return fmt.Sprintf("(#%d - %s)", v, v.Name())
}

func (v Variation) CodeTypeNameWithReferences(refs ...string) string {
    return fmt.Sprintf("(#%d - %s) %s", v, v.Name(), strings.Join(refs, ", "))
}

func (v Variation) IsValid() bool {
    return v > NoError && v < MaxError
}

// Structure returns the full VariantStructure from the registry.
func (v Variation) Structure() VariantStructure {
    return variantRegistry[v]
}
```

---

## Global Variant Registry

Single source-of-truth map for all error variants — enables iteration, lookup, and serialization:

```go
// apperrtype/variant_registry.go
package apperrtype

var variantRegistry = map[Variation]VariantStructure{
    // ── E1xxx — Configuration ──
    ConfigFileMissing:  {Name: "ConfigFileMissing",  Code: "E1001", Message: "configuration file not found",            Variant: ConfigFileMissing},
    ConfigParseFailure: {Name: "ConfigParseFailure", Code: "E1002", Message: "failed to parse configuration",           Variant: ConfigParseFailure},
    ConfigKeyMissing:   {Name: "ConfigKeyMissing",   Code: "E1003", Message: "required configuration key missing",      Variant: ConfigKeyMissing},
    EnvVarMissing:      {Name: "EnvVarMissing",      Code: "E1004", Message: "required environment variable not set",   Variant: EnvVarMissing},

    // ── E2xxx — Database / Site / Plugin ──
    DBConnectionFailed: {Name: "DBConnectionFailed", Code: "E2001", Message: "database connection failed",     Variant: DBConnectionFailed},
    DBQueryFailed:      {Name: "DBQueryFailed",      Code: "E2002", Message: "database query failed",          Variant: DBQueryFailed},
    DBRecordNotFound:   {Name: "DBRecordNotFound",   Code: "E2003", Message: "record not found",               Variant: DBRecordNotFound},
    DBDuplicateKey:     {Name: "DBDuplicateKey",     Code: "E2004", Message: "duplicate key violation",        Variant: DBDuplicateKey},
    DBMigrationFailed:  {Name: "DBMigrationFailed",  Code: "E2005", Message: "database migration failed",     Variant: DBMigrationFailed},
    SiteNotFound:       {Name: "SiteNotFound",       Code: "E2010", Message: "site not found",                 Variant: SiteNotFound},
    SiteBlocked:        {Name: "SiteBlocked",        Code: "E2011", Message: "site is blocked",                Variant: SiteBlocked},
    PluginSlugMissing:  {Name: "PluginSlugMissing",  Code: "E2012", Message: "plugin slug required",           Variant: PluginSlugMissing},
    PluginNotFound:     {Name: "PluginNotFound",     Code: "E2013", Message: "plugin not found",               Variant: PluginNotFound},
    PluginAlreadyActive:{Name: "PluginAlreadyActive", Code: "E2014", Message: "plugin is already active",     Variant: PluginAlreadyActive},

    // ── E3xxx — WordPress API ──
    WPConnectionFailed: {Name: "WPConnectionFailed", Code: "E3001", Message: "WordPress connection failed",    Variant: WPConnectionFailed},
    WPAuthFailed:       {Name: "WPAuthFailed",       Code: "E3002", Message: "WordPress authentication failed", Variant: WPAuthFailed},
    WPEndpointNotFound: {Name: "WPEndpointNotFound", Code: "E3003", Message: "WordPress endpoint not found",  Variant: WPEndpointNotFound},
    WPRateLimited:      {Name: "WPRateLimited",      Code: "E3004", Message: "WordPress API rate limited",    Variant: WPRateLimited},
    WPResponseInvalid:  {Name: "WPResponseInvalid",  Code: "E3005", Message: "invalid WordPress API response", Variant: WPResponseInvalid},

    // ... (remaining domains follow identical pattern)
    // Full registry continues for E4xxx–E18xxx with same structure

    // ── E4xxx — File System (path variants) ──
    PathMissing:        {Name: "PathMissing",        Code: "E4016", Message: "required path is missing",   Variant: PathMissing},
    PathFailedToCreate: {Name: "PathFailedToCreate",  Code: "E4017", Message: "failed to create path",     Variant: PathFailedToCreate},
    PathFailedToRead:   {Name: "PathFailedToRead",    Code: "E4018", Message: "failed to read path",       Variant: PathFailedToRead},
    PathFailedToWrite:  {Name: "PathFailedToWrite",   Code: "E4019", Message: "failed to write to path",   Variant: PathFailedToWrite},
    PathFailedToDelete: {Name: "PathFailedToDelete",  Code: "E4020", Message: "failed to delete path",     Variant: PathFailedToDelete},
}
```

### Bounds

```go
// apperrtype/consts.go
package apperrtype

const (
    minValue = int(NoError)
    maxValue = int(MaxError)
)
```

### StringToVariantMap (reverse-lookup)

Auto-generated reverse map from variant name strings to `Variation` values — enables deserialization, CLI input parsing, and config-driven error selection:

```go
// apperrtype/string_to_variant_map.go
package apperrtype

// StringToVariantMap maps PascalCase variant names to their Variation enum values.
// Built once at init from variantRegistry — never modified at runtime.
var StringToVariantMap map[string]Variation

func init() {
    StringToVariantMap = make(map[string]Variation, len(variantRegistry))
    for v, vs := range variantRegistry {
        StringToVariantMap[vs.Name] = v
    }
}

// VariationFromName looks up a Variation by its PascalCase name.
// Returns (variation, true) if found, (NoError, false) otherwise.
func VariationFromName(name string) (Variation, bool) {
    v, ok := StringToVariantMap[name]
    return v, ok
}

// MustVariationFromName looks up a Variation by name, panics if not found.
// Use only during initialization / config parsing.
func MustVariationFromName(name string) Variation {
    v, ok := StringToVariantMap[name]
    if !ok {
        panic(fmt.Sprintf("apperrtype: unknown variant name %q", name))
    }
    return v
}
```

**Usage examples:**

```go
// Reverse-lookup from string name
v, ok := apperrtype.VariationFromName("SiteNotFound")
// v → apperrtype.SiteNotFound, ok → true

// Config-driven error type selection
errTypeName := cfg.Get("default_error_type") // "DBConnectionFailed"
v, ok := apperrtype.VariationFromName(errTypeName)
if !ok {
    log.Fatalf("unknown error type in config: %s", errTypeName)
}

// Deserialization from JSON/API payload
v := apperrtype.MustVariationFromName(payload.ErrorType) // panics if invalid

// Iterate all known variant names
for name, variation := range apperrtype.StringToVariantMap {
    fmt.Printf("%s → %d\n", name, variation)
}
```

### CodeToVariantMap (code-based reverse-lookup)

Auto-generated reverse map from string error codes to `Variation` values — enables lookup from API responses, log entries, and error code references:

```go
// apperrtype/code_to_variant_map.go
package apperrtype

// CodeToVariantMap maps string error codes (e.g. "E2010") to their Variation enum values.
// Built once at init from variantRegistry — never modified at runtime.
var CodeToVariantMap map[string]Variation

func init() {
    CodeToVariantMap = make(map[string]Variation, len(variantRegistry))
    for v, vs := range variantRegistry {
        CodeToVariantMap[vs.Code] = v
    }
}

// VariationFromCode looks up a Variation by its string error code.
// Returns (variation, true) if found, (NoError, false) otherwise.
func VariationFromCode(code string) (Variation, bool) {
    v, ok := CodeToVariantMap[code]
    return v, ok
}

// MustVariationFromCode looks up a Variation by code, panics if not found.
// Use only during initialization / config parsing.
func MustVariationFromCode(code string) Variation {
    v, ok := CodeToVariantMap[code]
    if !ok {
        panic(fmt.Sprintf("apperrtype: unknown variant code %q", code))
    }
    return v
}
```

**Usage examples:**

```go
// Reverse-lookup from error code string
v, ok := apperrtype.VariationFromCode("E2010")
// v → apperrtype.SiteNotFound, ok → true

// Parse error code from API response or log line
v, ok := apperrtype.VariationFromCode(apiResp.ErrorCode)
if ok {
    fmt.Println(v.Name(), v.Message())
}

// Strict lookup during initialization
v := apperrtype.MustVariationFromCode("E3001") // panics if unknown
```

---

The `apperror` constructors accept `Variation` values directly through the `ErrorType` interface:

```go
// apperror/constructors.go — NewType uses the enum's built-in code + message
func NewType(errType apperrtype.ErrorType) *AppError {
    return New(errType.Code(), errType.Message())
}
```

### Usage Examples

```go
// ❌ Level 1 — raw strings (flagged by CODE-RED-008 lint rule)
apperror.New("E2010", "site not found")

// ✅ Level 2 — enum code, manual message
apperror.New(apperrtype.SiteNotFound.Code(), "site not found")

// ✅✅ Level 3 — enum with built-in message (best)
apperror.NewType(apperrtype.SiteNotFound)

// ✅ WrapType — wrap a raw error with enum code + message
apperror.WrapType(err, apperrtype.WPConnectionFailed).
    WithValue("url", siteURL).
    WithStatusCode(resp.StatusCode)
```

### Display Method Examples

```go
// String() output:
// "SiteNotFound (Code - 15) : site not found"
fmt.Println(apperrtype.SiteNotFound.String())

// CodeTypeName() output:
// "(#15 - SiteNotFound)"
fmt.Println(apperrtype.SiteNotFound.CodeTypeName())

// Structure() access:
vs := apperrtype.SiteNotFound.Structure()
fmt.Println(vs.Code)    // "E2010"
fmt.Println(vs.Name)    // "SiteNotFound"
fmt.Println(vs.Message) // "site not found"

// Direct error creation from VariantStructure:
err := apperrtype.SiteNotFound.Structure().Error("domain example.com")
// → "[e2010] site not found: domain example.com"
```

### With Convenience Constructors (Type Aliases)

```go
return apperror.FailBool(apperror.NewType(apperrtype.SiteNotFound))
return apperror.FailSettings(apperror.NewType(apperrtype.ConfigKeyMissing))
```

---

## Rules

| Rule | Description |
|------|-------------|
| Single enum type | `type Variation uint16` — all variants in one file |
| `ErrorType` interface | Variation implements `Code()` + `Message()` + `Name()` |
| Global registry map | `variantRegistry` maps Variation → `VariantStructure{Name, Code, Message, Variant}` |
| No raw string codes | Use `apperrtype.Xxx` variants; `CODE-RED-008` lint enforces this |
| Package location | All types live in `types/apperrtype/` |
| `iota` start at 0 | `NoError = 0` (zero value), domain variants start at 1+ |
| `MaxError` sentinel | Must remain last constant for bounds checking |
| Domain comments | Each domain block is delimited with `// ── Exxxx — Domain ──` comments |

---

## Migration from v1.x (per-domain byte enums)

| v1.x (old) | v2.0 (new) |
|-------------|------------|
| `type ConfigError byte` | `Variation uint16` (all domains combined) |
| `configErrorDetails map[ConfigError]ErrorDetail` | `variantRegistry map[Variation]VariantStructure` |
| `ErrorDetail{Code, Message}` | `VariantStructure{Name, Code, Message, Variant}` |
| Per-file `Code()` + `Message()` | Centralized on `Variation` type via registry lookup |
| No display methods | `String()`, `TypeNameCodeMessage()`, `CodeTypeName()`, `Error()`, `Panic()` |

**Breaking changes:**
- Import path unchanged (`types/apperrtype`)
- Variant names unchanged (`apperrtype.SiteNotFound`, `apperrtype.DBQueryFailed`)
- `ErrorType` interface gains `Name() string` method
- `ErrorDetail` struct replaced by `VariantStructure`

---

## Relationship to the Project Error Code Registry

The ecosystem uses **two complementary error code systems**:

| System | Format | Scope | Example | Used By |
|--------|--------|-------|---------|---------|
| **apperrtype Variation** | `E{x}xxx` (string) | `apperror` package internals — domain-level Go errors | `E2010`, `E5003` | Go backend handlers via `apperror.NewType()` |
| **Project Registry** | `PREFIX-NNNN` or `PREFIX-NNN-NN` (prefixed integer) | Project-level allocation — cross-stack error tracking | `AB-9301`, `GEN-100-02` | All projects, all stacks (Go, React, PHP) |

### Key Distinctions

1. **No overlap risk** — `E{x}xxx` codes are string-typed and never collide with the numeric project registry codes
2. **Different granularity** — Variation covers low-level Go domain errors; the registry covers project-scoped allocations
3. **Complementary usage** — a Go handler returns `E2010` via `apperror.NewType(apperrtype.SiteNotFound)`; the project registry tracks broader allocation
4. **Both must be registered** — new variants go in this file AND in the [Error Code Registry](../../../03-error-code-registry/01-registry.md)

---

## Adding a New Variant

1. Add the constant to the appropriate domain block in `variation.go` (before `MaxError`)
2. Add the `VariantStructure` entry to `variantRegistry` in `variant_registry.go`
3. Register the code in the [Error Code Registry](../../../03-error-code-registry/01-registry.md)
4. Update this spec's domain block if adding a new domain range

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Error Code Convention | [04-codes-and-policy.md](./04-codes-and-policy.md) |
| Error Code Registry | [03-error-code-registry/01-registry.md](../../../03-error-code-registry/01-registry.md) |
| Enum Specification | [03-golang/01-enum-specification](../../../../02-coding-guidelines/03-golang/01-enum-specification/00-overview.md) |
| Lint Rule CODE-RED-008 | `linter-scripts/validate-guidelines.go` |
| Source Inspiration | [evatix-go/errorwrapper/errtype](https://gitlab.com/auk-go/errorwrapper/-/tree/develop/errtype) |

---
