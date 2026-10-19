# Execution Plan: Database REST Command Implementation

**Slug:** `database-rest-command-implementation`  
**Task Name:** Database REST Command Implementation  
**Status:** ACTIVE  
**Lifecycle Budget:** N = 400 Steps (Phase 1: 1-100, Phase 2: 101-250, Phase 3: 251-400)  
**Concurrency:** A = 3, H = 2, Capacity = 6  

---

## 1. User Request (Verbatim)

```text
gitmap db-rest "db path" as  alias
gitmap db-rest aliasName $port,$port2
gitmap db-rest aliasName1,aliash2 $port,$port2
gitmap db-rest aliasName1,dbPath2 $port,$port2


# Database REST Command Implementation: high priority instruction, non-negotiable task

Similar to the Database (Database (Db)) add and read command, we can have Database (Database (Db)) rest command. We enable this and give a SQL Database (Database (Db)) path, or we can give the alias and then use that alias to basically use... Then we have port, so we can give a port number here, port one or port two, either one of the port. The reason is if the first port is not available, it will try for the second port. So with comma, we can provide multiple ports. And whatever port is available, the rest will actually create some endpoints on the fly, and it will also create the Swagger so that they could test the endpoint if the user wants. So there would be a Swagger endpoint and also the communication endpoint from that database. So in this endpoint, we should be able to see all these tables. We could also see the Database (Database (Db)) rest for multiple alias, alias one, alias two. Or we could do for like Database (Database (Db)). Both cases, they should open the rest endpoint, and using this we can do so many things. For example, we can see the table, we can see the-- First of all, if we have multiple databases selected, then the Database (Database (Db)) or Database (Database (Db)) list or database list should show us all the databases that is there. And if we query inside or look inside the database, we can see the database columns, tables, and each table has which of the columns, and column has which type of data inside it. So we should also be able to modify the table if we wanted to, modify the table, like alter column, create column, or also can migrate at a migration for that database as well. And that migration can be saved inside. We can have a split Database (Database (Db)) or Database (Database (Db)) rest. And that can actually contain a caching information for a separate caching split Database (Database (Db)) for the database path we are selecting to try to have that information into that split Database (Database (Db)) first, root Database (Database (Db)), and then keep other information like the alteration or migration information. Or probably we can clone the database to a new path as well. This type of rest command needs to be available. Also, we can modify and alter the data by creating a new database on a new path, and that path would also be shared from the rest API endpoint when it is completed. We can have port one, port two, whichever port first available will be showed up in the terminal. And also it would show the Swagger endpoint where we can go and test it. And we should be able to see the table information. We should be able to find column name in any table that is out there in the databases. We can map or create views on the fly with multiple tables, queries, and things like that. Also, we should be able to perform the add, edit, delete. Also add the Database (Database (Db)) record delete command. So the previous one, previously we have done the add, edit. We didn't do it for Database (Database (Db)) record delete. We also need to have RM command for Database (Database (Db)). That would remove the database. Things like that. I think I mentioned most of the things. So try to understand these commands, try to create test cases for these commands, and then try to complete it. The rest endpoint should have more type of endpoints that we could do add record, edit record, delete record, create view for complex join. Also, direct SQL queries can be passed inside the get command for that rest to modify the database and understand the data very quickly. Also, we can have a Database (Database (Db)) rest UI. So UI will give us, again, the UI view, not like Swagger, but a little bit advanced custom view that would allow the user to interact with the database, see the tables, columns, add, edit database column information, insert row, remove row, edit row. So all kind of things should be there. It should feel like a database full visualization from the tool itself.
```

---

## 2. Architecture & Design Blueprint

### Subsystem Topology
```
CLI Dispatch:
  gitmap db-rest <targets> <ports>  ->  cmddb.RunDbRest()
  gitmap db-record-del <args>       ->  cmddb.RunDbRecordDel()
  gitmap db-rm <args>               ->  cmddb.RunDbRm()

REST Engine Architecture:
  MultiDbRestServer (cli/cmddb/rest/)
  ├── Target Resolver & Multi-DB Manager (Resolves aliases & raw file paths)
  ├── Port Scanner & Fallback Engine (Probes comma-separated ports)
  ├── Standard HTTP Mux Router (net/http, zero external dependencies)
  ├── Universal Response Envelope (PascalCase keys, Status, Attributes, Results)
  ├── REST Handlers:
  │   ├── Database Catalog: GET /api/v1/databases
  │   ├── Schema & Tables: GET /api/v1/db/{alias}/tables, GET /api/v1/db/{alias}/schema
  │   ├── Column Discovery: GET /api/v1/db/{alias}/find-columns?name={col}
  │   ├── Record CRUD: GET/POST /api/v1/db/{alias}/{table}, GET/PUT/DELETE /api/v1/db/{alias}/{table}/{id}
  │   ├── Direct SQL: GET/POST /api/v1/db/{alias}/query
  │   ├── Dynamic Views: POST/GET/DELETE /api/v1/db/{alias}/views
  │   ├── DDL Alterations: POST /api/v1/db/{alias}/alter, POST /api/v1/db/{alias}/tables/create
  │   ├── Migrations: POST /api/v1/db/{alias}/migrate
  │   └── Clone: POST /api/v1/db/{alias}/clone
  ├── Split-DB Caching:
  │   └── Stores migration logs and alteration history in dedicated cache DB
  ├── Dynamic OpenAPI 3.0 / Swagger UI:
  │   ├── Spec: GET /swagger/doc.json (synthesized on-the-fly from schema)
  │   └── UI: GET /swagger (embedded responsive Swagger UI)
  └── Custom Advanced Database Visualizer UI:
      └── UI: GET /ui (full interactive database dashboard for tables, columns, rows, views, and SQL)
```

---

## 3. Subtask Work Breakdown Matrix

| Subtask ID | Title | Disjoint Files Owned | Assigned Role | Status |
|---|---|---|---|---|
| **Subtask-01** | Git Sync, Preflight & Task Initialization | `ledger.md`, `.ai-memory/plans/` | Lead | DONE |
| **Subtask-02** | Architecture & Specification Authoring | `02-spec/21-app/database-rest-command-implementation/` | Spec Authors (A=3) | IN_PROGRESS |
| **Subtask-03** | CLI Commands (`db-rest`, `db-record-del`, `db-rm`) & Port Fallback | `cli/cmddb/cmd_rest.go`, `cli/cmddb/cmd_record_del.go`, `cli/cmddb/cmd_rm.go`, `cli/cmddb/rest_port_picker.go` | Worker 01 | PENDING |
| **Subtask-04** | REST Multi-DB Engine, CRUD, DDL, Views & SQL Querying | `cli/cmddb/rest_server.go`, `cli/cmddb/rest_handlers_*.go`, `cli/cmddb/rest_migration_cache.go` | Worker 02 | PENDING |
| **Subtask-05** | OpenAPI/Swagger Generator & Custom Database Visualizer UI | `cli/cmddb/rest_swagger.go`, `cli/cmddb/rest_ui.go`, `cli/cmddb/rest_ui_assets.go` | Worker 03 | PENDING |
| **Subtask-06** | Targeted Unit Tests & Acceptance Verification | `cli/cmddb/cmd_rest_test.go`, `cli/cmddb/cmd_record_del_test.go`, `cli/cmddb/cmd_rm_test.go` | Worker Swarm | PENDING |

---

## 4. Verification Evidence & Quality Gates
- **G1 (Port Probing):** Comma-separated ports correctly bind to the first available port.
- **G2 (Multi-DB Alias Resolution):** Single and multiple aliases/paths are resolved and mounted.
- **G3 (Universal Response Envelope):** All REST responses return PascalCase `{ Status, Attributes, Results }`.
- **G4 (Swagger & UI):** `/swagger/doc.json`, `/swagger`, and `/ui` respond with HTTP 200.
- **G5 (CRUD & DDL):** Record insert/edit/delete, view creation, clone, and migration endpoints work reliably without database deadlocks.
- **G6 (CLI Parity):** `gitmap db-record-del` and `gitmap db-rm` operate properly from terminal.
