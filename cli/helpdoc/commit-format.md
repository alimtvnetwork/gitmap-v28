# Commit Message Formatting & Modular Syntax Guide

GitMap supports modular commit message formatting adhering to the **Hyphen-Separated Uppercase Module Standard (Zero Third Brackets)** across all semantic commit commands.

## Supported Syntaxes

### 1. Hyphen-Separated Uppercase Module Syntax (`MODULE - <item>`)

The primary and recommended standard for domain segmentation in repository commit logs:

```bash
gitmap cpf "AUTH - implement refresh token rotation"
# Formatted: Feature: AUTH - implement refresh token rotation

gitmap cbf "SQLITE - handle connection pool lock"
# Formatted: Bug: SQLITE - handle connection pool lock

gitmap cpc "SKILLS - sync gitmap skill across repositories"
# Formatted: Chore: SKILLS - sync gitmap skill across repositories

gitmap cpr "V6.600.0 - folder-specific workdirs & ssh automation"
# Formatted: Release: V6.600.0 - folder-specific workdirs & ssh automation
```

### 2. Conventional Commit Syntax (`type(module): <item>`)

Standard conventional commit types (`feat`, `fix`, `chore`, `release`, `refactor`, `test`, `docs`, `perf`) without brackets are preserved verbatim:

```bash
gitmap cpf "feat(workgroup): add folder ssh key binding"
# Formatted: feat(workgroup): add folder ssh key binding

gitmap cbf "fix(ssh): handle missing private key"
# Formatted: fix(ssh): handle missing private key

gitmap cpc "chore(deps): update golang.org/x/crypto"
# Formatted: chore(deps): update golang.org/x/crypto
```

### 3. Automated Bracket Sanitization

Legacy bracketed inputs (e.g. `[module]: <item>`) are automatically sanitized to the uppercase hyphen standard:

```bash
gitmap cpf "[auth]: add jwt validation"
# Formatted: Feature: AUTH - add jwt validation
```

## Examples

```bash
# Feature commit with uppercase module
gitmap cpf "AUTH - implement refresh token rotation"

# Bug fix with uppercase module using cbf alias
gitmap cbf "SQLITE - handle connection pool lock"

# Chore commit with uppercase module
gitmap cpc "SKILLS - sync gitmap skill across repositories"

# Release commit with uppercase module
gitmap cpr "V6.600.0 - folder-specific workdirs & ssh automation"

# Conventional syntax with plain scope
gitmap cpf "feat(workgroup): add folder ssh key binding"
```

## Normalization Pipeline

1. **Outer Quote Sanitization:** Strips wrapping quotes (`"..."`, `'...'`) injected by shell wrappers.
2. **Duplicate Prefix Trimming:** Removes redundant leading prefixes (e.g. `cpf "Feature: AUTH - test"` normalizes cleanly to `Feature: AUTH - test`).
3. **Bracket Sanitization:** Strips any square brackets around module tokens.
4. **Module Uppercase Normalization:** Transforms module identifier to uppercase (`strings.ToUpper(mod)`).
5. **Whitespace Normalization:** Trims leading and trailing whitespace and preserves single spaces between tokens.
