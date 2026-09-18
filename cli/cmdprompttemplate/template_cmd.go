package cmdprompttemplate

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunPromptsTemplateCLI routes gitmap prompts-template subcommands.
func RunPromptsTemplateCLI(args []string) error {
	if len(args) == 0 || args[0] == "ls" || args[0] == "list" {
		return runListTemplates()
	}
	sub := strings.ToLower(args[0])
	switch sub {
	case "add":
		return runAddTemplate(args[1:])
	case "edit":
		return runEditTemplate(args[1:])
	case "rm", "delete", "remove":
		return runDeleteTemplate(args[1:])
	case "export":
		return runExportTemplate(args[1:])
	case "import":
		return runImportTemplate(args[1:])
	case "export-all", "exportall":
		return runExportAllTemplates(args[1:])
	case "import-all", "importall":
		return runImportAllTemplates(args[1:])
	case "-h", "--help", "help":
		return printTemplateHelp()
	default:
		return apperror.NewValidationError("unknown prompts-template subcommand: " + sub)
	}
}

func printTemplateHelp() error {
	fmt.Printf("\n%s  Usage: gitmap prompts-template <subcommand> [args]%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Subcommands:")
	fmt.Println("    ls, list                         List all registered prompt templates")
	fmt.Println("    add <name> <content>             Create a new prompt template")
	fmt.Println("    edit <name|id> <content>         Edit existing template content")
	fmt.Println("    rm, delete <name|id>             Remove a template")
	fmt.Println("    export <name|id> <dest.json>     Export single template to JSON")
	fmt.Println("    import <src.json>                Import single template from JSON")
	fmt.Println("    export-all [dest.json]           Export all templates to JSON")
	fmt.Println("    import-all <src.json>            Bulk import templates from JSON")
	fmt.Println()

	return nil
}
