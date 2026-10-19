package cmddb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	_ "modernc.org/sqlite"
)

func createTestDbForRest(t *testing.T, name string) (string, *sql.DB) {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, name+".db")
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	conn.SetMaxOpenConns(1)
	seedItemsTable(t, conn)
	return dbPath, conn
}

func seedItemsTable(t *testing.T, conn *sql.DB) {
	t.Helper()
	ddl := `CREATE TABLE Items (
		ItemId INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT NOT NULL,
		Price REAL DEFAULT 0.0,
		Status TEXT NOT NULL
	);
	INSERT INTO Items (Name, Price, Status) VALUES ('Widget A', 10.5, 'ACTIVE'), ('Widget B', 25.0, 'ACTIVE');`
	if _, err := conn.Exec(ddl); err != nil {
		t.Fatalf("failed to seed items table: %v", err)
	}
}

func setupRestTestServer(t *testing.T) (*MultiDbManager, *http.ServeMux, string) {
	t.Helper()
	path1, conn1 := createTestDbForRest(t, "alpha")
	conn1.Close()
	path2, conn2 := createTestDbForRest(t, "beta")
	conn2.Close()

	mgr := NewMultiDbManager(map[string]string{"alpha": path1, "beta": path2})
	cache, _ := OpenMigrationCache()
	mux := BuildRestMux(mgr, cache)

	return mgr, mux, path1
}

func execHttpReq(mux http.Handler, method, target, body string) *httptest.ResponseRecorder {
	var rdr *strings.Reader
	if body != "" {
		rdr = strings.NewReader(body)
	} else {
		rdr = strings.NewReader("")
	}
	req := httptest.NewRequest(method, target, rdr)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestPortFallback_SequentialProbing(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to bind listener: %v", err)
	}
	defer ln.Close()
	occupiedPort := ln.Addr().(*net.TCPAddr).Port

	lnNext, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	freePort := lnNext.Addr().(*net.TCPAddr).Port
	_ = lnNext.Close()

	picked, err := PickAvailablePort([]int{occupiedPort, freePort})
	if err != nil {
		t.Fatalf("expected port picker to succeed, got %v", err)
	}
	if picked != freePort {
		t.Fatalf("expected picked port %d, got %d", freePort, picked)
	}
}

func TestPortFallback_AllPortsOccupied(t *testing.T) {
	ln1, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln1.Close()
	p1 := ln1.Addr().(*net.TCPAddr).Port

	ln2, _ := net.Listen("tcp", "127.0.0.1:0")
	defer ln2.Close()
	p2 := ln2.Addr().(*net.TCPAddr).Port

	_, err := PickAvailablePort([]int{p1, p2})
	if err == nil {
		t.Fatalf("expected error when all candidate ports are busy, got nil")
	}
}

func TestTargetResolver_MultiAndMixed(t *testing.T) {
	path1, conn1 := createTestDbForRest(t, "db_target1")
	conn1.Close()
	targets, err := ResolveMountTargets(fmt.Sprintf("master,%s", path1))
	if err != nil {
		t.Fatalf("failed to resolve mount targets: %v", err)
	}
	if len(targets) < 2 {
		t.Fatalf("expected at least 2 resolved targets, got %d", len(targets))
	}
}

func TestRestServer_CatalogAndSchema(t *testing.T) {
	mgr, mux, _ := setupRestTestServer(t)
	defer mgr.Close()

	rec := execHttpReq(mux, "GET", "/api/v1/databases", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from catalog, got %d", rec.Code)
	}

	rec = execHttpReq(mux, "GET", "/api/v1/db/alpha/tables", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from tables, got %d", rec.Code)
	}

	rec = execHttpReq(mux, "GET", "/api/v1/db/alpha/find-columns?name=Price", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from find-columns, got %d", rec.Code)
	}
}

func TestRestServer_CRUDLifecycle(t *testing.T) {
	mgr, mux, _ := setupRestTestServer(t)
	defer mgr.Close()

	rec := execHttpReq(mux, "POST", "/api/v1/db/alpha/Items", `{"Name": "Gadget X", "Price": 99.9, "Status": "NEW"}`)
	if rec.Code != http.StatusCreated && rec.Code != http.StatusOK {
		t.Fatalf("expected 201/200 on insert, got %d: %s", rec.Code, rec.Body.String())
	}

	rec = execHttpReq(mux, "GET", "/api/v1/db/alpha/Items", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on list items, got %d", rec.Code)
	}

	rec = execHttpReq(mux, "GET", "/api/v1/db/alpha/Items/1", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on get item 1, got %d", rec.Code)
	}

	assertCrudUpdateAndDelete(t, mux)
}

func assertCrudUpdateAndDelete(t *testing.T, mux http.Handler) {
	t.Helper()
	rec := execHttpReq(mux, "PUT", "/api/v1/db/alpha/Items/1", `{"Price": 15.0}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on update item 1, got %d", rec.Code)
	}

	rec = execHttpReq(mux, "DELETE", "/api/v1/db/alpha/Items/1", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on delete item 1, got %d", rec.Code)
	}

	rec = execHttpReq(mux, "GET", "/api/v1/db/alpha/Items/1", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after deleting item 1, got %d", rec.Code)
	}
}

func TestRestServer_ViewsAndDirectSQL(t *testing.T) {
	mgr, mux, _ := setupRestTestServer(t)
	defer mgr.Close()

	rec := execHttpReq(mux, "POST", "/api/v1/db/alpha/views", `{"ViewName": "v_active_items", "Query": "SELECT * FROM Items WHERE Status = 'ACTIVE'"}`)
	if rec.Code != http.StatusCreated && rec.Code != http.StatusOK {
		t.Fatalf("expected 201/200 on create view, got %d: %s", rec.Code, rec.Body.String())
	}

	rec = execHttpReq(mux, "GET", "/api/v1/db/alpha/query?sql=SELECT+count(*)+FROM+Items", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on GET query, got %d", rec.Code)
	}

	rec = execHttpReq(mux, "POST", "/api/v1/db/alpha/query", `{"Sql": "UPDATE Items SET Price = 30.0 WHERE ItemId = 2"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on POST query mutation, got %d", rec.Code)
	}
}

func TestRestServer_CloneAndAlter(t *testing.T) {
	mgr, mux, _ := setupRestTestServer(t)
	defer mgr.Close()

	rec := execHttpReq(mux, "POST", "/api/v1/db/alpha/alter", `{"Table": "Items", "Column": "Sku", "Type": "TEXT", "Action": "ADD COLUMN"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on alter add column, got %d: %s", rec.Code, rec.Body.String())
	}

	cloneTarget := filepath.Join(t.TempDir(), "cloned_alpha.db")
	rec = execHttpReq(mux, "POST", "/api/v1/db/alpha/clone", fmt.Sprintf(`{"TargetPath": "%s", "TargetAlias": "alpha_clone"}`, filepath.ToSlash(cloneTarget)))
	if rec.Code != http.StatusCreated && rec.Code != http.StatusOK {
		t.Fatalf("expected 201/200 on clone, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestSwaggerDocAndUIDelivery(t *testing.T) {
	mgr, mux, _ := setupRestTestServer(t)
	defer mgr.Close()

	rec := execHttpReq(mux, "GET", "/swagger/doc.json", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /swagger/doc.json, got %d", rec.Code)
	}

	var spec OpenAPISpec
	if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
		t.Fatalf("failed to parse OpenAPI spec: %v", err)
	}
	if spec.OpenAPI != "3.0.3" {
		t.Fatalf("expected OpenAPI 3.0.3, got %s", spec.OpenAPI)
	}

	assertZeroCdnHtmlDelivery(t, mux)
}

func assertZeroCdnHtmlDelivery(t *testing.T, mux http.Handler) {
	t.Helper()
	rec := execHttpReq(mux, "GET", "/swagger", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /swagger, got %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "<script src=\"http") {
		t.Fatalf("Swagger UI must have zero external CDN scripts")
	}

	rec = execHttpReq(mux, "GET", "/ui", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /ui, got %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "<script src=\"http") {
		t.Fatalf("Visualizer UI must have zero external CDN scripts")
	}
}

func TestConcurrency_DeadlockSafety(t *testing.T) {
	mgr, mux, _ := setupRestTestServer(t)
	defer mgr.Close()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				execHttpReq(mux, "GET", "/api/v1/db/alpha/Items", "")
			} else {
				execHttpReq(mux, "POST", "/api/v1/db/alpha/Items", fmt.Sprintf(`{"Name": "Item %d", "Price": 5.0, "Status": "OK"}`, idx))
			}
		}(i)
	}
	wg.Wait()
}
