package cmdprompttemplate

import (
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// FindTemplate looks up a template by ID or name prefix.
func FindTemplate(identifier string) (*PromptTemplate, bool) {
	list, err := LoadTemplates()
	if err != nil {
		return nil, false
	}
	target := strings.ToLower(strings.TrimSpace(identifier))
	for _, item := range list {
		if strings.ToLower(item.ID) == target || strings.ToLower(item.Name) == target {
			return &item, true
		}
		if strings.HasPrefix(strings.ToLower(item.Name), target) {
			return &item, true
		}
	}

	return nil, false
}

// AddTemplate creates and saves a new prompt template.
func AddTemplate(name, content string) (*PromptTemplate, error) {
	list, err := LoadTemplates()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	tpl := PromptTemplate{
		ID:        strings.ToLower(strings.ReplaceAll(name, " ", "-")),
		Name:      name,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	list = append(list, tpl)

	return &tpl, SaveTemplates(list)
}

// EditTemplate updates content of an existing template.
func EditTemplate(identifier, content string) (*PromptTemplate, error) {
	list, err := LoadTemplates()
	if err != nil {
		return nil, err
	}
	target := strings.ToLower(strings.TrimSpace(identifier))
	for i := range list {
		if isMatchIdentifier(list[i], target) {
			list[i].Content = content
			list[i].UpdatedAt = time.Now().UTC()

			return &list[i], SaveTemplates(list)
		}
	}

	return nil, apperror.NewNotFoundError("template not found: " + identifier)
}

// DeleteTemplate removes a template by identifier.
func DeleteTemplate(identifier string) error {
	list, err := LoadTemplates()
	if err != nil {
		return err
	}
	target := strings.ToLower(strings.TrimSpace(identifier))
	var updated []PromptTemplate
	found := false
	for _, item := range list {
		if isMatchIdentifier(item, target) {
			found = true
			continue
		}
		updated = append(updated, item)
	}
	if !found {
		return apperror.NewNotFoundError("template not found: " + identifier)
	}

	return SaveTemplates(updated)
}

func isMatchIdentifier(item PromptTemplate, target string) bool {
	return strings.ToLower(item.ID) == target || strings.ToLower(item.Name) == target
}
