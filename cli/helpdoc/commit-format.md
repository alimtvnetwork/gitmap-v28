# Commit Message Formatting & Modular Syntax Guide

GitMap supports modular commit message formatting embracing **bracketed module syntax** and **Conventional Commit syntax** across all semantic commit commands.

## Supported Syntaxes

### 1. Bracketed Module Syntax (`[module]: <item>`)

Recommended syntax for clear domain segmentation in repository commit logs:

```bash
gitmap cpf "[auth]: implement refresh token rotation"
# Formatted: Feature: [auth]: implement refresh token rotation

gitmap cbf "[sqlite]: handle connection pool lock"
# Formatted: Bug: [sqlite]: handle connection pool lock

gitmap cpc "[skills]: sync gitmap skill across repositories"
# Formatted: Chore: [skills]: sync gitmap skill across repositories

gitmap cpr "[v6.600.0]: folder-specific workdirs & ssh automation"
# Formatted: Release: [v6.600.0]: folder-specific workdirs & ssh automation
```

### 2. Conventional Commit Syntax (`type([module]): <item>`)

Conventional commit types (`feat`, `fix`, `chore`, `release`, `refactor`, `test`, `docs`, `perf`) are preserved verbatim:

```bash
gitmap cpf "feat([workgroup]): add folder ssh key binding"
# Formatted: feat([workgroup]): add folder ssh key binding

gitmap cbf "fix([ssh]): handle missing private key"
# Formatted: fix([ssh]): handle missing private key

gitmap cpc "chore(deps): update golang.org/x/crypto"
# Formatted: chore(deps): update golang.org/x/crypto
```

### 3. Legacy Hyphen Conversion (`<module> - <item>`)

Legacy hyphen-delimited messages are automatically normalized to bracketed format:

```bash
gitmap cpf "auth - add jwt validation"
# Formatted: Feature: [auth]: add jwt validation
```

## Examples

```bash
# Feature commit with bracketed module
gitmap cpf "[auth]: implement refresh token rotation"

# Bug fix with bracketed module using cbf alias
gitmap cbf "[sqlite]: handle connection pool lock"

# Chore commit with bracketed module
gitmap cpc "[skills]: sync gitmap skill across repositories"

# Release commit with bracketed module
gitmap cpr "[v6.600.0]: folder-specific workdirs & ssh automation"

# Conventional syntax with bracketed scope
gitmap cpf "feat([workgroup]): add folder ssh key binding"
```

## Normalization Pipeline

1. **Outer Quote Sanitization:** Strips wrapping quotes (`"..."`, `'...'`) injected by shell wrappers.
2. **Duplicate Prefix Trimming:** Removes redundant leading prefixes (e.g. `cpf "Feature: [auth]: test"` normalizes cleanly to `Feature: [auth]: test`).
3. **Module Extraction:** Validates module names consisting of alphanumeric characters, hyphens, underscores, dots, and slashes (`[A-Za-z0-9_\-\./]+`).
4. **Whitespace Trimming:** Trims leading and trailing whitespace and preserves single spaces between tokens.
