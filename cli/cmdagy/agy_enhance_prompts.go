package cmdagy

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/spf13/cobra"
)

var (
	enhancePromptsFolder string
)

// AgyEnhancePromptsCmd decomposes a complex prompt into structured micro-task markdown files.
var AgyEnhancePromptsCmd = &cobra.Command{
	Use:     "enhance-prompts <text-or-path>",
	Aliases: []string{"ep"},
	Short:   "Decompose complex prompt into structured markdown task files (01-index, 02-verbatim, 03-...)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("usage: gitmap agy enhance-prompts <text-or-path> [-folder $variable]")
		}
		return runAgyEnhancePrompts(args[0], enhancePromptsFolder)
	},
}

func init() {
	AgyEnhancePromptsCmd.Flags().StringVarP(&enhancePromptsFolder, "folder", "f", "", "Output directory path (supports $VARIABLE expansion)")
	AgyCmd.AddCommand(AgyEnhancePromptsCmd)
}

type decomposedSection struct {
	Index    int
	Title    string
	Slug     string
	Content  string
	FileName string
}

var nonAlphaNumRegex = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(text string) string {
	lower := strings.ToLower(text)
	cleaned := nonAlphaNumRegex.ReplaceAllString(lower, "-")
	return strings.Trim(cleaned, "-")
}

func runAgyEnhancePrompts(inputArg, targetFolder string) error {
	rawContent, sourceLabel := resolveEnhancePromptContent(inputArg)
	if strings.TrimSpace(rawContent) == "" {
		return apperror.NewSimple("prompt content cannot be empty", "E_EMPTY_PROMPT")
	}

	destFolder := resolveEnhanceDestFolder(targetFolder)
	if err := os.MkdirAll(destFolder, 0755); err != nil {
		return apperror.WrapSimple(err, "create destination folder")
	}

	sections := decomposePromptContent(rawContent)

	if err := writeVerbatimFile(destFolder, rawContent, sourceLabel); err != nil {
		return err
	}

	for _, sec := range sections {
		if err := writeSectionFile(destFolder, sec); err != nil {
			return err
		}
	}

	if err := writeIndexFile(destFolder, sections, sourceLabel); err != nil {
		return err
	}

	renderEnhanceSummary(destFolder, len(sections)+2)
	return nil
}

func resolveEnhancePromptContent(arg string) (string, string) {
	fi, err := os.Stat(arg)
	if err != nil || fi.IsDir() {
		return arg, "Raw Text Input"
	}
	data, readErr := os.ReadFile(arg)
	if readErr == nil {
		return string(data), filepath.Base(arg)
	}
	return arg, "Raw Text Input"
}

func resolveEnhanceDestFolder(folderFlag string) string {
	if folderFlag != "" {
		expanded := config.ExpandVariables(folderFlag, "agy")
		return filepath.Clean(expanded)
	}

	home, _ := os.UserHomeDir()
	cwd, _ := os.Getwd()
	aiMemPrompts := filepath.Join(cwd, ".ai-memory", "prompts")
	if fi, err := os.Stat(filepath.Join(cwd, ".ai-memory")); err == nil && fi.IsDir() {
		return aiMemPrompts
	}

	defaultDir := filepath.Join(home, ".gitmap", "prompts", "enhanced-"+time.Now().Format("20060102-150405"))
	return defaultDir
}

func decomposePromptContent(content string) []decomposedSection {
	lines := strings.Split(content, "\n")
	var sections []decomposedSection
	var currentLines []string
	currentTitle := ""
	secIndex := 3

	headerRegex := regexp.MustCompile(`^(#{1,3})\s+(.+)`)
	listHeaderRegex := regexp.MustCompile(`^(\d+\.|\*|\-)\s+\*\*(.+?)\*\*`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		match := headerRegex.FindStringSubmatch(trimmed)
		hasHeaderMatch := len(match) > 2
		if hasHeaderMatch {
			currentLines = appendFlushedSection(&sections, currentLines, currentTitle, &secIndex)
			currentTitle = match[2]
			continue
		}

		listMatch := listHeaderRegex.FindStringSubmatch(trimmed)
		if len(listMatch) > 2 && currentTitle == "" {
			currentTitle = listMatch[2]
			continue
		}

		currentLines = append(currentLines, line)
	}

	if len(currentLines) > 0 {
		title := resolveSectionTitle(currentTitle, "Core Task Implementation")
		sections = append(sections, makeSection(secIndex, title, strings.Join(currentLines, "\n")))
	}

	if len(sections) == 0 {
		sections = []decomposedSection{
			makeSection(3, "Task Requirements", content),
		}
	}

	return sections
}

func makeSection(idx int, title, body string) decomposedSection {
	slug := slugify(title)
	if len(slug) > 30 {
		slug = slug[:30]
	}
	if slug == "" {
		slug = "task"
	}
	fileName := fmt.Sprintf("%02d-%s.md", idx, slug)
	return decomposedSection{
		Index:    idx,
		Title:    title,
		Slug:     slug,
		Content:  strings.TrimSpace(body),
		FileName: fileName,
	}
}

func writeVerbatimFile(destDir, rawContent, sourceLabel string) error {
	path := filepath.Join(destDir, "02-verbatim-instructions.md")
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(rawContent)))

	var sb strings.Builder
	sb.WriteString("# 02 Verbatim Instructions\n\n")
	sb.WriteString(fmt.Sprintf("- **Generated:** %s\n", time.Now().UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("- **Source:** `%s`\n", sourceLabel))
	sb.WriteString(fmt.Sprintf("- **SHA256:** `%s`\n\n", hash))
	sb.WriteString("```markdown\n")
	sb.WriteString(rawContent)
	sb.WriteString("\n```\n")

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func writeSectionFile(destDir string, sec decomposedSection) error {
	path := filepath.Join(destDir, sec.FileName)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %02d %s\n\n", sec.Index, sec.Title))
	sb.WriteString("## 1. Objectives & Scope\n\n")
	sb.WriteString(sec.Content)
	sb.WriteString("\n\n## 2. Implementation Steps\n\n")
	sb.WriteString("1. Inspect current workspace files and dependencies.\n")
	sb.WriteString("2. Apply required code modifications following repository standards.\n")
	sb.WriteString("3. Ensure all tests and linters pass without errors.\n\n")
	sb.WriteString("## 3. Verification & Quality Gate\n\n")
	sb.WriteString("- [ ] Unit tests pass.\n")
	sb.WriteString("- [ ] Coding guidelines adhered to.\n")
	sb.WriteString("- [ ] No regression introduced.\n")

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func writeIndexFile(destDir string, sections []decomposedSection, sourceLabel string) error {
	path := filepath.Join(destDir, "01-index.md")
	var sb strings.Builder
	sb.WriteString("# Executive Task Decomposition Index\n\n")
	sb.WriteString(fmt.Sprintf("> Generated by GitMap AGY Prompt Enhancer on %s\n\n", time.Now().UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("- **Source:** `%s`\n", sourceLabel))
	sb.WriteString(fmt.Sprintf("- **Total Files:** %d\n\n", len(sections)+2))
	sb.WriteString("## Task Manifest\n\n")
	sb.WriteString("| Order | File | Description | Status |\n")
	sb.WriteString("| :---: | :--- | :--- | :---: |\n")
	sb.WriteString("| 01 | `01-index.md` | Executive overview and manifest | Done |\n")
	sb.WriteString("| 02 | `02-verbatim-instructions.md` | Original unaltered prompt instructions | Reference |\n")

	for _, sec := range sections {
		sb.WriteString(fmt.Sprintf("| %02d | `%s` | %s | Pending |\n", sec.Index, sec.FileName, sec.Title))
	}

	sb.WriteString("\n## Execution Guide\n\n")
	sb.WriteString("To inject and execute these micro-tasks sequentially with watch mode:\n\n")
	sb.WriteString(fmt.Sprintf("```bash\ngitmap agy ip %q --watch\n```\n", destDir))

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func renderEnhanceSummary(destDir string, fileCount int) {
	fmt.Println()
	fmt.Printf("  %s╔════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║            PROMPT SUCCESSFULLY ENHANCED & DECOMPOSED       ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  • Target Directory: %s%s%s\n", constants.ColorYellow, destDir, constants.ColorReset)
	fmt.Printf("  • Total Files:      %s%d files%s\n", constants.ColorGreen, fileCount, constants.ColorReset)
	fmt.Printf("  • Entry Point:      %s%s%s\n\n", constants.ColorWhite, filepath.Join(destDir, "01-index.md"), constants.ColorReset)
	fmt.Printf("  %s Run with: %sgitmap agy ip %q --watch%s\n\n", constants.ColorCyan+"▶"+constants.ColorReset, constants.ColorYellow, destDir, constants.ColorReset)
}

func flushCurrentSection(currentLines []string, currentTitle string, secIndex int) ([]decomposedSection, bool) {
	if len(currentLines) > 0 && currentTitle != "" {
		return []decomposedSection{makeSection(secIndex, currentTitle, strings.Join(currentLines, "\n"))}, true
	}
	return nil, false
}

func resolveSectionTitle(rawTitle, fallback string) string {
	if rawTitle != "" {
		return rawTitle
	}
	return fallback
}

func appendFlushedSection(sections *[]decomposedSection, currentLines []string, currentTitle string, secIndex *int) []string {
	flushed, ok := flushCurrentSection(currentLines, currentTitle, *secIndex)
	if ok {
		*sections = append(*sections, flushed...)
		*secIndex++
		return nil
	}
	return currentLines
}
