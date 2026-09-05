package dbengine

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// DatabaseDialectType identifies supported database engines.
type DatabaseDialectType string

const (
	DatabaseDialectSQLite     DatabaseDialectType = "sqlite"
	DatabaseDialectPostgreSQL DatabaseDialectType = "postgres"
	DatabaseDialectMySQL      DatabaseDialectType = "mysql"
	DatabaseDialectMariaDB    DatabaseDialectType = "mariadb"
	DatabaseDialectMSSQL      DatabaseDialectType = "mssql"
	DatabaseDialectOracle     DatabaseDialectType = "oracle"
	DatabaseDialectMongoDB    DatabaseDialectType = "mongodb"
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
