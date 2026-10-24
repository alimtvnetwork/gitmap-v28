package constants

import (
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store/enums"
)

// DriverSQLite defines the canonical driver name for SQLite databases.
const DriverSQLite = "sqlite"

var defaultPragmaList = enums.DefaultPragmas()

// ApplyDefaultPragmas configures standard SQLite PRAGMAs on the provided connection.
func ApplyDefaultPragmas(conn *sql.DB) error {
	for _, pragma := range defaultPragmaList {
		if _, err := conn.Exec(pragma.String()); err != nil {
			return apperror.WrapSimple(err, fmt.Sprintf("apply sqlite pragma %q", pragma))
		}
	}

	return nil
}
