package dbengine

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

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

// IsCompare checks equality against another DbType, string, or fmt.Stringer.
func (d DbType) IsCompare(target any) bool {
	switch v := target.(type) {
	case DbType:
		return d == v
	case string:
		return string(d) == v
	case fmt.Stringer:
		return string(d) == v.String()
	default:
		return false
	}
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
