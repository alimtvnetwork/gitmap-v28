package cmdrm

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"fmt"
	"os"
)

// PrintRepoSuggestions queries the database for suggestions and prints them.
func PrintRepoSuggestions(db *store.DB, target string) {
	if db == nil {
		return
	}

	suggs, _ := db.GetRepoSuggestions(target)
	if len(suggs) == 0 {
		return
	}

	fmt.Fprintf(os.Stderr, "Did you mean:\n")
	for _, s := range suggs {
		fmt.Fprintf(os.Stderr, "  %s\n", s)
	}
}

// ResolveMultiRepos resolves a slice of targets (including globs) across DB and JSON pools.
func ResolveMultiRepos(db *store.DB, targets []string) ([]model.ScanRecord, []string) {
	if ResolveMultiReposFn != nil {
		return ResolveMultiReposFn(db, targets)
	}

	return nil, targets
}

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
