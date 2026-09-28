package osclean

// EnhancedCategoryDef defines a developer tool cache category for scanning.
type EnhancedCategoryDef struct {
	ID        string
	Label     string
	Aliases   []string
	CleanFunc func(bool) CategoryCleanStats
}

// CleanTableEntry stores formatted column data for aligned table rendering.
type CleanTableEntry struct {
	Category string
	Status   string
	Files    string
	Dirs     string
	Size     string
}
