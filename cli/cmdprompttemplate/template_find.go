package cmdprompttemplate

import "strings"

// FindTemplate looks up a template by ID or name prefix.
func FindTemplate(identifier string) (*PromptTemplate, bool) {
	list, err := LoadTemplates()
	if err != nil {
		return nil, false
	}

	return searchTemplateInList(list, strings.ToLower(strings.TrimSpace(identifier)))
}

func searchTemplateInList(list []PromptTemplate, target string) (*PromptTemplate, bool) {
	for _, item := range list {
		if isMatchTemplate(item, target) {
			return &item, true
		}
	}

	return nil, false
}

func isMatchTemplate(item PromptTemplate, target string) bool {
	idLower := strings.ToLower(item.ID)
	nameLower := strings.ToLower(item.Name)
	if idLower == target || nameLower == target {
		return true
	}

	return strings.HasPrefix(nameLower, target) || strings.HasPrefix(idLower, target)
}
