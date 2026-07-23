package runlog

import (
	"database/sql"
	"fmt"
)

// lookupEnumID resolves an enum-mirror Name → Id. Centralized so every
// caller uses the same query shape (and so the test suite can mock it
// out via an in-memory DB seeded with the same enum rows).
func lookupEnumID(db *sql.DB, table, idCol, name string) (int64, error) {
	// #nosec G201 -- table/idCol are trusted package-internal constants, not user input.
	q := fmt.Sprintf("SELECT %s FROM %s WHERE Name = ?", idCol, table)
	var id int64
	if err := db.QueryRow(q, name).Scan(&id); err != nil {
		return 0, fmt.Errorf("lookup %s.%s where Name=%q: %w", table, idCol, name, err)
	}
	return id, nil
}
