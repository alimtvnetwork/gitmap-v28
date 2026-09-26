package store

import (
	"strings"
)

// CompiledTemplate holds a template item with all static $VAR and ${VAR} placeholders expanded in memory.
type CompiledTemplate struct {
	ID          string         `json:"id"`
	Category    string         `json:"category"`
	SubCategory string         `json:"subCategory,omitempty"`
	Slug        string         `json:"slug"`
	Title       string         `json:"title"`
	Text        string         `json:"text"`
	Additional  map[string]any `json:"additional,omitempty"`
}

// PrecompileTemplates opens gitmap-templates.db and pre-compiles matching templates with stored and extra variables.
func PrecompileTemplates(categoryOrSlug string, extraVars map[string]string) ([]CompiledTemplate, error) {
	db, err := OpenTemplatesSplitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db.PrecompileTemplates(categoryOrSlug, extraVars)
}

// PrecompileCategoryTemplates is an alias for PrecompileTemplates.
func PrecompileCategoryTemplates(categoryOrSlug string, extraVars map[string]string) ([]CompiledTemplate, error) {
	return PrecompileTemplates(categoryOrSlug, extraVars)
}

// PrecompileTemplates loads matching templates and variables from s and expands all $VAR and ${VAR} occurrences once.
func (s *TemplatesSplitDB) PrecompileTemplates(categoryOrSlug string, extraVars map[string]string) ([]CompiledTemplate, error) {
	mergedVars, err := s.loadMergedVariables(extraVars)
	if err != nil {
		return nil, err
	}

	items, err := s.resolveTemplatesByCategoryOrSlug(categoryOrSlug)
	if err != nil {
		return nil, err
	}

	out := make([]CompiledTemplate, 0, len(items))
	for _, it := range items {
		out = append(out, compileSingleTemplate(it, mergedVars))
	}

	return out, nil
}

func (s *TemplatesSplitDB) loadMergedVariables(extraVars map[string]string) (map[string]string, error) {
	storedVars, err := s.ListVariables("")
	if err != nil {
		return nil, err
	}

	merged := make(map[string]string, len(storedVars)+len(extraVars))
	for k, v := range storedVars {
		merged[k] = v
	}
	for k, v := range extraVars {
		merged[k] = v
	}

	return merged, nil
}

func (s *TemplatesSplitDB) resolveTemplatesByCategoryOrSlug(categoryOrSlug string) ([]StateTemplateItem, error) {
	key := strings.TrimSpace(categoryOrSlug)
	items, err := s.ListTemplateItems(key)
	if err != nil {
		return nil, err
	}

	if len(items) > 0 || key == "" {
		return items, nil
	}

	single, getErr := s.GetTemplateItem(key)
	if getErr != nil || single == nil {
		return nil, nil
	}

	return []StateTemplateItem{*single}, nil
}

func compileSingleTemplate(it StateTemplateItem, vars map[string]string) CompiledTemplate {
	return CompiledTemplate{
		ID:          it.ID,
		Category:    it.Category,
		SubCategory: it.SubCategory,
		Slug:        it.Slug,
		Title:       ExpandTemplateVariables(it.Title, vars),
		Text:        ExpandTemplateVariables(it.Text, vars),
		Additional:  it.Additional,
	}
}

// ExpandTemplateVariables replaces $VAR and ${VAR} placeholders in content using vars.
func ExpandTemplateVariables(content string, vars map[string]string) string {
	if content == "" || len(vars) == 0 {
		return content
	}

	expanded := expandDottedVars(content, vars)

	return templateVarRefPattern.ReplaceAllStringFunc(expanded, func(match string) string {
		key := extractVarKeyFromMatch(match)
		if val, hasKey := vars[key]; hasKey {
			return val
		}

		return match
	})
}

func expandDottedVars(content string, vars map[string]string) string {
	out := content
	for k, v := range vars {
		if strings.ContainsAny(k, ".-") {
			out = strings.ReplaceAll(out, "${"+k+"}", v)
			out = strings.ReplaceAll(out, "$"+k, v)
		}
	}

	return out
}

func extractVarKeyFromMatch(match string) string {
	if strings.HasPrefix(match, "${") && strings.HasSuffix(match, "}") {
		return match[2 : len(match)-1]
	}

	return strings.TrimPrefix(match, "$")
}
