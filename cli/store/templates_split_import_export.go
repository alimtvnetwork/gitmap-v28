package store

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var templateVarRefPattern = regexp.MustCompile(`\$\{([A-Za-z0-9_.-]+)\}|\$([A-Za-z0-9_]+)`)

// TemplateExportPayload defines the portable JSON schema for exporting and importing templates and variables.
type TemplateExportPayload struct {
	ExportID   string                  `json:"exportId"`
	Version    string                  `json:"version"`
	Categories []StateTemplateCategory `json:"categories,omitempty"`
	Variables  map[string]string       `json:"variables,omitempty"`
	Templates  []StateTemplateItem     `json:"templates"`
}

type canonicalExportDigest struct {
	Categories []StateTemplateCategory `json:"categories,omitempty"`
	Variables  map[string]string       `json:"variables,omitempty"`
	Templates  []StateTemplateItem     `json:"templates"`
}

// ComputeExportHashID computes a deterministic sha256-... digest over Categories, Variables, and Templates.
func ComputeExportHashID(payload TemplateExportPayload) string {
	canon := canonicalExportDigest{
		Categories: stripCategoryTimestamps(payload.Categories),
		Variables:  payload.Variables,
		Templates:  stripTemplateTimestamps(payload.Templates),
	}
	raw, _ := json.Marshal(canon)
	sum := sha256.Sum256(raw)

	return "sha256-" + hex.EncodeToString(sum[:])
}

func stripCategoryTimestamps(cats []StateTemplateCategory) []StateTemplateCategory {
	if len(cats) == 0 {
		return nil
	}

	out := make([]StateTemplateCategory, len(cats))
	for i, c := range cats {
		c.CategoryID = 0
		c.CreatedAt = ""
		c.UpdatedAt = ""
		out[i] = c
	}

	return out
}

func stripTemplateTimestamps(items []StateTemplateItem) []StateTemplateItem {
	if len(items) == 0 {
		return nil
	}

	out := make([]StateTemplateItem, len(items))
	for i, it := range items {
		it.CreatedAt = ""
		it.UpdatedAt = ""
		out[i] = it
	}

	return out
}

// ExportTemplatesToFile exports matching templates and their referenced variables to destPath.
func ExportTemplatesToFile(destPath, categoryFilter, idOrSlugFilter string) (*TemplateExportPayload, error) {
	db, err := OpenTemplatesSplitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db.ExportTemplatesToFile(destPath, categoryFilter, idOrSlugFilter)
}

// ExportTemplatesSuite is an alias for ExportTemplatesToFile.
func ExportTemplatesSuite(destPath, categoryFilter, idOrSlugFilter string) (*TemplateExportPayload, error) {
	return ExportTemplatesToFile(destPath, categoryFilter, idOrSlugFilter)
}

// ExportTemplatesToFile queries matching templates from s, extracts referenced variables, and writes JSON to destPath.
func (s *TemplatesSplitDB) ExportTemplatesToFile(destPath, categoryFilter, idOrSlugFilter string) (*TemplateExportPayload, error) {
	items, err := s.resolveExportTemplates(categoryFilter, idOrSlugFilter)
	if err != nil {
		return nil, err
	}

	allVars, err := s.ListVariables("")
	if err != nil {
		return nil, err
	}

	cats, _ := s.filterExportCategories(items)
	payload := TemplateExportPayload{
		Version:    "1.0",
		Categories: cats,
		Variables:  extractReferencedVars(items, allVars),
		Templates:  items,
	}
	payload.ExportID = ComputeExportHashID(payload)
	if err := writeExportJSONFile(destPath, payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (s *TemplatesSplitDB) resolveExportTemplates(categoryFilter, idOrSlugFilter string) ([]StateTemplateItem, error) {
	if cleanID := strings.TrimSpace(idOrSlugFilter); cleanID != "" {
		item, err := s.GetTemplateItem(cleanID)
		if err != nil {
			return nil, err
		}

		return []StateTemplateItem{*item}, nil
	}

	return s.ListTemplateItems(categoryFilter)
}

func (s *TemplatesSplitDB) filterExportCategories(items []StateTemplateItem) ([]StateTemplateCategory, error) {
	allCats, err := s.ListCategories()
	if err != nil {
		return nil, err
	}

	used := make(map[string]bool, len(items))
	for _, it := range items {
		used[it.Category] = true
	}

	var out []StateTemplateCategory
	for _, c := range allCats {
		if used[c.Slug] {
			out = append(out, c)
		}
	}

	return out, nil
}

func extractReferencedVars(items []StateTemplateItem, allVars map[string]string) map[string]string {
	if len(allVars) == 0 {
		return nil
	}

	referenced := make(map[string]string)
	for _, item := range items {
		collectVarRefsFromText(item.Title, allVars, referenced)
		collectVarRefsFromText(item.Text, allVars, referenced)
	}

	if len(referenced) == 0 {
		return nil
	}

	return referenced
}

func collectVarRefsFromText(content string, allVars, out map[string]string) {
	matches := templateVarRefPattern.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		key := m[1]
		if key == "" {
			key = m[2]
		}
		if val, hasKey := allVars[key]; hasKey {
			out[key] = val
		}
	}
}

func writeExportJSONFile(destPath string, payload TemplateExportPayload) error {
	if dir := filepath.Dir(destPath); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return apperror.WrapSimple(err, "templates_split.export_mkdir")
		}
	}

	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "templates_split.export_marshal")
	}

	if err := os.WriteFile(destPath, append(raw, '\n'), 0644); err != nil {
		return apperror.WrapSimple(err, "templates_split.export_write")
	}

	return nil
}

// ImportTemplatesFromFile imports templates and variables from srcPath with exportId hash deduplication.
func ImportTemplatesFromFile(srcPath string, isForce bool) (int, bool, string, error) {
	db, err := OpenTemplatesSplitDB()
	if err != nil {
		return 0, false, "", err
	}
	defer db.Close()

	return db.ImportTemplatesFromFile(srcPath, isForce)
}

// ImportTemplatesSuite is an alias for ImportTemplatesFromFile.
func ImportTemplatesSuite(srcPath string, isForce bool) (int, bool, string, error) {
	return ImportTemplatesFromFile(srcPath, isForce)
}

// ImportTemplatesFromFile reads a JSON payload, checks TemplateImportHistory, and upserts records.
func (s *TemplatesSplitDB) ImportTemplatesFromFile(srcPath string, isForce bool) (int, bool, string, error) {
	payload, exportID, err := readImportPayload(srcPath)
	if err != nil {
		return 0, false, "", err
	}

	if !isForce && s.isExportHashImported(exportID) {
		return 0, true, exportID, nil
	}

	if err := s.applyImportPayload(payload); err != nil {
		return 0, false, exportID, err
	}

	if err := s.recordImportHistory(exportID, srcPath, len(payload.Templates), len(payload.Variables)); err != nil {
		return 0, false, exportID, err
	}

	return len(payload.Templates), false, exportID, nil
}

func readImportPayload(srcPath string) (TemplateExportPayload, string, error) {
	var payload TemplateExportPayload
	raw, err := os.ReadFile(srcPath)
	if err != nil {
		return payload, "", apperror.WrapSimple(err, "templates_split.import_read")
	}

	if err := json.Unmarshal(raw, &payload); err != nil {
		return payload, "", apperror.WrapSimple(err, "templates_split.import_unmarshal")
	}

	exportID := strings.TrimSpace(payload.ExportID)
	if exportID == "" {
		exportID = ComputeExportHashID(payload)
	}

	return payload, exportID, nil
}

func (s *TemplatesSplitDB) isExportHashImported(exportID string) bool {
	var found string
	err := s.conn.QueryRow(`SELECT ExportHashId FROM TemplateImportHistory WHERE ExportHashId = ? LIMIT 1`, exportID).Scan(&found)

	return err == nil && found != ""
}

func (s *TemplatesSplitDB) applyImportPayload(payload TemplateExportPayload) error {
	for _, cat := range payload.Categories {
		if err := s.UpsertCategory(cat); err != nil {
			return err
		}
	}

	for k, v := range payload.Variables {
		if err := s.UpsertVariable(k, "global", v, ""); err != nil {
			return err
		}
	}

	for _, tpl := range payload.Templates {
		if _, err := s.UpsertTemplateItem(tpl); err != nil {
			return err
		}
	}

	return nil
}

func (s *TemplatesSplitDB) recordImportHistory(exportID, srcPath string, itemCount, varCount int) error {
	query := `INSERT INTO TemplateImportHistory (ExportHashId, SourcePath, ItemCount, VarCount, ImportedAt)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(ExportHashId) DO UPDATE SET
			SourcePath = excluded.SourcePath,
			ItemCount = excluded.ItemCount,
			VarCount = excluded.VarCount,
			ImportedAt = CURRENT_TIMESTAMP`
	if _, err := s.conn.Exec(query, exportID, srcPath, itemCount, varCount); err != nil {
		return apperror.WrapSimple(err, "templates_split.record_import")
	}

	return nil
}
