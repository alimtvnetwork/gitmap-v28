package cmdaudit

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return true
		}
	}

	return false
}

// openTasksDB opens the split tasks root database at .gitmap/data/tasks/sql.db.
func openTasksDB() (*store.DB, error) {
	tasksDB, err := store.OpenTasksRootSplitDB()
	if err != nil {
		return nil, err
	}

	return tasksDB.DB, nil
}
