---
name: cg-branch-immutability-and-clean-construction
description: Autonomously audits, refactors, and enforces branch immutability, clean constructor helper returns, and discrete condition decomposition across polyglot codebases.
---

# Skill: Branch Immutability, Clean Construction & Condition Decomposition (`cg-branch-immutability`)

This skill governs autonomous execution for eliminating piecemeal struct/object mutations inside conditional branches (`if`/`else` ladders or `switch` cases), compound 3+ clause boolean conditions, and mixed polarity logic chains across all source files.

## Mandatory Architectural Rules

1. **Branch Immutability (TOTAL BAN on In-Place Branch Mutation):**
   - NEVER assign or overwrite fields on an existing struct or object instance inside `if`, `else`, or `switch` branches.
   - Variables and structs must follow the **Single Assignment Principle**: assign once upon construction and never mutate afterwards.

2. **Pure Constructor Return Pattern:**
   - When different control-flow branches produce different struct configurations, each branch MUST call a dedicated constructor or builder helper (e.g. `buildAlternateConfig`, `buildStandardConfig`).
   - The helper function MUST return a newly constructed, fully-populated, immutable struct literal directly.
   - Each helper function must be compact: target <= 8 lines, hard cap <= 15 lines.

3. **Condition Decomposition & Max-2 Clause Rule:**
   - An `if` condition MUST NOT combine 3 or more logical clauses (`if a && !b && c` is strictly banned).
   - NEVER mix positive and negative checks in the same condition expression.
   - Decompose all checks into affirmative boolean variables (`is`, `has`).
   - Pre-compute the composite intent into a single affirmative boolean variable (`isAlternateOrder := hasMultipleTokens && isSecondHost && isFirstNotHost`).
   - The `if` statement evaluates ONLY the single affirmative intent boolean: `if isAlternateOrder { ... }`.

4. **Zero-Test / Zero-Build Execution Mandate:**
   - NEVER execute `go test`, `pytest`, `npm test`, or full build commands during routine execution turns.
   - Verify code using targeted file-level linters / autofixers on the specifically modified files (`exit 0`).

## Canonical Code Transformation

### ❌ Anti-Pattern: In-Place Mutation & 3+ Mixed Conditions

```go
// ❌ WRONG — mutating struct across branching conditions with compound negative checks
if len(positionals) >= 2 && !isTargetAddress(positionals[0]) && isTargetAddress(positionals[1]) {
    p.targetRaw = positionals[1]
    if p.alias == "" {
        p.alias = positionals[0]
    }
    if len(positionals) > 2 && p.password == "" {
        p.password = positionals[2]
    }
    return p, nil
}

p.targetRaw = positionals[0]
if len(positionals) > 1 && p.password == "" {
    p.password = positionals[1]
}
return p, nil
```

### ✅ Canonical Pattern: Affirmative Intent & Pure Constructor Helpers

```go
// ✅ REQUIRED — Pre-computed affirmative boolean + pure constructor helpers
hasMultipleTokens := len(positionals) >= 2
isFirstTarget := hasMultipleTokens && isTargetAddress(positionals[0])
isSecondTarget := hasMultipleTokens && isTargetAddress(positionals[1])

isFirstNotTarget := !isFirstTarget
isAlternateOrder := hasMultipleTokens && isSecondTarget && isFirstNotTarget

if isAlternateOrder {
    return buildAlternateEnrollParams(positionals, baseParams), nil
}

return buildStandardEnrollParams(positionals, baseParams), nil
```

Each helper assigns all fields in a single struct literal:

```go
func buildAlternateEnrollParams(positionals []string, base EnrollParams) EnrollParams {
    return EnrollParams{
        TargetRaw: positionals[1],
        Alias:     resolveTokenOrFallback(positionals, 0, base.Alias),
        Password:  resolveTokenOrFallback(positionals, 2, base.Password),
    }
}

func buildStandardEnrollParams(positionals []string, base EnrollParams) EnrollParams {
    return EnrollParams{
        TargetRaw: positionals[0],
        Password:  resolveTokenOrFallback(positionals, 1, base.Password),
        Alias:     resolveTokenOrFallback(positionals, 2, base.Alias),
    }
}
```

## Cross-References

- [32-branch-immutability-and-clean-construction.md](02-spec/02-coding-guidelines/01-cross-language/32-branch-immutability-and-clean-construction.md)
- [18-code-mutation-avoidance.md](02-spec/02-coding-guidelines/01-cross-language/18-code-mutation-avoidance.md)
- [04-parameters-and-conditions.md](02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/04-parameters-and-conditions.md)
