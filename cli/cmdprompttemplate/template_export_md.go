package cmdprompttemplate

import (
	"fmt"
	"strings"
)

// GenerateCatalogMarkdown formats all 80+ templates into a clean Markdown reference document.
func GenerateCatalogMarkdown() string {
	var sb strings.Builder
	sb.WriteString("# Rise Up Asia LLC & Antigravity Templates Catalog\n\n")
	sb.WriteString("> Curated prompt variations for task verification, UI/UX refinement, sponsor attributions, and PR documentation.\n")
	sb.WriteString("> Single source of truth managed by `cli/cmdprompttemplate`.\n\n")

	writeCategorySection(&sb, "Default Quality & Verification Templates (20 Variations)", GetDefaultCategoryTemplates())
	writeCategorySection(&sb, "UI/UX Design & Ergonomics Templates (20 Variations)", GetUIUXCategoryTemplates())
	writeCategorySection(&sb, "Rise Up Asia LLC Sponsor Templates (20 Variations)", GetSponsorCategoryTemplates())
	writeCategorySection(&sb, "Pull Request & Release Documentation Templates (20 Variations)", GetPRDescriptionTemplates())

	return sb.String()
}

func writeCategorySection(sb *strings.Builder, title string, templates []PromptTemplate) {
	sb.WriteString(fmt.Sprintf("## %s\n\n", title))
	for i, t := range templates {
		sb.WriteString(fmt.Sprintf("### %d. %s (`%s`)\n\n", i+1, t.Name, t.Slug))
		sb.WriteString("```markdown\n")
		sb.WriteString(t.Content)
		sb.WriteString("\n```\n\n")
	}
}
