package dbengine

import (
	"encoding/json"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// Variant represents any dynamic value acceptable for enum comparisons.
type Variant = any

// DbType identifies supported database engines.
type DbType string

const (
	DbSQLite     DbType = "sqlite"
	DbPostgreSQL DbType = "postgres"
	DbMySQL      DbType = "mysql"
	DbMariaDB    DbType = "mariadb"
	DbMSSQL      DbType = "mssql"
	DbOracle     DbType = "oracle"
	DbMongoDB    DbType = "mongodb"
)

// Name returns the identifier name of the database type.
func (d DbType) Name() string {
	return string(d)
}

// String returns the string representation of the database type.
func (d DbType) String() string {
	return string(d)
}

// Value returns the raw string value of the database type.
func (d DbType) Value() string {
	return string(d)
}

// dbTypeValidMap provides O(1) map validation for database types.
var dbTypeValidMap = map[DbType]bool{
	DbSQLite:     true,
	DbPostgreSQL: true,
	DbMySQL:      true,
	DbMariaDB:    true,
	DbMSSQL:      true,
	DbOracle:     true,
	DbMongoDB:    true,
}

// IsCompare checks equality against another DbType object.
func (d DbType) IsCompare(target DbType) bool {
	return d == target
}

// IsEnum checks whether this database type exists in the valid map.
func (d DbType) IsEnum() bool {
	return dbTypeValidMap[d]
}

// IsSQLite checks whether this database type is SQLite.
func (d DbType) IsSQLite() bool {
	return d == DbSQLite
}

// IsPostgreSQL checks whether this database type is PostgreSQL.
func (d DbType) IsPostgreSQL() bool {
	return d == DbPostgreSQL
}

// IsPostgres checks whether this database type is PostgreSQL.
func (d DbType) IsPostgres() bool {
	return d == DbPostgreSQL
}

// IsMySQL checks whether this database type is MySQL.
func (d DbType) IsMySQL() bool {
	return d == DbMySQL
}

// IsMariaDB checks whether this database type is MariaDB.
func (d DbType) IsMariaDB() bool {
	return d == DbMariaDB
}

// IsMSSQL checks whether this database type is MSSQL.
func (d DbType) IsMSSQL() bool {
	return d == DbMSSQL
}

// IsOracle checks whether this database type is Oracle.
func (d DbType) IsOracle() bool {
	return d == DbOracle
}

// IsMongoDB checks whether this database type is MongoDB.
func (d DbType) IsMongoDB() bool {
	return d == DbMongoDB
}

// MarshalJSON implements json.Marshaler.
func (d DbType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(d))
}

// UnmarshalJSON implements json.Unmarshaler with strict map validation.
func (d *DbType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	target := DbType(s)
	if !dbTypeValidMap[target] {
		return fmt.Errorf("invalid db type: %s", s)
	}
	*d = target
	return nil
}

// ToJSON converts the database type to a JSON string representation, returning an AppError on failure.
func (d DbType) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(string(d))
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize db type to json")
	}
	return string(b), nil
}

// FromJSON parses a database type from a JSON string, returning an AppError on failure.
func (d *DbType) FromJSON(s string) *apperror.AppError {
	var str string
	if err := json.Unmarshal([]byte(s), &str); err != nil {
		return apperror.WrapSimple(err, "deserialize db type from json")
	}
	target := DbType(str)
	if !dbTypeValidMap[target] {
		return apperror.WrapSimple(fmt.Errorf("invalid db type: %s", str), "validate db type from json")
	}
	*d = target
	return nil
}

// dbTypeRegistry provides scoped access to supported database types.
type dbTypeRegistry struct {
	SQLite     DbType
	PostgreSQL DbType
	Postgres   DbType
	MySQL      DbType
	MariaDB    DbType
	MSSQL      DbType
	Oracle     DbType
	MongoDB    DbType
}

// All returns a slice of all supported database types.
func (r dbTypeRegistry) All() []DbType {
	return []DbType{
		r.SQLite,
		r.PostgreSQL,
		r.MySQL,
		r.MariaDB,
		r.MSSQL,
		r.Oracle,
		r.MongoDB,
	}
}

// Names returns a slice of string names for all supported database types.
func (r dbTypeRegistry) Names() []string {
	return []string{
		"sqlite",
		"postgres",
		"mysql",
		"mariadb",
		"mssql",
		"oracle",
		"mongodb",
	}
}

// IsEnum checks whether the target object matches any registered database type.
func (r dbTypeRegistry) IsEnum(target DbType) bool {
	return dbTypeValidMap[target]
}

// IsSQLite checks whether the target matches SQLite.
func (r dbTypeRegistry) IsSQLite(target DbType) bool {
	return target == r.SQLite
}

// IsPostgreSQL checks whether the target matches PostgreSQL.
func (r dbTypeRegistry) IsPostgreSQL(target DbType) bool {
	return target == r.PostgreSQL
}

// IsPostgres checks whether the target matches PostgreSQL.
func (r dbTypeRegistry) IsPostgres(target DbType) bool {
	return target == r.Postgres
}

// IsMySQL checks whether the target matches MySQL.
func (r dbTypeRegistry) IsMySQL(target DbType) bool {
	return target == r.MySQL
}

// IsMariaDB checks whether the target matches MariaDB.
func (r dbTypeRegistry) IsMariaDB(target DbType) bool {
	return target == r.MariaDB
}

// IsMSSQL checks whether the target matches MSSQL.
func (r dbTypeRegistry) IsMSSQL(target DbType) bool {
	return target == r.MSSQL
}

// IsOracle checks whether the target matches Oracle.
func (r dbTypeRegistry) IsOracle(target DbType) bool {
	return target == r.Oracle
}

// IsMongoDB checks whether the target matches MongoDB.
func (r dbTypeRegistry) IsMongoDB(target DbType) bool {
	return target == r.MongoDB
}

// ToJSON converts the database type registry to a JSON string representation, returning an AppError on failure.
func (r dbTypeRegistry) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize db type registry to json")
	}
	return string(b), nil
}

// DbTypes provides scoped access to database type enums without repeating prefixes.
var DbTypes = dbTypeRegistry{
	SQLite:     DbSQLite,
	PostgreSQL: DbPostgreSQL,
	Postgres:   DbPostgreSQL,
	MySQL:      DbMySQL,
	MariaDB:    DbMariaDB,
	MSSQL:      DbMSSQL,
	Oracle:     DbOracle,
	MongoDB:    DbMongoDB,
}

// DatabaseDialectType is an alias to DbType for backward compatibility.
type DatabaseDialectType = DbType

const (
	DatabaseDialectSQLite     = DbSQLite
	DatabaseDialectPostgreSQL = DbPostgreSQL
	DatabaseDialectMySQL      = DbMySQL
	DatabaseDialectMariaDB    = DbMariaDB
	DatabaseDialectMSSQL      = DbMSSQL
	DatabaseDialectOracle     = DbOracle
	DatabaseDialectMongoDB    = DbMongoDB
)

// DialectCompiler translates generic query definitions into database-specific SQL.
type DialectCompiler interface {
	Dialect() DatabaseDialectType
	Placeholder(paramIndex int) string
	QuoteIdentifier(name string) string
	CompilePagination(limit, offset int) string
	CompileSearch(table string, fields []string, limit int) string
}

// ResolveCompiler returns the appropriate DialectCompiler for a dialect enum.
func ResolveCompiler(dialect DatabaseDialectType) (DialectCompiler, *apperror.AppError) {
	switch dialect {
	case DatabaseDialectSQLite:
		return &SQLiteCompiler{}, nil
	case DatabaseDialectPostgreSQL:
		return &PostgresCompiler{}, nil
	case DatabaseDialectMySQL, DatabaseDialectMariaDB:
		return &MySQLCompiler{dialect: dialect}, nil
	case DatabaseDialectMSSQL:
		return &MSSQLCompiler{}, nil
	case DatabaseDialectOracle:
		return &OracleCompiler{}, nil
	case DatabaseDialectMongoDB:
		return &MongoDBCompiler{}, nil
	default:
		return nil, apperror.WrapSimple(fmt.Errorf("unsupported dialect: %s", dialect), "resolve dialect compiler")
	}
}
