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
	if sub == "-h" || sub == "--help" || sub == "help" {
		return printTemplateHelp()
	}

	return dispatchTemplateSubcommand(sub, args[1:])
}

func dispatchTemplateSubcommand(sub string, rest []string) error {
	switch sub {
	case "add":
		return runAddTemplate(rest)
	case "edit":
		return runEditTemplate(rest)
	case "rm", "delete", "remove":
		return runDeleteTemplate(rest)
	default:
		return dispatchTemplateTransferSubcommand(sub, rest)
	}
}

func dispatchTemplateTransferSubcommand(sub string, rest []string) error {
	switch sub {
	case "export":
		return runExportTemplate(rest)
	case "import":
		return runImportTemplate(rest)
	case "export-all", "exportall":
		return runExportAllTemplates(rest)
	case "import-all", "importall":
		return runImportAllTemplates(rest)
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
