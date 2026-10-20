# 04-rest-endpoints-crud-and-ddl: REST Catalog, Schema & Record CRUD

- **Spec ID:** `02-spec/21-app/database-rest-command-implementation/04-rest-endpoints-crud-and-ddl.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** REST API Catalog, Schema Introspection, Column Search, Record CRUD
- **Parent Plan:** `.ai-memory/plans/database-rest-command-implementation.md`
- **Dependencies:** `cli/cmddb`, `cli/store`, `cli/appfault`

---

## 1. Universal Response Envelope Specification

Every JSON response returned by the REST API MUST adhere to the Universal Response Envelope:
```json
{
  "Status": {
    "IsSuccess": true,
    "IsFailed": false,
    "Code": 200,
    "Message": "OK",
    "Timestamp": "2026-10-20T10:15:30Z"
  },
  "Attributes": {
    "RequestedAt": "2026-10-20T10:15:30Z",
    "HasAnyErrors": false,
    "IsSingle": false,
    "IsMultiple": true,
    "IsEmpty": false,
    "TotalRecords": 2,
    "DatabaseAlias": "master",
    "TableName": "Repositories",
    "ExecutionTimeMs": 2
  },
  "Results": []
}
```

---

## 2. Database Catalog Endpoint

### `GET /api/v1/databases`
- **Method**: `GET`
- **Purpose**: Global catalog listing all mounted databases and their health.
- **Response**: List of mounted database objects containing `Alias`, `Path`, `TableCount`, `FileSizeBytes`, and `IsReadOnly`.

```json
{
  "Status": { "IsSuccess": true, "Code": 200, "Message": "OK" },
  "Attributes": { "TotalMounted": 2, "IsMultiDb": true },
  "Results": [
    {
      "Alias": "master",
      "Path": "store/gitmap.db",
      "TableCount": 18,
      "FileSizeBytes": 1228800,
      "IsReadOnly": false
    }
  ]
}
```

---

## 3. Database Schema & Table Listing Endpoints

### 3.1 `GET /api/v1/db/{alias}/tables`
- **Method**: `GET`
- **Purpose**: Returns array of all user tables and views with row counts.
- **Response Sample**:
```json
{
  "Results": [
    {
      "Name": "Repositories",
      "Type": "table",
      "ColumnCount": 12,
      "RowCount": 42
    }
  ]
}
```

### 3.2 `GET /api/v1/db/{alias}/schema`
- **Method**: `GET`
- **Query Parameter**: `table` (optional): Filter schema to specific table.
- **Purpose**: Full DDL and column metadata introspection via `PRAGMA table_info`.
- **Response Sample**:
```json
{
  "Results": [
    {
      "TableName": "Repositories",
      "Ddl": "CREATE TABLE Repositories (RepoId INTEGER PRIMARY KEY AUTOINCREMENT, ...);",
      "Columns": [
        {
          "ColumnId": 0,
          "Name": "RepoId",
          "DataType": "INTEGER",
          "IsNullable": false,
          "DefaultValue": null,
          "IsPrimaryKey": true
        }
      ]
    }
  ]
}
```

---

## 4. Cross-Table Column Search Endpoint

### `GET /api/v1/db/{alias}/find-columns`
- **Method**: `GET`
- **Query Parameters**:
  - `name`: Target column name to find (supports wildcard suffix e.g. `path*` or exact match).
- **Purpose**: Scans all tables in `{alias}` to discover which tables contain the specified column name.
- **Response Sample**:
```json
{
  "Results": [
    {
      "TableName": "Repositories",
      "ColumnName": "LocalPath",
      "DataType": "TEXT",
      "IsPrimaryKey": false
    }
  ]
}
```

---

## 5. Record CRUD Operations

Provides automated REST resources for any SQLite table without requiring manual route code.

### 5.1 `GET /api/v1/db/{alias}/{table}` (List Records)
- **Query Parameters**:
  - `limit` (default: 50, max: 1000): Number of rows to return.
  - `offset` (default: 0): Row offset for pagination.
  - `sort` (optional): Column name to order by.
  - `order` (optional, `asc` or `desc`, default: `asc`).
  - `{columnName}={val}`: Any additional query parameter is treated as an equality filter (`WHERE columnName = ?`).
- **Response**: List of row objects with PascalCase column names.

### 5.2 `GET /api/v1/db/{alias}/{table}/{id}` (Get Single Record)
- **URL Parameter**: `id`: Primary key value.
- **Response**: Single record object under `Results`. Returns `404` if not found.

### 5.3 `POST /api/v1/db/{alias}/{table}` (Create Record)
- **Request Body**: JSON object representing column-value map.
- **Execution**:
  1. Validates table exists and column names match table columns.
  2. Constructs parameterized `INSERT INTO <table> (...) VALUES (...)`.
  3. Executes insert and retrieves last insert ID via `res.LastInsertId()`.
- **Response Status**: `201 Created`

### 5.4 `PUT /api/v1/db/{alias}/{table}/{id}` (Update Record)
- **Request Body**: JSON object of columns to update.
- **Execution**:
  1. Locates primary key column via `PRAGMA table_info`.
  2. Builds parameterized `UPDATE <table> SET col1=?, col2=? WHERE <pkCol>=?`.
  3. Returns updated record object.
- **Response Status**: `200 OK`

### 5.5 `DELETE /api/v1/db/{alias}/{table}/{id}` (Delete Record)
- **Execution**: Executes `DELETE FROM <table> WHERE <pkCol> = ?`.
- **Response Status**: `200 OK` with execution summary:
```json
{
  "Status": { "IsSuccess": true, "Code": 200, "Message": "Record deleted" },
  "Results": {
    "Deleted": true,
    "RowsAffected": 1,
    "Id": 101
  }
}
```
