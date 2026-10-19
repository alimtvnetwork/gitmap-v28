# gitmap db-rest

Launch dynamic REST API server, Swagger UI, and custom Database Visualizer across SQLite databases.

## Aliases

`dbrs`, `db rest`

## Usage

```bash
gitmap db-rest "<db-path>" as <alias>
gitmap db-rest <alias-or-path> [ports]
gitmap db-rest <alias1>,<alias2> [ports]
gitmap db-rest <alias1>,<path2> [ports]
```

## Description

`gitmap db-rest` launches a high-speed, zero-dependency REST API server and interactive UI for one or more SQLite databases. It supports dynamic OpenAPI 3.0 (Swagger) generation, an embedded Swagger UI, and a dedicated dark-themed Database Visualizer UI.

Sequential port probing automatically tests candidate ports (e.g. `8080,8081,3000`) and binds to the first open port.

## Endpoints

| Endpoint | Method | Description |
|---|---|---|
| `/api/v1/databases` | GET | List mounted databases with metadata and table counts |
| `/api/v1/db/:alias/tables` | GET | List tables and views with row and column counts |
| `/api/v1/db/:alias/schema` | GET | Full DDL schema introspection |
| `/api/v1/db/:alias/find-columns` | GET | Search for columns across tables by name |
| `/api/v1/db/:alias/:table` | GET, POST | Query paginated records or insert a new record |
| `/api/v1/db/:alias/:table/:id` | GET, PUT, DELETE | Fetch, update, or delete record by primary key |
| `/api/v1/db/:alias/query` | GET, POST | Direct SQL query and mutation execution |
| `/api/v1/db/:alias/views` | GET, POST, DELETE | Dynamic multi-table join view builder |
| `/api/v1/db/:alias/alter` | POST | Alter table schema (add columns) |
| `/api/v1/db/:alias/migrate` | POST | Execute SQL migration script |
| `/api/v1/db/:alias/clone` | POST | Clone database to new target path and mount dynamically |
| `/swagger` | GET | Interactive Swagger UI (zero external CDN dependencies) |
| `/swagger/doc.json` | GET | Dynamic OpenAPI 3.0 specification JSON |
| `/ui` | GET | Custom full interactive Database Visualizer UI |

## Examples

### Bind a database path to an alias
```bash
gitmap db-rest "./data/app.db" as appdb
```

### Start REST server with fallback ports
```bash
gitmap db-rest appdb 8080,8081,3000
```

### Mount multiple databases
```bash
gitmap db-rest master,pipeline 8080,8081
```

### Mount mixed alias and filesystem path
```bash
gitmap db-rest appdb,./external/store.db 8080
```
