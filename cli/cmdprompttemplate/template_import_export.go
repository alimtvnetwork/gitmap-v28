package cmdprompttemplate

import (
	"encoding/json"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

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

	return os.WriteFile(destPath, data, 0644)
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

	return AddTemplate(tpl.Name, tpl.Content)
}

// ExportAllTemplates exports all registered templates into a suite JSON file.
func ExportAllTemplates(destPath string) error {
	list, err := LoadTemplates()
	if err != nil {
		return err
	}

	return writeTemplateSuiteFile(list, destPath)
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

	return parseAndImportSuite(data)
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
