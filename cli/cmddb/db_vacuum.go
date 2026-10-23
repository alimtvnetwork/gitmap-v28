package cmddb

// runDBVacuum performs SQLite VACUUM and PRAGMA optimize across databases.
func runDBVacuum(args []string) error {
	return runDBOptimize(args)
}
