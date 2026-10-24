package enums

import (
	"encoding/json"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DBTypeRegistry provides scoped access to supported database types.
type DBTypeRegistry struct {
	SQLite     DBType
	PostgreSQL DBType
	Postgres   DBType
	MySQL      DBType
	MariaDB    DBType
	MSSQL      DBType
	Oracle     DBType
	MongoDB    DBType
}

// All returns a slice of all supported database types.
func (r DBTypeRegistry) All() []DBType {
	return []DBType{
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
func (r DBTypeRegistry) Names() []string {
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
func (r DBTypeRegistry) IsEnum(target DBType) bool {
	return dbTypeValidMap[target]
}

// IsSQLite checks whether the target matches SQLite.
func (r DBTypeRegistry) IsSQLite(target DBType) bool {
	return target == r.SQLite
}

// IsPostgreSQL checks whether the target matches PostgreSQL.
func (r DBTypeRegistry) IsPostgreSQL(target DBType) bool {
	return target == r.PostgreSQL
}

// IsPostgres checks whether the target matches PostgreSQL.
func (r DBTypeRegistry) IsPostgres(target DBType) bool {
	return target == r.Postgres
}

// IsMySQL checks whether the target matches MySQL.
func (r DBTypeRegistry) IsMySQL(target DBType) bool {
	return target == r.MySQL
}

// IsMariaDB checks whether the target matches MariaDB.
func (r DBTypeRegistry) IsMariaDB(target DBType) bool {
	return target == r.MariaDB
}

// IsMSSQL checks whether the target matches MSSQL.
func (r DBTypeRegistry) IsMSSQL(target DBType) bool {
	return target == r.MSSQL
}

// IsOracle checks whether the target matches Oracle.
func (r DBTypeRegistry) IsOracle(target DBType) bool {
	return target == r.Oracle
}

// IsMongoDB checks whether the target matches MongoDB.
func (r DBTypeRegistry) IsMongoDB(target DBType) bool {
	return target == r.MongoDB
}

// ToJSON converts the database type registry to a JSON string representation, returning an AppError on failure.
func (r DBTypeRegistry) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize db type registry to json")
	}

	return string(b), nil
}

// DBTypes provides scoped access to database type enums without repeating prefixes.
var DBTypes = DBTypeRegistry{
	SQLite:     DBSQLite,
	PostgreSQL: DbPostgreSQL,
	Postgres:   DbPostgreSQL,
	MySQL:      DbMySQL,
	MariaDB:    DbMariaDB,
	MSSQL:      DbMSSQL,
	Oracle:     DbOracle,
	MongoDB:    DbMongoDB,
}
