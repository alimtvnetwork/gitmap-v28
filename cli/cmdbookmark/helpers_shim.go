package cmdbookmark

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func openDB() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return store.OpenGlobalDefault()
	}
	_ = db.Migrate()
	if countRegisteredSSHHosts(db) == 0 {
		return resolveGlobalSSHDBFallback(db), nil
	}
	return db, nil
}

func countRegisteredSSHHosts(db *store.DB) int {
	var count int
	row := db.SQL().QueryRow("SELECT count(*) FROM ssh_hosts")
	if scanErr := row.Scan(&count); scanErr != nil {
		return 0
	}
	return count
}
func resolveGlobalSSHDBFallback(localDB *store.DB) *store.DB {
	globalDB, err := store.OpenGlobalDefault()
	if err != nil {
		return localDB
	}
	_ = globalDB.Migrate()
	if countRegisteredSSHHosts(globalDB) > 0 {
		localDB.Close()
		return globalDB
	}
	globalDB.Close()
	return localDB
}
