---
name: cg-database-conventions
description: Autonomously audits, refactors, and validates repository-wide database conventions using SQLite split-db patterns, PascalCase rules, and ORM standards.
---
# Database Conventions Coding Guidelines (`cg-database-conventions`)

This skill provides autonomous audit, refactoring, and validation of repository-wide database management based on `02-spec/04-database-conventions/`.

## Core Invariants
1. **PascalCase Singular Tables**: `User`, `Project`, not `Users`.
2. **PascalCase Columns**: `UserId`, `CreatedAt`.
3. **Primary Keys**: `{TableName}Id` (`INTEGER PRIMARY KEY AUTOINCREMENT`).
4. **Foreign Keys**: Must match the exact PK name of the referenced table.
5. **Split DB Architecture**: Use SQLite as the default.
6. **No Raw SQL**: ORM must be used.
7. **Booleans**: Must use positive `Is`/`Has` prefixes.
