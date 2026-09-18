package cmdprompttemplate

import (
	"time"
)

// DefaultIsDoneContent is the pre-installed default verification prompt content.
const DefaultIsDoneContent = "Is it done properly? Can we check properly the missing items from the task that is mentioned below? Please check it carefully. Do not make any mistakes."

// DefaultTemplateID is the canonical ID for the default verification template.
const DefaultTemplateID = "is-done"

// PromptTemplate stores a reusable AGY prompt prefix or task template.
type PromptTemplate struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TemplateSuite holds a collection of prompt templates for bulk export and import.
type TemplateSuite struct {
	Version   string           `json:"version"`
	Templates []PromptTemplate `json:"templates"`
}

// BuildDefaultTemplate constructs the standard pre-installed is-done template.
func BuildDefaultTemplate() PromptTemplate {
	now := time.Now().UTC()

	return PromptTemplate{
		ID:          DefaultTemplateID,
		Name:        DefaultTemplateID,
		Description: "Standard task completion and quality verification check",
		Content:     DefaultIsDoneContent,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
