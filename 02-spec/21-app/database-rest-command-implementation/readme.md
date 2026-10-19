# Database REST Command Implementation Specification Suite

**Specification Slug:** `database-rest-command-implementation`  
**Root Path:** `02-spec/21-app/database-rest-command-implementation/`  
**Status:** DRAFT / IN_PROGRESS  
**Associated Plan:** `.ai-memory/plans/database-rest-command-implementation.md`  

---

## 1. Executive Summary

This specification defines the complete technical architecture, CLI interfaces, REST API contracts, dynamic OpenAPI 3.0 / Swagger documentation generator, embedded interactive visualizer UI, and database migration/cloning caching subsystems for the GitMap `db-rest` command suite.

### Core Commands Introduced
1. `gitmap db-rest "<db-path>" as <alias>`: Bind a custom SQLite database file to an alias in the system registry.
2. `gitmap db-rest <alias-or-path> [ports]`: Start the REST server for a specific database with sequential port fallback.
3. `gitmap db-rest <alias1>,<alias2> [ports]`: Start multi-database REST service exposing all specified databases under a unified catalog.
4. `gitmap db-rest <alias1>,<path2> [ports]`: Seamlessly mix existing registered aliases and direct SQLite file paths.
5. `gitmap db-record-del <db> -table <table> -id <id>`: CLI record deletion by primary key.
6. `gitmap db-rm <alias>`: Remove an alias from the registry or safely delete an unneeded database.

---

## 2. Specification Documents Directory

| Document | Focus Area |
|---|---|
| [01-architecture-spec.md](01-architecture-spec.md) | High-level system topology, lifecycle, port binding engine, and Split-DB caching architecture. |
| [02-cli-and-rest-api-spec.md](02-cli-and-rest-api-spec.md) | CLI syntax parsing, REST API routes, Universal Response Envelope, CRUD, DDL, Views & SQL execution. |
| [03-swagger-ui-and-verification-spec.md](03-swagger-ui-and-verification-spec.md) | Dynamic OpenAPI 3.0 generation, embedded Swagger UI, custom interactive Database Visualizer UI, and acceptance tests. |

---

## 3. Mandatory Compliance Invariants
- **Universal Response Envelope**: Every REST endpoint returns `{ Status, Attributes, Results }` with PascalCase JSON keys.
- **Pure Go CGO-Free SQLite Driver**: Operates strictly on `modernc.org/sqlite`.
- **Single-Writer Safety**: `conn.SetMaxOpenConns(1)` with all rows scanned into memory before running subsequent queries to prevent deadlocks.
- **Zero External Web Frameworks**: Pure standard library `net/http` router.
- **Port Probing**: Sequential fallback over comma-separated candidate ports.
