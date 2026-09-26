package store

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// StateTemplateCategory represents a template category in gitmap-templates.db.
type StateTemplateCategory struct {
	CategoryID  int64  `json:"categoryId,omitempty"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	ParentSlug  string `json:"parentSlug,omitempty"`
	Description string `json:"description,omitempty"`
	IsDefault   bool   `json:"isDefault,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

// StateTemplateItem represents a template record in gitmap-templates.db.
type StateTemplateItem struct {
	ID          string         `json:"id"`
	Category    string         `json:"category"`
	SubCategory string         `json:"subCategory,omitempty"`
	Slug        string         `json:"slug"`
	Title       string         `json:"title"`
	Text        string         `json:"text"`
	Additional  map[string]any `json:"additional,omitempty"`
	CreatedAt   string         `json:"createdAt,omitempty"`
	UpdatedAt   string         `json:"updatedAt,omitempty"`
}

// StateTemplateVariable represents a key-value variable in gitmap-templates.db.
type StateTemplateVariable struct {
	Key         string `json:"key"`
	Scope       string `json:"scope,omitempty"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

// UpsertCategory inserts or updates a template category by Slug.
func (s *TemplatesSplitDB) UpsertCategory(cat StateTemplateCategory) error {
	slug := normalizeTemplateSlug(cat.Slug, cat.Name)
	name := strings.TrimSpace(cat.Name)
	if name == "" {
		name = slug
	}

	query := `INSERT INTO TemplateCategory (Slug, Name, ParentSlug, Description, IsDefault, UpdatedAt)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(Slug) DO UPDATE SET
			Name = excluded.Name,
			ParentSlug = excluded.ParentSlug,
			Description = excluded.Description,
			UpdatedAt = CURRENT_TIMESTAMP`
	_, err := s.conn.Exec(query, slug, name, cat.ParentSlug, cat.Description, boolToInt(cat.IsDefault))
	if err != nil {
		return apperror.WrapSimple(err, "templates_split.upsert_category")
	}

	return nil
}

// ListCategories returns all template categories ordered by default status and slug.
func (s *TemplatesSplitDB) ListCategories() ([]StateTemplateCategory, error) {
	query := `SELECT CategoryId, Slug, Name, ParentSlug, Description, IsDefault, CreatedAt, UpdatedAt
		FROM TemplateCategory ORDER BY IsDefault DESC, Slug ASC`
	rows, err := s.conn.Query(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "templates_split.list_categories")
	}
	defer rows.Close()

	return scanTemplateCategories(rows)
}

func scanTemplateCategories(rows *sql.Rows) ([]StateTemplateCategory, error) {
	var out []StateTemplateCategory
	for rows.Next() {
		var c StateTemplateCategory
		var isDef int
		if err := rows.Scan(&c.CategoryID, &c.Slug, &c.Name, &c.ParentSlug, &c.Description, &isDef, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, apperror.WrapSimple(err, "templates_split.scan_category")
		}
		c.IsDefault = isDef == 1
		out = append(out, c)
	}

	return out, nil
}

// UpsertTemplateItem matches by ID first, then by Slug, or inserts a new template item.
func (s *TemplatesSplitDB) UpsertTemplateItem(item StateTemplateItem) (StateTemplateItem, error) {
	norm := normalizeTemplateItemInput(item)
	_ = s.ensureCategoryExists(norm.Category)
	addJSON := marshalAdditionalJSON(norm.Additional)

	resolvedID, err := s.persistTemplateItem(norm, addJSON)
	if err != nil {
		return norm, err
	}

	norm.ID = resolvedID
	if saved, getErr := s.GetTemplateItem(resolvedID); getErr == nil && saved != nil {
		return *saved, nil
	}

	return norm, nil
}

func (s *TemplatesSplitDB) persistTemplateItem(item StateTemplateItem, addJSON string) (string, error) {
	if item.ID != "" && s.hasTemplateID(item.ID) {
		return item.ID, s.updateTemplateByID(item, addJSON)
	}

	if existingID := s.findTemplateIDBySlug(item.Slug); existingID != "" {
		return s.updateTemplateBySlug(existingID, item, addJSON)
	}

	if item.ID == "" {
		item.ID = "tpl-" + item.Slug
	}

	return item.ID, s.insertTemplateRow(item, addJSON)
}

func (s *TemplatesSplitDB) hasTemplateID(id string) bool {
	var found string
	err := s.conn.QueryRow(`SELECT ItemId FROM TemplateItem WHERE ItemId = ? LIMIT 1`, id).Scan(&found)

	return err == nil && found != ""
}

func (s *TemplatesSplitDB) findTemplateIDBySlug(slug string) string {
	var found string
	_ = s.conn.QueryRow(`SELECT ItemId FROM TemplateItem WHERE Slug = ? LIMIT 1`, slug).Scan(&found)

	return found
}

func (s *TemplatesSplitDB) updateTemplateByID(item StateTemplateItem, addJSON string) error {
	query := `UPDATE TemplateItem SET CategorySlug = ?, SubCategorySlug = ?, Slug = ?,
		Title = ?, Text = ?, AdditionalJson = ?, UpdatedAt = CURRENT_TIMESTAMP WHERE ItemId = ?`
	_, err := s.conn.Exec(query, item.Category, item.SubCategory, item.Slug, item.Title, item.Text, addJSON, item.ID)
	if err != nil {
		return apperror.WrapSimple(err, "templates_split.update_by_id")
	}

	return nil
}

func (s *TemplatesSplitDB) updateTemplateBySlug(existingID string, item StateTemplateItem, addJSON string) (string, error) {
	targetID := existingID
	if item.ID != "" {
		targetID = item.ID
	}

	query := `UPDATE TemplateItem SET ItemId = ?, CategorySlug = ?, SubCategorySlug = ?,
		Title = ?, Text = ?, AdditionalJson = ?, UpdatedAt = CURRENT_TIMESTAMP WHERE Slug = ?`
	_, err := s.conn.Exec(query, targetID, item.Category, item.SubCategory, item.Title, item.Text, addJSON, item.Slug)
	if err != nil {
		return "", apperror.WrapSimple(err, "templates_split.update_by_slug")
	}

	return targetID, nil
}

func (s *TemplatesSplitDB) insertTemplateRow(item StateTemplateItem, addJSON string) error {
	query := `INSERT INTO TemplateItem (ItemId, CategorySlug, SubCategorySlug, Slug, Title, Text, AdditionalJson)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.conn.Exec(query, item.ID, item.Category, item.SubCategory, item.Slug, item.Title, item.Text, addJSON)
	if err != nil {
		return apperror.WrapSimple(err, "templates_split.insert_item")
	}

	return nil
}

func (s *TemplatesSplitDB) ensureCategoryExists(slug string) error {
	_, err := s.conn.Exec(`INSERT OR IGNORE INTO TemplateCategory (Slug, Name) VALUES (?, ?)`, slug, slug)

	return err
}

// DeleteTemplateItem deletes a template item matching either ItemId or Slug.
func (s *TemplatesSplitDB) DeleteTemplateItem(idOrSlug string) error {
	key := strings.TrimSpace(idOrSlug)
	if _, err := s.conn.Exec(`DELETE FROM TemplateItem WHERE ItemId = ? OR Slug = ?`, key, key); err != nil {
		return apperror.WrapSimple(err, "templates_split.delete_item")
	}

	return nil
}

// GetTemplateItem retrieves a single template item by ItemId or Slug.
func (s *TemplatesSplitDB) GetTemplateItem(idOrSlug string) (*StateTemplateItem, error) {
	key := strings.TrimSpace(idOrSlug)
	query := `SELECT ItemId, CategorySlug, SubCategorySlug, Slug, Title, Text, AdditionalJson, CreatedAt, UpdatedAt
		FROM TemplateItem WHERE ItemId = ? OR Slug = ? LIMIT 1`
	row := s.conn.QueryRow(query, key, key)

	return scanSingleTemplateRow(row)
}

// ListTemplateItems lists templates, optionally filtered by CategorySlug or SubCategorySlug.
func (s *TemplatesSplitDB) ListTemplateItems(categoryFilter string) ([]StateTemplateItem, error) {
	cat := strings.TrimSpace(categoryFilter)
	rows, err := s.queryTemplateRows(cat)
	if err != nil {
		return nil, apperror.WrapSimple(err, "templates_split.list_items")
	}
	defer rows.Close()

	return scanTemplateItems(rows)
}

func (s *TemplatesSplitDB) queryTemplateRows(cat string) (*sql.Rows, error) {
	if cat == "" {
		query := `SELECT ItemId, CategorySlug, SubCategorySlug, Slug, Title, Text, AdditionalJson, CreatedAt, UpdatedAt
			FROM TemplateItem ORDER BY CategorySlug ASC, Slug ASC`

		return s.conn.Query(query)
	}

	query := `SELECT ItemId, CategorySlug, SubCategorySlug, Slug, Title, Text, AdditionalJson, CreatedAt, UpdatedAt
		FROM TemplateItem WHERE CategorySlug = ? OR SubCategorySlug = ? ORDER BY Slug ASC`

	return s.conn.Query(query, cat, cat)
}

func scanTemplateItems(rows *sql.Rows) ([]StateTemplateItem, error) {
	var out []StateTemplateItem
	for rows.Next() {
		item, err := scanSingleTemplateRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *item)
	}

	return out, nil
}

func scanSingleTemplateRow(sc rowScanner) (*StateTemplateItem, error) {
	var item StateTemplateItem
	var addJSON string
	err := sc.Scan(
		&item.ID, &item.Category, &item.SubCategory, &item.Slug,
		&item.Title, &item.Text, &addJSON, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, apperror.WrapSimple(err, "templates_split.scan_item")
	}

	item.Additional = parseAdditionalJSON(addJSON)

	return &item, nil
}

// UpsertVariable inserts or updates a template variable in gitmap-templates.db.
func (s *TemplatesSplitDB) UpsertVariable(key, scope, value, desc string) error {
	cleanKey := strings.TrimSpace(key)
	cleanScope := normalizeVarScope(scope)
	query := `INSERT INTO TemplateVariable (VarKey, Scope, VarValue, Description, UpdatedAt)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(VarKey, Scope) DO UPDATE SET
			VarValue = excluded.VarValue,
			Description = CASE WHEN excluded.Description != '' THEN excluded.Description ELSE TemplateVariable.Description END,
			UpdatedAt = CURRENT_TIMESTAMP`
	if _, err := s.conn.Exec(query, cleanKey, cleanScope, value, desc); err != nil {
		return apperror.WrapSimple(err, "templates_split.upsert_variable")
	}

	return nil
}

// ListVariables returns a key-value map of variables for the requested scope (merged with global).
func (s *TemplatesSplitDB) ListVariables(scope string) (map[string]string, error) {
	records, err := s.ListVariableRecords(scope)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string, len(records))
	for _, rec := range records {
		out[rec.Key] = rec.Value
	}

	return out, nil
}

// ListVariableRecords returns structured variable records for the given scope.
func (s *TemplatesSplitDB) ListVariableRecords(scope string) ([]StateTemplateVariable, error) {
	cleanScope := strings.TrimSpace(scope)
	query := `SELECT VarKey, Scope, VarValue, Description, UpdatedAt FROM TemplateVariable
		WHERE ? = '' OR Scope = 'global' OR Scope = ?
		ORDER BY CASE WHEN Scope = 'global' THEN 0 ELSE 1 END ASC, VarKey ASC`
	rows, err := s.conn.Query(query, cleanScope, cleanScope)
	if err != nil {
		return nil, apperror.WrapSimple(err, "templates_split.list_variables")
	}
	defer rows.Close()

	return scanVariableRows(rows)
}

func scanVariableRows(rows *sql.Rows) ([]StateTemplateVariable, error) {
	var out []StateTemplateVariable
	for rows.Next() {
		var v StateTemplateVariable
		if err := rows.Scan(&v.Key, &v.Scope, &v.Value, &v.Description, &v.UpdatedAt); err != nil {
			return nil, apperror.WrapSimple(err, "templates_split.scan_variable")
		}
		out = append(out, v)
	}

	return out, nil
}

// DeleteVariable removes a variable key from the specified scope.
func (s *TemplatesSplitDB) DeleteVariable(key, scope string) error {
	cleanKey := strings.TrimSpace(key)
	cleanScope := normalizeVarScope(scope)
	if _, err := s.conn.Exec(`DELETE FROM TemplateVariable WHERE VarKey = ? AND Scope = ?`, cleanKey, cleanScope); err != nil {
		return apperror.WrapSimple(err, "templates_split.delete_variable")
	}

	return nil
}

func normalizeVarScope(scope string) string {
	clean := strings.TrimSpace(scope)
	if clean == "" {
		return "global"
	}

	return clean
}

func normalizeTemplateItemInput(item StateTemplateItem) StateTemplateItem {
	item.ID = strings.TrimSpace(item.ID)
	item.Category = strings.TrimSpace(item.Category)
	if item.Category == "" {
		item.Category = "seo"
	}

	item.Slug = normalizeTemplateSlug(item.Slug, item.ID, item.Title)

	return item
}

func normalizeTemplateSlug(candidates ...string) string {
	for _, c := range candidates {
		slug := slugifyTemplateToken(c)
		if slug != "" {
			return slug
		}
	}

	return "template-item"
}

func slugifyTemplateToken(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return ""
	}

	var b strings.Builder
	isPrevDash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			isPrevDash = false
		} else if !isPrevDash && b.Len() > 0 {
			b.WriteByte('-')
			isPrevDash = true
		}
	}

	return strings.Trim(b.String(), "-")
}

func marshalAdditionalJSON(m map[string]any) string {
	if len(m) == 0 {
		return "{}"
	}

	raw, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}

	return string(raw)
}

func parseAdditionalJSON(raw string) map[string]any {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed == "{}" {
		return nil
	}

	var out map[string]any
	if err := json.Unmarshal([]byte(trimmed), &out); err != nil {
		return nil
	}

	return out
}
