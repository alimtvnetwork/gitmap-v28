# HttpMethod Enum — `pkg/enums/httpmethodtype/variant.go` (AI Execution Prompt)

> **/goal** Eliminate all magic string HTTP method literals (`"GET"`, `"POST"`, `"PUT"`, `"DELETE"`) across Go services and CLIs by standardizing on the shared `pkg/enums/httpmethodtype` enum.
> **/learn** Master the Go enum architecture: typed `byte` variant with iota, PascalCase string lookup tables, `.HttpVerb()` representation, domain predicate helpers (`.HasBody()`, `.IsSafe()`, `.IsIdempotent()`), and JSON un/marshaling.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace all raw HTTP method string literals in `http.NewRequest` and API calls with `httpmethodtype.Variant` or `.HttpVerb()`.
- [ ] `/learn` Never define local ad-hoc HTTP method constants or strings; import from the shared `pkg/enums/httpmethodtype` package.
- [ ] `/goal` Use variant checker methods (`.IsGet()`, `.IsPost()`) and semantic predicates (`.HasBody()`, `.IsSafe()`) instead of manual string comparisons.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version**: 3.0.0
> **Last updated**: 2026-02-28
> **Shared package**: All CLIs import from `pkg/enums/httpmethodtype`

---

## Purpose

Universal HTTP method enum replacing all magic string literals (`"GET"`, `"POST"`, etc.) across every CLI. This is a **shared `pkg/` package** — not internal to any single CLI.

---

## Reference Implementation

```go
// pkg/enums/httpmethodtype/variant.go
package httpmethodtype

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Variant represents an HTTP method
type Variant byte

const (
	// Invalid is the zero value (invalid/unset)
	Invalid Variant = iota

	// Get is HTTP GET
	Get

	// Head is HTTP HEAD
	Head

	// Post is HTTP POST
	Post

	// Put is HTTP PUT
	Put

	// Patch is HTTP PATCH
	Patch

	// Delete is HTTP DELETE
	Delete

	// Options is HTTP OPTIONS
	Options
)

// --- Single lookup table (PascalCase) ---

var variantLabels = [...]string{
	Invalid: "Invalid",
	Get:     "Get",
	Head:    "Head",
	Post:    "Post",
	Put:     "Put",
	Patch:   "Patch",
	Delete:  "Delete",
	Options: "Options",
}

// --- Core methods ---

func (v Variant) String() string {
	if !v.IsValid() {
		return variantLabels[Invalid]
	}

	return variantLabels[v]
}

func (v Variant) Label() string {
	return v.String()
}

func (v Variant) IsValid() bool {
	return v > Invalid && v < Variant(len(variantLabels))
}

// --- Variant checkers ---

func (v Variant) IsInvalid() bool { return v == Invalid }
func (v Variant) IsGet() bool     { return v == Get }
func (v Variant) IsHead() bool    { return v == Head }
func (v Variant) IsPost() bool    { return v == Post }
func (v Variant) IsPut() bool     { return v == Put }
func (v Variant) IsPatch() bool   { return v == Patch }
func (v Variant) IsDelete() bool  { return v == Delete }
func (v Variant) IsOptions() bool { return v == Options }

// --- Collection helpers ---

func All() []Variant {
	return []Variant{Get, Head, Post, Put, Patch, Delete, Options}
}

func ByIndex(i int) Variant {
	if i < 0 || i >= len(variantLabels) {
		return Invalid
	}

	return Variant(i)
}

func Parse(s string) apperror.Result[Variant] {
	normalized := strings.TrimSpace(s)

	for i, label := range variantLabels {
		if strings.EqualFold(label, normalized) {
			return Variant(i), nil
		}
	}

	return Invalid, fmt.Errorf("invalid HTTP method: %q", s)
}

func Values() []string {
	result := make([]string, 0, len(variantLabels)-1)

	for _, s := range variantLabels[1:] {
		result = append(result, s)
	}

	return result
}

// --- JSON serialization ---

func (v Variant) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

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

// --- Domain-specific methods ---

// HttpVerb returns the uppercase HTTP verb string for use in http.NewRequest
func (v Variant) HttpVerb() string {
	switch v {
	case Get:
		return "GET"
	case Head:
		return "HEAD"
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Patch:
		return "PATCH"
	case Delete:
		return "DELETE"
	case Options:
		return "OPTIONS"
	default:
		return ""
	}
}

// HasBody returns true if the method typically carries a request body
func (v Variant) HasBody() bool {
	return v == Post || v == Put || v == Patch
}

// IsSafe returns true if the method is safe (no side effects)
func (v Variant) IsSafe() bool {
	return v == Get || v == Head || v == Options
}

// IsIdempotent returns true if the method is idempotent
func (v Variant) IsIdempotent() bool {
	return v == Get || v == Head || v == Put || v == Delete || v == Options
}
```

---

## Usage Pattern

```go
import "myproject/pkg/enums/httpmethodtype"

// ❌ FORBIDDEN — magic string
req, err := http.NewRequestWithContext(context, "GET", url, nil)

// ✅ REQUIRED — enum usage
req, err := http.NewRequestWithContext(context, httpmethodtype.Get.HttpVerb(), url, nil)

// ✅ Variant checks
method := httpmethodtype.Post

if method.HasBody() {
    // attach body
}

// ✅ Parse from config/input
method, err := httpmethodtype.Parse("Post")

// ✅ JSON serialization — no redundant tag needed
type Request struct {
    Method httpmethodtype.Variant
}
```

---

## Variant Reference

| Constant | `.String()` | `.HttpVerb()` | `.HasBody()` | `.IsSafe()` | `.IsIdempotent()` |
|----------|-------------|---------------|--------------|-------------|---------------------|
| `Invalid` | `"Invalid"` | `""` | `false` | `false` | `false` |
| `Get` | `"Get"` | `"GET"` | `false` | `true` | `true` |
| `Head` | `"Head"` | `"HEAD"` | `false` | `true` | `true` |
| `Post` | `"Post"` | `"POST"` | `true` | `false` | `false` |
| `Put` | `"Put"` | `"PUT"` | `true` | `false` | `true` |
| `Patch` | `"Patch"` | `"PATCH"` | `true` | `false` | `false` |
| `Delete` | `"Delete"` | `"DELETE"` | `false` | `false` | `true` |
| `Options` | `"Options"` | `"OPTIONS"` | `false` | `true` | `true` |

---

## Cross-References

- Enum architecture: `02-spec/10-brun-cli/01-backend/19-enum-architecture.md` §7
- Audit report: `02-spec/23-how-app-issues-track/07-magic-string-tuple-return-audit.md`
- Enum standard memory: `architecture/enum-standard`

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-GO-003: HttpMethod Enum Standard and Magic String Elimination

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** All HTTP request construction, routing, and method comparisons utilize `pkg/enums/httpmethodtype` rather than magic string literals, with zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.
