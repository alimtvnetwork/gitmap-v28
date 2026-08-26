// Package cmd — prompt_status_layout.go calculates column dimensions for prompt status table.
package cmd

import (
	"fmt"

)

type PromptStatusTableLayout struct {
	MaxRepo    int
	MaxStatus  int
	MaxVersion int
	MaxDate    int
}

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
