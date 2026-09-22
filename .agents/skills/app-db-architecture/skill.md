---
name: app-db-architecture
description: Design, implement, and maintain Split SQLite databases, Casbin RBAC schemas, PascalCase tables, and positive boolean fields.
---

# Split SQLite & Database Architecture Guidelines

This skill guides the design, schema management, and querying of database layers in alignment with `02-spec/04-database-conventions/` and `02-spec/05-split-db-architecture/`.

## Architecture Principles

1. **Split SQLite Architecture:**
   - **Root DB (`root.db`):** Global configurations, tenants, and system-wide metadata.
   - **App DB (`app.db`):** Core application entities, domain models, and business data.
   - **Session DB (`session.db`):** Transient user sessions, ephemeral tokens, and cache entries.
   - Integrated with Casbin RBAC model for fine-grained authorization checks.

2. **Schema & Table Naming Standards:**
   - Table names must strictly use **PascalCase** (e.g. `User`, `UserRole`, `AuditLog`).
   - Primary keys must follow the `{Table}Id` format (e.g. `UserId`, `UserRoleId`).
   - Column names must use **camelCase** (e.g. `createdAt`, `updatedAt`, `displayName`).

3. **Boolean Column Discipline:**
   - Always prefix booleans with positive prefixes `is` or `has` (e.g. `isActive`, `hasPermission`).
   - For database columns representing negative states, provide computed/inversed properties rather than negative column names.

4. **Foreign Keys & Explicit Joins:**
   - Enforce explicit foreign key constraints.
   - Avoid implicit ORM joins across split database boundaries.
