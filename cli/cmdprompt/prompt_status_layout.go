// Package cmdprompt — prompt_status_layout.go calculates column dimensions for prompt status table.
package cmdprompt

import (
	"fmt"
)

func NewPromptStatusTableLayout() *PromptStatusTableLayout {
	return &PromptStatusTableLayout{
		MaxRepo:    20,
		MaxStatus:  15,
		MaxVersion: 10,
		MaxDate:    20,
	}
}

func (l *PromptStatusTableLayout) PrintHeader() {
	fmt.Printf("  %-*s   %-*s   %-*s   %s\n",
		l.MaxRepo, "REPO",
		l.MaxStatus, "STATUS",
		l.MaxVersion, "VERSION",
		"INSTALLED AT",
	)
	fmt.Println("  --------------------------------------------------------------------------------")
}
