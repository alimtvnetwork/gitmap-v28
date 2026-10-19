package cmddb

import (
	"net/http"
	"os"
)

// DatabaseCatalogItem represents catalog listing details for an alias.
type DatabaseCatalogItem struct {
	Alias         string `json:"Alias"`
	Path          string `json:"Path"`
	FileSizeBytes int64  `json:"FileSizeBytes"`
	TableCount    int    `json:"TableCount"`
	IsActive      bool   `json:"IsActive"`
}

// handleCatalog returns an HTTP handler for GET /api/v1/databases.
func handleCatalog(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := buildCatalogItems(mgr)
		writeRestEnvelope(w, http.StatusOK, items, nil, nil)
	}
}

func buildCatalogItems(mgr *MultiDbManager) []DatabaseCatalogItem {
	aliases := mgr.ListAliases()
	items := make([]DatabaseCatalogItem, 0, len(aliases))
	for _, alias := range aliases {
		items = append(items, buildSingleCatalogItem(mgr, alias))
	}

	return items
}

func buildSingleCatalogItem(mgr *MultiDbManager, alias string) DatabaseCatalogItem {
	path, _ := mgr.GetPath(alias)
	size := getFileSizeSafe(path)
	tableCount := countTablesSafe(mgr, alias)

	return DatabaseCatalogItem{
		Alias:         alias,
		Path:          path,
		FileSizeBytes: size,
		TableCount:    tableCount,
		IsActive:      true,
	}
}

func getFileSizeSafe(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return info.Size()
}

func countTablesSafe(mgr *MultiDbManager, alias string) int {
	tables, err := mgr.GetTables(alias)
	if err != nil {
		return 0
	}

	return len(tables)
}
