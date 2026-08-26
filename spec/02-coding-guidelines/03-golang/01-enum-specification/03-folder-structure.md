# Folder Structure

**Version:** 3.2.0  
**Status:** Complete  
**Updated:** 2026-04-16

---

## Package Naming Convention — `type` Suffix

> **Convention over configuration**: By looking at a package name, you must immediately know it is an enum.

All Go enum packages MUST end with the `type` suffix. This makes enum packages instantly recognizable without inspecting their contents.

| ✅ Correct (with `type` suffix) | ❌ Wrong (no suffix) |
|----------------------------------|----------------------|
| `providertype` | `provider` |
| `platformtype` | `platform` |
| `httpmethodtype` | `httpmethod` |
| `executionstatustype` | `executionstatus` |
| `searchmodetype` | `search_mode` |
| `outputformattype` | `output_format` |

**Go rule:** Package names are lowercase with no underscores — the `type` suffix is appended directly (e.g., `providertype`, not `provider_type`).

**TypeScript rule:** File names use kebab-case with `-type` suffix (e.g., `http-method-type.ts`). The enum name itself keeps PascalCase without suffix (e.g., `HttpMethod`), since the `enum` keyword already signals the construct.

---

## Standard Layout

All enums MUST be placed in the `internal/enums/` directory (or `pkg/enums/` for shared packages) at the project root. Each package name ends with `type`.

```
{cli-root}/
├── cmd/
│   └── root.go
├── internal/
│   ├── enums/
│   │   ├── providertype/
│   │   │   └── variant.go
│   │   ├── platformtype/
│   │   │   └── variant.go
│   │   ├── enginetype/
│   │   │   └── variant.go
│   │   ├── searchmodetype/
│   │   │   └── variant.go
│   │   ├── outputformattype/
│   │   │   └── variant.go
│   │   └── registry.go       # Optional: Central registry
│   ├── models/
│   ├── services/
│   └── api/
└── pkg/
    └── enums/
        └── httpmethodtype/
            └── variant.go
```

---

## Naming Conventions

### Package Names

- **Must** end with `type` suffix
- Lowercase, no underscores (Go convention)
- One enum per package

| ✅ Correct | ❌ Wrong | Reason |
|-----------|----------|--------|
| `providertype` | `provider` | Missing `type` suffix |
| `searchmodetype` | `search_mode` | Underscore + missing suffix |
| `outputformattype` | `OutputFormat` | PascalCase + missing suffix |
| `movieprovidertype` | `movieProvider` | camelCase + missing suffix |
| `httpmethodtype` | `httpmethod` | Missing `type` suffix |

### File Names

- Always name the file `variant.go`
- Additional helper files allowed: `helpers.go`, `validation.go`

```
internal/enums/providertype/
├── variant.go      # Main enum definition
├── helpers.go      # Optional: Helper functions
└── validation.go   # Optional: Validation logic
```

---

## Import Pattern

```go
import (
    "myapp/internal/enums/providertype"
    "myapp/internal/enums/platformtype"
    "myapp/internal/enums/enginetype"
)

func main() {
    p := providertype.SerpApi

    if p.IsSerpApi() {
        // ...
    }
    
    platforms := []platformtype.Variant{
        platformtype.YouTube,
        platformtype.Reddit,
    }
}
```

---

## Enum Categories by CLI

### GSearch CLI

```
internal/enums/
├── providertype/           # SerpApi, MapsScraper, Colly
├── platformtype/           # YouTube, Reddit, LinkedIn, etc.
├── enginetype/             # Google, Bing, DuckDuckGo
├── searchmodetype/         # Sequential, Parallel, RoundRobin
├── outputformattype/       # JSON, CSV, Table, Markdown
├── movieprovidertype/      # Tmdb, Omdb, Trakt, ImdbScraper
├── socialmediatype/        # LinkedIn, Twitter, Instagram, etc.
└── contenttype/            # Web, Image, Video, News
```

### BRun CLI

```
internal/enums/
├── buildtype/              # Debug, Release, Test
├── runmodetype/            # Foreground, Background, Watch
├── logleveltype/           # Debug, Info, Warn, Error
└── profiletype/            # Development, Staging, Production
```

### AI Bridge CLI

```
internal/enums/
├── modelprovidertype/      # Ollama, OpenAI, Anthropic
├── reasoningmodetype/      # SinglePrompt, TwoStage, Research
├── steptype/               # Search, Fetch, Parse, Embed, etc.
├── executionstatustype/    # Pending, Running, Success, Failed
├── checkpointtype/         # Auto, Manual, Rollback
└── memoryflagtype/         # IsCritical, IsImportant, Standard
```

### Nexus Flow CLI

```
internal/enums/
├── nodetype/               # Start, End, Task, Decision, Fork, Join
├── flowstatustype/         # Draft, Active, Paused, Completed
├── triggertypetype/        # Manual, Scheduled, Webhook, Event
└── executionmodetype/      # Sequential, Parallel
```

### Spec Reverse CLI

```
internal/enums/
├── outputformattype/       # Markdown, JSON, YAML
├── parsertype/             # Go, TypeScript, Python
└── extractionmodetype/     # Full, Summary, Skeleton
```

### WP SEO Publish CLI

```
internal/enums/
├── contenttype/            # Post, Page, Product
├── publishstatustype/      # Draft, Pending, Published
├── seoscoretype/           # Poor, Fair, Good, Excellent
└── mediatype/              # Image, Video, Document
```

### AI Transcribe CLI

```
internal/enums/
├── audioformattype/        # MP3, WAV, FLAC, OGG
├── transcribeprovidertype/ # Whisper, DeepGram, AssemblyAI
├── outputformattype/       # SRT, VTT, TXT, JSON
└── languagetype/           # EN, ES, FR, DE, etc.
```

---

## Registry Pattern (Optional)

For CLIs with many enums, create a central registry:

```go
// internal/enums/registry.go
package enums

import (
    "myapp/internal/enums/providertype"
    "myapp/internal/enums/platformtype"
    "myapp/internal/enums/enginetype"
)

// Re-export for convenience
type (
    Provider   = providertype.Variant
    Platform   = platformtype.Variant
    Engine     = enginetype.Variant
)

// Constants re-export
const (
    ProviderSerpApi     = providertype.SerpApi
    ProviderMapsScraper = providertype.MapsScraper
    ProviderColly       = providertype.Colly
    
    PlatformYouTube     = platformtype.YouTube
    PlatformReddit      = platformtype.Reddit
    
    EngineGoogle        = enginetype.Google
    EngineBing          = enginetype.Bing
)
```

**Usage:**
```go
import "myapp/internal/enums"

p := enums.ProviderSerpApi
```

---

## File Template

```go
// internal/enums/{category}type/variant.go
package {category}type

import (
    "encoding/json"
    "fmt"
    "strings"
)

// Variant represents a {category} type
type Variant byte

const (
    // Invalid is the zero value
    Invalid Variant = iota
    
    // Add variants here...
)

// variantLabels maps each variant to its PascalCase display string.
// Label values MUST match the constant name exactly.
var variantLabels = [...]string{
    Invalid: "Invalid",
    // Add mappings — PascalCase matching constant names...
}

// String returns the PascalCase string representation
func (v Variant) String() string {
    if v.IsInvalid() {
        return variantLabels[Invalid]
    }

    return variantLabels[v]
}

// Label delegates to String
func (v Variant) Label() string {
    return v.String()
}

// IsValid checks if the variant is within the valid range
func (v Variant) IsValid() bool {
    return v > Invalid && v < Variant(len(variantLabels))
}

// IsInvalid checks if the variant is the zero value or out of range
func (v Variant) IsInvalid() bool {
    return v <= Invalid || v >= Variant(len(variantLabels))
}

// Add Is{Value}() methods for each variant...

// All returns all valid variants
func All() []Variant {
    result := make([]Variant, 0, len(variantLabels)-1)

    for i := 1; i < len(variantLabels); i++ {
        result = append(result, Variant(i))
    }

    return result
}

// ByIndex returns variant by index
func ByIndex(i int) Variant {
    isOutOfRange := i < 0 || i >= len(variantLabels)

    if isOutOfRange {
        return Invalid
    }

    return Variant(i)
}

// Parse parses a string to variant using case-insensitive matching
func Parse(s string) apperror.Result[Variant] {
    trimmed := strings.TrimSpace(s)

    for i, label := range variantLabels {
        if strings.EqualFold(label, trimmed) {
            return Variant(i), nil
        }
    }

    return Invalid, fmt.Errorf("invalid {category}: %q", s)
}

// Values returns all valid string values
func Values() []string {
    result := make([]string, 0, len(variantLabels)-1)

    for _, s := range variantLabels[1:] {
        result = append(result, s)
    }

    return result
}

// MarshalJSON implements json.Marshaler
func (v Variant) MarshalJSON() ([]byte, error) {
    return json.Marshal(v.String())
}

// UnmarshalJSON implements json.Unmarshaler
func (v *Variant) UnmarshalJSON(data []byte) error {
    var s string

    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }

    parsed, err := Parse(s)

    if err != nil {
        return err
    }

    *v = parsed

    return nil
}
```

---

*Folder structure standard for enum organization.*
