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

// IsEnum checks equality against another DbType, string, or fmt.Stringer (alias for IsCompare).
func (d DbType) IsEnum(target any) bool {
	return d.IsCompare(target)
}

// MarshalJSON implements json.Marshaler.
func (d DbType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(d))
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DbType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*d = DbType(s)
	return nil
}

// ToJSON converts the database type to a JSON string representation.
func (d DbType) ToJSON() (string, error) {
	b, err := json.Marshal(string(d))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromJSON parses a database type from a JSON string.
func (d *DbType) FromJSON(s string) error {
	return json.Unmarshal([]byte(s), d)
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

// IsEnum checks whether the given variant matches any registered database type.
func (r dbTypeRegistry) IsEnum(variant Variant) bool {
	for _, d := range r.All() {
		if d.IsCompare(variant) {
			return true
		}
	}
	return false
}

// IsSQLite checks whether the variant represents SQLite.
func (r dbTypeRegistry) IsSQLite(variant Variant) bool {
	return r.SQLite.IsCompare(variant)
}

// IsPostgreSQL checks whether the variant represents PostgreSQL.
func (r dbTypeRegistry) IsPostgreSQL(variant Variant) bool {
	return r.PostgreSQL.IsCompare(variant)
}

// IsPostgres checks whether the variant represents PostgreSQL.
func (r dbTypeRegistry) IsPostgres(variant Variant) bool {
	return r.PostgreSQL.IsCompare(variant)
}

// IsMySQL checks whether the variant represents MySQL.
func (r dbTypeRegistry) IsMySQL(variant Variant) bool {
	return r.MySQL.IsCompare(variant)
}

// IsMariaDB checks whether the variant represents MariaDB.
func (r dbTypeRegistry) IsMariaDB(variant Variant) bool {
	return r.MariaDB.IsCompare(variant)
}

// IsMSSQL checks whether the variant represents MSSQL.
func (r dbTypeRegistry) IsMSSQL(variant Variant) bool {
	return r.MSSQL.IsCompare(variant)
}

// IsOracle checks whether the variant represents Oracle.
func (r dbTypeRegistry) IsOracle(variant Variant) bool {
	return r.Oracle.IsCompare(variant)
}

// IsMongoDB checks whether the variant represents MongoDB.
func (r dbTypeRegistry) IsMongoDB(variant Variant) bool {
	return r.MongoDB.IsCompare(variant)
}

// ToJSON converts the database type registry to a JSON string representation.
func (r dbTypeRegistry) ToJSON() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
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
