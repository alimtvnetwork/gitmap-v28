package cmddb

import (
	"fmt"
	"net"
	"net/http"
)

// StartMultiDbRestServer starts an HTTP REST server on the provided listener.
func StartMultiDbRestServer(dbs map[string]string, listener net.Listener) (*http.Server, error) {
	if listener == nil {
		return nil, fmt.Errorf("listener cannot be nil")
	}

	mgr := NewMultiDbManager(dbs)
	cache, _ := OpenMigrationCache()
	mux := BuildRestMux(mgr, cache)

	server := &http.Server{
		Handler: mux,
	}
	server.RegisterOnShutdown(func() {
		_ = mgr.Close()
		if cache != nil {
			_ = cache.Close()
		}
	})

	go func() {
		_ = server.Serve(listener)
	}()

	return server, nil
}

// BuildRestMux configures and registers all REST API endpoints.
func BuildRestMux(mgr *MultiDbManager, cache *MigrationCache) *http.ServeMux {
	mux := http.NewServeMux()
	registerCatalogRoutes(mux, mgr)
	registerSchemaRoutes(mux, mgr)
	registerViewRoutes(mux, mgr, cache)
	registerDdlRoutes(mux, mgr, cache)
	registerQueryRoutes(mux, mgr)
	registerCrudRoutes(mux, mgr)
	RegisterSwaggerRoutes(mux, mgr)
	RegisterVisualizerRoutes(mux, mgr)

	return mux
}

func registerCatalogRoutes(mux *http.ServeMux, mgr *MultiDbManager) {
	mux.HandleFunc("GET /api/v1/databases", handleCatalog(mgr))
}

func registerSchemaRoutes(mux *http.ServeMux, mgr *MultiDbManager) {
	mux.HandleFunc("GET /api/v1/db/{alias}/tables", handleTables(mgr))
	mux.HandleFunc("GET /api/v1/db/{alias}/schema", handleSchema(mgr))
	mux.HandleFunc("GET /api/v1/db/{alias}/find-columns", handleFindColumns(mgr))
}

func registerViewRoutes(mux *http.ServeMux, mgr *MultiDbManager, cache *MigrationCache) {
	mux.HandleFunc("GET /api/v1/db/{alias}/views", handleListViews(mgr))
	mux.HandleFunc("POST /api/v1/db/{alias}/views", handleCreateView(mgr, cache))
	mux.HandleFunc("DELETE /api/v1/db/{alias}/views/{name}", handleDeleteView(mgr))
}

func registerDdlRoutes(mux *http.ServeMux, mgr *MultiDbManager, cache *MigrationCache) {
	mux.HandleFunc("POST /api/v1/db/{alias}/tables/create", handleCreateTable(mgr, cache))
	mux.HandleFunc("POST /api/v1/db/{alias}/alter", handleAlter(mgr, cache))
	mux.HandleFunc("POST /api/v1/db/{alias}/migrate", handleMigrate(mgr, cache))
	mux.HandleFunc("POST /api/v1/db/{alias}/clone", handleClone(mgr, cache))
}

func registerQueryRoutes(mux *http.ServeMux, mgr *MultiDbManager) {
	mux.HandleFunc("GET /api/v1/db/{alias}/query", handleQuery(mgr))
	mux.HandleFunc("POST /api/v1/db/{alias}/query", handleQuery(mgr))
}

func registerCrudRoutes(mux *http.ServeMux, mgr *MultiDbManager) {
	mux.HandleFunc("GET /api/v1/db/{alias}/{table}", handleGetRecords(mgr))
	mux.HandleFunc("POST /api/v1/db/{alias}/{table}", handleInsertRecord(mgr))
	mux.HandleFunc("GET /api/v1/db/{alias}/{table}/{id}", handleGetRecordById(mgr))
	mux.HandleFunc("PUT /api/v1/db/{alias}/{table}/{id}", handleUpdateRecord(mgr))
	mux.HandleFunc("PATCH /api/v1/db/{alias}/{table}/{id}", handleUpdateRecord(mgr))
	mux.HandleFunc("DELETE /api/v1/db/{alias}/{table}/{id}", handleDeleteRecord(mgr))
}
