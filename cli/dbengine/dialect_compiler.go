package dbengine

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DialectCompiler translates generic query definitions into database-specific SQL.
type DialectCompiler interface {
	Dialect() DatabaseDialectType
	Placeholder(paramIndex int) string
	QuoteIdentifier(name string) string
	CompilePagination(limit, offset int) string
	CompileSearch(table string, fields []string, limit int) string
	CompileCreateView(name string, selectSql string) string
	CompileDropView(name string) string
	CompileFunctionCall(name string, argCount int) string
	CompileCount(table string, field string) string
	CompileDelete(table string, field string) string
	CompileInspectColumns(tableOrView string) string
	CompileInspectViewExists(viewName string) string
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
		return nil, apperror.NewValidationError(fmt.Sprintf("unsupported dialect: %s", dialect))
	}
}
