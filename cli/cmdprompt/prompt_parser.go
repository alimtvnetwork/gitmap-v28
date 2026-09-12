// Package cmdprompt — prompt_parser.go implements markdown prompt frontmatter and body parsing.
package cmdprompt

import (
	"strings"
)

func parsePromptMarkdown(content string) PromptTemplate {
	pt := PromptTemplate{
		Version: "1.0.0",
		Tags:    make([]string, 0),
	}

	lines := strings.Split(content, "\n")
	hasFrontmatter := len(lines) > 2 && strings.TrimSpace(lines[0]) == "---"
	if !hasFrontmatter {
		pt.Body = strings.TrimSpace(content)

		return finalizePromptTemplate(pt)
	}

	parseFrontmatterLines(lines, &pt)

	return finalizePromptTemplate(pt)
}

func parseFrontmatterLines(lines []string, pt *PromptTemplate) {
	bodyStart := 1
	for i := 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "---" {
			bodyStart = i + 1
			break
		}

		parseFrontmatterLine(trimmed, pt)
	}

	assignPromptBody(lines, bodyStart, pt)
}

func assignPromptBody(lines []string, bodyStart int, pt *PromptTemplate) {
	if bodyStart < len(lines) {
		pt.Body = strings.TrimSpace(strings.Join(lines[bodyStart:], "\n"))
	}
}

func finalizePromptTemplate(pt PromptTemplate) PromptTemplate {
	if pt.Title == "" {
		pt.Title = "Custom Prompt"
	}

	return pt
}

func parseFrontmatterLine(line string, pt *PromptTemplate) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}

	key := strings.ToLower(strings.TrimSpace(parts[0]))
	val := strings.Trim(strings.TrimSpace(parts[1]), "\"'")

	switch key {
	case "title":
		pt.Title = val
	case "slug":
		pt.Slug = val
	case "version":
		pt.Version = val
	case "description":
		pt.Description = val
	}
}
