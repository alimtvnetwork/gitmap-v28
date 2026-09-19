---
name: cg-enum-standards
description: Autonomously audits and enforces cross-language enum standards (Go, TS, PHP, Rust), ensuring PascalCase formatting, Type suffixes, and generator workflow usage.
---
# Enum Standards Coding Guidelines (`cg-enum-standards`)

This skill provides autonomous audit and validation of repository-wide enum definitions based on `02-spec/17-consolidated-guidelines/07-enum-standards.md`.

## Core Invariants
1. **Source of Truth**: Enums are defined ONCE in `02-spec/<module>/enums/<EnumName>.yaml`.
2. **Generators**: Must use `03-ai-scripts/30-enum-generator.py` to generate enums for Go, TS, PHP, Rust. Never hand-edit `*_generated.*` files.
3. **Naming**: Enums must use PascalCase. PHP enums must have `Type` suffix.
4. **Rust Derives**: Rust enums must include `Debug, Clone, PartialEq, Serialize, Deserialize`.
