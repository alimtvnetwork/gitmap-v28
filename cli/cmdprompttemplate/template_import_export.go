package cmdprompttemplate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// TemplateSearchIndex models the catalog index for instant template discovery.
type TemplateSearchIndex struct {
	TotalTemplates int                   `json:"totalTemplates"`
	UpdatedAt      string                `json:"updatedAt"`
	Items          []TemplateSearchEntry `json:"items"`
}

// TemplateSearchEntry represents one indexed template item.
type TemplateSearchEntry struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Length    int    `json:"length"`
	Snippet   string `json:"snippet"`
	UpdatedAt string `json:"updatedAt"`
}

// ResolveTemplateIndexFilePath returns the filepath where the template search index is persisted.
func ResolveTemplateIndexFilePath() string {
	storePath := ResolveTemplateStorePath()
	return filepath.Join(filepath.Dir(storePath), "prompts_templates_index.json")
}

// RebuildTemplateIndex builds and saves the template search index from current templates.
func RebuildTemplateIndex() error {
	list, err := LoadTemplates()
	if err != nil {
		return err
	}

	idx := buildTemplateSearchIndex(list)
	return saveTemplateSearchIndex(idx)
}

func buildTemplateSearchIndex(list []PromptTemplate) TemplateSearchIndex {
	items := make([]TemplateSearchEntry, 0, len(list))
	for _, tpl := range list {
		snippet := tpl.Content
		if len(snippet) > 100 {
			snippet = snippet[:100] + "..."
		}
		items = append(items, TemplateSearchEntry{
			ID:        tpl.ID,
			Name:      tpl.Name,
			Length:    len(tpl.Content),
			Snippet:   snippet,
			UpdatedAt: tpl.UpdatedAt.Format(time.RFC3339),
		})
	}

	return TemplateSearchIndex{
		TotalTemplates: len(list),
		UpdatedAt:      time.Now().UTC().Format(time.RFC3339),
		Items:          items,
	}
}

func saveTemplateSearchIndex(idx TemplateSearchIndex) error {
	destPath := ResolveTemplateIndexFilePath()
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir template index dir")
	}

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal template index")
	}

	return os.WriteFile(destPath, data, 0644)
}

// ExportSingleTemplate exports one template as formatted JSON to a file.
func ExportSingleTemplate(identifier, destPath string) error {
	tpl, found := FindTemplate(identifier)
	if !found {
		return apperror.NewNotFoundError("template not found: " + identifier)
	}
	data, err := json.MarshalIndent(tpl, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal template")
	}

	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return err
	}
	_ = RebuildTemplateIndex()
	return nil
}

// ImportSingleTemplate imports a single template from JSON into storage.
func ImportSingleTemplate(srcPath string) (*PromptTemplate, error) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read template file")
	}
	var tpl PromptTemplate
	if parseErr := json.Unmarshal(data, &tpl); parseErr != nil {
		return nil, apperror.WrapSimple(parseErr, "parse template JSON")
	}

	res, err := AddTemplate(tpl.Name, tpl.Content)
	if err != nil {
		return nil, err
	}
	_ = RebuildTemplateIndex()
	return res, nil
}

// ExportAllTemplates exports all registered templates into a suite JSON file.
func ExportAllTemplates(destPath string) error {
	list, err := LoadTemplates()
	if err != nil {
		return err
	}

	if err := writeTemplateSuiteFile(list, destPath); err != nil {
		return err
	}
	_ = RebuildTemplateIndex()
	return nil
}

func writeTemplateSuiteFile(list []PromptTemplate, destPath string) error {
	suite := TemplateSuite{Version: "1.0.0", Templates: list}
	data, err := json.MarshalIndent(suite, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal template suite")
	}

	return os.WriteFile(destPath, data, 0644)
}

// ImportAllTemplates bulk imports templates from a suite JSON file.
func ImportAllTemplates(srcPath string) (int, error) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return 0, apperror.WrapSimple(err, "read template suite file")
	}

	count, err := parseAndImportSuite(data)
	if err != nil {
		return 0, err
	}
	_ = RebuildTemplateIndex()
	return count, nil
}

func parseAndImportSuite(data []byte) (int, error) {
	var suite TemplateSuite
	if parseErr := json.Unmarshal(data, &suite); parseErr != nil {
		return 0, apperror.WrapSimple(parseErr, "parse template suite JSON")
	}

	return insertSuiteTemplates(suite.Templates), nil
}

func insertSuiteTemplates(templates []PromptTemplate) int {
	count := 0
	for _, item := range templates {
		if _, addErr := AddTemplate(item.Name, item.Content); addErr == nil {
			count++
		}
	}

	return count
}
