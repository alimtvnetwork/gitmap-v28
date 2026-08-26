# Braces, Nesting & Exemptions

> **Version:** 4.0.0  
> **Updated:** 2026-03-31  
> **Applies to:** PHP, TypeScript, Go  
> **Rules covered:** 1, 2, 7

---

## Rule 1: Always Use Braces — No Single-Line Statements

Every `if`, `for`, `foreach`/`for...of`, `while` block **must** use curly braces `{}`, even for single-statement bodies.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN
if ($this->initialized) return;

if ($error === null) return false;

// ✅ REQUIRED
if ($this->initialized) {
    return;
}

if ($error === null) {
    return false;
}
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN
if (isLoading) return null;

// ✅ REQUIRED
if (isLoading) {
    return null;
}
```

```go
// ── Go ───────────────────────────────────────────────────────
// Go enforces braces by syntax — this rule is already satisfied.
```

---

## Rule 2: Zero Nested `if` — Absolute Ban

Nested `if` blocks are **absolutely forbidden** — zero tolerance, no exceptions. Every nested `if` must be flattened using one of: (a) combined conditions, (b) early returns, (c) extracted helper functions. If a helper function already handles the null/empty check internally, rely on it — don't wrap it in a redundant outer guard.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: Nested if — redundant null guard
if ($error !== null) {
    if (ErrorChecker::isFatalError($error)) {
        $this->logger->fatal($error);
    }
}

// ✅ REQUIRED: Flat — isFatalError() handles null internally
if (ErrorChecker::isFatalError($error)) {
    $this->logger->fatal($error);
}

// ✅ ALSO OK: Early return to flatten
if ($request === null) {
    return;
}

if ($request->has_param('file')) {
    $this->process($request);
}
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN: Nested if
if (response) {
    if (response.status >= 400) {
        handleError(response);
    }
}

// ✅ REQUIRED: Early return or combined condition
if (!response) {
    return;
}

if (response.status >= 400) {
    handleError(response);
}
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN: Nested if
if err != nil {
    if resp != nil {
        handleError(resp)
    }
}

// ❌ ALSO FORBIDDEN: Multiple negative checks combined
if err != nil && resp != nil {
    handleError(resp)
}

// ✅ REQUIRED: Positive named booleans, flat
hasError := err != nil
hasResponse := resp != nil
hasIssue := hasError && hasResponse

if hasIssue {
    handleError(resp)
}
```

### Exemption: `switch`-Based Enum Parsers

`FromString` / `fromString` enum-parsing functions that use a `switch` statement with a `default` fallback are **exempt** from the zero-nesting rule. The `switch` is the idiomatic, most readable pattern for exhaustive enum mapping — flattening into sequential `if` blocks would reduce clarity without improving safety.

> **AI Rule — No String Literals in Switch Cases:**  
> Switch cases **must never** use raw string literals (`"production"`, `"staging"`, etc.).  
> Instead, each enum value **must** define a corresponding string constant, and the `case` labels **must** reference those constants. This prevents typos, enables compile-time safety, and ensures the string representation is defined once alongside the enum — not scattered across parser functions.
>
> **For full enum patterns, folder structure, and required methods, see:**  
> - **Go:** [Enum Specification](../../03-golang/01-enum-specification/00-overview.md) — pattern, required methods, folder layout  
> - **TypeScript:** [TypeScript Enums](../../02-typescript/00-overview.md) — string enum definitions, usage patterns  
> - **PHP:** [PHP Enums](../../04-php/01-enums.md) — backed enum classes

```go
// ── Go ───────────────────────────────────────────────────────
// Each enum lives in its own package (e.g., `internal/enums/environmenttype/`).
// The package name provides the grouping — no type-name prefix needed on constants.
// See: ../../03-golang/01-enum-specification/03-folder-structure.md

package environmenttype

type Variant byte

const (
    Invalid     Variant = iota
    Production
    Staging
    Development
)

// String constants — defined ONCE, co-located with the enum in the same package
const (
    ProductionStr  = "production"
    StagingStr     = "staging"
    DevelopmentStr = "development"
)

// ❌ FORBIDDEN: Raw string literals in switch cases
func FromString(s string) (Variant, error) {
    switch strings.ToLower(s) {
    case "production":
        return Production, nil
    case "staging":
        return Staging, nil
    default:
        return Invalid, apperror.New(apperror.ErrInvalidEnum, "unknown environment: "+s)
    }
}

// ✅ REQUIRED: Package-scoped string constants in switch cases
func FromString(s string) (Variant, error) {
    switch strings.ToLower(s) {
    case ProductionStr:
        return Production, nil
    case StagingStr:
        return Staging, nil
    case DevelopmentStr:
        return Development, nil
    default:
        return Invalid, apperror.New(
            apperror.ErrInvalidEnum,
            "unknown environment: "+s,
        )
    }
}

// Callers use the package name as the grouping prefix:
//   environmenttype.Production   (not EnvironmentProduction)
//   environmenttype.FromString() (not EnvironmentFromString())
```

```php
// ── PHP (8.1+ backed enum) ──────────────────────────────────
// See: ../../04-php/01-enums.md

// PHP backed enums auto-compile from string via ::from() / ::tryFrom().
// No manual switch/match with string literals needed.
enum Environment: string {
    case Production  = 'production';
    case Staging     = 'staging';
    case Development = 'development';
}

// ✅ REQUIRED: Let the backed enum handle parsing
$env = Environment::from($input);        // throws ValueError if invalid
$env = Environment::tryFrom($input);     // returns null if invalid
```

```typescript
// ── TypeScript ──────────────────────────────────────────────
// See: ../../02-typescript/00-overview.md for all enum definitions

// ❌ FORBIDDEN: Raw string literals in switch
switch (input) {
    case "production": return Environment.Production;
    case "staging":    return Environment.Staging;
}

// ✅ REQUIRED: Use enum values (string enums auto-compile)
enum Environment {
    Production  = "production",
    Staging     = "staging",
    Development = "development",
}

// TypeScript string enums allow direct comparison against enum members:
function parseEnvironment(input: string): Environment {
    const match = Object.values(Environment).find((v) => v === input);

    if (!match) {
        throw new Error(`Unknown environment: ${input}`);
    }

    return match;
}
```

**Conditions for exemption (from zero-nesting rule):**
- Function is a pure enum parser (`FromString`, `fromString`, `Parse<Enum>`)
- Uses `switch` (Go) or `match` (PHP 8.0+) — not chained `if/else`
- Has an explicit `default` / `default =>` branch that returns an error
- Contains no other logic beyond the mapping

**Mandatory for all languages:**
- Switch/match cases must use **enum constants or defined string constants** — never raw string literals
- String representations are defined **once**, co-located with the enum type definition
- **Go:** Constants use package-scoped names (`ProductionStr`), not type-prefixed names (`EnvironmentProductionStr`) — the package name (`environmenttype`) already provides the grouping

### Exemption: Optional-Field Validation (isDefined → isValid)

When validating an **optional input field**, a two-level nested `if` is permitted where:
- The **outer** `if` checks existence/presence (`isDefined`, `IsDefined`, `has_param`)
- The **inner** `if/else` checks validity and returns an error on the `else` branch

This pattern is exempt because flattening would require evaluating validity on undefined fields, which is semantically incorrect — the field should only be validated if it exists.

```go
// ✅ EXEMPT: Optional-field validation — outer checks existence, inner checks validity
if input.Config.IsDefined() {
    if input.Config.IsDefinedAndValid() {
        applyConfig(input.Config)
    } else {
        return apperror.FailNew[Site](apperror.ErrValidation, "invalid config")
    }
}
```

**Conditions for exemption:**
- Outer condition is a pure existence/presence check
- Inner condition validates the same field
- The `else` branch returns/propagates an error
- No additional nesting beyond 2 levels

---

## Rule 7: Zero Nested `if` — Absolute Ban (Reinforced)

This is a **reinforcement of Rule 2** with stricter language. Nested `if` blocks are the single biggest readability killer. There is **zero tolerance** — any code review finding a nested `if` is an automatic rejection.

### Flattening Techniques

| Nesting Pattern | Flattening Technique |
|----------------|---------------------|
| Null guard → logic | Early return for null |
| Permission → action | Early return for no permission |
| Multiple conditions | Combined `&&` (extract if 2+ operators) |
| If-inside-loop | Extract loop body to helper function |
| If-inside-if-inside-if | Extract to dedicated method |

```php
// ❌ FORBIDDEN: Triple nesting
if ($request !== null) {
    if ($request->hasParam('file')) {
        if ($this->isValidFile($request->getParam('file'))) {
            $this->process($request);
        }
    }
}

// ✅ REQUIRED: Flat with early returns
if ($request === null) {
    return;
}

$hasValidFile = $request->hasParam('file')
    && $this->isValidFile($request->getParam('file'));

if ($hasValidFile) {
    $this->process($request);
}
```

---

*Part of [Code Style](./00-overview.md) — Rules 1, 2, 7*
