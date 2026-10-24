package dbengine

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine/enums"
)

// Variant represents any dynamic value acceptable for enum comparisons.
type Variant = any

// DBType is the canonical database engine enum, defined in the dbengine enums package.
type DBType = enums.DBType

// DbType is retained for backward compatibility; new code uses DBType.
type DbType = enums.DBType

// DatabaseDialectType is an alias to DBType for backward compatibility.
type DatabaseDialectType = enums.DBType

const (
	DBSQLite     = enums.DBSQLite
	DbPostgreSQL = enums.DbPostgreSQL
	DbMySQL      = enums.DbMySQL
	DbMariaDB    = enums.DbMariaDB
	DbMSSQL      = enums.DbMSSQL
	DbOracle     = enums.DbOracle
	DbMongoDB    = enums.DbMongoDB
)

const (
	DatabaseDialectSQLite     = enums.DBSQLite
	DatabaseDialectPostgreSQL = enums.DbPostgreSQL
	DatabaseDialectMySQL      = enums.DbMySQL
	DatabaseDialectMariaDB    = enums.DbMariaDB
	DatabaseDialectMSSQL      = enums.DbMSSQL
	DatabaseDialectOracle     = enums.DbOracle
	DatabaseDialectMongoDB    = enums.DbMongoDB
)

// DbTypes provides scoped access to database type enums without repeating prefixes.
var DbTypes = enums.DBTypes
