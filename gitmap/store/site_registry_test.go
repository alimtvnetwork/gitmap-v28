package store

import (
	"path/filepath"
	"testing"
)

func TestSitesSplitDBCRUD(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "sites.db")

	db, err := OpenSitesSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenSitesSplitDBAt failed: %v", err)
	}
	defer db.Close()

	rec := SiteRecord{
		Domain:          "example.com",
		SiteType:        "wordpress",
		DocumentRoot:    "/var/www/example.com",
		NginxConfigPath: "/etc/nginx/sites-available/example.com",
		ListenPort:      80,
		IsSslEnabled:    false,
		IsActive:        true,
		Description:     "Primary WordPress site",
	}

	if err := db.UpsertSite(rec); err != nil {
		t.Fatalf("UpsertSite failed: %v", err)
	}

	got, err := db.GetSite("example.com")
	if err != nil {
		t.Fatalf("GetSite failed: %v", err)
	}
	if got == nil || got.Domain != "example.com" || got.SiteType != "wordpress" {
		t.Fatalf("unexpected site retrieved: %+v", got)
	}

	subRec := SiteRecord{
		Domain:          "sub.example.com",
		SiteType:        "laravel",
		DocumentRoot:    "/var/www/sub.example.com",
		NginxConfigPath: "/etc/nginx/sites-available/sub.example.com",
		ListenPort:      443,
		IsSslEnabled:    true,
		IsActive:        true,
	}
	if err := db.UpsertSite(subRec); err != nil {
		t.Fatalf("UpsertSite sub failed: %v", err)
	}

	all, err := db.ListSites()
	if err != nil {
		t.Fatalf("ListSites failed: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 sites, got %d", len(all))
	}

	if err := db.DeleteSite("example.com"); err != nil {
		t.Fatalf("DeleteSite failed: %v", err)
	}

	afterDel, err := db.GetSite("example.com")
	if err != nil {
		t.Fatalf("GetSite after del failed: %v", err)
	}
	if afterDel != nil {
		t.Fatalf("expected nil after delete, got %+v", afterDel)
	}
}
