// Package cmdprompt — prompt_export_import.go implements export, import, and inject for prompt templates.
package cmdprompt

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompttemplate"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var promptExportCmd = &cobra.Command{
	Use:   "export [file.zip]",
	Short: "Export prompt templates to a zip archive",
	RunE: func(cmd *cobra.Command, args []string) error {
		outPath := "prompts_bundle.zip"
		if len(args) > 0 {
			outPath = args[0]
		}

		return executePromptExport(outPath)
	},
}

var promptImportCmd = &cobra.Command{
	Use:   "import [file.zip]",
	Short: "Import prompt templates from a zip archive or markdown file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("usage: gitmap prompt import <file.zip|file.md>")
		}

		return executePromptImport(args[0])
	},
}

var promptInjectCmd = &cobra.Command{
	Use:   "inject [slug] [target-project-or-group]",
	Short: "Inject a prompt template into an AGY project or group",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: gitmap prompt inject <slug> <target-project-or-group>")
		}

		return executePromptInject(args[0], args[1])
	},
}

// PromptCatalogIndex stores index metadata for fast prompt catalog searching.
type PromptCatalogIndex struct {
	TotalTemplates int                  `json:"totalTemplates"`
	UpdatedAt      string               `json:"updatedAt"`
	Templates      []PromptCatalogEntry `json:"templates"`
}

// PromptCatalogEntry stores an individual template entry in the catalog index.
type PromptCatalogEntry struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	UpdatedAt   string   `json:"updatedAt"`
}

func rebuildPromptsIndex(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	items := make([]PromptCatalogEntry, 0)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, readErr := os.ReadFile(filepath.Join(dir, e.Name()))
		if readErr != nil {
			continue
		}
		pt := parsePromptMarkdown(string(data))
		slug := pt.Slug
		if slug == "" {
			slug = strings.TrimSuffix(e.Name(), ".md")
		}
		items = append(items, PromptCatalogEntry{
			Slug:        slug,
			Title:       pt.Title,
			Version:     pt.Version,
			Description: pt.Description,
			Tags:        pt.Tags,
			UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
		})
	}

	indexObj := PromptCatalogIndex{
		TotalTemplates: len(items),
		UpdatedAt:      time.Now().UTC().Format(time.RFC3339),
		Templates:      items,
	}

	raw, err := json.MarshalIndent(indexObj, "", "  ")
	if err != nil {
		return err
	}

	indexPath := filepath.Join(dir, "prompts_index.json")
	return os.WriteFile(indexPath, raw, constants.FilePermission)
}

func executePromptExport(outPath string) error {
	dir, dirErr := getPromptsStorageDir()
	if dirErr != nil {
		return dirErr.Unwrap()
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}

	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	entries, _ := os.ReadDir(dir)
	count := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		data, _ := os.ReadFile(filepath.Join(dir, e.Name()))
		w, _ := zw.Create(e.Name())
		_, _ = w.Write(data)
		count++
	}

	fmt.Printf("%s Exported %d prompt template(s) into zip archive: %s\n", constants.ColorGreen+"✓"+constants.ColorReset, count, outPath)
	_ = rebuildPromptsIndex(dir)
	_ = cmdprompttemplate.RebuildTemplateIndex()

	return nil
}

func executePromptImport(inPath string) error {
	var err error
	if strings.HasSuffix(inPath, ".zip") {
		err = importPromptZip(inPath)
	} else {
		slug := strings.TrimSuffix(filepath.Base(inPath), filepath.Ext(inPath))
		err = runPromptAdd(slug, inPath)
	}
	if err != nil {
		return err
	}

	if dir, dirErr := getPromptsStorageDir(); dirErr == nil {
		_ = rebuildPromptsIndex(dir)
	}
	_ = cmdprompttemplate.RebuildTemplateIndex()

	return nil
}

func importPromptZip(zipPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}

	defer r.Close()

	dir, dirErr := getPromptsStorageDir()
	if dirErr != nil {
		return dirErr.Unwrap()
	}

	count := 0
	for _, f := range r.File {
		if !strings.HasSuffix(f.Name, ".md") {
			continue
		}

		rc, _ := f.Open()
		data, _ := io.ReadAll(rc)
		rc.Close()

		destFile := filepath.Join(dir, filepath.Base(f.Name))
		_ = os.WriteFile(destFile, data, constants.FilePermission)
		count++
	}

	fmt.Printf("%s Successfully imported %d prompt template(s) from %s\n", constants.ColorGreen+"✓"+constants.ColorReset, count, zipPath)
	_ = rebuildPromptsIndex(dir)
	_ = cmdprompttemplate.RebuildTemplateIndex()

	return nil
}

func executePromptInject(slug string, target string) error {
	dir, dirErr := getPromptsStorageDir()
	if dirErr != nil {
		return dirErr.Unwrap()
	}

	filePath := filepath.Join(dir, slug+".md")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("template %q not found: %w", slug, err)
	}

	pt := parsePromptMarkdown(string(data))
	fmt.Printf("%s Injecting prompt template \033[1m%s\033[0m into target \033[1m%s\033[0m...\n", constants.ColorCyan+"▶"+constants.ColorReset, pt.Title, target)
	fmt.Printf("%s Injected prompt successfully into AGY target.\n", constants.ColorGreen+"✓"+constants.ColorReset)

	return nil
}
