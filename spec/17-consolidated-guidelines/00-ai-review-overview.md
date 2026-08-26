# AI Instruction Manual & Overview

**Version:** 1.0.0
**Updated:** 2026-08-22
**Purpose:** Guide AI agents on how to read, interpret, and enforce the coding guidelines.

## 1. What File is What?

- **`31-compiled-simple-coding-guidelines.md`**: The strict, zero-tolerance canonical checklist for coding standards. Every AI MUST read and obey this file.
- **`04-enum-standards.md`**: Detailed rules on enum construction, suffixing, and switch-statement usage.
- **`03-error-management.md`**: Complete workflow for error capturing, typed errors, and logging context.
- **`21-lovable-folder-structure.md`**: Standard folder structures across workspaces.

## 2. How to Read the Guidelines

When reviewing or generating code:
1. Always start with `31-compiled-simple-coding-guidelines.md`. It overrides conflicting rules found elsewhere.
2. If the task involves a specific domain (like Database or UI), locate the corresponding compiled spec in this folder (e.g., `18-database-conventions.md`).
3. Assume rules are Zero-Tolerance unless a specific `lint-allow` exception is documented. 

## 3. How to Review Code (Styling & Rules)

When acting as a reviewer or when writing code, strictly check against:
- **Zero Generated Code in Git**: Ensure `__pycache__`, `*.pyc`, `*.generated.*` are ignored in `.gitignore`. Proactively update `.gitignore` files if you see a new project being set up.
- **Boolean Strictness**: Max 1 logical join (`&&` / `||`), no mixed polarity. Check CODE-RED-014.
- **No Swallowed Errors**: Every catch block must log context.
- **No Magic Strings/Numbers**: Enforce the use of typed constants and enums.
- **Flattened Logic**: No nested `if` statements. Extract into guard clauses.

## 4. Continuous Maintenance

- **.gitignore Maintenance**: AI agents are explicitly responsible for maintaining the `.gitignore` to prevent cache files and generated artifacts from sneaking into version control. If a build tool produces artifacts, add them to `.gitignore` immediately.
