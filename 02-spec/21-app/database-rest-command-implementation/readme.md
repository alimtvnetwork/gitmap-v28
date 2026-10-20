# Database REST Command Implementation Specification Suite

**Specification Slug:** `database-rest-command-implementation`  
**Root Path:** `02-spec/21-app/database-rest-command-implementation/`  
**Status:** APPROVED / IMPLEMENTED  
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

| Document | Focus Area | Lines |
|---|---|---|
| [01-architecture-overview.md](01-architecture-overview.md) | High-level system topology, lifecycle, multi-DB engine, and port fallback. | $\le 300$ |
| [02-storage-and-split-db-cache.md](02-storage-and-split-db-cache.md) | Split-DB caching architecture, schema definitions, and concurrency safety. | $\le 300$ |
| [03-cli-command-contracts.md](03-cli-command-contracts.md) | CLI command grammar, syntax parsing, banners, and companion commands. | $\le 300$ |
| [04-rest-endpoints-crud-and-ddl.md](04-rest-endpoints-crud-and-ddl.md) | Universal Response Envelope, catalog, tables/schema, and record CRUD. | $\le 300$ |
| [05-rest-views-queries-and-cloning.md](05-rest-views-queries-and-cloning.md) | Dynamic views, direct SQL queries, table alterations, migrations, and cloning. | $\le 300$ |
| [06-swagger-openapi-and-schema.md](06-swagger-openapi-and-schema.md) | Dynamic OpenAPI 3.0 generation, introspection pipeline, and type mapping. | $\le 300$ |
| [07-database-visualizer-ui.md](07-database-visualizer-ui.md) | Embedded zero-CDN Swagger UI and custom dark-theme Database Visualizer UI. | $\le 300$ |
| [08-verification-and-acceptance-matrix.md](08-verification-and-acceptance-matrix.md) | Comprehensive acceptance verification strategy, test scenarios, and quality gates. | $\le 300$ |

---

## 3. Mandatory Compliance Invariants
- **Universal Response Envelope**: Every REST endpoint returns `{ Status, Attributes, Results }` with PascalCase JSON keys.
- **Pure Go CGO-Free SQLite Driver**: Operates strictly on `modernc.org/sqlite`.
- **Single-Writer Safety**: `conn.SetMaxOpenConns(1)` with all rows scanned into memory before running subsequent queries to prevent deadlocks.
- **Zero External Web Frameworks**: Pure standard library `net/http` router.
- **Port Probing**: Sequential fallback over comma-separated candidate ports.
