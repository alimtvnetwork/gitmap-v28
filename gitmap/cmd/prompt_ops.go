// Package cmd — prompt_ops.go implements list, show, add, and remove operations for prompt templates.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runPromptList() error {
	dir, dirErr := getPromptsStorageDir()
	if dirErr != nil {
		return dirErr.Unwrap()
	}

	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return readErr
	}

	templates := make([]PromptTemplate, 0)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		pt := parsePromptMarkdown(string(data))
		if pt.Slug == "" {
			pt.Slug = strings.TrimSuffix(e.Name(), ".md")
		}

		templates = append(templates, pt)
	}

	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Slug < templates[j].Slug
	})

	printPromptTable(templates)
	return nil
}

func printPromptTable(templates []PromptTemplate) {
	fmt.Printf("%s Prompt Templates (%d available):\n\n", constants.ColorCyan+"▶"+constants.ColorReset, len(templates))
	fmt.Printf("    %-20s %-8s %s\n", "SLUG", "VER", "DESCRIPTION")
	fmt.Printf("    %s\n", "──────────────────────────────────────────────────────────────────")

	for _, t := range templates {
		desc := t.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}

		fmt.Printf("  • %-20s %-8s %s\n", t.Slug, t.Version, desc)
	}
}

func runPromptShow(slug string) error {
	dir, dirErr := getPromptsStorageDir()
	if dirErr != nil {
		return dirErr.Unwrap()
	}

	filePath := filepath.Join(dir, slug+".md")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("prompt template %q not found", slug)
	}

	pt := parsePromptMarkdown(string(data))
	fmt.Printf("%s Template: \033[1m%s\033[0m (v%s)\n", constants.ColorCyan+"▶"+constants.ColorReset, pt.Title, pt.Version)
	if pt.Description != "" {
		fmt.Printf("  Description: %s\n", pt.Description)
	}

	fmt.Printf("\n%s\n", pt.Body)
	return nil
}

func runPromptAdd(slug string, srcPath string) error {
	data, readErr := os.ReadFile(srcPath)
	if readErr != nil {
		return readErr
	}

	dir, dirErr := getPromptsStorageDir()
	if dirErr != nil {
		return dirErr.Unwrap()
	}

	destPath := filepath.Join(dir, slug+".md")
	if writeErr := os.WriteFile(destPath, data, constants.FilePermission); writeErr != nil {
		return writeErr
	}

	fmt.Printf("%s Successfully saved prompt template %q to %s\n", constants.ColorGreen+"✓"+constants.ColorReset, slug, destPath)
	return nil
}

func runPromptRm(slug string) error {
	dir, dirErr := getPromptsStorageDir()
	if dirErr != nil {
		return dirErr.Unwrap()
	}

	destPath := filepath.Join(dir, slug+".md")
	if err := os.Remove(destPath); err != nil {
		return fmt.Errorf("remove template: %w", err)
	}

	fmt.Printf("%s Removed prompt template %q.\n", constants.ColorGreen+"✓"+constants.ColorReset, slug)
	return nil
}
