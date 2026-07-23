// Package store_test — unit tests for CommitTemplates CRUD operations.
package store_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// openTestDB creates an in-memory-style temp DB for testing.
func openTestDB(t *testing.T) *store.DB {
	t.Helper()
	dir := t.TempDir()
	db, err := store.Open(dir)
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}

	if err := db.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

// TestInsertTemplate_Title verifies inserting a title template.
func TestInsertTemplate_Title(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	err := db.InsertTemplate(constants.TemplateKindTitle, "Top {service} in {area}")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestInsertTemplate_Description verifies inserting a description template.
func TestInsertTemplate_Description(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	err := db.InsertTemplate(constants.TemplateKindDescription, "Best {service} provider in {area}")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestListTemplatesByKind_ReturnsCorrectKind verifies filtering by kind.
func TestListTemplatesByKind_ReturnsCorrectKind(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	db.InsertTemplate(constants.TemplateKindTitle, "Title A")
	db.InsertTemplate(constants.TemplateKindTitle, "Title B")
	db.InsertTemplate(constants.TemplateKindDescription, "Desc A")

	titles, err := db.ListTemplatesByKind(constants.TemplateKindTitle)
	if err != nil {
		t.Fatal(err)
	}
	if len(titles) != 2 {
		t.Errorf("expected 2 titles, got %d", len(titles))
	}

	descs, err := db.ListTemplatesByKind(constants.TemplateKindDescription)
	if err != nil {
		t.Fatal(err)
	}
	if len(descs) != 1 {
		t.Errorf("expected 1 description, got %d", len(descs))
	}
}

// TestListTemplatesByKind_EmptyTable returns empty slice.
func TestListTemplatesByKind_EmptyTable(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	titles, err := db.ListTemplatesByKind(constants.TemplateKindTitle)
	if err != nil {
		t.Fatal(err)
	}
	if len(titles) != 0 {
		t.Errorf("expected 0 titles, got %d", len(titles))
	}
}

// TestCountTemplates_Empty verifies zero count on fresh DB.
func TestCountTemplates_Empty(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	count, err := db.CountTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

// TestCountTemplates_AfterInserts verifies count after inserts.
func TestCountTemplates_AfterInserts(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	db.InsertTemplate(constants.TemplateKindTitle, "T1")
	db.InsertTemplate(constants.TemplateKindTitle, "T2")
	db.InsertTemplate(constants.TemplateKindDescription, "D1")

	count, err := db.CountTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("expected 3, got %d", count)
	}
}

// TestInsertTemplate_UniqueIDs verifies each insert gets a unique ID.
func TestInsertTemplate_UniqueIDs(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	db.InsertTemplate(constants.TemplateKindTitle, "T1")
	db.InsertTemplate(constants.TemplateKindTitle, "T2")

	titles, _ := db.ListTemplatesByKind(constants.TemplateKindTitle)
	if len(titles) < 2 {
		t.Fatal("expected at least 2 titles")
	}
	if titles[0].ID == titles[1].ID {
		t.Error("expected unique IDs for each template")
	}
}

// TestListTemplatesByKind_TemplateContent verifies content is stored correctly.
func TestListTemplatesByKind_TemplateContent(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	expected := "Top {service} in {area} - visit {url}"
	db.InsertTemplate(constants.TemplateKindTitle, expected)

	titles, _ := db.ListTemplatesByKind(constants.TemplateKindTitle)
	if len(titles) != 1 {
		t.Fatal("expected 1 title")
	}
	if titles[0].Template != expected {
		t.Errorf("expected %q, got %q", expected, titles[0].Template)
	}
}

// TestListTemplatesByKind_KindField verifies the Kind field is correct.
func TestListTemplatesByKind_KindField(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	db.InsertTemplate(constants.TemplateKindTitle, "Title")
	db.InsertTemplate(constants.TemplateKindDescription, "Desc")

	titles, _ := db.ListTemplatesByKind(constants.TemplateKindTitle)
	if len(titles) == 0 {
		t.Fatal("expected titles")
	}
	if titles[0].Kind != constants.TemplateKindTitle {
		t.Errorf("expected kind %q, got %q", constants.TemplateKindTitle, titles[0].Kind)
	}

	descs, _ := db.ListTemplatesByKind(constants.TemplateKindDescription)
	if len(descs) == 0 {
		t.Fatal("expected descriptions")
	}
	if descs[0].Kind != constants.TemplateKindDescription {
		t.Errorf("expected kind %q, got %q", constants.TemplateKindDescription, descs[0].Kind)
	}
}

// TestInsertTemplate_CreatedAtPopulated verifies CreatedAt is set.
func TestInsertTemplate_CreatedAtPopulated(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	db.InsertTemplate(constants.TemplateKindTitle, "Test")

	titles, _ := db.ListTemplatesByKind(constants.TemplateKindTitle)
	if len(titles) == 0 {
		t.Fatal("expected titles")
	}
	if titles[0].CreatedAt == "" {
		t.Error("expected CreatedAt to be populated")
	}
}
