package cmddb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func createEngineTestDbFile(t *testing.T, name string) string {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, name)
	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		t.Fatalf("create test db: %v", openRes.AppError())
	}
	conn := openRes.Value
	defer conn.Close()

	ddl := `CREATE TABLE User (Id INTEGER PRIMARY KEY AUTOINCREMENT, Name TEXT, Age INTEGER, IsActive INTEGER DEFAULT 1);
		INSERT INTO User (Name, Age, IsActive) VALUES ('Alice', 30, 1), ('Bob', 25, 1);`
	if _, err := conn.Exec(ddl); err != nil {
		t.Fatalf("seed test db: %v", err)
	}

	return dbPath
}

func setupTestMux(t *testing.T, dbPath string) (*http.ServeMux, *MultiDbManager, *MigrationCache) {
	t.Helper()
	mgr := NewMultiDbManager(map[string]string{"testdb": dbPath})
	cachePath := filepath.Join(t.TempDir(), "test_cache.db")
	cache, err := OpenMigrationCache(cachePath)
	if err != nil {
		t.Fatalf("open migration cache: %v", err)
	}
	t.Cleanup(func() {
		_ = cache.Close()
		_ = mgr.Close()
	})
	mux := BuildRestMux(mgr, cache)

	return mux, mgr, cache
}

func TestUniversalResponseEnvelope(t *testing.T) {
	rec := httptest.NewRecorder()
	writeRestEnvelope(rec, http.StatusOK, map[string]any{"Name": "Test"}, nil, nil)

	var env RestEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	if env.Status.Code != 200 || env.Status.IsFailed {
		t.Fatalf("unexpected status: %+v", env.Status)
	}
	results, ok := env.Results.([]any)
	if !ok || len(results) != 1 {
		t.Fatalf("results should be slice with 1 item: %+v", env.Results)
	}
}

func TestCatalogAndSchemaEndpoints(t *testing.T) {
	dbPath := createEngineTestDbFile(t, "catalog.db")
	mux, mgr, cache := setupTestMux(t, dbPath)
	defer mgr.Close()
	defer cache.Close()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/databases", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	recTables := httptest.NewRecorder()
	reqTables := httptest.NewRequest(http.MethodGet, "/api/v1/db/testdb/tables", nil)
	mux.ServeHTTP(recTables, reqTables)
	if recTables.Code != http.StatusOK {
		t.Fatalf("expected 200 for tables, got %d", recTables.Code)
	}
}

func TestCrudRecordLifecycle(t *testing.T) {
	dbPath := createEngineTestDbFile(t, "crud.db")
	mux, mgr, cache := setupTestMux(t, dbPath)
	defer mgr.Close()
	defer cache.Close()

	// Insert
	body := bytes.NewBufferString(`{"Name": "Charlie", "Age": 28, "IsActive": 1}`)
	recIns := httptest.NewRecorder()
	reqIns := httptest.NewRequest(http.MethodPost, "/api/v1/db/testdb/User", body)
	mux.ServeHTTP(recIns, reqIns)
	if recIns.Code != http.StatusCreated {
		t.Fatalf("insert record expected 201, got %d", recIns.Code)
	}

	// Query list
	recList := httptest.NewRecorder()
	reqList := httptest.NewRequest(http.MethodGet, "/api/v1/db/testdb/User?page=1&limit=10", nil)
	mux.ServeHTTP(recList, reqList)
	if recList.Code != http.StatusOK {
		t.Fatalf("query records expected 200, got %d", recList.Code)
	}

	// Update
	updBody := bytes.NewBufferString(`{"Age": 29}`)
	recUpd := httptest.NewRecorder()
	reqUpd := httptest.NewRequest(http.MethodPut, "/api/v1/db/testdb/User/1", updBody)
	mux.ServeHTTP(recUpd, reqUpd)
	if recUpd.Code != http.StatusOK {
		t.Fatalf("update record expected 200, got %d", recUpd.Code)
	}

	// Delete
	recDel := httptest.NewRecorder()
	reqDel := httptest.NewRequest(http.MethodDelete, "/api/v1/db/testdb/User/1", nil)
	mux.ServeHTTP(recDel, reqDel)
	if recDel.Code != http.StatusOK {
		t.Fatalf("delete record expected 200, got %d", recDel.Code)
	}
}

func TestViewsAndQueryEndpoints(t *testing.T) {
	dbPath := createEngineTestDbFile(t, "views.db")
	mux, mgr, cache := setupTestMux(t, dbPath)
	defer mgr.Close()
	defer cache.Close()

	// Direct query GET
	recQ := httptest.NewRecorder()
	reqQ := httptest.NewRequest(http.MethodGet, "/api/v1/db/testdb/query?sql=SELECT+Name+FROM+User+ORDER+BY+Id+ASC", nil)
	mux.ServeHTTP(recQ, reqQ)
	if recQ.Code != http.StatusOK {
		t.Fatalf("query GET expected 200, got %d", recQ.Code)
	}

	// Create view
	vBody := bytes.NewBufferString(`{"Name": "ActiveUsers", "Query": "SELECT * FROM User WHERE IsActive = 1"}`)
	recV := httptest.NewRecorder()
	reqV := httptest.NewRequest(http.MethodPost, "/api/v1/db/testdb/views", vBody)
	mux.ServeHTTP(recV, reqV)
	if recV.Code != http.StatusCreated {
		t.Fatalf("create view expected 201, got %d", recV.Code)
	}

	// List views
	recVList := httptest.NewRecorder()
	reqVList := httptest.NewRequest(http.MethodGet, "/api/v1/db/testdb/views", nil)
	mux.ServeHTTP(recVList, reqVList)
	if recVList.Code != http.StatusOK {
		t.Fatalf("list views expected 200, got %d", recVList.Code)
	}

	// Drop view
	recVDel := httptest.NewRecorder()
	reqVDel := httptest.NewRequest(http.MethodDelete, "/api/v1/db/testdb/views/ActiveUsers", nil)
	mux.ServeHTTP(recVDel, reqVDel)
	if recVDel.Code != http.StatusOK {
		t.Fatalf("drop view expected 200, got %d", recVDel.Code)
	}
}

func TestDdlAndCloneEndpoints(t *testing.T) {
	dbPath := createEngineTestDbFile(t, "ddl.db")
	mux, mgr, cache := setupTestMux(t, dbPath)
	defer mgr.Close()
	defer cache.Close()

	// Alter
	altBody := bytes.NewBufferString(`{"Table": "User", "Action": "ADD COLUMN", "Column": "Email", "Type": "TEXT"}`)
	recAlt := httptest.NewRecorder()
	reqAlt := httptest.NewRequest(http.MethodPost, "/api/v1/db/testdb/alter", altBody)
	mux.ServeHTTP(recAlt, reqAlt)
	if recAlt.Code != http.StatusOK {
		t.Fatalf("alter table expected 200, got %d: %s", recAlt.Code, recAlt.Body.String())
	}

	// Clone
	cloneTarget := filepath.Join(t.TempDir(), "cloned.db")
	cloneBody := bytes.NewBufferString(fmt.Sprintf(`{"TargetAlias": "clonedb", "TargetPath": %q}`, cloneTarget))
	recClone := httptest.NewRecorder()
	reqClone := httptest.NewRequest(http.MethodPost, "/api/v1/db/testdb/clone", cloneBody)
	mux.ServeHTTP(recClone, reqClone)
	if recClone.Code != http.StatusOK {
		t.Fatalf("clone db expected 200, got %d: %s", recClone.Code, recClone.Body.String())
	}
	if _, err := os.Stat(cloneTarget); err != nil {
		t.Fatalf("cloned file should exist: %v", err)
	}
}

func TestStartMultiDbRestServer(t *testing.T) {
	dbPath := createEngineTestDbFile(t, "server.db")
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()

	srv, srvErr := StartMultiDbRestServer(map[string]string{"srvdb": dbPath}, l)
	if srvErr != nil {
		t.Fatalf("start rest server: %v", srvErr)
	}
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	resp, getErr := http.Get(fmt.Sprintf("http://%s/api/v1/databases", l.Addr().String()))
	if getErr != nil {
		t.Fatalf("http get: %v", getErr)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
