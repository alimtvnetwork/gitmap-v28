package cmdprompttemplate

import (
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)


// AddTemplate creates and saves a new prompt template.
func AddTemplate(name, content string) (*PromptTemplate, error) {
	list, err := LoadTemplates()
	if err != nil {
		return nil, err
	}
	tpl := createNewTemplate(name, content)
	list = append(list, tpl)

	return &tpl, SaveTemplates(list)
}

func createNewTemplate(name, content string) PromptTemplate {
	now := time.Now().UTC()
	return PromptTemplate{
		ID:        strings.ToLower(strings.ReplaceAll(name, " ", "-")),
		Name:      name,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// EditTemplate updates content of an existing template.
func EditTemplate(identifier, content string) (*PromptTemplate, error) {
	list, err := LoadTemplates()
	if err != nil {
		return nil, err
	}

	return updateMatchingTemplate(list, strings.ToLower(strings.TrimSpace(identifier)), content)
}

func updateMatchingTemplate(list []PromptTemplate, target, content string) (*PromptTemplate, error) {
	for i := range list {
		if isMatchTemplate(list[i], target) {
			list[i].Content = content
			list[i].UpdatedAt = time.Now().UTC()
			return &list[i], SaveTemplates(list)
		}
	}

	return nil, apperror.NewNotFoundError("template not found: " + target)
}

// DeleteTemplate removes a template by identifier.
func DeleteTemplate(identifier string) error {
	list, err := LoadTemplates()
	if err != nil {
		return err
	}

	return removeMatchingTemplate(list, strings.ToLower(strings.TrimSpace(identifier)))
}

func removeMatchingTemplate(list []PromptTemplate, target string) error {
	updated, found := filterOutTemplate(list, target)
	if !found {
		return apperror.NewNotFoundError("template not found: " + target)
	}

	return SaveTemplates(updated)
}

func filterOutTemplate(list []PromptTemplate, target string) ([]PromptTemplate, bool) {
	var updated []PromptTemplate
	found := false
	for _, item := range list {
		if isMatchTemplate(item, target) {
			found = true
			continue
		}
		updated = append(updated, item)
	}

	return updated, found
}
