# Branch Immutability, Clean Construction & Condition Decomposition (AI Execution Prompt)

> **/goal** Ban piecemeal struct/object mutation across conditional branches and eliminate compound mixed-polarity conditions by enforcing pure constructor helper returns and discrete affirmative boolean decomposition.
> **/learn** Understand the fragility, cognitive load, and regression hazards of dirtying partially initialized objects in `if`/`else` ladders. Master discrete pre-computed affirmative intent booleans (`is`, `has`) and compact factory helper functions.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Enforce branch immutability: completely ban field mutation of existing object/struct instances within `if`/`else` branches.
- [ ] `/learn` Delegate conditional configuration and struct creation to dedicated pure constructor helpers returning immutable literals.
- [ ] `/goal` Decompose complex expressions into discrete affirmative booleans with a hard limit of max-2 clauses per condition and zero mixed polarity.
- [ ] `/learn` Verify branch immutability and clean construction patterns across all languages via guideline linters and CI checks.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 1.0.0
**Updated:** 2026-09-29
**Applies to:** All languages (Go, TypeScript, Python, PHP, Rust)
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Overview & Core Motivation

A pervasive source of subtle regressions and high cyclomatic complexity is the **piecemeal object mutation across branching logic** anti-pattern, often paired with **compound multi-clause mixed-polarity conditions**.

In this anti-pattern:
1. An object or struct is instantiated partially or zero-initialized.
2. A sprawling `if` statement tests 3 or more conditions chained with mixed polarity (`if len >= 2 && !isX && isY`).
3. Inside the conditional branches, fields of the existing object are incrementally overwritten or defaulted.
4. Additional nested conditions mutate further fields before returning the dirty object.

This creates code where:
- State transitions are unpredictable and spread across multiple branches.
- Branch reasoning requires mental tracking of which fields were overwritten and which remained default.
- Readers must reverse-engineer compound boolean algebra with inverted logic (`!isX`).
- Refactoring or adding flags introduces regression cascades.

---

## 2. Mandatory Rules

### Rule 1: Single Assignment & Branch Immutability (TOTAL BAN on In-Place Branch Mutation)

- **Total Ban:** NEVER mutate fields on an existing struct or object instance across conditional branches (`if`/`else` ladders or `switch` cases).
- Once constructed, an object or parameter struct MUST remain immutable.
- If branching paths require different field assignments, each branch MUST call a pure helper function that constructs and returns a complete, immutable struct instance directly.

### Rule 2: Pure Constructor Helper Return Pattern

- Every branching branch that determines distinct configuration paths MUST delegate to a dedicated constructor or factory helper (e.g. `buildAlternateConfig`, `buildStandardConfig`).
- The helper MUST:
  1. Receive only the input arguments it needs.
  2. Assign all fields in a single struct/object literal.
  3. Return the fully formed, immutable instance directly.
  4. Keep body length compact (target <= 8 lines, hard cap <= 15 lines).

### Rule 3: Discrete Condition Decomposition & Max-2 Condition Rule

- **Max 2 Clauses Per Expression:** An `if` condition MUST NOT combine 3 or more logical clauses with `&&` or `||`.
- **Zero Mixed Polarity:** NEVER combine a positive check and a negative check in the same `if` condition (e.g. `isA && !isB`).
- **Affirmative Pre-Computation:**
  1. Evaluate each discrete predicate into an affirmative boolean variable (`is`, `has`).
  2. If an inverted check is required, assign it to a positive semantic counterpart (`isInvalid := !isValid`).
  3. Pre-compute the composite decision intent into a single affirmative boolean variable (`isAlternateOrder := hasEnoughTokens && isSecondValid && isFirstInvalid`).
  4. The `if` statement evaluates ONLY the single affirmative intent boolean: `if isAlternateOrder { ... }`.

---

## 3. Concrete Generic Examples

### Go Implementation

#### ❌ ANTI-PATTERN: Compound Mixed Condition & In-Place Branch Mutation

```go
// ❌ WRONG — 3 compound conditions with mixed polarity AND piecemeal mutation
func parseConnectionConfig(tokens []string, base ConfigOptions) (ConnectionParams, *appfault.AppError) {
    p := ConnectionParams{
        Timeout: base.Timeout,
    }

    if len(tokens) == 0 {
        return ConnectionParams{}, appfault.NewValidationError("missing required host")
    }

    // ❌ VIOLATION: 3 conditions chained, mixed polarity, no explanatory boolean
    if len(tokens) >= 2 && !isHostAddress(tokens[0]) && isHostAddress(tokens[1]) {
        // ❌ VIOLATION: In-place mutation of fields across branch
        p.Host = tokens[1]
        if p.Label == "" {
            p.Label = tokens[0]
        }
        if len(tokens) > 2 && p.Secret == "" {
            p.Secret = tokens[2]
        }
        return p, nil
    }

    // ❌ VIOLATION: Further fallback mutations on the same struct
    p.Host = tokens[0]
    if len(tokens) > 1 && p.Secret == "" {
        p.Secret = tokens[1]
    }
    if len(tokens) > 2 && p.Label == "" {
        p.Label = tokens[2]
    }

    return p, nil
}
```

#### ✅ CANONICAL PATTERN: Pre-Computed Intent & Pure Constructor Helpers

```go
// ✅ REQUIRED — Pre-computed affirmative intent + pure constructor helpers
func parseConnectionConfig(tokens []string, base ConfigOptions) (ConnectionParams, *appfault.AppError) {
    hasNoTokens := len(tokens) == 0
    if hasNoTokens {
        return ConnectionParams{}, appfault.NewValidationError("missing required host")
    }

    hasMultipleTokens := len(tokens) >= 2
    isFirstHost := hasMultipleTokens && isHostAddress(tokens[0])
    isSecondHost := hasMultipleTokens && isHostAddress(tokens[1])

    isFirstNotHost := !isFirstHost
    isAlternateOrder := hasMultipleTokens && isSecondHost && isFirstNotHost

    if isAlternateOrder {
        return buildAlternateConnection(tokens, base), nil
    }

    return buildStandardConnection(tokens, base), nil
}

func buildAlternateConnection(tokens []string, base ConfigOptions) ConnectionParams {
    return ConnectionParams{
        Host:    tokens[1],
        Label:   resolveTokenOrFallback(tokens, 0, base.DefaultLabel),
        Secret:  resolveTokenOrFallback(tokens, 2, base.DefaultSecret),
        Timeout: base.Timeout,
    }
}

func buildStandardConnection(tokens []string, base ConfigOptions) ConnectionParams {
    return ConnectionParams{
        Host:    tokens[0],
        Secret:  resolveTokenOrFallback(tokens, 1, base.DefaultSecret),
        Label:   resolveTokenOrFallback(tokens, 2, base.DefaultLabel),
        Timeout: base.Timeout,
    }
}

func resolveTokenOrFallback(tokens []string, index int, fallback string) string {
    isAvailable := len(tokens) > index
    if isAvailable {
        return tokens[index]
    }

    return fallback
}
```

---

### TypeScript Implementation

#### ❌ ANTI-PATTERN: Object Reassignment in Branch Ladder

```typescript
// ❌ WRONG — mutating config object inside branching paths
function resolveEndpoint(tokens: string[], base: BaseOptions): EndpointConfig {
  const config: Partial<EndpointConfig> = { timeoutMs: base.timeoutMs };

  if (tokens.length >= 2 && !isHost(tokens[0]) && isHost(tokens[1])) {
    config.host = tokens[1];
    config.label = config.label ?? tokens[0];
    if (tokens.length > 2) {
      config.secret = tokens[2];
    }
    return config as EndpointConfig;
  }

  config.host = tokens[0];
  if (tokens.length > 1) {
    config.secret = tokens[1];
  }
  return config as EndpointConfig;
}
```

#### ✅ CANONICAL PATTERN: Immutable Construction & Clean Helpers

```typescript
// ✅ REQUIRED — affirmative booleans and immutable factory functions
interface EndpointConfig {
  readonly host: string;
  readonly label: string;
  readonly secret: string;
  readonly timeoutMs: number;
}

function resolveEndpoint(tokens: string[], base: BaseOptions): EndpointConfig {
  const hasMultipleTokens = tokens.length >= 2;
  const isFirstHost = hasMultipleTokens && isHost(tokens[0]);
  const isSecondHost = hasMultipleTokens && isHost(tokens[1]);

  const isFirstNotHost = !isFirstHost;
  const isAlternateOrder = hasMultipleTokens && isSecondHost && isFirstNotHost;

  if (isAlternateOrder) {
    return createAlternateEndpoint(tokens, base);
  }

  return createStandardEndpoint(tokens, base);
}

function createAlternateEndpoint(tokens: string[], base: BaseOptions): EndpointConfig {
  return {
    host: tokens[1],
    label: tokens[0] ?? base.defaultLabel,
    secret: tokens[2] ?? base.defaultSecret,
    timeoutMs: base.timeoutMs,
  };
}

function createStandardEndpoint(tokens: string[], base: BaseOptions): EndpointConfig {
  return {
    host: tokens[0],
    secret: tokens[1] ?? base.defaultSecret,
    label: tokens[2] ?? base.defaultLabel,
    timeoutMs: base.timeoutMs,
  };
}
```

---

### Python Implementation

#### ❌ ANTI-PATTERN: Dataclass Mutation Across Condition Branches

```python
# ❌ WRONG — mutating fields in place with compound mixed conditions
def build_client_profile(tokens: list[str], base: ProfileOptions) -> ClientProfile:
    profile = ClientProfile(timeout=base.timeout)

    if len(tokens) >= 2 and not is_host(tokens[0]) and is_host(tokens[1]):
        profile.host = tokens[1]
        profile.label = tokens[0]
        if len(tokens) > 2:
            profile.secret = tokens[2]
        return profile

    profile.host = tokens[0]
    if len(tokens) > 1:
        profile.secret = tokens[1]
    return profile
```

#### ✅ CANONICAL PATTERN: Pure Dataclass Factories & Affirmative Guards

```python
# ✅ REQUIRED — affirmative boolean intent and pure factory construction
from dataclasses import dataclass

@dataclass(frozen=True)
class ClientProfile:
    host: string
    label: string
    secret: string
    timeout: int

def build_client_profile(tokens: list[str], base: ProfileOptions) -> ClientProfile:
    has_multiple_tokens = len(tokens) >= 2
    is_first_host = has_multiple_tokens and is_host(tokens[0])
    is_second_host = has_multiple_tokens and is_host(tokens[1])

    is_first_not_host = not is_first_host
    is_alternate_order = has_multiple_tokens and is_second_host and is_first_not_host

    if is_alternate_order:
        return _create_alternate_profile(tokens, base)

    return _create_standard_profile(tokens, base)

def _create_alternate_profile(tokens: list[str], base: ProfileOptions) -> ClientProfile:
    return ClientProfile(
        host=tokens[1],
        label=tokens[0] if len(tokens) > 0 else base.default_label,
        secret=tokens[2] if len(tokens) > 2 else base.default_secret,
        timeout=base.timeout,
    )

def _create_standard_profile(tokens: list[str], base: ProfileOptions) -> ClientProfile:
    return ClientProfile(
        host=tokens[0],
        secret=tokens[1] if len(tokens) > 1 else base.default_secret,
        label=tokens[2] if len(tokens) > 2 else base.default_label,
        timeout=base.timeout,
    )
```

---

## 4. Comparison & Benefits Matrix

| Aspect | ❌ In-Place Branch Mutation | ✅ Clean Constructor Return |
| :--- | :--- | :--- |
| **State Predictability** | Unpredictable; fields mutated piecemeal | 100% predictable; single complete literal |
| **Immutability** | Mutable instance exposed to race conditions | Immutable / frozen instance by design |
| **Cyclomatic Complexity** | High (nested `if` inside branches) | Flat (early returns calling discrete helpers) |
| **Cognitive Load** | Reader tracks state overrides in brain | Reader inspects self-contained factory function |
| **Function Length** | Sprawls over 25–40 lines | Main function <= 15 lines; helpers <= 8 lines |
| **Condition Readability** | Compound negative chain (`len >= 2 && !a && b`) | Single affirmative intent boolean (`if isAlternateOrder`) |
| **Testability** | Hard to isolate intermediate state overrides | Each constructor helper is independently unit-testable |

---

## 5. Cross-References

- [Boolean Principles — Parameters & Conditions](./02-boolean-principles/04-parameters-and-conditions.md) — Principle 6 & 6.2
- [Code Mutation Avoidance](./18-code-mutation-avoidance.md) — Rule 1 & Rule 7
- [Magic Values & Immutability](./26-magic-values-and-immutability.md) — Immutable by default
- [Canonical Sizing Tier](../02-canonical-size-tier.md) — Function size caps (8–15 lines)
- [Nested If Elimination](../01-cross-language/20-nesting-resolution-patterns.md) — Guard clauses and flat returns

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TYPE-032: Branch Immutability and Clean Object Construction

**Given** Conditional branching logic and object construction across Go, TypeScript, Python, PHP, and Rust codebases.
**When** Codebases are analyzed by static analyzers and guideline linters.
**Then** Zero in-place object mutations occur inside conditional branches, all configuration paths delegate to pure constructor helpers, and all branching conditions evaluate pre-computed affirmative intent booleans with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.
