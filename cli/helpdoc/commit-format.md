# Conventional Commit Formatting & Module Syntax Guide

GitMap supports modular, structured commit formatting using **Uppercase Conventional Commit syntax** across all semantic commit commands.

## Formatting Standard

All GitMap semantic commit tools (`gitmap cpf`, `gitmap cbf`, `gitmap cfr`, `gitmap cpc`, `gitmap cpr`) enforce clean, standardized commit messages in the format:

```text
TYPE(MODULE): <summary>
```

- **Type**: Strictly uppercase semantic verb (`FEAT`, `FIX`, `REFACTOR`, `CHORE`, `RELEASE`, `DOCS`, `TEST`, `PERF`).
- **Module Scope**: Uppercase module or scope enclosed in round parentheses: `(MODULE)`. Square brackets (`[...]`, also known as "third brackets") are strictly prohibited in commit titles and are automatically sanitized to round parentheses.
- **Summary**: Concise description of changes.

## Supported Input Syntaxes

GitMap's commit formatter intelligently parses multiple input formats and standardizes them:

### 1. Parenthesized Scope Syntax (`(module): <item>`)

Recommended syntax for clear domain segmentation:

```bash
gitmap cpf "(auth): implement refresh token rotation"
# Output: FEAT(AUTH): implement refresh token rotation

gitmap cbf "(sqlite): handle connection pool lock"
# Output: FIX(SQLITE): handle connection pool lock

gitmap cfr "(cache): hoist compiled regex patterns"
# Output: REFACTOR(CACHE): hoist compiled regex patterns

gitmap cpc "(ci): update go version to 1.24"
# Output: CHORE(CI): update go version to 1.24

gitmap cpr "(v6.605.0): release native toolchain wrappers"
# Output: RELEASE(V6.605.0): release native toolchain wrappers
```

### 2. Bare Colon Module Syntax (`module: <item>`)

Bare module names before a colon are automatically enclosed in parentheses and uppercased:

```bash
gitmap cpf "auth: implement refresh token rotation"
# Output: FEAT(AUTH): implement refresh token rotation

gitmap cbf "sqlite: handle connection pool lock"
# Output: FIX(SQLITE): handle connection pool lock
```

### 3. Square Bracket Normalization (`[module]: <item>`)

Square brackets (`[...]`) are automatically sanitized to round parentheses:

```bash
gitmap cpf "[auth]: add jwt validation"
# Output: FEAT(AUTH): add jwt validation
```

### 4. Legacy Hyphen Conversion (`<module> - <item>`)

Legacy hyphen-delimited messages are automatically normalized:

```bash
gitmap cpf "auth - add jwt validation"
# Output: FEAT(AUTH): add jwt validation
```

## Normalization Pipeline

1. **Outer Quote Sanitization:** Strips wrapping quotes (`"..."`, `'...'`) injected by shell wrappers.
2. **Bracket Normalization:** Converts any square brackets `[...]` to standard parentheses `(...)`.
3. **Type Capitalization:** Guarantees semantic verb prefix is uppercase (`FEAT`, `FIX`, `REFACTOR`, `CHORE`, `RELEASE`).
4. **Duplicate Prefix Trimming:** Removes redundant leading prefixes (e.g. `cpf "FEAT(AUTH): test"` normalizes cleanly without duplicating `FEAT`).
5. **Whitespace Preservation:** Trims leading/trailing whitespace and ensures single space separation.

## Examples

```bash
# Feature commit
gitmap cpf "(toolchain): add cargo and rust native runners"

# Bug fix commit
gitmap cbf "(ssh): resolve missing identity error"

# Refactoring commit
gitmap cfr "(types): extract generic models to types.go"

# Chore commit
gitmap cpc "(deps): bump x/crypto"
```
