# 06-swagger-openapi-and-schema: Dynamic OpenAPI 3.0 Engine & Schema Mapping

- **Spec ID:** `02-spec/21-app/database-rest-command-implementation/06-swagger-openapi-and-schema.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Dynamic OpenAPI 3.0 Engine, Schema Introspection, Type Mapping
- **Parent Plan:** `.ai-memory/plans/database-rest-command-implementation.md`
- **Dependencies:** `cli/cmddb`, `cli/store`, `cli/appfault`

---

## 1. Executive Summary & Architecture Scope

This specification establishes technical requirements and schema synthesis logic for the dynamic OpenAPI 3.0 generator:
1. **Dynamic Introspection**: Inspects `sqlite_master` and `PRAGMA table_info` for all mounted databases.
2. **Zero External CDN Dependencies**: All assets are pre-compiled and embedded directly within the Go binary (`cli/cmddb/rest_ui_assets.go`), ensuring full offline functionality.
3. **Universal Response Envelope**: All API endpoints return the standardized envelope structure `{ Status, Attributes, Results }` using PascalCase JSON fields.
4. **Single-Writer Safety**: SQLite database connections maintain single-writer concurrency (`conn.SetMaxOpenConns(1)`).

---

## 2. Dynamic OpenAPI 3.0 Specification Engine

The dynamic OpenAPI engine generates an OpenAPI 3.0.3 compliant specification JSON document for all mounted SQLite databases upon request at `/swagger/doc.json`.

### 2.1 Introspection Pipeline

```text
HTTP GET /swagger/doc.json
  │
  ├── 1. Acquire MultiDbManager Read Lock
  ├── 2. Iterate Mounted Databases (aliases and direct paths)
  ├── 3. For each Database:
  │     ├── Query sqlite_master for Tables and Views
  │     └── For each Table/View:
  │           PRAGMA table_info("{tableName}")
  │           Extract: cid, name, type, notnull, dflt_value, pk
  ├── 4. Map SQLite Data Types to OpenAPI 3.0 Types & Formats
  ├── 5. Synthesize Paths, Parameters, Request Bodies & Component Schemas
  └── 6. Serialize to JSON and Serve with Content-Type: application/json
```

---

## 3. SQLite to OpenAPI 3.0 Type Mapping Matrix

SQLite uses dynamic type affinity. The generator maps declared column types to strict OpenAPI 3.0 data types:

| SQLite Declared Type | Normalization Rule | OpenAPI Type | OpenAPI Format | Example Default Value |
|---|---|---|---|---|
| `INTEGER`, `INT`, `BIGINT`, `SMALLINT`, `TINYINT` | Contains `INT` | `integer` | `int64` | `0` |
| `REAL`, `FLOAT`, `DOUBLE`, `DECIMAL`, `NUMERIC` | Contains `REAL`, `FLOA`, `DOUB` | `number` | `double` | `0.0` |
| `TEXT`, `VARCHAR`, `NVARCHAR`, `CHAR`, `CLOB` | Contains `TEXT`, `CHAR`, `CLOB` | `string` | *(none)* | `""` |
| `BOOLEAN`, `BOOL` | Contains `BOOL` | `boolean` | *(none)* | `false` |
| `DATETIME`, `TIMESTAMP`, `DATE` | Contains `DATE`, `TIME` | `string` | `date-time` | `2026-01-01T00:00:00Z` |
| `BLOB` | Contains `BLOB` or Empty | `string` | `byte` | `""` |
| Unspecified / Unknown | Fallback | `string` | *(none)* | `""` |

---

## 4. OpenAPI Document Structure Overview

The synthesized OpenAPI document incorporates all system endpoints and dynamic table endpoints:
- **Paths Generated**:
  - `/api/v1/databases` (Catalog)
  - `/api/v1/db/{alias}/tables` (Tables & Views)
  - `/api/v1/db/{alias}/schema` (Schema & DDL)
  - `/api/v1/db/{alias}/find-columns` (Column Search)
  - `/api/v1/db/{alias}/{table}` (GET List, POST Insert)
  - `/api/v1/db/{alias}/{table}/{id}` (GET, PUT, DELETE)
  - `/api/v1/db/{alias}/query` (GET Read-only, POST Raw SQL)
  - `/api/v1/db/{alias}/views` (GET List, POST Create, DELETE)
  - `/api/v1/db/{alias}/alter` (POST Alter Table)
  - `/api/v1/db/{alias}/tables/create` (POST Create Table)
  - `/api/v1/db/{alias}/migrate` (POST Migration Runner)
  - `/api/v1/db/{alias}/clone` (POST Database Clone)
- **Component Schemas**:
  - `UniversalResponseEnvelope`: Standardized response structure.
  - Per-table DTOs generated dynamically based on column names and types.
