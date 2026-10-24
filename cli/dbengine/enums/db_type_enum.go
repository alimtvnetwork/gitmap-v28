package enums

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DBType identifies supported database engines.
type DBType string

const (
	DBSQLite     DBType = "sqlite"
	DbPostgreSQL DBType = "postgres"
	DbMySQL      DBType = "mysql"
	DbMariaDB    DBType = "mariadb"
	DbMSSQL      DBType = "mssql"
	DbOracle     DBType = "oracle"
	DbMongoDB    DBType = "mongodb"
)

// Name returns the identifier name of the database type.
func (d DBType) Name() string {
	return string(d)
}

// String returns the string representation of the database type.
func (d DBType) String() string {
	return string(d)
}

// Value returns the raw string value of the database type.
func (d DBType) Value() string {
	return string(d)
}

// dbTypeValidMap provides O(1) map validation for database types.
var dbTypeValidMap = map[DBType]bool{
	DBSQLite:     true,
	DbPostgreSQL: true,
	DbMySQL:      true,
	DbMariaDB:    true,
	DbMSSQL:      true,
	DbOracle:     true,
	DbMongoDB:    true,
}

// IsCompare checks equality against another DBType object.
func (d DBType) IsCompare(target DBType) bool {
	return d == target
}

// IsEnum checks whether this database type exists in the valid map.
func (d DBType) IsEnum() bool {
	return dbTypeValidMap[d]
}

// IsValid reports whether this database type is a recognized enumeration value.
func (d DBType) IsValid() bool {
	return d.IsEnum()
}

// IsSQLite checks whether this database type is SQLite.
func (d DBType) IsSQLite() bool {
	return d == DBSQLite
}

// IsPostgreSQL checks whether this database type is PostgreSQL.
func (d DBType) IsPostgreSQL() bool {
	return d == DbPostgreSQL
}

// IsPostgres checks whether this database type is PostgreSQL.
func (d DBType) IsPostgres() bool {
	return d == DbPostgreSQL
}

// IsMySQL checks whether this database type is MySQL.
func (d DBType) IsMySQL() bool {
	return d == DbMySQL
}

// IsMariaDB checks whether this database type is MariaDB.
func (d DBType) IsMariaDB() bool {
	return d == DbMariaDB
}

// IsMSSQL checks whether this database type is MSSQL.
func (d DBType) IsMSSQL() bool {
	return d == DbMSSQL
}

// IsOracle checks whether this database type is Oracle.
func (d DBType) IsOracle() bool {
	return d == DbOracle
}

// IsMongoDB checks whether this database type is MongoDB.
func (d DBType) IsMongoDB() bool {
	return d == DbMongoDB
}

// ParseDBType parses a raw string into a validated DBType.
func ParseDBType(s string) (DBType, *apperror.AppError) {
	val := DBType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseDBType", "E1000", fmt.Sprintf("invalid db type: %q", s)).WithContext("value", s)
}

// AllDBTypes returns all supported database types.
func AllDBTypes() []DBType {
	return []DBType{
		DBSQLite,
		DbPostgreSQL,
		DbMySQL,
		DbMariaDB,
		DbMSSQL,
		DbOracle,
		DbMongoDB,
	}
}

// MarshalJSON implements json.Marshaler.
func (d DBType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(d))
}

// UnmarshalJSON implements json.Unmarshaler with strict map validation.
func (d *DBType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	target := DBType(s)
	if !dbTypeValidMap[target] {
		return apperror.NewValidationError(fmt.Sprintf("invalid db type: %s", s))
	}

	*d = target

	return nil
}

// ToJSON converts the database type to a JSON string representation, returning an AppError on failure.
func (d DBType) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(string(d))
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize db type to json")
	}

	return string(b), nil
}

// FromJSON parses a database type from a JSON string, returning an AppError on failure.
func (d *DBType) FromJSON(s string) *apperror.AppError {
	var str string
	if err := json.Unmarshal([]byte(s), &str); err != nil {
		return apperror.WrapSimple(err, "deserialize db type from json")
	}

	target := DBType(str)
	if !dbTypeValidMap[target] {
		return apperror.NewValidationError(fmt.Sprintf("invalid db type: %s", str))
	}

	*d = target

	return nil
}
