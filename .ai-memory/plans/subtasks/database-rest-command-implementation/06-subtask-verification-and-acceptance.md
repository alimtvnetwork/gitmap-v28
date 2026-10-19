# Subtask 06: Targeted Unit Tests & Acceptance Verification

> **Document Version:** 1.0.0  
> **Status:** PENDING  
> **Parent Plan:** `.ai-memory/plans/database-rest-command-implementation.md`  
> **Associated Spec:** `02-spec/21-app/database-rest-command-implementation/03-swagger-ui-and-verification-spec.md`  
> **Assigned Role:** Worker Swarm / QA & Test Specialist  
> **Disjoint Files Owned:**  
> - `cli/cmddb/cmd_rest_test.go`  
> - `cli/cmddb/cmd_record_del_test.go`  
> - `cli/cmddb/cmd_rm_test.go`  

---

## 1. Objective

Architect, implement, and execute a comprehensive automated unit and integration test suite verifying all aspects of the `gitmap db-rest`, `gitmap db-record-del`, and `gitmap db-rm` command implementations. Ensure strict verification of port fallback probing, multi-alias mounting, full CRUD lifecycle, DDL alterations, dynamic views, direct SQL execution, OpenAPI 3.0 specification generation, UI asset delivery, and deadlock prevention under high concurrency.

---

## 2. Test Architecture & Scope

```
Unit & Integration Test Suite (`cli/cmddb/`)
├── `cmd_rest_test.go`
│   ├── Test 1: Port Fallback Engine & Sequential Probing
│   ├── Test 2: Target Resolver & Multi-DB Mounting
│   ├── Test 3: Schema Inspection & Column Discovery (`find-columns`)
│   ├── Test 4: End-to-End CRUD Operations (Insert, Query, Update, Delete)
│   ├── Test 5: DDL Alterations & Table Creation
│   ├── Test 6: Dynamic Views & Direct SQL Console
│   ├── Test 7: Concurrency & Single-Writer Deadlock Prevention
│   └── Test 8: Swagger Spec Generation & UI Delivery (`/swagger`, `/ui`)
├── `cmd_record_del_test.go`
│   └── Test 9: CLI Command `gitmap db-record-del`
└── `cmd_rm_test.go`
    └── Test 10: CLI Command `gitmap db-rm`
```

---

## 3. Detailed Test Suites & Scenarios

### Test Suite 1: Port Fallback Engine (`cli/cmddb/cmd_rest_test.go`)
- **`TestPortFallback_SequentialProbing`**:
  - Spin up a temporary TCP listener on candidate port `18080`.
  - Pass candidate ports `18080,18081,18082` to the port picker.
  - Assert that the engine skips `18080` and binds to `18081`.
- **`TestPortFallback_AllPortsOccupied`**:
  - Spin up temporary TCP listeners on all specified candidate ports (`18090,18091`).
  - Pass `18090,18091` to the port picker.
  - Assert that a clear, structured error is returned indicating port exhaustion.
- **`TestPortFallback_DefaultPort`**:
  - When no port arguments are provided, verify the engine starts from default port `8080` and increments sequentially.

### Test Suite 2: Multi-Alias & Path Mounting (`cli/cmddb/cmd_rest_test.go`)
- **`TestMultiDbMount_MultipleAliases`**:
  - Create two isolated SQLite databases in temporary directories: `test_alpha.db` and `test_beta.db`.
  - Initialize tables in each database (`alpha_records` in alpha, `beta_records` in beta).
  - Launch REST server mounting `alpha,beta`.
  - Query `GET /api/v1/databases`.
  - Assert HTTP 200 OK and response envelope containing both mounted aliases.
- **`TestMultiDbMount_MixedAliasAndDirectPath`**:
  - Mount an existing alias alongside a direct file path: `alpha,./test_gamma.db`.
  - Verify both databases are registered and can be addressed via their respective endpoints.

### Test Suite 3: Schema Inspection & Column Discovery (`cli/cmddb/cmd_rest_test.go`)
- **`TestSchemaInspection_PragmaTableInfo`**:
  - Create table `products` with columns: `id INTEGER PRIMARY KEY`, `name TEXT NOT NULL`, `price REAL`, `is_active BOOLEAN`.
  - Call `GET /api/v1/db/alpha/schema`.
  - Verify that column types, nullability, primary key flags, and default values match SQLite schema definition.
- **`TestColumnDiscovery_FindColumns`**:
  - Create two tables `users` (has `email`) and `logs` (has `message` and `email`).
  - Query `GET /api/v1/db/alpha/find-columns?name=email`.
  - Assert that both `users` and `logs` are returned in the results list.

### Test Suite 4: End-to-End Record CRUD Operations (`cli/cmddb/cmd_rest_test.go`)
- **`TestRecordCRUD_InsertAndGet`**:
  - Call `POST /api/v1/db/alpha/products` with payload `{"name": "Widget", "price": 19.99, "is_active": true}`.
  - Assert HTTP 201 Created and envelope containing generated Primary Key `id=1`.
  - Call `GET /api/v1/db/alpha/products/1`.
  - Assert HTTP 200 OK and field values match inserted payload.
- **`TestRecordCRUD_Update`**:
  - Call `PUT /api/v1/db/alpha/products/1` with payload `{"price": 24.99}`.
  - Assert HTTP 200 OK.
  - Call `GET /api/v1/db/alpha/products/1` and assert `price` equals `24.99`.
- **`TestRecordCRUD_Delete`**:
  - Call `DELETE /api/v1/db/alpha/products/1`.
  - Assert HTTP 200 OK.
  - Call `GET /api/v1/db/alpha/products/1`.
  - Assert HTTP 404 Not Found within standard response envelope.
- **`TestRecordCRUD_PaginationAndSorting`**:
  - Insert 30 records into `products`.
  - Query `GET /api/v1/db/alpha/products?limit=10&offset=10&sort=id&order=desc`.
  - Assert exactly 10 records are returned in descending order.

### Test Suite 5: DDL Alterations & Table Creation (`cli/cmddb/cmd_rest_test.go`)
- **`TestDDL_AddColumn`**:
  - Call `POST /api/v1/db/alpha/alter` with payload:
    ```json
    { "table": "products", "operation": "add_column", "columnName": "sku", "columnType": "TEXT", "defaultValue": "DEFAULT-SKU" }
    ```
  - Assert HTTP 200 OK.
  - Verify column `sku` exists via `PRAGMA table_info` and records can store SKU data.
- **`TestDDL_CreateTable`**:
  - Call `POST /api/v1/db/alpha/tables/create` with table name and column definition list.
  - Assert HTTP 201 Created.
  - Verify table appears in `GET /api/v1/db/alpha/tables`.

### Test Suite 6: Dynamic Views & Direct SQL Engine (`cli/cmddb/cmd_rest_test.go`)
- **`TestDynamicViews_Lifecycle`**:
  - Call `POST /api/v1/db/alpha/views` with `{"name": "v_expensive_products", "query": "SELECT * FROM products WHERE price > 20"}`.
  - Assert HTTP 201 Created.
  - Query `GET /api/v1/db/alpha/v_expensive_products`.
  - Verify returned records match the view filter condition.
  - Call `DELETE /api/v1/db/alpha/views?name=v_expensive_products`.
  - Assert HTTP 200 OK and view is removed from catalog.
- **`TestDirectSQL_SelectAndMutation`**:
  - Call `GET /api/v1/db/alpha/query?q=SELECT+COUNT(*)+AS+total+FROM+products`.
  - Assert HTTP 200 OK and correct count in results.
  - Call `POST /api/v1/db/alpha/query` with `{"sql": "UPDATE products SET price = price + 5 WHERE id = ?", "params": [1]}`.
  - Assert HTTP 200 OK and rows affected count reported.

### Test Suite 7: Concurrency & Deadlock Prevention (`cli/cmddb/cmd_rest_test.go`)
- **`TestConcurrency_DeadlockPrevention`**:
  - Enforce `conn.SetMaxOpenConns(1)` single-writer pool.
  - Launch 50 concurrent goroutines continuously executing mixed operations for 3 seconds:
    - 25 goroutines reading (`GET /api/v1/db/alpha/products`).
    - 15 goroutines inserting (`POST /api/v1/db/alpha/products`).
    - 10 goroutines updating (`PUT /api/v1/db/alpha/products/{id}`).
  - Assert zero `database is locked` errors, zero deadlocks, and zero HTTP 500 errors.

### Test Suite 8: Swagger Spec Generation & UI Delivery (`cli/cmddb/cmd_rest_test.go`)
- **`TestSwagger_OpenAPISpecValidation`**:
  - Query `GET /swagger/doc.json`.
  - Assert HTTP 200 OK and `Content-Type: application/json`.
  - Parse JSON and verify OpenAPI 3.0.3 compliance (`openapi: "3.0.3"`, valid `paths`, valid `components/schemas`).
- **`TestSwagger_HtmlDelivery`**:
  - Query `GET /swagger`.
  - Assert HTTP 200 OK and `Content-Type: text/html; charset=utf-8`.
  - Verify zero external CDN links (`https://unpkg.com`, `https://cdnjs.cloudflare.com`).
- **`TestVisualizer_HtmlDelivery`**:
  - Query `GET /ui`.
  - Assert HTTP 200 OK and `Content-Type: text/html; charset=utf-8`.
  - Verify embedded CSS and JavaScript application code.

### Test Suite 9: CLI Command `gitmap db-record-del` (`cli/cmddb/cmd_record_del_test.go`)
- **`TestCli_DbRecordDel_Success`**:
  - Create test database with record `id=42`.
  - Execute `RunDbRecordDel([]string{"test_db", "-table", "items", "-id", "42"})`.
  - Assert exit code 0 and record is deleted from SQLite table.
- **`TestCli_DbRecordDel_MissingArgs`**:
  - Execute with missing `-table` or `-id`.
  - Assert structured CLI error with usage instructions.

### Test Suite 10: CLI Command `gitmap db-rm` (`cli/cmddb/cmd_rm_test.go`)
- **`TestCli_DbRm_AliasRemoval`**:
  - Register alias `temp_alias` in test registry.
  - Execute `RunDbRm([]string{"temp_alias"})`.
  - Assert exit code 0 and alias is deleted from registry without deleting the underlying database file.
- **`TestCli_DbRm_PurgeFlag`**:
  - Execute `RunDbRm([]string{"temp_alias", "-purge"})`.
  - Assert exit code 0 and database file is deleted if explicitly confirmed.

---

## 4. Test Execution & Quality Gates

### Execution Commands
```bash
# Run all database REST test suites
go test -v ./cli/cmddb/... -run "TestPortFallback|TestMultiDbMount|TestRecordCRUD|TestDDL|TestDynamicViews|TestDirectSQL|TestConcurrency|TestSwagger|TestVisualizer|TestCli"

# Run with race detector to verify concurrency safety
go test -race -v ./cli/cmddb/... -run "TestConcurrency_DeadlockPrevention"
```

### Quality Gate Pass Criteria
- **100% Pass Rate**: All 10 test suites must pass without skipped critical checks.
- **Race Condition Clean**: Zero race condition warnings reported by Go race detector (`-race`).
- **Zero CDN Dependencies**: Asserts verify that HTML bodies contain 0 external `<script src="http...">` tags.
- **Deadlock Free**: Concurrency stress test confirms zero SQLite busy errors under sustained load.
