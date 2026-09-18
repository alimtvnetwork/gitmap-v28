package cmdprompttemplate

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ResolveTemplateStorePath returns the JSON storage path for prompt templates.
func ResolveTemplateStorePath() string {
	dataDir := store.BinaryDataDir()

	return filepath.Join(dataDir, "templates", "prompts_templates.json")
}

// LoadTemplates loads all stored prompt templates, initializing defaults if absent.
func LoadTemplates() ([]PromptTemplate, error) {
	filePath := ResolveTemplateStorePath()
	data, readErr := os.ReadFile(filePath)
	if os.IsNotExist(readErr) {
		return initDefaultStore(filePath)
	}
	if readErr != nil {
		return nil, apperror.WrapSimple(readErr, "read prompt templates")
	}

	return parseTemplateData(data)
}

func parseTemplateData(data []byte) ([]PromptTemplate, error) {
	var list []PromptTemplate
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, apperror.WrapSimple(err, "parse prompt templates JSON")
	}

	return ensureIsDonePresent(list), nil
}

func ensureIsDonePresent(list []PromptTemplate) []PromptTemplate {
	for _, item := range list {
		if item.ID == DefaultTemplateID || item.Name == DefaultTemplateID {
			return list
		}
	}

	return append([]PromptTemplate{BuildDefaultTemplate()}, list...)
}

func initDefaultStore(filePath string) ([]PromptTemplate, error) {
	defaultList := []PromptTemplate{BuildDefaultTemplate()}
	if err := SaveTemplates(defaultList); err != nil {
		return defaultList, err
	}

	return defaultList, nil
}

// SaveTemplates writes prompt templates to atomic JSON storage.
func SaveTemplates(list []PromptTemplate) error {
	filePath := ResolveTemplateStorePath()
	dir := filepath.Dir(filePath)
	if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
		return apperror.WrapSimple(mkErr, "create template directory")
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal prompt templates")
	}

	return os.WriteFile(filePath, data, 0644)
}
