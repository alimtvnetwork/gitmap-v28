# Subtask 01: Database Architecture Specification & Ent Go Codegen Feasibility

> **Subtask ID:** `Subtask-01`  
> **Parent Plan:** `.ai-memory/plans/database-migration-feasibility-analysis.md`  
> **Target Path:** `.ai-memory/plans/subtasks/database-migration-feasibility-analysis/01-schema-and-codegen-feasibility.md`  
> **Status:** `PENDING`  
> **Assigned Agent:** `Author 01`  
> **Owned Files:**  
> - `02-spec/21-app/database-migration-feasibility-analysis/01-architecture-spec.md`  
> - `.ai-memory/plans/subtasks/database-migration-feasibility-analysis/01-schema-and-codegen-feasibility.md`  

---

## 1. Objectives & Scope

1. **Author the Master Architecture Specification:** Formulate and publish `02-spec/21-app/database-migration-feasibility-analysis/01-architecture-spec.md` detailing GitMap's multi-tier SQLite Split-DB topology, driver constraints, and query mechanics.
2. **Evaluate Schema-as-Code & Code Generation:** Analyze Ent Go's declarative Go DSL, AST code generation pipeline (`ent generate`), and predicate builders against GitMap's existing `cli/dbengine/` repository abstraction.
3. **Resolve Field Naming & Primary Key Collisions:** Document the architectural conflict between Ent's mandatory `.ID` struct field accessor and GitMap's `{TableName}Id` convention (Rule 29: `Id`, `Url`, `Api`), evaluating `StorageKey` overrides and custom template generation.
4. **Define Pure-Go Integration & Concurrency Contracts:** Establish the technical integration pattern for `modernc.org/sqlite` v1.34.5 via `entsql.OpenDB` under strict `CGO_ENABLED=0` constraints, single-connection pooling (`db.SetMaxOpenConns(1)`), and WAL mode.
5. **Quantify Code Bloat & Build Overhead:** Calculate precise projections for generated file count (194-238 files), generated code volume (48k-58k LoC), binary footprint expansion (+11MB to +17MB), and compilation latency.
6. **Analyze Split-DB Impedance:** Document package namespace fragmentation across GitMap's 6+ database tiers and the inability to establish graph edges across physically discrete SQLite files.

---

## 2. Detailed Technical Deliverables

### 2.1 Architecture Specification (`01-architecture-spec.md`)
- Detailed breakdown of all 6 SQLite database tiers (`gitmap.db`, `installation.db`, `repodb/sql.db`, `repodb/pipeline.db`, `heatmap.db`, `agent-task.db`, and specialized Split DBs).
- Full analysis of the `modernc.org/sqlite` v1.34.5 Pure-Go driver constraint and rejection of CGO.
- Detailed examination of current query mechanics (`cli/dbengine/Repository[T, F]`, `ModelScanner[T]`, and manual `rows.Scan` boilerplate).
- In-depth evaluation of Ent's AST generation, primary key collisions, pure-Go wrapping, concurrency deadlocks, code bloat, and Split-DB impedance.
- Definitive architectural verdict: Conditional Reject / Architectural Deferral, recommending lightweight AST-based scanner generation for `cli/dbengine/` instead.

### 2.2 Schema Mapping Contracts

The migration feasibility assessment establishes the following schema mapping contracts between GitMap's PascalCase conventions and Ent Go's DSL:

```
+-----------------------------------+-----------------------------------+-----------------------------------+
| Database / Domain Element         | GitMap Repository Convention      | Ent Go Schema DSL Equivalent      |
+-----------------------------------+-----------------------------------+-----------------------------------+
| Table Name                        | PascalCase, Singular (`RepoFile`) | `func (RepoFile) Type() string`   |
| Primary Key Column                | `{TableName}Id` (`RepoFileId`)    | `field.Int64("id").StorageKey()`  |
| Primary Key Go Field              | `model.RepoFileId` (uint64/int64) | `model.ID` (Hardcoded uppercase)  |
| Foreign Key Column                | Referenced PK (`ProjectId`)       | `field.Int64("project_id")`       |
| Boolean Column (Active)           | `IsActive` (INT 0/1, positive)    | `field.Bool("is_active").Default` |
| Boolean Column (State)            | `HasLicense` (INT 0/1, positive)  | `field.Bool("has_license").Def...`|
| Free-Text Entity Column           | `Description TEXT NULL`           | `field.Text("description").Opt...`|
| Transaction Notes Column          | `Notes TEXT NULL`                 | `field.Text("notes").Optional()`  |
| Timestamps (`CreatedAt`, `*At`)   | `INTEGER` (Epoch seconds, UTC)    | `field.Int64("created_at")`       |
+-----------------------------------+-----------------------------------+-----------------------------------+
```

#### Contract Violations & Workarounds
1. **Primary Key Collision:**
   - *Conflict:* Ent requires `.ID` for struct accessors. GitMap code, JSON serializers, and Casbin policies expect `.ProjectId`.
   - *Workaround:* Custom Ent codegen templates (`templates/fields.tmpl`) or runtime DTO transformation adapters.
2. **Edge Invalidation:**
   - *Conflict:* Ent edges (`edge.To`, `edge.From`) require relational SQL foreign keys and single-database joins.
   - *Result:* Edges between entities in different SQLite files (e.g., `Repo` in `gitmap.db` and `RepoScanLog` in `repodb.db`) are invalid and cannot be declared in Ent schemas.

### 2.3 Pure-Go Driver & Connection Pool Contract

Any future adoption of Ent Go within GitMap must strictly adhere to the following initialization lifecycle:

```go
// Required Pure-Go Driver Bridge Pattern
func OpenEntSQLite(dbPath string) (*ent.Client, error) {
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }
    // Mandatory single-connection serialize to prevent SQLite database lock contention
    db.SetMaxOpenConns(1)
    db.SetMaxIdleConns(1)
    db.SetConnMaxLifetime(0)

    // Execute required SQLite PRAGMAs prior to Ent query execution
    pragmas := []string{
        "PRAGMA busy_timeout = 5000;",
        "PRAGMA journal_mode = WAL;",
        "PRAGMA synchronous = NORMAL;",
        "PRAGMA foreign_keys = ON;",
    }
    for _, p := range pragmas {
        if _, err := db.Exec(p); err != nil {
            _ = db.Close()
            return nil, err
        }
    }

    drv := entsql.OpenDB(dialect.SQLite, db)
    return ent.NewClient(ent.Driver(drv)), nil
}
```

---

## 3. Acceptance Criteria & Invariants

1. **Strict Relative Git Paths:** All file paths cited across both specification and plan documents must be relative forward-slash paths (e.g., `02-spec/21-app/...`, `cli/repodb/...`). Zero absolute paths (`C:\...`, `/home/...`) or `file:///` URIs permitted.
2. **Pure-Go Driver Invariant:** Specification unambiguously mandates `modernc.org/sqlite` v1.34.5 and `CGO_ENABLED=0` cross-compilation across Windows, Linux, and macOS.
3. **Single-Connection SQLite Invariant:** Concurrency architecture guarantees `db.SetMaxOpenConns(1)` with WAL mode and `busy_timeout = 5000` to prevent database locks.
4. **Primary Key Convention Analysis:** Explicitly identifies and details the structural conflict between GitMap's `{TableName}Id` (Rule 29: single capital `Id`) and Ent's mandatory `.ID`.
5. **Split-DB Isolation Invariant:** Explicitly details why Ent edges cannot cross physical SQLite database boundaries without unsupported SQLite file attachments.
6. **Strict Disjoint File Ownership:** Author 01 modifies only its two assigned files:
   - `02-spec/21-app/database-migration-feasibility-analysis/01-architecture-spec.md`
   - `.ai-memory/plans/subtasks/database-migration-feasibility-analysis/01-schema-and-codegen-feasibility.md`

---

## 4. Verification Steps & Commands

1. **Verify Specification Structure & Required Headings:**
   Verify that `02-spec/21-app/database-migration-feasibility-analysis/01-architecture-spec.md` contains all mandatory metadata fields (Version, Updated, AI Confidence, Ambiguity, Keywords, Scoring table, Purpose, Cross-References).

2. **Verify Relative Path Hygiene:**
   Ensure zero occurrences of absolute filesystem paths (`d:\`, `c:\`, `/home/`, `/Users/`) or `file:///` URIs in authored files.

3. **Verify Pure-Go & Driver Citations:**
   Ensure `modernc.org/sqlite` v1.34.5 and `CGO_ENABLED=0` are explicitly validated as non-negotiable architectural constraints.

4. **Verify Cross-References Resolution:**
   Validate that all cross-references in `01-architecture-spec.md` target valid existing files:
   - `02-spec/04-database-conventions/02-naming-conventions.md`
   - `02-spec/04-database-conventions/03-schema-design.md`
   - `02-spec/04-database-conventions/04-orm-and-views.md`
   - `02-spec/04-database-conventions/08-split-db-pattern.md`
   - `02-spec/05-split-db-architecture/02-fundamentals.md`
   - `02-spec/23-app-db/01-index.md`
   - `.ai-memory/plans/subtasks/database-migration-feasibility-analysis/01-schema-and-codegen-feasibility.md`
