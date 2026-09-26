package store

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestTemplatesSplitDB_SchemaAndCategories(t *testing.T) {
	db := openTestTemplatesDB(t)
	defer db.Close()

	cats, err := db.ListCategories()
	if err != nil {
		t.Fatalf("ListCategories failed: %v", err)
	}
	if len(cats) < 5 {
		t.Fatalf("expected at least 5 default categories, got %d", len(cats))
	}

	verifyDefaultSeededPrompts(t, db)
	verifyDefaultSeededSponsors(t, db)
	verifyCustomCategoryUpsert(t, db)
}

func verifyDefaultSeededPrompts(t *testing.T, db *TemplatesSplitDB) {
	t.Helper()
	item, err := db.GetTemplateItem("tpl-prompt-ui-ux-audit")
	if err != nil || item == nil || item.Category != "prompts" || item.SubCategory != "ui-ux" {
		t.Fatalf("expected default ui-ux prompt template, got %+v (err=%v)", item, err)
	}
}

func verifyDefaultSeededSponsors(t *testing.T, db *TemplatesSplitDB) {
	t.Helper()
	item, err := db.GetTemplateItem("tpl-seo-sponsor-default")
	if err != nil || item == nil || item.Category != "seo" || item.SubCategory != "sponsor" {
		t.Fatalf("expected default seo sponsor template, got %+v (err=%v)", item, err)
	}
}

func openTestTemplatesDB(t *testing.T) *TemplatesSplitDB {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test-templates.db")
	db, err := OpenTemplatesSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenTemplatesSplitDBAt failed: %v", err)
	}

	return db
}

func verifyCustomCategoryUpsert(t *testing.T, db *TemplatesSplitDB) {
	t.Helper()
	err := db.UpsertCategory(StateTemplateCategory{Slug: "custom-cat", Name: "Custom", Description: "Desc"})
	if err != nil {
		t.Fatalf("UpsertCategory failed: %v", err)
	}

	cats, err := db.ListCategories()
	if err != nil || len(cats) < 6 {
		t.Fatalf("expected custom category in list, got %d (err=%v)", len(cats), err)
	}
}

func TestTemplatesSplitDB_UpsertByIDAndSlug(t *testing.T) {
	db := openTestTemplatesDB(t)
	defer db.Close()

	saved := insertInitialTemplateItem(t, db)
	verifyUpdateByID(t, db, saved.ID)
	verifyUpdateBySlugFallback(t, db, saved.ID)
	verifyAutoIDAndDelete(t, db)
}

func insertInitialTemplateItem(t *testing.T, db *TemplatesSplitDB) StateTemplateItem {
	t.Helper()
	saved, err := db.UpsertTemplateItem(StateTemplateItem{
		ID:       "seo-01",
		Category: "seo",
		Slug:     "why-canonical",
		Title:    "# Why use canonical tags?",
		Text:     "Because duplicate URLs dilute ranking signals by 38%.",
	})
	if err != nil {
		t.Fatalf("UpsertTemplateItem failed: %v", err)
	}

	return saved
}

func verifyUpdateByID(t *testing.T, db *TemplatesSplitDB, id string) {
	t.Helper()
	updated, err := db.UpsertTemplateItem(StateTemplateItem{
		ID:       id,
		Category: "seo",
		Slug:     "why-canonical-v2",
		Title:    "# Why use canonical tags v2?",
		Text:     "Because canonical consolidation recovers 42% crawl budget.",
	})
	if err != nil || updated.Slug != "why-canonical-v2" {
		t.Fatalf("update by ID failed: %+v, err=%v", updated, err)
	}
}

func verifyUpdateBySlugFallback(t *testing.T, db *TemplatesSplitDB, expectedID string) {
	t.Helper()
	updated, err := db.UpsertTemplateItem(StateTemplateItem{
		Category: "seo",
		Slug:     "why-canonical-v2",
		Title:    "# Why use canonical tags v3?",
		Text:     "Because canonical URLs prevent split equity.",
	})
	if err != nil || updated.ID != expectedID || !strings.Contains(updated.Title, "v3") {
		t.Fatalf("update by Slug failed: %+v, err=%v", updated, err)
	}
}

func verifyAutoIDAndDelete(t *testing.T, db *TemplatesSplitDB) {
	t.Helper()
	autoItem, err := db.UpsertTemplateItem(StateTemplateItem{
		Category: "prefix",
		Slug:     "fast-prefix",
		Title:    "# Why prefix?",
		Text:     "Because structured commits improve traceability.",
	})
	if err != nil || autoItem.ID != "tpl-fast-prefix" {
		t.Fatalf("expected auto ID tpl-fast-prefix, got %+v (err=%v)", autoItem, err)
	}

	if err := db.DeleteTemplateItem("fast-prefix"); err != nil {
		t.Fatalf("DeleteTemplateItem failed: %v", err)
	}
}

func TestTemplatesSplitDB_ExportImportAndPrecompile(t *testing.T) {
	db1 := openTestTemplatesDB(t)
	defer db1.Close()

	seedTemplateAndVarsForExport(t, db1)
	exportPath := filepath.Join(t.TempDir(), "exported-templates.json")
	payload, err := db1.ExportTemplatesToFile(exportPath, "seo", "seo-brand-01")
	if err != nil {
		t.Fatalf("ExportTemplatesToFile failed: %v", err)
	}
	if !strings.HasPrefix(payload.ExportID, "sha256-") || payload.Variables["BRAND"] != "GitMap" {
		t.Fatalf("unexpected export payload: %+v", payload)
	}
	if len(payload.Categories) == 0 || len(payload.Categories[0].Items) != 1 {
		t.Fatalf("expected category-sequenced items in export payload: %+v", payload.Categories)
	}

	verifyImportDeduplicationAndPrecompile(t, exportPath)
}

func seedTemplateAndVarsForExport(t *testing.T, db *TemplatesSplitDB) {
	t.Helper()
	if err := db.UpsertVariable("BRAND", "global", "GitMap", "Brand name"); err != nil {
		t.Fatalf("UpsertVariable failed: %v", err)
	}
	_, err := db.UpsertTemplateItem(StateTemplateItem{
		ID:       "seo-brand-01",
		Category: "seo",
		Slug:     "brand-velocity",
		Title:    "# Why does $BRAND optimize $TARGET?",
		Text:     "Because ${BRAND} reduces latency by 64% across $TARGET.",
	})
	if err != nil {
		t.Fatalf("UpsertTemplateItem failed: %v", err)
	}
}

func verifyImportDeduplicationAndPrecompile(t *testing.T, exportPath string) {
	t.Helper()
	db2 := openTestTemplatesDB(t)
	defer db2.Close()

	count1, isSkipped1, _, err := db2.ImportTemplatesFromFile(exportPath, false)
	if err != nil || isSkipped1 || count1 != 1 {
		t.Fatalf("first import failed: count=%d skipped=%v err=%v", count1, isSkipped1, err)
	}

	count2, isSkipped2, _, err := db2.ImportTemplatesFromFile(exportPath, false)
	if err != nil || !isSkipped2 || count2 != 0 {
		t.Fatalf("second import should skip: count=%d skipped=%v err=%v", count2, isSkipped2, err)
	}

	verifyPrecompiledOutput(t, db2)
}

func verifyPrecompiledOutput(t *testing.T, db *TemplatesSplitDB) {
	t.Helper()
	compiled, err := db.PrecompileTemplates("seo-brand-01", map[string]string{"TARGET": "pipelines"})
	if err != nil || len(compiled) != 1 {
		t.Fatalf("PrecompileTemplates failed: len=%d err=%v", len(compiled), err)
	}

	got := compiled[0]
	if !strings.Contains(got.Title, "GitMap") || !strings.Contains(got.Text, "pipelines") {
		t.Fatalf("variables not expanded in precompiled template: %+v", got)
	}
}
