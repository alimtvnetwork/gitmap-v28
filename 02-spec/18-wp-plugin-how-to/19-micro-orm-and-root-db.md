# Phase 19 — Micro-ORM and Cross-Plugin Root Database

> **Purpose:** Document the fluent micro-ORM query builder and the cross-plugin Root DB registration pattern for snapshot/backup coordination.

---

## 19.1 Micro-ORM — Fluent Query Builder

The micro-ORM is a lightweight, Idiorm-style fluent query builder for SQLite. It provides method chaining for SELECT, INSERT, UPDATE, and DELETE without requiring a full ORM framework.

### Architecture

```
Database/
├── Orm.php                      ← Shell class (static PDO, forTable factory)
└── Traits/
    ├── OrmWhereTrait.php        ← WHERE clause builders
    ├── OrmQueryTrait.php        ← SELECT, ORDER BY, GROUP BY, LIMIT, findOne/findMany
    └── OrmMutationTrait.php     ← create(), set(), save(), delete()
```

### Shell class — Orm.php

The Orm class holds shared static state and composes all three traits:

```php
class Orm {
    use OrmWhereTrait;
    use OrmQueryTrait;
    use OrmMutationTrait;

    private static ?PDO $pdo = null;
    private string $tableName = '';
    private array $data = [];
    private array $whereClauses = [];
    private array $whereParams = [];
    private array $orderBy = [];
    private ?int $limitValue = null;
    private ?int $offsetValue = null;
    private array $selectColumns = ['*'];
    private array $groupBy = [];
    private bool $isNew = false;
    private string $idColumn = 'Id';
    private static int $paramCounter = 0;

    public static function configure(PDO $pdo): void { self::$pdo = $pdo; }
    public static function getPdo(): ?PDO { return self::$pdo; }
    public static function forTable(string $tableName): self { /* factory */ }
    public static function rawExecute(string $sql, array $params = []): array { /* escape hatch */ }

    private function __construct() {}
}
```

### Key design decisions

| Decision | Rationale |
|----------|-----------|
| Static `$pdo` | Configured once via `Orm::configure($pdo)` during `Database::initDatabase()` — no injection needed |
| `forTable()` factory | Returns a fresh instance per query chain — prevents state leakage between queries |
| Private constructor | Forces use of `forTable()` — no accidental bare instantiation |
| `$paramCounter` static | Generates unique parameter names across all WHERE clauses to prevent collisions |

---

## 19.2 OrmWhereTrait — WHERE Clause Building

All WHERE methods return `$this` for chaining. Parameters are auto-named via `generateParamName()` to prevent collisions in complex queries.

### Methods

| Method | SQL Output |
|--------|-----------|
| `where($col, $val)` | `column = :param` |
| `whereEqual($col, $val)` | Alias for `where()` |
| `whereNotEqual($col, $val)` | `column != :param` |
| `whereGt($col, $val)` | `column > :param` |
| `whereGte($col, $val)` | `column >= :param` |
| `whereLt($col, $val)` | `column < :param` |
| `whereLte($col, $val)` | `column <= :param` |
| `whereLike($col, $val)` | `column LIKE :param` |
| `whereNull($col)` | `column IS NULL` |
| `whereNotNull($col)` | `column IS NOT NULL` |
| `whereIn($col, $values)` | `column IN (:p1, :p2, ...)` |
| `whereNotIn($col, $values)` | `column NOT IN (:p1, :p2, ...)` |
| `whereRaw($clause, $params)` | Raw SQL clause with manual params |

### Parameter generation

```php
private function generateParamName(string $column): string {
    self::$paramCounter++;
    $safeColumn = preg_replace('/[^a-zA-Z0-9_]/', '', $column);
    return ':' . $safeColumn . '_' . self::$paramCounter;
}
```

This ensures parameters like `:Status_1`, `:Status_2` never collide, even with multiple `where()` calls on the same column.

### Edge case — empty `whereIn()`

```php
public function whereIn(string $column, array $values) {
    if (empty($values)) {
        $this->whereClauses[] = '1 = 0';  // Always-false — returns no rows
        return $this;
    }
    // ...
}
```

---

## 19.3 OrmQueryTrait — SELECT Operations

### Column selection

| Method | Purpose |
|--------|---------|
| `select($columns)` | Set columns (array or variadic args) |
| `selectColumn($column)` | Single column shorthand |
| `selectCount($alias)` | `COUNT(*) as {alias}` |

### Ordering, grouping, pagination

| Method | Purpose |
|--------|---------|
| `orderByAsc($col)` | `ORDER BY column ASC` |
| `orderByDesc($col)` | `ORDER BY column DESC` |
| `orderBy($col, $dir)` | Custom direction with validation |
| `groupBy($col)` | `GROUP BY column` |
| `limit($n)` | `LIMIT n` |
| `offset($n)` | `OFFSET n` |

### Query execution

| Method | Returns | Purpose |
|--------|---------|---------|
| `findOne($id)` | `?array` | Single record by primary key |
| `findMany()` | `array` | Execute built query, return all rows |
| `count()` | `int` | Count matching records |

### SQL builder

`buildSelectSql()` assembles the final SQL from accumulated state:

```
SELECT {columns} FROM {table}
  [WHERE {clauses}]
  [GROUP BY {columns}]
  [ORDER BY {columns}]
  [LIMIT {n}]
  [OFFSET {n}]
```

All queries are wrapped in try-catch with `InitHelpers::errorLog()` fallback — the ORM never throws.

---

## 19.4 OrmMutationTrait — Write Operations

### Insert flow

```php
$id = Orm::forTable('Transactions')
    ->create()
    ->set('Action', ActionType::Upload->value)
    ->set('Status', 'Success')
    ->set('CreatedAt', DateHelper::nowUtc())
    ->save();
// Returns: int (last insert ID) or false on failure
```

Internally:
1. `create()` sets `$isNew = true` and clears `$data`
2. `set($col, $val)` accumulates column-value pairs
3. `save()` dispatches to `doInsert()` which builds `INSERT INTO ... VALUES (...)`

### Update flow

```php
$affected = Orm::forTable('Transactions')
    ->where('Id', $id)
    ->set('Status', 'Completed')
    ->save();
// Returns: int (rows affected) or false
```

When `$isNew` is false, `save()` dispatches to `doUpdate()` which builds `UPDATE ... SET ... WHERE ...`. **Requires at least one WHERE clause** — bare updates are blocked (returns 0).

### Delete flow

```php
$deleted = Orm::forTable('Transactions')
    ->where('Status', 'Failed')
    ->whereLt('CreatedAt', $cutoffDate)
    ->delete();
// Returns: int (rows deleted)
```

**Safety:** `delete()` refuses to execute without WHERE clauses (returns 0) — prevents accidental `DELETE FROM table`.

---

## 19.5 Usage Examples

### Paginated query with filters

```php
$results = Orm::forTable('Transactions')
    ->select('Id', 'Action', 'Status', 'CreatedAt')
    ->where('Status', 'Success')
    ->whereLike('Action', '%Upload%')
    ->orderByDesc('CreatedAt')
    ->limit(25)
    ->offset(50)
    ->findMany();
```

### Count with filters

```php
$total = Orm::forTable('Transactions')
    ->where('Status', 'Success')
    ->count();
```

### Raw SQL escape hatch

```php
$rows = Orm::rawExecute(
    "SELECT Action, COUNT(*) as Total FROM Transactions GROUP BY Action ORDER BY Total DESC",
);
```

Use `rawExecute()` only for complex queries that the builder cannot express (aggregations with HAVING, subqueries, JOINs).

---

## 19.6 Integration with Database class

The ORM is configured during `Database::initDatabase()` via the ConnectionTrait:

```php
// Inside DatabaseConnectionTrait::initDatabase()
$this->pdo = InitHelpers::initSqliteConnection($this->dbPath, $this->fileLogger);
Orm::configure($this->pdo);  // ← ORM now shares the same PDO
$this->createTables();        // ← migrations run
```

After this, `Orm::forTable()` is available anywhere in the plugin without passing a PDO instance.

---

## 19.7 Cross-Plugin Root Database (RootDb)

The Root Database (`a-root.db`) is a **cross-plugin coordination database** used during snapshot exports. It serves as a manifest that records which tables were exported, their checksums, dependency graphs, and plugin versions.

### Architecture

```
Database/
├── RootDb.php                        ← Shell class (singleton, create/read)
└── Traits/
    ├── RootDbSchemaTrait.php         ← Schema creation, metadata population, dependency graph
    └── RootDbRegistrationTrait.php   ← Table registration, stats, incrementals, plugin snapshots
```

### Shell class — RootDb.php

```php
class RootDb {
    use RootDbSchemaTrait;
    use RootDbRegistrationTrait;

    private FileLogger $logger;
    private DependencyAnalyzer $analyzer;
    private static ?self $instance = null;

    public static function getInstance(?FileLogger $logger = null, ?DependencyAnalyzer $analyzer = null): self;
    public function create(string $filepath): PDO;   // Creates a-root.db with schema
}
```

**Key:** Unlike the main `Database` class, `RootDb` creates a **separate PDO connection** per snapshot — it does not share the plugin's main PDO. Each `a-root.db` is an independent, self-contained manifest file.

### Schema (RootDbSchemaTrait)

The root database contains 5 tables:

#### SnapshotMeta — Export metadata

| Column | Type | Purpose |
|--------|------|---------|
| Id | INTEGER PK | Always 1 (single row) |
| Title | TEXT | Snapshot title |
| Type | TEXT | `Full` or `Incremental` (from `SnapshotModeType` enum) |
| CreatedAt | TEXT | ISO 8601 UTC timestamp |
| CreatedBy | TEXT | Hostname |
| MysqlVersion | TEXT | Source MySQL version |
| WpVersion | TEXT | Source WordPress version |
| PluginVersion | TEXT | Plugin version at export time |
| TableCount | INTEGER | Total tables exported |
| TotalRows | INTEGER | Total rows across all tables |
| ConfigJson | TEXT | Export configuration as JSON |

#### SnapshotTables — Per-table export records

| Column | Type | Purpose |
|--------|------|---------|
| TableName | TEXT UNIQUE | WordPress table name |
| RowCount | INTEGER | Rows exported |
| SqliteFile | TEXT | Relative path to the table's SQLite file |
| FileSizeBytes | INTEGER | File size |
| ChecksumMd5 | TEXT | MD5 checksum for integrity verification |
| ExportedAt | TEXT | When this table was exported |

#### TableDependencies — Foreign key dependency graph

| Column | Type | Purpose |
|--------|------|---------|
| ParentTable | TEXT | Referenced table |
| ChildTable | TEXT | Table with the foreign key |
| FkColumn | TEXT | Foreign key column name |
| RefColumn | TEXT | Referenced column name |

#### IncrementalBackups — Incremental backup log

| Column | Type | Purpose |
|--------|------|---------|
| SequenceNum | INTEGER | Monotonic sequence number |
| FolderName | TEXT | Backup folder name |
| CreatedAt | TEXT | When created |
| TablesChanged | INTEGER | Number of changed tables |
| TotalNewRows | INTEGER | New rows since last backup |
| RelativePath | TEXT | Path relative to snapshots dir |

#### PluginSnapshots — Plugin ZIP archives

| Column | Type | Purpose |
|--------|------|---------|
| PluginSlug | TEXT | Plugin identifier |
| PluginName | TEXT | Display name |
| PluginVersion | TEXT | Version at snapshot time |
| ZipFile | TEXT | Path to ZIP archive |
| FileSizeBytes | INTEGER | Archive size |
| ChecksumMd5 | TEXT | Archive checksum |

### Registration methods (RootDbRegistrationTrait)

| Method | Purpose |
|--------|---------|
| `registerTable($pdo, $tableName, $rowCount, $sqliteFile, ...)` | Record a table export with checksum |
| `updateStats($pdo, $tableCount, $totalRows)` | Update final counts in SnapshotMeta |
| `registerIncremental($pdo, $info)` | Log an incremental backup entry |
| `registerPluginSnapshot($pdo, $info)` | Record a plugin ZIP export |
| `readMetadata($filepath)` | Read full metadata from an existing `a-root.db` |

### Backward compatibility — Legacy table/column resolution

The RootDb supports reading older snapshots that used `snake_case` naming (before the PascalCase migration):

```php
public function resolveRootDbTableName(PDO $pdo, string $pascalName): string {
    // 1. Check if PascalCase table exists → use it
    // 2. Look up legacy map (SnapshotMeta → snapshot_meta)
    // 3. Check if legacy table exists → use it
    // 4. Fall back to PascalCase name
}
```

Legacy maps are maintained for both table names and column names, ensuring the plugin can read snapshots from any version.

### Metadata population flow

```php
$rootDb = RootDb::getInstance($logger, $analyzer);
$pdo = $rootDb->create($snapshotDir . '/a-root.db');

// 1. Populate metadata (MySQL/WP/plugin versions, title, config)
$rootDb->populateMetadata($pdo, $config);

// 2. Populate dependency graph (foreign key analysis)
$rootDb->populateDependencies($pdo, 'all');

// 3. Register each exported table
foreach ($exportedTables as $table) {
    $rootDb->registerTable($pdo, $table['name'], $table['rows'], $table['file'], ...);
}

// 4. Update final stats
$rootDb->updateStats($pdo, count($exportedTables), $totalRows);
```

---

## 19.8 Key Patterns Summary

| Pattern | Where used | Rule |
|---------|-----------|------|
| Static PDO configuration | `Orm::configure()` | Configure once during Database init; available globally after |
| Factory method | `Orm::forTable()` | Fresh instance per query — prevents state leakage |
| Refuse bare mutations | `doUpdate()`, `delete()` | Require WHERE clause — never allow unguarded UPDATE/DELETE |
| Unique param naming | `generateParamName()` | Static counter prevents parameter collisions |
| Separate PDO per manifest | `RootDb::create()` | Each snapshot gets its own independent `a-root.db` |
| Legacy compatibility | `resolveRootDbTableName()` | Always check PascalCase first, then fall back to snake_case |
| Error swallowing | `findMany()`, `count()`, `save()` | ORM catches exceptions and returns empty/zero/false — never throws |
| Typed result wrappers | `TypedQuery` | Never return raw arrays — wrap in `DbResult<T>` / `DbResultSet<T>` / `DbExecResult` |
| Closure mappers | `TypedQuery::queryOne()` | Caller supplies a `Closure(array): T` — like Go's scanner functions |
| Trait decomposition | `FileCache` | Shell class with `FileCacheScanTrait` + `FileCacheStoreTrait` |

---

## 19.9 TypedQuery — Go-Style Typed Database Results

`TypedQuery` provides a **type-safe alternative to the fluent ORM** for cases where you need explicit SQL control with structured result handling. It mirrors Go's `dbutil.Result[T]` pattern.

### Architecture

```
Database/
├── TypedQuery.php       ← Query executor with closure-based row mapping
├── DbResult.php         ← Single-row result wrapper (value | empty | error)
├── DbResultSet.php      ← Multi-row result wrapper (items | error)
└── DbExecResult.php     ← Mutation result wrapper (affectedRows, lastInsertId | error)
```

### TypedQuery class

Takes a PDO instance (not static like Orm) and provides three methods:

```php
final class TypedQuery {
    public function __construct(private readonly PDO $pdo) {}

    /** @return DbResult<T> */
    public function queryOne(string $sql, array $params, Closure $mapper): DbResult;

    /** @return DbResultSet<T> */
    public function queryMany(string $sql, array $params, Closure $mapper): DbResultSet;

    /** Non-query (INSERT/UPDATE/DELETE) */
    public function exec(string $sql, array $params = []): DbExecResult;
}
```

**Key difference from Orm:** TypedQuery uses raw SQL + closure mappers. The caller controls the SQL; TypedQuery handles error wrapping and row-to-object mapping.

### Result wrappers

#### DbResult\<T\> — Single row

| Method | Returns | Purpose |
|--------|---------|---------|
| `DbResult::of($value)` | `self<T>` | Successful result with mapped value |
| `DbResult::empty()` | `self<T>` | No row found (not an error) |
| `DbResult::error($e)` | `self<T>` | Query or mapping failed |
| `->isDefined()` | `bool` | True when a row was mapped |
| `->isEmpty()` | `bool` | True when no row found |
| `->hasError()` | `bool` | True when query failed |
| `->isSafe()` | `bool` | True when defined AND no error |
| `->value()` | `T\|null` | The mapped value |
| `->getError()` | `?Throwable` | The underlying error |
| `->stackTrace()` | `string` | Captured stack trace |

#### DbResultSet\<T\> — Multiple rows

| Method | Returns | Purpose |
|--------|---------|---------|
| `DbResultSet::of($items)` | `self<T>` | Successful result set |
| `DbResultSet::error($e)` | `self<T>` | Query or mapping failed |
| `->items()` | `array<T>` | The mapped items |
| `->count()` | `int` | Number of items |
| `->hasAny()` | `bool` | True when at least one item |
| `->isEmpty()` | `bool` | True when zero items |
| `->first()` | `DbResult<T>` | First item as a DbResult |
| `->isSafe()` | `bool` | True when no error |

#### DbExecResult — Mutations

| Method | Returns | Purpose |
|--------|---------|---------|
| `DbExecResult::of($affected, $lastId)` | `self` | Successful mutation |
| `DbExecResult::error($e)` | `self` | Mutation failed |
| `->affectedRows()` | `int` | Rows changed |
| `->lastInsertId()` | `int` | Auto-increment ID |
| `->isEmpty()` | `bool` | True when zero rows affected |
| `->isSafe()` | `bool` | True when no error |

### Usage example

```php
$query = new TypedQuery($pdo);

// Single row with mapper closure
$result = $query->queryOne(
    "SELECT Id, Name, Version FROM Plugins WHERE Slug = ?",
    [$slug],
    fn(array $row) => new PluginInfo(
        id: (int) $row['Id'],
        name: $row['Name'],
        version: $row['Version'],
    ),
);

if ($result->isSafe()) {
    $plugin = $result->value();  // PluginInfo instance
}

// Multi-row
$resultSet = $query->queryMany(
    "SELECT * FROM Transactions WHERE Status = ? ORDER BY CreatedAt DESC",
    ['Success'],
    fn(array $row) => TransactionDto::fromRow($row),
);

if ($resultSet->hasAny()) {
    foreach ($resultSet->items() as $tx) { /* ... */ }
}

// Mutation
$execResult = $query->exec(
    "UPDATE Plugins SET Version = ? WHERE Slug = ?",
    [$newVersion, $slug],
);

if ($execResult->isSafe()) {
    $affected = $execResult->affectedRows();
}
```

### When to use TypedQuery vs Orm

| Use case | Tool |
|----------|------|
| Simple CRUD, pagination, filters | `Orm::forTable()` fluent builder |
| Complex SQL (JOINs, subqueries, CTEs) | `TypedQuery` with raw SQL |
| Need typed return objects (not arrays) | `TypedQuery` with closure mappers |
| Quick counts, existence checks | `Orm::forTable()->count()` |
| Mutation with affected-row tracking | `TypedQuery::exec()` → `DbExecResult` |

---

## 19.10 FileCache — SQLite-Backed File Hash Cache

`FileCache` manages MD5-based file hash caching for efficient sync comparisons. It follows the standard **shell class + trait decomposition** pattern.

### Architecture

```
Database/
├── FileCache.php                     ← Shell class (singleton)
└── Traits/
    ├── FileCacheScanTrait.php        ← Directory scanning, manifest building, reconciliation
    └── FileCacheStoreTrait.php       ← Cache CRUD (load, upsert, delete, invalidate)
```

### Shell class

```php
class FileCache {
    use FileCacheScanTrait;
    use FileCacheStoreTrait;

    private static ?FileCache $instance = null;
    private FileLogger $logger;
    private Database $db;

    public static function getInstance(FileLogger $logger, Database $db): static;
}
```

**Dependencies:** Takes both `FileLogger` (for logging) and `Database` (for PDO access and readiness check).

### FileCacheScanTrait — Manifest building

The primary method is `getManifest()`, which builds a complete file manifest for a plugin:

```php
$manifest = $fileCache->getManifest($pluginSlug, $pluginDir, $ignoreRules);
// Returns: [
//   'Files'    => [ ['path' => '...', 'hash' => '...', 'modifiedAt' => '...', 'size' => 123], ... ],
//   'Cached'   => 42,    // Files resolved from cache (fast)
//   'Computed'  => 3,     // Files that needed fresh MD5 computation
//   'Removed'   => 1,     // Stale cache entries pruned
// ]
```

#### Reconciliation flow

1. **Load cached entries** from SQLite (keyed by `RelativePath`)
2. **Scan directory** recursively, respecting `.riseupuploadignore` rules
3. **For each file on disk:**
   - If cache hit (same `ModifiedAt` + `FileSize`) → use cached hash (fast path)
   - If cache miss → compute `md5_file()`, upsert cache entry (slow path)
4. **Prune stale entries** — cached files no longer on disk are deleted

#### Graceful degradation

```php
$isDbUnavailable = ($this->db->isReady() === false);
if ($isDbUnavailable) {
    return $this->fullScan($pluginDir, $ignore);  // No cache, compute all hashes
}
```

If the database isn't ready, `getManifest()` falls back to a full scan without caching — ensuring the feature always works.

### FileCacheStoreTrait — Cache CRUD

| Method | Purpose |
|--------|---------|
| `invalidate($pluginSlug)` | Delete all cache entries for a plugin (returns count) |
| `loadCachedEntries($pluginSlug)` | Load all cached entries as `path → row` map |
| `upsertCacheEntry(...)` | INSERT OR REPLACE a cache entry |
| `deleteCacheEntry($pluginSlug, $path)` | Remove a single stale entry |

**Note:** `invalidate()` and `loadCachedEntries()` use the **Orm fluent builder**, while `upsertCacheEntry()` uses **raw PDO** for the `INSERT OR REPLACE` syntax that the Orm doesn't support. This demonstrates the practical coexistence of both database access patterns.

### Cache table schema

```sql
CREATE TABLE IF NOT EXISTS FileCache (
    Id           INTEGER PRIMARY KEY AUTOINCREMENT,
    PluginSlug   TEXT NOT NULL,
    RelativePath TEXT NOT NULL,
    Md5Hash      TEXT NOT NULL,
    ModifiedAt   TEXT NOT NULL,
    FileSize     INTEGER DEFAULT 0,
    CachedAt     TEXT NOT NULL,
    UNIQUE(PluginSlug, RelativePath)
);
```

### Key patterns demonstrated

| Pattern | Implementation |
|---------|---------------|
| Graceful degradation | Falls back to full scan when DB unavailable |
| Cache invalidation | `invalidate($slug)` clears all entries for a plugin |
| Stale entry pruning | `pruneStaleEntries()` removes files no longer on disk |
| Mixed DB access | Orm for queries, raw PDO for `INSERT OR REPLACE` |
| Semantic guards | `BooleanHelpers::isKeyMissing()` for stale path detection |
| ResponseKeyType enums | All response keys use backed enums — no magic strings |
