# 08-verification-and-acceptance-matrix: Acceptance Verification & Quality Gates

- **Spec ID:** `02-spec/21-app/database-rest-command-implementation/08-verification-and-acceptance-matrix.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Acceptance Verification, Testing Strategy, Traceability Matrix
- **Parent Plan:** `.ai-memory/plans/database-rest-command-implementation.md`
- **Dependencies:** `cli/cmddb`, `cli/store`, `cli/appfault`

---

## 1. Test Suite Organization

| Test Suite | Implementation File | Verification Focus |
|---|---|---|
| **Port Probing & Fallback** | `cli/cmddb/cmd_rest_test.go` | Sequential port availability checks, occupied port skipping, error handling |
| **Multi-Alias Mounting** | `cli/cmddb/cmd_rest_test.go` | Single alias, comma-separated aliases, mixed alias/path resolution, catalog |
| **Schema & Discovery** | `cli/cmddb/cmd_rest_test.go` | `PRAGMA table_info` extraction, column search across tables (`find-columns`) |
| **Record CRUD Engine** | `cli/cmddb/cmd_rest_test.go` | Pagination, sorting, insertion, primary key update, deletion, 404 responses |
| **DDL & Migrations** | `cli/cmddb/cmd_rest_test.go` | `alter` table add column, create table, migration script caching in split-db |
| **Dynamic Views & SQL** | `cli/cmddb/cmd_rest_test.go` | View creation/query/drop, direct SQL GET/POST, SQL error envelopes |
| **Deadlock & Concurrency** | `cli/cmddb/cmd_rest_test.go` | `SetMaxOpenConns(1)` enforcement, parallel goroutines read/write stress |
| **Swagger & UI Delivery** | `cli/cmddb/cmd_rest_test.go` | `/swagger/doc.json` OpenAPI 3.0 validity, `/swagger` and `/ui` HTTP 200 delivery |
| **CLI Commands Parity** | `cli/cmddb/cmd_record_del_test.go`, `cli/cmddb/cmd_rm_test.go` | `gitmap db-record-del` and `gitmap db-rm` command execution |

---

## 2. Detailed Verification Scenarios

### Scenario 1: Sequential Port Fallback Engine
1. Occupy port with a temporary listener.
2. Invoke port picker with candidate ports.
3. Assert that the engine automatically binds to the next available candidate.
4. Occupy all candidate ports and assert descriptive port exhaustion error.

### Scenario 2: Multi-Alias & Path Resolution
1. Create two test SQLite databases.
2. Mount using command syntax: `main,analytics`.
3. Query `GET /api/v1/databases`.
4. Verify response contains both aliases with their file paths.
5. Mount mixed syntax: `main,./test_extra.db` and verify auto-derived alias.

### Scenario 3: End-to-End CRUD Lifecycle
1. Execute `POST /api/v1/db/main/users` with record payload.
2. Verify HTTP 201 Created and response containing generated Primary Key `id=1`.
3. Execute `GET /api/v1/db/main/users/1`. Verify all fields.
4. Execute `PUT /api/v1/db/main/users/1` and verify update.
5. Execute `DELETE /api/v1/db/main/users/1`. Verify HTTP 200 OK.
6. Verify subsequent `GET` returns HTTP 404 Not Found.

### Scenario 4: DDL Alteration & Column Discovery
1. Create table `inventory`.
2. Execute `GET /api/v1/db/main/find-columns?name=sku` and verify table returned.
3. Execute `POST /api/v1/db/main/alter` to add column `quantity`.
4. Confirm column `quantity` exists via `PRAGMA table_info`.

### Scenario 5: Dynamic Views & SQL Console
1. Execute `POST /api/v1/db/main/views` to create dynamic view.
2. Query `GET /api/v1/db/main/{viewName}` and verify results.
3. Execute direct query via `POST /api/v1/db/main/query` and verify telemetry and rows.

### Scenario 6: Deadlock Prevention Under Heavy Concurrency
1. Configure SQLite connection with `conn.SetMaxOpenConns(1)`.
2. Launch 50 concurrent goroutines executing simultaneous read/write operations.
3. Verify zero instances of `database is locked`, zero deadlocks, and zero goroutine panics.

### Scenario 7: OpenAPI 3.0 Document & UI Delivery
1. Query `GET /swagger/doc.json` and validate OpenAPI 3.0.3 schema.
2. Query `GET /swagger` and `GET /ui` and verify HTTP 200 OK and self-contained zero-CDN delivery.

### Scenario 8: CLI Command Parity
1. Run `gitmap db-record-del` and verify row is deleted by PK.
2. Run `gitmap db-rm` and verify alias is removed from registry.

---

## 3. Traceability Matrix & Quality Gates

| Quality Gate | Requirement | Verification Method | Pass Criteria |
|---|---|---|---|
| **G1: Port Probing** | Comma-separated port fallback | Unit Test `TestPortFallback` | Binds to first free port; errors when all busy |
| **G2: Multi-DB Mount** | Multiple aliases and direct paths | Unit Test `TestMultiDbMount` | All DBs mounted and queryable via `/api/v1/databases` |
| **G3: Universal Envelope** | PascalCase `{ Status, Attributes, Results }` | Integration Test `TestEnvelopeFormat` | 100% of REST responses match envelope format |
| **G4: OpenAPI 3.0 & Swagger** | Dynamic spec and embedded UI | Integration Test `TestSwaggerDocAndUI` | `/swagger/doc.json` valid OpenAPI 3.0; `/swagger` returns HTTP 200 |
| **G5: Visualizer UI** | Interactive dark mode dashboard | Integration Test `TestVisualizerUI` | `/ui` returns HTTP 200; zero CDN references |
| **G6: CRUD & DDL Engine** | Insert, read, update, delete, alter | Integration Test `TestCrudAndDdl` | Record mutations and table alterations succeed |
| **G7: Deadlock Prevention** | Single-writer safety under concurrency | Stress Test `TestConcurrentAccess` | 50 concurrent goroutines with 0 lock errors |
| **G8: CLI Parity** | `db-record-del` and `db-rm` commands | CLI Test `TestCliCommands` | Commands execute and modify SQLite state cleanly |
